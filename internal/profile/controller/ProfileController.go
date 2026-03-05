package controller

import (
	"cloud-kitchen/internal/profile/service"
	constant "cloud-kitchen/pkg/constants"

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
	
	util.WriteSuccessResponse(w,requestId, "Profile retrieved successfully", "", http.StatusOK, profile)
}

func (c *ProfileController) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	// Implement logic to update user profile
}

func (c *ProfileController) DeleteProfile(w http.ResponseWriter, r *http.Request) {
	// Implement logic to delete user profile
}
