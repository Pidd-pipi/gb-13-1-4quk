package dto

// UploadResponse 文件上传响应。
type UploadResponse struct {
	URL      string `json:"url"`
	Object   string `json:"object"`
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
}
