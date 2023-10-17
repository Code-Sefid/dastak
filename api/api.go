package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/soheilkhaledabdi/dastak/api/middlewares"
	routers "github.com/soheilkhaledabdi/dastak/api/routes"
	"github.com/soheilkhaledabdi/dastak/api/validation"
	"github.com/soheilkhaledabdi/dastak/config"
	"github.com/soheilkhaledabdi/dastak/pkg/logging"
)

var logger = logging.NewLogger(config.GetConfig())

func InitServer() {
	config := config.GetConfig()
	r := gin.New()
	r.Use(middlewares.Cors())
	RegisterValidators()
	r.Use(gin.Logger(), gin.Recovery())
	r.RedirectTrailingSlash = false

	router := r.Group("api/v1")
	{
		router.Static("/files", "./uploads")

		auth := router.Group("auth")
		routers.Auth(auth, config)

		categories := router.Group("categories")
		routers.Categories(categories, config)

		products := router.Group("products")
		routers.Products(products, config)

		factor := router.Group("factors")
		routers.Factors(factor, config)

		bank := router.Group("bank")
		routers.Bank(bank, config)

		customer := router.Group("factor-detail")
		routers.Customer(customer, config)

		payment := router.Group("payment")
		routers.Payment(payment, config)

		transactions := router.Group("transactions")
		routers.Transactions(transactions, config)

		checkOut := router.Group("check-out")
		routers.CheckOut(checkOut, config)
	}

	r.Run(fmt.Sprintf(":%s", config.Server.InternalPort))
}

func RegisterValidators() {
	val, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		err := val.RegisterValidation("mobile", validation.IranianMobileNumberValidator, true)
		if err != nil {
			logger.Error(logging.Validation, logging.Startup, err.Error(), nil)
		}
	}
}
