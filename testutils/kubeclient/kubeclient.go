package kubeclient

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TODO: Maybe change the name of the below receiver? It's the same as the package name
func (kubeclient *KubeClient) GetAllPodsFromAllNamespaces() (*v1.PodList, error) {
	pods, err := kubeclient.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		// TODO: The caller will also try to tell a similar error, right? Then it will be a repeat? Hmm
		return nil, fmt.Errorf("error getting pods from all namespaces: %v", err)
	}

	return pods, nil
}
