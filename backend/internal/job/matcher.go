package job

import (
	"strings"
)

// MatcherConfig holds user preferences for job matching
type MatcherConfig struct {
	Keywords           []string
	Skills             []string
	PreferredLocations []string
	SalaryMin          float64
	SalaryMax          float64
	IsRemote           bool
	ExperienceYears    float64
	PreferredCompanies []string
}

// ScoreBreakdown shows per-factor score contribution
type ScoreBreakdown struct {
	KeywordMatch    int `json:"keyword_match"`
	SkillMatch      int `json:"skill_match"`
	LocationMatch   int `json:"location_match"`
	SalaryMatch     int `json:"salary_match"`
	RemoteMatch     int `json:"remote_match"`
	ExperienceMatch int `json:"experience_match"`
	CompanyMatch    int `json:"company_match"`
	Total           int `json:"total"`
	MaxPossible     int `json:"max_possible"`
	Percentage      int `json:"percentage"`
}

const (
	scoreKeyword    = 20
	scoreSkill      = 15
	scoreLocation   = 15
	scoreSalary     = 15
	scoreRemote     = 10
	scoreExperience = 10
	scoreCompany    = 10
	maxScore        = scoreKeyword + scoreSkill + scoreLocation + scoreSalary + scoreRemote + scoreExperience + scoreCompany
)

// Matcher implements rule-based job matching — NO AI
type Matcher struct{}

func NewMatcher() *Matcher {
	return &Matcher{}
}

// Score computes a match percentage for a job against user preferences
func (m *Matcher) Score(job *Job, cfg *MatcherConfig) *ScoreBreakdown {
	breakdown := &ScoreBreakdown{MaxPossible: maxScore}

	jobText := strings.ToLower(job.Title + " " + job.Description + " " + strings.Join(job.Skills, " "))

	// +20: Keyword match — any keyword found in title/description
	for _, kw := range cfg.Keywords {
		if strings.Contains(jobText, strings.ToLower(kw)) {
			breakdown.KeywordMatch = scoreKeyword
			break
		}
	}

	// +15: Skill match — at least half of user's skills appear in job skills
	if len(cfg.Skills) > 0 && len(job.Skills) > 0 {
		matches := 0
		jobSkillsLower := toLower(job.Skills)
		for _, skill := range cfg.Skills {
			for _, js := range jobSkillsLower {
				if strings.Contains(js, strings.ToLower(skill)) || strings.Contains(strings.ToLower(skill), js) {
					matches++
					break
				}
			}
		}
		ratio := float64(matches) / float64(len(cfg.Skills))
		if ratio >= 0.5 {
			breakdown.SkillMatch = scoreSkill
		} else if ratio >= 0.25 {
			breakdown.SkillMatch = scoreSkill / 2
		}
	}

	// +15: Location match
	if len(cfg.PreferredLocations) > 0 {
		jobLocText := strings.ToLower(job.Location + " " + strings.Join(job.Locations, " "))
		for _, loc := range cfg.PreferredLocations {
			if strings.Contains(jobLocText, strings.ToLower(loc)) {
				breakdown.LocationMatch = scoreLocation
				break
			}
		}
	}

	// +15: Salary match — job salary range overlaps with user's expectations
	if cfg.SalaryMin > 0 && (job.SalaryMin > 0 || job.SalaryMax > 0) {
		jobMin := job.SalaryMin
		jobMax := job.SalaryMax
		if jobMax == 0 {
			jobMax = jobMin * 1.3
		}
		if cfg.SalaryMax == 0 {
			cfg.SalaryMax = cfg.SalaryMin * 1.5
		}
		// Ranges overlap
		if jobMax >= cfg.SalaryMin && jobMin <= cfg.SalaryMax {
			breakdown.SalaryMatch = scoreSalary
		}
	}

	// +10: Remote preference match
	if cfg.IsRemote && (job.IsRemote || job.IsHybrid) {
		breakdown.RemoteMatch = scoreRemote
	}

	// +10: Experience match — user's experience fits within job range
	if cfg.ExperienceYears > 0 {
		minOk := job.ExperienceMin == 0 || cfg.ExperienceYears >= job.ExperienceMin
		maxOk := job.ExperienceMax == 0 || cfg.ExperienceYears <= job.ExperienceMax+2
		if minOk && maxOk {
			breakdown.ExperienceMatch = scoreExperience
		}
	}

	// +10: Preferred company match
	if len(cfg.PreferredCompanies) > 0 {
		// Get company from context (passed via job.CompanyID lookup or job title context)
		// Here we check if any preferred company keyword is in the job description
		for _, co := range cfg.PreferredCompanies {
			if strings.Contains(jobText, strings.ToLower(co)) {
				breakdown.CompanyMatch = scoreCompany
				break
			}
		}
	}

	breakdown.Total = breakdown.KeywordMatch + breakdown.SkillMatch +
		breakdown.LocationMatch + breakdown.SalaryMatch +
		breakdown.RemoteMatch + breakdown.ExperienceMatch + breakdown.CompanyMatch

	if breakdown.MaxPossible > 0 {
		breakdown.Percentage = (breakdown.Total * 100) / breakdown.MaxPossible
	}

	return breakdown
}

func toLower(ss []string) []string {
	result := make([]string, len(ss))
	for i, s := range ss {
		result[i] = strings.ToLower(s)
	}
	return result
}
