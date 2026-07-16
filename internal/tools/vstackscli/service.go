package vstackscli

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/api"
)

// RunService dispatches service configuration commands (get, set-mode, enable,
// disable, list) against the service_config store.
//
// Parameters:
//   - store: The underlying PebbleDB storage
//   - args: Command-line arguments (sub-command and flags)
func RunService(store storage.BasicStorage, args []string) {
	if len(args) == 0 {
		printServiceUsage()
		os.Exit(1)
	}

	configStore := api.NewConfigStore(store)
	if err := configStore.LoadAll(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to load service configs: %v\n", err)
	}

	cmd := args[0]
	cmdArgs := args[1:]

	switch cmd {
	case "get":
		serviceGetCmd(configStore, cmdArgs)
	case "set-mode":
		serviceSetModeCmd(configStore, cmdArgs)
	case "enable":
		serviceEnableCmd(configStore, cmdArgs, true)
	case "disable":
		serviceEnableCmd(configStore, cmdArgs, false)
	case "list":
		serviceListCmd(configStore)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", cmd)
		printServiceUsage()
		os.Exit(1)
	}
}

func printServiceUsage() {
	fmt.Println("Service Commands (service_config):")
	fmt.Println("  vstacks service get       <name>")
	fmt.Println("  vstacks service set-mode  <name> -mode <IMPLEMENTED|FALLBACK|ERROR_INJECTION>")
	fmt.Println("  vstacks service enable    <name>")
	fmt.Println("  vstacks service disable   <name>")
	fmt.Println("  vstacks service list")
}

func serviceGetCmd(store *api.ConfigStore, args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Error: service name is required")
		os.Exit(1)
	}

	cfg, err := store.Get(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get service config: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(formatServiceConfigJSON(cfg.ServiceName, formatMode(int(cfg.Mode)), cfg.Enabled))
}

func serviceSetModeCmd(store *api.ConfigStore, args []string) {
	fs := flag.NewFlagSet("set-mode", flag.ExitOnError)
	mode := fs.String("mode", "", "Service mode (IMPLEMENTED, FALLBACK, ERROR_INJECTION)")
	_ = fs.Parse(args)

	serviceName := fs.Arg(0)
	if serviceName == "" || *mode == "" {
		fmt.Fprintln(os.Stderr, "Error: service name and -mode are required")
		os.Exit(1)
	}

	if err := store.SetMode(serviceName, api.ServiceMode(parseMode(*mode))); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to set mode: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Service '%s' mode set to %s\n", serviceName, strings.ToUpper(*mode))
}

func serviceEnableCmd(store *api.ConfigStore, args []string, enabled bool) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Error: service name is required")
		os.Exit(1)
	}

	serviceName := args[0]
	if err := store.SetEnabled(serviceName, enabled); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to %s service: %v\n", actionLabel(enabled), err)
		os.Exit(1)
	}

	fmt.Printf("Service '%s' %s\n", serviceName, actionLabel(enabled))
}

func actionLabel(enabled bool) string {
	if enabled {
		return "enabled"
	}
	return "disabled"
}

func serviceListCmd(store *api.ConfigStore) {
	fmt.Printf("Services (%d):\n", store.Count())
	_ = store.ForEach(func(cfg *api.ServiceConfig) error {
		fmt.Printf("  %-25s %-15s %v\n", cfg.ServiceName, formatMode(int(cfg.Mode)), cfg.Enabled)
		return nil
	})
}
