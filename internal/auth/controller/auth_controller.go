package controller

import (
	"cloud-kitchen/internal/auth/model"
	"cloud-kitchen/internal/auth/service"
	"cloud-kitchen/pkg/constants"
	"cloud-kitchen/pkg/util"
	"encoding/json"
	"net/http"

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

	ctx := r.Context()
	requestID := ctx.Value(constants.RequestIDKey).(string)

	var req model.SignupRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.Error(ctx, "controller.Signup invalid body: %v", err)

		util.WriteErrorResponse(w, requestID, constants.InvalidRequestBody, constants.ErrInvalidRequest, http.StatusBadRequest)
		return
	}

	user, accessToken, refreshToken, err := a.service.Signup(&req, ctx)
	if err != nil {
		util.Error(ctx, "controller.Signup service error: %v", err)

		util.WriteErrorResponse(w, requestID, err.Error(), constants.ErrSignupFailed, http.StatusBadRequest)
		return
	}

	userResponse := &model.UserResponse{
		ID:             user.ID,
		Name:           user.Name,
		Email:          user.Email,
		Provider:       user.Provider,
		MobileNumber:   user.MobileNumber,
		ProfilePicture: user.ProfilePicture,
		Addresses:      user.Addresses,
		CreatedAt:      user.CreatedAt,
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
}

func (a *AuthController) Login(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	requestID := ctx.Value(constants.RequestIDKey).(string)

	var req model.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.Error(ctx, "controller.Login invalid body: %v", err)

		util.WriteErrorResponse(w, requestID, constants.InvalidRequestBody, constants.ErrInvalidRequest, http.StatusBadRequest)
		return
	}

	user, accessToken, refreshToken, err := a.service.Login(&req, ctx)
	if err != nil {
		util.Error(ctx, "controller.Login service error: %v", err)

		util.WriteErrorResponse(w, requestID, constants.InvalidCredentials, constants.ErrInvalidCredsCode, http.StatusUnauthorized)
		return
	}

	userResponse := &model.UserResponse{
		ID:             user.ID,
		Name:           user.Name,
		Email:          user.Email,
		Provider:       user.Provider,
		MobileNumber:   user.MobileNumber,
		ProfilePicture: user.ProfilePicture,
		Addresses:      user.Addresses,
		CreatedAt:      user.CreatedAt,
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
}
func (a *AuthController) GoogleLogin(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	requestID := ctx.Value(constants.RequestIDKey).(string)

	var req model.GoogleLoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.Error(ctx, "controller.GoogleLogin invalid body: %v", err)

		util.WriteErrorResponse(w, requestID, constants.InvalidRequestBody, constants.ErrInvalidRequest, http.StatusBadRequest)
		return
	}

	user, accessToken, refreshToken, err := a.service.GoogleLogin(ctx, req.IDToken)
	if err != nil {
		util.Error(ctx, "controller.GoogleLogin service error: %v", err)

		util.WriteErrorResponse(w, requestID, err.Error(), constants.ErrGoogleLoginFailed, http.StatusUnauthorized)
		return
	}

	userResponse := &model.UserResponse{
		ID:             user.ID,
		Name:           user.Name,
		Email:          user.Email,
		Provider:       user.Provider,
		MobileNumber:   user.MobileNumber,
		ProfilePicture: user.ProfilePicture,
		Addresses:      user.Addresses,
		CreatedAt:      user.CreatedAt,
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
}
func (a *AuthController) Refresh(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	requestID := ctx.Value(constants.RequestIDKey).(string)

	var req model.RefreshRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.WriteErrorResponse(w, requestID, constants.InvalidRequestBody, constants.ErrInvalidRequest, http.StatusBadRequest)
		return
	}

	token, err := jwt.Parse(req.RefreshToken, func(t *jwt.Token) (interface{}, error) {
		return util.Secret, nil
	})

	if err != nil || !token.Valid {
		util.WriteErrorResponse(w, requestID, "invalid refresh token", constants.ErrInvalidRefreshToken, http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		util.WriteErrorResponse(w, requestID, "invalid token claims", constants.ErrInvalidTokenClaims, http.StatusUnauthorized)
		return
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		util.WriteErrorResponse(w, requestID, "invalid user id", constants.ErrInvalidTokenClaims, http.StatusUnauthorized)
		return
	}

	newAccessToken, err := util.GenerateAccessToken(userID)
	if err != nil {
		util.WriteErrorResponse(w, requestID, "token generation failed", constants.ErrTokenGenerationFailed, http.StatusInternalServerError)
		return
	}

	resp := util.APIResponse{
		RequestID: requestID,
		Success:   true,
		Message:   "token refreshed successfully",
		Data: map[string]string{
			"access_token": newAccessToken,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}