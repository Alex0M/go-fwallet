package accountsstatements

import (
	"fmt"
	"go-fwallet/internal/models"
	"go-fwallet/internal/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *Handler) GetAccountsStatements(c *gin.Context) {
	closingDate, isClosingDateSpecify := c.GetQuery("closing_date")
	ass, err := h.DB.GetAccountsStatements(closingDate, isClosingDateSpecify, c.Request.Context())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorResponse("Internal server error"))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("All accounts statemets", ass))
}

func (h *Handler) CreateAccountsStatements(c *gin.Context) {
	var rp models.AccountStatementRequestPayload
	if err := c.BindJSON(&rp); err != nil {
		h.DB.Logger.Error("error bindJSON", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorResponse(fmt.Sprintf("Bad request: %s", err)))
		return
	}

	ass, err := h.DB.CreateAccountsStatements(&rp, c.Request.Context())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorResponse("Internal server error"))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("accounts statemets were created successfully", ass))
}
