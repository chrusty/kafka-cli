package cli

import (
	"sync"

	"github.com/chrusty/kafka-cli/internal/configuration"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// CLI contains our dependencies:
type CLI struct {
	adminClient *kafka.Client
	commands    map[string]*cobra.Command
	config      *configuration.Config
	logger      *logrus.Logger
	mutex       sync.Mutex
}

// New returns a configured CLI command:
func New(config *configuration.Config, logger *logrus.Logger) *CLI {

	// Get an admin client:
	adminClient, err := config.Kafka.Admin(logger)
	if err != nil {
		logger.WithError(err).Fatal("Unable to prepare a Kafka admin client")
	}

	// Make a new CLI:
	c := &CLI{
		adminClient: adminClient,
		commands:    make(map[string]*cobra.Command),
		config:      config,
		logger:      logger,
	}

	// Add a root command:
	rootCmd := c.rootCommand()
	c.commands["root"] = rootCmd

	// Add subcommands:
	c.initAdmin()
	c.initAdminGroups()
	c.initAdminTopics()

	return c
}

// Execute the command:
func (cli *CLI) Execute() error {
	return cli.GetCommand("root").Execute()
}

// Retrieve a command (by name):
func (cli *CLI) GetCommand(name string) *cobra.Command {
	cli.mutex.Lock()
	defer cli.mutex.Unlock()

	// Check if the command exists, and return it:
	if command, ok := cli.commands[name]; ok {
		return command
	}

	cli.logger.Fatalf("Unable to find command %s", name)
	return nil
}

// Set a command (by name):
func (cli *CLI) SetCommand(name, parent string, cmd *cobra.Command) {

	cli.logger.Tracef("Adding command %s as a child of %s", name, parent)

	// Add the command to its parent:
	cli.GetCommand(parent).AddCommand(cmd)

	// Set the command in our map:
	cli.mutex.Lock()
	defer cli.mutex.Unlock()
	cli.commands[name] = cmd
}

// rootCommand is the CLI entry-point:
func (cli *CLI) rootCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "kafka-cli",
		Short: "CLI tools to work with Kafka",
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}
}
