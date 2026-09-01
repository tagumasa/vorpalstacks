package cloudtrail

import "testing"

// TestDeleteTrailCoreRejectsEmptyName pins the shared empty-name
// rejection: it fires before any store access, so both the HTTP API and
// the admin console answer InvalidParameterException for an omitted
// trail name instead of a not-found.
func TestDeleteTrailCoreRejectsEmptyName(t *testing.T) {
	svc := &CloudTrailService{}
	if err := svc.deleteTrailCore(nil, DeleteTrailInput{NameOrARN: ""}); err != ErrInvalidParameter {
		t.Fatalf("expected ErrInvalidParameter for an empty name, got %v", err)
	}
}
