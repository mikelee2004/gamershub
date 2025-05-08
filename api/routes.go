package api

import (
	"gamershub/internal/controllers"
	"gamershub/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter(authController *controllers.AuthController, friendshipCtrl *controllers.FriendshipController) *gin.Engine {
	router := gin.Default()
	api := router.Group("/api/v1")
	{
		public := api.Group("/auth")
		{
			public.POST("/register", authController.Register)
			public.POST("/login", authController.Login)
		}
		private := api.Group("/friendship")
		private.Use(middleware.AuthRequired())
		{
			private.POST("/add/:friend_id", friendshipCtrl.SendFriendRequest)
			private.PUT("/accept/:friend_id/accept", friendshipCtrl.AcceptFriendRequest)
			private.GET("/friendlist", friendshipCtrl.GetFriends)
			private.GET("/requests", friendshipCtrl.GetPendingRequests)
		}
	}

	return router
}
