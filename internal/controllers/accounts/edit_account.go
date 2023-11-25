package accounts

import (
	"fmt"
	"go-fwallet/internal/models"
	"go-fwallet/internal/response"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *Handler) EditAccount(c *gin.Context) {
	var a = models.Account{
		UpdatedAt: time.Now(),
	}
	//Mybe need to move this code to function
	if err := c.BindJSON(&a); err != nil {
		h.DB.Logger.Error("error bindJSON", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorResponse(fmt.Sprintf("Bad request: %s", err)))
		return
	}

	newAccount, err := h.DB.EditAccount(c.Param("id"), &a, c.Request.Context())
	if err != nil {
		if err == h.DB.GetErrValidateID() {
			c.JSON(http.StatusBadRequest, response.ErrorResponse(h.DB.GetErrValidateID().Error()))
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorResponse("Internal server error"))
		return
	}

	if newAccount == nil {
		c.JSON(http.StatusNotFound, response.NotFound())
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Single account", newAccount))
}
