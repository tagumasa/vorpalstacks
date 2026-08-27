package iot

import (
	"context"
	"encoding/json"

	"vorpalstacks/internal/common/request"
)

func parseStringListParam(params map[string]interface{}, key string) []string {
	v, ok := params[key]
	if !ok {
		return nil
	}
	switch val := v.(type) {
	case string:
		if val == "" {
			return nil
		}
		var list []string
		if json.Unmarshal([]byte(val), &list) == nil {
			return list
		}
		return nil
	case []interface{}:
		result := make([]string, 0, len(val))
		for _, item := range val {
			if s, ok := item.(string); ok {
				result = append(result, s)
			}
		}
		return result
	}
	return nil
}

// parseDomainConfigAuthorizerParam extracts the wire authorizer-config
// structure; ok is false when the member is absent.
func parseDomainConfigAuthorizerParam(params map[string]interface{}) (DomainConfigAuthorizerConfig, bool) {
	m := request.GetMapParamCaseInsensitive(params, "authorizerConfig")
	if m == nil {
		return DomainConfigAuthorizerConfig{}, false
	}
	cfg := DomainConfigAuthorizerConfig{
		DefaultAuthorizerName: request.GetParamCaseInsensitive(m, "defaultAuthorizerName"),
	}
	if request.HasParam(m, "allowAuthorizerOverride") {
		cfg.AllowAuthorizerOverride = request.GetBoolParam(m, "allowAuthorizerOverride")
		cfg.OverrideProvided = true
	}
	return cfg, true
}

func (s *IoTService) CreateDomainConfiguration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	authorizerConfig, authorizerProvided := parseDomainConfigAuthorizerParam(req.Parameters)

	result, err := s.createDomainConfigurationCore(store, CreateDomainConfigurationInput{
		DomainConfigurationName:  request.GetParamCaseInsensitive(req.Parameters, "domainConfigurationName"),
		DomainName:               request.GetParamCaseInsensitive(req.Parameters, "domainName"),
		ServerCertificateARNs:    parseStringListParam(req.Parameters, "serverCertificateArns"),
		ValidationCertificateARN: request.GetParamCaseInsensitive(req.Parameters, "validationCertificateArn"),
		AuthorizerConfig:         authorizerConfig,
		AuthorizerConfigProvided: authorizerProvided,
		ServiceType:              request.GetParamCaseInsensitive(req.Parameters, "serviceType"),
		AuthenticationType:       request.GetParamCaseInsensitive(req.Parameters, "authenticationType"),
		ApplicationProtocol:      request.GetParamCaseInsensitive(req.Parameters, "applicationProtocol"),
		Tags:                     tagListParam(req.Parameters),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"domainConfigurationName": result.DomainConfigurationName,
		"domainConfigurationArn":  result.DomainConfigurationARN,
	}, nil
}

func (s *IoTService) DescribeDomainConfiguration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	dc, err := s.describeDomainConfigurationCore(store, DescribeDomainConfigurationInput{
		DomainConfigurationName: request.GetParamCaseInsensitive(req.Parameters, "domainConfigurationName"),
	})
	if err != nil {
		return nil, err
	}

	return domainConfigDetailResponse(dc), nil
}

func (s *IoTService) UpdateDomainConfiguration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	authorizerConfig, authorizerProvided := parseDomainConfigAuthorizerParam(req.Parameters)

	result, err := s.updateDomainConfigurationCore(store, UpdateDomainConfigurationInput{
		DomainConfigurationName:   request.GetParamCaseInsensitive(req.Parameters, "domainConfigurationName"),
		AuthorizerConfig:          authorizerConfig,
		AuthorizerConfigProvided:  authorizerProvided,
		DomainConfigurationStatus: request.GetParamCaseInsensitive(req.Parameters, "domainConfigurationStatus"),
		RemoveAuthorizerConfig:    request.HasParam(req.Parameters, "removeAuthorizerConfig") && request.GetBoolParam(req.Parameters, "removeAuthorizerConfig"),
		AuthenticationType:        request.GetParamCaseInsensitive(req.Parameters, "authenticationType"),
		ApplicationProtocol:       request.GetParamCaseInsensitive(req.Parameters, "applicationProtocol"),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"domainConfigurationName": result.DomainConfigurationName,
		"domainConfigurationArn":  result.DomainConfigurationARN,
	}, nil
}

func (s *IoTService) DeleteDomainConfiguration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.deleteDomainConfigurationCore(store, request.GetParamCaseInsensitive(req.Parameters, "domainConfigurationName")); err != nil {
		return nil, err
	}

	return map[string]interface{}{}, nil
}

func (s *IoTService) ListDomainConfigurations(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	opts := parseListOptions(req.Parameters)
	result, err := s.listDomainConfigurationsCore(store, opts.Marker, opts.MaxItems)
	if err != nil {
		return nil, err
	}

	items := make([]map[string]interface{}, 0, len(result.DomainConfigurations))
	for _, dc := range result.DomainConfigurations {
		items = append(items, domainConfigResponse(dc))
	}

	return listResponse("domainConfigurations", items, result.NextMarker), nil
}
