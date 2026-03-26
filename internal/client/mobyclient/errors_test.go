package mobyclient

import (
	"errors"
	"testing"
)

func TestContainerError(t *testing.T) {
	err := &ContainerError{
		ContainerID: "test-container",
		Operation:   "start",
		Err:         errors.New("failed to start"),
	}

	if err.Error() != "container test-container: start failed: failed to start" {
		t.Errorf("Unexpected error message: %s", err.Error())
	}

	if !errors.Is(err, err.Err) {
		t.Error("Expected Unwrap to return underlying error")
	}
}

func TestImageError(t *testing.T) {
	err := &ImageError{
		ImageName: "test-image",
		Operation: "pull",
		Err:       errors.New("failed to pull"),
	}

	if err.Error() != "image test-image: pull failed: failed to pull" {
		t.Errorf("Unexpected error message: %s", err.Error())
	}

	if !errors.Is(err, err.Err) {
		t.Error("Expected Unwrap to return underlying error")
	}
}

func TestNetworkError(t *testing.T) {
	err := &NetworkError{
		NetworkID: "test-network",
		Operation: "create",
		Err:       errors.New("failed to create"),
	}

	if err.Error() != "network test-network: create failed: failed to create" {
		t.Errorf("Unexpected error message: %s", err.Error())
	}

	if !errors.Is(err, err.Err) {
		t.Error("Expected Unwrap to return underlying error")
	}
}

func TestVolumeError(t *testing.T) {
	err := &VolumeError{
		VolumeName: "test-volume",
		Operation:  "create",
		Err:        errors.New("failed to create"),
	}

	if err.Error() != "volume test-volume: create failed: failed to create" {
		t.Errorf("Unexpected error message: %s", err.Error())
	}

	if !errors.Is(err, err.Err) {
		t.Error("Expected Unwrap to return underlying error")
	}
}

func TestErrorVariables(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{
			name: "ErrContainerNotFound",
			err:  ErrContainerNotFound,
		},
		{
			name: "ErrImageNotFound",
			err:  ErrImageNotFound,
		},
		{
			name: "ErrNetworkNotFound",
			err:  ErrNetworkNotFound,
		},
		{
			name: "ErrVolumeNotFound",
			err:  ErrVolumeNotFound,
		},
		{
			name: "ErrContainerAlreadyExists",
			err:  ErrContainerAlreadyExists,
		},
		{
			name: "ErrImageAlreadyExists",
			err:  ErrImageAlreadyExists,
		},
		{
			name: "ErrContainerRunning",
			err:  ErrContainerRunning,
		},
		{
			name: "ErrContainerStopped",
			err:  ErrContainerStopped,
		},
		{
			name: "ErrInvalidConfig",
			err:  ErrInvalidConfig,
		},
		{
			name: "ErrConnectionFailed",
			err:  ErrConnectionFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err == nil {
				t.Errorf("Expected error, got nil")
			}
			if tt.err.Error() == "" {
				t.Errorf("Expected non-empty error message")
			}
		})
	}
}
