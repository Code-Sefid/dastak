package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/soheilkhaledabdi/dastak/api/dto"
	"github.com/soheilkhaledabdi/dastak/api/helper"
	"github.com/soheilkhaledabdi/dastak/config"
	"github.com/soheilkhaledabdi/dastak/constants"
	"github.com/soheilkhaledabdi/dastak/services"
)

type AuthHandler struct {
	service *services.AuthService
}

func NewAuthHandler(cfg *config.Config) *AuthHandler {
	service := services.NewAuthService(cfg)
	return &AuthHandler{
		service: service,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	req := new(dto.Login)
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithValidationError(false, err, "Validation error", "Please fill in the data correctly"))
		return
	}
	token, alert, err := h.service.Login(req)

	if alert != nil && err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized,
			helper.GenerateBaseResponseWithError(true, err, alert.Title, alert.Message))
		return
	}
	if err != nil && alert == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized,
			helper.GenerateBaseResponseWithError(false, err, "Internal error", "Please try again or contact support"))
		return
	}

	c.JSON(http.StatusOK, helper.GenerateBaseResponse(token, true, alert.Title, alert.Message))
}

func (h *AuthHandler) Register(c *gin.Context) {
	req := new(dto.Register)
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithValidationError(false, err, "Validation error", "Please fill in the data correctly"))
		return
	}
	token, alert, err := h.service.Register(c, req)

	if alert != nil && err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized,
			helper.GenerateBaseResponseWithError(true, err, alert.Title, alert.Message))
		return
	}
	if err != nil && alert == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized,
			helper.GenerateBaseResponseWithError(false, err, "Internal error", "Please try again or contact support"))
		return
	}

	c.JSON(http.StatusOK, helper.GenerateBaseResponse(token, true, alert.Title, alert.Message))
}

func (h *AuthHandler) ResendPassword(c *gin.Context) {
	req := new(dto.Mobile)
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithValidationError(false, err, "Validation error", "Please fill in the data correctly"))
		return
	}
	alert, err := h.service.ResendPassword(c, *req)

	if alert != nil && err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized,
			helper.GenerateBaseResponseWithError(true, err, alert.Title, alert.Message))
		return
	}
	if err != nil && alert == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized,
			helper.GenerateBaseResponseWithError(false, err, "Internal error", "Please try again or contact support"))
		return
	}

	c.JSON(http.StatusOK, helper.GenerateBaseResponse(nil, true, alert.Title, alert.Message))
}

func (h *AuthHandler) Logout(c *gin.Context) {
	auth := c.GetHeader(constants.AuthorizationHeaderKey)
	token := strings.Split(auth, " ")
	err := h.service.Logout(token[1])
	if err != nil {
		c.JSON(http.StatusUnauthorized, helper.GenerateBaseResponse(nil, false, "Your not Login", "Please login and try again"))

	}

	c.JSON(http.StatusCreated, helper.GenerateBaseResponse(nil, true, "Logout", "We are waiting for your return"))
}
