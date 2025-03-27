package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ErrorHandler is a middleware that logs errors and returns a JSON error response.
func ErrorHandler(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		// If there are errors attached to the context, log them and send an error response.
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				logger.Error("HTTP error", zap.Error(err.Err))
			}
			// If no response has been written, send a JSON error.
			if !c.Writer.Written() {
				c.JSON(http.StatusInternalServerError, gin.H{"error": c.Errors[0].Error()})
			}
		}
	}
}
