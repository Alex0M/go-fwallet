package users

import (
	"fmt"
	"go-fwallet/internal/models"
	"go-fwallet/internal/response"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *Handler) GetUsers(c *gin.Context) {
	acc, err := h.DB.GetUsers(c.Request.Context())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorResponse("Internal server error"))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("All users", acc))
}

func (h *Handler) GetUser(c *gin.Context) {
	u, err := h.DB.GetUser(c.Param("id"), c.Request.Context())
	if err != nil {
		if err == h.DB.GetErrValidateID() {
			c.JSON(http.StatusBadRequest, response.ErrorResponse(h.DB.GetErrValidateID().Error()))
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorResponse("Internal server error"))
		return
	}

	if u == nil {
		c.JSON(http.StatusNotFound, response.NotFound())
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Single user", u))
}

func (h *Handler) AddUser(c *gin.Context) {
	var u models.User
	if err := c.BindJSON(&u); err != nil {
		h.DB.Logger.Error("error bindJSON", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorResponse(fmt.Sprintf("Bad request: %s", err)))
		return
	}

	err := h.DB.AddUser(&u, c.Request.Context())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorResponse("Internal server error"))
		return
	}

	c.JSON(http.StatusCreated, response.SuccessResponse("User was created successfully", nil))
}

func (h *Handler) EditUser(c *gin.Context) {
	var u = models.User{
		UpdatedAt: time.Now(),
	}
	//Mybe need to move this code to function
	if err := c.BindJSON(&u); err != nil {
		h.DB.Logger.Error("error bindJSON", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorResponse(fmt.Sprintf("Bad request: %s", err)))
		return
	}

	newAccount, err := h.DB.EditUser(c.Param("id"), &u, c.Request.Context())
	if err != nil {
		if err == h.DB.GetErrValidateID() {
			c.JSON(http.StatusBadRequest, response.ErrorResponse(h.DB.GetErrValidateID().Error()))
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorResponse("Internal server error"))
		return
	}

	if newAccount == nil {
		c.JSON(http.StatusNotFound, response.NotFound())
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Single user", newAccount))
}

func (h *Handler) DeleteUser(c *gin.Context) {
	u, err := h.DB.DeleteUser(c.Param("id"), c.Request.Context())
	if err != nil {
		if err == h.DB.GetErrValidateID() {
			c.JSON(http.StatusBadRequest, response.ErrorResponse(h.DB.GetErrValidateID().Error()))
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorResponse("Internal server error"))
		return
	}

	if u == nil {
		c.JSON(http.StatusNotFound, response.NotFound())
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("User was deleted successfully", nil))
}
