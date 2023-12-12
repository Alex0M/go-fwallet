package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"go-fwallet/internal/controllers/accounts"
	"go-fwallet/internal/controllers/categories"
	"go-fwallet/internal/controllers/transactions"
	"go-fwallet/internal/controllers/transactiontypes"
	"go-fwallet/internal/controllers/users"
	"go-fwallet/internal/database"
	routes "go-fwallet/routes"

	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
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

	r := gin.New()
	r.Use(ginzap.Ginzap(l, time.RFC3339, true))
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	accounts.RegisterRoutes(r, db)
	transactions.RegisterRoutes(r, db)
	categories.RegisterRoutes(r, db)
	users.RegisterRoutes(r, db)
	transactiontypes.RegisterRoutes(r, db)

	routes.Routes(r)
	//Move Servier IP and Port to config
	log.Fatal(r.Run("0.0.0.0:9090"))
}
