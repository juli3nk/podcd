package cli

import (
	"encoding/json"
	"fmt"

	"github.com/juli3nk/podcd/internal/ipc"
	"github.com/juli3nk/podcd/internal/version"
	"github.com/spf13/cobra"
)

func NewVersionCommand(client *ipc.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show PodCD version",
		Long:  versionDescription,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runVersion(client)
		},
	}
}

func runVersion(client *ipc.Client) error {
	clientVersion := version.New()

	fmt.Println("Client:")
	clientVersion.Show()

	fmt.Println("Server:")
	resp, err := client.Send(ipc.Request{
		Command: "version",
	})
	if err != nil {
		fmt.Printf("%v\n", err)
	} else {
		dataBytes, _ := json.Marshal(resp.Data)
		var serverVersion version.VersionInfo
		_ = json.Unmarshal(dataBytes, &serverVersion)

		serverVersion.Show()
	}

	return nil
}

const versionDescription = `
Show PodCD version

`
