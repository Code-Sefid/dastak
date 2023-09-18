package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/soheilkhaledabdi/dastak/api/handlers"
	"github.com/soheilkhaledabdi/dastak/api/middlewares"
	"github.com/soheilkhaledabdi/dastak/config"
)

func Products(r *gin.RouterGroup, cfg *config.Config) {
	h := handlers.NewProductsHandler(cfg)

	r.POST("/", middlewares.Authentication(cfg), h.Create)
	r.PUT("/:id", middlewares.Authentication(cfg), h.Update)
	r.DELETE("/:id", middlewares.Authentication(cfg), h.Delete)
	r.GET("/:id", middlewares.Authentication(cfg), h.GetById)
	r.POST("/get-by-filter", middlewares.Authentication(cfg), h.GetByFilter)
}
