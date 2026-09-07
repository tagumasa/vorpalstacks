package dynamodb

import (
	"context"
	"time"

	"vorpalstacks/internal/common/request"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// ---------------------------------------------------------------------------
// PITR Core — single validation + persistence path for the continuous
// backups (point-in-time recovery) operations.
//
// Both the HTTP API handlers (pitr_operations.go) and any future admin
// handler delegate to these methods to ensure identical behaviour.
// ---------------------------------------------------------------------------

// pitrEarliestRestorable returns the trailing edge of the restorable
// window. Continuous backups maintain a recovery period of preceding days
// (RecoveryPeriodInDays, default 35, range 1-35), so the earliest restorable
// moment is the later of the recovery-enable time and now minus the
// configured period — records older than the trailing edge can never be
// replayed and the journal sweeper prunes them at this cutoff.
func pitrEarliestRestorable(pitr *dbstore.PointInTimeRecoveryDescription, now time.Time) time.Time {
	periodDays := pitrDefaultRecoveryPeriodDays
	if pitr.RecoveryPeriodInDays > 0 {
		periodDays = pitr.RecoveryPeriodInDays
	}
	trailingEdge := now.AddDate(0, 0, -periodDays)
	if pitr.EarliestRestorableDateTime.After(trailingEdge) {
		return pitr.EarliestRestorableDateTime
	}
	return trailingEdge
}

// describeContinuousBackupsInput carries the raw wire parameters for
// DescribeContinuousBackups.
type describeContinuousBackupsInput struct {
	Parameters map[string]interface{}
}

// describeContinuousBackupsCore returns the continuous backups description
// of the named table.
func (s *DynamoDBService) describeContinuousBackupsCore(ctx context.Context, reqCtx *request.RequestContext, in describeContinuousBackupsInput) (interface{}, error) {
	// Smithy declares TableNotFoundException for this op.
	table, err := s.validateAndGetTableWithErr(reqCtx, in.Parameters, ErrTableNotFoundException)
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
		// The restorable window is the trailing recovery period ending at
		// the present; mutations commit synchronously on this platform, so
		// now is restorable.
		now := time.Now()
		pitrDescription["EarliestRestorableDateTime"] = pitrEarliestRestorable(pitr, now).Unix()
		pitrDescription["LatestRestorableDateTime"] = now.Unix()
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

// updateContinuousBackupsInput carries the raw wire parameters for
// UpdateContinuousBackups.
type updateContinuousBackupsInput struct {
	Parameters map[string]interface{}
}

// updateContinuousBackupsCore enables or disables point-in-time recovery
// for the named table.
func (s *DynamoDBService) updateContinuousBackupsCore(ctx context.Context, reqCtx *request.RequestContext, in updateContinuousBackupsInput) (interface{}, error) {
	// Smithy declares TableNotFoundException for this op.
	table, err := s.validateAndGetTableWithErr(reqCtx, in.Parameters, ErrTableNotFoundException)
	if err != nil {
		return nil, err
	}
	tableName := table.Name

	pitrSpec, ok := in.Parameters["PointInTimeRecoverySpecification"].(map[string]interface{})
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

	// Continuous backups are enabled on every table at creation, so the
	// outer status is always ENABLED — matching the describe response; only
	// the point-in-time recovery status depends on the table's settings.
	return map[string]interface{}{
		"ContinuousBackupsDescription": map[string]interface{}{
			"ContinuousBackupsStatus": "ENABLED",
			"PointInTimeRecoveryDescription": map[string]interface{}{
				"PointInTimeRecoveryStatus": string(pitr.Status),
			},
		},
	}, nil
}
