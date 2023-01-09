package configuration

import (
	"github.com/sirupsen/logrus"
)

// LoggingConfig configures logging:
type LoggingConfig struct {
	Level      string `env:"LOGGING_LEVEL" envDefault:"INFO"`
	Timestamps bool   `env:"LOGGING_TIMESTAMPS" envDefault:"true"`
}

// ParseLevel sets the configured logging level with the provided logger:
func (l *LoggingConfig) ParseLevel(logger *logrus.Logger) {

	// Set the logging level specified in the config:
	loggingLevel, err := logrus.ParseLevel(l.Level)
	if err != nil {
		logger.WithError(err).Warn("Invalid log level")
		return
	}
	logger.SetLevel(loggingLevel)
}
