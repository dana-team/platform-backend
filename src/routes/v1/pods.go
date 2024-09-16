package v1

import (
	"github.com/dana-team/platform-backend/src/customerrors"
	"github.com/dana-team/platform-backend/src/middleware"
	"github.com/dana-team/platform-backend/src/routes"

	"github.com/dana-team/platform-backend/src/controllers"
	"github.com/dana-team/platform-backend/src/types"
	"github.com/dana-team/platform-backend/src/utils/pagination"
	"github.com/gin-gonic/gin"
	"net/http"
)

// podHandler wraps a handler function with context setup for PodController.
func podHandler(handler func(controller controllers.PodController, c *gin.Context) (interface{}, error)) gin.HandlerFunc {
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
		podController := controllers.NewPodController(kubeClient, context, logger)

		result, err := handler(podController, c)
		if middleware.AddErrorToContext(c, err) {
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

// GetPods returns a Gin handler function for retrieving pods of a specific capp.
func GetPods() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request types.PodRequestUri
		if err := c.BindUri(&request); err != nil {
			middleware.AddErrorToContext(c, customerrors.NewValidationError(err.Error()))
			return
		}

		limit, page, err := pagination.ExtractPaginationParamsFromCtx(c)
		if err != nil {
			middleware.AddErrorToContext(c, customerrors.NewValidationError(err.Error()))
			return
		}

		podHandler(func(controller controllers.PodController, c *gin.Context) (interface{}, error) {
			return controller.GetPods(request.NamespaceName, request.CappName, limit, page)
		})(c)
	}
}
