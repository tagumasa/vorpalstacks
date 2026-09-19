package kinesis

import (
	"fmt"
	"sync"
	"testing"

	"vorpalstacks/internal/core/storage"
)

// newLockTestStore opens a throwaway region storage and returns a store
// over it for lock-matrix interleaving tests.
func newLockTestStore(t *testing.T) *KinesisStore {
	t.Helper()
	mgr, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatalf("new region storage manager: %v", err)
	}
	bs, err := mgr.GetStorage("us-east-1")
	if err != nil {
		t.Fatalf("get region storage: %v", err)
	}
	return NewKinesisStore(bs, "000000000000", "us-east-1")
}

// TestConcurrentIteratorCreationIsSynchronised drives concurrent iterator
// creation, which advances the store's expiry-sweep bookkeeping: under
// -race an unsynchronised cleanup field fails here.
func TestConcurrentIteratorCreationIsSynchronised(t *testing.T) {
	store := newLockTestStore(t)
	if _, err := store.CreateStream("iter_race", 1, StreamModeProvisioned, 0, 0, nil); err != nil {
		t.Fatalf("create stream: %v", err)
	}

	var wg sync.WaitGroup
	for g := 0; g < 8; g++ {
		wg.Add(1)
		go func(g int) {
			defer wg.Done()
			for i := 0; i < 25; i++ {
				if _, err := store.CreateShardIterator("iter_race", "shardId-000000000000", "LATEST", "", nil); err != nil {
					t.Errorf("create iterator: %v", err)
					return
				}
			}
		}(g)
	}
	wg.Wait()
}

// TestConcurrentRegistrationsRespectQuota registers thirty distinct
// consumer names concurrently while configuration updates interleave: the
// locked registration path must admit exactly the quota (twenty), and the
// interleaved configuration writes must never revert the consumer count.
func TestConcurrentRegistrationsRespectQuota(t *testing.T) {
	store := newLockTestStore(t)
	if _, err := store.CreateStream("quota_race", 1, StreamModeProvisioned, 0, 0, nil); err != nil {
		t.Fatalf("create stream: %v", err)
	}

	stop := make(chan struct{})
	var configWG sync.WaitGroup
	for g := 0; g < 3; g++ {
		configWG.Add(1)
		go func() {
			defer configWG.Done()
			for {
				select {
				case <-stop:
					return
				default:
					_, _ = store.UpdateStreamFields("quota_race", func(stream *Stream) error {
						stream.EncryptionType = "KMS"
						return nil
					})
				}
			}
		}()
	}

	var mu sync.Mutex
	successes := 0
	var registerWG sync.WaitGroup
	for i := 0; i < 30; i++ {
		registerWG.Add(1)
		go func(i int) {
			defer registerWG.Done()
			if _, err := store.RegisterStreamConsumer(
				fmt.Sprintf("arn:aws:kinesis:us-east-1:000000000000:stream/quota_race"),
				fmt.Sprintf("reader-%02d", i),
				nil,
			); err == nil {
				mu.Lock()
				successes++
				mu.Unlock()
			}
		}(i)
	}
	registerWG.Wait()
	close(stop)
	configWG.Wait()

	if successes != MaxConsumersPerStream {
		t.Fatalf("concurrent registrations: %d succeeded, want exactly the quota %d", successes, MaxConsumersPerStream)
	}
	stream, err := store.GetStream("quota_race")
	if err != nil {
		t.Fatalf("get stream: %v", err)
	}
	if stream.ConsumerCount != int32(MaxConsumersPerStream) {
		t.Fatalf("consumer count after interleaved config writes: got %d, want %d (no reversion)", stream.ConsumerCount, MaxConsumersPerStream)
	}
}

// TestConfigUpdateCannotResurrectDeletedStream flips a configuration field
// concurrently with stream deletion: a read-modify-write that acquired the
// lock after the deletion must observe not-found, never write the stream
// back into existence.
func TestConfigUpdateCannotResurrectDeletedStream(t *testing.T) {
	store := newLockTestStore(t)
	if _, err := store.CreateStream("zombie", 1, StreamModeProvisioned, 0, 0, nil); err != nil {
		t.Fatalf("create stream: %v", err)
	}

	stop := make(chan struct{})
	var wg sync.WaitGroup
	for g := 0; g < 4; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-stop:
					return
				default:
					_, _ = store.UpdateStreamFields("zombie", func(stream *Stream) error {
						stream.RetentionPeriodHours = MinRetentionPeriodHours + 1
						return nil
					})
				}
			}
		}()
	}

	if err := store.DeleteStream("zombie"); err != nil {
		t.Fatalf("delete stream: %v", err)
	}
	close(stop)
	wg.Wait()

	if _, err := store.GetStream("zombie"); err == nil {
		t.Fatal("the deleted stream was resurrected by a concurrent configuration write")
	}
}
