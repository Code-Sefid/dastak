package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/soheilkhaledabdi/dastak/api/dto"
	"github.com/soheilkhaledabdi/dastak/api/helper"
	"github.com/soheilkhaledabdi/dastak/config"
	"github.com/soheilkhaledabdi/dastak/services"
	"gorm.io/gorm"
)

type PaymentHandler struct {
	service *services.PaymentService
}

func NewPaymentHandler(cfg *config.Config) *PaymentHandler {
	service := services.NewPaymentService(cfg)
	return &PaymentHandler{
		service: service,
	}
}

func (h *PaymentHandler) Create(c *gin.Context) {
	req := new(dto.Payment)
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithValidationError(false, err, "لطفا اطلاعات ها را به درستی پر کنید"))
		return
	}
	res, _, err := h.service.PaymentURL(c, req)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			helper.GenerateBaseResponseWithError(false, err, "لطفا دوباره مجدد امتحان بکنید یا با پشتیبانی ارتباط بگیرید"))
		return
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		c.AbortWithStatusJSON(http.StatusNotFound,
			helper.GenerateBaseResponseWithError(false, err, "لطفا اطلاعات ها را به درستی پر کنید"))
		return
	}

	c.JSON(http.StatusOK, helper.GenerateBaseResponse(res, true, ""))
}

func (h *PaymentHandler) Check(c *gin.Context) {
	req := new(dto.Verify)
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithValidationError(false, err, "لطفا اطلاعات ها را به درستی پر کنید"))
		return
	}
	status, alert, err := h.service.CheckPayment(c, req)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			helper.GenerateBaseResponseWithError(status, err, alert.Message))
		return
	}

	c.JSON(http.StatusOK, helper.GenerateBaseResponse(nil, status, ""))
}
