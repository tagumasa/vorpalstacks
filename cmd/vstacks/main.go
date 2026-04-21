package main

import (
	"flag"
	"fmt"
	"os"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/tools/vstackscli"
)

func main() {
	dataPath := flag.String("data", "./data", "Path to data directory")
	accountId := flag.String("account-id", "123456789012", "AWS Account ID")
	endpoint := flag.String("endpoint", "http://127.0.0.1:9090", "gRPC-Web admin endpoint")
	httpEndpoint := flag.String("http-endpoint", "http://127.0.0.1:8080", "HTTP server endpoint")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		printUsage()
		os.Exit(1)
	}

	group := args[0]
	groupArgs := args[1:]

	switch group {
	case "server":
		handleServer(groupArgs, *endpoint, *httpEndpoint)
	case "iam":
		handleIAM(*dataPath, *accountId, groupArgs)
	case "config":
		handleConfig(*dataPath, groupArgs)
	case "service":
		handleService(*dataPath, groupArgs)
	case "backup":
		vstackscli.RunBackup(groupArgs)
	default:
		fmt.Fprintf(os.Stderr, "Unknown group: %s\n", group)
		printUsage()
		os.Exit(1)
	}
}

func handleServer(args []string, endpoint, httpEndpoint string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: vstacks server <stop|status>")
		os.Exit(1)
	}

	switch args[0] {
	case "stop":
		if err := vstackscli.RunServerStop(endpoint); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
	case "status":
		vstackscli.RunServerStatus(httpEndpoint)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: server %s\n", args[0])
		os.Exit(1)
	}
}

func handleIAM(dataPath, accountId string, args []string) {
	store, err := storage.NewPebbleStorage(&storage.Config{Path: dataPath})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open storage: %v\n", err)
		os.Exit(1)
	}
	defer store.Close()

	vstackscli.RunIAM(store, accountId, args)
}

func handleConfig(dataPath string, args []string) {
	store, err := storage.NewPebbleStorage(&storage.Config{Path: dataPath})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open storage: %v\n", err)
		os.Exit(1)
	}
	defer store.Close()

	vstackscli.RunConfig(store, args)
}

func handleService(dataPath string, args []string) {
	store, err := storage.NewPebbleStorage(&storage.Config{Path: dataPath})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open storage: %v\n", err)
		os.Exit(1)
	}
	defer store.Close()

	vstackscli.RunService(store, args)
}

func printUsage() {
	fmt.Println("vstacks - vorpalstacks administration CLI")
	fmt.Println()
	fmt.Println("Usage: vstacks [options] <group> <command> [args]")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -data <path>           Path to data directory (default: ./data)")
	fmt.Println("  -account-id <id>       AWS Account ID (default: 123456789012)")
	fmt.Println("  -endpoint <url>        gRPC-Web admin endpoint (default: http://127.0.0.1:9090)")
	fmt.Println("  -http-endpoint <url>   HTTP server endpoint (default: http://127.0.0.1:8080)")
	fmt.Println()
	fmt.Println("Groups:")
	fmt.Println("  server   Server control (stop, status)")
	fmt.Println("  iam      IAM user and access key management")
	fmt.Println("  config   Application configuration (app_config)")
	fmt.Println("  service  Service mode configuration (service_config)")
	fmt.Println("  backup   Data backup and restore")
}
