package main

import (
	"github.com/chrusty/kafka-cli/internal/cli"
	"github.com/chrusty/kafka-cli/internal/configuration"
	"github.com/sirupsen/logrus"
)

func main() {

	// Use a Logrus logger:
	logger := logrus.New()

	// Load config from env-vars:
	config, err := configuration.Load(logger)
	if err != nil {
		logger.WithError(err).Fatal("Unable to load config")
	}

	// Prepare a new CLI:
	cli := cli.New(config, logger)
	if err != nil {
		logger.WithError(err).Fatal("Unable to prepare a new CLI")
	}

	// Execute the root command:
	if err := cli.Execute(); err != nil {
		logger.WithError(err).Trace("Unable to execute root command")
	}
}
