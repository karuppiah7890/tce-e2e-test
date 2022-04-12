package main

import (
	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
	"github.com/karuppiah7890/tce-e2e-test/testutils/tce"
)

// Different ways to install TCE. Same / Similar applies to TF too
// 1. Install from source - installing from source code using golang, make etc with `make release`
// 2. Install latest stable release - by detecting latest stable release, current OS and architecture and pulling the
// appropriate artifact and installing TCE
// 3. Install latest daily official build - by detecting latest daily official build, current OS and architecture and pulling the
// appropriate artifact and installing TCE
// 4. Install latest non-stable / non-GA pre-release (RCs, alpha, beta etc) - by detecting latest non-stable / non-GA pre-release, current OS and architecture and pulling the
// appropriate artifact and installing TCE
// 5. Install from given tar ball URL. Note that this may fail / not work in some cases, for example, if one tries to install linux tar ball using URL in Mac OS by mistake
// 6. Install from given tar ball file path in local. Note that this may fail / not work in some cases, for example, if one tries to install linux tar ball using URL in Mac OS by mistake
// 7. Install from given TCE version - which can be stable (latest or not), prerelease(latest or not). Detect current OS and architecture and pull the
// appropriate artifact and install TCE
// 8. Install from given daily official build date - the build date can be latest or not. Detect current OS and architecture and pull the
// appropriate artifact and install TCE

func main() {
	log.InitLogger("tce-install")
	// TODO: Get version from flags (--version) or arguments
	version := "0.11.0"
	err := tce.Install(version)
	if err != nil {
		log.Fatalf("error occurred while installing TCE version %s: %v", version, err)
	}
}