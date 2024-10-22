package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/dana-team/platform-backend/internal/controllers"

	"github.com/dana-team/platform-backend/internal/utils/testutils"
	"github.com/dana-team/platform-backend/internal/utils/testutils/mocks"
	"github.com/stretchr/testify/assert"
)

func TestCreateToken(t *testing.T) {
	namespaceName := testutils.TestNamespace + "-create-token"
	type args struct {
		serviceAccountName   string
		namespaceName        string
		createServiceAccount bool
	}
	type want struct {
		statusCode int
		response   map[string]interface{}
	}
	tests := map[string]struct {
		args args
		want want
	}{
		"ShouldCreateToken": {
			args: args{
				serviceAccountName:   testutils.ServiceAccountName,
				namespaceName:        namespaceName,
				createServiceAccount: true,
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					"token":   "",
					"expires": time.Time{},
				},
			},
		},
		"ShouldNotCreateTokenWhenServiceAccountDoesNotExist": {
			args: args{
				serviceAccountName:   testutils.ServiceAccountName + testutils.NonExistentSuffix,
				namespaceName:        namespaceName,
				createServiceAccount: false,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					testutils.ErrorKey: fmt.Sprintf("%s, %s",
						fmt.Sprintf(controllers.ErrCouldNotGetServiceAccount, testutils.ServiceAccountName+testutils.NonExistentSuffix, namespaceName),
						fmt.Sprintf("serviceaccounts %q not found", testutils.ServiceAccountName+testutils.NonExistentSuffix),
					),
					testutils.ReasonKey: testutils.ReasonNotFound,
				},
			},
		},
	}
	setup()
	mocks.CreateTestNamespace(fakeClient, namespaceName)

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			params := url.Values{}

			if tc.args.createServiceAccount {
				mocks.CreateTestServiceAccount(fakeClient, tc.args.namespaceName, tc.args.serviceAccountName, "")
			}

			baseURI := fmt.Sprintf("/v1/namespaces/%s/serviceaccounts/%s/token", tc.args.namespaceName, tc.args.serviceAccountName)
			request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s?%s", baseURI, params.Encode()), nil)
			assert.NoError(t, err)
			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, request)

			assert.Equal(t, tc.want.statusCode, writer.Code)

			var response map[string]interface{}
			err = json.Unmarshal(writer.Body.Bytes(), &response)
			assert.NoError(t, err)

			wantResponseJSON, err := json.Marshal(tc.want.response)
			assert.NoError(t, err)
			var wantResponseNormalized map[string]interface{}
			err = json.Unmarshal(wantResponseJSON, &wantResponseNormalized)
			assert.NoError(t, err)
			assert.Equal(t, wantResponseNormalized, response)
		})
	}
}

func TestRevokeToken(t *testing.T) {
	namespaceName := testutils.TestNamespace + "-revoke-token"
	type args struct {
		serviceAccountName   string
		namespaceName        string
		secretName           string
		createServiceAccount bool
	}
	type want struct {
		statusCode int
		response   map[string]interface{}
	}
	tests := map[string]struct {
		args args
		want want
	}{
		"ShouldRevokeToken": {
			args: args{
				serviceAccountName:   testutils.ServiceAccountName,
				namespaceName:        namespaceName,
				secretName:           fmt.Sprintf("%s-%s", testutils.ServiceAccountName, testutils.TokenRequestSuffix),
				createServiceAccount: true,
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					"message": fmt.Sprintf("Revoked tokens for ServiceAccount %q", testutils.ServiceAccountName),
				},
			},
		},
		"ShouldNotRevokeTokenWhenSecretNotFound": {
			args: args{
				serviceAccountName:   testutils.ServiceAccountName + "-no-secret",
				namespaceName:        namespaceName,
				secretName:           "",
				createServiceAccount: true,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					testutils.ErrorKey: fmt.Sprintf("%s, %s",
						fmt.Sprintf(controllers.ErrCouldNotGetTokenRequestSecret, testutils.ServiceAccountName+"-no-secret", namespaceName),
						fmt.Sprintf("secrets %q not found", fmt.Sprintf("%s-%s", testutils.ServiceAccountName+"-no-secret", testutils.TokenRequestSuffix)),
					),
					testutils.ReasonKey: testutils.ReasonNotFound,
				},
			},
		},
		"ShouldNotRevokeTokenWhenServiceAccountDoesNotExist": {
			args: args{
				serviceAccountName:   testutils.ServiceAccountName + testutils.NonExistentSuffix,
				namespaceName:        namespaceName,
				secretName:           fmt.Sprintf("%s-%s", testutils.ServiceAccountName+testutils.NonExistentSuffix, testutils.TokenRequestSuffix),
				createServiceAccount: false,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					testutils.ErrorKey: fmt.Sprintf("%s, %s",
						fmt.Sprintf(controllers.ErrCouldNotGetServiceAccount, testutils.ServiceAccountName+testutils.NonExistentSuffix, namespaceName),
						fmt.Sprintf("serviceaccounts %q not found", testutils.ServiceAccountName+testutils.NonExistentSuffix),
					),
					testutils.ReasonKey: testutils.ReasonNotFound,
				},
			},
		},
	}
	setup()
	mocks.CreateTestNamespace(fakeClient, namespaceName)

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			params := url.Values{}

			if tc.args.createServiceAccount {
				mocks.CreateTestServiceAccount(fakeClient, tc.args.namespaceName, tc.args.serviceAccountName, "")
			}
			if tc.args.secretName != "" {
				mocks.CreateTestSecret(fakeClient, tc.args.secretName, tc.args.namespaceName)
			}
			baseURI := fmt.Sprintf("/v1/namespaces/%s/serviceaccounts/%s/token", tc.args.namespaceName, tc.args.serviceAccountName)
			request, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s?%s", baseURI, params.Encode()), nil)
			assert.NoError(t, err)
			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, request)

			assert.Equal(t, tc.want.statusCode, writer.Code)

			var response map[string]interface{}
			err = json.Unmarshal(writer.Body.Bytes(), &response)
			assert.NoError(t, err)

			wantResponseJSON, err := json.Marshal(tc.want.response)
			assert.NoError(t, err)
			var wantResponseNormalized map[string]interface{}
			err = json.Unmarshal(wantResponseJSON, &wantResponseNormalized)
			assert.NoError(t, err)
			assert.Equal(t, wantResponseNormalized, response)
		})
	}
}
