package routes

import (
	"spinLuck/config"
	"spinLuck/config/db"
	"spinLuck/internal/features/storage"
	"spinLuck/internal/shared/middleware"
	"spinLuck/internal/shared/utils"

	"github.com/gin-gonic/gin"
)

func StorageRoutes(rg *gin.RouterGroup) {
	rd := config.Rdb
	maker := utils.NewPasetoMaker()
	rpStg := storage.NewPostgresRepository(db.DB)
	stgUs := storage.NewStorageUseCase(rpStg)
	stH := storage.NewStorageHandler(stgUs)

	rg.GET("/storage/file/:fileId", stH.GetByFileId)

	protected := rg.Group("/storage")
	protected.Use(middleware.AuthMiddleware(maker, rd))
	{
		protected.POST("/upload", stH.Create)
		protected.PUT("/upload/:fileId", stH.Update)
		protected.DELETE("/file/:fileId", stH.Delete)
	}
}
