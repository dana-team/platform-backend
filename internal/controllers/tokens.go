package controllers

import (
	"context"
	"fmt"
	"strconv"

	"k8s.io/apimachinery/pkg/api/errors"

	"github.com/dana-team/platform-backend/internal/types"

	corev1 "k8s.io/api/core/v1"

	"github.com/dana-team/platform-backend/internal/customerrors"
	"github.com/dana-team/platform-backend/internal/utils"
	authenticationv1 "k8s.io/api/authentication/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
)

const (
	secretKind         = "Secret"
	tokenRequestSuffix = "token-request"
	serviceAccountKind = "ServiceAccount"
	defaultExpiration  = "36000"
)

const (
	ErrCouldNotCreateTokenRequestSecret = "Could not create token request secret for serviceaccount %q in namespace %q"
	ErrCouldNotGetTokenRequestSecret    = "Could not get token request secret for serviceaccount %q in namespace %q"
	ErrCouldNotDeleteTokenRequestSecret = "Could not delete token request secret for serviceaccount %q in namespace %q"
	ErrCouldNotCreateTokenRequest       = "Could not create token request for serviceaccount %q in namespace %q"
	ErrInvalidExpirationSeconds         = "Invalid expiration seconds: %s"
)

type tokenController struct {
	client kubernetes.Interface
	ctx    context.Context
	logger *zap.Logger
}

type TokenController interface {
	// RevokeToken revokes a token by name of the serviceaccount.
	RevokeToken(serviceAccountName, namespace string) error

	// CreateToken creates a token for a serviceaccount.
	CreateToken(serviceAccountName, namespace string, expirationSeconds string) (types.TokenRequestResponse, error)
}

// NewTokenController creates a new token controller.
func NewTokenController(client kubernetes.Interface, ctx context.Context, logger *zap.Logger) TokenController {
	return &tokenController{
		client: client,
		ctx:    ctx,
		logger: logger,
	}
}

func (t *tokenController) RevokeToken(serviceAccountName, namespace string) error {
	serviceAccount, err := t.client.CoreV1().ServiceAccounts(namespace).Get(t.ctx, serviceAccountName, metav1.GetOptions{})
	if err != nil {
		t.logger.Error(fmt.Sprintf("%s with error: %s", fmt.Sprintf(ErrCouldNotGetServiceAccount, serviceAccountName, namespace), err.Error()))
		return customerrors.NewAPIError(fmt.Sprintf(ErrCouldNotGetServiceAccount, serviceAccountName, namespace), err)
	}

	_, err = t.client.CoreV1().Secrets(namespace).Get(t.ctx, fmt.Sprintf("%s-%s", serviceAccount.Name, tokenRequestSuffix), metav1.GetOptions{})
	if err != nil {
		t.logger.Error(fmt.Sprintf("%s with error: %s", fmt.Sprintf(ErrCouldNotGetTokenRequestSecret, serviceAccountName, namespace), err.Error()))
		return customerrors.NewAPIError(fmt.Sprintf(ErrCouldNotGetTokenRequestSecret, serviceAccountName, namespace), err)
	}
	if err = DeleteTokenRequestSecret(t.ctx, t.client, serviceAccountName, namespace); err != nil {
		t.logger.Error(fmt.Sprintf("%s with error: %s", fmt.Sprintf(ErrCouldNotDeleteTokenRequestSecret, serviceAccountName, namespace), err.Error()))
		return customerrors.NewAPIError(fmt.Sprintf(ErrCouldNotDeleteTokenRequestSecret, serviceAccountName, namespace), err)
	}
	return nil
}

func (t *tokenController) CreateToken(serviceAccountName, namespace string, expirationSeconds string) (types.TokenRequestResponse, error) {
	tokenRequestSecret, err := t.client.CoreV1().Secrets(namespace).Get(t.ctx, fmt.Sprintf("%s-%s", serviceAccountName, tokenRequestSuffix), metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			createdSecret, createSecretErr := CreateTokenRequestSecret(t.ctx, t.client, serviceAccountName, namespace, t.logger)
			if createSecretErr != nil {
				t.logger.Error(fmt.Sprintf("%s with error: %s", fmt.Sprintf(ErrCouldNotCreateTokenRequestSecret, serviceAccountName, namespace), createSecretErr.Error()))
				return types.TokenRequestResponse{}, customerrors.NewAPIError(fmt.Sprintf(ErrCouldNotCreateTokenRequestSecret, serviceAccountName, namespace), createSecretErr)
			}
			tokenRequestSecret = createdSecret
		} else {
			t.logger.Error(fmt.Sprintf("%s with error: %s", fmt.Sprintf(ErrCouldNotGetTokenRequestSecret, serviceAccountName, namespace), err.Error()))
			return types.TokenRequestResponse{}, customerrors.NewAPIError(fmt.Sprintf(ErrCouldNotGetTokenRequestSecret, serviceAccountName, namespace), err)
		}
	}
	expireIn, err := getExpirationSeconds(expirationSeconds)
	if err != nil {
		t.logger.Error(fmt.Sprintf("%s with error: %s", fmt.Sprintf(ErrCouldNotCreateTokenRequest, serviceAccountName, namespace), err.Error()))
		return types.TokenRequestResponse{}, customerrors.NewValidationError(err.Error())
	}

	tokenRequest := prepareTokenRequest(serviceAccountName, namespace, tokenRequestSecret, expireIn)

	createdRequest, err := t.client.CoreV1().ServiceAccounts(namespace).CreateToken(t.ctx, serviceAccountName, tokenRequest, metav1.CreateOptions{})
	if err != nil {
		t.logger.Error(fmt.Sprintf("%s with error: %s", fmt.Sprintf(ErrCouldNotCreateTokenRequest, serviceAccountName, namespace), err.Error()))
		return types.TokenRequestResponse{}, customerrors.NewAPIError(fmt.Sprintf(ErrCouldNotCreateTokenRequest, serviceAccountName, namespace), err)
	}
	return types.TokenRequestResponse{Token: createdRequest.Status.Token, ExpirationTimestamp: createdRequest.Status.ExpirationTimestamp.Time}, nil
}

// CreateTokenRequestSecret creates a secret for the token request.
// All token requests for a service account are bound to this secret.
// When this secret is deleted, all tokens bound to it are revoked.
func CreateTokenRequestSecret(ctx context.Context, client kubernetes.Interface, serviceAccountName, namespace string, log *zap.Logger) (*corev1.Secret, error) {
	_, err := client.CoreV1().Secrets(namespace).Get(ctx, fmt.Sprintf("%s-%s", serviceAccountName, tokenRequestSuffix), metav1.GetOptions{})
	if err != nil {
		if !errors.IsNotFound(err) {
			log.Error(fmt.Sprintf("%s with error: %s", fmt.Sprintf(ErrCouldNotGetTokenRequestSecret, serviceAccountName, namespace), err.Error()))
			return nil, customerrors.NewAPIError(fmt.Sprintf(ErrCouldNotGetTokenRequestSecret, serviceAccountName, namespace), err)
		}
	}
	serviceAccount, err := client.CoreV1().ServiceAccounts(namespace).Get(ctx, serviceAccountName, metav1.GetOptions{})
	if err != nil {
		log.Error(fmt.Sprintf("%s with error: %s", fmt.Sprintf(ErrCouldNotGetServiceAccount, serviceAccountName, namespace), err.Error()))
		return nil, customerrors.NewAPIError(fmt.Sprintf(ErrCouldNotGetServiceAccount, serviceAccountName, namespace), err)
	}
	serviceAccountSecret := prepareTokenRequestSecret(serviceAccount)
	requestSecret, err := client.CoreV1().Secrets(namespace).Create(ctx, serviceAccountSecret, metav1.CreateOptions{})
	if err != nil {
		log.Error(fmt.Sprintf("%s with error: %s", fmt.Sprintf(ErrCouldNotCreateTokenRequestSecret, serviceAccountName, namespace), err.Error()))
		return nil, customerrors.NewAPIError(fmt.Sprintf(ErrCouldNotCreateTokenRequestSecret, serviceAccountName, namespace), err)
	}
	return requestSecret, nil
}

// DeleteTokenRequestSecret deletes the secret for token requests.
func DeleteTokenRequestSecret(ctx context.Context, client kubernetes.Interface, serviceAccountName, namespace string) error {
	err := client.CoreV1().Secrets(namespace).Delete(ctx, fmt.Sprintf("%s-%s", serviceAccountName, tokenRequestSuffix), metav1.DeleteOptions{})
	if err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
	}
	return nil
}

// getExpirationSeconds returns the expiration seconds for a token as a pointer.
func getExpirationSeconds(expirationSeconds string) (*int64, error) {
	if expirationSeconds == "" {
		expirationSeconds = defaultExpiration
	}
	expireIn, err := strconv.Atoi(expirationSeconds)
	if err != nil {
		return nil, errors.NewBadRequest(fmt.Sprintf(ErrInvalidExpirationSeconds, expirationSeconds))
	}
	seconds := new(int64)
	*seconds = int64(expireIn)
	return seconds, nil
}

// prepareTokenRequest creates a token request object.
func prepareTokenRequest(name, namespace string, tokenRequestSecret *corev1.Secret, expirationSeconds *int64) *authenticationv1.TokenRequest {
	return &authenticationv1.TokenRequest{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: authenticationv1.TokenRequestSpec{
			BoundObjectRef: &authenticationv1.BoundObjectReference{
				Kind:       secretKind,
				Name:       tokenRequestSecret.Name,
				APIVersion: corev1.SchemeGroupVersion.String(),
				UID:        tokenRequestSecret.UID,
			},
			ExpirationSeconds: expirationSeconds,
		},
	}
}

// prepareTokenRequestSecret creates a secret for the token requests.
func prepareTokenRequestSecret(serviceAccount *corev1.ServiceAccount) *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s", serviceAccount.Name, tokenRequestSuffix),
			Namespace: serviceAccount.Namespace,
			Labels: map[string]string{
				utils.ManagedLabel: utils.ManagedLabelValue,
			},
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: corev1.SchemeGroupVersion.String(),
					Kind:       serviceAccountKind,
					Name:       serviceAccount.Name,
					UID:        serviceAccount.UID,
				},
			},
		},
		Type: corev1.SecretTypeOpaque,
	}
}
