package configuration

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/ec2rolecreds"
	"github.com/segmentio/kafka-go/sasl"
	"github.com/segmentio/kafka-go/sasl/aws_msk_iam_v2"
)

func AWSSaslMechanismV1() (sasl.Mechanism, error) {

	// Look up the AWS region:
	awsRegion, ok := os.LookupEnv("AWS_REGION")
	if !ok {
		return nil, fmt.Errorf("unable to read AWS_REGION env-var")
	}

	// Get an AWS config:
	awsConfig, err := config.LoadDefaultConfig(context.TODO())
	awsConfig.Region = awsRegion
	awsConfig.Credentials = ec2rolecreds.New()
	if err != nil {
		return nil, err
	}

	// Define an SASL mechanism from an AWS client config:
	saslMechanism := aws_msk_iam_v2.NewMechanism(awsConfig)
	saslMechanism.Region = awsRegion
	saslMechanism.Start(context.TODO())

	return saslMechanism, nil
}

func AWSSaslMechanismV2() (sasl.Mechanism, error) {

	// Look up the AWS region:
	awsRegion, ok := os.LookupEnv("AWS_REGION")
	if !ok {
		return nil, fmt.Errorf("unable to read AWS_REGION env-var")
	}

	// Get an AWS config:
	awsConfig, err := config.LoadDefaultConfig(context.TODO())
	awsConfig.Region = awsRegion
	awsConfig.Credentials = ec2rolecreds.New()
	if err != nil {
		return nil, err
	}

	// Define an SASL mechanism from an AWS client config:
	saslMechanism := aws_msk_iam_v2.NewMechanism(awsConfig)
	saslMechanism.Region = awsRegion
	saslMechanism.Start(context.TODO())

	return saslMechanism, nil
}
