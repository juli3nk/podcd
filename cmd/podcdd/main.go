package main

import (
	"github.com/juli3nk/go-utils"
)

func main() {
	cmd := newCommand()

    if err := cmd.Execute(); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
}
