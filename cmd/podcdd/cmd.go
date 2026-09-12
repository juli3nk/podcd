package main

import (
	"github.com/spf13/cobra"
)

func newCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "podcdd",
		Short: "PodCD daemon",
		Long:  "PodCD daemon - starts the main daemon with Unix socket",
		Run:   runDaemon,
	}

	return cmd
}
