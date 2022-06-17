package aws

import (
	"os"

	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
	"github.com/karuppiah7890/tce-e2e-test/testutils/utils"
)

const AccountID = "AWS_ACCOUNT_ID"
const AccessKey = "AWS_ACCESS_KEY_ID"
const SecretKey = "AWS_SECRET_ACCESS_KEY"
const B64Creds = "AWS_B64ENCODED_CREDENTIALS"
const Region = "AWS_REGION"
const SshPublicName = "AWS_SSH_KEY_NAME"

type TestSecrets struct {
	AccountID     string
	AccessKey     string
	SecretKey     string
	B64Creds      string
	Region        string
	SshPublicName string
}

func ExtractAwsTestSecretsFromEnvVars() TestSecrets {
	utils.CheckRequiredEnvVars(PROVIDER)

	log.Info("Extracting AWS test secrets from environment variables")

	return TestSecrets{
		AccountID:     os.Getenv(AccountID),
		AccessKey:     os.Getenv(AccessKey),
		SecretKey:     os.Getenv(SecretKey),
		B64Creds:      os.Getenv(B64Creds),
		Region:        os.Getenv(Region),
		SshPublicName: os.Getenv(SshPublicName),
	}
}
