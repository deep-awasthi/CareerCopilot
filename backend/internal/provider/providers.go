package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/deepawasthi/careercopilot/internal/job"
)

// GreenhouseProvider scrapes Greenhouse-hosted job boards via their public API
type GreenhouseProvider struct {
	BaseProvider
	client *http.Client
}

func NewGreenhouseProvider() *GreenhouseProvider {
	return &GreenhouseProvider{
		BaseProvider: BaseProvider{MaxRetries: 3, RetryDelay: time.Second},
		client:       &http.Client{Timeout: 30 * time.Second},
	}
}

func (g *GreenhouseProvider) Name() string { return "greenhouse" }

func (g *GreenhouseProvider) Search(ctx context.Context, params SearchParams) ([]*job.NormalizedJob, error) {
	// Greenhouse public job board API: https://boards-api.greenhouse.io/v1/boards/{company}/jobs
	// We aggregate across well-known companies using Greenhouse
	knownBoards := []string{
		"airbnb", "stripe", "twilio", "datadog", "okta", "confluent",
		"figma", "notion", "linear", "vercel", "hashicorp", "grafana",
	}

	var allJobs []*job.NormalizedJob
	for _, board := range knownBoards {
		boardJobs, err := g.fetchBoard(ctx, board, params)
		if err != nil {
			continue // Skip failed boards, don't abort
		}
		allJobs = append(allJobs, boardJobs...)
	}
	return allJobs, nil
}

func (g *GreenhouseProvider) fetchBoard(ctx context.Context, company string, params SearchParams) ([]*job.NormalizedJob, error) {
	apiURL := fmt.Sprintf("https://boards-api.greenhouse.io/v1/boards/%s/jobs?content=true", company)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "CareerCopilot/1.0")

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("greenhouse board %s returned %d", company, resp.StatusCode)
	}

	var raw struct {
		Jobs []struct {
			ID         int    `json:"id"`
			Title      string `json:"title"`
			Location   struct {
				Name string `json:"name"`
			} `json:"location"`
			AbsoluteURL string `json:"absolute_url"`
			UpdatedAt   string `json:"updated_at"`
			Content     string `json:"content"`
			Metadata    []struct {
				Name  string `json:"name"`
				Value interface{} `json:"value"`
			} `json:"metadata"`
		} `json:"jobs"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	var normalized []*job.NormalizedJob
	for _, j := range raw.Jobs {
		if !matchesParams(j.Title, params) {
			continue
		}

		nj := &job.NormalizedJob{
			ExternalID:     fmt.Sprintf("greenhouse-%d", j.ID),
			Provider:       g.Name(),
			Company:        company,
			Title:          j.Title,
			Description:    j.Content,
			Location:       j.Location.Name,
			ApplicationURL: j.AbsoluteURL,
			SourceURL:      j.AbsoluteURL,
			PostedAt:       j.UpdatedAt,
			IsRemote:       strings.Contains(strings.ToLower(j.Location.Name), "remote"),
			SalaryCurrency: "USD",
		}
		normalized = append(normalized, nj)
	}
	return normalized, nil
}

func (g *GreenhouseProvider) FetchJob(ctx context.Context, externalID string) (*job.NormalizedJob, error) {
	return nil, fmt.Errorf("greenhouse: FetchJob not implemented for ID %s", externalID)
}

func (g *GreenhouseProvider) IsAvailable(ctx context.Context) bool {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		"https://boards-api.greenhouse.io/v1/boards/stripe/jobs", nil)
	if err != nil {
		return false
	}
	resp, err := g.client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

// LeverProvider scrapes Lever-hosted job boards
type LeverProvider struct {
	BaseProvider
	client *http.Client
}

func NewLeverProvider() *LeverProvider {
	return &LeverProvider{
		BaseProvider: BaseProvider{MaxRetries: 3, RetryDelay: time.Second},
		client:       &http.Client{Timeout: 30 * time.Second},
	}
}

func (l *LeverProvider) Name() string { return "lever" }

func (l *LeverProvider) Search(ctx context.Context, params SearchParams) ([]*job.NormalizedJob, error) {
	knownBoards := []string{
		"netflix", "shopify", "reddit", "lyft", "cloudflare",
		"mongodb", "elastic", "databricks", "atlassian",
	}

	var allJobs []*job.NormalizedJob
	for _, company := range knownBoards {
		jobs, err := l.fetchBoard(ctx, company, params)
		if err != nil {
			continue
		}
		allJobs = append(allJobs, jobs...)
	}
	return allJobs, nil
}

func (l *LeverProvider) fetchBoard(ctx context.Context, company string, params SearchParams) ([]*job.NormalizedJob, error) {
	apiURL := fmt.Sprintf("https://api.lever.co/v0/postings/%s?mode=json", company)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "CareerCopilot/1.0")

	resp, err := l.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("lever %s returned %d", company, resp.StatusCode)
	}

	var raw []struct {
		ID          string `json:"id"`
		Text        string `json:"text"`
		Categories  struct {
			Location   string `json:"location"`
			Team       string `json:"team"`
			Commitment string `json:"commitment"`
		} `json:"categories"`
		HostedURL   string `json:"hostedUrl"`
		ApplyURL    string `json:"applyUrl"`
		CreatedAt   int64  `json:"createdAt"`
		Description struct {
			Description string `json:"description"`
			Lists       []struct {
				Text  string   `json:"text"`
				Content []string `json:"content"`
			} `json:"lists"`
		} `json:"descriptionPlain"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	var normalized []*job.NormalizedJob
	for _, j := range raw {
		if !matchesParams(j.Text, params) {
			continue
		}
		postedAt := time.Unix(j.CreatedAt/1000, 0).Format(time.RFC3339)
		nj := &job.NormalizedJob{
			ExternalID:     "lever-" + j.ID,
			Provider:       l.Name(),
			Company:        company,
			Title:          j.Text,
			Location:       j.Categories.Location,
			ApplicationURL: j.ApplyURL,
			SourceURL:      j.HostedURL,
			PostedAt:       postedAt,
			IsRemote:       strings.Contains(strings.ToLower(j.Categories.Location), "remote"),
			SalaryCurrency: "USD",
		}
		normalized = append(normalized, nj)
	}
	return normalized, nil
}

func (l *LeverProvider) FetchJob(ctx context.Context, externalID string) (*job.NormalizedJob, error) {
	return nil, nil
}

func (l *LeverProvider) IsAvailable(ctx context.Context) bool { return true }

// WorkdayProvider fetches jobs from Workday career pages (company-specific)
type WorkdayProvider struct {
	BaseProvider
	client *http.Client
}

func NewWorkdayProvider() *WorkdayProvider {
	return &WorkdayProvider{
		BaseProvider: BaseProvider{MaxRetries: 2, RetryDelay: 2 * time.Second},
		client:       &http.Client{Timeout: 30 * time.Second},
	}
}

func (w *WorkdayProvider) Name() string { return "workday" }

func (w *WorkdayProvider) Search(ctx context.Context, params SearchParams) ([]*job.NormalizedJob, error) {
	// Workday search API endpoints for common companies
	endpoints := []struct {
		Company string
		URL     string
	}{
		{"Amazon", "https://www.amazon.jobs/en/search.json?base_query=%s&loc_query=&job_count=10&result_limit=10&sort=relevant&category%5B%5D=software-development"},
		{"Google", "https://careers.google.com/api/jobs/jobs-v1/search/?q=%s&page=1"},
	}

	var allJobs []*job.NormalizedJob
	keyword := ""
	if len(params.Keywords) > 0 {
		keyword = url.QueryEscape(strings.Join(params.Keywords, " "))
	}

	for _, ep := range endpoints {
		apiURL := fmt.Sprintf(ep.URL, keyword)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
		if err != nil {
			continue
		}
		req.Header.Set("User-Agent", "Mozilla/5.0 CareerCopilot")
		resp, err := w.client.Do(req)
		if err != nil || resp.StatusCode != 200 {
			if resp != nil {
				resp.Body.Close()
			}
			continue
		}
		resp.Body.Close()
		// Note: actual parsing depends on specific company API format
		// Placeholder returns empty — Playwright scrapers handle these
		_ = allJobs
	}

	return allJobs, nil
}

func (w *WorkdayProvider) FetchJob(ctx context.Context, externalID string) (*job.NormalizedJob, error) {
	return nil, nil
}

func (w *WorkdayProvider) IsAvailable(ctx context.Context) bool { return true }

// IndeedProvider fetches from Indeed public API
type IndeedProvider struct {
	BaseProvider
	client *http.Client
}

func NewIndeedProvider() *IndeedProvider {
	return &IndeedProvider{
		BaseProvider: BaseProvider{MaxRetries: 3, RetryDelay: 2 * time.Second},
		client:       &http.Client{Timeout: 30 * time.Second},
	}
}

func (i *IndeedProvider) Name() string { return "indeed" }

func (i *IndeedProvider) Search(ctx context.Context, params SearchParams) ([]*job.NormalizedJob, error) {
	// Indeed requires a publisher API key; stub returns empty
	// Playwright-based scraper handles actual scraping
	return []*job.NormalizedJob{}, nil
}

func (i *IndeedProvider) FetchJob(ctx context.Context, externalID string) (*job.NormalizedJob, error) {
	return nil, nil
}

func (i *IndeedProvider) IsAvailable(ctx context.Context) bool { return true }

// LinkedInProvider fetches from LinkedIn Jobs
type LinkedInProvider struct {
	BaseProvider
	client *http.Client
}

func NewLinkedInProvider() *LinkedInProvider {
	return &LinkedInProvider{
		BaseProvider: BaseProvider{MaxRetries: 3, RetryDelay: 3 * time.Second},
		client:       &http.Client{Timeout: 30 * time.Second},
	}
}

func (l *LinkedInProvider) Name() string { return "linkedin" }

func (l *LinkedInProvider) Search(ctx context.Context, params SearchParams) ([]*job.NormalizedJob, error) {
	// LinkedIn scraping requires Playwright — stub returns empty
	// The PlaywrightProvider handles actual browser-based scraping
	return []*job.NormalizedJob{}, nil
}

func (l *LinkedInProvider) FetchJob(ctx context.Context, externalID string) (*job.NormalizedJob, error) {
	return nil, nil
}

func (l *LinkedInProvider) IsAvailable(ctx context.Context) bool { return true }

// WellfoundProvider fetches startup jobs from Wellfound (AngelList)
type WellfoundProvider struct {
	BaseProvider
	client *http.Client
}

func NewWellfoundProvider() *WellfoundProvider {
	return &WellfoundProvider{
		BaseProvider: BaseProvider{MaxRetries: 3, RetryDelay: time.Second},
		client:       &http.Client{Timeout: 30 * time.Second},
	}
}

func (w *WellfoundProvider) Name() string { return "wellfound" }

func (w *WellfoundProvider) Search(ctx context.Context, params SearchParams) ([]*job.NormalizedJob, error) {
	return []*job.NormalizedJob{}, nil
}

func (w *WellfoundProvider) FetchJob(ctx context.Context, externalID string) (*job.NormalizedJob, error) {
	return nil, nil
}

func (w *WellfoundProvider) IsAvailable(ctx context.Context) bool { return true }

// NaukriProvider fetches from Naukri (Indian job portal)
type NaukriProvider struct {
	BaseProvider
	client *http.Client
}

func NewNaukriProvider() *NaukriProvider {
	return &NaukriProvider{
		BaseProvider: BaseProvider{MaxRetries: 3, RetryDelay: 2 * time.Second},
		client:       &http.Client{Timeout: 30 * time.Second},
	}
}

func (n *NaukriProvider) Name() string { return "naukri" }

func (n *NaukriProvider) Search(ctx context.Context, params SearchParams) ([]*job.NormalizedJob, error) {
	return []*job.NormalizedJob{}, nil
}

func (n *NaukriProvider) FetchJob(ctx context.Context, externalID string) (*job.NormalizedJob, error) {
	return nil, nil
}

func (n *NaukriProvider) IsAvailable(ctx context.Context) bool { return true }

// CareerPageProvider monitors company career pages for new postings
type CareerPageProvider struct {
	BaseProvider
	client *http.Client
}

func NewCareerPageProvider() *CareerPageProvider {
	return &CareerPageProvider{
		BaseProvider: BaseProvider{MaxRetries: 2, RetryDelay: 5 * time.Second},
		client:       &http.Client{Timeout: 60 * time.Second},
	}
}

func (c *CareerPageProvider) Name() string { return "career_page" }

func (c *CareerPageProvider) Search(ctx context.Context, params SearchParams) ([]*job.NormalizedJob, error) {
	return []*job.NormalizedJob{}, nil
}

func (c *CareerPageProvider) FetchJob(ctx context.Context, externalID string) (*job.NormalizedJob, error) {
	return nil, nil
}

func (c *CareerPageProvider) IsAvailable(ctx context.Context) bool { return true }

// Helper: check if job title matches search params
func matchesParams(title string, params SearchParams) bool {
	if len(params.Keywords) == 0 {
		return true
	}
	titleLower := strings.ToLower(title)
	for _, kw := range params.Keywords {
		if strings.Contains(titleLower, strings.ToLower(kw)) {
			return true
		}
	}
	return false
}

// NewDefaultRegistry creates a registry with all default providers
func NewDefaultRegistry() *Registry {
	reg := NewRegistry()
	reg.Register(NewGreenhouseProvider())
	reg.Register(NewLeverProvider())
	reg.Register(NewWorkdayProvider())
	reg.Register(NewIndeedProvider())
	reg.Register(NewLinkedInProvider())
	reg.Register(NewWellfoundProvider())
	reg.Register(NewNaukriProvider())
	reg.Register(NewCareerPageProvider())
	return reg
}
