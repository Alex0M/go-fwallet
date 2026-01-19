package transactionhashes

import (
	"go-fwallet/internal/middleware"
	"go-fwallet/internal/models"
	"go-fwallet/internal/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetSingleTransactionHash(c *gin.Context) {
	t, err := h.DB.GetSingleTransactionHash(c.Param("hash"), c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Transaction hash", t))
}

func (h *Handler) AddTransactionHash(c *gin.Context) {
	var th models.TransactionHash
	if err := c.BindJSON(&th); err != nil {
		e := middleware.NewHttpError("bad request", err.Error(), http.StatusBadRequest)
		c.Error(e)
		return
	}

	err := h.DB.AddTransactionHash(&th, c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, response.SuccessResponse("Transaction hash was created successfully", nil))
}
