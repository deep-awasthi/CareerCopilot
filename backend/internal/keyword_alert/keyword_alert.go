package keyword_alert

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/deepawasthi/careercopilot/pkg/config"
	"github.com/deepawasthi/careercopilot/pkg/middleware"
	"github.com/deepawasthi/careercopilot/pkg/response"
	"github.com/deepawasthi/careercopilot/pkg/validator"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ---- Entity ----

type KeywordAlert struct {
	ID              uint           `gorm:"primarykey" json:"id"`
	UserID          uint           `gorm:"not null;index" json:"user_id"`
	Keyword         string         `gorm:"not null" json:"keyword"`
	IsActive        bool           `gorm:"default:true" json:"is_active"`
	EmailNotify     bool           `gorm:"default:true" json:"email_notify"`
	BrowserNotify   bool           `gorm:"default:true" json:"browser_notify"`
	MatchCount      int            `gorm:"default:0" json:"match_count"`
	LastMatchedAt   *time.Time     `json:"last_matched_at"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

func (KeywordAlert) TableName() string { return "keyword_alerts" }

// ---- DTOs ----

type CreateAlertRequest struct {
	Keyword       string `json:"keyword" validate:"required,min=2,max=100"`
	EmailNotify   bool   `json:"email_notify"`
	BrowserNotify bool   `json:"browser_notify"`
}

type UpdateAlertRequest struct {
	IsActive      *bool `json:"is_active"`
	EmailNotify   *bool `json:"email_notify"`
	BrowserNotify *bool `json:"browser_notify"`
}

// ---- Repository & Service ----

type Repository interface {
	Create(ctx context.Context, alert *KeywordAlert) error
	FindByID(ctx context.Context, id, userID uint) (*KeywordAlert, error)
	ListByUser(ctx context.Context, userID uint) ([]*KeywordAlert, error)
	ListActive(ctx context.Context) ([]*KeywordAlert, error)
	Update(ctx context.Context, alert *KeywordAlert) error
	Delete(ctx context.Context, id, userID uint) error
	IncrementMatchCount(ctx context.Context, id uint) error
}

type Service interface {
	Create(ctx context.Context, userID uint, req *CreateAlertRequest) (*KeywordAlert, error)
	List(ctx context.Context, userID uint) ([]*KeywordAlert, error)
	Update(ctx context.Context, id, userID uint, req *UpdateAlertRequest) (*KeywordAlert, error)
	Delete(ctx context.Context, id, userID uint) error
	ProcessJobText(ctx context.Context, jobText string) ([]uint, error) // returns matching user IDs
}

type repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) Repository { return &repository{db: db} }

func (r *repository) Create(ctx context.Context, alert *KeywordAlert) error {
	return r.db.WithContext(ctx).Create(alert).Error
}

func (r *repository) FindByID(ctx context.Context, id, userID uint) (*KeywordAlert, error) {
	var a KeywordAlert
	if err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&a).Error; err != nil {
		return nil, fmt.Errorf("alert not found")
	}
	return &a, nil
}

func (r *repository) ListByUser(ctx context.Context, userID uint) ([]*KeywordAlert, error) {
	var alerts []*KeywordAlert
	r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Find(&alerts)
	return alerts, nil
}

func (r *repository) ListActive(ctx context.Context) ([]*KeywordAlert, error) {
	var alerts []*KeywordAlert
	r.db.WithContext(ctx).Where("is_active = true").Find(&alerts)
	return alerts, nil
}

func (r *repository) Update(ctx context.Context, alert *KeywordAlert) error {
	return r.db.WithContext(ctx).Save(alert).Error
}

func (r *repository) Delete(ctx context.Context, id, userID uint) error {
	return r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&KeywordAlert{}).Error
}

func (r *repository) IncrementMatchCount(ctx context.Context, id uint) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&KeywordAlert{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"match_count":     gorm.Expr("match_count + 1"),
			"last_matched_at": now,
		}).Error
}

type service struct{ repo Repository }

func NewService(repo Repository) Service { return &service{repo: repo} }

func (s *service) Create(ctx context.Context, userID uint, req *CreateAlertRequest) (*KeywordAlert, error) {
	alert := &KeywordAlert{
		UserID:        userID,
		Keyword:       strings.TrimSpace(req.Keyword),
		IsActive:      true,
		EmailNotify:   req.EmailNotify,
		BrowserNotify: req.BrowserNotify,
	}
	return alert, s.repo.Create(ctx, alert)
}

func (s *service) List(ctx context.Context, userID uint) ([]*KeywordAlert, error) {
	return s.repo.ListByUser(ctx, userID)
}

func (s *service) Update(ctx context.Context, id, userID uint, req *UpdateAlertRequest) (*KeywordAlert, error) {
	alert, err := s.repo.FindByID(ctx, id, userID)
	if err != nil { return nil, err }
	if req.IsActive != nil { alert.IsActive = *req.IsActive }
	if req.EmailNotify != nil { alert.EmailNotify = *req.EmailNotify }
	if req.BrowserNotify != nil { alert.BrowserNotify = *req.BrowserNotify }
	return alert, s.repo.Update(ctx, alert)
}

func (s *service) Delete(ctx context.Context, id, userID uint) error {
	return s.repo.Delete(ctx, id, userID)
}

// ProcessJobText scans a job's text against ALL active alerts
// Returns slice of (alertID, userID) pairs that matched
func (s *service) ProcessJobText(ctx context.Context, jobText string) ([]uint, error) {
	alerts, err := s.repo.ListActive(ctx)
	if err != nil { return nil, err }

	lower := strings.ToLower(jobText)
	var matchedUserIDs []uint
	seen := make(map[uint]bool)

	for _, alert := range alerts {
		if strings.Contains(lower, strings.ToLower(alert.Keyword)) {
			if !seen[alert.UserID] {
				seen[alert.UserID] = true
				matchedUserIDs = append(matchedUserIDs, alert.UserID)
			}
			_ = s.repo.IncrementMatchCount(ctx, alert.ID)
		}
	}
	return matchedUserIDs, nil
}

// ---- Controller ----

type Controller struct{ svc Service }

func NewController(svc Service) *Controller { return &Controller{svc: svc} }

func (c *Controller) Create(ctx *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok { response.Unauthorized(ctx, "unauthorized"); return }
	var req CreateAlertRequest
	if err := ctx.ShouldBindJSON(&req); err != nil { response.BadRequest(ctx, err.Error()); return }
	if err := validator.Validate(&req); err != nil { response.BadRequest(ctx, err.Error()); return }
	alert, err := c.svc.Create(ctx.Request.Context(), userID, &req)
	if err != nil { response.Error(ctx, err); return }
	response.Created(ctx, "keyword alert created", alert)
}

func (c *Controller) List(ctx *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok { response.Unauthorized(ctx, "unauthorized"); return }
	alerts, err := c.svc.List(ctx.Request.Context(), userID)
	if err != nil { response.Error(ctx, err); return }
	response.Success(ctx, "alerts retrieved", alerts)
}

func (c *Controller) Update(ctx *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok { response.Unauthorized(ctx, "unauthorized"); return }
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)
	var req UpdateAlertRequest
	if err := ctx.ShouldBindJSON(&req); err != nil { response.BadRequest(ctx, err.Error()); return }
	alert, err := c.svc.Update(ctx.Request.Context(), uint(id), userID, &req)
	if err != nil { response.Error(ctx, err); return }
	response.Success(ctx, "alert updated", alert)
}

func (c *Controller) Delete(ctx *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok { response.Unauthorized(ctx, "unauthorized"); return }
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err := c.svc.Delete(ctx.Request.Context(), uint(id), userID); err != nil {
		response.Error(ctx, err); return
	}
	response.Success(ctx, "alert deleted", nil)
}

func RegisterRoutes(r *gin.RouterGroup, ctrl *Controller, cfg *config.Config) {
	alerts := r.Group("/alerts")
	alerts.Use(middleware.JWTAuth(&cfg.JWT))
	{
		alerts.POST("", ctrl.Create)
		alerts.GET("", ctrl.List)
		alerts.PUT("/:id", ctrl.Update)
		alerts.DELETE("/:id", ctrl.Delete)
	}
}
