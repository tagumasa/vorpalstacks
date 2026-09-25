package dynamodb

import (
	"time"

	"vorpalstacks/internal/core/logs"
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
	s.startIntervalSweeper(&s.streamSweepOnce, retentionSweepInterval, "dynamodb retention sweep", s.sweepStoreRetentions)
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
