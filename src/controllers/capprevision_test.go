package controllers_test

import (
	"context"
	"fmt"
	cappv1alpha1 "github.com/dana-team/container-app-operator/api/v1alpha1"
	"github.com/dana-team/platform-backend/src/controllers"
	"github.com/dana-team/platform-backend/src/controllers/mocks"
	"github.com/dana-team/platform-backend/src/types"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"testing"
)

const (
	cappRevisionNamespace = testName + "-capp-revision-ns"
	cappRevisionName      = testName + "-capp-revision"
	labelKey              = "key"
	labelValue            = "value"
)

var controller controllers.CappRevisionController

func setupCappRevisions() {
	controller = controllers.NewCappRevisionController(dynClient, context.TODO(), logger)
}

// createTestCappRevision creates a test CappRevision object.
func createTestCappRevision(name, namespace string, labels, annotations map[string]string) {
	cappRevision := mocks.PrepareCappRevision(name, namespace, labels, annotations)
	err := dynClient.Create(context.TODO(), &cappRevision)
	if err != nil {
		panic(err)
	}
}

func TestGetCappRevision(t *testing.T) {
	namespaceName := cappRevisionNamespace + "-GetOne"

	type requestParams struct {
		name      string
		namespace string
	}

	type want struct {
		cappRevision types.CappRevision
		errorStatus  v1.StatusReason
	}

	cases := map[string]struct {
		requestParams requestParams
		want          want
	}{
		"ShouldSucceedGettingCappRevision": {
			requestParams: requestParams{
				namespace: namespaceName,
				name:      cappRevisionName + "1",
			},
			want: want{
				cappRevision: types.CappRevision{
					Metadata: types.Metadata{Name: cappRevisionName + "1", Namespace: namespaceName},
					Labels:   []types.KeyValue{{Key: labelKey + "-1", Value: labelValue + "-1"}},
					Spec: cappv1alpha1.CappRevisionSpec{
						RevisionNumber: 1,
						CappTemplate:   cappv1alpha1.CappTemplate{Spec: cappv1alpha1.CappSpec{}},
					},
					Status: cappv1alpha1.CappRevisionStatus{},
				},
				errorStatus: v1.StatusSuccess,
			},
		},
		"ShouldFailGettingNonExistingCappRevision": {
			requestParams: requestParams{
				namespace: namespaceName,
				name:      doesNotExist,
			},
			want: want{
				cappRevision: types.CappRevision{},
				errorStatus:  v1.StatusReasonNotFound,
			},
		},
		"ShouldFailGettingNonExistingNamespace": {
			requestParams: requestParams{
				namespace: doesNotExist,
				name:      cappRevisionName,
			},
			want: want{
				cappRevision: types.CappRevision{},
				errorStatus:  v1.StatusReasonNotFound,
			},
		},
	}
	createTestNamespace(namespaceName)
	createTestCappRevision(cappRevisionName+"1", namespaceName, map[string]string{labelKey + "-1": labelValue + "-1"}, map[string]string{})
	createTestCappRevision(cappRevisionName+"2", namespaceName, map[string]string{labelKey + "-2": labelValue + "-2"}, map[string]string{})
	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			response, err := controller.GetCappRevision(test.requestParams.namespace, test.requestParams.name)
			if test.want.errorStatus != v1.StatusSuccess {
				reason := err.(errors.APIStatus).Status().Reason

				assert.Equal(t, test.want.errorStatus, reason)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, test.want.cappRevision, response)
		})

	}

}

func TestGetCappRevisions(t *testing.T) {
	namespaceName := cappRevisionNamespace + "-GetMany"

	type requestParams struct {
		cappQuery types.CappRevisionQuery
		namespace string
	}

	type want struct {
		cappRevisions types.CappRevisionList
		errorStatus   v1.StatusReason
	}

	cases := map[string]struct {
		requestParams requestParams
		want          want
	}{
		"ShouldSucceedGettingAllCappRevisions": {
			requestParams: requestParams{
				namespace: namespaceName,
				cappQuery: types.CappRevisionQuery{},
			},
			want: want{
				cappRevisions: types.CappRevisionList{CappRevisions: []string{cappRevisionName + "1", cappRevisionName + "2"}, Count: 2},
				errorStatus:   v1.StatusSuccess,
			}},
		"ShouldSucceedGettingCappRevisionsByLabels": {
			requestParams: requestParams{
				namespace: namespaceName,
				cappQuery: types.CappRevisionQuery{LabelSelector: fmt.Sprintf("%s-2=%s-2", labelKey, labelValue)},
			},
			want: want{
				cappRevisions: types.CappRevisionList{CappRevisions: []string{cappRevisionName + "2"}, Count: 1},
				errorStatus:   v1.StatusSuccess,
			},
		},
		"ShouldThrowErrorWithInvalidSelector": {
			requestParams: requestParams{
				namespace: namespaceName,
				cappQuery: types.CappRevisionQuery{LabelSelector: ":--"},
			},
			want: want{
				cappRevisions: types.CappRevisionList{},
				errorStatus:   v1.StatusReasonBadRequest,
			},
		},
		"ShouldFailGettingNonExistingNamespace": {
			requestParams: requestParams{
				namespace: doesNotExist,
				cappQuery: types.CappRevisionQuery{},
			},
			want: want{
				cappRevisions: types.CappRevisionList{},
				errorStatus:   v1.StatusSuccess,
			},
		},
	}
	createTestNamespace(namespaceName)
	createTestCappRevision(cappRevisionName+"1", namespaceName, map[string]string{labelKey + "-1": labelValue + "-1"}, map[string]string{})
	createTestCappRevision(cappRevisionName+"2", namespaceName, map[string]string{labelKey + "-2": labelValue + "-2"}, map[string]string{})
	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			response, err := controller.GetCappRevisions(test.requestParams.namespace, test.requestParams.cappQuery)
			if test.want.errorStatus != v1.StatusSuccess {
				reason := err.(errors.APIStatus).Status().Reason

				assert.Equal(t, test.want.errorStatus, reason)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, test.want.cappRevisions, response)
		})

	}

}
