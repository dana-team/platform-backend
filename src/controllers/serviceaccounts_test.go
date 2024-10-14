package controllers

import (
	"fmt"
	"testing"

	"github.com/dana-team/platform-backend/src/types"
	"github.com/dana-team/platform-backend/src/utils"
	"github.com/dana-team/platform-backend/src/utils/testutils"
	"github.com/dana-team/platform-backend/src/utils/testutils/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetServiceAccount(t *testing.T) {
	namespaceName := testutils.TestNamespace + "-getServiceAccount"
	type args struct {
		namespace                  string
		existingServiceAccountName string
		name                       string
	}
	type want struct {
		serviceAccount types.ServiceAccount
		error          string
	}
	cases := map[string]struct {
		args args
		want want
	}{
		"ShouldSucceedGettingServiceAccount": {
			args: args{
				namespace:                  namespaceName,
				name:                       testutils.ServiceAccountName,
				existingServiceAccountName: testutils.ServiceAccountName,
			},
			want: want{
				serviceAccount: types.ServiceAccount{Name: testutils.ServiceAccountName},
			},
		},
		"ShouldNotFindServiceAccountInNonExistingNamespace": {
			args: args{
				name:                       testutils.ServiceAccountName + testutils.NonExistentSuffix,
				existingServiceAccountName: "",
				namespace:                  namespaceName + testutils.NonExistentSuffix,
			},
			want: want{
				serviceAccount: types.ServiceAccount{},
				error:          fmt.Sprintf(ErrCouldNotGetServiceAccount, testutils.ServiceAccountName+testutils.NonExistentSuffix, namespaceName+testutils.NonExistentSuffix),
			},
		},
	}

	setup()
	createTestNamespace(namespaceName, utils.AddManagedLabel(map[string]string{}))

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			if test.args.existingServiceAccountName != "" {
				mocks.CreateTestServiceAccount(fakeClient, namespaceName, test.args.existingServiceAccountName, "")
			}
			c := mocks.GinContext()
			serviceAccountController := NewServiceAccountController(fakeClient, c, logger)
			response, err := serviceAccountController.GetServiceAccount(test.args.name, test.args.namespace)

			if test.want.error != "" {
				assert.ErrorContains(t, err, test.want.error)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, test.want.serviceAccount.Name, response.Name)
		})
	}
}

func TestGetServiceAccountToken(t *testing.T) {
	namespaceName := testutils.TestNamespace + "-getServiceAccount"
	type args struct {
		namespace string
		name      string
	}
	type want struct {
		tokenResponse types.TokenResponse
		error         string
	}
	cases := map[string]struct {
		args args
		want want
	}{
		"ShouldSucceedGettingServiceAccount": {
			args: args{
				namespace: namespaceName,
				name:      testutils.ServiceAccountName,
			},
			want: want{
				tokenResponse: types.TokenResponse{
					Token: "value",
				},
			},
		},
		"ShouldNotFindServiceAccountInNonExistingNamespace": {
			args: args{
				name:      testutils.CappName + testutils.NonExistentSuffix,
				namespace: namespaceName + testutils.NonExistentSuffix,
			},
			want: want{
				error: fmt.Sprintf(ErrCouldNotGetServiceAccount, testutils.CappName+testutils.NonExistentSuffix, namespaceName+testutils.NonExistentSuffix),
			},
		},
	}

	setup()
	createTestNamespace(namespaceName, utils.AddManagedLabel(map[string]string{}))
	mocks.CreateTestServiceAccountWithToken(fakeClient, namespaceName, testutils.ServiceAccountName, "token-secret", "value", "docker-cfg")

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			c := mocks.GinContext()
			serviceAccountController := NewServiceAccountController(fakeClient, c, logger)
			response, err := serviceAccountController.GetServiceAccountToken(test.args.name, test.args.namespace)

			if test.want.error != "" {
				assert.ErrorContains(t, err, test.want.error)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, test.want.tokenResponse, response)
		})
	}
}

func TestCreateServiceAccount(t *testing.T) {
	namespaceName := testutils.TestNamespace + "-createServiceAccount"
	type args struct {
		namespace                  string
		name                       string
		existingServiceAccountName string
	}
	type want struct {
		response types.ServiceAccount
		error    string
	}
	cases := map[string]struct {
		args args
		want want
	}{
		"ShouldSucceedCreatingServiceAccount": {
			args: args{
				namespace:                  namespaceName,
				name:                       testutils.ServiceAccountName,
				existingServiceAccountName: "",
			},
			want: want{
				response: types.ServiceAccount{Name: testutils.ServiceAccountName},
				error:    "",
			},
		},
		"ShouldHandleAlreadyExistingServiceAccount": {
			args: args{
				name:                       testutils.ServiceAccountName + "-new",
				namespace:                  namespaceName,
				existingServiceAccountName: testutils.ServiceAccountName + "-new",
			},
			want: want{
				response: types.ServiceAccount{},
				error:    fmt.Sprintf("serviceaccounts %q already exists", testutils.ServiceAccountName+"-new"),
			},
		},
	}

	setup()
	createTestNamespace(namespaceName, utils.AddManagedLabel(map[string]string{}))

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			if test.args.existingServiceAccountName != "" {
				mocks.CreateTestServiceAccount(fakeClient, namespaceName, test.args.existingServiceAccountName, "")
			}
			c := mocks.GinContext()
			serviceAccountController := NewServiceAccountController(fakeClient, c, logger)
			response, err := serviceAccountController.CreateServiceAccount(test.args.name, test.args.namespace)

			if test.want.error != "" {
				assert.ErrorContains(t, err, test.want.error)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, test.want.response, response)
		})
	}
}

func TestDeleteServiceAccount(t *testing.T) {
	namespaceName := testutils.TestNamespace + "-deleteServiceAccount"
	type args struct {
		namespace string
		name      string
	}
	type want struct {
		error string
	}
	cases := map[string]struct {
		args args
		want want
	}{
		"ShouldSucceedDeletingServiceAccount": {
			args: args{
				namespace: namespaceName,
				name:      testutils.ServiceAccountName,
			},
			want: want{
				error: "",
			},
		},
		"ShouldHandleNonExistingServiceAccount": {
			args: args{
				namespace: namespaceName,
				name:      testutils.ServiceAccountName + testutils.NonExistentSuffix,
			},
			want: want{
				error: fmt.Sprintf(ErrCouldNotDeleteServiceAccount, testutils.ServiceAccountName+testutils.NonExistentSuffix, namespaceName),
			},
		},
	}

	setup()
	createTestNamespace(namespaceName, utils.AddManagedLabel(map[string]string{}))
	mocks.CreateTestServiceAccount(fakeClient, namespaceName, testutils.ServiceAccountName, "")

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			c := mocks.GinContext()
			serviceAccountController := NewServiceAccountController(fakeClient, c, logger)
			err := serviceAccountController.DeleteServiceAccount(test.args.name, test.args.namespace)

			if test.want.error != "" {
				assert.ErrorContains(t, err, test.want.error)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
