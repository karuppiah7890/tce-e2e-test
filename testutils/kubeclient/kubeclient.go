package kubeclient

import (
	"context"
	"fmt"

	"github.com/karuppiah7890/tce-e2e-test/testutils/clirunner"
	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TODO: Should we return []Pod instead of *v1.PodList ? Something to think about
// TODO: Maybe change the name of the below receiver? It's the same as the package name
func (kubeclient *KubeClient) GetAllPodsFromAllNamespaces() (*v1.PodList, error) {
	pods, err := kubeclient.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		// TODO: The caller will also try to tell a similar error, right? Then it will be a repeat? Hmm
		return nil, fmt.Errorf("error getting all pods from all namespaces: %v", err)
	}

	return pods, nil
}

// TODO: Should we return []Node instead of *v1.PodList ? Something to think about
func (kubeclient *KubeClient) GetAllNodes() (*v1.NodeList, error) {
	nodes, err := kubeclient.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		// TODO: The caller will also try to tell a similar error, right? Then it will be a repeat? Hmm
		return nil, fmt.Errorf("error getting all nodes: %v", err)
	}

	return nodes, nil
}

func UseKubeConfigContext(workloadClusterKubeContext string) error {
	exitCode, err := clirunner.Run(clirunner.Cmd{
		Name: "kubectl",
		Args: []string{
			"config",
			"use-context",
			workloadClusterKubeContext,
		},
		Stdout: log.InfoWriter,
		// TODO: Should we log standard errors as errors in the log? Because tanzu prints other information also
		// to standard error, which are kind of like information, apart from actual errors, so showing
		// everything as error is misleading. Gotta think what to do about this. The main problem is
		// console has only standard output and standard error, and tanzu is using standard output only for
		// giving output for things like --dry-run when it needs to print yaml content, but everything else
		// is printed to standard error
		Stderr: log.ErrorWriter,
	})
	if err != nil {
		return fmt.Errorf("error occurred while setting cluster kube config context. exit code: %v. error: %v", exitCode, err)
	}
	return nil
}
