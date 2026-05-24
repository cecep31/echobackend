package repository

import "gorm.io/gorm"

func preloadUserBrief(db *gorm.DB) *gorm.DB {
	return db.Select("id", "username", "image")
}
