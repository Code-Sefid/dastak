package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/soheilkhaledabdi/dastak/config"
	"github.com/soheilkhaledabdi/dastak/services"
)

type ProductsHandler struct {
	service *services.ProductsService
}

func NewProductsHandler(cfg *config.Config) *ProductsHandler {
	return &ProductsHandler{
		service: services.NewProductsService(cfg),
	}
}

func (h *ProductsHandler) Create(c *gin.Context) {
	CreateByUserId(c, h.service.CreateByUserId)
}

func (h *ProductsHandler) Update(c *gin.Context) {
	Update(c, h.service.Update)
}

func (h *ProductsHandler) Delete(c *gin.Context) {
	Delete(c, h.service.Delete)
}

func (h *ProductsHandler) GetById(c *gin.Context) {
	GetById(c, h.service.GetById)
}

func (h *ProductsHandler) GetByFilter(c *gin.Context) {
	GetByFilter(c, h.service.GetByFilter)
}
