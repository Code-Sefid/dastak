package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/soheilkhaledabdi/dastak/api/dto"
	"github.com/soheilkhaledabdi/dastak/api/helper"
	"github.com/soheilkhaledabdi/dastak/config"
	"github.com/soheilkhaledabdi/dastak/constants"
	"github.com/soheilkhaledabdi/dastak/services"
)

type FactorService struct {
	service *services.FactorService
}

func NewFactorsHandler(cfg *config.Config) *FactorService {
	service := services.NewFactorService(cfg)
	return &FactorService{
		service: service,
	}
}

func (h *FactorService) Create(c *gin.Context) {
	userID := int(c.Value(constants.UserIdKey).(float64))
	req := new(dto.CreateFactor)
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithValidationError(false, err, "Validation error", "Please fill in the data correctly"))
		return
	}
	err = h.service.Create(c,userID, *req)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			helper.GenerateBaseResponseWithError(false, err, "Internal error", "Please try again or contact support"))
		return
	}

	c.JSON(http.StatusOK, helper.GenerateBaseResponse(nil, true, "", ""))
}


func (h *FactorService) GetAll(c *gin.Context) {
	userID := int(c.Value(constants.UserIdKey).(float64))

	res , err := h.service.GetAll(c,userID)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			helper.GenerateBaseResponseWithError(false, err, "Internal error", "Please try again or contact support"))
		return
	}

	c.JSON(http.StatusOK, helper.GenerateBaseResponse(res, true, "", ""))
}

func (h *FactorService) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Params.ByName("id"))
	req := new(dto.UpdateFactor)
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithValidationError(false, err, "Validation error", "Please fill in the data correctly"))
		return
	}
	err = h.service.Update(c,id,*req)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			helper.GenerateBaseResponseWithError(false, err, "Internal error", "Please try again or contact support"))
		return
	}

	c.JSON(http.StatusOK, helper.GenerateBaseResponse(nil, true, "", ""))
}

func (h *FactorService) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Params.ByName("id"))

	err := h.service.Delete(c,id,)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			helper.GenerateBaseResponseWithError(false, err, "Internal error", "Please try again or contact support"))
		return
	}

	c.JSON(http.StatusOK, helper.GenerateBaseResponse(nil, true, "", ""))
}

func (h *FactorService) GetByCode(c *gin.Context) {
    code := c.Params.ByName("code")

	res , err := h.service.GetByCode(c,code)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			helper.GenerateBaseResponseWithError(false, err, "Internal error", "Please try again or contact support"))
		return
	}

    c.JSON(http.StatusOK, helper.GenerateBaseResponse(res, true, "", ""))
}
