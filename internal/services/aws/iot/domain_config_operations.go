package iot

import (
	"context"
	"encoding/json"
	"time"

	"vorpalstacks/internal/common/request"
	iotstore "vorpalstacks/internal/store/aws/iot"
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

func (s *IoTService) CreateDomainConfiguration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "domainConfigurationName")
	if name == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if _, err := store.GetDomainConfiguration(name); err == nil {
		return nil, iotstore.ErrDomainConfigurationAlreadyExists
	}

	now := time.Now().UTC()
	dc := &iotstore.DomainConfiguration{
		DomainConfigurationName:   name,
		DomainName:                request.GetParamCaseInsensitive(req.Parameters, "domainName"),
		ServerCertificateARNs:     parseStringListParam(req.Parameters, "serverCertificateArns"),
		AuthorizerConfig:          request.GetParamCaseInsensitive(req.Parameters, "authorizerConfig"),
		ServiceType:               request.GetParamCaseInsensitive(req.Parameters, "serviceType"),
		DomainConfigurationStatus: "ENABLED",
		CreationDate:              now,
		LastModifiedDate:          now,
	}

	created, err := store.CreateDomainConfiguration(dc)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"domainConfigurationName": created.DomainConfigurationName,
		"domainConfigurationArn":  created.DomainConfigurationARN,
	}, nil
}

func (s *IoTService) DescribeDomainConfiguration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "domainConfigurationName")
	if name == "" {
		name = "default"
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	dc, err := store.GetDomainConfiguration(name)
	if err != nil {
		return nil, iotstore.ErrDomainConfigurationNotFound
	}

	return domainConfigDetailResponse(dc), nil
}

func (s *IoTService) UpdateDomainConfiguration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "domainConfigurationName")
	if name == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	dc, err := store.GetDomainConfiguration(name)
	if err != nil {
		return nil, iotstore.ErrDomainConfigurationNotFound
	}

	if domainName := request.GetParamCaseInsensitive(req.Parameters, "domainName"); domainName != "" {
		dc.DomainName = domainName
	}
	if authCfg := request.GetParamCaseInsensitive(req.Parameters, "authorizerConfig"); authCfg != "" {
		dc.AuthorizerConfig = authCfg
	}
	if svcType := request.GetParamCaseInsensitive(req.Parameters, "serviceType"); svcType != "" {
		dc.ServiceType = svcType
	}
	if certARNs := parseStringListParam(req.Parameters, "serverCertificateArns"); certARNs != nil {
		dc.ServerCertificateARNs = certARNs
	}
	dc.LastModifiedDate = time.Now().UTC()

	if err := store.UpdateDomainConfiguration(name, dc); err != nil {
		return nil, err
	}

	return map[string]interface{}{}, nil
}

func (s *IoTService) DeleteDomainConfiguration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "domainConfigurationName")
	if name == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.DeleteDomainConfiguration(name); err != nil {
		return nil, err
	}

	return map[string]interface{}{}, nil
}

func (s *IoTService) ListDomainConfigurations(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	dcs, err := store.ListDomainConfigurations(parseListOptions(req.Parameters))
	if err != nil {
		return nil, err
	}

	items := make([]map[string]interface{}, 0, len(dcs.Items))
	for _, dc := range dcs.Items {
		items = append(items, domainConfigResponse(dc))
	}

	return listResponse("domainConfigurations", items, dcs.NextMarker), nil
}
