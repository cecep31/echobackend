package test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestPostQueryFilterValidation tests the filter validation methods
func TestPostQueryFilterValidation(t *testing.T) {
	filter := &PostQueryFilter{}

	// Test default sort field
	assert.Equal(t, "posts.created_at", filter.GetSortField(), "Default sort field should be created_at")

	// Test valid sort field
	filter.SortBy = "title"
	assert.Equal(t, "posts.title", filter.GetSortField(), "Sort field should be posts.title")

	// Test invalid sort field (should default to created_at)
	filter.SortBy = "invalid_field"
	assert.Equal(t, "posts.created_at", filter.GetSortField(), "Invalid sort field should default to created_at")

	// Test default sort order
	assert.Equal(t, "desc", filter.GetSortOrder(), "Default sort order should be desc")

	// Test valid sort order
	filter.SortOrder = "asc"
	assert.Equal(t, "asc", filter.GetSortOrder(), "Sort order should be asc")

	// Test invalid sort order (should default to desc)
	filter.SortOrder = "invalid"
	assert.Equal(t, "desc", filter.GetSortOrder(), "Invalid sort order should default to desc")
}

// TestQueryParameterParsing tests that query parameters are parsed correctly
func TestQueryParameterParsing(t *testing.T) {
	// Test tags parsing
	tags := "golang, webdev, api"
	parsedTags := parseTagsString(tags)

	expectedTags := []string{"golang", "webdev", "api"}
	assert.Equal(t, expectedTags, parsedTags, "Tags should be parsed correctly")
}

// parseTagsString helper function to parse comma-separated tags
func parseTagsString(tagsString string) []string {
	if tagsString == "" {
		return []string{}
	}

	tags := make([]string, 0)
	for _, tag := range strings.Split(tagsString, ",") {
		trimmed := strings.TrimSpace(tag)
		if trimmed != "" {
			tags = append(tags, trimmed)
		}
	}
	return tags
}

// PostQueryFilter represents the filter structure for testing
type PostQueryFilter struct {
	Limit     int    `json:"limit" query:"limit"`
	Offset    int    `json:"offset" query:"offset"`
	Search    string `json:"search" query:"search"`
	SortBy    string `json:"sort_by" query:"sort_by"`
	SortOrder string `json:"sort_order" query:"sort_order"`
	StartDate string `json:"start_date" query:"start_date"`
	EndDate   string `json:"end_date" query:"end_date"`
	CreatedBy string `json:"created_by" query:"created_by"`
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
