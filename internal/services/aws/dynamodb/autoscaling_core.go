package dynamodb

import (
	"context"
	"sort"

	"vorpalstacks/internal/common/request"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
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
//
// Typed store records flow through this file: requests parse into
// *dbstore.AutoScalingSettingsDescription values, merges happen on the typed
// form, and the wire shape is rendered only when a response is assembled.
// ---------------------------------------------------------------------------

// autoScalingSettingsToWire renders one capacity dimension's auto-scaling
// description in the response shape. Only set members are rendered; a nil
// description renders as nil.
func autoScalingSettingsToWire(s *dbstore.AutoScalingSettingsDescription) map[string]interface{} {
	if s == nil {
		return nil
	}
	desc := map[string]interface{}{}
	if s.MinimumUnits != nil {
		desc["MinimumUnits"] = *s.MinimumUnits
	}
	if s.MaximumUnits != nil {
		desc["MaximumUnits"] = *s.MaximumUnits
	}
	if s.AutoScalingDisabled != nil {
		desc["AutoScalingDisabled"] = *s.AutoScalingDisabled
	}
	if s.AutoScalingRoleArn != nil {
		desc["AutoScalingRoleArn"] = *s.AutoScalingRoleArn
	}
	if len(s.ScalingPolicies) > 0 {
		policies := make([]interface{}, 0, len(s.ScalingPolicies))
		for _, policy := range s.ScalingPolicies {
			entry := map[string]interface{}{}
			if policy.PolicyName != nil {
				entry["PolicyName"] = *policy.PolicyName
			}
			if tt := policy.TargetTrackingScalingPolicyConfiguration; tt != nil {
				ttDesc := map[string]interface{}{"TargetValue": tt.TargetValue}
				if tt.DisableScaleIn != nil {
					ttDesc["DisableScaleIn"] = *tt.DisableScaleIn
				}
				if tt.ScaleInCooldown != nil {
					ttDesc["ScaleInCooldown"] = *tt.ScaleInCooldown
				}
				if tt.ScaleOutCooldown != nil {
					ttDesc["ScaleOutCooldown"] = *tt.ScaleOutCooldown
				}
				entry["TargetTrackingScalingPolicyConfiguration"] = ttDesc
			}
			policies = append(policies, entry)
		}
		desc["ScalingPolicies"] = policies
	}
	return desc
}

// indexAutoScalingSettingsToWire renders one index's settings entry; only
// set members are rendered.
func indexAutoScalingSettingsToWire(index dbstore.IndexAutoScalingSettings) map[string]interface{} {
	entry := map[string]interface{}{"IndexName": index.IndexName}
	if index.ProvisionedReadCapacityUnits != nil {
		entry["ProvisionedReadCapacityUnits"] = *index.ProvisionedReadCapacityUnits
	}
	if index.ProvisionedWriteCapacityUnits != nil {
		entry["ProvisionedWriteCapacityUnits"] = *index.ProvisionedWriteCapacityUnits
	}
	if read := autoScalingSettingsToWire(index.Read); read != nil {
		entry["ProvisionedReadCapacityAutoScalingSettings"] = read
	}
	if write := autoScalingSettingsToWire(index.Write); write != nil {
		entry["ProvisionedWriteCapacityAutoScalingSettings"] = write
	}
	return entry
}

// replicaAutoScalingDescriptionsToWire renders the per-replica descriptions
// in the response shape. The result is always non-nil so serialisation
// renders an empty list, never null.
func replicaAutoScalingDescriptionsToWire(replicas []dbstore.ReplicaAutoScalingDescription) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(replicas))
	for _, replica := range replicas {
		entry := map[string]interface{}{"RegionName": replica.RegionName}
		if read := autoScalingSettingsToWire(replica.Read); read != nil {
			entry["ReplicaProvisionedReadCapacityAutoScalingSettings"] = read
		}
		if write := autoScalingSettingsToWire(replica.Write); write != nil {
			entry["ReplicaProvisionedWriteCapacityAutoScalingSettings"] = write
		}
		if len(replica.GlobalSecondaryIndexes) > 0 {
			gsi := make([]interface{}, 0, len(replica.GlobalSecondaryIndexes))
			for _, index := range replica.GlobalSecondaryIndexes {
				gsi = append(gsi, indexAutoScalingSettingsToWire(index))
			}
			entry["GlobalSecondaryIndexes"] = gsi
		}
		result = append(result, entry)
	}
	return result
}

// replicaDescriptionsFromSettings reads the stored per-region replica
// descriptions. The result is always non-nil so merges start from an empty
// list, never from a nil slice.
func replicaDescriptionsFromSettings(settings *dbstore.TableReplicaAutoScalingSettings) []dbstore.ReplicaAutoScalingDescription {
	if settings == nil {
		return []dbstore.ReplicaAutoScalingDescription{}
	}
	return settings.Replicas
}

// mergeReplicaDescriptions upserts the updated per-region replica
// descriptions into the stored list: a region named in the update replaces
// its stored description, unmentioned regions are preserved, and new
// regions are appended in the update's order. ReplicaUpdates carries
// modification semantics only — the model has no delete action.
func mergeReplicaDescriptions(existing, updated []dbstore.ReplicaAutoScalingDescription) []dbstore.ReplicaAutoScalingDescription {
	merged := make([]dbstore.ReplicaAutoScalingDescription, len(existing))
	copy(merged, existing)
	byRegion := make(map[string]int, len(merged))
	for i, desc := range merged {
		byRegion[desc.RegionName] = i
	}
	for _, desc := range updated {
		if idx, ok := byRegion[desc.RegionName]; ok {
			merged[idx] = desc
			continue
		}
		merged = append(merged, desc)
	}
	return merged
}

// tableAutoScalingDescription builds the TableAutoScalingDescription
// response shared by the describe and update operations.
func tableAutoScalingDescription(table *dbstore.Table, replicas []map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"TableAutoScalingDescription": map[string]interface{}{
			"TableName":   table.Name,
			"TableStatus": string(table.Status),
			"Replicas":    replicas,
		},
	}
}

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
	return tableAutoScalingDescription(table, replicaAutoScalingDescriptionsToWire(replicaDescriptionsFromSettings(settings))), nil
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
	replicas := []dbstore.ReplicaAutoScalingDescription{}

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
		replicaByRegion := make(map[string]int)

		for _, update := range replicaUpdates {
			updateMap, ok := update.(map[string]interface{})
			if !ok {
				return nil, ErrInvalidParameter
			}
			regionName, _ := updateMap["RegionName"].(string)
			if regionName == "" {
				return nil, ErrInvalidParameter
			}

			idx, exists := replicaByRegion[regionName]
			if !exists {
				replicas = append(replicas, dbstore.ReplicaAutoScalingDescription{RegionName: regionName})
				idx = len(replicas) - 1
				replicaByRegion[regionName] = idx
			}

			if readAS, ok := updateMap["ReplicaProvisionedReadCapacityAutoScalingUpdate"].(map[string]interface{}); ok {
				settings, err := parseAutoScalingSettings(readAS)
				if err != nil {
					return nil, err
				}
				replicas[idx].Read = settings
			}
			if tableWriteAS != nil {
				replicas[idx].Write = tableWriteAS
			}
			gsiReadAS, err := parseGSIAutoScalingRead(updateMap["ReplicaGlobalSecondaryIndexUpdates"])
			if err != nil {
				return nil, err
			}
			if gsi := mergeReplicaGSIAutoScaling(tableGSIWriteAS, gsiReadAS); len(gsi) > 0 {
				replicas[idx].GlobalSecondaryIndexes = gsi
			}
		}
	}

	// Merge the parsed per-region descriptions into the stored settings
	// (upsert by region) so an update without ReplicaUpdates — e.g. a
	// table-level write-capacity-only update — keeps the previously stored
	// replica descriptions.
	existing, err := store.Tables().GetAutoScalingSettings(table.Name)
	if err != nil {
		return nil, err
	}
	replicas = mergeReplicaDescriptions(replicaDescriptionsFromSettings(existing), replicas)

	if err := store.Tables().SetAutoScalingSettings(table.Name, &dbstore.TableReplicaAutoScalingSettings{Replicas: replicas}); err != nil {
		return nil, err
	}

	return tableAutoScalingDescription(table, replicaAutoScalingDescriptionsToWire(replicas)), nil
}

// wirePositiveLong converts a JSON number into a non-negative integer
// capacity value; the model types these members as Long.
func wirePositiveLong(v interface{}) (*int64, error) {
	f, ok := v.(float64)
	if !ok || f < 0 || f != float64(int64(f)) {
		return nil, ErrInvalidParameter
	}
	n := int64(f)
	return &n, nil
}

// parseAutoScalingSettings extracts AutoScaling settings from a request
// parameter map into a typed description, validating the Smithy member
// constraints on the values it carries. Returns ErrInvalidParameter when a
// member violates its documented bounds.
func parseAutoScalingSettings(m map[string]interface{}) (*dbstore.AutoScalingSettingsDescription, error) {
	desc := &dbstore.AutoScalingSettingsDescription{}
	if v, ok := m["MinimumUnits"]; ok {
		units, err := wirePositiveLong(v)
		if err != nil {
			return nil, err
		}
		desc.MinimumUnits = units
	}
	if v, ok := m["MaximumUnits"]; ok {
		units, err := wirePositiveLong(v)
		if err != nil {
			return nil, err
		}
		desc.MaximumUnits = units
	}
	if v, ok := m["AutoScalingDisabled"]; ok {
		if disabled, isBool := v.(bool); isBool {
			desc.AutoScalingDisabled = &disabled
		}
	}
	if v, ok := m["AutoScalingRoleArn"]; ok {
		if roleArn, isStr := v.(string); isStr {
			if !validateAutoScalingRoleArn(roleArn) {
				return nil, ErrInvalidParameter
			}
			desc.AutoScalingRoleArn = &roleArn
		}
	}
	if pol, ok := m["ScalingPolicyUpdate"].(map[string]interface{}); ok {
		policy := dbstore.AutoScalingPolicyDescription{}
		if name, ok := pol["PolicyName"]; ok {
			if policyName, isStr := name.(string); isStr {
				if !validateAutoScalingPolicyName(policyName) {
					return nil, ErrInvalidParameter
				}
				policy.PolicyName = &policyName
			}
		}
		// The description form carries the target tracking configuration
		// alongside the policy name; its TargetValue member is required and
		// bounded by the documented metric range.
		if tt, ok := pol["TargetTrackingScalingPolicyConfiguration"].(map[string]interface{}); ok {
			target, ok := tt["TargetValue"].(float64)
			if !ok {
				return nil, ErrInvalidParameter
			} else if target < autoScalingTargetValueMin || target > autoScalingTargetValueMax {
				return nil, ErrInvalidParameter
			}
			ttDesc := &dbstore.TargetTrackingScalingPolicyConfiguration{TargetValue: target}
			if v, ok := tt["DisableScaleIn"].(bool); ok {
				ttDesc.DisableScaleIn = &v
			}
			if v, ok := tt["ScaleInCooldown"].(float64); ok {
				cooldown := int32(v)
				ttDesc.ScaleInCooldown = &cooldown
			}
			if v, ok := tt["ScaleOutCooldown"].(float64); ok {
				cooldown := int32(v)
				ttDesc.ScaleOutCooldown = &cooldown
			}
			policy.TargetTrackingScalingPolicyConfiguration = ttDesc
		}
		desc.ScalingPolicies = []dbstore.AutoScalingPolicyDescription{policy}
	}
	return desc, nil
}

// parseOptionalAutoScalingSettings parses an optional top-level
// AutoScalingSettingsUpdate member, returning nil when the member is absent.
func parseOptionalAutoScalingSettings(parameters map[string]interface{}, key string) (*dbstore.AutoScalingSettingsDescription, error) {
	if m, ok := parameters[key].(map[string]interface{}); ok {
		return parseAutoScalingSettings(m)
	}
	return nil, nil
}

// parseGSIAutoScalingWrite parses a GlobalSecondaryIndexUpdates list into
// per-index write-capacity AutoScaling descriptions keyed by index name.
func parseGSIAutoScalingWrite(updates interface{}) (map[string]*dbstore.AutoScalingSettingsDescription, error) {
	result := map[string]*dbstore.AutoScalingSettingsDescription{}
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
func parseGSIAutoScalingRead(updates interface{}) (map[string]*dbstore.AutoScalingSettingsDescription, error) {
	result := map[string]*dbstore.AutoScalingSettingsDescription{}
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
// replica-level read-side per-index AutoScaling settings into the per-index
// settings list, index names in a stable order.
func mergeReplicaGSIAutoScaling(write, read map[string]*dbstore.AutoScalingSettingsDescription) []dbstore.IndexAutoScalingSettings {
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
	result := make([]dbstore.IndexAutoScalingSettings, 0, len(indexes))
	for _, name := range indexes {
		result = append(result, dbstore.IndexAutoScalingSettings{
			IndexName: name,
			Read:      read[name],
			Write:     write[name],
		})
	}
	return result
}
