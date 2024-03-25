package middleware

import (
	"go-fwallet/internal/auth"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		err := auth.ValidateToken(tokenString)
		if err != nil {
			c.Error(NewHttpError("Unauthorized", "Unauthorized", http.StatusForbidden))
			return
		}

		c.Next()
	}
}
