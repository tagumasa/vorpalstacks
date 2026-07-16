package iot

import (
	"testing"
)

func TestComputeDeltaNoDesired(t *testing.T) {
	delta := computeDelta(nil, map[string]interface{}{"temp": 25})
	if delta != nil {
		t.Errorf("expected nil delta for empty desired, got %v", delta)
	}
}

func TestComputeDeltaEmptyMaps(t *testing.T) {
	delta := computeDelta(map[string]interface{}{}, map[string]interface{}{})
	if delta != nil {
		t.Errorf("expected nil for empty maps, got %v", delta)
	}
}

func TestComputeDeltaNoReported(t *testing.T) {
	desired := map[string]interface{}{"temp": "30"}
	delta := computeDelta(desired, nil)
	if delta == nil {
		t.Fatal("expected non-nil delta")
	}
	if delta["temp"] != "30" {
		t.Errorf("delta[temp] = %v", delta["temp"])
	}
}

func TestComputeDeltaMatching(t *testing.T) {
	desired := map[string]interface{}{"temp": "25"}
	reported := map[string]interface{}{"temp": "25"}
	delta := computeDelta(desired, reported)
	if delta != nil {
		t.Errorf("expected nil delta when desired==reported, got %v", delta)
	}
}

func TestComputeDeltaDifferent(t *testing.T) {
	desired := map[string]interface{}{"temp": "30"}
	reported := map[string]interface{}{"temp": "25"}
	delta := computeDelta(desired, reported)
	if delta == nil {
		t.Fatal("expected non-nil delta")
	}
	if delta["temp"] != "30" {
		t.Errorf("delta[temp] = %v, want '30'", delta["temp"])
	}
}

func TestComputeDeltaNestedDifferent(t *testing.T) {
	desired := map[string]interface{}{
		"config": map[string]interface{}{"mode": "active"},
	}
	reported := map[string]interface{}{
		"config": map[string]interface{}{"mode": "idle"},
	}
	delta := computeDelta(desired, reported)
	if delta == nil {
		t.Fatal("expected non-nil delta")
	}
	configDelta, ok := delta["config"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected nested map in delta, got %T", delta["config"])
	}
	if configDelta["mode"] != "active" {
		t.Errorf("delta[config][mode] = %v", configDelta["mode"])
	}
}

func TestComputeDeltaNestedMatching(t *testing.T) {
	desired := map[string]interface{}{
		"config": map[string]interface{}{"mode": "active"},
	}
	reported := map[string]interface{}{
		"config": map[string]interface{}{"mode": "active"},
	}
	delta := computeDelta(desired, reported)
	if delta != nil {
		t.Errorf("expected nil delta for matching nested, got %v", delta)
	}
}

func TestComputeDeltaExtraInReported(t *testing.T) {
	desired := map[string]interface{}{"temp": 25}
	reported := map[string]interface{}{"temp": 25, "humidity": 60}
	delta := computeDelta(desired, reported)
	if delta != nil {
		t.Errorf("expected nil delta (extra reported fields ignored), got %v", delta)
	}
}

func TestMergeShadowStateReported(t *testing.T) {
	target := &ShadowState{
		Reported: map[string]interface{}{"temp": "20"},
	}
	incoming := ShadowState{
		Reported: map[string]interface{}{"temp": "25", "humidity": "60"},
	}
	mergeShadowState(target, incoming)

	if target.Reported["temp"] != "25" {
		t.Errorf("reported[temp] = %v", target.Reported["temp"])
	}
	if target.Reported["humidity"] != "60" {
		t.Errorf("reported[humidity] = %v", target.Reported["humidity"])
	}
}

func TestMergeShadowStateDesired(t *testing.T) {
	target := &ShadowState{}
	incoming := ShadowState{
		Desired: map[string]interface{}{"mode": "active"},
	}
	mergeShadowState(target, incoming)

	if target.Desired == nil {
		t.Fatal("expected desired to be created")
	}
	if target.Desired["mode"] != "active" {
		t.Errorf("desired[mode] = %v", target.Desired["mode"])
	}
}

func TestMergeShadowStateNilTarget(t *testing.T) {
	target := &ShadowState{}
	incoming := ShadowState{
		Reported: map[string]interface{}{"val": "r"},
		Desired:  map[string]interface{}{"val": "d"},
	}
	mergeShadowState(target, incoming)

	if target.Reported["val"] != "r" {
		t.Errorf("reported[val] = %v", target.Reported["val"])
	}
	if target.Desired["val"] != "d" {
		t.Errorf("desired[val] = %v", target.Desired["val"])
	}
}

func TestMergeMapNilDeletesKey(t *testing.T) {
	target := map[string]interface{}{"a": "1", "b": "2"}
	source := map[string]interface{}{"b": nil}

	mergeMap(target, source)

	if target["a"] != "1" {
		t.Errorf("a = %v, want '1'", target["a"])
	}
	if _, exists := target["b"]; exists {
		t.Error("expected 'b' to be deleted (nil value in source)")
	}
}

func TestMergeMapNestedMerge(t *testing.T) {
	target := map[string]interface{}{
		"config": map[string]interface{}{"mode": "old", "keep": true},
	}
	source := map[string]interface{}{
		"config": map[string]interface{}{"mode": "new"},
	}

	mergeMap(target, source)

	config, ok := target["config"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected config to be map, got %T", target["config"])
	}
	if config["mode"] != "new" {
		t.Errorf("config[mode] = %v, want 'new'", config["mode"])
	}
	if config["keep"] != true {
		t.Errorf("config[keep] = %v, want true", config["keep"])
	}
}

func TestMergeMapReplaceNonMap(t *testing.T) {
	target := map[string]interface{}{"val": "old"}
	source := map[string]interface{}{"val": "new"}

	mergeMap(target, source)

	if target["val"] != "new" {
		t.Errorf("val = %v, want 'new'", target["val"])
	}
}

func TestDeepCopyValueMap(t *testing.T) {
	original := map[string]interface{}{
		"a": "1",
		"b": map[string]interface{}{"c": "2"},
	}
	copied := deepCopyValue(original)

	cp, ok := copied.(map[string]interface{})
	if !ok {
		t.Fatalf("expected map, got %T", copied)
	}
	cp["a"] = "changed"
	if original["a"] != "1" {
		t.Error("expected original to be unchanged after modifying copy")
	}
}

func TestDeepCopyValueSlice(t *testing.T) {
	original := []interface{}{1, 2, map[string]interface{}{"k": "v"}}
	copied := deepCopyValue(original)

	cp, ok := copied.([]interface{})
	if !ok {
		t.Fatalf("expected slice, got %T", copied)
	}
	if len(cp) != 3 {
		t.Fatalf("len = %d, want 3", len(cp))
	}
	cp[0] = 999
	if original[0] == 999 {
		t.Error("expected original slice to be unchanged")
	}
}

func TestDeepCopyValuePrimitive(t *testing.T) {
	if v := deepCopyValue(42); v != 42 {
		t.Errorf("deepCopyValue(42) = %v", v)
	}
	if v := deepCopyValue("hello"); v != "hello" {
		t.Errorf("deepCopyValue(hello) = %v", v)
	}
	if v := deepCopyValue(nil); v != nil {
		t.Errorf("deepCopyValue(nil) = %v", v)
	}
}
