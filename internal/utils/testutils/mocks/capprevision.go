package mocks

import (
	cappv1alpha1 "github.com/dana-team/container-app-operator/api/v1alpha1"
	"github.com/dana-team/platform-backend/internal/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PrepareCappRevision returns a mock CappRevision object.
func PrepareCappRevision(name, namespace, site string, labels, annotations map[string]string) cappv1alpha1.CappRevision {
	cappRevision := cappv1alpha1.CappRevision{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   namespace,
			Annotations: annotations,
			Labels:      labels,
		},
		Spec:   PrepareCappRevisionSpec(site, labels, annotations),
		Status: PrepareCappRevisionStatus(),
	}

	return cappRevision
}

// PrepareCappRevisionSpec returns a mock CappRevision Spec object.
func PrepareCappRevisionSpec(site string, labels, annotations map[string]string) cappv1alpha1.CappRevisionSpec {
	cappRevisionSpec := cappv1alpha1.CappRevisionSpec{
		RevisionNumber: 1,
		CappTemplate: cappv1alpha1.CappTemplate{
			Spec: PrepareCappSpec(site),
		},
	}

	if annotations != nil {
		cappRevisionSpec.CappTemplate.Annotations = annotations
	}

	if labels != nil {
		cappRevisionSpec.CappTemplate.Labels = labels
	}

	return cappRevisionSpec
}

// PrepareCappRevisionStatus returns a mock CappRevision Status object.
func PrepareCappRevisionStatus() cappv1alpha1.CappRevisionStatus {
	return cappv1alpha1.CappRevisionStatus{}
}

func ConvertKeyValueSliceToMap(keyValueSlice []types.KeyValue) map[string]string {
	keyValueMap := make(map[string]string)
	for _, kv := range keyValueSlice {
		keyValueMap[kv.Key] = kv.Value
	}

	return keyValueMap
}
