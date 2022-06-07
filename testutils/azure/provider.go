package azure

import "github.com/karuppiah7890/tce-e2e-test/testutils/utils"

// TODO: Change name?
type Provider struct {
}

func (provider Provider) CheckRequiredEnvVars() bool {
	CheckRequiredAzureEnvVars()
	return true
}

func (provider Provider) Name() string {
	return "aws"
}

// TODO: Change name?
var PROVIDER utils.Provider = Provider{}
