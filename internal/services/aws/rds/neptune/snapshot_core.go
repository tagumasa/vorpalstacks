package neptune

import (
	"context"
	"fmt"
	"net/http"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/protocol"
	"vorpalstacks/internal/core/logs"
	rdssvc "vorpalstacks/internal/services/aws/rds"
	neptunestore "vorpalstacks/internal/store/aws/rds/neptune"
	arnutil "vorpalstacks/internal/utils/aws/arn"
)

// CreateDBClusterSnapshotInput carries the wire-parsed
// CreateDBClusterSnapshot request.
type CreateDBClusterSnapshotInput struct {
	DBClusterSnapshotIdentifier string
	DBClusterIdentifier         string
	AccountID                   string
	Region                      string
}

// DeleteDBClusterSnapshotInput carries the wire-parsed
// DeleteDBClusterSnapshot request.
type DeleteDBClusterSnapshotInput struct {
	DBClusterSnapshotIdentifier string
}

// CopyDBClusterSnapshotInput carries the wire-parsed CopyDBClusterSnapshot
// request.
type CopyDBClusterSnapshotInput struct {
	SourceDBClusterSnapshotIdentifier string
	TargetDBClusterSnapshotIdentifier string
	AccountID                         string
	Region                            string
}

// DescribeDBClusterSnapshotAttributesInput carries the wire-parsed
// DescribeDBClusterSnapshotAttributes request.
type DescribeDBClusterSnapshotAttributesInput struct {
	DBClusterSnapshotIdentifier string
}

// ModifyDBClusterSnapshotAttributeInput carries the wire-parsed
// ModifyDBClusterSnapshotAttribute request.
type ModifyDBClusterSnapshotAttributeInput struct {
	DBClusterSnapshotIdentifier string
	AttributeName               string
	ValuesToAdd                 []string
	ValuesToRemove              []string
}

// RestoreDBClusterInput carries the wire-parsed restore request members
// shared by RestoreDBClusterFromSnapshot and RestoreDBClusterToPointInTime.
// Only the members the two operations actually consume are carried;
// KmsKeyId, PreferredBackupWindow and PreferredMaintenanceWindow reach the
// restored cluster only through the snapshot/source inheritance paths, so
// the wire values are carried solely to gate those inheritances.
type RestoreDBClusterInput struct {
	DBClusterIdentifier         string
	SnapshotIdentifier          string
	SourceDBClusterIdentifier   string
	Engine                      string
	EngineVersion               string
	BackupRetentionPeriod       int
	Port                        int
	DBClusterParameterGroupName string
	DBSubnetGroupName           string
	StorageEncrypted            bool
	DeletionProtection          bool
	PreferredBackupWindow       string
	PreferredMaintenanceWindow  string
	AccountID                   string
	Region                      string
}

// PromoteReadReplicaDBClusterInput carries the wire-parsed
// PromoteReadReplicaDBCluster request.
type PromoteReadReplicaDBClusterInput struct {
	DBClusterIdentifier string
}

// buildRestoredCluster assembles the cluster record for a restore operation.
// Cluster snapshots are metadata-only — they do not capture row-level data.
func buildRestoredCluster(in *RestoreDBClusterInput, now *time.Time) *neptunestore.DBCluster {
	backupRetention := in.BackupRetentionPeriod
	if backupRetention == 0 {
		backupRetention = 1
	}
	return &neptunestore.DBCluster{
		DBClusterIdentifier:         in.DBClusterIdentifier,
		Engine:                      in.Engine,
		EngineVersion:               in.EngineVersion,
		Status:                      "available",
		Port:                        in.Port,
		BackupRetentionPeriod:       backupRetention,
		DBClusterParameterGroupName: in.DBClusterParameterGroupName,
		DBSubnetGroupName:           in.DBSubnetGroupName,
		StorageEncrypted:            in.StorageEncrypted,
		DeletionProtection:          in.DeletionProtection,
		ClusterCreateTime:           now,
		EarliestRestorableTime:      now,
		LatestRestorableTime:        now,
		AccountID:                   in.AccountID,
		Region:                      in.Region,
		DBClusterArn:                arnutil.NewARNBuilder(in.AccountID, in.Region).RDS().Cluster(in.DBClusterIdentifier),
	}
}

// createDBClusterSnapshotCore validates and persists a new cluster snapshot,
// records the creation event and transitions it to available with a
// safety-net goroutine.
func (s *NeptuneService) createDBClusterSnapshotCore(ctx context.Context, store neptunestore.NeptuneStoreInterface, in *CreateDBClusterSnapshotInput) (interface{}, error) {
	snapshotID := in.DBClusterSnapshotIdentifier
	if snapshotID == "" {
		return nil, awserrors.NewMissingParameter("DBClusterSnapshotIdentifier is required")
	}
	if err := rdssvc.ValidateDBSnapshotIdentifier(snapshotID); err != nil {
		return nil, awserrors.NewAWSError("InvalidParameterValue", err.Error(), http.StatusBadRequest)
	}
	clusterID := in.DBClusterIdentifier
	if clusterID == "" {
		return nil, awserrors.NewMissingParameter("DBClusterIdentifier is required")
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
		DBSnapshotArn:               arnutil.NewARNBuilder(in.AccountID, in.Region).RDS().ClusterSnapshot(snapshotID),
		AccountID:                   in.AccountID,
		Region:                      in.Region,
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
	s.scheduleTransition(in.Region, 500*time.Millisecond, func(st neptunestore.NeptuneStoreInterface) error {
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

// deleteDBClusterSnapshotCore deletes a cluster snapshot, clears its tags and
// records the deletion event.
func (s *NeptuneService) deleteDBClusterSnapshotCore(ctx context.Context, store neptunestore.NeptuneStoreInterface, in *DeleteDBClusterSnapshotInput) (interface{}, error) {
	snapshotID := in.DBClusterSnapshotIdentifier
	if snapshotID == "" {
		return nil, awserrors.NewMissingParameter("DBClusterSnapshotIdentifier is required")
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

// copyDBClusterSnapshotCore copies an existing cluster snapshot under a new
// identifier.
func (s *NeptuneService) copyDBClusterSnapshotCore(ctx context.Context, store neptunestore.NeptuneStoreInterface, in *CopyDBClusterSnapshotInput) (interface{}, error) {
	sourceID := in.SourceDBClusterSnapshotIdentifier
	if sourceID == "" {
		return nil, awserrors.NewMissingParameter("SourceDBClusterSnapshotIdentifier is required")
	}
	targetID := in.TargetDBClusterSnapshotIdentifier
	if targetID == "" {
		return nil, awserrors.NewMissingParameter("TargetDBClusterSnapshotIdentifier is required")
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
	copy.DBSnapshotArn = arnutil.NewARNBuilder(in.AccountID, in.Region).RDS().ClusterSnapshot(targetID)

	if err := store.CreateSnapshot(&copy); err != nil {
		return nil, translateStoreError(err)
	}

	return map[string]interface{}{
		"DBClusterSnapshot": &copy,
	}, nil
}

// getDBClusterSnapshotCore resolves a cluster snapshot by name.
func (s *NeptuneService) getDBClusterSnapshotCore(ctx context.Context, store neptunestore.NeptuneStoreInterface, in *DescribeDBClusterSnapshotAttributesInput) (*neptunestore.DBClusterSnapshot, error) {
	snapshotID := in.DBClusterSnapshotIdentifier
	if snapshotID == "" {
		return nil, awserrors.NewMissingParameter("DBClusterSnapshotIdentifier is required")
	}
	return store.GetSnapshot(snapshotID)
}

// modifyDBClusterSnapshotAttributeCore applies add/remove value changes to
// the snapshot's restore attribute and persists the result.
func (s *NeptuneService) modifyDBClusterSnapshotAttributeCore(ctx context.Context, store neptunestore.NeptuneStoreInterface, in *ModifyDBClusterSnapshotAttributeInput) (interface{}, error) {
	snapshotID := in.DBClusterSnapshotIdentifier
	if snapshotID == "" {
		return nil, awserrors.NewMissingParameter("DBClusterSnapshotIdentifier is required")
	}

	snap, err := store.GetSnapshot(snapshotID)
	if err != nil {
		return nil, translateStoreError(err)
	}

	attrName := in.AttributeName
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
	for _, v := range in.ValuesToRemove {
		delete(existing, v)
	}
	for _, v := range in.ValuesToAdd {
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

// restoreDBClusterFromSnapshotCore validates and persists a cluster restored
// from a cluster snapshot, then opens the engine and sets the endpoint.
func (s *NeptuneService) restoreDBClusterFromSnapshotCore(ctx context.Context, store neptunestore.NeptuneStoreInterface, in *RestoreDBClusterInput) (interface{}, error) {
	clusterID := in.DBClusterIdentifier
	if clusterID == "" {
		return nil, awserrors.NewMissingParameter("DBClusterIdentifier is required")
	}
	if err := rdssvc.ValidateDBClusterIdentifier(clusterID); err != nil {
		return nil, awserrors.NewAWSError("InvalidParameterValue", err.Error(), http.StatusBadRequest)
	}
	engine := in.Engine
	if engine == "" {
		engine = "neptune"
	}
	if err := rdssvc.ValidateEngine(engine); err != nil {
		return nil, awserrors.NewAWSError("InvalidParameterValue", err.Error(), http.StatusBadRequest)
	}
	snapshotID := in.SnapshotIdentifier
	if snapshotID == "" {
		return nil, awserrors.NewMissingParameter("SnapshotIdentifier is required")
	}

	snapshot, err := store.GetSnapshot(snapshotID)
	if err != nil {
		return nil, translateStoreError(err)
	}

	engineVersion := in.EngineVersion
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
	cluster := buildRestoredCluster(in, &now)
	cluster.Engine = engine
	cluster.EngineVersion = engineVersion
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
		if port, err := eng.Open(in.Region, clusterID); err != nil {
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

// restoreDBClusterToPointInTimeCore validates and persists a point-in-time
// restore of a source cluster, inheriting omitted members from the source,
// then opens the engine and sets the endpoint.
func (s *NeptuneService) restoreDBClusterToPointInTimeCore(ctx context.Context, store neptunestore.NeptuneStoreInterface, in *RestoreDBClusterInput) (interface{}, error) {
	clusterID := in.DBClusterIdentifier
	if clusterID == "" {
		return nil, awserrors.NewMissingParameter("DBClusterIdentifier is required")
	}
	sourceID := in.SourceDBClusterIdentifier
	if sourceID == "" {
		return nil, awserrors.NewMissingParameter("SourceDBClusterIdentifier is required")
	}

	source, err := store.GetCluster(sourceID)
	if err != nil {
		return nil, translateStoreError(err)
	}

	now := time.Now()
	cluster := buildRestoredCluster(in, &now)
	cluster.Engine = source.Engine
	cluster.EngineVersion = source.EngineVersion
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
	if source.BackupRetentionPeriod > 0 && in.BackupRetentionPeriod == 0 {
		cluster.BackupRetentionPeriod = source.BackupRetentionPeriod
	}
	if source.PreferredBackupWindow != "" && in.PreferredBackupWindow == "" {
		cluster.PreferredBackupWindow = source.PreferredBackupWindow
	}
	if source.PreferredMaintenanceWindow != "" && in.PreferredMaintenanceWindow == "" {
		cluster.PreferredMaintenanceWindow = source.PreferredMaintenanceWindow
	}

	if err := store.CreateCluster(cluster); err != nil {
		return nil, translateStoreError(err)
	}

	var enginePort int
	if eng := s.engineFor(cluster.Engine); eng != nil {
		if port, err := eng.Open(in.Region, clusterID); err != nil {
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

// promoteReadReplicaDBClusterCore promotes a read replica cluster to a
// standalone writable cluster, detaching it from any global cluster.
func (s *NeptuneService) promoteReadReplicaDBClusterCore(ctx context.Context, store neptunestore.NeptuneStoreInterface, in *PromoteReadReplicaDBClusterInput) (interface{}, error) {
	id := in.DBClusterIdentifier
	if id == "" {
		return nil, awserrors.NewMissingParameter("DBClusterIdentifier is required")
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
