package kubeclient

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
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
		return nil, fmt.Errorf("could not get Kubernetes config for kubeconfig path %s and context %s: %v", kubeConfigPath, context, err)
	}
	return config, nil
}

// new package / file for kubeconfig stuff?
// DeleteContext deletes the context from the kubeconfig file and also deletes the corresponding
// cluster and user from the kubeconfig file
func DeleteContext(kubeConfigPath string, contextName string) error {
	rawConfig, err := getRawConfig(kubeConfigPath)
	if err != nil {
		return fmt.Errorf("could not get raw kubernetes config for kubeconfig path %s: %v", kubeConfigPath, err)
	}

	// modify raw config to delete the context with the name contextName.
	context, ok := rawConfig.Contexts[contextName]

	if !ok {
		return fmt.Errorf("could not find context named %s in kubeconfig file at path %s", contextName, kubeConfigPath)
	}

	clusterName := context.Cluster
	authInfoName := context.AuthInfo

	delete(rawConfig.Contexts, contextName)
	delete(rawConfig.Clusters, clusterName)
	delete(rawConfig.AuthInfos, authInfoName)

	configAccess := getConfigAccess(kubeConfigPath)

	err = clientcmd.ModifyConfig(configAccess, rawConfig, true)

	if err != nil {
		return fmt.Errorf("error in modifying the kubeconfig file at path %s: %v", kubeConfigPath, err)
	}

	return nil
}

// getRawConfig creates a raw Kubernetes configuration for a given kubeconfig path
func getRawConfig(kubeConfigPath string) (clientcmdapi.Config, error) {
	config, err := getConfig(kubeConfigPath, "").RawConfig()
	if err != nil {
		return clientcmdapi.Config{}, fmt.Errorf("could not get raw Kubernetes config for kubeconfig path %s: %v", kubeConfigPath, err)
	}
	return config, nil
}

// getConfigAccess creates a raw Kubernetes configuration for a given kubeconfig path
func getConfigAccess(kubeConfigPath string) clientcmd.ConfigAccess {
	return getConfig(kubeConfigPath, "").ConfigAccess()
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
