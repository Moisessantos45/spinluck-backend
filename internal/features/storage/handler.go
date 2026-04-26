package storage

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"spinLuck/internal/shared/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type StorageHandler struct {
	uc StorageServiceInterface
}

func NewStorageHandler(useCase StorageServiceInterface) *StorageHandler {
	return &StorageHandler{uc: useCase}
}

func (h *StorageHandler) deleteExistingFile(ctx context.Context, fileID string) (string, error) {

	fileData, err := h.uc.GetByFileId(ctx, fileID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("No archivo anterior para fileId %s", fileID)
		return fileID, nil
	}
	if err != nil {
		log.Printf("Error fetching existing photo data: %v", err)
		return fileID, fmt.Errorf("get file: %w", err)
	}

	if fileData != nil && fileData.Path != "" {
		if delErr := utils.DeleteFile(fileData.Path); delErr != nil {
			log.Printf("Warning: no pudo borrar %s: %v", fileData.Path, delErr)
		}
	}

	return fileID, nil
}

func (h *StorageHandler) GetByFileId(c *gin.Context) {
	fileId := c.Param("fileId")

	file, err := h.uc.GetByFileId(c.Request.Context(), fileId)
	if err != nil {
		c.Status(404)
		return
	}

	c.Header("Cache-Control", "public, max-age=3600")

	c.File(file.Path)
}

func (h *StorageHandler) Create(c *gin.Context) {
	_, userID, err := utils.ExtractedParamsJwt(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized: " + err.Error()})
		return
	}

	file, size, mimeType, err := utils.ExtractImage(c, "image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	idFile, path, err := utils.SaveFile("prizes/photos", "", file, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error al guardar la imagen: " + err.Error()})
		return
	}

	result, err := h.uc.Create(c.Request.Context(), userID, idFile, path, mimeType, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error al crear el registro de almacenamiento: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result, "message": "Archivo subido exitosamente"})
}

func (h *StorageHandler) Update(c *gin.Context) {
	_, userID, err := utils.ExtractedParamsJwt(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized: " + err.Error()})
		return
	}

	fileId := c.Param("fileId")
	if fileId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "fileId es requerido"})
		return
	}

	existingFileId, err := h.deleteExistingFile(c.Request.Context(), fileId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error al eliminar el archivo existente: " + err.Error()})
		return
	}

	file, size, mimeType, err := utils.ExtractImage(c, "image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	idFile, path, err := utils.SaveFile("prizes/photos", existingFileId, file, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error al guardar la imagen: " + err.Error()})
		return
	}

	result, err := h.uc.Update(c.Request.Context(), userID, idFile, path, mimeType, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error al actualizar el registro de almacenamiento: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result, "message": "Archivo actualizado exitosamente"})
}

func (h *StorageHandler) Delete(c *gin.Context) {
	_, userID, err := utils.ExtractedParamsJwt(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized: " + err.Error()})
		return
	}

	fileId := c.Param("fileId")

	if err := h.uc.Delete(c.Request.Context(), userID, fileId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error al eliminar el registro de almacenamiento: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Archivo eliminado exitosamente"})
}
