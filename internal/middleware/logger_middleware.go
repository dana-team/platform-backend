package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LoggerMiddleware is a middleware that initializes a logger and sets it in the request context.
// The logger is initialized with a default "unknown" user. This middleware should be applied
// globally to ensure the logger is available in the context for all requests.
func LoggerMiddleware(baseLogger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("logger", baseLogger)
		c.Next()
	}
}
