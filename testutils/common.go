package testutils

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
)

// TODO: Rename this file / refactor and move this method to an appropriate place

func CheckRequiredEnvVars(requiredEnvVars []string) []error {
	log.Info("Checking required environment variables")

	errors := make([]error, 0, len(requiredEnvVars))
	for _, envVar := range requiredEnvVars {
		// TODO: Do we also error out when the env var value is defined but empty??
		_, isDefined := os.LookupEnv(envVar)

		if !isDefined {
			errors = append(errors, fmt.Errorf("environment variable `%s` is required but not defined", envVar))
		}
	}

	return errors
}

func SplitYAML(resources []byte) ([][]byte, error) {
	dec := yaml.NewDecoder(bytes.NewReader(resources))

	var res [][]byte
	for {
		var value interface{}
		err := dec.Decode(&value)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		valueBytes, err := yaml.Marshal(value)
		if err != nil {
			return nil, err
		}
		res = append(res, valueBytes)
	}
	return res, nil
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
