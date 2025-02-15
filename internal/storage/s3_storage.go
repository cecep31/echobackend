package storage

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
)

type MinioStorage struct {
	client *minio.Client
	bucket string
}

func NewMinioStorage(client *minio.Client, bucket string) *MinioStorage {
	return &MinioStorage{
		client: client,
		bucket: bucket,
	}
}

func (s *MinioStorage) Save(ctx context.Context, path string, file io.Reader) error {
	_, err := s.client.PutObject(ctx, s.bucket, path, file, -1, minio.PutObjectOptions{})
	return err
}

func (s *MinioStorage) Get(ctx context.Context, path string) (io.ReadCloser, error) {
	obj, err := s.client.GetObject(ctx, s.bucket, path, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func (s *MinioStorage) Delete(ctx context.Context, path string) error {
	return s.client.RemoveObject(ctx, s.bucket, path, minio.RemoveObjectOptions{})
}
