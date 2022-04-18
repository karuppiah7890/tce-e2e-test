package tanzu

import (
	"fmt"
	"os"

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
