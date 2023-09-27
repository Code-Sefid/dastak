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
	router.POST("/get-by-filter", middlewares.Authentication(cfg), f.GetAll)
	router.PUT("/:id", middlewares.Authentication(cfg), f.Update)
	router.GET("/:code", f.GetByCode)
	router.DELETE("/:id", middlewares.Authentication(cfg), f.Delete)
	router.POST("/item/:id", middlewares.Authentication(cfg), f.AddItem)
	router.DELETE("/item/:id", middlewares.Authentication(cfg), f.DeleteItem)
}
