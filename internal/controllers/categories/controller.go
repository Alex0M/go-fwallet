package categories

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

	routes := r.Group("/categories")
	routes.GET("", h.GetCategories)
	routes.GET("/:id", h.GetCategory)
	routes.GET("/name/:name", h.GetCategoryByName)
	routes.POST("/", h.AddCategory)
	routes.PUT("/:id", h.EditCategory)
	routes.DELETE("/:id", h.DeleteCategory)
}
