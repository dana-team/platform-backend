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

type ContainerAppController interface {
	// CreateContainerApp creates a new container app in the specified namespace.
	CreateContainerApp(namespace string, capp types.Capp) (types.Capp, error)

	// GetContainerApps gets all container apps from a specific namespace.
	GetContainerApps(namespace string, cappQuery types.CappQuery) (types.CappList, error)

	// GetContainerApp gets a specific container app from the specified namespace.
	GetContainerApp(namespace, name string) (types.Capp, error)

	// PatchContainerApp patches a specific container app in the specified namespace.
	PatchContainerApp(namespace, name string, capp types.Capp) (types.Capp, error)

	// DeleteContainerApp deletes a specific container app in the specified namespace.
	DeleteContainerApp(namespace, name string) (types.CappError, error)
}

type containerAppController struct {
	client client.Client
	ctx    context.Context
	logger *zap.Logger
}

func NewContainerAppController(client client.Client, context context.Context, logger *zap.Logger) (ContainerAppController, error) {
	return &containerAppController{
		client: client,
		ctx:    context,
		logger: logger,
	}, nil
}

func (c *containerAppController) CreateContainerApp(namespace string, capp types.Capp) (types.Capp, error) {
	c.logger.Debug(fmt.Sprintf("Trying to create capp in namespace: %q", namespace))

	newCapp := createV1alpha1Capp(capp)
	err := c.client.Create(c.ctx, &newCapp)
	if err != nil {
		c.logger.Error(fmt.Sprintf("Could not create capp in namespace %q with error: %s", namespace, err.Error()))
		return types.Capp{}, err
	}

	return capp, nil
}

func (c *containerAppController) GetContainerApps(namespace string, cappQuery types.CappQuery) (types.CappList, error) {
	c.logger.Debug(fmt.Sprintf("Trying to fetch all capps in namespace: %q", namespace))

	cappList := &v1alpha1.CappList{}
	selector, err := metav1.LabelSelectorAsSelector(&metav1.LabelSelector{
		MatchLabels: convertLabelsToMap(cappQuery.Labels),
	})
	if err != nil {
		c.logger.Error(fmt.Sprintf("Could not create label selector with error: %s", err.Error()))
		return types.CappList{}, err
	}

	err = c.client.List(c.ctx, cappList, &client.ListOptions{
		Namespace:     namespace,
		LabelSelector: selector,
	})
	if err != nil {
		c.logger.Error(fmt.Sprintf("Could not fetch capps in namespace %q with error: %s", namespace, err.Error()))
		return types.CappList{}, err
	}

	result := types.CappList{}
	for _, item := range cappList.Items {
		result.Capps = append(result.Capps, convertV1alpha1CappToTypesCapp(item))
	}
	result.Count = len(cappList.Items)
	return result, nil
}

func (c *containerAppController) GetContainerApp(namespace, name string) (types.Capp, error) {
	c.logger.Debug(fmt.Sprintf("Trying to fetch capp %q in namespace %q", name, namespace))

	capp := &v1alpha1.Capp{}
	err := c.client.Get(c.ctx, client.ObjectKey{Namespace: namespace, Name: name}, capp)
	if err != nil {
		c.logger.Error(fmt.Sprintf("Could not fetch capp %q in namespace %q with error: %s", name, namespace, err.Error()))
		return types.Capp{}, err
	}

	return convertV1alpha1CappToTypesCapp(*capp), nil
}

func (c *containerAppController) PatchContainerApp(namespace, name string, capp types.Capp) (types.Capp, error) {
	c.logger.Debug(fmt.Sprintf("Trying to patch capp %q in namespace %q", name, namespace))

	oldCapp := &v1alpha1.Capp{}
	err := c.client.Get(c.ctx, client.ObjectKey{Namespace: namespace, Name: name}, oldCapp)
	if err != nil {
		c.logger.Error(fmt.Sprintf("Could not fetch capp %q in namespace %q with error: %s", name, namespace, err.Error()))
		return types.Capp{}, err
	}

	newCapp := createV1alpha1Capp(capp)
	err = c.client.Patch(c.ctx, oldCapp, client.MergeFrom(newCapp.DeepCopy()))
	if err != nil {
		c.logger.Error(fmt.Sprintf("Could not patch capp %q in namespace %q with error: %s", name, namespace, err.Error()))
		return types.Capp{}, err
	}

	return convertV1alpha1CappToTypesCapp(newCapp), nil
}

func (c *containerAppController) DeleteContainerApp(namespace, name string) (types.CappError, error) {
	c.logger.Debug(fmt.Sprintf("Trying to delete capp %q in namespace %q", name, namespace))

	capp := &v1alpha1.Capp{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
		},
	}
	err := c.client.Delete(c.ctx, capp)
	if err != nil {
		c.logger.Error(fmt.Sprintf("Could not delete capp %q in namespace %q with error: %s", name, namespace, err.Error()))
		return types.CappError{}, nil
	}

	return types.CappError{
		Message: fmt.Sprintf("Deleted capp %q in namespace %q successfully", name, namespace),
	}, nil
}

func createV1alpha1Capp(capp types.Capp) v1alpha1.Capp {
	return v1alpha1.Capp{
		ObjectMeta: metav1.ObjectMeta{
			Name:        capp.Metadata.Name,
			Namespace:   capp.Metadata.Namespace,
			Annotations: convertAnnotationsToMap(capp.Annotations),
			Labels:      convertLabelsToMap(capp.Labels),
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

func convertV1alpha1CappToTypesCapp(v1Capp v1alpha1.Capp) types.Capp {
	return types.Capp{
		Metadata: types.Metadata{
			Name:      v1Capp.Name,
			Namespace: v1Capp.Namespace,
		},
		Annotations: convertMapToAnnotations(v1Capp.Annotations),
		Labels:      convertMapToLabels(v1Capp.Labels),
		Spec: v1alpha1.CappSpec{
			ScaleMetric:       v1Capp.Spec.ScaleMetric,
			Site:              v1Capp.Spec.Site,
			State:             v1Capp.Spec.State,
			ConfigurationSpec: v1Capp.Spec.ConfigurationSpec,
			RouteSpec:         v1Capp.Spec.RouteSpec,
			LogSpec:           v1Capp.Spec.LogSpec,
			VolumesSpec:       v1Capp.Spec.VolumesSpec,
		},
		Status: v1alpha1.CappStatus{
			ApplicationLinks:    v1Capp.Status.ApplicationLinks,
			KnativeObjectStatus: v1Capp.Status.KnativeObjectStatus,
			RevisionInfo:        v1Capp.Status.RevisionInfo,
			StateStatus:         v1Capp.Status.StateStatus,
			LoggingStatus:       v1Capp.Status.LoggingStatus,
			RouteStatus:         v1Capp.Status.RouteStatus,
			VolumesStatus:       v1Capp.Status.VolumesStatus,
			Conditions:          v1Capp.Status.Conditions,
		},
	}
}

func convertAnnotationsToMap(kvList []types.KeyValue) map[string]string {
	annotations := make(map[string]string)
	for _, kv := range kvList {
		annotations[kv.Key] = kv.Value
	}
	return annotations
}

func convertMapToAnnotations(annotations map[string]string) []types.KeyValue {
	var kvList []types.KeyValue
	for key, value := range annotations {
		kvList = append(kvList, types.KeyValue{Key: key, Value: value})
	}
	return kvList
}

func convertMapToLabels(labels map[string]string) []types.KeyValue {
	var kvList []types.KeyValue
	for key, value := range labels {
		kvList = append(kvList, types.KeyValue{Key: key, Value: value})
	}
	return kvList
}

func convertLabelsToMap(kvList []types.KeyValue) map[string]string {
	labels := make(map[string]string)
	for _, kv := range kvList {
		labels[kv.Key] = kv.Value
	}
	return labels
}
