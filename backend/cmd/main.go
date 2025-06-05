// prometheus/backend/cmd/main.go
package main

import (
	"fmt"
	"log"
	"prometheus/backend/config"
	"prometheus/backend/database"
	"prometheus/backend/internal/auth" // Import auth package for User model
	"prometheus/backend/internal/role" // Import role package for Role model
	"prometheus/backend/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	_ = godotenv.Load("../.env")

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error: Failed to load configuration: %v", err)
	}

	db, err := database.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Error: Failed to connect to the database: %v", err)
	}
	log.Println("Database connected successfully.")

	log.Println("Running database auto-migrations...")
	err = db.AutoMigrate(
		&auth.User{},
		&role.Role{},
	)
	if err != nil {
		log.Fatalf("Error: Failed to auto-migrate database schema: %v", err)
	}
	log.Println("Database auto-migrations completed successfully.")

	// Seed the database with initial data (roles, god admin)
	// This should run after migrations to ensure tables exist.
	log.Println("Starting database seeding process...")
	if err := database.SeedRoles(db); err != nil {
		// Log the error but don't necessarilyFatalf, as the app might still run
		// depending on how critical initial roles are for startup vs. dynamic creation.
		log.Printf("Error during role seeding: %v", err)
	} else {
		log.Println("Role seeding completed.")
	}

	if err := database.SeedGodAdmin(db, cfg); err != nil {
		log.Printf("Error during god admin seeding: %v", err)
	} else {
		log.Println("God Admin user seeding process completed.")
	}
	log.Println("Database seeding process finished.")

	router := gin.Default()
	routes.SetupRoutes(router, db, cfg)

	serverAddr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Server starting on http://localhost%s (AppEnv: %s)", serverAddr, cfg.AppEnv)

	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Error: Failed to start server: %v", err)
	}
}
