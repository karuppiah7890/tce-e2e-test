package download

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
)

// TODO: Replace this with a third party library maybe? So that we have less code to maintain
func DownloadFileFromUrl(fileUrl string, targetFilePath string) error {
	log.Infof("Starting download of %s to %s", fileUrl, targetFilePath)

	output, err := os.Create(targetFilePath)
	if err != nil {
		return fmt.Errorf("error while creating %s: %v", targetFilePath, err)
	}
	defer output.Close()

	response, err := http.Get(fileUrl)
	if err != nil {
		return fmt.Errorf("error while downloading %s: %v", fileUrl, err)
	}
	defer response.Body.Close()

	n, err := io.Copy(output, response.Body)
	if err != nil {
		return fmt.Errorf("error while downloading %s: %v", fileUrl, err)
	}

	log.Infof("%d bytes downloaded for %s", n, fileUrl)

	return nil
}
