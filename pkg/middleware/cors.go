package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware enables browser calls from the Next.js frontend.
// By default it allows http://localhost:3000 but can be overridden via CORS_ORIGIN env var.
func CORSMiddleware() gin.HandlerFunc {
	allowedOrigin := os.Getenv("CORS_ORIGIN")
	if allowedOrigin == "" {
		allowedOrigin = "http://localhost:3000"
	}

	return func(c *gin.Context) {
		h := c.Writer.Header()
		h.Set("Access-Control-Allow-Origin", allowedOrigin)
		h.Set("Access-Control-Allow-Credentials", "true")
		h.Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept, X-Requested-With")
		h.Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
