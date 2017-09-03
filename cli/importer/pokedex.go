package importer

import (
	"fmt"

	"github.com/gquental/pokedex/persistence"

	"github.com/gquental/pokedex/config"
	"github.com/gquental/pokedex/data"
	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

func ImportPokedex() {
	pList := &pokedexList{}
	fetchFromAPI(pList)
}

type pokedexEntry struct {
	Name           string
	PokemonEntries []struct {
		PokemonSpecies struct {
			Name string
		} `json:"pokemon_species"`
	} `json:"pokemon_entries"`
}

func (p *pokedexEntry) Store() {
	pokedex := data.Pokedex{Name: p.Name}
	for _, item := range p.PokemonEntries {
		pokedex.Pokemon = append(
			pokedex.Pokemon,
			item.PokemonSpecies.Name,
		)
	}

	session, collection, err := persistence.GetCollection("pokedex")
	if err != nil {
		return
	}
	defer session.Close()

	oldPokedex := data.Pokedex{}
	err = collection.Find(bson.M{"name": pokedex.Name}).One(&oldPokedex)

	if err == nil {
		logrus.Error(fmt.Errorf("Pokedex %s already inserted.", pokedex.Name))
		return
	}

	err = collection.Insert(pokedex)

	logrus.Info(fmt.Sprintf("Pokedex %s inserted", pokedex.Name))

	if err != nil {
		logrus.Error("Couldn't insert pokedex %s: %v", pokedex.Name, err)
		return
	}
}

func (p *pokedexEntry) ExpandDetails() []ExpandEntry {
	return nil
}

type pokedexList struct {
	ListDefinition
}

func (p *pokedexList) GetEndpoint() string {
	return fmt.Sprintf("%s%s", config.Config.APIEndpoint, "pokedex")
}

func (p *pokedexList) GetCount() int {
	return p.Count
}

func (p *pokedexList) GetNext() string {
	return p.Next
}

func (p *pokedexList) EraseNext() {
	p.Next = ""
}

func (p *pokedexList) List() []ItemEntry {
	entries := []ItemEntry{}

	for _, item := range p.Results {
		entry := &pokedexEntry{}
		entries = append(
			entries,
			ItemEntry{Url: item.Url, Type: entry},
		)
	}

	return entries
}
