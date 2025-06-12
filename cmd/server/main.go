package main

import (
	"gamershub/api"
	_ "gamershub/cmd/server/docs"
	"gamershub/internal/controllers"
	"gamershub/internal/models"
	"gamershub/internal/repositories"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

// @title Gamershub
// @version 1.0
// @description Description of your API
// @host localhost:8080
// @BasePath /api/v1
func main() {
	dsn := "host=db user=postgres password=passpass dbname=gamershub port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	//	auto-migration
	if err := db.AutoMigrate(
		&models.User{},
		&models.Friendship{},
		&models.PlayerForm{},
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
	matchmakingController := controllers.NewMatchmakingController(userRepository, db)
	formController := controllers.NewPlayerFormController(formRepository)

	//	routes
	router := api.SetupRouter(authController, friendshipController, userController, formController, matchmakingController)

	//	launch server
	router.Run(":4040")
}
