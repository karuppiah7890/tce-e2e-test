package tce

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/karuppiah7890/tce-e2e-test/testutils/download"
	"github.com/karuppiah7890/tce-e2e-test/testutils/extract"
	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
	"github.com/karuppiah7890/tce-e2e-test/testutils/platforms"
)

const SHELL = "sh"
const BAT = "bat"

// TODO: Should we support getTceArtifactUrl("v0.11.0") too? Or just have one of them? Which one?
// Example: getTceArtifactUrl("0.11.0")
func getTceArtifactUrl(version string) (string, error) {
	log.Infof("Getting TCE artifact URL for version %s", version)

	// TODO: Convert magic strings like "linux/amd64", "darwin/amd64", "windows/amd64" to constants
	supportedPlatforms := []string{"linux/amd64", "darwin/amd64", "windows/amd64"}

	// TODO: Maybe merge the supported OSes and artifact extension data as artifact extension should be
	// present for each supported operating system in this case and currently there's a duplication of
	// data here
	artifactExtensions := map[string]string{platforms.LINUX: extract.TARGZ, platforms.DARWIN: extract.TARGZ, platforms.WINDOWS: extract.ZIP}

	architecture := runtime.GOARCH
	operatingSystem := runtime.GOOS
	platform := fmt.Sprintf("%s/%s", operatingSystem, architecture)
	if !isPresentIn(platform, supportedPlatforms) {
		return "", fmt.Errorf("platform %s is not supported by TCE. Supported platforms: %v", platform, supportedPlatforms)
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

	artifactName := getArtifactNameFromUrl(artifactUrl)

	// TODO: Maybe avoid downloading again if there is a file already present locally with same checksum?
	// We could use data like etag header etc. This could act like a cache for us :) To avoid redownloading :D

	// TODO: Maybe do an integrity check for downloaded artifact with the tce-checksums.txt file which has to be downloaded separately -
	// but we don't have to download it into a file, just pulling it / downloading it into the program memory and getting
	// the checksum works too!

	// TODO: Maybe change this naming? The package (download) or function name (DownloadFileFromUrl).
	// It reads weird when it says download twice
	err = download.DownloadFileFromUrl(artifactUrl, artifactName)
	if err != nil {
		return fmt.Errorf("error downloading TCE artifact: %v", err)
	}

	targetDirectory := getTargetDirectory()
	// extract tar ball or zip based on previous step
	extract.Extract(artifactName, targetDirectory)

	return invokeTceInstallScript(targetDirectory)
}

func getTargetDirectory() string {
	return filepath.Join(os.TempDir(), fmt.Sprintf("tce-install-%d", time.Now().Unix()))
}

func invokeTceInstallScript(targetDirectory string) error {
	dirEntries, err := os.ReadDir(targetDirectory)

	if err != nil {
		return fmt.Errorf("error reading target directory %s containing extracted files: %v", targetDirectory, err)
	}

	if len(dirEntries) != 1 {
		return fmt.Errorf("expected target directory %s to contain only 1 directory but it contains %d directories", targetDirectory, len(dirEntries))
	}

	tceDir := dirEntries[0]

	operatingSystem := runtime.GOOS
	installScriptExtensions := map[string]string{platforms.LINUX: SHELL, platforms.DARWIN: SHELL, platforms.WINDOWS: BAT}
	// TODO: Maybe merge the supported OSes and script extension data as script extension should be
	// present for each supported operating system in this case and currently there's a duplication of
	// data here

	installScriptExtension := installScriptExtensions[operatingSystem]
	installScript := filepath.Join(targetDirectory, tceDir.Name(), fmt.Sprintf("install.%s", installScriptExtension))

	// invoke install.sh or install.bat based on OS
	cmd := exec.Command(installScript)
	cmd.Dir = filepath.Join(targetDirectory, tceDir.Name())
	cmd.Stdout = log.InfoWriter
	cmd.Stderr = log.ErrorWriter

	if operatingSystem == platforms.LINUX || operatingSystem == platforms.DARWIN {
		cmd.Env = append(os.Environ(), "ALLOW_INSTALL_AS_ROOT=true")
	}

	log.Infof("Running the command `%v`", cmd.String())

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("error occurred while running the command `%v`: %v. Exit code: %d", cmd.String(), err, cmd.ProcessState.ExitCode())
	}

	return nil
}

func getArtifactNameFromUrl(artifactUrl string) string {
	// TODO: Change name? tokens may not be great.
	tokens := strings.Split(artifactUrl, "/")
	// TODO: Change name? this is more of the artifact's file path but with just the name and nothing else,
	// so it will get downloaded to the current working directory of the program
	artifactName := tokens[len(tokens)-1]

	return artifactName
}
