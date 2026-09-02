package wafv2

import (
	"time"

	"vorpalstacks/internal/core/logs"
)

// sampleSweepInterval is how often the retention sweep drops expired
// sampled requests and trims per-rule depth. This is a platform tuning
// value: the API-visible retention bound is wafstore.SampleRetention,
// and Query enforces it on every read regardless of the sweep.
const sampleSweepInterval = time.Minute

// startSampleRetentionSweep lazily starts the periodic retention sweep
// over every regional sampling store created so far. Stores created
// after a sweep tick are picked up by the next one because the sweep
// ranges over the live store map.
func (s *WAFv2Service) startSampleRetentionSweep() {
	s.sampleSweepOnce.Do(func() {
		go func() {
			ticker := time.NewTicker(sampleSweepInterval)
			defer ticker.Stop()
			for now := range ticker.C {
				s.stores.Range(func(_, value any) bool {
					stores, ok := value.(*wafv2Stores)
					if !ok {
						return true
					}
					if err := stores.samples.PurgeExpired(now); err != nil {
						logs.Warn("wafv2 sampled-request retention sweep failed, retrying next interval", logs.Err(err))
					}
					return true
				})
			}
		}()
	})
}
