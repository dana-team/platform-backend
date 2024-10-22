package middleware

import (
	"errors"
	"github.com/dana-team/platform-backend/internal/customerrors"
	"github.com/dana-team/platform-backend/internal/types"
	"github.com/gin-gonic/gin"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"net/http"
)

const (
	errorKey      = "error"
	reason        = "reason"
	reasonUnknown = "Unknown"
)

// ErrorHandlingMiddleware handles errors and sends appropriate responses.
func ErrorHandlingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if !hasErrors(c) {
			return
		}

		statusCode := http.StatusInternalServerError
		lastError := c.Errors.Last()
		errorMessage := lastError.Err.Error()

		var customError customerrors.ErrorWithStatusCode
		var k8sErr *k8serrors.StatusError

		errorReason := reasonUnknown

		if errors.As(lastError.Err, &customError) {
			statusCode = customError.StatusCode()
			errorMessage = customError.Error()
			errorReason = string(customError.StatusReason())
		} else if errors.As(lastError.Err, &k8sErr) {
			statusCode = int(k8sErr.ErrStatus.Code)
			errorMessage = k8sErr.ErrStatus.Message
			errorReason = string(k8sErr.ErrStatus.Reason)
		}

		errorResponse := types.ErrorResponse{
			Error:  errorMessage,
			Reason: errorReason,
		}

		c.JSON(statusCode, errorResponse)
	}
}

// hasErrors returns true if there are any errors in the Gin context.
func hasErrors(c *gin.Context) bool {
	return len(c.Errors) > 0 && c.Errors.Last() != nil
}
