package iot

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// ---- Logging / Event / Encryption config --------------------------
// Persisted via GenericKV under "config/<name>". A missing key means "not yet
// configured"; the Cores return a default/empty shape for that case and
// propagate genuine store errors.

func (s *IoTService) GetV2LoggingOptions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rec, err := s.getV2LoggingOptionsCore(store)
	if err != nil {
		return nil, err
	}
	// GetV2LoggingOptions output shape is flat (roleArn, defaultLogLevel,
	// disableAllLogs at the top level). Wrapping in "loggingOptions" causes
	// the AWS SDK parser to discard all fields.
	return rec, nil
}

// rawListParam extracts a raw wire list member from the request parameters,
// returning nil when the member is absent or not a list.
func rawListParam(params map[string]interface{}, key string) []interface{} {
	if raw, ok := params[key].([]interface{}); ok {
		return raw
	}
	return nil
}

func (s *IoTService) SetV2LoggingOptions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	in := SetV2LoggingOptionsInput{
		RoleArn:             request.GetParamCaseInsensitive(req.Parameters, "roleArn"),
		DefaultLogLevel:     request.GetParamCaseInsensitive(req.Parameters, "defaultLogLevel"),
		DisableAllLogs:      request.GetBoolParam(req.Parameters, "disableAllLogs"),
		EventConfigurations: rawListParam(req.Parameters, "eventConfigurations"),
	}
	if err := s.setV2LoggingOptionsCore(store, in); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}
func (s *IoTService) DeleteV2LoggingLevel(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	targetType := request.GetParamCaseInsensitive(req.Parameters, "targetType")
	targetName := request.GetParamCaseInsensitive(req.Parameters, "targetName")
	if err := s.deleteV2LoggingLevelCore(store, targetType, targetName); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}
func (s *IoTService) ListV2LoggingLevels(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	configs, err := s.listV2LoggingLevelsCore(store, request.GetParamCaseInsensitive(req.Parameters, "targetType"))
	if err != nil {
		return nil, err
	}
	return paginatedMaps("logTargetConfigurations", configs, req.Parameters)
}
func (s *IoTService) SetV2LoggingLevel(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	logTarget := request.GetMapParamCaseInsensitive(req.Parameters, "logTarget")
	logLevel := request.GetParamCaseInsensitive(req.Parameters, "logLevel")
	if err := s.setV2LoggingLevelCore(store, logTarget, logLevel); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}
func (s *IoTService) GetLoggingOptions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return s.getLoggingOptionsCore(store)
}
func (s *IoTService) SetLoggingOptions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	// AWS uses httpPayload on LoggingOptionsPayload so the body is wrapped.
	props := request.GetMapParamCaseInsensitive(req.Parameters, "loggingOptionsPayload")
	if props == nil {
		// Some SDKs send the members flat (no wrapper); accept either form.
		props = req.Parameters
	}
	roleArn := request.GetParamCaseInsensitive(props, "roleArn")
	logLevel := request.GetParamCaseInsensitive(props, "logLevel")
	if err := s.setLoggingOptionsCore(store, roleArn, logLevel); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}
func (s *IoTService) DescribeEventConfigurations(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rec, err := s.describeEventConfigurationsCore(store)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{"eventConfigurations": rec}, nil
}
func (s *IoTService) UpdateEventConfigurations(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	incoming, _ := req.Parameters["eventConfigurations"].(map[string]interface{})
	if err := s.updateEventConfigurationsCore(store, incoming); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}
func (s *IoTService) DescribeEncryptionConfiguration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rec, err := s.describeEncryptionConfigurationCore(store)
	if err != nil {
		return nil, err
	}
	// DescribeEncryptionConfiguration output shape is flat (encryptionType,
	// kmsKeyArn at the top level). Wrapping in "encryptionConfiguration"
	// causes the AWS SDK parser to discard all fields.
	return rec, nil
}
func (s *IoTService) UpdateEncryptionConfiguration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	in := UpdateEncryptionConfigurationInput{
		KmsKeyArn:        request.GetParamCaseInsensitive(req.Parameters, "kmsKeyArn"),
		KmsAccessRoleArn: request.GetParamCaseInsensitive(req.Parameters, "kmsAccessRoleArn"),
		EncryptionType:   request.GetParamCaseInsensitive(req.Parameters, "encryptionType"),
	}
	if err := s.updateEncryptionConfigurationCore(store, in); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}
