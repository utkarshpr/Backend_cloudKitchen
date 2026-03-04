package util

type APIResponse struct {
	RequestID string      `json:"request_id"`
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
	ErrorCode string      `json:"error_code"`
}
