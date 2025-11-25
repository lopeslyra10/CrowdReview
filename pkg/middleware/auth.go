package middleware

import (
	"net/http"
	"strings"

	"crowdreview/config"
	"crowdreview/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuthRequired ensures a valid access token is present.
func AuthRequired(cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			utils.JSONError(c, http.StatusUnauthorized, "missing token")
			c.Abort()
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ParseToken(token, cfg.JWTSecret)
		if err != nil {
			utils.JSONError(c, http.StatusUnauthorized, "invalid token")
			c.Abort()
			return
		}
		userID, err := uuid.Parse(claims.Subject)
		if err != nil {
			utils.JSONError(c, http.StatusUnauthorized, "invalid token subject")
			c.Abort()
			return
		}
		role := claims.Role
		if role == "" {
			role = "user"
		}
		c.Set("userID", userID)
		c.Set("role", role)
		c.Next()
	}
}
