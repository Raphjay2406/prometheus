// prometheus/backend/routes/router.go
package routes

import (
	"net/http"
	"prometheus/backend/config"
	"prometheus/backend/internal/auth"
	"prometheus/backend/internal/utils" // For the placeholder handler & responses
	"prometheus/backend/middleware"     // Ensure your middleware package is correctly referenced

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRoutes initializes all API routes including authentication and protected routes.
func SetupRoutes(r *gin.Engine, db *gorm.DB, cfg *config.Config) {
	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Prometheus backend is healthy and running!"})
	})

	// Initialize services and handlers
	// Auth
	authService := auth.NewAuthService(db, cfg)
	authHandler := auth.NewAuthHandler(authService)

	// API v1 Group
	apiV1 := r.Group("/api/v1")
	{
		// --- Authentication Routes (Public) ---
		authRoutes := apiV1.Group("/auth")
		{
			authRoutes.POST("/register", authHandler.Register)
			authRoutes.POST("/login", authHandler.Login)
			// TODO: Add future auth routes: /refresh-token, /logout, /forgot-password, /reset-password
		}

		// --- Protected Routes (Require Authentication via JWT) ---
		protected := apiV1.Group("/")
		protected.Use(middleware.AuthMiddleware(cfg.JWTSecret)) // Apply JWT authentication
		{
			// Example: Get current authenticated user's profile
			protected.GET("/me", func(c *gin.Context) {
				userID, _ := c.Get("userID")
				username, _ := c.Get("username")
				email, _ := c.Get("email")
				role, _ := c.Get("role")

				utils.SendSuccessResponse(c, http.StatusOK, "Current user profile fetched successfully", gin.H{
					"id":       userID,
					"username": username,
					"email":    email,
					"role":     role,
				})
			})

			// --- Admin Only Routes (Example of RBAC) ---
			// These routes require authentication AND 'admin' or 'god-admin' role.
			adminRoutes := protected.Group("/admin")
			// Apply RBACMiddleware for admin roles AFTER AuthMiddleware
			adminRoutes.Use(middleware.RBACMiddleware("admin", "god-admin"))
			{
				adminRoutes.GET("/dashboard", func(c *gin.Context) {
					username, _ := c.Get("username") // Username is set by AuthMiddleware
					utils.SendSuccessResponse(c, http.StatusOK, "Admin dashboard data loaded.", gin.H{
						"message": "Welcome to the admin dashboard, " + username.(string) + "!",
					})
				})
				// TODO: Add more admin-specific routes: user management, system settings, audit logs etc.
				// adminRoutes.GET("/users", userHandler.ListUsers)
				// adminRoutes.PUT("/users/:userID/status", userHandler.UpdateUserStatus)
			}

			// --- HR Routes (Example of RBAC) ---
			hrRoutes := protected.Group("/hr")
			// HR, Admin, and GodAdmin can access these routes
			hrRoutes.Use(middleware.RBACMiddleware("hr", "admin", "god-admin"))
			{
				hrRoutes.GET("/employee-data", func(c *gin.Context) {
					utils.SendSuccessResponse(c, http.StatusOK, "Sensitive Employee Data (Mock)", gin.H{
						"data": "This is mock HR-specific employee data accessible by HR, Admin, GodAdmin.",
					})
				})
				// TODO: Add more HR-specific routes: manage employee profiles, leave requests, payroll previews etc.
			}

			// --- Manager Routes (Example of RBAC) ---
			managerRoutes := protected.Group("/manager")
			// Managers, HR, Admin, and GodAdmin can access these routes
			managerRoutes.Use(middleware.RBACMiddleware("manager", "hr", "admin", "god-admin"))
			{
				managerRoutes.GET("/team-overview", func(c *gin.Context) {
					utils.SendSuccessResponse(c, http.StatusOK, "Team Overview Data (Mock)", gin.H{
						"data": "This is mock data for a manager's team.",
					})
				})
				// TODO: Add routes for approving leave, overtime for team members.
			}

			// --- Staff Routes (Example of RBAC) ---
			// Example for a 'staff' accessible route (most permissive after login)
			// All authenticated users (staff, manager, hr, admin, god-admin) can access these.
			staffAccessibleRoutes := protected.Group("/staff-area") // Using a more descriptive group name
			staffAccessibleRoutes.Use(middleware.RBACMiddleware("staff", "manager", "hr", "admin", "god-admin"))
			{
				staffAccessibleRoutes.GET("/my-tasks", func(c *gin.Context) {
					utils.SendSuccessResponse(c, http.StatusOK, "List of my tasks (Mock)", gin.H{
						"tasks": []string{"Complete TPS reports", "Attend mandatory fun session"},
					})
				})
			}

			// TODO: Add other protected routes for different modules (user, division, attendance, etc.)
			// Ensure each group has appropriate RBACMiddleware.
		}
	}

	// Fallback for undefined routes (404 Not Found)
	r.NoRoute(func(c *gin.Context) {
		utils.SendErrorResponse(c, http.StatusNotFound, "The requested resource was not found on this server.")
	})
}
