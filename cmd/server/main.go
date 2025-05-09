package main

import (
	"gamershub/api"
	"gamershub/internal/controllers"
	"gamershub/internal/models"
	"gamershub/internal/repositories"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func main() {
	dsn := "host=localhost user=postgres password=3791 dbname=gamershub port=5432"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	//	auto-migration
	if err := db.AutoMigrate(
		&models.User{},
		&models.Friendship{},
	); err != nil {
		log.Fatal(err)
	}

	//	repo init
	var userRepository = repositories.NewUserRepository(db)
	var friendshipRepository = repositories.NewFriendshipRepository(db)

	//	controllers init
	authController := controllers.NewAuthController(userRepository)
	friendshipController := controllers.NewFriendshipController(friendshipRepository, userRepository)

	//	routes
	router := api.SetupRouter(authController, friendshipController)

	//	launch server
	if err := router.Run(":5050"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
