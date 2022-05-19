package tf_test

import (
	"testing"

	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
	"github.com/karuppiah7890/tce-e2e-test/testutils/tf"
)

func TestTfLegacyInstall(t *testing.T) {
	version := "0.21.1"

	err := tf.LegacyInstall(version)
	if err != nil {
		log.Fatalf("error occurred while installing TC version %s: %v", version, err)
	}
}