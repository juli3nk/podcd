package main

import (
	"github.com/spf13/cobra"
)

const appName = "podcdd"

func newCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   appName,
		Short: "PodCD daemon",
		Long:  "PodCD daemon - starts the main daemon with Unix socket",
		Run:   runDaemon,
	}

	return cmd
}
