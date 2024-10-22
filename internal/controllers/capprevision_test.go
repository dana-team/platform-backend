package controllers

import (
	"context"
	cappv1alpha1 "github.com/dana-team/container-app-operator/api/v1alpha1"
	"github.com/dana-team/platform-backend/internal/customerrors"
	"github.com/dana-team/platform-backend/internal/types"
	"github.com/dana-team/platform-backend/internal/utils/pagination"
	"github.com/dana-team/platform-backend/internal/utils/testutils"
	"github.com/dana-team/platform-backend/internal/utils/testutils/mocks"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

// createTestCappRevision creates a test CappRevision object.
func createTestCappRevision(name, namespace string, labels, annotations map[string]string) {
	cappRevision := mocks.PrepareCappRevision(name, namespace, testutils.SiteName, labels, annotations)
	err := dynClient.Create(context.TODO(), &cappRevision)
	if err != nil {
		panic(err)
	}
}

func TestGetCappRevision(t *testing.T) {
	namespaceName := testutils.CappRevisionNamespace + "-get"
	labels := []types.KeyValue{{Key: testutils.LabelKey + "-1", Value: testutils.LabelValue + "-1"}}

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
					Labels:   labels,
					Spec: cappv1alpha1.CappRevisionSpec{
						RevisionNumber: 1,
						CappTemplate:   cappv1alpha1.CappTemplate{Spec: mocks.PrepareCappSpec(testutils.SiteName), Labels: mocks.ConvertKeyValueSliceToMap(labels)},
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

	cappRevisionController := NewCappRevisionController(dynClient, mocks.GinContext(), logger)
	createTestNamespace(namespaceName, map[string]string{})
	createTestCappRevision(testutils.CappRevisionName+"-1", namespaceName, map[string]string{testutils.LabelKey + "-1": testutils.LabelValue + "-1"}, map[string]string{})
	createTestCappRevision(testutils.CappRevisionName+"-2", namespaceName, map[string]string{testutils.LabelKey + "-2": testutils.LabelValue + "-2"}, map[string]string{})

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			response, err := cappRevisionController.GetCappRevision(test.requestParams.namespace, test.requestParams.name)
			if test.want.errorStatus != metav1.StatusSuccess {
				reason := err.(customerrors.ErrorWithStatusCode).StatusReason()
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
		cappName  string
		namespace string
		limit     int
		page      int
	}

	type want struct {
		cappRevisions types.CappRevisionList
		errorStatus   metav1.StatusReason
	}

	cases := map[string]struct {
		requestParams requestParams
		want          want
	}{
		"ShouldSucceedGettingCappRevisionsOfCapp": {
			requestParams: requestParams{
				namespace: namespaceName,
				cappName:  testutils.CappName + "-1",
			},
			want: want{
				cappRevisions: types.CappRevisionList{CappRevisions: []string{testutils.CappRevisionName + "-1"}, ListMetadata: types.ListMetadata{Count: 1}},
				errorStatus:   metav1.StatusSuccess,
			},
		},
		"ShouldSucceedGettingCappRevisions": {
			requestParams: requestParams{
				namespace: namespaceName,
				cappName:  "",
			},
			want: want{
				cappRevisions: types.CappRevisionList{CappRevisions: []string{testutils.CappRevisionName + "-1", testutils.CappRevisionName + "-2"}, ListMetadata: types.ListMetadata{Count: 2}},
				errorStatus:   metav1.StatusSuccess,
			},
		},
		"ShouldFailGettingNonExistingNamespace": {
			requestParams: requestParams{
				namespace: testutils.CappRevisionNamespace + testutils.NonExistentSuffix,
				cappName:  testutils.CappName,
			},
			want: want{
				cappRevisions: types.CappRevisionList{},
				errorStatus:   metav1.StatusSuccess,
			},
		},
	}
	setup()
	createTestNamespace(namespaceName, map[string]string{})
	createTestCappRevision(testutils.CappRevisionName+"-1", namespaceName, map[string]string{testutils.LabelCappName: testutils.CappName + "-1"}, map[string]string{})
	createTestCappRevision(testutils.CappRevisionName+"-2", namespaceName, map[string]string{testutils.LabelCappName: testutils.CappName + "-2"}, map[string]string{})

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			c := mocks.GinContext()
			mocks.SetPaginationValues(c, test.requestParams.limit, test.requestParams.page)
			cappRevisionController := NewCappRevisionController(dynClient, c, logger)

			limit, page, _ := pagination.ExtractPaginationParamsFromCtx(c)
			response, err := cappRevisionController.GetCappRevisions(test.requestParams.namespace, limit, page, test.requestParams.cappName)
			if test.want.errorStatus != metav1.StatusSuccess {
				reason := err.(customerrors.ErrorWithStatusCode).StatusReason()
				assert.Equal(t, test.want.errorStatus, reason)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, test.want.cappRevisions, response)
		})

	}

}
