package utils

import (
	"context"
	"fmt"
	"io"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var (
	tailLines int64 = 10
)

// GetPodsByLabel returns the pods in a namespace using a given label selector.
func GetPodsByLabel(ctx context.Context, client kubernetes.Interface, namespace, labelSelector string, listOptions metav1.ListOptions) (*corev1.PodList, error) {
	listOptions.LabelSelector = labelSelector
	return client.CoreV1().Pods(namespace).List(ctx, listOptions)
}

// GetPodLogStream returns the logs of a container in a pod.
func GetPodLogStream(ctx context.Context, client kubernetes.Interface, namespace, podName, containerName string) (io.ReadCloser, error) {
	if containerName == "" {
		var err error
		containerName, err = getDefaultContainerName(ctx, client, namespace, podName)
		if err != nil {
			return nil, err
		}
	}

	ok, err := isContainerInPod(ctx, client, namespace, podName, containerName)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, fmt.Errorf("container %q not found in the pod %q", containerName, podName)
	}

	req := client.CoreV1().Pods(namespace).GetLogs(podName, &corev1.PodLogOptions{
		Container: containerName,
		Follow:    true,
		TailLines: &tailLines,
	})

	return req.Stream(ctx)
}

// isContainerInPod checks if a container with the given name exists in the specified pod.
func isContainerInPod(ctx context.Context, client kubernetes.Interface, namespace, podName, containerName string) (bool, error) {
	pod, err := client.CoreV1().Pods(namespace).Get(ctx, podName, metav1.GetOptions{})
	if err != nil {
		return false, fmt.Errorf("failed to get pod: %v", err)
	}

	for _, container := range pod.Spec.Containers {
		if container.Name == containerName {
			return true, nil
		}
	}

	return false, nil
}

// IsPodInPodList checks if a pod with the given name exists in the provided list of pods.
func IsPodInPodList(podName string, pods *corev1.PodList) bool {
	for _, pod := range pods.Items {
		if pod.Name == podName {
			return true
		}
	}

	return false
}

// getDefaultContainerName returns the name of the only container in the pod if there is exactly one container, otherwise returns an error.
func getDefaultContainerName(ctx context.Context, client kubernetes.Interface, namespace, podName string) (string, error) {
	pod, err := client.CoreV1().Pods(namespace).Get(ctx, podName, metav1.GetOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to get pod: %v", err)
	}

	if len(pod.Spec.Containers) == 1 {
		return pod.Spec.Containers[0].Name, nil
	}

	return "", fmt.Errorf("pod %q has multiple containers, please specify the container name", podName)
}
