package kms

import (
	"context"
	"testing"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/storage"
	kmsstore "vorpalstacks/internal/store/aws/kms"
)

// TestResolveKeyAcceptsAliasARNForm pins the resolver's routing of the
// documented alias-ARN key-identifier form: a full alias ARN routes
// through the alias resolution (which already normalises ARNs to alias
// names) and settles to the target key, on both the service's own
// resolution path and the key resolver the consumer services wire — a
// prefix-only alias gate would route the ARN to the key lookup and answer
// KeyNotFound. A key ARN and a bare key id keep resolving through the key
// lookup on the same store.
func TestResolveKeyAcceptsAliasARNForm(t *testing.T) {
	sm, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatalf("storage manager: %v", err)
	}
	t.Cleanup(func() { sm.Close() })
	svc := NewKMSService("123456789012", "us-east-1", nil)
	svc.SetStorageManager(sm)
	reqCtx := request.NewRequestContext(context.Background(), sm, "123456789012", "us-east-1")

	stores, err := svc.store(reqCtx)
	if err != nil {
		t.Fatalf("kms stores: %v", err)
	}
	key, err := stores.keys.Create("1234abcd-12ab-34cd-56ef-1234567890ab", kmsstore.KeyUsageEncryptDecrypt, kmsstore.KeySpecSymmetricDefault, "resolve pin", kmsstore.OriginTypeAWSKMS, false)
	if err != nil {
		t.Fatalf("create key: %v", err)
	}
	if _, err := stores.aliases.Create("alias/resolve-pin", key.KeyID); err != nil {
		t.Fatalf("create alias: %v", err)
	}

	aliasArn := "arn:aws:kms:us-east-1:123456789012:alias/resolve-pin"

	// The service's own resolution path: the alias ARN settles to the key.
	resolved, err := svc.resolveKey(stores, map[string]interface{}{"KeyId": aliasArn})
	if err != nil {
		t.Fatalf("resolveKey with an alias ARN: %v", err)
	}
	if resolved.KeyID != key.KeyID {
		t.Fatalf("resolveKey with an alias ARN = key %q, want %q", resolved.KeyID, key.KeyID)
	}

	// The wired consumer face: the resolver answers the target key's ARN,
	// which is what the consumer services' SSE key members store.
	resolver := svc.NewKeyResolver()
	gotArn, err := resolver.ResolveKeyArn(context.Background(), "us-east-1", aliasArn)
	if err != nil {
		t.Fatalf("ResolveKeyArn with an alias ARN: %v", err)
	}
	if gotArn != key.Arn {
		t.Fatalf("ResolveKeyArn with an alias ARN = %q, want %q", gotArn, key.Arn)
	}

	// The alias-name form answers the same target on the wired face.
	gotArn, err = resolver.ResolveKeyArn(context.Background(), "us-east-1", "alias/resolve-pin")
	if err != nil || gotArn != key.Arn {
		t.Fatalf("ResolveKeyArn with an alias name = (%q, %v), want (%q, nil)", gotArn, err, key.Arn)
	}

	// The non-alias identifier forms keep routing through the key lookup.
	gotArn, err = resolver.ResolveKeyArn(context.Background(), "us-east-1", key.Arn)
	if err != nil || gotArn != key.Arn {
		t.Fatalf("ResolveKeyArn with a key ARN = (%q, %v), want (%q, nil)", gotArn, err, key.Arn)
	}
}
