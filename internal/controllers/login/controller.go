package login

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

	r.POST("/login", h.Login)
}
