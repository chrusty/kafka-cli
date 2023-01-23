package configuration

import (
	"crypto/tls"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
	"github.com/sirupsen/logrus"
)

// Consumer returns a Kafka Consumer based on our config:
func (kc *KafkaConfig) Consumer(logger *logrus.Logger, groupId, topicName string) (*kafka.Reader, error) {

	// Prepare a reader config:
	readerConfig := kafka.ReaderConfig{
		Brokers:        kc.BootstrapServers,
		CommitInterval: 5 * time.Second,
		ErrorLogger:    &kafkaErrorLogger{logger: logger},
		Logger:         &kafkaLogger{logger: logger},
		Topic:          topicName,
	}

	// Add the groupId if one was provided:
	if groupId != "" {
		readerConfig.GroupID = groupId
	}

	// Prepare a dialer (with our custom auth settings):
	dialer := &kafka.Dialer{
		ClientID: "kafka-cli",
		Timeout:  10 * time.Second,
	}

	switch kc.SecurityProtocol {

	case "PLAINTEXT":

	case "SSL":

		// Configure our dialer to use TLS:
		dialer.TLS = &tls.Config{}

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

		// Add it to our dialer:
		dialer.SASLMechanism = saslMechanism

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

		// Configure our dialer to use SCRAM and TLS:
		dialer.SASLMechanism = saslMechanism
		dialer.TLS = &tls.Config{}

	default:
		return nil, fmt.Errorf("Unsupported security protocol %s", kc.SecurityProtocol)
	}

	// Put a reader together with our config:
	readerConfig.Dialer = dialer
	reader := kafka.NewReader(readerConfig)
	return reader, nil
}
