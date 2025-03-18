package storage

import (
	"context"
	"echobackend/config"
	"io"
	"net/http"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioStorage struct {
	client *minio.Client
	bucket string
}

func NewMinioStorage(config *config.Config) *MinioStorage {
	// Set custom transport with timeouts
	transport := &http.Transport{
		ResponseHeaderTimeout: 30 * time.Second,
		IdleConnTimeout:       90 * time.Second,
		MaxIdleConnsPerHost:   100,
		MaxConnsPerHost:       100,
	}

	client, err := minio.New(config.MINIO_ENDPOINT, &minio.Options{
		Creds:     credentials.NewStaticV4(config.MINIO_ACCESS_KEY, config.MINIO_SECRET_KEY, ""),
		Secure:    config.MINIO_USE_SSL,
		Transport: transport,
	})
	if err != nil {
		panic(err)
	}

	// Create bucket if it doesn't exist
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	exists, err := client.BucketExists(ctx, config.MINIO_BUCKET)
	if err != nil {
		panic(err)
	}

	if !exists {
		err = client.MakeBucket(ctx, config.MINIO_BUCKET, minio.MakeBucketOptions{})
		if err != nil {
			// Check if another process created the bucket before us
			exists, errCheck := client.BucketExists(ctx, config.MINIO_BUCKET)
			if errCheck != nil || !exists {
				panic(err)
			}
		}
	}

	return &MinioStorage{
		client: client,
		bucket: config.MINIO_BUCKET,
	}
}

func (s *MinioStorage) Save(ctx context.Context, path string, file io.Reader) error {
	// Add timeout to context
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// If file is a ReadCloser, ensure it's closed after the operation
	if rc, ok := file.(io.ReadCloser); ok {
		defer rc.Close()
	}

	_, err := s.client.PutObject(ctx, s.bucket, path, file, -1, minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})
	return err
}

func (s *MinioStorage) Get(ctx context.Context, path string) (io.ReadCloser, error) {
	// Add timeout to context
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	
	obj, err := s.client.GetObject(ctx, s.bucket, path, minio.GetObjectOptions{})
	if err != nil {
		cancel() // Cancel context on error
		return nil, err
	}
	
	// Validate object exists by reading a header
	_, err = obj.Stat()
	if err != nil {
		cancel() // Cancel context on error
		obj.Close() // Close on error
		return nil, err
	}
	
	// Return a wrapper that will cancel the context when closed
	return &readCloserWithCancel{obj, cancel}, nil
}

func (s *MinioStorage) Delete(ctx context.Context, path string) error {
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
