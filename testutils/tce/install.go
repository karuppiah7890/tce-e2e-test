package tce

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/karuppiah7890/tce-e2e-test/testutils/extract"
	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
)

// TODO: Should we support getTceArtifactUrl("v0.11.0") too? Or just have one of them? Which one?
// Example: getTceArtifactUrl("0.11.0")
func getTceArtifactUrl(version string) (string, error) {
	log.Infof("Getting TCE artifact URL for version %s", version)
	// TODO: Convert magic strings like "amd64" to constants
	supportedArchitectures := []string{"amd64"}
	// TODO: Convert magic strings like "linux", "darwin", "windows" to constants
	supportedOperatingSystems := []string{"linux", "darwin", "windows"}
	// TODO: Maybe merge the supported OSes and artifact extension data as artifact extension should be
	// present for each supported operating system in this case and currently there's a duplication of
	// data here
	// TODO: Convert magic strings like "tar.gz", "zip" to constants
	// TODO: Convert magic strings like "linux", "darwin", "windows" to constants
	artifactExtensions := map[string]string{"linux": "tar.gz", "darwin": "tar.gz", "windows": "zip"}

	architecture := runtime.GOARCH
	if !isPresentIn(architecture, supportedArchitectures) {
		return "", fmt.Errorf("architecture %s is not supported by TCE. Supported architectures: %v", architecture, supportedArchitectures)
	}

	operatingSystem := runtime.GOOS

	if !isPresentIn(operatingSystem, supportedOperatingSystems) {
		return "", fmt.Errorf("operating System %s is not supported by TCE. Supported operating systems: %v", operatingSystem, supportedOperatingSystems)
	}

	artifactExtension := artifactExtensions[operatingSystem]
	// Example: https://github.com/vmware-tanzu/community-edition/releases/download/v0.11.0/tce-darwin-amd64-v0.11.0.tar.gz
	return fmt.Sprintf("https://github.com/vmware-tanzu/community-edition/releases/download/v%s/tce-%s-%s-v%s.%s", version, operatingSystem, architecture, version, artifactExtension), nil
}

// searches if needle is present in haystack
func isPresentIn(needle string, haystack []string) bool {
	for _, thing := range haystack {
		if thing == needle {
			return true
		}
	}

	return false
}

// TODO: Should we support Install("v0.11.0") too? Or just have one of them? Which one?
// Example: Install("0.11.0")
func Install(version string) error {
	log.Infof("Starting install of TCE version %s", version)
	artifactUrl, err := getTceArtifactUrl(version)
	if err != nil {
		return fmt.Errorf("error getting TCE artifact URL: %v", err)
	}

	// TODO: Change name? tokens may not be great.
	tokens := strings.Split(artifactUrl, "/")
	// TODO: Change name? this is more of the artifact's file path but with just the name and nothing else,
	// so it will get downloaded to the current working directory of the program
	artifactName := tokens[len(tokens)-1]

	// TODO: Maybe change this naming? The package (download) or function name (DownloadFileFromUrl).
	// It reads weird when it says download twice
	// err = download.DownloadFileFromUrl(artifactUrl, artifactName)
	// if err != nil {
	// 	return fmt.Errorf("error downloading TCE artifact: %v", err)
	// }

	targetDirectory := filepath.Join(os.TempDir(), fmt.Sprintf("tce-install-%d", time.Now().Unix()))
	// extract tar ball or zip based on previous step
	extract.Extract(artifactName, targetDirectory)

	// invoke install.sh or install.bat based on previous step
	dirEntries, err := os.ReadDir(targetDirectory)

	if err != nil {
		return fmt.Errorf("error reading target directory %s containing extracted files: %v", targetDirectory, err)
	}

	if len(dirEntries) != 1 {
		return fmt.Errorf("expected target directory %s to contain only 1 directory but it contains %d directories", targetDirectory, len(dirEntries))
	}

	tceDir := dirEntries[0]

	operatingSystem := runtime.GOOS
	installScriptExtensions := map[string]string{"linux": "sh", "darwin": "sh", "windows": ".bat"}
	// TODO: Maybe merge the supported OSes and script extension data as script extension should be
	// present for each supported operating system in this case and currently there's a duplication of
	// data here

	installScriptExtension := installScriptExtensions[operatingSystem]

	installScript := filepath.Join(targetDirectory, tceDir.Name(), fmt.Sprintf("install.%s", installScriptExtension))

	cmd := exec.Command(installScript)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if operatingSystem == "linux" || operatingSystem == "darwin" {
		cmd.Env = append(os.Environ(), "ALLOW_INSTALL_AS_ROOT=true")
	}

	log.Infof("Running the command `%v`", cmd.String())

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("error occurred while running the command `%v`: %v. Exit code: %d", cmd.String(), err, cmd.ProcessState.ExitCode())
	}

	return nil
}
