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
	initialConfigFile := filepath.Join(tmpKubeConfigPath, "testdata", "test_config.yaml")
	testConfigFile := filepath.Join(tmpKubeConfigPath, "testdata", "temp_test_config.yaml")

	// TODO: Check what t.Run() returns.
	// TODO: Check how we can run multiple t.Run() in parallel by using t.Parallel() and making modifications to our testing
	// logic. Currently the tests run in sequence
	t.Run("when deleting a context that does not exist it should throw error", func(t *testing.T) {
		err := CopyFile(initialConfigFile, testConfigFile)
		if err != nil {
			os.Remove(testConfigFile)
			log.Fatalf("expected no error while copying test kubeconfig to a temporary place but got error: %v", err)
		}
		nonExistentContext := "contextDoesNotExist"
		err = kubeclient.DeleteContext(testConfigFile, nonExistentContext)
		if err == nil {
			os.Remove(testConfigFile)
			log.Fatalf("expected error while deleting non-existent context %s but got no error", nonExistentContext)
		}

		if !strings.Contains(err.Error(), fmt.Sprintf("could not find context named %s in kubeconfig file at path", nonExistentContext)) {
			os.Remove(testConfigFile)
			log.Fatalf("expected error around finding non-existent context %s but got some other error: %v", nonExistentContext, err)
		}

		os.Remove(testConfigFile)
	})

	t.Run("when deleting a context that exists it should delete the context without any errors", func(t *testing.T) {
		err := CopyFile(initialConfigFile, testConfigFile)
		if err != nil {
			os.Remove(testConfigFile)
			log.Fatalf("expected no error while copying test kubeconfig to a temporary place but got error: %v", err)
		}

		existingContext := "test-admin@test"
		err = kubeclient.DeleteContext(testConfigFile, existingContext)
		if err != nil {
			os.Remove(testConfigFile)
			log.Fatalf("expected no error while deleting existing context but got error: %v", err)
		}

		configFileData, err := ioutil.ReadFile(testConfigFile)
		if err != nil {
			os.Remove(testConfigFile)
			log.Fatalf("expected no error while reading temp config file, but got error: %v", err)
		}

		if strings.Contains(string(configFileData), "test-admin@test") {
			os.Remove(testConfigFile)
			log.Fatalf("cluster context didn't delete successfully")
		}

		os.Remove(testConfigFile)
	})

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
