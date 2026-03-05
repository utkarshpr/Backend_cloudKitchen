package middleware

import (
	"context"
	"strings"

	"cloud-kitchen/pkg/constants"
	"cloud-kitchen/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

func AuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		auth := c.GetHeader("Authorization")

		if auth == "" {
			c.JSON(401, gin.H{"error": "no token"})
			c.Abort()
			return
		}

		tokenString := strings.Split(auth, " ")[1]

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return util.Secret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(401, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		c.Set("user_id", claims["user_id"])

		c.Next()
	}
}
