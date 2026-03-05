package repository

import (
	"cloud-kitchen/internal/auth/model"
	profilemodel "cloud-kitchen/internal/profile/model"
	"context"
)

type ProfileRepositoryInterface interface {
	GetProfile(ctx context.Context, userID int) (*profilemodel.Profile, error)
	UpdateProfile(ctx context.Context, profile *profilemodel.Profile) error
	DeleteProfile(ctx context.Context, userID int) error
	UpdateUserProfile(ctx context.Context, userID string, req *profilemodel.UpdateProfileRequest) error
	CreateAddress(ctx context.Context, address *model.AddressModel) error
	UpdateAddress(ctx context.Context, address *model.AddressModel) error
	DeleteAddress(ctx context.Context, addressID string) error
	
}
