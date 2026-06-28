package resume

import "context"

type Repository interface {
	Upsert(ctx context.Context, resume *Resume) error
	FindByUserID(ctx context.Context, userID uint) (*Resume, error)
}

type Service interface {
	Submit(ctx context.Context, userID uint, req *SubmitResumeRequest) (*ResumeResponse, error)
	Get(ctx context.Context, userID uint) (*ResumeResponse, error)
}
