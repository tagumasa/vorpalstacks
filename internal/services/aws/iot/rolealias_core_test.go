package iot

import "testing"

// TestRoleAliasCoreRejectsDurationRange pins the CredentialDurationSeconds
// bounds on both role-alias write paths: out-of-range durations are
// rejected before any store access (a nil store stands in for the real
// one; the guards run first).
func TestRoleAliasCoreRejectsDurationRange(t *testing.T) {
	svc := &IoTService{}

	createCases := []int64{MinRoleAliasCredentialDuration - 1, MaxRoleAliasCredentialDuration + 1}
	for _, duration := range createCases {
		if _, err := svc.createRoleAliasCore(nil, CreateRoleAliasInput{
			RoleAlias:                 "alias",
			RoleARN:                   "arn:aws:iam::123456789012:role/alias-role",
			CredentialDurationSeconds: duration,
			DurationProvided:          true,
		}); err == nil {
			t.Fatalf("createRoleAliasCore accepted out-of-range duration %d", duration)
		}
	}

	updateCases := []int64{0, MinRoleAliasCredentialDuration - 1, MaxRoleAliasCredentialDuration + 1}
	for _, duration := range updateCases {
		if _, err := svc.updateRoleAliasCore(nil, UpdateRoleAliasInput{
			RoleAlias:                 "alias",
			CredentialDurationSeconds: duration,
			DurationProvided:          true,
		}); err == nil {
			t.Fatalf("updateRoleAliasCore accepted out-of-range duration %d", duration)
		}
	}

	// An absent duration must not be rejected on either path (the default
	// applies); the call then proceeds to store access, so only the
	// missing-member guard is observable with a nil store.
	if _, err := svc.updateRoleAliasCore(nil, UpdateRoleAliasInput{}); err == nil {
		t.Fatal("updateRoleAliasCore accepted an empty role alias")
	}
}
