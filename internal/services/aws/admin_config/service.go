// Package admin_config provides runtime configuration management for vorpalstacks.
package admin_config

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime"
	"strconv"
	"time"

	"connectrpc.com/connect"

	pb "vorpalstacks/internal/pb/aws/admin_config"
	"vorpalstacks/internal/pb/aws/common"
	"vorpalstacks/internal/serviceconfig"
	storeconfig "vorpalstacks/internal/store/config"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
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
	startTime    time.Time
	dataPath     string
	version      string
}

func NewAdminConfigService(configStore ConfigStore, shutdownFunc func(), dataPath string, version string) *AdminConfigService {
	cpu.Percent(0, false)
	return &AdminConfigService{
		configStore:  configStore,
		shutdownFunc: shutdownFunc,
		startTime:    time.Now(),
		dataPath:     dataPath,
		version:      version,
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

	value = s.coerceValue(req.Msg.Key, value)

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
	allServices := serviceconfig.Services
	result := make([]*pb.ServiceInfo, 0, len(allServices))

	for i := range allServices {
		svc := &allServices[i]
		result = append(result, &pb.ServiceInfo{
			Name:     svc.Name,
			Enabled:  s.isServiceEnabled(svc.Name),
			PortMode: s.resolvedPortMode(svc),
		})
	}

	return connect.NewResponse(&pb.ListServicesResponse{
		Services: result,
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
	svc := serviceconfig.ByName(req.Msg.Name)
	if svc == nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("service not found: %s", req.Msg.Name))
	}

	return connect.NewResponse(&pb.ServiceStatus{
		Name:    svc.Name,
		Enabled: s.isServiceEnabled(svc.Name),
	}), nil
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

// GetServerMetrics returns process and machine-level telemetry metrics.
func (s *AdminConfigService) GetServerMetrics(ctx context.Context, req *connect.Request[pb.GetServerMetricsRequest]) (*connect.Response[pb.GetServerMetricsResponse], error) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	resp := &pb.GetServerMetricsResponse{
		UptimeSeconds:         int64(time.Since(s.startTime).Seconds()),
		ProcessMemorySysBytes: int64(m.Sys),
		ProcessHeapAllocBytes: int64(m.HeapAlloc),
		GoroutineCount:        int32(runtime.NumGoroutine()),
		GcCount:               int64(m.NumGC),
		GcPauseTotalNs:        int64(m.PauseTotalNs),
		NumCpu:                int32(runtime.NumCPU()),
		GoVersion:             runtime.Version(),
		Version:               s.version,
	}

	if vm, err := mem.VirtualMemory(); err == nil {
		resp.MachineTotalMemoryBytes = int64(vm.Total)
		resp.MachineAvailableMemoryBytes = int64(vm.Available)
	}

	if pct, err := cpu.Percent(0, false); err == nil && len(pct) > 0 {
		resp.MachineCpuUsagePercent = pct[0]
	}

	if du, err := disk.Usage(s.dataPath); err == nil {
		resp.MachineDiskTotalBytes = int64(du.Total)
		resp.MachineDiskFreeBytes = int64(du.Free)
	}

	return connect.NewResponse(resp), nil
}

// EnableService enables a service by setting services.{name}.enabled = true.
func (s *AdminConfigService) EnableService(ctx context.Context, req *connect.Request[pb.EnableServiceRequest]) (*connect.Response[pb.ServiceStatus], error) {
	svc := serviceconfig.ByName(req.Msg.ServiceName)
	if svc == nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("service not found: %s", req.Msg.ServiceName))
	}

	key := "services." + svc.Name + ".enabled"
	if err := s.configStore.Set(key, true); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&pb.ServiceStatus{
		Name:    svc.Name,
		Enabled: true,
	}), nil
}

func (s *AdminConfigService) DisableService(ctx context.Context, req *connect.Request[pb.DisableServiceRequest]) (*connect.Response[pb.ServiceStatus], error) {
	svc := serviceconfig.ByName(req.Msg.ServiceName)
	if svc == nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("service not found: %s", req.Msg.ServiceName))
	}

	key := "services." + svc.Name + ".enabled"
	if err := s.configStore.Set(key, false); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&pb.ServiceStatus{
		Name:    svc.Name,
		Enabled: false,
	}), nil
}

// GetPortMode returns the current port mode for a service.
func (s *AdminConfigService) GetPortMode(ctx context.Context, req *connect.Request[pb.GetPortModeRequest]) (*connect.Response[pb.PortModeResponse], error) {
	svc := serviceconfig.ByName(req.Msg.ServiceName)
	if svc == nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("service not found: %s", req.Msg.ServiceName))
	}

	mode := s.resolvedPortMode(svc)
	resp := &pb.PortModeResponse{
		ServiceName: svc.Name,
		Mode:        mode,
	}

	if mode == "individual" {
		if p, err := s.configStore.GetResourcePort(svc.PortKey, "default"); err == nil && p > 0 {
			resp.Port = int32(p)
		}
	}

	return connect.NewResponse(resp), nil
}

func (s *AdminConfigService) SetPortMode(ctx context.Context, req *connect.Request[pb.SetPortModeRequest]) (*connect.Response[pb.PortModeResponse], error) {
	svc := serviceconfig.ByName(req.Msg.ServiceName)
	if svc == nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("service not found: %s", req.Msg.ServiceName))
	}

	if req.Msg.Mode != "fqdn" && req.Msg.Mode != "individual" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid port mode: %s (must be fqdn or individual)", req.Msg.Mode))
	}

	key := svc.PortKey + ".mode"
	if err := s.configStore.Set(key, req.Msg.Mode); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	resp := &pb.PortModeResponse{
		ServiceName: svc.Name,
		Mode:        req.Msg.Mode,
	}

	if req.Msg.Mode == "individual" {
		if p, err := s.configStore.GetResourcePort(svc.PortKey, "default"); err == nil && p > 0 {
			resp.Port = int32(p)
		}
	}

	return connect.NewResponse(resp), nil
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
	return &pb.ConfigEntry{
		Key:         entry.Key,
		Value:       formatValue(entry.Value, entry.Type),
		Type:        string(entry.Type),
		Source:      string(entry.Source),
		Description: entry.Description,
		ReadOnly:    entry.ReadOnly,
		UpdatedAt:   entry.UpdatedAt,
		EnvVar:      entry.EnvVar,
		Category:    string(entry.Category),
	}
}

func formatValue(v interface{}, typ storeconfig.ConfigType) string {
	switch typ {
	case storeconfig.ConfigTypeInt, storeconfig.ConfigTypePort:
		switch n := v.(type) {
		case int:
			return strconv.Itoa(n)
		case float64:
			return strconv.Itoa(int(n))
		case int64:
			return strconv.Itoa(int(n))
		}
	case storeconfig.ConfigTypeBool:
		return fmt.Sprintf("%v", v)
	}
	return fmt.Sprintf("%v", v)
}

func (s *AdminConfigService) isServiceEnabled(name string) bool {
	if entry, err := s.configStore.Get("services." + name + ".enabled"); err == nil {
		if v, ok := entry.Value.(bool); ok {
			return v
		}
	}
	return true
}

func (s *AdminConfigService) resolvedPortMode(svc *serviceconfig.ServiceDef) string {
	if svc.PortKey == "" {
		return svc.DefaultPortMode()
	}
	if entry, err := s.configStore.Get(svc.PortKey + ".mode"); err == nil {
		if v, ok := entry.Value.(string); ok && v != "" {
			return v
		}
	}
	return svc.DefaultPortMode()
}

func (s *AdminConfigService) coerceValue(key string, value interface{}) interface{} {
	entry, err := s.configStore.Get(key)
	if err != nil {
		return value
	}
	switch entry.Type {
	case storeconfig.ConfigTypePort, storeconfig.ConfigTypeInt:
		switch v := value.(type) {
		case float64:
			return int(v)
		case string:
			if i, e := strconv.Atoi(v); e == nil {
				return i
			}
		}
	case storeconfig.ConfigTypeBool:
		switch v := value.(type) {
		case string:
			return v == "true" || v == "1"
		case float64:
			return v != 0
		}
	}
	return value
}
