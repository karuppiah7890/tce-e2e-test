package tanzu

import (
	"fmt"
	"os"

	// TODO: Better name for this package?
	tf "github.com/vmware-tanzu/tanzu-framework/apis/config/v1alpha1"
	// TODO: Better name for this package?
	tfconfig "github.com/vmware-tanzu/tanzu-framework/pkg/v1/config"

	"github.com/karuppiah7890/tce-e2e-test/testutils/clirunner"
	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
)

// TODO: Consider using Tanzu golang client library instead of running tanzu as a CLI.
// We could invoke plugins using their names and have tight integration. It comes with it's
// own pros and cons. Pro - tight and smooth integration with Tanzu Framework.
// Con - same as Pro - tight and smooth integration with Tanzu Framework - why? because Tanzu
// Framework does not provide any guarantee for API support, also it's in 0.x.y series which
// means they can break a lot of things which can break already generally fragile E2E tests
// more easily and more often. Also, if we import Tanzu Framework as a library, to test different
// versions of Tanzu Framework, we have import different versions of it, unlike CLI where we can
// just install the appropriate CLI version before testing it. For example, test v0.11.4 TF that
// TCE currently uses and also test v0.20.0 TF which is the latest version of TF. Of course it's not
// easy to concurrently / simultaneously test both versions, at least not in CLI, and with library, idk,
// it might be possible and easy? not sure for now, gotta experiment. We can also consider dynamically linked
// libraries and similar concept, we currently instead have tanzu CLI tool which is dynamically invoked and linked
// to this test program

// TODO: Maybe create a wrapper function called Tanzu() around clirunner.Run?

type TanzuConfig map[string]string

type EnvVars []string

// TODO: Move this to a tanzu specific lib
func TanzuAzureConfig(clusterName string) TanzuConfig {
	// TODO: Ideas:
	// We could also represent this config in a test data yaml file,
	// but config as code is more powerful - we can do more over here
	// for example, run tests for multiple plans - dev and prod very
	// easily instead of duplicating the whole config yaml file just to
	// run same test with different plan. We can then easily run many tests
	// with different set of config values by defining the range / possible set
	// of test values for each config.
	// Some configs that can be changed -
	// 1. Infra provider can change if it's AWS, vSphere etc
	// but this function is named as tanzuAzureConfig so it's okay.
	// 2. Cluster plan - dev and prod
	// 3. Azure location - the whole big list of azure locations
	// 4. Azure control plane machine type - the whole big list of azure VM machine types. Note: https://github.com/vmware-tanzu/community-edition/issues/1749. Also note, we might need VMs of some minimum size for cluster creation to work
	// 5. Azure worker node machine type - the whole big list of azure VM machine types. Note: https://github.com/vmware-tanzu/community-edition/issues/1749. Also note, we might need VMs of some minimum size for cluster creation to work
	// 6. OS_ARCH - amd64, arm64 . There's talks around ARM support at different levels now
	// 7. OS_VERSION - 20.04 or 18.04 as of now
	// 8. AZURE_VNET_CIDR - any CIDR range
	// 9. AZURE_CONTROL_PLANE_SUBNET_CIDR - any CIDR range
	// 10. AZURE_NODE_SUBNET_CIDR - any CIDR range
	// 11. CLUSTER_CIDR - any CIDR range
	// 12. CLUSTER_CIDR - any CIDR range
	// 13. ENABLE_CEIP_PARTICIPATION - true or false
	// 14. ENABLE_MHC - true or false. In one issue, someone had to set this to false or else their cluster creation was failing
	// 15. IDENTITY_MANAGEMENT_TYPE - none or some particular set of identity management types

	// TODO: In our bash E2E test, we control the value of the below env vars using the cluster name along with some suffix
	// AZURE_RESOURCE_GROUP
	// AZURE_VNET_RESOURCE_GROUP
	// AZURE_VNET_NAME
	// AZURE_CONTROL_PLANE_SUBNET_NAME
	// AZURE_NODE_SUBNET_NAME
	return TanzuConfig{
		"CLUSTER_NAME":                     clusterName,
		"INFRASTRUCTURE_PROVIDER":          "azure",
		"CLUSTER_PLAN":                     "dev",
		"AZURE_LOCATION":                   "australiaeast",
		"AZURE_CONTROL_PLANE_MACHINE_TYPE": "Standard_D4s_v3",
		"AZURE_NODE_MACHINE_TYPE":          "Standard_D4s_v3",
		"OS_ARCH":                          "amd64",
		"OS_NAME":                          "ubuntu",
		"OS_VERSION":                       "20.04",
		"AZURE_VNET_CIDR":                  "10.0.0.0/16",
		"AZURE_CONTROL_PLANE_SUBNET_CIDR":  "10.0.0.0/24",
		"AZURE_NODE_SUBNET_CIDR":           "10.0.1.0/24",
		"CLUSTER_CIDR":                     "100.96.0.0/11",
		"SERVICE_CIDR":                     "100.64.0.0/13",
		"ENABLE_CEIP_PARTICIPATION":        "false",
		"ENABLE_MHC":                       "true",
		"IDENTITY_MANAGEMENT_TYPE":         "none",
	}
}

func TanzuConfigToEnvVars(tanzuConfig TanzuConfig) EnvVars {
	envVars := make(EnvVars, 0, len(tanzuConfig))

	for key, value := range tanzuConfig {
		envVars = append(envVars, fmt.Sprintf("%s=%s", key, value))
	}

	return envVars
}

// Runs `tanzu config set features.global.context-aware-cli-for-plugins false` command
func DisableContextAwareCliForPluginsGlobally() error {
	exitCode, err := clirunner.Run(clirunner.Cmd{
		Name: "tanzu",
		Args: []string{
			"config",
			"set",
			"features.global.context-aware-cli-for-plugins",
			"false",
		},
		Env:    os.Environ(),
		Stdout: log.InfoWriter,
		Stderr: log.ErrorWriter,
	})
	if err != nil {
		return fmt.Errorf("error occurred while disabling context aware cli for plugins globally. Exit code: %v. Error: %v", exitCode, err)
	}
	return nil
}

func GetClientConfig() (*tf.ClientConfig, error) {
	return tfconfig.GetClientConfig()
}
