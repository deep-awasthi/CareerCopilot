package referral

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
	"gorm.io/gorm"
)

type ReferralStatus string

const (
	StatusNotContacted    ReferralStatus = "not_contacted"
	StatusContacted       ReferralStatus = "contacted"
	StatusFollowUp        ReferralStatus = "follow_up"
	StatusReferralReceived ReferralStatus = "referral_received"
	StatusApplied         ReferralStatus = "applied"
	StatusRejected        ReferralStatus = "rejected"
)

type Referral struct {
	ID                    uint           `gorm:"primarykey" json:"id"`
	UserID                uint           `gorm:"not null;index" json:"user_id"`
	CompanyID             *uint          `json:"company_id"`
	ReferrerName          string         `gorm:"not null" json:"referrer_name"`
	ReferrerDesignation   string         `json:"referrer_designation"`
	ReferrerDepartment    string         `json:"referrer_department"`
	ReferrerOfficeLocation string        `json:"referrer_office_location"`
	ReferrerProfileURL    string         `json:"referrer_profile_url"`
	Status                ReferralStatus `gorm:"type:referral_status;default:'not_contacted'" json:"status"`
	Notes                 string         `gorm:"type:text" json:"notes"`
	ContactedAt           *time.Time     `json:"contacted_at"`
	FollowUpAt            *time.Time     `json:"follow_up_at"`
	ReferralReceivedAt    *time.Time     `json:"referral_received_at"`
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
	DeletedAt             gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Referral) TableName() string { return "referrals" }

type CreateReferralRequest struct {
	CompanyID             *uint  `json:"company_id"`
	ReferrerName          string `json:"referrer_name" validate:"required"`
	ReferrerDesignation   string `json:"referrer_designation"`
	ReferrerDepartment    string `json:"referrer_department"`
	ReferrerOfficeLocation string `json:"referrer_office_location"`
	ReferrerProfileURL    string `json:"referrer_profile_url"`
	Notes                 string `json:"notes"`
}

type UpdateReferralRequest struct {
	Status     string `json:"status" validate:"omitempty,oneof=not_contacted contacted follow_up referral_received applied rejected"`
	Notes      string `json:"notes"`
	FollowUpAt *time.Time `json:"follow_up_at"`
}

type ReferralFilter struct {
	CompanyID uint   `form:"company_id"`
	Status    string `form:"status"`
	Page      int    `form:"page,default=1"`
	PerPage   int    `form:"per_page,default=20"`
}

type Repository interface {
	Create(ctx context.Context, r *Referral) error
	FindByID(ctx context.Context, id, userID uint) (*Referral, error)
	List(ctx context.Context, userID uint, filter *ReferralFilter) ([]*Referral, int64, error)
	Update(ctx context.Context, r *Referral) error
	Delete(ctx context.Context, id, userID uint) error
}

type Service interface {
	Create(ctx context.Context, userID uint, req *CreateReferralRequest) (*Referral, error)
	Get(ctx context.Context, id, userID uint) (*Referral, error)
	List(ctx context.Context, userID uint, filter *ReferralFilter) ([]*Referral, int64, error)
	Update(ctx context.Context, id, userID uint, req *UpdateReferralRequest) (*Referral, error)
	Delete(ctx context.Context, id, userID uint) error
}

type repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) Repository { return &repository{db: db} }

func (r *repository) Create(ctx context.Context, ref *Referral) error {
	return r.db.WithContext(ctx).Create(ref).Error
}

func (r *repository) FindByID(ctx context.Context, id, userID uint) (*Referral, error) {
	var ref Referral
	if err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&ref).Error; err != nil {
		return nil, fmt.Errorf("referral not found")
	}
	return &ref, nil
}

func (r *repository) List(ctx context.Context, userID uint, filter *ReferralFilter) ([]*Referral, int64, error) {
	var refs []*Referral
	var total int64
	q := r.db.WithContext(ctx).Model(&Referral{}).Where("user_id = ?", userID)
	if filter.CompanyID > 0 { q = q.Where("company_id = ?", filter.CompanyID) }
	if filter.Status != "" { q = q.Where("status = ?", filter.Status) }
	q.Count(&total)
	page := filter.Page; if page < 1 { page = 1 }
	perPage := filter.PerPage; if perPage < 1 { perPage = 20 }
	q.Order("created_at DESC").Offset((page-1)*perPage).Limit(perPage).Find(&refs)
	return refs, total, nil
}

func (r *repository) Update(ctx context.Context, ref *Referral) error {
	return r.db.WithContext(ctx).Save(ref).Error
}

func (r *repository) Delete(ctx context.Context, id, userID uint) error {
	return r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&Referral{}).Error
}

type service struct{ repo Repository }

func NewService(repo Repository) Service { return &service{repo: repo} }

func (s *service) Create(ctx context.Context, userID uint, req *CreateReferralRequest) (*Referral, error) {
	ref := &Referral{
		UserID:                 userID,
		CompanyID:              req.CompanyID,
		ReferrerName:           req.ReferrerName,
		ReferrerDesignation:    req.ReferrerDesignation,
		ReferrerDepartment:     req.ReferrerDepartment,
		ReferrerOfficeLocation: req.ReferrerOfficeLocation,
		ReferrerProfileURL:     req.ReferrerProfileURL,
		Status:                 StatusNotContacted,
		Notes:                  req.Notes,
	}
	return ref, s.repo.Create(ctx, ref)
}

func (s *service) Get(ctx context.Context, id, userID uint) (*Referral, error) {
	return s.repo.FindByID(ctx, id, userID)
}

func (s *service) List(ctx context.Context, userID uint, filter *ReferralFilter) ([]*Referral, int64, error) {
	return s.repo.List(ctx, userID, filter)
}

func (s *service) Update(ctx context.Context, id, userID uint, req *UpdateReferralRequest) (*Referral, error) {
	ref, err := s.repo.FindByID(ctx, id, userID)
	if err != nil { return nil, err }
	if req.Status != "" {
		ref.Status = ReferralStatus(req.Status)
		now := time.Now()
		switch ref.Status {
		case StatusContacted: ref.ContactedAt = &now
		case StatusReferralReceived: ref.ReferralReceivedAt = &now
		}
	}
	if req.Notes != "" { ref.Notes = req.Notes }
	if req.FollowUpAt != nil { ref.FollowUpAt = req.FollowUpAt }
	return ref, s.repo.Update(ctx, ref)
}

func (s *service) Delete(ctx context.Context, id, userID uint) error {
	return s.repo.Delete(ctx, id, userID)
}

type Controller struct{ svc Service }

func NewController(svc Service) *Controller { return &Controller{svc: svc} }

func (c *Controller) Create(ctx *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok { response.Unauthorized(ctx, "unauthorized"); return }
	var req CreateReferralRequest
	if err := ctx.ShouldBindJSON(&req); err != nil { response.BadRequest(ctx, err.Error()); return }
	if err := validator.Validate(&req); err != nil { response.BadRequest(ctx, err.Error()); return }
	ref, err := c.svc.Create(ctx.Request.Context(), userID, &req)
	if err != nil { response.Error(ctx, err); return }
	response.Created(ctx, "referral created", ref)
}

func (c *Controller) List(ctx *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok { response.Unauthorized(ctx, "unauthorized"); return }
	var filter ReferralFilter
	_ = ctx.ShouldBindQuery(&filter)
	refs, total, _ := c.svc.List(ctx.Request.Context(), userID, &filter)
	response.Paginated(ctx, "referrals retrieved", refs, filter.Page, filter.PerPage, total)
}

func (c *Controller) Get(ctx *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok { response.Unauthorized(ctx, "unauthorized"); return }
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)
	ref, err := c.svc.Get(ctx.Request.Context(), uint(id), userID)
	if err != nil { response.NotFound(ctx, "referral not found"); return }
	response.Success(ctx, "referral retrieved", ref)
}

func (c *Controller) Update(ctx *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok { response.Unauthorized(ctx, "unauthorized"); return }
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)
	var req UpdateReferralRequest
	if err := ctx.ShouldBindJSON(&req); err != nil { response.BadRequest(ctx, err.Error()); return }
	ref, err := c.svc.Update(ctx.Request.Context(), uint(id), userID, &req)
	if err != nil { response.Error(ctx, err); return }
	response.Success(ctx, "referral updated", ref)
}

func (c *Controller) Delete(ctx *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok { response.Unauthorized(ctx, "unauthorized"); return }
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err := c.svc.Delete(ctx.Request.Context(), uint(id), userID); err != nil {
		response.Error(ctx, err); return
	}
	response.Success(ctx, "referral deleted", nil)
}

func RegisterRoutes(r *gin.RouterGroup, ctrl *Controller, cfg *config.Config) {
	refs := r.Group("/referrals")
	refs.Use(middleware.JWTAuth(&cfg.JWT))
	{
		refs.POST("", ctrl.Create)
		refs.GET("", ctrl.List)
		refs.GET("/:id", ctrl.Get)
		refs.PUT("/:id", ctrl.Update)
		refs.DELETE("/:id", ctrl.Delete)
	}
}
