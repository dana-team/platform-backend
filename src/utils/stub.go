package utils

import (
	cappv1alpha1 "github.com/dana-team/container-app-operator/api/v1alpha1"
	"github.com/dana-team/platform-backend/src/types"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetStubCappRevision returns a CappRevision object with only ObjectMeta set.
func GetStubCappRevision(name, namespace string, labels, annotations map[string]string) cappv1alpha1.CappRevision {
	cappRevision := cappv1alpha1.CappRevision{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   namespace,
			Annotations: annotations,
			Labels:      labels,
		},
		Spec:   cappv1alpha1.CappRevisionSpec{},
		Status: cappv1alpha1.CappRevisionStatus{},
	}

	return cappRevision
}

// GetStubCappRevisionType returns a CappRevision type object.
func GetStubCappRevisionType(name, namespace string, labels, annotations []types.KeyValue) types.CappRevision {
	cappRevision := types.CappRevision{
		Annotations: annotations,
		Labels:      labels,
		Metadata: types.Metadata{
			Name:      name,
			Namespace: namespace,
		},
		Spec:   cappv1alpha1.CappRevisionSpec{},
		Status: cappv1alpha1.CappRevisionStatus{},
	}

	return cappRevision
}

// GetStubConfigMap returns a ConfigMap object with only ObjectMeta set.
func GetStubConfigMap(name, namespace string, data map[string]string) corev1.ConfigMap {
	configMap := corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: data,
	}

	return configMap
}

// GetStubConfigMapType returns a ConfigMap type object.
func GetStubConfigMapType(data map[string]string) types.ConfigMap {
	configMap := types.ConfigMap{
		Data: ConvertMapToKeyValue(data),
	}

	return configMap
}
