package pubkey

import (
	"github.com/juli3nk/podcd/internal/ipc"
	"github.com/spf13/cobra"
)

type Commands struct {
	client *ipc.Client
}

func NewCommand(client *ipc.Client) *cobra.Command {
	c := &Commands{
		client: client,
	}

	cmd := &cobra.Command{
		Use:   "pubkey",
		Short: "Show public keys",
		Long:  pubkeyDescription,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Usage()
		},
	}

	cmd.AddCommand(
		c.newAgeCommand(),
		c.newSshCommand(),
	)

	return cmd
}

const pubkeyDescription = `
The **podcd pubkey** command has subcommands for showing public keys.

To see help for a subcommand, use:

    podcd pubkey [command] --help

`
