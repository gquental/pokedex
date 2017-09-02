package importer

import (
	"fmt"

	"encoding/json"

	"github.com/gquental/pokedex/config"
)

func ImportPokemons() {
	pList := &pokemonList{}
	fetchFromAPI(pList)
}

type pokemonEntry struct {
	Name      string
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
	fmt.Println(p)
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
