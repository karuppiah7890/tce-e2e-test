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

	err := CopyFile(initialConfigFile, testConfigFile)
	if err != nil {
		os.Remove(testConfigFile)
		log.Fatalf("expected no error while copying test kubeconfig to a temporary place but got error: %v", err)
	}

	nonExistentContext := "contextDoesNotExist"
	err = kubeclient.DeleteContext(testConfigFile, nonExistentContext)
	if err == nil {
		os.Remove(testConfigFile)
		log.Fatal("expected error while deleting non-existent context but got no error")
	}

	if !strings.Contains(err.Error(), fmt.Sprintf("could not find context named %s in kubeconfig file at path", nonExistentContext)) {
		log.Fatalf("expected error around finding non-existent context but got some other error: %v", err)
	}

	existingContext := "standalone-v0.10.0-rc.3-admin@standalone-v0.10.0-rc.3"
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

	if strings.Contains(string(configFileData), "standalone-v0.10.0-rc.3-admin@standalone-v0.10.0-rc.3") {
		os.Remove(testConfigFile)
		log.Fatalf("cluster context didn't delete successfully")
	}

	os.Remove(testConfigFile)

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
