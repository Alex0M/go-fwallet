package categories

import (
	"fmt"
	"go-fwallet/internal/models"
	"go-fwallet/internal/response"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *Handler) GetCategories(c *gin.Context) {
	cats, err := h.DB.GetCategories(c.Request.Context())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorResponse("Internal server error"))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("All categories", cats))
}

func (h *Handler) GetCategory(c *gin.Context) {
	cat, err := h.DB.GetCategory(c.Param("id"), c.Request.Context())
	if err != nil {
		if err == h.DB.GetErrValidateID() {
			c.JSON(http.StatusBadRequest, response.ErrorResponse(h.DB.GetErrValidateID().Error()))
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorResponse("Internal server error"))
		return
	}

	if cat == nil {
		c.JSON(http.StatusNotFound, response.NotFound())
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Single category", cat))
}

func (h *Handler) GetCategoryByName(c *gin.Context) {
	cat, err := h.DB.GetCategoryByName(c.Param("name"), c.Request.Context())
	if err != nil {
		if err == h.DB.GetErrValidateID() {
			c.JSON(http.StatusBadRequest, response.ErrorResponse(h.DB.GetErrValidateID().Error()))
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorResponse("Internal server error"))
		return
	}

	if cat == nil {
		c.JSON(http.StatusNotFound, response.NotFound())
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Single category", cat))
}

func (h *Handler) AddCategory(c *gin.Context) {
	var cat models.Category
	if err := c.BindJSON(&cat); err != nil {
		h.DB.Logger.Error("error bindJSON", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorResponse(fmt.Sprintf("Bad request: %s", err)))
		return
	}

	err := h.DB.AddCategory(&cat, c.Request.Context())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorResponse("Internal server error"))
		return
	}

	c.JSON(http.StatusCreated, response.SuccessResponse("Category was created successfully", cat))
}

func (h *Handler) EditCategory(c *gin.Context) {
	var cat = models.Category{
		UpdatedAt: time.Now(),
	}

	if err := c.BindJSON(&cat); err != nil {
		h.DB.Logger.Error("error bindJSON", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorResponse(fmt.Sprintf("Bad request: %s", err)))
		return
	}

	newCategory, err := h.DB.EditCategory(c.Param("id"), &cat, c.Request.Context())
	if err != nil {
		if err == h.DB.GetErrValidateID() {
			c.JSON(http.StatusBadRequest, response.ErrorResponse(h.DB.GetErrValidateID().Error()))
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorResponse("Internal server error"))
		return
	}

	if newCategory == nil {
		c.JSON(http.StatusNotFound, response.NotFound())
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Category was updated successfully", newCategory))
}

func (h *Handler) DeleteCategory(c *gin.Context) {
	cat, err := h.DB.DeleteCategory(c.Param("id"), c.Request.Context())
	if err != nil {
		if err == h.DB.GetErrValidateID() {
			c.JSON(http.StatusBadRequest, response.ErrorResponse(h.DB.GetErrValidateID().Error()))
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorResponse("Internal server error"))
		return
	}

	if cat == nil {
		c.JSON(http.StatusNotFound, response.NotFound())
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Category was deleted successfully", nil))
}
