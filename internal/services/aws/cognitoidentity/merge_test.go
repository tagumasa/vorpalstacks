package cognitoidentity

import (
	"context"
	"errors"
	"sync"
	"testing"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/storage"
	cognitoidentitystore "vorpalstacks/internal/store/aws/cognitoidentity"
)

func newMergeTestService(t *testing.T) (*CognitoIdentityService, *cognitoidentitystore.CognitoIdentityStore) {
	t.Helper()
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	real := cognitoidentitystore.NewCognitoIdentityStore(st, "000000000000", "us-east-1")
	svc := NewCognitoIdentityService("000000000000", "us-east-1")
	svc.stores.Store("us-east-1", real)
	return svc, real
}

type mergeFixture struct {
	poolID        string
	provider      string
	sourceUser    string
	destUser      string
	sourceIdentID string
	destIdentID   string
}

func seedMergeFixture(t *testing.T, store *cognitoidentitystore.CognitoIdentityStore) mergeFixture {
	t.Helper()
	pool := cognitoidentitystore.NewIdentityPool("merge-pool", false, "us-east-1")
	if _, err := store.CreateIdentityPool(pool); err != nil {
		t.Fatalf("create pool: %v", err)
	}
	sourceIdentity := cognitoidentitystore.NewIdentity(pool.ID)
	sourceIdentity.Logins = map[string]string{"graph.facebook.com": "source-token"}
	if err := store.CreateIdentity(sourceIdentity); err != nil {
		t.Fatalf("create source identity: %v", err)
	}
	destIdentity := cognitoidentitystore.NewIdentity(pool.ID)
	if err := store.CreateIdentity(destIdentity); err != nil {
		t.Fatalf("create dest identity: %v", err)
	}
	fixture := mergeFixture{
		poolID:        pool.ID,
		provider:      "login.example.com",
		sourceUser:    "source-user",
		destUser:      "dest-user",
		sourceIdentID: sourceIdentity.ID,
		destIdentID:   destIdentity.ID,
	}
	for _, link := range []*cognitoidentitystore.DeveloperIdentity{
		{DeveloperUserIdentifier: fixture.sourceUser, DeveloperProviderName: fixture.provider, IdentityPoolID: pool.ID, IdentityID: sourceIdentity.ID},
		{DeveloperUserIdentifier: fixture.destUser, DeveloperProviderName: fixture.provider, IdentityPoolID: pool.ID, IdentityID: destIdentity.ID},
	} {
		if err := store.LinkDeveloperIdentity(link); err != nil {
			t.Fatalf("link developer identity: %v", err)
		}
	}
	return fixture
}

// Concurrent token requests for one developer user must resolve to a single
// identity: the first caller creates it and the remaining callers reuse it.
func TestGetOpenIdTokenForDeveloperIdentityConcurrentSingleIdentity(t *testing.T) {
	svc, real := newMergeTestService(t)
	pool := cognitoidentitystore.NewIdentityPool("dev-pool", false, "us-east-1")
	if _, err := real.CreateIdentityPool(pool); err != nil {
		t.Fatalf("create pool: %v", err)
	}

	const goroutines = 32
	var wg sync.WaitGroup
	start := make(chan struct{})
	identityIDs := make([]string, goroutines)
	reqErrs := make([]error, goroutines)
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-start
			resp, err := svc.GetOpenIdTokenForDeveloperIdentity(context.Background(),
				&request.RequestContext{Region: "us-east-1"},
				&request.ParsedRequest{Parameters: map[string]interface{}{
					"IdentityPoolId": pool.ID,
					"Logins":         map[string]interface{}{"login.example.com": "user-1"},
				}})
			reqErrs[i] = err
			if err == nil {
				mapped, ok := resp.(map[string]interface{})
				if !ok {
					reqErrs[i] = errors.New("unexpected response shape")
					return
				}
				id, _ := mapped["IdentityId"].(string)
				identityIDs[i] = id
			}
		}(i)
	}
	close(start)
	wg.Wait()

	for i, err := range reqErrs {
		if err != nil {
			t.Fatalf("concurrent request %d failed: %v", i, err)
		}
	}
	for i, id := range identityIDs {
		if id != identityIDs[0] {
			t.Fatalf("concurrent request %d resolved identity %s, want the single identity %s", i, id, identityIDs[0])
		}
	}
	identities, _, err := real.ListIdentitiesByPool(pool.ID, 60, "")
	if err != nil {
		t.Fatalf("list identities: %v", err)
	}
	if len(identities) != 1 {
		t.Fatalf("pool holds %d identities after concurrent requests, want 1", len(identities))
	}
}

// A successful merge moves the developer identity link to the destination
// identity, merges the source logins into it and removes the source identity.
func TestMergeDeveloperIdentitiesSuccessMovesLinkAndLogins(t *testing.T) {
	svc, real := newMergeTestService(t)
	fixture := seedMergeFixture(t, real)

	resp, err := svc.MergeDeveloperIdentities(context.Background(),
		&request.RequestContext{Region: "us-east-1"},
		&request.ParsedRequest{Parameters: map[string]interface{}{
			"IdentityPoolId":            fixture.poolID,
			"DeveloperProviderName":     fixture.provider,
			"SourceUserIdentifier":      fixture.sourceUser,
			"DestinationUserIdentifier": fixture.destUser,
		}})
	if err != nil {
		t.Fatalf("MergeDeveloperIdentities: %v", err)
	}
	mapped, ok := resp.(map[string]interface{})
	if !ok {
		t.Fatal("unexpected response shape")
	}
	if mapped["IdentityId"] != fixture.destIdentID {
		t.Fatalf("merge returned identity %v, want %s", mapped["IdentityId"], fixture.destIdentID)
	}

	sourceLink, err := real.GetDeveloperIdentity(fixture.poolID, fixture.provider, fixture.sourceUser)
	if err != nil {
		t.Fatalf("source developer identity vanished: %v", err)
	}
	if sourceLink.IdentityID != fixture.destIdentID {
		t.Fatalf("source developer identity still points at %s, want %s", sourceLink.IdentityID, fixture.destIdentID)
	}

	if _, err := real.GetIdentity(fixture.poolID, fixture.sourceIdentID); err == nil {
		t.Fatal("source identity still exists after a successful merge")
	}
	destIdentity, err := real.GetIdentity(fixture.poolID, fixture.destIdentID)
	if err != nil {
		t.Fatalf("destination identity vanished: %v", err)
	}
	if destIdentity.Logins["graph.facebook.com"] != "source-token" {
		t.Fatalf("source logins were not merged into the destination identity: %v", destIdentity.Logins)
	}
}

// The lookup response carries the model's members only — the developer user
// identifiers plus the optional identity ID and page token.
func TestLookupDeveloperIdentityResultCarriesOnlyModelMembers(t *testing.T) {
	result := lookupDeveloperIdentityResult("us-east-1:abc", []string{"user-1"}, "")
	if _, hasPoolID := result["IdentityPoolId"]; hasPoolID {
		t.Fatal("response carries the non-model IdentityPoolId member")
	}
	if result["IdentityId"] != "us-east-1:abc" {
		t.Fatalf("IdentityId mismatch: %v", result["IdentityId"])
	}

	bare := lookupDeveloperIdentityResult("", []string{"user-1"}, "next")
	if _, hasID := bare["IdentityId"]; hasID {
		t.Fatal("empty match must omit the IdentityId member")
	}
	if bare["NextToken"] != "next" {
		t.Fatalf("NextToken mismatch: %v", bare["NextToken"])
	}
}
