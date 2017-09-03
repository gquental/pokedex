package importer

import (
	"fmt"

	"github.com/gquental/pokedex/config"
	"github.com/gquental/pokedex/data"
	"github.com/gquental/pokedex/persistence"
	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

func ImportTypes() {
	tList := &typeList{}
	fetchFromAPI(tList)
}

type typeEntry struct {
	data.DamageType
	FinalData data.DamageType
}

func (t *typeEntry) ExpandDetails() []ExpandEntry {
	return nil
}

func (t *typeEntry) Store() {
	t.FinalData.Name = t.Name
	t.FinalData.Damage = t.Damage

	session, collection, err := persistence.GetCollection("types")
	if err != nil {
		return
	}
	defer session.Close()

	oldType := data.DamageType{}
	err = collection.Find(bson.M{"name": t.FinalData.Name}).One(&oldType)

	if err == nil {
		logrus.Error(fmt.Errorf("Type %s already inserted.", t.FinalData.Name))
		return
	}

	err = collection.Insert(t.FinalData)

	if err != nil {
		logrus.Error("Couldn't insert type %s: %v", t.FinalData.Name, err)
		return
	}
}

type typeList struct {
	ListDefinition
}

func (t *typeList) GetEndpoint() string {
	return fmt.Sprintf("%s%s", config.Config.APIEndpoint, "type")
}

func (t *typeList) GetCount() int {
	return t.Count
}

func (t *typeList) GetNext() string {
	return t.Next
}

func (t *typeList) EraseNext() {
	t.Next = ""
}

func (t *typeList) List() []ItemEntry {
	entries := []ItemEntry{}

	for _, item := range t.Results {
		entry := &typeEntry{}
		entries = append(
			entries,
			ItemEntry{Url: item.Url, Type: entry},
		)
	}

	return entries
}
