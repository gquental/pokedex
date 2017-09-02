package importer

import (
	"encoding/json"
	"fmt"

	"github.com/gquental/pokedex/config"
	"github.com/gquental/pokedex/data"
	"github.com/gquental/pokedex/persistence"
	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

func ImportTypes() {
	fetchTypesFromAPI()
}

func fetchTypesFromAPI() {
	typeCountCh := make(chan int)
	typesCh := make(chan data.DamageType)

	var totalTypes int
	var countTypes int

	go fetchTypeList("", typeCountCh, typesCh)
	totalTypes = <-typeCountCh
	close(typeCountCh)

	for {
		select {
		case typeDetail := <-typesCh:
			storeType(typeDetail)
			countTypes++
		}

		if countTypes == totalTypes {
			close(typesCh)
			break
		}
	}

}

func fetchTypeList(url string, typeCountCh chan int, typesCh chan data.DamageType) {
	firstTime := false
	if url == "" {
		firstTime = true
		url = fmt.Sprintf("%s%s", config.Config.APIEndpoint, "type")
	}

	body, err := requestAPI(url)
	if err != nil {
		logrus.Error(fmt.Errorf("Problem in type list request: %v", err))
		close(typesCh)
		if firstTime {
			typeCountCh <- 0
		}
		return
	}

	types := listType{}
	json.Unmarshal(body, &types)

	if firstTime {
		typeCountCh <- types.Count
	}

	if types.Next != "" {
		go fetchTypeList(types.Next, typeCountCh, typesCh)
	}

	for _, item := range types.Results {
		go fetchTypeDetail(item.Url, typesCh)
	}
}

func fetchTypeDetail(url string, typesCh chan data.DamageType) {
	body, err := requestAPI(url)
	if err != nil {
		logrus.Error(fmt.Errorf("Problem in type detail request: %v", err))
		close(typesCh)
		return
	}

	typeDetail := data.DamageType{}
	json.Unmarshal(body, &typeDetail)

	typesCh <- typeDetail
}

func storeType(typeDetail data.DamageType) {
	session, collection, err := persistence.GetCollection("types")
	if err != nil {
		return
	}
	defer session.Close()

	oldType := data.DamageType{}
	err = collection.Find(bson.M{"name": typeDetail.Name}).One(&oldType)

	if err == nil {
		logrus.Error(fmt.Errorf("Type %s already inserted.", typeDetail.Name))
		return
	}

	err = collection.Insert(typeDetail)

	if err != nil {
		logrus.Error("Couldn't insert type %s: %v", typeDetail.Name, err)
		return
	}
}
