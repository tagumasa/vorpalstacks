// Package admin_config provides runtime configuration management for vorpalstacks.
package admin_config

import (
	"context"
	"encoding/json"
	"fmt"

	"connectrpc.com/connect"

	pb "vorpalstacks/internal/pb/aws/admin_config"
	"vorpalstacks/internal/pb/aws/common"
	storeconfig "vorpalstacks/internal/store/config"
)

// ConfigStore defines the interface for configuration storage operations.
type ConfigStore interface {
	Get(key string) (*storeconfig.ConfigEntry, error)
	Set(key string, value interface{}) error
	Reset(key string) (*storeconfig.ConfigEntry, error)
	GetAll() ([]*storeconfig.ConfigEntry, error)
	ListByCategory(category storeconfig.ConfigCategory) ([]*storeconfig.ConfigEntry, error)
	GetResourcePort(servicePortKey, resourceID string) (int, error)
	SetResourcePort(servicePortKey, resourceID string, port int) error
}

// AdminConfigService provides admin configuration management and server control.
type AdminConfigService struct {
	configStore  ConfigStore
	shutdownFunc func()
}

// NewAdminConfigService creates a new admin config service instance.
func NewAdminConfigService(configStore ConfigStore, shutdownFunc func()) *AdminConfigService {
	return &AdminConfigService{
		configStore:  configStore,
		shutdownFunc: shutdownFunc,
	}
}

// GetConfig retrieves a configuration value by key.
//
// Parameters:
//   - ctx: The request context
//   - req: The get config request containing the key
//
// Returns:
//   - *connect.Response[pb.GetConfigResponse]: The configuration entry if found
//   - error: An error if the key is not found or operation fails
func (s *AdminConfigService) GetConfig(ctx context.Context, req *connect.Request[pb.GetConfigRequest]) (*connect.Response[pb.GetConfigResponse], error) {
	entry, err := s.configStore.Get(req.Msg.Key)
	if err != nil {
		if err == storeconfig.ErrConfigNotFound {
			return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("configuration not found: %s", req.Msg.Key))
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&pb.GetConfigResponse{
		Entry: toPbEntry(entry),
	}), nil
}

// ListConfig retrieves configuration entries, optionally filtered by category.
//
// Parameters:
//   - ctx: The request context
//   - req: The list config request, optionally containing a category filter
//
// Returns:
//   - *connect.Response[pb.ListConfigResponse]: The list of configuration entries
//   - error: An error if the operation fails
func (s *AdminConfigService) ListConfig(ctx context.Context, req *connect.Request[pb.ListConfigRequest]) (*connect.Response[pb.ListConfigResponse], error) {
	var entries []*storeconfig.ConfigEntry
	var err error

	if req.Msg.Category != "" {
		entries, err = s.configStore.ListByCategory(storeconfig.ConfigCategory(req.Msg.Category))
	} else {
		entries, err = s.configStore.GetAll()
	}

	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	pbEntries := make([]*pb.ConfigEntry, len(entries))
	for i, entry := range entries {
		pbEntries[i] = toPbEntry(entry)
	}

	return connect.NewResponse(&pb.ListConfigResponse{
		Entries: pbEntries,
	}), nil
}

// UpdateConfig updates a configuration value by key.
//
// Parameters:
//   - ctx: The request context
//   - req: The update config request containing the key and new value
//
// Returns:
//   - *connect.Response[pb.ConfigEntry]: The updated configuration entry
//   - error: An error if the update fails
func (s *AdminConfigService) UpdateConfig(ctx context.Context, req *connect.Request[pb.UpdateConfigRequest]) (*connect.Response[pb.ConfigEntry], error) {
	var value interface{}
	if err := json.Unmarshal([]byte(req.Msg.Value), &value); err != nil {
		value = req.Msg.Value
	}

	if err := s.configStore.Set(req.Msg.Key, value); err != nil {
		if err == storeconfig.ErrConfigReadOnly {
			return nil, connect.NewError(connect.CodePermissionDenied, fmt.Errorf("configuration is read-only: %s", req.Msg.Key))
		}
		if err == storeconfig.ErrConfigNotFound {
			return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("configuration not found: %s", req.Msg.Key))
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	entry, err := s.configStore.Get(req.Msg.Key)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(toPbEntry(entry)), nil
}

// ResetConfig resets a configuration value to its default.
//
// Parameters:
//   - ctx: The request context
//   - req: The reset config request containing the key
//
// Returns:
//   - *connect.Response[pb.ConfigEntry]: The reset configuration entry
//   - error: An error if the key is not found or reset fails
func (s *AdminConfigService) ResetConfig(ctx context.Context, req *connect.Request[pb.ResetConfigRequest]) (*connect.Response[pb.ConfigEntry], error) {
	entry, err := s.configStore.Reset(req.Msg.Key)
	if err != nil {
		if err == storeconfig.ErrConfigNotFound {
			return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("configuration not found: %s", req.Msg.Key))
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(toPbEntry(entry)), nil
}

// ListServices lists available services.
//
// Parameters:
//   - ctx: The request context
//   - req: The list services request
//
// Returns:
//   - *connect.Response[pb.ListServicesResponse]: The list of available services
//   - error: An error if the operation fails
func (s *AdminConfigService) ListServices(ctx context.Context, req *connect.Request[pb.ListServicesRequest]) (*connect.Response[pb.ListServicesResponse], error) {
	return connect.NewResponse(&pb.ListServicesResponse{
		Services: []*pb.ServiceInfo{},
	}), nil
}

// GetServiceStatus retrieves the status of a service.
//
// Parameters:
//   - ctx: The request context
//   - req: The get service status request containing the service name
//
// Returns:
//   - *connect.Response[pb.ServiceStatus]: The service status if found
//   - error: An error if the service is not found
func (s *AdminConfigService) GetServiceStatus(ctx context.Context, req *connect.Request[pb.GetServiceStatusRequest]) (*connect.Response[pb.ServiceStatus], error) {
	return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("service not found: %s", req.Msg.Name))
}

// GetResourcePort retrieves the port for a specific resource.
//
// Parameters:
//   - ctx: The request context
//   - req: The get resource port request containing service port key and resource ID
//
// Returns:
//   - *connect.Response[pb.GetResourcePortResponse]: The resource port information
//   - error: An error if the port configuration is not found
func (s *AdminConfigService) GetResourcePort(ctx context.Context, req *connect.Request[pb.GetResourcePortRequest]) (*connect.Response[pb.GetResourcePortResponse], error) {
	port, err := s.configStore.GetResourcePort(req.Msg.ServicePortKey, req.Msg.ResourceId)
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("port configuration not found"))
	}

	entry, _ := s.configStore.Get(req.Msg.ServicePortKey + "." + req.Msg.ResourceId)
	source := "default"
	if entry != nil {
		source = string(entry.Source)
	}

	return connect.NewResponse(&pb.GetResourcePortResponse{
		Port:   int32(port),
		Source: source,
	}), nil
}

// SetResourcePort sets the port for a specific resource.
//
// Parameters:
//   - ctx: The request context
//   - req: The set resource port request containing service port key, resource ID and port
//
// Returns:
//   - *connect.Response[common.Empty]: An empty response on success
//   - error: An error if the operation fails
func (s *AdminConfigService) SetResourcePort(ctx context.Context, req *connect.Request[pb.SetResourcePortRequest]) (*connect.Response[common.Empty], error) {
	if err := s.configStore.SetResourcePort(req.Msg.ServicePortKey, req.Msg.ResourceId, int(req.Msg.Port)); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&common.Empty{}), nil
}

// ShutdownServer triggers a graceful server shutdown. The response is sent
// immediately; the actual shutdown happens asynchronously in a goroutine.
func (s *AdminConfigService) ShutdownServer(ctx context.Context, req *connect.Request[pb.ShutdownServerRequest]) (*connect.Response[pb.ShutdownServerResponse], error) {
	resp := connect.NewResponse(&pb.ShutdownServerResponse{
		Message: "Server is shutting down",
	})

	if s.shutdownFunc != nil {
		go s.shutdownFunc()
	}

	return resp, nil
}

func toPbEntry(entry *storeconfig.ConfigEntry) *pb.ConfigEntry {
	valueBytes, err := json.Marshal(entry.Value)
	if err != nil {
		valueBytes = []byte("{}")
	}
	return &pb.ConfigEntry{
		Key:         entry.Key,
		Value:       string(valueBytes),
		Type:        string(entry.Type),
		Source:      string(entry.Source),
		Description: entry.Description,
		ReadOnly:    entry.ReadOnly,
		UpdatedAt:   entry.UpdatedAt,
		EnvVar:      entry.EnvVar,
		Category:    string(entry.Category),
	}
}
