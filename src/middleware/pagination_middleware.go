package middleware

import (
	"fmt"
	"github.com/dana-team/platform-backend/src/customerrors"
	"github.com/dana-team/platform-backend/src/types"
	"github.com/gin-gonic/gin"
)

const (
	PageCtxKey  = "page"
	LimitCtxKey = "limit"
)

// PaginationMiddleware is a middleware for extracting and setting pagination parameters from the request context
func PaginationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger, err := GetLogger(c)
		if AddErrorToContext(c, err) {
			return
		}

		logger.Info("setting up pagination middleware")
		var paginationParams types.PaginationParams
		if err := c.BindQuery(&paginationParams); err != nil {
			AddErrorToContext(c, customerrors.NewValidationError(fmt.Sprintf("invalid request, %s", err)))
			c.Abort()
			return
		}

		c.Set(LimitCtxKey, paginationParams.Limit)
		c.Set(PageCtxKey, paginationParams.Page)
		c.Next()
	}
}
