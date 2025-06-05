// prometheus/backend/internal/role/model.go
package role

import "gorm.io/gorm"

// Role represents a user role in the system.
type Role struct {
	gorm.Model
	Name        string `gorm:"type:varchar(50);uniqueIndex;not null" json:"name" example:"admin"`
	Description string `gorm:"type:varchar(255)" json:"description" example:"Administrator with full access"`

	// Users []auth.User `gorm:"foreignKey:RoleID"` // Example of a Has Many relationship if needed later
}
