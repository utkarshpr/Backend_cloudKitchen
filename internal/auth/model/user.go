package model

import "time"

type User struct {
	ID        string
	Name      string
	Email     string
	Password  string
	Provider  string
	CreatedAt time.Time
}

type SignupRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Provider  string    `json:"provider"`
	CreatedAt time.Time `json:"createdAt"`
}

type SignupResponse struct {
	User  *UserResponse `json:"user"`
	Token string        `json:"token"`
}

type LoginResponse struct {
	User  *UserResponse `json:"user"`
	Token string        `json:"token"`
}

type GoogleLoginRequest struct {
	IDToken string `json:"id_token"`
}
