package neptune

import (
	"context"
	"fmt"
	"net/http"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/protocol"
	"vorpalstacks/internal/common/request"
	appconfig "vorpalstacks/internal/config"
	neptunestore "vorpalstacks/internal/store/aws/neptune"
	arnutil "vorpalstacks/internal/utils/aws/arn"
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
		Status:                      "available",
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
		return nil, translateStoreError(err)
	}

	if err := store.DeleteSnapshot(snapshotID); err != nil {
		return nil, translateStoreError(err)
	}

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

	snapshotID := request.GetStringParam(params, "DBClusterSnapshotIdentifier")
	if snapshotID != "" {
		snapshot, err := store.GetSnapshot(snapshotID)
		if err != nil {
			return nil, translateStoreError(err)
		}
		return map[string]interface{}{
			"DBClusterSnapshots": protocol.XMLElements{ElementName: "DBClusterSnapshot", Items: []interface{}{snapshot}},
		}, nil
	}

	snapshots, err := store.ListSnapshots()
	if err != nil {
		return nil, translateStoreError(err)
	}

	clusterID := request.GetStringParam(params, "DBClusterIdentifier")
	items := make([]interface{}, 0, len(snapshots))
	for _, snap := range snapshots {
		if clusterID != "" && snap.DBClusterIdentifier != clusterID {
			continue
		}
		items = append(items, snap)
	}

	marker := request.GetStringParam(params, "Marker")
	maxRecords := request.GetIntParam(params, "MaxRecords")
	resultItems, nextMarker, isTruncated := paginateItems(items, marker, maxRecords, func(item interface{}) string {
		return item.(*neptunestore.DBClusterSnapshot).DBClusterSnapshotIdentifier
	})

	result := map[string]interface{}{
		"DBClusterSnapshots": protocol.XMLElements{ElementName: "DBClusterSnapshot", Items: resultItems},
	}
	if isTruncated {
		result["Marker"] = nextMarker
	}
	return result, nil
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
		return nil, fmt.Errorf("neptune: DBClusterSnapshotIdentifier is required")
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
		return nil, fmt.Errorf("neptune: DBClusterSnapshotIdentifier is required")
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
	valuesToAdd := request.GetStringList(params, "ValuesToAdd")
	valuesToRemove := request.GetStringList(params, "ValuesToRemove")

	if attrName == "" {
		attrName = "restore"
	}
	if attrName != "restore" {
		return nil, fmt.Errorf("neptune: unsupported attribute name '%s', only 'restore' is supported", attrName)
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
		return nil, translateStoreError(err)
	}

	now := time.Now()
	port := request.GetIntParam(params, "Port")
	if port == 0 {
		port, _ = appconfig.GetResourcePort("ports.neptune", clusterID)
	}
	backupRetention := request.GetIntParam(params, "BackupRetentionPeriod")
	if backupRetention == 0 {
		backupRetention = 1
	}
	cluster := &neptunestore.DBCluster{
		DBClusterIdentifier:         clusterID,
		Engine:                      engine,
		EngineVersion:               snapshot.EngineVersion,
		Status:                      "available",
		Port:                        port,
		BackupRetentionPeriod:       backupRetention,
		DBClusterParameterGroupName: request.GetStringParam(params, "DBClusterParameterGroupName"),
		DBSubnetGroupName:           request.GetStringParam(params, "DBSubnetGroupName"),
		StorageEncrypted:            request.GetBoolParam(params, "StorageEncrypted"),
		DeletionProtection:          request.GetBoolParam(params, "DeletionProtection"),
		ClusterCreateTime:           &now,
		EarliestRestorableTime:      &now,
		LatestRestorableTime:        &now,
		AccountID:                   reqCtx.GetAccountID(),
		Region:                      reqCtx.GetRegion(),
		DBClusterArn:                arnutil.NewARNBuilder(reqCtx.GetAccountID(), reqCtx.GetRegion()).RDS().Cluster(clusterID),
	}

	if err := store.CreateCluster(cluster); err != nil {
		return nil, translateStoreError(err)
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
		return nil, translateStoreError(err)
	}

	now := time.Now()
	port := request.GetIntParam(params, "Port")
	if port == 0 {
		port, _ = appconfig.GetResourcePort("ports.neptune", clusterID)
	}
	backupRetention := request.GetIntParam(params, "BackupRetentionPeriod")
	if backupRetention == 0 {
		backupRetention = 1
	}
	cluster := &neptunestore.DBCluster{
		DBClusterIdentifier:         clusterID,
		Engine:                      source.Engine,
		EngineVersion:               source.EngineVersion,
		Status:                      "available",
		Port:                        port,
		BackupRetentionPeriod:       backupRetention,
		DBClusterParameterGroupName: request.GetStringParam(params, "DBClusterParameterGroupName"),
		DBSubnetGroupName:           request.GetStringParam(params, "DBSubnetGroupName"),
		StorageEncrypted:            source.StorageEncrypted,
		DeletionProtection:          request.GetBoolParam(params, "DeletionProtection"),
		ClusterCreateTime:           &now,
		EarliestRestorableTime:      &now,
		LatestRestorableTime:        &now,
		AccountID:                   reqCtx.GetAccountID(),
		Region:                      reqCtx.GetRegion(),
		DBClusterArn:                arnutil.NewARNBuilder(reqCtx.GetAccountID(), reqCtx.GetRegion()).RDS().Cluster(clusterID),
	}

	if err := store.CreateCluster(cluster); err != nil {
		return nil, translateStoreError(err)
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
		return nil, translateStoreError(err)
	}

	if cluster.ReplicationSourceIdentifier == "" {
		return nil, awserrors.NewAWSError("InvalidDBClusterStateFault", fmt.Sprintf("cluster %s is not a read replica", id), http.StatusBadRequest)
	}

	if cluster.Status != "available" {
		return nil, awserrors.NewAWSError("InvalidDBClusterStateFault", fmt.Sprintf("cluster %s is not in available state", id), http.StatusBadRequest)
	}

	cluster.ReplicationSourceIdentifier = ""
	cluster.GlobalClusterIdentifier = ""
	cluster.Status = "promoting"
	if err := store.UpdateCluster(cluster); err != nil {
		return nil, translateStoreError(err)
	}

	cluster.Status = "available"
	if err := store.UpdateCluster(cluster); err != nil {
		return nil, translateStoreError(err)
	}

	return map[string]interface{}{
		"DBCluster": cluster,
	}, nil
}
