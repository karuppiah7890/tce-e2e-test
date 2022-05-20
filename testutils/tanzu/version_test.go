package tanzu_test

import (
	"testing"

	"github.com/karuppiah7890/tce-e2e-test/testutils/tanzu"
	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
)

func TestPrintTanzuVersion(t *testing.T) {
	log.InitLogger("tanzu-version")
	err := tanzu.PrintTanzuVersion()

	if err != nil {
		log.Fatalf("expected no error but error occurred while printing tanzu CLI version: %v", err)
	}
}
