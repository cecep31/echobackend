package model

import (
	"time"

	"gorm.io/gorm"
)

// User represents the user model in the database
type User struct {
	// gorm.Model can be embedded for ID, CreatedAt, UpdatedAt, DeletedAt
	// However, since ID is a string (likely UUID) and other fields are already defined,
	// we will define them explicitly with GORM tags.
	ID             string         `json:"id" gorm:"type:uuid;primaryKey"` // Assuming UUID, adjust if different
	Email          string         `json:"email" gorm:"uniqueIndex;not null"`
	FirstName      string         `json:"first_name" gorm:"not null"`
	LastName       string         `json:"last_name" gorm:"not null"`
	Username       string         `json:"username" gorm:"uniqueIndex;not null"`
	Image          string         `json:"image"`
	Password       string         `json:"-" gorm:"not null"`
	IsSuperAdmin   bool           `json:"is_super_admin"`
	FollowersCount int64          `json:"followers_count" gorm:"default:0"`
	FollowingCount int64          `json:"following_count" gorm:"default:0"`
	CreatedAt      time.Time      `json:"created_at"`     // GORM handles this automatically
	UpdatedAt      time.Time      `json:"updated_at"`     // GORM handles this automatically
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"` // GORM specific type for soft delete

	// Relationships
	Followers []UserFollow `json:"followers,omitempty" gorm:"foreignKey:FollowingID"`
	Following []UserFollow `json:"following,omitempty" gorm:"foreignKey:FollowerID"`
}

// TableName specifies the table name for GORM
func (User) TableName() string {
	return "users"
}

// UserResponse represents the user data that can be safely sent to clients
type UserResponse struct {
	ID             string     `json:"id"`
	Email          string     `json:"email"`
	Name           string     `json:"name"`
	Username       string     `json:"username"`
	Image          string     `json:"image"`
	IsSuperAdmin   bool       `json:"is_super_admin"`
	FirstName      string     `json:"first_name"`
	LastName       string     `json:"last_name"`
	FollowersCount int64      `json:"followers_count"`
	FollowingCount int64      `json:"following_count"`
	IsFollowing    *bool      `json:"is_following,omitempty"` // Whether current user follows this user
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"` // Keep as *time.Time for response flexibility
}

// ToResponse converts a User model to a UserResponse
func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:             u.ID,
		Email:          u.Email,
		Name:           u.FirstName + " " + u.LastName,
		Username:       u.Username,
		Image:          u.Image,
		IsSuperAdmin:   u.IsSuperAdmin,
		FirstName:      u.FirstName,
		LastName:       u.LastName,
		FollowersCount: u.FollowersCount,
		FollowingCount: u.FollowingCount,
		CreatedAt:      u.CreatedAt,
		UpdatedAt:      u.UpdatedAt,
		// Convert gorm.DeletedAt to *time.Time for the response
		DeletedAt: convertDeletedAtToTime(u.DeletedAt),
	}
}

// convertDeletedAtToTime helper function
func convertDeletedAtToTime(deletedAt gorm.DeletedAt) *time.Time {
	if deletedAt.Valid {
		return &deletedAt.Time
	}
	return nil
}
