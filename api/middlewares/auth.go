package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/soheilkhaledabdi/dastak/api/helper"
	"github.com/soheilkhaledabdi/dastak/config"
	"github.com/soheilkhaledabdi/dastak/constants"
	"github.com/soheilkhaledabdi/dastak/pkg/service_errors"
	"github.com/soheilkhaledabdi/dastak/services"
)

func Authentication(cfg *config.Config) gin.HandlerFunc {
	var tokenService = services.NewTokenService(cfg)

	return func(c *gin.Context) {
		var err error
		claimMap := map[string]interface{}{}
		auth := c.GetHeader(constants.AuthorizationHeaderKey)
		token := strings.Split(auth, " ")
		if auth == "" {
			err = &service_errors.ServiceError{EndUserMessage: service_errors.TokenRequired}
		} else {
			claimMap, err = tokenService.GetClaims(token[1])
			if err != nil {
				switch e := err.(type) {
				case *jwt.ValidationError:
					switch e.Errors {
					case jwt.ValidationErrorExpired:
						err = &service_errors.ServiceError{EndUserMessage: service_errors.TokenExpired}
					default:
						err = &service_errors.ServiceError{EndUserMessage: service_errors.TokenInvalid}
					}
				default:
					err = &service_errors.ServiceError{EndUserMessage: service_errors.TokenInvalid}
				}
			}
		}
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, helper.GenerateBaseResponseWithError(
				false, err, "Unauthorized", "Please Login to Access Applications",
			))
			return
		}

		c.Set(constants.UserIdKey, claimMap[constants.UserIdKey])
		c.Set(constants.ExpireTimeKey, claimMap[constants.ExpireTimeKey])

		c.Next()
	}
}
