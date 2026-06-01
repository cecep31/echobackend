package service

import (
	"context"
	"errors"
	"mime/multipart"
	"testing"

	apperrors "echobackend/internal/errors"
)

func TestUploadImagePostsRejectsFilesLargerThanOneMiB(t *testing.T) {
	svc := NewPostService(&mockPostRepo{}, nil, nil, nil)

	err := svc.UploadImagePosts(context.Background(), &multipart.FileHeader{
		Filename: "large.jpg",
		Size:     maxPostImageSize + 1,
	})

	if !errors.Is(err, apperrors.ErrFileTooLarge) {
		t.Fatalf("expected ErrFileTooLarge, got %v", err)
	}
}
