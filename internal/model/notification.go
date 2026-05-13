package model

import "time"

type Notification struct {
	ID        string     `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID    string     `json:"user_id" gorm:"type:uuid;not null;index:idx_notifications_user_id"`
	Type      string     `json:"type" gorm:"type:varchar(50);not null"`
	Title     string     `json:"title" gorm:"type:varchar(255);not null"`
	Message   *string    `json:"message"`
	Read      bool       `json:"read" gorm:"default:false;not null;index:idx_notifications_read"`
	Data      *string    `json:"data"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`

	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

func (Notification) TableName() string {
	return "notifications"
}
