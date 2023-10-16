package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/soheilkhaledabdi/dastak/api/handlers"
	"github.com/soheilkhaledabdi/dastak/api/middlewares"
	"github.com/soheilkhaledabdi/dastak/config"
)

func Transactions(r *gin.RouterGroup, cfg *config.Config) {
	h := handlers.NewTransactionsHandler(cfg)
	r.POST("", middlewares.Authentication(cfg), h.GetAll)
}
