package cli

import (
	"github.com/juli3nk/podcd/internal/cli/pubkey"
	"github.com/juli3nk/podcd/internal/ipc"
	"github.com/spf13/cobra"
)

func newCommand(client *ipc.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "podcd",
		Short: "PodCD CLI",
		Long:  "PodCD CLI",
	}

	cmd.AddCommand(pubkey.NewCommand(client))
	cmd.AddCommand(NewReqsCommand(client))
	cmd.AddCommand(NewVersionCommand(client))

	return cmd
}
