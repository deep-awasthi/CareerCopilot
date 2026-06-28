package job

import (
	"context"
	"crypto/sha256"
	"fmt"
	"strings"
	"time"

	"github.com/deepawasthi/careercopilot/internal/user"
	"github.com/lib/pq"
)

type service struct {
	repo     Repository
	userRepo user.Repository
	matcher  *Matcher
	dedup    *DeduplicationEngine
}

func NewService(repo Repository, userRepo user.Repository) Service {
	return &service{
		repo:     repo,
		userRepo: userRepo,
		matcher:  NewMatcher(),
		dedup:    NewDeduplicationEngine(),
	}
}

func (s *service) GetJob(ctx context.Context, id uint) (*JobResponse, error) {
	job, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return ToJobResponse(job), nil
}

func (s *service) ListJobs(ctx context.Context, filter *JobFilter, userID uint) ([]*JobResponse, int64, error) {
	jobs, total, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	matchCfg := s.getMatcherConfig(ctx, userID)

	var result []*JobResponse
	for _, j := range jobs {
		resp := ToJobResponse(j)
		if matchCfg != nil {
			breakdown := s.matcher.Score(j, matchCfg)
			resp.MatchScore = breakdown.Percentage
			resp.MatchBreakdown = breakdown
		}
		result = append(result, resp)
	}
	return result, total, nil
}

func (s *service) UpsertFromProvider(ctx context.Context, normalized *NormalizedJob) error {
	hash := s.dedup.ComputeHash(normalized.Company, normalized.Title, normalized.Location)
	normalized.DedupHash = hash

	postedAt := parseDate(normalized.PostedAt)

	job := &Job{
		ExternalID:     normalized.ExternalID,
		Title:          normalized.Title,
		Description:    normalized.Description,
		Location:       normalized.Location,
		Locations:      pq.StringArray(normalized.Locations),
		IsRemote:       normalized.IsRemote,
		IsHybrid:       normalized.IsHybrid,
		EmploymentType: normalized.EmploymentType,
		ExperienceMin:  normalized.ExperienceMin,
		ExperienceMax:  normalized.ExperienceMax,
		SalaryMin:      normalized.SalaryMin,
		SalaryMax:      normalized.SalaryMax,
		SalaryCurrency: normalized.SalaryCurrency,
		Skills:         pq.StringArray(normalized.Skills),
		ApplicationURL: normalized.ApplicationURL,
		PostedAt:       postedAt,
		DedupHash:      hash,
		IsActive:       true,
	}

	source := &JobSource{
		Provider:   normalized.Provider,
		ExternalID: normalized.ExternalID,
		SourceURL:  normalized.SourceURL,
	}

	return s.repo.UpsertWithSource(ctx, job, source)
}

func (s *service) GetMatchedJobs(ctx context.Context, userID uint, filter *JobFilter) ([]*JobResponse, int64, error) {
	jobs, total, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	matchCfg := s.getMatcherConfig(ctx, userID)

	var result []*JobResponse
	for _, j := range jobs {
		resp := ToJobResponse(j)
		if matchCfg != nil {
			breakdown := s.matcher.Score(j, matchCfg)
			resp.MatchScore = breakdown.Percentage
			resp.MatchBreakdown = breakdown
		}
		result = append(result, resp)
	}
	return result, total, nil
}

func (s *service) getMatcherConfig(ctx context.Context, userID uint) *MatcherConfig {
	if userID == 0 {
		return nil
	}
	profile, err := s.userRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil
	}

	isRemote := strings.Contains(strings.ToLower(profile.PreferredWorkType), "remote")
	for _, loc := range profile.PreferredLocations {
		if strings.EqualFold(loc, "remote") {
			isRemote = true
			break
		}
	}

	return &MatcherConfig{
		Keywords:           profile.PreferredRoles,
		Skills:             profile.PreferredSkills,
		PreferredLocations: profile.PreferredLocations,
		SalaryMin:          profile.ExpectedCTC,
		IsRemote:           isRemote,
		ExperienceYears:    profile.ExperienceYears,
	}
}

func (s *service) TodayCount(ctx context.Context) (int64, error) {
	return s.repo.TodayCount(ctx)
}

func parseDate(dateStr string) *time.Time {
	if dateStr == "" {
		return nil
	}
	formats := []string{
		time.RFC3339,
		"2006-01-02",
		"01/02/2006",
		"2006-01-02T15:04:05Z",
	}
	for _, f := range formats {
		if t, err := time.Parse(f, dateStr); err == nil {
			return &t
		}
	}
	return nil
}

func computeHash(parts ...string) string {
	combined := strings.Join(parts, "|")
	h := sha256.Sum256([]byte(strings.ToLower(combined)))
	return fmt.Sprintf("%x", h)[:16]
}
