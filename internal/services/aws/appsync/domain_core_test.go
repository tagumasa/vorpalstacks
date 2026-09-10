package appsync

import (
	"testing"

	"vorpalstacks/internal/core/storage"
	appsyncstore "vorpalstacks/internal/store/aws/appsync"
)

func newDomainCoreTestStore(t *testing.T) *appsyncstore.AppSyncStore {
	t.Helper()
	ps, err := storage.NewPebbleStorage(&storage.Config{
		Path:       t.TempDir(),
		TTLEnabled: false,
	})
	if err != nil {
		t.Fatalf("NewPebbleStorage: %v", err)
	}
	t.Cleanup(func() { ps.Close() })
	return appsyncstore.NewAppSyncStore(ps, "123456789012", "us-east-1")
}

// Domain tags live in the tag store — the same store the tag operations
// serve — so create-time tags are visible to ListTagsForResource and tag
// operations are reflected back in the domain responses.
func TestDomainTagsUseTheTagStore(t *testing.T) {
	store := newDomainCoreTestStore(t)
	svc := &AppSyncService{}

	cfg, tags, err := svc.createDomainNameCore(store, createDomainNameInput{
		DomainName:     "example.com",
		CertificateArn: "arn:aws:acm:us-east-1:123456789012:certificate/abc-def0123456789",
		Tags:           map[string]string{"env": "dev"},
	})
	if err != nil {
		t.Fatalf("createDomainNameCore: %v", err)
	}
	if tags["env"] != "dev" {
		t.Fatalf("create must return the tag-store view, got %v", tags)
	}

	stored, err := store.TagStore.List(cfg.DomainNameArn)
	if err != nil {
		t.Fatalf("TagStore.List: %v", err)
	}
	if stored["env"] != "dev" {
		t.Fatalf("create-time tags must be persisted in the tag store, got %v", stored)
	}

	if err := store.TagStore.Untag(cfg.DomainNameArn, []string{"env"}); err != nil {
		t.Fatalf("TagStore.Untag: %v", err)
	}
	_, tags, err = svc.getDomainNameCore(store, "example.com")
	if err != nil {
		t.Fatalf("getDomainNameCore: %v", err)
	}
	if len(tags) != 0 {
		t.Fatalf("domain responses must reflect the tag store, got %v", tags)
	}
}
