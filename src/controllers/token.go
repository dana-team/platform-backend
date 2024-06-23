package controllers

import (
	"context"
	"fmt"

	"github.com/dana-team/platform-backend/src/types"
	"go.uber.org/zap"
	authenticationv1 "k8s.io/api/authentication/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type TokenController interface {
	// CreateServiceAccount creates a new service account in the specified namespace.
	CreateServiceAccount(namespace string, request types.ServiceAccount) (types.ServiceAccount, error)

	// GetToken gets a new token from the service account.
	GetToken(namespace string, request types.ServiceAccount) (types.Token, error)
}

type tokenController struct {
	client kubernetes.Interface
	ctx    context.Context
	logger *zap.Logger
}

// NewTokenController creates a new token controller to get a new token from a specific
// service account.
func NewTokenController(client kubernetes.Interface, context context.Context, logger *zap.Logger) TokenController {
	return &tokenController{
		client: client,
		ctx:    context,
		logger: logger,
	}
}

func (t *tokenController) CreateServiceAccount(namespace string, request types.ServiceAccount) (types.ServiceAccount, error) {
	t.logger.Debug(fmt.Sprintf("Trying to create a new service account: %q", request.ServiceAccountName))

	serviceAccount := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      request.ServiceAccountName,
			Namespace: namespace,
		},
	}
	serviceAccount, err := t.client.CoreV1().ServiceAccounts(namespace).Create(t.ctx, serviceAccount, metav1.CreateOptions{})
	if err != nil {
		t.logger.Error(fmt.Sprintf("Could not create service account %q with error: %v", request.ServiceAccountName, err.Error()))
		return types.ServiceAccount{}, err
	}

	t.logger.Debug(fmt.Sprintf("Created a new service account %q successfully", serviceAccount.Name))

	return types.ServiceAccount{
		ServiceAccountName: serviceAccount.Name,
	}, err
}

func (t *tokenController) GetToken(namespace string, request types.ServiceAccount) (types.Token, error) {
	t.logger.Debug(fmt.Sprintf("Trying to get a new token from service account: %q", request.ServiceAccountName))

	expirationSeconds := int64(100 * 365 * 24 * 60 * 60) // 100 years
	tokenRequest := &authenticationv1.TokenRequest{
		Spec: authenticationv1.TokenRequestSpec{
			Audiences:         nil,
			ExpirationSeconds: &expirationSeconds,
			BoundObjectRef:    nil,
		},
	}
	token, err := t.client.CoreV1().ServiceAccounts(namespace).CreateToken(t.ctx, request.ServiceAccountName, tokenRequest, metav1.CreateOptions{})
	if err != nil {
		t.logger.Error(fmt.Sprintf("Could not create token %q with error: %v", request.ServiceAccountName, err.Error()))
		return types.Token{}, err
	}

	t.logger.Debug(fmt.Sprintf("Created a new token from service account %q successfully", request.ServiceAccountName))

	return types.Token{
		Token: token.Status.Token,
	}, err
}
