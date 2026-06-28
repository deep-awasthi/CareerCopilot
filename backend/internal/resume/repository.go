package resume

import (
	"context"
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Upsert(ctx context.Context, resume *Resume) error {
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"raw_text", "parsed_skills", "companies", "projects",
			"education", "experience", "certifications", "parsed_at", "updated_at",
		}),
	}).Create(resume).Error
}

func (r *repository) FindByUserID(ctx context.Context, userID uint) (*Resume, error) {
	var resume Resume
	result := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&resume)
	if result.Error != nil {
		return nil, result.Error
	}
	return &resume, nil
}

type service struct {
	repo   Repository
	parser *Parser
}

func NewService(repo Repository) Service {
	return &service{repo: repo, parser: NewParser()}
}

func (s *service) Submit(ctx context.Context, userID uint, req *SubmitResumeRequest) (*ResumeResponse, error) {
	result := s.parser.Parse(req.RawText)
	now := time.Now()

	resume := &Resume{
		UserID:         userID,
		RawText:        req.RawText,
		ParsedSkills:   pq.StringArray(result.Skills),
		Companies:      pq.StringArray(result.Companies),
		Projects:       pq.StringArray(result.Projects),
		Education:      JSONBArray(result.Education),
		Experience:     JSONBArray(result.Experience),
		Certifications: pq.StringArray(result.Certifications),
		ParsedAt:       &now,
	}

	if err := s.repo.Upsert(ctx, resume); err != nil {
		return nil, err
	}

	return ToResumeResponse(resume), nil
}

func (s *service) Get(ctx context.Context, userID uint) (*ResumeResponse, error) {
	resume, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return ToResumeResponse(resume), nil
}
