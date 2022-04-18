// TODO: Should we rename the package to the full form like tanzuframework?
package tf

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/karuppiah7890/tce-e2e-test/testutils/download"
	"github.com/karuppiah7890/tce-e2e-test/testutils/extract"
	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
	"github.com/karuppiah7890/tce-e2e-test/testutils/platforms"
	"github.com/karuppiah7890/tce-e2e-test/testutils/search"
	"github.com/karuppiah7890/tce-e2e-test/testutils/tanzu"
)

// TODO: Should we support getTfArtifactUrl("v0.20.0") too? Or just have one of them? Which one?
// Example: getTfArtifactUrl("0.20.0")
func getTfArtifactUrl(version string) (string, error) {
	log.Infof("Getting TF artifact URL for version %s", version)

	// TODO: Convert magic strings like "linux/amd64", "darwin/amd64", "windows/amd64" to constants
	supportedPlatforms := []string{"linux/amd64", "darwin/amd64", "windows/amd64"}

	// TODO: Maybe merge the supported OSes and artifact extension data as artifact extension should be
	// present for each supported operating system in this case and currently there's a duplication of
	// data here
	artifactExtensions := map[string]string{platforms.LINUX: extract.TARGZ, platforms.DARWIN: extract.TARGZ, platforms.WINDOWS: extract.ZIP}

	architecture := runtime.GOARCH
	operatingSystem := runtime.GOOS
	platform := fmt.Sprintf("%s/%s", operatingSystem, architecture)
	if !search.IsPresentIn(platform, supportedPlatforms) {
		return "", fmt.Errorf("platform %s is not supported by TF. Supported platforms: %v", platform, supportedPlatforms)
	}

	artifactExtension := artifactExtensions[operatingSystem]
	// Example: https://github.com/vmware-tanzu/tanzu-framework/releases/download/v0.20.0/tanzu-framework-linux-amd64.tar.gz
	return fmt.Sprintf("https://github.com/vmware-tanzu/tanzu-framework/releases/download/v%s/tanzu-framework-%s-%s.%s", version, operatingSystem, architecture, artifactExtension), nil
}

// TODO: Support latest installation method
// where API driven plugin discovery is activated which is the default in latest versions of TF
// This is based on option 1 here https://github.com/vmware-tanzu/tanzu-framework/blob/main/docs/cli/getting-started.md#option-1-manual-download-cli-binary-from-github-releases

// TODO: Support install script equivalent for TF
// This is based on option 2 here https://github.com/vmware-tanzu/tanzu-framework/blob/main/docs/cli/getting-started.md#option-2-using-install-script

// TODO: Should we support LegacyInstall("v0.20.0") too? Or just have one of them? Which one?
// Example: LegacyInstall("0.20.0")
// This is based on legacy installation method https://github.com/vmware-tanzu/tanzu-framework/blob/main/docs/cli/getting-started.md#legacy-method-to-install-plugins-with-api-driven-plugin-discovery-deactivated
func LegacyInstall(version string) error {
	log.Infof("Starting legacy install of TF version %s", version)

	if runtime.GOOS == platforms.WINDOWS {
		return fmt.Errorf("automated installation of TF on windows is not yet supported")
	}

	artifactUrl, err := getTfArtifactUrl(version)
	if err != nil {
		return fmt.Errorf("error getting TF artifact URL: %v", err)
	}

	artifactName := getArtifactNameFromUrl(artifactUrl)

	// TODO: Maybe avoid downloading again if there is a file already present locally with same checksum?
	// We could use data like etag header etc. This could act like a cache for us :) To avoid redownloading :D

	// TODO: Maybe do an integrity check for downloaded artifact with the tanzu-framework-executables-checksums.txt and tanzu-framework-executables-checksums.txt.asc files which has to be downloaded separately -
	// but we don't have to download it into a file, just pulling it / downloading it into the program memory and getting
	// the checksum and verifying it works too!

	// TODO: Maybe change this naming? The package (download) or function name (DownloadFileFromUrl).
	// It reads weird when it says download twice
	err = download.DownloadFileFromUrl(artifactUrl, artifactName)
	if err != nil {
		return fmt.Errorf("error downloading TF artifact: %v", err)
	}

	targetDirectory := getTargetDirectory()
	// extract tar ball or zip based on previous step
	extract.Extract(artifactName, targetDirectory)

	// TODO: install the TF core `tanzu` CLI to /usr/local/bin in Linux and MacOS
	// Example for MacOS - `install /var/folders/4z/09jpfvfj6c19lxl7ch78pzvc0000gn/T/tf-install-1650248624/tanzu-core-darwin_amd64 /usr/local/bin/tanzu`
	// Example for Linux - `sudo install /tmp/tf-install-1650248624/tanzu-core-linux_amd64 /usr/local/bin/tanzu`
	coreCliBinPath := filepath.Join(targetDirectory, "cli", "core", fmt.Sprintf("v%s", version), fmt.Sprintf("tanzu-core-%s_%s", runtime.GOOS, runtime.GOARCH))
	targetCoreCliBinPath := filepath.Join("/", "usr", "local", "bin", "tanzu")

	err = os.Rename(coreCliBinPath, targetCoreCliBinPath)
	if err != nil {
		return fmt.Errorf("error installing tanzu core CLI from %s to %s: %v", coreCliBinPath, targetCoreCliBinPath, err)
	}

	// TODO: Handle TF core `tanzu` CLI installation in Windows

	// Run `tanzu config set features.global.context-aware-cli-for-plugins false` command
	err = tanzu.DisableContextAwareCliForPluginsGlobally()
	if err != nil {
		return fmt.Errorf("error disabling context aware cli for plugins globally: %v", err)
	}

	// If tanzu is installed and then -
	// If you have a previous version of tanzu CLI already installed and the config file ~/.config/tanzu/config.yaml is present, run this command to make sure the default plugin repo points to the right path.
	// `tanzu plugin repo update -b tanzu-cli-framework core`
	// TODO: For now running the above command always!
	// TODO: Convert magic strings "core" and "tanzu-cli-framework" to constants inside tanzu package maybe
	err = tanzu.PluginRepoUpdate("core", "tanzu-cli-framework")
	if err != nil {
		return fmt.Errorf("error updating plugin repository: %v", err)
	}

	// Run `tanzu plugin install --local <target-directory>/cli all` command to install all the plugins in the local directory
	err = tanzu.PluginInstall("all", filepath.Join(targetDirectory, "cli"))
	if err != nil {
		return fmt.Errorf("error installing all plugins: %v", err)
	}

	// Run `tanzu plugin list` to list the plugins
	err = tanzu.PluginList()
	if err != nil {
		return fmt.Errorf("error listing plugins: %v", err)
	}

	return nil
}

func getTargetDirectory() string {
	return filepath.Join(os.TempDir(), fmt.Sprintf("tf-install-%d", time.Now().Unix()))
}

// TODO: Duplicate of getArtifactNameFromUrl in testutils/tce/install.go .
// Is this duplication okay?
func getArtifactNameFromUrl(artifactUrl string) string {
	// TODO: Change name? tokens may not be great.
	tokens := strings.Split(artifactUrl, "/")
	// TODO: Change name? this is more of the artifact's file path but with just the name and nothing else,
	// so it will get downloaded to the current working directory of the program
	artifactName := tokens[len(tokens)-1]

	return artifactName
}
