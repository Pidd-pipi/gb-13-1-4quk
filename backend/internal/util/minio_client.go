package util

import (
	"context"
	"log/slog"
	"time"

	"github.com/campusbooks/campusbooks/internal/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinIOClient wraps the minio-go client with the platform bucket.
type MinIOClient struct {
	Client *minio.Client
	Bucket string
}

// NewMinIOClient connects to MinIO and ensures the bucket exists.
func NewMinIOClient(cfg *config.Config, logger *slog.Logger) (*MinIOClient, error) {
	client, err := minio.New(cfg.MinIOEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinIOUser, cfg.MinIOPassword, ""),
		Secure: cfg.MinIOUseSSL,
	})
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	exists, err := client.BucketExists(ctx, cfg.MinIOBucket)
	if err != nil {
		return nil, err
	}
	if !exists {
		if err := client.MakeBucket(ctx, cfg.MinIOBucket, minio.MakeBucketOptions{}); err != nil {
			return nil, err
		}
	}
	return &MinIOClient{Client: client, Bucket: cfg.MinIOBucket}, nil
}
