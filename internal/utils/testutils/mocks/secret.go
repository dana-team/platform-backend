package mocks

import (
	"github.com/dana-team/platform-backend/internal/types"
	"github.com/dana-team/platform-backend/internal/utils/testutils"
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

// PrepareTokenSecret returns a mock token secret object.
func PrepareTokenSecret(name, namespace, tokenValue, serviceAccountName string) corev1.Secret {
	return corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Annotations: map[string]string{
				testutils.ServiceAccountAnnotation: serviceAccountName,
			},
		},
		Type: corev1.SecretTypeServiceAccountToken,
		Data: map[string][]byte{
			testutils.TokenKey: []byte(tokenValue),
		},
	}
}

// PrepareDockerConfigSecret returns a mock docker config secret object.
func PrepareDockerConfigSecret(name, namespace, serviceAccountTokenSecretName string) corev1.Secret {
	ownerReference := metav1.OwnerReference{
		APIVersion: testutils.V1,
		Kind:       testutils.Secret,
		Name:       serviceAccountTokenSecretName,
	}

	return corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:            name,
			Namespace:       namespace,
			OwnerReferences: []metav1.OwnerReference{ownerReference},
		},
		Type: corev1.SecretTypeDockercfg,
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
