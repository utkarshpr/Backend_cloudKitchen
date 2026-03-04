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
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {

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

	data := map[string]interface{}{
		"user":  user,
		"token": token,
	}

	resp := util.APIResponse{
		RequestID: requestID,
		Success:   true,
		Message:   constants.SignupSuccess,
		Data:      data,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (a *AuthController) Login(w http.ResponseWriter, r *http.Request) {

	var req model.LoginRequest
	ctx := r.Context()
	requestID := ctx.Value(constants.RequestIDKey).(string)

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {

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

	data := map[string]interface{}{
		"user":  user,
		"token": token,
	}

	resp := util.APIResponse{
		RequestID: requestID,
		Success:   true,
		Message:   constants.LoginSuccess,
		Data:      data,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
