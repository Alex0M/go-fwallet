package main

import (
	"fmt"
	"log"
	"os"

	"go-fwallet/internal/controllers/accounts"
	"go-fwallet/internal/database"
	routes "go-fwallet/routes"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	//Move to config package
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))

	l, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("error creating logger: %s", err)
	}

	db := database.Init(dsn, l)
	err = db.Ping()
	if err != nil {
		log.Fatalf("Cannot connect to DB. %s", err)
	}

	r := gin.Default()
	accounts.RegisterRoutes(r, db)

	routes.Routes(r)
	//Move Servier IP and Port to config
	log.Fatal(r.Run("0.0.0.0:9090"))
}
