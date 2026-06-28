package resume

type SubmitResumeRequest struct {
	RawText string `json:"raw_text" validate:"required,min=100"`
}

type ResumeResponse struct {
	ID             uint                     `json:"id"`
	UserID         uint                     `json:"user_id"`
	RawText        string                   `json:"raw_text"`
	ParsedSkills   []string                 `json:"parsed_skills"`
	Companies      []string                 `json:"companies"`
	Projects       []string                 `json:"projects"`
	Education      []map[string]interface{} `json:"education"`
	Experience     []map[string]interface{} `json:"experience"`
	Certifications []string                 `json:"certifications"`
}

func ToResumeResponse(r *Resume) *ResumeResponse {
	resp := &ResumeResponse{
		ID:             r.ID,
		UserID:         r.UserID,
		RawText:        r.RawText,
		ParsedSkills:   r.ParsedSkills,
		Companies:      r.Companies,
		Projects:       r.Projects,
		Education:      r.Education,
		Experience:     r.Experience,
		Certifications: r.Certifications,
	}
	if resp.ParsedSkills == nil {
		resp.ParsedSkills = []string{}
	}
	if resp.Companies == nil {
		resp.Companies = []string{}
	}
	if resp.Projects == nil {
		resp.Projects = []string{}
	}
	if resp.Education == nil {
		resp.Education = []map[string]interface{}{}
	}
	if resp.Experience == nil {
		resp.Experience = []map[string]interface{}{}
	}
	if resp.Certifications == nil {
		resp.Certifications = []string{}
	}
	return resp
}
