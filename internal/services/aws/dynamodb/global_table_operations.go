// Package dynamodb provides DynamoDB service operations for vorpalstacks.
package dynamodb

import (
	"context"

	"vorpalstacks/internal/common/request"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

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
	// AWS requires at least one replica in the replication group.
	if len(regions) == 0 {
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

// validateGlobalTableReplica enforces the documented conditions a table must
// satisfy to join a global table: a table with the same name as the global
// table exists in the replica region, streams both item images, holds no
// data, and matches the reference replica's key and index schemas. The
// reference may be nil for the first replica.
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

// CreateGlobalTable creates a new global table.
func (s *DynamoDBService) CreateGlobalTable(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	globalTableName := request.GetStringParam(req.Parameters, "GlobalTableName")
	if !validateGlobalTableName(globalTableName) {
		return nil, ErrInvalidParameter
	}

	regions, err := globalTableReplicaRegions(req.Parameters)
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

// DescribeGlobalTable returns information about a global table.
func (s *DynamoDBService) DescribeGlobalTable(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	globalTableName := request.GetStringParam(req.Parameters, "GlobalTableName")
	if !validateGlobalTableName(globalTableName) {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	globalTable, err := store.GlobalTables().Get(globalTableName)
	if err != nil {
		return nil, ErrGlobalTableNotFound
	}

	return map[string]interface{}{
		"GlobalTableDescription": buildGlobalTableDescription(globalTable),
	}, nil
}

// DescribeGlobalTableSettings returns the settings of a global table.
func (s *DynamoDBService) DescribeGlobalTableSettings(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	globalTableName := request.GetStringParam(req.Parameters, "GlobalTableName")
	if !validateGlobalTableName(globalTableName) {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	globalTable, err := store.GlobalTables().Get(globalTableName)
	if err != nil {
		return nil, ErrGlobalTableNotFound
	}

	return map[string]interface{}{
		"GlobalTableName": globalTable.GlobalTableName,
		"ReplicaSettings": buildGlobalTableReplicaSettings(globalTable),
	}, nil
}

// buildGlobalTableReplicaSettings renders the per-replica settings list
// shared by the global table settings operations.
func buildGlobalTableReplicaSettings(globalTable *dbstore.GlobalTable) []map[string]interface{} {
	replicaSettings := make([]map[string]interface{}, 0, len(globalTable.ReplicationGroup))
	for _, replica := range globalTable.ReplicationGroup {
		settings := map[string]interface{}{
			"RegionName":                           replica.RegionName,
			"ReplicaStatus":                        replica.ReplicaStatus,
			"ReplicaProvisionedReadCapacityUnits":  replica.ProvisionedReadCapacityUnits,
			"ReplicaProvisionedWriteCapacityUnits": replica.ProvisionedWriteCapacityUnits,
		}
		if replica.BillingMode != "" {
			settings["ReplicaBillingModeSummary"] = map[string]interface{}{
				"BillingMode": replica.BillingMode,
			}
		}
		replicaSettings = append(replicaSettings, settings)
	}
	return replicaSettings
}

// ListGlobalTables lists the global tables for a given account.
func (s *DynamoDBService) ListGlobalTables(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	regionName := request.GetStringParam(req.Parameters, "RegionName")
	limit := request.GetIntParam(req.Parameters, "Limit")
	if limit == 0 {
		limit = listGlobalTablesDefaultLimit
	} else {
		if !validateListGlobalTablesLimit(limit) {
			return nil, ErrInvalidParameter
		}
	}
	exclusiveStartGlobalTableName := request.GetStringParam(req.Parameters, "ExclusiveStartGlobalTableName")
	if exclusiveStartGlobalTableName != "" {
		if !validateGlobalTableName(exclusiveStartGlobalTableName) {
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

// UpdateGlobalTable updates a global table.
func (s *DynamoDBService) UpdateGlobalTable(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	globalTableName := request.GetStringParam(req.Parameters, "GlobalTableName")
	if !validateGlobalTableName(globalTableName) {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	globalTable, err := store.GlobalTables().Get(globalTableName)
	if err != nil {
		return nil, ErrGlobalTableNotFound
	}

	updates, ok := req.Parameters["ReplicaUpdates"].([]interface{})
	if !ok {
		return nil, ErrInvalidParameter
	}
	// A joining replica must satisfy the same conditions as one named by
	// CreateGlobalTable, compared against an existing replica's table.
	reference := s.referenceReplicaTable(globalTableName, globalTable.ReplicationGroup)
	for _, update := range updates {
		updateMap, ok := update.(map[string]interface{})
		if !ok {
			return nil, ErrInvalidParameter
		}

		if createMap, ok := updateMap["Create"].(map[string]interface{}); ok {
			regionName, _ := createMap["RegionName"].(string)
			if regionName == "" {
				return nil, ErrInvalidParameter
			}
			for _, r := range globalTable.ReplicationGroup {
				if r.RegionName == regionName {
					return nil, ErrReplicaAlreadyExists
				}
			}
			if err := s.validateGlobalTableReplica(globalTableName, regionName, reference); err != nil {
				return nil, err
			}
			globalTable.ReplicationGroup = append(globalTable.ReplicationGroup, &dbstore.Replica{
				RegionName:    regionName,
				ReplicaStatus: "ACTIVE",
			})
		}

		if deleteMap, ok := updateMap["Delete"].(map[string]interface{}); ok {
			regionName, _ := deleteMap["RegionName"].(string)
			if regionName == "" {
				return nil, ErrInvalidParameter
			}
			found := false
			var newReplicas []*dbstore.Replica
			for _, r := range globalTable.ReplicationGroup {
				if r.RegionName == regionName {
					found = true
					continue
				}
				newReplicas = append(newReplicas, r)
			}
			if !found {
				return nil, ErrReplicaNotFound
			}
			globalTable.ReplicationGroup = newReplicas
		}
	}

	if err := store.GlobalTables().Put(globalTable); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"GlobalTableDescription": buildGlobalTableDescription(globalTable),
	}, nil
}

// UpdateGlobalTableSettings updates the settings of a global table.
func (s *DynamoDBService) UpdateGlobalTableSettings(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	globalTableName := request.GetStringParam(req.Parameters, "GlobalTableName")
	if !validateGlobalTableName(globalTableName) {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	globalTable, err := store.GlobalTables().Get(globalTableName)
	if err != nil {
		return nil, ErrGlobalTableNotFound
	}

	replicaSettingsUpdates, ok := req.Parameters["ReplicaSettingsUpdate"].([]interface{})
	if ok {
		for _, update := range replicaSettingsUpdates {
			updateMap, ok := update.(map[string]interface{})
			if !ok {
				continue
			}

			regionName, _ := updateMap["RegionName"].(string)
			if regionName == "" {
				return nil, ErrInvalidParameter
			}

			for _, replica := range globalTable.ReplicationGroup {
				if replica.RegionName == regionName {
					if billingMode, ok := updateMap["ReplicaBillingMode"].(string); ok {
						replica.BillingMode = billingMode
					}
					if readUnits, ok := updateMap["ReplicaProvisionedReadCapacityUnits"].(float64); ok {
						replica.ProvisionedReadCapacityUnits = int64(readUnits)
					}
					if writeUnits, ok := updateMap["ReplicaProvisionedWriteCapacityUnits"].(float64); ok {
						replica.ProvisionedWriteCapacityUnits = int64(writeUnits)
					}
					break
				}
			}
		}

		if err := store.GlobalTables().Put(globalTable); err != nil {
			return nil, err
		}
	}

	return map[string]interface{}{
		"GlobalTableName": globalTable.GlobalTableName,
		"ReplicaSettings": buildGlobalTableReplicaSettings(globalTable),
	}, nil
}

func buildGlobalTableDescription(gt *dbstore.GlobalTable) map[string]interface{} {
	return map[string]interface{}{
		"GlobalTableName":   gt.GlobalTableName,
		"GlobalTableArn":    gt.GlobalTableArn,
		"GlobalTableStatus": gt.GlobalTableStatus,
		"CreationDateTime":  gt.CreationDateTime.Unix(),
		"ReplicationGroup":  buildReplicationGroup(gt.ReplicationGroup),
	}
}

func buildReplicationGroup(replicas []*dbstore.Replica) []map[string]interface{} {
	var result []map[string]interface{}
	for _, r := range replicas {
		result = append(result, map[string]interface{}{
			"RegionName":    r.RegionName,
			"ReplicaStatus": r.ReplicaStatus,
		})
	}
	return result
}
