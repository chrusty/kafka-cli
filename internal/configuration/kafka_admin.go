package configuration

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/chrusty/kafka-cli/internal/types"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/aws_msk_iam_v2"
	"github.com/segmentio/kafka-go/sasl/scram"
	"github.com/sirupsen/logrus"
)

// Admin returns a Kafka Client based on our config:
func (kc *KafkaConfig) Admin(logger *logrus.Logger) (*kafka.Client, error) {

	// Prepare a low-level client:
	client := &kafka.Client{
		Addr: kafka.TCP(kc.BootstrapServers...),
	}

	switch kc.SecurityProtocol {

	case types.SecProtocolPlaintext:

		client.Transport = &kafka.Transport{}

	case types.SecProtocolAWSMSKIAM:

		awsConfig, err := config.LoadDefaultConfig(context.TODO())
		if err != nil {
			return nil, err
		}
		fmt.Printf("Config: %v\n", awsConfig)

		// Define an SASL mechanism from an AWS client config:
		saslMechanism := aws_msk_iam_v2.NewMechanism(awsConfig)
		saslMechanism.Start(context.TODO())

		// Transport:
		client.Transport = &kafka.Transport{
			SASL: saslMechanism,
			TLS:  &tls.Config{},
		}

	case types.SecProtocolSSL:

		client.Transport = &kafka.Transport{
			TLS: &tls.Config{},
		}

	case types.SecProtocolSaslPlaintext:

		// Define an SASL mechanism:
		saslMechanism, err := scram.Mechanism(
			kc.saslAlgorithm(logger),
			kc.Username,
			kc.Password,
		)
		if err != nil {
			return nil, err
		}

		// Transport:
		client.Transport = &kafka.Transport{
			SASL: saslMechanism,
		}

	case types.SecProtocolSaslSSL:

		// Define an SASL mechanism:
		saslMechanism, err := scram.Mechanism(
			kc.saslAlgorithm(logger),
			kc.Username,
			kc.Password,
		)
		if err != nil {
			return nil, err
		}

		// Transport:
		client.Transport = &kafka.Transport{
			SASL: saslMechanism,
			TLS:  &tls.Config{},
		}

	default:
		return nil, fmt.Errorf("unsupported security protocol %s", kc.SecurityProtocol)
	}

	return client, nil
}
