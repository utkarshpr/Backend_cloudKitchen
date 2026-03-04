package constants

// Success Messages
const (
	SignupSuccess = "signup successful"
	LoginSuccess  = "login successful"
)

// Error Messages
const (
	InvalidRequestBody = "invalid request body"
	InvalidCredentials = "invalid credentials"
	UserAlreadyExists  = "user already exists"
)

// Error Codes
const (
	ErrInvalidRequest   = "INVALID_REQUEST"
	ErrSignupFailed     = "SIGNUP_FAILED"
	ErrInvalidCredsCode = "INVALID_CREDENTIALS"
)

type ContextKey string

const RequestIDKey ContextKey = "request_id"

const (
	ErrGoogleLoginFailed = "GOOGLE_LOGIN_FAILED"
)
