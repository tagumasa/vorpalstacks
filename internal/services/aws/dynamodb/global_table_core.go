package dynamodb

import (
	"context"
	"errors"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	commonstore "vorpalstacks/internal/store/aws/common"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// ---------------------------------------------------------------------------
// Global Table Core — single validation + persistence path for the global
// table (2017-style API) operations.
//
// Both the HTTP API handlers (global_table_operations.go) and any future
// admin handler delegate to these methods to ensure identical behaviour.
// ---------------------------------------------------------------------------
//
// This file carries the membership plane: replica validation, the
// Create/Describe/List/Update operations, and the replication-group update
// core the modern UpdateTable ReplicaUpdates member shares. The settings
// operations (Describe/UpdateGlobalTableSettings) with their provisioned
// capacity parse/merge machinery live in global_settings_core.go.

// replicaUpdatesCarryCreate reports whether any action in a
// ReplicationGroupUpdateList is a Create action.
func replicaUpdatesCarryCreate(updates []interface{}) bool {
	for _, update := range updates {
		if updateMap, ok := update.(map[string]interface{}); ok {
			if _, ok := updateMap["Create"].(map[string]interface{}); ok {
				return true
			}
		}
	}
	return false
}

// multiRegionConsistencyEventual is the consistency mode this platform's
// replication provides: asynchronous, multi-active, eventual — the mode the
// model documents as the default for a new global table. The STRONG mode
// has no substrate here and UpdateTable refuses it.
const multiRegionConsistencyEventual = "EVENTUAL"

// globalTableReplicaRegions extracts the replica regions from a parsed
// ReplicationGroup member. A region may appear at most once.
func globalTableReplicaRegions(params map[string]interface{}) ([]string, error) {
	replicationGroupParams, ok := params["ReplicationGroup"].([]interface{})
	if !ok {
		return nil, ErrInvalidParameter
	}
	seen := make(map[string]bool)
	var regions []string
	for _, r := range replicationGroupParams {
		rMap, ok := r.(map[string]interface{})
		if !ok {
			return nil, ErrInvalidParameter
		}
		regionName, ok := rMap["RegionName"].(string)
		if !ok || regionName == "" {
			return nil, ErrInvalidParameter
		}
		if seen[regionName] {
			return nil, ErrInvalidParameter
		}
		seen[regionName] = true
		regions = append(regions, regionName)
	}
	// A global table is a replication relationship between two or more
	// tables ("A global table creates a replication relationship between
	// two or more DynamoDB tables with the same table name in the
	// provided Regions"), so a single-region group is rejected outright.
	if len(regions) < 2 {
		return nil, ErrInvalidParameter
	}
	return regions, nil
}

// gsiKeySchemasEqual reports whether two GSI lists agree on every index
// name and key schema, in order.
func gsiKeySchemasEqual(a, b []*dbstore.GlobalSecondaryIndex) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].IndexName != b[i].IndexName || !keySchemasEqual(a[i].KeySchema, b[i].KeySchema) {
			return false
		}
	}
	return true
}

// lsiKeySchemasEqual reports whether two LSI lists agree on every index
// name and key schema, in order.
func lsiKeySchemasEqual(a, b []*dbstore.LocalSecondaryIndex) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].IndexName != b[i].IndexName || !keySchemasEqual(a[i].KeySchema, b[i].KeySchema) {
			return false
		}
	}
	return true
}

// tableWriteCapacityMatches reports whether two replica tables carry the
// same write capacity settings: equal provisioned write capacity units,
// with an absent (on-demand) throughput record counting as zero. That one
// comparison is complete for the documented joining condition: the billing
// mode invariants (PROVISIONED requires throughput with both units at least
// 1; PAY_PER_REQUEST forbids a throughput record) make every billing-mode
// difference show up as a unit difference, while two on-demand tables both
// compare as zero and two provisioned tables compare their units directly.
func tableWriteCapacityMatches(a, b *dbstore.Table) bool {
	var aUnits, bUnits int64
	if a.ProvisionedThroughput != nil {
		aUnits = a.ProvisionedThroughput.WriteCapacityUnits
	}
	if b.ProvisionedThroughput != nil {
		bUnits = b.ProvisionedThroughput.WriteCapacityUnits
	}
	return aUnits == bUnits
}

// validateGlobalTableReplica enforces the documented conditions a table must
// satisfy to join a global table: a table with the same name as the global
// table exists in the replica region, streams both item images, holds no
// data, and matches the reference replica's key and index schemas and write
// capacity settings. The reference may be nil for the first replica.
func (s *DynamoDBService) validateGlobalTableReplica(globalTableName, region string, reference *dbstore.Table) error {
	replicaStore, err := s.GetStoreForRegion(region)
	if err != nil {
		return ErrTableNotFoundException
	}
	table, err := replicaStore.Tables().Get(globalTableName)
	if err != nil || table == nil {
		return ErrTableNotFoundException
	}
	if table.StreamSpecification == nil || !table.StreamSpecification.StreamEnabled ||
		table.StreamSpecification.StreamViewType != dbstore.StreamViewTypeNewAndOldImages {
		return ErrInvalidParameter
	}
	if table.ItemCount > 0 {
		return ErrInvalidParameter
	}
	// The item count is maintained transactionally; a bounded scan guards
	// against any drift between the counter and stored items.
	hasItems := false
	_, _ = replicaStore.Items().ScanWithOptions(globalTableName, dbstore.ScanOptions{Limit: 1}, func(item *dbstore.Item) error {
		hasItems = true
		return nil
	})
	if hasItems {
		return ErrInvalidParameter
	}
	if reference != nil {
		if !keySchemasEqual(reference.KeySchema, table.KeySchema) ||
			!gsiKeySchemasEqual(reference.GlobalSecondaryIndexes, table.GlobalSecondaryIndexes) ||
			!lsiKeySchemasEqual(reference.LocalSecondaryIndexes, table.LocalSecondaryIndexes) {
			return ErrInvalidParameter
		}
		// The documented joining conditions also require the same
		// provisioned (and maximum) write capacity units as the existing
		// replicas.
		if !tableWriteCapacityMatches(reference, table) {
			return ErrInvalidParameter
		}
	}
	return nil
}

// referenceReplicaTable fetches the table backing the first recorded replica
// region, for schema comparison when a new replica joins.
func (s *DynamoDBService) referenceReplicaTable(globalTableName string, replicas []*dbstore.Replica) *dbstore.Table {
	for _, replica := range replicas {
		store, err := s.GetStoreForRegion(replica.RegionName)
		if err != nil {
			continue
		}
		table, err := store.Tables().Get(globalTableName)
		if err == nil && table != nil {
			return table
		}
	}
	return nil
}

// createGlobalTableInput carries the raw wire parameters for
// CreateGlobalTable.
type createGlobalTableInput struct {
	Parameters map[string]interface{}
}

// createGlobalTableCore validates the request, then creates a new global
// table record over the qualifying replica tables.
func (s *DynamoDBService) createGlobalTableCore(ctx context.Context, reqCtx *request.RequestContext, in createGlobalTableInput) (interface{}, error) {
	globalTableName := request.GetStringParam(in.Parameters, "GlobalTableName")
	if !validateResourceName(globalTableName) {
		return nil, ErrInvalidParameter
	}

	regions, err := globalTableReplicaRegions(in.Parameters)
	if err != nil {
		return nil, err
	}

	// Every replica region must hold a qualifying table before the global
	// table record is created; the first validated table is the schema
	// reference for the remaining regions.
	var reference *dbstore.Table
	for _, region := range regions {
		if err := s.validateGlobalTableReplica(globalTableName, region, reference); err != nil {
			return nil, err
		}
		if reference == nil {
			if store, storeErr := s.GetStoreForRegion(region); storeErr == nil {
				if table, tableErr := store.Tables().Get(globalTableName); tableErr == nil {
					reference = table
				}
			}
		}
	}

	var replicationGroup []*dbstore.Replica
	for _, region := range regions {
		replicationGroup = append(replicationGroup, &dbstore.Replica{
			RegionName:    region,
			ReplicaStatus: "ACTIVE",
		})
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	globalTable, err := store.GlobalTables().Create(globalTableName, replicationGroup)
	if err != nil {
		if dbstore.IsTableAlreadyExists(err) {
			return nil, ErrGlobalTableAlreadyExists
		}
		return nil, err
	}

	return map[string]interface{}{
		"GlobalTableDescription": buildGlobalTableDescription(globalTable),
	}, nil
}

// describeGlobalTableInput carries the raw wire parameters for
// DescribeGlobalTable.
type describeGlobalTableInput struct {
	Parameters map[string]interface{}
}

// describeGlobalTableCore validates the request, then returns the global
// table description.
func (s *DynamoDBService) describeGlobalTableCore(ctx context.Context, reqCtx *request.RequestContext, in describeGlobalTableInput) (interface{}, error) {
	globalTableName := request.GetStringParam(in.Parameters, "GlobalTableName")
	if !validateResourceName(globalTableName) {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	globalTable, err := store.GlobalTables().Get(globalTableName)
	if err != nil {
		// Absence is the documented GlobalTableNotFound; any other store
		// failure (bucket I/O, unmarshal) is reported as the storage error
		// it is, the same split the Update cores apply.
		if commonstore.IsNotFound(err) {
			return nil, ErrGlobalTableNotFound
		}
		return nil, err
	}

	return map[string]interface{}{
		"GlobalTableDescription": buildGlobalTableDescription(globalTable),
	}, nil
}

// listGlobalTablesInput carries the raw wire parameters for
// ListGlobalTables.
type listGlobalTablesInput struct {
	Parameters map[string]interface{}
}

// listGlobalTablesCore validates the request, then lists the global tables,
// optionally filtered to a replica region.
func (s *DynamoDBService) listGlobalTablesCore(ctx context.Context, reqCtx *request.RequestContext, in listGlobalTablesInput) (interface{}, error) {
	regionName := request.GetStringParam(in.Parameters, "RegionName")
	// The Limit member targets PositiveIntegerObject (minimum 1): an
	// explicit zero is a validation failure, while an omitted member keeps
	// the default page size — presence decides, not the zero value.
	limit, limitSet, limitErr := intParamWithPresence(in.Parameters, "Limit")
	if limitErr != nil {
		return nil, limitErr
	}
	if limitSet {
		if !validateListGlobalTablesLimit(limit) {
			return nil, ErrInvalidParameter
		}
	} else {
		limit = listGlobalTablesDefaultLimit
	}
	exclusiveStartGlobalTableName := request.GetStringParam(in.Parameters, "ExclusiveStartGlobalTableName")
	if exclusiveStartGlobalTableName != "" {
		if !validateResourceName(exclusiveStartGlobalTableName) {
			return nil, ErrInvalidParameter
		}
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	var filteredTables []*dbstore.GlobalTable
	marker := exclusiveStartGlobalTableName
	pageSize := limit
	if pageSize < listGlobalTablesMinPageSize {
		pageSize = listGlobalTablesMinPageSize
	}

	for len(filteredTables) < limit {
		page, nextMarker, err := store.GlobalTables().List(marker, pageSize)
		if err != nil {
			return nil, err
		}
		for _, gt := range page {
			if regionName != "" {
				found := false
				for _, r := range gt.ReplicationGroup {
					if r.RegionName == regionName {
						found = true
						break
					}
				}
				if !found {
					continue
				}
			}
			filteredTables = append(filteredTables, gt)
			if len(filteredTables) >= limit {
				break
			}
		}
		if nextMarker == "" {
			break
		}
		marker = nextMarker
	}

	hasMore := len(filteredTables) >= limit

	if len(filteredTables) > limit {
		filteredTables = filteredTables[:limit]
	}

	var globalTableNames []map[string]interface{}
	for _, gt := range filteredTables {
		globalTableNames = append(globalTableNames, map[string]interface{}{
			"GlobalTableName":  gt.GlobalTableName,
			"ReplicationGroup": buildReplicationGroup(gt.ReplicationGroup),
		})
	}

	resp := map[string]interface{}{
		"GlobalTables": globalTableNames,
	}

	if hasMore && len(globalTableNames) > 0 {
		resp["LastEvaluatedGlobalTableName"] = globalTableNames[len(globalTableNames)-1]["GlobalTableName"]
	}

	return resp, nil
}

// checkReplicaTableDeletable refuses a Delete action whose destination
// region's replica table carries deletion protection: the action deletes
// that table, and the protection must hold through the replica path as it
// holds through a direct table deletion. A region whose store the platform
// cannot resolve, or a table already absent, leaves the action to proceed
// on the record alone — the membership removal is then the consistency
// repair for that region.
func (s *DynamoDBService) checkReplicaTableDeletable(tableName, regionName string) error {
	replicaStore, err := s.GetStoreForRegion(regionName)
	if err != nil {
		return nil
	}
	table, err := replicaStore.Tables().Get(tableName)
	if err != nil || table == nil {
		return nil
	}
	if table.DeletionProtectionEnabled {
		return ErrTableDeletionProtected
	}
	return nil
}

// deleteReplicaTable performs the destination-Region table deletion a
// Delete action owes after the record commits: the removed region's table
// goes through the same cascade a direct table deletion takes, and the
// archived table description the cascade returns travels back for the
// response. A table already absent leaves the region consistent with the
// record, so the cascade's not-found is the success it reports here; a
// region whose store cannot be resolved — after the same bounded retry a
// replication delivery takes — cannot serve an orphaned table either, and
// is logged rather than failing the committed membership change; both
// cases report a nil description.
func (s *DynamoDBService) deleteReplicaTable(ctx context.Context, regionName, tableName string) (*dbstore.Table, error) {
	var destStore dbstore.DynamoDBStoreInterface
	resolveErr := withReplicationRetry(ctx, func() error {
		store, err := s.GetStoreForRegion(regionName)
		if err != nil {
			return err
		}
		destStore = store
		return nil
	})
	if resolveErr != nil {
		logs.Warn("DynamoDB Global Tables: failed to resolve the removed replica region's store",
			logs.String("table", tableName),
			logs.String("region", regionName),
			logs.Err(resolveErr))
		return nil, nil
	}
	archived, err := s.cascadeDeleteTable(ctx, destStore, tableName)
	if err != nil && !errors.Is(err, ErrTableNotFound) {
		return nil, err
	}
	return archived, nil
}

// wireReplicaReadOverrideUnits extracts the replica Update action's
// ProvisionedThroughputOverride. The model shapes the member as the
// read-side override alone — a single ReadCapacityUnits of
// PositiveLongObject: write capacity is global across a replication
// group and carries no per-replica value, so the write member the shape
// does not define is not read and the read member must fall inside the
// positive-Long range under the table read-capacity maximum.
func wireReplicaReadOverrideUnits(overrideMap map[string]interface{}) (int64, error) {
	raw, ok := overrideMap["ReadCapacityUnits"]
	if !ok {
		return 0, ErrInvalidParameter
	}
	f, ok := raw.(float64)
	if !ok || f != float64(int64(f)) {
		return 0, ErrInvalidParameter
	}
	rcu := int64(f)
	if rcu < 1 || rcu > dbstore.TableMaxReadCapacityUnits {
		return 0, ErrInvalidParameter
	}
	return rcu, nil
}

// wireReplicaOnDemandOverride extracts the replica Update action's
// OnDemandThroughputOverride — a single MaxReadRequestUnits member of the
// same bounds the table-level pair enforces (a limit of at least the
// minimum or the removal sentinel).
func wireReplicaOnDemandOverride(odtMap map[string]interface{}) (*dbstore.OnDemandThroughput, error) {
	raw, ok := odtMap["MaxReadRequestUnits"]
	if !ok {
		return &dbstore.OnDemandThroughput{}, nil
	}
	f, ok := raw.(float64)
	if !ok || f != float64(int64(f)) {
		return nil, ErrInvalidParameter
	}
	v := int64(f)
	if v != dbstore.OnDemandThroughputRemoveValue && v < dbstore.OnDemandThroughputMinUnits {
		return nil, ErrInvalidParameter
	}
	return &dbstore.OnDemandThroughput{MaxReadRequestUnits: v}, nil
}

// applyReplicaGSIOverrides applies the replica Update action's
// GlobalSecondaryIndexes member to the replica record: every entry names
// an index the reference table carries (an unknown name is the index
// not-found error) and records the entry's read-side provisioned override
// and on-demand maximum override. A repeated index name replaces its
// recorded entry. The reference is part of the contract — an override
// validates against the reference table's indexes — so a record whose
// member tables no longer resolves a reference is rejected here, never
// dereferenced.
func applyReplicaGSIOverrides(replica *dbstore.Replica, gsiList []interface{}, reference *dbstore.Table) error {
	if reference == nil {
		return ErrInvalidParameter
	}
	known := make(map[string]bool, len(reference.GlobalSecondaryIndexes))
	for _, gsi := range reference.GlobalSecondaryIndexes {
		known[gsi.IndexName] = true
	}
	kept := make([]*dbstore.ReplicaGlobalSecondaryIndex, 0, len(replica.GlobalSecondaryIndexOverrides)+len(gsiList))
	for _, entry := range replica.GlobalSecondaryIndexOverrides {
		if known[entry.IndexName] {
			kept = append(kept, entry)
		}
	}
	for _, raw := range gsiList {
		entryMap, ok := raw.(map[string]interface{})
		if !ok {
			return ErrInvalidParameter
		}
		indexName, _ := entryMap["IndexName"].(string)
		if indexName == "" || !known[indexName] {
			return ErrIndexNotFound
		}
		override := &dbstore.ReplicaGlobalSecondaryIndex{IndexName: indexName}
		if ptMap, ok := entryMap["ProvisionedThroughputOverride"].(map[string]interface{}); ok {
			rcu, rcuErr := wireReplicaReadOverrideUnits(ptMap)
			if rcuErr != nil {
				return rcuErr
			}
			override.ProvisionedReadCapacityUnits = rcu
		}
		if odtMap, ok := entryMap["OnDemandThroughputOverride"].(map[string]interface{}); ok {
			odt, odtErr := wireReplicaOnDemandOverride(odtMap)
			if odtErr != nil {
				return odtErr
			}
			override.OnDemandThroughputOverride = odt
		}
		replaced := false
		for i, existing := range kept {
			if existing.IndexName == indexName {
				kept[i] = override
				replaced = true
				break
			}
		}
		if !replaced {
			kept = append(kept, override)
		}
	}
	replica.GlobalSecondaryIndexOverrides = kept
	return nil
}

// applyReplicaMembershipUpdates applies ReplicationGroupUpdate actions
// (Create, Update, Delete) to the table's global-table record under the
// record's lock — the membership core the legacy UpdateGlobalTable operation
// and the modern UpdateTable ReplicaUpdates member share. With
// seedOnMissing, the modern flow's first Create action on a table with no
// global-table record seeds the group with the requesting region's table
// (which must itself satisfy the joining conditions), converting the
// regional table into a global table. The Update action applies the
// per-replica overrides the record carries — the KMS key identifier, the
// provisioned read-throughput override, the on-demand read maximum
// override, the per-index capacity overrides, and the table class
// override.
// The Delete action carries its destination region's table deletion: the
// record commits first — membership is what directs replication — and the
// removed region's table then goes through the same cascade a direct table
// deletion takes, the DeleteTableReplica contract the model states for the
// action. A replication group the Delete actions empty takes the
// global-table record with it, the same end state deleting the final member
// table directly produces.
//
// The result reports the committed record, the regions whose tables the
// Delete actions removed, and those cascades' archived table descriptions
// (nil where a destination table was already absent or its region
// unresolvable) — the caller serving a table-shaped response needs the
// description of a departed requesting-region table.
func (s *DynamoDBService) applyReplicaMembershipUpdates(ctx context.Context, store dbstore.DynamoDBStoreInterface, tableName, requestRegion string, updates []interface{}, seedOnMissing bool) (*replicaMembershipResult, error) {
	if seedOnMissing {
		hasCreate := false
		for _, update := range updates {
			if updateMap, ok := update.(map[string]interface{}); ok {
				if _, ok := updateMap["Create"].(map[string]interface{}); ok {
					hasCreate = true
				}
			}
		}
		if hasCreate {
			existing, err := store.GlobalTables().Get(tableName)
			if err != nil && !commonstore.IsNotFound(err) {
				return nil, err
			}
			if existing == nil {
				if err := s.validateGlobalTableReplica(tableName, requestRegion, nil); err != nil {
					return nil, err
				}
				if _, err := store.GlobalTables().Create(tableName, []*dbstore.Replica{
					{RegionName: requestRegion, ReplicaStatus: "ACTIVE"},
				}); err != nil {
					// A concurrent creator winning the seed race surfaces as
					// the family's AlreadyExists error, the same mapping the
					// CreateGlobalTable path applies to this store error —
					// never a raw storage error.
					if dbstore.IsTableAlreadyExists(err) {
						return nil, ErrGlobalTableAlreadyExists
					}
					return nil, err
				}
			}
		}
	}
	// The Delete actions collect their destination regions here as they
	// validate; the table deletions run after the record commits below.
	var deletedRegions []string

	globalTable, err := store.GlobalTables().Update(tableName, func(gt *dbstore.GlobalTable) error {
		// A joining replica must satisfy the same conditions as one named by
		// CreateGlobalTable, compared against an existing replica's table.
		reference := s.referenceReplicaTable(tableName, gt.ReplicationGroup)
		for _, update := range updates {
			updateMap, ok := update.(map[string]interface{})
			if !ok {
				return ErrInvalidParameter
			}

			if createMap, ok := updateMap["Create"].(map[string]interface{}); ok {
				regionName, _ := createMap["RegionName"].(string)
				if regionName == "" {
					return ErrInvalidParameter
				}
				for _, r := range gt.ReplicationGroup {
					if r.RegionName == regionName {
						return ErrReplicaAlreadyExists
					}
				}
				if err := s.validateGlobalTableReplica(tableName, regionName, reference); err != nil {
					return err
				}
				gt.ReplicationGroup = append(gt.ReplicationGroup, &dbstore.Replica{
					RegionName:    regionName,
					ReplicaStatus: "ACTIVE",
				})
			}

			if updateMap, ok := updateMap["Update"].(map[string]interface{}); ok {
				regionName, _ := updateMap["RegionName"].(string)
				if regionName == "" {
					return ErrInvalidParameter
				}
				var replica *dbstore.Replica
				for _, r := range gt.ReplicationGroup {
					if r.RegionName == regionName {
						replica = r
					}
				}
				if replica == nil {
					return ErrReplicaNotFound
				}
				if kmsKey, present, kmsErr := wireStringMember(updateMap["KMSMasterKeyId"]); kmsErr != nil {
					return kmsErr
				} else if present && kmsKey != "" {
					// The identifier validates to the same key-identifier
					// standard the SSE plane's resolver accepts: a full
					// ARN, an alias, or a bare key id.
					if !validateKMSKeyIdentifier(kmsKey) {
						return ErrInvalidParameter
					}
					replica.KMSMasterKeyId = kmsKey
				}
				if odtMap, ok := updateMap["OnDemandThroughputOverride"].(map[string]interface{}); ok {
					odt, odtErr := wireReplicaOnDemandOverride(odtMap)
					if odtErr != nil {
						return odtErr
					}
					replica.OnDemandThroughputOverride = odt
				}
				if gsiList, ok := updateMap["GlobalSecondaryIndexes"].([]interface{}); ok {
					if gsiErr := applyReplicaGSIOverrides(replica, gsiList, reference); gsiErr != nil {
						return gsiErr
					}
				}
				if overrideMap, ok := updateMap["ProvisionedThroughputOverride"].(map[string]interface{}); ok {
					rcu, rcuErr := wireReplicaReadOverrideUnits(overrideMap)
					if rcuErr != nil {
						return rcuErr
					}
					replica.ProvisionedReadCapacityUnits = rcu
				}
				if tableClass, _ := updateMap["TableClassOverride"].(string); tableClass != "" {
					if !validateTableClassValue(tableClass) {
						return ErrInvalidParameter
					}
					replica.TableClass = dbstore.TableClass(tableClass)
					now := time.Now().UTC()
					replica.TableClassLastUpdated = &now
				}
			}

			if deleteMap, ok := updateMap["Delete"].(map[string]interface{}); ok {
				regionName, _ := deleteMap["RegionName"].(string)
				if regionName == "" {
					return ErrInvalidParameter
				}
				found := false
				var newReplicas []*dbstore.Replica
				for _, r := range gt.ReplicationGroup {
					if r.RegionName == regionName {
						found = true
						continue
					}
					newReplicas = append(newReplicas, r)
				}
				if !found {
					return ErrReplicaNotFound
				}
				// The action deletes the destination region's table, so a
				// deletion-protected replica table refuses it here — before
				// the record commits — keeping the membership and the
				// surviving table from disagreeing.
				if err := s.checkReplicaTableDeletable(tableName, regionName); err != nil {
					return err
				}
				gt.ReplicationGroup = newReplicas
				deletedRegions = append(deletedRegions, regionName)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Each Delete action's destination-region table deletion runs after the
	// record commits: membership is what directs replication, so the record
	// goes first and the removed region's table follows. A cascade failure
	// is non-fatal — the committed membership is the truth, and erroring
	// the response would hand the caller an unrecoverable retry (the
	// member is already gone from the record, so a retry answers
	// ReplicaNotFound) — so the failure is logged naming the region and
	// the pass continues, the settings-propagation contract.
	deletedTables := make(map[string]*dbstore.Table, len(deletedRegions))
	for _, regionName := range deletedRegions {
		archived, err := s.deleteReplicaTable(ctx, regionName, tableName)
		if err != nil {
			logs.Error("DynamoDB Global Tables: destination-region table deletion failed after the membership commit",
				logs.String("table", tableName),
				logs.String("region", regionName),
				logs.Err(err))
			continue
		}
		deletedTables[regionName] = archived
	}

	// The final member's departure takes the global-table record with it —
	// the same end state deleting that last member table directly produces,
	// where the record's own delete follows the cascade.
	if err := deleteGlobalTableRecordIfEmpty(store, tableName); err != nil {
		return nil, err
	}
	return &replicaMembershipResult{
		GlobalTable:    globalTable,
		DeletedRegions: deletedRegions,
		DeletedTables:  deletedTables,
	}, nil
}

// replicaMembershipResult reports a membership pass's outcome to its
// callers: the committed record, the regions whose tables the Delete
// actions removed, and those cascades' archived table descriptions.
type replicaMembershipResult struct {
	GlobalTable    *dbstore.GlobalTable
	DeletedRegions []string
	DeletedTables  map[string]*dbstore.Table
}

// deleteGlobalTableRecordIfEmpty takes the global-table record away when
// the replication group has no members left. Both departure paths reach
// this end state — the Delete replica action emptying the group, and a
// direct delete of the final member table — so they share the step. The
// emptiness check and the delete run atomically under the record lock:
// a member re-added concurrently with the departure must never find its
// record erased.
func deleteGlobalTableRecordIfEmpty(store dbstore.DynamoDBStoreInterface, tableName string) error {
	_, err := store.GlobalTables().DeleteIfEmpty(tableName)
	return err
}

// updateGlobalTableInput carries the raw wire parameters for
// UpdateGlobalTable.
type updateGlobalTableInput struct {
	Parameters map[string]interface{}
}

// updateGlobalTableCore validates the request, then applies replica
// membership updates to the global table.
func (s *DynamoDBService) updateGlobalTableCore(ctx context.Context, reqCtx *request.RequestContext, in updateGlobalTableInput) (interface{}, error) {
	globalTableName := request.GetStringParam(in.Parameters, "GlobalTableName")
	if !validateResourceName(globalTableName) {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	updates, ok := in.Parameters["ReplicaUpdates"].([]interface{})
	if !ok {
		return nil, ErrInvalidParameter
	}
	result, err := s.applyReplicaMembershipUpdates(ctx, store, globalTableName, "", updates, false)
	if err != nil {
		if commonstore.IsNotFound(err) {
			return nil, ErrGlobalTableNotFound
		}
		return nil, err
	}

	return map[string]interface{}{
		"GlobalTableDescription": buildGlobalTableDescription(result.GlobalTable),
	}, nil
}

func buildGlobalTableDescription(gt *dbstore.GlobalTable) map[string]interface{} {
	return map[string]interface{}{
		"GlobalTableName":   gt.GlobalTableName,
		"GlobalTableArn":    gt.GlobalTableArn,
		"GlobalTableStatus": string(gt.GlobalTableStatus),
		"CreationDateTime":  gt.CreationDateTime.Unix(),
		"ReplicationGroup":  buildReplicationGroup(gt.ReplicationGroup),
	}
}

// buildReplicationGroup renders the record's replicas in the
// ReplicaDescription shape: the recorded per-replica settings — the KMS
// key identifier, the read-throughput overrides, the per-index capacity
// overrides, and the table class override — are part of the description
// every response face echoes, each unset member omitted per the
// protocol's unset-member rule.
func buildReplicationGroup(replicas []*dbstore.Replica) []map[string]interface{} {
	var result []map[string]interface{}
	for _, r := range replicas {
		entry := map[string]interface{}{
			"RegionName":    r.RegionName,
			"ReplicaStatus": string(r.ReplicaStatus),
		}
		if r.KMSMasterKeyId != "" {
			entry["KMSMasterKeyId"] = r.KMSMasterKeyId
		}
		if r.ProvisionedReadCapacityUnits > 0 {
			entry["ProvisionedThroughputOverride"] = map[string]interface{}{
				"ReadCapacityUnits": r.ProvisionedReadCapacityUnits,
			}
		}
		if r.OnDemandThroughputOverride != nil {
			entry["OnDemandThroughputOverride"] = map[string]interface{}{
				"MaxReadRequestUnits": r.OnDemandThroughputOverride.MaxReadRequestUnits,
			}
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
				if gsi.OnDemandThroughputOverride != nil {
					gsiEntry["OnDemandThroughputOverride"] = map[string]interface{}{
						"MaxReadRequestUnits": gsi.OnDemandThroughputOverride.MaxReadRequestUnits,
					}
				}
				gsis = append(gsis, gsiEntry)
			}
			entry["GlobalSecondaryIndexes"] = gsis
		}
		if r.TableClass != "" {
			summary := map[string]interface{}{"TableClass": string(r.TableClass)}
			if r.TableClassLastUpdated != nil {
				summary["LastUpdateDateTime"] = *r.TableClassLastUpdated
			}
			entry["ReplicaTableClassSummary"] = summary
		}
		result = append(result, entry)
	}
	return result
}
