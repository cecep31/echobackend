package model

import (
	"time"

	"gorm.io/gorm" // Import GORM
)

type Post struct {
	ID            string         `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	CreatedAt     *time.Time     `json:"created_at"`
	UpdatedAt     *time.Time     `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
	Title         *string        `json:"title" gorm:"type:varchar(255);not null"`
	CreatedBy     *string        `json:"created_by" gorm:"type:uuid;not null;uniqueIndex:creator_and_slug_unique"`
	Body          *string        `json:"body"`
	Slug          *string        `json:"slug" gorm:"type:varchar(255);not null;uniqueIndex:creator_and_slug_unique"`
	Photo_url     *string        `json:"photo_url"`
	Published     *bool          `json:"published" gorm:"default:true"`
	PublishedAt   *time.Time     `json:"published_at"`
	ViewCount     int64          `json:"view_count" gorm:"type:bigint;default:0;check:view_count >= 0"`
	LikeCount     int64          `json:"like_count" gorm:"type:bigint;default:0;check:like_count >= 0"`
	BookmarkCount int64          `json:"bookmark_count" gorm:"type:bigint;default:0;check:bookmark_count >= 0"`
	PostComments  []PostComment  `gorm:"foreignKey:PostID"`
	PostLikes     []PostLike     `gorm:"foreignKey:PostID"`
	PostBookmarks []PostBookmark `gorm:"foreignKey:PostID"`
	Creator       *User          `gorm:"foreignKey:CreatedBy;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Tags          []Tag          `gorm:"many2many:posts_to_tags;"`
}

func (Post) TableName() string {
	return "posts"
}

type PostResponse struct {
	ID            string        `json:"id"`
	Title         *string       `json:"title"`
	Photo_url     *string       `json:"photo_url"`
	Body          *string       `json:"body"`
	Slug          *string       `json:"slug"`
	ViewCount     int64         `json:"view_count"`
	LikeCount     int64         `json:"like_count"`
	BookmarkCount int64         `json:"bookmark_count"`
	Published     *bool         `json:"published"`
	PublishedAt   *time.Time    `json:"published_at"`
	Creator       *UserResponse `json:"creator,omitempty"`
	Tags          []TagResponse `json:"tags,omitempty"`
	CreatedAt     *time.Time    `json:"created_at"`
	UpdatedAt     *time.Time    `json:"updated_at"`
	DeletedAt     *time.Time    `json:"deleted_at,omitempty"`
}

func (p *Post) ToResponse() *PostResponse {
	var creatorResp *UserResponse
	if p.Creator.ID != "" {
		creatorResp = p.Creator.ToResponse()
	}

	var tagResponses []TagResponse
	if p.Tags != nil {
		tagResponses = make([]TagResponse, len(p.Tags))
		for i, tag := range p.Tags {
			tagResponses[i] = *tag.ToResponse()
		}
	}

	var deletedAtTime *time.Time
	if p.DeletedAt.Valid {
		deletedAtTime = &p.DeletedAt.Time
	}

	return &PostResponse{
		ID:            p.ID,
		Title:         p.Title,
		Photo_url:     p.Photo_url,
		Body:          p.Body,
		Slug:          p.Slug,
		ViewCount:     p.ViewCount,
		LikeCount:     p.LikeCount,
		BookmarkCount: p.BookmarkCount,
		Published:     p.Published,
		PublishedAt:   p.PublishedAt,
		Creator:       creatorResp,
		Tags:          tagResponses,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
		DeletedAt:     deletedAtTime,
	}
}

type CreatePostDTO struct {
	Title     string   `json:"title" validate:"required,min=7"`
	Photo_url string   `json:"photo_url"`
	Slug      string   `json:"slug" validate:"required,min=7"`
	Body      string   `json:"body" validate:"required,min=10"`
	Published bool     `json:"published"`
	Tags      []string `json:"tags" `
}

type UpdatePostDTO struct {
	Title     string   `json:"title"`
	Photo_url string   `json:"photo_url"`
	Slug      string   `json:"slug"`
	Body      string   `json:"body"`
	Published *bool    `json:"published"`
	Tags      []string `json:"tags"`
}

type PostQueryFilter struct {
	Limit     int      `json:"limit" query:"limit"`
	Offset    int      `json:"offset" query:"offset"`
	Search    string   `json:"search" query:"search"`
	SortBy    string   `json:"sort_by" query:"sort_by"`
	SortOrder string   `json:"sort_order" query:"sort_order"`
	StartDate string   `json:"start_date" query:"start_date"`
	EndDate   string   `json:"end_date" query:"end_date"`
	Published *bool    `json:"published" query:"published"`
	CreatedBy string   `json:"created_by" query:"created_by"`
	Tags      []string `json:"tags" query:"tags"`
}

// ValidSortFields defines allowed sort fields
func (f *PostQueryFilter) ValidSortFields() map[string]string {
	return map[string]string{
		"id":         "posts.id",
		"title":      "posts.title",
		"created_at": "posts.created_at",
		"updated_at": "posts.updated_at",
		"view_count": "posts.view_count",
		"like_count": "posts.like_count",
	}
}

// ValidSortOrders defines allowed sort orders
func (f *PostQueryFilter) ValidSortOrders() []string {
	return []string{"asc", "desc"}
}

// GetSortField returns the database field for sorting, defaults to created_at
func (f *PostQueryFilter) GetSortField() string {
	if field, exists := f.ValidSortFields()[f.SortBy]; exists {
		return field
	}
	return "posts.created_at" // Default sort field
}

// GetSortOrder returns the sort order, defaults to desc
func (f *PostQueryFilter) GetSortOrder() string {
	for _, order := range f.ValidSortOrders() {
		if order == f.SortOrder {
			return order
		}
	}
	return "desc" // Default sort order
}

// SitemapPost represents a minimal post structure for sitemap
type SitemapPost struct {
	Username  *string    `json:"username"`
	Slug      *string    `json:"slug"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}
