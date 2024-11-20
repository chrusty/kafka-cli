package configuration

import (
	"crypto/tls"
	"fmt"

	"github.com/chrusty/kafka-cli/internal/types"
	"github.com/segmentio/kafka-go"
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

		// Get an AWS-loaded SASL mechanism:
		saslMechanism, err := AWSSaslMechanismV1()
		if err != nil {
			return nil, err
		}

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
