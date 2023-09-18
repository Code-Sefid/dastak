package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/soheilkhaledabdi/dastak/api/handlers"
	"github.com/soheilkhaledabdi/dastak/api/middlewares"
	"github.com/soheilkhaledabdi/dastak/config"
)

func Auth(router *gin.RouterGroup, cfg *config.Config) {
	h := handlers.NewAuthHandler(cfg)
	router.POST("/register", h.Register)
	router.POST("/login", h.Login)
	router.POST("/resend-password", h.ResendPassword)
	router.POST("/logout", middlewares.Authentication(cfg), h.Logout)
}
