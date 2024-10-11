package controllers

import (
	"context"
	"fmt"
	"github.com/dana-team/platform-backend/src/customerrors"
	"go.uber.org/zap"
	"io"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/dana-team/platform-backend/src/utils"
	"k8s.io/client-go/kubernetes"
)

const (
	errCouldNotOpenLogStream = "error opening log stream"
	errFetchingCappPods      = "error fetching Capp pods"
	errNoPodsFound           = "no pods found for Capp %q in namespace %q"
	errPodNotFound           = "pod %q not found for Capp %q in namespace %q"
)

// FetchPodLogs retrieves the logs of a specific container in a pod.
// It opens a log stream, reads the logs, and returns them as a string.
func FetchPodLogs(ctx context.Context, client kubernetes.Interface, namespace, podName, containerName string, previous bool, logger *zap.Logger) (io.ReadCloser, error) {
	logStream, err := utils.GetPodLogStream(ctx, client, namespace, podName, containerName, previous)
	if err != nil {
		logger.Error(fmt.Sprintf("%v: %v", errCouldNotOpenLogStream, err))
		return nil, customerrors.NewAPIError(errCouldNotOpenLogStream, err)
	}

	return logStream, nil
}

// FetchCappLogs retrieves the logs of a Capp's Knative service.
// It fetches the pods associated with the service, selects the first pod, and retrieves its logs.
func FetchCappLogs(ctx context.Context, client kubernetes.Interface, namespace, cappName, containerName, podName string, previous bool, logger *zap.Logger) (io.ReadCloser, error) {
	pods, err := utils.GetPodsByLabel(ctx, client, namespace, fmt.Sprintf(utils.ParentCappLabelSelector, cappName), metav1.ListOptions{})
	if err != nil {
		logger.Error(fmt.Sprintf("%v: %v", errFetchingCappPods, err))
		return nil, customerrors.NewAPIError(errFetchingCappPods, err)
	}

	if len(pods.Items) == 0 {
		logger.Error(fmt.Sprintf(errNoPodsFound, cappName, namespace))
		return nil, customerrors.NewAPIError(fmt.Sprintf(errNoPodsFound, cappName, namespace), err)
	}

	podName, ok := FetchCappPodName(podName, pods)
	if !ok {
		logger.Error(fmt.Sprintf(errPodNotFound, podName, cappName, namespace))
		return nil, customerrors.NewAPIError(fmt.Sprintf(errPodNotFound, podName, cappName, namespace), err)
	}

	if containerName == "" {
		containerName = cappName
	}

	return FetchPodLogs(ctx, client, namespace, podName, containerName, previous, logger)
}

// FetchCappPodName returns the validated pod name from the provided list of pods.
// If name is empty, it returns the name of the first pod in the list.
// It also returns a boolean indicating if the pod name was found in the list.
func FetchCappPodName(podName string, pods *corev1.PodList) (string, bool) {
	if podName == "" {
		return pods.Items[0].Name, true
	}

	if utils.IsPodInPodList(podName, pods) {
		return podName, true
	}

	return podName, false
}
