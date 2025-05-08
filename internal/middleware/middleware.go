package middleware

import (
	"gamershub/pkg/utils"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Извлекаем токен из заголовка
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is missing",
			})
			return
		}

		// 2. Проверяем формат "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token format. Expected: Bearer <token>",
			})
			return
		}

		tokenString := tokenParts[1]

		// 3. Парсим токен
		claims, err := utils.ParseJWT(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid token",
				"details": err.Error(), // Добавляем детали ошибки для отладки
			})
			return
		}

		// 4. Логируем для отладки (удалите в продакшене)
		log.Printf("[DEBUG] Authenticated user: ID=%d, Role=%s", claims.UserId, claims.Role)
		// 5. Добавляем данные в контекст
		c.Set("userID", claims.UserId)
		c.Set("userRole", claims.Role)

		c.Next()
	}
}
