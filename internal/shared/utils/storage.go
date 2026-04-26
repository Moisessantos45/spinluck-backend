package utils

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func AppDir() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		dir, err := os.Getwd()
		if err != nil {
			return "", err
		}
		return dir, nil
	}

	dir := filepath.Dir(exe)
	if strings.Contains(dir, "go-build") {
		dir, err = os.Getwd()
		if err != nil {
			return "", err
		}

	}
	return dir, nil
}

func SaveFile(dir string, fileId string, file *multipart.FileHeader, c *gin.Context) (string, string, error) {
	appDir, err := AppDir()
	if err != nil {
		return "", "", err
	}

	uploadsDir := filepath.Join(appDir, "cmd", "uploads", dir)
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		return "", "", fmt.Errorf("mkdir: %w", err)
	}

	id := fileId
	if strings.TrimSpace(id) == "" {
		id = fmt.Sprintf("%d", time.Now().UnixNano())
	}

	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s_%s", id, ext)
	dst := filepath.Join(uploadsDir, filename)

	if err := c.SaveUploadedFile(file, dst); err != nil {
		return "", "", fmt.Errorf("save file: %w", err)
	}

	return id, dst, nil
}

func DeleteFile(filePath string) error {
	if strings.TrimSpace(filePath) == "" {
		return fmt.Errorf("file path is required")
	}

	err := os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("delete file: %w", err)
	}

	return nil
}
