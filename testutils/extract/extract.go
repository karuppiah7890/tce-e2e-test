package extract

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
)

const TARGZ = "tar.gz"
const ZIP = "zip"

// TODO: Replace this with a third party library maybe? So that we have less code to maintain
func Extract(compressedFile string, targetDirectoryToExtract string) error {
	log.Infof("Starting to extract %s to %s", compressedFile, targetDirectoryToExtract)
	supportedCompressedFormats := []string{TARGZ, ZIP}
	compressionFormat, ok := getCompressionFormat(compressedFile, supportedCompressedFormats)
	if !ok {
		return fmt.Errorf("compressed file %s is not supported by extractor to extract. Supported compressed formats: %v", compressedFile, supportedCompressedFormats)
	}
	extractionFunctions := map[string]func(compressedFile string, targetDirectoryToExtract string) error{
		TARGZ: extractTarGz,
		ZIP:   extractZip,
	}
	extractionFunction, ok := extractionFunctions[compressionFormat]
	if !ok {
		// TODO: Ideally this should NOT happen at this point. We should be able to get
		// extraction function given the compression format
		log.DPanicf("No extractor found for compression format %s", compressionFormat)
		return fmt.Errorf("error occurred while getting extractor for extracting file %s: No extractor found for compression format %s", compressedFile, compressionFormat)
	}

	err := extractionFunction(compressedFile, targetDirectoryToExtract)

	if err != nil {
		return fmt.Errorf("error occurred while extracting %s to %s: %v", compressedFile, targetDirectoryToExtract, err)
	}

	log.Infof("Done extracting %s to %s", compressedFile, targetDirectoryToExtract)

	return nil
}

// TODO: Should we use file content to get the compression format? By looking at file metadata in the file data
// getCompressionFormat gets compression format from the compressed file's name or returns empty
func getCompressionFormat(compressedFile string, supportedCompressedFormats []string) (string, bool) {
	for _, format := range supportedCompressedFormats {
		if strings.HasSuffix(compressedFile, format) {
			return format, true
		}
	}

	return "", false
}

// TODO: Replace this with a third party library maybe? So that we have less code to maintain
func extractZip(compressedFile string, targetDirectoryToExtract string) error {
	return fmt.Errorf("extracting zip files is not yet implemented")
}

// TODO: Replace this with a third party library maybe? So that we have less code to maintain
func extractTarGz(compressedFile string, targetDirectoryToExtract string) error {
	targzReader, err := os.Open(compressedFile)
	if err != nil {
		return fmt.Errorf("error occurred while opening %s: %v", compressedFile, err)
	}
	defer targzReader.Close()

	tarArchive, err := gzip.NewReader(targzReader)
	if err != nil {
		return fmt.Errorf("error occurred while trying to process %s: %v", compressedFile, err)
	}
	defer tarArchive.Close()

	tarReader := tar.NewReader(tarArchive)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("error occurred while trying to process %s: %v", compressedFile, err)
		}

		path := filepath.Join(targetDirectoryToExtract, header.Name)
		info := header.FileInfo()
		if info.IsDir() {
			if err = os.MkdirAll(path, info.Mode()); err != nil {
				return fmt.Errorf("error occurred while trying to process %s: error occurred while trying to create directory %s with permission %v: %v", compressedFile, path, info.Mode(), err)
			}
			continue
		}

		file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return fmt.Errorf("error occurred while trying to process %s: error occurred while trying to open file %s: %v", compressedFile, path, err)
		}
		defer file.Close()
		_, err = io.Copy(file, tarReader)
		if err != nil {
			return fmt.Errorf("error occurred while trying to process %s: error occurred while trying to extract %s: %v", compressedFile, compressedFile, err)
		}
	}

	return nil
}
