package main

import (
	rout "ginchat/router"
	"ginchat/utils"
)

func main() {
	utils.InitConfig()
	utils.InitMySQL()
	utils.InitRedis()

	r := rout.Router()
	r.Run("localhost:8080")
}
