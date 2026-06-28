package job

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

// DeduplicationEngine merges duplicate jobs from multiple providers
type DeduplicationEngine struct{}

func NewDeduplicationEngine() *DeduplicationEngine {
	return &DeduplicationEngine{}
}

// ComputeHash generates a deduplication hash for a job
// Based on: company + normalized title + primary location
func (d *DeduplicationEngine) ComputeHash(companyName, title, location string) string {
	normalized := strings.ToLower(strings.TrimSpace(companyName)) +
		"|" + normalizeTitle(title) +
		"|" + normalizeLocation(location)
	h := sha256.Sum256([]byte(normalized))
	return fmt.Sprintf("%x", h)[:16]
}

// AreDuplicates checks if two jobs are duplicates using multiple signals
func (d *DeduplicationEngine) AreDuplicates(a, b *NormalizedJob) bool {
	// Same dedup hash is definitive
	if a.DedupHash != "" && a.DedupHash == b.DedupHash {
		return true
	}
	// Same company + same title (normalized) + same location
	if sameCompany(a.Company, b.Company) &&
		titlesMatch(a.Title, b.Title) &&
		locationsMatch(a.Location, b.Location) {
		return true
	}
	// Same application URL
	if a.ApplicationURL != "" && a.ApplicationURL == b.ApplicationURL {
		return true
	}
	return false
}

// NormalizedJob is the common schema for all providers
type NormalizedJob struct {
	ExternalID     string
	Provider       string
	Company        string
	Title          string
	Description    string
	Location       string
	Locations      []string
	IsRemote       bool
	IsHybrid       bool
	EmploymentType string
	ExperienceMin  float64
	ExperienceMax  float64
	SalaryMin      float64
	SalaryMax      float64
	SalaryCurrency string
	Skills         []string
	ApplicationURL string
	SourceURL      string
	PostedAt       string
	DedupHash      string
}

func normalizeTitle(title string) string {
	title = strings.ToLower(title)
	// Remove level indicators for dedup
	replacements := []string{
		"senior ", "sr. ", "sr ", "junior ", "jr. ", "jr ",
		"lead ", "principal ", "staff ", "associate ",
		"i ", "ii ", "iii ", "iv ", "v ",
		"- ", "– ",
	}
	for _, r := range replacements {
		title = strings.ReplaceAll(title, r, "")
	}
	return strings.TrimSpace(title)
}

func normalizeLocation(loc string) string {
	return strings.ToLower(strings.TrimSpace(loc))
}

func sameCompany(a, b string) bool {
	a = strings.ToLower(strings.TrimSpace(a))
	b = strings.ToLower(strings.TrimSpace(b))
	if a == b {
		return true
	}
	// Check if one contains the other
	if strings.Contains(a, b) || strings.Contains(b, a) {
		return true
	}
	return false
}

func titlesMatch(a, b string) bool {
	return normalizeTitle(a) == normalizeTitle(b)
}

func locationsMatch(a, b string) bool {
	a = normalizeLocation(a)
	b = normalizeLocation(b)
	if a == b || a == "" || b == "" {
		return true
	}
	// Handle remote
	if strings.Contains(a, "remote") && strings.Contains(b, "remote") {
		return true
	}
	// Match city names
	aParts := strings.Split(a, ",")
	bParts := strings.Split(b, ",")
	if len(aParts) > 0 && len(bParts) > 0 {
		return strings.TrimSpace(aParts[0]) == strings.TrimSpace(bParts[0])
	}
	return false
}
