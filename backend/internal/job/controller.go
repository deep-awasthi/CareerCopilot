package job

import (
	"strconv"

	"github.com/deepawasthi/careercopilot/pkg/config"
	"github.com/deepawasthi/careercopilot/pkg/middleware"
	"github.com/deepawasthi/careercopilot/pkg/response"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	svc Service
}

func NewController(svc Service) *Controller {
	return &Controller{svc: svc}
}

// ListJobs godoc
// @Summary List jobs with filters
// @Tags jobs
// @Produce json
// @Security BearerAuth
// @Param q query string false "Search query"
// @Param location query string false "Location"
// @Param remote query bool false "Remote only"
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Success 200 {object} response.PaginatedResponse
// @Router /api/v1/jobs [get]
func (c *Controller) ListJobs(ctx *gin.Context) {
	userID, _ := middleware.GetUserIDFromContext(ctx)
	var filter JobFilter
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		response.BadRequest(ctx, "invalid query parameters")
		return
	}
	jobs, total, err := c.svc.ListJobs(ctx.Request.Context(), &filter, userID)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	page := filter.Page
	if page < 1 {
		page = 1
	}
	perPage := filter.PerPage
	if perPage < 1 {
		perPage = 20
	}
	response.Paginated(ctx, "jobs retrieved", jobs, page, perPage, total)
}

// GetJob godoc
// @Summary Get a single job by ID
// @Tags jobs
// @Produce json
// @Security BearerAuth
// @Param id path int true "Job ID"
// @Success 200 {object} response.Response
// @Router /api/v1/jobs/{id} [get]
func (c *Controller) GetJob(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(ctx, "invalid job id")
		return
	}
	job, err := c.svc.GetJob(ctx.Request.Context(), uint(id))
	if err != nil {
		response.NotFound(ctx, "job not found")
		return
	}
	response.Success(ctx, "job retrieved", job)
}

// GetMatchedJobs godoc
// @Summary Get jobs matched to user profile
// @Tags jobs
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.PaginatedResponse
// @Router /api/v1/jobs/matched [get]
func (c *Controller) GetMatchedJobs(ctx *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok {
		response.Unauthorized(ctx, "unauthorized")
		return
	}
	var filter JobFilter
	_ = ctx.ShouldBindQuery(&filter)
	jobs, total, err := c.svc.GetMatchedJobs(ctx.Request.Context(), userID, &filter)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Paginated(ctx, "matched jobs retrieved", jobs, filter.Page, filter.PerPage, total)
}

func RegisterRoutes(r *gin.RouterGroup, ctrl *Controller, cfg *config.Config) {
	jobs := r.Group("/jobs")
	jobs.Use(middleware.JWTAuth(&cfg.JWT))
	{
		jobs.GET("", ctrl.ListJobs)
		jobs.GET("/matched", ctrl.GetMatchedJobs)
		jobs.GET("/:id", ctrl.GetJob)
	}
}
