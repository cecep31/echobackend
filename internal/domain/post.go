package domain

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	ID        string         `json:"id" gorm:"primaryKey"`
	Title     string         `json:"title"`
	Photo_url string         `json:"photo_url" gorm:"type:text"`
	Slug      string         `json:"slug" gorm:"unique"`
	Body      string         `json:"body" gorm:"type=text"`
	Published bool           `json:"published"`
	CreatedBy string         `json:"created_by"`
	Creator   User           `json:"creator" gorm:"foreignKey:CreatedBy"`
	Tags      []Tag          `json:"tags" gorm:"many2many:posts_to_tags"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `sql:"index" json:"deleted_at"`
}
