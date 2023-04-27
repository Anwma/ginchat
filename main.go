package main

func main() {
	utils.InitConfig()
	utils.InitMySQL()

	r:=router.Router()
	r.Run("localhost:8081")
}
