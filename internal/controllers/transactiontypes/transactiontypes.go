package transactiontypes

import (
	"go-fwallet/internal/middleware"
	"go-fwallet/internal/models"
	"go-fwallet/internal/response"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetTransactionTypes(c *gin.Context) {
	tts, err := h.DB.GetTransactionTypes(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, response.SuccessResponse("All transactions types", tts))
}

func (h *Handler) GetTransactionType(c *gin.Context) {
	tt, err := h.DB.GetTransactionType(c.Param("id"), c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Single transaction type", tt))
}

func (h *Handler) AddTransactionType(c *gin.Context) {
	var tt models.TransactionType
	if err := c.BindJSON(&tt); err != nil {
		e := middleware.NewHttpError("bad request", err.Error(), http.StatusBadRequest)
		c.Error(e)
		return
	}

	err := h.DB.AddTransactionType(&tt, c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, response.SuccessResponse("Transaction type was created successfully", tt))
}

func (h *Handler) EditTransactionType(c *gin.Context) {
	var tt = models.TransactionType{
		UpdatedAt: time.Now(),
	}

	if err := c.BindJSON(&tt); err != nil {
		e := middleware.NewHttpError("bad request", err.Error(), http.StatusBadRequest)
		c.Error(e)
		return
	}

	newTransactionType, err := h.DB.EditTransactionType(c.Param("id"), &tt, c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Transaction type was updated successfully", newTransactionType))
}

func (h *Handler) DeleteTransactionType(c *gin.Context) {
	err := h.DB.DeleteTransactionType(c.Param("id"), c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Transaction type was deleted successfully", nil))
}
