package service

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/campusbooks/campusbooks/internal/constants"
	"github.com/campusbooks/campusbooks/internal/util"
	"github.com/minio/minio-go/v7"
)

// UploadService stores files into MinIO.
type UploadService struct {
	minio  *util.MinIOClient
	logger *slog.Logger
}

// NewUploadService creates an UploadService.
func NewUploadService(minio *util.MinIOClient, logger *slog.Logger) *UploadService {
	return &UploadService{minio: minio, logger: logger}
}

// Upload stores a file and returns a public URL path.
func (s *UploadService) Upload(ctx context.Context, filename string, reader interface {
	Read([]byte) (int, error)
	ReadAt([]byte, int64) (int, error)
	Seek(int64, int) (int64, error)
}, size int64, contentType string) (*string, error) {
	ext := strings.ToLower(filepath.Ext(filename))
	if !allowedExt(ext) {
		return nil, util.NewAppError(http.StatusBadRequest, constants.CodeBadRequest, "仅支持 jpg/png/webp/gif 图片")
	}
	objectKey := fmt.Sprintf("%d-%s%s", time.Now().UnixNano(), randomName(), ext)
	_, err := s.minio.Client.PutObject(ctx, s.minio.Bucket, objectKey, reader, size, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		s.logger.Error(constants.LogUploadFailed, "filename", filename, "error", err)
		return nil, util.NewAppError(http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
	}
	url := "/api/v1/files/" + objectKey
	s.logger.Info(constants.LogUploadSuccess, "object", objectKey, "bucket", s.minio.Bucket)
	return &url, nil
}

func allowedExt(ext string) bool {
	switch ext {
	case ".jpg", ".jpeg", ".png", ".webp", ".gif":
		return true
	default:
		return false
	}
}

func randomName() string {
	return fmt.Sprintf("%x", time.Now().UnixNano())
}
