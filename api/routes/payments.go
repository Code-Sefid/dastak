package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/soheilkhaledabdi/dastak/api/handlers"
	"github.com/soheilkhaledabdi/dastak/config"
)

func Payment(r *gin.RouterGroup, cfg *config.Config) {
	h := handlers.NewPaymentHandler(cfg)
	r.POST("/", h.Create)
	r.POST("/check", h.Check)
}
