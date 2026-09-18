package sts

import (
	"os"
	"testing"

	"vorpalstacks/internal/core/storage"
)

func newSessionTestStore(t *testing.T) *SessionStore {
	t.Helper()
	tmpDir := "./tmp/sts-session-test"
	os.RemoveAll(tmpDir)
	t.Cleanup(func() { os.RemoveAll(tmpDir) })

	s, err := storage.Open(tmpDir)
	if err != nil {
		t.Fatalf("failed to open test storage: %v", err)
	}
	t.Cleanup(func() { s.Close() })

	return NewSessionStore(s, "us-east-1")
}

// TestCreateRecordsIssuanceTime pins the session's CreatedAt: it is set at
// creation and carried through ResolveSession, because the CloudTrail
// sessionContext of a temporary credential reports it as the session's
// creationDate.
func TestCreateRecordsIssuanceTime(t *testing.T) {
	store := newSessionTestStore(t)

	session, err := store.Create(CreateSessionParams{
		PrincipalType:   "AssumedRole",
		PrincipalName:   "alice",
		PrincipalArn:    "arn:aws:sts::123456789012:assumed-role/DevRole/Dev1",
		RoleArn:         "arn:aws:iam::123456789012:role/DevRole",
		RoleSessionName: "Dev1",
		DurationSeconds: 3600,
	})
	if err != nil {
		t.Fatalf("create session: %v", err)
	}
	if session.CreatedAt.IsZero() {
		t.Fatal("created session carries no issuance time")
	}

	creds, err := store.ResolveSession(session.AccessKeyId)
	if err != nil {
		t.Fatalf("resolve session: %v", err)
	}
	if creds.CreatedAt.IsZero() {
		t.Error("resolved credentials carry no issuance time")
	}
	if creds.RoleArn != "arn:aws:iam::123456789012:role/DevRole" {
		t.Errorf("resolved RoleArn = %q, want the assumed role's ARN", creds.RoleArn)
	}
	if creds.RoleSessionName != "Dev1" {
		t.Errorf("resolved RoleSessionName = %q, want Dev1", creds.RoleSessionName)
	}
}
