package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/soheilkhaledabdi/dastak/api/handlers"
	"github.com/soheilkhaledabdi/dastak/api/middlewares"
	"github.com/soheilkhaledabdi/dastak/config"
)

func Bank(r *gin.RouterGroup, cfg *config.Config) {
	h := handlers.NewBankHandler(cfg)

	r.POST("/", middlewares.Authentication(cfg), h.Create)
	r.PUT("/:id", middlewares.Authentication(cfg), h.Update)
	r.GET("/", middlewares.Authentication(cfg), h.GetByUserId)
}
