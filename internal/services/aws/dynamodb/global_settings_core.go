package dynamodb

import (
	"context"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	commonstore "vorpalstacks/internal/store/aws/common"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// ---------------------------------------------------------------------------
// Global Settings Core — single validation + persistence path for the
// global table settings operations (DescribeGlobalTableSettings /
// UpdateGlobalTableSettings): the per-replica provisioned capacity and
// auto-scaling settings, the global write-side GSI settings, and the
// parse/merge machinery both operations share. The global table
// membership plane (Create/Describe/List/UpdateGlobalTable) lives in
// global_table_core.go.
// ---------------------------------------------------------------------------

// describeGlobalTableSettingsInput carries the raw wire parameters for
// DescribeGlobalTableSettings.
type describeGlobalTableSettingsInput struct {
	Parameters map[string]interface{}
}

// describeGlobalTableSettingsCore validates the request, then returns the
// settings of the global table.
func (s *DynamoDBService) describeGlobalTableSettingsCore(ctx context.Context, reqCtx *request.RequestContext, in describeGlobalTableSettingsInput) (interface{}, error) {
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
		"GlobalTableName": globalTable.GlobalTableName,
		"ReplicaSettings": buildGlobalTableReplicaSettings(globalTable),
	}, nil
}

// buildGlobalTableReplicaSettings renders the per-replica settings list
// shared by the global table settings operations. The global-level write
// auto-scaling and per-index write settings echo on every replica; the
// per-replica read settings merge with the per-index write side.
func buildGlobalTableReplicaSettings(globalTable *dbstore.GlobalTable) []map[string]interface{} {
	replicaSettings := make([]map[string]interface{}, 0, len(globalTable.ReplicationGroup))
	for _, replica := range globalTable.ReplicationGroup {
		settings := map[string]interface{}{
			"RegionName":                           replica.RegionName,
			"ReplicaStatus":                        string(replica.ReplicaStatus),
			"ReplicaProvisionedReadCapacityUnits":  replica.ProvisionedReadCapacityUnits,
			"ReplicaProvisionedWriteCapacityUnits": replica.ProvisionedWriteCapacityUnits,
		}
		if replica.BillingMode != "" {
			settings["ReplicaBillingModeSummary"] = map[string]interface{}{
				"BillingMode": string(replica.BillingMode),
			}
		}
		if replica.ReadAutoScalingSettings != nil {
			settings["ReplicaProvisionedReadCapacityAutoScalingSettings"] = autoScalingSettingsToWire(replica.ReadAutoScalingSettings)
		}
		if globalTable.WriteAutoScalingSettings != nil {
			settings["ReplicaProvisionedWriteCapacityAutoScalingSettings"] = autoScalingSettingsToWire(globalTable.WriteAutoScalingSettings)
		}
		if gsi := mergeGSISettings(globalTable, replica); len(gsi) > 0 {
			settings["ReplicaGlobalSecondaryIndexSettings"] = gsi
		}
		if replica.TableClass != "" {
			summary := map[string]interface{}{"TableClass": string(replica.TableClass)}
			if replica.TableClassLastUpdated != nil {
				summary["LastUpdateDateTime"] = *replica.TableClassLastUpdated
			}
			settings["ReplicaTableClassSummary"] = summary
		}
		replicaSettings = append(replicaSettings, settings)
	}
	return replicaSettings
}

// parseGSISettingsUpdates walks one GSI settings update list — the global
// table's write side and the replica's read side share the entry shape:
// the model marks each entry's IndexName required, bounds the list at
// 1-20 entries, and leaves the capacity member to the side-specific
// extractor the caller passes. The list and its entries hold the
// typed-member presence contract: a present value of the wrong type is
// the invalid-parameter error, never a silently skipped entry.
func parseGSISettingsUpdates(updates interface{}, extractSide func(uMap map[string]interface{}, entry *dbstore.IndexAutoScalingSettings) error) ([]dbstore.IndexAutoScalingSettings, error) {
	if updates == nil {
		return nil, nil
	}
	gsiUpdates, ok := updates.([]interface{})
	if !ok {
		return nil, ErrInvalidParameter
	}
	// The model bounds the settings update list at 1-20 entries.
	if len(gsiUpdates) < 1 || len(gsiUpdates) > dbstore.MaxReplicaGSISettingsUpdates {
		return nil, ErrInvalidParameter
	}
	var result []dbstore.IndexAutoScalingSettings
	for _, u := range gsiUpdates {
		uMap, ok := u.(map[string]interface{})
		if !ok {
			return nil, ErrInvalidParameter
		}
		indexName, _ := uMap["IndexName"].(string)
		if indexName == "" {
			return nil, ErrInvalidParameter
		}
		entry := dbstore.IndexAutoScalingSettings{IndexName: indexName}
		if err := extractSide(uMap, &entry); err != nil {
			return nil, err
		}
		result = append(result, entry)
	}
	return result, nil
}

// parseGlobalGSIWriteSettings converts GlobalTableGlobalSecondaryIndexSettingsUpdate
// entries into the per-index write settings stored on the global table. The
// model marks each entry's IndexName required and the capacity member a
// positive long.
func parseGlobalGSIWriteSettings(updates interface{}) ([]dbstore.IndexAutoScalingSettings, error) {
	return parseGSISettingsUpdates(updates, func(uMap map[string]interface{}, entry *dbstore.IndexAutoScalingSettings) error {
		if raw, present := uMap["ProvisionedWriteCapacityUnits"]; present {
			units, err := wirePositiveLong(raw)
			if err != nil || *units < 1 {
				return ErrInvalidParameter
			}
			entry.ProvisionedWriteCapacityUnits = units
		}
		if writeAS, present, err := wireStructMember(uMap["ProvisionedWriteCapacityAutoScalingSettingsUpdate"]); err != nil {
			return err
		} else if present {
			settings, err := parseAutoScalingSettings(writeAS)
			if err != nil {
				return err
			}
			entry.Write = settings
		}
		return nil
	})
}

// parseReplicaGSIReadSettings converts ReplicaGlobalSecondaryIndexSettingsUpdate
// entries into the per-index read settings stored on the replica. The model
// marks each entry's IndexName required and the capacity member a positive
// long.
func parseReplicaGSIReadSettings(updates interface{}) ([]dbstore.IndexAutoScalingSettings, error) {
	return parseGSISettingsUpdates(updates, func(uMap map[string]interface{}, entry *dbstore.IndexAutoScalingSettings) error {
		if raw, present := uMap["ProvisionedReadCapacityUnits"]; present {
			units, err := wirePositiveLong(raw)
			if err != nil || *units < 1 {
				return ErrInvalidParameter
			}
			entry.ProvisionedReadCapacityUnits = units
		}
		if readAS, present, err := wireStructMember(uMap["ProvisionedReadCapacityAutoScalingSettingsUpdate"]); err != nil {
			return err
		} else if present {
			settings, err := parseAutoScalingSettings(readAS)
			if err != nil {
				return err
			}
			entry.Read = settings
		}
		return nil
	})
}

// mergeIndexSettingsLists applies update semantics to a per-index settings
// list: an update whose IndexName matches a stored entry replaces it, any
// other update is appended.
func mergeIndexSettingsLists(stored, updates []dbstore.IndexAutoScalingSettings) []dbstore.IndexAutoScalingSettings {
	merged := make([]dbstore.IndexAutoScalingSettings, len(stored))
	copy(merged, stored)
	for _, update := range updates {
		replaced := false
		for i, entry := range merged {
			if entry.IndexName == update.IndexName {
				merged[i] = update
				replaced = true
				break
			}
		}
		if !replaced {
			merged = append(merged, update)
		}
	}
	return merged
}

// mergeGSISettings merges the global write-side and replica read-side
// per-index settings into the ReplicaGlobalSecondaryIndexSettings echo.
func mergeGSISettings(globalTable *dbstore.GlobalTable, replica *dbstore.Replica) []map[string]interface{} {
	if len(globalTable.GlobalSecondaryIndexWriteSettings) == 0 && len(replica.GlobalSecondaryIndexReadSettings) == 0 {
		return nil
	}
	merged := make(map[string]*dbstore.IndexAutoScalingSettings)
	var order []string
	add := func(s dbstore.IndexAutoScalingSettings) {
		if s.IndexName == "" {
			return
		}
		entry, exists := merged[s.IndexName]
		if !exists {
			entry = &dbstore.IndexAutoScalingSettings{IndexName: s.IndexName}
			merged[s.IndexName] = entry
			order = append(order, s.IndexName)
		}
		if s.ProvisionedReadCapacityUnits != nil {
			entry.ProvisionedReadCapacityUnits = s.ProvisionedReadCapacityUnits
		}
		if s.ProvisionedWriteCapacityUnits != nil {
			entry.ProvisionedWriteCapacityUnits = s.ProvisionedWriteCapacityUnits
		}
		if s.Read != nil {
			entry.Read = s.Read
		}
		if s.Write != nil {
			entry.Write = s.Write
		}
	}
	for _, settings := range globalTable.GlobalSecondaryIndexWriteSettings {
		add(settings)
	}
	for _, settings := range replica.GlobalSecondaryIndexReadSettings {
		add(settings)
	}
	result := make([]map[string]interface{}, 0, len(order))
	for _, name := range order {
		result = append(result, indexAutoScalingSettingsToWire(*merged[name]))
	}
	return result
}

// updateGlobalTableSettingsInput carries the raw wire parameters for
// UpdateGlobalTableSettings.
type updateGlobalTableSettingsInput struct {
	Parameters map[string]interface{}
}

// replicaSettingsUpdate is one parsed ReplicaSettingsUpdate entry: the
// region it addresses and the per-replica members the request carries.
type replicaSettingsUpdate struct {
	region     string
	readUnits  *int64
	readAS     *dbstore.AutoScalingSettingsDescription
	gsiRead    []dbstore.IndexAutoScalingSettings
	tableClass dbstore.TableClass
}

// resultingGlobalBillingMode resolves the billing mode the settings update
// leaves the replica tables in: the explicit GlobalTableBillingMode, or —
// because every capacity member of this request is a provisioned-mode
// member on the table plane, and the model documents PROVISIONED as the
// billing-mode member's default — PROVISIONED when any capacity member is
// present. Empty means the request carries no capacity member at all and
// the tables keep their mode.
func resultingGlobalBillingMode(explicit string, hasCapacityMember bool) string {
	if explicit != "" {
		return explicit
	}
	if hasCapacityMember {
		return string(dbstore.BillingModeProvisioned)
	}
	return ""
}

// gsiExistsOnTable reports whether the table holds a global secondary
// index under the name.
func gsiExistsOnTable(table *dbstore.Table, indexName string) bool {
	for _, gsi := range table.GlobalSecondaryIndexes {
		if gsi.IndexName == indexName {
			return true
		}
	}
	return false
}

// updateGlobalTableSettingsCore validates the request, applies the replica
// settings updates to the global table record, and propagates the capacity
// settings to every member region's table record — the settings are the
// replica tables' state, and the global record is their echo: leaving the
// tables untouched would make DescribeTable and DescribeGlobalTableSettings
// report different state for the same replica. Cross-region validation runs
// before the global record commits so a refused request mutates nothing.
func (s *DynamoDBService) updateGlobalTableSettingsCore(ctx context.Context, reqCtx *request.RequestContext, in updateGlobalTableSettingsInput) (interface{}, error) {
	globalTableName := request.GetStringParam(in.Parameters, "GlobalTableName")
	if !validateResourceName(globalTableName) {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	// The billing mode and provisioned write capacity are global settings
	// in the model — they apply to every replica of the global table, not
	// to a per-replica member.
	globalBillingMode := request.GetStringParam(in.Parameters, "GlobalTableBillingMode")
	if globalBillingMode != "" &&
		globalBillingMode != string(dbstore.BillingModeProvisioned) &&
		globalBillingMode != string(dbstore.BillingModePayPerRequest) {
		return nil, ErrInvalidParameter
	}
	globalWriteUnits := int64(0)
	hasGlobalWriteUnits := false
	if raw, present := in.Parameters["GlobalTableProvisionedWriteCapacityUnits"]; present {
		// The member is a positive Long: fractional input is a
		// ValidationException, the same contract the GSI capacity members
		// of this request enforce.
		units, err := wirePositiveLong(raw)
		if err != nil || *units < 1 {
			return nil, ErrInvalidParameter
		}
		globalWriteUnits = *units
		hasGlobalWriteUnits = true
	}
	globalWriteAS, err := parseOptionalAutoScalingSettings(in.Parameters, "GlobalTableProvisionedWriteCapacityAutoScalingSettingsUpdate")
	if err != nil {
		return nil, err
	}
	globalGSIUpdates, err := parseGlobalGSIWriteSettings(in.Parameters["GlobalTableGlobalSecondaryIndexSettingsUpdate"])
	if err != nil {
		return nil, err
	}

	var replicaUpdates []replicaSettingsUpdate
	if rawSettings := in.Parameters["ReplicaSettingsUpdate"]; rawSettings != nil {
		updates, ok := rawSettings.([]interface{})
		if !ok {
			return nil, ErrInvalidParameter
		}
		// The model bounds the replica settings update list at 1-50 entries.
		if len(updates) < 1 || len(updates) > dbstore.MaxReplicaSettingsUpdates {
			return nil, ErrInvalidParameter
		}
		for _, update := range updates {
			updateMap, ok := update.(map[string]interface{})
			if !ok {
				return nil, ErrInvalidParameter
			}
			regionName, _ := updateMap["RegionName"].(string)
			if regionName == "" {
				return nil, ErrInvalidParameter
			}
			entry := replicaSettingsUpdate{region: regionName}
			if rawUnits, present := updateMap["ReplicaProvisionedReadCapacityUnits"]; present {
				// A positive Long member: fractional input is a
				// ValidationException, matching the GSI read capacity
				// members of the same request.
				units, err := wirePositiveLong(rawUnits)
				if err != nil || *units < 1 {
					return nil, ErrInvalidParameter
				}
				entry.readUnits = units
			}
			if readAS, present, err := wireStructMember(updateMap["ReplicaProvisionedReadCapacityAutoScalingSettingsUpdate"]); err != nil {
				return nil, err
			} else if present {
				settings, err := parseAutoScalingSettings(readAS)
				if err != nil {
					return nil, err
				}
				entry.readAS = settings
			}
			if entry.gsiRead, err = parseReplicaGSIReadSettings(updateMap["ReplicaGlobalSecondaryIndexSettingsUpdate"]); err != nil {
				return nil, err
			}
			if classStr, present, err := wireStringMember(updateMap["ReplicaTableClass"]); err != nil {
				return nil, err
			} else if present && classStr != "" {
				class := dbstore.TableClass(classStr)
				if class != dbstore.TableClassStandard && class != dbstore.TableClassStandardInfrequentAccess {
					return nil, ErrInvalidParameter
				}
				entry.tableClass = class
			}
			replicaUpdates = append(replicaUpdates, entry)
		}
	}
	readUpdateByRegion := make(map[string]replicaSettingsUpdate, len(replicaUpdates))
	for _, entry := range replicaUpdates {
		readUpdateByRegion[entry.region] = entry
	}

	// The member regions' table records are the validation truth: the
	// settings apply to those records, so their schema and capacity state
	// decides what the request may set. They are loaded before the global
	// record commits so a refused request mutates nothing anywhere.
	group, err := store.GlobalTables().Get(globalTableName)
	if err != nil {
		if commonstore.IsNotFound(err) {
			return nil, ErrGlobalTableNotFound
		}
		return nil, err
	}
	memberTables := make(map[string]*dbstore.Table, len(group.ReplicationGroup))
	for _, replica := range group.ReplicationGroup {
		regionStore, err := s.GetCachedStoreForRegion(replica.RegionName)
		if err != nil {
			return nil, err
		}
		table, err := regionStore.Tables().Get(globalTableName)
		if err != nil {
			if commonstore.IsNotFound(err) {
				return nil, ErrTableNotFoundException
			}
			return nil, err
		}
		memberTables[replica.RegionName] = table
	}

	// A per-index settings update must name a global secondary index the
	// table actually holds: the model's IndexNotFoundException. The global
	// write side applies to every member table; a replica's read side to
	// that member's table alone.
	for _, update := range globalGSIUpdates {
		for _, table := range memberTables {
			if !gsiExistsOnTable(table, update.IndexName) {
				return nil, ErrIndexNotFound
			}
		}
	}
	for _, entry := range replicaUpdates {
		if table := memberTables[entry.region]; table != nil {
			for _, update := range entry.gsiRead {
				if !gsiExistsOnTable(table, update.IndexName) {
					return nil, ErrIndexNotFound
				}
			}
		}
	}

	hasCapacityMember := hasGlobalWriteUnits
	for _, update := range globalGSIUpdates {
		if update.ProvisionedWriteCapacityUnits != nil {
			hasCapacityMember = true
		}
	}
	for _, entry := range replicaUpdates {
		if entry.readUnits != nil || len(entry.gsiRead) > 0 {
			hasCapacityMember = true
		}
	}
	resultingMode := resultingGlobalBillingMode(globalBillingMode, hasCapacityMember)
	if resultingMode == string(dbstore.BillingModePayPerRequest) && hasCapacityMember {
		// Capacity members are provisioned-mode members: carrying them
		// into an on-demand result is the same contradiction the table
		// plane's billing-mode consistency check refuses.
		return nil, ErrInvalidParameter
	}

	if resultingMode != "" {
		for region, table := range memberTables {
			if err := validateReplicaTableCapacityResult(table, resultingMode, hasGlobalWriteUnits, globalWriteUnits, globalGSIUpdates, readUpdateByRegion[region]); err != nil {
				return nil, err
			}
		}
	}

	globalTable, err := store.GlobalTables().Update(globalTableName, func(gt *dbstore.GlobalTable) error {
		if globalBillingMode != "" || hasGlobalWriteUnits {
			for _, replica := range gt.ReplicationGroup {
				if globalBillingMode != "" {
					replica.BillingMode = dbstore.BillingMode(globalBillingMode)
				}
				if hasGlobalWriteUnits {
					replica.ProvisionedWriteCapacityUnits = globalWriteUnits
				}
			}
		}
		if globalWriteAS != nil {
			gt.WriteAutoScalingSettings = globalWriteAS
		}
		if globalGSIUpdates != nil {
			gt.GlobalSecondaryIndexWriteSettings = mergeIndexSettingsLists(gt.GlobalSecondaryIndexWriteSettings, globalGSIUpdates)
		}

		for _, entry := range replicaUpdates {
			// A settings update naming a region that is no longer part
			// of the global table is the documented replica-not-found
			// error rather than a silently ignored entry.
			matched := false
			for _, replica := range gt.ReplicationGroup {
				if replica.RegionName == entry.region {
					matched = true
					if entry.readUnits != nil {
						replica.ProvisionedReadCapacityUnits = *entry.readUnits
					}
					if entry.readAS != nil {
						replica.ReadAutoScalingSettings = entry.readAS
					}
					if entry.gsiRead != nil {
						replica.GlobalSecondaryIndexReadSettings = mergeIndexSettingsLists(replica.GlobalSecondaryIndexReadSettings, entry.gsiRead)
					}
					if entry.tableClass != "" {
						replica.TableClass = entry.tableClass
						now := time.Now().UTC()
						replica.TableClassLastUpdated = &now
					}
					break
				}
			}
			if !matched {
				return ErrReplicaNotFound
			}
		}
		return nil
	})
	if err != nil {
		if commonstore.IsNotFound(err) {
			return nil, ErrGlobalTableNotFound
		}
		return nil, err
	}

	// The capacity settings propagate to every member region's table
	// record — the same state the global record echoes. Membership changes
	// committed by the update above are followed here: the loop walks the
	// committed group. The propagation is best-effort: the global record
	// has already committed, so aborting on one region's failure cannot
	// roll anything back — it would only abandon the remaining regions and
	// report an error for a request whose committed state stands. Every
	// region is attempted, a failure is logged naming the region (the
	// DescribeTable/DescribeGlobalTableSettings divergence is observable
	// and repairable — re-issuing the same settings update re-propagates
	// to the lagging members as a no-op for the current ones), and the
	// response reports the committed group.
	if resultingMode != "" {
		for _, replica := range globalTable.ReplicationGroup {
			regionStore, err := s.GetCachedStoreForRegion(replica.RegionName)
			if err != nil {
				logs.Warn("DynamoDB: global table capacity propagation could not open the member region's store",
					logs.String("globalTable", globalTableName),
					logs.String("region", replica.RegionName),
					logs.Err(err))
				continue
			}
			entry := readUpdateByRegion[replica.RegionName]
			if _, err := regionStore.Tables().Update(globalTableName, func(table *dbstore.Table) error {
				return applyGlobalCapacityToTable(table, resultingMode, hasGlobalWriteUnits, globalWriteUnits, globalGSIUpdates, entry)
			}); err != nil {
				logs.Warn("DynamoDB: global table capacity propagation failed for a member region",
					logs.String("globalTable", globalTableName),
					logs.String("region", replica.RegionName),
					logs.Err(err))
			}
		}
	}

	return map[string]interface{}{
		"GlobalTableName": globalTable.GlobalTableName,
		"ReplicaSettings": buildGlobalTableReplicaSettings(globalTable),
	}, nil
}

// validateReplicaTableCapacityResult checks the state the settings update
// would leave one member table in before anything commits: a switch to
// on-demand passes the table plane's rolling limit, and a provisioned
// result carries a complete capacity pair within the table capacity caps —
// the base table's write side from the global setting and read side from
// this request's replica setting or the existing pair, and likewise for
// every index the request or the switch leaves provisioned. The
// documentation fixes the switch limits and the pair requirement on the
// table plane; this request applies the same state, so it passes the same
// rules — including the caps validateProvisionedThroughputValues would
// enforce at propagation, caught here before the global record commits.
func validateReplicaTableCapacityResult(table *dbstore.Table, resultingMode string, hasWriteUnits bool, writeUnits int64, globalGSIUpdates []dbstore.IndexAutoScalingSettings, entry replicaSettingsUpdate) error {
	newMode := dbstore.BillingMode(resultingMode)
	switching := newMode != table.BillingMode
	if switching && newMode == dbstore.BillingModePayPerRequest && table.BillingMode == dbstore.BillingModeProvisioned {
		// Only the provisioned-to-on-demand direction is rate limited
		// ("You can switch tables from on-demand mode to provisioned
		// capacity mode at any time"), and the window is rolling: the
		// stamps older than 24 hours drop out before the count is
		// taken. The propagation that lands the switch on the member
		// record re-checks the same quota under the region's lock.
		if len(windowBillingModeSwitches(table.BillingModeSwitches, time.Now().UTC())) >= dbstore.BillingModeSwitchesPerDay {
			return ErrInvalidParameter
		}
	}

	if newMode != dbstore.BillingModeProvisioned {
		return nil
	}

	// The resulting base pair: write from the global setting or the
	// existing pair, read from this request's replica setting or the
	// existing pair.
	write := int64(0)
	if hasWriteUnits {
		write = writeUnits
	} else if table.ProvisionedThroughput != nil {
		write = table.ProvisionedThroughput.WriteCapacityUnits
	}
	read := int64(0)
	if entry.readUnits != nil {
		read = *entry.readUnits
	} else if table.ProvisionedThroughput != nil {
		read = table.ProvisionedThroughput.ReadCapacityUnits
	}
	if !capacityPairWithinTableCaps(write, read) {
		return ErrInvalidParameter
	}

	// A provisioned table's every index carries a complete pair: an index
	// the request touches takes its write side from the global write
	// setting and its read side from this request's replica index setting;
	// an index the request does not touch keeps its existing pair — which
	// an index gained by a mode switch from on-demand never has, so a
	// switching table must be given settings for every index it holds.
	for _, gsi := range table.GlobalSecondaryIndexes {
		write := int64(0)
		read := int64(0)
		if gsi.ProvisionedThroughput != nil {
			write = gsi.ProvisionedThroughput.WriteCapacityUnits
			read = gsi.ProvisionedThroughput.ReadCapacityUnits
		}
		for _, update := range globalGSIUpdates {
			if update.IndexName == gsi.IndexName && update.ProvisionedWriteCapacityUnits != nil {
				write = *update.ProvisionedWriteCapacityUnits
			}
		}
		for _, update := range entry.gsiRead {
			if update.IndexName == gsi.IndexName && update.ProvisionedReadCapacityUnits != nil {
				read = *update.ProvisionedReadCapacityUnits
			}
		}
		if !capacityPairWithinTableCaps(write, read) {
			return ErrInvalidParameter
		}
	}
	return nil
}

// capacityPairWithinTableCaps reports whether a composed provisioned pair is
// complete and inside the table capacity caps — the pair rule
// validateProvisionedThroughputValues states on the table plane, evaluated
// on the state the settings request would leave the member table in.
func capacityPairWithinTableCaps(write, read int64) bool {
	return write >= 1 && read >= 1 &&
		write <= dbstore.TableMaxWriteCapacityUnits &&
		read <= dbstore.TableMaxReadCapacityUnits
}

// applyGlobalCapacityToTable writes the resolved capacity state onto one
// member region's table record: the billing mode (with the on-demand
// switch discarding provisioned settings, the same semantics the table
// plane's own update applies), the global write units, the per-index write
// units, and this region's replica read units. The post-state passes the
// same consistency checks the table plane enforces.
func applyGlobalCapacityToTable(table *dbstore.Table, resultingMode string, hasWriteUnits bool, writeUnits int64, globalGSIUpdates []dbstore.IndexAutoScalingSettings, entry replicaSettingsUpdate) error {
	newMode := dbstore.BillingMode(resultingMode)
	// The switch lands through the shared helper, so the quota validated
	// before the global record committed is re-enforced here under the
	// region's own lock: a member table that crossed the window between
	// the validation read and this propagation is refused rather than
	// stamped past the quota, and the on-demand discard carries the same
	// semantics the table plane's own update applies.
	if err := registerBillingModeSwitch(table, newMode); err != nil {
		return err
	}

	if newMode == dbstore.BillingModeProvisioned {
		if table.ProvisionedThroughput == nil {
			table.ProvisionedThroughput = &dbstore.ProvisionedThroughput{}
		}
		if hasWriteUnits {
			table.ProvisionedThroughput.WriteCapacityUnits = writeUnits
		}
		if entry.readUnits != nil {
			table.ProvisionedThroughput.ReadCapacityUnits = *entry.readUnits
		}
		for _, update := range globalGSIUpdates {
			if update.ProvisionedWriteCapacityUnits == nil {
				continue
			}
			for _, gsi := range table.GlobalSecondaryIndexes {
				if gsi.IndexName != update.IndexName {
					continue
				}
				if gsi.ProvisionedThroughput == nil {
					gsi.ProvisionedThroughput = &dbstore.ProvisionedThroughput{}
				}
				gsi.ProvisionedThroughput.WriteCapacityUnits = *update.ProvisionedWriteCapacityUnits
			}
		}
		for _, update := range entry.gsiRead {
			if update.ProvisionedReadCapacityUnits == nil {
				continue
			}
			for _, gsi := range table.GlobalSecondaryIndexes {
				if gsi.IndexName != update.IndexName {
					continue
				}
				if gsi.ProvisionedThroughput == nil {
					gsi.ProvisionedThroughput = &dbstore.ProvisionedThroughput{}
				}
				gsi.ProvisionedThroughput.ReadCapacityUnits = *update.ProvisionedReadCapacityUnits
			}
		}
	}

	if !validateBillingModeConsistency(table.BillingMode, table.ProvisionedThroughput) {
		return ErrInvalidParameter
	}
	if !validateIndexThroughputModeConsistency(table.BillingMode, table.GlobalSecondaryIndexes) {
		return ErrInvalidParameter
	}
	return nil
}
