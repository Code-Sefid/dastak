package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/soheilkhaledabdi/dastak/config"
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
	CreateByUserId(c, h.service.CreateByUserId)
}

func (h *BankHandler) Update(c *gin.Context) {
	Update(c, h.service.Update)
}


func (h *BankHandler) GetByUserId(c *gin.Context) {
	GetByUserId(c, h.service.GetByUserId)
}
