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
	instance := &neptunestore.DBInstance{
		DBInstanceIdentifier:             id,
		DBClusterIdentifier:              clusterID,
		Engine:                           request.GetStringParam(params, "Engine"),
		EngineVersion:                    request.GetStringParam(params, "EngineVersion"),
		DBInstanceClass:                  request.GetStringParam(params, "DBInstanceClass"),
		DBInstanceStatus:                 "available",
		AvailabilityZone:                 request.GetStringParam(params, "AvailabilityZone"),
		DBParameterGroupName:             request.GetStringParam(params, "DBParameterGroupName"),
		DBSubnetGroupName:                request.GetStringParam(params, "DBSubnetGroupName"),
		IAMDatabaseAuthenticationEnabled: request.GetBoolParam(params, "EnableIAMDatabaseAuthentication"),
		PubliclyAccessible:               request.GetBoolParam(params, "PubliclyAccessible"),
		AutoMinorVersionUpgrade:          request.GetBoolParam(params, "AutoMinorVersionUpgrade"),
		InstanceCreateTime:               &now,
		AccountID:                        reqCtx.GetAccountID(),
		Region:                           reqCtx.GetRegion(),
		DBInstanceArn:                    arnutil.NewARNBuilder(reqCtx.GetAccountID(), reqCtx.GetRegion()).RDS().DBInstance(id),
	}

	if err := store.CreateInstance(instance); err != nil {
		return nil, translateStoreError(err)
	}

	if clusterID != "" {
		if s.porter != nil {
			if port, err := s.porter.GetPort(clusterID); err == nil && port > 0 {
				instance.Endpoint = &neptunestore.Endpoint{
					Address: fmt.Sprintf("%s.%s.%s.neptune.amazonaws.com", id, s.accountID, s.region),
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

	if err := store.DeleteInstance(id); err != nil {
		return nil, translateStoreError(err)
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
