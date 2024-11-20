package cli

import (
	"context"

	"github.com/segmentio/kafka-go"
	"github.com/spf13/cobra"
)

func (cli *CLI) initProduce() {
	produceCommand := cli.produceCommand()
	cli.SetCommand("produce", "root", produceCommand)
}

// consumeCommand deals with producing to topics:
func (cli *CLI) produceCommand() *cobra.Command {

	return &cobra.Command{
		Use:        "produce <topic>",
		Short:      "Produce messages to a topic",
		Args:       cobra.ExactArgs(1),
		ArgAliases: []string{"topic"},
		Run: func(cmd *cobra.Command, args []string) {

			// Get the topic name:
			topicName := args[0]

			// Config:
			cli.logger.
				WithField("sasl", cli.config.Kafka.SaslMechanism).
				WithField("security", cli.config.Kafka.SecurityProtocol).
				WithField("servers", cli.config.Kafka.BootstrapServers).
				WithField("username", cli.config.Kafka.Username).
				Debugf("Producing to topic: %s", topicName)

			// Get a producer client:
			producer, err := cli.config.Kafka.Producer(cli.logger, topicName)
			if err != nil {
				cli.logger.WithError(err).WithField("topic", topicName).Fatal("Unable to prepare a producer")
			}

			// Build up a message:
			message := kafka.Message{
				Value: []byte("cruft"),
			}

			// Produce a message:
			if err := producer.WriteMessages(context.TODO(), message); err != nil {
				cli.logger.WithError(err).WithField("topic", topicName).Fatal("Unable to produce message")
			}
		},
	}
}
