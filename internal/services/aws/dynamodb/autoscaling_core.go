package dynamodb

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// ---------------------------------------------------------------------------
// Autoscaling Core — single validation + persistence path for the table
// replica AutoScaling operations.
//
// AWS auto-scaling uses Application Auto Scaling; the settings are stored
// for API round-trip compatibility without implementing the scaling engine
// itself. Both the HTTP API handlers (autoscaling_operations.go) and any
// future admin handler delegate to these methods to ensure identical
// behaviour.
// ---------------------------------------------------------------------------

// describeTableReplicaAutoScalingInput carries the raw wire parameters for
// DescribeTableReplicaAutoScaling.
type describeTableReplicaAutoScalingInput struct {
	Parameters map[string]interface{}
}

// describeTableReplicaAutoScalingCore returns the stored AutoScaling
// settings description of the named table.
func (s *DynamoDBService) describeTableReplicaAutoScalingCore(ctx context.Context, reqCtx *request.RequestContext, in describeTableReplicaAutoScalingInput) (interface{}, error) {
	table, err := s.validateAndGetTable(reqCtx, in.Parameters)
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

// updateTableReplicaAutoScalingInput carries the raw wire parameters for
// UpdateTableReplicaAutoScaling.
type updateTableReplicaAutoScalingInput struct {
	Parameters map[string]interface{}
}

// updateTableReplicaAutoScalingCore applies replica AutoScaling settings
// updates to the named table.
func (s *DynamoDBService) updateTableReplicaAutoScalingCore(ctx context.Context, reqCtx *request.RequestContext, in updateTableReplicaAutoScalingInput) (interface{}, error) {
	table, err := s.validateAndGetTable(reqCtx, in.Parameters)
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
	replicas := []map[string]interface{}{}

	// Parse ReplicaUpdates from the request. Each entry contains a
	// RegionName and per-replica AutoScaling settings. We store these for
	// API round-trip compatibility; the scaling engine itself is not
	// implemented.
	replicaUpdates, hasUpdates := in.Parameters["ReplicaUpdates"].([]interface{})
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
				settings, err := parseAutoScalingSettings(readAS)
				if err != nil {
					return nil, err
				}
				desc["ReplicaProvisionedReadCapacityAutoScalingSettings"] = settings
			}
			if writeAS, ok := updateMap["ReplicaProvisionedWriteCapacityAutoScalingUpdate"].(map[string]interface{}); ok {
				settings, err := parseAutoScalingSettings(writeAS)
				if err != nil {
					return nil, err
				}
				desc["ReplicaProvisionedWriteCapacityAutoScalingSettings"] = settings
			}
			if gsiUpdates, ok := updateMap["ReplicaGlobalSecondaryIndexUpdates"].([]interface{}); ok {
				gsiSettings, err := parseGSIAutoScalingSettings(gsiUpdates)
				if err != nil {
					return nil, err
				}
				desc["ReplicaGlobalSecondaryIndexSettings"] = gsiSettings
			}
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
// parameter map into a response-compatible description map, validating the
// Smithy member constraints on the values it carries. Returns
// ErrInvalidParameter when a member violates its documented bounds.
func parseAutoScalingSettings(m map[string]interface{}) (map[string]interface{}, error) {
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
		if roleArn, isStr := v.(string); isStr && !validateAutoScalingRoleArn(roleArn) {
			return nil, ErrInvalidParameter
		}
		desc["AutoScalingRoleArn"] = v
	}
	if pol, ok := m["ScalingPolicyUpdate"].(map[string]interface{}); ok {
		policy := map[string]interface{}{}
		if name, ok := pol["PolicyName"]; ok {
			if policyName, isStr := name.(string); isStr && !validateAutoScalingPolicyName(policyName) {
				return nil, ErrInvalidParameter
			}
			policy["PolicyName"] = name
		}
		desc["ScalingPolicies"] = []interface{}{policy}
	}
	return desc, nil
}

// parseGSIAutoScalingSettings extracts per-GSI AutoScaling settings from
// a ReplicaGlobalSecondaryIndexSettingsUpdate list.
func parseGSIAutoScalingSettings(updates []interface{}) ([]map[string]interface{}, error) {
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
			settings, err := parseAutoScalingSettings(readAS)
			if err != nil {
				return nil, err
			}
			for k, v := range settings {
				entry[k] = v
			}
		}
		result = append(result, entry)
	}
	return result, nil
}
