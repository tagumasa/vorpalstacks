// Package mobyclient provides Docker/Moby API client types and constants.
package mobyclient

import (
	"time"
)

// Platform represents a container platform.
type Platform string

// Platform constants define supported container platforms.
const (
	PlatformLinuxAmd64 Platform = "linux/amd64"
	PlatformLinuxArm64 Platform = "linux/arm64"
)

// Isolation represents container isolation level.
type Isolation string

// Isolation constants define container isolation levels.
const (
	IsolationDefault Isolation = ""
	IsolationProcess Isolation = "process"
	IsolationHyperV  Isolation = "hyperv"
)

// PortProtocol represents the protocol for a port.
type PortProtocol string

// PortProtocol constants define supported port protocols.
const (
	PortProtocolTCP PortProtocol = "tcp"
	PortProtocolUDP PortProtocol = "udp"
)

// ContainerConfig represents basic container configuration.
type ContainerConfig struct {
	Name       string
	Image      string
	Cmd        []string
	Entrypoint []string
	Env        []string
	Labels     map[string]string
	WorkingDir string

	Memory     int64
	MemorySwap int64
	CPUs       int64

	NetworkMode string
	Hostname    string

	Mounts []Mount

	UseGPU bool
}

// Mount represents a mount point.
type Mount struct {
	Type     string
	Source   string
	Target   string
	ReadOnly bool
}

// PortBinding represents a port binding.
type PortBinding struct {
	HostIP   string
	HostPort string
}

// ContainerInfo represents container information.
type ContainerInfo struct {
	ID         string
	Name       string
	Image      string
	State      string
	Status     string
	Created    time.Time
	Ports      []Port
	Labels     map[string]string
	NetworkIPs map[string]string
}

// PortMapping represents a port mapping between host and container.
type PortMapping struct {
	ContainerPort uint16
	HostPort      uint16
	Protocol      PortProtocol
	HostIP        string
}

// BindMount represents a bind mount.
type BindMount struct {
	Source   string
	Target   string
	ReadOnly bool
}

// TmpfsMount represents a tmpfs mount.
type TmpfsMount struct {
	Target string
	Size   int64
}

// VolumeMount represents a volume mount.
type VolumeMount struct {
	Source   string
	Target   string
	ReadOnly bool
}

// DeviceRequest represents a request for a device.
type DeviceRequest struct {
	Driver       string
	Count        int
	Capabilities [][]string
	DeviceIDs    []string
}

// Ulimit represents a ulimit setting.
type Ulimit struct {
	Name string
	Hard int64
	Soft int64
}

// AdvancedContainerConfig represents advanced container configuration.
type AdvancedContainerConfig struct {
	Name string

	Image      string
	PullImage  bool
	Platform   Platform
	Isolation  Isolation
	Entrypoint []string
	Cmd        []string
	Env        map[string]string
	EnvFiles   []string
	Labels     map[string]string
	WorkingDir string
	Hostname   string
	User       string

	Network      string
	NetworkAlias []string
	DNS          []string
	DNSSearch    []string
	DNSOptions   []string
	ExtraHosts   []string
	PortMappings []PortMapping
	HostPorts    map[uint16]uint16

	BindMounts  []BindMount
	VolumeMount []VolumeMount
	TmpfsMounts []TmpfsMount

	Memory     int64
	MemorySwap int64
	CPUShares  int64
	CPUQuota   int64
	CPUPeriod  int64
	CPUSetCPUs string
	CPUSetMems string
	PidsLimit  int64
	Ulimits    []Ulimit

	Privileged     bool
	ReadonlyRootfs bool
	SecurityOpt    []string
	CapAdd         []string
	CapDrop        []string

	DeviceRequests []DeviceRequest

	AutoRemove    bool
	RestartPolicy string
	StopTimeout   int

	StdinOnce bool
	OpenStdin bool
	TTY       bool

	AdditionalFlags []string
}

// CreateContainerResult represents the result of a container creation.
type CreateContainerResult struct {
	ID       string
	Warnings []string
}

// ContainerStatus represents the status of a container.
type ContainerStatus string

// ContainerStatus constants define possible container statuses.
const (
	ContainerStatusCreated    ContainerStatus = "created"
	ContainerStatusRunning    ContainerStatus = "running"
	ContainerStatusPaused     ContainerStatus = "paused"
	ContainerStatusRestarting ContainerStatus = "restarting"
	ContainerStatusRemoving   ContainerStatus = "removing"
	ContainerStatusExited     ContainerStatus = "exited"
	ContainerStatusDead       ContainerStatus = "dead"
)

// Port represents a port mapping.
type Port struct {
	IP          string
	PrivatePort uint16
	PublicPort  uint16
	Protocol    PortProtocol
}

// DetailedContainerInfo represents detailed information about a container.
type DetailedContainerInfo struct {
	ID         string
	Name       string
	Image      string
	ImageID    string
	State      ContainerStatus
	Status     string
	Created    time.Time
	StartedAt  time.Time
	FinishedAt time.Time
	ExitCode   int
	Error      string
	Ports      []Port
	Labels     map[string]string
	NetworkIPs map[string]string
	MacAddress string
	Gateway    string
	IPAddress  string
}

// ContainerEvent represents an event from a container.
type ContainerEvent struct {
	Type   string
	Action string
	Actor  struct {
		ID         string
		Attributes map[string]string
	}
	Time     int64
	TimeNano int64
}

// ContainerStatsDetailed represents detailed container statistics.
type ContainerStatsDetailed struct {
	CPUUsage      float64
	CPUPercent    float64
	MemoryUsage   float64
	MemoryLimit   float64
	MemoryPercent float64
	NetworkRx     uint64
	NetworkTx     uint64
	BlockRead     uint64
	BlockWrite    uint64
	PIDs          uint64
}

// ImageBuildResult represents the result of an image build.
type ImageBuildResult struct {
	ImageID  string
	Warnings []string
}

// CommitContainerOptions represents options for committing a container.
type CommitContainerOptions struct {
	Repository string
	Tag        string
	Author     string
	Message    string
	Pause      bool
	Changes    []string
}

// CommitContainerResult represents the result of committing a container.
type CommitContainerResult struct {
	ImageID string
}

// ExecConfig represents configuration for executing a command in a container.
type ExecConfig struct {
	Cmd          []string
	AttachStdout bool
	AttachStderr bool
	Env          []string
	WorkingDir   string
	Privileged   bool
}

// ExecResult represents the result of executing a command in a container.
type ExecResult struct {
	ExitCode int
	Stdout   string
	Stderr   string
}

// ContainerStats represents container statistics.
type ContainerStats struct {
	CPUUsage    float64
	MemoryUsage float64
	MemoryLimit float64
	NetworkRx   uint64
	NetworkTx   uint64
}

// ImageInfo represents information about an image.
type ImageInfo struct {
	ID       string
	RepoTags []string
	Size     int64
	Created  time.Time
}

// NetworkInfo represents information about a network.
type NetworkInfo struct {
	ID      string
	Name    string
	Driver  string
	Scope   string
	Subnets []SubnetInfo
}

// SubnetInfo represents information about a subnet.
type SubnetInfo struct {
	Subnet  string
	Gateway string
}

// VolumeInfo represents information about a volume.
type VolumeInfo struct {
	Name       string
	Driver     string
	Mountpoint string
	CreatedAt  time.Time
	Size       int64
}
