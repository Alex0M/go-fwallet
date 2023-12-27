package categories

import (
	"go-fwallet/internal/middleware"
	"go-fwallet/internal/models"
	"go-fwallet/internal/response"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetCategories(c *gin.Context) {
	cats, err := h.DB.GetCategories(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("All categories", cats))
}

func (h *Handler) GetCategory(c *gin.Context) {
	cat, err := h.DB.GetCategory(c.Param("id"), c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Single category", cat))
}

func (h *Handler) GetCategoryByName(c *gin.Context) {
	cat, err := h.DB.GetCategoryByName(c.Param("name"), c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Single category", cat))
}

func (h *Handler) AddCategory(c *gin.Context) {
	var cat models.Category
	if err := c.BindJSON(&cat); err != nil {
		e := middleware.NewHttpError("bad request", err.Error(), http.StatusBadRequest)
		c.Error(e)
		return
	}

	err := h.DB.AddCategory(&cat, c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, response.SuccessResponse("Category was created successfully", cat))
}

func (h *Handler) EditCategory(c *gin.Context) {
	var cat = models.Category{
		UpdatedAt: time.Now(),
	}

	if err := c.BindJSON(&cat); err != nil {
		e := middleware.NewHttpError("bad request", err.Error(), http.StatusBadRequest)
		c.Error(e)
		return
	}

	newCategory, err := h.DB.EditCategory(c.Param("id"), &cat, c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Category was updated successfully", newCategory))
}

func (h *Handler) DeleteCategory(c *gin.Context) {
	err := h.DB.DeleteCategory(c.Param("id"), c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Category was deleted successfully", nil))
}
