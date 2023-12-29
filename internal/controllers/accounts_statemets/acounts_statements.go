package accountsstatements

import (
	"go-fwallet/internal/middleware"
	"go-fwallet/internal/models"
	"go-fwallet/internal/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetAccountsStatements(c *gin.Context) {
	closingDate, isClosingDateSpecify := c.GetQuery("closing_date")
	ass, err := h.DB.GetAccountsStatements(closingDate, isClosingDateSpecify, c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("All accounts statemets", ass))
}

func (h *Handler) CreateAccountsStatements(c *gin.Context) {
	var rp models.AccountStatementRequestPayload
	if err := c.BindJSON(&rp); err != nil {
		e := middleware.NewHttpError("bad request", err.Error(), http.StatusBadRequest)
		c.Error(e)
		return
	}

	ass, err := h.DB.CreateAccountsStatements(&rp, c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("accounts statemets were created successfully", ass))
}

func (h *Handler) GetAccountStatement(c *gin.Context) {
	closingDate, isClosingDateSpecify := c.GetQuery("closing_date")
	ass, err := h.DB.GetAccountStatement(c.Param("accountID"), closingDate, isClosingDateSpecify, c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Account statemet", ass))
}

func (h *Handler) CreateAccountStatement(c *gin.Context) {
	var rp models.AccountStatementRequestPayload
	if err := c.BindJSON(&rp); err != nil {
		e := middleware.NewHttpError("bad request", err.Error(), http.StatusBadRequest)
		c.Error(e)
		return
	}

	as, err := h.DB.CreateAccountStatement(c.Param("accountID"), &rp, c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("accounts statemets were created successfully", as))
}
