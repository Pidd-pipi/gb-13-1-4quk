package handler

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"

	"github.com/campusbooks/campusbooks/internal/dto"
	"github.com/campusbooks/campusbooks/internal/service"
	"github.com/campusbooks/campusbooks/internal/util"
)

// UploadHandler handles file upload/download.
type UploadHandler struct {
	uploadService *service.UploadService
	minio         *util.MinIOClient
	logger        *slog.Logger
}

// NewUploadHandler creates an UploadHandler.
func NewUploadHandler(uploadService *service.UploadService, minio *util.MinIOClient, logger *slog.Logger) *UploadHandler {
	return &UploadHandler{uploadService: uploadService, minio: minio, logger: logger}
}

// Upload stores an image file (avatar / book photos).
func (h *UploadHandler) Upload(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, util.NewAppError(http.StatusBadRequest, 40000, "缺少文件字段 file"))
		return
	}
	if fileHeader.Size > 5*1024*1024 {
		c.AbortWithError(http.StatusBadRequest, util.NewAppError(http.StatusBadRequest, 40000, "文件大小不能超过 5MB"))
		return
	}
	file, err := fileHeader.Open()
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, util.NewAppError(http.StatusBadRequest, 40000, "读取文件失败"))
		return
	}
	defer file.Close()
	contentType := fileHeader.Header.Get("Content-Type")
	url, err := h.uploadService.Upload(c.Request.Context(), fileHeader.Filename, file, fileHeader.Size, contentType)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.UploadResponse{URL: *url, Object: *url, Filename: fileHeader.Filename, Size: fileHeader.Size}))
}

// Get proxies a stored object back to the client.
func (h *UploadHandler) Get(c *gin.Context) {
	key := c.Param("key")
	obj, err := h.minio.Client.GetObject(c.Request.Context(), h.minio.Bucket, key, minio.GetObjectOptions{})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, dto.Fail(40400, "文件不存在"))
		return
	}
	defer obj.Close()
	stat, err := obj.Stat()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, dto.Fail(40400, "文件不存在"))
		return
	}
	c.Header("Content-Type", stat.ContentType)
	c.Header("Cache-Control", "public, max-age=86400")
	http.ServeContent(c.Writer, c.Request, key, stat.LastModified, obj)
}
