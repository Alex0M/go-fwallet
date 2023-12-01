package transactions

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

	routes := r.Group("/transactions")
	routes.GET("/", h.GetAllTransaction)
	routes.GET("/:id", h.GetTransaction)
	routes.GET("/SumByCategory", h.GetSumOfTransactionsByCategory)
	routes.POST("/", h.AddTransaction)
	routes.PUT("/:id", h.EditTransaction)
	routes.DELETE("/:id", h.DeleteTransaction)
	routes.POST("/updateCategory", h.UpdateTransactionCategories)
}
