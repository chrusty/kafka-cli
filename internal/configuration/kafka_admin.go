package configuration

import (
	"crypto/tls"
	"fmt"

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

	case "PLAINTEXT":

		client.Transport = &kafka.Transport{}

	case "SSL":

		client.Transport = &kafka.Transport{
			TLS: &tls.Config{},
		}

	case "SASL_PLAINTEXT":

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

	case "SASL_SSL":

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
		return nil, fmt.Errorf("Unsupported security protocol %s", kc.SecurityProtocol)
	}

	return client, nil
}
