package controller

import (
	"cloud-kitchen/internal/auth/model"
	"cloud-kitchen/internal/auth/service"
	"cloud-kitchen/pkg/constants"
	"cloud-kitchen/pkg/util"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthController struct {
	service service.AuthService
}



func NewAuthController(service service.AuthService) *AuthController {
	return &AuthController{
		service: service,
	}
}

func (a *AuthController) Signup(w http.ResponseWriter, r *http.Request) {

	var req model.SignupRequest
	ctx := r.Context()
	requestID := ctx.Value(constants.RequestIDKey).(string)
	util.Info(ctx, "controller.Signup received request id=%s", requestID)

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		util.Error(ctx, "controller.Signup invalid body: %v", err)

		resp := util.APIResponse{
			RequestID: requestID,
			Success:   false,
			Message:   constants.InvalidRequestBody,
			ErrorCode: constants.ErrInvalidRequest,
			Data:      nil,
		}

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	user, accessToken, refreshToken, err := a.service.Signup(&req, ctx)
	if err != nil {
		util.Error(ctx, "controller.Signup service error: %v", err)

		resp := util.APIResponse{
			RequestID: requestID,
			Success:   false,
			Message:   err.Error(),
			ErrorCode: constants.ErrSignupFailed,
			Data:      nil,
		}

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}
	util.Info(ctx, "controller.Signup service returned user=%s", user.ID)

	userResponse := &model.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Provider:  user.Provider,
		CreatedAt: user.CreatedAt,
	}

	data := &model.SignupResponse{
		User:         userResponse,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	resp := util.APIResponse{
		RequestID: requestID,
		Success:   true,
		Message:   constants.SignupSuccess,
		Data:      data,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
	util.Info(ctx, "controller.Signup response sent for user=%s", user.ID)
}

func (a *AuthController) Login(w http.ResponseWriter, r *http.Request) {

	var req model.LoginRequest
	ctx := r.Context()
	requestID := ctx.Value(constants.RequestIDKey).(string)
	util.Info(ctx, "controller.Login received request id=%s", requestID)

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		util.Error(ctx, "controller.Login invalid body: %v", err)

		resp := util.APIResponse{
			RequestID: requestID,
			Success:   false,
			Message:   constants.InvalidRequestBody,
			ErrorCode: constants.ErrInvalidRequest,
			Data:      nil,
		}

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	user, accessToken, refreshToken, err := a.service.Login(&req, ctx)
	if err != nil {
		util.Error(ctx, "controller.Login service error: %v", err)

		resp := util.APIResponse{
			RequestID: requestID,
			Success:   false,
			Message:   constants.InvalidCredentials,
			ErrorCode: constants.ErrInvalidCredsCode,
			Data:      nil,
		}

		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(resp)
		return
	}
	util.Info(ctx, "controller.Login service returned user=%s", user.ID)

	userResponse := &model.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Provider:  user.Provider,
		CreatedAt: user.CreatedAt,
	}

	data := &model.LoginResponse{
		User:         userResponse,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	resp := util.APIResponse{
		RequestID: requestID,
		Success:   true,
		Message:   constants.LoginSuccess,
		Data:      data,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
	util.Info(ctx, "controller.Login response sent for user=%s", user.ID)
}

func (a *AuthController) GoogleLogin(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	requestID := ctx.Value(constants.RequestIDKey).(string)
	util.Info(ctx, "controller.GoogleLogin received request id=%s", requestID)

	var req model.GoogleLoginRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		util.Error(ctx, "controller.GoogleLogin invalid body: %v", err)
		resp := util.APIResponse{
			RequestID: requestID,
			Success:   false,
			Message:   constants.InvalidRequestBody,
			ErrorCode: constants.ErrInvalidRequest,
			Data:      nil,
		}

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)

		return
	}

	user, accessToken, refreshToken, err := a.service.GoogleLogin(ctx, req.IDToken)

	if err != nil {
		util.Error(ctx, "controller.GoogleLogin service error: %v", err)
		resp := util.APIResponse{
			RequestID: requestID,
			Success:   false,
			Message:   err.Error(),
			ErrorCode: constants.ErrGoogleLoginFailed,
			Data:      nil,
		}

		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(resp)
		return
	}

	util.Info(ctx, "controller.GoogleLogin service returned user=%s", user.ID)

	userResponse := &model.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Provider:  user.Provider,
		CreatedAt: user.CreatedAt,
	}

	data := &model.LoginResponse{
		User:         userResponse,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	resp := util.APIResponse{
		RequestID: requestID,
		Success:   true,
		Message:   "google login successful",
		Data:      data,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
	util.Info(ctx, "controller.GoogleLogin response sent for user=%s", user.ID)
}

func (a *AuthController) Refresh(w http.ResponseWriter, r *http.Request) {
	var req *model.RefreshRequest
	ctx := r.Context()
	requestID := ctx.Value(constants.RequestIDKey).(string)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.Error(ctx, "controller.Refresh invalid body: %v", err)
		resp := util.APIResponse{
			RequestID: requestID,
			Success:   false,
			Message:   constants.InvalidRequestBody,
			ErrorCode: constants.ErrInvalidRequest,
			Data:      nil,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	token, err := jwt.Parse(req.RefreshToken, func(t *jwt.Token) (interface{}, error) {
		return util.Secret, nil
	})

	if err != nil || !token.Valid {
		util.Error(ctx, "controller.Refresh invalid refresh token: %v", err)
		resp := util.APIResponse{
			RequestID: requestID,
			Success:   false,
			Message:   "invalid refresh token",
			ErrorCode: constants.ErrInvalidRefreshToken,
			Data:      nil,
		}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(resp)
		return
	}

	claims := token.Claims.(jwt.MapClaims)
	userID, ok := claims["user_id"].(string)
	if !ok {
		util.Error(ctx, "controller.Refresh invalid token claims")
		resp := util.APIResponse{
			RequestID: requestID,
			Success:   false,
			Message:   "invalid token claims",
			ErrorCode: constants.ErrInvalidTokenClaims,
			Data:      nil,
		}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(resp)
		return
	}

	newAccessToken, err := util.GenerateAccessToken(userID)
	if err != nil {
		util.Error(ctx, "controller.Refresh token generation failed: %v", err)
		resp := util.APIResponse{
			RequestID: requestID,
			Success:   false,
			Message:   "token generation failed",
			ErrorCode: constants.ErrTokenGenerationFailed,
			Data:      nil,
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp := util.APIResponse{
		RequestID: requestID,
		Success:   true,
		Message:   "token refreshed successfully",
		Data:      gin.H{"access_token": newAccessToken},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
