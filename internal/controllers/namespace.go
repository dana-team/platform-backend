package controllers

import (
	"context"
	"fmt"
	"github.com/dana-team/platform-backend/internal/customerrors"
	"github.com/dana-team/platform-backend/internal/types"
	"github.com/dana-team/platform-backend/internal/utils"
	"github.com/dana-team/platform-backend/internal/utils/pagination"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	ErrCouldNotGetNamespaces   = "Could not get namespaces"
	ErrCouldNotFetchNamespace  = "Could not fetch namespace %q"
	ErrCouldNotCreateNamespace = "Could not create namespace %q"
	ErrCouldNotDeleteNamespace = "Could not delete namespace %q"
)

type NamespaceController interface {
	GetNamespaces(limit, page int) (types.NamespaceList, error)
	GetNamespace(name string) (types.Namespace, error)
	CreateNamespace(name string) (types.Namespace, error)
	DeleteNamespace(name string) error
}

type namespaceController struct {
	client kubernetes.Interface
	ctx    context.Context
	logger *zap.Logger
}

// NamespacePaginator paginates through secrets in a specified namespace.
type NamespacePaginator struct {
	pagination.GenericPaginator
	client kubernetes.Interface
}

func NewNamespaceController(client kubernetes.Interface, context context.Context, logger *zap.Logger) NamespaceController {
	return &namespaceController{
		logger: logger,
		client: client,
		ctx:    context,
	}
}

func (n *namespaceController) GetNamespaces(limit, page int) (types.NamespaceList, error) {
	namespaceList := types.NamespaceList{}
	n.logger.Debug("Trying to fetch all namespaces")

	namespacePaginator := &NamespacePaginator{
		GenericPaginator: pagination.CreatePaginator(n.ctx, n.logger),
		client:           n.client,
	}

	namespaces, err := pagination.FetchPage[corev1.Namespace](limit, page, namespacePaginator)
	if err != nil {
		n.logger.Error(fmt.Sprintf("%v with error: %v", ErrCouldNotGetNamespaces, err))
		return types.NamespaceList{}, customerrors.NewAPIError(ErrCouldNotGetNamespaces, err)
	}

	for _, namespace := range namespaces {
		namespaceList.Namespaces = append(namespaceList.Namespaces, types.Namespace{Name: namespace.Name})
	}
	namespaceList.Count = len(namespaceList.Namespaces)
	return namespaceList, nil
}

func (n *namespaceController) GetNamespace(name string) (types.Namespace, error) {
	n.logger.Debug(fmt.Sprintf("Trying to fetch namespace: %q", name))

	namespace, err := n.client.CoreV1().Namespaces().Get(n.ctx, name, metav1.GetOptions{})
	if err != nil {
		n.logger.Error(fmt.Sprintf("%v with error: %s", fmt.Sprintf(ErrCouldNotFetchNamespace, name), err.Error()))
		return types.Namespace{}, customerrors.NewAPIError(fmt.Sprintf(ErrCouldNotFetchNamespace, name), err)
	}

	n.logger.Debug(fmt.Sprintf("Fetched namespace %q successfully", name))
	return types.Namespace{Name: namespace.Name}, nil
}

func (n *namespaceController) CreateNamespace(name string) (types.Namespace, error) {
	n.logger.Debug(fmt.Sprintf("Trying to create namespace: %q", name))

	newNamespace := corev1.Namespace{}
	newNamespace.Name = name
	newNamespace.Labels = utils.AddManagedLabel(map[string]string{})
	namespace, err := n.client.CoreV1().Namespaces().Create(n.ctx, &newNamespace, metav1.CreateOptions{})
	if err != nil {
		n.logger.Error(fmt.Sprintf("%v with error: %s", fmt.Sprintf(ErrCouldNotCreateNamespace, name), err.Error()))
		return types.Namespace{}, customerrors.NewAPIError(fmt.Sprintf(ErrCouldNotCreateNamespace, name), err)
	}

	n.logger.Debug(fmt.Sprintf("Created namespace %q successfully", name))
	return types.Namespace{Name: namespace.Name}, nil
}

func (n *namespaceController) DeleteNamespace(name string) error {
	n.logger.Debug(fmt.Sprintf("Trying to delete namespace: %q", name))

	if err := n.client.CoreV1().Namespaces().Delete(n.ctx, name, metav1.DeleteOptions{}); err != nil {
		n.logger.Debug(fmt.Sprintf(ErrCouldNotDeleteNamespace, name))
		return customerrors.NewAPIError(fmt.Sprintf(ErrCouldNotDeleteNamespace, name), err)
	}
	return nil
}

// FetchList retrieves a list of secrets from the specified namespace with given options.
func (p *NamespacePaginator) FetchList(listOptions metav1.ListOptions) (*types.List[corev1.Namespace], error) {
	namespaces, err := p.client.CoreV1().Namespaces().List(p.Ctx, listOptions)
	if err != nil {
		p.Logger.Error(fmt.Sprintf("Could not fetch namespaces with error: %s", err.Error()))
		return nil, err
	}

	p.Logger.Debug("Fetched namespaces successfully")
	return (*types.List[corev1.Namespace])(namespaces), nil
}
