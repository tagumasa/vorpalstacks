// Package mobyclient provides Docker/Moby API client operations for images.
package mobyclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"
	"vorpalstacks/internal/core/logs"

	"github.com/moby/moby/client"
)

// PullImage pulls an image from registry.
func (c *Client) PullImage(ctx context.Context, imageName string) error {
	c.logger.Info("Pulling image", logs.String("image", imageName))

	reader, err := c.cli.ImagePull(ctx, imageName, client.ImagePullOptions{})
	if err != nil {
		return fmt.Errorf("failed to pull image: %w", err)
	}
	defer reader.Close()

	// Read and log progress
	decoder := json.NewDecoder(reader)
	for {
		var progress struct {
			Status   string `json:"status"`
			Progress string `json:"progress"`
			ID       string `json:"id"`
		}
		if err := decoder.Decode(&progress); err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("error reading pull progress: %w", err)
		}

		c.logger.Debug("Pull progress",
			logs.String("id", progress.ID),
			logs.String("status", progress.Status),
			logs.String("progress", progress.Progress),
		)
	}

	c.logger.Info("Image pulled successfully", logs.String("image", imageName))
	return nil
}

// PullImageWithAuth pulls an image with authentication.
func (c *Client) PullImageWithAuth(ctx context.Context, imageName string, authConfig string) error {
	c.logger.Info("Pulling image with auth", logs.String("image", imageName))

	options := client.ImagePullOptions{
		RegistryAuth: authConfig,
	}

	reader, err := c.cli.ImagePull(ctx, imageName, options)
	if err != nil {
		return fmt.Errorf("failed to pull image with auth: %w", err)
	}
	defer reader.Close()

	// Read and log progress
	decoder := json.NewDecoder(reader)
	for {
		var progress struct {
			Status   string `json:"status"`
			Progress string `json:"progress"`
			ID       string `json:"id"`
		}
		if err := decoder.Decode(&progress); err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("error reading pull progress: %w", err)
		}

		c.logger.Debug("Pull progress",
			logs.String("id", progress.ID),
			logs.String("status", progress.Status),
			logs.String("progress", progress.Progress),
		)
	}

	c.logger.Info("Image pulled successfully with auth", logs.String("image", imageName))
	return nil
}

// PullImageStream pulls an image and returns a stream reader.
func (c *Client) PullImageStream(ctx context.Context, imageName string) (io.ReadCloser, error) {
	c.logger.Info("Pulling image stream", logs.String("image", imageName))

	reader, err := c.cli.ImagePull(ctx, imageName, client.ImagePullOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to pull image stream: %w", err)
	}

	return reader, nil
}

// ImageExists checks if an image exists locally.
func (c *Client) ImageExists(ctx context.Context, imageName string) (bool, error) {
	_, err := c.cli.ImageInspect(ctx, imageName)
	if err != nil {
		if strings.Contains(err.Error(), "no such image") {
			return false, nil
		}
		return false, fmt.Errorf("failed to inspect image: %w", err)
	}
	return true, nil
}

// GetImage gets image information.
func (c *Client) GetImage(ctx context.Context, imageName string) (*ImageInfo, error) {
	inspect, err := c.cli.ImageInspect(ctx, imageName)
	if err != nil {
		return nil, fmt.Errorf("failed to inspect image: %w", err)
	}

	createdTime, err := time.Parse(time.RFC3339, inspect.Created)
	if err != nil {
		createdTime = time.Now()
	}

	return &ImageInfo{
		ID:       inspect.ID,
		RepoTags: inspect.RepoTags,
		Size:     inspect.Size,
		Created:  createdTime,
	}, nil
}

// ListImages lists all images.
func (c *Client) ListImages(ctx context.Context) ([]ImageInfo, error) {
	images, err := c.cli.ImageList(ctx, client.ImageListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list images: %w", err)
	}

	result := make([]ImageInfo, 0, len(images.Items))
	for _, img := range images.Items {
		createdTime := time.Unix(img.Created, 0)

		result = append(result, ImageInfo{
			ID:       img.ID,
			RepoTags: img.RepoTags,
			Size:     img.Size,
			Created:  createdTime,
		})
	}

	return result, nil
}

// ListImagesWithFilters lists images with filters.
func (c *Client) ListImagesWithFilters(ctx context.Context, filters client.Filters) ([]ImageInfo, error) {
	options := client.ImageListOptions{
		Filters: filters,
	}

	images, err := c.cli.ImageList(ctx, options)
	if err != nil {
		return nil, fmt.Errorf("failed to list images with filters: %w", err)
	}

	result := make([]ImageInfo, 0, len(images.Items))
	for _, img := range images.Items {
		createdTime := time.Unix(img.Created, 0)

		result = append(result, ImageInfo{
			ID:       img.ID,
			RepoTags: img.RepoTags,
			Size:     img.Size,
			Created:  createdTime,
		})
	}

	return result, nil
}

// RemoveImage removes an image.
func (c *Client) RemoveImage(ctx context.Context, imageName string, force bool) error {
	c.logger.Debug("Removing image", logs.String("image", imageName))

	options := client.ImageRemoveOptions{
		Force:         force,
		PruneChildren: true,
	}

	_, err := c.cli.ImageRemove(ctx, imageName, options)
	if err != nil {
		return fmt.Errorf("failed to remove image: %w", err)
	}

	c.logger.Info("Image removed", logs.String("image", imageName))
	return nil
}

// BuildImage builds an image from a Dockerfile.
func (c *Client) BuildImage(ctx context.Context, buildContext io.Reader, options client.ImageBuildOptions) (string, error) {
	c.logger.Info("Building image", logs.String("context", options.Dockerfile))

	resp, err := c.cli.ImageBuild(ctx, buildContext, options)
	if err != nil {
		return "", fmt.Errorf("failed to build image: %w", err)
	}
	defer resp.Body.Close()

	// Read and log progress
	decoder := json.NewDecoder(resp.Body)
	var imageID string
	for {
		var progress struct {
			Stream string `json:"stream"`
			Aux    struct {
				ID string `json:"ID"`
			} `json:"aux"`
		}
		if err := decoder.Decode(&progress); err != nil {
			if err == io.EOF {
				break
			}
			return "", fmt.Errorf("error reading build progress: %w", err)
		}

		if progress.Stream != "" {
			c.logger.Debug("Build progress", logs.String("stream", progress.Stream))
		}

		if progress.Aux.ID != "" {
			imageID = progress.Aux.ID
		}
	}

	c.logger.Info("Image built successfully", logs.String("id", imageID))
	return imageID, nil
}

// TagImage tags an image.
func (c *Client) TagImage(ctx context.Context, source, target string) error {
	c.logger.Debug("Tagging image",
		logs.String("source", source),
		logs.String("target", target),
	)

	options := client.ImageTagOptions{
		Source: source,
		Target: target,
	}
	_, err := c.cli.ImageTag(ctx, options)
	if err != nil {
		return fmt.Errorf("failed to tag image: %w", err)
	}

	c.logger.Info("Image tagged successfully",
		logs.String("source", source),
		logs.String("target", target),
	)
	return nil
}

// PushImage pushes an image to a registry.
func (c *Client) PushImage(ctx context.Context, imageName string) error {
	c.logger.Info("Pushing image", logs.String("image", imageName))

	reader, err := c.cli.ImagePush(ctx, imageName, client.ImagePushOptions{})
	if err != nil {
		return fmt.Errorf("failed to push image: %w", err)
	}
	defer reader.Close()

	// Read and log progress
	decoder := json.NewDecoder(reader)
	for {
		var progress struct {
			Status   string `json:"status"`
			Progress string `json:"progress"`
			ID       string `json:"id"`
		}
		if err := decoder.Decode(&progress); err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("error reading push progress: %w", err)
		}

		c.logger.Debug("Push progress",
			logs.String("id", progress.ID),
			logs.String("status", progress.Status),
			logs.String("progress", progress.Progress),
		)
	}

	c.logger.Info("Image pushed successfully", logs.String("image", imageName))
	return nil
}

// PruneImages removes unused images.
func (c *Client) PruneImages(ctx context.Context, filters client.Filters) (int64, error) {
	c.logger.Info("Pruning images")

	pruneOptions := client.ImagePruneOptions{
		Filters: filters,
	}
	result, err := c.cli.ImagePrune(ctx, pruneOptions)
	if err != nil {
		return 0, fmt.Errorf("failed to prune images: %w", err)
	}

	spaceReclaimed := int64(result.Report.SpaceReclaimed)
	c.logger.Info("Images pruned",
		logs.Int64("space_reclaimed", spaceReclaimed),
		logs.Int("images_deleted", len(result.Report.ImagesDeleted)),
	)

	return spaceReclaimed, nil
}
