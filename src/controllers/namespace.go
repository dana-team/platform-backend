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

type NamespaceController interface {
	GetNamespaces() (types.NamespaceList, error)
	GetNamespace(name string) (types.Namespace, error)
	CreateNamespace(name string) (types.Namespace, error)
	DeleteNamespace(name string) error
}

type namespaceController struct {
	client *kubernetes.Clientset
	ctx    context.Context
	logger *zap.Logger
}

func New(client *kubernetes.Clientset, context context.Context, logger *zap.Logger) NamespaceController {
	return &namespaceController{
		logger: logger,
		client: client,
		ctx:    context,
	}
}

func (n *namespaceController) GetNamespaces() (types.NamespaceList, error) {
	n.logger.Debug("Trying to fetch all namespaces")

	namespaceList := types.NamespaceList{}
	namespaces, err := n.client.CoreV1().Namespaces().List(n.ctx, metav1.ListOptions{})
	if err != nil {
		n.logger.Error(fmt.Sprintf("Could not fetch namespaces with error: %s", err.Error()))
		return namespaceList, err
	}

	n.logger.Debug("Fetched namespaces successfully")

	for _, namespace := range namespaces.Items {
		namespaceList.Namespaces = append(namespaceList.Namespaces, types.Namespace{Name: namespace.Name})
	}
	namespaceList.Count = len(namespaceList.Namespaces)
	return namespaceList, nil
}

func (n *namespaceController) GetNamespace(name string) (types.Namespace, error) {
	n.logger.Debug(fmt.Sprintf("Trying to fetch namespace: %q", name))

	namespace, err := n.client.CoreV1().Namespaces().Get(n.ctx, name, metav1.GetOptions{})
	if err != nil {
		n.logger.Error(fmt.Sprintf("Could not fetch namespace %q with error: %s", name, err.Error()))
	}

	n.logger.Debug(fmt.Sprintf("Fetched namespace %q successfully", name))
	return types.Namespace{Name: namespace.Name}, err
}

func (n *namespaceController) CreateNamespace(name string) (types.Namespace, error) {
	n.logger.Debug(fmt.Sprintf("Trying to create namespace: %q", name))

	newNamespace := v1.Namespace{}
	newNamespace.Name = name
	namespace, err := n.client.CoreV1().Namespaces().Create(n.ctx, &newNamespace, metav1.CreateOptions{})
	if err != nil {
		n.logger.Error(fmt.Sprintf("Could not create namespace %q with error: %s", name, err.Error()))
		return types.Namespace{}, err
	}

	n.logger.Debug(fmt.Sprintf("Created namespace %q successfully", name))
	return types.Namespace{Name: namespace.Name}, err
}

func (n *namespaceController) DeleteNamespace(name string) error {
	n.logger.Debug(fmt.Sprintf("Trying to delete namespace: %q", name))

	if err := n.client.CoreV1().Namespaces().Delete(n.ctx, name, metav1.DeleteOptions{}); err != nil {
		n.logger.Debug(fmt.Sprintf("Deleted namespace %q successfully", name))
	}
	return nil
}
