package iotevents

import (
	"context"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

// Package iotevents hosts both the iotevents control plane and the
// iotevents-data data plane. Per AGENTS.md #19 it is a declared sub-service
// of iot and shares the iotstore.

// BatchAcknowledgeAlarm acknowledges one or more alarms. Alarms in ACTIVE or
// LATCHED state transition to ACKNOWLEDGED. Non-existent or non-active
// alarms are silently accepted (AWS returns no error for those cases).
func (s *IoTEventsService) BatchAcknowledgeAlarm(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	asm := s.alarmStateMachine(store)
	acknowledgeActions := getAcknowledgeActionRequests(req.Parameters)
	for _, a := range acknowledgeActions {
		if asm != nil {
			asm.AcknowledgeAlarm(a.alarmModelName, a.keyValue)
		}
	}
	return batchAlarmResponse(), nil
}

// BatchDisableAlarm transitions the specified alarms to DISABLED state.
// Non-existent alarms are created in DISABLED state (matching AWS behaviour
// where BatchDisableAlarm on a non-existent alarm succeeds silently).
func (s *IoTEventsService) BatchDisableAlarm(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	asm := s.alarmStateMachine(store)
	for _, a := range getAlarmActionRequests(req.Parameters, "disableActionRequests") {
		if asm != nil {
			asm.DisableAlarm(a.alarmModelName, a.keyValue)
		}
	}
	return batchAlarmResponse(), nil
}

// BatchEnableAlarm creates or enables the specified alarms. Newly created
// alarms start in NORMAL state with the current time as creationTime.
func (s *IoTEventsService) BatchEnableAlarm(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	asm := s.alarmStateMachine(store)
	for _, a := range getAlarmActionRequests(req.Parameters, "enableActionRequests") {
		if asm != nil {
			asm.EnableAlarm(a.alarmModelName, a.keyValue)
		}
	}
	return batchAlarmResponse(), nil
}

// BatchResetAlarm transitions the specified alarms to NORMAL state.
func (s *IoTEventsService) BatchResetAlarm(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	asm := s.alarmStateMachine(store)
	for _, a := range getAlarmActionRequests(req.Parameters, "resetActionRequests") {
		if asm != nil {
			asm.ResetAlarm(a.alarmModelName, a.keyValue)
		}
	}
	return batchAlarmResponse(), nil
}

// BatchSnoozeAlarm transitions the specified alarms to SNOOZE state.
func (s *IoTEventsService) BatchSnoozeAlarm(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	asm := s.alarmStateMachine(store)
	for _, a := range getAlarmActionRequests(req.Parameters, "snoozeActionRequests") {
		if asm != nil {
			asm.SnoozeAlarm(a.alarmModelName, a.keyValue)
		}
	}
	return batchAlarmResponse(), nil
}

// BatchDeleteDetector deletes detector instances identified by
// (detectorModelName, keyValue) pairs. Non-existent instances are silently
// accepted (AWS returns no error for those cases). Returns an error entry
// for each request whose detectorModelName or keyValue is missing (H-SM4).
func (s *IoTEventsService) BatchDeleteDetector(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	concrete, ok := store.(*iotstore.IotStore)
	if !ok || concrete.StateMachine() == nil {
		return batchDetectorResponse("batchDeleteDetectorErrorEntries"), nil
	}

	detectors := getDetectorRequests(req.Parameters, "detectors")
	var errorEntries []map[string]interface{}
	for _, d := range detectors {
		if d.modelName == "" || d.keyValue == "" {
			errorEntries = append(errorEntries, map[string]interface{}{
				"errorCode":    "InvalidRequestException",
				"errorMessage": "detectorModelName and keyValue are required",
			})
			continue
		}
		concrete.DeleteDetector(d.modelName, d.keyValue)
	}

	if errorEntries == nil {
		errorEntries = []map[string]interface{}{}
	}
	return map[string]interface{}{
		"batchDeleteDetectorErrorEntries": errorEntries,
	}, nil
}

// BatchUpdateDetector updates the state of detector instances identified by
// (detectorModelName, keyValue) pairs. Returns an error entry for each
// request that is missing required fields (H-SM4).
func (s *IoTEventsService) BatchUpdateDetector(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	concrete, ok := store.(*iotstore.IotStore)
	if !ok || concrete.StateMachine() == nil {
		return batchDetectorResponse("batchUpdateDetectorErrorEntries"), nil
	}

	updates := getDetectorUpdateRequests(req.Parameters)
	var errorEntries []map[string]interface{}
	for _, u := range updates {
		if u.modelName == "" || u.keyValue == "" || u.stateName == "" {
			errorEntries = append(errorEntries, map[string]interface{}{
				"errorCode":    "InvalidRequestException",
				"errorMessage": "detectorModelName, keyValue, and state.stateName are required",
			})
			continue
		}
		concrete.UpdateDetectorState(u.modelName, u.keyValue, u.stateName)
	}

	if errorEntries == nil {
		errorEntries = []map[string]interface{}{}
	}
	return map[string]interface{}{
		"batchUpdateDetectorErrorEntries": errorEntries,
	}, nil
}

// DescribeAlarm returns the current state of the specified alarm instance.
// Returns ResourceNotFoundException if the alarm has not been created.
func (s *IoTEventsService) DescribeAlarm(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "alarmModelName")
	if name == "" {
		return nil, iotstore.ErrMissingParam
	}
	key := request.GetParamCaseInsensitive(req.Parameters, "keyValue")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	asm := s.alarmStateMachine(store)
	if asm == nil {
		return nil, awserrors.NewAWSError("InternalFailureException", "alarm state machine not initialised", 500)
	}

	stateName, creationTime, lastUpdate, ok := asm.DescribeAlarm(name, key)
	if !ok {
		return nil, iotstore.ErrAlarmModelNotFound
	}

	return map[string]interface{}{
		"alarmModelName": name,
		"keyValue":       key,
		"state": map[string]interface{}{
			"stateName": stateName,
		},
		"creationTime": creationTime.Unix(),
		"lastUpdate":   lastUpdate.Unix(),
	}, nil
}

// DescribeDetector returns the current state of the specified detector instance.
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
	var creationTime time.Time
	if concrete, ok := store.(*iotstore.IotStore); ok && concrete.StateMachine() != nil {
		if state, ct, _, found := concrete.StateMachine().GetDetectorDetail(name, key); found {
			stateName = state
			creationTime = ct
		}
	}

	resp := map[string]interface{}{
		"detectorModelName": name,
		"keyValue":          key,
		"state": map[string]interface{}{
			"stateName": stateName,
		},
	}
	// Use the instance's actual creation time; fall back to time.Now() only
	// when the detector has never been instantiated (AWS returns the current
	// time for a freshly queried detector that has not yet received input).
	if !creationTime.IsZero() {
		resp["creationTime"] = creationTime.Unix()
	} else {
		resp["creationTime"] = time.Now().Unix()
	}
	return resp, nil
}

// ListAlarms returns all alarm instances for the specified alarm model.
func (s *IoTEventsService) ListAlarms(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	modelName := request.GetParamCaseInsensitive(req.Parameters, "alarmModelName")
	if modelName == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	asm := s.alarmStateMachine(store)
	var summaries []map[string]interface{}
	if asm != nil {
		instances := asm.ListAlarms(modelName)
		summaries = make([]map[string]interface{}, 0, len(instances))
		for _, inst := range instances {
			summaries = append(summaries, map[string]interface{}{
				"alarmModelName": modelName,
				"keyValue":       inst.KeyValue,
				"state": map[string]interface{}{
					"stateName": inst.StateName,
				},
				"creationTime": inst.CreationTime.Unix(),
				"lastUpdate":   inst.LastUpdate.Unix(),
			})
		}
	} else {
		summaries = []map[string]interface{}{}
	}

	return listResponse("alarmSummaries", summaries, ""), nil
}

// ListDetectors returns detector instances for the specified detector model.
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

// alarmStateMachine extracts the AlarmStateMachine from the store, if available.
func (s *IoTEventsService) alarmStateMachine(store interface{}) *iotstore.AlarmStateMachine {
	concrete, ok := store.(*iotstore.IotStore)
	if !ok || concrete == nil {
		return nil
	}
	return concrete.AlarmStateMachine()
}

// alarmActionRequest holds the parsed fields from a batch alarm action request entry.
type alarmActionRequest struct {
	alarmModelName string
	keyValue       string
}

// getAlarmActionRequests extracts alarm action requests from the parameters
// for the given parameter key (e.g. "enableActionRequests", "disableActionRequests").
func getAlarmActionRequests(params map[string]interface{}, paramKey string) []alarmActionRequest {
	raw, ok := params[paramKey]
	if !ok || raw == nil {
		return nil
	}
	arr, ok := raw.([]interface{})
	if !ok {
		return nil
	}
	result := make([]alarmActionRequest, 0, len(arr))
	for _, item := range arr {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		result = append(result, alarmActionRequest{
			alarmModelName: getStringFromMap(m, "alarmModelName"),
			keyValue:       getStringFromMap(m, "keyValue"),
		})
	}
	return result
}

// getAcknowledgeActionRequests extracts acknowledge action requests.
func getAcknowledgeActionRequests(params map[string]interface{}) []alarmActionRequest {
	return getAlarmActionRequests(params, "acknowledgeActionRequests")
}

// getStringFromMap safely extracts a string value from a map.
func getStringFromMap(m map[string]interface{}, key string) string {
	v, ok := m[key]
	if !ok || v == nil {
		return ""
	}
	s, ok := v.(string)
	if !ok {
		return ""
	}
	return s
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

// detectorRequest holds the parsed fields from a BatchDeleteDetector entry.
type detectorRequest struct {
	modelName string
	keyValue  string
}

// getDetectorRequests extracts detector requests from the parameters for
// the given parameter key (typically "detectors").
func getDetectorRequests(params map[string]interface{}, paramKey string) []detectorRequest {
	raw, ok := params[paramKey]
	if !ok || raw == nil {
		return nil
	}
	arr, ok := raw.([]interface{})
	if !ok {
		return nil
	}
	result := make([]detectorRequest, 0, len(arr))
	for _, item := range arr {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		result = append(result, detectorRequest{
			modelName: getStringFromMap(m, "detectorModelName"),
			keyValue:  getStringFromMap(m, "keyValue"),
		})
	}
	return result
}

// detectorUpdateRequest holds the parsed fields from a BatchUpdateDetector entry.
type detectorUpdateRequest struct {
	modelName string
	keyValue  string
	stateName string
}

// getDetectorUpdateRequests extracts detector update requests including the
// target state name from the nested state object.
func getDetectorUpdateRequests(params map[string]interface{}) []detectorUpdateRequest {
	raw, ok := params["detectors"]
	if !ok || raw == nil {
		return nil
	}
	arr, ok := raw.([]interface{})
	if !ok {
		return nil
	}
	result := make([]detectorUpdateRequest, 0, len(arr))
	for _, item := range arr {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		update := detectorUpdateRequest{
			modelName: getStringFromMap(m, "detectorModelName"),
			keyValue:  getStringFromMap(m, "keyValue"),
		}
		if state, ok := m["state"].(map[string]interface{}); ok {
			update.stateName = getStringFromMap(state, "stateName")
		}
		result = append(result, update)
	}
	return result
}
