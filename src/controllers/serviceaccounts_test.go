package controllers

import (
	"fmt"
	"github.com/dana-team/platform-backend/src/types"
	"github.com/dana-team/platform-backend/src/utils"
	"github.com/dana-team/platform-backend/src/utils/testutils"
	"github.com/dana-team/platform-backend/src/utils/testutils/mocks"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestGetServiceAccount(t *testing.T) {
	namespaceName := testutils.TestNamespace + "-getServiceAccount"
	type args struct {
		namespace string
		name      string
	}
	type want struct {
		serviceAccount *corev1.ServiceAccount
		error          string
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
				serviceAccount: &corev1.ServiceAccount{
					ObjectMeta: metav1.ObjectMeta{
						Name:      testutils.ServiceAccountName,
						Namespace: namespaceName,
					},
					Secrets: []corev1.ObjectReference{
						{
							Kind:      testutils.Secret,
							Name:      "docker-cfg",
							Namespace: namespaceName,
						},
					},
				},
			},
		},
		"ShouldNotFindServiceAccountInNonExistingNamespace": {
			args: args{
				name:      testutils.CappName + testutils.NonExistentSuffix,
				namespace: namespaceName + testutils.NonExistentSuffix,
			},
			want: want{
				serviceAccount: &corev1.ServiceAccount{},
				error:          fmt.Sprintf(ErrCouldNotGetServiceAccount, testutils.CappName+testutils.NonExistentSuffix, namespaceName+testutils.NonExistentSuffix),
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
			response, err := serviceAccountController.GetServiceAccount(test.args.name, test.args.namespace)

			if test.want.error != "" {
				assert.ErrorContains(t, err, test.want.error)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, test.want.serviceAccount, response)
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
