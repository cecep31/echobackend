package repository

import (
	"context"
	"errors" // For gorm.ErrRecordNotFound checks
	"fmt"    // For error wrapping

	"echobackend/internal/model"

	"gorm.io/gorm"
	// "gorm.io/gorm/clause" // For Preload with conditions if needed - will add if used
)

// Define common errors or use gorm.ErrRecordNotFound directly
var (
	ErrPostNotFound = errors.New("post not found")
)

type PostRepository interface {
	CreatePost(ctx context.Context, post *model.CreatePostDTO, creator_id string) (*model.Post, error)
	CreatePostWithTags(ctx context.Context, post *model.CreatePostDTO, creator_id string, tags []model.Tag) (*model.Post, error)
	GetPosts(ctx context.Context, limit int, offset int) ([]*model.Post, int64, error)
	GetPostByUsername(ctx context.Context, username string, offset int, limit int) ([]*model.Post, int64, error)
	GetPostsRandom(ctx context.Context, limit int) ([]*model.Post, error)
	GetPostByID(ctx context.Context, id string) (*model.Post, error)
	GetPostBySlugAndUsername(ctx context.Context, slug string, username string) (*model.Post, error)
	GetPostsByCreatedBy(ctx context.Context, createdBy string, offset int, limit int) ([]*model.Post, int64, error)
	DeletePostByID(ctx context.Context, id string) error
	UpdatePost(ctx context.Context, id string, post *model.UpdatePostDTO) (*model.Post, error)

	// Additional functions
	SearchPosts(ctx context.Context, keyword string, limit int, offset int) ([]*model.Post, int64, error)
	GetPostsByTag(ctx context.Context, tag string, limit int, offset int) ([]*model.Post, int64, error)
	ExistsByID(ctx context.Context, id string) (bool, error)
}

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) UpdatePost(ctx context.Context, id string, postDTO *model.UpdatePostDTO) (*model.Post, error) {
	// Check if the post exists first (optional, Updates will return RowsAffected = 0 if not found)
	// var existingCheck model.Post
	// if err := r.db.WithContext(ctx).Select("id").First(&existingCheck, "id = ?", id).Error; err != nil {
	// 	if errors.Is(err, gorm.ErrRecordNotFound) {
	// 		return nil, ErrPostNotFound
	// 	}
	// 	return nil, fmt.Errorf("error checking post existence before update: %w", err)
	// }

	// Create a map for updates to handle partial updates and zero values correctly if needed.
	// If UpdatePostDTO fields are pointers, checking for nil is good.
	// If they are value types, GORM's Updates() on a struct will only update non-zero fields by default.
	// Using a map gives more control.
	updates := make(map[string]interface{})

	// UpdatePostDTO fields are value types (string), so check for non-empty.
	if postDTO.Title != "" {
		updates["title"] = postDTO.Title
	}
	if postDTO.Body != "" {
		updates["body"] = postDTO.Body
	}
	if postDTO.Slug != "" {
		updates["slug"] = postDTO.Slug
	}
	if postDTO.Photo_url != "" {
		updates["photo_url"] = postDTO.Photo_url
	}
	if postDTO.Published != nil {
		updates["published"] = *postDTO.Published
	}

	// Handling Tags update is more complex and usually done via Associations.
	// For now, focusing on simple field updates.
	// if len(postDTO.Tags) > 0 { /* logic to update tags */ }

	if len(updates) == 0 && len(postDTO.Tags) == 0 { // Check if DTO is effectively empty for updates
		// No actual fields to update based on DTO content, and no tags to update.
		// Fetch and return the current post.
		// For now, let's assume an empty DTO means no operation or an error.
		// Or, if DTO might only contain tags to update, handle that separately.
		// We should fetch the post to return it, even if no fields changed.
		var currentPost model.Post
		err := r.db.WithContext(ctx).Preload("Creator").Preload("Tags").First(&currentPost, "id = ?", id).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrPostNotFound
			}
			return nil, fmt.Errorf("failed to fetch post (no updates provided): %w", err)
		}
		return &currentPost, nil
	}

	result := r.db.WithContext(ctx).Model(&model.Post{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to update post: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, ErrPostNotFound // Post with ID not found
	}

	// After updating, fetch the post again to get the full model with associations
	var updatedPost model.Post
	err := r.db.WithContext(ctx).Preload("Creator").Preload("Tags").First(&updatedPost, "id = ?", id).Error
	if err != nil {
		// This case (update succeeded but fetch failed) should be rare but handled.
		return nil, fmt.Errorf("post updated, but failed to retrieve updated record: %w", err)
	}

	// TODO: Handle Tags update if postDTO.Tags is provided.
	// This typically involves using GORM's association mode:
	// if postDTO.Tags != nil {
	//   // Convert DTO tags to []model.Tag or list of IDs
	//   // e.g., r.db.Model(&updatedPost).Association("Tags").Replace(newTags)
	// }

	return &updatedPost, nil
}

func (r *postRepository) CreatePost(ctx context.Context, postDTO *model.CreatePostDTO, creator_id string) (*model.Post, error) {
	newpost := &model.Post{
		Title:     &postDTO.Title,
		Slug:      &postDTO.Slug,
		Body:      &postDTO.Body,
		CreatedBy: &creator_id,
		Photo_url: &postDTO.Photo_url,
		Published: &postDTO.Published,
	}

	err := r.db.WithContext(ctx).Create(newpost).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}
	return newpost, nil // Returning the instance passed to Create. ID should be populated.
}

func (r *postRepository) CreatePostWithTags(ctx context.Context, postDTO *model.CreatePostDTO, creator_id string, tags []model.Tag) (*model.Post, error) {
	newpost := &model.Post{
		Title:     &postDTO.Title,
		Slug:      &postDTO.Slug,
		Body:      &postDTO.Body,
		CreatedBy: &creator_id,
		Photo_url: &postDTO.Photo_url,
		Published: &postDTO.Published,
		Tags:      tags, // Associate tags with the post
	}

	// Create the post with associated tags
	err := r.db.WithContext(ctx).Create(newpost).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create post with tags: %w", err)
	}

	// Load the created post with all associations for return
	err = r.db.WithContext(ctx).Preload("Creator").Preload("Tags").First(newpost, "id = ?", newpost.ID).Error
	if err != nil {
		return nil, fmt.Errorf("failed to load created post with associations: %w", err)
	}

	return newpost, nil
}

func (r *postRepository) GetPostByUsername(ctx context.Context, username string, offset int, limit int) ([]*model.Post, int64, error) {
	var posts []*model.Post
	var count int64

	// Base query for counting and fetching
	// Assuming Post model has CreatorID field linking to User model's ID,
	// and User model has Username field.
	// Table names "posts" and "users" are assumed based on GORM defaults or TableName() methods.
	query := r.db.WithContext(ctx).Model(&model.Post{}).
		Joins("JOIN users ON users.id = posts.created_by"). // Ensure 'posts.created_by' is the correct foreign key column name
		Where("users.username = ?", username)

	// Count total records
	err := query.Count(&count).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count posts for username %s: %w", username, err)
	}

	// Get paginated records
	// Re-apply Joins or ensure Preload works as expected.
	// For Preload("Creator"), GORM uses the defined foreign keys.
	// The explicit Join above is mainly for the WHERE clause on users.username.
	err = r.db.WithContext(ctx).Model(&model.Post{}).
		Preload("Creator"). // GORM will fetch Creator based on associations
		Preload("Tags").
		Joins("JOIN users ON users.id = posts.created_by"). // Keep join for filtering
		Where("users.username = ?", username).
		Order("posts.created_at DESC"). // Assuming 'posts' is the table name for Post model
		Offset(offset).
		Limit(limit).
		Find(&posts).Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to get posts for username %s: %w", username, err)
	}

	return posts, count, nil
}

func (r *postRepository) DeletePostByID(ctx context.Context, id string) error {
	// Assumes soft delete if model.Post has gorm.DeletedAt field.
	// If hard delete is needed, use Unscoped().Delete()
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Post{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete post: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrPostNotFound // Or specific error
	}
	return nil
}

func (r *postRepository) GetPosts(ctx context.Context, limit int, offset int) ([]*model.Post, int64, error) {
	var posts []*model.Post
	var count int64

	// Count total published records
	err := r.db.WithContext(ctx).Model(&model.Post{}).Where("published = ?", true).Count(&count).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count posts: %w", err)
	}

	// Get paginated published records
	err = r.db.WithContext(ctx).
		Preload("Creator").
		Preload("Tags").
		Where("published = ?", true).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&posts).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get posts: %w", err)
	}

	return posts, count, nil
}

func (r *postRepository) GetPostBySlugAndUsername(ctx context.Context, slug string, username string) (*model.Post, error) {
	var post model.Post
	// We need to find a post with a given slug AND created by a user with the given username.
	err := r.db.WithContext(ctx).
		Preload("Creator").                                 // Preload the Creator
		Preload("Tags").                                    // Preload Tags
		Joins("JOIN users ON users.id = posts.created_by"). // Join with users table
		Where("posts.slug = ? AND users.username = ?", slug, username).
		First(&post).Error // Find the first matching record

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPostNotFound // Or a more specific "post with slug/username not found"
		}
		return nil, fmt.Errorf("failed to get post by slug '%s' and username '%s': %w", slug, username, err)
	}
	// The JOIN and WHERE clause should ensure post.Creator.Username matches,
	// but an explicit check after loading can be added for extra safety if Creator is preloaded.
	// if post.Creator == nil || post.Creator.Username != username {
	//  return nil, ErrPostNotFound // Should not happen if JOIN is correct
	// }
	return &post, nil
}

func (r *postRepository) GetPostByID(ctx context.Context, id string) (*model.Post, error) {
	var post model.Post
	err := r.db.WithContext(ctx).
		Preload("Creator").              // Assuming Post model has a Creator field (struct or ID)
		Preload("Tags").                 // Assuming Post model has a Tags field (slice of Tag)
		First(&post, "id = ?", id).Error // GORM uses primary key by default if just `id` is passed to First

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPostNotFound
		}
		return nil, fmt.Errorf("failed to get post by ID: %w", err)
	}
	return &post, nil
}

func (r *postRepository) GetPostsRandom(ctx context.Context, limit int) ([]*model.Post, error) {
	var randomPosts []*model.Post
	err := r.db.WithContext(ctx).
		Preload("Creator").
		Preload("Tags").
		Where("published = ?", true).
		Order("RANDOM()"). // Works for PostgreSQL and SQLite. For others, might need specific syntax.
		Limit(limit).
		Find(&randomPosts).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get random posts: %w", err)
	}
	return randomPosts, nil
}

func (r *postRepository) GetPostsByCreatedBy(ctx context.Context, createdBy string, offset int, limit int) ([]*model.Post, int64, error) {
	var posts []*model.Post
	var count int64

	// Count posts by created_by (user ID)
	// Assuming model.Post has a 'CreatedBy' field that stores the user ID.
	err := r.db.WithContext(ctx).Model(&model.Post{}).Where("created_by = ?", createdBy).Count(&count).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count posts by creator ID %s: %w", createdBy, err)
	}

	// Get paginated records
	err = r.db.WithContext(ctx).
		Preload("Creator"). // Preload creator details
		Preload("Tags").    // Preload tags
		Where("created_by = ?", createdBy).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&posts).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get posts by creator ID %s: %w", createdBy, err)
	}
	return posts, count, nil
}

// SearchPosts allows searching posts by keyword in title or body.
func (r *postRepository) SearchPosts(ctx context.Context, keyword string, limit int, offset int) ([]*model.Post, int64, error) {
	var posts []*model.Post
	var count int64
	likePattern := "%" + keyword + "%"

	// Count total matching records
	err := r.db.WithContext(ctx).Model(&model.Post{}).
		Where("(title ILIKE ? OR body ILIKE ?) AND published = ?", likePattern, likePattern, true).
		Count(&count).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count posts for search: %w", err)
	}

	// Get paginated records
	err = r.db.WithContext(ctx).
		Preload("Creator").
		Preload("Tags").
		Where("(title ILIKE ? OR body ILIKE ?) AND published = ?", likePattern, likePattern, true).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&posts).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search posts: %w", err)
	}
	return posts, count, nil
}

// GetPostsByTag fetches posts with a specific tag name.
func (r *postRepository) GetPostsByTag(ctx context.Context, tag string, limit int, offset int) ([]*model.Post, int64, error) {
	var posts []*model.Post
	var count int64

	// Join posts and tags tables
	query := r.db.WithContext(ctx).Model(&model.Post{}).
		Joins("JOIN posts_to_tags ON posts_to_tags.post_id = posts.id").
		Joins("JOIN tags ON tags.id = posts_to_tags.tag_id").
		Where("tags.name = ? AND posts.published = ?", tag, true)

	// Count total records
	err := query.Count(&count).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count posts by tag: %w", err)
	}

	// Get paginated records
	err = r.db.WithContext(ctx).Model(&model.Post{}).
		Preload("Creator").
		Preload("Tags").
		Joins("JOIN posts_to_tags ON posts_to_tags.post_id = posts.id").
		Joins("JOIN tags ON tags.id = posts_to_tags.tag_id").
		Where("tags.name = ? AND posts.published = ?", tag, true).
		Order("posts.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&posts).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get posts by tag: %w", err)
	}
	return posts, count, nil
}

// ExistsByID checks if a post exists by its ID.
func (r *postRepository) ExistsByID(ctx context.Context, id string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Post{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check post existence: %w", err)
	}
	return count > 0, nil
}
