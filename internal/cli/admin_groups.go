package cli

import (
	"context"

	"github.com/segmentio/kafka-go"
	"github.com/spf13/cobra"
)

func (cli *CLI) initAdminGroups() {
	cli.SetCommand("adminGroups", "admin", cli.adminGroupsCommand())
	cli.SetCommand("adminGroupsDelete", "adminGroups", cli.adminGroupsDeleteCommand())
	cli.SetCommand("adminGroupsDescribe", "adminGroups", cli.adminGroupsDescribeCommand())
	cli.SetCommand("adminGroupsList", "adminGroups", cli.adminGroupsListCommand())
}

// adminGroupsCommand deals with groups:
func (cli *CLI) adminGroupsCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "groups",
		Short: "Work with consumer groups",
	}
}

// adminGroupsDeleteCommand deals with deleting groups:
func (cli *CLI) adminGroupsDeleteCommand() *cobra.Command {

	return &cobra.Command{
		Use:   "delete <group>",
		Short: "Delete a group",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			// Get the groupId:
			groupId := args[0]

			// Config:
			cli.logger.
				WithField("sasl", cli.config.Kafka.SaslMechanism).
				WithField("security", cli.config.Kafka.SecurityProtocol).
				WithField("servers", cli.config.Kafka.BootstrapServers).
				WithField("username", cli.config.Kafka.Username).
				Debugf("Describing group: %s", groupId)

			// Delete the group:
			response, err := cli.adminClient.OffsetDelete(
				context.TODO(),
				&kafka.OffsetDeleteRequest{
					GroupID: groupId,
				},
			)
			if err != nil {
				cli.logger.WithError(err).WithField("group", groupId).Fatal("Unable to delete group")
			}

			cli.logger.WithError(response.Error).WithField("group", groupId).Info("Group deleted")
		},
	}
}

// adminGroupsDescribeCommand deals with describing groups:
func (cli *CLI) adminGroupsDescribeCommand() *cobra.Command {

	return &cobra.Command{
		Use:   "describe <group>",
		Short: "Describe a group",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			// Get the groupId:
			groupId := args[0]

			// Config:
			cli.logger.
				WithField("sasl", cli.config.Kafka.SaslMechanism).
				WithField("security", cli.config.Kafka.SecurityProtocol).
				WithField("servers", cli.config.Kafka.BootstrapServers).
				WithField("username", cli.config.Kafka.Username).
				Debugf("Describing group: %s", groupId)

			// Retrieve group config:
			response, err := cli.adminClient.DescribeGroups(
				context.TODO(),
				&kafka.DescribeGroupsRequest{
					GroupIDs: []string{groupId},
				},
			)
			if err != nil {
				cli.logger.WithError(err).Fatal("Unable to retrieve group config")
			}

			// List the config:
			for _, group := range response.Groups {
				cli.logger.WithError(group.Error).
					WithField("members", len(group.Members)).
					WithField("state", group.GroupState).
					Infof("Group config [%s]", group.GroupID)
			}
		},
	}
}

// adminGroupsListCommand deals with listing groups:
func (cli *CLI) adminGroupsListCommand() *cobra.Command {

	return &cobra.Command{
		Use:   "list",
		Short: "List groups",
		Run: func(cmd *cobra.Command, args []string) {

			// Config:
			cli.logger.
				WithField("sasl", cli.config.Kafka.SaslMechanism).
				WithField("security", cli.config.Kafka.SecurityProtocol).
				WithField("servers", cli.config.Kafka.BootstrapServers).
				WithField("username", cli.config.Kafka.Username).
				Debug("Listing groups")

			// Retrieve groups:
			response, err := cli.adminClient.ListGroups(context.TODO(), &kafka.ListGroupsRequest{})
			if err != nil {
				cli.logger.WithError(err).Fatal("Unable to retrieve cluster metadata")
			}

			// List the topics:
			for _, group := range response.Groups {
				cli.logger.WithField("id", group.GroupID).Info("Group")
			}
		},
	}
}
