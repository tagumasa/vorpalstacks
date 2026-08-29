package dynamodb

import "testing"

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
