package kubeclient_test

import (
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/karuppiah7890/tce-e2e-test/testutils/kubeclient"
	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
)

func TestConfigDeletion(t *testing.T) {
	log.InitLogger("config-deletion")
	tmpKubeConfigPath, _ := os.Getwd()

	Copy(tmpKubeConfigPath+"/testdata/test_config", tmpKubeConfigPath+"/testdata/temp_test_config")

	err := kubeclient.DeleteContext(tmpKubeConfigPath+"/testdata/temp_test_config", "aman")
	if err == nil {
		os.Remove(tmpKubeConfigPath + "/testdata/temp_test_config")
		log.Fatal("expected error while deleting non-existent context but got no error")
	}

	if !strings.Contains(err.Error(), "could not find context named aman in kubeconfig file at path") {
		log.Fatalf("expected error around finding non-existent context but got some other error: %v", err)
	}

	err = kubeclient.DeleteContext(tmpKubeConfigPath+"/testdata/temp_test_config", "standalone-v0.10.0-rc.3-admin@standalone-v0.10.0-rc.3")
	if err != nil {
		os.Remove(tmpKubeConfigPath + "/testdata/temp_test_config")
		log.Fatalf("expected no error while deleting existent context but got error: %v", err)
	}

	configFileData, err := ioutil.ReadFile(tmpKubeConfigPath + "/testdata/temp_test_config")
	if err != nil {
		os.Remove(tmpKubeConfigPath + "/testdata/temp_test_config")
		log.Fatalf("expected no error while reading temp config file, but got error: %v", err)
	}

	if strings.Contains(string(configFileData), "standalone-v0.10.0-rc.3-admin@standalone-v0.10.0-rc.3") {
		os.Remove(tmpKubeConfigPath + "/testdata/temp_test_config")
		log.Fatalf("cluster context didn't delete successfully")
	}

	os.Remove(tmpKubeConfigPath + "/testdata/temp_test_config")

}

func Copy(src, dst string) error {
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
