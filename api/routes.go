package api

import (
	_ "gamershub/cmd/server/docs"
	"gamershub/internal/controllers"
	"gamershub/internal/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(
	authController *controllers.AuthController,
	friendshipCtrl *controllers.FriendshipController,
	userCtrl *controllers.UserController,
	formCtrl *controllers.PlayerFormController,
	matchmakingCtrl *controllers.MatchmakingController,
) *gin.Engine {

	router := gin.Default()
	// add swagger
	url := ginSwagger.URL("/swagger/doc.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	api := router.Group("/api/v1")
	{
		public := api.Group("/auth")
		{
			public.POST("/register", authController.Register)
			public.POST("/login", authController.Login)
		}
		userRoutes := api.Group("/user")
		userRoutes.Use(middleware.AuthRequired())
		{
			userRoutes.GET("/profile", userCtrl.GetProfile)
			userRoutes.POST("/logout", authController.Logout)
			userRoutes.POST("/forgotpassword", userCtrl.ChangePassword)
		}
		friendshipRoutes := api.Group("/friends")
		friendshipRoutes.Use(middleware.AuthRequired())
		{
			friendshipRoutes.POST("/add/:friend_id", friendshipCtrl.SendFriendRequest)
			friendshipRoutes.PUT("/accept/:friend_id/accept", friendshipCtrl.AcceptFriendRequest)
			friendshipRoutes.GET("/friendlist", friendshipCtrl.GetFriends)
			friendshipRoutes.GET("/requests", friendshipCtrl.GetPendingRequests)
		}

		formRoutes := api.Group("/form")
		formRoutes.Use(middleware.AuthRequired())
		{
			formRoutes.POST("/newform", formCtrl.CreateForm)
			formRoutes.GET("/myform", formCtrl.Get)
			formRoutes.POST("/delete", formCtrl.Delete)
		}
		//matchmakingRoutes := api.Group("/matchmaking")
		//matchmakingRoutes.Use(middleware.AuthRequired())
		//{
		//	matchmakingRoutes.GET("find", matchmakingCtrl.FindMatches)
		//	matchmakingRoutes.POST("accept")
		//
		//}

	}
	return router
}
