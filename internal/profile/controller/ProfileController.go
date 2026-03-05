package controller

import (
	profilemodel "cloud-kitchen/internal/profile/model"
	"cloud-kitchen/internal/profile/service"
	constant "cloud-kitchen/pkg/constants"
	"encoding/json"

	"cloud-kitchen/pkg/util"
	"net/http"
)

type ProfileController struct {
	service *service.ProfileService
}

func NewProfileController(service *service.ProfileService) *ProfileController {
	return &ProfileController{service: service}
}

func (c *ProfileController) GetProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := ctx.Value(constant.RequestIDKey).(string)
	email := r.URL.Query().Get("email")
	profile, err := c.service.GetProfile(ctx, email)
	if err != nil {
		http.Error(w, "Profile not found", http.StatusNotFound)
		return
	}
	util.Info(ctx, "New Profile Controller [Started]")

	util.WriteSuccessResponse(w, requestId, "Profile retrieved successfully", "", http.StatusOK, profile)
}

func (c *ProfileController) UpdateProfile(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	requestID := ctx.Value(constant.RequestIDKey).(string)

	var req profilemodel.UpdateProfileRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		util.WriteErrorResponse(w, requestID, "Invalid request", "invalid body", http.StatusBadRequest)
		return
	}

	userID := ctx.Value("user_id").(string)

	err = c.service.UpdateProfile(ctx, userID, &req)
	if err != nil {
		util.WriteErrorResponse(w, requestID, "Update failed", err.Error(), http.StatusInternalServerError)
		return
	}

	util.WriteSuccessResponse(w, requestID, "Profile updated successfully", "", http.StatusOK, nil)
}

func (c *ProfileController) DeleteProfile(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	requestID := ctx.Value(constant.RequestIDKey).(string)

	userIDVal := ctx.Value("user_id")
	userID, ok := userIDVal.(string)
	if !ok {
		util.WriteErrorResponse(
			w,
			requestID,
			"Unauthorized",
			"user id missing",
			http.StatusUnauthorized,
		)
		return
	}

	err := c.service.DeleteProfile(ctx, userID)
	if err != nil {
		util.WriteErrorResponse(
			w,
			requestID,
			"Delete failed",
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	util.WriteSuccessResponse(
		w,
		requestID,
		"Profile deleted successfully",
		"",
		http.StatusOK,
		nil,
	)
}
