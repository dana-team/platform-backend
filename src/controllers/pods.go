package controllers

import (
	"context"
	"fmt"

	"github.com/dana-team/platform-backend/src/customerrors"
	"github.com/dana-team/platform-backend/src/types"
	"github.com/dana-team/platform-backend/src/utils"
	"github.com/dana-team/platform-backend/src/utils/pagination"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	ErrCouldNotGetPods = "Could not get pods"
)

// PodController defines methods to interact with pods.
type PodController interface {
	GetPods(namespace, cappName string, limit, page int) (types.GetPodsResponse, error)
}

// podController implements the PodController interface.
type podController struct {
	client kubernetes.Interface
	ctx    context.Context
	logger *zap.Logger
}

// PodPaginator paginates through pods in a specified namespace.
type PodPaginator struct {
	pagination.GenericPaginator
	namespace string
	client    kubernetes.Interface
	cappName  string
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
func (n *podController) GetPods(namespace, cappName string, limit, page int) (types.GetPodsResponse, error) {
	n.logger.Debug(fmt.Sprintf("Trying to get all pods in %q namespace", namespace))

	podPaginator := &PodPaginator{
		GenericPaginator: pagination.CreatePaginator(n.ctx, n.logger),
		namespace:        namespace,
		client:           n.client,
		cappName:         cappName,
	}

	pods, err := pagination.FetchPage[corev1.Pod](limit, page, podPaginator)
	if err != nil {
		n.logger.Error(fmt.Sprintf("%v with error: %v", ErrCouldNotGetPods, err))
		return types.GetPodsResponse{}, customerrors.NewAPIError(ErrCouldNotGetPods, err)
	}

	response := types.GetPodsResponse{}
	response.Count = len(pods)
	for _, pod := range pods {
		response.Pods = append(
			response.Pods,
			types.Pod{
				PodName: pod.Name,
			})
	}

	n.logger.Debug("Fetched all pods successfully")
	return response, nil
}

// FetchList retrieves a list of pods from the specified namespace with given options.
func (p *PodPaginator) FetchList(listOptions metav1.ListOptions) (*types.List[corev1.Pod], error) {
	pods, err := utils.GetPodsByLabel(p.Ctx, p.client, p.namespace, fmt.Sprintf(utils.ParentCappLabelSelector, p.cappName), metav1.ListOptions{
		Limit:    listOptions.Limit,
		Continue: listOptions.Continue,
	})
	if err != nil {
		p.Logger.Error(fmt.Sprintf("%v: %s", errFetchingCappPods, err.Error()))
		return nil, customerrors.NewAPIError(errFetchingCappPods, err)
	}

	return (*types.List[corev1.Pod])(pods), nil
}
