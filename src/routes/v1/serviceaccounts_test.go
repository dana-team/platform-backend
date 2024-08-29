package v1

import (
	"encoding/json"
	"fmt"
	"github.com/dana-team/platform-backend/src/controllers"
	"github.com/dana-team/platform-backend/src/utils/testutils"
	"github.com/dana-team/platform-backend/src/utils/testutils/mocks"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	serviceAccountNamespace = testutils.TestNamespace + testutils.TokenKey
)

func TestGetServiceAccountToken(t *testing.T) {
	testNamespaceName := serviceAccountNamespace + "-get-token"

	type args struct {
		serviceAccountName string
		namespace          string
	}

	type want struct {
		statusCode int
		response   map[string]interface{}
	}

	cases := map[string]struct {
		args args
		want want
	}{
		"ShouldSucceedGettingToken": {
			args: args{
				namespace:          testNamespaceName,
				serviceAccountName: testutils.ServiceAccountName,
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					testutils.TokenKey: testutils.Value,
				},
			},
		},
		"ShouldNotSucceedGettingTokenForNonExistingNamespace": {
			args: args{
				namespace:          testNamespaceName + testutils.NonExistentSuffix,
				serviceAccountName: testutils.ServiceAccountName,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					testutils.ErrorKey: fmt.Sprintf("%s, %s",
						fmt.Sprintf(controllers.ErrCouldNotGetServiceAccount, testutils.ServiceAccountName, testNamespaceName+testutils.NonExistentSuffix),
						fmt.Sprintf("serviceaccounts %q not found", testutils.ServiceAccountName),
					),
					testutils.ReasonKey: testutils.ReasonNotFound,
				},
			},
		},
		"ShouldNotSucceedGettingTokenForNonExistingName": {
			args: args{
				namespace:          testNamespaceName,
				serviceAccountName: testutils.ServiceAccountName + testutils.NonExistentSuffix,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					testutils.ErrorKey: fmt.Sprintf("%s, %s",
						fmt.Sprintf(controllers.ErrCouldNotGetServiceAccount, testutils.ServiceAccountName+testutils.NonExistentSuffix, testNamespaceName),
						fmt.Sprintf("serviceaccounts %q not found", testutils.ServiceAccountName+testutils.NonExistentSuffix),
					),
					testutils.ReasonKey: testutils.ReasonNotFound,
				},
			},
		},
	}

	setup()
	mocks.CreateTestNamespace(fakeClient, testNamespaceName)
	mocks.CreateTestServiceAccountWithToken(fakeClient, testNamespaceName, testutils.ServiceAccountName, "token-secret", "value", "docker-cfg")

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			params := url.Values{}

			baseURI := fmt.Sprintf("/v1/namespaces/%s/serviceaccounts/%s/token", test.args.namespace, test.args.serviceAccountName)
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
