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
	containerNamespace = testutils.TestNamespace + testutils.ContainersKey
	failedToGetPodErr  = "failed to get pod %q, in the namespace %q with error: pods %q not found"
)

func TestGetContainer(t *testing.T) {
	testNamespaceName := containerNamespace + "-get"

	type args struct {
		podName   string
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
		"ShouldSucceedGettingContainers": {
			args: args{
				namespace: testNamespaceName,
				podName:   pod2,
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					testutils.ContainersKey: []interface{}{
						map[string]interface{}{testutils.ContainerNameKey: testutils.TestContainerName},
						map[string]interface{}{testutils.ContainerNameKey: testutils.CappName},
					},
					testutils.CountKey: 2,
				},
			},
		},
		"ShouldHandleNotFoundPod": {
			args: args{
				namespace: testNamespaceName,
				podName:   pod1 + testutils.NonExistentSuffix,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					testutils.DetailsKey: fmt.Sprintf(failedToGetPodErr, pod1+testutils.NonExistentSuffix, testNamespaceName, pod1+testutils.NonExistentSuffix),
					testutils.ErrorKey:   testutils.OperationFailed,
				},
			},
		},
		"ShouldHandleNotFoundNamespace": {
			args: args{
				namespace: testNamespaceName + testutils.NonExistentSuffix,
				podName:   pod1,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					testutils.DetailsKey: fmt.Sprintf(failedToGetPodErr, pod1, testNamespaceName+testutils.NonExistentSuffix, pod1),
					testutils.ErrorKey:   testutils.OperationFailed,
				},
			},
		},
	}

	setup()
	mocks.CreateTestNamespace(fakeClient, testNamespaceName)
	mocks.CreateTestPod(fakeClient, testNamespaceName, pod1, "", false)
	mocks.CreateTestPod(fakeClient, testNamespaceName, pod2, testutils.CappName, true)

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			baseURI := fmt.Sprintf("/v1/namespaces/%s/pods/%s/containers", test.args.namespace, test.args.podName)
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
