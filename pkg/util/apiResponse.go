package util

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	RequestID string      `json:"request_id"`
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
	ErrorCode string      `json:"error_code"`
}

// WriteErrorResponse writes a standardized error response to the http.ResponseWriter
func WriteErrorResponse(w http.ResponseWriter, requestID, message, code string, statusCode int) {
	resp := APIResponse{
		RequestID: requestID,
		Success:   false,
		Message:   message,
		ErrorCode: code,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(resp)
}
