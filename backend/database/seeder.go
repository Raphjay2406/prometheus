// prometheus/backend/database/seeder.go
package database

import (
	"errors"
	"fmt"
	"log"
	"prometheus/backend/config"
	"prometheus/backend/internal/auth" // For auth.User model and HashPassword
	"prometheus/backend/internal/role" // For role.Role model

	"gorm.io/gorm"
)

// SeedRoles creates predefined roles in the database if they don't already exist.
func SeedRoles(db *gorm.DB) error {
	rolesToSeed := []role.Role{
		{Name: "staff", Description: "Regular employee with basic access."},
		{Name: "manager", Description: "Managerial role with oversight of a team/department."},
		{Name: "hr", Description: "Human Resources personnel with access to employee data and HR functions."},
		{Name: "admin", Description: "System administrator with broad access, excluding god-level operations."},
		{Name: "god-admin", Description: "Super administrator with unrestricted access to all system functionalities."},
	}

	log.Println("Seeding roles...")
	var count int64
	for _, r := range rolesToSeed {
		// Check if role already exists
		err := db.Model(&role.Role{}).Where("name = ?", r.Name).Count(&count).Error
		if err != nil {
			log.Printf("Error counting role %s: %v\n", r.Name, err)
			continue // Skip to next role on error
		}

		if count == 0 {
			// Role does not exist, create it
			if err := db.Create(&r).Error; err != nil {
				log.Printf("Error creating role %s: %v\n", r.Name, err)
			} else {
				log.Printf("Role '%s' seeded successfully with ID %d.\n", r.Name, r.ID)
			}
		} else {
			log.Printf("Role '%s' already exists. Skipping.\n", r.Name)
		}
	}
	log.Println("Role seeding process completed.")
	return nil // Can be enhanced to return aggregated errors
}

// SeedGodAdmin creates a god-level administrator user if one doesn't exist.
// This function assumes roles have already been seeded, especially the "god-admin" role.
func SeedGodAdmin(db *gorm.DB, cfg *config.Config) error {
	log.Println("Attempting to seed God Admin user...")

	// 1. Check if god admin email is configured
	if cfg.GodAdminEmail == "" || cfg.GodAdminPassword == "" {
		log.Println("GodAdminEmail or GodAdminPassword not configured in .env. Skipping God Admin seed.")
		return nil
	}

	// 2. Find the "god-admin" role
	var godAdminRole role.Role
	if err := db.Where("name = ?", "god-admin").First(&godAdminRole).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("Error: 'god-admin' role not found. Ensure roles are seeded before seeding God Admin.")
			return fmt.Errorf("'god-admin' role not found: %w", err)
		}
		log.Printf("Error fetching 'god-admin' role: %v\n", err)
		return fmt.Errorf("error fetching 'god-admin' role: %w", err)
	}
	log.Printf("'god-admin' role found with ID: %d\n", godAdminRole.ID)

	// 3. Check if a user with the god admin email already exists
	var existingUser auth.User
	err := db.Model(&auth.User{}).Where("email = ?", cfg.GodAdminEmail).First(&existingUser).Error
	if err == nil {
		// User with this email already exists
		log.Printf("User with email '%s' (ID: %d) already exists. Ensuring it has 'god-admin' role.", cfg.GodAdminEmail, existingUser.ID)
		// Optionally, ensure this existing user has the god-admin role
		if existingUser.RoleID != godAdminRole.ID {
			log.Printf("Updating user %s (ID: %d) to 'god-admin' role (ID: %d)", existingUser.Username, existingUser.ID, godAdminRole.ID)
			existingUser.RoleID = godAdminRole.ID
			if err := db.Save(&existingUser).Error; err != nil {
				log.Printf("Failed to update existing user %s to 'god-admin' role: %v", existingUser.Username, err)
				return fmt.Errorf("failed to update existing user to 'god-admin': %w", err)
			}
		}
		return nil // God admin (or user with that email) already exists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// A different database error occurred
		log.Printf("Database error when checking for existing god admin user: %v\n", err)
		return fmt.Errorf("db error checking existing god admin: %w", err)
	}

	// 4. User does not exist, create the god admin user
	hashedPassword, err := auth.HashPassword(cfg.GodAdminPassword)
	if err != nil {
		log.Printf("Error hashing god admin password: %v\n", err)
		return fmt.Errorf("error hashing god admin password: %w", err)
	}

	godAdminUser := auth.User{
		Username: "godadmin", // Or derive from email, or make configurable
		Email:    cfg.GodAdminEmail,
		Password: hashedPassword,
		RoleID:   godAdminRole.ID,
		IsActive: true,
	}

	if err := db.Create(&godAdminUser).Error; err != nil {
		log.Printf("Error creating god admin user: %v\n", err)
		return fmt.Errorf("error creating god admin user: %w", err)
	}

	log.Printf("God Admin user '%s' (Email: %s) seeded successfully with ID %d and Role ID %d.\n", godAdminUser.Username, godAdminUser.Email, godAdminUser.ID, godAdminUser.RoleID)
	return nil
}
