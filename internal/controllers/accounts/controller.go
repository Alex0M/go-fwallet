package accounts

import (
	database "go-fwallet/internal/db"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	DB *database.Database
}

func RegisterRoutes(r *gin.Engine, db *database.Database) {
	h := &Handler{
		DB: db,
	}

	routes := r.Group("/accounts")
	routes.GET("/", h.GetAccounts)
	routes.POST("/", h.AddAcount)
	routes.GET("/:id", h.GetSingleAccount)
	routes.PUT("/:id", h.EditAccount)
	routes.DELETE("/:id", h.DeleteAccount)
}
