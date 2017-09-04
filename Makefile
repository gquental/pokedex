build:
	docker-compose build

run:
	docker-compose up -d

import-types:
	docker-compose run pokedex /app/pokedex-cli importTypes

import-pokedex:
	docker-compose run pokedex /app/pokedex-cli importPokedex

import-pokemons:
	docker-compose run pokedex /app/pokedex-cli importPokemons