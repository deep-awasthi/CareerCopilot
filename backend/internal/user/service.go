package user

import (
	"context"
)

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetProfile(ctx context.Context, userID uint) (*ProfileResponse, error) {
	profile, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		// Return empty profile if not found yet
		return ToProfileResponse(&Profile{UserID: userID}), nil
	}
	return ToProfileResponse(profile), nil
}

func (s *service) UpdateProfile(ctx context.Context, userID uint, req *UpdateProfileRequest) (*ProfileResponse, error) {
	profile, err := s.repo.Upsert(ctx, userID, req)
	if err != nil {
		return nil, err
	}
	return ToProfileResponse(profile), nil
}
