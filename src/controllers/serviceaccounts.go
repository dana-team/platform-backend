package controllers

import (
	"context"
	"fmt"

	"github.com/dana-team/platform-backend/src/utils"

	"github.com/dana-team/platform-backend/src/customerrors"
	"github.com/dana-team/platform-backend/src/types"
	"github.com/dana-team/platform-backend/src/utils/pagination"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	tokenKey = "token"
)

const (
	ErrCouldNotGetServiceAccount    = "Could not get ServiceAccount %q in namespace %q"
	ErrCouldNotGetServiceAccounts   = "Could not list ServiceAccounts"
	ErrNoTokenFound                 = "No token found for ServiceAccount %q"
	ErrCouldNotCreateServiceAccount = "Could not create ServiceAccount %q in namespace %q"
	ErrCouldNotDeleteServiceAccount = "Could not delete ServiceAccount %q in namespace %q"
)

// NewServiceAccountController creates a new instance of ServiceAccountController.
func NewServiceAccountController(client kubernetes.Interface, context context.Context, logger *zap.Logger) ServiceAccountController {
	return &serviceAccountController{
		client: client,
		ctx:    context,
		logger: logger,
	}
}

// serviceAccount implements the ServiceAccountController interface.
type serviceAccountController struct {
	client kubernetes.Interface
	ctx    context.Context
	logger *zap.Logger
}

// ServiceAccountController defines methods to interact with ServiceAccounts.
type ServiceAccountController interface {
	// GetServiceAccount retrieves a ServiceAccount by name and namespace, returning the ServiceAccount and any error encountered.
	GetServiceAccount(name, namespace string) (types.ServiceAccount, error)

	// GetServiceAccountToken retrieves the token for a given ServiceAccount by name and namespace, returning the token and any error encountered.
	GetServiceAccountToken(serviceAccountName, namespace string) (types.TokenResponse, error)

	// CreateServiceAccount creates a new ServiceAccount with the given name and namespace.
	CreateServiceAccount(name, namespace string) (types.ServiceAccount, error)

	// DeleteServiceAccount deletes a ServiceAccount by name and namespace.
	DeleteServiceAccount(name, namespace string) error

	// GetServiceAccounts retrieves all ServiceAccounts in a given namespace, returning the ServiceAccounts and any error encountered.
	GetServiceAccounts(namespace string, limit, page int) (types.ServiceAccountOutput, error)
}

// ServiceAccountPaginator paginates through secrets in a specified namespace.
type ServiceAccountPaginator struct {
	pagination.GenericPaginator
	client    kubernetes.Interface
	namespace string
}

// FetchList retrieves a list of serviceAccounts from the specified namespace with given options.
func (p *ServiceAccountPaginator) FetchList(listOptions metav1.ListOptions) (*types.List[corev1.ServiceAccount], error) {
	serviceAccounts, err := p.client.CoreV1().ServiceAccounts(p.namespace).List(p.Ctx, listOptions)
	if err != nil {
		p.Logger.Error(fmt.Sprintf("%v with error: %v", ErrCouldNotGetServiceAccounts, err.Error()))
		return nil, customerrors.NewAPIError(ErrCouldNotGetServiceAccounts, err)
	}

	p.Logger.Debug("Fetched all serviceAccounts successfully")
	return (*types.List[corev1.ServiceAccount])(serviceAccounts), nil
}

func (c *serviceAccountController) GetServiceAccount(name, namespace string) (types.ServiceAccount, error) {
	serviceAccount, err := c.client.CoreV1().ServiceAccounts(namespace).Get(c.ctx, name, metav1.GetOptions{})
	if err != nil {
		c.logger.Error(fmt.Sprintf("%s with error: %s", fmt.Sprintf(ErrCouldNotGetServiceAccount, name, namespace), err.Error()))
		return types.ServiceAccount{}, customerrors.NewAPIError(fmt.Sprintf(ErrCouldNotGetServiceAccount, name, namespace), err)
	}

	return types.ServiceAccount{Name: serviceAccount.Name}, nil
}

func (c *serviceAccountController) GetServiceAccountToken(serviceAccountName, namespace string) (types.TokenResponse, error) {
	serviceAccount, err := c.client.CoreV1().ServiceAccounts(namespace).Get(c.ctx, serviceAccountName, metav1.GetOptions{})
	if err != nil {
		c.logger.Error(fmt.Sprintf("%s with error: %s", fmt.Sprintf(ErrCouldNotGetServiceAccount, serviceAccountName, namespace), err.Error()))
		return types.TokenResponse{}, customerrors.NewAPIError(fmt.Sprintf(ErrCouldNotGetServiceAccount, serviceAccountName, namespace), err)
	}

	token, err := c.getServiceAccountToken(serviceAccount, namespace)
	if err != nil {
		return types.TokenResponse{}, err
	}

	return types.TokenResponse{
		Token: token,
	}, nil
}

// getServiceAccountToken extracts the token from a ServiceAccount's associated secrets, returning the token and any error encountered.
func (c *serviceAccountController) getServiceAccountToken(serviceAccount *corev1.ServiceAccount, namespace string) (string, error) {
	for _, ref := range serviceAccount.Secrets {
		secret, err := c.client.CoreV1().Secrets(namespace).Get(c.ctx, ref.Name, metav1.GetOptions{})
		if err != nil {
			return "", err
		}

		if secret.Type == corev1.SecretTypeDockercfg {
			for _, ownerRef := range secret.OwnerReferences {
				tokenSecret, err := c.client.CoreV1().Secrets(namespace).Get(c.ctx, ownerRef.Name, metav1.GetOptions{})
				if err != nil {
					return "", err
				}
				if tokenSecret.Type == corev1.SecretTypeServiceAccountToken {
					return string(tokenSecret.Data[tokenKey]), nil
				}
			}
		}
	}

	return "", customerrors.NewValidationError(fmt.Sprintf(ErrNoTokenFound, serviceAccount.Name))
}

func (c *serviceAccountController) CreateServiceAccount(name, namespace string) (types.ServiceAccount, error) {
	serviceAccount := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				utils.ManagedLabel: utils.ManagedLabelValue,
			},
		},
	}
	_, err := c.client.CoreV1().ServiceAccounts(namespace).Create(c.ctx, serviceAccount, metav1.CreateOptions{})
	if err != nil {
		return types.ServiceAccount{}, customerrors.NewAPIError(fmt.Sprintf(ErrCouldNotCreateServiceAccount, name, namespace), err)
	}
	return types.ServiceAccount{Name: name}, nil
}

func (c *serviceAccountController) DeleteServiceAccount(name, namespace string) error {
	err := c.client.CoreV1().ServiceAccounts(namespace).Delete(c.ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return customerrors.NewAPIError(fmt.Sprintf(ErrCouldNotDeleteServiceAccount, name, namespace), err)
	}

	return nil
}

func (c *serviceAccountController) GetServiceAccounts(namespace string, limit, page int) (types.ServiceAccountOutput, error) {
	serviceAccountOutput := types.ServiceAccountOutput{}
	c.logger.Debug(fmt.Sprintf("Trying to get all serviceaccounts in namespace %q", namespace))

	serviceAccountPaginator := &ServiceAccountPaginator{
		GenericPaginator: pagination.CreatePaginator(c.ctx, c.logger),
		namespace:        namespace,
		client:           c.client,
	}

	serviceAccounts, err := pagination.FetchPage[corev1.ServiceAccount](limit, page, serviceAccountPaginator)
	if err != nil {
		c.logger.Error(fmt.Sprintf("%v with error: %v", ErrCouldNotGetServiceAccounts, err))
		return types.ServiceAccountOutput{}, customerrors.NewAPIError(ErrCouldNotGetServiceAccounts, err)
	}
	for _, serviceAccount := range serviceAccounts {
		serviceAccountOutput.ServiceAccounts = append(serviceAccountOutput.ServiceAccounts, serviceAccount.Name)
	}
	serviceAccountOutput.Count = len(serviceAccounts)

	return serviceAccountOutput, nil
}
