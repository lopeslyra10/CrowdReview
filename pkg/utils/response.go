package utils

import "github.com/gin-gonic/gin"

// JSONError sends a standardized error payload.
func JSONError(c *gin.Context, code int, msg string) {
	c.JSON(code, gin.H{"error": msg})
}

// JSONSuccess sends a standardized success payload.
func JSONSuccess(c *gin.Context, code int, data interface{}) {
	c.JSON(code, gin.H{"data": data})
}
