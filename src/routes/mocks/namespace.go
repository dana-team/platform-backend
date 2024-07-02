package mocks

import (
	"github.com/dana-team/platform-backend/src/types"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PrepareNamespace returns a mock Namespace object.
func PrepareNamespace(name string) corev1.Namespace {
	return corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
}

// PrepareNamespaceType returns a mock Namespace type object.
func PrepareNamespaceType(name string) types.Namespace {
	return types.Namespace{
		Name: name,
	}
}
