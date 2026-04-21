package vstackscli

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"vorpalstacks/internal/core/storage"
	storeconfig "vorpalstacks/internal/store/config"
)

// RunConfig dispatches application configuration commands (get, set, reset,
// list, schema) against the app_config store.
//
// Parameters:
//   - store: The underlying PebbleDB storage
//   - args: Command-line arguments (sub-command and flags)
func RunConfig(store storage.BasicStorage, args []string) {
	if len(args) == 0 {
		printConfigUsage()
		os.Exit(1)
	}

	configStore := storeconfig.NewStore(store)
	cmd := args[0]
	cmdArgs := args[1:]

	switch cmd {
	case "get":
		configGetCmd(configStore, cmdArgs)
	case "set":
		configSetCmd(configStore, cmdArgs)
	case "reset":
		configResetCmd(configStore, cmdArgs)
	case "list":
		configListCmd(configStore, cmdArgs)
	case "schema":
		configSchemaCmd()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", cmd)
		printConfigUsage()
		os.Exit(1)
	}
}

func printConfigUsage() {
	fmt.Println("Config Commands (app_config):")
	fmt.Println("  vstacks config get     <key>")
	fmt.Println("  vstacks config set     <key> <value>")
	fmt.Println("  vstacks config reset   <key>")
	fmt.Println("  vstacks config list    [-category <cat>]")
	fmt.Println("  vstacks config schema")
}

func configGetCmd(store *storeconfig.Store, args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Error: key is required")
		os.Exit(1)
	}

	entry, err := store.Get(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get config: %v\n", err)
		os.Exit(1)
	}

	printConfigEntry(entry)
}

func configSetCmd(store *storeconfig.Store, args []string) {
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "Error: key and value are required")
		os.Exit(1)
	}

	key := args[0]
	valueStr := args[1]

	entry, err := store.Get(key)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unknown config key: %s\n", key)
		os.Exit(1)
	}

	var value interface{}
	switch entry.Type {
	case storeconfig.ConfigTypeBool:
		value = valueStr == "true" || valueStr == "1" || valueStr == "yes"
	case storeconfig.ConfigTypeInt, storeconfig.ConfigTypePort:
		var i int
		for _, c := range valueStr {
			if c < '0' || c > '9' {
				fmt.Fprintf(os.Stderr, "Invalid integer value: %s\n", valueStr)
				os.Exit(1)
			}
			i = i*10 + int(c-'0')
		}
		value = i
	default:
		value = valueStr
	}

	if err := store.Set(key, value); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to set config: %v\n", err)
		os.Exit(1)
	}

	updated, _ := store.Get(key)
	fmt.Printf("Config updated:\n")
	printConfigEntry(updated)
}

func configResetCmd(store *storeconfig.Store, args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Error: key is required")
		os.Exit(1)
	}

	entry, err := store.Reset(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to reset config: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Config reset to default:\n")
	printConfigEntry(entry)
}

func configListCmd(store *storeconfig.Store, args []string) {
	fs := flag.NewFlagSet("list", flag.ExitOnError)
	category := fs.String("category", "", "Filter by category")
	_ = fs.Parse(args)

	var entries []*storeconfig.ConfigEntry
	var err error

	if *category != "" {
		entries, err = store.ListByCategory(storeconfig.ConfigCategory(*category))
	} else {
		entries, err = store.GetAll()
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to list config: %v\n", err)
		os.Exit(1)
	}

	for _, entry := range entries {
		printConfigEntry(entry)
		fmt.Println()
	}
}

func configSchemaCmd() {
	schema := storeconfig.GetSchema()
	for _, s := range schema {
		readOnly := ""
		if s.ReadOnly {
			readOnly = " [read-only]"
		}
		fmt.Printf("%-35s %-8s %s%s\n", s.Key, s.Type, s.Description, readOnly)
		if s.EnvVar != "" {
			fmt.Printf("  env: %s  default: %v\n", s.EnvVar, s.Default)
		}
	}
}

func printConfigEntry(entry *storeconfig.ConfigEntry) {
	fmt.Printf("  Key:         %s\n", entry.Key)
	fmt.Printf("  Value:       %v\n", entry.Value)
	fmt.Printf("  Type:        %s\n", entry.Type)
	fmt.Printf("  Source:      %s\n", entry.Source)
	fmt.Printf("  Category:    %s\n", entry.Category)
	if entry.Description != "" {
		fmt.Printf("  Description: %s\n", entry.Description)
	}
	if entry.EnvVar != "" {
		fmt.Printf("  EnvVar:      %s\n", entry.EnvVar)
	}
}

func printServiceConfig(cfg *serviceConfigInfo) {
	fmt.Printf("  Service:  %s\n", cfg.ServiceName)
	fmt.Printf("  Mode:     %s\n", cfg.Mode)
	fmt.Printf("  Enabled:  %v\n", cfg.Enabled)
}

type serviceConfigInfo struct {
	ServiceName string `json:"service_name"`
	Mode        string `json:"mode"`
	Enabled     bool   `json:"enabled"`
}

func formatServiceConfigJSON(serviceName, mode string, enabled bool) string {
	data, _ := json.MarshalIndent(serviceConfigInfo{
		ServiceName: serviceName,
		Mode:        mode,
		Enabled:     enabled,
	}, "  ", "  ")
	return string(data)
}
