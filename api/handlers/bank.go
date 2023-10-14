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

type BankHandler struct {
	service *services.BankService
}

func NewBankHandler(cfg *config.Config) *BankHandler {
	return &BankHandler{
		service: services.NewBankService(cfg),
	}
}

func (h *BankHandler) Create(c *gin.Context) {
	userID := int(c.Value(constants.UserIdKey).(float64))
	req := new(dto.CreateBankRequest)
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithValidationError(false, err, "لطفا اطلاعات ها را به درستی پر کنید"))
		return
	}
	res, err := h.service.CreateOrUpdate(c, req, userID)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			helper.GenerateBaseResponseWithError(false, err, "لطفا دوباره مجدد امتحان بکنید یا با پشتیبانی ارتباط بگیرید"))
		return
	}

	c.JSON(http.StatusOK, helper.GenerateBaseResponse(res, true, ""))
}

func (h *BankHandler) Update(c *gin.Context) {
	Update(c, h.service.Update)
}

func (h *BankHandler) GetByUserId(c *gin.Context) {
	GetByUserId(c, h.service.GetByUserId)
}
