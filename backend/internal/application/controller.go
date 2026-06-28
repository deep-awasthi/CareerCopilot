package application

import (
	"strconv"

	"github.com/deepawasthi/careercopilot/pkg/config"
	"github.com/deepawasthi/careercopilot/pkg/middleware"
	"github.com/deepawasthi/careercopilot/pkg/response"
	"github.com/deepawasthi/careercopilot/pkg/validator"
	"github.com/gin-gonic/gin"
)

type Controller struct{ svc Service }

func NewController(svc Service) *Controller { return &Controller{svc: svc} }

func (c *Controller) Create(ctx *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok { response.Unauthorized(ctx, "unauthorized"); return }
	var req CreateApplicationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil { response.BadRequest(ctx, err.Error()); return }
	if err := validator.Validate(&req); err != nil { response.BadRequest(ctx, err.Error()); return }
	result, err := c.svc.Create(ctx.Request.Context(), userID, &req)
	if err != nil { response.Error(ctx, err); return }
	response.Created(ctx, "application created", result)
}

func (c *Controller) List(ctx *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok { response.Unauthorized(ctx, "unauthorized"); return }
	var filter ApplicationFilter
	_ = ctx.ShouldBindQuery(&filter)
	apps, total, err := c.svc.List(ctx.Request.Context(), userID, &filter)
	if err != nil { response.Error(ctx, err); return }
	response.Paginated(ctx, "applications retrieved", apps, filter.Page, filter.PerPage, total)
}

func (c *Controller) Get(ctx *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok { response.Unauthorized(ctx, "unauthorized"); return }
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)
	app, err := c.svc.Get(ctx.Request.Context(), uint(id), userID)
	if err != nil { response.NotFound(ctx, "application not found"); return }
	response.Success(ctx, "application retrieved", app)
}

func (c *Controller) Update(ctx *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok { response.Unauthorized(ctx, "unauthorized"); return }
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)
	var req UpdateApplicationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil { response.BadRequest(ctx, err.Error()); return }
	app, err := c.svc.Update(ctx.Request.Context(), uint(id), userID, &req)
	if err != nil { response.Error(ctx, err); return }
	response.Success(ctx, "application updated", app)
}

func (c *Controller) Delete(ctx *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok { response.Unauthorized(ctx, "unauthorized"); return }
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err := c.svc.Delete(ctx.Request.Context(), uint(id), userID); err != nil {
		response.Error(ctx, err); return
	}
	response.Success(ctx, "application deleted", nil)
}

func (c *Controller) Stats(ctx *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok { response.Unauthorized(ctx, "unauthorized"); return }
	stats, err := c.svc.GetStats(ctx.Request.Context(), userID)
	if err != nil { response.Error(ctx, err); return }
	response.Success(ctx, "application stats", stats)
}

func RegisterRoutes(r *gin.RouterGroup, ctrl *Controller, cfg *config.Config) {
	apps := r.Group("/applications")
	apps.Use(middleware.JWTAuth(&cfg.JWT))
	{
		apps.POST("", ctrl.Create)
		apps.GET("", ctrl.List)
		apps.GET("/stats", ctrl.Stats)
		apps.GET("/:id", ctrl.Get)
		apps.PUT("/:id", ctrl.Update)
		apps.PATCH("/:id", ctrl.Update)
		apps.DELETE("/:id", ctrl.Delete)
	}
}
