// Package mobyclient provides Docker/Moby API client operations for networks.
package mobyclient

import (
	"context"
	"fmt"

	"vorpalstacks/internal/core/logs"

	"github.com/moby/moby/api/types/network"
	"github.com/moby/moby/client"
)

// CreateNetwork creates a network.
func (c *Client) CreateNetwork(ctx context.Context, name string, driver string) (string, error) {
	c.logger.Info("Creating network", logs.String("name", name), logs.String("driver", driver))

	options := client.NetworkCreateOptions{
		Driver: driver,
	}

	resp, err := c.cli.NetworkCreate(ctx, name, options)
	if err != nil {
		return "", fmt.Errorf("failed to create network: %w", err)
	}

	c.logger.Info("Network created",
		logs.String("id", resp.ID),
		logs.String("name", name),
	)

	return resp.ID, nil
}

// CreateNetworkWithOptions creates a network with options.
func (c *Client) CreateNetworkWithOptions(ctx context.Context, name string, opts NetworkCreateOptions) (string, error) {
	c.logger.Info("Creating network with options", logs.String("name", name))

	options := client.NetworkCreateOptions{
		Driver:     opts.Driver,
		Internal:   opts.Internal,
		Attachable: opts.Attachable,
		Labels:     opts.Labels,
	}

	resp, err := c.cli.NetworkCreate(ctx, name, options)
	if err != nil {
		return "", fmt.Errorf("failed to create network with options: %w", err)
	}

	c.logger.Info("Network created with options",
		logs.String("id", resp.ID),
		logs.String("name", name),
	)

	return resp.ID, nil
}

// ConnectNetwork connects a container to a network.
func (c *Client) ConnectNetwork(ctx context.Context, networkID, containerID string) error {
	c.logger.Debug("Connecting container to network",
		logs.String("container", containerID),
		logs.String("network", networkID),
	)

	options := client.NetworkConnectOptions{
		Container:      containerID,
		EndpointConfig: &network.EndpointSettings{},
	}

	_, err := c.cli.NetworkConnect(ctx, networkID, options)
	if err != nil {
		return fmt.Errorf("failed to connect network: %w", err)
	}

	c.logger.Info("Container connected to network",
		logs.String("container", containerID),
		logs.String("network", networkID),
	)
	return nil
}

// ConnectNetworkWithEndpointConfig connects a container to a network with endpoint configuration.
func (c *Client) ConnectNetworkWithEndpointConfig(ctx context.Context, networkID, containerID string, endpointConfig *network.EndpointSettings) error {
	c.logger.Debug("Connecting container to network with endpoint config",
		logs.String("container", containerID),
		logs.String("network", networkID),
	)

	options := client.NetworkConnectOptions{
		Container:      containerID,
		EndpointConfig: endpointConfig,
	}

	_, err := c.cli.NetworkConnect(ctx, networkID, options)
	if err != nil {
		return fmt.Errorf("failed to connect network with endpoint config: %w", err)
	}

	c.logger.Info("Container connected to network with endpoint config",
		logs.String("container", containerID),
		logs.String("network", networkID),
	)
	return nil
}

// DisconnectNetwork disconnects a container from a network.
func (c *Client) DisconnectNetwork(ctx context.Context, networkID, containerID string, force bool) error {
	c.logger.Debug("Disconnecting container from network",
		logs.String("container", containerID),
		logs.String("network", networkID),
	)

	options := client.NetworkDisconnectOptions{
		Container: containerID,
		Force:     force,
	}

	_, err := c.cli.NetworkDisconnect(ctx, networkID, options)
	if err != nil {
		return fmt.Errorf("failed to disconnect network: %w", err)
	}

	c.logger.Info("Container disconnected from network",
		logs.String("container", containerID),
		logs.String("network", networkID),
	)
	return nil
}

// RemoveNetwork removes a network.
func (c *Client) RemoveNetwork(ctx context.Context, networkID string) error {
	c.logger.Debug("Removing network", logs.String("id", networkID))

	options := client.NetworkRemoveOptions{}

	_, err := c.cli.NetworkRemove(ctx, networkID, options)
	if err != nil {
		return fmt.Errorf("failed to remove network: %w", err)
	}

	c.logger.Info("Network removed", logs.String("id", networkID))
	return nil
}

// ListNetworks lists all networks.
func (c *Client) ListNetworks(ctx context.Context) ([]NetworkInfo, error) {
	result, err := c.cli.NetworkList(ctx, client.NetworkListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list networks: %w", err)
	}

	networks := result.Items
	networkList := make([]NetworkInfo, 0, len(networks))
	for _, net := range networks {
		subnets := make([]SubnetInfo, 0, len(net.IPAM.Config))
		for _, ipam := range net.IPAM.Config {
			subnets = append(subnets, SubnetInfo{
				Subnet:  ipam.Subnet.String(),
				Gateway: ipam.Gateway.String(),
			})
		}

		networkList = append(networkList, NetworkInfo{
			ID:      net.ID,
			Name:    net.Name,
			Driver:  net.Driver,
			Scope:   net.Scope,
			Subnets: subnets,
		})
	}

	return networkList, nil
}

// ListNetworksWithFilters lists networks with filters.
func (c *Client) ListNetworksWithFilters(ctx context.Context, filters client.Filters) ([]NetworkInfo, error) {
	options := client.NetworkListOptions{
		Filters: filters,
	}

	result, err := c.cli.NetworkList(ctx, options)
	if err != nil {
		return nil, fmt.Errorf("failed to list networks with filters: %w", err)
	}

	networks := result.Items
	networkList := make([]NetworkInfo, 0, len(networks))
	for _, net := range networks {
		subnets := make([]SubnetInfo, 0, len(net.IPAM.Config))
		for _, ipam := range net.IPAM.Config {
			subnets = append(subnets, SubnetInfo{
				Subnet:  ipam.Subnet.String(),
				Gateway: ipam.Gateway.String(),
			})
		}

		networkList = append(networkList, NetworkInfo{
			ID:      net.ID,
			Name:    net.Name,
			Driver:  net.Driver,
			Scope:   net.Scope,
			Subnets: subnets,
		})
	}

	return networkList, nil
}

// InspectNetwork inspects a network.
func (c *Client) InspectNetwork(ctx context.Context, networkID string) (*NetworkInfo, error) {
	options := client.NetworkInspectOptions{}

	result, err := c.cli.NetworkInspect(ctx, networkID, options)
	if err != nil {
		return nil, fmt.Errorf("failed to inspect network: %w", err)
	}

	net := result.Network

	subnets := make([]SubnetInfo, 0, len(net.IPAM.Config))
	for _, ipam := range net.IPAM.Config {
		subnets = append(subnets, SubnetInfo{
			Subnet:  ipam.Subnet.String(),
			Gateway: ipam.Gateway.String(),
		})
	}

	return &NetworkInfo{
		ID:      net.ID,
		Name:    net.Name,
		Driver:  net.Driver,
		Scope:   net.Scope,
		Subnets: subnets,
	}, nil
}

// PruneNetworks removes unused networks.
func (c *Client) PruneNetworks(ctx context.Context, filters client.Filters) (int64, error) {
	c.logger.Info("Pruning networks")

	options := client.NetworkPruneOptions{
		Filters: filters,
	}

	result, err := c.cli.NetworkPrune(ctx, options)
	if err != nil {
		return 0, fmt.Errorf("failed to prune networks: %w", err)
	}

	c.logger.Info("Networks pruned",
		logs.Int("networks_deleted", len(result.Report.NetworksDeleted)),
	)

	return 0, nil
}
