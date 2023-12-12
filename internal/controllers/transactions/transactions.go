package transactions

import (
	"fmt"
	"go-fwallet/internal/models"
	"go-fwallet/internal/response"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *Handler) GetAllTransaction(c *gin.Context) {
	t, err := h.DB.GetAllTransaction(getQueryParametersToStruct(c), c.Request.Context())
	if err != nil {
		if err == h.DB.GetErrValidateID() {
			c.JSON(http.StatusBadRequest, response.ErrorResponse(h.DB.GetErrValidateID().Error()))
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorResponse("Internal server error"))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("All transactions", t))
}

func (h *Handler) GetSumOfTransactionsByCategory(c *gin.Context) {
	t, err := h.DB.GetSumOfTransactionsByCategory(getQueryParametersToStruct(c), c.Request.Context())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorResponse("Internal server error"))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Sum of transactions by category", t))
}

func (h *Handler) GetTransaction(c *gin.Context) {
	t, err := h.DB.GetTransaction(c.Param("id"), c.Request.Context())
	if err != nil {
		if err == h.DB.GetErrValidateID() {
			c.JSON(http.StatusBadRequest, response.ErrorResponse(h.DB.GetErrValidateID().Error()))
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorResponse("Internal server error"))
		return
	}

	if t == nil {
		c.JSON(http.StatusNotFound, response.NotFound())
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Single transaction", t))
}

func (h *Handler) AddTransaction(c *gin.Context) {
	var t models.Transaction
	if err := c.BindJSON(&t); err != nil {
		h.DB.Logger.Error("error bindJSON", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorResponse(fmt.Sprintf("Bad request: %s", err)))
		return
	}

	err := h.DB.AddTransaction(&t, c.Request.Context())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorResponse("Internal server error"))
		return
	}

	c.JSON(http.StatusCreated, response.SuccessResponse("Transaction was created successfully", nil))
}

func (h *Handler) EditTransaction(c *gin.Context) {
	var t = models.Transaction{
		UpdatedAt: time.Now(),
	}

	if err := c.BindJSON(&t); err != nil {
		h.DB.Logger.Error("error bindJSON", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorResponse(fmt.Sprintf("Bad request: %s", err)))
		return
	}

	newTransaction, err := h.DB.EditTransaction(c.Param("id"), &t, c.Request.Context())
	if err != nil {
		if err == h.DB.GetErrValidateID() {
			c.JSON(http.StatusBadRequest, response.ErrorResponse(h.DB.GetErrValidateID().Error()))
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorResponse("Internal server error"))
		return
	}

	if newTransaction == nil {
		c.JSON(http.StatusNotFound, response.NotFound())
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Transaction was updated successfully", newTransaction))
}

func (h *Handler) DeleteTransaction(c *gin.Context) {
	t, err := h.DB.DeleteTransaction(c.Param("id"), c.Request.Context())
	if err != nil {
		if err == h.DB.GetErrValidateID() {
			c.JSON(http.StatusBadRequest, response.ErrorResponse(h.DB.GetErrValidateID().Error()))
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorResponse("Internal server error"))
		return
	}

	if t == nil {
		c.JSON(http.StatusNotFound, response.NotFound())
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Transaction was deleted successfully", nil))
}

func (h *Handler) UpdateTransactionCategories(c *gin.Context) {
	var tap models.TransactionAPIRequestParams
	if err := c.BindJSON(&tap); err != nil {
		h.DB.Logger.Error("error bindJSON", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorResponse(fmt.Sprintf("Bad request: %s", err)))
		return
	}

	rows, err := h.DB.UpdateTransactionCategories(&tap, c.Request.Context())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorResponse("Internal server error"))
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
