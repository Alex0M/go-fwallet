package transactionhashes

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

	routes := r.Group("/transaction-hash")
	routes.GET("/:hash", h.GetSingleTransactionHash)
	routes.POST("", h.AddTransactionHash)
}
