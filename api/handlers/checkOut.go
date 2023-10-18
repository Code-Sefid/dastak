package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/soheilkhaledabdi/dastak/api/dto"
	"github.com/soheilkhaledabdi/dastak/api/helper"
	"github.com/soheilkhaledabdi/dastak/config"
	"github.com/soheilkhaledabdi/dastak/constants"
	"github.com/soheilkhaledabdi/dastak/services"
)

type CheckOutHandler struct {
	service *services.CheckOutService
}

func NewCheckOutHandler(cfg *config.Config) *CheckOutHandler {
	return &CheckOutHandler{
		service: services.NewCheckOutService(cfg),
	}
}

func (p *CheckOutHandler) CheckOut(c *gin.Context) {
	userID := int(c.Value(constants.UserIdKey).(float64))
	req := new(dto.CheckOut)
	err := c.ShouldBind(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithValidationError(false, err, "لطفا اطلاعات ها را به درستی پر کنید"))
		return
	}

	alert, status, err := p.service.CheckOutMony(c, userID, *req)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			helper.GenerateBaseResponseWithAnyError(false, err, "لطفا دوباره مجدد امتحان بکنید یا با پشتیبانی ارتباط بگیرید"))
		return
	} else if err != nil && alert != nil {
		c.JSON(http.StatusOK, helper.GenerateBaseResponse(nil, status, alert.Message))
		return
	} else if err == nil && alert != nil {
		c.JSON(http.StatusOK, helper.GenerateBaseResponse(nil, status, alert.Message))
		return
	}

	c.JSON(http.StatusOK, helper.GenerateBaseResponse(nil, status, ""))
}
