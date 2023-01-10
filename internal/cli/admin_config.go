package cli

import (
	"github.com/segmentio/kafka-go"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

func (cli *CLI) initAdminConfig() {
	cli.SetCommand("adminConfig", "admin", cli.adminConfigCommand())
	cli.SetCommand("adminConfigMetadata", "adminConfig", cli.adminConfigMetadataCommand())
	cli.SetCommand("adminConfigParams", "adminConfig", cli.adminConfigParamsCommand())
}

// adminConfigCommand deals with config:
func (cli *CLI) adminConfigCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "Work with cluster config",
	}
}

// adminConfigMetadataCommand deals with cluster metadata:
func (cli *CLI) adminConfigMetadataCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "metadata",
		Short: "Retrieve cluster metadata",
		Run: func(cmd *cobra.Command, args []string) {

			// Config:
			cli.logger.
				WithField("sasl", cli.config.Kafka.SaslMechanism).
				WithField("security", cli.config.Kafka.SecurityProtocol).
				WithField("servers", cli.config.Kafka.BootstrapServers).
				WithField("username", cli.config.Kafka.Username).
				Debugf("Retrieving cluster metadata")

			// Retrieve cluster metadata:
			kafkaMetadata, err := cli.adminClient.Metadata(context.TODO(), &kafka.MetadataRequest{})
			if err != nil {
				cli.logger.WithError(err).Fatal("Unable to retrieve cluster metadata")
			}

			// List the general cluster metadata:
			cli.logger.
				WithField("cluster_id", kafkaMetadata.ClusterID).
				WithField("controller_id", kafkaMetadata.Controller.ID).
				WithField("throttle_duration", kafkaMetadata.Throttle.String()).
				Info("Cluster")

			// List the brokers:
			for _, broker := range kafkaMetadata.Brokers {
				cli.logger.
					WithField("host", broker.Host).
					WithField("id", broker.ID).
					WithField("port", broker.Port).
					WithField("rack", broker.Rack).
					Info("Broker")
			}
		},
	}
}

// adminConfigParamsCommand deals with cluster config parameters:
func (cli *CLI) adminConfigParamsCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Short: "Retrieve cluster params",
		Run: func(cmd *cobra.Command, args []string) {

			// Config:
			cli.logger.
				WithField("sasl", cli.config.Kafka.SaslMechanism).
				WithField("security", cli.config.Kafka.SecurityProtocol).
				WithField("servers", cli.config.Kafka.BootstrapServers).
				WithField("username", cli.config.Kafka.Username).
				Debugf("Retrieving cluster parameters")

			// Retrieve topic config:
			response, err := cli.adminClient.DescribeConfigs(
				context.TODO(),
				&kafka.DescribeConfigsRequest{
					Resources: []kafka.DescribeConfigRequestResource{
						{
							ResourceName: "any",
							ResourceType: kafka.ResourceTypeAny,
						},
					},
					IncludeSynonyms:      true,
					IncludeDocumentation: true,
				},
			)
			if err != nil {
				cli.logger.WithError(err).Fatal("Unable to retrieve config parameters")
			}

			// List the config:
			for _, resource := range response.Resources {
				for _, configEntry := range resource.ConfigEntries {
					cli.logger.
						WithField("name", configEntry.ConfigName).
						WithField("value", configEntry.ConfigValue).
						Infof("Config parameter [%s]", resource.ResourceName)
				}
			}
		},
	}
}
