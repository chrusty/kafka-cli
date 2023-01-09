package cli

import "github.com/spf13/cobra"

func (cli *CLI) initAdmin() {
	cli.SetCommand("admin", "root", cli.adminCommand())
}

func (cli *CLI) adminCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "admin",
		Short: "Administration of topics",
	}
}
