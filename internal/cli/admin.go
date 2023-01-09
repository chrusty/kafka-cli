package cli

import "github.com/spf13/cobra"

func (cli *CLI) initAdmin() {
	adminCommand := cli.adminCommand()
	cli.SetCommand("admin", "root", adminCommand)
}

func (cli *CLI) adminCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "admin",
		Short: "Administration of topics",
	}
}
