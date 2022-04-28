package tanzu

import (
	"fmt"
	"strings"

	tf "github.com/vmware-tanzu/tanzu-framework/apis/config/v1alpha1"
)

func GetBootstrapClusterDockerContainerNameForManagementCluster(managementClusterName string) (string, error) {
	clientConfig, err := GetClientConfig()
	if err != nil {
		return "", fmt.Errorf("error getting tanzu client config: %v", err)
	}
	for _, server := range clientConfig.KnownServers {
		if server.Type == tf.ManagementClusterServerType && server.Name == managementClusterName {
			if server.ManagementClusterOpts == nil {
				return "", fmt.Errorf("managementClusterOpts for management cluster %s is empty", managementClusterName)
			}

			bootstrapClusterContext := server.ManagementClusterOpts.Context

			bootstrapClusterUniqueSuffix := strings.Replace(bootstrapClusterContext, "kind-tkg-kind-", "", 1)

			tanzuBootstrapClusterDockerContainerName := fmt.Sprintf("tkg-kind-%s-control-plane", bootstrapClusterUniqueSuffix)

			return tanzuBootstrapClusterDockerContainerName, nil
		}
	}
	return "", fmt.Errorf("could not find the management cluster %s in the tanzu config yaml", managementClusterName)
}
