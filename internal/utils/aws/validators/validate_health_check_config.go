// Package validators provides AWS-specific utility functions for vorpalstacks.
package validators

import (
	"fmt"
	"net"
)

// HealthCheckConfig represents a Route53 health check configuration.
type HealthCheckConfig struct {
	IPAddress                string
	Port                     int
	Type                     string
	ResourcePath             string
	FullyQualifiedDomainName string
	RequestInterval          int
	FailureThreshold         int
}

// ValidateHealthCheckConfig validates a Route53 health check configuration.
//
// Parameters:
//   - config: The health check configuration to validate
//
// Returns:
//   - error: An error if validation fails, nil otherwise
func ValidateHealthCheckConfig(config interface{}) error {
	if config == nil {
		return fmt.Errorf("health check config cannot be nil")
	}

	healthCheck, ok := config.(*HealthCheckConfig)
	if !ok {
		return fmt.Errorf("invalid health check config type")
	}

	if healthCheck.IPAddress != "" {
		ip := net.ParseIP(healthCheck.IPAddress)
		if ip == nil {
			return fmt.Errorf("invalid IP address: %s", healthCheck.IPAddress)
		}
	}

	if healthCheck.Port < 0 || healthCheck.Port > 65535 {
		return fmt.Errorf("port must be between 0 and 65535")
	}

	if healthCheck.Port == 0 {
		healthCheck.Port = 80
	}

	validTypes := map[string]bool{
		"HTTP":                true,
		"HTTPS":               true,
		"HTTP_STR_MATCH":      true,
		"HTTPS_STR_MATCH":     true,
		"TCP":                 true,
		"SSL":                 true,
		"CALCULATED":          true,
		"CALCULATED_FAILURES": true,
	}
	if !validTypes[healthCheck.Type] {
		return fmt.Errorf("invalid health check type: %s", healthCheck.Type)
	}

	if healthCheck.RequestInterval < 1 || healthCheck.RequestInterval > 30 {
		return fmt.Errorf("request interval must be between 1 and 30 seconds")
	}

	if healthCheck.FailureThreshold < 1 || healthCheck.FailureThreshold > 10 {
		return fmt.Errorf("failure threshold must be between 1 and 10")
	}

	return nil
}

// ParseHealthCheckConfig parses health check config from map.
//
// Parameters:
//   - m: The map to parse from
//
// Returns:
//   - *HealthCheckConfig: The parsed health check config
//   - error: An error if parsing fails
func ParseHealthCheckConfig(m map[string]interface{}) (*HealthCheckConfig, error) {
	config := &HealthCheckConfig{}

	if ip, ok := m["ipAddress"].(string); ok {
		config.IPAddress = ip
	}

	if port, ok := m["port"].(float64); ok {
		config.Port = int(port)
	}

	if healthType, ok := m["type"].(string); ok {
		config.Type = healthType
	}

	if path, ok := m["resourcePath"].(string); ok {
		config.ResourcePath = path
	}

	if domain, ok := m["fullyQualifiedDomainName"].(string); ok {
		config.FullyQualifiedDomainName = domain
	}

	if interval, ok := m["requestInterval"].(float64); ok {
		config.RequestInterval = int(interval)
	}

	if threshold, ok := m["failureThreshold"].(float64); ok {
		config.FailureThreshold = int(threshold)
	}

	return config, nil
}
