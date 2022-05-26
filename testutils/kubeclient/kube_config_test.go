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
	tmpKubeConfigPath, _ := os.Getwd()
	// Creating a temp test config file
	Copy(tmpKubeConfigPath+"/tmp/test_config", tmpKubeConfigPath+"/tmp/temp_test_config")

	// deleting the context
	err := kubeclient.DeleteContext(tmpKubeConfigPath+"/tmp/test_config", "aman")
	if err != nil {
		os.Remove(tmpKubeConfigPath + "/tmp/temp_test_config")
		log.Fatalf("%v", err)
	}
	configFileData, err := ioutil.ReadFile(tmpKubeConfigPath + "/tmp/test_config")
	if err != nil {
		os.Remove(tmpKubeConfigPath + "/tmp/temp_test_config")
		log.Fatalf("error while reading temp config file")
	}

	// checking context is deleted or not
	if strings.Contains(string(configFileData), "standalone-v0.10.0-rc.3-admin@standalone-v0.10.0-rc.3") {
		os.Remove(tmpKubeConfigPath + "/tmp/temp_test_config")
		log.Fatalf("cluster context didn't delete sucessfully")
	}

	os.Remove(tmpKubeConfigPath + "/tmp/temp_test_config")

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
