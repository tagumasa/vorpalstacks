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

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	view, err := s.getThingShadowCore(store, thingName, shadowName)
	if err != nil {
		return nil, err
	}

	return shadowViewResponse(view), nil
}

func (s *IoTService) UpdateThingShadow(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	thingName := request.GetParamCaseInsensitive(req.Parameters, "thingName")
	shadowName := request.GetParamCaseInsensitive(req.Parameters, "shadowName")

	var clientVersion int64
	if raw := request.GetParamCaseInsensitive(req.Parameters, "version"); raw != "" {
		if n, _ := fmt.Sscanf(raw, "%d", &clientVersion); n != 1 {
			return nil, iotstore.ErrInvalidRequest
		}
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	view, err := s.updateThingShadowCore(store, reqCtx.GetRegion(), UpdateThingShadowInput{
		ThingName:     thingName,
		ShadowName:    shadowName,
		Body:          req.Body,
		ClientVersion: clientVersion,
	})
	if err != nil {
		return nil, err
	}

	return shadowViewResponse(view), nil
}

func (s *IoTService) DeleteThingShadow(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	thingName := request.GetParamCaseInsensitive(req.Parameters, "thingName")
	shadowName := request.GetParamCaseInsensitive(req.Parameters, "shadowName")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.deleteThingShadowCore(store, reqCtx.GetRegion(), thingName, shadowName); err != nil {
		return nil, err
	}

	return map[string]interface{}{}, nil
}

func (s *IoTService) ListNamedShadowsForThing(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	thingName := request.GetParamCaseInsensitive(req.Parameters, "thingName")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	opts := parseListOptions(req.Parameters)
	names, nextToken, err := s.listNamedShadowsForThingCore(store, thingName, opts)
	if err != nil {
		return nil, err
	}

	resp := map[string]interface{}{
		"results":   names,
		"timestamp": time.Now().Unix(),
	}
	if nextToken != "" {
		resp["nextToken"] = nextToken
	}
	return resp, nil
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

// shadowViewResponse serialises a Core shadow view onto the wire shape.
func shadowViewResponse(view *ShadowView) map[string]interface{} {
	return map[string]interface{}{
		"state": map[string]interface{}{
			"reported": view.Reported,
			"desired":  view.Desired,
			"delta":    view.Delta,
		},
		"version":   view.Version,
		"timestamp": view.TimestampSeconds,
		"metadata":  view.Metadata,
	}
}
