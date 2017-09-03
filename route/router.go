package route

import (
	"github.com/gquental/pokedex/server"

	"github.com/gin-gonic/contrib/jwt"
	"github.com/gin-gonic/gin"
)

func Load() *gin.Engine {
	router := gin.Default()

	pokemon := router.Group("/pokemon")
	{
		pokemon.GET("/:pokemon", server.GetPokemonDetail)
		pokemon.GET("/", server.GetPokemonList)
	}

	pokedex := router.Group("/pokedex")
	{
		pokedex.Use(jwt.Auth("secret"))
		pokedex.GET("/:region", server.GetPokedex)
	}

	return router
}
