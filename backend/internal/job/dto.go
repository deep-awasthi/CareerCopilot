package job

import "time"

type JobFilter struct {
	Query          string  `form:"q"`
	Company        string  `form:"company"`
	Location       string  `form:"location"`
	IsRemote       *bool   `form:"remote"`
	IsHybrid       *bool   `form:"hybrid"`
	EmploymentType string  `form:"type"`
	ExperienceMin  float64 `form:"exp_min"`
	ExperienceMax  float64 `form:"exp_max"`
	SalaryMin      float64 `form:"salary_min"`
	SalaryMax      float64 `form:"salary_max"`
	Skills         string  `form:"skills"`
	Provider       string  `form:"provider"`
	Page           int     `form:"page,default=1"`
	PerPage        int     `form:"per_page,default=20"`
	SortBy         string  `form:"sort_by,default=posted_at"`
	SortOrder      string  `form:"sort_order,default=desc"`
}

type JobResponse struct {
	ID               uint         `json:"id"`
	Title            string       `json:"title"`
	Company          string       `json:"company,omitempty"`
	CompanyID        *uint        `json:"company_id"`
	Location         string       `json:"location"`
	Locations        []string     `json:"locations"`
	IsRemote         bool         `json:"is_remote"`
	IsHybrid         bool         `json:"is_hybrid"`
	EmploymentType   string       `json:"employment_type"`
	ExperienceMin    float64      `json:"experience_min"`
	ExperienceMax    float64      `json:"experience_max"`
	SalaryMin        float64      `json:"salary_min"`
	SalaryMax        float64      `json:"salary_max"`
	SalaryCurrency   string       `json:"salary_currency"`
	Skills           []string     `json:"skills"`
	ApplicationURL   string       `json:"application_url"`
	PostedAt         *time.Time   `json:"posted_at"`
	Sources          []SourceDTO  `json:"sources"`
	MatchScore       int          `json:"match_score,omitempty"`
	MatchBreakdown   *ScoreBreakdown `json:"match_breakdown,omitempty"`
	IsBookmarked     bool         `json:"is_bookmarked,omitempty"`
	CreatedAt        time.Time    `json:"created_at"`
}

type SourceDTO struct {
	Provider  string `json:"provider"`
	SourceURL string `json:"source_url"`
}

type CreateJobRequest struct {
	CompanyID      *uint    `json:"company_id"`
	Title          string   `json:"title" validate:"required,min=3,max=500"`
	Description    string   `json:"description"`
	Location       string   `json:"location"`
	IsRemote       bool     `json:"is_remote"`
	IsHybrid       bool     `json:"is_hybrid"`
	EmploymentType string   `json:"employment_type"`
	ExperienceMin  float64  `json:"experience_min"`
	ExperienceMax  float64  `json:"experience_max"`
	SalaryMin      float64  `json:"salary_min"`
	SalaryMax      float64  `json:"salary_max"`
	Skills         []string `json:"skills"`
	ApplicationURL string   `json:"application_url"`
}

func ToJobResponse(j *Job) *JobResponse {
	resp := &JobResponse{
		ID:             j.ID,
		Title:          j.Title,
		CompanyID:      j.CompanyID,
		Location:       j.Location,
		Locations:      j.Locations,
		IsRemote:       j.IsRemote,
		IsHybrid:       j.IsHybrid,
		EmploymentType: j.EmploymentType,
		ExperienceMin:  j.ExperienceMin,
		ExperienceMax:  j.ExperienceMax,
		SalaryMin:      j.SalaryMin,
		SalaryMax:      j.SalaryMax,
		SalaryCurrency: j.SalaryCurrency,
		Skills:         j.Skills,
		ApplicationURL: j.ApplicationURL,
		PostedAt:       j.PostedAt,
		MatchScore:     j.MatchScore,
		MatchBreakdown: j.MatchBreakdown,
		CreatedAt:      j.CreatedAt,
	}
	if resp.Skills == nil {
		resp.Skills = []string{}
	}
	if resp.Locations == nil {
		resp.Locations = []string{}
	}
	for _, s := range j.Sources {
		resp.Sources = append(resp.Sources, SourceDTO{
			Provider:  s.Provider,
			SourceURL: s.SourceURL,
		})
	}
	return resp
}
