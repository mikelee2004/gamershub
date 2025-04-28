package api

import (
	"gamershub/internal/controllers"
	"gamershub/internal/middleware"
	"gamershub/internal/types"
	"github.com/gin-gonic/gin"
)

func SetupRouter(authController *controllers.AuthController) *gin.Engine {
	router := gin.Default()

	api := router.Group("/api/v1")
	{
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/register", authController.Register)
			authGroup.POST("/login", authController.Login)
		}

		// Защищенные маршруты
		protected := api.Group("/admin")
		protected.Use(middleware.AuthMiddleware(types.RoleAdmin))
		{
			// todo: admin endpoints
		}
	}

	return router
}
