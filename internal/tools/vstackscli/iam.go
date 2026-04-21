package vstackscli

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"vorpalstacks/internal/core/storage"
	iamstore "vorpalstacks/internal/store/aws/iam"
)

// RunIAM dispatches IAM administration commands (user, access-key, login-profile
// CRUD) against the PebbleDB-backed IAM stores.
//
// Parameters:
//   - store: The underlying PebbleDB storage
//   - accountId: The AWS account ID used when creating users
//   - args: Command-line arguments (sub-command and flags)
func RunIAM(store storage.BasicStorage, accountId string, args []string) {
	if len(args) == 0 {
		printIAMUsage()
		os.Exit(1)
	}

	userStore := iamstore.NewUserStore(store, accountId)
	accessKeyStore := iamstore.NewAccessKeyStore(store)
	loginProfileStore := iamstore.NewLoginProfileStore(store)

	cmd := args[0]
	cmdArgs := args[1:]

	switch cmd {
	case "create-user":
		createUserCmd(userStore, accountId, cmdArgs)
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
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", cmd)
		printIAMUsage()
		os.Exit(1)
	}
}

func printIAMUsage() {
	fmt.Println("IAM Commands:")
	fmt.Println("  vstacks iam create-user -user <name> [-path <path>]")
	fmt.Println("  vstacks iam delete-user -user <name>")
	fmt.Println("  vstacks iam list-users [-path-prefix <prefix>]")
	fmt.Println("  vstacks iam get-user -user <name>")
	fmt.Println("  vstacks iam create-access-key -user <name>")
	fmt.Println("  vstacks iam list-access-keys -user <name>")
	fmt.Println("  vstacks iam delete-access-key -access-key-id <id>")
	fmt.Println("  vstacks iam create-login-profile -user <name> -password <password>")
	fmt.Println("  vstacks iam delete-login-profile -user <name>")
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
	pathPrefix := fs.String("path-prefix", "", "User path prefix")
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
	_ = fs.Parse(args)

	if *userName == "" || *password == "" {
		fmt.Fprintln(os.Stderr, "Error: -user and -password are required")
		os.Exit(1)
	}

	if !userStore.Exists(*userName) {
		fmt.Fprintf(os.Stderr, "Error: User '%s' not found\n", *userName)
		os.Exit(1)
	}

	profile, err := store.Create(*userName, *password, false)
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

func formatMode(mode int) string {
	switch mode {
	case 0:
		return "IMPLEMENTED"
	case 1:
		return "FALLBACK"
	case 2:
		return "ERROR_INJECTION"
	default:
		return fmt.Sprintf("UNKNOWN(%d)", mode)
	}
}

func parseMode(s string) int {
	switch strings.ToUpper(s) {
	case "IMPLEMENTED":
		return 0
	case "FALLBACK":
		return 1
	case "ERROR_INJECTION":
		return 2
	default:
		return 1
	}
}
