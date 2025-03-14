package storage

import (
	"context"
	"echobackend/config"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioStorage struct {
	client *minio.Client
	bucket string
}

func NewMinioStorage(config *config.Config) *MinioStorage {
	client, err := minio.New(config.MINIO_ENDPOINT, &minio.Options{
		Creds: credentials.NewStaticV4(config.MINIO_ACCESS_KEY, config.MINIO_SECRET_KEY, ""),
	})
	if err != nil {
		panic(err)
	}
	return &MinioStorage{
		client: client,
		bucket: config.MINIO_BUCKET,
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
