package v1

import (
	"fmt"
	"net/http"

	"github.com/dana-team/platform-backend/internal/utils/pagination"

	"github.com/dana-team/platform-backend/internal/controllers"
	"github.com/dana-team/platform-backend/internal/customerrors"
	"github.com/dana-team/platform-backend/internal/middleware"
	"github.com/dana-team/platform-backend/internal/types"
	"github.com/gin-gonic/gin"
)

// serviceAccountHandler wraps a handler function with context setup for serviceAccountController.
func serviceAccountHandler(handler func(controller controllers.ServiceAccountController, c *gin.Context) (interface{}, error)) gin.HandlerFunc {
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
		serviceAccountController := controllers.NewServiceAccountController(kubeClient, context, logger)

		result, err := handler(serviceAccountController, c)
		if middleware.AddErrorToContext(c, err) {
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

// GetToken returns a Gin handler function for retrieving token of a specific service account.
func GetToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request types.ServiceAccountRequestUri
		if err := c.BindUri(&request); err != nil {
			middleware.AddErrorToContext(c, customerrors.NewValidationError(err.Error()))
			return
		}

		serviceAccountHandler(func(controller controllers.ServiceAccountController, c *gin.Context) (interface{}, error) {
			return controller.GetServiceAccountToken(request.ServiceAccountName, request.NamespaceName)
		})(c)
	}
}

// CreateServiceAccount returns a Gin handler function for creating a service account.
func CreateServiceAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request types.ServiceAccountRequestUri
		if err := c.BindUri(&request); err != nil {
			middleware.AddErrorToContext(c, customerrors.NewValidationError(err.Error()))
			return
		}

		serviceAccountHandler(func(controller controllers.ServiceAccountController, c *gin.Context) (interface{}, error) {
			return controller.CreateServiceAccount(request.ServiceAccountName, request.NamespaceName)
		})(c)
	}
}

// DeleteServiceAccount returns a Gin handler function for deleting a service account.
func DeleteServiceAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request types.ServiceAccountRequestUri
		if err := c.BindUri(&request); err != nil {
			middleware.AddErrorToContext(c, customerrors.NewValidationError(err.Error()))
			return
		}

		serviceAccountHandler(func(controller controllers.ServiceAccountController, c *gin.Context) (interface{}, error) {
			name := request.ServiceAccountName
			message := fmt.Sprintf("Deleted serviceAccount successfully %q", name)
			return types.MessageResponse{Message: message}, controller.DeleteServiceAccount(request.ServiceAccountName, request.NamespaceName)
		})(c)
	}
}

// GetServiceAccounts returns a Gin handler function for retrieving service accounts in a namespace.
func GetServiceAccounts() gin.HandlerFunc {
	return func(c *gin.Context) {
		var namespace types.NamespaceUri
		if err := c.BindUri(&namespace); err != nil {
			middleware.AddErrorToContext(c, customerrors.NewValidationError(err.Error()))
			return
		}

		limit, page, err := pagination.ExtractPaginationParamsFromCtx(c)
		if err != nil {
			middleware.AddErrorToContext(c, customerrors.NewValidationError(err.Error()))
			return
		}

		serviceAccountHandler(func(controller controllers.ServiceAccountController, c *gin.Context) (interface{}, error) {
			return controller.GetServiceAccounts(namespace.NamespaceName, limit, page)
		})(c)
	}
}

// GetServiceAccount returns a Gin handler function for retrieving a specific service account from a namespace.
func GetServiceAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request types.ServiceAccountRequestUri
		if err := c.BindUri(&request); err != nil {
			middleware.AddErrorToContext(c, customerrors.NewValidationError(err.Error()))
			return
		}

		serviceAccountHandler(func(controller controllers.ServiceAccountController, c *gin.Context) (interface{}, error) {
			return controller.GetServiceAccount(request.ServiceAccountName, request.NamespaceName)
		})(c)
	}
}
