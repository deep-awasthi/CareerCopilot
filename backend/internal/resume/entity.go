package resume

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

// JSONB helper type
type JSONB map[string]interface{}

func (j JSONB) Value() (driver.Value, error) {
	valueString, err := json.Marshal(j)
	return string(valueString), err
}

func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	s, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan JSONB")
	}
	return json.Unmarshal(s, j)
}

type JSONBArray []map[string]interface{}

func (j JSONBArray) Value() (driver.Value, error) {
	valueString, err := json.Marshal(j)
	return string(valueString), err
}

func (j *JSONBArray) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	s, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan JSONBArray")
	}
	return json.Unmarshal(s, j)
}

type Resume struct {
	ID             uint           `gorm:"primarykey" json:"id"`
	UserID         uint           `gorm:"uniqueIndex;not null" json:"user_id"`
	RawText        string         `gorm:"type:text;not null" json:"raw_text"`
	ParsedSkills   pq.StringArray `gorm:"type:text[]" json:"parsed_skills"`
	Companies      pq.StringArray `gorm:"type:text[]" json:"companies"`
	Projects       pq.StringArray `gorm:"type:text[]" json:"projects"`
	Education      JSONBArray     `gorm:"type:jsonb" json:"education"`
	Experience     JSONBArray     `gorm:"type:jsonb" json:"experience"`
	Certifications pq.StringArray `gorm:"type:text[]" json:"certifications"`
	ParsedAt       *time.Time     `json:"parsed_at"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Resume) TableName() string {
	return "resumes"
}
