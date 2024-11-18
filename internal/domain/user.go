package domain

import "time"

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name" validate:"required,min=3"`
	Email     string    `json:"email" validate:"required,email"`
	Username  string    `json:"username" validate:"required,min=3"`
	Password  string    `json:"-,omitempty" validate:"required,min=6"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserResponse struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
