package middleware

import (
	"github.com/dana-team/platform-backend/src/types"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

const (
	PageCtxKey  = "page"
	LimitCtxKey = "limit"
)

// PaginationMiddleware is a middleware for extracting and setting pagination parameters from the request context
func PaginationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctxLogger, exists := c.Get(LoggerCtxKey)
		if !exists {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "logger not set in context"})
			return
		}
		logger := ctxLogger.(*zap.Logger)

		logger.Info("setting up pagination middleware")

		var paginationParams types.PaginationParams
		if err := c.BindQuery(&paginationParams); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
			return
		}

		c.Set(LimitCtxKey, paginationParams.Limit)
		c.Set(PageCtxKey, paginationParams.Page)
		c.Next()
	}
}
