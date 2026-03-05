package controller

import (
	profilemodel "cloud-kitchen/internal/profile/model"
	"cloud-kitchen/internal/profile/service"
	constant "cloud-kitchen/pkg/constants"
	"cloud-kitchen/pkg/util"
	"encoding/json"
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
	requestID := ctx.Value(constant.RequestIDKey).(string)

	util.Info(ctx, "ProfileController.GetProfile started request_id=%s", requestID)

	userIDVal := ctx.Value("user_id")
	userID, ok := userIDVal.(string)
	if !ok {
		util.Error(ctx, "ProfileController.GetProfile user_id missing request_id=%s", requestID)

		util.WriteErrorResponse(
			w,
			requestID,
			"Unauthorized",
			"user id missing",
			http.StatusUnauthorized,
		)
		return
	}

	util.Info(ctx, "ProfileController.GetProfile fetching profile user_id=%s", userID)

	profile, err := c.service.GetProfile(ctx, userID)
	if err != nil {

		util.Error(ctx, "ProfileController.GetProfile failed user_id=%s err=%v", userID, err)

		util.WriteErrorResponse(
			w,
			requestID,
			"Profile not found",
			"PROFILE_NOT_FOUND",
			http.StatusNotFound,
		)
		return
	}

	util.Info(ctx, "ProfileController.GetProfile success user_id=%s", userID)

	util.WriteSuccessResponse(
		w,
		requestID,
		"Profile retrieved successfully",
		"",
		http.StatusOK,
		profile,
	)
}

func (c *ProfileController) UpdateProfile(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	requestID := ctx.Value(constant.RequestIDKey).(string)

	util.Info(ctx, "ProfileController.UpdateProfile started request_id=%s", requestID)

	defer r.Body.Close()

	var req profilemodel.UpdateProfileRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {

		util.Error(ctx, "ProfileController.UpdateProfile invalid body request_id=%s err=%v", requestID, err)

		util.WriteErrorResponse(
			w,
			requestID,
			"Invalid request",
			"INVALID_BODY",
			http.StatusBadRequest,
		)
		return
	}

	userIDVal := ctx.Value("user_id")
	userID, ok := userIDVal.(string)
	if !ok {

		util.Error(ctx, "ProfileController.UpdateProfile user_id missing request_id=%s", requestID)

		util.WriteErrorResponse(
			w,
			requestID,
			"Unauthorized",
			"user id missing",
			http.StatusUnauthorized,
		)
		return
	}

	util.Info(ctx, "ProfileController.UpdateProfile updating profile user_id=%s", userID)

	err = c.service.UpdateProfile(ctx, userID, &req)
	if err != nil {

		util.Error(ctx, "ProfileController.UpdateProfile failed user_id=%s err=%v", userID, err)

		util.WriteErrorResponse(
			w,
			requestID,
			"Update failed",
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	util.Info(ctx, "ProfileController.UpdateProfile success user_id=%s", userID)

	util.WriteSuccessResponse(
		w,
		requestID,
		"Profile updated successfully",
		"",
		http.StatusOK,
		nil,
	)
}

func (c *ProfileController) DeleteProfile(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	requestID := ctx.Value(constant.RequestIDKey).(string)

	util.Info(ctx, "ProfileController.DeleteProfile started request_id=%s", requestID)

	userIDVal := ctx.Value("user_id")
	userID, ok := userIDVal.(string)
	if !ok {

		util.Error(ctx, "ProfileController.DeleteProfile user_id missing request_id=%s", requestID)

		util.WriteErrorResponse(
			w,
			requestID,
			"Unauthorized",
			"user id missing",
			http.StatusUnauthorized,
		)
		return
	}

	util.Info(ctx, "ProfileController.DeleteProfile deleting user_id=%s", userID)

	err := c.service.DeleteProfile(ctx, userID)
	if err != nil {

		util.Error(ctx, "ProfileController.DeleteProfile failed user_id=%s err=%v", userID, err)

		util.WriteErrorResponse(
			w,
			requestID,
			"Delete failed",
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	util.Info(ctx, "ProfileController.DeleteProfile success user_id=%s", userID)

	util.WriteSuccessResponse(
		w,
		requestID,
		"Profile deleted successfully",
		"",
		http.StatusOK,
		nil,
	)
}
