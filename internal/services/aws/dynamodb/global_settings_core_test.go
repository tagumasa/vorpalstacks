package dynamodb

import (
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
}
