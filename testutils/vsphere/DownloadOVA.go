package vsphere

import (
	"encoding/json"
	"fmt"
	"github.com/karuppiah7890/tce-e2e-test/testutils/download"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
	"gopkg.in/yaml.v3"
)

type Files struct {
	// TODO: Extract out this struct into a type for brevity and clarity
	DownloadFiles []struct {
		FileName       string `json:"fileName"`
		Sha1Checksum   string `json:"sha1checksum"`
		Sha256Checksum string `json:"sha256checksum"`
		Md5Checksum    string `json:"md5checksum"`
		Build          string `json:"build"`
		ReleaseDate    string `json:"releaseDate"`
		FileType       string `json:"fileType"`
		Description    string `json:"description"`
		FileSize       string `json:"fileSize"`
		Title          string `json:"title"`
		Version        string `json:"version"`
		Status         string `json:"status"`
		Uuid           string `json:"uuid"`
		Header         bool   `json:"header"`
		DisplayOrder   int    `json:"displayOrder"`
		Relink         bool   `json:"relink"`
		Rsync          bool   `json:"rsync"`
	} `json:"downloadFiles"`
}

type Tkr struct {
	// TODO: Extract out this struct into a type for brevity and clarity
	Ova []struct {
		Name   string `yaml:"name"`
		Osinfo struct {
			Name    string `yaml:"name"`
			Version string `yaml:"version"`
			Arch    string `yaml:"arch"`
		} `json:"osinfo"`
		Version string `yaml:"version"`
	} `yaml:"ova"`
	// TODO: Extract out this struct into a type for brevity and clarity
	Azure []struct {
		Sku             string `yaml:"sku"`
		Publisher       string `yaml:"publisher"`
		Offer           string `yaml:"offer"`
		Version         string `yaml:"version"`
		ThirdPartyImage bool   `yaml:"thirdPartyImage"`
		// TODO: Extract out this struct into a type for brevity and clarity
		Osinfo struct {
			Name    string `yaml:"name"`
			Version string `yaml:"version"`
			Arch    string `yaml:"arch"`
		} `yaml:"osinfo"`
	} `yaml:"azure"`
}
type Cc []string

const (
	configDir      = ".config"
	tanzuConfigDir = "tanzu"
	tkgConfigDir   = "tkg"
	bom            = "bom"
)

// Rename to RetrieveAndDownload
func RetrieveAndDownload(version, dir, fileName string) {

	url := fmt.Sprintf("https://download3.vmware.com/software/TCE-%s/%s", version, fileName)
	log.Infof(url)
	downloadFile := dir + fileName
	exist := fileExists(downloadFile)

	if exist {
		log.Infof("File exist, Skipping download")
	} else {
		log.Infof("Downloading file at %s from %s", downloadFile, url)
		download.DownloadFileFromUrl(url, downloadFile)
	}

}
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
func RetriveVersion(version string) []string {
	url := fmt.Sprintf("https://customerconnect.vmware.com/channel/public/api/v1.0/dlg/details?downloadGroup=TCE-%s", version)
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	jsonMap := Files{}
	json.Unmarshal(responseData, &jsonMap)

	ovaFiles := []string{}
	for i := range jsonMap.DownloadFiles {
		if jsonMap.DownloadFiles[i].FileType == "ova" {
			log.Info(jsonMap.DownloadFiles[i].FileName)
			ovaFiles = append(ovaFiles, jsonMap.DownloadFiles[i].FileName)
		}
	}
	return ovaFiles
}
func GetOvaFileNameFromTanzuFramework() []string {
	bomDir, err := GetTanzuBomConfigPath()
	if err != nil {
		log.Errorf("failed to resolve Tanzu Bom config path. Error: %s", err.Error())
	}
	log.Infof("Tanzu Framework tkg bom home dir %s", bomDir)
	bomFiles, err := ioutil.ReadDir(bomDir)
	if err != nil {
		log.Fatal(err)
	}
	ovaFiles := []string{}
	for _, file := range bomFiles {
		log.Infof("looking in files for Ova files %s", file.Name())
		fileData, err := ioutil.ReadFile(filepath.Join(bomDir, "/", file.Name()))
		if err != nil {
			log.Errorf("failed to resolve tanzu Bom config path. Error: %s", err.Error())
		}
		OvaFilesMap := &Tkr{}
		err2 := yaml.Unmarshal(fileData, &OvaFilesMap)
		if err2 != nil {
			log.Fatal(err2)
		}
		for _, value := range OvaFilesMap.Ova {
			log.Infof("%s-%s-%s-%s", value.Osinfo.Name, value.Osinfo.Version, "kube", value.Version)
			ovaFiles = append(ovaFiles, fmt.Sprintf("%s-%s-%s-%s", value.Osinfo.Name, value.Osinfo.Version, "kube", value.Version))
		}
	}
	return ovaFiles
}

//Todo: to be moved to common utils file
func GetTanzuConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to resolve tanzu config path. Error: %s", err.Error())
	}

	return filepath.Join(home, configDir, tanzuConfigDir), nil
}

//Todo: to be moved to common utils file
func GetTanzuBomConfigPath() (string, error) {
	path, err := GetTanzuConfigPath()
	if err != nil {
		return "", fmt.Errorf("failed to resolve tanzu TKG config path. Error: %s", err.Error())
	}

	return filepath.Join(path, tkgConfigDir, bom), nil
}
