package iot

import (
	"encoding/json"
	"fmt"
	"time"

	pb "vorpalstacks/internal/pb/storage/storage_iot"
	"vorpalstacks/internal/store/aws/common"
)

type ShadowState struct {
	Reported map[string]interface{} `json:"reported,omitempty"`
	Desired  map[string]interface{} `json:"desired,omitempty"`
	Delta    map[string]interface{} `json:"delta,omitempty"`
}

// ShadowUpdateResult is returned by UpdateShadow to give the service layer
// everything it needs to build a response and publish the delta MQTT message.
type ShadowUpdateResult struct {
	Document *ShadowDocument
	Delta    map[string]interface{}
}

func mergeShadowState(target *ShadowState, incoming ShadowState) {
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

func deepCopyValue(v interface{}) interface{} {
	switch val := v.(type) {
	case map[string]interface{}:
		cp := make(map[string]interface{}, len(val))
		for k, child := range val {
			cp[k] = deepCopyValue(child)
		}
		return cp
	case []interface{}:
		cp := make([]interface{}, len(val))
		for i, child := range val {
			cp[i] = deepCopyValue(child)
		}
		return cp
	default:
		return v
	}
}

// UpdateShadow performs an atomic read-merge-version-delta-write for a device
// shadow. It reads the existing shadow (if any), merges the incoming reported
// and desired state, increments the version, computes the delta between desired
// and reported, and writes the result back with optimistic concurrency control.
func (s *IotStore) UpdateShadow(thingName, shadowName string, incoming ShadowState, clientVersion int64) (*ShadowUpdateResult, error) {
	lockKey := shadowKey(thingName, shadowName)
	s.shadowLocker.Lock(lockKey)
	defer s.shadowLocker.Unlock(lockKey)

	var state ShadowState
	var version int64
	var metaJSON string

	pbDoc := &pb.ShadowDocument{}
	err := s.shadowsBase.GetProto(shadowKey(thingName, shadowName), pbDoc)
	if err == nil {
		existing := ProtoToShadow(pbDoc)
		if err := json.Unmarshal([]byte(existing.State), &state); err != nil {
			return nil, fmt.Errorf("failed to parse existing shadow state: %w", err)
		}
		version = existing.VersionNumber
		metaJSON = existing.Metadata
	} else if !common.IsNotFound(err) {
		return nil, err
	} else {
		state = ShadowState{
			Reported: make(map[string]interface{}),
			Desired:  make(map[string]interface{}),
		}
	}

	if clientVersion > 0 && version != clientVersion {
		return nil, ErrVersionConflict
	}

	mergeShadowState(&state, incoming)
	version++

	delta := computeDelta(state.Desired, state.Reported)
	state.Delta = delta

	stateJSON, err := json.Marshal(state)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal shadow state: %w", err)
	}

	doc := &ShadowDocument{
		ThingName:     thingName,
		VersionNumber: version,
		Timestamp:     time.Now().UTC(),
		State:         string(stateJSON),
		Metadata:      metaJSON,
	}

	if err := s.shadowsBase.PutProto(shadowKey(thingName, shadowName), ShadowToProto(doc)); err != nil {
		return nil, err
	}

	return &ShadowUpdateResult{Document: doc, Delta: delta}, nil
}
