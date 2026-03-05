package model

import "time"

type User struct {
	ID             string
	Name           string
	Email          string
	Password       string
	Provider       string
	CreatedAt      time.Time
	MobileNumber   string
	ProfilePicture string
	Addresses      []Address
}

// ...existing code...

type SignupRequest struct {
	Name           string    `json:"name" binding:"required"`
	Email          string    `json:"email" binding:"required,email"`
	Password       string    `json:"password" binding:"required,min=6"`
	MobileNumber   string    `json:"mobile_number" binding:"required"`
	ProfilePicture string    `json:"profile_picture"`
	Addresses      []Address `json:"addresses" binding:"required"`
}

type Address struct {
	Label     string  `json:"label" binding:"required"` // e.g., "Home", "Office", "Other"
	Street    string  `json:"street" binding:"required"`
	City      string  `json:"city" binding:"required"`
	State     string  `json:"state" binding:"required"`
	ZipCode   string  `json:"zip_code" binding:"required"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	IsDefault bool    `json:"is_default"`
}

// ...existing code...

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// ...existing code...

type UserResponse struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	Provider       string    `json:"provider"`
	MobileNumber   string    `json:"mobile_number"`
	ProfilePicture string    `json:"profile_picture"`
	Addresses      []Address `json:"addresses"`
	CreatedAt      time.Time `json:"created_at"`
}

// ...existing code...

type SignupResponse struct {
	User         *UserResponse `json:"user"`
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
}

type LoginResponse struct {
	User         *UserResponse `json:"user"`
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
}

type GoogleLoginRequest struct {
	IDToken string `json:"id_token"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type AddressModel struct {
	ID        string
	UserID    string
	Label     string
	Street    string
	City      string
	State     string
	ZipCode   string
	Latitude  float64
	Longitude float64
	IsDefault bool
}

type UserRepo struct {
	ID             string
	Name           string
	Email          string
	Password       string
	Provider       string
	MobileNumber   string
	ProfilePicture string
	CreatedAt      time.Time
}