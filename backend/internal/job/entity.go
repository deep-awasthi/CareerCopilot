package job

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Job struct {
	ID               uint           `gorm:"primarykey" json:"id"`
	ExternalID       string         `json:"external_id"`
	CompanyID        *uint          `json:"company_id"`
	Title            string         `gorm:"not null" json:"title"`
	Description      string         `gorm:"type:text" json:"description"`
	ShortDescription string         `json:"short_description"`
	Location         string         `json:"location"`
	Locations        pq.StringArray `gorm:"type:text[]" json:"locations"`
	IsRemote         bool           `json:"is_remote"`
	IsHybrid         bool           `json:"is_hybrid"`
	EmploymentType   string         `json:"employment_type"`
	ExperienceMin    float64        `gorm:"type:decimal(4,1)" json:"experience_min"`
	ExperienceMax    float64        `gorm:"type:decimal(4,1)" json:"experience_max"`
	SalaryMin        float64        `gorm:"type:decimal(15,2)" json:"salary_min"`
	SalaryMax        float64        `gorm:"type:decimal(15,2)" json:"salary_max"`
	SalaryCurrency   string         `gorm:"default:'INR'" json:"salary_currency"`
	Skills           pq.StringArray `gorm:"type:text[]" json:"skills"`
	ApplicationURL   string         `json:"application_url"`
	PostedAt         *time.Time     `json:"posted_at"`
	ExpiresAt        *time.Time     `json:"expires_at"`
	IsActive         bool           `gorm:"default:true" json:"is_active"`
	IsVerified       bool           `gorm:"default:false" json:"is_verified"`
	DedupHash        string         `gorm:"uniqueIndex" json:"dedup_hash"`
	ViewCount        int            `gorm:"default:0" json:"view_count"`
	Sources          []JobSource    `gorm:"foreignKey:JobID" json:"sources,omitempty"`
	MatchScore       int            `gorm:"-" json:"match_score,omitempty"`
	MatchBreakdown   *ScoreBreakdown `gorm:"-" json:"match_breakdown,omitempty"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Job) TableName() string {
	return "jobs"
}

type JobSource struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	JobID      uint      `gorm:"not null" json:"job_id"`
	Provider   string    `gorm:"not null" json:"provider"`
	ExternalID string    `json:"external_id"`
	SourceURL  string    `json:"source_url"`
	ScrapedAt  time.Time `json:"scraped_at"`
	CreatedAt  time.Time `json:"created_at"`
}

func (JobSource) TableName() string {
	return "job_sources"
}
