package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/soheilkhaledabdi/dastak/api/helper"
	"github.com/soheilkhaledabdi/dastak/config"
	"github.com/soheilkhaledabdi/dastak/constants"
	"github.com/soheilkhaledabdi/dastak/services"
)

type TransactionsService struct {
	service *services.TransactionsService
}

func NewTransactionsHandler(cfg *config.Config) *TransactionsService {
	service := services.NewTransactionsService(cfg)
	return &TransactionsService{
		service: service,
	}
}

func (h *TransactionsService) GetAll(c *gin.Context) {
	userID := int(c.Value(constants.UserIdKey).(float64))

	res, err := h.service.GetByFilter(c, userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			helper.GenerateBaseResponseWithError(false, err, "لطفا دوباره مجدد امتحان بکنید یا با پشتیبانی ارتباط بگیرید"))
		return
	}

	c.JSON(http.StatusOK, helper.GenerateBaseResponse(res, true, ""))
}
