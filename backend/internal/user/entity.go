package user

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
	"time"
)

type Profile struct {
	ID                 uint           `gorm:"primarykey" json:"id"`
	UserID             uint           `gorm:"uniqueIndex;not null" json:"user_id"`
	Name               string         `gorm:"not null;default:''" json:"name"`
	Phone              string         `json:"phone"`
	ExperienceYears    float64        `gorm:"type:decimal(4,1);default:0" json:"experience_years"`
	CurrentCompany     string         `json:"current_company"`
	CurrentCTC         float64        `gorm:"type:decimal(15,2)" json:"current_ctc"`
	ExpectedCTC        float64        `gorm:"type:decimal(15,2)" json:"expected_ctc"`
	NoticePeriodDays   int            `gorm:"default:0" json:"notice_period_days"`
	PreferredLocations pq.StringArray `gorm:"type:text[]" json:"preferred_locations"`
	PreferredRoles     pq.StringArray `gorm:"type:text[]" json:"preferred_roles"`
	PreferredSkills    pq.StringArray `gorm:"type:text[]" json:"preferred_skills"`
	Bio                string         `json:"bio"`
	LinkedinURL        string         `json:"linkedin_url"`
	GithubURL          string         `json:"github_url"`
	PortfolioURL       string         `json:"portfolio_url"`
	AvatarURL          string         `json:"avatar_url"`
	IsOpenToWork       bool           `gorm:"default:true" json:"is_open_to_work"`
	PreferredWorkType  string         `gorm:"default:'any'" json:"preferred_work_type"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Profile) TableName() string {
	return "profiles"
}
