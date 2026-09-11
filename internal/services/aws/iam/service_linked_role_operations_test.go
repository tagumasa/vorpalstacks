package iam

import (
	"errors"
	"net/http"
	"testing"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/core/storage"
	iamstore "vorpalstacks/internal/store/aws/iam"
)

// The failure Reason must follow the modelled DeletionTaskFailureReasonType
// structure: a string Reason member (plus RoleUsageList, unpopulated — no
// usage tracking exists on this platform). The earlier invented
// ReasonCode/ReasonMessage members match no model shape.
func TestSLRoleDeletionStatusResponseReason(t *testing.T) {
	failed := deletionStatusToResponse(&iamstore.SLRoleDeletionTask{
		Status:         "FAILED",
		DeletionFailed: true,
		ErrorReason:    "role in use",
	})

	if got := failed["Status"]; got != "FAILED" {
		t.Errorf("Status: got %v, want FAILED", got)
	}
	reason, ok := failed["Reason"].(map[string]interface{})
	if !ok {
		t.Fatalf("Reason must serialise as the DeletionTaskFailureReasonType structure, got %T", failed["Reason"])
	}
	if got := reason["Reason"]; got != "role in use" {
		t.Errorf("Reason.Reason: got %v, want the task's ErrorReason", got)
	}
	for _, key := range []string{"ReasonCode", "ReasonMessage"} {
		if _, present := reason[key]; present {
			t.Errorf("Reason must not carry the invented %s member", key)
		}
	}

	success := deletionStatusToResponse(&iamstore.SLRoleDeletionTask{
		Status: "SUCCEEDED",
	})
	if _, present := success["Reason"]; present {
		t.Error("a SUCCEEDED task must not carry a Reason member")
	}
}

// A failed role read during DeleteServiceLinkedRole must surface as an
// internal failure carrying the cause — the role may exist while the read
// itself failed — and a genuinely absent role still answers NoSuchRole
// through the resolved Get, with no Exists pre-probe doubling the read.
func TestDeleteServiceLinkedRoleCoreReadFailureSemantics(t *testing.T) {
	s, store := faultTestStore(t, "iam_roles", &readFault{
		failGet: true,
		err:     errors.New("simulated role read failure"),
	})

	_, err := s.deleteServiceLinkedRoleCore(store, "AWSServiceRoleForECS")
	assertInternalFailure(t, err, "simulated role read failure")

	plainSt, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	defer plainSt.Close()
	plainStore := iamstore.NewIAMStore(plainSt, "123456789012")

	_, err = s.deleteServiceLinkedRoleCore(plainStore, "ghost-role")
	if err == nil {
		t.Fatal("expected NoSuchRole for a missing role")
	}
	awsErr, ok := err.(*awserrors.AWSError)
	if !ok {
		t.Fatalf("expected *awserrors.AWSError, got %T", err)
	}
	if awsErr.GetHTTPStatusCode() != http.StatusNotFound {
		t.Fatalf("got status %d, want 404 for a missing role", awsErr.GetHTTPStatusCode())
	}
}
