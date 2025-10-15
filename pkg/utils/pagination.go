package utils

import (
	"fmt"
	"math"
)

// Pagination represents pagination parameters and metadata
type Pagination struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
	Offset     int   `json:"offset"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

// NewPagination creates a new pagination object
func NewPagination(page, limit int, totalItems int64) *Pagination {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	} else if limit > 100 {
		limit = 100 // cap at 100 to prevent excessive resource usage
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(limit)))
	offset := (page - 1) * limit

	return &Pagination{
		Page:       page,
		Limit:      limit,
		TotalItems: totalItems,
		TotalPages: totalPages,
		Offset:     offset,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}
}

// ValidatePagination validates page and limit parameters
func ValidatePagination(page, limit int) error {
	if page < 1 {
		return fmt.Errorf("page must be greater than 0")
	}
	if limit < 1 {
		return fmt.Errorf("limit must be greater than 0")
	}
	if limit > 100 {
		return fmt.Errorf("limit must not exceed 100")
	}
	return nil
}

// GetOffset calculates the offset for pagination queries
func GetOffset(page, limit int) int {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	return (page - 1) * limit
}
