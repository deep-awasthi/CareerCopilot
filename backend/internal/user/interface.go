package user

import "context"

type Repository interface {
	Create(ctx context.Context, profile *Profile) error
	FindByUserID(ctx context.Context, userID uint) (*Profile, error)
	Update(ctx context.Context, profile *Profile) error
	Upsert(ctx context.Context, userID uint, req *UpdateProfileRequest) (*Profile, error)
}

type Service interface {
	GetProfile(ctx context.Context, userID uint) (*ProfileResponse, error)
	UpdateProfile(ctx context.Context, userID uint, req *UpdateProfileRequest) (*ProfileResponse, error)
}
