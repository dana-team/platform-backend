package controllers

import (
	"context"
	"fmt"

	"github.com/dana-team/container-app-operator/api/v1alpha1"
	"github.com/dana-team/platform-backend/src/types"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ContainerAppRevisionController interface {
	// GetContainerAppRevisions gets all container app revisions from a specific namespace.
	GetContainerAppRevisions(namespace string, cappRevisionQuery types.CappRevisionQuery) (types.CappRevisionList, error)

	// GetContainerAppRevision gets a specific container app revision from the specified namespace.
	GetContainerAppRevision(namespace, name string) (types.CappRevision, error)
}

type containerAppRevisionController struct {
	client client.Client
	ctx    context.Context
	logger *zap.Logger
}

func NewContainerAppRevisionController(client client.Client, context context.Context, logger *zap.Logger) (ContainerAppRevisionController, error) {
	return &containerAppRevisionController{
		client: client,
		ctx:    context,
		logger: logger,
	}, nil
}

func (c *containerAppRevisionController) GetContainerAppRevision(namespace string, name string) (types.CappRevision, error) {
	c.logger.Debug(fmt.Sprintf("Trying to fetch capp revision %q in namespace %q", name, namespace))

	cappRevision := v1alpha1.CappRevision{}
	if err := c.client.Get(c.ctx, client.ObjectKey{Namespace: namespace, Name: name}, &cappRevision); err != nil {
		c.logger.Error(fmt.Sprintf("Could not fetch capp revision %q in namespace %q with error: %s", name, namespace, err.Error()))
		return types.CappRevision{}, err
	}

	return convertV1alpha1CappRevisionToTypesCappRevision(cappRevision), nil
}

func (c *containerAppRevisionController) GetContainerAppRevisions(namespace string, cappQuery types.CappRevisionQuery) (types.CappRevisionList, error) {
	c.logger.Debug(fmt.Sprintf("Trying to fetch all capp revisions in namespace: %q", namespace))

	cappRevisionList := &v1alpha1.CappRevisionList{}
	selector, err := metav1.LabelSelectorAsSelector(&metav1.LabelSelector{
		MatchLabels: convertLabelsToMap(cappQuery.Labels),
	})
	if err != nil {
		c.logger.Error(fmt.Sprintf("Could not create label selector with error: %v", err.Error()))
		return types.CappRevisionList{}, err
	}

	err = c.client.List(c.ctx, cappRevisionList, &client.ListOptions{
		Namespace:     namespace,
		LabelSelector: selector,
	})
	if err != nil {
		c.logger.Error(fmt.Sprintf("Could not fetch capp revisions in namespace %q with error: %v", namespace, err.Error()))
		return types.CappRevisionList{}, err
	}

	result := types.CappRevisionList{}
	result.CappRevisions = append(result.CappRevisions, cappRevisionList.Items...)
	result.Count = len(cappRevisionList.Items)
	return result, nil
}

func convertV1alpha1CappRevisionToTypesCappRevision(v1CappRevision v1alpha1.CappRevision) types.CappRevision {
	return types.CappRevision{
		Metadata: types.Metadata{
			Name:      v1CappRevision.Name,
			Namespace: v1CappRevision.Namespace,
		},
		Annotations: convertMapToAnnotations(v1CappRevision.Annotations),
		Labels:      convertMapToLabels(v1CappRevision.Labels),
		Spec: v1alpha1.CappRevisionSpec{
			RevisionNumber: v1CappRevision.Spec.RevisionNumber,
			CappTemplate:   v1CappRevision.Spec.CappTemplate,
		},
		Status: v1alpha1.CappRevisionStatus{},
	}
}
