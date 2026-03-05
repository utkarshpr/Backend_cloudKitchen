package service

import (
	"cloud-kitchen/internal/auth/model"
	"cloud-kitchen/internal/auth/repository"
	"cloud-kitchen/pkg/util"
	"context"
	"errors"
	"os"

	"github.com/google/uuid"
	"google.golang.org/api/idtoken"
)

type AuthService interface {
	Signup(req *model.SignupRequest, ctx context.Context) (*model.User, string, string, error)
	Login(req *model.LoginRequest, ctx context.Context) (*model.User, string, string, error)
	GoogleLogin(ctx context.Context, idToken string) (*model.User, string, string, error)
}

func (s *authService) GoogleLogin(ctx context.Context, idToken string) (*model.User, string, string, error) {
	util.Info(ctx, "service.GoogleLogin start")

	var GoogleClientID = os.Getenv("GOOGLE_CLIENT_ID")
	payload, err := idtoken.Validate(ctx, idToken, GoogleClientID)
	if err != nil {
		util.Error(ctx, "service.GoogleLogin token validation failed: %v", err)
		return nil, "", "", errors.New("invalid google token")
	}

	email := payload.Claims["email"].(string)
	name := payload.Claims["name"].(string)
	util.Info(ctx, "service.GoogleLogin token validated email=%s name=%s", email, name)

	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		util.Error(ctx, "service.GoogleLogin GetUserByEmail error: %v", err)
		return nil, "", "", err
	}

	if user == nil {
		util.Info(ctx, "service.GoogleLogin user not found, creating new user email=%s", email)

		user = &model.User{
			ID:       uuid.New().String(),
			Name:     name,
			Email:    email,
			Provider: "google",
		}

		err := s.repo.CreateUser(ctx, user)
		if err != nil {
			util.Error(ctx, "service.GoogleLogin CreateUser failed: %v", err)
			return nil, "", "", err
		}
		util.Info(ctx, "service.GoogleLogin user created id=%s", user.ID)
	} else {
		util.Info(ctx, "service.GoogleLogin existing user found id=%s", user.ID)
	}

	accessToken, err := util.GenerateAccessToken(user.ID)
	if err != nil {
		util.Error(ctx, "service.GoogleLogin GenerateAccessToken failed: %v", err)
		return nil, "", "", err
	}

	refreshToken, err := util.GenerateRefreshToken(user.ID)
	if err != nil {
		util.Error(ctx, "service.GoogleLogin GenerateRefreshToken failed: %v", err)
		return nil, "", "", err
	}

	util.Info(ctx, "service.GoogleLogin success user=%s", user.ID)
	return user, accessToken, refreshToken, nil
}

type authService struct {
	repo repository.AuthRepository
}

func NewAuthService(repo repository.AuthRepository) AuthService {
	return &authService{
		repo: repo,
	}
}

func (s *authService) Signup(req *model.SignupRequest, ctx context.Context) (*model.User, string, string, error) {
	util.Info(ctx, "service.Signup start email=%s", req.Email)
	name := req.Name
	email := req.Email
	password := req.Password

	existingUser, _ := s.repo.GetUserByEmail(ctx, req.Email)

	if existingUser != nil {
		util.Error(ctx, "service.Signup user already exists email=%s", req.Email)
		return nil, "", "", errors.New("user already exists")
	}

	hashedPassword, err := util.HashPassword(password)
	if err != nil {
		util.Error(ctx, "service.Signup HashPassword failed: %v", err)
		return nil, "", "", err
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
		return nil, "", "", err
	}

	accessToken, err := util.GenerateAccessToken(user.ID)
	if err != nil {
		util.Error(ctx, "service.Signup GenerateAccessToken failed: %v", err)
		return nil, "", "", err
	}

	refreshToken, err := util.GenerateRefreshToken(user.ID)
	if err != nil {
		util.Error(ctx, "service.Signup GenerateRefreshToken failed: %v", err)
		return nil, "", "", err
	}

	util.Info(ctx, "service.Signup success user=%s", user.ID)
	return user, accessToken, refreshToken, nil
}

func (s *authService) Login(req *model.LoginRequest, ctx context.Context) (*model.User, string, string, error) {
	util.Info(ctx, "service.Login start email=%s", req.Email)
	email := req.Email
	password := req.Password

	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		util.Error(ctx, "service.Login repo error: %v", err)
		return nil, "", "", err
	}

	if user == nil {
		util.Error(ctx, "service.Login user not found email=%s", email)
		return nil, "", "", errors.New("invalid credentials")
	}

	err = util.CheckPassword(password, user.Password)
	if err != nil {
		util.Error(ctx, "service.Login password check failed: %v", err)
		return nil, "", "", errors.New("invalid credentials")
	}

	accessToken, err := util.GenerateAccessToken(user.ID)
	if err != nil {
		util.Error(ctx, "service.Login GenerateAccessToken failed: %v", err)
		return nil, "", "", err
	}

	refreshToken, err := util.GenerateRefreshToken(user.ID)
	if err != nil {
		util.Error(ctx, "service.Login GenerateRefreshToken failed: %v", err)
		return nil, "", "", err
	}

	util.Info(ctx, "service.Login success user=%s", user.ID)
	return user, accessToken, refreshToken, nil
}
