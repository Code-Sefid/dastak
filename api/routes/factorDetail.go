package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/soheilkhaledabdi/dastak/api/handlers"
	"github.com/soheilkhaledabdi/dastak/config"
)

func Customer(r *gin.RouterGroup, cfg *config.Config) {
	h := handlers.NewFactorDetailHandler(cfg)

	r.POST("/", h.Create)
	r.PUT("/:id", h.Update)
	r.GET("/:id", h.GetById)
	r.POST("/get-by-filter", h.GetByFilter)
	r.POST("/payment", h.FactorPayment)
}
