package tanzu

import (
	"fmt"
	"os"

	"github.com/karuppiah7890/tce-e2e-test/testutils/clirunner"
	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
)

// TODO: Should we rename this to UpdatePluginRepo?
// Update a repository configuration
// Runs `tanzu plugin repo update --gcp-bucket-name <gcp-bucket-name> <repo-name>`
func PluginRepoUpdate(repoName string, gcpBucketName string) error {
	exitCode, err := clirunner.Run(clirunner.Cmd{
		Name: "tanzu",
		Args: []string{
			"plugin",
			"repo",
			"update",
			"--gcp-bucket-name",
			gcpBucketName,
			repoName,
		},
		Env:    os.Environ(),
		Stdout: log.InfoWriter,
		Stderr: log.ErrorWriter,
	})
	if err != nil {
		return fmt.Errorf("error occurred while updating repository configuration of `%s` repository with `%s` GCP bucket. Exit code: %v. Error: %v", repoName, gcpBucketName, exitCode, err)
	}
	return nil
}

// TODO: Should we rename this to ListPlugins?
// TODO: Maybe in the future we can parse the data (with -o json JSON output) and return it
// List available plugins
// Runs `tanzu plugin list`
func PluginList() error {
	exitCode, err := clirunner.Run(clirunner.Cmd{
		Name: "tanzu",
		Args: []string{
			"plugin",
			"list",
		},
		Env:    os.Environ(),
		Stdout: log.InfoWriter,
		Stderr: log.ErrorWriter,
	})
	if err != nil {
		return fmt.Errorf("error occurred while listing available plugins. Exit code: %v. Error: %v", exitCode, err)
	}
	return nil
}

// TODO: Should we rename this to InstallPlugin?
// TODO: Should we have a better and smaller name for pathToLocalDiscoveryOrDistributionSource?
// Install a plugin
// Runs `tanzu plugin install --local <path-to-local-discovery-or-distribution-source> <plugin-name>`
func PluginInstall(pluginName string, pathToLocalDiscoveryOrDistributionSource string) error {
	exitCode, err := clirunner.Run(clirunner.Cmd{
		Name: "tanzu",
		Args: []string{
			"plugin",
			"install",
			"--local",
			pathToLocalDiscoveryOrDistributionSource,
			pluginName,
		},
		Env:    os.Environ(),
		Stdout: log.InfoWriter,
		Stderr: log.ErrorWriter,
	})
	if err != nil {
		return fmt.Errorf("error occurred while install `%s` plugin from `%s` local source. Exit code: %v. Error: %v", pluginName, pathToLocalDiscoveryOrDistributionSource, exitCode, err)
	}
	return nil
}
