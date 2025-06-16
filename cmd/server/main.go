package main

import (
	"gamershub/api"
	"gamershub/config"
	_ "gamershub/docs"
	"gamershub/internal/controllers"
	"gamershub/internal/models"
	"gamershub/internal/repositories"
	"github.com/gin-contrib/cors"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"time"
)

// @title GamersHub Matchmaking API
// @version 1.0
// @description API for finding gaming partners and managing matches

// @contact.name API Support
// @contact.email strober2004@gmail.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:4040
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token
// @security BearerAuth
// @name Authorization
func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	db := config.ConnectToDatabase()
	//	auto-migration
	if err := db.AutoMigrate(
		&models.User{},
		&models.Friendship{},
		&models.PlayerForm{},
		&models.Rank{},
		&models.GameSession{},
	); err != nil {
		log.Fatal(err)
	}
	//	repo init
	var userRepository = repositories.NewUserRepository(db)
	var friendshipRepository = repositories.NewFriendshipRepository(db)
	var formRepository = repositories.NewFormRepository(db)

	//	controllers init
	authController := controllers.NewAuthController(userRepository)
	friendshipController := controllers.NewFriendshipController(friendshipRepository, userRepository)
	userController := controllers.NewUserController(userRepository, db)
	formController := controllers.NewPlayerFormController(formRepository)
	matchmakingController := controllers.NewMatchmakingController(formRepository, db)

	//	routes
	router := api.SetupRouter(authController, friendshipController, userController, formController, matchmakingController)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	router.Run(":4040")
}
