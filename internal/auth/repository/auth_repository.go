package repository

import (
	"context"

	"cloud-kitchen/internal/auth/model"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	CreateAddress(ctx context.Context, address *model.AddressModel) error
	GetAddressesByUserID(ctx context.Context, userID string) ([]model.Address, error)
}

