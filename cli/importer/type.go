package importer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gquental/pokedex/config"
	"github.com/gquental/pokedex/data"
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
			fmt.Println(typeDetail.Name)
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

	}

	types := listType{}
	json.Unmarshal(body, &types)

	if firstTime {
		typeCountCh <- types.Count
	}

	for _, item := range types.Results {
		go fetchTypeDetail(item.Url, typesCh)
	}
}

func fetchTypeDetail(url string, typesCh chan data.DamageType) {
	body, err := requestAPI(url)
	if err != nil {

	}

	typeDetail := data.DamageType{}
	json.Unmarshal(body, &typeDetail)

	typesCh <- typeDetail
}

