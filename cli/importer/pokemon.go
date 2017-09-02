package importer

import (
	"encoding/json"
	"fmt"

	"github.com/gquental/pokedex/config"
	"github.com/sirupsen/logrus"
)

func ImportPokemons() {
	fetchPokemonsFromAPI()
}

func fetchPokemonsFromAPI() {
	pokemonCountCh := make(chan int)
	pokemonCh := make(chan pokemonType)

	var totalPokemons int
	var countPokemons int

	go fetchPokemonList("", pokemonCountCh, pokemonCh)
	totalPokemons = <-pokemonCountCh
	close(pokemonCountCh)

	for {
		select {
		case pokemon := <-pokemonCh:
			fmt.Println(pokemon)
			countPokemons++
		}

		if countPokemons == totalPokemons {
			close(pokemonCh)
			break
		}
	}
}

func fetchPokemonList(url string, pokemonCountCh chan int, pokemonCh chan pokemonType) {
	firstTime := false
	if url == "" {
		firstTime = true
		url = fmt.Sprintf("%s%s", config.Config.APIEndpoint, "pokemon?limit=100")
	}

	body, err := requestAPI(url)
	if err != nil {
		logrus.Error(fmt.Errorf("Problem in pokemon list request: %v", err))
		close(pokemonCh)
		if firstTime {
			pokemonCountCh <- 0
		}
		return
	}

	pokemons := listType{}
	json.Unmarshal(body, &pokemons)

	if firstTime {
		pokemonCountCh <- pokemons.Count
	}

	if pokemons.Next != "" {
		go fetchPokemonList(pokemons.Next, pokemonCountCh, pokemonCh)
	}

	for _, item := range pokemons.Results {
		go fetchPokemonDetail(item.Url, pokemonCh)
	}
}

type pokemonType struct {
	Name  string
	Types []struct {
		Type struct {
			Name string
		}
	}
}

func fetchPokemonDetail(url string, pokemonCh chan pokemonType) {
	body, err := requestAPI(url)
	if err != nil {
		logrus.Error(fmt.Errorf("Problem in pokemon detail request: %v", err))
		close(pokemonCh)
		return
	}

	pokemon := pokemonType{}
	json.Unmarshal(body, &pokemon)

	pokemonCh <- pokemon
}
