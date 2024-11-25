package model

type Tag struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
}
