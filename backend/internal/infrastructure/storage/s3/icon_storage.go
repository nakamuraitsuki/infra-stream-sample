package s3

import (
	"bytes"
	"context"
	"fmt"

	"example.com/m/internal/domain/user"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type iconStorage struct {
	client     *s3.Client
	bucketName string
}

func NewIconStorage(clientSet *S3ClientSet, cfg Config) user.IconStorage {
	return &iconStorage{
		client:     clientSet.Client,
		bucketName: cfg.BucketName,
	}
}

func (s *iconStorage) Upload(ctx context.Context, key string, data []byte) error {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String("image/png"),
	})
	if err != nil {
		return fmt.Errorf("failed to upload icon to s3: %w", err)
	}
	return nil
}

func (s *iconStorage) Delete(ctx context.Context, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete icon from s3: %w", err)
	}
	return nil
}
