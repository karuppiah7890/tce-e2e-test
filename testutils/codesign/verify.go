package codesign

import (
	"fmt"
	"runtime"

	"github.com/karuppiah7890/tce-e2e-test/testutils/clirunner"
	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
	"github.com/karuppiah7890/tce-e2e-test/testutils/search"
)

func Verify(binPath string) error {
	// TODO: Supported platforms code is duplicated. And this code is based on TCE support only. This function is kind of generic though,
	// it simply checks if the binary is signed correctly. Or do we assume only supported binaries are passed?

	// TODO: Convert magic strings like "linux/amd64", "darwin/amd64", "windows/amd64" to constants
	supportedPlatforms := []string{"linux/amd64", "darwin/amd64", "windows/amd64"}

	architecture := runtime.GOARCH
	operatingSystem := runtime.GOOS
	platform := fmt.Sprintf("%s/%s", operatingSystem, architecture)
	if !search.IsPresentIn(platform, supportedPlatforms) {
		return fmt.Errorf("platform %s is not supported by TCE. Supported platforms: %v", platform, supportedPlatforms)
	}

	// For MacOS binary, use spctl and check it
	if runtime.GOOS == "darwin" {
		cmd := spctlCommandForBinary(binPath)
		exitCode, err := clirunner.Run(cmd)
		if err != nil {
			return fmt.Errorf("error occurred while verifying %s binary signature. Exit code: %v. Error: %v", binPath, exitCode, err)
		}
	}

	// For WindowsOS binary, use SignTool https://docs.microsoft.com/en-us/windows/win32/seccrypto/signtool
	if runtime.GOOS == "windows" {
		cmd := signtoolCommandForBinary(binPath)
		exitCode, err := clirunner.Run(cmd)
		if err != nil {
			return fmt.Errorf("error occurred while verifying %s binary signature. Exit code: %v. Error: %v", binPath, exitCode, err)
		}
	}

	return nil
}

func spctlCommandForBinary(binPath string) clirunner.Cmd {
	return clirunner.Cmd{
		Name: "spctl",
		Args: []string{
			"-vv",
			"--type",
			"install",
			"--asses",
			binPath,
		},
		Stdout: log.InfoWriter,
		Stderr: log.ErrorWriter,
	}
}

func signtoolCommandForBinary(binPath string) clirunner.Cmd {
	return clirunner.Cmd{
		Name: "signtool",
		Args: []string{
			"verify",
			"/pa",
			binPath,
		},
		Stdout: log.InfoWriter,
		Stderr: log.ErrorWriter,
	}
}
