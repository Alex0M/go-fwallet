package transactiontypes

import (
	"go-fwallet/internal/database"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	DB *database.Database
}

func RegisterRoutes(r *gin.Engine, db *database.Database) {
	h := &Handler{
		DB: db,
	}

	routes := r.Group("/transactiontypes")
	routes.GET("", h.GetTransactionTypes)
	routes.GET("/:code", h.GetTransactionType)
	routes.POST("", h.AddTransactionType)
	routes.PUT("/:code", h.EditTransactionType)
	routes.DELETE("/:code", h.DeleteTransactionType)
}
