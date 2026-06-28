package job

import "context"

type Repository interface {
	Create(ctx context.Context, job *Job) error
	FindByID(ctx context.Context, id uint) (*Job, error)
	FindByDedupHash(ctx context.Context, hash string) (*Job, error)
	List(ctx context.Context, filter *JobFilter) ([]*Job, int64, error)
	Update(ctx context.Context, job *Job) error
	AddSource(ctx context.Context, source *JobSource) error
	UpsertWithSource(ctx context.Context, job *Job, source *JobSource) error
	CountByDate(ctx context.Context, days int) ([]map[string]interface{}, error)
	TopSkills(ctx context.Context, limit int) ([]map[string]interface{}, error)
	TodayCount(ctx context.Context) (int64, error)
}

type Service interface {
	GetJob(ctx context.Context, id uint) (*JobResponse, error)
	ListJobs(ctx context.Context, filter *JobFilter, userID uint) ([]*JobResponse, int64, error)
	UpsertFromProvider(ctx context.Context, normalized *NormalizedJob) error
	GetMatchedJobs(ctx context.Context, userID uint, filter *JobFilter) ([]*JobResponse, int64, error)
	TodayCount(ctx context.Context) (int64, error)
}
