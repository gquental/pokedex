package server

import (
	"net/http"

	"github.com/gquental/pokedex/persistence"

	"strconv"

	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gquental/pokedex/data"
	"gopkg.in/mgo.v2/bson"
)

func GetPokemonDetail(c *gin.Context) {
	identifier := c.Param("pokemon")

	pokemon := data.Pokemon{}

	session, collection, err := persistence.GetCollection("pokemons")
	defer session.Close()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "It was not possible to connect to pokemon collection"})
		return
	}

	err = collection.Find(bson.M{"pokemonid": identifier}).One(&pokemon)
	if err == nil {
		types, err := getTypes(pokemon)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"pokemon": pokemon, "damages": types})
	}

	err = collection.Find(bson.M{"name": identifier}).One(&pokemon)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	types, err := getTypes(pokemon)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"pokemon": pokemon, "damages": types})
}

func GetPokemonList(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		page = 1
	}

	session, collection, err := persistence.GetCollection("pokemons")
	defer session.Close()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "It was not possible to connect to pokemon collection"})
		return
	}

	pokemons := []data.Pokemon{}
	err = collection.Find(bson.M{}).Limit(20).Skip(20 * page).All(&pokemons)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	count, err := collection.Count()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	pages := count / 20
	nextPage := page + 1
	next := fmt.Sprintf("%s%s%s%d", c.Request.URL.Scheme, c.Request.URL.Host, "pokemon?page=", nextPage)

	if page >= pages {
		next = ""
	}

	c.JSON(http.StatusOK, gin.H{"count": count, "pokemons": pokemons, "next": next})

}

func getTypes(pokemon data.Pokemon) ([]data.DamageType, error) {
	session, collection, err := persistence.GetCollection("types")
	defer session.Close()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	types := []data.DamageType{}
	err = collection.Find(bson.M{"name": bson.M{"$in": pokemon.Types}}).All(&types)

	if err != nil {
		return nil, err
	}

	return types, nil
}
