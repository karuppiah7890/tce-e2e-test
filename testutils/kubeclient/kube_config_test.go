package kubeclient_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/karuppiah7890/tce-e2e-test/testutils/kubeclient"
	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
)

func TestConfigDeletion(t *testing.T) {
	log.InitLogger("config-deletion")
	tmpKubeConfigPath, _ := os.Getwd()
	initialConfigFilePath := filepath.Join(tmpKubeConfigPath, "testdata", "test-kubeconfig.yaml")

	// TODO: Check what t.Run() returns.
	t.Run("when deleting a context that does not exist it should throw error", func(t *testing.T) {
		t.Parallel()
		configFilePath, cleanup := setupTestKubeConfigFile(initialConfigFilePath)
		defer cleanup()

		nonExistentContext := "contextDoesNotExist"
		err := kubeclient.DeleteContext(configFilePath, nonExistentContext)
		if err == nil {
			log.Fatalf("expected error while deleting non-existent context %s but got no error", nonExistentContext)
		}

		if !strings.Contains(err.Error(), fmt.Sprintf("could not find context named %s in kubeconfig file at path", nonExistentContext)) {
			log.Fatalf("expected error around finding non-existent context %s but got some other error: %v", nonExistentContext, err)
		}
	})

	t.Run("when deleting a context that exists it should delete the context without any errors", func(t *testing.T) {
		t.Parallel()
		configFilePath, cleanup := setupTestKubeConfigFile(initialConfigFilePath)
		defer cleanup()

		existingContext := "test-admin@test"
		err := kubeclient.DeleteContext(configFilePath, existingContext)
		if err != nil {
			log.Fatalf("expected no error while deleting existing context but got error: %v", err)
		}

		configFileData, err := ioutil.ReadFile(configFilePath)
		if err != nil {
			log.Fatalf("expected no error while reading temp config file, but got error: %v", err)
		}

		if strings.Contains(string(configFileData), "test-admin@test") {
			log.Fatalf("cluster context didn't delete successfully")
		}

	})

}

// setupTestKubeConfigFile clones the kubeconfig file at path initialConfigFilePath in OS temp directory
func setupTestKubeConfigFile(initialConfigFilePath string) (string, func() error) {
	configFile, err := os.CreateTemp(os.TempDir(), "test-kubeconfig-*.yaml")
	if err != nil {
		log.Fatalf("expected no error while creating temp file for test kubeconfig but got error: %v", err)
	}

	configFilePath := configFile.Name()

	cleanupFunc := func() error {
		// We should cleanup the file only when the test passes. When test fails, we need the file to check any problems
		return os.Remove(configFilePath)
	}

	err = CopyFile(initialConfigFilePath, configFilePath)
	if err != nil {
		log.Fatalf("expected no error while copying test kubeconfig to a temporary place but got error: %v", err)
	}

	return configFilePath, cleanupFunc
}

func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}
