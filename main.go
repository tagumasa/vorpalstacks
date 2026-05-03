// Vorpalstacks is an AWS-compatible cloud platform for edge and on-premises environments.
//
// It provides 32 service APIs covering 27 AWS services with a single binary,
// using CockroachDB Pebble for persistent storage and supporting both
// JSON and Query AWS API protocols.
package main

import (
	"fmt"
	"os"
	"strings"

	"vorpalstacks/internal/config"
	"vorpalstacks/internal/server/apps"
)

func main() {
	bc := config.LoadBootstrapConfig()
	cfg := apps.FromBootstrap(bc)

	app, err := apps.New(cfg)
	if err != nil {
		if strings.Contains(err.Error(), "resource temporarily unavailable") ||
			strings.Contains(err.Error(), "lock file") {
			fmt.Fprintln(os.Stderr, "Another vorpalstacks process is running. Kill it first: pkill -9 vorpalstacks")
		}
		fmt.Fprintf(os.Stderr, "Failed to initialise: %v\n", err)
		os.Exit(1)
	}

	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}
