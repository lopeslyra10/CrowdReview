package middleware

import (
	"net/http"

	"crowdreview/pkg/utils"

	"github.com/gin-gonic/gin"
)

// AdminRequired ensures user role is admin.
func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "admin" {
			utils.JSONError(c, http.StatusForbidden, "admin only")
			c.Abort()
			return
		}
		c.Next()
	}
}
