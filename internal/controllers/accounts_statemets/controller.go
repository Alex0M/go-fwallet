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

	routes := r.Group("/accountstatements")
	routes.GET("", h.GetAccountsStatements)
	routes.POST("", h.CreateAccountsStatements)
	routes.GET("/:accountID", h.GetAccountStatement)
	routes.POST("/:accountID", h.CreateAccountStatement)
}
