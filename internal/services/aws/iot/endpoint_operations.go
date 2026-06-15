package iot

import (
	"context"
	"fmt"

	"vorpalstacks/internal/common/request"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

// DescribeEndpoint returns the endpoint address for the specified IoT
// endpoint type. Endpoint format follows AWS convention:
// {prefix}.iot.{region}.amazonaws.com (varies by type).
func (s *IoTService) DescribeEndpoint(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	endpointType := request.GetParamCaseInsensitive(req.Parameters, "endpointType")
	if endpointType == "" {
		endpointType = "iot:Data-ATS"
	}

	region := reqCtx.GetRegion()
	var host string
	switch endpointType {
	case "iot:Data-ATS", "iot:Data":
		host = fmt.Sprintf("%s-ats.iot.%s.amazonaws.com", s.accountID, region)
	case "iot:Data-ALPN":
		host = fmt.Sprintf("data-alpn.iot.%s.amazonaws.com", region)
	case "iot:CredentialProvider":
		host = fmt.Sprintf("%s.cred.iot.%s.amazonaws.com", s.accountID, region)
	case "iot:Jobs":
		host = fmt.Sprintf("%s.jobs.iot.%s.amazonaws.com", s.accountID, region)
	default:
		host = fmt.Sprintf("%s.iot.%s.amazonaws.com", s.accountID, region)
	}

	return map[string]interface{}{
		"endpointAddress": host,
	}, nil
}

// GetIndexingConfiguration retrieves the fleet indexing configuration from the store.
func (s *IoTService) GetIndexingConfiguration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	ic, err := store.GetIndexingConfiguration()
	if err != nil {
		return map[string]interface{}{
			"thingIndexingConfiguration": map[string]interface{}{
				"thingIndexingMode": "OFF",
			},
		}, nil
	}

	return map[string]interface{}{
		"thingIndexingConfiguration": map[string]interface{}{
			"thingIndexingMode":             ic.ThingIndexingMode,
			"thingGroupIndexingMode":        ic.ThingGroupIndexingMode,
			"thingConnectivityIndexingMode": ic.ThingConnectivityIndexingMode,
		},
	}, nil
}

// UpdateIndexingConfiguration persists the fleet indexing configuration to the store.
func (s *IoTService) UpdateIndexingConfiguration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	ic, err := store.GetIndexingConfiguration()
	if err != nil {
		ic = &iotstore.IndexingConfiguration{}
	}

	if tic, ok := req.Parameters["thingIndexingConfiguration"]; ok {
		if ticMap, ok := tic.(map[string]interface{}); ok {
			if v, ok := ticMap["thingIndexingMode"].(string); ok {
				ic.ThingIndexingMode = v
			}
			if v, ok := ticMap["thingGroupIndexingMode"].(string); ok {
				ic.ThingGroupIndexingMode = v
			}
			if v, ok := ticMap["thingConnectivityIndexingMode"].(string); ok {
				ic.ThingConnectivityIndexingMode = v
			}
		}
	}

	if err := store.UpdateIndexingConfiguration(ic); err != nil {
		return nil, err
	}

	return map[string]interface{}{}, nil
}
