package utils

import (
	"github.com/gin-gonic/gin"
)

// SuccessResponse defines the structure for a successful API response.
type SuccessResponse struct {
	Status  string      `json:"status"`         // e.g., "success"
	Message string      `json:"message"`        // Descriptive message
	Data    interface{} `json:"data,omitempty"` // The actual data payload (optional)
}

// ErrorResponse defines the structure for an error API response.
type ErrorResponse struct {
	Status  string `json:"status"`  // e.g., "error"
	Message string `json:"message"` // Detailed error message
}

// SendSuccessResponse sends a standardized success JSON response.
func SendSuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, SuccessResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

// SendErrorResponse sends a standardized error JSON response.
func SendErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, ErrorResponse{
		Status:  "error",
		Message: message,
	})
}
