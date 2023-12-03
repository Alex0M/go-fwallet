package users

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

	routes := r.Group("/users")
	routes.GET("", h.GetUsers)
	routes.POST("/", h.AddUser)
	routes.GET("/:id", h.GetUser)
	routes.PUT("/:id", h.EditUser)
	routes.DELETE("/:id", h.DeleteUser)
}
