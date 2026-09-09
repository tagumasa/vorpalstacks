package s3

import (
	"context"
	"errors"
	"strings"
	"testing"
)

// Suspension stops new versions but rewrites nothing: the layout-based
// list and count filters must still yield one entry per key, the
// pre-suspension versions must stay suppressed, and the count must reach
// zero once every record is deleted — the state DeleteBucket checks.
func TestSuspendedVersioningListsAndCountsOncePerKey(t *testing.T) {
	store, buckets := newMultipartTestStore(t)
	ctx := context.Background()

	if _, err := buckets.Create("vb", "test-region"); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if err := buckets.SetVersioning("vb", BucketVersioningEnabled, ""); err != nil {
		t.Fatalf("SetVersioning enabled: %v", err)
	}

	v1, err := store.PutWithVersioning(ctx, "vb", "k", strings.NewReader("one"), "", nil, false, StorageClassStandard, nil)
	if err != nil {
		t.Fatalf("put v1: %v", err)
	}
	if _, err := store.PutWithVersioning(ctx, "vb", "k", strings.NewReader("two"), "", nil, false, StorageClassStandard, nil); err != nil {
		t.Fatalf("put v2: %v", err)
	}

	if err := buckets.SetVersioning("vb", BucketVersioningSuspended, ""); err != nil {
		t.Fatalf("SetVersioning suspended: %v", err)
	}

	// One entry for the twice-versioned key, not two versions plus a
	// pointer copy.
	list, err := store.List("vb", "", "", "", 1000)
	if err != nil {
		t.Fatalf("List after suspend: %v", err)
	}
	if len(list.Objects) != 1 || list.Objects[0].Key != "k" {
		t.Fatalf("list after suspend = %d objects, want exactly [k]", len(list.Objects))
	}
	if n, err := store.CountByBucket("vb"); err != nil || n != 1 {
		t.Fatalf("count after suspend = %d, %v; want 1", n, err)
	}

	// A new write in the suspended bucket lands as a single null record
	// and counts once.
	if _, err := store.PutWithVersioning(ctx, "vb", "k2", strings.NewReader("suspended write"), "", nil, false, StorageClassStandard, nil); err != nil {
		t.Fatalf("put suspended: %v", err)
	}
	if n, err := store.CountByBucket("vb"); err != nil || n != 2 {
		t.Fatalf("count after suspended write = %d, %v; want 2", n, err)
	}

	// Deleting every record brings the count to zero — the state whose
	// absence used to block bucket deletion forever.
	if _, err := store.DeleteWithVersion(ctx, "vb", "k", v1.VersionID); err != nil {
		// "k" holds two versions; deleting the resolved current one is
		// enough for the count only if it is the latest — delete by the
		// head's version id.
		t.Fatalf("DeleteWithVersion(k): %v", err)
	}
	if obj, err := store.Head(ctx, "vb", "k"); err == nil {
		if _, err := store.DeleteWithVersion(ctx, "vb", "k", obj.VersionID); err != nil {
			t.Fatalf("DeleteWithVersion(k remaining): %v", err)
		}
	}
	if marker, err := store.DeleteWithVersion(ctx, "vb", "k2", ""); err != nil || marker == nil || !marker.IsDeleteMarker {
		t.Fatalf("suspended delete = %+v, %v; want a delete marker as the new current version", marker, err)
	}
	if n, err := store.CountByBucket("vb"); err != nil || n != 0 {
		t.Fatalf("count after deletes = %d, %v; want 0", n, err)
	}
	_ = v1
}

// An object written before versioning was enabled keeps a single null
// record; the first versioned overwrite must flip it to non-latest so it
// does not list alongside the new version.
func TestLegacyNullRecordFlippedOnFirstVersionedWrite(t *testing.T) {
	store, buckets := newMultipartTestStore(t)
	ctx := context.Background()

	if _, err := buckets.Create("vb2", "test-region"); err != nil {
		t.Fatalf("Create: %v", err)
	}

	if _, err := store.PutWithVersioning(ctx, "vb2", "legacy", strings.NewReader("pre-versioning"), "", nil, false, StorageClassStandard, nil); err != nil {
		t.Fatalf("legacy put: %v", err)
	}
	if err := buckets.SetVersioning("vb2", BucketVersioningEnabled, ""); err != nil {
		t.Fatalf("SetVersioning enabled: %v", err)
	}
	if _, err := store.PutWithVersioning(ctx, "vb2", "legacy", strings.NewReader("versioned"), "", nil, false, StorageClassStandard, nil); err != nil {
		t.Fatalf("versioned put: %v", err)
	}

	list, err := store.List("vb2", "", "", "", 1000)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(list.Objects) != 1 {
		t.Fatalf("list = %d objects, want 1 (legacy null must not double with the new version)", len(list.Objects))
	}
	if n, err := store.CountByBucket("vb2"); err != nil || n != 1 {
		t.Fatalf("count = %d, %v; want 1", n, err)
	}
}

// The mutation skeleton resolves by layout, not by the bucket's current
// status: a suspended bucket with versioned-only keys still mutates the
// current version (the status-driven copies used to fail with not-found),
// and a suspended null record never resurrects the deleted "_latest"
// pointer.
func TestMutateObjectRecordByLayoutOnSuspendedBucket(t *testing.T) {
	store, buckets := newMultipartTestStore(t)
	ctx := context.Background()

	if _, err := buckets.Create("vb3", "test-region"); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if err := buckets.SetVersioning("vb3", BucketVersioningEnabled, ""); err != nil {
		t.Fatalf("SetVersioning enabled: %v", err)
	}
	if _, err := store.PutWithVersioning(ctx, "vb3", "k", strings.NewReader("v1"), "", nil, false, StorageClassStandard, nil); err != nil {
		t.Fatalf("put v1: %v", err)
	}
	v2, err := store.PutWithVersioning(ctx, "vb3", "k", strings.NewReader("v2"), "", nil, false, StorageClassStandard, nil)
	if err != nil {
		t.Fatalf("put v2: %v", err)
	}

	if err := buckets.SetVersioning("vb3", BucketVersioningSuspended, ""); err != nil {
		t.Fatalf("SetVersioning suspended: %v", err)
	}

	if err := store.SetStorageClass("vb3", "k", "", StorageClassDeepArchive); err != nil {
		t.Fatalf("SetStorageClass on suspended versioned-only key: %v", err)
	}
	head, err := store.Head(ctx, "vb3", "k")
	if err != nil {
		t.Fatalf("Head: %v", err)
	}
	if head.StorageClass != StorageClassDeepArchive {
		t.Fatalf("current version storage class = %q, want DEEP_ARCHIVE", head.StorageClass)
	}
	if head.VersionID != v2.VersionID {
		t.Fatalf("mutated version = %q, want the current %q", head.VersionID, v2.VersionID)
	}

	// A suspended write replaces the pointer with its null record; mutating
	// that record must not bring the pointer back.
	if _, err := store.PutWithVersioning(ctx, "vb3", "k2", strings.NewReader("suspended"), "", nil, false, StorageClassStandard, nil); err != nil {
		t.Fatalf("suspended put: %v", err)
	}
	if err := store.SetStorageClass("vb3", "k2", "", StorageClassGlacierIR); err != nil {
		t.Fatalf("SetStorageClass on suspended null record: %v", err)
	}
	if store.BaseStore.Exists(store.latestKeyStorageKey("vb3", "k2")) {
		t.Fatal("mutation resurrected the _latest pointer of a suspended key")
	}
	head2, err := store.Head(ctx, "vb3", "k2")
	if err != nil {
		t.Fatalf("Head k2: %v", err)
	}
	if head2.StorageClass != StorageClassGlacierIR {
		t.Fatalf("suspended record storage class = %q, want GLACIER_IR", head2.StorageClass)
	}
}

// Delete markers are not readable objects: metadata mutations and ACL reads
// on a marker version surface ErrObjectNotFound instead of silently
// mutating or returning the marker's fields.
func TestMutateObjectRecordRejectsDeleteMarker(t *testing.T) {
	store, buckets := newMultipartTestStore(t)
	ctx := context.Background()

	if _, err := buckets.Create("vb4", "test-region"); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if err := buckets.SetVersioning("vb4", BucketVersioningEnabled, ""); err != nil {
		t.Fatalf("SetVersioning enabled: %v", err)
	}
	if _, err := store.PutWithVersioning(ctx, "vb4", "k", strings.NewReader("v1"), "", nil, false, StorageClassStandard, nil); err != nil {
		t.Fatalf("put v1: %v", err)
	}
	marker, err := store.DeleteWithVersion(ctx, "vb4", "k", "")
	if err != nil || marker == nil || !marker.IsDeleteMarker {
		t.Fatalf("delete = %+v, %v; want a delete marker", marker, err)
	}

	if err := store.SetStorageClass("vb4", "k", marker.VersionID, StorageClassDeepArchive); err != ErrObjectNotFound {
		t.Fatalf("SetStorageClass on marker = %v, want ErrObjectNotFound", err)
	}
	if _, err := store.GetACLWithVersion("vb4", "k", marker.VersionID); err != ErrObjectNotFound {
		t.Fatalf("GetACLWithVersion on marker = %v, want ErrObjectNotFound", err)
	}
}

// Mutating a legacy null record on an enabled bucket addresses the record
// where it lives and never grows a "_latest" pointer for it — the layout
// stays "no pointer until the first versioned write".
func TestMutateLegacyNullCreatesNoPointer(t *testing.T) {
	store, buckets := newMultipartTestStore(t)
	ctx := context.Background()

	if _, err := buckets.Create("vb5", "test-region"); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if _, err := store.PutWithVersioning(ctx, "vb5", "k", strings.NewReader("pre-versioning"), "", nil, false, StorageClassStandard, nil); err != nil {
		t.Fatalf("legacy put: %v", err)
	}
	if err := buckets.SetVersioning("vb5", BucketVersioningEnabled, ""); err != nil {
		t.Fatalf("SetVersioning enabled: %v", err)
	}

	if err := store.SetStorageClass("vb5", "k", "", StorageClassGlacier); err != nil {
		t.Fatalf("SetStorageClass on legacy null: %v", err)
	}
	if store.BaseStore.Exists(store.latestKeyStorageKey("vb5", "k")) {
		t.Fatal("mutation grew a _latest pointer for a legacy null record")
	}
	obj, err := store.GetMetadata("vb5", "k")
	if err != nil {
		t.Fatalf("GetMetadata: %v", err)
	}
	if obj.StorageClass != StorageClassGlacier {
		t.Fatalf("legacy record storage class = %q, want GLACIER", obj.StorageClass)
	}

	// The first versioned write still flips the legacy record as before.
	if _, err := store.PutWithVersioning(ctx, "vb5", "k", strings.NewReader("versioned"), "", nil, false, StorageClassStandard, nil); err != nil {
		t.Fatalf("versioned put: %v", err)
	}
	if n, err := store.CountByBucket("vb5"); err != nil || n != 1 {
		t.Fatalf("count after versioned write = %d, %v; want 1", n, err)
	}
}

// An encrypted put on a suspended bucket runs through the same persistence
// tail as a plain one: the null record becomes current and the "_latest"
// pointer left by the earlier enabled period is demoted. Before the put
// paths were unified, the encrypted variant skipped the demotion and the
// stale pointer kept serving the old version.
func TestEncryptedPutOnSuspendedBucketDemotesPointer(t *testing.T) {
	store, buckets := newMultipartTestStore(t)
	ctx := context.Background()

	if _, err := buckets.Create("vb6", "test-region"); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if err := buckets.SetVersioning("vb6", BucketVersioningEnabled, ""); err != nil {
		t.Fatalf("SetVersioning enabled: %v", err)
	}
	if _, err := store.PutWithVersioning(ctx, "vb6", "k", strings.NewReader("old"), "", nil, false, StorageClassStandard, nil); err != nil {
		t.Fatalf("versioned put: %v", err)
	}

	if err := buckets.SetVersioning("vb6", BucketVersioningSuspended, ""); err != nil {
		t.Fatalf("SetVersioning suspended: %v", err)
	}

	sse := &SSEObjectMetadata{EncryptionType: SSETypeAES256, UnencryptedSize: 3}
	if _, err := store.PutEncryptedWithVersioning(ctx, "vb6", "k", []byte("new"), "text/plain", nil, sse, false, StorageClassStandard, nil); err != nil {
		t.Fatalf("encrypted put: %v", err)
	}

	if store.BaseStore.Exists(store.latestKeyStorageKey("vb6", "k")) {
		t.Fatal("encrypted suspended put left the stale _latest pointer in place")
	}
	obj, err := store.GetMetadata("vb6", "k")
	if err != nil {
		t.Fatalf("GetMetadata: %v", err)
	}
	if obj.SSEMetadata == nil {
		t.Fatal("encrypted put did not stamp SSE metadata on the record")
	}
	// Blob content is outside this harness (the test blob store's puts are
	// inert stubs); the SDK suite covers the encrypted round-trip against
	// the real blob tier.
}

// The null version is a legal version address on a non-enabled bucket: the
// unversioned object IS the null version, so DeleteWithVersion("null")
// must remove its record permanently instead of reporting that the version
// does not exist — the delete the SDK suite's cleanup issues from
// ListObjectVersions output.
func TestNullVersionAddressedDeleteOnNonEnabledBucket(t *testing.T) {
	store, buckets := newMultipartTestStore(t)
	ctx := context.Background()

	if _, err := buckets.Create("nb", "test-region"); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if _, err := store.PutWithVersioning(ctx, "nb", "a.txt", strings.NewReader("plain"), "", nil, false, StorageClassStandard, nil); err != nil {
		t.Fatalf("put unversioned: %v", err)
	}

	if _, err := store.DeleteWithVersion(ctx, "nb", "a.txt", "null"); err != nil {
		t.Fatalf("DeleteWithVersion(null) on never-versioned bucket: %v", err)
	}
	if n, err := store.CountByBucket("nb"); err != nil || n != 0 {
		t.Fatalf("count after null delete = %d, %v; want 0", n, err)
	}

	// A genuinely missing version id on a never-versioned bucket still
	// reports the not-enabled sentinel the API layer maps to NoSuchVersion.
	if _, err := store.DeleteWithVersion(ctx, "nb", "a.txt", "does-not-exist"); err == nil {
		t.Fatalf("missing version id on never-versioned bucket unexpectedly succeeded")
	}

	// A suspended bucket's null record deletes the same way — permanently,
	// without adding a delete marker.
	if _, err := buckets.Create("sb", "test-region"); err != nil {
		t.Fatalf("Create sb: %v", err)
	}
	if err := buckets.SetVersioning("sb", BucketVersioningEnabled, ""); err != nil {
		t.Fatalf("SetVersioning enabled: %v", err)
	}
	if _, err := store.PutWithVersioning(ctx, "sb", "real.txt", strings.NewReader("real"), "", nil, false, StorageClassStandard, nil); err != nil {
		t.Fatalf("put real version: %v", err)
	}
	if err := buckets.SetVersioning("sb", BucketVersioningSuspended, ""); err != nil {
		t.Fatalf("SetVersioning suspended: %v", err)
	}
	if _, err := store.PutWithVersioning(ctx, "sb", "null-rec.txt", strings.NewReader("suspended write"), "", nil, false, StorageClassStandard, nil); err != nil {
		t.Fatalf("put suspended null record: %v", err)
	}

	if _, err := store.DeleteWithVersion(ctx, "sb", "null-rec.txt", "null"); err != nil {
		t.Fatalf("DeleteWithVersion(null) on suspended bucket: %v", err)
	}
	if n, err := store.CountByBucket("sb"); err != nil || n != 1 {
		t.Fatalf("count after suspended null delete = %d, %v; want 1 (the real version only)", n, err)
	}
}

// A read that resolves to a delete marker reports the marker record
// alongside not-found on every read path — Head matches the Get and range
// paths — so the service layer can surface the marker's metadata.
func TestHeadWithVersionSurfacesDeleteMarker(t *testing.T) {
	store, buckets := newMultipartTestStore(t)
	ctx := context.Background()

	if _, err := buckets.Create("mb", "test-region"); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if err := buckets.SetVersioning("mb", BucketVersioningEnabled, ""); err != nil {
		t.Fatalf("SetVersioning enabled: %v", err)
	}
	if _, err := store.PutWithVersioning(ctx, "mb", "k", strings.NewReader("one"), "", nil, false, StorageClassStandard, nil); err != nil {
		t.Fatalf("put: %v", err)
	}
	marker, err := store.DeleteWithVersion(ctx, "mb", "k", "")
	if err != nil || marker == nil || !marker.IsDeleteMarker {
		t.Fatalf("delete = %+v, %v; want a delete marker", marker, err)
	}

	obj, err := store.HeadWithVersion(ctx, "mb", "k", marker.VersionID)
	if !errors.Is(err, ErrObjectNotFound) {
		t.Fatalf("head marker version err = %v, want ErrObjectNotFound", err)
	}
	if obj == nil || !obj.IsDeleteMarker || obj.VersionID != marker.VersionID {
		t.Fatalf("head marker version = %+v, want the marker record", obj)
	}
}
