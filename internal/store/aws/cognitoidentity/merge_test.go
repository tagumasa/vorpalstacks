package cognitoidentity

import (
	"errors"
	"testing"
)

// failingLink always errors, standing in for a developer identity link write
// that fails mid-merge.
func failingLink(di *DeveloperIdentity) error {
	return errors.New("injected link failure")
}

// A failed merge must not leave the source identity destroyed: the developer
// identity link is the authoritative association and must move before any
// identity record is removed, so a failure at any step leaves every developer
// identity pointing at a live identity.
func TestMergeDeveloperIdentitiesFailurePreservesSourceIdentity(t *testing.T) {
	s := newTestStore(t)

	pool := NewIdentityPool("merge-pool", false, "us-east-1")
	if _, err := s.CreateIdentityPool(pool); err != nil {
		t.Fatalf("create pool: %v", err)
	}
	sourceIdentity := NewIdentity(pool.ID)
	sourceIdentity.Logins = map[string]string{"graph.facebook.com": "source-token"}
	if err := s.CreateIdentity(sourceIdentity); err != nil {
		t.Fatalf("create source identity: %v", err)
	}
	destIdentity := NewIdentity(pool.ID)
	if err := s.CreateIdentity(destIdentity); err != nil {
		t.Fatalf("create dest identity: %v", err)
	}
	for _, link := range []*DeveloperIdentity{
		{DeveloperUserIdentifier: "source-user", DeveloperProviderName: "login.example.com", IdentityPoolID: pool.ID, IdentityID: sourceIdentity.ID},
		{DeveloperUserIdentifier: "dest-user", DeveloperProviderName: "login.example.com", IdentityPoolID: pool.ID, IdentityID: destIdentity.ID},
	} {
		if err := s.LinkDeveloperIdentity(link); err != nil {
			t.Fatalf("link developer identity: %v", err)
		}
	}

	if _, err := s.mergeDeveloperIdentities(pool.ID, "login.example.com", "source-user", "dest-user", failingLink); err == nil {
		t.Fatal("expected the injected link failure to surface")
	}

	if _, err := s.GetIdentity(pool.ID, sourceIdentity.ID); err != nil {
		t.Fatalf("source identity was destroyed by the failed merge: %v", err)
	}
}

// A successful merge moves the developer identity link to the destination
// identity, merges the source logins into it and removes the source identity.
func TestMergeDeveloperIdentitiesSuccessMovesLinkAndLogins(t *testing.T) {
	s := newTestStore(t)

	pool := NewIdentityPool("merge-ok-pool", false, "us-east-1")
	if _, err := s.CreateIdentityPool(pool); err != nil {
		t.Fatalf("create pool: %v", err)
	}
	sourceIdentity := NewIdentity(pool.ID)
	sourceIdentity.Logins = map[string]string{"graph.facebook.com": "source-token"}
	if err := s.CreateIdentity(sourceIdentity); err != nil {
		t.Fatalf("create source identity: %v", err)
	}
	destIdentity := NewIdentity(pool.ID)
	if err := s.CreateIdentity(destIdentity); err != nil {
		t.Fatalf("create dest identity: %v", err)
	}
	for _, link := range []*DeveloperIdentity{
		{DeveloperUserIdentifier: "source-user", DeveloperProviderName: "login.example.com", IdentityPoolID: pool.ID, IdentityID: sourceIdentity.ID},
		{DeveloperUserIdentifier: "dest-user", DeveloperProviderName: "login.example.com", IdentityPoolID: pool.ID, IdentityID: destIdentity.ID},
	} {
		if err := s.LinkDeveloperIdentity(link); err != nil {
			t.Fatalf("link developer identity: %v", err)
		}
	}

	destIdentityID, err := s.MergeDeveloperIdentities(pool.ID, "login.example.com", "source-user", "dest-user")
	if err != nil {
		t.Fatalf("MergeDeveloperIdentities: %v", err)
	}
	if destIdentityID != destIdentity.ID {
		t.Fatalf("merge returned identity %s, want %s", destIdentityID, destIdentity.ID)
	}

	sourceLink, err := s.GetDeveloperIdentity(pool.ID, "login.example.com", "source-user")
	if err != nil {
		t.Fatalf("source developer identity vanished: %v", err)
	}
	if sourceLink.IdentityID != destIdentity.ID {
		t.Fatalf("source developer identity still points at %s, want %s", sourceLink.IdentityID, destIdentity.ID)
	}
	if _, err := s.GetIdentity(pool.ID, sourceIdentity.ID); err == nil {
		t.Fatal("source identity still exists after a successful merge")
	}
	merged, err := s.GetIdentity(pool.ID, destIdentity.ID)
	if err != nil {
		t.Fatalf("destination identity vanished: %v", err)
	}
	if merged.Logins["graph.facebook.com"] != "source-token" {
		t.Fatalf("source logins were not merged into the destination identity: %v", merged.Logins)
	}
}
