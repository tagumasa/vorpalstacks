// Package mobyclient provides Docker/Moby API client operations for volumes.
package mobyclient

import (
	"context"
	"fmt"
	"time"

	"vorpalstacks/internal/core/logs"

	"github.com/moby/moby/client"
)

// CreateVolume creates a volume.
func (c *Client) CreateVolume(ctx context.Context, name string, driver string, labels map[string]string) (string, error) {
	c.logger.Info("Creating volume", logs.String("name", name), logs.String("driver", driver))

	options := client.VolumeCreateOptions{
		Name:   name,
		Driver: driver,
		Labels: labels,
	}

	result, err := c.cli.VolumeCreate(ctx, options)
	if err != nil {
		return "", fmt.Errorf("failed to create volume: %w", err)
	}

	c.logger.Info("Volume created", logs.String("name", result.Volume.Name))
	return result.Volume.Name, nil
}

// CreateVolumeWithOptions creates a volume with options.
func (c *Client) CreateVolumeWithOptions(ctx context.Context, opts VolumeCreateOptions) (string, error) {
	c.logger.Info("Creating volume with options", logs.String("name", opts.Name))

	options := client.VolumeCreateOptions{
		Name:       opts.Name,
		Driver:     opts.Driver,
		DriverOpts: opts.DriverOpts,
		Labels:     opts.Labels,
	}

	result, err := c.cli.VolumeCreate(ctx, options)
	if err != nil {
		return "", fmt.Errorf("failed to create volume with options: %w", err)
	}

	c.logger.Info("Volume created with options", logs.String("name", result.Volume.Name))
	return result.Volume.Name, nil
}

// RemoveVolume removes a volume.
func (c *Client) RemoveVolume(ctx context.Context, volumeID string, force bool) error {
	c.logger.Debug("Removing volume", logs.String("id", volumeID))

	options := client.VolumeRemoveOptions{
		Force: force,
	}

	_, err := c.cli.VolumeRemove(ctx, volumeID, options)
	if err != nil {
		return fmt.Errorf("failed to remove volume: %w", err)
	}

	c.logger.Info("Volume removed", logs.String("id", volumeID))
	return nil
}

// ListVolumes lists all volumes.
func (c *Client) ListVolumes(ctx context.Context) ([]VolumeInfo, error) {
	result, err := c.cli.VolumeList(ctx, client.VolumeListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list volumes: %w", err)
	}

	volumes := result.Items
	volumeList := make([]VolumeInfo, 0, len(volumes))
	for _, vol := range volumes {
		volumeList = append(volumeList, VolumeInfo{
			Name:       vol.Name,
			Driver:     vol.Driver,
			Mountpoint: vol.Mountpoint,
			CreatedAt:  parseTime(vol.CreatedAt),
		})
	}

	return volumeList, nil
}

// ListVolumesWithFilters lists volumes with filters.
func (c *Client) ListVolumesWithFilters(ctx context.Context, filters client.Filters) ([]VolumeInfo, error) {
	options := client.VolumeListOptions{
		Filters: filters,
	}

	result, err := c.cli.VolumeList(ctx, options)
	if err != nil {
		return nil, fmt.Errorf("failed to list volumes with filters: %w", err)
	}

	volumes := result.Items
	volumeList := make([]VolumeInfo, 0, len(volumes))
	for _, vol := range volumes {
		volumeList = append(volumeList, VolumeInfo{
			Name:       vol.Name,
			Driver:     vol.Driver,
			Mountpoint: vol.Mountpoint,
			CreatedAt:  parseTime(vol.CreatedAt),
		})
	}

	return volumeList, nil
}

// InspectVolume inspects a volume.
func (c *Client) InspectVolume(ctx context.Context, volumeID string) (*VolumeInfo, error) {
	options := client.VolumeInspectOptions{}

	result, err := c.cli.VolumeInspect(ctx, volumeID, options)
	if err != nil {
		return nil, fmt.Errorf("failed to inspect volume: %w", err)
	}

	vol := result.Volume

	return &VolumeInfo{
		Name:       vol.Name,
		Driver:     vol.Driver,
		Mountpoint: vol.Mountpoint,
		CreatedAt:  parseTime(vol.CreatedAt),
	}, nil
}

// PruneVolumes removes unused volumes.
func (c *Client) PruneVolumes(ctx context.Context, filters client.Filters) (int64, error) {
	c.logger.Info("Pruning volumes")

	options := client.VolumePruneOptions{
		Filters: filters,
	}

	result, err := c.cli.VolumePrune(ctx, options)
	if err != nil {
		return 0, fmt.Errorf("failed to prune volumes: %w", err)
	}

	spaceReclaimed := int64(result.Report.SpaceReclaimed)
	c.logger.Info("Volumes pruned",
		logs.Int64("space_reclaimed", spaceReclaimed),
		logs.Int("volumes_deleted", len(result.Report.VolumesDeleted)),
	)

	return spaceReclaimed, nil
}

// Helper function to parse time string
func parseTime(timeStr string) time.Time {
	if timeStr == "" {
		return time.Time{}
	}
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return time.Time{}
	}
	return t
}
