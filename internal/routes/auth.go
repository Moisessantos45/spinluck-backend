package routes

import (
	"spinLuck/config"
	"spinLuck/config/db"
	"spinLuck/internal/features/auth"
	"spinLuck/internal/shared/middleware"
	"spinLuck/internal/shared/utils"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(rg *gin.RouterGroup) {
	rd := config.Rdb
	maker := utils.NewPasetoMaker()

	rp := auth.NewPostgresRepository(db.DB)
	s := auth.NewAuthUseCase(rp, rd, maker)
	h := auth.NewAuthHandler(s)

	rg.POST("/login", h.Login)
	rg.POST("/forward-email-verification", h.ForwardEmailVerification)
	rg.POST("/forgot-password", h.SendPasswordReset)
	rg.POST("/register", h.Register)
	rg.POST("/logout", h.Logout)
	rg.POST("/refresh-token", h.RefreshToken)

	preAuth := rg.Group("/")
	preAuth.Use(middleware.PreAuthMiddleware(rd))
	{
		preAuth.POST("/2fa/verify", h.Verify2FALogin)
	}

	protected := rg.Group("/")
	protected.Use(middleware.AuthMiddleware(maker, rd))
	{
		protected.GET("/confirm-account", h.ConfirmAccount)
		protected.GET("/session", h.GetSession)
		protected.POST("/verify-email", h.VerifyEmail)
		protected.POST("/verify-two-factor", h.TestVerifyTOTP)
		protected.PATCH("/reset-password", h.ResetPassword)
		protected.PATCH("/change-password", h.UpdatePassword)
		protected.PATCH("/enable-two-factor", h.EnableTWOFA)
		protected.PATCH("/disable-two-factor", h.DisableTWOFA)

		protected.GET("/generate-two-factor", h.TestingCrateTWOFA)
	}
}

