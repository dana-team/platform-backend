package controllers

import (
	"fmt"
	"testing"

	"github.com/dana-team/platform-backend/internal/types"

	"github.com/dana-team/platform-backend/internal/customerrors"
	"github.com/dana-team/platform-backend/internal/utils"
	"github.com/dana-team/platform-backend/internal/utils/testutils"
	"github.com/dana-team/platform-backend/internal/utils/testutils/mocks"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestCreateToken(t *testing.T) {
	namespaceName := testutils.TokenNamespace + "-create-token"
	type requestParams struct {
		name                 string
		namespace            string
		createServiceAccount bool
		expiration           string
	}
	type want struct {
		response    types.TokenRequestResponse
		errorStatus metav1.StatusReason
	}
	cases := map[string]struct {
		request requestParams
		want    want
	}{
		"ShouldSucceedCreatingToken": {
			request: requestParams{
				name:                 testutils.ServiceAccountName,
				namespace:            namespaceName,
				createServiceAccount: true,
				expiration:           "3600",
			},
			want: want{
				response:    types.TokenRequestResponse{Token: ""},
				errorStatus: metav1.StatusSuccess,
			},
		},
		"ShouldFailCreatingTokenWithoutServiceAccount": {
			request: requestParams{
				name:                 testutils.ServiceAccountName + testutils.NonExistentSuffix,
				namespace:            namespaceName,
				createServiceAccount: false,
				expiration:           "3600",
			},
			want: want{
				response:    types.TokenRequestResponse{Token: ""},
				errorStatus: metav1.StatusReasonNotFound,
			},
		},
		"ShouldFailCreatingTokenWithInvalidExpiration": {
			request: requestParams{
				name:                 testutils.ServiceAccountName + "-invalid-expiration",
				namespace:            namespaceName,
				createServiceAccount: true,
				expiration:           testutils.TestName,
			},
			want: want{
				response:    types.TokenRequestResponse{Token: ""},
				errorStatus: metav1.StatusReasonBadRequest,
			},
		},
	}
	setup()
	tokenController := NewTokenController(fakeClient, mocks.GinContext(), logger)
	createTestNamespace(namespaceName, utils.AddManagedLabel(map[string]string{}))

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if tc.request.createServiceAccount {
				mocks.CreateTestServiceAccount(fakeClient, tc.request.namespace, tc.request.name, "")
			}
			response, err := tokenController.CreateToken(tc.request.name, tc.request.namespace, tc.request.expiration)
			if tc.want.errorStatus != metav1.StatusSuccess {
				reason := err.(customerrors.ErrorWithStatusCode).StatusReason()

				assert.Equal(t, tc.want.errorStatus, reason)
			} else {
				assert.NoError(t, err)
				// Since there is no TokenController running, no token is created so this is not worth much.
				// TODO: Rework this when https://github.com/kubernetes-sigs/controller-runtime/pull/2969 is released.
				assert.Equal(t, tc.want.response, response)

			}
		})
	}
}

func TestRevokeToken(t *testing.T) {
	namespaceName := testutils.TokenNamespace + "-revoke-token"
	type requestParams struct {
		name                 string
		namespace            string
		secretName           string
		createServiceAccount bool
	}
	type want struct {
		errorStatus metav1.StatusReason
	}
	cases := map[string]struct {
		requestParams requestParams
		want          want
	}{
		"ShouldSucceedRevokingToken": {
			requestParams: requestParams{
				name:                 testutils.ServiceAccountName,
				namespace:            namespaceName,
				secretName:           fmt.Sprintf("%s-%s", testutils.ServiceAccountName, testutils.TokenRequestSuffix),
				createServiceAccount: true,
			},
			want: want{
				errorStatus: metav1.StatusSuccess,
			},
		},
		"ShouldFailRevokingTokenWithoutServiceAccount": {
			requestParams: requestParams{
				name:                 testutils.ServiceAccountName + testutils.NonExistentSuffix,
				namespace:            namespaceName,
				secretName:           fmt.Sprintf("%s-%s", testutils.ServiceAccountName+testutils.NonExistentSuffix, testutils.TokenRequestSuffix),
				createServiceAccount: false,
			},
			want: want{
				errorStatus: metav1.StatusReasonNotFound,
			},
		},
		"ShouldFailRevokingTokenWithoutSecret": {
			requestParams: requestParams{
				name:                 testutils.ServiceAccountName + "-1",
				namespace:            namespaceName,
				secretName:           "",
				createServiceAccount: true,
			},
			want: want{
				errorStatus: metav1.StatusReasonNotFound,
			},
		},
	}
	setup()
	tokenController := NewTokenController(fakeClient, mocks.GinContext(), logger)
	createTestNamespace(namespaceName, utils.AddManagedLabel(map[string]string{}))
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if tc.requestParams.createServiceAccount {
				mocks.CreateTestServiceAccount(fakeClient, tc.requestParams.namespace, tc.requestParams.name, "")
			}
			if tc.requestParams.secretName != "" {
				createTestSecret(tc.requestParams.secretName, tc.requestParams.namespace, utils.AddManagedLabel(map[string]string{}))
			}
			err := tokenController.RevokeToken(tc.requestParams.name, tc.requestParams.namespace)
			if tc.want.errorStatus != metav1.StatusSuccess {
				reason := err.(customerrors.ErrorWithStatusCode).StatusReason()
				assert.Equal(t, tc.want.errorStatus, reason)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
