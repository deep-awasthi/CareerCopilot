package user

import (
	"context"

	"github.com/deepawasthi/careercopilot/pkg/errors"
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

func (r *repository) Create(ctx context.Context, profile *Profile) error {
	return r.db.WithContext(ctx).Create(profile).Error
}

func (r *repository) FindByUserID(ctx context.Context, userID uint) (*Profile, error) {
	var profile Profile
	result := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&profile)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.NotFound("profile not found")
		}
		return nil, result.Error
	}
	return &profile, nil
}

func (r *repository) Update(ctx context.Context, profile *Profile) error {
	return r.db.WithContext(ctx).Save(profile).Error
}

func (r *repository) Upsert(ctx context.Context, userID uint, req *UpdateProfileRequest) (*Profile, error) {
	profile := Profile{
		UserID: userID,
	}

	r.db.WithContext(ctx).Where("user_id = ?", userID).First(&profile)

	if req.Name != "" {
		profile.Name = req.Name
	}
	if req.Phone != "" {
		profile.Phone = req.Phone
	}
	if req.ExperienceYears > 0 {
		profile.ExperienceYears = req.ExperienceYears
	}
	if req.CurrentCompany != "" {
		profile.CurrentCompany = req.CurrentCompany
	}
	if req.CurrentCTC > 0 {
		profile.CurrentCTC = req.CurrentCTC
	}
	if req.ExpectedCTC > 0 {
		profile.ExpectedCTC = req.ExpectedCTC
	}
	if req.NoticePeriodDays > 0 {
		profile.NoticePeriodDays = req.NoticePeriodDays
	}
	if len(req.PreferredLocations) > 0 {
		profile.PreferredLocations = pq.StringArray(req.PreferredLocations)
	}
	if len(req.PreferredRoles) > 0 {
		profile.PreferredRoles = pq.StringArray(req.PreferredRoles)
	}
	if len(req.PreferredSkills) > 0 {
		profile.PreferredSkills = pq.StringArray(req.PreferredSkills)
	}
	if req.Bio != "" {
		profile.Bio = req.Bio
	}
	if req.LinkedinURL != "" {
		profile.LinkedinURL = req.LinkedinURL
	}
	if req.GithubURL != "" {
		profile.GithubURL = req.GithubURL
	}
	if req.PortfolioURL != "" {
		profile.PortfolioURL = req.PortfolioURL
	}
	if req.IsOpenToWork != nil {
		profile.IsOpenToWork = *req.IsOpenToWork
	}
	if req.PreferredWorkType != "" {
		profile.PreferredWorkType = req.PreferredWorkType
	}

	result := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "phone", "experience_years", "current_company", "current_ctc", "expected_ctc", "notice_period_days", "preferred_locations", "preferred_roles", "preferred_skills", "bio", "linkedin_url", "github_url", "portfolio_url", "is_open_to_work", "preferred_work_type", "updated_at"}),
	}).Create(&profile)

	if result.Error != nil {
		return nil, result.Error
	}
	return &profile, nil
}
