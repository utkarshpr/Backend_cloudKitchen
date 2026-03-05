package middleware

import (
	"context"
	"net/http"
	"strings"

	constant "cloud-kitchen/pkg/constants"
	"cloud-kitchen/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func RequestIDMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		requestID := uuid.New().String()

		ctx := context.WithValue(c.Request.Context(), constant.RequestIDKey, requestID)

		c.Request = c.Request.WithContext(ctx)

		c.Writer.Header().Set("X-Request-ID", requestID)

		// log request start
		util.Info(ctx, "incoming %s %s", c.Request.Method, c.Request.URL.Path)

		c.Next()
	}
}

func AuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		auth := c.GetHeader("Authorization")

		if auth == "" {
			util.WriteErrorResponse(
				c.Writer,
				c.Request.Context().Value(constant.RequestIDKey).(string),
				"Invalid token",
				"Authorization header missing",
				http.StatusUnauthorized,
			)
			c.Abort()
			return
		}

		tokenString := strings.Split(auth, " ")[1]

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return util.Secret, nil
		})

		if err != nil || !token.Valid {
			util.WriteErrorResponse(c.Writer, c.Request.Context().Value(constant.RequestIDKey).(string), "Invalid token", "invalid token", 401)
			c.Abort()
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		c.Set("user_id", claims["user_id"])

		c.Next()
	}
}
