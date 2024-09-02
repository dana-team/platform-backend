package v1

import (
	"fmt"
	"github.com/dana-team/platform-backend/src/customerrors"
	"github.com/dana-team/platform-backend/src/middleware"
	"github.com/dana-team/platform-backend/src/utils/pagination"
	"net/http"

	"github.com/dana-team/platform-backend/src/controllers"
	"github.com/dana-team/platform-backend/src/types"
	"github.com/gin-gonic/gin"
)

func namespaceHandler(handler func(controller controllers.NamespaceController, c *gin.Context) (interface{}, error)) gin.HandlerFunc {
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
		namespaceController := controllers.NewNamespaceController(kubeClient, context, logger)

		result, err := handler(namespaceController, c)
		if middleware.AddErrorToContext(c, err) {
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func GetNamespaces() gin.HandlerFunc {
	return func(c *gin.Context) {
		limit, page, err := pagination.ExtractPaginationParamsFromCtx(c)
		if err != nil {
			middleware.AddErrorToContext(c, customerrors.NewValidationError(err.Error()))
			return
		}

		namespaceHandler(func(controller controllers.NamespaceController, c *gin.Context) (interface{}, error) {
			return controller.GetNamespaces(limit, page)
		})(c)
	}
}

func GetNamespace() gin.HandlerFunc {
	return func(c *gin.Context) {
		var namespaceUri types.NamespaceUri
		if err := c.BindUri(&namespaceUri); err != nil {
			middleware.AddErrorToContext(c, customerrors.NewValidationError(err.Error()))
			return
		}

		namespaceHandler(func(controller controllers.NamespaceController, c *gin.Context) (interface{}, error) {
			return controller.GetNamespace(namespaceUri.NamespaceName)
		})(c)
	}
}

func CreateNamespace() gin.HandlerFunc {
	return func(c *gin.Context) {
		var namespace types.Namespace
		if err := c.BindJSON(&namespace); err != nil {
			middleware.AddErrorToContext(c, customerrors.NewValidationError(err.Error()))
			return
		}

		namespaceHandler(func(controller controllers.NamespaceController, c *gin.Context) (interface{}, error) {
			return controller.CreateNamespace(namespace.Name)
		})(c)
	}
}

func DeleteNamespace() gin.HandlerFunc {
	return func(c *gin.Context) {
		var namespaceUri types.NamespaceUri
		if err := c.BindUri(&namespaceUri); err != nil {
			middleware.AddErrorToContext(c, customerrors.NewValidationError(err.Error()))
			return
		}

		namespaceHandler(func(controller controllers.NamespaceController, c *gin.Context) (interface{}, error) {
			name := namespaceUri.NamespaceName
			message := fmt.Sprintf("Deleted namespace successfully %q", name)
			return gin.H{"message": message}, controller.DeleteNamespace(name)
		})(c)
	}
}
