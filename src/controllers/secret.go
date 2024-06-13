package controllers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/dana-team/platform-backend/src/types"
	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
)

type SecretController interface {
	// CreateSecret creates a new secret in the specified namespace.
	CreateSecret(namespace string, request types.CreateSecretRequest) (types.CreateSecretResponse, error)

	// GetSecrets gets all secretes from the specified namespace.
	GetSecrets(namespace string) (types.GetSecretsResponse, error)

	// GetSecret gets a specific secret from the specified namespace.
	GetSecret(namespace string, name string) (types.GetSecretResponse, error)

	// PatchSecret patches a specific secret in the specified namespace.
	PatchSecret(namespace string, name string, request types.PatchSecretRequest) (types.PatchSecretResponse, error)

	// DeleteSecret deletes a specific secret in the specified namespace.
	DeleteSecret(namespace string, name string) (types.DeleteSecretResponse, error)
}

type secretController struct {
	client kubernetes.Interface
	ctx    context.Context
	logger *zap.Logger
}

// NewSecretController creates a new secrets controller to make requests of
// CRUD operations using Kubernetes API.
func NewSecretController(client kubernetes.Interface, context context.Context, logger *zap.Logger) SecretController {
	return &secretController{
		client: client,
		ctx:    context,
		logger: logger,
	}
}

func (n *secretController) CreateSecret(namespace string, request types.CreateSecretRequest) (types.CreateSecretResponse, error) {
	n.logger.Debug(fmt.Sprintf("Trying to create secret: %q", request.SecretName))

	response := types.CreateSecretResponse{}
	name := request.SecretName

	secret, err := newSecretFromRequest(namespace, request)
	if err != nil {
		return response, err
	}

	newSecret, err := n.client.CoreV1().Secrets(namespace).Create(n.ctx, secret, metav1.CreateOptions{})
	if err != nil {
		n.logger.Error(fmt.Sprintf("Could not create secret %q with error: %v", name, err.Error()))
		return response, err
	}

	n.logger.Debug(fmt.Sprintf("Created secret %q successfully", name))

	response.Type = string(newSecret.Type)
	response.SecretName = newSecret.Name
	response.NamespaceName = newSecret.Namespace
	return response, err
}

func (n *secretController) GetSecrets(namespace string) (types.GetSecretsResponse, error) {
	n.logger.Debug(fmt.Sprintf("Trying to get all secrets in %q namespace", namespace))

	response := types.GetSecretsResponse{}
	secrets, err := n.client.CoreV1().Secrets(namespace).List(n.ctx, metav1.ListOptions{})
	if err != nil {
		n.logger.Error(fmt.Sprintf("Could not get secrets with error: %v", err.Error()))
		return response, err
	}

	n.logger.Debug("Fetched all secrets successfully")

	response.Count = len(secrets.Items)
	for _, secret := range secrets.Items {
		response.Secrets = append(
			response.Secrets,
			types.Secret{
				Type:          string(secret.Type),
				SecretName:    secret.Name,
				NamespaceName: secret.Namespace,
			})
	}
	return response, nil
}

func (n *secretController) GetSecret(namespace string, name string) (types.GetSecretResponse, error) {
	n.logger.Debug(fmt.Sprintf("Trying to get %q secret in %q namespace", name, namespace))

	response := types.GetSecretResponse{}
	secret, err := n.client.CoreV1().Secrets(namespace).Get(n.ctx, name, metav1.GetOptions{})
	if err != nil {
		n.logger.Error(fmt.Sprintf("Could not get secret %q with error: %v", name, err.Error()))
		return response, err
	}

	n.logger.Debug("Fetched secret successfully")

	response.Id = string(secret.UID)
	response.Type = string(secret.Type)
	response.SecretName = secret.Name
	for k, v := range secret.Data {
		value, err := base64.StdEncoding.DecodeString(string(v))
		if err != nil {
			n.logger.Error(fmt.Sprintf("Error decoding secret %q value: %v", name, err.Error()))
			continue
		}

		response.Data = append(response.Data, types.KeyValue{
			Key:   k,
			Value: string(value),
		})
	}
	return response, nil
}

func (n *secretController) PatchSecret(namespace string, name string, request types.PatchSecretRequest) (types.PatchSecretResponse, error) {
	n.logger.Debug(fmt.Sprintf("Trying to patch an existing secret in %q namespace", namespace))

	response := types.PatchSecretResponse{}
	secret := v1.Secret{
		Data: map[string][]byte{},
	}
	for _, kv := range request.Data {
		secret.Data[kv.Key] = []byte(kv.Value)
	}

	data, err := json.Marshal(secret)
	if err != nil {
		n.logger.Error(fmt.Sprintf("Could not marshal secret %q with error: %v", name, err.Error()))
		return response, err
	}

	result, err := n.client.CoreV1().Secrets(namespace).Patch(n.ctx, name, k8stypes.StrategicMergePatchType, data, metav1.PatchOptions{})
	if err != nil {
		n.logger.Error(fmt.Sprintf("Could not patch secret %q with error: %v", name, err.Error()))
		return response, err
	}

	response.Id = string(result.UID)
	response.Type = string(result.Type)
	response.SecretName = result.Name
	for k, v := range secret.Data {
		value, err := base64.StdEncoding.DecodeString(string(v))
		if err != nil {
			n.logger.Error(fmt.Sprintf("Error decoding secret %q value: %v", name, err.Error()))
			continue
		}

		response.Data = append(response.Data, types.KeyValue{
			Key:   k,
			Value: string(value),
		})
	}
	return response, nil
}

func (n *secretController) DeleteSecret(namespace string, name string) (types.DeleteSecretResponse, error) {
	n.logger.Debug(fmt.Sprintf("Trying to delete an existing secret in %q namespace", namespace))

	response := types.DeleteSecretResponse{}
	if err := n.client.CoreV1().Secrets(namespace).Delete(n.ctx, name, metav1.DeleteOptions{}); err != nil {
		n.logger.Error(fmt.Sprintf("Could not delete secret %q with error: %v", name, err.Error()))
		return response, err
	}

	response.Message = fmt.Sprintf("Secret %q was deleted successfully", name)
	return response, nil
}

// createSecretFromRequest returns a new secret based on different secret
// types, either TLS or Opaque.
func newSecretFromRequest(namespace string, request types.CreateSecretRequest) (*v1.Secret, error) {
	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      request.SecretName,
			Namespace: namespace,
		},
	}

	switch request.Type {
	case string(v1.SecretTypeTLS):
		if request.Cert == "" || request.Key == "" {
			return secret, fmt.Errorf("cert and key are required for TLS secrets")
		}
		secret.Type = v1.SecretTypeTLS
		secret.Data = map[string][]byte{
			"tls.crt": []byte(base64.StdEncoding.EncodeToString([]byte(request.Cert))),
			"tls.key": []byte(base64.StdEncoding.EncodeToString([]byte(request.Key))),
		}
	case string(v1.SecretTypeOpaque):
		if len(request.Data) == 0 {
			return secret, fmt.Errorf("data is required for Opaque secrets")
		}
		secret.Type = v1.SecretTypeOpaque
		secret.Data = map[string][]byte{}
		for _, kv := range request.Data {
			secret.Data[kv.Key] = []byte(base64.StdEncoding.EncodeToString([]byte(kv.Value)))
		}
	default:
		return secret, fmt.Errorf("unsupported secret type: %q", request.Type)
	}
	return secret, nil
}
