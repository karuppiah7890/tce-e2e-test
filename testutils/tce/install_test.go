package tce_test

import (
	"testing"

	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
	"github.com/karuppiah7890/tce-e2e-test/testutils/tanzu"
	"github.com/karuppiah7890/tce-e2e-test/testutils/tce"
)

func TestTceInstall(t *testing.T) {
	log.InitLogger("tce-install-test")
	version := "0.12.1"

	err := tce.Install(version, "")
	if err != nil {
		log.Fatalf("error occurred while installing TCE version %s: %v", version, err)
	}

	err = tanzu.PrintTanzuVersion()
	if err != nil {
		log.Fatalf("expected no error but error occurred while printing tanzu CLI version: %v", err)
	}
}
