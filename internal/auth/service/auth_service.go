package service

import (
	"cloud-kitchen/internal/auth/model"
	"cloud-kitchen/internal/auth/repository"
	"cloud-kitchen/pkg/util"
	"errors"

	"github.com/google/uuid"
)

type AuthService interface {
	Signup(name string, email string, password string) (*model.User, string, error)
	Login(email string, password string) (*model.User, string, error)
}

type authService struct {
	repo repository.AuthRepository
}

func NewAuthService(repo repository.AuthRepository) AuthService {
	return &authService{
		repo: repo,
	}
}

func (s *authService) Signup(name string, email string, password string) (*model.User, string, error) {

	existingUser, _ := s.repo.GetUserByEmail(email)

	if existingUser != nil {
		return nil, "", errors.New("user already exists")
	}

	hashedPassword, err := util.HashPassword(password)
	if err != nil {
		return nil, "", err
	}

	user := &model.User{
		ID:       uuid.New().String(),
		Name:     name,
		Email:    email,
		Password: hashedPassword,
		Provider: "local",
	}

	err = s.repo.CreateUser(user)
	if err != nil {
		return nil, "", err
	}

	token, err := util.GenerateJWT(user.ID)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *authService) Login(email string, password string) (*model.User, string, error) {

	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return nil, "", err
	}

	if user == nil {
		return nil, "", errors.New("invalid credentials")
	}

	err = util.CheckPassword(password, user.Password)
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	token, err := util.GenerateJWT(user.ID)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}
