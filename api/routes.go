package api

import (
	"gamershub/internal/controllers"
	"gamershub/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter(
	authController *controllers.AuthController,
	friendshipCtrl *controllers.FriendshipController,
	userCtrl *controllers.UserController,
	formCtrl *controllers.PlayerFormController,
) *gin.Engine {

	router := gin.Default()
	api := router.Group("/api/v1")
	{
		public := api.Group("/auth")
		{
			public.POST("/register", authController.Register)
			public.POST("/login", authController.Login)
		}
		friendshipRoutes := api.Group("/friends")
		friendshipRoutes.Use(middleware.AuthRequired())
		{
			friendshipRoutes.POST("/add/:friend_id", friendshipCtrl.SendFriendRequest)
			friendshipRoutes.PUT("/accept/:friend_id/accept", friendshipCtrl.AcceptFriendRequest)
			friendshipRoutes.GET("/friendlist", friendshipCtrl.GetFriends)
			friendshipRoutes.GET("/requests", friendshipCtrl.GetPendingRequests)
		}
		userRoutes := api.Group("/user")
		userRoutes.Use(middleware.AuthRequired())
		{
			userRoutes.GET("/profile", userCtrl.GetProfile)
			userRoutes.POST("/logout", authController.Logout)
			userRoutes.POST("/forgotpassword", userCtrl.ChangePassword)
		}
		formRoutes := api.Group("/form")
		formRoutes.Use(middleware.AuthRequired())
		{
			formRoutes.POST("/newform", formCtrl.CreateForm)
			formRoutes.GET("/myform", formCtrl.Get)
			formRoutes.POST("/delete", formCtrl.Delete)
		}
	}
	return router
}
