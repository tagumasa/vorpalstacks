// Package mobyclient provides Docker/Moby API client error definitions.
package mobyclient

import (
	"errors"
	"fmt"
)

// Error variables for container operations.
var (
	// ErrContainerNotFound is returned when a container is not found.
	ErrContainerNotFound = errors.New("container not found")

	// ErrImageNotFound is returned when an image is not found.
	ErrImageNotFound = errors.New("image not found")

	// ErrNetworkNotFound is returned when a network is not found.
	ErrNetworkNotFound = errors.New("network not found")

	// ErrVolumeNotFound is returned when a volume is not found.
	ErrVolumeNotFound = errors.New("volume not found")

	// ErrContainerAlreadyExists is returned when attempting to create a container
	// that already exists.
	ErrContainerAlreadyExists = errors.New("container already exists")

	// ErrImageAlreadyExists is returned when attempting to create an image
	// that already exists.
	ErrImageAlreadyExists = errors.New("image already exists")

	// ErrContainerRunning is returned when an operation requires a stopped container.
	ErrContainerRunning = errors.New("container is running")

	// ErrContainerStopped is returned when an operation requires a running container.
	ErrContainerStopped = errors.New("container is stopped")

	// ErrInvalidConfig is returned when the configuration is invalid.
	ErrInvalidConfig = errors.New("invalid configuration")

	// ErrConnectionFailed is returned when the connection to the Docker daemon fails.
	ErrConnectionFailed = errors.New("connection to Docker daemon failed")
)

// ContainerError represents an error that occurs during container operations.
type ContainerError struct {
	ContainerID string
	Operation   string
	Err         error
}

// Error returns the error message for ContainerError.
func (e *ContainerError) Error() string {
	return fmt.Sprintf("container %s: %s failed: %v", e.ContainerID, e.Operation, e.Err)
}

// Unwrap returns the underlying error for ContainerError.
func (e *ContainerError) Unwrap() error {
	return e.Err
}

// ImageError represents an error that occurs during image operations.
type ImageError struct {
	ImageName string
	Operation string
	Err       error
}

// Error returns the error message for ImageError.
func (e *ImageError) Error() string {
	return fmt.Sprintf("image %s: %s failed: %v", e.ImageName, e.Operation, e.Err)
}

// Unwrap returns the underlying error for ImageError.
func (e *ImageError) Unwrap() error {
	return e.Err
}

// NetworkError represents an error that occurs during network operations.
type NetworkError struct {
	NetworkID string
	Operation string
	Err       error
}

// Error returns the error message for NetworkError.
func (e *NetworkError) Error() string {
	return fmt.Sprintf("network %s: %s failed: %v", e.NetworkID, e.Operation, e.Err)
}

// Unwrap returns the underlying error for NetworkError.
func (e *NetworkError) Unwrap() error {
	return e.Err
}

// VolumeError represents an error that occurs during volume operations.
type VolumeError struct {
	VolumeName string
	Operation  string
	Err        error
}

// Error returns the error message for VolumeError.
func (e *VolumeError) Error() string {
	return fmt.Sprintf("volume %s: %s failed: %v", e.VolumeName, e.Operation, e.Err)
}

// Unwrap returns the underlying error for VolumeError.
func (e *VolumeError) Unwrap() error {
	return e.Err
}
