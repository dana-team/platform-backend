package v1

import (
	"encoding/json"
	"fmt"
	"github.com/dana-team/platform-backend/src/utils/testutils"
	"github.com/dana-team/platform-backend/src/utils/testutils/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	podNamespace = testutils.TestNamespace + testutils.PodsKey
)

func TestGetPod(t *testing.T) {
	testNamespaceName := podNamespace + "-get"

	type args struct {
		cappName  string
		namespace string
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
		"ShouldNotGetPods": {
			args: args{
				namespace: testNamespaceName,
				cappName:  testutils.CappName + testutils.NonExistentSuffix,
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					testutils.CountKey: 0,
					testutils.PodsKey:  interface{}(nil),
				},
			},
		},
		"ShouldNotGetPodsOnNotFoundNamespace": {
			args: args{
				namespace: testNamespaceName + testutils.NonExistentSuffix,
				cappName:  testutils.CappName,
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					testutils.CountKey: 0,
					testutils.PodsKey:  interface{}(nil),
				},
			},
		},
	}

	setup()
	mocks.CreateTestNamespace(fakeClient, testNamespaceName)
	mocks.CreateTestPod(fakeClient, testNamespaceName, pod1, testutils.CappName, false)
	mocks.CreateTestPod(fakeClient, testNamespaceName, pod2, testutils.CappName, true)
	mocks.CreateTestPod(fakeClient, testNamespaceName, pod3, "", false)

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			baseURI := fmt.Sprintf("/v1/namespaces/%s/capps/%s/pods", test.args.namespace, test.args.cappName)
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
