package util

import (
	"context"
	"log"

	"cloud-kitchen/pkg/constants"
)

// LogWithContext logs a formatted message prefixing it with the request id stored
// in the context (if available) so that every entry can be correlated.
func LogWithContext(ctx context.Context, format string, args ...interface{}) {
	requestID, _ := ctx.Value(constants.RequestIDKey).(string)
	if requestID != "" {
		log.Printf("[request_id=%s] "+format, append([]interface{}{requestID}, args...)...)
	} else {
		log.Printf(format, args...)
	}
}

// Convenience helpers follow if you need specific levels in future.
// For now everything goes through a single logger but the wrappers make
// it easier to change implementation later.
func Info(ctx context.Context, format string, args ...interface{}) {
	LogWithContext(ctx, "INFO: "+format, args...)
}

func Error(ctx context.Context, format string, args ...interface{}) {
	LogWithContext(ctx, "ERROR: "+format, args...)
}
