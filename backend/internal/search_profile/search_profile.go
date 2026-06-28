package search_profile

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/deepawasthi/careercopilot/pkg/config"
	"github.com/deepawasthi/careercopilot/pkg/middleware"
	"github.com/deepawasthi/careercopilot/pkg/response"
	"github.com/deepawasthi/careercopilot/pkg/validator"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type SearchProfile struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	UserID        uint           `gorm:"not null;index" json:"user_id"`
	Name          string         `gorm:"not null" json:"name"`
	Keywords      pq.StringArray `gorm:"type:text[]" json:"keywords"`
	ExperienceMin float64        `gorm:"type:decimal(4,1)" json:"experience_min"`
	ExperienceMax float64        `gorm:"type:decimal(4,1)" json:"experience_max"`
	Locations     pq.StringArray `gorm:"type:text[]" json:"locations"`
	SalaryMin     float64        `gorm:"type:decimal(15,2)" json:"salary_min"`
	SalaryMax     float64        `gorm:"type:decimal(15,2)" json:"salary_max"`
	IsRemote      bool           `json:"is_remote"`
	IsHybrid      bool           `json:"is_hybrid"`
	JobType       string         `json:"job_type"`
	IsActive      bool           `gorm:"default:true" json:"is_active"`
	LastRunAt     *time.Time     `json:"last_run_at"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (SearchProfile) TableName() string { return "search_profiles" }

type CreateProfileRequest struct {
	Name          string   `json:"name" validate:"required,min=2,max=255"`
	Keywords      []string `json:"keywords"`
	ExperienceMin float64  `json:"experience_min" validate:"omitempty,gte=0"`
	ExperienceMax float64  `json:"experience_max" validate:"omitempty,gte=0"`
	Locations     []string `json:"locations"`
	SalaryMin     float64  `json:"salary_min" validate:"omitempty,gte=0"`
	SalaryMax     float64  `json:"salary_max" validate:"omitempty,gte=0"`
	IsRemote      bool     `json:"is_remote"`
	IsHybrid      bool     `json:"is_hybrid"`
	JobType       string   `json:"job_type"`
}

type UpdateProfileRequest struct {
	Name          string   `json:"name" validate:"omitempty,min=2,max=255"`
	Keywords      []string `json:"keywords"`
	ExperienceMin float64  `json:"experience_min"`
	ExperienceMax float64  `json:"experience_max"`
	Locations     []string `json:"locations"`
	SalaryMin     float64  `json:"salary_min"`
	SalaryMax     float64  `json:"salary_max"`
	IsRemote      *bool    `json:"is_remote"`
	IsHybrid      *bool    `json:"is_hybrid"`
	JobType       string   `json:"job_type"`
	IsActive      *bool    `json:"is_active"`
}

type Repository interface {
	Create(ctx context.Context, sp *SearchProfile) error
	FindByID(ctx context.Context, id, userID uint) (*SearchProfile, error)
	ListByUser(ctx context.Context, userID uint) ([]*SearchProfile, error)
	ListActive(ctx context.Context) ([]*SearchProfile, error)
	Update(ctx context.Context, sp *SearchProfile) error
	Delete(ctx context.Context, id, userID uint) error
}

type Service interface {
	Create(ctx context.Context, userID uint, req *CreateProfileRequest) (*SearchProfile, error)
	Get(ctx context.Context, id, userID uint) (*SearchProfile, error)
	List(ctx context.Context, userID uint) ([]*SearchProfile, error)
	Update(ctx context.Context, id, userID uint, req *UpdateProfileRequest) (*SearchProfile, error)
	Delete(ctx context.Context, id, userID uint) error
}

type repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) Repository { return &repository{db: db} }

func (r *repository) Create(ctx context.Context, sp *SearchProfile) error {
	return r.db.WithContext(ctx).Create(sp).Error
}

func (r *repository) FindByID(ctx context.Context, id, userID uint) (*SearchProfile, error) {
	var sp SearchProfile
	if err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&sp).Error; err != nil {
		return nil, fmt.Errorf("search profile not found")
	}
	return &sp, nil
}

func (r *repository) ListByUser(ctx context.Context, userID uint) ([]*SearchProfile, error) {
	var sps []*SearchProfile
	r.db.WithContext(ctx).Where("user_id = ?", userID).Order("name ASC").Find(&sps)
	return sps, nil
}

func (r *repository) ListActive(ctx context.Context) ([]*SearchProfile, error) {
	var sps []*SearchProfile
	r.db.WithContext(ctx).Where("is_active = true").Find(&sps)
	return sps, nil
}

func (r *repository) Update(ctx context.Context, sp *SearchProfile) error {
	return r.db.WithContext(ctx).Save(sp).Error
}

func (r *repository) Delete(ctx context.Context, id, userID uint) error {
	return r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&SearchProfile{}).Error
}

type service struct{ repo Repository }

func NewService(repo Repository) Service { return &service{repo: repo} }

func (s *service) Create(ctx context.Context, userID uint, req *CreateProfileRequest) (*SearchProfile, error) {
	sp := &SearchProfile{
		UserID:        userID,
		Name:          req.Name,
		Keywords:      pq.StringArray(req.Keywords),
		ExperienceMin: req.ExperienceMin,
		ExperienceMax: req.ExperienceMax,
		Locations:     pq.StringArray(req.Locations),
		SalaryMin:     req.SalaryMin,
		SalaryMax:     req.SalaryMax,
		IsRemote:      req.IsRemote,
		IsHybrid:      req.IsHybrid,
		JobType:       req.JobType,
		IsActive:      true,
	}
	return sp, s.repo.Create(ctx, sp)
}

func (s *service) Get(ctx context.Context, id, userID uint) (*SearchProfile, error) {
	return s.repo.FindByID(ctx, id, userID)
}

func (s *service) List(ctx context.Context, userID uint) ([]*SearchProfile, error) {
	return s.repo.ListByUser(ctx, userID)
}

func (s *service) Update(ctx context.Context, id, userID uint, req *UpdateProfileRequest) (*SearchProfile, error) {
	sp, err := s.repo.FindByID(ctx, id, userID)
	if err != nil { return nil, err }
	if req.Name != "" { sp.Name = req.Name }
	if len(req.Keywords) > 0 { sp.Keywords = pq.StringArray(req.Keywords) }
	if len(req.Locations) > 0 { sp.Locations = pq.StringArray(req.Locations) }
	if req.ExperienceMin > 0 { sp.ExperienceMin = req.ExperienceMin }
	if req.ExperienceMax > 0 { sp.ExperienceMax = req.ExperienceMax }
	if req.SalaryMin > 0 { sp.SalaryMin = req.SalaryMin }
	if req.SalaryMax > 0 { sp.SalaryMax = req.SalaryMax }
	if req.IsRemote != nil { sp.IsRemote = *req.IsRemote }
	if req.IsHybrid != nil { sp.IsHybrid = *req.IsHybrid }
	if req.JobType != "" { sp.JobType = req.JobType }
	if req.IsActive != nil { sp.IsActive = *req.IsActive }
	return sp, s.repo.Update(ctx, sp)
}

func (s *service) Delete(ctx context.Context, id, userID uint) error {
	return s.repo.Delete(ctx, id, userID)
}

type Controller struct{ svc Service }

func NewController(svc Service) *Controller { return &Controller{svc: svc} }

func (c *Controller) Create(ctx *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok { response.Unauthorized(ctx, "unauthorized"); return }
	var req CreateProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil { response.BadRequest(ctx, err.Error()); return }
	if err := validator.Validate(&req); err != nil { response.BadRequest(ctx, err.Error()); return }
	sp, err := c.svc.Create(ctx.Request.Context(), userID, &req)
	if err != nil { response.Error(ctx, err); return }
	response.Created(ctx, "search profile created", sp)
}

func (c *Controller) List(ctx *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok { response.Unauthorized(ctx, "unauthorized"); return }
	sps, _ := c.svc.List(ctx.Request.Context(), userID)
	response.Success(ctx, "search profiles retrieved", sps)
}

func (c *Controller) Get(ctx *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok { response.Unauthorized(ctx, "unauthorized"); return }
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)
	sp, err := c.svc.Get(ctx.Request.Context(), uint(id), userID)
	if err != nil { response.NotFound(ctx, "search profile not found"); return }
	response.Success(ctx, "search profile retrieved", sp)
}

func (c *Controller) Update(ctx *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok { response.Unauthorized(ctx, "unauthorized"); return }
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)
	var req UpdateProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil { response.BadRequest(ctx, err.Error()); return }
	sp, err := c.svc.Update(ctx.Request.Context(), uint(id), userID, &req)
	if err != nil { response.Error(ctx, err); return }
	response.Success(ctx, "search profile updated", sp)
}

func (c *Controller) Delete(ctx *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok { response.Unauthorized(ctx, "unauthorized"); return }
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err := c.svc.Delete(ctx.Request.Context(), uint(id), userID); err != nil {
		response.Error(ctx, err); return
	}
	response.Success(ctx, "search profile deleted", nil)
}

func RegisterRoutes(r *gin.RouterGroup, ctrl *Controller, cfg *config.Config) {
	sp := r.Group("/search-profiles")
	sp.Use(middleware.JWTAuth(&cfg.JWT))
	{
		sp.POST("", ctrl.Create)
		sp.GET("", ctrl.List)
		sp.GET("/:id", ctrl.Get)
		sp.PUT("/:id", ctrl.Update)
		sp.DELETE("/:id", ctrl.Delete)
	}
}
