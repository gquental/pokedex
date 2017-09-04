package server

import (
	"net/http"

	"github.com/gquental/pokedex/persistence"

	"github.com/gin-gonic/gin"
	"github.com/gquental/pokedex/data"
	"gopkg.in/mgo.v2/bson"
)

func GetPokedex(c *gin.Context) {
	region := c.Param("region")

	session, collection, err := persistence.GetCollection("pokedex")
	defer session.Close()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "It was not possible to connect to pokemon collection"})
		return
	}

	pokedex := data.Pokedex{}
	err = collection.Find(bson.M{"name": region}).One(&pokedex)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pokedex": pokedex})

}
