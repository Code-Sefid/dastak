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

type FactorDetailHandler struct {
	service *services.FactorDetailService
}

func NewFactorDetailHandler(cfg *config.Config) *FactorDetailHandler {
	return &FactorDetailHandler{
		service: services.NewFactorDetailService(cfg),
	}
}

func (h *FactorDetailHandler) Create(c *gin.Context) {
	Create(c, h.service.CreateOrUpdate)
}

func (h *FactorDetailHandler) Update(c *gin.Context) {
	Update(c, h.service.Update)
}

func (h *FactorDetailHandler) GetById(c *gin.Context) {
	GetById(c, h.service.GetById)
}

func (h *FactorDetailHandler) GetByFilter(c *gin.Context) {
	GetByFilter(c, h.service.GetByFilter)
}

func (p *FactorDetailHandler) FactorPayment(c *gin.Context) {
	req := new(dto.FactorPayment)
	err := c.ShouldBind(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithValidationError(false, err, "لطفا اطلاعات ها را به درستی پر کنید"))
		return
	}

	err = p.service.FactorPayment(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			helper.GenerateBaseResponseWithAnyError(false, err, "لطفا دوباره مجدد امتحان بکنید یا با پشتیبانی ارتباط بگیرید"))
		return
	}

	c.JSON(http.StatusOK, helper.GenerateBaseResponse(nil, true, ""))
}

func (p *FactorDetailHandler) AddTrackingCode(c *gin.Context) {
	req := new(dto.AddTrackingCode)
	err := c.ShouldBind(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithValidationError(false, err, "لطفا اطلاعات ها را به درستی پر کنید"))
		return
	}

	err = p.service.AddTrackingCode(c, req)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound,
			helper.GenerateBaseResponseWithAnyError(false, err, "فاکتور مورد نظر یافت نشد"))
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError,
			helper.GenerateBaseResponseWithAnyError(false, err, "لطفا دوباره مجدد امتحان بکنید یا با پشتیبانی ارتباط بگیرید"))
		return
	}

	c.JSON(http.StatusOK, helper.GenerateBaseResponse(nil, true, ""))
}
