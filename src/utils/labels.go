package utils

import (
	"fmt"
	cappv1alpha1 "github.com/dana-team/container-app-operator/api/v1alpha1"
)

var cappAPIGroup = cappv1alpha1.GroupVersion.Group
var ManagedLabel = cappAPIGroup + "/managed"
var ManagedLabelValue = "true"
var ManagedLabelSelector = fmt.Sprintf("%s=%s", ManagedLabel, ManagedLabelValue)

// AddManagedLabel adds the managed label to the given labels map.
func AddManagedLabel(labels map[string]string) map[string]string {
	labels[ManagedLabel] = "true"
	return labels
}
