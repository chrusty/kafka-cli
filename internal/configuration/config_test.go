package configuration

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {

	// Load config from env-vars:
	_, err := Load(logrus.New())
	assert.NoError(t, err, "Error while loading config")
}
