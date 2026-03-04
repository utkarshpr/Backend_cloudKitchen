package repository

import (
	"context"

	"cloud-kitchen/internal/auth/model"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
}
