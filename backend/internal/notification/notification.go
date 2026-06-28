package notification

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/deepawasthi/careercopilot/pkg/config"
	"github.com/deepawasthi/careercopilot/pkg/middleware"
	"github.com/deepawasthi/careercopilot/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/datatypes"
)

type NotificationType string
type NotificationChannel string

const (
	TypeNewJob          NotificationType = "new_job"
	TypeKeywordMatch    NotificationType = "keyword_match"
	TypeCompanyOpening  NotificationType = "company_opening"
	TypeReferralUpdate  NotificationType = "referral_update"
	TypeInterviewReminder NotificationType = "interview_reminder"
	TypeDailyDigest     NotificationType = "daily_digest"
)

const (
	ChannelEmail   NotificationChannel = "email"
	ChannelBrowser NotificationChannel = "browser"
)

type Notification struct {
	ID        uint                `gorm:"primarykey" json:"id"`
	UserID    uint                `gorm:"not null;index" json:"user_id"`
	Type      NotificationType    `gorm:"type:notification_type;not null" json:"type"`
	Channel   NotificationChannel `gorm:"type:notification_channel;not null" json:"channel"`
	Title     string              `gorm:"not null" json:"title"`
	Body      string              `gorm:"type:text" json:"body"`
	Metadata  datatypes.JSON      `json:"metadata"`
	IsRead    bool                `gorm:"default:false" json:"is_read"`
	SentAt    *time.Time          `json:"sent_at"`
	CreatedAt time.Time           `json:"created_at"`
}

func (Notification) TableName() string { return "notifications" }

type CreateNotificationRequest struct {
	UserID  uint                `json:"user_id"`
	Type    NotificationType    `json:"type"`
	Channel NotificationChannel `json:"channel"`
	Title   string              `json:"title"`
	Body    string              `json:"body"`
}

type NotificationFilter struct {
	IsRead  *bool  `form:"is_read"`
	Type    string `form:"type"`
	Page    int    `form:"page,default=1"`
	PerPage int    `form:"per_page,default=20"`
}

type Repository interface {
	Create(ctx context.Context, n *Notification) error
	FindByID(ctx context.Context, id, userID uint) (*Notification, error)
	List(ctx context.Context, userID uint, filter *NotificationFilter) ([]*Notification, int64, error)
	MarkRead(ctx context.Context, id, userID uint) error
	MarkAllRead(ctx context.Context, userID uint) error
	UnreadCount(ctx context.Context, userID uint) (int64, error)
	Delete(ctx context.Context, id, userID uint) error
}

type Service interface {
	Create(ctx context.Context, req *CreateNotificationRequest) (*Notification, error)
	List(ctx context.Context, userID uint, filter *NotificationFilter) ([]*Notification, int64, error)
	MarkRead(ctx context.Context, id, userID uint) error
	MarkAllRead(ctx context.Context, userID uint) error
	UnreadCount(ctx context.Context, userID uint) (int64, error)
	Delete(ctx context.Context, id, userID uint) error
}

type repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) Repository { return &repository{db: db} }

func (r *repository) Create(ctx context.Context, n *Notification) error {
	now := time.Now()
	n.SentAt = &now
	return r.db.WithContext(ctx).Create(n).Error
}

func (r *repository) FindByID(ctx context.Context, id, userID uint) (*Notification, error) {
	var n Notification
	if err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&n).Error; err != nil {
		return nil, fmt.Errorf("notification not found")
	}
	return &n, nil
}

func (r *repository) List(ctx context.Context, userID uint, filter *NotificationFilter) ([]*Notification, int64, error) {
	var ns []*Notification
	var total int64
	q := r.db.WithContext(ctx).Model(&Notification{}).Where("user_id = ?", userID)
	if filter.IsRead != nil { q = q.Where("is_read = ?", *filter.IsRead) }
	if filter.Type != "" { q = q.Where("type = ?", filter.Type) }
	q.Count(&total)
	page := filter.Page; if page < 1 { page = 1 }
	perPage := filter.PerPage; if perPage < 1 { perPage = 20 }
	q.Order("created_at DESC").Offset((page-1)*perPage).Limit(perPage).Find(&ns)
	return ns, total, nil
}

func (r *repository) MarkRead(ctx context.Context, id, userID uint) error {
	return r.db.WithContext(ctx).Model(&Notification{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("is_read", true).Error
}

func (r *repository) MarkAllRead(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).Model(&Notification{}).
		Where("user_id = ? AND is_read = false", userID).
		Update("is_read", true).Error
}

func (r *repository) UnreadCount(ctx context.Context, userID uint) (int64, error) {
	var count int64
	r.db.WithContext(ctx).Model(&Notification{}).Where("user_id = ? AND is_read = false", userID).Count(&count)
	return count, nil
}

func (r *repository) Delete(ctx context.Context, id, userID uint) error {
	return r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&Notification{}).Error
}

type service struct{ repo Repository }

func NewService(repo Repository) Service { return &service{repo: repo} }

func (s *service) Create(ctx context.Context, req *CreateNotificationRequest) (*Notification, error) {
	n := &Notification{
		UserID:  req.UserID,
		Type:    req.Type,
		Channel: req.Channel,
		Title:   req.Title,
		Body:    req.Body,
	}
	return n, s.repo.Create(ctx, n)
}

func (s *service) List(ctx context.Context, userID uint, filter *NotificationFilter) ([]*Notification, int64, error) {
	return s.repo.List(ctx, userID, filter)
}

func (s *service) MarkRead(ctx context.Context, id, userID uint) error {
	return s.repo.MarkRead(ctx, id, userID)
}

func (s *service) MarkAllRead(ctx context.Context, userID uint) error {
	return s.repo.MarkAllRead(ctx, userID)
}

func (s *service) UnreadCount(ctx context.Context, userID uint) (int64, error) {
	return s.repo.UnreadCount(ctx, userID)
}

func (s *service) Delete(ctx context.Context, id, userID uint) error {
	return s.repo.Delete(ctx, id, userID)
}

type Controller struct{ svc Service }

func NewController(svc Service) *Controller { return &Controller{svc: svc} }

func (c *Controller) List(ctx *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok { response.Unauthorized(ctx, "unauthorized"); return }
	var filter NotificationFilter
	_ = ctx.ShouldBindQuery(&filter)
	ns, total, _ := c.svc.List(ctx.Request.Context(), userID, &filter)
	response.Paginated(ctx, "notifications retrieved", ns, filter.Page, filter.PerPage, total)
}

func (c *Controller) MarkRead(ctx *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok { response.Unauthorized(ctx, "unauthorized"); return }
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err := c.svc.MarkRead(ctx.Request.Context(), uint(id), userID); err != nil {
		response.Error(ctx, err); return
	}
	response.Success(ctx, "notification marked as read", nil)
}

func (c *Controller) MarkAllRead(ctx *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok { response.Unauthorized(ctx, "unauthorized"); return }
	_ = c.svc.MarkAllRead(ctx.Request.Context(), userID)
	response.Success(ctx, "all notifications marked as read", nil)
}

func (c *Controller) UnreadCount(ctx *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok { response.Unauthorized(ctx, "unauthorized"); return }
	count, _ := c.svc.UnreadCount(ctx.Request.Context(), userID)
	response.Success(ctx, "unread count", gin.H{"count": count})
}

func (c *Controller) Delete(ctx *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok { response.Unauthorized(ctx, "unauthorized"); return }
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)
	_ = c.svc.Delete(ctx.Request.Context(), uint(id), userID)
	response.Success(ctx, "notification deleted", nil)
}

func RegisterRoutes(r *gin.RouterGroup, ctrl *Controller, cfg *config.Config) {
	ns := r.Group("/notifications")
	ns.Use(middleware.JWTAuth(&cfg.JWT))
	{
		ns.GET("", ctrl.List)
		ns.GET("/unread-count", ctrl.UnreadCount)
		ns.PUT("/read-all", ctrl.MarkAllRead)
		ns.PUT("/:id/read", ctrl.MarkRead)
		ns.DELETE("/:id", ctrl.Delete)
	}
}
