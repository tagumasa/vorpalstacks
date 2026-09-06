package apigateway

import (
	"fmt"
	"testing"

	"vorpalstacks/internal/core/storage"
)

// TestRemoveBasePathMappingsForApiAllPages pins the dangling-reference sweep:
// mappings on domains beyond the first listing page must be removed too. The
// domain count (520) exceeds both the ForEachAll page size (100) and the
// fixed 500-item listing the previous implementation used, so a regression to
// any single-page walk leaves mappings behind and fails this test.
func TestRemoveBasePathMappingsForApiAllPages(t *testing.T) {
	ps, err := storage.NewPebbleStorage(&storage.Config{
		Path:       t.TempDir(),
		TTLEnabled: false,
	})
	if err != nil {
		t.Fatalf("NewPebbleStorage: %v", err)
	}
	t.Cleanup(func() { ps.Close() })

	store := NewDomainStore(ps, "123456789012", "us-east-1")

	const domainCount = 520
	for i := 0; i < domainCount; i++ {
		domain := fmt.Sprintf("sweep-%04d.example.com", i)
		if _, err := store.CreateDomainName(&DomainName{DomainName: domain}); err != nil {
			t.Fatalf("CreateDomainName(%s): %v", domain, err)
		}
		if _, err := store.CreateBasePathMapping(domain, &BasePathMapping{
			BasePath:  "v1",
			RestApiId: "apiBeingDeleted",
			Stage:     "prod",
		}); err != nil {
			t.Fatalf("CreateBasePathMapping(%s): %v", domain, err)
		}
	}
	// A mapping for a different API must survive the sweep.
	if _, err := store.CreateBasePathMapping("sweep-0000.example.com", &BasePathMapping{
		BasePath:  "keep",
		RestApiId: "apiOther",
		Stage:     "prod",
	}); err != nil {
		t.Fatalf("CreateBasePathMapping(keep): %v", err)
	}

	if err := store.RemoveBasePathMappingsForApi("apiBeingDeleted"); err != nil {
		t.Fatalf("RemoveBasePathMappingsForApi: %v", err)
	}

	for i := 0; i < domainCount; i++ {
		domain := fmt.Sprintf("sweep-%04d.example.com", i)
		if _, err := store.GetBasePathMapping(domain, "v1"); err == nil {
			t.Fatalf("mapping for deleted API still present on %s", domain)
		}
	}
	if _, err := store.GetBasePathMapping("sweep-0000.example.com", "keep"); err != nil {
		t.Fatalf("mapping for other API must survive, got: %v", err)
	}
}
