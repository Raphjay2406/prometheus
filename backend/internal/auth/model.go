// prometheus/backend/internal/auth/model.go
package auth

import (
	"time"

	"prometheus/backend/internal/role" // Import the role package

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// User represents a user account in the system.
type User struct {
	gorm.Model
	Username string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"username" binding:"required" example:"johndoe"`
	Email    string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"email" binding:"required,email" example:"john.doe@example.com"`
	Password string    `gorm:"type:varchar(255);not null" json:"-" binding:"required"` // Store hashed password, '-' to omit from JSON
	IsActive bool      `gorm:"default:true;not null" json:"is_active" example:"true"`
	RoleID   uint      `json:"role_id" example:"1"`                                                          // example:"1" ; removed binding:"required" to allow default role assignment
	Role     role.Role `gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"role"` // Belongs To relationship with Role

	LastLogin *time.Time `json:"last_login,omitempty"`
	// RefreshToken string `gorm:"type:varchar(512);index" json:"-"` // If refresh tokens are implemented, consider length and indexing
}

// LoginRequest defines the structure for user login requests.
type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"johndoe"` // Can be username or email
	Password string `json:"password" binding:"required" example:"password123"`
}

// RegisterRequest defines the structure for new user registration requests.
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=100" example:"janedoe"`
	Email    string `json:"email" binding:"required,email" example:"jane.doe@example.com"`
	Password string `json:"password" binding:"required,min=6,max=72" example:"SecurePassword123"` // Max 72 for bcrypt compatibility
	RoleID   uint   `json:"role_id,omitempty" example:"2"`                                        // Optional: if not provided, service might assign a default role
}

// Claims defines the JWT claims structure
type Claims struct {
	jwt.RegisteredClaims
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"` // Role name (e.g., "admin", "staff")
}

// AuthResponse defines the structure for authentication responses (e.g., login success)
type AuthResponse struct {
	User         UserCompact `json:"user"`
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token,omitempty"` // Only if refresh tokens are implemented
}

// UserCompact defines a compact user structure for API responses
type UserCompact struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	RoleName string `json:"role_name"`
	IsActive bool   `json:"is_active"`
}

// TokenDetails was present in your initial files but not used.
// If you plan to use it for more complex token management (e.g. with Redis), keep it.
// Otherwise, it can be removed if only simple access/refresh tokens are in AuthResponse.
/*
type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUUID   string
	RefreshUUID  string
	AtExpires    int64
	RtExpires    int64
}
*/
