package configuration

import (
	"fmt"

	"github.com/caarlos0/env"
	"github.com/sirupsen/logrus"
)

// Config for as many generic deps as we can handle:
type Config struct {
	Kafka   KafkaConfig
	Logging LoggingConfig
}

// Load prepares a new config and populates it from environment variables:
func Load(logger *logrus.Logger) (*Config, error) {
	newConfig := &Config{}

	for _, configSection := range []interface{}{
		newConfig,
		&newConfig.Kafka,
		&newConfig.Logging,
	} {
		if err := env.Parse(configSection); err != nil {
			return nil, fmt.Errorf("unable to load the config: %v", err)
		}
	}

	// Parse and set the log-level config in the given logger:
	newConfig.Logging.ParseLevel(logger)

	return newConfig, nil
}
