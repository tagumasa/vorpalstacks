// Package audit provides AWS CloudTrail audit logging functionality for vorpalstacks.
package audit

import (
	"os"
	"strconv"
)

// AuditConfig contains configuration for audit logging.
type AuditConfig struct {
	Enabled bool
}

// LoadConfig loads the audit configuration from environment variables.
func LoadConfig() *AuditConfig {
	return &AuditConfig{
		Enabled: getEnvBool("VS_AUDIT_ENABLED", false),
	}
}

func getEnvBool(key string, defaultValue bool) bool {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	b, err := strconv.ParseBool(val)
	if err != nil {
		return defaultValue
	}
	return b
}

var defaultConfig *AuditConfig

// DefaultConfig returns the default audit configuration.
func DefaultConfig() *AuditConfig {
	if defaultConfig == nil {
		defaultConfig = LoadConfig()
	}
	return defaultConfig
}
