package iot

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	storecommon "vorpalstacks/internal/store/aws/common"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

// Core functions for the device-shadow family. Handlers on both protocol
// planes are thin adapters: they parse the JSON payload and shadow name
// off the wire, call these functions, and serialise the view. Shadow-name
// validation, size limits, persistence and MQTT side effects live here.

// ShadowView is the transport-agnostic shadow document view.
type ShadowView struct {
	Reported         map[string]interface{}
	Desired          map[string]interface{}
	Delta            map[string]interface{}
	Version          int64
	TimestampSeconds int64
	Metadata         interface{}
}

// validateShadowName enforces the shadow-name charset for named shadows.
func validateShadowName(shadowName string) error {
	if shadowName != "" && !shadowNameRegex.MatchString(shadowName) {
		return iotstore.ErrInvalidRequest
	}
	return nil
}

// getThingShadowCore loads a shadow document and projects it onto the
// response view.
func (s *IoTService) getThingShadowCore(store iotstore.IotStoreInterface, thingName, shadowName string) (*ShadowView, error) {
	if thingName == "" {
		return nil, iotstore.ErrMissingParam
	}
	if err := validateShadowName(shadowName); err != nil {
		return nil, err
	}
	doc, err := store.GetShadow(thingName, shadowName)
	if err != nil {
		return nil, err
	}
	var state iotstore.ShadowState
	if err := json.Unmarshal([]byte(doc.State), &state); err != nil {
		return nil, iotstore.ErrInvalidRequest
	}
	return &ShadowView{
		Reported:         ensureGenericMap(state.Reported),
		Desired:          ensureGenericMap(state.Desired),
		Delta:            ensureGenericMap(state.Delta),
		Version:          doc.VersionNumber,
		TimestampSeconds: doc.Timestamp.Unix(),
		Metadata:         unmarshalMetadata(doc.Metadata),
	}, nil
}

// UpdateThingShadowInput carries the parsed update payload. The handler
// decodes the JSON body and the optional version member; semantics
// (version conflict, delta computation) live in the store and this Core.
type UpdateThingShadowInput struct {
	ThingName     string
	ShadowName    string
	Body          []byte
	ClientVersion int64
}

// updateThingShadowCore validates the payload size, persists the state
// update and publishes the delta notification to the region's broker.
func (s *IoTService) updateThingShadowCore(store iotstore.IotStoreInterface, region string, in UpdateThingShadowInput) (*ShadowView, error) {
	if in.ThingName == "" {
		return nil, iotstore.ErrMissingParam
	}
	if err := validateShadowName(in.ShadowName); err != nil {
		return nil, err
	}
	if len(in.Body) > maxShadowSizeBytes {
		return nil, iotstore.ErrShadowTooLarge
	}
	var payload struct {
		State   iotstore.ShadowState `json:"state"`
		Version *int64               `json:"version"`
	}
	if err := json.Unmarshal(in.Body, &payload); err != nil {
		return nil, iotstore.ErrInvalidRequest
	}
	clientVersion := in.ClientVersion
	if payload.Version != nil {
		clientVersion = *payload.Version
	}

	result, err := store.UpdateShadow(in.ThingName, in.ShadowName, payload.State, clientVersion)
	if err != nil {
		return nil, err
	}

	if brk := s.brokerForRegion(region); len(result.Delta) > 0 && brk != nil {
		deltaTopic := buildDeltaTopic(in.ThingName, in.ShadowName)
		deltaPayload, err := json.Marshal(map[string]interface{}{
			"state": map[string]interface{}{
				"delta": result.Delta,
			},
			"version":   result.Document.VersionNumber,
			"timestamp": time.Now().Unix(),
			"thingName": in.ThingName,
		})
		if err != nil {
			slog.Error("shadow delta marshal error", "error", err)
		} else if pubErr := brk.Publish(deltaTopic, deltaPayload); pubErr != nil {
			slog.Error("shadow delta publish error", "error", pubErr)
		}
	}

	var respState iotstore.ShadowState
	if unmarshalErr := json.Unmarshal([]byte(result.Document.State), &respState); unmarshalErr != nil {
		slog.Error("shadow state unmarshal error", "error", unmarshalErr)
	}

	return &ShadowView{
		Reported:         ensureGenericMap(respState.Reported),
		Desired:          ensureGenericMap(respState.Desired),
		Delta:            ensureGenericMap(respState.Delta),
		Version:          result.Document.VersionNumber,
		TimestampSeconds: result.Document.Timestamp.Unix(),
		Metadata:         unmarshalMetadata(result.Document.Metadata),
	}, nil
}

// deleteThingShadowCore removes a shadow (verifying it exists first:
// DeleteShadow uses Pebble's idempotent Delete which returns nil for
// missing keys, so without this check a non-existent shadow would
// silently succeed) and publishes the delete/accepted notification.
func (s *IoTService) deleteThingShadowCore(store iotstore.IotStoreInterface, region, thingName, shadowName string) error {
	if thingName == "" {
		return iotstore.ErrMissingParam
	}
	if err := validateShadowName(shadowName); err != nil {
		return err
	}

	// Verify the shadow exists before deleting; DeleteShadow uses Pebble's
	// idempotent Delete which returns nil for missing keys, so without this
	// check a non-existent shadow would silently succeed.
	if _, err := store.GetShadow(thingName, shadowName); err != nil {
		return err
	}
	if err := store.DeleteShadow(thingName, shadowName); err != nil {
		return err
	}

	// Publish delete/accepted notification (AWS IoT behaviour).
	if brk := s.brokerForRegion(region); brk != nil {
		deleteTopic := fmt.Sprintf("$aws/things/%s/shadow/delete/accepted", thingName)
		if shadowName != "" {
			deleteTopic = fmt.Sprintf("$aws/things/%s/shadow/name/%s/delete/accepted", thingName, shadowName)
		}
		notif, _ := json.Marshal(map[string]interface{}{
			"timestamp": time.Now().Unix(),
			"thingName": thingName,
		})
		if pubErr := brk.Publish(deleteTopic, notif); pubErr != nil {
			slog.Error("shadow delete notification publish error", "error", pubErr)
		}
	}
	return nil
}

// listNamedShadowsForThingCore lists a thing's named shadows page by page.
func (s *IoTService) listNamedShadowsForThingCore(store iotstore.IotStoreInterface, thingName string, opts storecommon.ListOptions) ([]string, string, error) {
	if thingName == "" {
		return nil, "", iotstore.ErrMissingParam
	}
	return store.ListShadowNames(thingName, opts)
}
