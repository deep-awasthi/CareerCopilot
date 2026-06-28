package user

type UpdateProfileRequest struct {
	Name               string   `json:"name" validate:"omitempty,min=2,max=100"`
	Phone              string   `json:"phone" validate:"omitempty,max=20"`
	ExperienceYears    float64  `json:"experience_years" validate:"omitempty,gte=0,lte=50"`
	CurrentCompany     string   `json:"current_company" validate:"omitempty,max=255"`
	CurrentCTC         float64  `json:"current_ctc" validate:"omitempty,gte=0"`
	ExpectedCTC        float64  `json:"expected_ctc" validate:"omitempty,gte=0"`
	NoticePeriodDays   int      `json:"notice_period_days" validate:"omitempty,gte=0,lte=365"`
	PreferredLocations []string `json:"preferred_locations"`
	PreferredRoles     []string `json:"preferred_roles"`
	PreferredSkills    []string `json:"preferred_skills"`
	Bio                string   `json:"bio" validate:"omitempty,max=1000"`
	LinkedinURL        string   `json:"linkedin_url" validate:"omitempty,url"`
	GithubURL          string   `json:"github_url" validate:"omitempty,url"`
	PortfolioURL       string   `json:"portfolio_url" validate:"omitempty,url"`
	IsOpenToWork       *bool    `json:"is_open_to_work"`
	PreferredWorkType  string   `json:"preferred_work_type" validate:"omitempty,oneof=remote hybrid onsite any"`
}

type ProfileResponse struct {
	ID                 uint     `json:"id"`
	UserID             uint     `json:"user_id"`
	Name               string   `json:"name"`
	Phone              string   `json:"phone"`
	ExperienceYears    float64  `json:"experience_years"`
	CurrentCompany     string   `json:"current_company"`
	CurrentCTC         float64  `json:"current_ctc"`
	ExpectedCTC        float64  `json:"expected_ctc"`
	NoticePeriodDays   int      `json:"notice_period_days"`
	PreferredLocations []string `json:"preferred_locations"`
	PreferredRoles     []string `json:"preferred_roles"`
	PreferredSkills    []string `json:"preferred_skills"`
	Bio                string   `json:"bio"`
	LinkedinURL        string   `json:"linkedin_url"`
	GithubURL          string   `json:"github_url"`
	PortfolioURL       string   `json:"portfolio_url"`
	AvatarURL          string   `json:"avatar_url"`
	IsOpenToWork       bool     `json:"is_open_to_work"`
	PreferredWorkType  string   `json:"preferred_work_type"`
}

func ToProfileResponse(p *Profile) *ProfileResponse {
	return &ProfileResponse{
		ID:                 p.ID,
		UserID:             p.UserID,
		Name:               p.Name,
		Phone:              p.Phone,
		ExperienceYears:    p.ExperienceYears,
		CurrentCompany:     p.CurrentCompany,
		CurrentCTC:         p.CurrentCTC,
		ExpectedCTC:        p.ExpectedCTC,
		NoticePeriodDays:   p.NoticePeriodDays,
		PreferredLocations: p.PreferredLocations,
		PreferredRoles:     p.PreferredRoles,
		PreferredSkills:    p.PreferredSkills,
		Bio:                p.Bio,
		LinkedinURL:        p.LinkedinURL,
		GithubURL:          p.GithubURL,
		PortfolioURL:       p.PortfolioURL,
		AvatarURL:          p.AvatarURL,
		IsOpenToWork:       p.IsOpenToWork,
		PreferredWorkType:  p.PreferredWorkType,
	}
}
