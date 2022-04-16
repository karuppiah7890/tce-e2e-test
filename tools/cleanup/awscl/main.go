package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/karuppiah7890/tce-e2e-test/testutils/aws"
	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
	"os"
)

// S3ListBucketsAPI defines the interface for the ListBuckets function.
// We use this interface to test the function using a mocked service.
type S3ListBucketsAPI interface {
	ListBuckets(ctx context.Context,
		params *s3.ListBucketsInput,
		optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error)
}

// GetAllBuckets retrieves a list of your Amazon Simple Storage Service (Amazon S3) buckets.
// Inputs:
//     c is the context of the method call, which includes the AWS Region.
//     api is the interface that defines the method call.
//     input defines the input arguments to the service call.
// Output:
//     If success, a ListBucketsOutput object containing the result of the service call and nil.
//     Otherwise, nil and an error from the call to ListBuckets.
func GetAllBuckets(c context.Context, api S3ListBucketsAPI, input *s3.ListBucketsInput) (*s3.ListBucketsOutput, error) {
	return api.ListBuckets(c, input)
}

func main() {

	log.InitLogger("aws")

	// TODO: Support providing multiple resource group names to delete.
	// TODO: Support running delete on multiple resource group names concurrently / parallely vs sequentially based on the order of
	// occurrence in the list in the CLI command.
	// Use urfave/cli for handling variadic arguments? or use plain golang std library as usual?

	if len(os.Args) != 2 {
		log.Fatal("Usage: ./aws <resource-group-name>")
	}

	awsTestSecrets := aws.ExtractAwsTestSecretsFromEnvVars()
	log.Infof("%s", awsTestSecrets)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := s3.NewFromConfig(cfg)

	input := &s3.ListBucketsInput{}

	result, err := GetAllBuckets(context.TODO(), client, input)
	if err != nil {
		fmt.Println("Got an error retrieving buckets:")
		fmt.Println(err)
		return
	}

	for _, bucket := range result.Buckets {
		fmt.Println(*bucket.Name + ": " + bucket.CreationDate.Format("2006-01-02 15:04:05 Monday"))
	}
}
