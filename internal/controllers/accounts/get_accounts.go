package accounts

import (
	"go-fwallet/internal/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetAccounts(c *gin.Context) {
	acc, err := h.DB.GetAccounts(c.Request.Context())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorResponse("Internal server error"))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("All accounts", acc))
}
