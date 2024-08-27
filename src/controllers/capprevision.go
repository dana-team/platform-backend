package controllers

import (
	"context"
	"fmt"
	"github.com/dana-team/platform-backend/src/customerrors"
	"github.com/dana-team/platform-backend/src/utils"
	"github.com/dana-team/platform-backend/src/utils/pagination"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	cappv1alpha1 "github.com/dana-team/container-app-operator/api/v1alpha1"
	"github.com/dana-team/platform-backend/src/types"
	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	ErrCouldNotListCappRevisions = "Could not list capp revisions"
	ErrCouldNotGetCappRevision   = "Could not get capp revision %q in namespace %q"
)

type CappRevisionController interface {
	// GetCappRevisions gets all CappRevision names and count from a specific namespace.
	GetCappRevisions(namespace string, limit, page int, cappRevisionQuery types.CappRevisionQuery) (types.CappRevisionList, error)

	// GetCappRevision gets a specific CappRevision from the specified namespace.
	GetCappRevision(namespace, name string) (types.CappRevision, error)
}

type cappRevisionController struct {
	client client.Client
	ctx    context.Context
	logger *zap.Logger
}

// CappRevisionPaginator paginates through capprevisions in a specified namespace.
type CappRevisionPaginator struct {
	pagination.GenericPaginator
	namespace string
	client    client.Client
	cappQuery types.CappRevisionQuery
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

	cappRevision := cappv1alpha1.CappRevision{}
	if err := c.client.Get(c.ctx, client.ObjectKey{Namespace: namespace, Name: name}, &cappRevision); err != nil {
		c.logger.Error(fmt.Sprintf("%s with error: %s", fmt.Sprintf(ErrCouldNotGetCappRevision, name, namespace), err.Error()))
		return types.CappRevision{}, customerrors.NewAPIError(fmt.Sprintf(ErrCouldNotGetCappRevision, name, namespace), err)
	}

	return convertCappRevisionToType(cappRevision), nil
}

func (c *cappRevisionController) GetCappRevisions(namespace string, limit, page int, cappQuery types.CappRevisionQuery) (types.CappRevisionList, error) {
	c.logger.Debug(fmt.Sprintf("Trying to fetch all capp revisions in namespace: %q", namespace))

	cappRevisionPaginator := &CappRevisionPaginator{
		GenericPaginator: pagination.CreatePaginator(c.ctx, c.logger),
		namespace:        namespace,
		client:           c.client,
		cappQuery:        cappQuery,
	}

	cappRevisionList, err := pagination.FetchPage[cappv1alpha1.CappRevision](limit, page, cappRevisionPaginator)
	if err != nil {
		c.logger.Error(fmt.Sprintf("%s with error: %s", ErrCouldNotListCappRevisions, err.Error()))
		return types.CappRevisionList{}, customerrors.NewAPIError(ErrCouldNotListCappRevisions, err)
	}

	result := types.CappRevisionList{}
	for _, revision := range cappRevisionList {
		result.CappRevisions = append(result.CappRevisions, revision.Name)
	}
	result.Count = len(cappRevisionList)

	return result, nil
}

// FetchList retrieves a list of capps from the specified namespace with given options.
func (p *CappRevisionPaginator) FetchList(listOptions metav1.ListOptions) (*types.List[cappv1alpha1.CappRevision], error) {
	cappRevisionList := &cappv1alpha1.CappRevisionList{}
	selector, err := labels.Parse(p.cappQuery.LabelSelector)
	if err != nil {
		p.Logger.Error(fmt.Sprintf("%s with error: %v", ErrParsingLabelSelector, err.Error()))
		return nil, customerrors.NewValidationError(ErrParsingLabelSelector)
	}

	err = p.client.List(p.Ctx, cappRevisionList, &client.ListOptions{
		Limit:         listOptions.Limit,
		Continue:      listOptions.Continue,
		Namespace:     p.namespace,
		LabelSelector: selector,
	})
	if err != nil {
		p.Logger.Error(fmt.Sprintf("%v with error: %v", ErrCouldNotListCappRevisions, err.Error()))
		return nil, err
	}

	return (*types.List[cappv1alpha1.CappRevision])(cappRevisionList), nil
}

// convertCappRevisionToType converts an API CappRevision to a Type CappRevision.
func convertCappRevisionToType(cappRevision cappv1alpha1.CappRevision) types.CappRevision {
	return types.CappRevision{
		Metadata: types.Metadata{
			Name:      cappRevision.Name,
			Namespace: cappRevision.Namespace,
		},
		Annotations: utils.ConvertMapToKeyValue(cappRevision.Annotations),
		Labels:      utils.ConvertMapToKeyValue(cappRevision.Labels),
		Spec: cappv1alpha1.CappRevisionSpec{
			RevisionNumber: cappRevision.Spec.RevisionNumber,
			CappTemplate:   cappRevision.Spec.CappTemplate,
		},
		Status: cappv1alpha1.CappRevisionStatus{},
	}
}
