package middlewares

import (
    "github.com/gin-gonic/gin"
)

func Cors() gin.HandlerFunc {
    return func(c *gin.Context) {

        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Headers", "*")
        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Header("Access-Control-Max-Age", "21600")
        c.Set("content-type", "application/json")
        
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        
        c.Next()
    }
}
