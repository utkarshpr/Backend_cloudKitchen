package service

import (
	profilemodel "cloud-kitchen/internal/profile/model"
	"cloud-kitchen/internal/profile/repository"
	"context"
)

type ProfileService struct {
	repo *repository.ProfileRepository
}

func NewProfileService(repo *repository.ProfileRepository) *ProfileService {
	return &ProfileService{repo: repo}
}

// Implement methods for profile service (e.g., GetProfile, UpdateProfile, DeleteProfile)
func (s *ProfileService) GetProfile(ctx context.Context, email string) (*profilemodel.Profile, error) {
	return s.repo.GetProfile(ctx, email)
}

func (s *ProfileService) UpdateProfile(ctx context.Context, profile *profilemodel.Profile) error {
	return s.repo.UpdateProfile(ctx, profile)
}

func (s *ProfileService) DeleteProfile(ctx context.Context, userID int) error {
	return s.repo.DeleteProfile(ctx, userID)
}
