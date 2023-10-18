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
			helper.GenerateBaseResponseWithValidationError(false, err, "لطفا اطلاعات ها را به درستی پر کنید"))
		return
	}
	token, alert, status, err := h.service.Login(req)

	if alert != nil && err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized,
			helper.GenerateBaseResponseWithError(status, err, alert.Message))
		return
	}
	if err != nil && alert == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized,
			helper.GenerateBaseResponseWithError(status, err, "لطفا دوباره مجدد امتحان بکنید یا با پشتیبانی ارتباط بگیرید"))
		return
	}

	c.JSON(http.StatusOK, helper.GenerateBaseResponse(token, status, alert.Message))
}

func (h *AuthHandler) Register(c *gin.Context) {
	req := new(dto.Register)
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithValidationError(false, err, "لطفا اطلاعات ها را به درستی پر کنید"))
		return
	}
	token, alert, status, err := h.service.Register(c, req)

	if alert != nil && err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized,
			helper.GenerateBaseResponseWithError(status, err, alert.Message))
		return
	}
	if err != nil && alert == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized,
			helper.GenerateBaseResponseWithError(status, err, "لطفا دوباره مجدد امتحان بکنید یا با پشتیبانی ارتباط بگیرید"))
		return
	}

	if alert == nil {
		c.JSON(http.StatusOK, helper.GenerateBaseResponse(token, status, ""))
	} else {
		c.JSON(http.StatusOK, helper.GenerateBaseResponse(token, status, alert.Message))
	}
}

func (h *AuthHandler) ResendPassword(c *gin.Context) {
	req := new(dto.Mobile)
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithValidationError(false, err, "لطفا اطلاعات ها را به درستی پر کنید"))
		return
	}
	alert, status, err := h.service.ResendPassword(c, *req)

	if alert != nil && err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized,
			helper.GenerateBaseResponseWithError(status, err, alert.Message))
		return
	}
	if err != nil && alert == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized,
			helper.GenerateBaseResponseWithError(status, err, "لطفا دوباره مجدد امتحان بکنید یا با پشتیبانی ارتباط بگیرید"))
		return
	}

	c.JSON(http.StatusOK, helper.GenerateBaseResponse(status, true, alert.Message))
}

func (h *AuthHandler) CheckMobile(c *gin.Context) {
	req := new(*dto.Mobile)
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithValidationError(false, err, "لطفا اطلاعات ها را به درستی پر کنید"))
		return
	}
	alert, status, err := h.service.CheckMobile(*req)

	if alert != nil && err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized,
			helper.GenerateBaseResponseWithError(status, err, alert.Message))
		return
	} else if err != nil && alert == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized,
			helper.GenerateBaseResponseWithError(status, err, "لطفا دوباره مجدد امتحان بکنید یا با پشتیبانی ارتباط بگیرید"))
		return
	}

	c.JSON(http.StatusOK, helper.GenerateBaseResponse(status, true, alert.Message))
}

func (h *AuthHandler) Logout(c *gin.Context) {
	auth := c.GetHeader(constants.AuthorizationHeaderKey)
	token := strings.Split(auth, " ")
	err := h.service.Logout(token[1])
	if err != nil {
		c.JSON(http.StatusUnauthorized, helper.GenerateBaseResponse(nil, false, "لطفا دوباره وارد شوید"))

	}

	c.JSON(http.StatusCreated, helper.GenerateBaseResponse(nil, true, "با موفقیت خارج شدید"))
}
