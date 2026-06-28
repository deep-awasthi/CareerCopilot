package job

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, job *Job) error {
	return r.db.WithContext(ctx).Create(job).Error
}

func (r *repository) FindByID(ctx context.Context, id uint) (*Job, error) {
	var job Job
	result := r.db.WithContext(ctx).Preload("Sources").First(&job, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("job not found")
		}
		return nil, result.Error
	}
	return &job, nil
}

func (r *repository) FindByDedupHash(ctx context.Context, hash string) (*Job, error) {
	var job Job
	result := r.db.WithContext(ctx).Where("dedup_hash = ?", hash).Preload("Sources").First(&job)
	if result.Error != nil {
		return nil, result.Error
	}
	return &job, nil
}

func (r *repository) List(ctx context.Context, filter *JobFilter) ([]*Job, int64, error) {
	var jobs []*Job
	var total int64

	q := r.db.WithContext(ctx).Model(&Job{}).Where("is_active = true")

	if filter.Query != "" {
		q = q.Where("to_tsvector('english', title || ' ' || COALESCE(description,'')) @@ plainto_tsquery('english', ?)", filter.Query)
	}
	if filter.Location != "" {
		q = q.Where("LOWER(location) LIKE ?", "%"+strings.ToLower(filter.Location)+"%")
	}
	if filter.IsRemote != nil && *filter.IsRemote {
		q = q.Where("is_remote = true")
	}
	if filter.IsHybrid != nil && *filter.IsHybrid {
		q = q.Where("is_hybrid = true")
	}
	if filter.EmploymentType != "" {
		q = q.Where("employment_type = ?", filter.EmploymentType)
	}
	if filter.ExperienceMin > 0 {
		q = q.Where("experience_max >= ? OR experience_max = 0", filter.ExperienceMin)
	}
	if filter.ExperienceMax > 0 {
		q = q.Where("experience_min <= ?", filter.ExperienceMax)
	}
	if filter.SalaryMin > 0 {
		q = q.Where("salary_max >= ? OR salary_max = 0", filter.SalaryMin)
	}
	if filter.Skills != "" {
		skillList := strings.Split(filter.Skills, ",")
		skillArray := pq.StringArray(skillList)
		q = q.Where("skills && ?", skillArray)
	}
	if filter.Provider != "" {
		q = q.Joins("JOIN job_sources js ON js.job_id = jobs.id").Where("js.provider = ?", filter.Provider)
	}

	q.Count(&total)

	page := filter.Page
	if page < 1 {
		page = 1
	}
	perPage := filter.PerPage
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	sortBy := "posted_at"
	validSorts := map[string]bool{"posted_at": true, "created_at": true, "salary_max": true, "experience_min": true}
	if validSorts[filter.SortBy] {
		sortBy = filter.SortBy
	}
	sortOrder := "DESC"
	if strings.ToUpper(filter.SortOrder) == "ASC" {
		sortOrder = "ASC"
	}

	result := q.Order(fmt.Sprintf("%s %s NULLS LAST", sortBy, sortOrder)).
		Offset(offset).Limit(perPage).
		Preload("Sources").
		Find(&jobs)

	return jobs, total, result.Error
}

func (r *repository) Update(ctx context.Context, job *Job) error {
	return r.db.WithContext(ctx).Save(job).Error
}

func (r *repository) AddSource(ctx context.Context, source *JobSource) error {
	return r.db.WithContext(ctx).
		Where(JobSource{JobID: source.JobID, Provider: source.Provider}).
		Assign(*source).
		FirstOrCreate(source).Error
}

func (r *repository) UpsertWithSource(ctx context.Context, job *Job, source *JobSource) error {
	tx := r.db.WithContext(ctx).Begin()

	// Try to find existing by dedup hash
	var existing Job
	if err := tx.Where("dedup_hash = ?", job.DedupHash).First(&existing).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create new job
			if err := tx.Create(job).Error; err != nil {
				tx.Rollback()
				return err
			}
			source.JobID = job.ID
		} else {
			tx.Rollback()
			return err
		}
	} else {
		// Job exists, just add/update source
		source.JobID = existing.ID
	}

	// Upsert source
	source.ScrapedAt = time.Now()
	if err := tx.Where(JobSource{JobID: source.JobID, Provider: source.Provider}).
		Assign(JobSource{SourceURL: source.SourceURL, ExternalID: source.ExternalID, ScrapedAt: source.ScrapedAt}).
		FirstOrCreate(source).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *repository) CountByDate(ctx context.Context, days int) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	r.db.WithContext(ctx).Raw(`
		SELECT DATE(posted_at) as date, COUNT(*) as count
		FROM jobs
		WHERE posted_at >= NOW() - INTERVAL '? days' AND is_active = true
		GROUP BY DATE(posted_at)
		ORDER BY date ASC
	`, days).Scan(&results)
	return results, nil
}

func (r *repository) TopSkills(ctx context.Context, limit int) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	r.db.WithContext(ctx).Raw(`
		SELECT skill, COUNT(*) as count
		FROM jobs, unnest(skills) AS skill
		WHERE is_active = true
		GROUP BY skill
		ORDER BY count DESC
		LIMIT ?
	`, limit).Scan(&results)
	return results, nil
}

func (r *repository) TodayCount(ctx context.Context) (int64, error) {
	var count int64
	r.db.WithContext(ctx).Model(&Job{}).
		Where("DATE(created_at) = DATE(NOW()) AND is_active = true").
		Count(&count)
	return count, nil
}
