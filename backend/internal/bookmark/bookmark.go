package bookmark

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
)

type BookmarkType string

const (
	TypeJob      BookmarkType = "job"
	TypeCompany  BookmarkType = "company"
	TypeReferral BookmarkType = "referral"
)

type Bookmark struct {
	ID        uint         `gorm:"primarykey" json:"id"`
	UserID    uint         `gorm:"not null;index" json:"user_id"`
	Type      BookmarkType `gorm:"type:bookmark_type;not null" json:"type"`
	TargetID  uint         `gorm:"not null" json:"target_id"`
	Notes     string       `json:"notes"`
	CreatedAt time.Time    `json:"created_at"`
}

func (Bookmark) TableName() string { return "bookmarks" }

type Repository interface {
	Create(ctx context.Context, b *Bookmark) error
	Delete(ctx context.Context, userID uint, bType BookmarkType, targetID uint) error
	List(ctx context.Context, userID uint, bType string, page, perPage int) ([]*Bookmark, int64, error)
	Exists(ctx context.Context, userID uint, bType BookmarkType, targetID uint) bool
}

type Service interface {
	Bookmark(ctx context.Context, userID uint, bType, targetIDStr string) (*Bookmark, error)
	Unbookmark(ctx context.Context, userID uint, bType, targetIDStr string) error
	List(ctx context.Context, userID uint, bType string, page, perPage int) ([]*Bookmark, int64, error)
}

type repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) Repository { return &repository{db: db} }

func (r *repository) Create(ctx context.Context, b *Bookmark) error {
	return r.db.WithContext(ctx).
		Where(Bookmark{UserID: b.UserID, Type: b.Type, TargetID: b.TargetID}).
		FirstOrCreate(b).Error
}

func (r *repository) Delete(ctx context.Context, userID uint, bType BookmarkType, targetID uint) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND type = ? AND target_id = ?", userID, bType, targetID).
		Delete(&Bookmark{}).Error
}

func (r *repository) List(ctx context.Context, userID uint, bType string, page, perPage int) ([]*Bookmark, int64, error) {
	var bs []*Bookmark
	var total int64
	q := r.db.WithContext(ctx).Model(&Bookmark{}).Where("user_id = ?", userID)
	if bType != "" { q = q.Where("type = ?", bType) }
	q.Count(&total)
	if page < 1 { page = 1 }
	if perPage < 1 { perPage = 20 }
	q.Order("created_at DESC").Offset((page-1)*perPage).Limit(perPage).Find(&bs)
	return bs, total, nil
}

func (r *repository) Exists(ctx context.Context, userID uint, bType BookmarkType, targetID uint) bool {
	var count int64
	r.db.WithContext(ctx).Model(&Bookmark{}).
		Where("user_id = ? AND type = ? AND target_id = ?", userID, bType, targetID).
		Count(&count)
	return count > 0
}

type service struct{ repo Repository }

func NewService(repo Repository) Service { return &service{repo: repo} }

func (s *service) Bookmark(ctx context.Context, userID uint, bType, targetIDStr string) (*Bookmark, error) {
	targetID, err := strconv.ParseUint(targetIDStr, 10, 64)
	if err != nil { return nil, fmt.Errorf("invalid target_id") }
	b := &Bookmark{
		UserID:   userID,
		Type:     BookmarkType(bType),
		TargetID: uint(targetID),
	}
	return b, s.repo.Create(ctx, b)
}

func (s *service) Unbookmark(ctx context.Context, userID uint, bType, targetIDStr string) error {
	targetID, err := strconv.ParseUint(targetIDStr, 10, 64)
	if err != nil { return fmt.Errorf("invalid target_id") }
	return s.repo.Delete(ctx, userID, BookmarkType(bType), uint(targetID))
}

func (s *service) List(ctx context.Context, userID uint, bType string, page, perPage int) ([]*Bookmark, int64, error) {
	return s.repo.List(ctx, userID, bType, page, perPage)
}

type Controller struct{ svc Service }

func NewController(svc Service) *Controller { return &Controller{svc: svc} }

func (c *Controller) List(ctx *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok { response.Unauthorized(ctx, "unauthorized"); return }
	bType := ctx.Query("type")
	page := 1; perPage := 20
	bs, total, _ := c.svc.List(ctx.Request.Context(), userID, bType, page, perPage)
	response.Paginated(ctx, "bookmarks retrieved", bs, page, perPage, total)
}

func (c *Controller) Create(ctx *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok { response.Unauthorized(ctx, "unauthorized"); return }
	bType := ctx.Param("type")
	targetID := ctx.Param("id")
	b, err := c.svc.Bookmark(ctx.Request.Context(), userID, bType, targetID)
	if err != nil { response.Error(ctx, err); return }
	response.Created(ctx, "bookmarked", b)
}

func (c *Controller) Delete(ctx *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok { response.Unauthorized(ctx, "unauthorized"); return }
	bType := ctx.Param("type")
	targetID := ctx.Param("id")
	_ = c.svc.Unbookmark(ctx.Request.Context(), userID, bType, targetID)
	response.Success(ctx, "bookmark removed", nil)
}

func RegisterRoutes(r *gin.RouterGroup, ctrl *Controller, cfg *config.Config) {
	bm := r.Group("/bookmarks")
	bm.Use(middleware.JWTAuth(&cfg.JWT))
	{
		bm.GET("", ctrl.List)
		bm.POST("/:type/:id", ctrl.Create)
		bm.DELETE("/:type/:id", ctrl.Delete)
	}
}
