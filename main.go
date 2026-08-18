// Vorpalstacks is an AWS-compatible cloud platform for edge and on-premises environments.
//
// It provides the supported AWS service APIs (authoritative list and count
// in docs/services.md) with a single binary, using CockroachDB Pebble for
// persistent storage and supporting both JSON and Query AWS API protocols.
package main

import (
	"embed"
	"fmt"
	"os"
	"strings"

	"vorpalstacks/internal/config"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/server/apps"
)

//go:embed webconsole/dist
var webconsoleFS embed.FS

func main() {
	dataPath := os.Getenv("DATA_PATH")
	if dataPath == "" {
		dataPath = "./data"
	}

	storageMgr, err := storage.NewRegionStorageManager(&storage.Config{Path: dataPath})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialise storage: %v\n", err)
		os.Exit(1)
	}

	globalStorage, err := storageMgr.GetGlobalStorage()
	if err != nil {
		storageMgr.Close()
		fmt.Fprintf(os.Stderr, "Failed to get global storage: %v\n", err)
		os.Exit(1)
	}

	config.Initialise(globalStorage)

	bc := config.LoadBootstrapConfig()
	cfg := apps.FromBootstrap(bc)
	cfg.ConsoleAssets = webconsoleFS

	app, err := apps.NewWithStorage(cfg, storageMgr)
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
