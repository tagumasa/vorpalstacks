package neptune

import (
	"context"
	"fmt"
	"time"

	"vorpalstacks/internal/services/aws/common/protocol"
	"vorpalstacks/internal/services/aws/common/request"
	neptunestore "vorpalstacks/internal/store/aws/neptune"
)

// CreateDBClusterSnapshot creates a new snapshot of the specified DB cluster.
func (s *NeptuneService) CreateDBClusterSnapshot(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	snapshotID := request.GetStringParam(params, "DBClusterSnapshotIdentifier")
	if snapshotID == "" {
		return nil, fmt.Errorf("neptune: DBClusterSnapshotIdentifier is required")
	}
	clusterID := request.GetStringParam(params, "DBClusterIdentifier")
	if clusterID == "" {
		return nil, fmt.Errorf("neptune: DBClusterIdentifier is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if _, err := store.GetSnapshot(snapshotID); err == nil {
		return nil, fmt.Errorf("neptune: DBClusterSnapshot %s already exists", snapshotID)
	}

	cluster, err := store.GetCluster(clusterID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	snapshot := &neptunestore.DBClusterSnapshot{
		DBClusterSnapshotIdentifier: snapshotID,
		DBClusterIdentifier:         clusterID,
		SnapshotCreateTime:          &now,
		Engine:                      cluster.Engine,
		EngineVersion:               cluster.EngineVersion,
		Status:                      "creating",
		Port:                        cluster.Port,
		StorageEncrypted:            cluster.StorageEncrypted,
		KmsKeyId:                    cluster.KmsKeyId,
		DBSnapshotArn:               fmt.Sprintf("arn:aws:rds:%s:%s:cluster-snapshot:%s", reqCtx.GetRegion(), reqCtx.GetAccountID(), snapshotID),
		AccountID:                   reqCtx.GetAccountID(),
		Region:                      reqCtx.GetRegion(),
	}
	snapshot.Status = "available"

	if err := store.CreateSnapshot(snapshot); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"DBClusterSnapshot": snapshot,
	}, nil
}

// DeleteDBClusterSnapshot deletes the specified DB cluster snapshot.
func (s *NeptuneService) DeleteDBClusterSnapshot(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	snapshotID := request.GetStringParam(params, "DBClusterSnapshotIdentifier")
	if snapshotID == "" {
		return nil, fmt.Errorf("neptune: DBClusterSnapshotIdentifier is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	snapshot, err := store.GetSnapshot(snapshotID)
	if err != nil {
		return nil, err
	}

	_ = store.DeleteSnapshot(snapshotID)

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

	snapshotID := request.GetStringParam(params, "DBClusterSnapshotIdentifier")
	if snapshotID != "" {
		snapshot, err := store.GetSnapshot(snapshotID)
		if err != nil {
			return nil, err
		}
		return map[string]interface{}{
			"DBClusterSnapshots": protocol.XMLElements{ElementName: "DBClusterSnapshot", Items: []interface{}{snapshot}},
		}, nil
	}

	snapshots, err := store.ListSnapshots()
	if err != nil {
		return nil, err
	}

	result := make([]interface{}, 0, len(snapshots))
	for _, snap := range snapshots {
		result = append(result, snap)
	}

	return map[string]interface{}{
		"DBClusterSnapshots": protocol.XMLElements{ElementName: "DBClusterSnapshot", Items: result},
	}, nil
}

// CopyDBClusterSnapshot creates a copy of the specified DB cluster snapshot.
func (s *NeptuneService) CopyDBClusterSnapshot(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	sourceID := request.GetStringParam(params, "SourceDBClusterSnapshotIdentifier")
	if sourceID == "" {
		return nil, fmt.Errorf("neptune: SourceDBClusterSnapshotIdentifier is required")
	}
	targetID := request.GetStringParam(params, "TargetDBClusterSnapshotIdentifier")
	if targetID == "" {
		return nil, fmt.Errorf("neptune: TargetDBClusterSnapshotIdentifier is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	source, err := store.GetSnapshot(sourceID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	copy := *source
	copy.DBClusterSnapshotIdentifier = targetID
	copy.SnapshotCreateTime = &now
	copy.DBSnapshotArn = copy.ARN(reqCtx.GetAccountID(), reqCtx.GetRegion())

	if err := store.CreateSnapshot(&copy); err != nil {
		return nil, err
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
		return nil, fmt.Errorf("neptune: DBClusterSnapshotIdentifier is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	_, err = store.GetSnapshot(snapshotID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"DBClusterSnapshotIdentifier": snapshotID,
		"DBClusterSnapshotAttributes": protocol.XMLElements{ElementName: "DBClusterSnapshotAttribute", Items: []interface{}{
			map[string]interface{}{"AttributeName": "restore", "AttributeValues": protocol.XMLElements{ElementName: "AttributeValue", Items: []interface{}{}}},
		}},
	}, nil
}

// ModifyDBClusterSnapshotAttribute modifies an attribute of the specified DB
// cluster snapshot.
func (s *NeptuneService) ModifyDBClusterSnapshotAttribute(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	snapshotID := request.GetStringParam(params, "DBClusterSnapshotIdentifier")
	if snapshotID == "" {
		return nil, fmt.Errorf("neptune: DBClusterSnapshotIdentifier is required")
	}

	return map[string]interface{}{
		"DBClusterSnapshotIdentifier": snapshotID,
		"DBClusterSnapshotAttributes": protocol.XMLElements{ElementName: "DBClusterSnapshotAttribute", Items: []interface{}{
			map[string]interface{}{"AttributeName": "restore", "AttributeValues": protocol.XMLElements{ElementName: "AttributeValue", Items: []interface{}{}}},
		}},
	}, nil
}

// RestoreDBClusterFromSnapshot creates a new DB cluster from a DB cluster snapshot.
func (s *NeptuneService) RestoreDBClusterFromSnapshot(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	clusterID := request.GetStringParam(params, "DBClusterIdentifier")
	if clusterID == "" {
		return nil, fmt.Errorf("neptune: DBClusterIdentifier is required")
	}
	engine := request.GetStringParam(params, "Engine")
	if engine == "" {
		engine = "neptune"
	}
	snapshotID := request.GetStringParam(params, "SnapshotIdentifier")
	if snapshotID == "" {
		return nil, fmt.Errorf("neptune: SnapshotIdentifier is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	snapshot, err := store.GetSnapshot(snapshotID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	cluster := &neptunestore.DBCluster{
		DBClusterIdentifier:         clusterID,
		Engine:                      engine,
		EngineVersion:               snapshot.EngineVersion,
		Status:                      "creating",
		Port:                        request.GetIntParam(params, "Port"),
		BackupRetentionPeriod:       request.GetIntParam(params, "BackupRetentionPeriod"),
		DBClusterParameterGroupName: request.GetStringParam(params, "DBClusterParameterGroupName"),
		DBSubnetGroupName:           request.GetStringParam(params, "DBSubnetGroupName"),
		StorageEncrypted:            request.GetBoolParam(params, "StorageEncrypted"),
		DeletionProtection:          request.GetBoolParam(params, "DeletionProtection"),
		ClusterCreateTime:           &now,
		EarliestRestorableTime:      &now,
		LatestRestorableTime:        &now,
		AccountID:                   reqCtx.GetAccountID(),
		Region:                      reqCtx.GetRegion(),
	}
	cluster.Status = "available"

	if err := store.CreateCluster(cluster); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"DBCluster": cluster,
	}, nil
}

// RestoreDBClusterToPointInTime restores a DB cluster to a point in time from
// a source cluster.
func (s *NeptuneService) RestoreDBClusterToPointInTime(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	clusterID := request.GetStringParam(params, "DBClusterIdentifier")
	if clusterID == "" {
		return nil, fmt.Errorf("neptune: DBClusterIdentifier is required")
	}
	sourceID := request.GetStringParam(params, "SourceDBClusterIdentifier")
	if sourceID == "" {
		return nil, fmt.Errorf("neptune: SourceDBClusterIdentifier is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	source, err := store.GetCluster(sourceID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	cluster := &neptunestore.DBCluster{
		DBClusterIdentifier:         clusterID,
		Engine:                      source.Engine,
		EngineVersion:               source.EngineVersion,
		Status:                      "creating",
		Port:                        request.GetIntParam(params, "Port"),
		BackupRetentionPeriod:       request.GetIntParam(params, "BackupRetentionPeriod"),
		DBClusterParameterGroupName: request.GetStringParam(params, "DBClusterParameterGroupName"),
		DBSubnetGroupName:           request.GetStringParam(params, "DBSubnetGroupName"),
		StorageEncrypted:            source.StorageEncrypted,
		DeletionProtection:          request.GetBoolParam(params, "DeletionProtection"),
		ClusterCreateTime:           &now,
		EarliestRestorableTime:      &now,
		LatestRestorableTime:        &now,
		AccountID:                   reqCtx.GetAccountID(),
		Region:                      reqCtx.GetRegion(),
	}
	cluster.Status = "available"

	if err := store.CreateCluster(cluster); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"DBCluster": cluster,
	}, nil
}

// PromoteReadReplicaDBCluster promotes a read replica cluster to a standalone
// writable cluster.
func (s *NeptuneService) PromoteReadReplicaDBCluster(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	id := request.GetStringParam(params, "DBClusterIdentifier")
	if id == "" {
		return nil, fmt.Errorf("neptune: DBClusterIdentifier is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	cluster, err := store.GetCluster(id)
	if err != nil {
		return nil, err
	}

	cluster.ReplicationSourceIdentifier = ""
	cluster.Status = "available"
	_ = store.UpdateCluster(cluster)

	return map[string]interface{}{
		"DBCluster": cluster,
	}, nil
}
