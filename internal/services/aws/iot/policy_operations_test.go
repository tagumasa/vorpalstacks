package iot

import "testing"

func TestPolicyNamePattern(t *testing.T) {
	valid := []string{"MyPolicy", "policy-1", "policy_name", "policy.name", "p@-+=,.x"}
	for _, name := range valid {
		if !policyNamePattern.MatchString(name) {
			t.Errorf("valid name %q rejected by pattern", name)
		}
	}

	invalid := []string{"INVALID NAME", "name with spaces", "name/with/slashes", "name#with#hash", ""}
	for _, name := range invalid {
		if policyNamePattern.MatchString(name) {
			t.Errorf("invalid name %q accepted by pattern", name)
		}
	}
}
