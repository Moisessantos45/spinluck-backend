package routes

import (
	"spinLuck/config"
	"spinLuck/config/db"
	"spinLuck/internal/features/organizer"
	"spinLuck/internal/features/ticket"
	"spinLuck/internal/shared/middleware"
	"spinLuck/internal/shared/utils"

	"github.com/gin-gonic/gin"
)

func TicketRoutes(rg *gin.RouterGroup) {
	rd := config.Rdb
	maker := utils.NewPasetoMaker()

	rpOrg := organizer.NewPostgresRepository(db.DB)
	rpTkt := ticket.NewPostgresRepository(db.DB)
	uc := ticket.NewTicketUseCase(rpTkt, rpOrg)
	h := ticket.NewTicketHandler(uc)

	rg.GET("/ticket/raffle/:raffleID/available", h.GetRandomAvailableTicket)
	rg.GET("/ticket/raffle/:raffleID/search", h.FindByNumberAndRaffleID)

	protected := rg.Group("/ticket")
	protected.Use(middleware.AuthMiddleware(maker, rd))
	{
		protected.GET("/status", h.GetAllStatus)
		protected.GET("/raffle/:raffleID", h.GetAllByRaffleIDToOrgaizer)
		protected.GET("/raffle/:raffleID/sold", h.GetRandomSoldTicket)
		protected.GET("/voucher/:id", h.GenerateVoucher)
		protected.GET("/:id", h.GetByID)
		protected.POST("", h.Create)
		protected.PUT("/:id", h.Update)
		protected.PATCH("/raffles/:raffleID/tickets/:id/winner", h.UpdateWinner)
	}
}
