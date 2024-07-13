package utils

import cappv1 "github.com/dana-team/container-app-operator/api/v1alpha1"

var domain = cappv1.GroupVersion.Group
var managedLabel = domain + "/managed"
var ManagedLabelSelctor = managedLabel + "=true"

// AddManagedLabel adds the managed label to the given labels map.
func AddManagedLabel(labels map[string]string) map[string]string {
	labels[managedLabel] = "true"
	return labels
}
