package routes

import (
	"spinLuck/config"
	"spinLuck/config/db"
	"spinLuck/internal/features/auth"
	"spinLuck/internal/features/organizer"
	"spinLuck/internal/shared/middleware"
	"spinLuck/internal/shared/utils"

	"github.com/gin-gonic/gin"
)

func OrganizerRoutes(rg *gin.RouterGroup) {

	rd := config.Rdb
	maker := utils.NewPasetoMaker()
	rpAuth := auth.NewPostgresRepository(db.DB)
	rp := organizer.NewPostgresRepository(db.DB)
	s := organizer.NewOrganizerUseCase(rp, rpAuth)
	h := organizer.NewOrganizerHandler(s)

	protected := rg.Group("/organizer")
	protected.Use(middleware.AuthMiddleware(maker, rd))
	{
		protected.GET("/dashboard-metrics", h.GetDashboardInfoDetailGeneric)
		protected.POST("", h.Create)
		protected.GET("", h.GetByID)
		protected.PUT("", h.Update)
	}
}
