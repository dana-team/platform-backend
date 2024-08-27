package v1

import (
	"github.com/dana-team/platform-backend/src/customerrors"
	"github.com/dana-team/platform-backend/src/routes"
	"net/http"

	"github.com/dana-team/platform-backend/src/controllers"
	"github.com/dana-team/platform-backend/src/types"
	"github.com/gin-gonic/gin"
)

// configMapHandler handles the request of the client to the Kubernetes cluster.
func configMapHandler(handler func(controller controllers.ConfigMapController, c *gin.Context) (interface{}, error)) gin.HandlerFunc {
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
		configMapController := controllers.NewConfigMapController(kubeClient, context, logger)

		result, err := handler(configMapController, c)
		if routes.AddErrorToContext(c, err) {
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

// GetConfigMap gets a specific config map from the specified namespace.
func GetConfigMap() gin.HandlerFunc {
	return func(c *gin.Context) {
		var uriRequest types.ConfigMapUri
		if err := c.BindUri(&uriRequest); err != nil {
			routes.AddErrorToContext(c, customerrors.NewValidationError(err.Error()))
			return
		}

		configMapHandler(func(controller controllers.ConfigMapController, c *gin.Context) (interface{}, error) {
			return controller.GetConfigMap(uriRequest.NamespaceName, uriRequest.ConfigMapName)
		})(c)
	}
}
