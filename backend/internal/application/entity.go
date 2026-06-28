package application

import (
	"time"
	"gorm.io/gorm"
)

type ApplicationStatus string

const (
	StatusSaved     ApplicationStatus = "saved"
	StatusApplied   ApplicationStatus = "applied"
	StatusInterview ApplicationStatus = "interview"
	StatusOffer     ApplicationStatus = "offer"
	StatusRejected  ApplicationStatus = "rejected"
	StatusArchived  ApplicationStatus = "archived"
)

type Application struct {
	ID            uint              `gorm:"primarykey" json:"id"`
	UserID        uint              `gorm:"not null;index" json:"user_id"`
	JobID         uint              `gorm:"not null;index" json:"job_id"`
	Status        ApplicationStatus `gorm:"type:application_status;default:'saved'" json:"status"`
	Notes         string            `gorm:"type:text" json:"notes"`
	FollowUpDate  *time.Time        `json:"follow_up_date"`
	SalaryOffered float64           `gorm:"type:decimal(15,2)" json:"salary_offered"`
	ReferralUsed  bool              `gorm:"default:false" json:"referral_used"`
	ReferralID    *uint             `json:"referral_id"`
	AppliedAt     *time.Time        `json:"applied_at"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
	DeletedAt     gorm.DeletedAt    `gorm:"index" json:"-"`
}

func (Application) TableName() string { return "applications" }
