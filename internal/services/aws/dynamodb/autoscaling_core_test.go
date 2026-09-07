package dynamodb

import (
	"testing"

	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// TestParseAutoScalingSettingsTargetValueRange pins the documented
// TargetTrackingScalingPolicyConfiguration.TargetValue metric range on the
// shared auto-scaling settings parser: the member is required and its value
// must fall within the documented Base-10 bounds.
func TestParseAutoScalingSettingsTargetValueRange(t *testing.T) {
	settings := func(target float64) map[string]interface{} {
		return map[string]interface{}{
			"MinimumUnits": float64(5),
			"MaximumUnits": float64(50),
			"ScalingPolicyUpdate": map[string]interface{}{
				"PolicyName": "tracking-policy",
				"TargetTrackingScalingPolicyConfiguration": map[string]interface{}{
					"TargetValue": target,
				},
			},
		}
	}

	for _, valid := range []float64{50, autoScalingTargetValueMin, autoScalingTargetValueMax} {
		if _, err := parseAutoScalingSettings(settings(valid)); err != nil {
			t.Fatalf("TargetValue %v: unexpected error: %v", valid, err)
		}
	}
	for _, invalid := range []float64{0, -1, 1e-200, 2e108} {
		if _, err := parseAutoScalingSettings(settings(invalid)); err == nil {
			t.Fatalf("TargetValue %v: expected rejection", invalid)
		}
	}

	// The member is required: a target tracking configuration without a
	// TargetValue is rejected.
	noTarget := map[string]interface{}{
		"ScalingPolicyUpdate": map[string]interface{}{
			"PolicyName": "tracking-policy",
			"TargetTrackingScalingPolicyConfiguration": map[string]interface{}{
				"DisableScaleIn": true,
			},
		},
	}
	if _, err := parseAutoScalingSettings(noTarget); err == nil {
		t.Fatal("missing TargetValue: expected rejection")
	}
}

// TestParseAutoScalingSettingsTypedMembers pins the typed parse: present
// members land in the description, absent members stay nil, and capacity
// members typed as Long reject fractional or negative values.
func TestParseAutoScalingSettingsTypedMembers(t *testing.T) {
	roleArn := "arn:aws:iam::123456789012:role/auto-scaling"
	desc, err := parseAutoScalingSettings(map[string]interface{}{
		"MinimumUnits":        float64(5),
		"MaximumUnits":        float64(50),
		"AutoScalingDisabled": true,
		"AutoScalingRoleArn":  roleArn,
		"ScalingPolicyUpdate": map[string]interface{}{
			"PolicyName": "tracking-policy",
			"TargetTrackingScalingPolicyConfiguration": map[string]interface{}{
				"TargetValue":      70.0,
				"DisableScaleIn":   true,
				"ScaleInCooldown":  float64(300),
				"ScaleOutCooldown": float64(60),
			},
		},
	})
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if desc.MinimumUnits == nil || *desc.MinimumUnits != 5 {
		t.Fatalf("MinimumUnits = %v, want 5", desc.MinimumUnits)
	}
	if desc.MaximumUnits == nil || *desc.MaximumUnits != 50 {
		t.Fatalf("MaximumUnits = %v, want 50", desc.MaximumUnits)
	}
	if desc.AutoScalingDisabled == nil || !*desc.AutoScalingDisabled {
		t.Fatalf("AutoScalingDisabled = %v, want true", desc.AutoScalingDisabled)
	}
	if desc.AutoScalingRoleArn == nil || *desc.AutoScalingRoleArn != roleArn {
		t.Fatalf("AutoScalingRoleArn = %v", desc.AutoScalingRoleArn)
	}
	if len(desc.ScalingPolicies) != 1 {
		t.Fatalf("ScalingPolicies = %v, want one policy", desc.ScalingPolicies)
	}
	policy := desc.ScalingPolicies[0]
	if policy.PolicyName == nil || *policy.PolicyName != "tracking-policy" {
		t.Fatalf("PolicyName = %v", policy.PolicyName)
	}
	tt := policy.TargetTrackingScalingPolicyConfiguration
	if tt == nil || tt.TargetValue != 70 || tt.DisableScaleIn == nil || !*tt.DisableScaleIn {
		t.Fatalf("target tracking = %+v", tt)
	}
	if tt.ScaleInCooldown == nil || *tt.ScaleInCooldown != 300 || tt.ScaleOutCooldown == nil || *tt.ScaleOutCooldown != 60 {
		t.Fatalf("cooldowns = %v/%v", tt.ScaleInCooldown, tt.ScaleOutCooldown)
	}

	empty, err := parseAutoScalingSettings(map[string]interface{}{})
	if err != nil {
		t.Fatalf("empty parse: %v", err)
	}
	if empty.MinimumUnits != nil || empty.MaximumUnits != nil || empty.AutoScalingDisabled != nil ||
		empty.AutoScalingRoleArn != nil || empty.ScalingPolicies != nil {
		t.Fatalf("absent members must stay nil, got %+v", empty)
	}

	for _, invalid := range []interface{}{float64(1.5), float64(-1), "5"} {
		if _, err := parseAutoScalingSettings(map[string]interface{}{"MinimumUnits": invalid}); err == nil {
			t.Fatalf("MinimumUnits %v: expected rejection", invalid)
		}
	}
}

func replicaDesc(region string) dbstore.ReplicaAutoScalingDescription {
	return dbstore.ReplicaAutoScalingDescription{RegionName: region}
}

func TestMergeReplicaDescriptionsUpsertsPerRegion(t *testing.T) {
	existing := []dbstore.ReplicaAutoScalingDescription{
		replicaDesc("us-east-1"),
		replicaDesc("eu-west-1"),
	}
	// An update for one region and a new region: us-east-1 is replaced,
	// eu-west-1 is preserved untouched, ap-southeast-1 is appended.
	units := int64(1)
	updatedEast := dbstore.ReplicaAutoScalingDescription{
		RegionName: "us-east-1",
		Read:       &dbstore.AutoScalingSettingsDescription{MinimumUnits: &units},
	}
	merged := mergeReplicaDescriptions(existing, []dbstore.ReplicaAutoScalingDescription{updatedEast, replicaDesc("ap-southeast-1")})

	if len(merged) != 3 {
		t.Fatalf("merged = %v, want 3 replicas", merged)
	}
	if merged[0].Read == nil {
		t.Fatalf("us-east-1 must carry the updated settings, got %v", merged[0])
	}
	if merged[1].RegionName != "eu-west-1" {
		t.Fatalf("unmentioned region must be preserved in place, got %v", merged[1])
	}
	if merged[2].RegionName != "ap-southeast-1" {
		t.Fatalf("new region must be appended, got %v", merged[2])
	}

	// An empty update leaves the stored list unchanged — an update without
	// ReplicaUpdates must not erase the stored replicas.
	same := mergeReplicaDescriptions(existing, nil)
	if len(same) != 2 || same[0].RegionName != "us-east-1" || same[1].RegionName != "eu-west-1" {
		t.Fatalf("empty update must preserve the stored replicas, got %v", same)
	}

	// The input list is never mutated.
	if existing[0].Read != nil {
		t.Fatal("merge must not mutate the stored descriptions")
	}
}

func TestReplicaDescriptionsFromSettingsNeverNil(t *testing.T) {
	if got := replicaDescriptionsFromSettings(nil); got == nil || len(got) != 0 {
		t.Fatalf("nil settings must yield an empty non-nil list, got %v", got)
	}
	got := replicaDescriptionsFromSettings(&dbstore.TableReplicaAutoScalingSettings{
		Replicas: []dbstore.ReplicaAutoScalingDescription{replicaDesc("us-east-1")},
	})
	if len(got) != 1 || got[0].RegionName != "us-east-1" {
		t.Fatalf("stored replicas must pass through, got %v", got)
	}
}

// TestAutoScalingSettingsWireRoundTrip pins the wire rendering: only set
// members appear, and the render covers the full member set.
func TestAutoScalingSettingsWireRoundTrip(t *testing.T) {
	if got := autoScalingSettingsToWire(nil); got != nil {
		t.Fatalf("nil description must render nil, got %v", got)
	}
	units := int64(5)
	wire := autoScalingSettingsToWire(&dbstore.AutoScalingSettingsDescription{
		MinimumUnits: &units,
		ScalingPolicies: []dbstore.AutoScalingPolicyDescription{{
			TargetTrackingScalingPolicyConfiguration: &dbstore.TargetTrackingScalingPolicyConfiguration{TargetValue: 70},
		}},
	})
	if wire["MinimumUnits"] != int64(5) {
		t.Fatalf("MinimumUnits = %v", wire["MinimumUnits"])
	}
	if _, ok := wire["MaximumUnits"]; ok {
		t.Fatal("absent members must not render")
	}
	policies, ok := wire["ScalingPolicies"].([]interface{})
	if !ok || len(policies) != 1 {
		t.Fatalf("ScalingPolicies = %v", wire["ScalingPolicies"])
	}
	policy := policies[0].(map[string]interface{})
	if _, ok := policy["PolicyName"]; ok {
		t.Fatal("absent PolicyName must not render")
	}
	tt, ok := policy["TargetTrackingScalingPolicyConfiguration"].(map[string]interface{})
	if !ok || tt["TargetValue"] != float64(70) {
		t.Fatalf("target tracking = %v", policy)
	}
}
