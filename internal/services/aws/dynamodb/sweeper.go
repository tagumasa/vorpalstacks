package dynamodb

import (
	"sync"
	"time"

	"vorpalstacks/internal/core/resilience"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// startIntervalSweeper starts, behind once, the single background
// goroutine that runs sweep on every regional store cached by the
// service, once per interval, until the service context is cancelled.
// Every periodic pruner of the service (journal, retention,
// system-backup, idempotency) shares this starter; the pruners
// themselves differ in what they do to a store, not in how they are
// scheduled. Each sweep invocation carries its own panic boundary, so a
// panicking pruner ends that invocation alone and the ticker schedules
// the next one; the goroutine's outermost recover stays as the backstop
// for anything that panics outside a sweep call.
func (s *DynamoDBService) startIntervalSweeper(once *sync.Once, interval time.Duration, label string, sweep func(store dbstore.DynamoDBStoreInterface)) {
	once.Do(func() {
		s.bgWg.Add(1)
		go func() {
			defer func() {
				if r := recover(); r != nil {
					resilience.LogPanic(label, r)
				}
			}()
			defer s.bgWg.Done()
			ticker := time.NewTicker(interval)
			defer ticker.Stop()
			runSweep := func(store dbstore.DynamoDBStoreInterface) {
				defer func() {
					if r := recover(); r != nil {
						resilience.LogPanic(label, r)
					}
				}()
				sweep(store)
			}
			for {
				select {
				case <-ticker.C:
					s.stores.Range(func(_, v any) bool {
						store, ok := v.(dbstore.DynamoDBStoreInterface)
						if !ok {
							return true
						}
						runSweep(store)
						return true
					})
				case <-s.bgCtx.Done():
					return
				}
			}
		}()
	})
}
