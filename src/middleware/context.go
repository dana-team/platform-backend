package middleware

import (
	"github.com/dana-team/platform-backend/src/customerrors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// GetKubeClient retrieves the Kubernetes client from the gin.Context.
func GetKubeClient(c *gin.Context) (kubernetes.Interface, error) {
	kube, exists := c.Get(KubeClientCtxKey)
	if !exists {
		return nil, c.Error(customerrors.NewNotFoundError("kubernetes client not found in context"))
	}
	return kube.(kubernetes.Interface), nil
}

// GetDynClient retrieves the dynamic client from the gin.Context.
func GetDynClient(c *gin.Context) (client.Client, error) {
	kube, exists := c.Get(DynamicClientCtxKey)
	if !exists {
		return nil, c.Error(customerrors.NewNotFoundError("dynamic client not found in context"))
	}
	return kube.(client.Client), nil
}

// GetConfig retrieves the config from the gin.Context.
func GetConfig(c *gin.Context) (*rest.Config, error) {
	config, exists := c.Get(ConfigKey)
	if !exists {
		return nil, c.Error(customerrors.NewNotFoundError("config not found in context"))
	}
	return config.(*rest.Config), nil
}

// GetLogger retrieves the logger from the gin.Context.
func GetLogger(c *gin.Context) (*zap.Logger, error) {
	logger, exists := c.Get(LoggerCtxKey)
	if !exists {
		return nil, c.Error(customerrors.NewNotFoundError("logger not found in context"))
	}
	return logger.(*zap.Logger), nil
}

// GetCluster retrieves the cluster from the gin.Context.
func GetCluster(c *gin.Context) (string, bool) {
	cluster, exists := c.Get(ClusterCtxKey)
	if !exists {
		return "", false
	}

	return cluster.(string), true
}

// AddErrorToContext checks if the error is non-nil and adds it to the Gin context if so.
func AddErrorToContext(c *gin.Context, err error) bool {
	if err != nil {
		// Add the error to the Gin context for further handling. Ignoring the return value as no action is needed based on it.
		_ = c.Error(err)
		return true
	}
	return false
}
