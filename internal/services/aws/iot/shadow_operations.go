package iot

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"vorpalstacks/internal/common/request"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

const maxShadowSizeBytes = 8192

// shadowState represents the state section of an IoT Thing Shadow document
// containing reported, desired, and delta values. The AWS IoT Shadow
// protocol wraps this under a top-level "state" key in the request body.
type shadowState struct {
	Reported map[string]interface{} `json:"reported,omitempty"`
	Desired  map[string]interface{} `json:"desired,omitempty"`
	Delta    map[string]interface{} `json:"delta,omitempty"`
}

func (s *IoTService) GetThingShadow(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	thingName := request.GetParamCaseInsensitive(req.Parameters, "thingName")
	shadowName := request.GetParamCaseInsensitive(req.Parameters, "shadowName")

	if thingName == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	doc, err := store.GetShadow(thingName, shadowName)
	if err != nil {
		return nil, iotstore.ErrShadowNotFound
	}

	var state shadowState
	if err := json.Unmarshal([]byte(doc.State), &state); err != nil {
		return nil, iotstore.ErrInvalidRequest
	}

	return map[string]interface{}{
		"state": map[string]interface{}{
			"reported": ensureGenericMap(state.Reported),
			"desired":  ensureGenericMap(state.Desired),
			"delta":    ensureGenericMap(state.Delta),
		},
		"version":   doc.VersionNumber,
		"timestamp": doc.Timestamp.Unix(),
		"metadata":  unmarshalMetadata(doc.Metadata),
	}, nil
}

func (s *IoTService) UpdateThingShadow(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	thingName := request.GetParamCaseInsensitive(req.Parameters, "thingName")
	shadowName := request.GetParamCaseInsensitive(req.Parameters, "shadowName")

	if thingName == "" {
		return nil, iotstore.ErrMissingParam
	}

	if len(req.Body) > maxShadowSizeBytes {
		return nil, iotstore.ErrShadowTooLarge
	}

	var payload struct {
		State   shadowState `json:"state"`
		Version *int64      `json:"version"`
	}
	if err := json.Unmarshal(req.Body, &payload); err != nil {
		return nil, iotstore.ErrInvalidRequest
	}
	incoming := payload.State

	var clientVersion int64
	// Prefer the version from the parsed body JSON (correct source for
	// REST-JSON protocol). Fall back to req.Parameters for clients that
	// send version as a query parameter or via awsJson1_1 protocol.
	if payload.Version != nil {
		clientVersion = *payload.Version
	} else if raw := request.GetParamCaseInsensitive(req.Parameters, "version"); raw != "" {
		if n, _ := fmt.Sscanf(raw, "%d", &clientVersion); n != 1 {
			return nil, iotstore.ErrInvalidRequest
		}
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	existing, getErr := store.GetShadow(thingName, shadowName)
	var state shadowState
	var version int64
	if getErr == nil {
		if err := json.Unmarshal([]byte(existing.State), &state); err != nil {
			return nil, iotstore.ErrInvalidRequest
		}
		version = existing.VersionNumber
	} else {
		state = shadowState{
			Reported: make(map[string]interface{}),
			Desired:  make(map[string]interface{}),
		}
		version = 0
	}

	mergeShadowState(&state, incoming)
	version++

	delta := computeDelta(state.Desired, state.Reported)
	state.Delta = delta

	stateJSON, _ := json.Marshal(state)
	if len(stateJSON) > maxShadowSizeBytes {
		return nil, iotstore.ErrShadowTooLarge
	}
	metaJSON := ""
	if existing != nil {
		metaJSON = existing.Metadata
	}

	doc := buildShadowDocument(thingName, shadowName, string(stateJSON), metaJSON, version, time.Unix(time.Now().Unix(), 0))
	if err := store.PutShadowWithVersion(thingName, shadowName, doc, clientVersion); err != nil {
		return nil, err
	}

	if len(delta) > 0 && s.broker != nil {
		deltaTopic := buildDeltaTopic(thingName, shadowName)
		deltaPayload, _ := json.Marshal(map[string]interface{}{
			"state": map[string]interface{}{
				"delta": delta,
			},
			"version":   version,
			"timestamp": time.Now().Unix(),
			"thingName": thingName,
		})
		_ = s.broker.Publish(deltaTopic, deltaPayload)
	}

	return map[string]interface{}{
		"state": map[string]interface{}{
			"reported": ensureGenericMap(state.Reported),
			"desired":  ensureGenericMap(state.Desired),
			"delta":    ensureGenericMap(state.Delta),
		},
		"version":   version,
		"timestamp": time.Now().Unix(),
		"metadata":  unmarshalMetadata(metaJSON),
	}, nil
}

func (s *IoTService) DeleteThingShadow(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	thingName := request.GetParamCaseInsensitive(req.Parameters, "thingName")
	shadowName := request.GetParamCaseInsensitive(req.Parameters, "shadowName")

	if thingName == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.DeleteShadow(thingName, shadowName); err != nil {
		return nil, iotstore.ErrShadowNotFound
	}

	return map[string]interface{}{}, nil
}

func (s *IoTService) ListNamedShadowsForThing(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	thingName := request.GetParamCaseInsensitive(req.Parameters, "thingName")
	if thingName == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	names, err := store.ListShadowNames(thingName)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"results": names,
	}, nil
}

func mergeShadowState(target *shadowState, incoming shadowState) {
	if incoming.Reported != nil {
		if target.Reported == nil {
			target.Reported = make(map[string]interface{})
		}
		mergeMap(target.Reported, incoming.Reported)
	}
	if incoming.Desired != nil {
		if target.Desired == nil {
			target.Desired = make(map[string]interface{})
		}
		mergeMap(target.Desired, incoming.Desired)
	}
}

func mergeMap(target, source map[string]interface{}) {
	for k, v := range source {
		if v == nil {
			delete(target, k)
			continue
		}
		subSrc, srcIsMap := v.(map[string]interface{})
		subTgt, tgtIsMap := target[k]
		if srcIsMap && tgtIsMap {
			if tgtMap, ok := subTgt.(map[string]interface{}); ok {
				mergeMap(tgtMap, subSrc)
				continue
			}
		}
		target[k] = deepCopyValue(v)
	}
}

func computeDelta(desired, reported map[string]interface{}) map[string]interface{} {
	if len(desired) == 0 {
		return nil
	}
	delta := make(map[string]interface{})
	computeDeltaRecursive(desired, reported, delta)
	if len(delta) == 0 {
		return nil
	}
	return delta
}

func computeDeltaRecursive(desired, reported, delta map[string]interface{}) {
	for key, desiredVal := range desired {
		reportedVal, exists := reported[key]
		if !exists {
			delta[key] = desiredVal
			continue
		}
		desiredMap, desiredIsMap := desiredVal.(map[string]interface{})
		reportedMap, reportedIsMap := reportedVal.(map[string]interface{})
		if desiredIsMap && reportedIsMap {
			subDelta := make(map[string]interface{})
			computeDeltaRecursive(desiredMap, reportedMap, subDelta)
			if len(subDelta) > 0 {
				delta[key] = subDelta
			}
			continue
		}
		aj, _ := json.Marshal(desiredVal)
		bj, _ := json.Marshal(reportedVal)
		if string(aj) != string(bj) {
			delta[key] = desiredVal
		}
	}
}

func buildShadowDocument(thingName, shadowName, stateJSON, metaJSON string, version int64, timestamp time.Time) *iotstore.ShadowDocument {
	return &iotstore.ShadowDocument{
		ThingName:     thingName,
		VersionNumber: version,
		Timestamp:     timestamp,
		State:         stateJSON,
		Metadata:      metaJSON,
	}
}

func ensureGenericMap(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		return map[string]interface{}{}
	}
	return m
}

func unmarshalMetadata(raw string) interface{} {
	if raw == "" {
		return nil
	}
	var result interface{}
	if json.Unmarshal([]byte(raw), &result) == nil {
		return result
	}
	return nil
}

func buildDeltaTopic(thingName, shadowName string) string {
	if shadowName == "" {
		return fmt.Sprintf("$aws/things/%s/shadow/update/delta", thingName)
	}
	return fmt.Sprintf("$aws/things/%s/shadow/name/%s/update/delta", thingName, shadowName)
}

func deepCopyValue(v interface{}) interface{} {
	switch val := v.(type) {
	case map[string]interface{}:
		cp := make(map[string]interface{}, len(val))
		for k, child := range val {
			cp[k] = deepCopyValue(child)
		}
		return cp
	default:
		return v
	}
}
