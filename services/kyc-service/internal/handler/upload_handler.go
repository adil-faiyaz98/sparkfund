package handler

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	maxFileSize  = 5 * 1024 * 1024 // 5MB
	allowedTypes = ".jpg,.jpeg,.png,.pdf"
)

type UploadHandler struct {
	uploadDir string
}

func NewUploadHandler(uploadDir string) *UploadHandler {
	return &UploadHandler{
		uploadDir: uploadDir,
	}
}

func (h *UploadHandler) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": "No file uploaded"})
		return
	}

	// Validate file size
	if file.Size > maxFileSize {
		c.JSON(400, gin.H{"error": "File size exceeds 5MB limit"})
		return
	}

	// Validate file type
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !isAllowedFileType(ext) {
		c.JSON(400, gin.H{"error": fmt.Sprintf("File type not allowed. Allowed types: %s", allowedTypes)})
		return
	}

	// Generate unique filename
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	filepath := filepath.Join(h.uploadDir, filename)

	// Save file
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(500, gin.H{"error": "Failed to save file"})
		return
	}

	c.JSON(200, gin.H{
		"filename": filename,
		"path":     filepath,
		"size":     file.Size,
		"type":     file.Header.Get("Content-Type"),
	})
}

func (h *UploadHandler) ValidateFile(file *multipart.FileHeader) error {
	// Validate file size
	if file.Size > maxFileSize {
		return fmt.Errorf("file size exceeds 5MB limit")
	}

	// Validate file type
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !isAllowedFileType(ext) {
		return fmt.Errorf("file type not allowed. Allowed types: %s", allowedTypes)
	}

	return nil
}

func isAllowedFileType(ext string) bool {
	allowed := strings.Split(allowedTypes, ",")
	for _, t := range allowed {
		if ext == t {
			return true
		}
	}
	return false
}
