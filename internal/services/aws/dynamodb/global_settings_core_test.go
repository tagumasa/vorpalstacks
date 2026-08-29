package dynamodb

import (
	"reflect"
	"testing"
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
	policies, ok := desc["ScalingPolicies"].([]interface{})
	if !ok || len(policies) != 1 {
		t.Fatalf("expected one scaling policy, got %#v", desc["ScalingPolicies"])
	}
	policy, ok := policies[0].(map[string]interface{})
	if !ok {
		t.Fatalf("expected a policy map, got %#v", policies[0])
	}
	if policy["PolicyName"] != "write-tracking" {
		t.Errorf("policy name = %#v", policy["PolicyName"])
	}
	want := map[string]interface{}{
		"TargetValue":      50.0,
		"DisableScaleIn":   true,
		"ScaleInCooldown":  60.0,
		"ScaleOutCooldown": 30.0,
	}
	tt, ok := policy["TargetTrackingScalingPolicyConfiguration"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected a target tracking map, got %#v", policy["TargetTrackingScalingPolicyConfiguration"])
	}
	if !reflect.DeepEqual(tt, want) {
		t.Errorf("target tracking = %#v, want %#v", tt, want)
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
}
