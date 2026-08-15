package iot

import (
	"context"

	"vorpalstacks/internal/common/request"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

// ---- Logging / Event / Encryption config --------------------------
// Persisted via GenericKV under "config/<name>". A missing key means "not yet
// configured"; the handlers return a default/empty shape for that case and
// propagate genuine store errors.

func (s *IoTService) GetV2LoggingOptions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
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
	// GetV2LoggingOptions output shape is flat (roleArn, defaultLogLevel,
	// disableAllLogs at the top level). Wrapping in "loggingOptions" causes
	// the AWS SDK parser to discard all fields.
	return rec, nil
}
func (s *IoTService) SetV2LoggingOptions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	defaultLogLevel := request.GetParamCaseInsensitive(req.Parameters, "defaultLogLevel")
	if defaultLogLevel == "" {
		defaultLogLevel = "INFO"
	}
	rec := map[string]interface{}{
		"roleArn":         request.GetParamCaseInsensitive(req.Parameters, "roleArn"),
		"defaultLogLevel": defaultLogLevel,
		"disableAllLogs":  request.GetBoolParam(req.Parameters, "disableAllLogs"),
	}
	if err := store.PutGeneric("config/v2Logging", rec); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}
func (s *IoTService) DeleteV2LoggingLevel(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	targetType := request.GetParamCaseInsensitive(req.Parameters, "targetType")
	targetName := request.GetParamCaseInsensitive(req.Parameters, "targetName")
	if targetType == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.DeleteGeneric("v2LoggingLevel/" + targetType + "/" + targetName); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}
func (s *IoTService) ListV2LoggingLevels(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	prefix := "v2LoggingLevel/"
	if tt := request.GetParamCaseInsensitive(req.Parameters, "targetType"); tt != "" {
		prefix = "v2LoggingLevel/" + tt + "/"
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
	return paginatedMaps("logTargetConfigurations", configs, req.Parameters), nil
}
func (s *IoTService) SetV2LoggingLevel(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	logTarget := request.GetMapParamCaseInsensitive(req.Parameters, "logTarget")
	logLevel := request.GetParamCaseInsensitive(req.Parameters, "logLevel")
	targetType, _ := logTarget["targetType"].(string)
	targetName, _ := logTarget["targetName"].(string)
	if targetType == "" || logLevel == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.PutGeneric("v2LoggingLevel/"+targetType+"/"+targetName, map[string]interface{}{
		"logTarget": logTarget,
		"logLevel":  logLevel,
	}); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}
func (s *IoTService) GetLoggingOptions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists("config/legacyLogging", &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		// AWS returns an empty response when logging has never been configured.
		return map[string]interface{}{}, nil
	}
	return map[string]interface{}{
		"roleArn":  rec["roleArn"],
		"logLevel": rec["logLevel"],
	}, nil
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
	if roleArn == "" {
		return nil, iotstore.ErrMissingParam
	}
	rec := map[string]interface{}{
		"roleArn":  roleArn,
		"logLevel": request.GetParamCaseInsensitive(props, "logLevel"),
	}
	if err := store.PutGeneric("config/legacyLogging", rec); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}
func (s *IoTService) DescribeEventConfigurations(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rec := map[string]interface{}{}
	if _, err := store.GetGenericExists("config/eventConfigurations", &rec); err != nil {
		return nil, err
	}
	return map[string]interface{}{"eventConfigurations": rec}, nil
}
func (s *IoTService) UpdateEventConfigurations(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rec := map[string]interface{}{}
	if _, err := store.GetGenericExists("config/eventConfigurations", &rec); err != nil {
		return nil, err
	}
	// Merge incoming configuration attributes into the persisted map so that
	// partial updates behave like AWS IoT (per-event-type toggles).
	if incoming, ok := req.Parameters["eventConfigurations"].(map[string]interface{}); ok {
		for k, v := range incoming {
			rec[k] = v
		}
	}
	if err := store.PutGeneric("config/eventConfigurations", rec); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}
func (s *IoTService) DescribeEncryptionConfiguration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists("config/encryptionConfiguration", &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return map[string]interface{}{
			"encryptionType": "TLS",
		}, nil
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
	rec := map[string]interface{}{
		"kmsKeyArn":        request.GetParamCaseInsensitive(req.Parameters, "kmsKeyArn"),
		"kmsAccessRoleArn": request.GetParamCaseInsensitive(req.Parameters, "kmsAccessRoleArn"),
	}
	et := request.GetParamCaseInsensitive(req.Parameters, "encryptionType")
	if et == "" {
		et = "TLS"
	}
	rec["encryptionType"] = et
	if err := store.PutGeneric("config/encryptionConfiguration", rec); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}
