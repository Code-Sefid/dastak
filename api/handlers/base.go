package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/soheilkhaledabdi/dastak/api/helper"
	"github.com/soheilkhaledabdi/dastak/config"
	"github.com/soheilkhaledabdi/dastak/constants"
	"github.com/soheilkhaledabdi/dastak/pkg/logging"
)

var logger = logging.NewLogger(config.GetConfig())

func Create[Ti any, To any](c *gin.Context, caller func(ctx context.Context, req *Ti) (*To, error)) {
	req := new(Ti)
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithValidationError(false, err, "Internal error"))
		return
	}

	res, err := caller(c, req)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(false, err, "Internal error"))
		return
	}
	c.JSON(http.StatusCreated, helper.GenerateBaseResponse(res, true, ""))
}

func CreateByUserId[Ti any, To any](c *gin.Context, caller func(ctx context.Context, req Ti, userID int) (To, error)) {
	userID := int(c.Value(constants.UserIdKey).(float64))
	req := new(Ti)
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithValidationError(false, err, "Internal error"))
		return
	}

	res, err := caller(c, *req, userID)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(false, err, "Internal error"))
		return
	}
	c.JSON(http.StatusCreated, helper.GenerateBaseResponse(res, true, ""))
}

func Update[Ti any, To any](c *gin.Context, caller func(ctx context.Context, id int, req *Ti) (*To, error)) {
	id, _ := strconv.Atoi(c.Params.ByName("id"))
	req := new(Ti)
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithValidationError(false, err, "Internal error"))
		return
	}

	res, err := caller(c, id, req)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(false, err, "Internal error"))
		return
	}
	c.JSON(http.StatusOK, helper.GenerateBaseResponse(res, true, ""))
}

func Delete(c *gin.Context, caller func(ctx context.Context, id int) error) {
	id, _ := strconv.Atoi(c.Params.ByName("id"))
	if id == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound,
			helper.GenerateBaseResponse(nil, false, "Internal error"))
		return
	}

	err := caller(c, id)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(false, err, "Internal error"))
		return
	}
	c.JSON(http.StatusOK, helper.GenerateBaseResponse(nil, true, ""))
}

func GetById[To any](c *gin.Context, caller func(c context.Context, id int) (*To, error)) {
	id, _ := strconv.Atoi(c.Params.ByName("id"))
	if id == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound,
			helper.GenerateBaseResponse(nil, false, "Internal error"))
		return
	}

	res, err := caller(c, id)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(false, err, "Internal error"))
		return
	}
	c.JSON(http.StatusOK, helper.GenerateBaseResponse(res, true, ""))
}

func GetByUserId[To any](c *gin.Context, caller func(c context.Context, id int) (*To, error)) {
	userID := int(c.Value(constants.UserIdKey).(float64))
	if userID == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound,
			helper.GenerateBaseResponse(nil, false, "Internal error"))
		return
	}

	res, err := caller(c, userID)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(false, err, "Internal error"))
		return
	}
	c.JSON(http.StatusOK, helper.GenerateBaseResponse(res, true, ""))
}

func GetByFilter[Ti any, To any](c *gin.Context, caller func(c context.Context, req *Ti, userId int) (*To, error)) {
	userID := int(c.Value(constants.UserIdKey).(float64))
	req := new(Ti)
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithValidationError(false, err, "Internal error"))
		return
	}

	res, err := caller(c, req, userID)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(false, err, "Internal error"))
		return
	}
	c.JSON(http.StatusOK, helper.GenerateBaseResponse(res, true, ""))
}
