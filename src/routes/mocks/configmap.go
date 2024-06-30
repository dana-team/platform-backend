package mocks

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PrepareConfigMap returns a mock ConfigMap object.
func PrepareConfigMap(name, namespace string, data map[string]string) corev1.ConfigMap {
	configMap := corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: data,
	}

	return configMap
}
