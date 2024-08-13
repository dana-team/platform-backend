package v1

import (
	"github.com/dana-team/platform-backend/src/controllers"
	"github.com/dana-team/platform-backend/src/middleware"
	"github.com/dana-team/platform-backend/src/types"
	"github.com/dana-team/platform-backend/src/utils/pagination"
	"github.com/gin-gonic/gin"
	multicluster "github.com/oam-dev/cluster-gateway/pkg/apis/cluster/transport"
	"go.uber.org/zap"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
	"net/http"
)

// podHandler wraps a handler function with context setup for PodController.
func podHandler(handler func(controller controllers.PodController, c *gin.Context) (interface{}, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		client, exists := c.Get(middleware.KubeClientCtxKey)
		if !exists {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Kubernetes client not found"})
			return
		}

		ctxLogger, exists := c.Get(middleware.LoggerCtxKey)
		if !exists {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Logger not found in context"})
			return
		}

		logger := ctxLogger.(*zap.Logger)
		kubeClient := client.(kubernetes.Interface)
		clusterName := c.Param(clusterNameParam)
		context := multicluster.WithMultiClusterContext(c.Request.Context(), clusterName)

		podController := controllers.NewPodController(kubeClient, context, logger)
		result, err := handler(podController, c)
		if err != nil {
			c.AbortWithStatusJSON(int(err.(*k8serrors.StatusError).ErrStatus.Code), gin.H{"error": "Operation failed", "details": err.Error()})
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
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
			return
		}

		limit, page, err := pagination.ExtractPaginationParamsFromCtx(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
			return
		}

		podHandler(func(controller controllers.PodController, c *gin.Context) (interface{}, error) {
			return controller.GetPods(request.NamespaceName, request.CappName, limit, page)
		})(c)
	}
}
