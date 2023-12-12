package commoncontroller

import (
	"go-fwallet/internal/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Welcome(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Welcome to FWallet API",
	})
}

func (h *Handler) NotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, response.NotFound())
}
