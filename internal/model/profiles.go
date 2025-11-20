package model

import (
	"time"
)

type Profile struct {
	ID        int        `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    string     `json:"user_id" gorm:"uniqueIndex;not null;type:uuid"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	Bio       *string    `json:"bio"`
	Website   *string    `json:"website"`
	Phone     *string    `json:"phone" gorm:"type:varchar(50)"`
	Location  *string    `json:"location" gorm:"type:varchar(255)"`
	User      User       `gorm:"foreignKey:UserID"`
}

func (Profile) TableName() string {
	return "profiles"
}
