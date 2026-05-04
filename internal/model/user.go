package model

import (
	"time"

	"gorm.io/gorm"
)

// User represents the user model in the database
type User struct {
	ID             string         `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	CreatedAt      *time.Time     `json:"created_at"`
	UpdatedAt      *time.Time     `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
	FirstName      *string        `json:"first_name" gorm:"type:varchar(255)"`
	LastName       *string        `json:"last_name" gorm:"type:varchar(255)"`
	Email          string         `json:"email" gorm:"uniqueIndex;not null;type:varchar(255)"`
	Password       *string        `json:"-" gorm:"type:varchar(255)"`
	Image          *string        `json:"image"`
	IsSuperAdmin   *bool          `json:"-" gorm:"default:false"`
	Username       *string        `json:"username" gorm:"uniqueIndex;type:varchar(255)"`
	GithubID       *int64         `json:"github_id" gorm:"uniqueIndex:users_github_id_unique"`
	FollowersCount int64          `json:"followers_count" gorm:"type:bigint;default:0"`
	FollowingCount int64          `json:"following_count" gorm:"type:bigint;default:0"`
	LastLoggedAt   *time.Time     `json:"last_logged_at"`

	// Relationships
	Files           []File           `gorm:"foreignKey:CreatedBy"`
	PostComments    []PostComment    `gorm:"foreignKey:CreatedBy"`
	PostLikes       []PostLike       `gorm:"foreignKey:UserID"`
	PostViews       []PostView       `gorm:"foreignKey:UserID"`
	PostBookmarks   []PostBookmark   `gorm:"foreignKey:UserID"`
	BookmarkFolders []BookmarkFolder `gorm:"foreignKey:UserID"`
	Posts           []Post           `gorm:"foreignKey:CreatedBy"`
	Profile         *Profile         `gorm:"foreignKey:UserID"`
	Sessions        []Session        `gorm:"foreignKey:UserID"`
	Followers       []UserFollow     `gorm:"foreignKey:FollowingID"`
	Following       []UserFollow     `gorm:"foreignKey:FollowerID"`
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
	Username       *string    `json:"username"`
	Image          *string    `json:"image"`
	FirstName      *string    `json:"first_name"`
	LastName       *string    `json:"last_name"`
	FollowersCount int64      `json:"followers_count"`
	FollowingCount int64      `json:"following_count"`
	IsFollowing    *bool      `json:"is_following,omitempty"` // Whether current user follows this user
	Profile        *Profile   `json:"profile,omitempty"`
	CreatedAt      *time.Time `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
}

// PublicUserResponse represents user data for public endpoints
type PublicUserResponse struct {
	ID             string     `json:"id"`
	Email          string     `json:"email"`
	Name           string     `json:"name"`
	Username       *string    `json:"username"`
	Image          *string    `json:"image"`
	FirstName      *string    `json:"first_name"`
	LastName       *string    `json:"last_name"`
	FollowersCount int64      `json:"followers_count"`
	FollowingCount int64      `json:"following_count"`
	IsFollowing    *bool      `json:"is_following,omitempty"`
	Profile        *Profile   `json:"profile,omitempty"`
	CreatedAt      *time.Time `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
}

// ToPublicResponse converts a User model to a PublicUserResponse
func (u *User) ToPublicResponse() *PublicUserResponse {
	name := ""
	if u.FirstName != nil && u.LastName != nil {
		name = *u.FirstName + " " + *u.LastName
	}
	return &PublicUserResponse{
		ID:             u.ID,
		Email:          u.Email,
		Name:           name,
		Username:       u.Username,
		Image:          u.Image,
		FirstName:      u.FirstName,
		LastName:       u.LastName,
		FollowersCount: u.FollowersCount,
		FollowingCount: u.FollowingCount,
		Profile:        u.Profile,
		CreatedAt:      u.CreatedAt,
		UpdatedAt:      u.UpdatedAt,
	}
}

// ToResponse converts a User model to a UserResponse
func (u *User) ToResponse() *UserResponse {
	name := ""
	if u.FirstName != nil && u.LastName != nil {
		name = *u.FirstName + " " + *u.LastName
	}
	return &UserResponse{
		ID:             u.ID,
		Email:          u.Email,
		Name:           name,
		Username:       u.Username,
		Image:          u.Image,
		FirstName:      u.FirstName,
		LastName:       u.LastName,
		FollowersCount: u.FollowersCount,
		FollowingCount: u.FollowingCount,
		Profile:        u.Profile,
		CreatedAt:      u.CreatedAt,
		UpdatedAt:      u.UpdatedAt,
	}
}
