package v1

import (
	"fmt"
	"net/http"

	"github.com/dana-team/platform-backend/internal/customerrors"
	"github.com/dana-team/platform-backend/internal/types"

	"github.com/dana-team/platform-backend/internal/controllers"
	"github.com/dana-team/platform-backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

// tokenHandler wraps a handler function with context setup for tokenController.
func tokenHandler(handler func(controller controllers.TokenController, c *gin.Context) (interface{}, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		kubeClient, err := middleware.GetKubeClient(c)
		if middleware.AddErrorToContext(c, err) {
			return
		}

		logger, err := middleware.GetLogger(c)
		if middleware.AddErrorToContext(c, err) {
			return
		}

		context := c.Request.Context()
		tokenController := controllers.NewTokenController(kubeClient, context, logger)

		result, err := handler(tokenController, c)
		if middleware.AddErrorToContext(c, err) {
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

// CreateToken returns a Gin handler function for creating a service account token.
func CreateToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request types.ServiceAccountRequestUri
		if err := c.BindUri(&request); err != nil {
			middleware.AddErrorToContext(c, customerrors.NewValidationError(err.Error()))
			return
		}

		var query types.CreateTokenQuery
		if err := c.BindQuery(&query); err != nil {
			middleware.AddErrorToContext(c, customerrors.NewValidationError(err.Error()))
			return
		}

		tokenHandler(func(controller controllers.TokenController, c *gin.Context) (interface{}, error) {
			return controller.CreateToken(request.ServiceAccountName, request.NamespaceName, query.ExpirationSeconds)
		})(c)
	}
}

// RevokeToken returns a Gin handler function for revoking the tokens for a service account.
func RevokeToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request types.ServiceAccountRequestUri
		if err := c.BindUri(&request); err != nil {
			middleware.AddErrorToContext(c, customerrors.NewValidationError(err.Error()))
			return
		}

		tokenHandler(func(controller controllers.TokenController, c *gin.Context) (interface{}, error) {
			name := request.ServiceAccountName
			message := fmt.Sprintf("Revoked tokens for ServiceAccount %q", name)
			return types.MessageResponse{Message: message}, controller.RevokeToken(request.ServiceAccountName, request.NamespaceName)
		})(c)
	}
}
