package cli

import (
	"fmt"

	"github.com/juli3nk/podcd/internal/ipc"
	"github.com/spf13/cobra"
)

func NewStatusCommand(client *ipc.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show PodCD status",
		Long:  statusDescription,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runStatus(client)
		},
	}
}

func runStatus(client *ipc.Client) error {
	resp, err := client.Send(ipc.Request{
		Command: "status",
	})
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", resp.Data)

	return nil
}

const statusDescription = `
Show PodCD status

`
