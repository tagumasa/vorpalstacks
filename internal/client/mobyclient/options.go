// Package mobyclient provides Docker/Moby API client configuration options.
package mobyclient

// Config represents mobyclient configuration.
type Config struct {
	// Host is the Docker daemon socket (unix:///var/run/docker.sock or tcp://localhost:2375)
	Host string

	// Version is the Docker API version
	Version string

	// TLSVerify enables TLS verification
	TLSVerify bool

	// CertPath is the path to TLS certificates
	CertPath string

	// APIVersion is the API version to use
	APIVersion string

	// Timeout is the timeout for operations
	Timeout int
}

// DefaultConfig returns the default configuration.
func DefaultConfig() Config {
	return Config{
		Host:       "unix:///var/run/docker.sock",
		Version:    "1.41",
		TLSVerify:  false,
		APIVersion: "1.41",
		Timeout:    30,
	}
}

// ContainerCreateOptions represents options for creating a container.
type ContainerCreateOptions struct {
	// Name is the container name
	Name string

	// AutoRemove enables automatic removal of the container when it exits
	AutoRemove bool

	// RemoveVolumes removes volumes associated with the container when it is removed
	RemoveVolumes bool

	// ForceRemove forces removal of the container
	ForceRemove bool

	// StopTimeout is the timeout in seconds to stop the container before removing it
	StopTimeout int
}

// DefaultContainerCreateOptions returns the default container create options.
func DefaultContainerCreateOptions() ContainerCreateOptions {
	return ContainerCreateOptions{
		AutoRemove:    false,
		RemoveVolumes: true,
		ForceRemove:   false,
		StopTimeout:   10,
	}
}

// ContainerStartOptions represents options for starting a container.
type ContainerStartOptions struct {
	// CheckpointID is the checkpoint ID to restore from
	CheckpointID string

	// CheckpointDir is the directory containing the checkpoint
	CheckpointDir string
}

// DefaultContainerStartOptions returns the default container start options.
func DefaultContainerStartOptions() ContainerStartOptions {
	return ContainerStartOptions{}
}

// ContainerStopOptions represents options for stopping a container.
type ContainerStopOptions struct {
	// Timeout is the timeout in seconds to wait for the container to stop
	Timeout int

	// Force forces the container to stop immediately
	Force bool
}

// DefaultContainerStopOptions returns the default container stop options.
func DefaultContainerStopOptions() ContainerStopOptions {
	return ContainerStopOptions{
		Timeout: 10,
		Force:   false,
	}
}

// ImagePullOptions represents options for pulling an image.
type ImagePullOptions struct {
	// All pulls all tags
	All bool

	// Platform specifies the platform in the format os[/arch[/variant]]
	Platform string

	// RegistryAuth is the base64-encoded credentials for the registry
	RegistryAuth string
}

// DefaultImagePullOptions returns the default image pull options.
func DefaultImagePullOptions() ImagePullOptions {
	return ImagePullOptions{
		All:          false,
		Platform:     "",
		RegistryAuth: "",
	}
}

// ImageRemoveOptions represents options for removing an image.
type ImageRemoveOptions struct {
	// Force forces removal of the image
	Force bool

	// PruneChildren removes dangling child images
	PruneChildren bool
}

// DefaultImageRemoveOptions returns the default image remove options.
func DefaultImageRemoveOptions() ImageRemoveOptions {
	return ImageRemoveOptions{
		Force:         false,
		PruneChildren: true,
	}
}

// LogsOptions represents options for getting container logs.
type LogsOptions struct {
	// ShowStdout shows stdout
	ShowStdout bool

	// ShowStderr shows stderr
	ShowStderr bool

	// Timestamps shows timestamps
	Timestamps bool

	// Follow follows the log output
	Follow bool

	// Tail specifies the number of lines to show from the end of the logs
	Tail string

	// Since shows logs since a given time
	Since string

	// Until shows logs until a given time
	Until string
}

// DefaultLogsOptions returns the default logs options.
func DefaultLogsOptions() LogsOptions {
	return LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Timestamps: false,
		Follow:     false,
		Tail:       "all",
		Since:      "",
		Until:      "",
	}
}

// ExecOptions represents options for executing a command in a container.
type ExecOptions struct {
	// AttachStdout attaches stdout
	AttachStdout bool

	// AttachStderr attaches stderr
	AttachStderr bool

	// AttachStdin attaches stdin
	AttachStdin bool

	// Detach detaches from the exec command
	Detach bool

	// TTY allocates a pseudo-TTY
	TTY bool

	// WorkingDir is the working directory
	WorkingDir string

	// Env is the environment variables
	Env []string

	// Privileged gives extended privileges
	Privileged bool

	// User is the user to run the command as
	User string
}

// DefaultExecOptions returns the default exec options.
func DefaultExecOptions() ExecOptions {
	return ExecOptions{
		AttachStdout: true,
		AttachStderr: true,
		AttachStdin:  false,
		Detach:       false,
		TTY:          false,
		WorkingDir:   "",
		Env:          nil,
		Privileged:   false,
		User:         "",
	}
}

// NetworkCreateOptions represents options for creating a network.
type NetworkCreateOptions struct {
	// Driver is the network driver
	Driver string

	// Internal restricts external access to the network
	Internal bool

	// Attachable enables manual container attachment
	Attachable bool

	// Labels are the network labels
	Labels map[string]string

	// EnableIPv6 enables IPv6 networking
	EnableIPv6 bool
}

// DefaultNetworkCreateOptions returns the default network create options.
func DefaultNetworkCreateOptions() NetworkCreateOptions {
	return NetworkCreateOptions{
		Driver:     "bridge",
		Internal:   false,
		Attachable: true,
		Labels:     nil,
		EnableIPv6: false,
	}
}

// VolumeCreateOptions represents options for creating a volume.
type VolumeCreateOptions struct {
	// Driver is the volume driver
	Driver string

	// DriverOpts are driver specific options
	DriverOpts map[string]string

	// Labels are the volume labels
	Labels map[string]string

	// Name is the volume name
	Name string
}

// DefaultVolumeCreateOptions returns the default volume create options.
func DefaultVolumeCreateOptions() VolumeCreateOptions {
	return VolumeCreateOptions{
		Driver:     "local",
		DriverOpts: nil,
		Labels:     nil,
		Name:       "",
	}
}
