package accountsstatements

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

	r.GET("/accountstatements", h.GetAccountsStatements)
	r.POST("/accountstatements", h.CreateAccountsStatements)
}
