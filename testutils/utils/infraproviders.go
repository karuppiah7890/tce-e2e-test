package utils

// Question: Move this to a package named infrastructure?
// Say something like providers or infra providers? Or it looks when calling utils.AWS from a caller perspective

const AWS = "aws"
const VSPHERE = "vsphere"
const AZURE = "azure"
const Docker = "docker"

// TODO: Change name?
type Provider interface {
	Name() string
	// TODO: Change CheckRequiredEnvVars to GetListOfRequiredEnvVars ? And do check in a common manner?
	CheckRequiredEnvVars() bool
}
