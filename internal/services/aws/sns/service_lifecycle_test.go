package sns

import (
	"testing"

	snsstore "vorpalstacks/internal/store/aws/sns"
)

// closeCountingStore counts Close calls; every other interface method
// rides the embedded nil interface, so the test fails loudly if the
// lifecycle sweep ever calls anything but Close.
type closeCountingStore struct {
	snsstore.SNSStoreInterface
	closes int
}

func (c *closeCountingStore) Close() { c.closes++ }

// TestServiceCloseClosesEveryCachedStore pins the lifecycle sweep: every
// regional store the resolver cached carries a background deduplication
// sweeper only its own Close stops, so the service's Close must close the
// whole cache — lazily created regional stores included, which nothing
// else owns — exactly once each.
func TestServiceCloseClosesEveryCachedStore(t *testing.T) {
	svc := NewSNSService(nil, "123456789012", "us-east-1")
	injected := &closeCountingStore{}
	lazy := &closeCountingStore{}
	svc.stores.Store("us-east-1", injected)
	svc.stores.Store("eu-west-1", lazy)

	svc.Close()

	if injected.closes != 1 || lazy.closes != 1 {
		t.Fatalf("Close swept us-east-1 %d times and eu-west-1 %d times; want exactly once each — regional stores' sweep goroutines leak otherwise",
			injected.closes, lazy.closes)
	}
}
