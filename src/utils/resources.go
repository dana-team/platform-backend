package utils

import (
	cappv1alpha1 "github.com/dana-team/container-app-operator/api/v1alpha1"
	"github.com/dana-team/platform-backend/src/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetBareCappRevision returns a CappRevision object with only ObjectMeta set.
func GetBareCappRevision(name, namespace string, labels, annotations map[string]string) cappv1alpha1.CappRevision {
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

// GetBareCappRevisionType returns a CappRevision type object.
func GetBareCappRevisionType(name, namespace string, labels, annotations []types.KeyValue) types.CappRevision {
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
