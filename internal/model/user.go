package model

import (
	"time"

	"github.com/uptrace/bun"
)

// User represents the user model in the database
type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID           string     `json:"id" bun:"id,pk"`
	Email        string     `json:"email" bun:"email,unique,notnull"`
	FirstName    string     `json:"first_name" bun:"first_name,notnull"`
	LastName     string     `json:"last_name" bun:"last_name,notnull"`
	Username     string     `json:"username" bun:"username,unique,notnull"`
	Password     string     `json:"-" bun:"password,notnull"` // "-" means this field won't be included in JSON
	IsSuperAdmin bool       `json:"is_super_admin" bun:"is_super_admin"`
	CreatedAt    time.Time  `json:"created_at" bun:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" bun:"updated_at"`
	DeletedAt    *time.Time `json:"-" bun:"deleted_at,soft_delete"`
}

// No need for TableName with Bun as it's specified in the struct tag

// UserResponse represents the user data that can be safely sent to clients
type UserResponse struct {
	ID           string     `json:"id"`
	Email        string     `json:"email"`
	Name         string     `json:"name"`
	Username     string     `json:"username"`
	IsSuperAdmin bool       `json:"is_super_admin"`
	FirstName    string     `json:"first_name"`
	LastName     string     `json:"last_name"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
}

// ToResponse converts a User model to a UserResponse
func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:           u.ID,
		Email:        u.Email,
		Name:         u.FirstName + " " + u.LastName,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		IsSuperAdmin: u.IsSuperAdmin,
		Username:     u.Username,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
		DeletedAt:    u.DeletedAt,
	}
}
