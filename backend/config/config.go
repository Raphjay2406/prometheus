// prometheus/backend/config/config.go
package config

import (
	"os"
	"strconv" // For converting string to int

	"github.com/joho/godotenv"
)

// Config holds the application configuration
type Config struct {
	AppEnv             string
	Port               string
	DBHost             string
	DBPort             string
	DBUser             string
	DBPassword         string
	DBName             string
	JWTSecret          string
	JWTExpirationHours int // Added for JWT expiration
	GodAdminEmail      string
	GodAdminPassword   string
}

// LoadConfig reads configuration from environment variables or .env file
func LoadConfig() (*Config, error) {
	// Load .env file if it exists.
	// These paths cover running `go run cmd/main.go` from `backend/` or from project root `prometheus/`.
	if _, err := os.Stat(".env"); err == nil {
		godotenv.Load(".env") // Load from current directory (e.g., backend/)
	} else if _, err := os.Stat("../.env"); err == nil {
		godotenv.Load("../.env") // Load from project root if in backend/
	}
	// If running from `prometheus/` directly, `godotenv.Load("backend/.env")` might be needed if .env is in backend.
	// For simplicity, it's often best to have one .env at the project root or ensure it's found by the binary's CWD.

	jwtExpHours, err := strconv.Atoi(getEnv("JWT_EXPIRATION_HOURS", "168")) // Default to 7 days (24*7)
	if err != nil {
		jwtExpHours = 168 // Fallback default if conversion fails
	}

	return &Config{
		AppEnv:             getEnv("APP_ENV", "development"),
		Port:               getEnv("PORT", "8080"),
		DBHost:             getEnv("DB_HOST", "localhost"),
		DBPort:             getEnv("DB_PORT", "5432"),
		DBUser:             getEnv("DB_USER", "prometheus_user"),
		DBPassword:         getEnv("DB_PASSWORD", "prometheus_password"),
		DBName:             getEnv("DB_NAME", "prometheus_db"),
		JWTSecret:          getEnv("JWT_SECRET", "your_super_secret_jwt_key_that_is_very_long_and_secure"),
		JWTExpirationHours: jwtExpHours, // Added
		GodAdminEmail:      getEnv("GOD_ADMIN_EMAIL", "godadmin@example.com"),
		GodAdminPassword:   getEnv("GOD_ADMIN_PASSWORD", "SecureGodAdminP@ssw0rd123!"),
	}, nil
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
