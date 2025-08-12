package storage

import (
	"context"
	"echobackend/config"
	"io"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Storage struct {
	client *s3.Client
	bucket string
}

func NewS3Storage(cfg *config.Config) *S3Storage {
	// Load AWS configuration
	awsConfig, err := awsconfig.LoadDefaultConfig(context.TODO(),
		awsconfig.WithRegion("ap-southeast-1"), // You might want to make this configurable
	)
	if err != nil {
		log.Printf("Failed to load AWS config: %v", err)
		return nil
	}

	client := s3.NewFromConfig(awsConfig)

	s3Client := &S3Storage{
		client: client,
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

	uploader := manager.NewUploader(s.client)
	_, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
		Body:   file,
	})

	return err
}

func (s *S3Storage) Get(ctx context.Context, path string) (io.ReadCloser, error) {
	// Add timeout to context
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	output, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return nil, err
	}

	// Return a wrapper that will cancel the context when closed
	return &readCloserWithCancel{output.Body, cancel}, nil
}

func (s *S3Storage) Delete(ctx context.Context, path string) error {
	// Add timeout to context
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
	})
	return err
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
