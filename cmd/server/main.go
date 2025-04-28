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
	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatal(err)
	}

	//	repo init
	var userRepository = repositories.NewUserRepository(db)

	//	controllers init
	authController := controllers.NewAuthController(userRepository)

	//	routes
	router := api.SetupRouter(authController)

	//	launch server
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
