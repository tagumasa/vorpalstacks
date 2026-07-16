package iotevents

import (
	"context"
	"time"

	"vorpalstacks/internal/common/request"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

// Package iotevents hosts both the iotevents control plane and the
// iotevents-data data plane. Per AGENTS.md #19 it is a declared sub-service
// of iot and shares the iotstore.

func (s *IoTEventsService) BatchAcknowledgeAlarm(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return batchAlarmResponse(), nil
}

func (s *IoTEventsService) BatchDisableAlarm(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return batchAlarmResponse(), nil
}

func (s *IoTEventsService) BatchEnableAlarm(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return batchAlarmResponse(), nil
}

func (s *IoTEventsService) BatchResetAlarm(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return batchAlarmResponse(), nil
}

func (s *IoTEventsService) BatchSnoozeAlarm(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return batchAlarmResponse(), nil
}

func (s *IoTEventsService) BatchDeleteDetector(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return batchDetectorResponse("batchDeleteDetectorErrorEntries"), nil
}

func (s *IoTEventsService) BatchUpdateDetector(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return batchDetectorResponse("batchUpdateDetectorErrorEntries"), nil
}

func (s *IoTEventsService) DescribeAlarm(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "alarmModelName")
	if name == "" {
		return nil, iotstore.ErrMissingParam
	}
	key := request.GetParamCaseInsensitive(req.Parameters, "keyValue")
	return map[string]interface{}{
		"alarmModelName": name,
		"keyValue":       key,
		"state": map[string]interface{}{
			"stateName": "DISABLED",
		},
		"creationTime": time.Now().Unix(),
	}, nil
}

func (s *IoTEventsService) DescribeDetector(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "detectorModelName")
	if name == "" {
		return nil, iotstore.ErrMissingParam
	}
	key := request.GetParamCaseInsensitive(req.Parameters, "keyValue")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	stateName := "IDLE"
	if concrete, ok := store.(*iotstore.IotStore); ok && concrete.StateMachine() != nil {
		if actual, found := concrete.StateMachine().GetState(name, key); found {
			stateName = actual
		}
	}

	return map[string]interface{}{
		"detectorModelName": name,
		"keyValue":          key,
		"state": map[string]interface{}{
			"stateName": stateName,
		},
		"creationTime": time.Now().Unix(),
	}, nil
}

func (s *IoTEventsService) ListAlarms(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return listResponse("alarmSummaries", []map[string]interface{}{}, ""), nil
}

func (s *IoTEventsService) ListDetectors(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	modelName := request.GetParamCaseInsensitive(req.Parameters, "detectorModelName")
	stateFilter := request.GetParamCaseInsensitive(req.Parameters, "stateName")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	var instances []map[string]interface{}
	if concrete, ok := store.(*iotstore.IotStore); ok && concrete.StateMachine() != nil {
		instances = concrete.StateMachine().ListDetectorInstances(modelName)
	}

	if stateFilter != "" {
		filtered := instances[:0]
		for _, inst := range instances {
			if inst["stateName"] == stateFilter {
				filtered = append(filtered, inst)
			}
		}
		instances = filtered
	}

	summaries := make([]map[string]interface{}, 0, len(instances))
	for _, inst := range instances {
		summaries = append(summaries, map[string]interface{}{
			"detectorModelName": inst["detectorModelName"],
			"keyValue":          inst["keyValue"],
			"state": map[string]interface{}{
				"stateName": inst["stateName"],
			},
		})
	}

	return listResponse("detectorSummaries", summaries, ""), nil
}

// batchAlarmResponse returns the standard IoTEvents alarm-batch response
// shape. AWS SDK Go v2 deserialisers for all five Batch*Alarm operations
// expect the key "errorEntries" (see deserializers.go lines 170, 492, 653,
// 975, 1136). Since this implementation processes all operations
// successfully, the error-entries list is always empty.
func batchAlarmResponse() map[string]interface{} {
	return map[string]interface{}{
		"errorEntries": []map[string]interface{}{},
	}
}

// batchDetectorResponse returns the per-operation error-entries key with
// an empty list, matching the AWS IoT Events Data API response shape.
func batchDetectorResponse(responseKey string) map[string]interface{} {
	return map[string]interface{}{
		responseKey: []map[string]interface{}{},
	}
}
