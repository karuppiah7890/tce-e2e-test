package tce

import (
	"fmt"
	"os"

	"github.com/karuppiah7890/tce-e2e-test/testutils/clirunner"
	"github.com/karuppiah7890/tce-e2e-test/testutils/kubeclient"
	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
)

type Package struct {
	Name    string
	Version string
}

func PackageE2Etest(packageDetails Package, workloadClusterKubeContext string) error {
	err := kubeclient.UseKubeConfigContext(workloadClusterKubeContext)
	if err != nil {
		return fmt.Errorf("error occurred while using the workload cluster context. error: %v", err)
	}
	exitCode, err := clirunner.Run(clirunner.Cmd{
		Name: "tanzu",
		Args: []string{
			"package",
			"repository",
			"add",
			"tce-repo",
			"--url",
			"projects.registry.vmware.com/tce/main:0.12.0",
			"--namespace",
			"tanzu-package-repo-global",
		},
		Stdout: log.InfoWriter,
		// TODO: Should we log standard errors as errors in the log? Because tanzu prints other information also
		// to standard error, which are kind of like information, apart from actual errors, so showing
		// everything as error is misleading. Gotta think what to do about this. The main problem is
		// console has only standard output and standard error, and tanzu is using standard output only for
		// giving output for things like --dry-run when it needs to print yaml content, but everything else
		// is printed to standard error
		Stderr: log.ErrorWriter,
	})
	if err != nil {
		return fmt.Errorf("error occurred while adding package repo. exit code: %v. error: %v", exitCode, err)
	}

	err = os.Chdir("community-edition/addons/packages/" + packageDetails.Name + "/" + packageDetails.Version + "/test")
	if err != nil {
		return fmt.Errorf("error while changing directory to community-edition: %v", err)
	}
	exitCode, err = clirunner.Run(clirunner.Cmd{
		Name: "make",
		Args: []string{
			"e2e-test",
		},
		Stdout: log.InfoWriter,
		// TODO: Should we log standard errors as errors in the log? Because tanzu prints other information also
		// to standard error, which are kind of like information, apart from actual errors, so showing
		// everything as error is misleading. Gotta think what to do about this. The main problem is
		// console has only standard output and standard error, and tanzu is using standard output only for
		// giving output for things like --dry-run when it needs to print yaml content, but everything else
		// is printed to standard error
		Stderr: log.ErrorWriter,
	})

	if err != nil {
		return fmt.Errorf("error occurred while E2E test for %v. Exit code: %v. Error: %v", packageDetails.Name, exitCode, err)
	}
	return nil
}
