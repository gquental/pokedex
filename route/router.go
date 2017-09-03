package route

import (
	"github.com/gquental/pokedex/server"

	"github.com/gin-gonic/gin"
)

func Load() *gin.Engine {
	router := gin.Default()

	pokemon := router.Group("/pokemon")
	{
		pokemon.GET("/:pokemon", server.GetPokemonDetail)
		pokemon.GET("/", server.GetPokemonList)
	}

	return router
}
