package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/soheilkhaledabdi/dastak/api/handlers"
	"github.com/soheilkhaledabdi/dastak/api/middlewares"
	"github.com/soheilkhaledabdi/dastak/config"
)

func Factors(router *gin.RouterGroup, cfg *config.Config) {
	f := handlers.NewFactorsHandler(cfg)
	router.POST("/", middlewares.Authentication(cfg), f.Create)
	router.GET("/", middlewares.Authentication(cfg), f.GetAll)
}
