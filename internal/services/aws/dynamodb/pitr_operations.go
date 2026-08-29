// Package dynamodb provides DynamoDB service operations for vorpalstacks.
package dynamodb

import (
	"context"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/resilience"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// DescribeContinuousBackups returns the continuous backup settings for a table.
func (s *DynamoDBService) DescribeContinuousBackups(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.describeContinuousBackupsCore(ctx, reqCtx, describeContinuousBackupsInput{
		Parameters: req.Parameters,
	})
}

// UpdateContinuousBackups enables or disables continuous backup for a table.
func (s *DynamoDBService) UpdateContinuousBackups(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.updateContinuousBackupsCore(ctx, reqCtx, updateContinuousBackupsInput{
		Parameters: req.Parameters,
	})
}

// journalSweepInterval controls how often the journal pruner runs. The
// journal only shrinks through pruning, so without a sweep every journaled
// write would accumulate for the lifetime of the process.
const journalSweepInterval = time.Minute

// ensureJournalSweeper starts the background pruner that keeps journals
// bounded: records at or before a table's earliest restorable time can
// never be replayed, and tables without recovery enabled have no use for a
// journal at all.
func (s *DynamoDBService) ensureJournalSweeper() {
	s.journalSweepOnce.Do(func() {
		s.bgWg.Add(1)
		go func() {
			defer func() { resilience.RecoverPanic("dynamodb journal sweep") }()
			defer s.bgWg.Done()
			ticker := time.NewTicker(journalSweepInterval)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					s.stores.Range(func(_, v any) bool {
						store, ok := v.(dbstore.DynamoDBStoreInterface)
						if !ok {
							return true
						}
						s.sweepStoreJournals(store)
						return true
					})
				case <-s.bgCtx.Done():
					return
				}
			}
		}()
	})
}

// sweepStoreJournals prunes the journals of every table in one regional
// store according to its point-in-time recovery state.
func (s *DynamoDBService) sweepStoreJournals(store dbstore.DynamoDBStoreInterface) {
	marker := ""
	for {
		tables, next, err := store.Tables().List(marker, 0)
		if err != nil {
			logs.Error("Failed to list tables for journal sweep", logs.Err(err))
			return
		}
		for _, table := range tables {
			pitr := table.PointInTimeRecovery
			if pitr == nil || pitr.Status != dbstore.PITRStatusEnabled {
				if err := store.Journal().DeleteAllForTable(table.Name); err != nil {
					logs.Error("Failed to drop journal for table without recovery",
						logs.String("table", table.Name), logs.Err(err))
				}
				continue
			}
			if _, err := store.Journal().DeleteOlderThan(table.Name, pitr.EarliestRestorableDateTime); err != nil {
				logs.Error("Failed to prune journal before earliest restorable time",
					logs.String("table", table.Name), logs.Err(err))
			}
		}
		if next == "" {
			return
		}
		marker = next
	}
}
