// Package mobyclient provides Docker/Moby API client operations for containers.
package mobyclient

import (
	"context"
	"fmt"
	"io"
	"net/netip"
	"strings"
	"time"

	"vorpalstacks/internal/core/logs"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/events"
	"github.com/moby/moby/api/types/mount"
	"github.com/moby/moby/api/types/network"
	"github.com/moby/moby/client"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

// CreateContainer creates a new container.
func (c *Client) CreateContainer(ctx context.Context, cfg ContainerConfig) (string, error) {
	c.logger.Debug("Creating container",
		logs.String("image", cfg.Image),
		logs.String("name", cfg.Name),
	)

	containerCfg := &container.Config{
		Image:      cfg.Image,
		Cmd:        cfg.Cmd,
		Entrypoint: cfg.Entrypoint,
		Env:        cfg.Env,
		Labels:     cfg.Labels,
		WorkingDir: cfg.WorkingDir,
		Hostname:   cfg.Hostname,
	}

	hostCfg := &container.HostConfig{
		NetworkMode: container.NetworkMode(cfg.NetworkMode),
		Mounts:      convertMounts(cfg.Mounts),
		Resources: container.Resources{
			Memory:     cfg.Memory,
			MemorySwap: cfg.MemorySwap,
			NanoCPUs:   cfg.CPUs,
		},
	}

	if cfg.UseGPU {
		hostCfg.DeviceRequests = []container.DeviceRequest{
			{
				Driver:       "nvidia",
				Count:        -1,
				Capabilities: [][]string{{"gpu"}},
			},
		}
	}

	networkCfg := &network.NetworkingConfig{}

	options := client.ContainerCreateOptions{
		Config:           containerCfg,
		HostConfig:       hostCfg,
		NetworkingConfig: networkCfg,
		Name:             cfg.Name,
	}

	resp, err := c.cli.ContainerCreate(ctx, options)
	if err != nil {
		return "", fmt.Errorf("failed to create container: %w", err)
	}

	c.logger.Info("Container created",
		logs.String("id", resp.ID),
		logs.String("name", cfg.Name),
	)

	return resp.ID, nil
}

// StartContainer starts a container.
func (c *Client) StartContainer(ctx context.Context, containerID string) error {
	c.logger.Debug("Starting container", logs.String("id", containerID))

	_, err := c.cli.ContainerStart(ctx, containerID, client.ContainerStartOptions{})
	if err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}

	c.logger.Info("Container started", logs.String("id", containerID))
	return nil
}

// StartContainerWithCheckpoint starts a container from a checkpoint.
func (c *Client) StartContainerWithCheckpoint(ctx context.Context, containerID, checkpointDir string) error {
	options := client.ContainerStartOptions{
		CheckpointID:  checkpointDir,
		CheckpointDir: checkpointDir,
	}
	_, err := c.cli.ContainerStart(ctx, containerID, options)
	if err != nil {
		return fmt.Errorf("failed to start container with checkpoint: %w", err)
	}
	return nil
}

// StopContainer stops a container.
func (c *Client) StopContainer(ctx context.Context, containerID string, timeout int) error {
	c.logger.Debug("Stopping container", logs.String("id", containerID))

	stopTimeout := timeout
	if stopTimeout == 0 {
		stopTimeout = 10
	}

	options := client.ContainerStopOptions{
		Timeout: &stopTimeout,
	}
	_, err := c.cli.ContainerStop(ctx, containerID, options)
	if err != nil {
		return fmt.Errorf("failed to stop container: %w", err)
	}

	c.logger.Info("Container stopped", logs.String("id", containerID))
	return nil
}

// StopContainerForce stops a container forcefully.
func (c *Client) StopContainerForce(ctx context.Context, containerID string) error {
	timeout := 0
	options := client.ContainerStopOptions{
		Timeout: &timeout,
	}
	_, err := c.cli.ContainerStop(ctx, containerID, options)
	if err != nil {
		return fmt.Errorf("failed to stop container force: %w", err)
	}
	return nil
}

// RestartContainer restarts a container.
func (c *Client) RestartContainer(ctx context.Context, containerID string, timeout int) error {
	c.logger.Debug("Restarting container", logs.String("id", containerID))

	restartTimeout := timeout
	if restartTimeout == 0 {
		restartTimeout = 10
	}

	options := client.ContainerRestartOptions{
		Timeout: &restartTimeout,
	}
	_, err := c.cli.ContainerRestart(ctx, containerID, options)
	if err != nil {
		return fmt.Errorf("failed to restart container: %w", err)
	}

	c.logger.Info("Container restarted", logs.String("id", containerID))
	return nil
}

// RemoveContainer removes a container.
func (c *Client) RemoveContainer(ctx context.Context, containerID string, force bool) error {
	c.logger.Debug("Removing container", logs.String("id", containerID))

	opts := client.ContainerRemoveOptions{
		Force:         force,
		RemoveVolumes: true,
	}

	_, err := c.cli.ContainerRemove(ctx, containerID, opts)
	if err != nil {
		return fmt.Errorf("failed to remove container: %w", err)
	}

	c.logger.Info("Container removed", logs.String("id", containerID))
	return nil
}

// GetContainer gets container information.
func (c *Client) GetContainer(ctx context.Context, containerID string) (*ContainerInfo, error) {
	result, err := c.cli.ContainerInspect(ctx, containerID, client.ContainerInspectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to inspect container: %w", err)
	}

	inspect := result.Container

	networkIPs := make(map[string]string)
	if inspect.NetworkSettings != nil && inspect.NetworkSettings.Networks != nil {
		for netName, netSettings := range inspect.NetworkSettings.Networks {
			networkIPs[netName] = netSettings.IPAddress.String()
		}
	}

	createdTime, err := time.Parse(time.RFC3339, inspect.Created)
	if err != nil {
		createdTime = time.Now()
	}

	return &ContainerInfo{
		ID:         inspect.ID,
		Name:       inspect.Name,
		Image:      inspect.Image,
		State:      string(inspect.State.Status),
		Status:     string(inspect.State.Status),
		Created:    createdTime,
		Labels:     inspect.Config.Labels,
		NetworkIPs: networkIPs,
	}, nil
}

// ListContainers lists all containers.
func (c *Client) ListContainers(ctx context.Context, all bool) ([]ContainerInfo, error) {
	result, err := c.cli.ContainerList(ctx, client.ContainerListOptions{All: all})
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %w", err)
	}

	containers := result.Items
	containersList := make([]ContainerInfo, 0, len(containers))
	for _, ctr := range containers {
		ports := make([]Port, 0, len(ctr.Ports))
		for _, p := range ctr.Ports {
			ports = append(ports, Port{
				IP:          p.IP.String(),
				PrivatePort: p.PrivatePort,
				PublicPort:  p.PublicPort,
				Protocol:    PortProtocolTCP,
			})
		}

		name := ""
		if len(ctr.Names) > 0 {
			name = ctr.Names[0]
		}
		containersList = append(containersList, ContainerInfo{
			ID:      ctr.ID,
			Name:    name,
			Image:   ctr.Image,
			State:   string(ctr.State),
			Status:  ctr.Status,
			Created: time.Unix(ctr.Created, 0),
			Ports:   ports,
			Labels:  ctr.Labels,
		})
	}

	return containersList, nil
}

// ListContainersWithFilter lists containers with filters.
func (c *Client) ListContainersWithFilter(ctx context.Context, all bool, filterArgs client.Filters) ([]ContainerInfo, error) {
	opts := client.ContainerListOptions{
		All:     all,
		Filters: filterArgs,
	}

	result, err := c.cli.ContainerList(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list containers with filter: %w", err)
	}

	containers := result.Items
	containersList := make([]ContainerInfo, 0, len(containers))
	for _, ctr := range containers {
		ports := make([]Port, 0, len(ctr.Ports))
		for _, p := range ctr.Ports {
			ports = append(ports, Port{
				IP:          p.IP.String(),
				PrivatePort: p.PrivatePort,
				PublicPort:  p.PublicPort,
				Protocol:    PortProtocolTCP,
			})
		}

		name2 := ""
		if len(ctr.Names) > 0 {
			name2 = ctr.Names[0]
		}
		containersList = append(containersList, ContainerInfo{
			ID:      ctr.ID,
			Name:    name2,
			Image:   ctr.Image,
			State:   string(ctr.State),
			Status:  ctr.Status,
			Created: time.Unix(ctr.Created, 0),
			Ports:   ports,
			Labels:  ctr.Labels,
		})
	}

	return containersList, nil
}

// WaitContainer waits for a container to stop.
func (c *Client) WaitContainer(ctx context.Context, containerID string) (int64, error) {
	result := c.cli.ContainerWait(ctx, containerID, client.ContainerWaitOptions{Condition: container.WaitConditionNotRunning})

	select {
	case err := <-result.Error:
		if err != nil {
			return 0, fmt.Errorf("error waiting for container: %w", err)
		}
	case status := <-result.Result:
		return status.StatusCode, nil
	}

	return 0, nil
}

// KillContainer kills a container.
func (c *Client) KillContainer(ctx context.Context, containerID string, signal string) error {
	c.logger.Debug("Killing container",
		logs.String("id", containerID),
		logs.String("signal", signal),
	)

	options := client.ContainerKillOptions{
		Signal: signal,
	}
	_, err := c.cli.ContainerKill(ctx, containerID, options)
	if err != nil {
		return fmt.Errorf("failed to kill container: %w", err)
	}

	c.logger.Info("Container killed",
		logs.String("id", containerID),
		logs.String("signal", signal),
	)
	return nil
}

// GetContainerStats gets container statistics.
func (c *Client) GetContainerStats(ctx context.Context, containerID string, stream bool) (io.ReadCloser, error) {
	options := client.ContainerStatsOptions{
		Stream: stream,
	}
	result, err := c.cli.ContainerStats(ctx, containerID, options)
	if err != nil {
		return nil, fmt.Errorf("failed to get container stats: %w", err)
	}

	return result.Body, nil
}

// CreateContainerFromConfig creates a container with advanced configuration.
func (c *Client) CreateContainerFromConfig(ctx context.Context, cfg AdvancedContainerConfig) (*CreateContainerResult, error) {
	c.logger.Debug("Creating container from config",
		logs.String("image", cfg.Image),
		logs.String("name", cfg.Name),
		logs.String("platform", string(cfg.Platform)),
	)

	containerCfg := &container.Config{
		Image:      cfg.Image,
		Cmd:        cfg.Cmd,
		Entrypoint: cfg.Entrypoint,
		Labels:     cfg.Labels,
		WorkingDir: cfg.WorkingDir,
		Hostname:   cfg.Hostname,
		User:       cfg.User,
		OpenStdin:  cfg.OpenStdin,
		Tty:        cfg.TTY,
		StdinOnce:  cfg.StdinOnce,
		Env:        convertEnvMap(cfg.Env),
	}

	hostCfg := &container.HostConfig{
		NetworkMode:    container.NetworkMode(cfg.Network),
		AutoRemove:     cfg.AutoRemove,
		Privileged:     cfg.Privileged,
		ReadonlyRootfs: cfg.ReadonlyRootfs,
		SecurityOpt:    cfg.SecurityOpt,
		CapAdd:         cfg.CapAdd,
		CapDrop:        cfg.CapDrop,
		DNSOptions:     cfg.DNSOptions,
		DNSSearch:      cfg.DNSSearch,
		ExtraHosts:     cfg.ExtraHosts,
		Resources: container.Resources{
			Memory:     cfg.Memory,
			MemorySwap: cfg.MemorySwap,
			NanoCPUs:   cfg.CPUShares,
		},
	}

	if len(cfg.DNS) > 0 {
		hostCfg.DNS = convertDNSAddresses(cfg.DNS)
	}

	if len(cfg.PortMappings) > 0 || len(cfg.HostPorts) > 0 {
		hostCfg.PortBindings = convertPortBindings(cfg.PortMappings, cfg.HostPorts)
		containerCfg.ExposedPorts = convertExposedPorts(cfg.PortMappings, cfg.HostPorts)
	}

	if len(cfg.BindMounts) > 0 || len(cfg.VolumeMount) > 0 || len(cfg.TmpfsMounts) > 0 {
		hostCfg.Mounts = convertAdvancedMounts(cfg.BindMounts, cfg.VolumeMount, cfg.TmpfsMounts)
	}

	if len(cfg.DeviceRequests) > 0 {
		hostCfg.DeviceRequests = convertDeviceRequests(cfg.DeviceRequests)
	}

	if len(cfg.Ulimits) > 0 {
		hostCfg.Ulimits = convertUlimits(cfg.Ulimits)
	}

	if cfg.RestartPolicy != "" {
		hostCfg.RestartPolicy = container.RestartPolicy{
			Name: container.RestartPolicyMode(cfg.RestartPolicy),
		}
	}

	networkCfg := &network.NetworkingConfig{}
	if len(cfg.NetworkAlias) > 0 && cfg.Network != "" {
		networkCfg.EndpointsConfig = map[string]*network.EndpointSettings{
			cfg.Network: {
				Aliases: cfg.NetworkAlias,
			},
		}
	}

	var platform *ocispec.Platform
	if cfg.Platform != "" {
		parts := strings.Split(string(cfg.Platform), "/")
		if len(parts) >= 2 {
			p := ocispec.Platform{
				OS:           parts[0],
				Architecture: parts[1],
			}
			platform = &p
		}
	}

	options := client.ContainerCreateOptions{
		Config:           containerCfg,
		HostConfig:       hostCfg,
		NetworkingConfig: networkCfg,
		Name:             cfg.Name,
		Platform:         platform,
	}

	resp, err := c.cli.ContainerCreate(ctx, options)
	if err != nil {
		return nil, fmt.Errorf("failed to create container: %w", err)
	}

	c.logger.Info("Container created from config",
		logs.String("id", resp.ID),
		logs.String("name", cfg.Name),
	)

	return &CreateContainerResult{
		ID:       resp.ID,
		Warnings: resp.Warnings,
	}, nil
}

// GetContainerIPForNetwork gets the IP address of a container in a specific network.
func (c *Client) GetContainerIPForNetwork(ctx context.Context, containerID string, networkName string) (string, error) {
	result, err := c.cli.ContainerInspect(ctx, containerID, client.ContainerInspectOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to inspect container: %w", err)
	}

	inspect := result.Container
	if inspect.NetworkSettings == nil || inspect.NetworkSettings.Networks == nil {
		return "", fmt.Errorf("container has no network settings")
	}

	netSettings, ok := inspect.NetworkSettings.Networks[networkName]
	if !ok {
		return "", fmt.Errorf("container not connected to network %s", networkName)
	}

	return netSettings.IPAddress.String(), nil
}

// PauseContainer pauses a container.
func (c *Client) PauseContainer(ctx context.Context, containerID string) error {
	c.logger.Debug("Pausing container", logs.String("id", containerID))

	_, err := c.cli.ContainerPause(ctx, containerID, client.ContainerPauseOptions{})
	if err != nil {
		return fmt.Errorf("failed to pause container: %w", err)
	}

	c.logger.Info("Container paused", logs.String("id", containerID))
	return nil
}

// UnpauseContainer unpauses a container.
func (c *Client) UnpauseContainer(ctx context.Context, containerID string) error {
	c.logger.Debug("Unpausing container", logs.String("id", containerID))

	_, err := c.cli.ContainerUnpause(ctx, containerID, client.ContainerUnpauseOptions{})
	if err != nil {
		return fmt.Errorf("failed to unpause container: %w", err)
	}

	c.logger.Info("Container unpaused", logs.String("id", containerID))
	return nil
}

// GetContainerStatus gets the status of a container.
func (c *Client) GetContainerStatus(ctx context.Context, containerID string) (ContainerStatus, error) {
	result, err := c.cli.ContainerInspect(ctx, containerID, client.ContainerInspectOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to inspect container: %w", err)
	}

	return ContainerStatus(string(result.Container.State.Status)), nil
}

// InspectContainer gets detailed container information.
func (c *Client) InspectContainer(ctx context.Context, containerID string) (*DetailedContainerInfo, error) {
	result, err := c.cli.ContainerInspect(ctx, containerID, client.ContainerInspectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to inspect container: %w", err)
	}

	inspect := result.Container

	networkIPs := make(map[string]string)
	var macAddress, gateway, ipAddress string
	if inspect.NetworkSettings != nil && inspect.NetworkSettings.Networks != nil {
		for netName, netSettings := range inspect.NetworkSettings.Networks {
			networkIPs[netName] = netSettings.IPAddress.String()
			if macAddress == "" {
				macAddress = netSettings.MacAddress.String()
			}
			if gateway == "" {
				gateway = netSettings.Gateway.String()
			}
			if ipAddress == "" {
				ipAddress = netSettings.IPAddress.String()
			}
		}
	}

	createdTime, _ := time.Parse(time.RFC3339, inspect.Created)
	startedTime, _ := time.Parse(time.RFC3339, inspect.State.StartedAt)
	finishedTime, _ := time.Parse(time.RFC3339, inspect.State.FinishedAt)

	return &DetailedContainerInfo{
		ID:         inspect.ID,
		Name:       inspect.Name,
		Image:      inspect.Image,
		ImageID:    inspect.Image,
		State:      ContainerStatus(string(inspect.State.Status)),
		Status:     string(inspect.State.Status),
		Created:    createdTime,
		StartedAt:  startedTime,
		FinishedAt: finishedTime,
		ExitCode:   int(inspect.State.ExitCode),
		Error:      inspect.State.Error,
		Labels:     inspect.Config.Labels,
		NetworkIPs: networkIPs,
		MacAddress: macAddress,
		Gateway:    gateway,
		IPAddress:  ipAddress,
	}, nil
}

// Events streams container events.
func (c *Client) Events(ctx context.Context, filters client.Filters) (<-chan ContainerEvent, <-chan error) {
	eventChan := make(chan ContainerEvent, 100)
	errChan := make(chan error, 1)

	opts := client.EventsListOptions{
		Filters: filters,
	}

	go func() {
		defer close(eventChan)
		defer close(errChan)

		result := c.cli.Events(ctx, opts)

		for {
			select {
			case <-ctx.Done():
				return
			case err := <-result.Err:
				if err != nil {
					errChan <- fmt.Errorf("events error: %w", err)
					return
				}
				return
			case event, ok := <-result.Messages:
				if !ok {
					return
				}
				eventChan <- convertEvent(event)
			}
		}
	}()

	return eventChan, errChan
}

// Commit creates an image from a container.
func (c *Client) Commit(ctx context.Context, containerID string, opts CommitContainerOptions) (*CommitContainerResult, error) {
	c.logger.Debug("Committing container to image",
		logs.String("container", containerID),
		logs.String("repo", opts.Repository),
	)

	options := client.ContainerCommitOptions{
		Reference: opts.Repository + ":" + opts.Tag,
		Author:    opts.Author,
		Comment:   opts.Message,
		NoPause:   !opts.Pause,
		Changes:   opts.Changes,
	}

	resp, err := c.cli.ContainerCommit(ctx, containerID, options)
	if err != nil {
		return nil, fmt.Errorf("failed to commit container: %w", err)
	}

	c.logger.Info("Container committed",
		logs.String("container", containerID),
		logs.String("image", resp.ID),
	)

	return &CommitContainerResult{
		ImageID: resp.ID,
	}, nil
}

// Helper functions

func convertMounts(mounts []Mount) []mount.Mount {
	if mounts == nil {
		return nil
	}
	result := make([]mount.Mount, 0, len(mounts))
	for _, m := range mounts {
		result = append(result, mount.Mount{
			Type:     mount.Type(m.Type),
			Source:   m.Source,
			Target:   m.Target,
			ReadOnly: m.ReadOnly,
		})
	}
	return result
}

func convertEnvMap(env map[string]string) []string {
	if env == nil {
		return nil
	}
	result := make([]string, 0, len(env))
	for k, v := range env {
		result = append(result, k+"="+v)
	}
	return result
}

func convertDNSAddresses(dns []string) []netip.Addr {
	if dns == nil {
		return nil
	}
	result := make([]netip.Addr, 0, len(dns))
	for _, d := range dns {
		addr, err := netip.ParseAddr(d)
		if err == nil {
			result = append(result, addr)
		}
	}
	return result
}

func convertPortBindings(portMappings []PortMapping, hostPorts map[uint16]uint16) network.PortMap {
	portMap := make(network.PortMap)

	for _, pm := range portMappings {
		port := network.MustParsePort(fmt.Sprintf("%d/%s", pm.ContainerPort, pm.Protocol))
		if portMap[port] == nil {
			portMap[port] = []network.PortBinding{}
		}
		hostIP := pm.HostIP
		if hostIP == "" {
			hostIP = "0.0.0.0"
		}
		ip, _ := netip.ParseAddr(hostIP)
		portMap[port] = append(portMap[port], network.PortBinding{
			HostIP:   ip,
			HostPort: fmt.Sprintf("%d", pm.HostPort),
		})
	}

	for containerPort, hostPort := range hostPorts {
		port := network.MustParsePort(fmt.Sprintf("%d/tcp", containerPort))
		if portMap[port] == nil {
			portMap[port] = []network.PortBinding{}
		}
		ip, _ := netip.ParseAddr("0.0.0.0")
		portMap[port] = append(portMap[port], network.PortBinding{
			HostIP:   ip,
			HostPort: fmt.Sprintf("%d", hostPort),
		})
	}

	return portMap
}

func convertExposedPorts(portMappings []PortMapping, hostPorts map[uint16]uint16) network.PortSet {
	portSet := make(network.PortSet)

	for _, pm := range portMappings {
		port := network.MustParsePort(fmt.Sprintf("%d/%s", pm.ContainerPort, pm.Protocol))
		portSet[port] = struct{}{}
	}

	for containerPort := range hostPorts {
		port := network.MustParsePort(fmt.Sprintf("%d/tcp", containerPort))
		portSet[port] = struct{}{}
	}

	return portSet
}

func convertAdvancedMounts(bindMounts []BindMount, volumeMounts []VolumeMount, tmpfsMounts []TmpfsMount) []mount.Mount {
	var result []mount.Mount

	for _, bm := range bindMounts {
		result = append(result, mount.Mount{
			Type:     mount.TypeBind,
			Source:   bm.Source,
			Target:   bm.Target,
			ReadOnly: bm.ReadOnly,
		})
	}

	for _, vm := range volumeMounts {
		result = append(result, mount.Mount{
			Type:     mount.TypeVolume,
			Source:   vm.Source,
			Target:   vm.Target,
			ReadOnly: vm.ReadOnly,
		})
	}

	for _, tm := range tmpfsMounts {
		result = append(result, mount.Mount{
			Type:   mount.TypeTmpfs,
			Target: tm.Target,
			TmpfsOptions: &mount.TmpfsOptions{
				SizeBytes: tm.Size,
			},
		})
	}

	return result
}

func convertDeviceRequests(deviceRequests []DeviceRequest) []container.DeviceRequest {
	if deviceRequests == nil {
		return nil
	}
	result := make([]container.DeviceRequest, 0, len(deviceRequests))
	for _, dr := range deviceRequests {
		result = append(result, container.DeviceRequest{
			Driver:       dr.Driver,
			Count:        dr.Count,
			Capabilities: dr.Capabilities,
			DeviceIDs:    dr.DeviceIDs,
		})
	}
	return result
}

func convertUlimits(ulimits []Ulimit) []*container.Ulimit {
	if ulimits == nil {
		return nil
	}
	result := make([]*container.Ulimit, 0, len(ulimits))
	for _, u := range ulimits {
		result = append(result, &container.Ulimit{
			Name: u.Name,
			Hard: u.Hard,
			Soft: u.Soft,
		})
	}
	return result
}

func convertEvent(e events.Message) ContainerEvent {
	return ContainerEvent{
		Type:   string(e.Type),
		Action: string(e.Action),
		Actor: struct {
			ID         string
			Attributes map[string]string
		}{
			ID:         e.Actor.ID,
			Attributes: e.Actor.Attributes,
		},
		Time:     e.Time,
		TimeNano: e.TimeNano,
	}
}
