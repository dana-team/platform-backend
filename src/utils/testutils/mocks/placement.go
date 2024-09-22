package mocks

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1beta1 "open-cluster-management.io/api/cluster/v1beta1"
)

const (
	clusterSetName = "mock-clusterset"
)

// PreparePlacement returns a mock Placement object.
func PreparePlacement(name, namespace string, labels map[string]string) clusterv1beta1.Placement {
	placement := clusterv1beta1.Placement{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: clusterv1beta1.PlacementSpec{
			ClusterSets: []string{clusterSetName},
		},
	}

	return placement
}
