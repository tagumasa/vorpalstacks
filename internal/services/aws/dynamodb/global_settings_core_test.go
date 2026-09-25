package dynamodb

import (
	"context"
	"errors"
	"testing"

	"vorpalstacks/internal/common/request"
)

// TestParseAutoScalingSettingsTargetTracking pins the description form of
// the scaling policy: the target tracking configuration travels with the
// policy name, and its required TargetValue member is enforced.
func TestParseAutoScalingSettingsTargetTracking(t *testing.T) {
	desc, err := parseAutoScalingSettings(map[string]interface{}{
		"MinimumUnits": 5.0,
		"MaximumUnits": 50.0,
		"ScalingPolicyUpdate": map[string]interface{}{
			"PolicyName": "write-tracking",
			"TargetTrackingScalingPolicyConfiguration": map[string]interface{}{
				"TargetValue":      50.0,
				"DisableScaleIn":   true,
				"ScaleInCooldown":  60.0,
				"ScaleOutCooldown": 30.0,
			},
		},
	})
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	if len(desc.ScalingPolicies) != 1 {
		t.Fatalf("expected one scaling policy, got %#v", desc.ScalingPolicies)
	}
	policy := desc.ScalingPolicies[0]
	if policy.PolicyName == nil || *policy.PolicyName != "write-tracking" {
		t.Errorf("policy name = %#v", policy.PolicyName)
	}
	tt := policy.TargetTrackingScalingPolicyConfiguration
	if tt == nil {
		t.Fatalf("expected a target tracking configuration, got nil")
	}
	if tt.TargetValue != 50.0 {
		t.Errorf("target value = %#v", tt.TargetValue)
	}
	if tt.DisableScaleIn == nil || !*tt.DisableScaleIn {
		t.Errorf("disable scale-in = %#v", tt.DisableScaleIn)
	}
	if tt.ScaleInCooldown == nil || *tt.ScaleInCooldown != 60 {
		t.Errorf("scale-in cooldown = %#v", tt.ScaleInCooldown)
	}
	if tt.ScaleOutCooldown == nil || *tt.ScaleOutCooldown != 30 {
		t.Errorf("scale-out cooldown = %#v", tt.ScaleOutCooldown)
	}

	// A missing TargetValue is rejected.
	if _, err := parseAutoScalingSettings(map[string]interface{}{
		"ScalingPolicyUpdate": map[string]interface{}{
			"PolicyName": "no-target",
			"TargetTrackingScalingPolicyConfiguration": map[string]interface{}{
				"DisableScaleIn": true,
			},
		},
	}); err == nil {
		t.Error("expected an error for a missing TargetValue")
	}
}

// TestGlobalGSIWriteSettingsListBounds pins the model's 1-20 entry bound
// on the global GSI settings update list.
func TestGlobalGSIWriteSettingsListBounds(t *testing.T) {
	if _, err := parseGlobalGSIWriteSettings([]interface{}{}); err == nil {
		t.Error("expected an empty global GSI settings list to be rejected")
	}
	entries := make([]interface{}, 21)
	for i := range entries {
		entries[i] = map[string]interface{}{"IndexName": "gsi"}
	}
	if _, err := parseGlobalGSIWriteSettings(entries); err == nil {
		t.Error("expected a 21-entry global GSI settings list to be rejected")
	}
}

// TestReplicaGSIReadSettingsListBounds pins the model's 1-20 entry bound
// on the per-replica GSI settings update list.
func TestReplicaGSIReadSettingsListBounds(t *testing.T) {
	if _, err := parseReplicaGSIReadSettings([]interface{}{}); err == nil {
		t.Error("expected an empty replica GSI settings list to be rejected")
	}
	entries := make([]interface{}, 21)
	for i := range entries {
		entries[i] = map[string]interface{}{"IndexName": "gsi"}
	}
	if _, err := parseReplicaGSIReadSettings(entries); err == nil {
		t.Error("expected a 21-entry replica GSI settings list to be rejected")
	}
}

// TestGSIWriteReadSettingsMemberExtraction pins the per-member contract of
// both settings parsers: a missing IndexName is rejected, the capacity
// member must be an integral positive long, and each side reads only its
// own capacity and auto-scaling members into its own fields.
func TestGSIWriteReadSettingsMemberExtraction(t *testing.T) {
	if _, err := parseGlobalGSIWriteSettings([]interface{}{map[string]interface{}{
		"ProvisionedWriteCapacityUnits": 5.0,
	}}); err == nil {
		t.Error("expected a write entry without IndexName to be rejected")
	}
	if _, err := parseReplicaGSIReadSettings([]interface{}{map[string]interface{}{
		"IndexName": "gsi", "ProvisionedReadCapacityUnits": 0.5,
	}}); err == nil {
		t.Error("expected a fractional capacity member to be rejected")
	}
	if _, err := parseReplicaGSIReadSettings([]interface{}{map[string]interface{}{
		"IndexName": "gsi", "ProvisionedReadCapacityUnits": 0.0,
	}}); err == nil {
		t.Error("expected a zero capacity member to be rejected")
	}

	w, err := parseGlobalGSIWriteSettings([]interface{}{map[string]interface{}{
		"IndexName":                     "gsi",
		"ProvisionedWriteCapacityUnits": 7.0,
		"ProvisionedWriteCapacityAutoScalingSettingsUpdate": map[string]interface{}{"MinimumUnits": 1.0, "MaximumUnits": 9.0},
	}})
	if err != nil || len(w) != 1 {
		t.Fatalf("write-side happy path: entries = %v, err = %v", w, err)
	}
	if w[0].IndexName != "gsi" || w[0].ProvisionedWriteCapacityUnits == nil || *w[0].ProvisionedWriteCapacityUnits != 7 {
		t.Errorf("write entry = %+v", w[0])
	}
	if w[0].Write == nil || w[0].Write.MinimumUnits == nil || *w[0].Write.MinimumUnits != 1 {
		t.Errorf("write auto-scaling = %+v", w[0].Write)
	}

	r, err := parseReplicaGSIReadSettings([]interface{}{map[string]interface{}{
		"IndexName":                    "gsi",
		"ProvisionedReadCapacityUnits": 6.0,
		"ProvisionedReadCapacityAutoScalingSettingsUpdate": map[string]interface{}{"MinimumUnits": 2.0, "MaximumUnits": 8.0},
	}})
	if err != nil || len(r) != 1 {
		t.Fatalf("read-side happy path: entries = %v, err = %v", r, err)
	}
	if r[0].IndexName != "gsi" || r[0].ProvisionedReadCapacityUnits == nil || *r[0].ProvisionedReadCapacityUnits != 6 {
		t.Errorf("read entry = %+v", r[0])
	}
	if r[0].Read == nil || r[0].Read.MaximumUnits == nil || *r[0].Read.MaximumUnits != 8 {
		t.Errorf("read auto-scaling = %+v", r[0].Read)
	}
}

// TestGlobalSettingsUpdatesRejectWrongTypedMembers pins the typed-member
// presence contract on the global table settings face: the settings lists
// and their AutoScalingSettingsUpdate members reject a present wrong-typed
// value — the invalid-parameter error — instead of silently skipping it,
// while an absent member and a member-less entry keep the documented skip
// semantics.
func TestGlobalSettingsUpdatesRejectWrongTypedMembers(t *testing.T) {
	// The per-index settings lists: a present non-list carrier, a
	// non-map entry, and a wrong-typed settings member are all rejected.
	if _, err := parseGlobalGSIWriteSettings("not-a-list"); !errors.Is(err, ErrInvalidParameter) {
		t.Errorf("non-list GSI settings member: expected the invalid-parameter rejection, got %v", err)
	}
	if _, err := parseGlobalGSIWriteSettings([]interface{}{"not-a-map"}); !errors.Is(err, ErrInvalidParameter) {
		t.Errorf("non-map GSI settings entry: expected the invalid-parameter rejection, got %v", err)
	}
	if _, err := parseGlobalGSIWriteSettings([]interface{}{map[string]interface{}{
		"IndexName": "gsi",
		"ProvisionedWriteCapacityAutoScalingSettingsUpdate": "not-a-structure",
	}}); !errors.Is(err, ErrInvalidParameter) {
		t.Errorf("non-structure write settings member: expected the invalid-parameter rejection, got %v", err)
	}
	if _, err := parseReplicaGSIReadSettings([]interface{}{map[string]interface{}{
		"IndexName": "gsi",
		"ProvisionedReadCapacityAutoScalingSettingsUpdate": 42.0,
	}}); !errors.Is(err, ErrInvalidParameter) {
		t.Errorf("non-structure read settings member: expected the invalid-parameter rejection, got %v", err)
	}

	// An absent member and a member-less entry keep the documented skip
	// semantics: nothing to apply, no error.
	if updates, err := parseGlobalGSIWriteSettings(nil); err != nil || updates != nil {
		t.Fatalf("absent GSI settings member: %v, %+v", err, updates)
	}
	if updates, err := parseGlobalGSIWriteSettings([]interface{}{map[string]interface{}{"IndexName": "gsi"}}); err != nil || len(updates) != 1 || updates[0].IndexName != "gsi" {
		t.Fatalf("entry without a settings member: %v, %+v", err, updates)
	}

	// The replica settings list: a present non-list carrier, a non-map
	// entry, and a wrong-typed read settings member are rejected before
	// any store access — the global table need not exist for the parse
	// rejection to answer.
	svc, reqCtx := billingModePlaneFixture(t)
	settings := func(replicaUpdates interface{}) error {
		_, err := svc.UpdateGlobalTableSettings(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"GlobalTableName":       "GtSettingsTypeFaceTable",
			"ReplicaSettingsUpdate": replicaUpdates,
		}})
		return err
	}
	if err := settings("not-a-list"); !errors.Is(err, ErrInvalidParameter) {
		t.Errorf("non-list ReplicaSettingsUpdate: expected the invalid-parameter rejection, got %v", err)
	}
	if err := settings([]interface{}{"not-a-map"}); !errors.Is(err, ErrInvalidParameter) {
		t.Errorf("non-map replica settings entry: expected the invalid-parameter rejection, got %v", err)
	}
	if err := settings([]interface{}{map[string]interface{}{
		"RegionName": reqCtx.GetRegion(),
		"ReplicaProvisionedReadCapacityAutoScalingSettingsUpdate": "not-a-structure",
	}}); !errors.Is(err, ErrInvalidParameter) {
		t.Errorf("non-structure replica read settings member: expected the invalid-parameter rejection, got %v", err)
	}

	// The table class member is a string member of the same contract: a
	// present non-string value is rejected at the parse, and a well-typed
	// value passes the parse on to the global-table lookup.
	if err := settings([]interface{}{map[string]interface{}{
		"RegionName":        reqCtx.GetRegion(),
		"ReplicaTableClass": 42.0,
	}}); !errors.Is(err, ErrInvalidParameter) {
		t.Errorf("non-string ReplicaTableClass: expected the invalid-parameter rejection, got %v", err)
	}
	if err := settings([]interface{}{map[string]interface{}{
		"RegionName":        reqCtx.GetRegion(),
		"ReplicaTableClass": "STANDARD",
	}}); !errors.Is(err, ErrGlobalTableNotFound) {
		t.Errorf("well-typed ReplicaTableClass: expected the parse to pass and the lookup to answer, got %v", err)
	}
}
