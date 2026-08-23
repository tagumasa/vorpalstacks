package s3

import (
	"context"
	"sync"
	"testing"

	"vorpalstacks/internal/core/storage"
	s3store "vorpalstacks/internal/store/aws/s3"
)

// newInvokerTestService builds an S3Service backed by a throwaway store so
// the invoker methods can be exercised without the full server.
func newInvokerTestService(t *testing.T) (*S3Service, s3store.S3StoreInterface) {
	t.Helper()
	mgr, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatalf("open region storage: %v", err)
	}
	t.Cleanup(func() { _ = mgr.Close() })
	store := s3store.NewS3Store(mgr, nil, "000000000000")
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
