package resume

import (
	"github.com/deepawasthi/careercopilot/pkg/config"
	"github.com/deepawasthi/careercopilot/pkg/middleware"
	"github.com/deepawasthi/careercopilot/pkg/response"
	"github.com/deepawasthi/careercopilot/pkg/validator"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	svc Service
}

func NewController(svc Service) *Controller {
	return &Controller{svc: svc}
}

// Submit godoc
// @Summary Submit or update resume (plain text)
// @Tags resume
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body SubmitResumeRequest true "Resume text"
// @Success 200 {object} response.Response
// @Router /api/v1/resume [post]
func (c *Controller) Submit(ctx *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok {
		response.Unauthorized(ctx, "unauthorized")
		return
	}
	var req SubmitResumeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "invalid request body")
		return
	}
	if err := validator.Validate(&req); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	result, err := c.svc.Submit(ctx.Request.Context(), userID, &req)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, "resume parsed and saved", result)
}

// Get godoc
// @Summary Get parsed resume
// @Tags resume
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Router /api/v1/resume [get]
func (c *Controller) Get(ctx *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok {
		response.Unauthorized(ctx, "unauthorized")
		return
	}
	result, err := c.svc.Get(ctx.Request.Context(), userID)
	if err != nil {
		response.NotFound(ctx, "no resume found")
		return
	}
	response.Success(ctx, "resume retrieved", result)
}

func RegisterRoutes(r *gin.RouterGroup, ctrl *Controller, cfg *config.Config) {
	resume := r.Group("/resume")
	resume.Use(middleware.JWTAuth(&cfg.JWT))
	{
		resume.POST("", ctrl.Submit)
		resume.GET("", ctrl.Get)
	}
}
