package data

import "gopkg.in/mgo.v2/bson"

type Type struct {
	ID   bson.ObjectId `bson:"_id,omitempty"`
	Name string        `bson:"name"`
}

type DamageType struct {
	Type
	Damage struct {
		NoDamageTo       []Type `bson:"noDamageTo"`
		HalfDamageTo     []Type `bson:"halfDamageTo"`
		DoubleDamageTo   []Type `bson:"doubleDamageTo"`
		NoDamageFrom     []Type `bson:"noDamageFrom"`
		HalfDamageFrom   []Type `bson:"halfDamageFrom"`
		DoubleDamageFrom []Type `bson:"doubleDamageFrom"`
	} `bson:"damage"`
}

