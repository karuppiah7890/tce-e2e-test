package tanzu

import (
	"fmt"
	"os"

	"github.com/karuppiah7890/tce-e2e-test/testutils/clirunner"
	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
)

func CollectManagementClusterDiagnostics(managementClusterName string) error {
	log.Infof("Collecting diagnostics of `%s` management cluster", managementClusterName)
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
		return fmt.Errorf("error occurred while collecting diagnostics of `%s` management cluster. exit code: %v. error: %v", managementClusterName, exitCode, err)
	}
	return nil
}

// TODO: Convert workload cluster infra from string to a type - say iota or similar to get pre-defined (compile time) constants like azure, aws, vsphere, docker
func CollectManagementClusterAndWorkloadClusterDiagnostics(managementClusterName string, workloadClusterName string, workloadClusterInfra string) error {
	log.Infof("Collecting diagnostics of `%s` management cluster and `%s` workload cluster (in `%s` infra)", managementClusterName, workloadClusterName, workloadClusterInfra)
	// Run the command
	// `tanzu diagnostics collect --bootstrap-cluster-skip \
	//         --management-cluster-name <management-cluster-name> \
	//         --workload-cluster-infra <workload-cluster-infra> \
	//         --workload-cluster-name <workload-cluster-name>`

	exitCode, err := clirunner.Run(clirunner.Cmd{
		Name: "tanzu",
		Args: []string{
			"diagnostics",
			"collect",
			"--bootstrap-cluster-skip",
			"--management-cluster-name",
			managementClusterName,
			"--workload-cluster-name",
			workloadClusterName,
			"--workload-cluster-infra",
			workloadClusterInfra,
		},
		Env:    os.Environ(),
		Stdout: log.InfoWriter,
		Stderr: log.ErrorWriter,
	})
	if err != nil {
		return fmt.Errorf("error occurred while collecting diagnostics of `%s` management cluster and `%s` workload cluster (in `%s` infra). exit code: %v. error: %v",
			managementClusterName, workloadClusterName, workloadClusterInfra, exitCode, err)
	}
	return nil
}
