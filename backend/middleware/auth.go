// prometheus/backend/middleware/auth.go
package middleware

import (
	"errors" // Make sure 'errors' is imported
	// Make sure 'fmt' is imported for potential future use, though not strictly needed for this fix
	"net/http"
	"prometheus/backend/internal/auth" // For auth.Claims
	"prometheus/backend/internal/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware creates a Gin middleware for JWT authentication.
// It verifies the token and sets user information in the context if valid.
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.SendErrorResponse(c, http.StatusUnauthorized, "Authorization header is required")
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			utils.SendErrorResponse(c, http.StatusUnauthorized, "Authorization header format must be Bearer {token}")
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims := &auth.Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// Validate the alg is what you expect (e.g., HMAC).
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				// If the signing method is not HMAC, then our secret is not applicable.
				// Return jwt.ErrSignatureInvalid to indicate this problem.
				// The parser will then wrap this in a jwt.ValidationError.
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtSecret), nil
		})

		if err != nil {
			var errMsg string
			// Use errors.Is to correctly check for wrapped error types provided by the jwt/v5 library.
			if errors.Is(err, jwt.ErrTokenMalformed) {
				errMsg = "Token is malformed."
			} else if errors.Is(err, jwt.ErrTokenExpired) {
				errMsg = "Token has expired."
			} else if errors.Is(err, jwt.ErrTokenNotValidYet) {
				errMsg = "Token not yet valid."
			} else if errors.Is(err, jwt.ErrSignatureInvalid) {
				// This will catch the jwt.ErrSignatureInvalid returned from the Keyfunc
				// if the signing method was unexpected, or if the signature itself is invalid.
				errMsg = "Token signature is invalid or signing method is not supported."
			} else {
				// For any other errors, including other jwt.ValidationError types not explicitly checked above.
				errMsg = "Invalid token: " + err.Error()
			}
			utils.SendErrorResponse(c, http.StatusUnauthorized, errMsg)
			c.Abort()
			return
		}

		// This check is technically redundant if the error handling above is comprehensive,
		// as an invalid token would have resulted in an error from ParseWithClaims.
		// However, it's a good safeguard.
		if !token.Valid {
			utils.SendErrorResponse(c, http.StatusUnauthorized, "Token is invalid.")
			c.Abort()
			return
		}

		// Token is valid, set user claims in context for downstream handlers
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)

		c.Next()
	}
}
