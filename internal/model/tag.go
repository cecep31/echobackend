package model

import (
	"github.com/uptrace/bun"
)

type Tag struct {
	bun.BaseModel `bun:"table:tags,alias:t"`
	ID            uint   `json:"id" bun:"id,pk"`
	Name          string `json:"name" bun:"name"`
	// Posts         []Post `json:"-" bun:"m2m:posts_to_tags,join:Post=Tag"`
}

type PostsToTags struct {
	bun.BaseModel `bun:"table:posts_to_tags,alias:pt"`
	ID            string `bun:"id,pk,type:uuid" json:"id"`

	// Post reference
	PostID string `bun:"post_id,pk,type:uuid"`
	Post   *Post  `bun:"rel:belongs-to,join:post_id=id"`
	TagID  uint   `bun:"tag_id,pk"`
	Tag    *Tag   `bun:"rel:belongs-to,join:tag_id=id"`
}
