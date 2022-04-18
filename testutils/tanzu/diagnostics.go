package tanzu

import (
	"fmt"
	"os"

	"github.com/karuppiah7890/tce-e2e-test/testutils/clirunner"
	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
)

func CollectManagementClusterDiagnostics(managementClusterName string) error {
	log.Infof("Collecting diagnostics `%s` management cluster", managementClusterName)
	// Run `tanzu diagnostics collect --management-cluster-name <management-cluster-name>`

	exitCode, err := clirunner.Run(clirunner.Cmd{
		Name: "tanzu",
		Args: []string{
			"diagnostics",
			"collect",
			"--management-cluster-name",
			managementClusterName,
		},
		Env:    os.Environ(),
		Stdout: log.InfoWriter,
		Stderr: log.ErrorWriter,
	})
	if err != nil {
		return fmt.Errorf("error occurred while collecting diagnostics of management cluster. exit code: %v. error: %v", exitCode, err)
	}
	return nil
}
