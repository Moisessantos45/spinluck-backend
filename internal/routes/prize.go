package routes

import (
	"spinLuck/config"
	"spinLuck/config/db"
	"spinLuck/internal/features/auth"
	"spinLuck/internal/features/organizer"
	"spinLuck/internal/features/prize"
	"spinLuck/internal/features/storage"
	"spinLuck/internal/shared/middleware"
	"spinLuck/internal/shared/utils"

	"github.com/gin-gonic/gin"
)

func PrizeRoutes(rg *gin.RouterGroup) {
	rd := config.Rdb
	maker := utils.NewPasetoMaker()

	rp := prize.NewPostgresRepository(db.DB)
	rpOrg := organizer.NewPostgresRepository(db.DB)
	rpAuth := auth.NewPostgresRepository(db.DB)
	rpStg := storage.NewPostgresRepository(db.DB)

	ucOrg := organizer.NewOrganizerUseCase(rpOrg, rpAuth)
	ucStg := storage.NewStorageUseCase(rpStg)
	uc := prize.NewPrizeUseCase(rp, ucOrg, ucStg)
	h := prize.NewPrizeHandler(uc)

	rg.GET("/prize/public/raffle/:slug", h.GetAllByIdRafflePublic)

	protected := rg.Group("/prize")
	protected.Use(middleware.AuthMiddleware(maker, rd))
	{
		protected.GET("/raffle/:raffleID", h.GetAllByIdRaffleByOrganizerID)
		protected.GET("/prizes/:id", h.GetByID)
		protected.POST("", h.Create)
		protected.PUT("/:id", h.Update)
		protected.DELETE("/raffle/:raffleID/prize/:prizeID", h.Delete)
	}
}
