package iot

import (
	"context"

	"vorpalstacks/internal/common/request"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

// DescribeEndpoint returns the endpoint address for the specified IoT
// endpoint type.
func (s *IoTService) DescribeEndpoint(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	endpointType := request.GetParamCaseInsensitive(req.Parameters, "endpointType")
	if endpointType == "" {
		endpointType = "iot:Data-ATS"
	}

	host := reqCtx.GetRegion() + ".iot." + s.accountID + ".amazonaws.com"

	return map[string]interface{}{
		"endpointAddress": host,
	}, nil
}

// GetIndexingConfiguration retrieves the fleet indexing configuration.
func (s *IoTService) GetIndexingConfiguration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return map[string]interface{}{
		"thingIndexingConfiguration": map[string]interface{}{
			"thingIndexingMode": "OFF",
		},
	}, nil
}

// UpdateIndexingConfiguration updates the fleet indexing configuration.
func (s *IoTService) UpdateIndexingConfiguration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return map[string]interface{}{}, nil
}

// DescribeDomainConfiguration retrieves a domain configuration.
func (s *IoTService) DescribeDomainConfiguration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "domainConfigurationName")
	if name == "" {
		name = "default"
	}

	return map[string]interface{}{
		"domainConfigurationName": name,
		"domainConfigurationArn":  "arn:aws:iot:" + reqCtx.GetRegion() + ":" + s.accountID + ":domainconfiguration/" + name,
	}, nil
}

// CreateDomainConfiguration creates a domain configuration.
func (s *IoTService) CreateDomainConfiguration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "domainConfigurationName")
	if name == "" {
		return nil, iotstore.ErrMissingParam
	}

	return map[string]interface{}{
		"domainConfigurationName": name,
		"domainConfigurationArn":  "arn:aws:iot:" + reqCtx.GetRegion() + ":" + s.accountID + ":domainconfiguration/" + name,
	}, nil
}

// UpdateDomainConfiguration modifies a domain configuration.
func (s *IoTService) UpdateDomainConfiguration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return map[string]interface{}{}, nil
}

// DeleteDomainConfiguration removes a domain configuration.
func (s *IoTService) DeleteDomainConfiguration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "domainConfigurationName")
	if name == "" {
		return nil, iotstore.ErrMissingParam
	}

	return map[string]interface{}{}, nil
}

// ListDomainConfigurations returns all domain configurations.
func (s *IoTService) ListDomainConfigurations(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return map[string]interface{}{
		"domainConfigurations": []map[string]interface{}{},
	}, nil
}
