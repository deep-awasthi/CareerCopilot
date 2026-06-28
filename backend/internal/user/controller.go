package user

import (
	"github.com/deepawasthi/careercopilot/pkg/middleware"
	"github.com/deepawasthi/careercopilot/pkg/response"
	"github.com/deepawasthi/careercopilot/pkg/validator"
	"github.com/gin-gonic/gin"
	"github.com/deepawasthi/careercopilot/pkg/config"
)

type Controller struct {
	svc Service
}

func NewController(svc Service) *Controller {
	return &Controller{svc: svc}
}

// GetProfile godoc
// @Summary Get user profile
// @Tags profile
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Router /api/v1/profile [get]
func (c *Controller) GetProfile(ctx *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok {
		response.Unauthorized(ctx, "unauthorized")
		return
	}
	profile, err := c.svc.GetProfile(ctx.Request.Context(), userID)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, "profile retrieved", profile)
}

// UpdateProfile godoc
// @Summary Update user profile
// @Tags profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body UpdateProfileRequest true "Profile payload"
// @Success 200 {object} response.Response
// @Router /api/v1/profile [put]
func (c *Controller) UpdateProfile(ctx *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok {
		response.Unauthorized(ctx, "unauthorized")
		return
	}
	var req UpdateProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "invalid request body")
		return
	}
	if err := validator.Validate(&req); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	profile, err := c.svc.UpdateProfile(ctx.Request.Context(), userID, &req)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, "profile updated", profile)
}

func RegisterRoutes(r *gin.RouterGroup, ctrl *Controller, cfg *config.Config) {
	profile := r.Group("/profile")
	profile.Use(middleware.JWTAuth(&cfg.JWT))
	{
		profile.GET("", ctrl.GetProfile)
		profile.PUT("", ctrl.UpdateProfile)
		profile.PATCH("", ctrl.UpdateProfile)
	}
}
