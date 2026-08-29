package dynamodb

import (
	"context"
	"sort"

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
// updates to the named table. The write-capacity members are table-level in
// the model — they apply to every replica description — while the
// read-capacity members are carried on the per-replica updates.
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

	tableWriteAS, err := parseOptionalAutoScalingSettings(in.Parameters, "ProvisionedWriteCapacityAutoScalingUpdate")
	if err != nil {
		return nil, err
	}
	tableGSIWriteAS, err := parseGSIAutoScalingWrite(in.Parameters["GlobalSecondaryIndexUpdates"])
	if err != nil {
		return nil, err
	}

	// Parse ReplicaUpdates from the request. Each entry contains a
	// RegionName and per-replica read AutoScaling settings. We store these
	// for API round-trip compatibility; the scaling engine itself is not
	// implemented.
	replicaUpdates, hasUpdates := in.Parameters["ReplicaUpdates"].([]interface{})
	if hasUpdates {
		// Deduplicate replica descriptions by region so repeated updates
		// for one region merge into a single description.
		replicaByRegion := make(map[string]map[string]interface{})

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
			if tableWriteAS != nil {
				desc["ReplicaProvisionedWriteCapacityAutoScalingSettings"] = tableWriteAS
			}
			gsiReadAS, err := parseGSIAutoScalingRead(updateMap["ReplicaGlobalSecondaryIndexUpdates"])
			if err != nil {
				return nil, err
			}
			if gsi := mergeReplicaGSIAutoScaling(tableGSIWriteAS, gsiReadAS); len(gsi) > 0 {
				desc["GlobalSecondaryIndexes"] = gsi
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
		// The description form carries the target tracking configuration
		// alongside the policy name; its TargetValue member is required.
		if tt, ok := pol["TargetTrackingScalingPolicyConfiguration"].(map[string]interface{}); ok {
			if target, ok := tt["TargetValue"].(float64); !ok {
				return nil, ErrInvalidParameter
			} else if target <= 0 {
				return nil, ErrInvalidParameter
			}
			ttDesc := map[string]interface{}{}
			for _, member := range []string{"DisableScaleIn", "ScaleInCooldown", "ScaleOutCooldown", "TargetValue"} {
				if v, ok := tt[member]; ok {
					ttDesc[member] = v
				}
			}
			policy["TargetTrackingScalingPolicyConfiguration"] = ttDesc
		}
		desc["ScalingPolicies"] = []interface{}{policy}
	}
	return desc, nil
}

// parseOptionalAutoScalingSettings parses an optional top-level
// AutoScalingSettingsUpdate member, returning nil when the member is absent.
func parseOptionalAutoScalingSettings(parameters map[string]interface{}, key string) (map[string]interface{}, error) {
	if m, ok := parameters[key].(map[string]interface{}); ok {
		return parseAutoScalingSettings(m)
	}
	return nil, nil
}

// parseGSIAutoScalingWrite parses a GlobalSecondaryIndexUpdates list into
// per-index write-capacity AutoScaling descriptions keyed by index name.
func parseGSIAutoScalingWrite(updates interface{}) (map[string]map[string]interface{}, error) {
	result := map[string]map[string]interface{}{}
	gsiUpdates, ok := updates.([]interface{})
	if !ok {
		return result, nil
	}
	for _, u := range gsiUpdates {
		uMap, ok := u.(map[string]interface{})
		if !ok {
			continue
		}
		indexName, _ := uMap["IndexName"].(string)
		if indexName == "" {
			continue
		}
		if writeAS, ok := uMap["ProvisionedWriteCapacityAutoScalingUpdate"].(map[string]interface{}); ok {
			settings, err := parseAutoScalingSettings(writeAS)
			if err != nil {
				return nil, err
			}
			result[indexName] = settings
		}
	}
	return result, nil
}

// parseGSIAutoScalingRead parses a ReplicaGlobalSecondaryIndexUpdates list
// into per-index read-capacity AutoScaling descriptions keyed by index name.
func parseGSIAutoScalingRead(updates interface{}) (map[string]map[string]interface{}, error) {
	result := map[string]map[string]interface{}{}
	gsiUpdates, ok := updates.([]interface{})
	if !ok {
		return result, nil
	}
	for _, u := range gsiUpdates {
		uMap, ok := u.(map[string]interface{})
		if !ok {
			continue
		}
		indexName, _ := uMap["IndexName"].(string)
		if indexName == "" {
			continue
		}
		if readAS, ok := uMap["ProvisionedReadCapacityAutoScalingUpdate"].(map[string]interface{}); ok {
			settings, err := parseAutoScalingSettings(readAS)
			if err != nil {
				return nil, err
			}
			result[indexName] = settings
		}
	}
	return result, nil
}

// mergeReplicaGSIAutoScaling merges the table-level write-side and
// replica-level read-side per-index AutoScaling settings into the
// GlobalSecondaryIndexes echo list, index names in a stable order.
func mergeReplicaGSIAutoScaling(write, read map[string]map[string]interface{}) []map[string]interface{} {
	if len(write) == 0 && len(read) == 0 {
		return nil
	}
	indexes := make([]string, 0, len(write)+len(read))
	seen := make(map[string]bool, len(write)+len(read))
	for name := range write {
		indexes = append(indexes, name)
		seen[name] = true
	}
	for name := range read {
		if !seen[name] {
			indexes = append(indexes, name)
		}
	}
	sort.Strings(indexes)
	result := make([]map[string]interface{}, 0, len(indexes))
	for _, name := range indexes {
		entry := map[string]interface{}{"IndexName": name}
		if settings, ok := read[name]; ok {
			entry["ProvisionedReadCapacityAutoScalingSettings"] = settings
		}
		if settings, ok := write[name]; ok {
			entry["ProvisionedWriteCapacityAutoScalingSettings"] = settings
		}
		result = append(result, entry)
	}
	return result
}
