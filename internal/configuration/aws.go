package configuration

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	ec2rolecredsv2 "github.com/aws/aws-sdk-go-v2/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/credentials"
	ec2rolecredsv1 "github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	sigv4 "github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/segmentio/kafka-go/sasl"
	"github.com/segmentio/kafka-go/sasl/aws_msk_iam"
	"github.com/segmentio/kafka-go/sasl/aws_msk_iam_v2"
)

func AWSSaslMechanismV1() (sasl.Mechanism, error) {

	// Get an AWS session:
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	// Attempt to retrieve EC2 role-provider credentials:
	awsCredentials := credentials.NewChainCredentials(
		[]credentials.Provider{
			&credentials.EnvProvider{},
			&ec2rolecredsv1.EC2RoleProvider{
				Client: ec2metadata.New(sess),
			},
		})

	// Look up the AWS region:
	awsRegion, ok := os.LookupEnv("AWS_REGION")
	if !ok {
		return nil, fmt.Errorf("unable to read AWS_REGION env-var")
	}

	// Define an SASL mechanism from an AWS client config:
	saslMechanism := &aws_msk_iam.Mechanism{
		Signer: sigv4.NewSigner(awsCredentials),
		Region: awsRegion,
	}

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
	awsConfig.Credentials = ec2rolecredsv2.New()
	if err != nil {
		return nil, err
	}

	// Define an SASL mechanism from an AWS client config:
	saslMechanism := aws_msk_iam_v2.NewMechanism(awsConfig)
	saslMechanism.Region = awsRegion
	saslMechanism.Start(context.TODO())

	return saslMechanism, nil
}
