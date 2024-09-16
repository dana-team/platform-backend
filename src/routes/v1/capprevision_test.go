package v1

import (
	"encoding/json"
	"fmt"
	cappv1alpha1 "github.com/dana-team/container-app-operator/api/v1alpha1"
	"github.com/dana-team/platform-backend/src/controllers"
	"github.com/dana-team/platform-backend/src/middleware"
	"github.com/dana-team/platform-backend/src/types"
	"github.com/dana-team/platform-backend/src/utils/testutils"
	"github.com/dana-team/platform-backend/src/utils/testutils/mocks"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

const (
	cappRevisionName      = testutils.TestName + "-capp-revision"
	capprevisionsKey      = "capprevisions"
	cappRevisionNamespace = testutils.TestNamespace + "-" + capprevisionsKey
)

func TestGetCappRevisions(t *testing.T) {
	testNamespaceName := cappRevisionNamespace + "-get"

	type pagination struct {
		limit string
		page  string
	}

	type requestURI struct {
		namespace        string
		cappName         string
		paginationParams pagination
		clusterName      string
	}

	type want struct {
		statusCode int
		response   map[string]interface{}
	}

	cases := map[string]struct {
		requestURI requestURI
		want       want
	}{
		"ShouldSucceedGettingCappRevisionsOfCapp": {
			requestURI: requestURI{
				namespace: testNamespaceName,
				cappName:  testutils.CappName + "-1",
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					capprevisionsKey:   []string{cappRevisionName + "-1", cappRevisionName + "-2"},
					testutils.CountKey: 2,
				},
			},
		},
		"ShouldSucceedGettingCappRevisions": {
			requestURI: requestURI{
				namespace:   testNamespaceName,
				cappName:    "",
				clusterName: cluster,
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					capprevisionsKey:   []string{cappRevisionName + "-1", cappRevisionName + "-2", cappRevisionName + "-3"},
					testutils.CountKey: 3,
				},
			},
		},
		"ShouldSucceedGettingAllCappRevisionsWithLimitOf2": {
			requestURI: requestURI{
				namespace:        testNamespaceName,
				paginationParams: pagination{limit: "2", page: "1"},
				cappName:         testutils.CappName + "-1",
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					capprevisionsKey:   []string{cappRevisionName + "-1", cappRevisionName + "-2"},
					testutils.CountKey: 2,
				},
			},
		},
		"ShouldSucceedGettingCappRevisionsWithLabelSelector": {
			requestURI: requestURI{
				namespace: testNamespaceName,
				cappName:  testutils.CappName + "-1",
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					capprevisionsKey:   []string{cappRevisionName + "-1", cappRevisionName + "-2"},
					testutils.CountKey: 2,
				},
			},
		},
	}

	setup()
	mocks.CreateTestNamespace(fakeClient, testNamespaceName)
	mocks.CreateTestCapp(dynClient, testutils.CappName+"-1", testNamespaceName, testutils.Domain, map[string]string{}, map[string]string{})
	mocks.CreateTestCappRevision(dynClient, cappRevisionName+"-1", testNamespaceName, map[string]string{testutils.LabelCappName: testutils.CappName + "-1"}, nil)
	mocks.CreateTestCappRevision(dynClient, cappRevisionName+"-2", testNamespaceName, map[string]string{testutils.LabelCappName: testutils.CappName + "-1"}, nil)
	mocks.CreateTestCappRevision(dynClient, cappRevisionName+"-3", testNamespaceName, map[string]string{testutils.LabelCappName: testutils.CappName + "-2"}, nil)

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			params := url.Values{}

			if test.requestURI.paginationParams.limit != "" {
				params.Add(middleware.LimitCtxKey, test.requestURI.paginationParams.limit)
			}

			if test.requestURI.paginationParams.page != "" {
				params.Add(middleware.PageCtxKey, test.requestURI.paginationParams.page)
			}

			baseURI := fmt.Sprintf("/v1/namespaces/%s/capps/%s/capprevisions", test.requestURI.namespace, test.requestURI.cappName)

			if test.requestURI.clusterName != "" {
				baseURI = fmt.Sprintf("/v1/clusters/%s/namespaces/%s/capprevisions", test.requestURI.clusterName, test.requestURI.namespace)
			}

			request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s?%s", baseURI, params.Encode()), nil)
			assert.NoError(t, err)
			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, request)

			assert.Equal(t, test.want.statusCode, writer.Code)

			var response map[string]interface{}
			err = json.Unmarshal(writer.Body.Bytes(), &response)
			assert.NoError(t, err)

			wantResponseJSON, err := json.Marshal(test.want.response)
			assert.NoError(t, err)
			var wantResponseNormalized map[string]interface{}
			err = json.Unmarshal(wantResponseJSON, &wantResponseNormalized)
			assert.NoError(t, err)
			assert.Equal(t, wantResponseNormalized, response)
		})
	}
}

func TestGetCappRevision(t *testing.T) {
	testNamespaceName := cappRevisionNamespace + "-get-one"
	labels := []types.KeyValue{{Key: testutils.LabelCappName, Value: testutils.CappName + "-2"}}

	type requestURI struct {
		name      string
		namespace string
		cappName  string
	}

	type want struct {
		statusCode int
		response   map[string]interface{}
	}

	cases := map[string]struct {
		requestURI requestURI
		want       want
	}{
		"ShouldSucceedGettingCappRevision": {
			requestURI: requestURI{
				namespace: testNamespaceName,
				name:      cappRevisionName,
				cappName:  testutils.CappName + "2",
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					testutils.MetadataKey:    types.Metadata{Name: cappRevisionName, Namespace: testNamespaceName},
					testutils.LabelsKey:      labels,
					testutils.AnnotationsKey: nil,
					testutils.SpecKey:        mocks.PrepareCappRevisionSpec(mocks.ConvertKeyValueSliceToMap(labels), map[string]string{}),
					testutils.StatusKey:      mocks.PrepareCappRevisionStatus(),
				},
			},
		},
		"ShouldHandleNotFoundCappRevision": {
			requestURI: requestURI{
				namespace: testNamespaceName,
				name:      cappRevisionName + testutils.NonExistentSuffix,
				cappName:  testutils.CappName + "2",
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					testutils.ErrorKey: fmt.Sprintf("%v, %v",
						fmt.Sprintf(controllers.ErrCouldNotGetCappRevision, cappRevisionName+testutils.NonExistentSuffix, testNamespaceName),
						fmt.Sprintf("%s.%s %q not found", capprevisionsKey, cappv1alpha1.GroupVersion.Group, cappRevisionName+testutils.NonExistentSuffix)),
					testutils.ReasonKey: metav1.StatusReasonNotFound,
				},
			},
		},
		"ShouldHandleNotFoundNamespace": {
			requestURI: requestURI{
				namespace: testNamespaceName + testutils.NonExistentSuffix,
				name:      cappRevisionName,
				cappName:  testutils.CappName + "2",
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					testutils.ErrorKey:  fmt.Sprintf("capps.%s %q not found", cappv1alpha1.GroupVersion.Group, testutils.CappName+"2"),
					testutils.ReasonKey: metav1.StatusReasonNotFound,
				},
			},
		},
	}

	setup()
	mocks.CreateTestNamespace(fakeClient, testNamespaceName)
	mocks.CreateTestCapp(dynClient, testutils.CappName+"2", testNamespaceName, testutils.Domain, map[string]string{}, map[string]string{})
	mocks.CreateTestCappRevision(dynClient, cappRevisionName, testNamespaceName, map[string]string{testutils.LabelCappName: testutils.CappName + "-2"}, nil)

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			baseURI := fmt.Sprintf("/v1/namespaces/%s/capps/%s/capprevisions/%s", test.requestURI.namespace, test.requestURI.cappName, test.requestURI.name)
			request, err := http.NewRequest(http.MethodGet, baseURI, nil)
			assert.NoError(t, err)
			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, request)

			assert.Equal(t, test.want.statusCode, writer.Code)

			var response map[string]interface{}
			err = json.Unmarshal(writer.Body.Bytes(), &response)
			assert.NoError(t, err)

			wantResponseJSON, err := json.Marshal(test.want.response)
			assert.NoError(t, err)
			var wantResponseNormalized map[string]interface{}
			err = json.Unmarshal(wantResponseJSON, &wantResponseNormalized)
			assert.NoError(t, err)
			assert.Equal(t, wantResponseNormalized, response)
		})
	}
}
