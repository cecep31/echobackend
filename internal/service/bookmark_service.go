package service

import (
	"context"
	"time"

	"echobackend/internal/dto"
	apperrors "echobackend/internal/errors"
	"echobackend/internal/model"
	"echobackend/internal/repository"
)

type BookmarkService interface {
	ToggleBookmark(ctx context.Context, postID, userID string, req *dto.ToggleBookmarkRequest) (*dto.ToggleBookmarkResponse, error)
	GetBookmarksByUser(ctx context.Context, userID string, folderID *string, limit, offset int) ([]*dto.BookmarkResponse, int64, error)
	UpdateBookmark(ctx context.Context, bookmarkID, userID string, req *dto.UpdateBookmarkRequest) (*dto.BookmarkResponse, error)
	MoveBookmark(ctx context.Context, bookmarkID, userID string, folderID *string) (*dto.BookmarkResponse, error)
	CreateFolder(ctx context.Context, userID string, req *dto.CreateBookmarkFolderRequest) (*dto.BookmarkFolderResponse, error)
	GetFoldersByUser(ctx context.Context, userID string) ([]*dto.BookmarkFolderResponse, error)
	UpdateFolder(ctx context.Context, folderID, userID string, req *dto.UpdateBookmarkFolderRequest) (*dto.BookmarkFolderResponse, error)
	DeleteFolder(ctx context.Context, folderID, userID string) error
}

type bookmarkService struct {
	bookmarkRepo BookmarkRepository
	postRepo     repository.PostRepository
}

type BookmarkRepository = repository.BookmarkRepository

func NewBookmarkService(bookmarkRepo repository.BookmarkRepository, postRepo repository.PostRepository) BookmarkService {
	return &bookmarkService{bookmarkRepo: bookmarkRepo, postRepo: postRepo}
}

func (s *bookmarkService) ToggleBookmark(ctx context.Context, postID, userID string, req *dto.ToggleBookmarkRequest) (*dto.ToggleBookmarkResponse, error) {
	if _, err := s.postRepo.GetPostByID(ctx, postID); err != nil {
		return nil, err
	}

	existing, err := s.bookmarkRepo.FindBookmarkByUserAndPost(ctx, userID, postID)
	if err == nil {
		if err := s.bookmarkRepo.DeleteBookmark(ctx, existing.ID); err != nil {
			return nil, err
		}
		return &dto.ToggleBookmarkResponse{Action: "removed"}, nil
	}
	if err != apperrors.ErrBookmarkNotFound {
		return nil, err
	}

	if req.FolderID != nil && *req.FolderID != "" {
		if _, err := s.bookmarkRepo.FindFolderByID(ctx, *req.FolderID, userID); err != nil {
			return nil, err
		}
	}

	bookmark := &model.PostBookmark{
		PostID:   postID,
		UserID:   userID,
		FolderID: req.FolderID,
		Name:     req.Name,
		Notes:    req.Notes,
	}
	if err := s.bookmarkRepo.CreateBookmark(ctx, bookmark); err != nil {
		return nil, err
	}

	created, err := s.bookmarkRepo.FindBookmarkByID(ctx, bookmark.ID, userID)
	if err != nil {
		return nil, err
	}

	return &dto.ToggleBookmarkResponse{
		Action:   "added",
		Bookmark: dto.BookmarkToResponse(created),
	}, nil
}

func (s *bookmarkService) GetBookmarksByUser(ctx context.Context, userID string, folderID *string, limit, offset int) ([]*dto.BookmarkResponse, int64, error) {
	bookmarks, total, err := s.bookmarkRepo.GetBookmarksByUser(ctx, userID, folderID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	result := make([]*dto.BookmarkResponse, 0, len(bookmarks))
	for _, bookmark := range bookmarks {
		result = append(result, dto.BookmarkToResponse(bookmark))
	}
	return result, total, nil
}

func (s *bookmarkService) UpdateBookmark(ctx context.Context, bookmarkID, userID string, req *dto.UpdateBookmarkRequest) (*dto.BookmarkResponse, error) {
	bookmark, err := s.bookmarkRepo.FindBookmarkByID(ctx, bookmarkID, userID)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		bookmark.Name = req.Name
	}
	if req.Notes != nil {
		bookmark.Notes = req.Notes
	}
	now := time.Now()
	bookmark.UpdatedAt = &now

	if err := s.bookmarkRepo.UpdateBookmark(ctx, bookmark); err != nil {
		return nil, err
	}
	updated, err := s.bookmarkRepo.FindBookmarkByID(ctx, bookmarkID, userID)
	if err != nil {
		return nil, err
	}
	return dto.BookmarkToResponse(updated), nil
}

func (s *bookmarkService) MoveBookmark(ctx context.Context, bookmarkID, userID string, folderID *string) (*dto.BookmarkResponse, error) {
	bookmark, err := s.bookmarkRepo.FindBookmarkByID(ctx, bookmarkID, userID)
	if err != nil {
		return nil, err
	}

	if folderID != nil && *folderID != "" {
		if _, err := s.bookmarkRepo.FindFolderByID(ctx, *folderID, userID); err != nil {
			return nil, err
		}
	}

	bookmark.FolderID = folderID
	now := time.Now()
	bookmark.UpdatedAt = &now
	if err := s.bookmarkRepo.UpdateBookmark(ctx, bookmark); err != nil {
		return nil, err
	}

	updated, err := s.bookmarkRepo.FindBookmarkByID(ctx, bookmarkID, userID)
	if err != nil {
		return nil, err
	}
	return dto.BookmarkToResponse(updated), nil
}

func (s *bookmarkService) CreateFolder(ctx context.Context, userID string, req *dto.CreateBookmarkFolderRequest) (*dto.BookmarkFolderResponse, error) {
	folder := &model.BookmarkFolder{
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
	}
	if err := s.bookmarkRepo.CreateFolder(ctx, folder); err != nil {
		return nil, err
	}
	created, err := s.bookmarkRepo.FindFolderByID(ctx, folder.ID, userID)
	if err != nil {
		return nil, err
	}
	return dto.BookmarkFolderToResponse(created), nil
}

func (s *bookmarkService) GetFoldersByUser(ctx context.Context, userID string) ([]*dto.BookmarkFolderResponse, error) {
	folders, err := s.bookmarkRepo.GetFoldersByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make([]*dto.BookmarkFolderResponse, 0, len(folders))
	for _, folder := range folders {
		result = append(result, dto.BookmarkFolderToResponse(folder))
	}
	return result, nil
}

func (s *bookmarkService) UpdateFolder(ctx context.Context, folderID, userID string, req *dto.UpdateBookmarkFolderRequest) (*dto.BookmarkFolderResponse, error) {
	folder, err := s.bookmarkRepo.FindFolderByID(ctx, folderID, userID)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		folder.Name = *req.Name
	}
	if req.Description != nil {
		folder.Description = req.Description
	}
	now := time.Now()
	folder.UpdatedAt = &now

	if err := s.bookmarkRepo.UpdateFolder(ctx, folder); err != nil {
		return nil, err
	}
	updated, err := s.bookmarkRepo.FindFolderByID(ctx, folderID, userID)
	if err != nil {
		return nil, err
	}
	return dto.BookmarkFolderToResponse(updated), nil
}

func (s *bookmarkService) DeleteFolder(ctx context.Context, folderID, userID string) error {
	return s.bookmarkRepo.DeleteFolder(ctx, folderID, userID)
}
