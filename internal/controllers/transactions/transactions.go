package transactions

import (
	"fmt"
	"go-fwallet/internal/middleware"
	"go-fwallet/internal/models"
	"go-fwallet/internal/response"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetAllTransaction(c *gin.Context) {
	t, err := h.DB.GetAllTransaction(getQueryParametersToStruct(c), c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("All transactions", t))
}

func (h *Handler) GetSumOfTransactionsByCategory(c *gin.Context) {
	t, err := h.DB.GetSumOfTransactionsByCategory(getQueryParametersToStruct(c), c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Sum of transactions by category", t))
}

func (h *Handler) GetTransaction(c *gin.Context) {
	t, err := h.DB.GetTransaction(c.Param("id"), c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Single transaction", t))
}

func (h *Handler) AddTransaction(c *gin.Context) {
	var t models.Transaction
	if err := c.BindJSON(&t); err != nil {
		e := middleware.NewHttpError("bad request", err.Error(), http.StatusBadRequest)
		c.Error(e)
		return
	}

	err := h.DB.AddTransaction(&t, c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, response.SuccessResponse("Transaction was created successfully", nil))
}

func (h *Handler) EditTransaction(c *gin.Context) {
	var t = models.Transaction{
		UpdatedAt: time.Now(),
	}

	if err := c.BindJSON(&t); err != nil {
		e := middleware.NewHttpError("bad request", err.Error(), http.StatusBadRequest)
		c.Error(e)
		return
	}

	newTransaction, err := h.DB.EditTransaction(c.Param("id"), &t, c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Transaction was updated successfully", newTransaction))
}

func (h *Handler) DeleteTransaction(c *gin.Context) {
	err := h.DB.DeleteTransaction(c.Param("id"), c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Transaction was deleted successfully", nil))
}

func (h *Handler) UpdateTransactionCategories(c *gin.Context) {
	var tap models.TransactionAPIRequestParams
	if err := c.BindJSON(&tap); err != nil {
		e := middleware.NewHttpError("bad request", err.Error(), http.StatusBadRequest)
		c.Error(e)
		return
	}

	rows, err := h.DB.UpdateTransactionCategories(&tap, c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, response.SuccessResponse(fmt.Sprintf("Row updated: %d", rows), nil))
}

func getQueryParametersToStruct(c *gin.Context) *models.TransactionAPIRequestParams {
	var tp models.TransactionAPIRequestParams
	tp.GteDate, tp.IsGteDateSpecify = c.GetQuery("date_gte")
	tp.LteDate, tp.IsLteDateSpecify = c.GetQuery("date_lte")
	tp.TransactionTypeCode, tp.IsTransactionTypeCodeSpecify = c.GetQuery("transaction_type_code")

	return &tp
}
