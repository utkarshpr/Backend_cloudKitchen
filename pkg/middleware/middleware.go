package middleware

import (
	"context"

	"cloud-kitchen/pkg/constants"
	"cloud-kitchen/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestIDMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		requestID := uuid.New().String()

		ctx := context.WithValue(c.Request.Context(), constants.RequestIDKey, requestID)

		c.Request = c.Request.WithContext(ctx)

		c.Writer.Header().Set("X-Request-ID", requestID)

		// log request start
		util.Info(ctx, "incoming %s %s", c.Request.Method, c.Request.URL.Path)

		c.Next()
	}
}
