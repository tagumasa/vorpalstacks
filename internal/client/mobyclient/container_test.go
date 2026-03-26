package mobyclient

import (
	"testing"

	"github.com/moby/moby/api/types/events"
	"github.com/moby/moby/api/types/network"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConvertMounts(t *testing.T) {
	mounts := []Mount{
		{
			Type:     "bind",
			Source:   "/host/path",
			Target:   "/container/path",
			ReadOnly: false,
		},
	}

	result := convertMounts(mounts)

	require.Len(t, result, 1)
	assert.Equal(t, "bind", string(result[0].Type))
	assert.Equal(t, "/host/path", result[0].Source)
	assert.Equal(t, "/container/path", result[0].Target)
	assert.False(t, result[0].ReadOnly)
}

func TestConvertMounts_Empty(t *testing.T) {
	result := convertMounts([]Mount{})
	assert.Empty(t, result)
}

func TestConvertMounts_Nil(t *testing.T) {
	result := convertMounts(nil)
	assert.Nil(t, result)
}

func TestConvertEnvMap(t *testing.T) {
	env := map[string]string{
		"KEY1": "value1",
		"KEY2": "value2",
	}

	result := convertEnvMap(env)

	require.Len(t, result, 2)
	assert.Contains(t, result, "KEY1=value1")
	assert.Contains(t, result, "KEY2=value2")
}

func TestConvertEnvMap_Empty(t *testing.T) {
	result := convertEnvMap(map[string]string{})
	assert.Empty(t, result)
}

func TestConvertEnvMap_Nil(t *testing.T) {
	result := convertEnvMap(nil)
	assert.Nil(t, result)
}

func TestConvertDNSAddresses(t *testing.T) {
	dns := []string{"8.8.8.8", "8.8.4.4"}

	result := convertDNSAddresses(dns)

	require.Len(t, result, 2)
	assert.Equal(t, "8.8.8.8", result[0].String())
	assert.Equal(t, "8.8.4.4", result[1].String())
}

func TestConvertDNSAddresses_Invalid(t *testing.T) {
	dns := []string{"invalid", "1.1.1.1"}

	result := convertDNSAddresses(dns)

	require.Len(t, result, 1)
	assert.Equal(t, "1.1.1.1", result[0].String())
}

func TestConvertDNSAddresses_Empty(t *testing.T) {
	result := convertDNSAddresses([]string{})
	assert.Empty(t, result)
}

func TestConvertDNSAddresses_Nil(t *testing.T) {
	result := convertDNSAddresses(nil)
	assert.Nil(t, result)
}

func TestConvertPortBindings(t *testing.T) {
	portMappings := []PortMapping{
		{ContainerPort: 80, HostPort: 8080, Protocol: PortProtocolTCP, HostIP: "0.0.0.0"},
		{ContainerPort: 443, HostPort: 8443, Protocol: PortProtocolTCP, HostIP: "127.0.0.1"},
	}
	hostPorts := map[uint16]uint16{3000: 3000}

	result := convertPortBindings(portMappings, hostPorts)

	port80 := network.MustParsePort("80/tcp")
	port443 := network.MustParsePort("443/tcp")
	port3000 := network.MustParsePort("3000/tcp")

	assert.Contains(t, result, port80)
	assert.Contains(t, result, port443)
	assert.Contains(t, result, port3000)
	assert.Equal(t, "8080", result[port80][0].HostPort)
	assert.Equal(t, "8443", result[port443][0].HostPort)
}

func TestConvertPortBindings_EmptyHostIP(t *testing.T) {
	portMappings := []PortMapping{
		{ContainerPort: 80, HostPort: 8080, Protocol: PortProtocolTCP},
	}

	result := convertPortBindings(portMappings, nil)

	port80 := network.MustParsePort("80/tcp")
	assert.Equal(t, "0.0.0.0", result[port80][0].HostIP.String())
}

func TestConvertExposedPorts(t *testing.T) {
	portMappings := []PortMapping{
		{ContainerPort: 80, HostPort: 8080, Protocol: PortProtocolTCP},
		{ContainerPort: 53, HostPort: 5353, Protocol: PortProtocolUDP},
	}
	hostPorts := map[uint16]uint16{3000: 3000}

	result := convertExposedPorts(portMappings, hostPorts)

	port80 := network.MustParsePort("80/tcp")
	port53 := network.MustParsePort("53/udp")
	port3000 := network.MustParsePort("3000/tcp")

	assert.Contains(t, result, port80)
	assert.Contains(t, result, port53)
	assert.Contains(t, result, port3000)
}

func TestConvertAdvancedMounts(t *testing.T) {
	bindMounts := []BindMount{
		{Source: "/host/path", Target: "/container/path", ReadOnly: true},
	}
	volumeMounts := []VolumeMount{
		{Source: "myvolume", Target: "/data", ReadOnly: false},
	}
	tmpfsMounts := []TmpfsMount{
		{Target: "/tmp", Size: 1024 * 1024},
	}

	result := convertAdvancedMounts(bindMounts, volumeMounts, tmpfsMounts)

	require.Len(t, result, 3)
	assert.Equal(t, "bind", string(result[0].Type))
	assert.Equal(t, "volume", string(result[1].Type))
	assert.Equal(t, "tmpfs", string(result[2].Type))
}

func TestConvertAdvancedMounts_Empty(t *testing.T) {
	result := convertAdvancedMounts(nil, nil, nil)
	assert.Empty(t, result)
}

func TestConvertDeviceRequests(t *testing.T) {
	deviceRequests := []DeviceRequest{
		{
			Driver:       "nvidia",
			Count:        -1,
			Capabilities: [][]string{{"gpu"}},
			DeviceIDs:    []string{"0", "1"},
		},
	}

	result := convertDeviceRequests(deviceRequests)

	require.Len(t, result, 1)
	assert.Equal(t, "nvidia", result[0].Driver)
	assert.Equal(t, -1, result[0].Count)
}

func TestConvertDeviceRequests_Empty(t *testing.T) {
	result := convertDeviceRequests([]DeviceRequest{})
	assert.Empty(t, result)
}

func TestConvertDeviceRequests_Nil(t *testing.T) {
	result := convertDeviceRequests(nil)
	assert.Nil(t, result)
}

func TestConvertUlimits(t *testing.T) {
	ulimits := []Ulimit{
		{Name: "nofile", Hard: 65536, Soft: 65536},
		{Name: "nproc", Hard: 4096, Soft: 2048},
	}

	result := convertUlimits(ulimits)

	require.Len(t, result, 2)
	assert.Equal(t, "nofile", result[0].Name)
	assert.Equal(t, int64(65536), result[0].Hard)
	assert.Equal(t, int64(2048), result[1].Soft)
}

func TestConvertUlimits_Empty(t *testing.T) {
	result := convertUlimits([]Ulimit{})
	assert.Empty(t, result)
}

func TestConvertUlimits_Nil(t *testing.T) {
	result := convertUlimits(nil)
	assert.Nil(t, result)
}

func TestConvertEvent(t *testing.T) {
	e := events.Message{
		Type:   events.ContainerEventType,
		Action: events.ActionStart,
		Actor: events.Actor{
			ID: "container-id",
			Attributes: map[string]string{
				"image": "nginx",
			},
		},
		Time:     1234567890,
		TimeNano: 1234567890000000000,
	}

	result := convertEvent(e)

	assert.Equal(t, "container", result.Type)
	assert.Equal(t, "start", result.Action)
	assert.Equal(t, "container-id", result.Actor.ID)
	assert.Equal(t, "nginx", result.Actor.Attributes["image"])
	assert.Equal(t, int64(1234567890), result.Time)
}
