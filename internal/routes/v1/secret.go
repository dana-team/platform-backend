package v1

import (
	"github.com/dana-team/platform-backend/internal/customerrors"
	"github.com/dana-team/platform-backend/internal/middleware"
	"github.com/dana-team/platform-backend/internal/utils/pagination"
	"net/http"

	"github.com/dana-team/platform-backend/internal/controllers"
	"github.com/dana-team/platform-backend/internal/types"
	"github.com/gin-gonic/gin"
)

// secretHandler handles the request of the client to the Kubernetes cluster.
func secretHandler(handler func(controller controllers.SecretController, c *gin.Context) (interface{}, error)) gin.HandlerFunc {
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
		secretController := controllers.NewSecretController(kubeClient, context, logger)

		result, err := handler(secretController, c)
		if middleware.AddErrorToContext(c, err) {
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

// CreateSecret creates a new secret in a specific namespace.
func CreateSecret() gin.HandlerFunc {
	return func(c *gin.Context) {
		var uriRequest types.SecretNamespaceUriRequest
		if err := c.BindUri(&uriRequest); err != nil {
			middleware.AddErrorToContext(c, customerrors.NewValidationError(err.Error()))
			return
		}
		var request types.CreateSecretRequest
		if err := c.BindJSON(&request); err != nil {
			middleware.AddErrorToContext(c, customerrors.NewValidationError(err.Error()))
			return
		}

		secretHandler(func(controller controllers.SecretController, c *gin.Context) (interface{}, error) {
			return controller.CreateSecret(uriRequest.NamespaceName, request)
		})(c)
	}
}

// GetSecrets gets all secrets in a specific namespace.
func GetSecrets() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request types.SecretNamespaceUriRequest
		if err := c.BindUri(&request); err != nil {
			middleware.AddErrorToContext(c, customerrors.NewValidationError(err.Error()))
			return
		}

		limit, page, err := pagination.ExtractPaginationParamsFromCtx(c)
		if err != nil {
			middleware.AddErrorToContext(c, customerrors.NewValidationError(err.Error()))
			return
		}

		secretHandler(func(controller controllers.SecretController, c *gin.Context) (interface{}, error) {
			return controller.GetSecrets(request.NamespaceName, limit, page)
		})(c)
	}
}

// GetSecret gets a specific secret from a specific namespace.
func GetSecret() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request types.SecretUriRequest
		if err := c.BindUri(&request); err != nil {
			middleware.AddErrorToContext(c, customerrors.NewValidationError(err.Error()))
			return
		}

		secretHandler(func(controller controllers.SecretController, c *gin.Context) (interface{}, error) {
			return controller.GetSecret(request.NamespaceName, request.SecretName)
		})(c)
	}
}

// UpdateSecret updates a specific secret in a specific namespace.
func UpdateSecret() gin.HandlerFunc {
	return func(c *gin.Context) {
		var uriRequest types.SecretUriRequest
		if err := c.BindUri(&uriRequest); err != nil {
			middleware.AddErrorToContext(c, customerrors.NewValidationError(err.Error()))
			return
		}
		var request types.UpdateSecretRequest
		if err := c.BindJSON(&request); err != nil {
			middleware.AddErrorToContext(c, customerrors.NewValidationError(err.Error()))
			return
		}

		secretHandler(func(controller controllers.SecretController, c *gin.Context) (interface{}, error) {
			return controller.UpdateSecret(uriRequest.NamespaceName, uriRequest.SecretName, request)
		})(c)
	}
}

// DeleteSecret deletes a specific secret in a specific namespace.
func DeleteSecret() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request types.SecretUriRequest
		if err := c.BindUri(&request); err != nil {
			middleware.AddErrorToContext(c, customerrors.NewValidationError(err.Error()))
			return
		}

		secretHandler(func(controller controllers.SecretController, c *gin.Context) (interface{}, error) {
			return controller.DeleteSecret(request.NamespaceName, request.SecretName)
		})(c)
	}
}
