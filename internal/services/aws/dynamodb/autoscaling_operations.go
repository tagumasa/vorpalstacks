// Package dynamodb provides DynamoDB service operations for vorpalstacks.
package dynamodb

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// DescribeTableReplicaAutoScaling returns the auto scaling settings for table replicas.
func (s *DynamoDBService) DescribeTableReplicaAutoScaling(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	table, err := s.validateAndGetTable(reqCtx, req.Parameters)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	settings, err := store.Tables().GetAutoScalingSettings(table.Name)
	if err != nil {
		return nil, err
	}

	var replicas []map[string]interface{}
	if settings != nil {
		if replicaList, ok := settings["replicas"].([]interface{}); ok {
			for _, r := range replicaList {
				if rMap, ok := r.(map[string]interface{}); ok {
					replicas = append(replicas, rMap)
				}
			}
		}
	}
	if replicas == nil {
		replicas = []map[string]interface{}{}
	}

	return map[string]interface{}{
		"TableAutoScalingDescription": map[string]interface{}{
			"TableName":   table.Name,
			"TableStatus": string(table.Status),
			"Replicas":    replicas,
		},
	}, nil
}

// UpdateTableReplicaAutoScaling updates the auto scaling settings for table replicas.
func (s *DynamoDBService) UpdateTableReplicaAutoScaling(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	table, err := s.validateAndGetTable(reqCtx, req.Parameters)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	// Parse replica settings updates and build the replica descriptions.
	// AWS auto-scaling uses Application Auto Scaling; we store settings
	// for API round-trip compatibility without implementing the scaling
	// engine itself.
	var replicas []map[string]interface{}

	if globalTableName := request.GetStringParam(req.Parameters, "GlobalTableName"); globalTableName != "" {
		if !validateGlobalTableName(globalTableName) {
			return nil, ErrInvalidParameter
		}
		gt, gtErr := store.GlobalTables().Get(globalTableName)
		if gtErr == nil && gt != nil {
			for _, replica := range gt.ReplicationGroup {
				replicaDesc := map[string]interface{}{
					"RegionName": replica.RegionName,
				}
				replicas = append(replicas, replicaDesc)
			}
		}
	}

	if replicas == nil {
		replicas = []map[string]interface{}{}
	}

	// Parse ReplicaUpdates from the request. Each entry contains a
	// RegionName and per-replica AutoScaling settings. We store these for
	// API round-trip compatibility; the scaling engine itself is not
	// implemented.
	replicaUpdates, hasUpdates := req.Parameters["ReplicaUpdates"].([]interface{})
	if hasUpdates {
		// Build a lookup of existing replica descriptions by region.
		replicaByRegion := make(map[string]map[string]interface{})
		for _, r := range replicas {
			if region, ok := r["RegionName"].(string); ok {
				replicaByRegion[region] = r
			}
		}

		for _, update := range replicaUpdates {
			updateMap, ok := update.(map[string]interface{})
			if !ok {
				return nil, ErrInvalidParameter
			}
			regionName, _ := updateMap["RegionName"].(string)
			if regionName == "" {
				return nil, ErrInvalidParameter
			}

			desc, exists := replicaByRegion[regionName]
			if !exists {
				desc = map[string]interface{}{"RegionName": regionName}
				replicaByRegion[regionName] = desc
				replicas = append(replicas, desc)
			}

			if readAS, ok := updateMap["ReplicaProvisionedReadCapacityAutoScalingUpdate"].(map[string]interface{}); ok {
				desc["ReplicaProvisionedReadCapacityAutoScalingSettings"] = parseAutoScalingSettings(readAS)
			}
			if writeAS, ok := updateMap["ReplicaProvisionedWriteCapacityAutoScalingUpdate"].(map[string]interface{}); ok {
				desc["ReplicaProvisionedWriteCapacityAutoScalingSettings"] = parseAutoScalingSettings(writeAS)
			}
			if gsiUpdates, ok := updateMap["ReplicaGlobalSecondaryIndexUpdates"].([]interface{}); ok {
				desc["ReplicaGlobalSecondaryIndexSettings"] = parseGSIAutoScalingSettings(gsiUpdates)
			}
		}
	}

	// Validate Smithy constraints on top-level AutoScaling fields.
	if policyName := request.GetStringParam(req.Parameters, "AutoScalingPolicyName"); policyName != "" {
		if !validateAutoScalingPolicyName(policyName) {
			return nil, ErrInvalidParameter
		}
	}
	if roleArn := request.GetStringParam(req.Parameters, "AutoScalingRoleArn"); roleArn != "" {
		if !validateAutoScalingRoleArn(roleArn) {
			return nil, ErrInvalidParameter
		}
	}

	settings := map[string]interface{}{
		"replicas": replicas,
	}
	if err := store.Tables().SetAutoScalingSettings(table.Name, settings); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"TableAutoScalingDescription": map[string]interface{}{
			"TableName":   table.Name,
			"TableStatus": string(table.Status),
			"Replicas":    replicas,
		},
	}, nil
}

// parseAutoScalingSettings extracts AutoScaling settings from a request
// parameter map into a response-compatible description map.
func parseAutoScalingSettings(m map[string]interface{}) map[string]interface{} {
	desc := map[string]interface{}{}
	if v, ok := m["MinimumUnits"]; ok {
		desc["MinimumUnits"] = v
	}
	if v, ok := m["MaximumUnits"]; ok {
		desc["MaximumUnits"] = v
	}
	if v, ok := m["AutoScalingDisabled"]; ok {
		desc["AutoScalingDisabled"] = v
	}
	if v, ok := m["AutoScalingRoleArn"]; ok {
		desc["AutoScalingRoleArn"] = v
	}
	if pol, ok := m["ScalingPolicyUpdate"].(map[string]interface{}); ok {
		policy := map[string]interface{}{}
		if name, ok := pol["PolicyName"]; ok {
			policy["PolicyName"] = name
		}
		desc["ScalingPolicies"] = []interface{}{policy}
	}
	return desc
}

// parseGSIAutoScalingSettings extracts per-GSI AutoScaling settings from
// a ReplicaGlobalSecondaryIndexSettingsUpdate list.
func parseGSIAutoScalingSettings(updates []interface{}) []map[string]interface{} {
	var result []map[string]interface{}
	for _, u := range updates {
		uMap, ok := u.(map[string]interface{})
		if !ok {
			continue
		}
		entry := map[string]interface{}{}
		if idx, ok := uMap["IndexName"].(string); ok {
			entry["IndexName"] = idx
		}
		if readAS, ok := uMap["ProvisionedReadCapacityAutoScalingSettingsUpdate"].(map[string]interface{}); ok {
			for k, v := range parseAutoScalingSettings(readAS) {
				entry[k] = v
			}
		}
		result = append(result, entry)
	}
	return result
}
