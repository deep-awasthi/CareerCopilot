package analytics

import (
	"context"

	"github.com/deepawasthi/careercopilot/pkg/config"
	"github.com/deepawasthi/careercopilot/pkg/middleware"
	"github.com/deepawasthi/careercopilot/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DashboardStats struct {
	JobsFoundToday     int64              `json:"jobs_found_today"`
	SavedJobs          int64              `json:"saved_jobs"`
	AppliedJobs        int64              `json:"applied_jobs"`
	Interviews         int64              `json:"interviews"`
	Offers             int64              `json:"offers"`
	ReferralOpportunities int64           `json:"referral_opportunities"`
	CompaniesFollowing int64              `json:"companies_following"`
	KeywordAlerts      int64              `json:"keyword_alerts"`
	UpcomingInterviews int64              `json:"upcoming_interviews"`
}

type ApplicationsPerMonth struct {
	Month string `json:"month"`
	Count int64  `json:"count"`
}

type TopCompany struct {
	Company string `json:"company"`
	Count   int64  `json:"count"`
}

type SkillDemand struct {
	Skill string `json:"skill"`
	Count int64  `json:"count"`
}

type SalaryBucket struct {
	Range string  `json:"range"`
	Count int64   `json:"count"`
	Avg   float64 `json:"avg"`
}

type AnalyticsData struct {
	ApplicationsPerMonth []ApplicationsPerMonth `json:"applications_per_month"`
	InterviewSuccessRate float64                `json:"interview_success_rate"`
	TopCompanies         []TopCompany           `json:"top_companies"`
	TopSkills            []SkillDemand          `json:"top_skills"`
	AverageSalary        float64                `json:"average_salary"`
	ResponseRate         float64                `json:"response_rate"`
	ReferralSuccessRate  float64                `json:"referral_success_rate"`
	SalaryDistribution   []SalaryBucket         `json:"salary_distribution"`
}

type Service interface {
	GetDashboardStats(ctx context.Context, userID uint) (*DashboardStats, error)
	GetAnalytics(ctx context.Context, userID uint) (*AnalyticsData, error)
}

type service struct{ db *gorm.DB }

func NewService(db *gorm.DB) Service { return &service{db: db} }

func (s *service) GetDashboardStats(ctx context.Context, userID uint) (*DashboardStats, error) {
	stats := &DashboardStats{}

	// Jobs found today
	s.db.WithContext(ctx).Table("jobs").
		Where("DATE(created_at) = DATE(NOW()) AND is_active = true").
		Count(&stats.JobsFoundToday)

	// Application stats
	s.db.WithContext(ctx).Table("applications").Where("user_id = ? AND status = 'saved'", userID).Count(&stats.SavedJobs)
	s.db.WithContext(ctx).Table("applications").Where("user_id = ? AND status = 'applied'", userID).Count(&stats.AppliedJobs)
	s.db.WithContext(ctx).Table("applications").Where("user_id = ? AND status = 'interview'", userID).Count(&stats.Interviews)
	s.db.WithContext(ctx).Table("applications").Where("user_id = ? AND status = 'offer'", userID).Count(&stats.Offers)

	// Referrals
	s.db.WithContext(ctx).Table("referrals").Where("user_id = ?", userID).Count(&stats.ReferralOpportunities)

	// Companies following
	s.db.WithContext(ctx).Table("company_watchlists").Where("user_id = ?", userID).Count(&stats.CompaniesFollowing)

	// Keyword alerts
	s.db.WithContext(ctx).Table("keyword_alerts").Where("user_id = ? AND is_active = true", userID).Count(&stats.KeywordAlerts)

	// Upcoming interviews
	s.db.WithContext(ctx).Table("interview_rounds").
		Joins("JOIN interviews iv ON iv.id = interview_rounds.interview_id").
		Where("iv.user_id = ? AND interview_rounds.scheduled_at > NOW()", userID).
		Count(&stats.UpcomingInterviews)

	return stats, nil
}

func (s *service) GetAnalytics(ctx context.Context, userID uint) (*AnalyticsData, error) {
	data := &AnalyticsData{}

	// Applications per month (last 12 months)
	var appsPerMonth []ApplicationsPerMonth
	s.db.WithContext(ctx).Raw(`
		SELECT TO_CHAR(DATE_TRUNC('month', created_at), 'Mon YYYY') as month,
		       COUNT(*) as count
		FROM applications
		WHERE user_id = ?
		  AND created_at >= NOW() - INTERVAL '12 months'
		GROUP BY DATE_TRUNC('month', created_at)
		ORDER BY DATE_TRUNC('month', created_at)
	`, userID).Scan(&appsPerMonth)
	data.ApplicationsPerMonth = appsPerMonth

	// Interview success rate
	var totalInterviews, passedInterviews int64
	s.db.WithContext(ctx).Table("interview_rounds").
		Joins("JOIN interviews iv ON iv.id = interview_rounds.interview_id").
		Where("iv.user_id = ?", userID).Count(&totalInterviews)
	s.db.WithContext(ctx).Table("interview_rounds").
		Joins("JOIN interviews iv ON iv.id = interview_rounds.interview_id").
		Where("iv.user_id = ? AND interview_rounds.result = 'passed'", userID).Count(&passedInterviews)
	if totalInterviews > 0 {
		data.InterviewSuccessRate = float64(passedInterviews) / float64(totalInterviews) * 100
	}

	// Top companies applying to
	var topCompanies []TopCompany
	s.db.WithContext(ctx).Raw(`
		SELECT c.name as company, COUNT(a.id) as count
		FROM applications a
		JOIN jobs j ON j.id = a.job_id
		JOIN companies c ON c.id = j.company_id
		WHERE a.user_id = ?
		GROUP BY c.name
		ORDER BY count DESC
		LIMIT 10
	`, userID).Scan(&topCompanies)
	data.TopCompanies = topCompanies

	// Top demanded skills (across all active jobs)
	var topSkills []SkillDemand
	s.db.WithContext(ctx).Raw(`
		SELECT skill, COUNT(*) as count
		FROM jobs, unnest(skills) AS skill
		WHERE is_active = true
		GROUP BY skill
		ORDER BY count DESC
		LIMIT 20
	`).Scan(&topSkills)
	data.TopSkills = topSkills

	// Average salary of offers
	var avgSalary struct{ Avg float64 }
	s.db.WithContext(ctx).Raw(`
		SELECT AVG(salary_offered) as avg
		FROM applications
		WHERE user_id = ? AND status = 'offer' AND salary_offered > 0
	`, userID).Scan(&avgSalary)
	data.AverageSalary = avgSalary.Avg

	// Response rate (applied → interview)
	var applied, gotInterview int64
	s.db.WithContext(ctx).Table("applications").Where("user_id = ? AND status != 'saved'", userID).Count(&applied)
	s.db.WithContext(ctx).Table("applications").Where("user_id = ? AND status IN ('interview','offer')", userID).Count(&gotInterview)
	if applied > 0 {
		data.ResponseRate = float64(gotInterview) / float64(applied) * 100
	}

	// Referral success rate
	var totalReferrals, successfulReferrals int64
	s.db.WithContext(ctx).Table("referrals").Where("user_id = ?", userID).Count(&totalReferrals)
	s.db.WithContext(ctx).Table("referrals").Where("user_id = ? AND status IN ('referral_received','applied')", userID).Count(&successfulReferrals)
	if totalReferrals > 0 {
		data.ReferralSuccessRate = float64(successfulReferrals) / float64(totalReferrals) * 100
	}

	return data, nil
}

type Controller struct{ svc Service }

func NewController(svc Service) *Controller { return &Controller{svc: svc} }

func (c *Controller) Dashboard(ctx *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok { response.Unauthorized(ctx, "unauthorized"); return }
	stats, err := c.svc.GetDashboardStats(ctx.Request.Context(), userID)
	if err != nil { response.Error(ctx, err); return }
	response.Success(ctx, "dashboard stats", stats)
}

func (c *Controller) Analytics(ctx *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok { response.Unauthorized(ctx, "unauthorized"); return }
	data, err := c.svc.GetAnalytics(ctx.Request.Context(), userID)
	if err != nil { response.Error(ctx, err); return }
	response.Success(ctx, "analytics data", data)
}

func RegisterRoutes(r *gin.RouterGroup, ctrl *Controller, cfg *config.Config) {
	analytics := r.Group("/analytics")
	analytics.Use(middleware.JWTAuth(&cfg.JWT))
	{
		analytics.GET("/dashboard", ctrl.Dashboard)
		analytics.GET("", ctrl.Analytics)
	}
}
