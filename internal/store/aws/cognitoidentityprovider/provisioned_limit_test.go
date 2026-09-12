package cognitoidentityprovider

import (
	"testing"

	"vorpalstacks/internal/core/storage"
)

// A saved provisioned limit must survive a store reopen (a fresh store over
// the same storage), and provisioned limits are regional: a store for
// another Region must not see the limit.
func TestProvisionedLimitPersistsAcrossReopenAndRegion(t *testing.T) {
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })

	s1 := NewCognitoStore(st, "000000000000", "us-east-1")
	if err := s1.SaveProvisionedLimit(&ProvisionedLimit{Category: "UserAuthentication", Value: 300}); err != nil {
		t.Fatalf("save provisioned limit: %v", err)
	}

	reopened := NewCognitoStore(st, "000000000000", "us-east-1")
	got, err := reopened.GetProvisionedLimit("UserAuthentication")
	if err != nil {
		t.Fatalf("get provisioned limit after reopen: %v", err)
	}
	if got.Category != "UserAuthentication" || got.Value != 300 {
		t.Fatalf("provisioned limit lost across reopen: got %+v", got)
	}
	if got.LastModifiedDate.IsZero() {
		t.Fatal("provisioned limit carries no LastModifiedDate")
	}

	if _, err := reopened.GetProvisionedLimit("UserCreation"); err == nil {
		t.Fatal("unset category returned a provisioned limit")
	}

	west := NewCognitoStore(st, "000000000000", "us-west-2")
	if _, err := west.GetProvisionedLimit("UserAuthentication"); err == nil {
		t.Fatal("provisioned limit leaked across regions")
	}
}
