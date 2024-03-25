package login

import (
	"go-fwallet/internal/auth"
	"go-fwallet/internal/middleware"
	"go-fwallet/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Login(c *gin.Context) {
	var u models.UserLogin
	if err := c.ShouldBindJSON(&u); err != nil {
		e := middleware.NewHttpError("bad request", err.Error(), http.StatusBadRequest)
		c.Error(e)
		return
	}

	user, err := h.DB.LoginCheck(&u, c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	token, err := auth.GenerateToken(user.Email, user.Username)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
