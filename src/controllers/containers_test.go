package controllers

import (
	"context"
	"fmt"
	"github.com/dana-team/platform-backend/src/types"
	"github.com/dana-team/platform-backend/src/utils"
	"github.com/dana-team/platform-backend/src/utils/testutils"
	"github.com/dana-team/platform-backend/src/utils/testutils/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetContainers(t *testing.T) {
	namespaceName := testutils.TestNamespace + "-getcontainers"
	type args struct {
		namespace string
		podName   string
	}
	type want struct {
		response types.GetContainersResponse
		error    string
	}
	cases := map[string]struct {
		args args
		want want
	}{
		"ShouldSucceedGettingAllContainers": {
			args: args{
				namespace: namespaceName,
				podName:   pod2,
			},
			want: want{
				response: types.GetContainersResponse{
					Count: 2,
					Containers: []types.Container{
						{ContainerName: testutils.TestContainerName},
						{ContainerName: testutils.CappName},
					},
				},
			},
		},
		"ShouldNotFindContainersInNonExistingNamespaces": {
			args: args{
				podName:   pod2,
				namespace: namespaceName + testutils.NonExistentSuffix,
			},
			want: want{
				response: types.GetContainersResponse{},
				error:    fmt.Sprintf("failed to get pod %q, in the namespace %q with error: %v", pod2, namespaceName+testutils.NonExistentSuffix, fmt.Sprintf(`pods %q not found`, pod2)),
			},
		},
	}

	setup()
	containerController := NewContainerController(fakeClient, context.TODO(), logger)
	createTestNamespace(namespaceName, utils.AddManagedLabel(map[string]string{}))
	mocks.CreateTestPod(fakeClient, namespaceName, pod1, "", false)
	mocks.CreateTestPod(fakeClient, namespaceName, pod2, testutils.CappName, true)

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			response, err := containerController.GetContainers(test.args.namespace, test.args.podName)
			if test.want.error != "" {
				assert.ErrorContains(t, err, test.want.error)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, test.want.response, response)
		})
	}
}
