package controllers

import (
	"context"
	"fmt"
	"github.com/dana-team/platform-backend/src/utils"
	"k8s.io/apimachinery/pkg/labels"

	"github.com/dana-team/container-app-operator/api/v1alpha1"
	"github.com/dana-team/platform-backend/src/types"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type CappController interface {
	// CreateCapp creates a new container app in the specified namespace.
	CreateCapp(namespace string, capp types.CreateCapp) (types.Capp, error)

	// GetCapps gets all container apps from a specific namespace.
	GetCapps(namespace string, cappQuery types.CappQuery) (types.CappList, error)

	// GetCapp gets a specific container app from the specified namespace.
	GetCapp(namespace, name string) (types.Capp, error)

	// UpdateCapp updates a specific container app in the specified namespace.
	UpdateCapp(namespace, name string, capp types.UpdateCapp) (types.Capp, error)

	// DeleteCapp deletes a specific container app in the specified namespace.
	DeleteCapp(namespace, name string) (types.CappError, error)
}

type cappController struct {
	client client.Client
	ctx    context.Context
	logger *zap.Logger
}

func NewCappController(client client.Client, context context.Context, logger *zap.Logger) (CappController, error) {
	return &cappController{
		client: client,
		ctx:    context,
		logger: logger,
	}, nil
}

func (c *cappController) CreateCapp(namespace string, capp types.CreateCapp) (types.Capp, error) {
	c.logger.Debug(fmt.Sprintf("Trying to create capp in namespace: %q", namespace))

	newCapp := createCappFromType(namespace, capp)
	if err := c.client.Create(c.ctx, &newCapp); err != nil {
		c.logger.Error(fmt.Sprintf("Could not create capp in namespace %q with error: %v", namespace, err.Error()))
		return types.Capp{}, err
	}

	return createCappFromV1Capp(newCapp), nil
}

func createCappFromV1Capp(capp v1alpha1.Capp) types.Capp {
	return types.Capp{
		Metadata: types.Metadata{
			Name:      capp.Name,
			Namespace: capp.Namespace,
		},
		Annotations: utils.ConvertMapToKeyValue(capp.Annotations),
		Labels:      utils.ConvertMapToKeyValue(capp.Labels),
		Spec: v1alpha1.CappSpec{
			ScaleMetric:       capp.Spec.ScaleMetric,
			Site:              capp.Spec.Site,
			State:             capp.Spec.State,
			ConfigurationSpec: capp.Spec.ConfigurationSpec,
			RouteSpec:         capp.Spec.RouteSpec,
			LogSpec:           capp.Spec.LogSpec,
			VolumesSpec:       capp.Spec.VolumesSpec,
		},
		Status: v1alpha1.CappStatus{
			ApplicationLinks:    capp.Status.ApplicationLinks,
			KnativeObjectStatus: capp.Status.KnativeObjectStatus,
			RevisionInfo:        capp.Status.RevisionInfo,
			StateStatus:         capp.Status.StateStatus,
			LoggingStatus:       capp.Status.LoggingStatus,
			RouteStatus:         capp.Status.RouteStatus,
			VolumesStatus:       capp.Status.VolumesStatus,
			Conditions:          capp.Status.Conditions,
		},
	}
}

func (c *cappController) GetCapps(namespace string, cappQuery types.CappQuery) (types.CappList, error) {
	c.logger.Debug(fmt.Sprintf("Trying to fetch all capps in namespace: %q", namespace))

	cappList := &v1alpha1.CappList{}
	selector, err := labels.Parse(cappQuery.LabelSelector)
	if err != nil {
		c.logger.Error(fmt.Sprintf("Could not parse labelSelector with error: %v", err.Error()))
		return types.CappList{}, err
	}

	err = c.client.List(c.ctx, cappList, &client.ListOptions{
		Namespace:     namespace,
		LabelSelector: selector,
	})
	if err != nil {
		c.logger.Error(fmt.Sprintf("Could not fetch capps in namespace %q with error: %v", namespace, err.Error()))
		return types.CappList{}, err
	}

	result := types.CappList{}
	for _, item := range cappList.Items {
		result.Capps = append(result.Capps, convertCappToType(item))
	}
	result.Count = len(cappList.Items)
	return result, nil
}

func (c *cappController) GetCapp(namespace, name string) (types.Capp, error) {
	c.logger.Debug(fmt.Sprintf("Trying to fetch capp %q in namespace %q", name, namespace))

	capp := &v1alpha1.Capp{}
	err := c.client.Get(c.ctx, client.ObjectKey{Namespace: namespace, Name: name}, capp)
	if err != nil {
		c.logger.Error(fmt.Sprintf("Could not fetch capp %q in namespace %q with error: %v", name, namespace, err.Error()))
		return types.Capp{}, err
	}

	return convertCappToType(*capp), nil
}

func (c *cappController) UpdateCapp(namespace, name string, newCapp types.UpdateCapp) (types.Capp, error) {
	c.logger.Debug(fmt.Sprintf("Trying to update capp %q in namespace %q", name, namespace))

	capp := &v1alpha1.Capp{}
	err := c.client.Get(c.ctx, client.ObjectKey{Namespace: namespace, Name: name}, capp)
	if err != nil {
		c.logger.Error(fmt.Sprintf("Could not fetch capp %q in namespace %q with error: %v", name, namespace, err.Error()))
		return types.Capp{}, err
	}

	capp.Annotations = utils.ConvertKeyValueToMap(newCapp.Annotations)
	capp.Labels = utils.ConvertKeyValueToMap(newCapp.Labels)
	capp.Spec = newCapp.Spec

	if err := c.client.Update(c.ctx, capp); err != nil {
		c.logger.Error(fmt.Sprintf("Could not update capp %q in namespace %q with error: %v", name, namespace, err.Error()))
		return types.Capp{}, err
	}

	return convertCappToType(*capp), nil
}

func (c *cappController) DeleteCapp(namespace, name string) (types.CappError, error) {
	c.logger.Debug(fmt.Sprintf("Trying to delete capp %q in namespace %q", name, namespace))

	capp := &v1alpha1.Capp{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
		},
	}
	if err := c.client.Delete(c.ctx, capp); err != nil {
		c.logger.Error(fmt.Sprintf("Could not delete capp %q in namespace %q with error: %v", name, namespace, err.Error()))
		return types.CappError{}, nil
	}

	return types.CappError{
		Message: fmt.Sprintf("Deleted capp %q in namespace %q successfully", name, namespace),
	}, nil
}

func createCappFromType(namespace string, capp types.CreateCapp) v1alpha1.Capp {
	return v1alpha1.Capp{
		ObjectMeta: metav1.ObjectMeta{
			Name:        capp.Metadata.Name,
			Namespace:   namespace,
			Annotations: utils.ConvertKeyValueToMap(capp.Annotations),
			Labels:      utils.ConvertKeyValueToMap(capp.Labels),
		},
		Spec: v1alpha1.CappSpec{
			ScaleMetric:       capp.Spec.ScaleMetric,
			Site:              capp.Spec.Site,
			State:             capp.Spec.State,
			ConfigurationSpec: capp.Spec.ConfigurationSpec,
			RouteSpec:         capp.Spec.RouteSpec,
			LogSpec:           capp.Spec.LogSpec,
			VolumesSpec:       capp.Spec.VolumesSpec,
		},
	}
}

func convertCappToType(capp v1alpha1.Capp) types.Capp {
	return types.Capp{
		Metadata: types.Metadata{
			Name:      capp.Name,
			Namespace: capp.Namespace,
		},
		Annotations: utils.ConvertMapToKeyValue(capp.Annotations),
		Labels:      utils.ConvertMapToKeyValue(capp.Labels),
		Spec: v1alpha1.CappSpec{
			ScaleMetric:       capp.Spec.ScaleMetric,
			Site:              capp.Spec.Site,
			State:             capp.Spec.State,
			ConfigurationSpec: capp.Spec.ConfigurationSpec,
			RouteSpec:         capp.Spec.RouteSpec,
			LogSpec:           capp.Spec.LogSpec,
			VolumesSpec:       capp.Spec.VolumesSpec,
		},
		Status: v1alpha1.CappStatus{
			ApplicationLinks:    capp.Status.ApplicationLinks,
			KnativeObjectStatus: capp.Status.KnativeObjectStatus,
			RevisionInfo:        capp.Status.RevisionInfo,
			StateStatus:         capp.Status.StateStatus,
			LoggingStatus:       capp.Status.LoggingStatus,
			RouteStatus:         capp.Status.RouteStatus,
			VolumesStatus:       capp.Status.VolumesStatus,
			Conditions:          capp.Status.Conditions,
		},
	}
}
