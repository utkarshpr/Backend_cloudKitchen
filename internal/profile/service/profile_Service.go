package service

import (
	"cloud-kitchen/internal/auth/model"
	profilemodel "cloud-kitchen/internal/profile/model"
	"cloud-kitchen/internal/profile/repository"
	"context"

	"github.com/google/uuid"
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

func (s *ProfileService) UpdateProfile(ctx context.Context, userID string, req *profilemodel.UpdateProfileRequest) error {

	// update user fields
	err := s.repo.UpdateUserProfile(ctx, userID, req)
	if err != nil {
		return err
	}

	for _, addr := range req.Addresses {

		if addr.ID == nil {
			// create new address
			newAddr := &model.AddressModel{
				ID:        uuid.New().String(),
				UserID:    userID,
				Label:     addr.Label,
				Street:    addr.Street,
				City:      addr.City,
				State:     addr.State,
				ZipCode:   addr.ZipCode,
				Latitude:  addr.Latitude,
				Longitude: addr.Longitude,
				IsDefault: addr.IsDefault,
			}

			err = s.repo.CreateAddress(ctx, newAddr)

		} else {
			// update existing address
			updateAddr := &model.AddressModel{
				ID:        *addr.ID,
				UserID:    userID,
				Label:     addr.Label,
				Street:    addr.Street,
				City:      addr.City,
				State:     addr.State,
				ZipCode:   addr.ZipCode,
				Latitude:  addr.Latitude,
				Longitude: addr.Longitude,
				IsDefault: addr.IsDefault,
			}

			err = s.repo.UpdateAddress(ctx, updateAddr)
		}

		if err != nil {
			return err
		}
	}

	return nil
}





func (s *ProfileService) DeleteProfile(ctx context.Context, userID string) error {
	return s.repo.DeleteProfile(ctx, userID)
}