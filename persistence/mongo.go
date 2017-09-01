package persistence

import (
	"fmt"

	"github.com/gquental/pokedex/config"

	"gopkg.in/mgo.v2"
)

func GetSession() (*mgo.Session, error) {
	session, err := mgo.Dial(config.Config.DBAddress)

	return session, err
}

func GetCollection(collection string) (*mgo.Session, *mgo.Collection, error) {
	session, err := GetSession()
	if err != nil {
		return &mgo.Session{}, &mgo.Collection{}, fmt.Errorf("Not possible to connect to %s collection: %v", collection, err)
	}

	c := session.DB(config.Config.Database).C(collection)
	return session, c, nil
}
