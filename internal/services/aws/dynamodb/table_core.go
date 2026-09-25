package dynamodb

import (
	"context"
	"fmt"
	"strings"
	"time"

	"vorpalstacks/internal/common/request"
	tagutil "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/core/logs"

	commonstore "vorpalstacks/internal/store/aws/common"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// ---------------------------------------------------------------------------
// Transport-agnostic DTO for CreateTable
// ---------------------------------------------------------------------------

// CreateTableInput carries every field that CreateTable needs, in a format
// independent of the wire protocol. Both the HTTP API handler and the
// admin gRPC handler build this struct and delegate to createTableCore,
// ensuring that validation and table creation follow a single code path.
type CreateTableInput struct {
	TableName              string
	KeySchema              []*dbstore.KeySchemaElement
	AttributeDefinitions   []*dbstore.AttributeDefinition
	BillingMode            dbstore.BillingMode
	ProvisionedThroughput  *dbstore.ProvisionedThroughput
	GlobalSecondaryIndexes []*dbstore.GlobalSecondaryIndex
	LocalSecondaryIndexes  []*dbstore.LocalSecondaryIndex
	VectorIndexes          []*dbstore.VectorIndex
	StreamSpecification    *dbstore.StreamSpecification
	Tags                   []tagutil.Tag
	// DeletionProtectionRaw is the wire member as extracted; the core
	// validates its type (absent nil yields the false default).
	DeletionProtectionRaw interface{}
	WarmThroughput        *dbstore.WarmThroughput
	OnDemandThroughput    *dbstore.OnDemandThroughput
	GlobalTableSourceArn  string
	// SSEDescription and SSEDisable carry the parsed SSESpecification
	// member: the description as parsed, plus whether the specification
	// explicitly asked to disable encryption. The core's SSE finaliser owns
	// what a disable means for the operation.
	SSEDescription *dbstore.SSEDescription
	SSEDisable     bool
	TableClass     string
	// ResourcePolicy carries an inline resource policy document: Set
	// distinguishes a present empty document (rejected by the core) from an
	// absent member; a non-empty document is applied through the resource
	// policy core after the table exists.
	ResourcePolicySet bool
	ResourcePolicy    string
}

// sseOp names the operation an SSESpecification member arrives on: the
// disable policy differs between creation and update.
type sseOp int

const (
	// sseOpTableCreate covers every surface that creates a table carrying
	// a specification (CreateTable, both restores, the import's table
	// creation parameters).
	sseOpTableCreate sseOp = iota
	// sseOpTableUpdate is UpdateTable's specification member.
	sseOpTableUpdate
)

// finaliseSSESpecification is the single owner of the SSESpecification
// policy and key resolution. Encryption at rest has no off state — every
// table is encrypted and the choice is between key types — so a disable
// request means one of two things by operation: on creation it states the
// AWS owned-key default the fresh table already carries, and there are no
// settings to apply; on UpdateTable it is refused. A specification that
// carries settings resolves its key identifier through resolveSSEKeyArn
// here, before anything is created, so an unresolvable key rejects the
// request outright instead of leaving a half-created table behind.
func (s *DynamoDBService) finaliseSSESpecification(ctx context.Context, reqCtx *request.RequestContext, desc *dbstore.SSEDescription, disable bool, op sseOp) (*dbstore.SSEDescription, error) {
	if disable {
		if op == sseOpTableUpdate {
			return nil, ErrInvalidParameter
		}
		return nil, nil
	}
	return s.resolveSSEKeyArn(ctx, reqCtx, desc)
}

// validateKMSKeyIdentifier reports whether an SSE key identifier carries
// one of the forms the key resolver accepts: a KMS key or alias ARN, an
// alias name, or a bare key id — an identifier carrying alias or ARN
// separators that is neither an alias nor an ARN names nothing, and an ARN
// of another service, or a KMS ARN whose resource names no key or alias,
// is not a key identifier at all. The ARN form is decided on the parsed
// ARN, never on the prefix alone.
func validateKMSKeyIdentifier(id string) bool {
	if strings.HasPrefix(id, "arn:") {
		_, service, _, _, resource := svcarn.SplitARN(id)
		return service == "kms" && (strings.HasPrefix(resource, "key/") || strings.HasPrefix(resource, "alias/"))
	}
	if strings.HasPrefix(id, "alias/") {
		return len(id) > len("alias/")
	}
	return id != "" && !strings.ContainsAny(id, "/:")
}

// resolveSSEKeyArn settles an SSE description's key identifier into the
// ARN form SSEDescription.KMSMasterKeyArn documents ("The KMS key ARN used
// for the KMS encryption"). An absent identifier names the managed
// DynamoDB key, whose reserved alias the platform's KMS does not host, so
// it never consults the resolver: it takes the deterministic alias ARN
// whatever is wired — and an explicit naming of that same reserved alias
// selects the same effective key, so it is normalised onto the identical
// path rather than rejected by a resolver that cannot know the reserved
// alias. With the resolver wired, any other explicit identifier resolves
// to the underlying key's ARN — an identifier the KMS store cannot
// resolve is rejected, the same platform extension the SQS and Kinesis
// key members carry. Without wiring, an explicit identifier falls back to
// the deterministic ARN construction: alias-ARN for alias identifiers,
// key-ARN for key identifiers, verbatim for a full ARN.
func (s *DynamoDBService) resolveSSEKeyArn(ctx context.Context, reqCtx *request.RequestContext, desc *dbstore.SSEDescription) (*dbstore.SSEDescription, error) {
	if desc == nil {
		return nil, nil
	}
	identifier := desc.KMSMasterKeyArn
	if identifier == "" || identifier == defaultSSEKeyAlias {
		identifier = defaultSSEKeyAlias
	} else if s.kmsResolver != nil {
		resolved, err := s.kmsResolver.ResolveKeyArn(ctx, reqCtx.GetRegion(), identifier)
		if err != nil {
			return nil, ErrInvalidParameter
		}
		desc.KMSMasterKeyArn = resolved
		return desc, nil
	}
	builder := svcarn.NewARNBuilder(reqCtx.GetAccountID(), reqCtx.GetRegion()).KMS()
	if strings.HasPrefix(identifier, "arn:") {
		desc.KMSMasterKeyArn = identifier
	} else if strings.HasPrefix(identifier, "alias/") {
		desc.KMSMasterKeyArn = builder.Alias(identifier)
	} else {
		desc.KMSMasterKeyArn = builder.Key(identifier)
	}
	return desc, nil
}

// ---------------------------------------------------------------------------
// Core function — single validation + persistence path
// ---------------------------------------------------------------------------

// validateAndGetTableWithErr behaves like validateAndGetTable but lets the
// caller pick the not-found error sentinel. Operations whose Smithy model
// declares TableNotFoundException (rather than the general
// ResourceNotFoundException) must pass ErrTableNotFoundException here so the
// client receives the individual error code.
func (s *DynamoDBService) validateAndGetTableWithErr(reqCtx *request.RequestContext, params map[string]interface{}, notFoundErr *APIError) (*dbstore.Table, error) {
	tableName := request.GetStringParam(params, "TableName")
	if tableName == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	table, err := store.Tables().Get(tableName)
	if err != nil {
		if dbstore.IsTableNotFound(err) {
			return nil, notFoundErr
		}
		return nil, err
	}
	return table, nil
}

// resolveBillingModeForCreation applies the billing-mode contract every
// table-creation entry point shares: an omitted mode defaults to
// provisioned (the documented CreateTable default), the value must be a
// defined mode, and the mode and the throughput pair must be consistent —
// a provisioned request carries a complete, quota-bounded pair and an
// on-demand request carries none. No creation surface resolves its own
// default or invents its own pair: each routes through this one path.
func resolveBillingModeForCreation(billingMode dbstore.BillingMode, throughput *dbstore.ProvisionedThroughput) (dbstore.BillingMode, error) {
	if billingMode == "" {
		billingMode = dbstore.BillingModeProvisioned
	}
	if !validateBillingModeValue(billingMode) {
		return "", ErrInvalidParameter
	}
	if !validateBillingModeConsistency(billingMode, throughput) {
		return "", ErrInvalidParameter
	}
	if throughput != nil && !validateProvisionedThroughputValues(throughput) {
		return "", ErrInvalidParameter
	}
	return billingMode, nil
}

// createTableCore is the single entry point for table creation shared by the
// HTTP API and the admin gRPC handler. It applies defaults, performs all
// validation, creates the table in the store, applies post-create field
// updates, tags the table, and returns the fully-created table. reqCtx is
// nil on the admin plane, whose create input carries no SSESpecification
// member — the SSE finaliser returns before reading a request context when
// no description is present.
func (s *DynamoDBService) createTableCore(ctx context.Context, reqCtx *request.RequestContext, store dbstore.DynamoDBStoreInterface, in CreateTableInput) (*dbstore.Table, error) {
	// 0. The SSE specification resolves first: an unresolvable key or a
	//    policy violation rejects the request before anything is created.
	sseDescription, sseErr := s.finaliseSSESpecification(ctx, reqCtx, in.SSEDescription, in.SSEDisable, sseOpTableCreate)
	if sseErr != nil {
		return nil, sseErr
	}
	deletionProtection, dpErr := validateBoolValue(in.DeletionProtectionRaw, false)
	if dpErr != nil {
		return nil, dpErr
	}
	// An inline resource policy must be a document, not an empty member;
	// the rejection runs before the table is created so an invalid request
	// leaves nothing behind.
	if in.ResourcePolicySet && in.ResourcePolicy == "" {
		return nil, ErrInvalidParameter
	}

	// 1. Table name validation (length 3-255, allowed characters).
	if !validateResourceName(in.TableName) {
		return nil, ErrInvalidParameter
	}

	// 2. Key schema must be non-empty and well-formed.
	if len(in.KeySchema) == 0 {
		return nil, ErrInvalidParameter
	}
	if !validateKeySchema(in.KeySchema) {
		return nil, ErrInvalidParameter
	}

	// 3. Attribute definitions must cover every key attribute.
	if !validateAttributeDefinitions(in.KeySchema, in.AttributeDefinitions) {
		return nil, ErrInvalidParameter
	}

	// 4. Billing mode default + consistency check (the contract every
	// table-creation entry point resolves through).
	resolvedMode, err := resolveBillingModeForCreation(in.BillingMode, in.ProvisionedThroughput)
	if err != nil {
		return nil, err
	}
	in.BillingMode = resolvedMode
	// A GSI that carries throughput settings must carry valid ones
	// (PositiveLongObject minimum 1, same as the table-level check), and the
	// settings must follow the table's billing mode: a provisioned table's
	// index carries them, an on-demand table's index carries none.
	for _, gsi := range in.GlobalSecondaryIndexes {
		if gsi.ProvisionedThroughput != nil && !validateProvisionedThroughputValues(gsi.ProvisionedThroughput) {
			return nil, ErrInvalidParameter
		}
	}
	if !validateIndexThroughputModeConsistency(dbstore.BillingMode(in.BillingMode), in.GlobalSecondaryIndexes) {
		return nil, ErrInvalidParameter
	}
	if !validateSecondaryIndexQuotas(in.GlobalSecondaryIndexes, in.LocalSecondaryIndexes) {
		return nil, ErrInvalidParameter
	}

	// 5. Cross-index validation: every key attribute across table + GSIs +
	//    LSIs must appear in AttributeDefinitions; LSI partition key must
	//    match the table partition key; index names must be unique.
	if !validateAllKeyAttributesInDefs(in.KeySchema, in.GlobalSecondaryIndexes, in.LocalSecondaryIndexes, in.AttributeDefinitions) {
		return nil, ErrInvalidParameter
	}
	if !validateLSIPartitionKey(in.KeySchema, in.LocalSecondaryIndexes) {
		return nil, ErrInvalidParameter
	}
	if err := validateIndexNameUniqueness(in.GlobalSecondaryIndexes, in.LocalSecondaryIndexes, in.VectorIndexes); err != nil {
		return nil, err
	}
	if err := validateVectorAttributeDimensions(in.VectorIndexes); err != nil {
		return nil, err
	}
	if len(in.VectorIndexes) > dbstore.VectorIndexesPerTable {
		return nil, ErrInvalidParameter
	}
	if in.TableClass != "" && !validateTableClassValue(in.TableClass) {
		return nil, ErrInvalidParameter
	}
	if err := validateStreamSpecification(in.StreamSpecification); err != nil {
		return nil, err
	}

	// 6. Persist.
	table, err := store.Tables().Create(dbstore.CreateTableParams{
		Name:                      in.TableName,
		KeySchema:                 in.KeySchema,
		AttributeDefinitions:      in.AttributeDefinitions,
		BillingMode:               in.BillingMode,
		ProvisionedThroughput:     in.ProvisionedThroughput,
		GlobalSecondaryIndexes:    in.GlobalSecondaryIndexes,
		LocalSecondaryIndexes:     in.LocalSecondaryIndexes,
		StreamSpecification:       in.StreamSpecification,
		DeletionProtectionEnabled: deletionProtection,
	})
	if err != nil {
		if dbstore.IsTableAlreadyExists(err) {
			return nil, ErrTableAlreadyExists
		}
		return nil, err
	}

	// 7. Post-create field updates (SSE, warm throughput, on-demand
	//    throughput, global table source ARN, table class).
	needsPersist := false
	if sseDescription != nil {
		table.SSEDescription = sseDescription
		needsPersist = true
	}
	if in.WarmThroughput != nil {
		table.WarmThroughput = in.WarmThroughput
		needsPersist = true
	}
	if in.OnDemandThroughput != nil {
		table.OnDemandThroughput = in.OnDemandThroughput
		needsPersist = true
	}
	if in.GlobalTableSourceArn != "" {
		table.GlobalTableSourceArn = in.GlobalTableSourceArn
		needsPersist = true
	}
	if in.TableClass != "" {
		table.TableClass = dbstore.TableClass(in.TableClass)
		needsPersist = true
	}
	if len(in.VectorIndexes) > 0 {
		arnBuilder := store.Tables().ARNBuilder()
		for _, vi := range in.VectorIndexes {
			vi.IndexArn = arnBuilder.Index(table.Name, vi.IndexName)
		}
		table.VectorIndexes = in.VectorIndexes
		needsPersist = true
	}
	if needsPersist {
		if err := store.Tables().Put(table); err != nil {
			return nil, err
		}
	}

	// 8. Tags — the TagStore is the single tag surface (ListTags and every
	// later TagResource write read and write it); the table record keeps no
	// copy.
	if len(in.Tags) > 0 {
		if err := store.Tables().Tags().Tag(in.TableName, tagutil.ToMap(in.Tags)); err != nil {
			return nil, err
		}
	}

	// 9. An inline resource policy rides through the same core the
	// PutResourcePolicy operation uses — one validation path for the policy
	// document. The admin plane sets no policy member, so its nil request
	// context never reaches this call.
	if in.ResourcePolicy != "" {
		if _, err := s.putResourcePolicyCore(ctx, reqCtx, PutResourcePolicyInput{
			ResourceArn: table.ARN,
			Policy:      in.ResourcePolicy,
		}); err != nil {
			return nil, err
		}
	}

	return table, nil
}

// ---------------------------------------------------------------------------
// Core functions for DeleteTable, DescribeTable, ListTables
// ---------------------------------------------------------------------------

// cascadeDeleteTable is the table-deletion transaction the delete paths
// share: under the table's record lock, in one transaction, it reads the
// table (mapping its absence to the documented not-found sentinel), refuses
// a deletion-protected table, leaves the delete-time system backup a
// recovery-enabled table owes, and sweeps the table's state through the
// cascade. The caller owns everything after the commit — the global-table
// membership consequences and the response. The cascade transaction runs
// under the table's record lock so a concurrent metric flush cannot write
// the record back after the cascade deletes it (the cascade queues no
// metric deltas, so the post-commit flush never re-takes the lock).
func (s *DynamoDBService) cascadeDeleteTable(ctx context.Context, store dbstore.DynamoDBStoreInterface, tableName string) (*dbstore.Table, error) {
	if !validateResourceName(tableName) {
		return nil, ErrInvalidParameter
	}

	var deletedTable *dbstore.Table

	err := store.Tables().WithTableLock(tableName, func() error {
		return store.Update(ctx, func(txn *dbstore.DynamoDBTxn) error {
			table, err := txn.GetTable(tableName)
			if err != nil {
				if dbstore.IsTableNotFound(err) {
					return ErrTableNotFound
				}
				return err
			}
			if table.DeletionProtectionEnabled {
				return ErrTableDeletionProtected
			}
			deletedTable = table
			// A table deleted with recovery enabled leaves the documented
			// system backup behind: the snapshot is written inside this
			// transaction, before the cascade sweeps the items it captures.
			// The naming convention and the 35-day retention are the
			// developer guide's "Delete a table with PITR enabled" contract.
			if table.PointInTimeRecovery != nil && table.PointInTimeRecovery.Status == dbstore.PITRStatusEnabled {
				if err := txn.CreateDeletedTableSystemBackup(table); err != nil {
					return fmt.Errorf("create deleted-table system backup for %s: %w", tableName, err)
				}
			}
			return txn.DeleteTableCascade(tableName)
		})
	})
	if err != nil {
		return nil, err
	}
	return deletedTable, nil
}

// deleteTableCore is the single entry point for table deletion shared by the
// HTTP API and the admin gRPC handler. It validates the table name, checks
// for deletion protection, performs the cascade delete, and returns the
// archived table description — through the shared cascade helper — before
// performing the global-table cleanup a member table's deletion owes. The
// region names the plane the deleted table lived on: the member entry that
// leaves the group is the deleted table's own region.
func (s *DynamoDBService) deleteTableCore(ctx context.Context, store dbstore.DynamoDBStoreInterface, region, tableName string) (*dbstore.Table, error) {
	deletedTable, err := s.cascadeDeleteTable(ctx, store, tableName)
	if err != nil {
		return nil, err
	}

	// The global-table record is account-global metadata in its own storage,
	// outside the regional cascade transaction, so the membership cleanup the
	// cascade owes runs here after the table delete commits. Deleting one
	// member's table removes that region's membership alone — replicas join
	// and leave a group individually, and the group survives for its
	// remaining members — and only the emptied group takes the record with
	// it, the same end state the Delete replica action reaches.
	_, gtErr := store.GlobalTables().Update(tableName, func(gt *dbstore.GlobalTable) error {
		kept := gt.ReplicationGroup[:0]
		for _, r := range gt.ReplicationGroup {
			if r.RegionName != region {
				kept = append(kept, r)
			}
		}
		gt.ReplicationGroup = kept
		return nil
	})
	if gtErr != nil {
		if !commonstore.IsNotFound(gtErr) {
			return nil, gtErr
		}
		// A table outside any global table owes no membership cleanup.
	} else if err := deleteGlobalTableRecordIfEmpty(store, tableName); err != nil {
		return nil, err
	}

	deletedTable.Status = dbstore.TableStatusDeleting
	return deletedTable, nil
}

// replicasForTable renders the Replicas member a TableDescription reports:
// every region of the table's global-table replication group, or nil for a
// standalone table. Replica tables are created synchronously on this
// platform, so each region renders as a live ACTIVE replica unless the
// record says otherwise.
func (s *DynamoDBService) replicasForTable(store dbstore.DynamoDBStoreInterface, tableName string) []interface{} {
	gt, err := store.GlobalTables().Get(tableName)
	if err != nil || gt == nil || len(gt.ReplicationGroup) == 0 {
		return nil
	}
	replicas := make([]interface{}, 0, len(gt.ReplicationGroup))
	for _, r := range gt.ReplicationGroup {
		status := string(r.ReplicaStatus)
		if status == "" {
			status = string(dbstore.ReplicaStatusActive)
		}
		entry := map[string]interface{}{
			"RegionName":    r.RegionName,
			"ReplicaStatus": status,
		}
		if r.KMSMasterKeyId != "" {
			entry["KMSMasterKeyId"] = r.KMSMasterKeyId
		}
		if r.ProvisionedReadCapacityUnits > 0 {
			entry["ProvisionedThroughputOverride"] = map[string]interface{}{
				"ReadCapacityUnits": r.ProvisionedReadCapacityUnits,
			}
		}
		if odt := buildOnDemandThroughputResponse(r.OnDemandThroughputOverride); odt != nil {
			entry["OnDemandThroughputOverride"] = odt
		}
		if len(r.GlobalSecondaryIndexOverrides) > 0 {
			gsis := make([]map[string]interface{}, 0, len(r.GlobalSecondaryIndexOverrides))
			for _, gsi := range r.GlobalSecondaryIndexOverrides {
				gsiEntry := map[string]interface{}{"IndexName": gsi.IndexName}
				if gsi.ProvisionedReadCapacityUnits > 0 {
					gsiEntry["ProvisionedThroughputOverride"] = map[string]interface{}{
						"ReadCapacityUnits": gsi.ProvisionedReadCapacityUnits,
					}
				}
				if odt := buildOnDemandThroughputResponse(gsi.OnDemandThroughputOverride); odt != nil {
					gsiEntry["OnDemandThroughputOverride"] = odt
				}
				gsis = append(gsis, gsiEntry)
			}
			entry["GlobalSecondaryIndexes"] = gsis
		}
		replicas = append(replicas, entry)
	}
	return replicas
}

// describeTableCore is the single entry point for table lookup shared by the
// HTTP API and the admin gRPC handler. It validates the table name and maps
// store errors to DynamoDB API errors.
func (s *DynamoDBService) describeTableCore(store dbstore.DynamoDBStoreInterface, tableName string) (*dbstore.Table, error) {
	if !validateResourceName(tableName) {
		return nil, ErrInvalidParameter
	}

	table, err := store.Tables().Get(tableName)
	if err != nil {
		if dbstore.IsTableNotFound(err) {
			return nil, ErrTableNotFound
		}
		return nil, err
	}
	return table, nil
}

// listTablesCore is the single entry point for table listing shared by the
// HTTP API and the admin gRPC handler. It validates the Limit range (the
// model's ListTablesInputLimit, 1-100) and the marker's table-name shape
// before delegating to the store; an absent limit arrives as the documented
// maximum from both planes.
func (s *DynamoDBService) listTablesCore(store dbstore.DynamoDBStoreInterface, marker string, limit int) ([]*dbstore.Table, string, error) {
	if !validateListTablesLimit(limit) {
		return nil, "", ErrInvalidParameter
	}
	if marker != "" && !validateResourceName(marker) {
		return nil, "", ErrInvalidParameter
	}
	return store.Tables().List(marker, limit)
}

// ---------------------------------------------------------------------------
// Transport-agnostic DTO for UpdateTable
// ---------------------------------------------------------------------------

// UpdateTableInput carries every field that UpdateTable can modify, in a
// format independent of the wire protocol. Fields with zero values are
// treated as "no change". Both the HTTP API handler and the admin gRPC
// handler build this struct and delegate to updateTableCore.
type UpdateTableInput struct {
	TableName             string
	BillingMode           string                         // "" = no change
	ProvisionedThroughput *dbstore.ProvisionedThroughput // nil = no change
	OnDemandThroughput    *dbstore.OnDemandThroughput    // nil = no change
	WarmThroughput        *dbstore.WarmThroughput        // nil = no change
	AttributeDefinitions  []*dbstore.AttributeDefinition // nil = no change
	GSIUpdates            []interface{}                  // nil = no change
	VectorIndexUpdates    []interface{}                  // nil = no change
	StreamSpecification   *dbstore.StreamSpecification   // nil = no change
	SSESpecification      *dbstore.SSEDescription        // nil = no change
	SSEDisable            bool                           // the member explicitly asked to disable encryption
	DeletionProtectionSet bool                           // whether DeletionProtectionEnabled was provided
	DeletionProtectionRaw interface{}                    // the wire member as extracted; validated by the core
	TableClass            string                         // "" = no change
	// ReplicaUpdates carries the modern replica-management actions the model
	// defines on UpdateTable; RequestRegion is the region the request
	// arrived through, the implicit first member when the first Create
	// action converts a regional table into a global table.
	ReplicaUpdates                     []interface{}
	RequestRegion                      string
	MultiRegionConsistency             string // "" = absent
	GlobalTableWitnessUpdates          []interface{}
	GlobalTableSettingsReplicationMode string // "" = absent
}

// windowBillingModeSwitches returns the billing-mode switch stamps inside
// the rolling 24-hour window at now — the stamps older than the window drop
// out before any count is taken.
func windowBillingModeSwitches(stamps []time.Time, now time.Time) []time.Time {
	recent := make([]time.Time, 0, len(stamps))
	for _, ts := range stamps {
		if now.Sub(ts) < 24*time.Hour {
			recent = append(recent, ts)
		}
	}
	return recent
}

// registerBillingModeSwitch applies a billing-mode switch to a table under
// the documented switch limit: only the provisioned-to-on-demand direction
// is rate limited ("You can switch tables from on-demand mode to provisioned
// capacity mode at any time", up to four switches in any rolling 24 hours),
// so a switch beyond the window's quota is refused and an allowed switch is
// recorded with the aged stamps dropped. The mode is then assigned, and
// switching to on-demand discards the provisioned capacity — an on-demand
// table keeps no throughput settings, and the indexes inherit the mode from
// the base table, so their provisioned settings go with it. An empty mode
// records nothing: the request carries no switch. Every surface that moves
// a table between the modes routes through this one helper, so the quota is
// enforced wherever the switch lands, not only where it was validated.
func registerBillingModeSwitch(table *dbstore.Table, newMode dbstore.BillingMode) error {
	if newMode == "" {
		return nil
	}
	if newMode == dbstore.BillingModePayPerRequest && table.BillingMode == dbstore.BillingModeProvisioned {
		now := time.Now().UTC()
		recent := windowBillingModeSwitches(table.BillingModeSwitches, now)
		if len(recent) >= dbstore.BillingModeSwitchesPerDay {
			return ErrInvalidParameter
		}
		table.BillingModeSwitches = append(recent, now)
	}
	table.BillingMode = newMode
	if newMode == dbstore.BillingModePayPerRequest {
		table.ProvisionedThroughput = nil
		for _, gsi := range table.GlobalSecondaryIndexes {
			gsi.ProvisionedThroughput = nil
		}
	}
	return nil
}

// updateTableCore is the single entry point for table updates shared by the
// HTTP API and the admin gRPC handler. It validates all update parameters,
// applies the changes to the table, persists them, and backfills any newly
// created GSIs. The record update runs through the store's locked
// read-modify-write path, so concurrent counter flushes and settings writes
// cannot be lost; the GSI backfill runs after the lock is released because
// its chunk transactions take the same lock for their own metric flushes.
func (s *DynamoDBService) updateTableCore(ctx context.Context, reqCtx *request.RequestContext, store dbstore.DynamoDBStoreInterface, in UpdateTableInput) (*dbstore.Table, error) {
	// The SSE policy runs first: a specification that asks to disable
	// encryption is refused outright, and a carrying one resolves its key
	// before any mutation begins.
	sseSpecification, sseErr := s.finaliseSSESpecification(ctx, reqCtx, in.SSESpecification, in.SSEDisable, sseOpTableUpdate)
	if sseErr != nil {
		return nil, sseErr
	}
	if !validateResourceName(in.TableName) {
		return nil, ErrInvalidParameter
	}
	if in.BillingMode != "" && !validateBillingModeValue(dbstore.BillingMode(in.BillingMode)) {
		return nil, ErrInvalidParameter
	}
	if in.TableClass != "" && !validateTableClassValue(in.TableClass) {
		return nil, ErrInvalidParameter
	}
	if err := validateStreamSpecification(in.StreamSpecification); err != nil {
		return nil, err
	}
	if in.ProvisionedThroughput != nil && !validateProvisionedThroughputValues(in.ProvisionedThroughput) {
		return nil, ErrInvalidParameter
	}
	// The model bounds the replica update list at a minimum of one entry:
	// a present-but-empty list is a malformed request, not an omission.
	if in.ReplicaUpdates != nil && len(in.ReplicaUpdates) == 0 {
		return nil, ErrInvalidParameter
	}

	// The replica-family members the modern model defines on UpdateTable.
	// EVENTUAL consistency alongside Create actions is the platform's own
	// replication mode (the documented default for a new global table), so
	// the member is accepted as the no-op it is here; STRONG consistency,
	// witness replicas, and the settings-replication mode all configure
	// behaviour this replication substrate cannot provide, and are refused
	// rather than accepted without effect. The consistency member is
	// documented as valid only with Create actions in the list.
	if in.MultiRegionConsistency != "" {
		if in.MultiRegionConsistency != multiRegionConsistencyEventual || !replicaUpdatesCarryCreate(in.ReplicaUpdates) {
			return nil, ErrInvalidParameter
		}
	}
	if len(in.GlobalTableWitnessUpdates) > 0 {
		return nil, ErrInvalidParameter
	}
	if in.GlobalTableSettingsReplicationMode != "" {
		return nil, ErrInvalidParameter
	}

	// Replica membership changes land before the table-record update: the
	// cross-region validation is the half most likely to refuse, and the
	// response renders the resulting group either way. The ACTIVE
	// requirement every other UpdateTable member meets inside the record
	// update applies here too — read up front, so a membership-only request
	// on a CREATING/UPDATING/DELETING table is refused rather than
	// mutating the group (a racing status change is still caught by the
	// record update for the members that reach it).
	if len(in.ReplicaUpdates) > 0 {
		current, getErr := store.Tables().Get(in.TableName)
		if getErr != nil {
			if dbstore.IsTableNotFound(getErr) {
				return nil, ErrResourceNotFound
			}
			return nil, getErr
		}
		if current.Status != dbstore.TableStatusActive {
			return nil, ErrTableNotActive
		}
		result, err := s.applyReplicaMembershipUpdates(ctx, store, in.TableName, in.RequestRegion, in.ReplicaUpdates, true)
		if err != nil {
			return nil, err
		}
		for _, region := range result.DeletedRegions {
			if region != in.RequestRegion {
				continue
			}
			// A Delete action naming the requesting region has cascaded
			// this region's table away behind the committed membership
			// change — the requested update is complete. The response
			// renders the departed table's archived description in the
			// DELETING state a direct DeleteTable reports; no table
			// record remains here to update.
			archived := result.DeletedTables[region]
			if archived == nil {
				archived = &dbstore.Table{Name: in.TableName}
			}
			archived.Status = dbstore.TableStatusDeleting
			return archived, nil
		}
	}

	existingGSINames := make(map[string]bool)
	deletedGSINames := []string{}
	vectorCreatedNames := []string{}
	vectorDeletedNames := []string{}
	table, err := store.Tables().Update(in.TableName, func(table *dbstore.Table) error {
		if table.Status != dbstore.TableStatusActive {
			return ErrTableNotActive
		}

		if in.BillingMode != "" {
			// The switch window, the stamp, and the on-demand discard
			// live in the shared helper; the consistency check below
			// then rejects a request that tries to carry
			// ProvisionedThroughput into the new mode.
			if err := registerBillingModeSwitch(table, dbstore.BillingMode(in.BillingMode)); err != nil {
				return err
			}
		}

		if in.ProvisionedThroughput != nil {
			table.ProvisionedThroughput = in.ProvisionedThroughput
		}

		// The on-demand maximum-throughput and warm-throughput pairs apply
		// on the update plane the same way the creation plane applies them
		// — a member the request carries replaces the stored pair.
		if in.OnDemandThroughput != nil {
			table.OnDemandThroughput = in.OnDemandThroughput
		}
		if in.WarmThroughput != nil {
			table.WarmThroughput = in.WarmThroughput
		}

		if !validateBillingModeConsistency(table.BillingMode, table.ProvisionedThroughput) {
			return ErrInvalidParameter
		}

		if len(in.AttributeDefinitions) > 0 {
			if !validateAttributeDefinitions(table.KeySchema, in.AttributeDefinitions) {
				return ErrInvalidParameter
			}
			table.AttributeDefinitions = mergeAttributeDefinitions(table.AttributeDefinitions, in.AttributeDefinitions)
		}

		for _, g := range table.GlobalSecondaryIndexes {
			existingGSINames[g.IndexName] = true
		}

		if len(in.GSIUpdates) > 0 {
			updatedGSIs, deleted, err := applyGSIUpdates(store.Tables().ARNBuilder(), table.Name, table.GlobalSecondaryIndexes, in.GSIUpdates)
			if err != nil {
				return err
			}
			table.GlobalSecondaryIndexes = updatedGSIs
			deletedGSINames = deleted
		}

		if len(in.VectorIndexUpdates) > 0 {
			updatedVector, created, deleted, err := applyVectorIndexUpdates(store.Tables().ARNBuilder(), table.Name, table.VectorIndexes, in.VectorIndexUpdates)
			if err != nil {
				return err
			}
			table.VectorIndexes = updatedVector
			vectorCreatedNames = created
			vectorDeletedNames = deleted
		}

		if !validateAllKeyAttributesInDefs(table.KeySchema, table.GlobalSecondaryIndexes, table.LocalSecondaryIndexes, table.AttributeDefinitions) {
			return ErrInvalidParameter
		}
		// The mode consistency runs on the post-update state, after the
		// index updates land: a switch back to provisioned mode must bring
		// the indexes' provisioned settings with it ("initial provisioned
		// capacity values must be set"), and an index update may not attach
		// provisioned settings to an on-demand table's index. The index
		// quotas run on the same post-update state — a GSI-creating update
		// is bounded by the per-table counts and the projection sum the
		// same way creation is.
		if !validateIndexThroughputModeConsistency(table.BillingMode, table.GlobalSecondaryIndexes) {
			return ErrInvalidParameter
		}
		if !validateSecondaryIndexQuotas(table.GlobalSecondaryIndexes, table.LocalSecondaryIndexes) {
			return ErrInvalidParameter
		}
		if err := validateIndexNameUniqueness(table.GlobalSecondaryIndexes, table.LocalSecondaryIndexes, table.VectorIndexes); err != nil {
			return err
		}
		if err := validateVectorAttributeDimensions(table.VectorIndexes); err != nil {
			return err
		}
		if len(table.VectorIndexes) > dbstore.VectorIndexesPerTable {
			return ErrInvalidParameter
		}

		if in.StreamSpecification != nil {
			table.StreamSpecification = in.StreamSpecification
			if in.StreamSpecification.StreamEnabled {
				streamLabel := time.Now().UTC().Format("2006-01-02T15:04:05.000")
				table.StreamArn = store.Tables().ARNBuilder().TableStream(table.Name, streamLabel)
				table.LatestStreamLabel = streamLabel
			} else {
				table.StreamArn = ""
				table.LatestStreamLabel = ""
			}
		}

		if sseSpecification != nil {
			table.SSEDescription = sseSpecification
		}

		if in.DeletionProtectionSet {
			enabled, dpErr := validateBoolValue(in.DeletionProtectionRaw, false)
			if dpErr != nil {
				return dpErr
			}
			table.DeletionProtectionEnabled = enabled
		}

		if in.TableClass != "" {
			table.TableClass = dbstore.TableClass(in.TableClass)
		}

		table.LastUpdatedDateTime = time.Now().UTC()
		return nil
	})
	if err != nil {
		if dbstore.IsTableNotFound(err) {
			return nil, ErrResourceNotFound
		}
		return nil, err
	}

	// A stream-specification update mints a new stream generation: the
	// re-enabled or re-asserted stream starts empty, so the superseded
	// generation's records leave the space. The sweep runs in its own pass
	// after the locked table update for the same reason the index-entry
	// sweeps below give — TableStore.Update exposes no transaction to the
	// mutate callback. A write that raced the generation change can still
	// commit a record of the superseded generation into the emptied
	// space; the ARN comparison on the stream read path keeps that record
	// from ever being served to the successor generation, so a failure of
	// this sweep leaves unserved residue rather than wrong reads and is
	// logged while the committed update stands.
	if in.StreamSpecification != nil {
		if err := store.Streams().DeleteTableRecords(in.TableName); err != nil {
			logs.Warn("DynamoDB: failed to empty the superseded stream generation's record space",
				logs.String("table", in.TableName),
				logs.Err(err))
		}
	}

	// A deleted index's entries must not outlive the index: a re-created
	// same-name index would inherit stale entries that resolve to items no
	// longer holding the index key attributes. The cleanup runs in its own
	// transaction after the locked table update — TableStore.Update exposes
	// no transaction to the mutate callback, and the storage layer is
	// read-committed, so folding it into the table write would not close
	// the stale-reader window either. Item writes between the two commits
	// write no entries for the deleted index: the table record they read no
	// longer lists it. Both deletion sweeps — this one and the vector twin
	// below — run ahead of the backfills: the metadata deletions have
	// already committed and a re-sent request carries no deletion to
	// re-run, so a backfill failure must not be able to strand the sweeps.
	if len(deletedGSINames) > 0 {
		err := store.Update(ctx, func(txn *dbstore.DynamoDBTxn) error {
			for _, name := range deletedGSINames {
				if err := txn.DeleteIndexEntriesForIndex(table.Name, name); err != nil {
					return fmt.Errorf("delete index entries for %s: %w", name, err)
				}
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	// A deleted vector index's entries must not outlive it — the same
	// lifecycle the GSI sweep above applies, through the vector twin
	// helper, and under the same ordering rule: ahead of every backfill.
	if len(vectorDeletedNames) > 0 {
		err := store.Update(ctx, func(txn *dbstore.DynamoDBTxn) error {
			for _, name := range vectorDeletedNames {
				if err := txn.DeleteVectorEntriesForIndex(table.Name, name); err != nil {
					return fmt.Errorf("delete vector entries for %s: %w", name, err)
				}
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	for _, g := range table.GlobalSecondaryIndexes {
		if existingGSINames[g.IndexName] {
			continue
		}
		if err := s.backfillGSI(ctx, store, table.Name, g.IndexName); err != nil {
			return nil, err
		}
	}

	// A newly created vector index is populated from the existing items —
	// the vector twin of the GSI backfill above, through the same helper.
	for _, name := range vectorCreatedNames {
		if err := s.backfillVectorIndex(ctx, store, table.Name, name); err != nil {
			return nil, err
		}
	}

	return table, nil
}
