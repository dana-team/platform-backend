package mocks

import (
	"github.com/dana-team/platform-backend/src/types"
	"github.com/dana-team/platform-backend/src/utils/testutils"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PrepareSecret returns a mock CappRevision object.
func PrepareSecret(name, namespace, dataKey, dataValue string) corev1.Secret {
	return corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    map[string]string{testutils.ManagedLabel: "true"},
		},
		Type: corev1.SecretTypeOpaque,
		Data: map[string][]byte{
			dataKey: []byte(dataValue),
		},
	}
}

// PrepareCreateSecretRequestType returns a Secret Request type object.
func PrepareCreateSecretRequestType(name, secretType, cert, key string, data []types.KeyValue) types.CreateSecretRequest {
	return types.CreateSecretRequest{
		Type:       secretType,
		SecretName: name,
		Cert:       cert,
		Key:        key,
		Data:       data,
	}
}

// PrepareSecretRequestType returns a Secret Request type object.
func PrepareSecretRequestType(data []types.KeyValue) types.UpdateSecretRequest {
	return types.UpdateSecretRequest{
		Data: data,
	}
}
