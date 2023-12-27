package users

import (
	"go-fwallet/internal/middleware"
	"go-fwallet/internal/models"
	"go-fwallet/internal/response"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetUsers(c *gin.Context) {
	acc, err := h.DB.GetUsers(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("All users", acc))
}

func (h *Handler) GetUser(c *gin.Context) {
	u, err := h.DB.GetUser(c.Param("id"), c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Single user", u))
}

func (h *Handler) AddUser(c *gin.Context) {
	var u models.User
	if err := c.BindJSON(&u); err != nil {
		e := middleware.NewHttpError("bad request", err.Error(), http.StatusBadRequest)
		c.Error(e)
		return
	}

	err := h.DB.AddUser(&u, c.Request.Context())
	if err != nil {
		c.Error(err)
	}

	c.JSON(http.StatusCreated, response.SuccessResponse("User was created successfully", nil))
}

func (h *Handler) EditUser(c *gin.Context) {
	var u = models.User{
		UpdatedAt: time.Now(),
	}

	if err := c.BindJSON(&u); err != nil {
		e := middleware.NewHttpError("bad request", err.Error(), http.StatusBadRequest)
		c.Error(e)
		return
	}

	newAccount, err := h.DB.EditUser(c.Param("id"), &u, c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Single user", newAccount))
}

func (h *Handler) DeleteUser(c *gin.Context) {
	err := h.DB.DeleteUser(c.Param("id"), c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("User was deleted successfully", nil))
}
