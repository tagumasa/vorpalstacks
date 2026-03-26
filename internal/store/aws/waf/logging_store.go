package waf

import (
	"encoding/json"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

const loggingConfigBucketName = "waf_logging_configs"

// LoggingStore provides storage for WAF logging configurations.
type LoggingStore struct {
	*common.BaseStore
}

// NewLoggingStore creates a new logging configuration store.
func NewLoggingStore(store storage.BasicStorage) *LoggingStore {
	return &LoggingStore{
		BaseStore: common.NewBaseStore(store.Bucket(loggingConfigBucketName), "waf"),
	}
}

// Create creates a new logging configuration.
func (s *LoggingStore) Create(resourceArn string, logDestinationConfigs []string, logScope, logType string, loggingFilter *LoggingFilter, managedByFirewallManager bool, redactedFields []*FieldToMatch) (*LoggingConfiguration, error) {
	config := &LoggingConfiguration{
		ResourceArn:              resourceArn,
		LogDestinationConfigs:    logDestinationConfigs,
		LogScope:                 logScope,
		LogType:                  logType,
		LoggingFilter:            loggingFilter,
		ManagedByFirewallManager: managedByFirewallManager,
		RedactedFields:           redactedFields,
		CreatedAt:                time.Now(),
	}

	if config.LogScope == "" {
		config.LogScope = "CUSTOMER"
	}
	if config.LogType == "" {
		config.LogType = "WAF_LOGS"
	}

	if err := s.BaseStore.Put(resourceArn, config); err != nil {
		return nil, NewStoreError("create_logging_configuration", err)
	}
	return config, nil
}

// Get retrieves a logging configuration by resource ARN.
func (s *LoggingStore) Get(resourceArn string) (*LoggingConfiguration, error) {
	var config LoggingConfiguration
	if err := s.BaseStore.Get(resourceArn, &config); err != nil {
		if common.IsNotFound(err) {
			return nil, NewStoreError("get_logging_configuration", ErrNotFound)
		}
		return nil, NewStoreError("get_logging_configuration", err)
	}
	return &config, nil
}

// Update updates an existing logging configuration.
func (s *LoggingStore) Update(resourceArn string, logDestinationConfigs []string, logScope, logType string, loggingFilter *LoggingFilter, managedByFirewallManager bool, redactedFields []*FieldToMatch) (*LoggingConfiguration, error) {
	config, err := s.Get(resourceArn)
	if err != nil {
		return nil, err
	}

	config.LogDestinationConfigs = logDestinationConfigs
	if logScope != "" {
		config.LogScope = logScope
	}
	if logType != "" {
		config.LogType = logType
	}
	config.LoggingFilter = loggingFilter
	config.ManagedByFirewallManager = managedByFirewallManager
	config.RedactedFields = redactedFields

	if err := s.BaseStore.Put(resourceArn, config); err != nil {
		return nil, NewStoreError("update_logging_configuration", err)
	}
	return config, nil
}

// Delete deletes a logging configuration.
func (s *LoggingStore) Delete(resourceArn string) error {
	if err := s.BaseStore.Delete(resourceArn); err != nil {
		return NewStoreError("delete_logging_configuration", err)
	}
	return nil
}

// List returns a paginated list of logging configurations.
func (s *LoggingStore) List(scope string, marker string, maxItems int) (*LoggingConfigurationListResult, error) {
	if maxItems <= 0 {
		maxItems = 100
	}

	var configs []*LoggingConfiguration
	count := 0
	started := marker == ""
	hasMore := false
	var lastKey string

	err := s.ForEach(func(key string, value []byte) error {
		var config LoggingConfiguration
		if err := json.Unmarshal(value, &config); err != nil {
			return err
		}

		if !started {
			if config.ResourceArn == marker {
				started = true
			}
			return nil
		}

		if count < maxItems {
			configs = append(configs, &config)
			lastKey = key
			count++
		} else {
			hasMore = true
		}
		return nil
	})

	if err != nil {
		return nil, NewStoreError("list_logging_configurations", err)
	}

	var nextMarker string
	if hasMore {
		nextMarker = lastKey
	}

	return &LoggingConfigurationListResult{
		LoggingConfigurations: configs,
		IsTruncated:           hasMore,
		NextMarker:            nextMarker,
	}, nil
}

// GetByResourceArn retrieves a logging configuration by resource ARN.
func (s *LoggingStore) GetByResourceArn(resourceArn string) (*LoggingConfiguration, error) {
	var config LoggingConfiguration
	if err := s.BaseStore.Get(resourceArn, &config); err != nil {
		if common.IsNotFound(err) {
			return nil, NewStoreError("get_logging_configuration_by_arn", ErrNotFound)
		}
		return nil, NewStoreError("get_logging_configuration_by_arn", err)
	}
	return &config, nil
}
