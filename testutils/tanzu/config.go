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
