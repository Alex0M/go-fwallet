package accounts

import (
	"go-fwallet/internal/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) DeleteAccount(c *gin.Context) {
	a, err := h.DB.DeleteAccount(c.Param("id"), c.Request.Context())
	if err != nil {
		if err == h.DB.GetErrValidateID() {
			c.JSON(http.StatusBadRequest, response.ErrorResponse(h.DB.GetErrValidateID().Error()))
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorResponse("Internal server error"))
		return
	}

	if a == nil {
		c.JSON(http.StatusNotFound, response.NotFound())
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Account deleted successfully", nil))
}
