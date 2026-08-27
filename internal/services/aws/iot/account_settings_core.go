package iot

import (
	iotstore "vorpalstacks/internal/store/aws/iot"
)

// ---------------------------------------------------------------------------
// Account settings Core (logging, event configurations, encryption).
// The account-level settings persist via the generic-KV store under
// "config/<name>" keys. A missing key means "not yet configured"; the read
// Cores return the documented default shape for that case and propagate
// genuine store errors.
// ---------------------------------------------------------------------------

// SetV2LoggingOptionsInput carries the parsed SetV2LoggingOptions request.
type SetV2LoggingOptionsInput struct {
	RoleArn         string
	DefaultLogLevel string
	DisableAllLogs  bool
}

// setV2LoggingOptionsCore persists the account default V2 logging options.
func (s *IoTService) setV2LoggingOptionsCore(store iotstore.IotStoreInterface, in SetV2LoggingOptionsInput) error {
	defaultLogLevel := in.DefaultLogLevel
	if defaultLogLevel == "" {
		defaultLogLevel = "INFO"
	}
	rec := map[string]interface{}{
		"roleArn":         in.RoleArn,
		"defaultLogLevel": defaultLogLevel,
		"disableAllLogs":  in.DisableAllLogs,
	}
	return store.PutGeneric("config/v2Logging", rec)
}

// getV2LoggingOptionsCore loads the account default V2 logging options and
// returns the persisted record (with the documented default shape when the
// configuration has never been written).
func (s *IoTService) getV2LoggingOptionsCore(store iotstore.IotStoreInterface) (map[string]interface{}, error) {
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists("config/v2Logging", &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return map[string]interface{}{
			"defaultLogLevel": "DISABLED",
			"disableAllLogs":  true,
		}, nil
	}
	if rec["defaultLogLevel"] == nil || rec["defaultLogLevel"] == "" {
		rec["defaultLogLevel"] = "INFO"
	}
	return rec, nil
}

// logLevels is the LogLevel enum member set.
var logLevels = map[string]bool{
	"DEBUG": true, "INFO": true, "ERROR": true, "WARN": true, "DISABLED": true,
}

// setV2LoggingLevelCore validates and persists a per-target V2 logging
// level. The logLevel member must be a LogLevel enum member.
func (s *IoTService) setV2LoggingLevelCore(store iotstore.IotStoreInterface, logTarget map[string]interface{}, logLevel string) error {
	targetType, _ := logTarget["targetType"].(string)
	targetName, _ := logTarget["targetName"].(string)
	if targetType == "" || logLevel == "" {
		return iotstore.ErrMissingParam
	}
	if !logLevels[logLevel] {
		return iotstore.ErrInvalidRequest
	}
	return store.PutGeneric("v2LoggingLevel/"+targetType+"/"+targetName, map[string]interface{}{
		"logTarget": logTarget,
		"logLevel":  logLevel,
	})
}

// deleteV2LoggingLevelCore removes a per-target V2 logging level. Both the
// targetType and targetName members are required by the model.
func (s *IoTService) deleteV2LoggingLevelCore(store iotstore.IotStoreInterface, targetType, targetName string) error {
	if targetType == "" || targetName == "" {
		return iotstore.ErrMissingParam
	}
	return store.DeleteGeneric("v2LoggingLevel/" + targetType + "/" + targetName)
}

// listV2LoggingLevelsCore lists the per-target V2 logging levels, optionally
// restricted to one target type, projected onto the LogTargetConfiguration
// members (logTarget, logLevel).
func (s *IoTService) listV2LoggingLevelsCore(store iotstore.IotStoreInterface, targetType string) ([]map[string]interface{}, error) {
	prefix := "v2LoggingLevel/"
	if targetType != "" {
		prefix = "v2LoggingLevel/" + targetType + "/"
	}
	items, err := store.ListGeneric(prefix)
	if err != nil {
		return nil, err
	}
	configs := make([]map[string]interface{}, 0, len(items))
	for _, rec := range items {
		logTarget, _ := rec["logTarget"].(map[string]interface{})
		if logTarget == nil {
			logTarget = map[string]interface{}{}
		}
		configs = append(configs, map[string]interface{}{
			"logTarget": logTarget,
			"logLevel":  rec["logLevel"],
		})
	}
	return configs, nil
}

// setLoggingOptionsCore validates and persists the legacy logging options.
func (s *IoTService) setLoggingOptionsCore(store iotstore.IotStoreInterface, roleArn, logLevel string) error {
	if roleArn == "" {
		return iotstore.ErrMissingParam
	}
	rec := map[string]interface{}{
		"roleArn":  roleArn,
		"logLevel": logLevel,
	}
	return store.PutGeneric("config/legacyLogging", rec)
}

// getLoggingOptionsCore loads the legacy logging options. AWS returns an
// empty response when logging has never been configured.
func (s *IoTService) getLoggingOptionsCore(store iotstore.IotStoreInterface) (map[string]interface{}, error) {
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists("config/legacyLogging", &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return map[string]interface{}{}, nil
	}
	return map[string]interface{}{
		"roleArn":  rec["roleArn"],
		"logLevel": rec["logLevel"],
	}, nil
}

// updateEventConfigurationsCore merges the incoming configuration attributes
// into the persisted map so that partial updates behave like AWS IoT
// (per-event-type toggles).
func (s *IoTService) updateEventConfigurationsCore(store iotstore.IotStoreInterface, incoming map[string]interface{}) error {
	rec := map[string]interface{}{}
	if _, err := store.GetGenericExists("config/eventConfigurations", &rec); err != nil {
		return err
	}
	for k, v := range incoming {
		rec[k] = v
	}
	return store.PutGeneric("config/eventConfigurations", rec)
}

// describeEventConfigurationsCore loads the event configuration map (empty
// when never configured).
func (s *IoTService) describeEventConfigurationsCore(store iotstore.IotStoreInterface) (map[string]interface{}, error) {
	rec := map[string]interface{}{}
	if _, err := store.GetGenericExists("config/eventConfigurations", &rec); err != nil {
		return nil, err
	}
	return rec, nil
}

// UpdateEncryptionConfigurationInput carries the parsed
// UpdateEncryptionConfiguration request.
type UpdateEncryptionConfigurationInput struct {
	KmsKeyArn        string
	KmsAccessRoleArn string
	EncryptionType   string
}

// encryptionTypes is the EncryptionType enum member set.
var encryptionTypes = map[string]bool{
	"CUSTOMER_MANAGED_KMS_KEY": true, "AWS_OWNED_KMS_KEY": true,
}

// updateEncryptionConfigurationCore persists the account encryption
// configuration. The encryptionType member is required and must be an
// EncryptionType enum member.
func (s *IoTService) updateEncryptionConfigurationCore(store iotstore.IotStoreInterface, in UpdateEncryptionConfigurationInput) error {
	if !encryptionTypes[in.EncryptionType] {
		return iotstore.ErrInvalidRequest
	}
	rec := map[string]interface{}{
		"kmsKeyArn":        in.KmsKeyArn,
		"kmsAccessRoleArn": in.KmsAccessRoleArn,
	}
	rec["encryptionType"] = in.EncryptionType
	return store.PutGeneric("config/encryptionConfiguration", rec)
}

// describeEncryptionConfigurationCore loads the account encryption
// configuration, returning the documented default when never configured.
func (s *IoTService) describeEncryptionConfigurationCore(store iotstore.IotStoreInterface) (map[string]interface{}, error) {
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists("config/encryptionConfiguration", &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		// AWS IoT accounts default to an AWS owned KMS key for data at
		// rest.
		return map[string]interface{}{
			"encryptionType": "AWS_OWNED_KMS_KEY",
		}, nil
	}
	return rec, nil
}
