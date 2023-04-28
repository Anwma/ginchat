package main

import (
	rout "ginchat/router"
	"ginchat/utils"
)

func main() {
	utils.InitConfig()
	utils.InitMySQL()

	r := rout.Router()
	r.Run("localhost:8080")
}
