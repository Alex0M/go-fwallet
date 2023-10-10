package main

import (
	"log"

	database "go-fwallet/database"
	routes "go-fwallet/routes"

	"github.com/gin-gonic/gin"
)

func main() {

	database.DatabaseConnection()

	router := gin.Default()
	routes.Routes(router)
	log.Fatal(router.Run("0.0.0.0:9090"))
}
