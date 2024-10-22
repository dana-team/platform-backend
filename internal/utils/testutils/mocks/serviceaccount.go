package mocks

import (
	"github.com/dana-team/platform-backend/internal/utils/testutils"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PrepareServiceAccount returns a mock ServiceAccount object.
func PrepareServiceAccount(name, namespace string, dockerCfgSecretName string) *corev1.ServiceAccount {
	serviceAccount := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				testutils.ManagedLabel: "true",
			},
		},
	}

	if dockerCfgSecretName != "" {
		serviceAccount.Secrets = []corev1.ObjectReference{
			{
				Kind:      testutils.Secret,
				Name:      dockerCfgSecretName,
				Namespace: namespace,
			},
		}
	}

	return serviceAccount
}
