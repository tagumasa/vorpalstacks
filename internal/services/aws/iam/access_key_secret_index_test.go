package iam

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"vorpalstacks/internal/core/storage"
	iamstore "vorpalstacks/internal/store/aws/iam"
)

// The secret→accessKeyId index serves GetBySecretKey without the
// full-bucket scan, while the record set stays authoritative: a lost
// index entry is healed from the scan backstop, an unknown secret keeps
// the not-found semantics, and Delete removes the index entry together
// with the record. The bucket name below mirrors the unexported store
// constant, duplicated here because the store keeps it unexported.
func TestGetBySecretKeyIndexServeRepairAndDelete(t *testing.T) {
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	t.Cleanup(func() { st.Close() })
	store := iamstore.NewIAMStore(st, "123456789012")

	if _, err := store.Users().Create("index-user", "/", "123456789012", nil); err != nil {
		t.Fatalf("create user: %v", err)
	}
	key, err := store.AccessKeys().Create("index-user")
	if err != nil {
		t.Fatalf("create access key: %v", err)
	}

	// Fast path: the created key is served through the index.
	found, err := store.AccessKeys().GetBySecretKey(key.SecretAccessKey)
	if err != nil {
		t.Fatalf("get by secret key: %v", err)
	}
	assert.Equal(t, key.AccessKeyId, found.AccessKeyId)

	// Backstop: with the index bucket emptied, the scan still finds the
	// record and repairs the entry.
	indexBucket := st.Bucket("iam_access_key_secret_index")
	var staleKeys [][]byte
	if err := indexBucket.ForEach(func(k, v []byte) error {
		staleKeys = append(staleKeys, k)
		return nil
	}); err != nil {
		t.Fatalf("scan secret index: %v", err)
	}
	for _, k := range staleKeys {
		if err := indexBucket.Delete(k); err != nil {
			t.Fatalf("empty secret index: %v", err)
		}
	}

	found, err = store.AccessKeys().GetBySecretKey(key.SecretAccessKey)
	if err != nil {
		t.Fatalf("get by secret key after index loss: %v", err)
	}
	assert.Equal(t, key.AccessKeyId, found.AccessKeyId)
	assert.Equal(t, 1, indexBucket.Count(), "the scan backstop must repair the index entry")

	// Unknown secrets keep the not-found semantics.
	if _, err := store.AccessKeys().GetBySecretKey("no-such-secret"); err == nil {
		t.Fatal("expected a not-found error for an unknown secret")
	} else {
		assert.Contains(t, err.Error(), "not found")
	}

	// Delete removes the record together with its index entry.
	if err := store.AccessKeys().Delete(key.AccessKeyId); err != nil {
		t.Fatalf("delete access key: %v", err)
	}
	assert.Equal(t, 0, indexBucket.Count(), "delete must remove the secret index entry")
	if _, err := store.AccessKeys().GetBySecretKey(key.SecretAccessKey); err == nil {
		t.Fatal("expected the deleted key to be unreachable by secret")
	}
}
