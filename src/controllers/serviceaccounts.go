package controllers

import (
	"context"
	"fmt"
	"github.com/dana-team/platform-backend/src/customerrors"
	"github.com/dana-team/platform-backend/src/types"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	tokenKey = "token"
)

const (
	ErrCouldNotGetServiceAccount = "Could not get ServiceAccount %q in namespace %q"
	ErrNoTokenFound              = "no token found for ServiceAccount %s"
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
	GetServiceAccount(name, namespace string) (*corev1.ServiceAccount, error)
	GetServiceAccountToken(serviceAccountName, namespace string) (types.TokenResponse, error)
	getServiceAccountToken(serviceAccount *corev1.ServiceAccount, namespace string) (string, error)
}

// GetServiceAccount retrieves a ServiceAccount by name and namespace, returning the ServiceAccount and any error encountered.
func (c *serviceAccountController) GetServiceAccount(name, namespace string) (*corev1.ServiceAccount, error) {
	c.logger.Debug(fmt.Sprintf("Trying to get service account: %q", name))

	serviceAccount, err := c.client.CoreV1().ServiceAccounts(namespace).Get(c.ctx, name, metav1.GetOptions{})
	if err != nil {
		c.logger.Error(fmt.Sprintf("%s with error: %s", fmt.Sprintf(ErrCouldNotGetServiceAccount, name, namespace), err.Error()))
		return &corev1.ServiceAccount{}, customerrors.NewAPIError(fmt.Sprintf(ErrCouldNotGetServiceAccount, name, namespace), err)
	}

	return serviceAccount, nil
}

// GetServiceAccountToken retrieves the token for a given ServiceAccount by name and namespace, returning the token and any error encountered.
func (c *serviceAccountController) GetServiceAccountToken(serviceAccountName, namespace string) (types.TokenResponse, error) {
	serviceAccount, err := c.GetServiceAccount(serviceAccountName, namespace)
	if err != nil {
		return types.TokenResponse{}, err
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
