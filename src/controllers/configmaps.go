package controllers

import (
	"context"
	"fmt"

	"github.com/dana-team/platform-backend/src/types"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type ConfigMapController interface {
	// GetConfigMap gets a specific config map from the specified namespace.
	GetConfigMap(namespace, name string) (types.ConfigMap, error)
}

type configMapController struct {
	client kubernetes.Interface
	ctx    context.Context
	logger *zap.Logger
}

// NewConfigMapController creates a new config map controller to get them from Kubernetes API.
func NewConfigMapController(client kubernetes.Interface, context context.Context, logger *zap.Logger) ConfigMapController {
	return &configMapController{
		client: client,
		ctx:    context,
		logger: logger,
	}
}

func (c *configMapController) GetConfigMap(namespace string, name string) (types.ConfigMap, error) {
	c.logger.Debug(fmt.Sprintf("Trying to get a config map: %q", name))

	configMap, err := c.client.CoreV1().ConfigMaps(namespace).Get(c.ctx, name, metav1.GetOptions{})
	if err != nil {
		c.logger.Error(fmt.Sprintf("Could not get config map %q with error: %v", name, err.Error()))
		return types.ConfigMap{}, err
	}

	c.logger.Debug(fmt.Sprintf("Got config map %q successfully", name))

	return types.ConfigMap{
		Data: convertMapToKeyValue(configMap.Data),
	}, nil
}
