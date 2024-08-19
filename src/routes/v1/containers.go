package v1

import (
	"github.com/dana-team/platform-backend/src/customerrors"

	"github.com/dana-team/platform-backend/src/controllers"
	"github.com/dana-team/platform-backend/src/routes"
	"github.com/dana-team/platform-backend/src/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

// containerHandler wraps a handler function with context setup for ContainerController.
func containerHandler(handler func(controller controllers.ContainerController, c *gin.Context) (interface{}, error)) gin.HandlerFunc {
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
		containerController := controllers.NewContainerController(kubeClient, context, logger)

		result, err := handler(containerController, c)
		if routes.AddErrorToContext(c, err) {
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

// GetContainers returns a Gin handler function for retrieving containers information.
func GetContainers() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request types.ContainerRequestUri
		if err := c.BindUri(&request); err != nil {
			routes.AddErrorToContext(c, customerrors.NewValidationError(err.Error()))
			return
		}

		containerHandler(func(controller controllers.ContainerController, c *gin.Context) (interface{}, error) {
			return controller.GetContainers(request.NamespaceName, request.PodName)
		})(c)
	}
}
