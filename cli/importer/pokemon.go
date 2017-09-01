package importer

import "github.com/gquental/pokedex/data"

func ImportPokemons() {

}

func fetchPokemonsFromAPI() {
	pokemonCountCh := make(chan int)
	pokemonCh := make(chan data.Pokemon)

	var totalTypes int
	var countTypes int

	go fetchTypeList("", pokemonCountCh, pokemonCh)
	totalTypes = <-pokemonCountCh
	close(pokemonCountCh)

	for {
		select {
		case pokemon := <-pokemonCh:
			countTypes++
		}

		if countTypes == totalTypes {
			close(pokemonCh)
			break
		}
	}
}
