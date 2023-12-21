package accounts

import (
	"go-fwallet/internal/middleware"
	"go-fwallet/internal/models"
	"go-fwallet/internal/response"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetAccounts(c *gin.Context) {
	acc, err := h.DB.GetAccounts(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("All accounts", acc))
}

func (h *Handler) GetSingleAccount(c *gin.Context) {
	a, err := h.DB.GetSingleAccount(c.Param("id"), c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Single account", a))
}

func (h *Handler) AddAcount(c *gin.Context) {
	var account models.Account
	if err := c.BindJSON(&account); err != nil {
		err := middleware.NewHttpError("bad request", err.Error(), http.StatusBadRequest)
		c.Error(err)
		return
	}

	err := h.DB.AddAccount(&account, c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, response.SuccessResponse("Account was created successfully", nil))
}

func (h *Handler) EditAccount(c *gin.Context) {
	var a = models.Account{
		UpdatedAt: time.Now(),
	}

	if err := c.BindJSON(&a); err != nil {
		err := middleware.NewHttpError("bad request", err.Error(), http.StatusBadRequest)
		c.Error(err)
		return
	}

	newAccount, err := h.DB.EditAccount(c.Param("id"), &a, c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Single account", newAccount))
}

func (h *Handler) DeleteAccount(c *gin.Context) {
	err := h.DB.DeleteAccount(c.Param("id"), c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Account was deleted successfully", nil))
}
