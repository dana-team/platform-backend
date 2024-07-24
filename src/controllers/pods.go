package controllers

import (
	"context"
	"fmt"
	"github.com/dana-team/platform-backend/src/types"
	"github.com/dana-team/platform-backend/src/utils"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
)

// PodController defines methods to interact with pod pods.
type PodController interface {
	GetPods(namespace, cappName string) (types.GetPodsResponse, error)
}

// podController implements the PodController interface.
type podController struct {
	client kubernetes.Interface
	ctx    context.Context
	logger *zap.Logger
}

// NewPodController creates a new instance of PodController.
func NewPodController(client kubernetes.Interface, context context.Context, logger *zap.Logger) PodController {
	return &podController{
		client: client,
		ctx:    context,
		logger: logger,
	}
}

// GetPods returns a list of pod names for a given capp in a specific namespace.
func (n *podController) GetPods(namespace, cappName string) (types.GetPodsResponse, error) {
	n.logger.Debug(fmt.Sprintf("Trying to get all pods in %q namespace", namespace))

	pods, err := utils.GetPodsByLabel(n.ctx, n.client, namespace, fmt.Sprintf(utils.ParentCappLabelSelector, cappName))
	if err != nil {
		n.logger.Error(fmt.Sprintf("error fetching Capp pods: %s", err.Error()))
		return types.GetPodsResponse{}, err
	}

	response := types.GetPodsResponse{}
	response.Count = len(pods.Items)
	for _, pod := range pods.Items {
		response.Pods = append(
			response.Pods,
			types.Pod{
				PodName: pod.Name,
			})
	}

	n.logger.Debug("Fetched all pods successfully")
	return response, nil
}
