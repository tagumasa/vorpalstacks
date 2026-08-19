package dynamodb

import (
	"time"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/resilience"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// retentionSweepInterval controls how often the retention pruner runs.
// Stream records and contributor access counters only leave the store
// through this sweep, so without it every captured change and every
// counted access would accumulate for the lifetime of the table.
const retentionSweepInterval = time.Minute

// ensureRetentionSweeper starts the background pruner that keeps stream
// records and contributor access counters inside the 24-hour retention
// window documented for DynamoDB Streams.
func (s *DynamoDBService) ensureRetentionSweeper() {
	s.streamSweepOnce.Do(func() {
		s.bgWg.Add(1)
		go func() {
			defer func() { resilience.RecoverPanic("dynamodb retention sweep") }()
			defer s.bgWg.Done()
			ticker := time.NewTicker(retentionSweepInterval)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					s.stores.Range(func(_, v any) bool {
						store, ok := v.(dbstore.DynamoDBStoreInterface)
						if !ok {
							return true
						}
						s.sweepStoreRetentions(store)
						return true
					})
				case <-s.bgCtx.Done():
					return
				}
			}
		}()
	})
}

// sweepStoreRetentions trims the stream records of every streaming table
// and the contributor access counters of every insights-enabled table in
// one regional store to the retention window.
func (s *DynamoDBService) sweepStoreRetentions(store dbstore.DynamoDBStoreInterface) {
	cutoff := streamTimeNow().Add(-dbstore.StreamRetention)
	marker := ""
	for {
		tables, next, err := store.Tables().List(marker, 0)
		if err != nil {
			logs.Error("Failed to list tables for retention sweep", logs.Err(err))
			return
		}
		for _, table := range tables {
			if table.StreamSpecification != nil && table.StreamSpecification.StreamEnabled {
				if err := store.Streams().TrimOlderThan(table.Name, cutoff); err != nil {
					logs.Error("Failed to trim stream records past retention",
						logs.String("table", table.Name), logs.Err(err))
				}
			}
			if table.ContributorInsightsEnabled {
				if err := store.Contributors().SweepTableOlderThan(table.Name, cutoff); err != nil {
					logs.Error("Failed to sweep contributor counters past retention",
						logs.String("table", table.Name), logs.Err(err))
				}
			}
		}
		if next == "" {
			return
		}
		marker = next
	}
}
