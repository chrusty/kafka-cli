package configuration

import (
	"github.com/chrusty/kafka-cli/internal/types"
	"github.com/segmentio/kafka-go/sasl/scram"
	"github.com/sirupsen/logrus"
)

var defaultSaslAlgorithm = scram.SHA512

// KafkaConfig configures the Kafka client:
type KafkaConfig struct {
	AWSRegion        string   `env:"KAFKA_AWSREGION" envDefault:"ap-southeast-2"`
	BootstrapServers []string `env:"KAFKA_BOOTSTRAPSERVERS" envDefault:"localhost:9092"`
	IAMAuth          bool     `env:"KAFKA_IAMAUTH" envDefault:"false"`               // Set this to true to enable IAM auth with SASL/SCRAM
	Password         string   `env:"KAFKA_PASSWORD"`                                 // SASL/SCRAM password
	RequiredAcks     int      `env:"KAFKA_REQUIREDACKS" envDefault:"1"`              // Required ACKS [1,2]
	SaslMechanism    string   `env:"KAFKA_SASLMECHANISM" envDefault:"SCRAM-SHA-512"` // [SCRAM-SHA-256, SCRAM-SHA-512]
	SecurityProtocol string   `env:"KAFKA_SECURITYPROTOCOL" envDefault:"PLAINTEXT"`  // [AWS_MSK_IAM, SASL_SSL, SASL_PLAINTEXT, SSL, PLAINTEXT]
	Username         string   `env:"KAFKA_USERNAME"`                                 // SASL/SCRAM username
}

// kafkaLogger implements the kafka.Logger interface:
type kafkaLogger struct {
	logger *logrus.Logger
}

// Printf regurgitates the log as DEBUG:
func (kl *kafkaLogger) Printf(format string, args ...interface{}) {
	kl.logger.Debugf(format, args...)
}

// kafkaErrorLogger implements the kafka.Logger interface:
type kafkaErrorLogger struct {
	logger *logrus.Logger
}

// Printf regurgitates the log as WARN:
func (kel *kafkaErrorLogger) Printf(format string, args ...interface{}) {
	kel.logger.Warnf(format, args...)
}

func (kc *KafkaConfig) saslAlgorithm(logger *logrus.Logger) scram.Algorithm {
	switch kc.SaslMechanism {
	case types.SaslMechanismScramSHA256:
		return scram.SHA256
	case types.SaslMechanismScramSHA512:
		return scram.SHA512
	default:
		logger.Warnf("Unsupported SASL mechanism (%s), assuming default (%s)", kc.SaslMechanism, defaultSaslAlgorithm)
		return defaultSaslAlgorithm
	}
}
