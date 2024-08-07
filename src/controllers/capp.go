package controllers

import (
	"context"
	"fmt"
	"github.com/dana-team/platform-backend/src/utils"
	"github.com/dana-team/platform-backend/src/utils/pagination"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"

	cappv1alpha1 "github.com/dana-team/container-app-operator/api/v1alpha1"
	"github.com/dana-team/platform-backend/src/types"

	dnsrecordv1alpha1 "github.com/dana-team/provider-dns/apis/record/v1alpha1"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	disabledState = "disabled"
	noRevision    = "No revision available"
	dnsLimit      = 10
)

type CappController interface {
	// CreateCapp creates a new Capp in the specified namespace.
	CreateCapp(namespace string, capp types.CreateCapp) (types.Capp, error)

	// GetCapps gets all Capps from a specific namespace.
	GetCapps(namespace string, limit, page int, cappQuery types.CappQuery) (types.CappList, error)

	// GetCapp gets a specific Capp from the specified namespace.
	GetCapp(namespace, name string) (types.Capp, error)

	// UpdateCapp updates a specific Capp in the specified namespace.
	UpdateCapp(namespace, name string, capp types.UpdateCapp) (types.Capp, error)

	// DeleteCapp deletes a specific Capp in the specified namespace.
	DeleteCapp(namespace, name string) (types.CappError, error)

	// EditCappState edits the state of a specific Capp in the specified namespace.
	EditCappState(namespace string, cappName string, state string) (types.CappStateReponse, error)

	// GetCappState gets the state of a specific Capp from the specified namespace.
	GetCappState(namespace, name string) (types.GetCappStateResponse, error)

	// GetCappDNS gets the dns records which are related to the Capp
	GetCappDNS(namespace, name string) (types.GetDNSResponse, error)
}

type cappController struct {
	client client.Client
	ctx    context.Context
	logger *zap.Logger
}

func NewCappController(client client.Client, context context.Context, logger *zap.Logger) CappController {
	return &cappController{
		client: client,
		ctx:    context,
		logger: logger,
	}
}

// CappPaginator paginates through capps in a specified namespace.
type CappPaginator struct {
	pagination.GenericPaginator
	namespace string
	client    client.Client
	cappQuery types.CappQuery
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

func createCappFromV1Capp(capp cappv1alpha1.Capp) types.Capp {
	return types.Capp{
		Metadata: types.Metadata{
			Name:      capp.Name,
			Namespace: capp.Namespace,
		},
		Annotations: utils.ConvertMapToKeyValue(capp.Annotations),
		Labels:      utils.ConvertMapToKeyValue(capp.Labels),
		Spec: cappv1alpha1.CappSpec{
			ScaleMetric:       capp.Spec.ScaleMetric,
			Site:              capp.Spec.Site,
			State:             capp.Spec.State,
			ConfigurationSpec: capp.Spec.ConfigurationSpec,
			RouteSpec:         capp.Spec.RouteSpec,
			LogSpec:           capp.Spec.LogSpec,
			VolumesSpec:       capp.Spec.VolumesSpec,
		},
		Status: cappv1alpha1.CappStatus{
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

func (c *cappController) GetCapps(namespace string, limit, page int, cappQuery types.CappQuery) (types.CappList, error) {
	c.logger.Debug(fmt.Sprintf("Trying to fetch all capps in namespace: %q", namespace))

	cappPaginator := &CappPaginator{
		GenericPaginator: pagination.CreatePaginator(c.ctx, c.logger),
		namespace:        namespace,
		client:           c.client,
		cappQuery:        cappQuery,
	}

	cappList, err := pagination.FetchPage[cappv1alpha1.Capp](limit, page, cappPaginator)
	if err != nil {
		c.logger.Error(fmt.Sprintf("Could not get capps with error: %v", err))
		return types.CappList{}, err
	}

	result := types.CappList{}
	for _, item := range cappList {
		summary := types.CappSummary{
			Name:   item.Name,
			URL:    getCappURL(item),
			Images: getCappImages(item),
		}
		result.Capps = append(result.Capps, summary)
	}
	result.Count = len(cappList)

	return result, nil
}

func (c *cappController) GetCapp(namespace, name string) (types.Capp, error) {
	c.logger.Debug(fmt.Sprintf("Trying to fetch capp %q in namespace %q", name, namespace))

	capp := &cappv1alpha1.Capp{}
	err := c.client.Get(c.ctx, client.ObjectKey{Namespace: namespace, Name: name}, capp)
	if err != nil {
		c.logger.Error(fmt.Sprintf("Could not fetch capp %q in namespace %q with error: %v", name, namespace, err.Error()))
		return types.Capp{}, err
	}

	return convertCappToType(*capp), nil
}

func (c *cappController) GetCappState(namespace, name string) (types.GetCappStateResponse, error) {
	c.logger.Debug(fmt.Sprintf("Trying to fetch capp %q in namespace %q", name, namespace))

	capp := &cappv1alpha1.Capp{}
	err := c.client.Get(c.ctx, client.ObjectKey{Namespace: namespace, Name: name}, capp)
	if err != nil {
		c.logger.Error(fmt.Sprintf("Could not fetch capp %q in namespace %q with error: %v", name, namespace, err.Error()))
		return types.GetCappStateResponse{}, err
	}

	var cappState types.GetCappStateResponse

	if capp.Status.StateStatus.State == disabledState {
		cappState = types.GetCappStateResponse{
			LastCreatedRevision: noRevision,
			LastReadyRevision:   noRevision,
			State:               capp.Status.StateStatus.State}
	} else {
		cappState = types.GetCappStateResponse{
			LastCreatedRevision: capp.Status.KnativeObjectStatus.LatestCreatedRevisionName,
			LastReadyRevision:   capp.Status.KnativeObjectStatus.LatestReadyRevisionName,
			State:               capp.Status.StateStatus.State}

	}

	return cappState, nil
}

func (c *cappController) GetCappDNS(namespace, name string) (types.GetDNSResponse, error) {
	c.logger.Debug(fmt.Sprintf("Trying to fetch dns related to capp %q in namespace %q", name, namespace))

	dnsRecords := &dnsrecordv1alpha1.CNAMERecordList{}

	listOptions := prepareDNSListOptions(namespace, name)

	err := c.client.List(c.ctx, dnsRecords, listOptions)
	if err != nil {
		c.logger.Error(fmt.Sprintf("Could not fetch dns related to capp %q in namespace %q with error: %v", name, namespace, err.Error()))
		return types.GetDNSResponse{}, err
	}

	if len(dnsRecords.Items) == 0 {
		_, err := c.GetCapp(namespace, name)
		if err != nil {
			return types.GetDNSResponse{}, err
		}
		return types.GetDNSResponse{}, nil
	}

	return types.GetDNSResponse{Records: parseDNSList(dnsRecords)}, nil
}

// parseDNSList converts DNSRecords from their Kubernetes representation to a custom type representation.
func parseDNSList(recordsK8S *dnsrecordv1alpha1.CNAMERecordList) []types.DNS {
	var records []types.DNS

	for _, record := range recordsK8S.Items {
		syncedStatus := record.GetCondition(xpv1.TypeSynced).Status
		readyStatus := record.GetCondition(xpv1.TypeReady).Status
		dnsStatus := corev1.ConditionFalse

		if syncedStatus == corev1.ConditionTrue && readyStatus == corev1.ConditionTrue {
			dnsStatus = corev1.ConditionTrue
		} else if syncedStatus == corev1.ConditionUnknown || readyStatus == corev1.ConditionUnknown {
			dnsStatus = corev1.ConditionUnknown
		}

		records = append(records, types.DNS{Status: dnsStatus, Name: computeDNSNameFromID(*record.Status.AtProvider.ID)})
	}
	return records
}

// computes dns record ID ready for presentation, it removes the  '.' at the end of the string
func computeDNSNameFromID(dnsID string) string {
	return dnsID[:len(dnsID)-1]
}

// prepareDNSListOptions prepares a list options for querying.
func prepareDNSListOptions(namespace, name string) *client.ListOptions {
	labelSet := map[string]string{utils.ParentCappNSLabel: namespace, utils.ParentCappLabel: name}
	labelSelector := labels.SelectorFromSet(labelSet)

	listOptions := &client.ListOptions{
		LabelSelector: labelSelector,
		Limit:         dnsLimit,
	}

	return listOptions
}

func (c *cappController) UpdateCapp(namespace, name string, newCapp types.UpdateCapp) (types.Capp, error) {
	c.logger.Debug(fmt.Sprintf("Trying to update capp %q in namespace %q", name, namespace))

	capp := &cappv1alpha1.Capp{}
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

func (c *cappController) EditCappState(namespace string, cappName string, state string) (types.CappStateReponse, error) {
	c.logger.Debug(fmt.Sprintf("Trying to update capp %q in namespace %q", cappName, namespace))

	capp := &cappv1alpha1.Capp{}
	err := c.client.Get(c.ctx, client.ObjectKey{Namespace: namespace, Name: cappName}, capp)
	if err != nil {
		c.logger.Error(fmt.Sprintf("Could not fetch capp %q in namespace %q with error: %v", cappName, namespace, err.Error()))
		return types.CappStateReponse{}, err
	}

	capp.Spec.State = state
	if err := c.client.Update(c.ctx, capp); err != nil {
		c.logger.Error(fmt.Sprintf("Could not update capp %q in namespace %q with error: %v", state, namespace, err.Error()))
		return types.CappStateReponse{}, err
	}

	return types.CappStateReponse{Name: capp.Name, State: capp.Spec.State}, nil
}

func (c *cappController) DeleteCapp(namespace, name string) (types.CappError, error) {
	c.logger.Debug(fmt.Sprintf("Trying to delete capp %q in namespace %q", name, namespace))

	capp := &cappv1alpha1.Capp{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
		},
	}
	if err := c.client.Delete(c.ctx, capp); err != nil {
		c.logger.Error(fmt.Sprintf("Could not delete capp %q in namespace %q with error: %v", name, namespace, err.Error()))
		return types.CappError{}, err
	}

	return types.CappError{
		Message: fmt.Sprintf("Deleted capp %q in namespace %q successfully", name, namespace),
	}, nil
}

// FetchList retrieves a list of capps from the specified namespace with given options.
func (p *CappPaginator) FetchList(listOptions metav1.ListOptions) (*types.List[cappv1alpha1.Capp], error) {
	cappList := &cappv1alpha1.CappList{}
	selector, err := labels.Parse(p.cappQuery.LabelSelector)
	if err != nil {
		p.Logger.Error(fmt.Sprintf("Could not parse labelSelector with error: %v", err.Error()))
		return nil, k8serrors.NewBadRequest(err.Error())
	}

	err = p.client.List(p.Ctx, cappList, &client.ListOptions{
		Limit:         listOptions.Limit,
		Continue:      listOptions.Continue,
		Namespace:     p.namespace,
		LabelSelector: selector,
	})
	if err != nil {
		p.Logger.Error(fmt.Sprintf("Could not fetch capps in namespace %q with error: %v", p.namespace, err.Error()))
		return nil, err
	}

	return (*types.List[cappv1alpha1.Capp])(cappList), nil
}

func createCappFromType(namespace string, capp types.CreateCapp) cappv1alpha1.Capp {
	return cappv1alpha1.Capp{
		ObjectMeta: metav1.ObjectMeta{
			Name:        capp.Metadata.Name,
			Namespace:   namespace,
			Annotations: utils.ConvertKeyValueToMap(capp.Annotations),
			Labels:      utils.ConvertKeyValueToMap(capp.Labels),
		},
		Spec: cappv1alpha1.CappSpec{
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

func convertCappToType(capp cappv1alpha1.Capp) types.Capp {
	return types.Capp{
		Metadata: types.Metadata{
			Name:      capp.Name,
			Namespace: capp.Namespace,
		},
		Annotations: utils.ConvertMapToKeyValue(capp.Annotations),
		Labels:      utils.ConvertMapToKeyValue(capp.Labels),
		Spec: cappv1alpha1.CappSpec{
			ScaleMetric:       capp.Spec.ScaleMetric,
			Site:              capp.Spec.Site,
			State:             capp.Spec.State,
			ConfigurationSpec: capp.Spec.ConfigurationSpec,
			RouteSpec:         capp.Spec.RouteSpec,
			LogSpec:           capp.Spec.LogSpec,
			VolumesSpec:       capp.Spec.VolumesSpec,
		},
		Status: cappv1alpha1.CappStatus{
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

// getCappURL returns the URL of Capp; the shortened hostname is returned
// if it exists, otherwise the default URL is returned.
func getCappURL(capp cappv1alpha1.Capp) string {
	if capp.Status.RouteStatus.DomainMappingObjectStatus.URL != nil {
		return capp.Status.RouteStatus.DomainMappingObjectStatus.URL.URL().String()
	} else if capp.Status.KnativeObjectStatus.URL != nil {
		return capp.Status.KnativeObjectStatus.URL.URL().String()
	}

	return ""
}

// getCappImages returns the images of all containers of Capp.
func getCappImages(capp cappv1alpha1.Capp) []string {
	var images []string
	for _, container := range capp.Spec.ConfigurationSpec.Template.Spec.Containers {
		images = append(images, container.Image)
	}

	return images
}
