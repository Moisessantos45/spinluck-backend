package utils

import (
	"fmt"
	"mime"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func ExtractFileIdFromURL(filename string) string {
	if strings.TrimSpace(filename) == "" {
		return ""
	}

	parts := strings.Split(filename, "/")
	if len(parts) == 0 {
		return ""
	}

	// Retorna el último elemento del slice, que es el ID del archivo
	return parts[len(parts)-1]
}

// Para IMÁGENES (JPG/JPEG/PNG, máx 2MB)
func ExtractImage(c *gin.Context, fieldName string) (*multipart.FileHeader, int64, string, error) {
	if err := c.Request.ParseMultipartForm(2 << 20); err != nil {
		return nil, 0, "", fmt.Errorf("parse form: %w", err)
	}

	file, err := c.FormFile(fieldName)
	if err != nil {
		return nil, 0, "", fmt.Errorf("archivo '%s' requerido", fieldName)
	}

	if file.Size > 2<<20 {
		return nil, 0, "", fmt.Errorf("archivo muy grande (máx 2MB)")
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return nil, 0, "", fmt.Errorf("solo JPG/JPEG/PNG permitidos")
	}

	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" || (!strings.HasPrefix(mimeType, "image/") && mimeType != "image/png" && mimeType != "image/jpeg") {
		return nil, 0, "", fmt.Errorf("archivo no es imagen válida")
	}

	return file, file.Size, mimeType, nil
}

// Para PDF (solo PDF, máx 5MB)
func ExtractPDF(c *gin.Context, fieldName string) (*multipart.FileHeader, int64, string, error) {
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
		return nil, 0, "", fmt.Errorf("parse form: %w", err)
	}

	file, err := c.FormFile(fieldName)
	if err != nil {
		return nil, 0, "", fmt.Errorf("archivo '%s' requerido", fieldName)
	}

	if file.Size > 5<<20 {
		return nil, 0, "", fmt.Errorf("archivo muy grande (máx 10MB)")
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".pdf" {
		return nil, 0, "", fmt.Errorf("solo PDF permitidos")
	}

	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "application/pdf"
	}

	return file, file.Size, mimeType, nil
}
