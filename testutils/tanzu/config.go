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

//TODO: Should we stick to env vars for cluster config or can we use yaml like tanzu cli consumes
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
