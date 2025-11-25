package middleware

import (
	"net"
	"net/http"
	"time"

	"crowdreview/config"
	"crowdreview/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// RateLimitMiddleware applies a basic fixed-window limit per IP using Redis.
func RateLimitMiddleware(rdb *redis.Client, cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := clientIP(c)
		key := "rate:" + ip
		ctx := c.Request.Context()

		pipe := rdb.TxPipeline()
		incr := pipe.Incr(ctx, key)
		pipe.Expire(ctx, key, cfg.RateLimitWindow)
		if _, err := pipe.Exec(ctx); err != nil {
			utils.JSONError(c, http.StatusInternalServerError, "rate limiter unavailable")
			c.Abort()
			return
		}
		if incr.Val() > int64(cfg.RateLimitRequests) {
			utils.JSONError(c, http.StatusTooManyRequests, "too many requests")
			c.Abort()
			return
		}
		c.Next()
	}
}

func clientIP(c *gin.Context) string {
	if ip := c.ClientIP(); ip != "" {
		return ip
	}
	host, _, err := net.SplitHostPort(c.Request.RemoteAddr)
	if err != nil {
		return c.Request.RemoteAddr
	}
	return host
}
