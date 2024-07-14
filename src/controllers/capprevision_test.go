package controllers

import (
	"context"
	"fmt"
	cappv1alpha1 "github.com/dana-team/container-app-operator/api/v1alpha1"
	"github.com/dana-team/platform-backend/src/types"
	"github.com/dana-team/platform-backend/src/utils/testutils"
	"github.com/dana-team/platform-backend/src/utils/testutils/mocks"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"testing"
)

// createTestCappRevision creates a test CappRevision object.
func createTestCappRevision(name, namespace string, labels, annotations map[string]string) {
	cappRevision := mocks.PrepareCappRevision(name, namespace, labels, annotations)
	err := dynClient.Create(context.TODO(), &cappRevision)
	if err != nil {
		panic(err)
	}
}

func TestGetCappRevision(t *testing.T) {
	namespaceName := testutils.CappRevisionNamespace + "-get"

	type requestParams struct {
		name      string
		namespace string
	}

	type want struct {
		cappRevision types.CappRevision
		errorStatus  metav1.StatusReason
	}

	cases := map[string]struct {
		requestParams requestParams
		want          want
	}{
		"ShouldSucceedGettingCappRevision": {
			requestParams: requestParams{
				namespace: namespaceName,
				name:      testutils.CappRevisionName + "-1",
			},
			want: want{
				cappRevision: types.CappRevision{
					Metadata: types.Metadata{Name: testutils.CappRevisionName + "-1", Namespace: namespaceName},
					Labels:   []types.KeyValue{{Key: testutils.LabelKey + "-1", Value: testutils.LabelValue + "-1"}},
					Spec: cappv1alpha1.CappRevisionSpec{
						RevisionNumber: 1,
						CappTemplate:   cappv1alpha1.CappTemplate{Spec: mocks.PrepareCappSpec()},
					},
					Status: cappv1alpha1.CappRevisionStatus{},
				},
				errorStatus: metav1.StatusSuccess,
			},
		},
		"ShouldFailGettingNonExistingCappRevision": {
			requestParams: requestParams{
				namespace: namespaceName,
				name:      testutils.CappRevisionName + testutils.NonExistentSuffix,
			},
			want: want{
				cappRevision: types.CappRevision{},
				errorStatus:  metav1.StatusReasonNotFound,
			},
		},
		"ShouldFailGettingNonExistingNamespace": {
			requestParams: requestParams{
				namespace: testutils.CappRevisionNamespace + testutils.NonExistentSuffix,
				name:      testutils.CappRevisionName,
			},
			want: want{
				cappRevision: types.CappRevision{},
				errorStatus:  metav1.StatusReasonNotFound,
			},
		},
	}
	setup()
	cappRevisionController := NewCappRevisionController(dynClient, context.TODO(), logger)
	createTestNamespace(namespaceName, map[string]string{})
	createTestCappRevision(testutils.CappRevisionName+"-1", namespaceName, map[string]string{testutils.LabelKey + "-1": testutils.LabelValue + "-1"}, map[string]string{})
	createTestCappRevision(testutils.CappRevisionName+"-2", namespaceName, map[string]string{testutils.LabelKey + "-2": testutils.LabelValue + "-2"}, map[string]string{})
	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			response, err := cappRevisionController.GetCappRevision(test.requestParams.namespace, test.requestParams.name)
			if test.want.errorStatus != metav1.StatusSuccess {
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
	namespaceName := testutils.CappRevisionNamespace + "-getmany"

	type requestParams struct {
		cappQuery types.CappRevisionQuery
		namespace string
	}

	type want struct {
		cappRevisions types.CappRevisionList
		errorStatus   metav1.StatusReason
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
				cappRevisions: types.CappRevisionList{CappRevisions: []string{testutils.CappRevisionName + "-1", testutils.CappRevisionName + "-2"}, Count: 2},
				errorStatus:   metav1.StatusSuccess,
			}},
		"ShouldSucceedGettingCappRevisionsByLabels": {
			requestParams: requestParams{
				namespace: namespaceName,
				cappQuery: types.CappRevisionQuery{LabelSelector: fmt.Sprintf("%s-2=%s-2", testutils.LabelKey, testutils.LabelValue)},
			},
			want: want{
				cappRevisions: types.CappRevisionList{CappRevisions: []string{testutils.CappRevisionName + "-2"}, Count: 1},
				errorStatus:   metav1.StatusSuccess,
			},
		},
		"ShouldThrowErrorWithInvalidSelector": {
			requestParams: requestParams{
				namespace: namespaceName,
				cappQuery: types.CappRevisionQuery{LabelSelector: testutils.InvalidLabelSelector},
			},
			want: want{
				cappRevisions: types.CappRevisionList{},
				errorStatus:   metav1.StatusReasonBadRequest,
			},
		},
		"ShouldFailGettingNonExistingNamespace": {
			requestParams: requestParams{
				namespace: testutils.CappRevisionNamespace + testutils.NonExistentSuffix,
				cappQuery: types.CappRevisionQuery{},
			},
			want: want{
				cappRevisions: types.CappRevisionList{},
				errorStatus:   metav1.StatusSuccess,
			},
		},
	}
	setup()
	cappRevisionController := NewCappRevisionController(dynClient, context.TODO(), logger)
	createTestNamespace(namespaceName, map[string]string{})
	createTestCappRevision(testutils.CappRevisionName+"-1", namespaceName, map[string]string{testutils.LabelKey + "-1": testutils.LabelValue + "-1"}, map[string]string{})
	createTestCappRevision(testutils.CappRevisionName+"-2", namespaceName, map[string]string{testutils.LabelKey + "-2": testutils.LabelValue + "-2"}, map[string]string{})
	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			response, err := cappRevisionController.GetCappRevisions(test.requestParams.namespace, test.requestParams.cappQuery)
			if test.want.errorStatus != metav1.StatusSuccess {
				reason := err.(errors.APIStatus).Status().Reason

				assert.Equal(t, test.want.errorStatus, reason)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, test.want.cappRevisions, response)
		})

	}

}
