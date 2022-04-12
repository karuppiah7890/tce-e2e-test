package kubeclient

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// With the below we can have sugar coat methods on KubeClient and also to the full power of
// kube client golang library as that field is exposed
type KubeClient struct {
	kubernetes.Interface
}

// GetKubeClient creates a Kubernetes client for a given kubeconfig path and kubeconfig context in the kubeconfig.
// When context is empty, the current context mentioned in the kubeconfig is used by the kube client
func GetKubeClient(kubeConfigPath string, context string) (*KubeClient, error) {
	config, err := configForContext(kubeConfigPath, context)
	if err != nil {
		return nil, err
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		// TODO: The caller will also try to tell a similar error, right? Then it will be a repeat? Hmm
		return nil, fmt.Errorf("could not get Kubernetes client: %v", err)
	}

	return &KubeClient{client}, nil
}

// configForContext creates a Kubernetes REST client configuration for a given kubeconfig path and kubeconfig context in the kubeconfig.
func configForContext(kubeConfigPath string, context string) (*rest.Config, error) {
	config, err := getConfig(kubeConfigPath, context).ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("could not get Kubernetes config for kubeconfig path %s and context %q: %v", kubeConfigPath, context, err)
	}
	return config, nil
}

// getConfig returns a Kubernetes client config for a given kubeconfig path and kubeconfig context in the kubeconfig.
func getConfig(kubeConfigPath string, context string) clientcmd.ClientConfig {
	rules := &clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeConfigPath}

	// TODO: Should we get rid of the cluster defaults? Maybe we don't need it
	overrides := &clientcmd.ConfigOverrides{ClusterDefaults: clientcmd.ClusterDefaults}

	if context != "" {
		overrides.CurrentContext = context
	}

	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, overrides)
}
