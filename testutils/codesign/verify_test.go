package codesign_test

import (
	"testing"

	"github.com/karuppiah7890/tce-e2e-test/testutils/codesign"
	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
)

func TestVerify(t *testing.T) {
	log.InitLogger("verify-binary")

	t.Run("binary with valid signature", func(t *testing.T) {
		// TODO: Use a binary path that can work on all machines
		binPath := "/Users/karuppiahn/Downloads/tce-darwin-amd64-v0.11.0/tanzu"
		err := codesign.Verify(binPath)
		if err != nil {
			log.Fatalf("expected no error but got error: %v", err)
		}
	})

	t.Run("binary with invalid signature", func(t *testing.T) {
		// TODO: Use a binary path that can work on all machines
		binPath := "/Users/karuppiahn/Downloads/tce-darwin-amd64-v0.12.0-dev.1/tanzu"
		err := codesign.Verify(binPath)
		if err == nil {
			log.Fatalf("expected error but got no error")
		}

		// TODO: verify the error message content
	})
}
