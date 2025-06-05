// prometheus/backend/internal/auth/service.go
package auth

import (
	"errors"
	"fmt"
	"prometheus/backend/config"
	"prometheus/backend/internal/role" // Ensure this path is correct for your role package
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthService defines the interface for authentication operations.
type AuthService interface {
	RegisterUser(req RegisterRequest) (*User, error)
	LoginUser(req LoginRequest) (*AuthResponse, error)
	GenerateJWT(user *User) (string, error)
	ValidatePassword(hashedPassword, plainPassword string) error
}

// authService implements the AuthService interface.
type authService struct {
	db  *gorm.DB
	cfg *config.Config
}

// NewAuthService creates a new instance of AuthService.
func NewAuthService(db *gorm.DB, cfg *config.Config) AuthService {
	return &authService{db: db, cfg: cfg}
}

// HashPassword hashes a given password using bcrypt.
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// ValidatePassword compares a hashed password with a plain password.
func (s *authService) ValidatePassword(hashedPassword, plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}

// RegisterUser handles new user registration.
func (s *authService) RegisterUser(req RegisterRequest) (*User, error) {
	// Check if username or email already exists
	var existingUser User
	// The error "relation 'users' does not exist" originated from this GORM query
	// because the table wasn't created yet. AutoMigrate in main.go fixes this.
	if err := s.db.Where("username = ? OR email = ?", req.Username, req.Email).First(&existingUser).Error; err == nil {
		return nil, errors.New("username or email already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// This means a real database error occurred, other than "not found"
		return nil, fmt.Errorf("database error while checking existing user: %w", err)
	}

	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Determine RoleID
	roleID := req.RoleID
	var userRole role.Role // To hold the role details

	if roleID == 0 {
		// Default to "staff" role if RoleID is not provided or is 0
		if err := s.db.Where("name = ?", "staff").First(&userRole).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// This error highlights the need for seeding roles after migration.
				return nil, errors.New("default 'staff' role not found. Please ensure roles are seeded")
			}
			return nil, fmt.Errorf("failed to fetch default 'staff' role: %w", err)
		}
		roleID = userRole.ID
	} else {
		// Validate if the provided RoleID exists
		if err := s.db.First(&userRole, roleID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, fmt.Errorf("role with ID %d not found", roleID)
			}
			return nil, fmt.Errorf("failed to verify role ID %d: %w", roleID, err)
		}
	}

	newUser := User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		RoleID:   roleID,
		IsActive: true, // Default to active, can be changed by admin later
	}

	if err := s.db.Create(&newUser).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// After creating the user, their ID is populated. Now, preload their Role.
	// It's good practice to return the newly created user with its associated role.
	// The 'newUser' variable here will have its Role field populated by this Preload.
	if err := s.db.Preload("Role").First(&newUser, newUser.ID).Error; err != nil {
		// Log error but proceed; role might not be critical for immediate response, but it's good to know.
		fmt.Printf("Warning: failed to preload role for new user %s (ID: %d): %v\n", newUser.Username, newUser.ID, err)
		// Even if preloading fails, the user was created.
		// You might decide to return an error here if Role is absolutely critical for the response.
	}

	return &newUser, nil
}

// LoginUser handles user login and JWT generation.
func (s *authService) LoginUser(req LoginRequest) (*AuthResponse, error) {
	var user User
	// Preload Role to get Role.Name for JWT claims and user response
	// Login can be by username or email.
	if err := s.db.Preload("Role").Where("username = ? OR email = ?", req.Username, req.Username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid username or password") // Keep error generic for security
		}
		return nil, fmt.Errorf("database error during login: %w", err)
	}

	if !user.IsActive {
		return nil, errors.New("user account is inactive")
	}

	if err := s.ValidatePassword(user.Password, req.Password); err != nil {
		return nil, errors.New("invalid username or password") // Keep error generic
	}

	// Update LastLogin
	now := time.Now().UTC() // Use UTC for consistency
	user.LastLogin = &now
	if err := s.db.Save(&user).Error; err != nil {
		// Log error but proceed with login as this is not critical enough to fail login
		fmt.Printf("Warning: failed to update last login for user %s: %v\n", user.Username, err)
	}

	accessToken, err := s.GenerateJWT(&user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	authResponse := &AuthResponse{
		User: UserCompact{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			RoleName: user.Role.Name, // Role.Name should be populated due to Preload
			IsActive: user.IsActive,
		},
		AccessToken: accessToken,
		// RefreshToken: // TODO: Implement refresh token generation if needed
	}

	return authResponse, nil
}

// GenerateJWT creates a new JWT for a given user.
func (s *authService) GenerateJWT(user *User) (string, error) {
	// Ensure Role.Name is available for the JWT claims.
	// It should typically be preloaded before calling GenerateJWT.
	// If not, attempt a last-minute load.
	if user.Role.Name == "" && user.RoleID != 0 {
		var roleFromDB role.Role
		if err := s.db.First(&roleFromDB, user.RoleID).Error; err != nil {
			return "", fmt.Errorf("could not retrieve role name (ID: %d) for JWT generation: %w", user.RoleID, err)
		}
		user.Role.Name = roleFromDB.Name // Populate the role name
	} else if user.Role.Name == "" && user.RoleID == 0 {
		return "", errors.New("user has no RoleID or Role.Name for JWT generation")
	}

	expirationTime := time.Now().Add(time.Duration(s.cfg.JWTExpirationHours) * time.Hour)
	if s.cfg.JWTExpirationHours == 0 { // Default if not set or zero
		expirationTime = time.Now().Add(24 * 7 * time.Hour) // Default to 7 days
	}

	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			NotBefore: jwt.NewNumericDate(time.Now().UTC()),
			Subject:   fmt.Sprintf("%d", user.ID),
		},
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role.Name, // Role name (e.g., "admin", "staff")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT token: %w", err)
	}

	return signedToken, nil
}
