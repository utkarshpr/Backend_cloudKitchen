package repository

import "cloud-kitchen/internal/auth/model"

type AuthRepository interface {
	CreateUser(user *model.User) error
	GetUserByEmail(email string) (*model.User, error)
}
