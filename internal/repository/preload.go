package repository

import "gorm.io/gorm"

func preloadUserBrief(db *gorm.DB) *gorm.DB {
	return db.Select("id", "username", "image")
}

func activePostUserJoin(db *gorm.DB) *gorm.DB {
	return db.Joins("JOIN users ON users.id = posts.created_by AND users.deleted_at IS NULL")
}
