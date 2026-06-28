package company

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

// ---- Entity ----

type Company struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	Name          string         `gorm:"not null" json:"name"`
	Slug          string         `gorm:"uniqueIndex;not null" json:"slug"`
	Domain        string         `json:"domain"`
	CareerPageURL string         `json:"career_page_url"`
	LogoURL       string         `json:"logo_url"`
	Industry      string         `json:"industry"`
	Size          string         `json:"size"`
	Headquarters  string         `json:"headquarters"`
	Description   string         `gorm:"type:text" json:"description"`
	LinkedinURL   string         `json:"linkedin_url"`
	GlassdoorURL  string         `json:"glassdoor_url"`
	FoundedYear   int            `json:"founded_year"`
	IsActive      bool           `gorm:"default:true" json:"is_active"`
	LastScrapedAt *time.Time     `json:"last_scraped_at"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Company) TableName() string { return "companies" }

type Watchlist struct {
	ID             uint      `gorm:"primarykey" json:"id"`
	UserID         uint      `gorm:"not null;index" json:"user_id"`
	CompanyID      uint      `gorm:"not null;index" json:"company_id"`
	NotifyNewJobs  bool      `gorm:"default:true" json:"notify_new_jobs"`
	LastNotifiedAt *time.Time `json:"last_notified_at"`
	Company        *Company  `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

func (Watchlist) TableName() string { return "company_watchlists" }

// ---- DTOs ----

type CreateCompanyRequest struct {
	Name          string `json:"name" validate:"required"`
	Domain        string `json:"domain"`
	CareerPageURL string `json:"career_page_url"`
	Industry      string `json:"industry"`
	Size          string `json:"size"`
	Headquarters  string `json:"headquarters"`
	Description   string `json:"description"`
}

type WatchlistFilter struct {
	Page    int `form:"page,default=1"`
	PerPage int `form:"per_page,default=20"`
}

// ---- Repository & Service ----

type Repository interface {
	FindOrCreate(ctx context.Context, company *Company) error
	FindByID(ctx context.Context, id uint) (*Company, error)
	FindBySlug(ctx context.Context, slug string) (*Company, error)
	Search(ctx context.Context, q string, page, perPage int) ([]*Company, int64, error)
	AddToWatchlist(ctx context.Context, userID, companyID uint) error
	RemoveFromWatchlist(ctx context.Context, userID, companyID uint) error
	GetWatchlist(ctx context.Context, userID uint, filter *WatchlistFilter) ([]*Watchlist, int64, error)
	IsWatching(ctx context.Context, userID, companyID uint) bool
}

type Service interface {
	GetCompany(ctx context.Context, id uint) (*Company, error)
	SearchCompanies(ctx context.Context, q string, page, perPage int) ([]*Company, int64, error)
	WatchCompany(ctx context.Context, userID, companyID uint) error
	UnwatchCompany(ctx context.Context, userID, companyID uint) error
	GetWatchlist(ctx context.Context, userID uint, filter *WatchlistFilter) ([]*Watchlist, int64, error)
	CreateCompany(ctx context.Context, req *CreateCompanyRequest) (*Company, error)
}

type repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) Repository { return &repository{db: db} }

func (r *repository) FindOrCreate(ctx context.Context, company *Company) error {
	return r.db.WithContext(ctx).
		Where(Company{Slug: company.Slug}).
		Assign(*company).
		FirstOrCreate(company).Error
}

func (r *repository) FindByID(ctx context.Context, id uint) (*Company, error) {
	var c Company
	if err := r.db.WithContext(ctx).First(&c, id).Error; err != nil {
		return nil, fmt.Errorf("company not found")
	}
	return &c, nil
}

func (r *repository) FindBySlug(ctx context.Context, slug string) (*Company, error) {
	var c Company
	if err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&c).Error; err != nil {
		return nil, fmt.Errorf("company not found")
	}
	return &c, nil
}

func (r *repository) Search(ctx context.Context, q string, page, perPage int) ([]*Company, int64, error) {
	var companies []*Company
	var total int64
	query := r.db.WithContext(ctx).Model(&Company{}).Where("is_active = true")
	if q != "" {
		query = query.Where("LOWER(name) LIKE ? OR LOWER(domain) LIKE ?", "%"+q+"%", "%"+q+"%")
	}
	query.Count(&total)
	offset := (page - 1) * perPage
	query.Order("name ASC").Offset(offset).Limit(perPage).Find(&companies)
	return companies, total, nil
}

func (r *repository) AddToWatchlist(ctx context.Context, userID, companyID uint) error {
	wl := Watchlist{UserID: userID, CompanyID: companyID, NotifyNewJobs: true}
	return r.db.WithContext(ctx).
		Where(Watchlist{UserID: userID, CompanyID: companyID}).
		FirstOrCreate(&wl).Error
}

func (r *repository) RemoveFromWatchlist(ctx context.Context, userID, companyID uint) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND company_id = ?", userID, companyID).
		Delete(&Watchlist{}).Error
}

func (r *repository) GetWatchlist(ctx context.Context, userID uint, filter *WatchlistFilter) ([]*Watchlist, int64, error) {
	var wls []*Watchlist
	var total int64
	q := r.db.WithContext(ctx).Model(&Watchlist{}).Where("user_id = ?", userID)
	q.Count(&total)
	page := filter.Page
	if page < 1 { page = 1 }
	perPage := filter.PerPage
	if perPage < 1 { perPage = 20 }
	q.Preload("Company").Offset((page-1)*perPage).Limit(perPage).Find(&wls)
	return wls, total, nil
}

func (r *repository) IsWatching(ctx context.Context, userID, companyID uint) bool {
	var count int64
	r.db.WithContext(ctx).Model(&Watchlist{}).
		Where("user_id = ? AND company_id = ?", userID, companyID).Count(&count)
	return count > 0
}

type service struct{ repo Repository }

func NewService(repo Repository) Service { return &service{repo: repo} }

func slugify(name string) string {
	result := ""
	for _, r := range name {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' {
			result += string(r)
		} else if r >= 'A' && r <= 'Z' {
			result += string(r + 32)
		} else if r == ' ' || r == '-' || r == '_' {
			result += "-"
		}
	}
	return result
}

func (s *service) CreateCompany(ctx context.Context, req *CreateCompanyRequest) (*Company, error) {
	company := &Company{
		Name:          req.Name,
		Slug:          slugify(req.Name),
		Domain:        req.Domain,
		CareerPageURL: req.CareerPageURL,
		Industry:      req.Industry,
		Size:          req.Size,
		Headquarters:  req.Headquarters,
		Description:   req.Description,
		IsActive:      true,
	}
	return company, s.repo.FindOrCreate(ctx, company)
}

func (s *service) GetCompany(ctx context.Context, id uint) (*Company, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *service) SearchCompanies(ctx context.Context, q string, page, perPage int) ([]*Company, int64, error) {
	return s.repo.Search(ctx, q, page, perPage)
}

func (s *service) WatchCompany(ctx context.Context, userID, companyID uint) error {
	return s.repo.AddToWatchlist(ctx, userID, companyID)
}

func (s *service) UnwatchCompany(ctx context.Context, userID, companyID uint) error {
	return s.repo.RemoveFromWatchlist(ctx, userID, companyID)
}

func (s *service) GetWatchlist(ctx context.Context, userID uint, filter *WatchlistFilter) ([]*Watchlist, int64, error) {
	return s.repo.GetWatchlist(ctx, userID, filter)
}

// ---- Controller ----

type Controller struct{ svc Service }

func NewController(svc Service) *Controller { return &Controller{svc: svc} }

func (c *Controller) Search(ctx *gin.Context) {
	q := ctx.Query("q")
	page := 1
	perPage := 20
	companies, total, _ := c.svc.SearchCompanies(ctx.Request.Context(), q, page, perPage)
	response.Paginated(ctx, "companies retrieved", companies, page, perPage, total)
}

func (c *Controller) Get(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)
	company, err := c.svc.GetCompany(ctx.Request.Context(), uint(id))
	if err != nil { response.NotFound(ctx, "company not found"); return }
	response.Success(ctx, "company retrieved", company)
}

func (c *Controller) Create(ctx *gin.Context) {
	var req CreateCompanyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil { response.BadRequest(ctx, err.Error()); return }
	company, err := c.svc.CreateCompany(ctx.Request.Context(), &req)
	if err != nil { response.Error(ctx, err); return }
	response.Created(ctx, "company created", company)
}

func (c *Controller) GetWatchlist(ctx *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok { response.Unauthorized(ctx, "unauthorized"); return }
	var filter WatchlistFilter
	_ = ctx.ShouldBindQuery(&filter)
	wls, total, _ := c.svc.GetWatchlist(ctx.Request.Context(), userID, &filter)
	response.Paginated(ctx, "watchlist retrieved", wls, filter.Page, filter.PerPage, total)
}

func (c *Controller) Watch(ctx *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok { response.Unauthorized(ctx, "unauthorized"); return }
	companyID, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err := c.svc.WatchCompany(ctx.Request.Context(), userID, uint(companyID)); err != nil {
		response.Error(ctx, err); return
	}
	response.Success(ctx, "company added to watchlist", nil)
}

func (c *Controller) Unwatch(ctx *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok { response.Unauthorized(ctx, "unauthorized"); return }
	companyID, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err := c.svc.UnwatchCompany(ctx.Request.Context(), userID, uint(companyID)); err != nil {
		response.Error(ctx, err); return
	}
	response.Success(ctx, "company removed from watchlist", nil)
}

func RegisterRoutes(r *gin.RouterGroup, ctrl *Controller, cfg *config.Config) {
	companies := r.Group("/companies")
	companies.Use(middleware.JWTAuth(&cfg.JWT))
	{
		companies.GET("", ctrl.Search)
		companies.GET("/:id", ctrl.Get)
		companies.POST("", ctrl.Create)
		companies.POST("/:id/watch", ctrl.Watch)
		companies.DELETE("/:id/watch", ctrl.Unwatch)
	}

	watchlist := r.Group("/watchlists")
	watchlist.Use(middleware.JWTAuth(&cfg.JWT))
	{
		watchlist.GET("", ctrl.GetWatchlist)
	}
}
