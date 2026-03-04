package controller

import (
	"cloud-kitchen/internal/auth/model"
	"cloud-kitchen/internal/auth/service"
	"cloud-kitchen/pkg/constants"
	"cloud-kitchen/pkg/util"
	"encoding/json"
	"net/http"
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

	user, token, err := a.service.Signup(&req, ctx)
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
		User:  userResponse,
		Token: token,
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

	user, token, err := a.service.Login(&req, ctx)
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
		User:  userResponse,
		Token: token,
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
