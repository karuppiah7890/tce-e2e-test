package codesign_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/karuppiah7890/tce-e2e-test/testutils/codesign"
	"github.com/karuppiah7890/tce-e2e-test/testutils/download"
	"github.com/karuppiah7890/tce-e2e-test/testutils/extract"
	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
)

func TestVerify(t *testing.T) {
	log.InitLogger("verify-binary")

	testDir, err := os.MkdirTemp(os.TempDir(), "verify test")
	if err != nil {
		t.Fatalf("expected no error while creating temporary directory for test, but got error: %v", err)
	}

	// Below is for Darwin binaries

	t.Run("binary with valid signature at a path containing space in file path", func(t *testing.T) {
		artifactPath := filepath.Join(testDir, "tanzu-cli-darwin-amd64.tar.gz")
		err = download.DownloadFileFromUrl("https://github.com/vmware-tanzu/tanzu-framework/releases/download/v0.21.0/tanzu-cli-darwin-amd64.tar.gz", artifactPath)
		if err != nil {
			t.Fatalf("expected no error while downloading artifact, but got error: %v", err)
		}
		extract.Extract(artifactPath, testDir)
		binPath := filepath.Join(testDir, "v0.21.0", "tanzu-core-darwin_amd64")
		err = codesign.Verify(binPath)
		if err != nil {
			t.Fatalf("expected no error but got error: %v", err)
		}
	})

	t.Run("binary with invalid signature at a path containing space in file path", func(t *testing.T) {
		artifactPath := filepath.Join(testDir, "kubectl")
		err = download.DownloadFileFromUrl("https://dl.k8s.io/v1.24.0/bin/darwin/amd64/kubectl", artifactPath)
		if err != nil {
			t.Fatalf("expected no error while downloading artifact, but got error: %v", err)
		}
		binPath := filepath.Join(testDir, "kubectl")
		err := codesign.Verify(binPath)
		if err == nil {
			t.Fatalf("expected error but got no error")
		}

		// TODO: verify the error message content because the error should be because of signing issues and not due
		// to any other error
	})

	// TODO: Test Windows binaries - give no error for signed binary and error for unsigned binary. Check for file paths with space in the name

	// TODO: Test Linux binaries - give no error. Check for file paths with space in the name
}
