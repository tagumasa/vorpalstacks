package common

import (
	"runtime"
	"sync"
	"testing"
)

// Deleting a key right after unlocking it was the historical break: a
// waiter stranded on the removed mutex lets the next Lock create a fresh
// one, so two goroutines enter the same critical section and lose updates.
// The counter must therefore survive a delete-after-unlock storm.
func TestKeyLockerMutualExclusionUnderConcurrentDelete(t *testing.T) {
	var kl KeyLocker
	var counter int
	const goroutines = 8
	const iters = 200

	var wg sync.WaitGroup
	for g := 0; g < goroutines; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < iters; i++ {
				kl.Lock("hot-key")
				c := counter
				runtime.Gosched()
				counter = c + 1
				kl.Unlock("hot-key")
				kl.Delete("hot-key")
			}
		}()
	}
	wg.Wait()

	if counter != goroutines*iters {
		t.Fatalf("counter = %d, want %d — mutual exclusion was broken", counter, goroutines*iters)
	}
}

// Entries self-clean at their last Unlock, and Delete only reaps idle
// entries: one held entry must keep excluding a second locker even when
// Delete runs in between.
func TestKeyLockerSelfCleansAndDeleteOnlyReapsIdle(t *testing.T) {
	var kl KeyLocker

	kl.Lock("k")
	acquired := make(chan struct{})
	go func() {
		kl.Lock("k")
		close(acquired)
		kl.Unlock("k")
	}()

	// While the entry is held (and waited on), Delete must not remove it;
	// the second locker stays blocked.
	kl.Delete("k")
	select {
	case <-acquired:
		t.Fatal("second locker entered while the first still held the lock")
	default:
	}
	kl.Unlock("k")
	<-acquired

	remaining := 0
	kl.Range(func(key string) bool {
		remaining++
		return true
	})
	if remaining != 0 {
		t.Fatalf("entries remaining after last unlock = %d, want 0", remaining)
	}
}
