package pubkey

import (
	"github.com/spf13/cobra"
)

func (c *Commands) newSshCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "age",
		Short: "Show SSH public key",
		Long:  sshDescription,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.runSsh(args)
		},
	}

	return cmd
}

func (c *Commands) runSsh(args []string) error {
	return nil
}

const sshDescription = `
Show SSH public key

`
