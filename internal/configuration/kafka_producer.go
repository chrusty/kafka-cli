package configuration

import (
	"crypto/tls"
	"fmt"
	"time"

	"github.com/chrusty/kafka-cli/internal/types"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
	"github.com/sirupsen/logrus"
)

// Producer returns a Kafka Producer based on our config:
func (kc *KafkaConfig) Producer(logger *logrus.Logger, topicName string) (*kafka.Writer, error) {

	// Prepare a reader config:
	writerConfig := kafka.WriterConfig{
		Brokers:     kc.BootstrapServers,
		ErrorLogger: &kafkaErrorLogger{logger: logger},
		Logger:      &kafkaLogger{logger: logger},
		Topic:       topicName,
	}

	// Prepare a dialer (with our custom auth settings):
	dialer := &kafka.Dialer{
		Timeout: 10 * time.Second,
	}

	switch kc.SecurityProtocol {

	case types.SecProtocolPlaintext:

	case types.SecProtocolAWSMSKIAM:

		// Get an AWS-loaded SASL mechanism:
		saslMechanism, err := AWSSaslMechanismV1()
		if err != nil {
			return nil, err
		}

		// Add it to our dialer:
		dialer.SASLMechanism = saslMechanism
		dialer.TLS = &tls.Config{}

	case types.SecProtocolSSL:

		// Configure our dialer to use TLS:
		dialer.TLS = &tls.Config{}

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

		// Add it to our dialer:
		dialer.SASLMechanism = saslMechanism

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

		// Configure our dialer to use SCRAM and TLS:
		dialer.SASLMechanism = saslMechanism
		dialer.TLS = &tls.Config{}

	default:
		return nil, fmt.Errorf("unsupported security protocol %s", kc.SecurityProtocol)
	}

	// Put a writer together with our config:
	writerConfig.Dialer = dialer
	writer := kafka.NewWriter(writerConfig)
	return writer, nil
}
