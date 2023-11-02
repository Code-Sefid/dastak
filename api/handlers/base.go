package handlers

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/soheilkhaledabdi/dastak/api/helper"
	"github.com/soheilkhaledabdi/dastak/constants"
	"gorm.io/gorm"
)


func Create[Ti any, To any](c *gin.Context, caller func(ctx context.Context, req *Ti) (*To, error)) {
	req := new(Ti)
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithValidationError(false, err, "لطفا اطلاعات ها را به درستی پر کنید"))
		return
	}

	res, err := caller(c, req)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(false, err, "مشکلی پیش امده لطفا مجدد امتحان کنید یا با پشتیبانی تماس بگیرید"))
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
			helper.GenerateBaseResponseWithValidationError(false, err, "لطفا اطلاعات ها را به درستی پر کنید"))
		return
	}

	res, err := caller(c, *req, userID)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(false, err, "مشکلی پیش امده لطفا مجدد امتحان کنید یا با پشتیبانی تماس بگیرید"))
		return
	}
	c.JSON(http.StatusCreated, helper.GenerateBaseResponse(res, true, ""))
}

func Update[Ti any, To any](c *gin.Context, caller func(ctx context.Context, id int, userID int, req *Ti) (*To, error)) {
	userID := int(c.Value(constants.UserIdKey).(float64))
	id, _ := strconv.Atoi(c.Params.ByName("id"))
	req := new(Ti)
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithValidationError(false, err, "لطفا اطلاعات ها را به درستی پر کنید"))
		return
	}

	res, err := caller(c, id, userID, req)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(false, err, "مشکلی پیش امده لطفا مجدد امتحان کنید یا با پشتیبانی تماس بگیرید"))
		return
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		c.AbortWithStatusJSON(http.StatusNotFound,
			helper.GenerateBaseResponseWithError(false, err, "موردی با این مشخصات یافت نشد"))
		return
	}
	c.JSON(http.StatusOK, helper.GenerateBaseResponse(res, true, ""))
}

func Delete(c *gin.Context, caller func(ctx context.Context, id int, userID int) error) {
	userID := int(c.Value(constants.UserIdKey).(float64))
	id, _ := strconv.Atoi(c.Params.ByName("id"))
	if id == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound,
			helper.GenerateBaseResponse(nil, false, "پیدا نشد"))
		return
	}

	err := caller(c, id, userID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(false, err, "مشکلی پیش امده لطفا مجدد امتحان کنید یا با پشتیبانی تماس بگیرید"))
		return
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		c.AbortWithStatusJSON(http.StatusNotFound,
			helper.GenerateBaseResponseWithError(false, err, "موردی با این مشخصات یافت نشد"))
		return
	}
	c.JSON(http.StatusOK, helper.GenerateBaseResponse(nil, true, ""))
}

func GetById[To any](c *gin.Context, caller func(c context.Context, id int, userID int) (*To, error)) {
	userID := int(c.Value(constants.UserIdKey).(float64))
	id, _ := strconv.Atoi(c.Params.ByName("id"))
	if id == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound,
			helper.GenerateBaseResponse(nil, false, "پیدا نشد"))
		return
	}

	res, err := caller(c, id, userID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(false, err, "مشکلی پیش امده لطفا مجدد امتحان کنید یا با پشتیبانی تماس بگیرید"))
		return
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		c.AbortWithStatusJSON(http.StatusNotFound,
			helper.GenerateBaseResponseWithError(false, err, "موردی با این مشخصات یافت نشد"))
		return
	}
	c.JSON(http.StatusOK, helper.GenerateBaseResponse(res, true, ""))
}

func GetByUserId[To any](c *gin.Context, caller func(c context.Context, id int) (*To, error)) {
	userID := int(c.Value(constants.UserIdKey).(float64))
	if userID == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound,
			helper.GenerateBaseResponse(nil, false, "پیدا نشد"))
		return
	}

	res, err := caller(c, userID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(false, err, "مشکلی پیش امده لطفا مجدد امتحان کنید یا با پشتیبانی تماس بگیرید"))
		return
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		c.AbortWithStatusJSON(http.StatusNotFound,
			helper.GenerateBaseResponseWithError(false, err, ""))
		return
	}
	if res != nil {
		c.JSON(http.StatusOK, helper.GenerateBaseResponse(res, true, ""))
		return
	}
	c.JSON(http.StatusOK, helper.GenerateBaseResponse(res, false, ""))
}

func GetByFilter[Ti any, To any](c *gin.Context, caller func(c context.Context, req *Ti, userId int) (*To, error)) {
	userID := int(c.Value(constants.UserIdKey).(float64))
	req := new(Ti)
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithValidationError(false, err, "لطفا اطلاعات ها را به درستی پر کنید"))
		return
	}

	res, err := caller(c, req, userID)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(false, err, "مشکلی پیش امده لطفا مجدد امتحان کنید یا با پشتیبانی تماس بگیرید"))
		return
	}
	c.JSON(http.StatusOK, helper.GenerateBaseResponse(res, true, ""))
}
