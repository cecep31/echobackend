package database

import (
	"echobackend/internal/model"

	"github.com/uptrace/bun"
)

// PostToTag represents the join table for the many-to-many relationship

// RegisterModels registers all models with Bun ORM
func RegisterModels(db *bun.DB) {
	// Register the join table first
	db.RegisterModel((*model.PostsToTags)(nil))
}
