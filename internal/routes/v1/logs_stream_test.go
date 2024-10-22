package v1

import (
	"fmt"
	"github.com/dana-team/platform-backend/internal/middleware"
	"github.com/dana-team/platform-backend/internal/utils/testutils"
	"github.com/dana-team/platform-backend/internal/utils/testutils/mocks"
	"github.com/gorilla/websocket"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testNamespaceGetCappLogs = testutils.TestNamespace + "-Test_GetCappLogs"
	testNamespaceGetPodLogs  = testutils.TestNamespace + "-Test_GetPodLogs"

	authorizationHeader = "Authorization"
	pod1                = testutils.PodName + "-1"
	pod2                = testutils.PodName + "-2"
	pod3                = testutils.PodName + "-3"
)

func Test_GetCappLogs(t *testing.T) {
	type args struct {
		token string
		wsUrl string
	}
	type want struct {
		statusCode    int
		expectedLines []string
	}

	cases := map[string]struct {
		args args
		want want
	}{
		"ShouldStreamLogsWithoutQueryParams": {
			args: args{
				token: "valid_token",
				wsUrl: fmt.Sprintf("/v1/namespaces/%s/capps/%s/logs", testNamespaceGetCappLogs, testutils.CappName),
			},
			want: want{
				statusCode:    http.StatusSwitchingProtocols,
				expectedLines: []string{fmt.Sprintf("Capp: %q line: fake logs", testutils.CappName)},
			},
		},
		"ShouldNotStreamLogsWithInvalidCappName": {
			args: args{
				token: "valid_token",
				wsUrl: fmt.Sprintf("/v1/namespaces/%s/capps/invalid-capp/logs", testNamespaceGetCappLogs),
			},
			want: want{
				statusCode:    http.StatusNotFound,
				expectedLines: []string{fmt.Sprintf("error: Error streaming %q logs: no pods found for Capp %q in namespace %q", "Capp", "invalid-capp", testNamespaceGetCappLogs)},
			},
		},
		"ShouldStreamLogsWithQueryParams": {
			args: args{
				token: "valid_token",
				wsUrl: fmt.Sprintf("/v1/namespaces/%s/capps/%s/logs?podName=test-pod-2", testNamespaceGetCappLogs, testutils.CappName),
			},
			want: want{
				statusCode:    http.StatusSwitchingProtocols,
				expectedLines: []string{fmt.Sprintf("Capp: %q line: fake logs", testutils.CappName)},
			},
		},
		"ShouldStreamLogsWithPreviousQueryParam": {
			args: args{
				token: "valid_token",
				wsUrl: fmt.Sprintf("/v1/namespaces/%s/capps/%s/logs?previous=true", testNamespaceGetCappLogs, testutils.CappName),
			},
			want: want{
				statusCode:    http.StatusSwitchingProtocols,
				expectedLines: []string{fmt.Sprintf("Capp: %q line: fake logs", testutils.CappName)},
			},
		},
		"ShouldNotStreamLogsWithNonExistingPodName": {
			args: args{
				token: "valid_token",
				wsUrl: fmt.Sprintf("/v1/namespaces/%s/capps/%s/logs?podName=pod%s", testNamespaceGetCappLogs, testutils.CappName, testutils.NonExistentSuffix),
			},
			want: want{
				statusCode:    http.StatusSwitchingProtocols,
				expectedLines: []string{fmt.Sprintf("error: Error streaming %q logs: pod %q not found for Capp %q in namespace %q", "Capp", "pod"+testutils.NonExistentSuffix, testutils.CappName, testNamespaceGetCappLogs)},
			},
		},
		"ShouldNotStreamLogsWithInvalidContainerName": {
			args: args{
				token: "valid_token",
				wsUrl: fmt.Sprintf("/v1/namespaces/%s/capps/%s/logs?containerName=container%s", testNamespaceGetCappLogs, testutils.CappName, testutils.NonExistentSuffix),
			},
			want: want{
				statusCode:    http.StatusSwitchingProtocols,
				expectedLines: []string{fmt.Sprintf("error: Error streaming %q logs: error opening log stream, container %q not found in the pod %q", "Capp", "container"+testutils.NonExistentSuffix, pod2)},
			},
		},
		"ShouldStreamLogsWithValidContainerName": {
			args: args{
				token: "valid_token",
				wsUrl: fmt.Sprintf("/v1/namespaces/%s/capps/%s/logs?containerName=test-container", testNamespaceGetCappLogs, testutils.CappName),
			},
			want: want{
				statusCode:    http.StatusSwitchingProtocols,
				expectedLines: []string{fmt.Sprintf("Capp: %q line: fake logs", testutils.CappName)},
			},
		},
	}

	setup()
	mocks.CreateTestNamespace(fakeClient, testNamespaceGetCappLogs)
	mocks.CreateTestCapp(dynClient, testutils.CappName, testNamespaceGetCappLogs, testutils.Domain, testutils.SiteName, map[string]string{}, nil)
	mocks.CreateTestPod(fakeClient, testNamespaceGetCappLogs, pod1, "", false)
	mocks.CreateTestPod(fakeClient, testNamespaceGetCappLogs, pod2, testutils.CappName, true)

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			token = tc.args.token
			server := httptest.NewServer(router)
			defer server.Close()

			wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + tc.args.wsUrl

			dialer := websocket.DefaultDialer
			headers := http.Header{}
			headers.Add(authorizationHeader, tc.args.token)
			headers.Add(middleware.WebsocketTokenHeader, tc.args.token)

			conn, resp, err := dialer.Dial(wsURL, headers)
			assert.Equal(t, tc.want.statusCode, resp.StatusCode)
			if tc.want.statusCode == http.StatusUnauthorized || tc.want.statusCode == http.StatusNotFound {
				return
			}

			if err != nil {
				t.Fatalf("Failed to dial WebSocket: %v", err)
			}

			defer conn.Close()
			for _, expectedLine := range tc.want.expectedLines {
				_, message, err := conn.ReadMessage()
				if err != nil {
					t.Fatalf("Error reading message from WebSocket: %v", err)
				}
				assert.Contains(t, string(message), expectedLine)
			}
		})
	}
}

func Test_GetPodLogs(t *testing.T) {
	type args struct {
		token string
		wsUrl string
	}
	type want struct {
		statusCode    int
		expectedLines []string
	}

	cases := map[string]struct {
		args args
		want want
	}{
		"ShouldStreamLogsWithoutQueryParams": {
			args: args{
				token: "valid_token",
				wsUrl: fmt.Sprintf("/v1/namespaces/%s/pods/%s/logs", testNamespaceGetPodLogs, pod1),
			},
			want: want{
				statusCode:    http.StatusSwitchingProtocols,
				expectedLines: []string{fmt.Sprintf("Pod: %q line: fake logs", pod1)},
			},
		},
		"ShouldStreamLogsWithQueryParams": {
			args: args{
				token: "valid_token",
				wsUrl: fmt.Sprintf("/v1/namespaces/%s/pods/%s/logs?container=test-container", testNamespaceGetPodLogs, pod1),
			},
			want: want{
				statusCode:    http.StatusSwitchingProtocols,
				expectedLines: []string{fmt.Sprintf("Pod: %q line: fake logs", pod1)},
			},
		},
		"ShouldStreamLogsWithPreviousQueryParam": {
			args: args{
				token: "valid_token",
				wsUrl: fmt.Sprintf("/v1/namespaces/%s/pods/%s/logs?previous=true", testNamespaceGetPodLogs, pod1),
			},
			want: want{
				statusCode:    http.StatusSwitchingProtocols,
				expectedLines: []string{fmt.Sprintf("Pod: %q line: fake logs", pod1)},
			},
		},
		"ShouldStreamLogsWithoutQueryParamsMultipleContainers": {
			args: args{
				token: "valid_token",
				wsUrl: fmt.Sprintf("/v1/namespaces/%s/pods/%s/logs", testNamespaceGetPodLogs, pod3),
			},
			want: want{
				statusCode:    http.StatusSwitchingProtocols,
				expectedLines: []string{fmt.Sprintf(`error: Error streaming "Pod" logs: error opening log stream, pod %q has multiple containers, please specify the container name`, pod3)},
			},
		},
		"ShouldNotStreamLogsWithNonExistingPodName": {
			args: args{
				token: "valid_token",
				wsUrl: fmt.Sprintf("/v1/namespaces/%s/pods/test-invalid-pod/logs", testNamespaceGetPodLogs),
			},
			want: want{
				statusCode:    http.StatusSwitchingProtocols,
				expectedLines: []string{fmt.Sprintf(`error: Error streaming "Pod" logs: error opening log stream, failed to get pod: pods %q not found`, "test-invalid-pod")},
			},
		},
		"ShouldNotStreamLogsWithInvalidContainerName": {
			args: args{
				token: "valid_token",
				wsUrl: fmt.Sprintf("/v1/namespaces/%s/pods/%s/logs?containerName=container%s", testNamespaceGetPodLogs, pod1, testutils.NonExistentSuffix),
			},
			want: want{
				statusCode:    http.StatusSwitchingProtocols,
				expectedLines: []string{fmt.Sprintf(`error: Error streaming "Pod" logs: error opening log stream, container %q not found in the pod %q`, "container"+testutils.NonExistentSuffix, pod1)},
			},
		},
	}

	setup()
	mocks.CreateTestNamespace(fakeClient, testNamespaceGetPodLogs)

	mocks.CreateTestPod(fakeClient, testNamespaceGetPodLogs, pod1, "", false)
	mocks.CreateTestPod(fakeClient, testNamespaceGetPodLogs, pod2, testutils.CappName, false)
	mocks.CreateTestPod(fakeClient, testNamespaceGetPodLogs, pod3, testutils.CappName, true)

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			token = tc.args.token
			server := httptest.NewServer(router)
			defer server.Close()

			wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + tc.args.wsUrl

			dialer := websocket.DefaultDialer
			headers := http.Header{}
			headers.Add(authorizationHeader, tc.args.token)
			headers.Add(middleware.WebsocketTokenHeader, tc.args.token)

			conn, resp, err := dialer.Dial(wsURL, headers)
			assert.Equal(t, tc.want.statusCode, resp.StatusCode)
			if tc.want.statusCode == http.StatusUnauthorized {
				return
			}

			if err != nil {
				t.Fatalf("Failed to dial WebSocket: %v", err)
			}

			defer conn.Close()
			for _, expectedLine := range tc.want.expectedLines {
				_, message, err := conn.ReadMessage()
				if err != nil {
					t.Fatalf("Error reading message from WebSocket: %v", err)
				}
				assert.Contains(t, string(message), expectedLine)
			}
		})
	}
}
