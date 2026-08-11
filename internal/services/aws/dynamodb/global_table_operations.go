// Package dynamodb provides DynamoDB service operations for vorpalstacks.
package dynamodb

import (
	"context"
	"vorpalstacks/internal/common/request"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// CreateGlobalTable creates a new global table.
func (s *DynamoDBService) CreateGlobalTable(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	globalTableName := request.GetStringParam(req.Parameters, "GlobalTableName")
	if err := validateGlobalTableName(globalTableName); err != nil {
		return nil, err
	}

	replicationGroupParams, ok := req.Parameters["ReplicationGroup"].([]interface{})
	if !ok {
		return nil, ErrInvalidParameter
	}

	var replicationGroup []*dbstore.Replica
	for _, r := range replicationGroupParams {
		rMap, ok := r.(map[string]interface{})
		if !ok {
			continue
		}
		regionName, ok := rMap["RegionName"].(string)
		if !ok {
			return nil, ErrInvalidParameter
		}
		if regionName != "" {
			replicationGroup = append(replicationGroup, &dbstore.Replica{
				RegionName:    regionName,
				ReplicaStatus: "ACTIVE",
			})
		}
	}

	// AWS requires at least one replica in the replication group
	if len(replicationGroup) == 0 {
		return nil, ErrInvalidParameter
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
	if err := validateGlobalTableName(globalTableName); err != nil {
		return nil, err
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
	if err := validateGlobalTableName(globalTableName); err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	globalTable, err := store.GlobalTables().Get(globalTableName)
	if err != nil {
		return nil, ErrGlobalTableNotFound
	}

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

	return map[string]interface{}{
		"GlobalTableName": globalTable.GlobalTableName,
		"ReplicaSettings": replicaSettings,
	}, nil
}

// ListGlobalTables lists the global tables for a given account.
func (s *DynamoDBService) ListGlobalTables(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	regionName := request.GetStringParam(req.Parameters, "RegionName")
	limit := request.GetIntParam(req.Parameters, "Limit")
	if limit == 0 {
		limit = 100
	} else {
		if err := validateListGlobalTablesLimit(limit); err != nil {
			return nil, err
		}
	}
	exclusiveStartGlobalTableName := request.GetStringParam(req.Parameters, "ExclusiveStartGlobalTableName")
	if exclusiveStartGlobalTableName != "" {
		if err := validateGlobalTableName(exclusiveStartGlobalTableName); err != nil {
			return nil, err
		}
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	var filteredTables []*dbstore.GlobalTable
	marker := exclusiveStartGlobalTableName
	pageSize := limit
	if pageSize < 10 {
		pageSize = 10
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
	if err := validateGlobalTableName(globalTableName); err != nil {
		return nil, err
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
	if ok {
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
	}

	return map[string]interface{}{
		"GlobalTableDescription": buildGlobalTableDescription(globalTable),
	}, nil
}

// UpdateGlobalTableSettings updates the settings of a global table.
func (s *DynamoDBService) UpdateGlobalTableSettings(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	globalTableName := request.GetStringParam(req.Parameters, "GlobalTableName")
	if err := validateGlobalTableName(globalTableName); err != nil {
		return nil, err
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
		"GlobalTableDescription": buildGlobalTableDescription(globalTable),
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
