package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/soheilkhaledabdi/dastak/api/dto"
	"github.com/soheilkhaledabdi/dastak/api/helper"
	"github.com/soheilkhaledabdi/dastak/config"
	"github.com/soheilkhaledabdi/dastak/services"
	"gorm.io/gorm"
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
	id, _ := strconv.Atoi(c.Params.ByName("id"))
	req := new(dto.UpdateProductsRequest)
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithValidationError(false, err, "لطفا اطلاعات ها را به درستی پر کنید"))
		return
	}
	res , err := h.service.Update(c,req,id)
	if errors.Is(err,gorm.ErrRecordNotFound){
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			helper.GenerateBaseResponseWithError(false, err, ":) چیزی پیدا نکردم گلم اشتباه زدی"))
		return
	}else if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			helper.GenerateBaseResponseWithError(false, err, "لطفا دوباره مجدد امتحان بکنید یا با پشتیبانی ارتباط بگیرید"))
		return
	}

	c.JSON(http.StatusOK, helper.GenerateBaseResponse(res, true, ""))
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
