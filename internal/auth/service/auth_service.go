package service

import (
	"cloud-kitchen/internal/auth/model"
	"cloud-kitchen/internal/auth/repository"
	"cloud-kitchen/pkg/util"
	"context"
	"errors"

	"github.com/google/uuid"
)

type AuthService interface {
	Signup(req *model.SignupRequest, ctx context.Context) (*model.User, string, error)
	Login(req *model.LoginRequest, ctx context.Context) (*model.User, string, error)
}

type authService struct {
	repo repository.AuthRepository
}

func NewAuthService(repo repository.AuthRepository) AuthService {
	return &authService{
		repo: repo,
	}
}

func (s *authService) Signup(req *model.SignupRequest, ctx context.Context) (*model.User, string, error) {
	util.Info(ctx, "service.Signup start email=%s", req.Email)
	name := req.Name
	email := req.Email
	password := req.Password

	existingUser, _ := s.repo.GetUserByEmail(ctx, req.Email)

	if existingUser != nil {
		util.Error(ctx, "service.Signup user already exists email=%s", req.Email)
		return nil, "", errors.New("user already exists")
	}

	hashedPassword, err := util.HashPassword(password)
	if err != nil {
		util.Error(ctx, "service.Signup HashPassword failed: %v", err)
		return nil, "", err
	}

	user := &model.User{
		ID:       uuid.New().String(),
		Name:     name,
		Email:    email,
		Password: hashedPassword,
		Provider: "local",
	}

	if err = s.repo.CreateUser(ctx, user); err != nil {
		util.Error(ctx, "service.Signup CreateUser failed: %v", err)
		return nil, "", err
	}

	util.Info(ctx, "service.Signup created user id=%s", user.ID)
	token, err := util.GenerateJWT(user.ID)
	if err != nil {
		util.Error(ctx, "service.Signup GenerateJWT failed: %v", err)
		return nil, "", err
	}

	util.Info(ctx, "service.Signup success user=%s tokenLen=%d", user.ID, len(token))
	return user, token, nil
}

func (s *authService) Login(req *model.LoginRequest, ctx context.Context) (*model.User, string, error) {
	util.Info(ctx, "service.Login start email=%s", req.Email)
	email := req.Email
	password := req.Password

	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		util.Error(ctx, "service.Login repo error: %v", err)
		return nil, "", err
	}

	if user == nil {
		util.Error(ctx, "service.Login user not found email=%s", email)
		return nil, "", errors.New("invalid credentials")
	}

	err = util.CheckPassword(password, user.Password)
	if err != nil {
		util.Error(ctx, "service.Login password check failed: %v", err)
		return nil, "", errors.New("invalid credentials")
	}

	token, err := util.GenerateJWT(user.ID)
	if err != nil {
		util.Error(ctx, "service.Login GenerateJWT failed: %v", err)
		return nil, "", err
	}

	util.Info(ctx, "service.Login success user=%s tokenLen=%d", user.ID, len(token))
	return user, token, nil
}
