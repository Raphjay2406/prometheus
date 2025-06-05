// prometheus/backend/middleware/rbac.go
package middleware

import (
	"net/http"
	"prometheus/backend/internal/utils"
	"slices" // Go 1.21+ for slices.Contains

	"github.com/gin-gonic/gin"
)

// RBACMiddleware creates a Gin middleware for Role-Based Access Control.
// It checks if the authenticated user's role is one of the allowedRoles.
// This middleware should be used AFTER AuthMiddleware.
func RBACMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Attempt to get user role from context (set by AuthMiddleware)
		userRoleInterface, exists := c.Get("role")
		if !exists {
			// This should ideally not happen if AuthMiddleware is applied first
			// and successfully authenticates the user.
			utils.SendErrorResponse(c, http.StatusForbidden, "Access Denied: User role not found in context. Ensure AuthMiddleware runs first.")
			c.Abort()
			return
		}

		userRole, ok := userRoleInterface.(string)
		if !ok {
			// Role in context is not a string, which is unexpected.
			utils.SendErrorResponse(c, http.StatusInternalServerError, "Server Error: User role in context is not of expected type.")
			c.Abort()
			return
		}

		if len(userRole) == 0 {
			utils.SendErrorResponse(c, http.StatusForbidden, "Access Denied: User role is empty.")
			c.Abort()
			return
		}

		// Check if the user's role is in the list of allowed roles
		// Using slices.Contains for Go 1.21+
		if !slices.Contains(allowedRoles, userRole) {
			utils.SendErrorResponse(c, http.StatusForbidden, "Access Denied: You do not have the required role for this resource.")
			c.Abort()
			return
		}

		// User has the required role, proceed to the next handler
		c.Next()
	}
}
