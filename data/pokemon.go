package data

import "gopkg.in/mgo.v2/bson"

type Pokemon struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	PokemonID int
	Name      string
	Abilities []string
	Stats     []PokemonStat
	Types     []string
}

type PokemonStat struct {
	Name       string
	Base       int
	BattleOnly bool
}
