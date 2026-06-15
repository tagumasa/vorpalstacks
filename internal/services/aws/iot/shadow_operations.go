package iot

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"vorpalstacks/internal/common/request"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

const maxShadowSizeBytes = 8192

var shadowNameRegex = regexp.MustCompile(`^[a-zA-Z0-9_\.\-]+$`)

func (s *IoTService) GetThingShadow(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	thingName := request.GetParamCaseInsensitive(req.Parameters, "thingName")
	shadowName := request.GetParamCaseInsensitive(req.Parameters, "shadowName")

	if thingName == "" {
		return nil, iotstore.ErrMissingParam
	}
	if shadowName != "" && !shadowNameRegex.MatchString(shadowName) {
		return nil, iotstore.ErrInvalidRequest
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	doc, err := store.GetShadow(thingName, shadowName)
	if err != nil {
		return nil, iotstore.ErrShadowNotFound
	}

	var state iotstore.ShadowState
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
	if shadowName != "" && !shadowNameRegex.MatchString(shadowName) {
		return nil, iotstore.ErrInvalidRequest
	}

	if len(req.Body) > maxShadowSizeBytes {
		return nil, iotstore.ErrShadowTooLarge
	}

	var payload struct {
		State   iotstore.ShadowState `json:"state"`
		Version *int64               `json:"version"`
	}
	if err := json.Unmarshal(req.Body, &payload); err != nil {
		return nil, iotstore.ErrInvalidRequest
	}

	var clientVersion int64
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

	result, err := store.UpdateShadow(thingName, shadowName, payload.State, clientVersion)
	if err != nil {
		return nil, err
	}

	if len(result.Delta) > 0 && s.deps.Broker != nil {
		deltaTopic := buildDeltaTopic(thingName, shadowName)
		deltaPayload, err := json.Marshal(map[string]interface{}{
			"state": map[string]interface{}{
				"delta": result.Delta,
			},
			"version":   result.Document.VersionNumber,
			"timestamp": time.Now().Unix(),
			"thingName": thingName,
		})
		if err != nil {
			fmt.Printf("[iot] shadow delta marshal error: %v\n", err)
		} else if pubErr := s.deps.Broker.Publish(deltaTopic, deltaPayload); pubErr != nil {
			fmt.Printf("[iot] shadow delta publish error: %v\n", pubErr)
		}
	}

	var respState iotstore.ShadowState
	if unmarshalErr := json.Unmarshal([]byte(result.Document.State), &respState); unmarshalErr != nil {
		fmt.Printf("[iot] shadow state unmarshal error: %v\n", unmarshalErr)
	}

	return map[string]interface{}{
		"state": map[string]interface{}{
			"reported": ensureGenericMap(respState.Reported),
			"desired":  ensureGenericMap(respState.Desired),
			"delta":    ensureGenericMap(respState.Delta),
		},
		"version":   result.Document.VersionNumber,
		"timestamp": result.Document.Timestamp.Unix(),
		"metadata":  unmarshalMetadata(result.Document.Metadata),
	}, nil
}

func (s *IoTService) DeleteThingShadow(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	thingName := request.GetParamCaseInsensitive(req.Parameters, "thingName")
	shadowName := request.GetParamCaseInsensitive(req.Parameters, "shadowName")

	if thingName == "" {
		return nil, iotstore.ErrMissingParam
	}
	if shadowName != "" && !shadowNameRegex.MatchString(shadowName) {
		return nil, iotstore.ErrInvalidRequest
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
