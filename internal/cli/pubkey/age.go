package pubkey

import (
	"fmt"

	"github.com/juli3nk/podcd/internal/ipc"
	"github.com/spf13/cobra"
)

func (c *Commands) newAgeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "age",
		Short: "Show Age public key",
		Long:  ageDescription,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.runAge(args)
		},
	}

	return cmd
}

func (c *Commands) runAge(args []string) error {
	resp, err := c.client.Send(ipc.Request{
		Command: "status",
	})
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", resp.Data)

	return nil
}

const ageDescription = `
Show Age public key

`
