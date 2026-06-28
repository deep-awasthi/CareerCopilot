package auth

import (
	"time"

	"gorm.io/gorm"
)

// User entity
type User struct {
	ID                   uint           `gorm:"primarykey" json:"id"`
	UUID                 string         `gorm:"uniqueIndex;not null" json:"uuid"`
	Email                string         `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash         string         `gorm:"not null" json:"-"`
	IsEmailVerified      bool           `gorm:"default:false" json:"is_email_verified"`
	EmailVerifyToken     string         `json:"-"`
	PasswordResetToken   string         `json:"-"`
	PasswordResetExpires *time.Time     `json:"-"`
	LastLoginAt          *time.Time     `json:"last_login_at"`
	IsActive             bool           `gorm:"default:true" json:"is_active"`
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at"`
	DeletedAt            gorm.DeletedAt `gorm:"index" json:"-"`
}

func (User) TableName() string {
	return "users"
}
