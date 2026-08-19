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
	// Smithy declares TableNotFoundException for this op.
	table, err := s.validateAndGetTableWithErr(reqCtx, req.Parameters, ErrTableNotFoundException)
	if err != nil {
		return nil, err
	}
	tableName := table.Name

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	pitr, err := store.Tables().GetPointInTimeRecovery(tableName)
	if err != nil {
		return nil, err
	}

	// Continuous backups are enabled on every table at creation, so the
	// outer status is always ENABLED; only the point-in-time recovery
	// status depends on the table's settings.
	pitrStatus := "DISABLED"
	pitrDescription := map[string]interface{}{
		"PointInTimeRecoveryStatus": pitrStatus,
	}
	recoveryPeriod := pitrDefaultRecoveryPeriodDays
	if pitr != nil && pitr.Status == dbstore.PITRStatusEnabled {
		pitrStatus = "ENABLED"
		pitrDescription["PointInTimeRecoveryStatus"] = pitrStatus
		// The restorable window reaches from the moment recovery was
		// enabled (the journal starts there) to the present; mutations
		// commit synchronously on this platform, so now is restorable.
		pitrDescription["EarliestRestorableDateTime"] = pitr.EarliestRestorableDateTime.Unix()
		pitrDescription["LatestRestorableDateTime"] = time.Now().Unix()
		if pitr.RecoveryPeriodInDays > 0 {
			recoveryPeriod = pitr.RecoveryPeriodInDays
			pitrDescription["RecoveryPeriodInDays"] = recoveryPeriod
		}
	}

	return map[string]interface{}{
		"ContinuousBackupsDescription": map[string]interface{}{
			"ContinuousBackupsStatus":        "ENABLED",
			"PointInTimeRecoveryDescription": pitrDescription,
		},
	}, nil
}

// UpdateContinuousBackups enables or disables continuous backup for a table.
func (s *DynamoDBService) UpdateContinuousBackups(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	// Smithy declares TableNotFoundException for this op.
	table, err := s.validateAndGetTableWithErr(reqCtx, req.Parameters, ErrTableNotFoundException)
	if err != nil {
		return nil, err
	}
	tableName := table.Name

	pitrSpec, ok := req.Parameters["PointInTimeRecoverySpecification"].(map[string]interface{})
	if !ok {
		return nil, ErrInvalidParameter
	}

	enabled, err := validateBoolParam(pitrSpec, "PointInTimeRecoveryEnabled", false)
	if err != nil {
		return nil, err
	}

	recoveryPeriod := pitrDefaultRecoveryPeriodDays
	if _, ok := pitrSpec["RecoveryPeriodInDays"]; ok {
		rp := request.GetIntParam(pitrSpec, "RecoveryPeriodInDays")
		if !validateRecoveryPeriodInDays(rp) {
			return nil, ErrInvalidParameter
		}
		recoveryPeriod = rp
	}

	pitr := &dbstore.PointInTimeRecoveryDescription{
		Status: dbstore.PITRStatusDisabled,
	}
	if enabled {
		pitr.Status = dbstore.PITRStatusEnabled
		// The restorable window starts at the moment recovery is enabled:
		// the journal only records changes from that point on, so nothing
		// earlier can be restored.
		pitr.EarliestRestorableDateTime = time.Now()
		pitr.LatestRestorableDateTime = time.Now()
		pitr.RecoveryPeriodInDays = recoveryPeriod
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.Tables().SetPointInTimeRecovery(tableName, pitr); err != nil {
		return nil, err
	}
	if !enabled {
		// Disabling recovery invalidates the journal: re-enabling starts a
		// fresh restorable window at the re-enable time.
		if err := store.Journal().DeleteAllForTable(tableName); err != nil {
			return nil, err
		}
	}

	return map[string]interface{}{
		"ContinuousBackupsDescription": map[string]interface{}{
			"ContinuousBackupsStatus": string(pitr.Status),
			"PointInTimeRecoveryDescription": map[string]interface{}{
				"PointInTimeRecoveryStatus": string(pitr.Status),
			},
		},
	}, nil
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
