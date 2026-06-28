package application

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, app *Application) error
	FindByID(ctx context.Context, id, userID uint) (*Application, error)
	FindByUserIDAndJobID(ctx context.Context, userID, jobID uint) (*Application, error)
	List(ctx context.Context, userID uint, filter *ApplicationFilter) ([]*Application, int64, error)
	Update(ctx context.Context, app *Application) error
	Delete(ctx context.Context, id, userID uint) error
	CountByStatus(ctx context.Context, userID uint) ([]StatusCount, error)
	CountAppliedThisMonth(ctx context.Context, userID uint) (int64, error)
}

type Service interface {
	Create(ctx context.Context, userID uint, req *CreateApplicationRequest) (*ApplicationResponse, error)
	Get(ctx context.Context, id, userID uint) (*ApplicationResponse, error)
	List(ctx context.Context, userID uint, filter *ApplicationFilter) ([]*ApplicationResponse, int64, error)
	Update(ctx context.Context, id, userID uint, req *UpdateApplicationRequest) (*ApplicationResponse, error)
	Delete(ctx context.Context, id, userID uint) error
	GetStats(ctx context.Context, userID uint) ([]StatusCount, error)
}

type repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) Repository { return &repository{db: db} }

func (r *repository) Create(ctx context.Context, app *Application) error {
	return r.db.WithContext(ctx).Create(app).Error
}

func (r *repository) FindByID(ctx context.Context, id, userID uint) (*Application, error) {
	var app Application
	result := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&app)
	if result.Error != nil {
		return nil, fmt.Errorf("application not found")
	}
	return &app, nil
}

func (r *repository) FindByUserIDAndJobID(ctx context.Context, userID, jobID uint) (*Application, error) {
	var app Application
	result := r.db.WithContext(ctx).Where("user_id = ? AND job_id = ?", userID, jobID).First(&app)
	if result.Error != nil {
		return nil, result.Error
	}
	return &app, nil
}

func (r *repository) List(ctx context.Context, userID uint, filter *ApplicationFilter) ([]*Application, int64, error) {
	var apps []*Application
	var total int64
	q := r.db.WithContext(ctx).Model(&Application{}).Where("user_id = ?", userID)
	if filter.Status != "" {
		q = q.Where("status = ?", filter.Status)
	}
	q.Count(&total)
	page := filter.Page
	if page < 1 { page = 1 }
	perPage := filter.PerPage
	if perPage < 1 || perPage > 100 { perPage = 20 }
	offset := (page - 1) * perPage
	result := q.Order("created_at DESC").Offset(offset).Limit(perPage).Find(&apps)
	return apps, total, result.Error
}

func (r *repository) Update(ctx context.Context, app *Application) error {
	return r.db.WithContext(ctx).Save(app).Error
}

func (r *repository) Delete(ctx context.Context, id, userID uint) error {
	return r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&Application{}).Error
}

func (r *repository) CountByStatus(ctx context.Context, userID uint) ([]StatusCount, error) {
	var counts []StatusCount
	r.db.WithContext(ctx).Model(&Application{}).
		Select("status, COUNT(*) as count").
		Where("user_id = ?", userID).
		Group("status").
		Scan(&counts)
	return counts, nil
}

func (r *repository) CountAppliedThisMonth(ctx context.Context, userID uint) (int64, error) {
	var count int64
	r.db.WithContext(ctx).Model(&Application{}).
		Where("user_id = ? AND status != 'saved' AND DATE_TRUNC('month', created_at) = DATE_TRUNC('month', NOW())", userID).
		Count(&count)
	return count, nil
}

type service struct{ repo Repository }

func NewService(repo Repository) Service { return &service{repo: repo} }

func (s *service) Create(ctx context.Context, userID uint, req *CreateApplicationRequest) (*ApplicationResponse, error) {
	status := ApplicationStatus(req.Status)
	if status == "" { status = StatusSaved }

	app := &Application{
		UserID:        userID,
		JobID:         req.JobID,
		Status:        status,
		Notes:         req.Notes,
		FollowUpDate:  req.FollowUpDate,
		SalaryOffered: req.SalaryOffered,
		ReferralUsed:  req.ReferralUsed,
		ReferralID:    req.ReferralID,
	}
	if strings.EqualFold(string(status), "applied") {
		now := time.Now()
		app.AppliedAt = &now
	}
	if err := s.repo.Create(ctx, app); err != nil {
		return nil, err
	}
	return ToResponse(app), nil
}

func (s *service) Get(ctx context.Context, id, userID uint) (*ApplicationResponse, error) {
	app, err := s.repo.FindByID(ctx, id, userID)
	if err != nil { return nil, err }
	return ToResponse(app), nil
}

func (s *service) List(ctx context.Context, userID uint, filter *ApplicationFilter) ([]*ApplicationResponse, int64, error) {
	apps, total, err := s.repo.List(ctx, userID, filter)
	if err != nil { return nil, 0, err }
	var result []*ApplicationResponse
	for _, a := range apps { result = append(result, ToResponse(a)) }
	return result, total, nil
}

func (s *service) Update(ctx context.Context, id, userID uint, req *UpdateApplicationRequest) (*ApplicationResponse, error) {
	app, err := s.repo.FindByID(ctx, id, userID)
	if err != nil { return nil, err }
	if req.Status != "" { app.Status = ApplicationStatus(req.Status) }
	if req.Notes != "" { app.Notes = req.Notes }
	if req.FollowUpDate != nil { app.FollowUpDate = req.FollowUpDate }
	if req.SalaryOffered > 0 { app.SalaryOffered = req.SalaryOffered }
	if req.ReferralUsed != nil { app.ReferralUsed = *req.ReferralUsed }
	if req.AppliedAt != nil { app.AppliedAt = req.AppliedAt }
	if err := s.repo.Update(ctx, app); err != nil { return nil, err }
	return ToResponse(app), nil
}

func (s *service) Delete(ctx context.Context, id, userID uint) error {
	return s.repo.Delete(ctx, id, userID)
}

func (s *service) GetStats(ctx context.Context, userID uint) ([]StatusCount, error) {
	return s.repo.CountByStatus(ctx, userID)
}
