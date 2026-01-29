package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth != "Bearer secret-token" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorSchema{http.StatusUnauthorized, "unauthorized"})
			return
		}
		c.Set("user_id", "123")
		c.Next()
	}
}
