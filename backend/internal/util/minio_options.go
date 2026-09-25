package util

import "github.com/minio/minio-go/v7"

// GetObjectOptions is a small wrapper for minio GetObject options.
type GetObjectOptions struct{}

// ToMinio converts the wrapper into minio.GetObjectOptions.
func (o GetObjectOptions) ToMinio() minio.GetObjectOptions { return minio.GetObjectOptions{} }
