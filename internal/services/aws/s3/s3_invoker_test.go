package s3

import (
	"context"
	"sync"
	"testing"

	"vorpalstacks/internal/core/storage"
	s3store "vorpalstacks/internal/store/aws/s3"
)

// newInvokerTestService builds an S3Service backed by a throwaway store so
// the invoker methods can be exercised without the full server. The blob
// store is the hybrid one production wires, so object writes are real.
func newInvokerTestService(t *testing.T) (*S3Service, s3store.S3StoreInterface) {
	t.Helper()
	dataPath := t.TempDir()
	mgr, err := storage.NewRegionStorageManager(&storage.Config{Path: dataPath})
	if err != nil {
		t.Fatalf("open region storage: %v", err)
	}
	t.Cleanup(func() { _ = mgr.Close() })
	globalStore, err := storage.Open(dataPath)
	if err != nil {
		t.Fatalf("open global storage: %v", err)
	}
	t.Cleanup(func() { _ = globalStore.Close() })
	blobStore, err := storage.NewHybridBlobStore(globalStore, dataPath)
	if err != nil {
		t.Fatalf("open blob store: %v", err)
	}
	store := s3store.NewS3Store(mgr, blobStore, "000000000000")
	return NewS3Service(store, nil, "000000000000"), store
}

// Concurrent callers racing to create the same internal bucket must all
// succeed: the loser of the create race receives the store's
// already-exists sentinel, which EnsureBucket tolerates.
func TestEnsureBucketConcurrentCreationIsTolerated(t *testing.T) {
	svc, store := newInvokerTestService(t)

	const callers = 32
	start := &sync.WaitGroup{}
	start.Add(1)
	var wg sync.WaitGroup
	wg.Add(callers)
	errs := make(chan error, callers)
	for i := 0; i < callers; i++ {
		go func() {
			defer wg.Done()
			start.Wait()
			if err := svc.EnsureBucket(context.Background(), "us-east-1", "race-bucket"); err != nil {
				errs <- err
			}
		}()
	}
	start.Done()
	wg.Wait()
	close(errs)
	for err := range errs {
		t.Errorf("EnsureBucket during concurrent creation: %v", err)
	}

	if !store.Buckets("us-east-1").Exists("race-bucket") {
		t.Fatal("bucket was not created by the winning caller")
	}
	if err := svc.EnsureBucket(context.Background(), "us-east-1", "race-bucket"); err != nil {
		t.Errorf("EnsureBucket on an existing bucket: %v", err)
	}
}

// A cross-service write addressed to a region that does not hold the named
// bucket lands in the bucket's own region — bucket names are unique
// platform-wide, so the addressed region's store would be a shadow
// namespace no bucket record references. The multi-region CloudTrail
// trail's away-region deliveries rely on this resolution.
func TestPutObjectAddressesBucketOwnRegion(t *testing.T) {
	svc, store := newInvokerTestService(t)

	if err := svc.EnsureBucket(context.Background(), "us-east-1", "owner-region-bucket"); err != nil {
		t.Fatalf("ensure bucket: %v", err)
	}
	if err := svc.PutObject(context.Background(), "eu-west-1", "owner-region-bucket", "away/key.json.gz", []byte("payload"), "application/gzip"); err != nil {
		t.Fatalf("put via away region: %v", err)
	}

	reader, _, err := store.Objects("us-east-1").Get(context.Background(), "owner-region-bucket", "away/key.json.gz")
	if err != nil {
		t.Fatalf("object not readable in the bucket's own region: %v", err)
	}
	reader.Close()
	// The blob store is shared across regions, so the away region's
	// object RECORD is the discriminator: no record, no object there.
	if _, _, err := store.Objects("eu-west-1").Get(context.Background(), "owner-region-bucket", "away/key.json.gz"); err == nil {
		t.Fatal("the away region must not carry the object's record")
	}

	// An unknown bucket keeps the addressed region's store: the plain
	// per-region behaviour, unchanged shape.
	if err := svc.PutObject(context.Background(), "eu-west-1", "never-created-bucket", "k", []byte("v"), "text/plain"); err != nil {
		t.Fatalf("put to an unknown bucket: %v", err)
	}
	if _, _, err := store.Objects("eu-west-1").Get(context.Background(), "never-created-bucket", "k"); err != nil {
		t.Fatal("an unknown bucket's write must keep the addressed region's store")
	}
}

// Reads, existence, policy and deletion follow the same owner-region
// resolution the writes established: a cross-service consumer addressing a
// foreign region still finds the bucket, its policy and its objects —
// the trail-destination and import-source validations depend on it.
func TestInvokerReadsFollowBucketOwnerRegion(t *testing.T) {
	svc, store := newInvokerTestService(t)
	ctx := context.Background()

	if err := svc.EnsureBucket(ctx, "us-east-1", "read-owner-bucket"); err != nil {
		t.Fatalf("ensure bucket: %v", err)
	}
	if err := svc.PutObject(ctx, "eu-west-1", "read-owner-bucket", "logs/one.json.gz", []byte("payload"), "application/gzip"); err != nil {
		t.Fatalf("put via away region: %v", err)
	}
	const policy = `{"Version":"2012-10-17","Statement":[]}`
	if err := store.Buckets("us-east-1").SetPolicy("read-owner-bucket", policy); err != nil {
		t.Fatalf("set policy: %v", err)
	}

	data, err := svc.GetObject(ctx, "eu-west-1", "read-owner-bucket", "logs/one.json.gz", 0)
	if err != nil {
		t.Fatalf("GetObject via away region: %v", err)
	}
	if string(data) != "payload" {
		t.Fatalf("GetObject returned %q", string(data))
	}
	exists, err := svc.BucketExists(ctx, "eu-west-1", "read-owner-bucket")
	if err != nil || !exists {
		t.Fatalf("BucketExists via away region = (%v, %v), want true", exists, err)
	}
	got, err := svc.GetBucketPolicy(ctx, "eu-west-1", "read-owner-bucket")
	if err != nil || got != policy {
		t.Fatalf("GetBucketPolicy via away region = (%q, %v)", got, err)
	}
	keys, err := svc.ListObjects(ctx, "eu-west-1", "read-owner-bucket", "logs/", 0)
	if err != nil || len(keys) != 1 || keys[0] != "logs/one.json.gz" {
		t.Fatalf("ListObjects via away region = (%v, %v)", keys, err)
	}

	if err := svc.DeleteObject(ctx, "eu-west-1", "read-owner-bucket", "logs/one.json.gz"); err != nil {
		t.Fatalf("DeleteObject via away region: %v", err)
	}
	if _, _, err := store.Objects("us-east-1").Get(ctx, "read-owner-bucket", "logs/one.json.gz"); err == nil {
		t.Fatal("the object must be gone from the owner region's store")
	}

	// An unknown bucket keeps the addressed region's plain behaviour:
	// existence is false and the policy read surfaces the region store's
	// not-found error, exactly as before the routing.
	missing, err := svc.BucketExists(ctx, "eu-west-1", "read-never-created")
	if err != nil || missing {
		t.Fatalf("BucketExists on an unknown bucket = (%v, %v), want false", missing, err)
	}
	if _, err := svc.GetBucketPolicy(ctx, "eu-west-1", "read-never-created"); err == nil {
		t.Fatal("GetBucketPolicy on an unknown bucket must keep the region store's not-found error")
	}
}

// A bounded read that the object exceeds errors instead of returning a
// silent prefix: the bounded invoker consumers parse whole objects, for
// which truncated data is corruption. A bound at or above the object size
// returns it whole, and the zero bound stays unlimited.
func TestGetObjectBoundedReadErrorsOnExcess(t *testing.T) {
	svc, _ := newInvokerTestService(t)
	ctx := context.Background()

	if err := svc.EnsureBucket(ctx, "us-east-1", "bounded-bucket"); err != nil {
		t.Fatalf("ensure bucket: %v", err)
	}
	if err := svc.PutObject(ctx, "us-east-1", "bounded-bucket", "obj", []byte("12345"), "text/plain"); err != nil {
		t.Fatalf("put object: %v", err)
	}

	if _, err := svc.GetObject(ctx, "us-east-1", "bounded-bucket", "obj", 4); err == nil {
		t.Fatal("a read bounded below the object size must error, not truncate")
	}
	data, err := svc.GetObject(ctx, "us-east-1", "bounded-bucket", "obj", 5)
	if err != nil || string(data) != "12345" {
		t.Fatalf("a read bounded at the object size = (%q, %v)", string(data), err)
	}
	data, err = svc.GetObject(ctx, "us-east-1", "bounded-bucket", "obj", 0)
	if err != nil || string(data) != "12345" {
		t.Fatalf("the unlimited read = (%q, %v)", string(data), err)
	}
}
