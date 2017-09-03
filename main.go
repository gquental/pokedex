package main

import (
	"github.com/gquental/pokedex/config"

	"github.com/gquental/pokedex/route"
)

func main() {

	gin := route.Load()
	gin.Run(":" + config.Config.Port)
}
