package controllers

import (
	"context"
	"fmt"
	"github.com/dana-team/platform-backend/src/utils"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"

	"github.com/dana-team/container-app-operator/api/v1alpha1"
	"github.com/dana-team/platform-backend/src/types"
	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type CappRevisionController interface {
	// GetCappRevisions gets all CappRevision names and count from a specific namespace.
	GetCappRevisions(namespace string, cappRevisionQuery types.CappRevisionQuery) (types.CappRevisionList, error)

	// GetCappRevision gets a specific CappRevision from the specified namespace.
	GetCappRevision(namespace, name string) (types.CappRevision, error)
}

type cappRevisionController struct {
	client client.Client
	ctx    context.Context
	logger *zap.Logger
}

func NewCappRevisionController(client client.Client, context context.Context, logger *zap.Logger) CappRevisionController {
	return &cappRevisionController{
		client: client,
		ctx:    context,
		logger: logger,
	}
}

func (c *cappRevisionController) GetCappRevision(namespace string, name string) (types.CappRevision, error) {
	c.logger.Debug(fmt.Sprintf("Trying to fetch capp revision %q in namespace %q", name, namespace))

	cappRevision := v1alpha1.CappRevision{}
	if err := c.client.Get(c.ctx, client.ObjectKey{Namespace: namespace, Name: name}, &cappRevision); err != nil {
		c.logger.Error(fmt.Sprintf("Could not fetch capp revision %q in namespace %q with error: %s", name, namespace, err.Error()))
		return types.CappRevision{}, err
	}

	return convertCappRevisionToType(cappRevision), nil
}

func (c *cappRevisionController) GetCappRevisions(namespace string, cappQuery types.CappRevisionQuery) (types.CappRevisionList, error) {
	c.logger.Debug(fmt.Sprintf("Trying to fetch all capp revisions in namespace: %q", namespace))

	cappRevisionList := &v1alpha1.CappRevisionList{}
	selector, err := labels.Parse(cappQuery.LabelSelector)
	if err != nil {
		c.logger.Error(fmt.Sprintf("Could not parse labelSelector with error: %v", err.Error()))
		return types.CappRevisionList{}, k8serrors.NewBadRequest(err.Error())
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
	for _, revision := range cappRevisionList.Items {
		result.CappRevisions = append(result.CappRevisions, revision.Name)
	}
	result.Count = len(cappRevisionList.Items)

	return result, nil
}

// convertCappRevisionToType converts an API CappRevision to a Type CappRevision.
func convertCappRevisionToType(cappRevision v1alpha1.CappRevision) types.CappRevision {
	return types.CappRevision{
		Metadata: types.Metadata{
			Name:      cappRevision.Name,
			Namespace: cappRevision.Namespace,
		},
		Annotations: utils.ConvertMapToKeyValue(cappRevision.Annotations),
		Labels:      utils.ConvertMapToKeyValue(cappRevision.Labels),
		Spec: v1alpha1.CappRevisionSpec{
			RevisionNumber: cappRevision.Spec.RevisionNumber,
			CappTemplate:   cappRevision.Spec.CappTemplate,
		},
		Status: v1alpha1.CappRevisionStatus{},
	}
}
