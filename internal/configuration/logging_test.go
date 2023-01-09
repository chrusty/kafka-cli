package configuration

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestLoggingConfig(t *testing.T) {

	// Make a Logrus logger:
	logger := logrus.New()

	// Override one of the defaults:
	err := os.Setenv("LOGGING_LEVEL", "trace")
	assert.NoError(t, err, "Couldn't set an env var")

	// Load config from env-vars:
	testConfig, err := Load(logger)
	assert.NoError(t, err, "Error while loading config")

	// Make sure our logging config was picked up:
	assert.Equal(t, "trace", testConfig.Logging.Level)

	// Make sure our logging config was applied to the logger we provided:
	assert.Equal(t, logrus.TraceLevel, logger.Level)
}
