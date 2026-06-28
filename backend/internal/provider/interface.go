package provider

import (
	"context"
	"time"

	"github.com/deepawasthi/careercopilot/internal/job"
)

// SearchParams defines what to search for
type SearchParams struct {
	Keywords      []string
	Locations     []string
	ExperienceMin float64
	ExperienceMax float64
	SalaryMin     float64
	IsRemote      bool
	IsHybrid      bool
	JobType       string
	Page          int
	PageSize      int
}

// Provider is the interface every job source must implement
type Provider interface {
	// Name returns the unique provider identifier
	Name() string

	// Search fetches jobs matching the given params (with retry + pagination)
	Search(ctx context.Context, params SearchParams) ([]*job.NormalizedJob, error)

	// FetchJob fetches a single job by external ID
	FetchJob(ctx context.Context, externalID string) (*job.NormalizedJob, error)

	// IsAvailable checks if the provider is reachable
	IsAvailable(ctx context.Context) bool
}

// BaseProvider provides common retry + rate limiting logic
type BaseProvider struct {
	MaxRetries int
	RetryDelay time.Duration
}

func (b *BaseProvider) withRetry(ctx context.Context, fn func() error) error {
	maxRetries := b.MaxRetries
	if maxRetries == 0 {
		maxRetries = 3
	}
	delay := b.RetryDelay
	if delay == 0 {
		delay = 2 * time.Second
	}

	var lastErr error
	for i := 0; i < maxRetries; i++ {
		if err := fn(); err != nil {
			lastErr = err
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay * time.Duration(i+1)):
			}
			continue
		}
		return nil
	}
	return lastErr
}

// Registry holds all registered providers
type Registry struct {
	providers map[string]Provider
}

func NewRegistry() *Registry {
	return &Registry{providers: make(map[string]Provider)}
}

func (r *Registry) Register(p Provider) {
	r.providers[p.Name()] = p
}

func (r *Registry) Get(name string) (Provider, bool) {
	p, ok := r.providers[name]
	return p, ok
}

func (r *Registry) All() []Provider {
	var all []Provider
	for _, p := range r.providers {
		all = append(all, p)
	}
	return all
}
