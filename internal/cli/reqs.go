package cli

import (
	"fmt"

	"github.com/juli3nk/podcd/internal/ipc"
	"github.com/spf13/cobra"
)

func NewReqsCommand(client *ipc.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "reqs",
		Short: "Show PodCD requirements",
		Long:  reqsDescription,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runReqs(client)
		},
	}
}

func runReqs(client *ipc.Client) error {
	resp, err := client.Send(ipc.Request{
		Command: "reqs",
	})
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", resp.Data)

	return nil
}

const reqsDescription = `
Show PodCD requirements

`
