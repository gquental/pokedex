# pokedex

The project has two main areas, the API and the importer.

## Requirements

You must have installed
- Make
- Docker
- Docker compose

## Install

Before executing the commands, be sure to copy the **config.yml.dev** to **config.yml**

To run the application we must execute two commands, one to build the images and other to run the server.

- ```make build```
- ```make run```


In order to have content in the API we must run the importer, in the following order

- ```make import-types```
- ```make import-pokedex```
- ```make import-pokemons```