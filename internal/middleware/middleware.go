package middleware

import (
	"gamershub/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is missing",
			})
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token format. Expected: Bearer <token>",
			})
			return
		}

		tokenString := tokenParts[1]

		claims, err := utils.ParseJWT(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid token",
				"details": err.Error(), // Добавляем детали ошибки для отладки
			})
			return
		}
		c.Set("userID", claims.UserId)
		c.Set("userRole", claims.Role)

		c.Next()
	}
}
