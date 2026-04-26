package routes

import (
	"spinLuck/config"
	"spinLuck/config/db"
	"spinLuck/internal/features/auth"
	"spinLuck/internal/features/organizer"
	"spinLuck/internal/features/prize"
	"spinLuck/internal/features/raffle"
	"spinLuck/internal/features/storage"
	"spinLuck/internal/shared/middleware"
	"spinLuck/internal/shared/utils"

	"github.com/gin-gonic/gin"
)

func RaffleRoutes(rg *gin.RouterGroup) {
	rd := config.Rdb
	maker := utils.NewPasetoMaker()

	rpOrg := organizer.NewPostgresRepository(db.DB)
	rpPrz := prize.NewPostgresRepository(db.DB)
	rpAuth := auth.NewPostgresRepository(db.DB)
	ucOrg := organizer.NewOrganizerUseCase(rpOrg, rpAuth)
	rpStg := storage.NewPostgresRepository(db.DB)
	rpRaff := raffle.NewPostgresRepository(db.DB)

	ucStg := storage.NewStorageUseCase(rpStg)
	ucPrz := prize.NewPrizeUseCase(rpPrz, ucOrg, ucStg)

	uc := raffle.NewRaffleUseCase(rpRaff, rpOrg, ucStg, ucPrz)
	h := raffle.NewRaffleHandler(uc)

	rg.GET("/raffle/slug/:slug", h.GetBySlug)

	protected := rg.Group("/raffle")
	protected.Use(middleware.AuthMiddleware(maker, rd))
	{
		protected.GET("/status", h.GetAllStatus)
		protected.GET("/organizer", h.GetAllInfoGeneric)
		protected.GET("/organizer/recent", h.GetAllRecentInfoGeneric)
		protected.GET("/:id", h.GetByID)
		protected.POST("", h.Create)
		protected.PUT("/:id", h.Update)
	}
}
