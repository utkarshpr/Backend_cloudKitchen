package profilemodel

import (
	"cloud-kitchen/internal/auth/model"
)

type Profile struct {
	ID             string               `json:"id"`
	Name           string               `json:"name"`
	Email          string               `json:"email"`
	MobileNumber   string               `json:"mobile_number"`
	ProfilePicture string               `json:"profile_picture"`
	Addresses      []model.AddressModel `json:"addresses"`
	Provider       string               `json:"provider"`
}

type UpdateProfileRequest struct {
	Name           string `json:"name"`
	MobileNumber   string `json:"mobile_number"`
	ProfilePicture string `json:"profile_picture"`
}
