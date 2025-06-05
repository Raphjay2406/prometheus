// prometheus/backend/internal/auth/handler.go
package auth

import (
	"errors"
	"net/http"
	"prometheus/backend/internal/utils" // For error responses
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AuthHandler handles HTTP requests for authentication.
type AuthHandler struct {
	service AuthService
}

// NewAuthHandler creates a new instance of AuthHandler.
func NewAuthHandler(service AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

// Register handles new user registration requests.
// @Summary Register a new user
// @Description Creates a new user account. Default role is 'staff' if not specified.
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body RegisterRequest true "User registration details"
// @Success 201 {object} UserResponse "User created successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid input or user already exists"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}

	// Basic validation example (can be expanded with a validation library)
	if req.Username == "" || req.Email == "" || req.Password == "" {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Username, email, and password are required")
		return
	}
	if len(req.Password) < 6 {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Password must be at least 6 characters long")
		return
	}

	user, err := h.service.RegisterUser(req)
	if err != nil {
		// Check for specific error types if needed, e.g., user already exists
		if err.Error() == "username or email already exists" {
			utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		if err.Error() == "default 'staff' role not found. Please ensure roles are seeded" {
			// This error implies roles should be seeded. The AutoMigrate will create the table,
			// but seeding data (like specific roles) is a separate step, often done after migration.
			utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		// A bit fragile check for role ID not found error from service layer
		if _, ok := err.(interface{ Error() string }); ok && len(err.Error()) > 18 && err.Error()[:18] == "role with ID" {
			utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to register user: "+err.Error())
		return
	}

	// Create a response struct that doesn't include the password
	userResponse := UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		IsActive:  user.IsActive,
		RoleID:    user.RoleID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	if user.Role.Name != "" { // if role was preloaded
		userResponse.RoleName = user.Role.Name
	}

	utils.SendSuccessResponse(c, http.StatusCreated, "User registered successfully", userResponse)
}

// Login handles user login requests.
// @Summary Log in a user
// @Description Authenticates a user and returns a JWT.
// @Tags Auth
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "User login credentials"
// @Success 200 {object} AuthResponse "Login successful, includes user details and access token"
// @Failure 400 {object} utils.ErrorResponse "Invalid input"
// @Failure 401 {object} utils.ErrorResponse "Invalid username or password, or inactive account"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}

	authResponse, err := h.service.LoginUser(req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || err.Error() == "invalid username or password" {
			utils.SendErrorResponse(c, http.StatusUnauthorized, "Invalid username or password")
			return
		}
		if err.Error() == "user account is inactive" {
			utils.SendErrorResponse(c, http.StatusUnauthorized, err.Error())
			return
		}
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Login failed: "+err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Login successful", authResponse)
}

// UserResponse is a subset of User for registration responses.
// Avoids exposing hashed password or too many internal details directly.
type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	IsActive  bool      `json:"is_active"`
	RoleID    uint      `json:"role_id"`
	RoleName  string    `json:"role_name,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
