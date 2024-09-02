package v1

import (
	"github.com/dana-team/platform-backend/src/customerrors"
	"github.com/dana-team/platform-backend/src/middleware"
	"github.com/dana-team/platform-backend/src/routes"

	"github.com/dana-team/platform-backend/src/controllers"
	"github.com/dana-team/platform-backend/src/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

// containerHandler wraps a handler function with context setup for ContainerController.
func containerHandler(handler func(controller controllers.ContainerController, c *gin.Context) (interface{}, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		kubeClient, err := middleware.GetKubeClient(c)
		if middleware.AddErrorToContext(c, err) {
			return
		}

		logger, err := middleware.GetLogger(c)
		if middleware.AddErrorToContext(c, err) {
			return
		}

		context := routes.GetContext(c)
		containerController := controllers.NewContainerController(kubeClient, context, logger)

		result, err := handler(containerController, c)
		if middleware.AddErrorToContext(c, err) {
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

// GetPodsContainers returns a Gin handler function for retrieving containers information.
func GetPodsContainers() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request types.ContainerRequestUri
		if err := c.BindUri(&request); err != nil {
			middleware.AddErrorToContext(c, customerrors.NewValidationError(err.Error()))
			return
		}

		containerHandler(func(controller controllers.ContainerController, c *gin.Context) (interface{}, error) {
			return controller.GetContainers(request.NamespaceName, request.PodName)
		})(c)
	}
}
