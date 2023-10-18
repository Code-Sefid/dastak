package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/soheilkhaledabdi/dastak/api/handlers"
	"github.com/soheilkhaledabdi/dastak/api/middlewares"
	"github.com/soheilkhaledabdi/dastak/config"
)

func CheckOut(r *gin.RouterGroup, cfg *config.Config) {
	h := handlers.NewCheckOutHandler(cfg)
	r.POST("", middlewares.Authentication(cfg), h.CheckOut)
}
