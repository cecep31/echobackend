package model

import (
	"github.com/uptrace/bun"
)

type Tag struct {
	bun.BaseModel `bun:"table:tags,alias:t"`

	ID   uint   `json:"id" bun:"id,pk"`
	Name string `json:"name" bun:"name"`
}
