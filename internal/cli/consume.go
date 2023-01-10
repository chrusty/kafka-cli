package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

func (cli *CLI) initConsume() {
	consumeCommand := cli.consumeCommand()
	consumeCommand.PersistentFlags().String("groupid", "", "Consumer group ID (if blank then groups won't be used, offsets won't be committed)")
	cli.SetCommand("consume", "root", consumeCommand)
}

// consumeCommand deals with consuming from topics:
func (cli *CLI) consumeCommand() *cobra.Command {

	return &cobra.Command{
		Use:        "consume <topic>",
		Short:      "Consume messages from a topic",
		Args:       cobra.ExactArgs(1),
		ArgAliases: []string{"topic"},
		Run: func(cmd *cobra.Command, args []string) {

			// Get the autoCommit flag:
			groupId, err := cmd.Flags().GetString("groupid")
			if err != nil {
				cli.logger.WithError(err).WithField("flag", "groupid").Fatal("Unable to get flag")
			}

			// Get the topic name:
			topicName := args[0]

			// Config:
			cli.logger.
				WithField("sasl", cli.config.Kafka.SaslMechanism).
				WithField("security", cli.config.Kafka.SecurityProtocol).
				WithField("servers", cli.config.Kafka.BootstrapServers).
				WithField("username", cli.config.Kafka.Username).
				Debugf("Consuming topic: %s", topicName)

			// Get an admin client:
			consumer, err := cli.config.Kafka.Consumer(cli.logger, groupId, topicName)
			if err != nil {
				cli.logger.WithError(err).WithField("topic", topicName).WithField("group", groupId).Fatal("Unable to prepare a consumer")
			}

			// Periodically report our progress:
			startTime := time.Now()
			var totalErrors, messagesConsumed int32
			go func() {
				for {
					time.Sleep(time.Second)
					cli.logger.WithField("errors", totalErrors).WithField("messages", messagesConsumed).WithField("messages/s", messagesConsumed/int32(time.Since(startTime).Seconds())).Info("Progress report")
				}
			}()

			// Consume:
			for {

				// Get a message:
				message, err := consumer.ReadMessage(context.TODO())
				if err != nil {
					cli.logger.WithError(err).Error("Unable to consume a message")
					totalErrors++
					continue
				}
				messagesConsumed++

				// Log the message metadata:
				cli.logger.
					WithField("headers", message.Headers).
					WithField("key", string(message.Key)).
					WithField("offset", message.Offset).
					WithField("partition", message.Partition).
					Debug("Got a message")

				// Print the payload:
				fmt.Println(string(message.Value))
			}
		},
	}
}
