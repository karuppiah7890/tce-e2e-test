package main

import (
	"os"

	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
	"github.com/karuppiah7890/tce-e2e-test/testutils/tf"
)

// Different ways to install TF
// 1. Install from source - installing from source code using golang, make etc with `make release`
// 2. Install latest stable release - by detecting latest stable release, current OS and architecture and pulling the
// appropriate artifact and installing TF
// 4. Install latest non-stable / non-GA pre-release (RCs, alpha, beta etc) - by detecting latest non-stable / non-GA pre-release, current OS and architecture and pulling the
// appropriate artifact and installing TF
// 5. Install from given tar ball URL. Note that this may fail / not work in some cases, for example, if one tries to install linux tar ball using URL in Mac OS by mistake
// 6. Install from given tar ball file path in local. Note that this may fail / not work in some cases, for example, if one tries to install linux tar ball using URL in Mac OS by mistake
// 7. Install from given TF version - which can be stable (latest or not), prerelease(latest or not). Detect current OS and architecture and pull the
// appropriate artifact and install TF

// TODO: Support doing uninstall to cleanup any existing installation first and then do fresh install, given a flag
// like `--uninstall` or similar

func main() {
	log.InitLogger("tf-install")
	// TODO: Get version from flags (--version) or arguments

	if len(os.Args) != 2 {
		log.Fatal("Provide exactly one argument with Tanzu Framework (TF) version. Example Usage: tf-installer 0.20.0")
	}

	version := os.Args[1]
	err := tf.LegacyInstall(version)
	if err != nil {
		log.Fatalf("error occurred while installing TF version %s: %v", version, err)
	}
}
