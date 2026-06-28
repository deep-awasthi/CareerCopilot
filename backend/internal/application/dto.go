package application

import "time"

type CreateApplicationRequest struct {
	JobID         uint      `json:"job_id" validate:"required"`
	Status        string    `json:"status" validate:"omitempty,oneof=saved applied interview offer rejected archived"`
	Notes         string    `json:"notes"`
	FollowUpDate  *time.Time `json:"follow_up_date"`
	SalaryOffered float64   `json:"salary_offered"`
	ReferralUsed  bool      `json:"referral_used"`
	ReferralID    *uint     `json:"referral_id"`
}

type UpdateApplicationRequest struct {
	Status        string    `json:"status" validate:"omitempty,oneof=saved applied interview offer rejected archived"`
	Notes         string    `json:"notes"`
	FollowUpDate  *time.Time `json:"follow_up_date"`
	SalaryOffered float64   `json:"salary_offered"`
	ReferralUsed  *bool     `json:"referral_used"`
	AppliedAt     *time.Time `json:"applied_at"`
}

type ApplicationFilter struct {
	Status  string `form:"status"`
	Page    int    `form:"page,default=1"`
	PerPage int    `form:"per_page,default=20"`
}

type ApplicationResponse struct {
	ID            uint      `json:"id"`
	UserID        uint      `json:"user_id"`
	JobID         uint      `json:"job_id"`
	Status        string    `json:"status"`
	Notes         string    `json:"notes"`
	FollowUpDate  *time.Time `json:"follow_up_date"`
	SalaryOffered float64   `json:"salary_offered"`
	ReferralUsed  bool      `json:"referral_used"`
	AppliedAt     *time.Time `json:"applied_at"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func ToResponse(a *Application) *ApplicationResponse {
	return &ApplicationResponse{
		ID:            a.ID,
		UserID:        a.UserID,
		JobID:         a.JobID,
		Status:        string(a.Status),
		Notes:         a.Notes,
		FollowUpDate:  a.FollowUpDate,
		SalaryOffered: a.SalaryOffered,
		ReferralUsed:  a.ReferralUsed,
		AppliedAt:     a.AppliedAt,
		CreatedAt:     a.CreatedAt,
		UpdatedAt:     a.UpdatedAt,
	}
}

type StatusCount struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
}
