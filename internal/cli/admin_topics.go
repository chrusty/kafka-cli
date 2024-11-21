package cli

import (
	"github.com/segmentio/kafka-go"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

func (cli *CLI) initAdminTopics() {
	cli.SetCommand("adminTopics", "admin", cli.adminTopicsCommand())
	cli.SetCommand("adminTopicCreate", "adminTopics", cli.adminTopicCreateCommand())
	cli.SetCommand("adminTopicsDelete", "adminTopics", cli.adminTopicsDeleteCommand())
	cli.SetCommand("adminTopicsDescribe", "adminTopics", cli.adminTopicsDescribeCommand())
	cli.SetCommand("adminTopicsList", "adminTopics", cli.adminTopicsListCommand())
}

// adminTopicsCommand deals with managing topics:
func (cli *CLI) adminTopicsCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "topics",
		Short: "Work with topics",
	}
}

// adminTopicCreateCommand deals with creating topics:
func (cli *CLI) adminTopicCreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:        "create <topic>",
		Short:      "Create a topic",
		Args:       cobra.ExactArgs(1),
		ArgAliases: []string{"topic"},
		Run: func(cmd *cobra.Command, args []string) {

			// Get the topic name:
			topicName := args[0]

			// Get the partitions flag:
			numPartitions, err := cmd.Flags().GetInt("partitions")
			if err != nil {
				cli.logger.WithError(err).WithField("flag", "partitions:").Fatal("Unable to get flag")
			}

			// Get the replicas flag:
			numReplicas, err := cmd.Flags().GetInt("replicas")
			if err != nil {
				cli.logger.WithError(err).WithField("flag", "replicas:").Fatal("Unable to get flag")
			}

			// Config:
			cli.logger.
				WithField("sasl", cli.config.Kafka.SaslMechanism).
				WithField("security", cli.config.Kafka.SecurityProtocol).
				WithField("servers", cli.config.Kafka.BootstrapServers).
				WithField("partitions", numPartitions).
				WithField("replicas", numReplicas).
				Debugf("Creating topic: %s", topicName)

			// Delete the topic:
			response, err := cli.adminClient.CreateTopics(context.TODO(), &kafka.CreateTopicsRequest{
				Topics: []kafka.TopicConfig{
					{
						Topic:             topicName,
						NumPartitions:     numPartitions,
						ReplicationFactor: numReplicas,
					},
				},
			})
			if err != nil {
				cli.logger.WithError(err).WithField("topic", topicName).Fatal("Unable to create topic")
			}

			cli.logger.WithError(response.Errors[topicName]).WithField("topic", topicName).Info("Topic created")
		},
	}

	cmd.Flags().IntP("partitions", "p", 1, "Number of partitions for the new topic")
	cmd.Flags().IntP("replicas", "r", 1, "Number of replicas for the new topic")

	return cmd
}

// adminTopicDeleteCommand deals with deleting topics:
func (cli *CLI) adminTopicsDeleteCommand() *cobra.Command {
	return &cobra.Command{
		Use:        "delete <topic>",
		Short:      "Delete a topic",
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
				Debugf("Deleting topic: %s", topicName)

			// Delete the topic:
			response, err := cli.adminClient.DeleteTopics(context.TODO(), &kafka.DeleteTopicsRequest{
				Topics: []string{topicName},
			})
			if err != nil {
				cli.logger.WithError(err).WithField("topic", topicName).Fatal("Unable to delete topic")
			}

			cli.logger.WithError(response.Errors[topicName]).WithField("topic", topicName).Info("Topic deleted")
		},
	}
}

// adminTopicsDescribeCommand deals with describing topics:
func (cli *CLI) adminTopicsDescribeCommand() *cobra.Command {

	return &cobra.Command{
		Use:   "describe <topic>",
		Short: "Describe a topic",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			// Get the topic name:
			topicName := args[0]

			// Config:
			cli.logger.
				WithField("sasl", cli.config.Kafka.SaslMechanism).
				WithField("security", cli.config.Kafka.SecurityProtocol).
				WithField("servers", cli.config.Kafka.BootstrapServers).
				WithField("username", cli.config.Kafka.Username).
				Debugf("Describing topic: %s", topicName)

			// Retrieve topic config:
			response, err := cli.adminClient.DescribeConfigs(
				context.TODO(),
				&kafka.DescribeConfigsRequest{
					Resources: []kafka.DescribeConfigRequestResource{
						{
							ResourceName: topicName,
							ResourceType: kafka.ResourceTypeTopic,
						},
					},
					IncludeSynonyms:      true,
					IncludeDocumentation: true,
				},
			)
			if err != nil {
				cli.logger.WithError(err).Fatal("Unable to retrieve topic config")
			}

			// List the config:
			for _, resource := range response.Resources {
				for _, configEntry := range resource.ConfigEntries {
					cli.logger.
						WithField("name", configEntry.ConfigName).
						WithField("value", configEntry.ConfigValue).
						Infof("Topic config [%s]", resource.ResourceName)
				}
			}
		},
	}
}

// adminTopicsListCommand deals with listing topics:
func (cli *CLI) adminTopicsListCommand() *cobra.Command {

	return &cobra.Command{
		Use:   "list",
		Short: "List topics",
		Run: func(cmd *cobra.Command, args []string) {

			// Config:
			cli.logger.
				WithField("sasl", cli.config.Kafka.SaslMechanism).
				WithField("security", cli.config.Kafka.SecurityProtocol).
				WithField("servers", cli.config.Kafka.BootstrapServers).
				WithField("username", cli.config.Kafka.Username).
				Debug("Listing topics")

			// Retrieve cluster metadata:
			kafkaMetadata, err := cli.adminClient.Metadata(context.TODO(), &kafka.MetadataRequest{})
			if err != nil {
				cli.logger.WithError(err).Fatal("Unable to retrieve cluster metadata")
			}

			// List the topics:
			for _, topic := range kafkaMetadata.Topics {
				cli.logger.
					WithField("name", topic.Name).
					WithField("partitions", len(topic.Partitions)).
					WithField("internal", topic.Internal).
					Info("Topic")
			}
		},
	}
}
