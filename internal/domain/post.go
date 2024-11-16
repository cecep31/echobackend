package domain

type Post struct {
	Id        string `json:"id" gorm:"primaryKey"`
	Title     string `json:"title"`
	Photo_url string `json:"photo_url" gorm:"type:text"`
	Slug      string `json:"slug" gorm:"unique"`
	Body      string `json:"body" gorm:"type=text"`
	CreatedBy string `json:"created_by"`
	Creator   User   `json:"creator" gorm:"foreignKey:CreatedBy"`
}
