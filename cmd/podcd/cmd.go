package main

import (
	"github.com/juli3nk/podcd/internal/cli"
	"github.com/juli3nk/podcd/internal/cli/pubkey"
	"github.com/juli3nk/podcd/internal/ipc"
	"github.com/spf13/cobra"
)

const appName = "podcd"

func newCommand(client *ipc.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   appName,
		Short: "PodCD CLI",
		Long:  "PodCD CLI",
	}

	cmd.AddCommand(pubkey.NewCommand(client))
	cmd.AddCommand(cli.NewReqsCommand(client))
	cmd.AddCommand(cli.NewVersionCommand(client))

	return cmd
}
