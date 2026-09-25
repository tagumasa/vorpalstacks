package dynamodb

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"vorpalstacks/internal/common/request"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// ---------------------------------------------------------------------------
// Backup Core — single validation + persistence path for backup operations.
//
// These methods encapsulate backup lifecycle logic. Both the HTTP API
// handlers (backup_operations.go) and any future admin handler delegate to
// these methods to ensure identical behaviour.
// ---------------------------------------------------------------------------

// createBackupInput carries the raw wire parameters for CreateBackup.
type createBackupInput struct {
	Parameters map[string]interface{}
}

// createBackupCore validates the request, persists a new backup record in
// CREATING state and launches the background snapshot goroutine. The table
// must be ACTIVE. Returns the persisted backup in its initial state.
func (s *DynamoDBService) createBackupCore(ctx context.Context, reqCtx *request.RequestContext, in createBackupInput) (*dbstore.Backup, error) {
	// Table must be ACTIVE to create a backup. CreateBackup's Smithy model
	// declares TableNotFoundException (not ResourceNotFoundException), so
	// surface the individual error sentinel here.
	table, err := s.validateAndGetActiveTableWithErr(reqCtx, in.Parameters, ErrTableNotFoundException)
	if err != nil {
		return nil, err
	}

	backupName := request.GetStringParam(in.Parameters, "BackupName")
	if !validateResourceName(backupName) {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	// The store's Create mints the record's identity: a generated id in
	// the ARN, the documented identity model — the name is a label and the
	// documented CreateBackup error set carries no duplicate-name
	// rejection, so same-named backups coexist.
	backup, err := store.Backups().Create(backupName, table.Name, table.ARN, table.TableSizeBytes)
	if err != nil {
		return nil, err
	}

	backup.SourceTableId = table.TableId
	backup.KeySchema = table.KeySchema
	backup.AttributeDefinitions = table.AttributeDefinitions
	backup.BillingMode = table.BillingMode
	backup.ProvisionedThroughput = table.ProvisionedThroughput
	backup.GlobalSecondaryIndexes = table.GlobalSecondaryIndexes
	backup.LocalSecondaryIndexes = table.LocalSecondaryIndexes
	backup.VectorIndexes = table.VectorIndexes
	backup.SourceTableCreationTime = table.CreationDateTime
	backup.SourceTableSizeBytes = table.TableSizeBytes
	backup.SourceTableItemCount = table.ItemCount
	backup.BackupSizeBytes = table.TableSizeBytes

	if err := store.Backups().Put(backup); err != nil {
		return nil, err
	}

	// The item snapshot is taken synchronously, before the response
	// returns: the backup must capture the table exactly as of the request
	// (CreateBackup is a point-in-time operation), not whatever the table
	// holds whenever a background scan would happen to reach it. The backup
	// is therefore complete — and AVAILABLE — by the time the caller hears
	// back. The scan runs inside a store View — a pebble snapshot — so the
	// whole collection observes one commit point: a concurrent write either
	// precedes the snapshot and appears in it in full, or follows it and is
	// absent; no write interleaves with the walk.
	var snapshotItems []*dbstore.Item
	if err := store.View(ctx, func(txn *dbstore.DynamoDBTxn) error {
		return txn.Scan(table.Name, func(item *dbstore.Item) error {
			snapshotItems = append(snapshotItems, &dbstore.Item{
				TableName:  table.Name,
				Key:        copyAttributes(item.Key),
				Attributes: copyAttributes(item.Attributes),
			})
			return nil
		})
	}); err != nil {
		// A half-created backup must not linger as CREATING forever: the
		// record and any partial snapshot are removed, and the caller
		// retries the whole creation.
		_ = store.Backups().Delete(backup.BackupArn)
		return nil, err
	}
	if err := store.Backups().SaveSnapshot(backup.BackupArn, snapshotItems); err != nil {
		_ = store.Backups().Delete(backup.BackupArn)
		return nil, err
	}
	backup.BackupStatus = dbstore.BackupStatusAvailable
	if err := store.Backups().Put(backup); err != nil {
		return nil, err
	}

	return backup, nil
}

// deleteBackupCore validates the request, then deletes a backup by ARN and
// returns the deleted backup record for response formatting.
func (s *DynamoDBService) deleteBackupCore(ctx context.Context, reqCtx *request.RequestContext, backupArn string) (*dbstore.Backup, error) {
	if !validateBackupArn(backupArn) {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	backup, err := store.Backups().Get(backupArn)
	if err != nil {
		// Absence maps to the family's not-found sentinel; a storage
		// fault (bucket I/O, an unreadable record) is reported as the
		// storage error it is — never masquerading as absence, the same
		// split the describe path applies.
		if storeRecordMissing(err) {
			return nil, ErrBackupNotFound
		}
		return nil, err
	}
	if err := store.Backups().Delete(backup.BackupArn); err != nil {
		return nil, err
	}
	return backup, nil
}

// describeBackupCore validates the request, then returns a backup by ARN.
func (s *DynamoDBService) describeBackupCore(ctx context.Context, reqCtx *request.RequestContext, backupArn string) (*dbstore.Backup, error) {
	return describeByArn(s, reqCtx, backupArn, validateBackupArn(backupArn), ErrBackupNotFound, func(store dbstore.DynamoDBStoreInterface) (*dbstore.Backup, error) {
		return store.Backups().Get(backupArn)
	})
}

// ListBackupsCoreInput is the service-layer DTO for ListBackups.
type ListBackupsCoreInput struct {
	TableName               string
	BackupTypeFilter        string
	TimeRangeLowerBound     *int64
	TimeRangeUpperBound     *int64
	Limit                   int
	ExclusiveStartBackupArn string
}

// ListBackupsCoreResult is the service-layer result of ListBackups.
type ListBackupsCoreResult struct {
	Backups                []*dbstore.Backup
	LastEvaluatedBackupArn string
}

// listBackupsCore validates the request, then returns a filtered, paginated
// list of backups. The walk is driven by the store's own continuation
// marker: filters are applied across store pages until the request's limit
// is filled or the store is exhausted, so an active filter can no longer
// terminate the walk early.
func (s *DynamoDBService) listBackupsCore(ctx context.Context, reqCtx *request.RequestContext, in ListBackupsCoreInput) (*ListBackupsCoreResult, error) {
	// The TableName member documents both forms — "the name or Amazon
	// Resource Name (ARN) of the table" — and is optional (an omitted
	// member lists every table's backups); the store filters on the bare
	// name, so a present member is normalised before the store call.
	if in.TableName != "" {
		tableName, terr := resolveTableNameMember(in.TableName)
		if terr != nil {
			return nil, terr
		}
		in.TableName = tableName
	}
	if in.Limit == 0 {
		in.Limit = listBackupsMaxLimit
	} else {
		if !validateListBackupsLimit(in.Limit) {
			return nil, ErrInvalidParameter
		}
	}
	if in.ExclusiveStartBackupArn != "" {
		if !validateBackupArn(in.ExclusiveStartBackupArn) {
			return nil, ErrInvalidParameter
		}
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	marker := ""
	if in.ExclusiveStartBackupArn != "" {
		marker = svcarn.ExtractBackupIdFromARN(in.ExclusiveStartBackupArn)
	}

	// BackupType ALL lists every on-demand backup type (model enum
	// USER|SYSTEM|AWS_BACKUP|ALL); only a specific type filters.
	typeFilter := in.BackupTypeFilter
	if typeFilter == "ALL" {
		typeFilter = ""
	}

	var filtered []*dbstore.Backup
	for {
		backups, nextMarker, err := store.Backups().List(marker, listBackupsMinFetchSize, in.TableName)
		if err != nil {
			return nil, err
		}
		for _, b := range backups {
			if typeFilter != "" && string(b.BackupType) != typeFilter {
				continue
			}
			backupTime := b.BackupCreationDateTime.Unix()
			if in.TimeRangeLowerBound != nil && backupTime < *in.TimeRangeLowerBound {
				continue
			}
			if in.TimeRangeUpperBound != nil && backupTime > *in.TimeRangeUpperBound {
				continue
			}
			filtered = append(filtered, b)
			// One item beyond the limit proves more results may exist.
			if len(filtered) > in.Limit {
				break
			}
		}
		if len(filtered) > in.Limit || nextMarker == "" {
			break
		}
		marker = nextMarker
	}

	result := &ListBackupsCoreResult{}
	if len(filtered) > in.Limit {
		result.Backups = filtered[:in.Limit]
		result.LastEvaluatedBackupArn = filtered[in.Limit-1].BackupArn
	} else {
		result.Backups = filtered
	}
	return result, nil
}

// RestoreTableFromBackupCoreInput is the service-layer DTO for
// RestoreTableFromBackup.
type RestoreTableFromBackupCoreInput struct {
	BackupArn       string
	TargetTableName string
	// Overrides carries the override member family the model defines on
	// RestoreTableFromBackupInput; the zero value restores the backup's
	// own settings.
	Overrides RestoreOverrides
}

// RestoreOverrides carries the override member family both restore
// surfaces (from-backup and point-in-time) accept, parsed from the wire by
// parseRestoreOverridesWire. The index members distinguish an absent
// member (nil — the base record's indexes are restored) from a present
// empty list (a non-nil empty slice — the documentation's "exclude some or
// all of the indexes at the time of restore").
type RestoreOverrides struct {
	BillingMode           dbstore.BillingMode
	ProvisionedThroughput *dbstore.ProvisionedThroughput
	OnDemandThroughput    *dbstore.OnDemandThroughput
	SSEDescription        *dbstore.SSEDescription
	// SSEDisable reports an SSESpecificationOverride that explicitly asked
	// to disable encryption; the restore core's SSE finaliser owns what
	// that means.
	SSEDisable             bool
	GlobalSecondaryIndexes []*dbstore.GlobalSecondaryIndex
	LocalSecondaryIndexes  []*dbstore.LocalSecondaryIndex
	VectorIndexes          []*dbstore.VectorIndex
}

// restoreSchema is the settings family both restore surfaces read from
// their base record (the backup or the source table) before the request's
// overrides are applied.
type restoreSchema struct {
	billingMode    dbstore.BillingMode
	provThroughput *dbstore.ProvisionedThroughput
	gsi            []*dbstore.GlobalSecondaryIndex
	lsi            []*dbstore.LocalSecondaryIndex
	vectorIdx      []*dbstore.VectorIndex
}

// applyRestoreOverrides overlays the request's override members on the base
// schema. Index overrides select from the base's existing indexes — the
// documentation offers exclusion, not creation ("The indexes provided
// should match existing secondary indexes. You can choose to exclude some
// or all of the indexes at the time of restore."), so an override naming
// an unknown index or restating a definition differently is a validation
// error. Throughput overrides replace wholesale, and the on-demand pair is
// rejected on a provisioned resulting mode — the mirror of the
// mode-consistency rule the creation path enforces, since on-demand
// maximum settings do not exist on a provisioned table.
func applyRestoreOverrides(base restoreSchema, overrides RestoreOverrides) (restoreSchema, error) {
	// BillingModeOverride is the closed BillingMode enum (Smithy:
	// PROVISIONED | PAY_PER_REQUEST), validated here — the core-side overlay
	// both restore surfaces share — before any member is applied, the same
	// guard the creation and update planes enforce through
	// validateBillingModeValue.
	if overrides.BillingMode != "" && !validateBillingModeValue(overrides.BillingMode) {
		return base, ErrInvalidParameter
	}
	if overrides.BillingMode != "" && overrides.BillingMode != base.billingMode {
		// A mode switch carries the capacity-mode contract: switching to
		// on-demand drops every provisioned setting (the table's and the
		// indexes' — secondary indexes inherit the table's mode), and
		// switching to provisioned drops the indexes' on-demand pair. The
		// index entries are copied before clearing so the base record the
		// restore read stays untouched.
		if overrides.BillingMode == dbstore.BillingModePayPerRequest {
			base.provThroughput = nil
			cleared := make([]*dbstore.GlobalSecondaryIndex, len(base.gsi))
			for i, g := range base.gsi {
				copied := *g
				copied.ProvisionedThroughput = nil
				cleared[i] = &copied
			}
			base.gsi = cleared
		} else if overrides.BillingMode == dbstore.BillingModeProvisioned {
			cleared := make([]*dbstore.GlobalSecondaryIndex, len(base.gsi))
			for i, g := range base.gsi {
				copied := *g
				copied.OnDemandThroughput = nil
				cleared[i] = &copied
			}
			base.gsi = cleared
		}
		base.billingMode = overrides.BillingMode
	}
	if overrides.ProvisionedThroughput != nil {
		base.provThroughput = overrides.ProvisionedThroughput
	}
	if overrides.GlobalSecondaryIndexes != nil {
		selected, err := selectGSIOverrides(base.gsi, overrides.GlobalSecondaryIndexes)
		if err != nil {
			return base, err
		}
		base.gsi = selected
	}
	if overrides.LocalSecondaryIndexes != nil {
		selected, err := selectLSIOverrides(base.lsi, overrides.LocalSecondaryIndexes)
		if err != nil {
			return base, err
		}
		base.lsi = selected
	}
	if overrides.VectorIndexes != nil {
		selected, err := selectVectorIndexOverrides(base.vectorIdx, overrides.VectorIndexes)
		if err != nil {
			return base, err
		}
		base.vectorIdx = selected
	}
	if base.billingMode == dbstore.BillingModeProvisioned && overrides.OnDemandThroughput != nil {
		return base, ErrInvalidParameter
	}
	return base, nil
}

// selectGSIOverrides narrows the restored table's global secondary indexes
// to those named by the override list, applying the provided projection and
// throughput settings. Restore overrides select from the existing indexes:
// an override naming an unknown index, or one replacing the key schema, is
// a validation error.
func selectGSIOverrides(existing []*dbstore.GlobalSecondaryIndex, overrides []*dbstore.GlobalSecondaryIndex) ([]*dbstore.GlobalSecondaryIndex, error) {
	byName := make(map[string]*dbstore.GlobalSecondaryIndex, len(existing))
	for _, g := range existing {
		byName[g.IndexName] = g
	}
	selected := make([]*dbstore.GlobalSecondaryIndex, 0, len(overrides))
	for _, ov := range overrides {
		base, ok := byName[ov.IndexName]
		if !ok {
			return nil, ErrInvalidParameter
		}
		if len(ov.KeySchema) > 0 && !keySchemasEqual(ov.KeySchema, base.KeySchema) {
			return nil, ErrInvalidParameter
		}
		restored := *base
		if ov.Projection != nil {
			restored.Projection = ov.Projection
		}
		if ov.ProvisionedThroughput != nil {
			restored.ProvisionedThroughput = ov.ProvisionedThroughput
		}
		if ov.OnDemandThroughput != nil {
			restored.OnDemandThroughput = ov.OnDemandThroughput
		}
		if ov.WarmThroughput != nil {
			restored.WarmThroughput = ov.WarmThroughput
		}
		selected = append(selected, &restored)
	}
	return selected, nil
}

// selectLSIOverrides narrows the restored table's local secondary indexes
// to those named by the override list, applying the provided projection.
// Local secondary index key schemas cannot change, so an override naming an
// unknown index or a different key schema is a validation error.
func selectLSIOverrides(existing []*dbstore.LocalSecondaryIndex, overrides []*dbstore.LocalSecondaryIndex) ([]*dbstore.LocalSecondaryIndex, error) {
	byName := make(map[string]*dbstore.LocalSecondaryIndex, len(existing))
	for _, l := range existing {
		byName[l.IndexName] = l
	}
	selected := make([]*dbstore.LocalSecondaryIndex, 0, len(overrides))
	for _, ov := range overrides {
		base, ok := byName[ov.IndexName]
		if !ok {
			return nil, ErrInvalidParameter
		}
		if len(ov.KeySchema) > 0 && !keySchemasEqual(ov.KeySchema, base.KeySchema) {
			return nil, ErrInvalidParameter
		}
		restored := *base
		if ov.Projection != nil {
			restored.Projection = ov.Projection
		}
		selected = append(selected, &restored)
	}
	return selected, nil
}

// selectVectorIndexOverrides narrows the restored table's vector indexes
// to those named by the override list. The documentation binds the member
// to exclusion: "The indexes provided must match existing vector indexes
// from the backup. You can choose to exclude some or all of the vector
// indexes at the time of restore." — an override naming an unknown index
// or restating its definition differently is a validation error.
func selectVectorIndexOverrides(existing []*dbstore.VectorIndex, overrides []*dbstore.VectorIndex) ([]*dbstore.VectorIndex, error) {
	byName := make(map[string]*dbstore.VectorIndex, len(existing))
	for _, v := range existing {
		byName[v.IndexName] = v
	}
	selected := make([]*dbstore.VectorIndex, 0, len(overrides))
	for _, ov := range overrides {
		base, ok := byName[ov.IndexName]
		if !ok {
			return nil, ErrInvalidParameter
		}
		if !vectorIndexDefinitionsEqual(ov, base) {
			return nil, ErrInvalidParameter
		}
		selected = append(selected, base)
	}
	return selected, nil
}

// projectionsEqual compares two projections element by element.
func projectionsEqual(a, b *dbstore.Projection) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	if a.ProjectionType != b.ProjectionType || len(a.NonKeyAttributes) != len(b.NonKeyAttributes) {
		return false
	}
	for i := range a.NonKeyAttributes {
		if a.NonKeyAttributes[i] != b.NonKeyAttributes[i] {
			return false
		}
	}
	return true
}

// vectorIndexDefinitionsEqual compares the definitional members of two
// vector indexes — everything except the runtime state (ARN, status,
// backfill flag, sizes) the restored table re-derives for itself.
func vectorIndexDefinitionsEqual(a, b *dbstore.VectorIndex) bool {
	if a.VectorAttributeName != b.VectorAttributeName || a.Dimensions != b.Dimensions ||
		a.DistanceFunction != b.DistanceFunction {
		return false
	}
	if len(a.SearchSchema) != len(b.SearchSchema) {
		return false
	}
	for i := range a.SearchSchema {
		if *a.SearchSchema[i] != *b.SearchSchema[i] {
			return false
		}
	}
	return projectionsEqual(a.Projection, b.Projection)
}

// restoreTableFromBackupCore validates the request, then creates a new
// table from a backup snapshot. The target table must not already exist.
func (s *DynamoDBService) restoreTableFromBackupCore(ctx context.Context, reqCtx *request.RequestContext, in RestoreTableFromBackupCoreInput) (*dbstore.Table, error) {
	if !validateBackupArn(in.BackupArn) {
		return nil, ErrInvalidParameter
	}
	if !validateResourceName(in.TargetTableName) {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	backup, err := store.Backups().Get(in.BackupArn)
	if err != nil {
		// Absence is the documented BackupNotFound; a storage fault on an
		// existing record is reported as the storage error it is — never
		// masquerading as absence, the same split the delete path applies.
		if storeRecordMissing(err) {
			return nil, ErrBackupNotFound
		}
		return nil, err
	}
	// A backup is restorable only in AVAILABLE state: the record is
	// observable in CREATING (a crash between the record's initial write
	// and the snapshot save leaves it there permanently), and restoring
	// such a backup would report success over zero items. The model
	// documents BackupInUseException for a control-plane conflict on the
	// backup ("being created, deleted or restored").
	if backup.BackupStatus != dbstore.BackupStatusAvailable {
		return nil, ErrBackupNotAvailable
	}

	if store.Tables().Exists(in.TargetTableName) {
		return nil, ErrTableAlreadyExistsException
	}

	var keySchema []*dbstore.KeySchemaElement
	var attrDefs []*dbstore.AttributeDefinition
	schema := restoreSchema{
		billingMode:    backup.BillingMode,
		provThroughput: backup.ProvisionedThroughput,
		gsi:            backup.GlobalSecondaryIndexes,
		lsi:            backup.LocalSecondaryIndexes,
		vectorIdx:      backup.VectorIndexes,
	}
	if len(backup.KeySchema) > 0 {
		keySchema = backup.KeySchema
		attrDefs = backup.AttributeDefinitions
	} else {
		sourceTable, err := store.Tables().Get(backup.SourceTableName)
		if err != nil {
			return nil, fmt.Errorf("backup %s has no key schema and source table %q not found: %w", in.BackupArn, backup.SourceTableName, err)
		}
		keySchema = sourceTable.KeySchema
		attrDefs = sourceTable.AttributeDefinitions
		schema = restoreSchema{
			billingMode:    sourceTable.BillingMode,
			provThroughput: sourceTable.ProvisionedThroughput,
			gsi:            sourceTable.GlobalSecondaryIndexes,
			lsi:            sourceTable.LocalSecondaryIndexes,
			vectorIdx:      sourceTable.VectorIndexes,
		}
	}
	schema, err = applyRestoreOverrides(schema, in.Overrides)
	if err != nil {
		return nil, err
	}

	if !validateBillingModeConsistency(schema.billingMode, schema.provThroughput) {
		return nil, ErrInvalidParameter
	}

	// The SSE override resolves before the target exists — the same point
	// the point-in-time surface does — so an unresolvable specification
	// rejects the request outright instead of leaving a half-created
	// target behind.
	sseDescription, sseErr := s.finaliseSSESpecification(ctx, reqCtx, in.Overrides.SSEDescription, in.Overrides.SSEDisable, sseOpTableCreate)
	if sseErr != nil {
		return nil, sseErr
	}

	return s.runRestorePipeline(ctx, store, restorePipelineInput{
		targetTableName:    in.TargetTableName,
		keySchema:          keySchema,
		attrDefs:           attrDefs,
		schema:             schema,
		sseDescription:     sseDescription,
		onDemandThroughput: in.Overrides.OnDemandThroughput,
		restoreSummary: &dbstore.RestoreSummary{
			SourceBackupArn:   backup.BackupArn,
			SourceTableArn:    backup.SourceTableArn,
			RestoreDateTime:   backup.BackupCreationDateTime,
			RestoreInProgress: false,
		},
		fetchItems: func(ctx context.Context) ([]*dbstore.Item, error) {
			// A failed snapshot read is fatal: an empty backup holds a
			// written empty snapshot record, so an unreadable one is a
			// corrupted or incomplete backup — proceeding would mark an
			// empty table ACTIVE and present it as a completed restore.
			// The error propagates as on the point-in-time restore path.
			return store.Backups().GetSnapshot(backup.BackupArn)
		},
	})
}

// restoreTableToPointInTimeInput carries the raw wire parameters for
// RestoreTableToPointInTime.
type restoreTableToPointInTimeInput struct {
	Parameters map[string]interface{}
}

// restoreTableToPointInTimeCore validates the request, then creates a new
// table by copying the source table's schema and items, applying any
// overrides. The target table must not already exist.
// restoreChunkSize is the number of items written per store.Update()
// transaction during restore. Each chunk is atomic; the table remains
// in CREATING status (invisible to clients) until all chunks complete.
const restoreChunkSize = 500

func (s *DynamoDBService) restoreTableToPointInTimeCore(ctx context.Context, reqCtx *request.RequestContext, in restoreTableToPointInTimeInput) (*dbstore.Table, error) {
	sourceTableName := request.GetStringParam(in.Parameters, "SourceTableName")
	sourceTableArn := request.GetStringParam(in.Parameters, "SourceTableArn")

	if sourceTableName == "" && sourceTableArn == "" {
		return nil, ErrInvalidParameter
	}
	if sourceTableName != "" && sourceTableArn != "" {
		return nil, ErrInvalidParameter
	}

	if sourceTableName != "" {
		if !validateResourceName(sourceTableName) {
			return nil, ErrInvalidParameter
		}
	}

	if sourceTableArn != "" {
		sourceTableName = svcarn.ParseTableARN(sourceTableArn)
		if sourceTableName == "" {
			return nil, ErrResourceNotFound
		}
	}

	targetTableName := request.GetStringParam(in.Parameters, "TargetTableName")
	if !validateResourceName(targetTableName) {
		return nil, ErrInvalidParameter
	}

	// RestoreTableToPointInTime declares TableNotFoundException (rather
	// than the general ResourceNotFoundException) in the Smithy model, so
	// use the individual error sentinel here.
	sourceTable, err := s.validateAndGetTableWithErr(reqCtx, map[string]interface{}{"TableName": sourceTableName}, ErrTableNotFoundException)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	// Point-in-time restore requires recovery to be enabled on the source
	// table; the request must then name a point inside the restorable
	// window.
	pitr, err := store.Tables().GetPointInTimeRecovery(sourceTableName)
	if err != nil {
		return nil, err
	}
	if pitr == nil || pitr.Status != dbstore.PITRStatusEnabled {
		return nil, ErrPITRNotEnabled
	}

	now := time.Now().UTC()
	restoreDateTime, hasRestoreDateTime := parseTimestampParam(in.Parameters, "RestoreDateTime")
	// UseLatestRestorableTime is a boolean member: present-but-mistyped is
	// a request error, never a silently ignored value that falls through
	// to the RestoreDateTime branch.
	useLatest := false
	if raw, present := in.Parameters["UseLatestRestorableTime"]; present {
		b, isBool := raw.(bool)
		if !isBool {
			return nil, NewAPIError("com.amazon.coral.validate#ValidationException",
				"UseLatestRestorableTime must be a boolean", http.StatusBadRequest)
		}
		useLatest = b
	}
	switch {
	case useLatest:
		restoreDateTime = now
	case !hasRestoreDateTime:
		return nil, ErrInvalidParameter
	}
	if restoreDateTime.Before(pitrEarliestRestorable(pitr, now)) || restoreDateTime.After(now) {
		return nil, ErrInvalidRestoreTime
	}

	// Parse the override member family shared with the from-backup restore
	// and overlay it on the source table's settings; every absent member
	// keeps the source's setting.
	overrides, err := parseRestoreOverridesWire(in.Parameters)
	if err != nil {
		return nil, err
	}
	schema, err := applyRestoreOverrides(restoreSchema{
		billingMode:    sourceTable.BillingMode,
		provThroughput: sourceTable.ProvisionedThroughput,
		gsi:            sourceTable.GlobalSecondaryIndexes,
		lsi:            sourceTable.LocalSecondaryIndexes,
		vectorIdx:      sourceTable.VectorIndexes,
	}, overrides)
	if err != nil {
		return nil, err
	}

	if store.Tables().Exists(targetTableName) {
		return nil, ErrTableAlreadyExistsException
	}

	if !validateBillingModeConsistency(schema.billingMode, schema.provThroughput) {
		return nil, ErrInvalidParameter
	}

	// The SSE override resolves before the target exists, the same point the
	// from-backup surface resolves it: an unresolvable specification rejects
	// the request outright instead of leaving a half-created target behind.
	sseDescription, sseErr := s.finaliseSSESpecification(ctx, reqCtx, overrides.SSEDescription, overrides.SSEDisable, sseOpTableCreate)
	if sseErr != nil {
		return nil, sseErr
	}

	return s.runRestorePipeline(ctx, store, restorePipelineInput{
		targetTableName:    targetTableName,
		keySchema:          sourceTable.KeySchema,
		attrDefs:           sourceTable.AttributeDefinitions,
		schema:             schema,
		sseDescription:     sseDescription,
		onDemandThroughput: overrides.OnDemandThroughput,
		restoreSummary: &dbstore.RestoreSummary{
			SourceTableArn:    sourceTable.ARN,
			RestoreDateTime:   restoreDateTime,
			RestoreInProgress: false,
		},
		fetchItems: func(ctx context.Context) ([]*dbstore.Item, error) {
			// Snapshot the source as of the restore point (the current state
			// with every journaled mutation newer than the restore time
			// undone) as the items to write to the target; if any later
			// chunk fails, the partially-populated CREATING table is left
			// for startup recovery.
			snapshot, err := snapshotItemsAsOf(store, sourceTable.Name, restoreDateTime)
			if err != nil {
				return nil, err
			}
			items := make([]*dbstore.Item, 0, len(snapshot))
			for _, item := range snapshot {
				items = append(items, &dbstore.Item{
					TableName:  targetTableName,
					Key:        copyAttributes(item.Key),
					Attributes: copyAttributes(item.Attributes),
				})
			}
			return items, nil
		},
	})
}

// restorePipelineInput carries everything the shared restore pipeline — the
// create→CREATING→copy-chunks→ACTIVE tail both restore surfaces run — needs
// from its caller: the target identity, the resolved schema and key material,
// the members CreateTableParams cannot carry (the SSE description, the
// table-level on-demand pair, the vector index metadata), the restore summary
// the ACTIVE transition stamps, and the item source.
type restorePipelineInput struct {
	targetTableName    string
	keySchema          []*dbstore.KeySchemaElement
	attrDefs           []*dbstore.AttributeDefinition
	schema             restoreSchema
	sseDescription     *dbstore.SSEDescription
	onDemandThroughput *dbstore.OnDemandThroughput
	restoreSummary     *dbstore.RestoreSummary
	fetchItems         func(ctx context.Context) ([]*dbstore.Item, error)
}

// runRestorePipeline creates the target table, downgrades it to CREATING,
// copies the fetched items in atomic chunks, and transitions to ACTIVE with
// the restore summary. Both status transitions run through the locked
// read-modify-write Update: the Create-time record is ACTIVE, and the
// downgrade re-reads it under the record lock so a concurrent counter flush
// or settings write that lands between Create and the transition is
// preserved rather than clobbered by a wholesale write of the stale
// Create-time record — the same discipline the final ACTIVE transition
// already applied for the ItemCount/TableSizeBytes the chunked copy flushes.
func (s *DynamoDBService) runRestorePipeline(ctx context.Context, store dbstore.DynamoDBStoreInterface, in restorePipelineInput) (*dbstore.Table, error) {
	if _, err := store.Tables().Create(dbstore.CreateTableParams{
		Name:                   in.targetTableName,
		KeySchema:              in.keySchema,
		AttributeDefinitions:   in.attrDefs,
		BillingMode:            in.schema.billingMode,
		ProvisionedThroughput:  in.schema.provThroughput,
		GlobalSecondaryIndexes: in.schema.gsi,
		LocalSecondaryIndexes:  in.schema.lsi,
	}); err != nil {
		return nil, err
	}

	// Set the table to CREATING so clients cannot see partial data during
	// the restore; AWS uses the same status transition. The SSE
	// description, the table-level on-demand pair and the vector index
	// metadata ride the same locked transition — they are restore-only
	// members CreateTableParams does not carry, and the vector index ARNs
	// are re-derived for the restored table's ARN (the entries themselves
	// are rebuilt by the same item write path that restores the items).
	if _, err := store.Tables().Update(in.targetTableName, func(table *dbstore.Table) error {
		table.Status = dbstore.TableStatusCreating
		if in.sseDescription != nil {
			table.SSEDescription = in.sseDescription
		}
		if in.onDemandThroughput != nil {
			table.OnDemandThroughput = in.onDemandThroughput
		}
		applyRestoredVectorIndexes(store.Tables().ARNBuilder(), table, in.schema.vectorIdx)
		return nil
	}); err != nil {
		return nil, err
	}

	items, err := in.fetchItems(ctx)
	if err != nil {
		return nil, err
	}
	buffer := make([]*dbstore.Item, 0, restoreChunkSize)
	for _, item := range items {
		buffer = append(buffer, item)
		if len(buffer) >= restoreChunkSize {
			if err := s.flushRestoreChunk(ctx, store, in.targetTableName, buffer); err != nil {
				return nil, err
			}
			buffer = buffer[:0]
		}
	}
	if len(buffer) > 0 {
		if err := s.flushRestoreChunk(ctx, store, in.targetTableName, buffer); err != nil {
			return nil, err
		}
	}

	// All items copied — transition to ACTIVE under the table's record
	// lock, re-reading the record so the ItemCount/TableSizeBytes the
	// chunked restore flushed are preserved in the final write.
	return store.Tables().Update(in.targetTableName, func(table *dbstore.Table) error {
		table.Status = dbstore.TableStatusActive
		table.RestoreSummary = in.restoreSummary
		return nil
	})
}

// flushRestoreChunk writes a batch of items to the target table in a
// single atomic transaction, updating indexes, item count, and table
// size for each item through the canonical item-write composition.
func (s *DynamoDBService) flushRestoreChunk(ctx context.Context, store dbstore.DynamoDBStoreInterface, tableName string, items []*dbstore.Item) error {
	return store.Update(ctx, func(txn *dbstore.DynamoDBTxn) error {
		for _, item := range items {
			if err := txn.StoreItemWrite(tableName, item.Key, item.Attributes, nil, false, 0); err != nil {
				return err
			}
		}
		return nil
	})
}
