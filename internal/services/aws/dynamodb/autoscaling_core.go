package dynamodb

import (
	"context"
	"math"
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
	// implemented. The carrier and the read-capacity member hold the
	// typed-member presence contract: a present value of the wrong type is
	// the invalid-parameter error, never a silently skipped update.
	if replicasRaw := in.Parameters["ReplicaUpdates"]; replicasRaw != nil {
		replicaUpdates, ok := replicasRaw.([]interface{})
		if !ok {
			return nil, ErrInvalidParameter
		}
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

			if readAS, present, err := wireStructMember(updateMap["ReplicaProvisionedReadCapacityAutoScalingUpdate"]); err != nil {
				return nil, err
			} else if present {
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
	// (upsert by region) under the table's record lock — the read, the
	// merge, and the write form one locked read-modify-write, so two
	// concurrent updates cannot each start from the same stored state and
	// silently discard each other's replicas. An update without
	// ReplicaUpdates — e.g. a table-level write-capacity-only update —
	// keeps the previously stored replica descriptions through the same
	// merge.
	merged, err := store.Tables().UpdateAutoScalingSettings(table.Name, func(existing *dbstore.TableReplicaAutoScalingSettings) *dbstore.TableReplicaAutoScalingSettings {
		return &dbstore.TableReplicaAutoScalingSettings{
			Replicas: mergeReplicaDescriptions(replicaDescriptionsFromSettings(existing), replicas),
		}
	})
	if err != nil {
		return nil, err
	}

	return tableAutoScalingDescription(table, replicaAutoScalingDescriptionsToWire(merged.Replicas)), nil
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

// wireStringMember converts a JSON String member into its value,
// reporting whether the member was carried: the model types these
// members as String, so a present value of another type is the
// invalid-parameter error, never a silently skipped update. A nil value
// is an absent member — the protocol omits unset members rather than
// carrying an explicit null.
func wireStringMember(v interface{}) (value string, present bool, err error) {
	if v == nil {
		return "", false, nil
	}
	s, ok := v.(string)
	if !ok {
		return "", true, ErrInvalidParameter
	}
	return s, true, nil
}

// wireBoolMember converts a JSON Boolean member into its value with the
// typed-member family's discipline: the model types these members as
// Boolean, so a present value of another type is the invalid-parameter
// error, never a silently skipped update. A nil value is an absent member.
func wireBoolMember(v interface{}) (value bool, present bool, err error) {
	if v == nil {
		return false, false, nil
	}
	b, ok := v.(bool)
	if !ok {
		return false, true, ErrInvalidParameter
	}
	return b, true, nil
}

// wireIntegerObject converts a JSON number into an int32 for a member the
// model types as IntegerObject: the modelled shape carries no range, so
// integrality is the documented bound and the int32 bounds are the
// platform's storage type — a non-integral number or one outside int32 is
// the invalid-parameter error, never a silently truncated value. A nil
// value is an absent member.
func wireIntegerObject(v interface{}) (value int32, present bool, err error) {
	if v == nil {
		return 0, false, nil
	}
	f, ok := v.(float64)
	if !ok || f != math.Trunc(f) || f < math.MinInt32 || f > math.MaxInt32 {
		return 0, true, ErrInvalidParameter
	}
	return int32(f), true, nil
}

// wireStructMember converts a JSON structure member into its map form,
// reporting whether the member was carried — the typed-member family's
// discipline for structure members: a present value of another type is the
// invalid-parameter error, never a silently skipped specification. A nil
// value is an absent member.
func wireStructMember(v interface{}) (value map[string]interface{}, present bool, err error) {
	if v == nil {
		return nil, false, nil
	}
	m, ok := v.(map[string]interface{})
	if !ok {
		return nil, true, ErrInvalidParameter
	}
	return m, true, nil
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
		disabled, present, err := wireBoolMember(v)
		if err != nil {
			return nil, err
		}
		if present {
			desc.AutoScalingDisabled = &disabled
		}
	}
	if v, ok := m["AutoScalingRoleArn"]; ok {
		roleArn, present, err := wireStringMember(v)
		if err != nil {
			return nil, err
		}
		if present {
			if !validateAutoScalingRoleArn(roleArn) {
				return nil, ErrInvalidParameter
			}
			desc.AutoScalingRoleArn = &roleArn
		}
	}
	if pol, present, err := wireStructMember(m["ScalingPolicyUpdate"]); err != nil {
		return nil, err
	} else if present {
		policy := dbstore.AutoScalingPolicyDescription{}
		if name, ok := pol["PolicyName"]; ok {
			policyName, present, err := wireStringMember(name)
			if err != nil {
				return nil, err
			}
			if present {
				if !validateAutoScalingPolicyName(policyName) {
					return nil, ErrInvalidParameter
				}
				policy.PolicyName = &policyName
			}
		}
		// The description form carries the target tracking configuration
		// alongside the policy name; its TargetValue member is required and
		// bounded by the documented metric range.
		if tt, present, err := wireStructMember(pol["TargetTrackingScalingPolicyConfiguration"]); err != nil {
			return nil, err
		} else if present {
			target, ok := tt["TargetValue"].(float64)
			if !ok {
				return nil, ErrInvalidParameter
			} else if target < autoScalingTargetValueMin || target > autoScalingTargetValueMax {
				return nil, ErrInvalidParameter
			}
			ttDesc := &dbstore.TargetTrackingScalingPolicyConfiguration{TargetValue: target}
			if v, ok := tt["DisableScaleIn"]; ok {
				disableScaleIn, present, err := wireBoolMember(v)
				if err != nil {
					return nil, err
				}
				if present {
					ttDesc.DisableScaleIn = &disableScaleIn
				}
			}
			if v, ok := tt["ScaleInCooldown"]; ok {
				cooldown, present, err := wireIntegerObject(v)
				if err != nil {
					return nil, err
				}
				if present {
					ttDesc.ScaleInCooldown = &cooldown
				}
			}
			if v, ok := tt["ScaleOutCooldown"]; ok {
				cooldown, present, err := wireIntegerObject(v)
				if err != nil {
					return nil, err
				}
				if present {
					ttDesc.ScaleOutCooldown = &cooldown
				}
			}
			policy.TargetTrackingScalingPolicyConfiguration = ttDesc
		}
		desc.ScalingPolicies = []dbstore.AutoScalingPolicyDescription{policy}
	}
	return desc, nil
}

// parseOptionalAutoScalingSettings parses an optional top-level
// AutoScalingSettingsUpdate member, returning nil when the member is
// absent; a present value of another type is the invalid-parameter error.
func parseOptionalAutoScalingSettings(parameters map[string]interface{}, key string) (*dbstore.AutoScalingSettingsDescription, error) {
	m, present, err := wireStructMember(parameters[key])
	if err != nil {
		return nil, err
	}
	if !present {
		return nil, nil
	}
	return parseAutoScalingSettings(m)
}

// parseGSIAutoScalingList walks one per-index AutoScaling updates list —
// the table-level write side and the replica-level read side share the
// entry shape: neither IndexName nor the capacity-side member is required,
// and an entry lacking either contributes nothing. The side names its own
// capacity member key. A present member of the wrong wire type — the list
// itself, an entry, IndexName, or the capacity-side member — is the
// invalid-parameter error, never a silently skipped update.
func parseGSIAutoScalingList(updates interface{}, capacityMember string) (map[string]*dbstore.AutoScalingSettingsDescription, error) {
	result := map[string]*dbstore.AutoScalingSettingsDescription{}
	if updates == nil {
		return result, nil
	}
	gsiUpdates, ok := updates.([]interface{})
	if !ok {
		return nil, ErrInvalidParameter
	}
	for _, u := range gsiUpdates {
		uMap, present, err := wireStructMember(u)
		if err != nil {
			return nil, err
		}
		if !present {
			continue
		}
		indexName, present, err := wireStringMember(uMap["IndexName"])
		if err != nil {
			return nil, err
		}
		if !present || indexName == "" {
			continue
		}
		as, present, err := wireStructMember(uMap[capacityMember])
		if err != nil {
			return nil, err
		}
		if !present {
			continue
		}
		settings, err := parseAutoScalingSettings(as)
		if err != nil {
			return nil, err
		}
		result[indexName] = settings
	}
	return result, nil
}

// parseGSIAutoScalingWrite parses a GlobalSecondaryIndexUpdates list into
// per-index write-capacity AutoScaling descriptions keyed by index name.
func parseGSIAutoScalingWrite(updates interface{}) (map[string]*dbstore.AutoScalingSettingsDescription, error) {
	return parseGSIAutoScalingList(updates, "ProvisionedWriteCapacityAutoScalingUpdate")
}

// parseGSIAutoScalingRead parses a ReplicaGlobalSecondaryIndexUpdates list
// into per-index read-capacity AutoScaling descriptions keyed by index name.
func parseGSIAutoScalingRead(updates interface{}) (map[string]*dbstore.AutoScalingSettingsDescription, error) {
	return parseGSIAutoScalingList(updates, "ProvisionedReadCapacityAutoScalingUpdate")
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
