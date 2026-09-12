package main

import (
	"github.com/juli3nk/go-utils"
	"github.com/juli3nk/podcd/internal/config"
	"github.com/juli3nk/podcd/internal/ipc"
)

func main() {
	userMode := config.IsUserMode()
	paths := config.DefaultPaths(userMode)

	client := ipc.NewClient(paths.Socket)

	cmd := newCommand(client)

    if err := cli.Execute(); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
}
