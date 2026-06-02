package rest

import (
	"go-service/internal/presentation/rest/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth != "Bearer secret-token" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errors.ErrorSchema{Code: http.StatusUnauthorized, Text: "unauthorized"})
			return
		}
		c.Set("user_id", "123")
		c.Next()
	}
}
