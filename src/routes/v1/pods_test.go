package v1

import (
	"encoding/json"
	"fmt"
	cappv1alpha1 "github.com/dana-team/container-app-operator/api/v1alpha1"
	"github.com/dana-team/platform-backend/src/middleware"
	"github.com/dana-team/platform-backend/src/utils/testutils"
	"github.com/dana-team/platform-backend/src/utils/testutils/mocks"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	podNamespace = testutils.TestNamespace + testutils.PodsKey
)

func TestGetPods(t *testing.T) {
	testNamespaceName := podNamespace + "-get"

	type pagination struct {
		limit string
		page  string
	}

	type args struct {
		cappName         string
		namespace        string
		paginationParams pagination
	}

	type want struct {
		statusCode int
		response   map[string]interface{}
	}

	cases := map[string]struct {
		args args
		want want
	}{
		"ShouldSucceedGettingPods": {
			args: args{
				namespace: testNamespaceName,
				cappName:  testutils.CappName,
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					testutils.PodsKey: []interface{}{
						map[string]interface{}{testutils.PodNameKey: pod1},
						map[string]interface{}{testutils.PodNameKey: pod2},
					},
					testutils.CountKey: 2,
				},
			},
		},
		"ShouldSucceedGettingAllPodsWithLimitOf2": {
			args: args{
				namespace:        testNamespaceName,
				cappName:         testutils.CappName,
				paginationParams: pagination{limit: "2", page: "1"},
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					testutils.PodsKey: []interface{}{
						map[string]interface{}{testutils.PodNameKey: pod1},
						map[string]interface{}{testutils.PodNameKey: pod2},
					},
					testutils.CountKey: 2,
				},
			},
		},
		"ShouldNotGetPodsOnNotFoundCapp": {
			args: args{
				namespace: testNamespaceName,
				cappName:  testutils.CappName + testutils.NonExistentSuffix,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					testutils.ErrorKey:  fmt.Sprintf("%s.%s %q not found", testutils.CappsKey, cappv1alpha1.GroupVersion.Group, testutils.CappName+testutils.NonExistentSuffix),
					testutils.ReasonKey: testutils.ReasonNotFound,
				},
			},
		},
		"ShouldNotGetPodsOnNotFoundNamespace": {
			args: args{
				namespace: testNamespaceName + testutils.NonExistentSuffix,
				cappName:  testutils.CappName,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					testutils.ErrorKey:  fmt.Sprintf("%s.%s %q not found", testutils.CappsKey, cappv1alpha1.GroupVersion.Group, testutils.CappName),
					testutils.ReasonKey: testutils.ReasonNotFound,
				},
			},
		},
	}

	setup()
	mocks.CreateTestNamespace(fakeClient, testNamespaceName)
	mocks.CreateTestCapp(dynClient, testutils.CappName, testNamespaceName, testutils.Domain, testutils.SiteName, map[string]string{}, map[string]string{})
	mocks.CreateTestPod(fakeClient, testNamespaceName, pod1, testutils.CappName, false)
	mocks.CreateTestPod(fakeClient, testNamespaceName, pod2, testutils.CappName, true)
	mocks.CreateTestPod(fakeClient, testNamespaceName, pod3, "", false)

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			params := url.Values{}
			if test.args.paginationParams.limit != "" {
				params.Add(middleware.LimitCtxKey, test.args.paginationParams.limit)
			}

			if test.args.paginationParams.page != "" {
				params.Add(middleware.PageCtxKey, test.args.paginationParams.page)
			}

			baseURI := fmt.Sprintf("/v1/namespaces/%s/capps/%s/pods", test.args.namespace, test.args.cappName)
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
