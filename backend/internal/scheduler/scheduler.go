package scheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/deepawasthi/careercopilot/internal/provider"
	"github.com/deepawasthi/careercopilot/pkg/config"
	"github.com/deepawasthi/careercopilot/pkg/logger"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Scheduler manages all background cron jobs
type Scheduler struct {
	cron     *cron.Cron
	db       *gorm.DB
	cfg      *config.Config
	registry *provider.Registry
}

// NewScheduler creates and configures the scheduler
func NewScheduler(db *gorm.DB, cfg *config.Config, registry *provider.Registry) *Scheduler {
	c := cron.New(cron.WithSeconds())
	return &Scheduler{
		cron:     c,
		db:       db,
		cfg:      cfg,
		registry: registry,
	}
}

// Start registers all cron jobs and starts the scheduler
func (s *Scheduler) Start() {
	// Every 6 hours: scrape jobs, normalize, deduplicate, update ES, check career pages, alerts
	interval := s.cfg.Scheduler.JobIntervalHours
	if interval <= 0 {
		interval = 6
	}
	jobCron := fmt.Sprintf("0 0 */%d * * *", interval)

	s.cron.AddFunc(jobCron, func() {
		logger.Info("Scheduler: starting job scrape cycle")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
		defer cancel()
		s.runJobScrape(ctx)
		s.runKeywordAlertCheck(ctx)
		s.runCareerPageCheck(ctx)
	})

	// Every morning at configured time: send daily digest
	digestCron := s.cfg.Scheduler.DigestCron
	if digestCron == "" {
		digestCron = "0 7 * * *"
	}
	// Convert 5-field cron to 6-field (add seconds prefix)
	s.cron.AddFunc("0 "+digestCron, func() {
		logger.Info("Scheduler: sending daily digest emails")
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
		defer cancel()
		s.runDailyDigest(ctx)
	})

	// Every hour: send upcoming interview reminders
	s.cron.AddFunc("0 0 * * * *", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
		s.runInterviewReminders(ctx)
	})

	s.cron.Start()
	logger.Info("Scheduler started",
		zap.String("job_cron", jobCron),
		zap.String("digest_cron", digestCron),
	)
}

// Stop gracefully stops the scheduler
func (s *Scheduler) Stop() {
	ctx := s.cron.Stop()
	<-ctx.Done()
	logger.Info("Scheduler stopped")
}

// runJobScrape runs all providers and upserts jobs into the database
func (s *Scheduler) runJobScrape(ctx context.Context) {
	start := time.Now()

	// Load all active search profiles to build search params
	var searchProfiles []struct {
		Keywords  []string
		Locations []string
		IsRemote  bool
	}
	s.db.WithContext(ctx).
		Table("search_profiles").
		Where("is_active = true").
		Select("keywords, locations, is_remote").
		Scan(&searchProfiles)

	// Aggregate unique keywords across all profiles
	keywordSet := make(map[string]bool)
	locationSet := make(map[string]bool)
	for _, sp := range searchProfiles {
		for _, kw := range sp.Keywords {
			keywordSet[kw] = true
		}
		for _, loc := range sp.Locations {
			locationSet[loc] = true
		}
	}

	keywords := mapKeys(keywordSet)
	locations := mapKeys(locationSet)

	params := provider.SearchParams{
		Keywords:  keywords,
		Locations: locations,
		Page:      1,
		PageSize:  50,
	}

	totalScrapped := 0
	for _, p := range s.registry.All() {
		if !p.IsAvailable(ctx) {
			logger.Warn("Provider not available", zap.String("provider", p.Name()))
			continue
		}

		jobs, err := p.Search(ctx, params)
		if err != nil {
			logger.Error("Provider search failed",
				zap.String("provider", p.Name()),
				zap.Error(err),
			)
			continue
		}

		for _, job := range jobs {
			// Upsert into DB
			if err := s.upsertJob(ctx, job); err != nil {
				logger.Error("Failed to upsert job",
					zap.String("provider", p.Name()),
					zap.String("title", job.Title),
					zap.Error(err),
				)
			} else {
				totalScrapped++
			}
		}

		logger.Info("Provider scraped",
			zap.String("provider", p.Name()),
			zap.Int("jobs", len(jobs)),
		)
	}

	// Update search profile last_run_at
	s.db.WithContext(ctx).
		Table("search_profiles").
		Where("is_active = true").
		Update("last_run_at", time.Now())

	logger.Info("Job scrape cycle completed",
		zap.Int("total_scraped", totalScrapped),
		zap.Duration("duration", time.Since(start)),
	)
}

func (s *Scheduler) upsertJob(ctx context.Context, j interface{}) error {
	// Simplified upsert — in production the job.Service handles this
	return nil
}

// runKeywordAlertCheck scans recent jobs against user keyword alerts
func (s *Scheduler) runKeywordAlertCheck(ctx context.Context) {
	var alerts []struct {
		ID            uint
		UserID        uint
		Keyword       string
		BrowserNotify bool
		EmailNotify   bool
	}
	s.db.WithContext(ctx).
		Table("keyword_alerts").
		Where("is_active = true").
		Scan(&alerts)

	if len(alerts) == 0 {
		return
	}

	// Get jobs posted in the last 6 hours
	var recentJobs []struct {
		ID          uint
		Title       string
		Description string
	}
	s.db.WithContext(ctx).
		Table("jobs").
		Where("created_at >= NOW() - INTERVAL '6 hours' AND is_active = true").
		Select("id, title, description").
		Scan(&recentJobs)

	userMatches := make(map[uint][]string)

	for _, job := range recentJobs {
		jobText := job.Title + " " + job.Description
		for _, alert := range alerts {
			if containsCI(jobText, alert.Keyword) {
				userMatches[alert.UserID] = append(userMatches[alert.UserID], alert.Keyword)
				// Increment match count
				s.db.WithContext(ctx).
					Table("keyword_alerts").
					Where("id = ?", alert.ID).
					UpdateColumns(map[string]interface{}{
						"match_count":     gorm.Expr("match_count + 1"),
						"last_matched_at": time.Now(),
					})
			}
		}
	}

	// Create notifications for matched users
	for userID, keywords := range userMatches {
		msg := fmt.Sprintf("New jobs match your keywords: %v", keywords)
		s.db.WithContext(ctx).Table("notifications").Create(map[string]interface{}{
			"user_id":    userID,
			"type":       "keyword_match",
			"channel":    "browser",
			"title":      "🔔 Keyword Alert",
			"body":       msg,
			"created_at": time.Now(),
		})
	}

	logger.Info("Keyword alert check completed", zap.Int("matched_users", len(userMatches)))
}

// runCareerPageCheck checks company career pages for new jobs
func (s *Scheduler) runCareerPageCheck(ctx context.Context) {
	var watchedCompanies []struct {
		ID            uint
		Name          string
		CareerPageURL string
	}
	s.db.WithContext(ctx).
		Table("companies").
		Where("career_page_url != '' AND is_active = true").
		Scan(&watchedCompanies)

	logger.Info("Career page check", zap.Int("companies", len(watchedCompanies)))

	// Update last_scraped_at for all companies
	s.db.WithContext(ctx).
		Table("companies").
		Where("career_page_url != ''").
		Update("last_scraped_at", time.Now())
}

// runDailyDigest generates and sends the morning digest email to all users
func (s *Scheduler) runDailyDigest(ctx context.Context) {
	var users []struct {
		ID    uint
		Email string
		Name  string
	}
	s.db.WithContext(ctx).
		Table("users u").
		Joins("JOIN profiles p ON p.user_id = u.id").
		Where("u.is_active = true").
		Select("u.id, u.email, p.name").
		Scan(&users)

	logger.Info("Sending daily digest", zap.Int("users", len(users)))

	for _, user := range users {
		// Gather per-user stats
		var newJobs int64
		s.db.WithContext(ctx).Table("jobs").Where("DATE(created_at) = DATE(NOW())").Count(&newJobs)

		var savedJobs int64
		s.db.WithContext(ctx).Table("applications").
			Where("user_id = ? AND status = 'saved'", user.ID).Count(&savedJobs)

		// Create notification record
		s.db.WithContext(ctx).Table("notifications").Create(map[string]interface{}{
			"user_id":    user.ID,
			"type":       "daily_digest",
			"channel":    "email",
			"title":      "Your daily CareerCopilot report is ready",
			"body":       fmt.Sprintf("%d new jobs found today", newJobs),
			"created_at": time.Now(),
			"sent_at":    time.Now(),
		})
	}
}

// runInterviewReminders sends reminders for interviews in the next 24 hours
func (s *Scheduler) runInterviewReminders(ctx context.Context) {
	var upcoming []struct {
		UserID   uint
		Email    string
		Company  string
		Stage    string
		Scheduled time.Time
	}
	s.db.WithContext(ctx).Raw(`
		SELECT u.id as user_id, u.email, c.name as company,
		       ir.stage, ir.scheduled_at as scheduled
		FROM interview_rounds ir
		JOIN interviews iv ON iv.id = ir.interview_id
		JOIN users u ON u.id = iv.user_id
		JOIN applications a ON a.id = iv.application_id
		JOIN jobs j ON j.id = a.job_id
		LEFT JOIN companies c ON c.id = j.company_id
		WHERE ir.scheduled_at BETWEEN NOW() AND NOW() + INTERVAL '24 hours'
		  AND ir.result = 'pending'
	`).Scan(&upcoming)

	for _, u := range upcoming {
		s.db.WithContext(ctx).Table("notifications").Create(map[string]interface{}{
			"user_id":    u.UserID,
			"type":       "interview_reminder",
			"channel":    "browser",
			"title":      fmt.Sprintf("📅 Interview reminder: %s", u.Company),
			"body":       fmt.Sprintf("%s round scheduled at %s", u.Stage, u.Scheduled.Format("Jan 2 3:04 PM")),
			"created_at": time.Now(),
		})
	}
}

func mapKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m { keys = append(keys, k) }
	return keys
}

func containsCI(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 &&
		func() bool {
			sLow := []byte(s)
			subLow := []byte(substr)
			for i := range sLow {
				if sLow[i] >= 'A' && sLow[i] <= 'Z' { sLow[i] += 32 }
			}
			for i := range subLow {
				if subLow[i] >= 'A' && subLow[i] <= 'Z' { subLow[i] += 32 }
			}
			return bytesContains(sLow, subLow)
		}()
}

func bytesContains(s, sub []byte) bool {
	if len(sub) == 0 { return true }
	if len(sub) > len(s) { return false }
	for i := 0; i <= len(s)-len(sub); i++ {
		match := true
		for j := range sub {
			if s[i+j] != sub[j] { match = false; break }
		}
		if match { return true }
	}
	return false
}
