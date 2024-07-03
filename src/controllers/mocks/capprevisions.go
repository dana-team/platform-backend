package mocks

import (
	cappv1alpha1 "github.com/dana-team/container-app-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PrepareCappRevision returns a mock CappRevision object.
func PrepareCappRevision(name, namespace string, labels, annotations map[string]string) cappv1alpha1.CappRevision {
	return cappv1alpha1.CappRevision{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   namespace,
			Annotations: annotations,
			Labels:      labels,
		},
		Spec:   PrepareCappRevisionSpec(),
		Status: PrepareCappRevisionStatus(),
	}
}

// PrepareCappRevisionSpec returns a mock CappRevision Spec object.
func PrepareCappRevisionSpec() cappv1alpha1.CappRevisionSpec {
	return cappv1alpha1.CappRevisionSpec{
		RevisionNumber: 1,
		CappTemplate:   cappv1alpha1.CappTemplate{},
	}
}

// PrepareCappRevisionStatus returns a mock CappRevision Status object.
func PrepareCappRevisionStatus() cappv1alpha1.CappRevisionStatus {
	return cappv1alpha1.CappRevisionStatus{}
}
