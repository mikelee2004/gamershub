package api

import (
	_ "gamershub/docs"
	"gamershub/internal/controllers"
	"gamershub/internal/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

func SetupRouter(
	authController *controllers.AuthController,
	friendshipCtrl *controllers.FriendshipController,
	userCtrl *controllers.UserController,
	formCtrl *controllers.PlayerFormController,
	mmCtrl *controllers.MatchmakingController,
) *gin.Engine {

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Ваш Next.js адрес
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
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
			formRoutes.POST("/create", formCtrl.CreateForm)
			formRoutes.GET("/me", formCtrl.Get)
			formRoutes.DELETE("/delete", formCtrl.Delete)
		}
		mmRoutes := api.Group("/lft")
		mmRoutes.Use(middleware.AuthRequired())
		{
			mmRoutes.GET("/find", mmCtrl.GetMatches)
			mmRoutes.GET("/invites", mmCtrl.GetInvites)
			mmRoutes.POST("/invite/", mmCtrl.SendInvite)
			mmRoutes.POST("/respond", mmCtrl.RespondToInvitation)
		}
	}
	return router
}
