package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/soheilkhaledabdi/dastak/config"
	"github.com/soheilkhaledabdi/dastak/services"
)

type CategoriesHandler struct {
	service *services.CategoriesService
}

func NewCategoriesHandler(cfg *config.Config) *CategoriesHandler {
	return &CategoriesHandler{
		service: services.NewCategoriesService(cfg),
	}
}

func (h *CategoriesHandler) Create(c *gin.Context) {
	Create(c, h.service.Create)
}

func (h *CategoriesHandler) Update(c *gin.Context) {
	Update(c, h.service.Update)
}

func (h *CategoriesHandler) Delete(c *gin.Context) {
	Delete(c, h.service.Delete)
}

func (h *CategoriesHandler) GetById(c *gin.Context) {
	GetById(c, h.service.GetById)
}

func (h *CategoriesHandler) GetByFilter(c *gin.Context) {
	GetByFilter(c, h.service.GetByFilter)
}
