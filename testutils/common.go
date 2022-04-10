package testutils

import (
	"fmt"
	"os"

	"go.uber.org/zap"
)

func CheckRequiredEnvVars(requiredEnvVars []string) []error {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()

	sugar.Info("Checking required environment variables")

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
