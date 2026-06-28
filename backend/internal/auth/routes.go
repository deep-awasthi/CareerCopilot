package auth

import (
	"github.com/deepawasthi/careercopilot/pkg/config"
	"github.com/deepawasthi/careercopilot/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, ctrl *Controller, cfg *config.Config) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", ctrl.Register)
		auth.POST("/login", ctrl.Login)
		auth.POST("/refresh", ctrl.RefreshToken)
		auth.POST("/password-reset", ctrl.RequestPasswordReset)
		auth.POST("/password-reset/confirm", ctrl.ResetPassword)

		// Protected routes
		protected := auth.Group("")
		protected.Use(middleware.JWTAuth(&cfg.JWT))
		{
			protected.GET("/me", ctrl.Me)
			protected.POST("/logout", ctrl.Logout)
			protected.POST("/change-password", ctrl.ChangePassword)
		}
	}
}
