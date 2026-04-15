// Command iam-admin provides IAM administration utilities for vorpalstacks.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/api"
	iamstore "vorpalstacks/internal/store/aws/iam"
)

func main() {
	dataPath := flag.String("data", "./data", "Path to data directory")
	accountId := flag.String("account-id", "123456789012", "AWS Account ID")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		printUsage()
		os.Exit(1)
	}

	store, err := storage.NewPebbleStorage(&storage.Config{Path: *dataPath})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open storage: %v\n", err)
		os.Exit(1)
	}
	defer store.Close()

	userStore := iamstore.NewUserStore(store, *accountId)
	accessKeyStore := iamstore.NewAccessKeyStore(store)
	loginProfileStore := iamstore.NewLoginProfileStore(store)
	configStore := api.NewConfigStore(store)

	cmd := args[0]
	cmdArgs := args[1:]

	switch cmd {
	case "create-user":
		createUserCmd(userStore, *accountId, cmdArgs)
	case "delete-user":
		deleteUserCmd(userStore, cmdArgs)
	case "list-users":
		listUsersCmd(userStore, cmdArgs)
	case "get-user":
		getUserCmd(userStore, cmdArgs)
	case "create-access-key":
		createAccessKeyCmd(accessKeyStore, userStore, cmdArgs)
	case "list-access-keys":
		listAccessKeysCmd(accessKeyStore, userStore, cmdArgs)
	case "delete-access-key":
		deleteAccessKeyCmd(accessKeyStore, userStore, cmdArgs)
	case "create-login-profile":
		createLoginProfileCmd(loginProfileStore, userStore, cmdArgs)
	case "delete-login-profile":
		deleteLoginProfileCmd(loginProfileStore, userStore, cmdArgs)
	case "set-config":
		setConfigCmd(configStore, cmdArgs)
	case "get-config":
		getConfigCmd(configStore, cmdArgs)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", cmd)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("IAM Admin CLI")
	fmt.Println()
	fmt.Println("Usage: iam-admin [options] <command> [command-args]")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -data <path>       Path to data directory (default: ./data)")
	fmt.Println("  -account-id <id>   AWS Account ID (default: 123456789012)")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  create-user -user <name> [-path <path>]")
	fmt.Println("  delete-user -user <name>")
	fmt.Println("  list-users [-path-prefix <prefix>]")
	fmt.Println("  get-user -user <name>")
	fmt.Println("  create-access-key -user <name>")
	fmt.Println("  list-access-keys -user <name>")
	fmt.Println("  delete-access-key -access-key-id <id>")
	fmt.Println("  create-login-profile -user <name> -password <password>")
	fmt.Println("  delete-login-profile -user <name>")
	fmt.Println("  set-config -service <name> -mode <IMPLEMENTED|MOCK_SUCCESS|MOCK_ERROR>")
	fmt.Println("  get-config -service <name>")
}

func createUserCmd(store *iamstore.UserStore, accountId string, args []string) {
	fs := flag.NewFlagSet("create-user", flag.ExitOnError)
	userName := fs.String("user", "", "User name")
	path := fs.String("path", "/", "User path")
	_ = fs.Parse(args)

	if *userName == "" {
		fmt.Fprintln(os.Stderr, "Error: -user is required")
		os.Exit(1)
	}

	user, err := store.Create(*userName, *path, accountId, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create user: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("User created:\n")
	fmt.Printf("  UserName: %s\n", user.UserName)
	fmt.Printf("  UserId: %s\n", user.ID)
	fmt.Printf("  Arn: %s\n", user.Arn)
	fmt.Printf("  Path: %s\n", user.Path)
}

func deleteUserCmd(store *iamstore.UserStore, args []string) {
	fs := flag.NewFlagSet("delete-user", flag.ExitOnError)
	userName := fs.String("user", "", "User name")
	_ = fs.Parse(args)

	if *userName == "" {
		fmt.Fprintln(os.Stderr, "Error: -user is required")
		os.Exit(1)
	}

	if err := store.Delete(*userName); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to delete user: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("User '%s' deleted\n", *userName)
}

func listUsersCmd(store *iamstore.UserStore, args []string) {
	fs := flag.NewFlagSet("list-users", flag.ExitOnError)
	pathPrefix := fs.String("path-prefix", "", "Path prefix")
	_ = fs.Parse(args)

	result, err := store.List(*pathPrefix, "", 1000)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to list users: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Users (%d):\n", len(result.Users))
	for _, user := range result.Users {
		fmt.Printf("  - %s (%s)\n", user.UserName, user.Arn)
	}
}

func getUserCmd(store *iamstore.UserStore, args []string) {
	fs := flag.NewFlagSet("get-user", flag.ExitOnError)
	userName := fs.String("user", "", "User name")
	_ = fs.Parse(args)

	if *userName == "" {
		fmt.Fprintln(os.Stderr, "Error: -user is required")
		os.Exit(1)
	}

	user, err := store.Get(*userName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get user: %v\n", err)
		os.Exit(1)
	}

	data, _ := json.MarshalIndent(user, "", "  ")
	fmt.Println(string(data))
}

func createAccessKeyCmd(store *iamstore.AccessKeyStore, userStore *iamstore.UserStore, args []string) {
	fs := flag.NewFlagSet("create-access-key", flag.ExitOnError)
	userName := fs.String("user", "", "User name")
	_ = fs.Parse(args)

	if *userName == "" {
		fmt.Fprintln(os.Stderr, "Error: -user is required")
		os.Exit(1)
	}

	if !userStore.Exists(*userName) {
		fmt.Fprintf(os.Stderr, "Error: User '%s' not found\n", *userName)
		os.Exit(1)
	}

	key, err := store.Create(*userName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create access key: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Access Key created:\n")
	fmt.Printf("  AccessKeyId: %s\n", key.AccessKeyId)
	fmt.Printf("  SecretAccessKey: %s\n", key.SecretAccessKey)
	fmt.Printf("  UserName: %s\n", key.UserName)
	fmt.Printf("  Status: %s\n", key.Status)
}

func listAccessKeysCmd(store *iamstore.AccessKeyStore, userStore *iamstore.UserStore, args []string) {
	fs := flag.NewFlagSet("list-access-keys", flag.ExitOnError)
	userName := fs.String("user", "", "User name")
	_ = fs.Parse(args)

	if *userName == "" {
		fmt.Fprintln(os.Stderr, "Error: -user is required")
		os.Exit(1)
	}

	if !userStore.Exists(*userName) {
		fmt.Fprintf(os.Stderr, "Error: User '%s' not found\n", *userName)
		os.Exit(1)
	}

	keys, err := store.ListByUserName(*userName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to list access keys: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Access Keys for %s (%d):\n", *userName, len(keys))
	for _, key := range keys {
		fmt.Printf("  - %s (%s)\n", key.AccessKeyId, key.Status)
	}
}

func deleteAccessKeyCmd(store *iamstore.AccessKeyStore, userStore *iamstore.UserStore, args []string) {
	fs := flag.NewFlagSet("delete-access-key", flag.ExitOnError)
	accessKeyId := fs.String("access-key-id", "", "Access Key ID")
	_ = fs.Parse(args)

	if *accessKeyId == "" {
		fmt.Fprintln(os.Stderr, "Error: -access-key-id is required")
		os.Exit(1)
	}

	if err := store.Delete(*accessKeyId); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to delete access key: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Access Key '%s' deleted\n", *accessKeyId)
}

func createLoginProfileCmd(store *iamstore.LoginProfileStore, userStore *iamstore.UserStore, args []string) {
	fs := flag.NewFlagSet("create-login-profile", flag.ExitOnError)
	userName := fs.String("user", "", "User name")
	password := fs.String("password", "", "Password")
	passwordResetRequired := fs.Bool("password-reset-required", false, "Password reset required")
	_ = fs.Parse(args)

	if *userName == "" || *password == "" {
		fmt.Fprintln(os.Stderr, "Error: -user and -password are required")
		os.Exit(1)
	}

	if !userStore.Exists(*userName) {
		fmt.Fprintf(os.Stderr, "Error: User '%s' not found\n", *userName)
		os.Exit(1)
	}

	profile, err := store.Create(*userName, *password, *passwordResetRequired)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create login profile: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Login Profile created:\n")
	fmt.Printf("  UserName: %s\n", profile.UserName)
	fmt.Printf("  PasswordResetRequired: %v\n", profile.PasswordResetRequired)
}

func deleteLoginProfileCmd(store *iamstore.LoginProfileStore, userStore *iamstore.UserStore, args []string) {
	fs := flag.NewFlagSet("delete-login-profile", flag.ExitOnError)
	userName := fs.String("user", "", "User name")
	_ = fs.Parse(args)

	if *userName == "" {
		fmt.Fprintln(os.Stderr, "Error: -user is required")
		os.Exit(1)
	}

	if err := store.Delete(*userName); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to delete login profile: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Login Profile for '%s' deleted\n", *userName)
}

func setConfigCmd(store *api.ConfigStore, args []string) {
	fs := flag.NewFlagSet("set-config", flag.ExitOnError)
	serviceName := fs.String("service", "", "Service name")
	mode := fs.String("mode", "", "Service mode (IMPLEMENTED, MOCK_SUCCESS, MOCK_ERROR)")
	_ = fs.Parse(args)

	if *serviceName == "" || *mode == "" {
		fmt.Fprintln(os.Stderr, "Error: -service and -mode are required")
		os.Exit(1)
	}

	serviceMode := api.ParseServiceMode(strings.ToUpper(*mode))
	cfg := api.NewServiceConfig(*serviceName, serviceMode)

	if err := store.Put(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to set config: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Config set:\n")
	fmt.Printf("  Service: %s\n", *serviceName)
	fmt.Printf("  Mode: %s\n", serviceMode)
}

func getConfigCmd(store *api.ConfigStore, args []string) {
	fs := flag.NewFlagSet("get-config", flag.ExitOnError)
	serviceName := fs.String("service", "", "Service name")
	_ = fs.Parse(args)

	if *serviceName == "" {
		fmt.Fprintln(os.Stderr, "Error: -service is required")
		os.Exit(1)
	}

	cfg, err := store.Get(*serviceName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get config: %v\n", err)
		os.Exit(1)
	}

	data, _ := json.MarshalIndent(cfg, "", "  ")
	fmt.Println(string(data))
}
