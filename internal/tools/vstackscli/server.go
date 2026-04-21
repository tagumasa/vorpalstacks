package vstackscli

import (
	"context"
	"fmt"
	"net/http"
	"time"

	connect "connectrpc.com/connect"
	pb "vorpalstacks/internal/pb/aws/admin_config"
	"vorpalstacks/internal/pb/aws/admin_config/admin_configconnect"
)

// RunServerStop sends a ShutdownServer RPC to the gRPC-Web admin endpoint and
// waits for the server to acknowledge the request.
//
// Parameters:
//   - endpoint: The gRPC-Web admin endpoint URL
//
// Returns:
//   - error: An error if the RPC fails
func RunServerStop(endpoint string) error {
	client := admin_configconnect.NewAdminConfigServiceClient(
		http.DefaultClient,
		endpoint,
		connect.WithGRPC(),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.ShutdownServer(ctx, connect.NewRequest(&pb.ShutdownServerRequest{}))
	if err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}

	fmt.Printf("%s\n", resp.Msg.Message)
	return nil
}

// RunServerStatus checks whether the vorpalstacks HTTP server is reachable by
// querying the health endpoint.
//
// Parameters:
//   - httpEndpoint: The HTTP server base URL
func RunServerStatus(httpEndpoint string) {
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get(httpEndpoint + "/.well-known/health")
	if err != nil {
		fmt.Println("Server is not running")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Server is running")
	} else {
		fmt.Printf("Server returned status %d\n", resp.StatusCode)
	}
}


