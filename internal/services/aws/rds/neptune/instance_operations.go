package neptune

import (
	"context"
	"fmt"
	"net/http"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/protocol"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	rdssvc "vorpalstacks/internal/services/aws/rds"
	neptunestore "vorpalstacks/internal/store/aws/rds/neptune"
	arnutil "vorpalstacks/internal/utils/aws/arn"
	"vorpalstacks/internal/utils/aws/types"
)

// CreateDBInstance creates a new Neptune DB instance within a cluster.
func (s *NeptuneService) CreateDBInstance(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	id := request.GetStringParam(params, "DBInstanceIdentifier")
	if id == "" {
		return nil, awserrors.NewMissingParameter("DBInstanceIdentifier is required")
	}
	if err := rdssvc.ValidateDBInstanceIdentifier(id); err != nil {
		return nil, awserrors.NewAWSError("InvalidParameterValue", err.Error(), http.StatusBadRequest)
	}
	if class := request.GetStringParam(params, "DBInstanceClass"); class != "" {
		if err := rdssvc.ValidateDBInstanceClass(class); err != nil {
			return nil, awserrors.NewAWSError("InvalidParameterValue", err.Error(), http.StatusBadRequest)
		}
	}
	engineType := request.GetStringParam(params, "Engine")
	if engineType == "" {
		engineType = "neptune"
	}
	if err := rdssvc.ValidateEngine(engineType); err != nil {
		return nil, awserrors.NewAWSError("InvalidParameterValue", err.Error(), http.StatusBadRequest)
	}
	engineVersion := request.GetStringParam(params, "EngineVersion")
	if engineVersion == "" {
		engineVersion = rdssvc.DefaultEngineVersion(engineType)
	}
	if err := rdssvc.ValidateEngineVersion(engineType, engineVersion); err != nil {
		return nil, awserrors.NewAWSError("InvalidParameterValue", err.Error(), http.StatusBadRequest)
	}

	if port := int32(request.GetIntParam(params, "Port")); port > 0 {
		if err := rdssvc.ValidatePort(port); err != nil {
			return nil, awserrors.NewAWSError("InvalidParameterValue", err.Error(), http.StatusBadRequest)
		}
	}
	if mi := int32(request.GetIntParam(params, "MonitoringInterval")); mi > 0 {
		if err := rdssvc.ValidateMonitoringInterval(mi); err != nil {
			return nil, awserrors.NewAWSError("InvalidParameterValue", err.Error(), http.StatusBadRequest)
		}
	}
	if st := request.GetStringParam(params, "StorageType"); st != "" {
		if err := rdssvc.ValidateStorageType(st, engineType); err != nil {
			return nil, awserrors.NewAWSError("InvalidParameterValue", err.Error(), http.StatusBadRequest)
		}
	}
	if brp := int32(request.GetIntParam(params, "BackupRetentionPeriod")); brp > 0 {
		if err := rdssvc.ValidateBackupRetentionPeriod(brp); err != nil {
			return nil, awserrors.NewAWSError("InvalidParameterValue", err.Error(), http.StatusBadRequest)
		}
	}
	if as := int32(request.GetIntParam(params, "AllocatedStorage")); as > 0 {
		if err := rdssvc.ValidateAllocatedStorage(as, engineType); err != nil {
			return nil, awserrors.NewAWSError("InvalidParameterValue", err.Error(), http.StatusBadRequest)
		}
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	clusterID := request.GetStringParam(params, "DBClusterIdentifier")
	if clusterID != "" {
		if _, err := store.GetCluster(clusterID); err != nil {
			return nil, awserrors.NewAWSError("DBClusterNotFoundFault", fmt.Sprintf("DBCluster %s not found", clusterID), http.StatusNotFound)
		}
	}

	now := time.Now()

	// Accept MasterUserPassword and store as bcrypt hash (write-only).
	masterPassword := request.GetStringParam(params, "MasterUserPassword")
	masterPasswordHash, err := hashMasterPassword(masterPassword)
	if err != nil {
		return nil, awserrors.NewAWSError("InvalidParameterValue", err.Error(), http.StatusBadRequest)
	}

	instance := &neptunestore.DBInstance{
		DBInstanceIdentifier:             id,
		DBClusterIdentifier:              clusterID,
		Engine:                           engineType,
		EngineVersion:                    engineVersion,
		DBInstanceClass:                  request.GetStringParam(params, "DBInstanceClass"),
		DBInstanceStatus:                 "available",
		AvailabilityZone:                 request.GetStringParam(params, "AvailabilityZone"),
		DBParameterGroupName:             request.GetStringParam(params, "DBParameterGroupName"),
		DBSubnetGroupName:                request.GetStringParam(params, "DBSubnetGroupName"),
		IAMDatabaseAuthenticationEnabled: request.GetBoolParam(params, "EnableIAMDatabaseAuthentication"),
		PubliclyAccessible:               request.GetBoolParam(params, "PubliclyAccessible"),
		AutoMinorVersionUpgrade:          request.GetBoolParam(params, "AutoMinorVersionUpgrade"),
		Port:                             int32(request.GetIntParam(params, "Port")),
		InstanceCreateTime:               &now,
		AccountID:                        reqCtx.GetAccountID(),
		Region:                           reqCtx.GetRegion(),
		DBInstanceArn:                    arnutil.NewARNBuilder(reqCtx.GetAccountID(), reqCtx.GetRegion()).RDS().DBInstance(id),
		MasterUserPasswordHash:           masterPasswordHash,
	}

	if err := store.CreateInstance(instance); err != nil {
		return nil, translateStoreError(err)
	}

	// Register the instance in the cluster's DBClusterMembers list
	// so DescribeDBClusters returns a complete member list.
	if clusterID != "" {
		if cluster, err := store.GetCluster(clusterID); err == nil {
			isWriter := len(cluster.DBClusterMembers) == 0
			cluster.DBClusterMembers = append(cluster.DBClusterMembers, neptunestore.DBClusterMember{
				DBInstanceIdentifier: id,
				IsClusterWriter:      isWriter,
				PromotionTier:        0,
			})
			if err := store.UpdateCluster(cluster); err != nil {
				logs.Warn("failed to add instance to cluster members", logs.String("instance", id), logs.Err(err))
			}
		}
	}

	if clusterID != "" {
		if s.porter != nil {
			if port, err := s.porter.GetPort(clusterID); err == nil && port > 0 {
				instance.Endpoint = &neptunestore.Endpoint{
					Address: s.endpointAddressFor(id, instance.Engine),
					Port:    port,
				}
				if err := store.UpdateInstance(instance); err != nil {
					logs.Warn("failed to persist instance endpoint", logs.String("instance", id), logs.Err(err))
				}
			}
		}
	} else {
		// Standalone instance (no cluster) — open the engine directly
		// so the instance gets its own port and endpoint. This is the
		// normal path for MySQL RDS instances.
		engineType := instance.Engine
		if engineType == "" {
			engineType = "neptune"
		}
		if eng := s.engineFor(engineType); eng != nil {
			if port, err := eng.Open(reqCtx.GetRegion(), id); err != nil {
				logs.Warn("failed to open instance engine", logs.String("instance", id), logs.Err(err))
			} else {
				instance.Endpoint = &neptunestore.Endpoint{
					Address: s.endpointAddressFor(id, engineType),
					Port:    port,
				}
				if err := store.UpdateInstance(instance); err != nil {
					logs.Warn("failed to persist instance endpoint", logs.String("instance", id), logs.Err(err))
				}
			}
		}
	}

	recordEvent(store, "db-instance", id, instance.DBInstanceArn,
		fmt.Sprintf("DB instance %s created", id), []string{"creation"})

	if tagList := getNeptuneTagList(params); len(tagList) > 0 {
		storeTags := make([]types.Tag, 0, len(tagList))
		for _, t := range tagList {
			key, _ := t["Key"].(string)
			value, _ := t["Value"].(string)
			if key != "" {
				storeTags = append(storeTags, types.Tag{Key: key, Value: value})
			}
		}
		if err := store.AddTags(instance.DBInstanceArn, storeTags); err != nil {
			logs.Warn("failed to tag instance on create", logs.String("instance", id), logs.Err(err))
		}
	}

	return map[string]interface{}{
		"DBInstance": instance,
	}, nil
}

// DeleteDBInstance deletes the specified Neptune DB instance.
func (s *NeptuneService) DeleteDBInstance(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	id := request.GetStringParam(params, "DBInstanceIdentifier")
	if id == "" {
		return nil, awserrors.NewMissingParameter("DBInstanceIdentifier is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	instance, err := store.GetInstance(id)
	if err != nil {
		return nil, translateStoreError(err)
	}

	instance.DBInstanceStatus = "deleting"
	if err := store.UpdateInstance(instance); err != nil {
		return nil, translateStoreError(err)
	}

	// Close the engine for standalone instances (cluster-based instances
	// are cleaned up when the cluster's engine is closed).
	if instance.DBClusterIdentifier == "" {
		engineType := instance.Engine
		if engineType == "" {
			engineType = "neptune"
		}
		if eng := s.engineFor(engineType); eng != nil {
			if err := eng.Close(id); err != nil {
				logs.Warn("failed to close instance engine on delete", logs.String("instance", id), logs.Err(err))
			}
		}
	}

	if err := store.DeleteInstance(id); err != nil {
		return nil, translateStoreError(err)
	}

	// Remove the instance from the cluster's DBClusterMembers list.
	if instance.DBClusterIdentifier != "" {
		if cluster, err := store.GetCluster(instance.DBClusterIdentifier); err == nil {
			filtered := make([]neptunestore.DBClusterMember, 0, len(cluster.DBClusterMembers))
			for _, mem := range cluster.DBClusterMembers {
				if mem.DBInstanceIdentifier != id {
					filtered = append(filtered, mem)
				}
			}
			cluster.DBClusterMembers = filtered
			if err := store.UpdateCluster(cluster); err != nil {
				logs.Warn("failed to remove instance from cluster members", logs.String("instance", id), logs.Err(err))
			}
		}
	}

	removeTagsForResource(store, instance.DBInstanceArn)

	recordEvent(store, "db-instance", id, instance.DBInstanceArn,
		fmt.Sprintf("DB instance %s deleted", id), []string{"deletion"})

	return map[string]interface{}{
		"DBInstance": instance,
	}, nil
}

// ModifyDBInstance updates the configuration of the specified DB instance.
func (s *NeptuneService) ModifyDBInstance(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	id := request.GetStringParam(params, "DBInstanceIdentifier")
	if id == "" {
		return nil, awserrors.NewMissingParameter("DBInstanceIdentifier is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	instance, err := store.GetInstance(id)
	if err != nil {
		return nil, translateStoreError(err)
	}

	if v := request.GetStringParam(params, "DBInstanceClass"); v != "" {
		instance.DBInstanceClass = v
	}
	if v := request.GetStringParam(params, "EngineVersion"); v != "" {
		instance.EngineVersion = v
	}
	if v := request.GetStringParam(params, "DBParameterGroupName"); v != "" {
		instance.DBParameterGroupName = v
	}
	if v := request.GetStringParam(params, "PreferredMaintenanceWindow"); v != "" {
		instance.PreferredMaintenanceWindow = v
	}
	if request.HasParam(params, "PubliclyAccessible") {
		instance.PubliclyAccessible = request.GetBoolParam(params, "PubliclyAccessible")
	}
	if request.HasParam(params, "AutoMinorVersionUpgrade") {
		instance.AutoMinorVersionUpgrade = request.GetBoolParam(params, "AutoMinorVersionUpgrade")
	}
	if request.HasParam(params, "EnableIAMDatabaseAuthentication") {
		instance.IAMDatabaseAuthenticationEnabled = request.GetBoolParam(params, "EnableIAMDatabaseAuthentication")
	}
	if request.HasParam(params, "CopyTagsToSnapshot") {
		instance.CopyTagsToSnapshot = request.GetBoolParam(params, "CopyTagsToSnapshot")
	}

	if err := store.UpdateInstance(instance); err != nil {
		return nil, translateStoreError(err)
	}

	// Support NewDBInstanceIdentifier (instance rename).
	// Follows the same create-new + delete-old pattern as ModifyDBCluster.
	if newID := request.GetStringParam(params, "NewDBInstanceIdentifier"); newID != "" && newID != id {
		oldArn := instance.DBInstanceArn
		instance.DBInstanceIdentifier = newID
		instance.DBInstanceArn = arnutil.NewARNBuilder(instance.AccountID, instance.Region).RDS().DBInstance(newID)
		if err := store.CreateInstance(instance); err != nil {
			instance.DBInstanceIdentifier = id
			instance.DBInstanceArn = oldArn
			return nil, translateStoreError(err)
		}
		// Copy tags from old to new ARN.
		if tags, err := store.GetTags(oldArn); err == nil && len(tags) > 0 {
			store.AddTags(instance.DBInstanceArn, tags)
		}
		if err := store.DeleteInstance(id); err != nil {
			logs.Warn("failed to delete old instance after rename", logs.String("oldID", id), logs.Err(err))
		}
		removeTagsForResource(store, oldArn)
		// Update DBClusterMembers if the instance belongs to a cluster.
		if instance.DBClusterIdentifier != "" {
			if cluster, err := store.GetCluster(instance.DBClusterIdentifier); err == nil {
				for i, mem := range cluster.DBClusterMembers {
					if mem.DBInstanceIdentifier == id {
						cluster.DBClusterMembers[i].DBInstanceIdentifier = newID
					}
				}
				store.UpdateCluster(cluster)
			}
		}
	}

	return map[string]interface{}{
		"DBInstance": instance,
	}, nil
}

// RestoreDBInstanceFromDBSnapshot creates a new DB instance from a DB instance
// snapshot.
func (s *NeptuneService) RestoreDBInstanceFromDBSnapshot(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	instanceID := request.GetStringParam(params, "DBInstanceIdentifier")
	if instanceID == "" {
		return nil, awserrors.NewMissingParameter("DBInstanceIdentifier is required")
	}
	snapshotID := request.GetStringParam(params, "DBSnapshotIdentifier")
	if snapshotID == "" {
		return nil, awserrors.NewMissingParameter("DBSnapshotIdentifier is required")
	}
	if err := rdssvc.ValidateDBSnapshotIdentifier(snapshotID); err != nil {
		return nil, awserrors.NewAWSError("InvalidParameterValue", err.Error(), http.StatusBadRequest)
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	snapshot, err := store.GetInstanceSnapshot(snapshotID)
	if err != nil {
		return nil, translateStoreError(err)
	}

	now := time.Now()
	instance := &neptunestore.DBInstance{
		DBInstanceIdentifier:             instanceID,
		Engine:                           snapshot.Engine,
		EngineVersion:                    snapshot.EngineVersion,
		DBInstanceClass:                  request.GetStringParam(params, "DBInstanceClass"),
		DBInstanceStatus:                 "available",
		AvailabilityZone:                 request.GetStringParam(params, "AvailabilityZone"),
		DBParameterGroupName:             request.GetStringParam(params, "DBParameterGroupName"),
		DBSubnetGroupName:                request.GetStringParam(params, "DBSubnetGroupName"),
		IAMDatabaseAuthenticationEnabled: request.GetBoolParam(params, "EnableIAMDatabaseAuthentication"),
		PubliclyAccessible:               request.GetBoolParam(params, "PubliclyAccessible"),
		AutoMinorVersionUpgrade:          request.GetBoolParam(params, "AutoMinorVersionUpgrade"),
		CopyTagsToSnapshot:               request.GetBoolParam(params, "CopyTagsToSnapshot"),
		InstanceCreateTime:               &now,
		AccountID:                        reqCtx.GetAccountID(),
		Region:                           reqCtx.GetRegion(),
		DBInstanceArn:                    arnutil.NewARNBuilder(reqCtx.GetAccountID(), reqCtx.GetRegion()).RDS().DBInstance(instanceID),
	}

	if instance.DBInstanceClass == "" {
		instance.DBInstanceClass = "db.r5.large"
	}

	if err := store.CreateInstance(instance); err != nil {
		return nil, translateStoreError(err)
	}

	// Open the engine for the restored instance and restore row-level
	// data from the snapshot when the snapshot was taken from a MySQL
	// instance. Without this, the restored instance would have no
	// running engine and no user data.
	engineType := snapshot.Engine
	if engineType == "" {
		engineType = "neptune"
	}
	if eng := s.engineFor(engineType); eng != nil {
		if port, err := eng.Open(reqCtx.GetRegion(), instanceID); err != nil {
			logs.Warn("failed to open engine for restored instance",
				logs.String("instance", instanceID),
				logs.Err(err))
		} else {
			instance.Endpoint = &neptunestore.Endpoint{
				Address: s.endpointAddressFor(instanceID, engineType),
				Port:    port,
			}
			if err := store.UpdateInstance(instance); err != nil {
				logs.Warn("failed to persist restored instance endpoint",
					logs.String("instance", instanceID),
					logs.Err(err))
			}
		}
	}
	if s.snapOp != nil && snapshot.Engine == "mysql" {
		if err := s.snapOp.RestoreData(snapshotID, instanceID); err != nil {
			logs.Warn("neptune: RestoreData failed",
				logs.String("snapshot", snapshotID),
				logs.String("instance", instanceID),
				logs.Err(err))
		}
	}

	return map[string]interface{}{
		"DBInstance": instance,
	}, nil
}

// DescribeDBInstances returns information about the specified DB instance or lists
// all instances when no identifier is provided.
func (s *NeptuneService) DescribeDBInstances(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	instanceID := request.GetStringParam(params, "DBInstanceIdentifier")
	if instanceID != "" {
		instance, err := store.GetInstance(instanceID)
		if err != nil {
			return nil, translateStoreError(err)
		}
		return map[string]interface{}{
			"DBInstances": protocol.XMLElements{ElementName: "DBInstance", Items: []interface{}{instance}},
		}, nil
	}

	instances, err := store.ListInstances()
	if err != nil {
		return nil, translateStoreError(err)
	}

	items := make([]interface{}, 0, len(instances))
	for _, i := range instances {
		items = append(items, i)
	}

	marker := request.GetStringParam(params, "Marker")
	maxRecords := request.GetIntParam(params, "MaxRecords")
	resultItems, nextMarker, isTruncated := paginateItems(items, marker, maxRecords, func(item interface{}) string {
		return item.(*neptunestore.DBInstance).DBInstanceIdentifier
	})

	result := map[string]interface{}{
		"DBInstances": protocol.XMLElements{ElementName: "DBInstance", Items: resultItems},
	}
	if isTruncated {
		result["Marker"] = nextMarker
	}
	return result, nil
}

// RestoreDBInstanceToPointInTime restores a DB instance to a specified point
// in time from a source instance.
func (s *NeptuneService) RestoreDBInstanceToPointInTime(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	targetID := request.GetStringParam(params, "TargetDBInstanceIdentifier")
	if targetID == "" {
		return nil, awserrors.NewMissingParameter("TargetDBInstanceIdentifier is required")
	}
	sourceID := request.GetStringParam(params, "SourceDBInstanceIdentifier")
	if sourceID == "" {
		return nil, awserrors.NewMissingParameter("SourceDBInstanceIdentifier is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	source, err := store.GetInstance(sourceID)
	if err != nil {
		return nil, translateStoreError(err)
	}

	now := time.Now()
	instance := &neptunestore.DBInstance{
		DBInstanceIdentifier:             targetID,
		Engine:                           source.Engine,
		EngineVersion:                    source.EngineVersion,
		DBInstanceClass:                  source.DBInstanceClass,
		DBInstanceStatus:                 "available",
		AvailabilityZone:                 source.AvailabilityZone,
		DBParameterGroupName:             source.DBParameterGroupName,
		DBSubnetGroupName:                source.DBSubnetGroupName,
		IAMDatabaseAuthenticationEnabled: source.IAMDatabaseAuthenticationEnabled,
		PubliclyAccessible:               source.PubliclyAccessible,
		AutoMinorVersionUpgrade:          source.AutoMinorVersionUpgrade,
		InstanceCreateTime:               &now,
		AccountID:                        reqCtx.GetAccountID(),
		Region:                           reqCtx.GetRegion(),
		DBInstanceArn:                    arnutil.NewARNBuilder(reqCtx.GetAccountID(), reqCtx.GetRegion()).RDS().DBInstance(targetID),
	}

	if v := request.GetStringParam(params, "DBInstanceClass"); v != "" {
		instance.DBInstanceClass = v
	}
	if v := request.GetStringParam(params, "AvailabilityZone"); v != "" {
		instance.AvailabilityZone = v
	}
	if v := request.GetStringParam(params, "DBSubnetGroupName"); v != "" {
		instance.DBSubnetGroupName = v
	}

	if err := store.CreateInstance(instance); err != nil {
		return nil, translateStoreError(err)
	}

	// Open the engine for the restored instance so it can serve
	// queries. Without this, the restored instance exists in metadata
	// only and has no running engine — the same gap that would exist
	// if RestoreDBInstanceFromDBSnapshot omitted the engine open.
	engineType := source.Engine
	if engineType == "" {
		engineType = "neptune"
	}
	if eng := s.engineFor(engineType); eng != nil {
		if port, err := eng.Open(reqCtx.GetRegion(), targetID); err != nil {
			logs.Warn("failed to open engine for PITR restored instance",
				logs.String("instance", targetID),
				logs.Err(err))
		} else {
			instance.Endpoint = &neptunestore.Endpoint{
				Address: s.endpointAddressFor(targetID, engineType),
				Port:    port,
			}
			if err := store.UpdateInstance(instance); err != nil {
				logs.Warn("failed to persist PITR restored instance endpoint",
					logs.String("instance", targetID),
					logs.Err(err))
			}
		}
	}

	return map[string]interface{}{
		"DBInstance": instance,
	}, nil
}

// CreateDBInstanceReadReplica creates a read replica of an existing DB
// instance.
func (s *NeptuneService) CreateDBInstanceReadReplica(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	replicaID := request.GetStringParam(params, "DBInstanceIdentifier")
	if replicaID == "" {
		return nil, awserrors.NewMissingParameter("DBInstanceIdentifier is required")
	}
	sourceID := request.GetStringParam(params, "SourceDBInstanceIdentifier")
	if sourceID == "" {
		return nil, awserrors.NewMissingParameter("SourceDBInstanceIdentifier is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	source, err := store.GetInstance(sourceID)
	if err != nil {
		return nil, translateStoreError(err)
	}

	now := time.Now()
	instance := &neptunestore.DBInstance{
		DBInstanceIdentifier:             replicaID,
		DBClusterIdentifier:              source.DBClusterIdentifier,
		Engine:                           source.Engine,
		EngineVersion:                    source.EngineVersion,
		DBInstanceClass:                  source.DBInstanceClass,
		DBInstanceStatus:                 "available",
		AvailabilityZone:                 request.GetStringParam(params, "AvailabilityZone"),
		DBParameterGroupName:             source.DBParameterGroupName,
		DBSubnetGroupName:                source.DBSubnetGroupName,
		IAMDatabaseAuthenticationEnabled: source.IAMDatabaseAuthenticationEnabled,
		PubliclyAccessible:               source.PubliclyAccessible,
		AutoMinorVersionUpgrade:          source.AutoMinorVersionUpgrade,
		InstanceCreateTime:               &now,
		AccountID:                        reqCtx.GetAccountID(),
		Region:                           reqCtx.GetRegion(),
		DBInstanceArn:                    arnutil.NewARNBuilder(reqCtx.GetAccountID(), reqCtx.GetRegion()).RDS().DBInstance(replicaID),
	}

	if v := request.GetStringParam(params, "DBInstanceClass"); v != "" {
		instance.DBInstanceClass = v
	}

	if err := store.CreateInstance(instance); err != nil {
		return nil, translateStoreError(err)
	}

	return map[string]interface{}{
		"DBInstance": instance,
	}, nil
}

// PromoteReadReplica promotes a read replica DB instance to a standalone
// primary instance.
func (s *NeptuneService) PromoteReadReplica(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	id := request.GetStringParam(params, "DBInstanceIdentifier")
	if id == "" {
		return nil, awserrors.NewMissingParameter("DBInstanceIdentifier is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	instance, err := store.GetInstance(id)
	if err != nil {
		return nil, translateStoreError(err)
	}

	instance.DBInstanceStatus = "available"
	if err := store.UpdateInstance(instance); err != nil {
		return nil, translateStoreError(err)
	}

	return map[string]interface{}{
		"DBInstance": instance,
	}, nil
}

// RebootDBInstance reboots the specified DB instance.
func (s *NeptuneService) RebootDBInstance(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	id := request.GetStringParam(params, "DBInstanceIdentifier")
	if id == "" {
		return nil, awserrors.NewMissingParameter("DBInstanceIdentifier is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	instance, err := store.GetInstance(id)
	if err != nil {
		return nil, translateStoreError(err)
	}

	if instance.DBInstanceStatus != "available" {
		return nil, awserrors.NewAWSError("InvalidDBInstanceStateFault", fmt.Sprintf("instance %s is not in available state", id), http.StatusBadRequest)
	}

	instance.DBInstanceStatus = "available"
	if err := store.UpdateInstance(instance); err != nil {
		return nil, translateStoreError(err)
	}

	return map[string]interface{}{
		"DBInstance": instance,
	}, nil
}
