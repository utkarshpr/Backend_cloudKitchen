package repository

import (
	profilemodel "cloud-kitchen/internal/profile/model"
	"context"
)

type ProfileRepositoryInterface interface {
	GetProfile(ctx context.Context, userID int) (*profilemodel.Profile, error)
	UpdateProfile(ctx context.Context, profile *profilemodel.Profile) error
	DeleteProfile(ctx context.Context, userID int) error
}
