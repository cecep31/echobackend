package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// Comment represents a comment on a post
type Comment struct {
	bun.BaseModel `bun:"table:comments,alias:c"`

	ID        uuid.UUID  `json:"id" bun:"id,pk,type:uuid,default:uuid_generate_v4()"`
	BlockID   uuid.UUID  `json:"block_id" bun:"block_id,notnull,type:uuid"`
	Content   string     `json:"content" bun:"content,notnull,type:text"`
	CreatedBy string     `json:"created_by" bun:"created_by,notnull"`
	CreatedAt time.Time  `json:"created_at" bun:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" bun:"updated_at"`
	DeletedAt *time.Time `bun:"deleted_at,soft_delete" json:"deleted_at,omitempty"`
}

// No need for TableName with Bun as it's specified in the struct tag
