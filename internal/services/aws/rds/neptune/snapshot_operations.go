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
)

// CreateDBClusterSnapshot creates a new snapshot of the specified DB cluster.
func (s *NeptuneService) CreateDBClusterSnapshot(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	snapshotID := request.GetStringParam(params, "DBClusterSnapshotIdentifier")
	if snapshotID == "" {
		return nil, awserrors.NewMissingParameter("DBClusterSnapshotIdentifier is required")
	}
	if err := rdssvc.ValidateDBSnapshotIdentifier(snapshotID); err != nil {
		return nil, awserrors.NewAWSError("InvalidParameterValue", err.Error(), http.StatusBadRequest)
	}
	clusterID := request.GetStringParam(params, "DBClusterIdentifier")
	if clusterID == "" {
		return nil, awserrors.NewMissingParameter("DBClusterIdentifier is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	cluster, err := store.GetCluster(clusterID)
	if err != nil {
		return nil, translateStoreError(err)
	}

	now := time.Now()
	snapshot := &neptunestore.DBClusterSnapshot{
		DBClusterSnapshotIdentifier: snapshotID,
		DBClusterIdentifier:         clusterID,
		SnapshotCreateTime:          &now,
		Engine:                      cluster.Engine,
		EngineVersion:               cluster.EngineVersion,
		SnapshotType:                "manual",
		Status:                      "creating",
		Port:                        cluster.Port,
		StorageEncrypted:            cluster.StorageEncrypted,
		KmsKeyId:                    cluster.KmsKeyId,
		DBSnapshotArn:               arnutil.NewARNBuilder(reqCtx.GetAccountID(), reqCtx.GetRegion()).RDS().ClusterSnapshot(snapshotID),
		AccountID:                   reqCtx.GetAccountID(),
		Region:                      reqCtx.GetRegion(),
	}

	if err := store.CreateSnapshot(snapshot); err != nil {
		return nil, translateStoreError(err)
	}

	recordEvent(store, "db-snapshot", snapshotID, snapshot.DBSnapshotArn,
		fmt.Sprintf("DB cluster snapshot %s created", snapshotID), []string{"creation"})

	// State machine: synchronous transition from 'creating' to 'available'
	// with safety-net goroutine.
	snapshot.Status = "available"
	if err := store.UpdateSnapshot(snapshot); err != nil {
		logs.Warn("failed to transition snapshot to available", logs.String("snapshot", snapshotID), logs.Err(err))
	}
	s.scheduleTransition(reqCtx.GetRegion(), 500*time.Millisecond, func(st neptunestore.NeptuneStoreInterface) error {
		snap, err := st.GetSnapshot(snapshotID)
		if err != nil || snap.Status != "creating" {
			return nil
		}
		snap.Status = "available"
		return st.UpdateSnapshot(snap)
	})

	return map[string]interface{}{
		"DBClusterSnapshot": snapshot,
	}, nil
}

// DeleteDBClusterSnapshot deletes the specified DB cluster snapshot.
func (s *NeptuneService) DeleteDBClusterSnapshot(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	snapshotID := request.GetStringParam(params, "DBClusterSnapshotIdentifier")
	if snapshotID == "" {
		return nil, awserrors.NewMissingParameter("DBClusterSnapshotIdentifier is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	snapshot, err := store.GetSnapshot(snapshotID)
	if err != nil {
		return nil, translateStoreError(err)
	}

	if err := store.DeleteSnapshot(snapshotID); err != nil {
		return nil, translateStoreError(err)
	}

	removeTagsForResource(store, snapshot.DBSnapshotArn)

	recordEvent(store, "db-snapshot", snapshotID, snapshot.DBSnapshotArn,
		fmt.Sprintf("DB cluster snapshot %s deleted", snapshotID), []string{"deletion"})

	return map[string]interface{}{
		"DBClusterSnapshot": snapshot,
	}, nil
}

// DescribeDBClusterSnapshots returns information about the specified cluster
// snapshot or lists all snapshots when no identifier is provided.
func (s *NeptuneService) DescribeDBClusterSnapshots(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	snapshots, nextMarker, err := rdssvc.QueryClusterSnapshots(store, rdssvc.DescribeDBClusterSnapshotsInput{
		DBClusterSnapshotIdentifier: request.GetStringParam(params, "DBClusterSnapshotIdentifier"),
		DBClusterIdentifier:         request.GetStringParam(params, "DBClusterIdentifier"),
		SnapshotType:                request.GetStringParam(params, "SnapshotType"),
		Filters:                     nil,
		Marker:                      request.GetStringParam(params, "Marker"),
		MaxRecords:                  int32(request.GetIntParam(params, "MaxRecords")),
	})
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0, len(snapshots))
	for _, snap := range snapshots {
		items = append(items, snap)
	}

	result := map[string]interface{}{
		"DBClusterSnapshots": protocol.XMLElements{ElementName: "DBClusterSnapshot", Items: items},
	}
	if nextMarker != "" {
		result["Marker"] = nextMarker
	}
	return result, nil
}

// CopyDBClusterSnapshot creates a copy of the specified DB cluster snapshot.
func (s *NeptuneService) CopyDBClusterSnapshot(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	sourceID := request.GetStringParam(params, "SourceDBClusterSnapshotIdentifier")
	if sourceID == "" {
		return nil, awserrors.NewMissingParameter("SourceDBClusterSnapshotIdentifier is required")
	}
	targetID := request.GetStringParam(params, "TargetDBClusterSnapshotIdentifier")
	if targetID == "" {
		return nil, awserrors.NewMissingParameter("TargetDBClusterSnapshotIdentifier is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	source, err := store.GetSnapshot(sourceID)
	if err != nil {
		return nil, translateStoreError(err)
	}

	now := time.Now()
	copy := *source
	copy.DBClusterSnapshotIdentifier = targetID
	copy.SnapshotCreateTime = &now
	if source.ClusterCreateTime != nil {
		ct := *source.ClusterCreateTime
		copy.ClusterCreateTime = &ct
	}
	copy.DBSnapshotArn = arnutil.NewARNBuilder(reqCtx.GetAccountID(), reqCtx.GetRegion()).RDS().ClusterSnapshot(targetID)

	if err := store.CreateSnapshot(&copy); err != nil {
		return nil, translateStoreError(err)
	}

	return map[string]interface{}{
		"DBClusterSnapshot": &copy,
	}, nil
}

// DescribeDBClusterSnapshotAttributes returns the attributes of the specified
// DB cluster snapshot.
func (s *NeptuneService) DescribeDBClusterSnapshotAttributes(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	snapshotID := request.GetStringParam(params, "DBClusterSnapshotIdentifier")
	if snapshotID == "" {
		return nil, awserrors.NewMissingParameter("DBClusterSnapshotIdentifier is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	snap, err := store.GetSnapshot(snapshotID)
	if err != nil {
		return nil, translateStoreError(err)
	}

	attrValues := snap.RestoreAttributeValues
	if attrValues == nil {
		attrValues = []string{}
	}
	attrItems := make([]interface{}, 0, len(attrValues))
	for _, v := range attrValues {
		attrItems = append(attrItems, v)
	}

	return map[string]interface{}{
		"DBClusterSnapshotAttributesResult": map[string]interface{}{
			"DBClusterSnapshotIdentifier": snapshotID,
			"DBClusterSnapshotAttributes": protocol.XMLElements{ElementName: "DBClusterSnapshotAttribute", Items: []interface{}{
				map[string]interface{}{"AttributeName": "restore", "AttributeValues": protocol.XMLElements{ElementName: "AttributeValue", Items: attrItems}},
			}},
		},
	}, nil
}

// ModifyDBClusterSnapshotAttribute modifies an attribute of the specified DB
// cluster snapshot.
func (s *NeptuneService) ModifyDBClusterSnapshotAttribute(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	snapshotID := request.GetStringParam(params, "DBClusterSnapshotIdentifier")
	if snapshotID == "" {
		return nil, awserrors.NewMissingParameter("DBClusterSnapshotIdentifier is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	snap, err := store.GetSnapshot(snapshotID)
	if err != nil {
		return nil, translateStoreError(err)
	}

	attrName := request.GetStringParam(params, "AttributeName")
	valuesToAdd := getNeptuneStringList(params, "ValuesToAdd", "AttributeValue", "member")
	valuesToRemove := getNeptuneStringList(params, "ValuesToRemove", "AttributeValue", "member")

	if attrName == "" {
		attrName = "restore"
	}
	if attrName != "restore" {
		return nil, awserrors.NewValidationException(fmt.Sprintf("Unsupported attribute name '%s', only 'restore' is supported", attrName))
	}

	existing := make(map[string]bool, len(snap.RestoreAttributeValues))
	for _, v := range snap.RestoreAttributeValues {
		existing[v] = true
	}
	for _, v := range valuesToRemove {
		delete(existing, v)
	}
	for _, v := range valuesToAdd {
		existing[v] = true
	}
	result := make([]string, 0, len(existing))
	for v := range existing {
		result = append(result, v)
	}
	snap.RestoreAttributeValues = result
	if err := store.UpdateSnapshot(snap); err != nil {
		return nil, translateStoreError(err)
	}

	attrItems := make([]interface{}, 0, len(snap.RestoreAttributeValues))
	for _, v := range snap.RestoreAttributeValues {
		attrItems = append(attrItems, v)
	}

	return map[string]interface{}{
		"DBClusterSnapshotAttributesResult": map[string]interface{}{
			"DBClusterSnapshotIdentifier": snapshotID,
			"DBClusterSnapshotAttributes": protocol.XMLElements{ElementName: "DBClusterSnapshotAttribute", Items: []interface{}{
				map[string]interface{}{"AttributeName": attrName, "AttributeValues": protocol.XMLElements{ElementName: "AttributeValue", Items: attrItems}},
			}},
		},
	}, nil
}

// RestoreDBClusterFromSnapshot creates a new DB cluster from a DB cluster
// snapshot. Cluster snapshots are metadata-only — they do not capture
// row-level data.
func (s *NeptuneService) RestoreDBClusterFromSnapshot(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	clusterID := request.GetStringParam(params, "DBClusterIdentifier")
	if clusterID == "" {
		return nil, awserrors.NewMissingParameter("DBClusterIdentifier is required")
	}
	if err := rdssvc.ValidateDBClusterIdentifier(clusterID); err != nil {
		return nil, awserrors.NewAWSError("InvalidParameterValue", err.Error(), http.StatusBadRequest)
	}
	engine := request.GetStringParam(params, "Engine")
	if engine == "" {
		engine = "neptune"
	}
	if err := rdssvc.ValidateEngine(engine); err != nil {
		return nil, awserrors.NewAWSError("InvalidParameterValue", err.Error(), http.StatusBadRequest)
	}
	snapshotID := request.GetStringParam(params, "SnapshotIdentifier")
	if snapshotID == "" {
		return nil, awserrors.NewMissingParameter("SnapshotIdentifier is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	snapshot, err := store.GetSnapshot(snapshotID)
	if err != nil {
		return nil, translateStoreError(err)
	}

	engineVersion := request.GetStringParam(params, "EngineVersion")
	if engineVersion == "" {
		engineVersion = snapshot.EngineVersion
	}
	if engineVersion == "" {
		engineVersion = rdssvc.DefaultEngineVersion(engine)
	}
	if err := rdssvc.ValidateEngineVersion(engine, engineVersion); err != nil {
		return nil, awserrors.NewAWSError("InvalidParameterValue", err.Error(), http.StatusBadRequest)
	}

	now := time.Now()
	cluster := buildRestoredCluster(clusterID, engine, engineVersion, params, &now, reqCtx)
	if cluster.Port == 0 {
		cluster.Port = snapshot.Port
	}
	if !cluster.StorageEncrypted && snapshot.StorageEncrypted {
		cluster.StorageEncrypted = snapshot.StorageEncrypted
	}
	if cluster.KmsKeyId == "" && snapshot.KmsKeyId != "" {
		cluster.KmsKeyId = snapshot.KmsKeyId
	}

	if err := store.CreateCluster(cluster); err != nil {
		return nil, translateStoreError(err)
	}

	var enginePort int
	if eng := s.engineFor(cluster.Engine); eng != nil {
		if port, err := eng.Open(reqCtx.GetRegion(), clusterID); err != nil {
			logs.Warn("failed to open cluster engine on snapshot restore", logs.String("cluster", clusterID), logs.Err(err))
		} else {
			enginePort = port
		}
	}

	s.setClusterEndpoint(store, cluster, enginePort)

	return map[string]interface{}{
		"DBCluster": enrichClusterWithTags(store, cluster),
	}, nil
}

// RestoreDBClusterToPointInTime restores a DB cluster to a point in time from
// a source cluster.
func (s *NeptuneService) RestoreDBClusterToPointInTime(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	clusterID := request.GetStringParam(params, "DBClusterIdentifier")
	if clusterID == "" {
		return nil, awserrors.NewMissingParameter("DBClusterIdentifier is required")
	}
	sourceID := request.GetStringParam(params, "SourceDBClusterIdentifier")
	if sourceID == "" {
		return nil, awserrors.NewMissingParameter("SourceDBClusterIdentifier is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	source, err := store.GetCluster(sourceID)
	if err != nil {
		return nil, translateStoreError(err)
	}

	now := time.Now()
	cluster := buildRestoredCluster(clusterID, source.Engine, source.EngineVersion, params, &now, reqCtx)
	if cluster.Port == 0 {
		cluster.Port = source.Port
	}
	if !cluster.StorageEncrypted && source.StorageEncrypted {
		cluster.StorageEncrypted = source.StorageEncrypted
	}
	if cluster.KmsKeyId == "" && source.KmsKeyId != "" {
		cluster.KmsKeyId = source.KmsKeyId
	}
	if cluster.DBClusterParameterGroupName == "" && source.DBClusterParameterGroupName != "" {
		cluster.DBClusterParameterGroupName = source.DBClusterParameterGroupName
	}
	if cluster.DBSubnetGroupName == "" && source.DBSubnetGroupName != "" {
		cluster.DBSubnetGroupName = source.DBSubnetGroupName
	}
	if len(source.VpcSecurityGroupIds) > 0 {
		cluster.VpcSecurityGroupIds = source.VpcSecurityGroupIds
	}
	if len(source.EnabledCloudwatchLogsExports) > 0 {
		cluster.EnabledCloudwatchLogsExports = source.EnabledCloudwatchLogsExports
	}
	if source.BackupRetentionPeriod > 0 && request.GetIntParam(params, "BackupRetentionPeriod") == 0 {
		cluster.BackupRetentionPeriod = source.BackupRetentionPeriod
	}
	if source.PreferredBackupWindow != "" && request.GetStringParam(params, "PreferredBackupWindow") == "" {
		cluster.PreferredBackupWindow = source.PreferredBackupWindow
	}
	if source.PreferredMaintenanceWindow != "" && request.GetStringParam(params, "PreferredMaintenanceWindow") == "" {
		cluster.PreferredMaintenanceWindow = source.PreferredMaintenanceWindow
	}

	if err := store.CreateCluster(cluster); err != nil {
		return nil, translateStoreError(err)
	}

	var enginePort int
	if eng := s.engineFor(cluster.Engine); eng != nil {
		if port, err := eng.Open(reqCtx.GetRegion(), clusterID); err != nil {
			logs.Warn("failed to open cluster engine on restore", logs.String("cluster", clusterID), logs.Err(err))
		} else {
			enginePort = port
		}
	}

	s.setClusterEndpoint(store, cluster, enginePort)

	return map[string]interface{}{
		"DBCluster": enrichClusterWithTags(store, cluster),
	}, nil
}

// PromoteReadReplicaDBCluster promotes a read replica cluster to a standalone
// writable cluster.
func (s *NeptuneService) PromoteReadReplicaDBCluster(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	id := request.GetStringParam(params, "DBClusterIdentifier")
	if id == "" {
		return nil, awserrors.NewMissingParameter("DBClusterIdentifier is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	cluster, err := store.GetCluster(id)
	if err != nil {
		return nil, translateStoreError(err)
	}

	if cluster.ReplicationSourceIdentifier == "" {
		return nil, awserrors.NewAWSError("InvalidDBClusterStateFault", fmt.Sprintf("cluster %s is not a read replica", id), http.StatusBadRequest)
	}

	if cluster.Status != "available" {
		return nil, awserrors.NewAWSError("InvalidDBClusterStateFault", fmt.Sprintf("cluster %s is not in available state", id), http.StatusBadRequest)
	}

	// Remove the cluster from its global cluster's member list before
	// clearing the reference. Without this, the global cluster retains
	// a stale member ARN that points back at this now-promoted cluster.
	if cluster.GlobalClusterIdentifier != "" {
		if gc, err := store.GetGlobalCluster(cluster.GlobalClusterIdentifier); err == nil {
			filtered := make([]neptunestore.GlobalClusterMember, 0, len(gc.GlobalClusterMembers))
			for _, m := range gc.GlobalClusterMembers {
				if m.DBClusterArn != cluster.DBClusterArn {
					filtered = append(filtered, m)
				}
			}
			gc.GlobalClusterMembers = filtered
			if err := store.UpdateGlobalCluster(gc); err != nil {
				logs.Warn("failed to remove cluster from global cluster members on promote", logs.String("cluster", id), logs.Err(err))
			}
		}
	}

	cluster.ReplicationSourceIdentifier = ""
	cluster.GlobalClusterIdentifier = ""
	cluster.Status = "available"
	if err := store.UpdateCluster(cluster); err != nil {
		return nil, translateStoreError(err)
	}

	return map[string]interface{}{
		"DBCluster": cluster,
	}, nil
}
