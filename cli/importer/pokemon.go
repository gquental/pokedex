package importer

import (
	"fmt"

	"github.com/gquental/pokedex/persistence"

	"gopkg.in/mgo.v2/bson"

	"encoding/json"

	"github.com/gquental/pokedex/config"
	"github.com/gquental/pokedex/data"
	"github.com/sirupsen/logrus"
)

func ImportPokemons() {
	pList := &pokemonList{}
	fetchFromAPI(pList)
}

type pokemonEntry struct {
	Name      string
	ID        int `json:"id"`
	Abilities []struct {
		Ability struct {
			Name string
		}
	}
	Stats []struct {
		Stat struct {
			Url      string
			Name     string
			Expanded struct {
				BattleOnly bool `json:"is_battle_only"`
			}
		}
		BaseStat int `json:"base_stat"`
	}
	Types []struct {
		Type struct {
			Name string
		}
	}
}

func (p *pokemonEntry) Store() {
	pokemon := data.Pokemon{Name: p.Name, PokemonID: p.ID}
	for _, item := range p.Abilities {
		pokemon.Abilities = append(
			pokemon.Abilities,
			item.Ability.Name,
		)
	}

	for _, item := range p.Stats {
		pokemon.Stats = append(
			pokemon.Stats,
			data.PokemonStat{
				Name:       item.Stat.Name,
				Base:       item.BaseStat,
				BattleOnly: item.Stat.Expanded.BattleOnly,
			},
		)
	}

	for _, item := range p.Types {
		pokemon.Types = append(
			pokemon.Types,
			item.Type.Name,
		)
	}

	session, collection, err := persistence.GetCollection("pokemons")
	if err != nil {
		return
	}
	defer session.Close()

	oldPokemon := data.Pokemon{}
	err = collection.Find(bson.M{"name": pokemon.Name}).One(&oldPokemon)

	if err == nil {
		logrus.Error(fmt.Errorf("Pokemon %s already inserted.", pokemon.Name))
		return
	}

	err = collection.Insert(pokemon)

	logrus.Info(fmt.Sprintf("Pokemon %s inserted", pokemon.Name))

	if err != nil {
		logrus.Error("Couldn't insert pokemon %s: %v", pokemon.Name, err)
		return
	}
}

func (p *pokemonEntry) ExpandDetails() []ExpandEntry {
	entries := []ExpandEntry{}
	for index, item := range p.Stats {
		entries = append(
			entries,
			ExpandEntry{
				Url:   item.Stat.Url,
				Index: index,
				Assign: func(entry interface{}, body []byte, index int) ImportableDetail {
					pokemon := entry.(*pokemonEntry)
					json.Unmarshal(body, &pokemon.Stats[index].Stat.Expanded)

					return pokemon
				},
			},
		)
	}

	return entries
}

type pokemonList struct {
	ListDefinition
}

func (p *pokemonList) GetEndpoint() string {
	return fmt.Sprintf("%s%s", config.Config.APIEndpoint, "pokemon?limit=100")
}

func (p *pokemonList) GetCount() int {
	return p.Count
}

func (p *pokemonList) GetNext() string {
	return p.Next
}

func (p *pokemonList) EraseNext() {
	p.Next = ""
}

func (p *pokemonList) List() []ItemEntry {
	entries := []ItemEntry{}

	for _, item := range p.Results {
		entry := &pokemonEntry{}
		entries = append(
			entries,
			ItemEntry{Url: item.Url, Type: entry},
		)
	}

	return entries
}
