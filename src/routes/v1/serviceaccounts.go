package v1

import (
	"github.com/dana-team/platform-backend/src/controllers"
	"github.com/dana-team/platform-backend/src/customerrors"
	"github.com/dana-team/platform-backend/src/routes"
	"github.com/dana-team/platform-backend/src/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

// serviceAccountHandler wraps a handler function with context setup for serviceAccountController.
func serviceAccountHandler(handler func(controller controllers.ServiceAccountController, c *gin.Context) (interface{}, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		kubeClient, err := routes.GetKubeClient(c)
		if routes.AddErrorToContext(c, err) {
			return
		}

		logger, err := routes.GetLogger(c)
		if routes.AddErrorToContext(c, err) {
			return
		}

		context := c.Request.Context()
		serviceAccountController := controllers.NewServiceAccountController(kubeClient, context, logger)

		result, err := handler(serviceAccountController, c)
		if routes.AddErrorToContext(c, err) {
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
			routes.AddErrorToContext(c, customerrors.NewValidationError(err.Error()))
			return
		}

		serviceAccountHandler(func(controller controllers.ServiceAccountController, c *gin.Context) (interface{}, error) {
			return controller.GetServiceAccountToken(request.ServiceAccountName, request.NamespaceName)
		})(c)
	}
}
