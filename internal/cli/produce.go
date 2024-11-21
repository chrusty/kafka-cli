package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/segmentio/kafka-go"
	"github.com/spf13/cobra"
)

func (cli *CLI) initProduce() {
	produceCommand := cli.produceCommand()
	cli.SetCommand("produce", "root", produceCommand)
}

// consumeCommand deals with producing to topics:
func (cli *CLI) produceCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:        "produce <topic> | product <topic> <payload>",
		Short:      "Produce messages to a topic",
		Args:       cobra.MinimumNArgs(1),
		ArgAliases: []string{"topic"},
		Run: func(cmd *cobra.Command, args []string) {

			// Get the topic name:
			topicName := args[0]

			// User-defined payload:
			payload := strings.Join(args[1:], " ")

			// Get the iterations flag:
			iterations, err := cmd.Flags().GetInt("iterations")
			if err != nil {
				cli.logger.WithError(err).WithField("flag", "iterations").Fatal("Unable to get flag")
			}

			// Config:
			cli.logger.
				WithField("iterations", iterations).
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

			// If we were given a custom payload:
			if len(payload) > 0 {

				// Build up a message:
				message := kafka.Message{
					Value: []byte(payload),
				}

				// Produce a message:
				cli.logger.WithField("topic", topicName).WithField("message", string(message.Value)).Trace("Producing message")
				if err := producer.WriteMessages(context.TODO(), message); err != nil {
					cli.logger.WithError(err).WithField("topic", topicName).Fatal("Unable to produce message")
				}

				return
			}

			// Generate messages:
			for i := 0; i < iterations; i++ {

				// Build up a message:
				message := kafka.Message{
					Value: []byte(fmt.Sprintf("Kafka-CLI message (%d)", i)),
				}

				// Produce a message:
				cli.logger.WithField("topic", topicName).WithField("message", string(message.Value)).Trace("Producing message")
				if err := producer.WriteMessages(context.TODO(), message); err != nil {
					cli.logger.WithError(err).WithField("topic", topicName).Fatal("Unable to produce message")
				}
			}
		},
	}

	cmd.Flags().IntP("iterations", "i", 1, "Number of iterations to send (generated payloads include an incrementing number)")

	return cmd
}
