build-images:
    docker-compose build

run-server:
    docker-compose up -d

import-types:
    docker-compose run pokedex /app/pokedex-cli importTypes

import-pokemons:
    docker-compose run pokedex /app/pokedex-cli importPokemons

import-pokedex:
    docker-compose run pokedex /app/pokedex-cli importPokedex