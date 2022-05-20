package tanzu

import (
	"fmt"

	"github.com/karuppiah7890/tce-e2e-test/testutils/clirunner"
	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
)

func PrintTanzuVersion() error {
	cmd := clirunner.Cmd{
		Name: "tanzu",
		Args: []string{
			"version",
		},
		Stdout: log.InfoWriter,
		Stderr: log.ErrorWriter,
	}

	exitCode, err := clirunner.Run(cmd)
	if err != nil {
		return fmt.Errorf("error occurred while printing tanzu CLI version. Exit code: %v. Error: %v", exitCode, err)
	}

	return nil
}

// Where to get the list of plugins from?
// 1. tanzu plugin list . This also has json and yaml output with `-o` flag
// 2. check artifact (tar ball, zip) for K8s resource yaml files of kind cli.tanzu.vmware.com/v1alpha1/CLIPlugin inside the
// directory default-local/discovery/standalone in the artifact
// 3. Manually put the list - not feasible as list can change for different version
func PrintAllPluginVersions() {
	// TODO: Implement this
}
