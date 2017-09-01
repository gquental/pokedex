package data

import "gopkg.in/mgo.v2/bson"

type Type struct {
	ID   bson.ObjectId `bson:"_id,omitempty"`
	Name string        `bson:"name"`
}

type DamageType struct {
	Type   `bson:",inline"`
	Damage struct {
		NoDamageTo       []Type `bson:"noDamageTo" json:"no_damage_to"`
		HalfDamageTo     []Type `bson:"halfDamageTo" json:"half_damage_to"`
		DoubleDamageTo   []Type `bson:"doubleDamageTo" json:"double_damage_to"`
		NoDamageFrom     []Type `bson:"noDamageFrom" json:"no_damage_from"`
		HalfDamageFrom   []Type `bson:"halfDamageFrom" json:"half_damage_from"`
		DoubleDamageFrom []Type `bson:"doubleDamageFrom" json:"double_damage_from"`
	} `bson:"damage" json:"damage_relations"`
}
