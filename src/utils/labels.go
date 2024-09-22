package utils

import (
	"fmt"
	cappv1alpha1 "github.com/dana-team/container-app-operator/api/v1alpha1"
)

var (
	cappAPIGroup         = cappv1alpha1.GroupVersion.Group
	ManagedLabel         = cappAPIGroup + "/managed"
	ManagedLabelValue    = "true"
	ManagedLabelSelector = fmt.Sprintf("%s=%s", ManagedLabel, ManagedLabelValue)

	ParentCappLabel         = cappAPIGroup + "/parent-capp"
	ParentCappNSLabel       = cappAPIGroup + "/parent-capp-ns"
	ParentCappLabelSelector = ParentCappLabel + "=%s"

	PlacementEnvironmentKey   = "environment"
	PlacementRegionKey        = "region"
	PlacementEnvironmentLabel = cappAPIGroup + "/" + PlacementEnvironmentKey
	PlacementRegionLabel      = cappAPIGroup + "/" + PlacementRegionKey

	CappNameLabel         = cappAPIGroup + "/cappName"
	CappNameLabelSelector = CappNameLabel + "=%s"
)

// AddManagedLabel adds the managed label to the given labels map.
func AddManagedLabel(labels map[string]string) map[string]string {
	labels[ManagedLabel] = "true"
	return labels
}
