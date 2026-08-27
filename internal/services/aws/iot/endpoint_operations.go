package iot

import (
	"context"
	"fmt"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/serviceports"
	"vorpalstacks/internal/config"
)

// DescribeEndpoint returns a connectable local endpoint address for the
// specified IoT endpoint type. MQTT types map to the request region's
// broker port (each region owns an independent broker with a dynamic
// port); iot:CredentialProvider maps to the configured HTTP port
// (config.ServerPort, honouring PORT env overrides).
func (s *IoTService) DescribeEndpoint(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	endpointType := request.GetParamCaseInsensitive(req.Parameters, "endpointType")
	if endpointType == "" {
		endpointType = "iot:Data-ATS"
	}

	var port int
	switch endpointType {
	case "iot:CredentialProvider":
		port = config.ServerPort()
	default:
		// MQTT endpoint types (iot:Data, iot:Data-ATS, iot:Data-ALPN,
		// iot:Jobs and unknown free-form values per the Smithy model)
		// resolve to the request region's broker port so each region
		// returns its own endpoint.
		port = s.brokerPortForRequest(reqCtx)
	}

	return map[string]interface{}{
		"endpointAddress": fmt.Sprintf("localhost:%d", port),
	}, nil
}

// brokerPortForRequest resolves the MQTT broker port for the request's
// region. Falls back to the first available broker (for callers without a
// request context) and finally to the legacy fixed port constant when no
// brokers are configured (e.g. unit tests that never call Init).
func (s *IoTService) brokerPortForRequest(reqCtx *request.RequestContext) int {
	region := ""
	if reqCtx != nil {
		region = reqCtx.GetRegion()
	}
	if brk := s.brokers[region]; brk != nil {
		return brk.Port()
	}
	for _, brk := range s.brokers {
		return brk.Port()
	}
	return serviceports.IotMQTT
}

// GetIndexingConfiguration retrieves the fleet indexing configuration from the store.
func (s *IoTService) GetIndexingConfiguration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	ic := s.getIndexingConfigurationCore(store)
	if ic == nil {
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

	in := UpdateIndexingConfigurationInput{}
	if tic, ok := req.Parameters["thingIndexingConfiguration"]; ok {
		if ticMap, ok := tic.(map[string]interface{}); ok {
			if v, ok := ticMap["thingIndexingMode"].(string); ok {
				in.ThingIndexingMode = &v
			}
			if v, ok := ticMap["thingGroupIndexingMode"].(string); ok {
				in.ThingGroupIndexingMode = &v
			}
			if v, ok := ticMap["thingConnectivityIndexingMode"].(string); ok {
				in.ThingConnectivityIndexingMode = &v
			}
		}
	}

	if err := s.updateIndexingConfigurationCore(store, in); err != nil {
		return nil, err
	}

	return map[string]interface{}{}, nil
}
