package controllers

import (
	"context"
	"fmt"
	"github.com/dana-team/platform-backend/src/types"
	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type NameSpaceController interface {
	GetNamespaces() (types.OutputNamespaces, error)
	GetNamespace(nsName string) (types.Namespace, error)
	CreateNamespace(nsName string) (types.Namespace, error)
	DeleteNamespace(nsName string) error
}

type namespaceController struct {
	client *kubernetes.Clientset
	ctx    context.Context
	logger *zap.Logger
}

func NewNSController(client *kubernetes.Clientset, logger *zap.Logger) NameSpaceController {
	return &namespaceController{
		logger: logger,
		client: client,
		ctx:    context.Background(),
	}
}

func (n *namespaceController) GetNamespaces() (types.OutputNamespaces, error) {

	namespaces, err := n.client.CoreV1().Namespaces().List(n.ctx, metav1.ListOptions{})
	n.logger.Debug("trying to fetch all namespaces")
	if err != nil {
		n.logger.Error(fmt.Sprintf("Could not fetch namespaces with error: %s", err.Error()))
		return types.OutputNamespaces{}, err
	}
	n.logger.Debug("fetched namespaces successfully")
	outputNamespaces := types.OutputNamespaces{}

	for _, namespace := range namespaces.Items {
		outputNamespaces.Namespaces = append(outputNamespaces.Namespaces, types.Namespace{Name: namespace.Name})
	}
	outputNamespaces.Count = len(outputNamespaces.Namespaces)
	return outputNamespaces, nil
}

func (n *namespaceController) GetNamespace(nsName string) (types.Namespace, error) {
	outputNS := types.Namespace{}
	n.logger.Debug(fmt.Sprintf("Trying to fetch namespace: %q", nsName))
	namespace, err := n.client.CoreV1().Namespaces().Get(n.ctx, nsName, metav1.GetOptions{})
	if err != nil {
		n.logger.Error(fmt.Sprintf("Could not fetch namespace %q with error: %s", nsName, err.Error()))
	}
	n.logger.Debug(fmt.Sprintf("fetched namespace %q successfully", nsName))
	outputNS.Name = namespace.Name
	return outputNS, err

}

func (n *namespaceController) CreateNamespace(nsName string) (types.Namespace, error) {
	outputNS := types.Namespace{}
	namespace := v1.Namespace{}
	namespace.Name = nsName
	n.logger.Debug(fmt.Sprintf("Trying to create namespace: %q", nsName))
	createdNamespace, err := n.client.CoreV1().Namespaces().Create(n.ctx, &namespace, metav1.CreateOptions{})
	if err != nil {
		n.logger.Error(fmt.Sprintf("Could not create namespace %q with error: %s", nsName, err.Error()))
		return types.Namespace{}, err
	}
	outputNS.Name = createdNamespace.Name
	n.logger.Debug(fmt.Sprintf("created namespace %q successfully", nsName))
	return outputNS, err
}

func (n *namespaceController) DeleteNamespace(nsName string) error {
	n.logger.Debug(fmt.Sprintf("Trying to delete namespace: %q", nsName))
	if err := n.client.CoreV1().Namespaces().Delete(n.ctx, nsName, metav1.DeleteOptions{}); err != nil {
		n.logger.Debug(fmt.Sprintf("deleted namespace %q successfully", nsName))
	}
	return nil
}
