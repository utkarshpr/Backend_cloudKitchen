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
	Name           string               `json:"name"`
	MobileNumber   string               `json:"mobile_number"`
	ProfilePicture string               `json:"profile_picture"`
	Addresses      []UpdateAddressInput `json:"addresses"`
}

type UpdateAddressInput struct {
	ID        *string `json:"id,omitempty"`
	Label     string  `json:"label"`
	Street    string  `json:"street"`
	City      string  `json:"city"`
	State     string  `json:"state"`
	ZipCode   string  `json:"zip_code"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	IsDefault bool    `json:"is_default"`
}