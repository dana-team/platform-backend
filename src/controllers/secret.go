package controllers

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/dana-team/platform-backend/src/types"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	tlsType    = "tls"
	opaqueType = "opaque"
)

type SecretController interface {
	// CreateSecret creates a new secret in the specified namespace.
	CreateSecret(namespace string, request types.CreateSecretRequest) (types.CreateSecretResponse, error)

	// GetSecrets gets all secretes from the specified namespace.
	GetSecrets(namespace string) (types.GetSecretsResponse, error)

	// GetSecret gets a specific secret from the specified namespace.
	GetSecret(namespace, name string) (types.GetSecretResponse, error)

	// UpdateSecret updates a specific secret in the specified namespace.
	UpdateSecret(namespace, name string, request types.UpdateSecretRequest) (types.UpdateSecretResponse, error)

	// DeleteSecret deletes a specific secret in the specified namespace.
	DeleteSecret(namespace, name string) (types.DeleteSecretResponse, error)
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

	secret, err := newSecretFromRequest(namespace, request)
	if err != nil {
		return types.CreateSecretResponse{}, err
	}

	newSecret, err := n.client.CoreV1().Secrets(namespace).Create(n.ctx, secret, metav1.CreateOptions{})
	if err != nil {
		n.logger.Error(fmt.Sprintf("Could not create secret %q with error: %v", request.SecretName, err.Error()))
		return types.CreateSecretResponse{}, err
	}

	n.logger.Debug(fmt.Sprintf("Created secret %q successfully", newSecret.Name))

	return types.CreateSecretResponse{
		Type:          string(newSecret.Type),
		SecretName:    newSecret.Name,
		NamespaceName: newSecret.Namespace,
	}, err
}

func (n *secretController) GetSecrets(namespace string) (types.GetSecretsResponse, error) {
	n.logger.Debug(fmt.Sprintf("Trying to get all secrets in %q namespace", namespace))

	secrets, err := n.client.CoreV1().Secrets(namespace).List(n.ctx, metav1.ListOptions{})
	if err != nil {
		n.logger.Error(fmt.Sprintf("Could not get secrets with error: %v", err.Error()))
		return types.GetSecretsResponse{}, err
	}

	n.logger.Debug("Fetched all secrets successfully")

	response := types.GetSecretsResponse{}
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

func (n *secretController) GetSecret(namespace, name string) (types.GetSecretResponse, error) {
	n.logger.Debug(fmt.Sprintf("Trying to get %q secret in %q namespace", name, namespace))

	secret, err := n.client.CoreV1().Secrets(namespace).Get(n.ctx, name, metav1.GetOptions{})
	if err != nil {
		n.logger.Error(fmt.Sprintf("Could not get secret %q with error: %v", name, err.Error()))
		return types.GetSecretResponse{}, err
	}

	n.logger.Debug("Fetched secret successfully")

	secretData, err := decodeData(secret.Data)
	if err != nil {
		n.logger.Error(fmt.Sprintf("Error decoding secret %q value: %v", secret.Name, err.Error()))
		return types.GetSecretResponse{}, k8serrors.NewInternalError(err)

	}
	response := types.GetSecretResponse{
		Id:         string(secret.UID),
		Type:       string(secret.Type),
		SecretName: secret.Name,
		Data:       secretData,
	}

	return response, nil
}

func (n *secretController) UpdateSecret(namespace, name string, request types.UpdateSecretRequest) (types.UpdateSecretResponse, error) {
	n.logger.Debug(fmt.Sprintf("Trying to update an existing secret in %q namespace", namespace))

	secret, err := n.client.CoreV1().Secrets(namespace).Get(n.ctx, name, metav1.GetOptions{})
	if err != nil {
		n.logger.Error(fmt.Sprintf("Could not get secret %q with error: %v", name, err.Error()))
		return types.UpdateSecretResponse{}, err
	}
	secret.Data = convertKeyValueToByteMap(request.Data)

	result, err := n.client.CoreV1().Secrets(namespace).Update(n.ctx, secret, metav1.UpdateOptions{})
	if err != nil {
		n.logger.Error(fmt.Sprintf("Could not update secret %q with error: %v", name, err.Error()))
		return types.UpdateSecretResponse{}, err
	}

	secretData, err := decodeData(result.Data)
	if err != nil {
		n.logger.Error(fmt.Sprintf("Error decoding secret %q value: %v", result.Name, err.Error()))
		return types.UpdateSecretResponse{}, k8serrors.NewInternalError(err)

	}

	response := types.UpdateSecretResponse{
		Type:          string(result.Type),
		SecretName:    result.Name,
		NamespaceName: result.Namespace,
		Data:          secretData,
	}

	return response, nil
}

func (n *secretController) DeleteSecret(namespace, name string) (types.DeleteSecretResponse, error) {
	n.logger.Debug(fmt.Sprintf("Trying to delete an existing secret in %q namespace", namespace))

	if err := n.client.CoreV1().Secrets(namespace).Delete(n.ctx, name, metav1.DeleteOptions{}); err != nil {
		n.logger.Error(fmt.Sprintf("Could not delete secret %q with error: %v", name, err.Error()))
		return types.DeleteSecretResponse{}, err
	}

	return types.DeleteSecretResponse{
		Message: fmt.Sprintf("Deleted secret %q in namespace %q successfully", name, namespace),
	}, nil
}

// createSecretFromRequest returns a new secret based on different secret
// types, either TLS or Opaque.
func newSecretFromRequest(namespace string, request types.CreateSecretRequest) (*corev1.Secret, error) {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      request.SecretName,
			Namespace: namespace,
		},
	}

	switch request.Type {
	case tlsType:
		if request.Cert == "" || request.Key == "" {
			return secret, k8serrors.NewBadRequest("cert and key are required for TLS secrets")
		}
		secret.Type = corev1.SecretTypeTLS
		secret.Data = map[string][]byte{
			"tls.crt": []byte(base64.StdEncoding.EncodeToString([]byte(request.Cert))),
			"tls.key": []byte(base64.StdEncoding.EncodeToString([]byte(request.Key))),
		}
	case opaqueType:
		if len(request.Data) == 0 {
			return secret, k8serrors.NewBadRequest("data is required for Opaque secrets")
		}
		secret.Type = corev1.SecretTypeOpaque
		secret.Data = map[string][]byte{}
		for _, kv := range request.Data {
			secret.Data[kv.Key] = []byte(base64.StdEncoding.EncodeToString([]byte(kv.Value)))
		}
	default:
		return secret, k8serrors.NewBadRequest("unsupported secret type")
	}
	return secret, nil
}

// convertKeyValueToByteMap converts a slice of KeyValue pairs
// to a map with string keys and byte slice values.
func convertKeyValueToByteMap(kvList []types.KeyValue) map[string][]byte {
	data := map[string][]byte{}
	for _, kv := range kvList {
		data[kv.Key] = []byte(base64.StdEncoding.EncodeToString([]byte(kv.Value)))
	}
	return data
}

// decodeData decodes secret data.
func decodeData(encodedData map[string][]byte) ([]types.KeyValue, error) {
	var decodedData []types.KeyValue
	for k, v := range encodedData {
		value, err := base64.StdEncoding.DecodeString(string(v))
		if err != nil {
			return decodedData, err
		}

		decodedData = append(decodedData, types.KeyValue{
			Key:   k,
			Value: string(value),
		})
	}
	return decodedData, nil
}
