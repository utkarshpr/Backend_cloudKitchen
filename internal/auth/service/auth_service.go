package service

import (
	"cloud-kitchen/internal/auth/model"
	"cloud-kitchen/internal/auth/repository"
	"cloud-kitchen/pkg/util"
	"context"
	"errors"
	"os"
	"time"

	"github.com/google/uuid"
	"google.golang.org/api/idtoken"
)

type AuthService interface {
	Signup(req *model.SignupRequest, ctx context.Context) (*model.User, string, string, error)
	Login(req *model.LoginRequest, ctx context.Context) (*model.User, string, string, error)
	GoogleLogin(ctx context.Context, idToken string) (*model.User, string, string, error)
}

type authService struct {
	repo repository.AuthRepository
}

func NewAuthService(repo repository.AuthRepository) AuthService {
	return &authService{
		repo: repo,
	}
}

func (s *authService) GoogleLogin(ctx context.Context, idToken string) (*model.User, string, string, error) {

	util.Info(ctx, "service.GoogleLogin start")

	clientID := os.Getenv("GOOGLE_CLIENT_ID")

	payload, err := idtoken.Validate(ctx, idToken, clientID)
	if err != nil {
		util.Error(ctx, "service.GoogleLogin token validation failed: %v", err)
		return nil, "", "", errors.New("invalid google token")
	}

	email, ok := payload.Claims["email"].(string)
	if !ok {
		return nil, "", "", errors.New("email not found in token")
	}

	name, _ := payload.Claims["name"].(string)

	util.Info(ctx, "service.GoogleLogin token validated email=%s", email)

	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		util.Error(ctx, "service.GoogleLogin repo error: %v", err)
		return nil, "", "", err
	}

	if user == nil {

		util.Info(ctx, "service.GoogleLogin creating new user email=%s", email)

		user = &model.User{
			ID:       uuid.New().String(),
			Name:     name,
			Email:    email,
			Provider: "google",
		}

		if err := s.repo.CreateUser(ctx, user); err != nil {
			util.Error(ctx, "service.GoogleLogin CreateUser failed: %v", err)
			return nil, "", "", err
		}
	}

	accessToken, refreshToken, err := generateTokens(user.ID)
	if err != nil {
		return nil, "", "", err
	}

	util.Info(ctx, "service.GoogleLogin success user=%s", user.ID)

	return user, accessToken, refreshToken, nil
}

func (s *authService) Signup(req *model.SignupRequest, ctx context.Context) (*model.User, string, string, error) {

	util.Info(ctx, "service.Signup start email=%s", req.Email)

	// check if user exists
	existingUser, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		util.Error(ctx, "service.Signup repo error: %v", err)
		return nil, "", "", err
	}

	if existingUser != nil {
		util.Error(ctx, "service.Signup user already exists email=%s", req.Email)
		return nil, "", "", errors.New("user already exists")
	}

	// hash password
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		util.Error(ctx, "service.Signup HashPassword failed: %v", err)
		return nil, "", "", err
	}

	userID := uuid.New().String()

	user := &model.User{
		ID:             userID,
		Name:           req.Name,
		Email:          req.Email,
		Password:       hashedPassword,
		Provider:       "local",
		MobileNumber:   req.MobileNumber,
		ProfilePicture: req.ProfilePicture,
		CreatedAt:      time.Now(),
	}

	// create user
	if err := s.repo.CreateUser(ctx, user); err != nil {
		util.Error(ctx, "service.Signup CreateUser failed: %v", err)
		return nil, "", "", err
	}

	// create addresses
	var savedAddresses []model.Address

	for _, addr := range req.Addresses {

		address := &model.AddressModel{
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

		err := s.repo.CreateAddress(ctx, address)
		if err != nil {
			util.Error(ctx, "service.Signup CreateAddress failed: %v", err)
			return nil, "", "", err
		}

		savedAddresses = append(savedAddresses, model.Address{
			Label:     addr.Label,
			Street:    addr.Street,
			City:      addr.City,
			State:     addr.State,
			ZipCode:   addr.ZipCode,
			Latitude:  addr.Latitude,
			Longitude: addr.Longitude,
			IsDefault: addr.IsDefault,
		})
	}

	user.Addresses = savedAddresses

	// generate tokens
	accessToken, refreshToken, err := generateTokens(user.ID)
	if err != nil {
		util.Error(ctx, "service.Signup token generation failed: %v", err)
		return nil, "", "", err
	}

	util.Info(ctx, "service.Signup success user=%s", user.ID)

	return user, accessToken, refreshToken, nil
}

func (s *authService) Login(req *model.LoginRequest, ctx context.Context) (*model.User, string, string, error) {

	util.Info(ctx, "service.Login start email=%s", req.Email)

	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		util.Error(ctx, "service.Login repo error: %v", err)
		return nil, "", "", err
	}

	if user == nil {
		util.Error(ctx, "service.Login user not found email=%s", req.Email)
		return nil, "", "", errors.New("invalid credentials")
	}

	if err := util.CheckPassword(req.Password, user.Password); err != nil {
		util.Error(ctx, "service.Login password check failed")
		return nil, "", "", errors.New("invalid credentials")
	}

	accessToken, refreshToken, err := generateTokens(user.ID)
	if err != nil {
		return nil, "", "", err
	}
	addresses, err := s.repo.GetAddressesByUserID(ctx, user.ID)
if err != nil {
	util.Error(ctx, "service.Login GetAddresses failed: %v", err)
	return nil, "", "", err
}

user.Addresses = addresses

	util.Info(ctx, "service.Login success user=%s", user.ID)

	return user, accessToken, refreshToken, nil
}

func generateTokens(userID string) (string, string, error) {

	accessToken, err := util.GenerateAccessToken(userID)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := util.GenerateRefreshToken(userID)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
