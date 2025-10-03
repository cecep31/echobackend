package storage

import (
	"bytes"
	"context"
	"echobackend/config"
	"io"
	"log"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3Storage struct {
	client *minio.Client
	bucket string
}

func NewS3Storage(cfg *config.Config) *S3Storage {
	// Initialize MinIO client
	minioClient, err := minio.New(cfg.S3.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.S3.AccessKey, cfg.S3.SecretKey, ""),
		Secure: cfg.S3.UseSSL,
	})
	if err != nil {
		log.Printf("Failed to create MinIO client: %v", err)
		return nil
	}

	s3Client := &S3Storage{
		client: minioClient,
		bucket: cfg.S3.Bucket,
	}

	return s3Client
}

func (s *S3Storage) Save(ctx context.Context, path string, file io.Reader) error {
	// Add timeout to context
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// If file is a ReadCloser, ensure it's closed after the operation
	if rc, ok := file.(io.ReadCloser); ok {
		defer rc.Close()
	}

	// Read the entire file into memory to get the size
	// For larger files, consider using a streaming approach with known size
	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	_, err = s.client.PutObject(ctx, s.bucket, path, bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{})
	return err
}

func (s *S3Storage) Get(ctx context.Context, path string) (io.ReadCloser, error) {
	// Add timeout to context
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	object, err := s.client.GetObject(ctx, s.bucket, path, minio.GetObjectOptions{})
	if err != nil {
		cancel()
		return nil, err
	}

	// Return a wrapper that will cancel the context when closed
	return &readCloserWithCancel{object, cancel}, nil
}

func (s *S3Storage) Delete(ctx context.Context, path string) error {
	// Add timeout to context
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return s.client.RemoveObject(ctx, s.bucket, path, minio.RemoveObjectOptions{})
}

// readCloserWithCancel wraps a ReadCloser with a context cancellation function
type readCloserWithCancel struct {
	io.ReadCloser
	cancel context.CancelFunc
}

func (r *readCloserWithCancel) Close() error {
	err := r.ReadCloser.Close()
	r.cancel() // Cancel the context when closing
	return err
}
