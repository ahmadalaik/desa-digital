package middlewares

import (
	"net/http"
	"strings"

	"github.com/ahmadalaik/desa-digital/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte(config.GetEnv("JWT_SECRET", "secret_key"))

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")

		if tokenStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token is required",
			})
			c.Abort()
			return
		}

		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

		claims := &jwt.RegisteredClaims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthenticated",
			})
			c.Abort()
			return
		}

		c.Set("username", claims.Subject)
		c.Next()
	}
}
