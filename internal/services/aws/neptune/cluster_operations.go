package neptune

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/protocol"
	"vorpalstacks/internal/common/request"
	appconfig "vorpalstacks/internal/config"
	"vorpalstacks/internal/core/logs"
	neptunestore "vorpalstacks/internal/store/aws/neptune"
	arnutil "vorpalstacks/internal/utils/aws/arn"
	"vorpalstacks/internal/utils/aws/types"
)

func clusterToResponseMap(cluster *neptunestore.DBCluster) map[string]interface{} {
	data, err := json.Marshal(cluster)
	if err != nil {
		return map[string]interface{}{"error": err.Error()}
	}
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return map[string]interface{}{"error": err.Error()}
	}
	for k, v := range m {
		if v == nil {
			delete(m, k)
		}
	}
	return m
}

func enrichClusterWithTags(store neptunestore.NeptuneStoreInterface, cluster *neptunestore.DBCluster) map[string]interface{} {
	m := clusterToResponseMap(cluster)
	tags, _ := store.GetTags(cluster.DBClusterArn)
	if len(tags) > 0 {
		tagItems := make([]interface{}, 0, len(tags))
		for _, t := range tags {
			tagItems = append(tagItems, map[string]interface{}{"Key": t.Key, "Value": t.Value})
		}
		m["TagList"] = protocol.XMLElements{ElementName: "Tag", Items: tagItems}
	}
	return m
}

// CreateDBCluster creates a new Neptune DB cluster with the specified configuration.
func (s *NeptuneService) CreateDBCluster(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	id := request.GetStringParam(params, "DBClusterIdentifier")
	if id == "" {
		return nil, fmt.Errorf("neptune: DBClusterIdentifier is required")
	}

	engine := request.GetStringParam(params, "Engine")
	if engine == "" {
		engine = "neptune"
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	replicationSource := request.GetStringParam(params, "ReplicationSourceIdentifier")

	if replicationSource != "" {
		_, err := store.GetCluster(replicationSource)
		if err != nil {
			return nil, awserrors.NewAWSError("InvalidParameterValue", fmt.Sprintf("replication source cluster %s not found", replicationSource), http.StatusBadRequest)
		}
	}

	now := time.Now()
	port := request.GetIntParam(params, "Port")
	if port == 0 {
		port, _ = appconfig.GetResourcePort("ports.neptune", id)
	}
	backupRetention := request.GetIntParam(params, "BackupRetentionPeriod")
	if backupRetention == 0 {
		backupRetention = 1
	}
	cluster := &neptunestore.DBCluster{
		DBClusterIdentifier:             id,
		Engine:                          engine,
		EngineVersion:                   request.GetStringParam(params, "EngineVersion"),
		Status:                          "available",
		Port:                            port,
		BackupRetentionPeriod:           backupRetention,
		PreferredBackupWindow:           request.GetStringParam(params, "PreferredBackupWindow"),
		PreferredMaintenanceWindow:      request.GetStringParam(params, "PreferredMaintenanceWindow"),
		MasterUsername:                  request.GetStringParam(params, "MasterUsername"),
		DatabaseName:                    request.GetStringParam(params, "DatabaseName"),
		DBClusterParameterGroupName:     request.GetStringParam(params, "DBClusterParameterGroupName"),
		DBSubnetGroupName:               request.GetStringParam(params, "DBSubnetGroupName"),
		StorageEncrypted:                request.GetBoolParam(params, "StorageEncrypted"),
		KmsKeyId:                        request.GetStringParam(params, "KmsKeyId"),
		CopyTagsToSnapshot:              request.GetBoolParam(params, "CopyTagsToSnapshot"),
		DeletionProtection:              request.GetBoolParam(params, "DeletionProtection"),
		EnableIAMDatabaseAuthentication: request.GetBoolParam(params, "EnableIAMDatabaseAuthentication"),
		ClusterCreateTime:               &now,
		EarliestRestorableTime:          &now,
		LatestRestorableTime:            &now,
		ReplicationSourceIdentifier:     replicationSource,
		GlobalClusterIdentifier:         request.GetStringParam(params, "GlobalClusterIdentifier"),
		StorageType:                     request.GetStringParam(params, "StorageType"),
		AccountID:                       reqCtx.GetAccountID(),
		Region:                          reqCtx.GetRegion(),
		DBClusterArn:                    arnutil.NewARNBuilder(reqCtx.GetAccountID(), reqCtx.GetRegion()).RDS().Cluster(id),
	}

	if azList := request.GetStringList(params, "AvailabilityZones"); len(azList) > 0 {
		cluster.AvailabilityZones = azList
	}
	if sgList := request.GetStringList(params, "VpcSecurityGroupIds"); len(sgList) > 0 {
		if _, err := s.resolveSecurityGroups(ctx, reqCtx.GetRegion(), sgList); err != nil {
			return nil, translateStoreError(err)
		}
		cluster.VpcSecurityGroupIds = sgList
	}
	if logExports := request.GetStringList(params, "EnableCloudwatchLogsExports"); len(logExports) > 0 {
		cluster.EnableCloudwatchLogsExports = logExports
	}

	if err := store.CreateCluster(cluster); err != nil {
		return nil, translateStoreError(err)
	}

	if tagList := getNeptuneTagList(params); len(tagList) > 0 {
		storeTags := make([]types.Tag, 0, len(tagList))
		for _, t := range tagList {
			key, _ := t["Key"].(string)
			value, _ := t["Value"].(string)
			if key != "" {
				storeTags = append(storeTags, types.Tag{Key: key, Value: value})
			}
		}
		if err := store.AddTags(cluster.DBClusterArn, storeTags); err != nil {
			logs.Warn("failed to tag cluster on create", logs.String("cluster", id), logs.Err(err))
		}
	}

	recordEvent(store, "db-cluster", id, cluster.DBClusterArn,
		fmt.Sprintf("DB cluster %s created", id), []string{"creation"})

	return map[string]interface{}{
		"DBCluster": enrichClusterWithTags(store, cluster),
	}, nil
}

// DeleteDBCluster deletes the specified Neptune DB cluster.
func (s *NeptuneService) DeleteDBCluster(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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

	if cluster.DeletionProtection {
		return nil, awserrors.NewAWSError("InvalidDBClusterStateFault", "Cannot delete cluster when DeletionProtection is enabled", http.StatusBadRequest)
	}

	skipFinal := request.GetBoolParam(params, "SkipFinalSnapshot")
	finalSnapshotID := request.GetStringParam(params, "FinalDBSnapshotIdentifier")
	if !skipFinal && finalSnapshotID == "" {
		return nil, awserrors.NewAWSError("InvalidParameterCombination", "SkipFinalSnapshot must be true or FinalDBSnapshotIdentifier must be specified", http.StatusBadRequest)
	}
	if skipFinal && finalSnapshotID != "" {
		return nil, awserrors.NewAWSError("InvalidParameterCombination", "Cannot specify both SkipFinalSnapshot and FinalDBSnapshotIdentifier", http.StatusBadRequest)
	}

	cluster.Status = "deleting"

	if !skipFinal {
		now := time.Now()
		snapshot := &neptunestore.DBClusterSnapshot{
			DBClusterSnapshotIdentifier: finalSnapshotID,
			DBClusterIdentifier:         id,
			SnapshotCreateTime:          &now,
			Engine:                      cluster.Engine,
			EngineVersion:               cluster.EngineVersion,
			Status:                      "available",
			Port:                        cluster.Port,
			StorageEncrypted:            cluster.StorageEncrypted,
			KmsKeyId:                    cluster.KmsKeyId,
			DBSnapshotArn:               arnutil.NewARNBuilder(reqCtx.GetAccountID(), reqCtx.GetRegion()).RDS().ClusterSnapshot(finalSnapshotID),
			AccountID:                   reqCtx.GetAccountID(),
			Region:                      reqCtx.GetRegion(),
		}
		if err := store.CreateSnapshot(snapshot); err != nil {
			return nil, translateStoreError(err)
		}
	}

	cascadeDeleteClusterResources(store, cluster)

	if err := store.DeleteCluster(id); err != nil {
		return nil, translateStoreError(err)
	}

	recordEvent(store, "db-cluster", id, cluster.DBClusterArn,
		fmt.Sprintf("DB cluster %s deleted", id), []string{"deletion"})

	return map[string]interface{}{
		"DBCluster": cluster,
	}, nil
}

// ModifyDBCluster updates the configuration of the specified DB cluster.
func (s *NeptuneService) ModifyDBCluster(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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

	if v := request.GetStringParam(params, "EngineVersion"); v != "" {
		cluster.EngineVersion = v
	}
	if v := request.GetStringParam(params, "DBClusterParameterGroupName"); v != "" {
		cluster.DBClusterParameterGroupName = v
	}
	if v := request.GetIntParam(params, "Port"); v > 0 {
		cluster.Port = v
	}
	if v := request.GetIntParam(params, "BackupRetentionPeriod"); v > 0 {
		cluster.BackupRetentionPeriod = v
	}
	if v := request.GetStringParam(params, "PreferredBackupWindow"); v != "" {
		cluster.PreferredBackupWindow = v
	}
	if v := request.GetStringParam(params, "PreferredMaintenanceWindow"); v != "" {
		cluster.PreferredMaintenanceWindow = v
	}
	if v := request.GetStringParam(params, "StorageType"); v != "" {
		cluster.StorageType = v
	}
	if request.HasParam(params, "DeletionProtection") {
		cluster.DeletionProtection = request.GetBoolParam(params, "DeletionProtection")
	}
	if request.HasParam(params, "EnableIAMDatabaseAuthentication") {
		cluster.EnableIAMDatabaseAuthentication = request.GetBoolParam(params, "EnableIAMDatabaseAuthentication")
	}
	if request.HasParam(params, "VpcSecurityGroupIds") {
		sgList := request.GetStringList(params, "VpcSecurityGroupIds")
		if len(sgList) > 0 {
			if _, err := s.resolveSecurityGroups(ctx, reqCtx.GetRegion(), sgList); err != nil {
				return nil, translateStoreError(err)
			}
		}
		cluster.VpcSecurityGroupIds = sgList
	}
	if request.HasParam(params, "EnableCloudwatchLogsExports") {
		cluster.EnableCloudwatchLogsExports = request.GetStringList(params, "EnableCloudwatchLogsExports")
	}

	if err := store.UpdateCluster(cluster); err != nil {
		return nil, translateStoreError(err)
	}

	newID := request.GetStringParam(params, "NewDBClusterIdentifier")
	if newID != "" && newID != id {
		oldID := cluster.DBClusterIdentifier
		cluster.DBClusterIdentifier = newID
		cluster.DBClusterArn = arnutil.NewARNBuilder(cluster.AccountID, cluster.Region).RDS().Cluster(newID)
		if err := store.CreateCluster(cluster); err != nil {
			return nil, translateStoreError(err)
		}
		if err := store.DeleteCluster(oldID); err != nil {
			_ = store.DeleteCluster(newID)
			return nil, translateStoreError(err)
		}
	}

	return map[string]interface{}{
		"DBCluster": cluster,
	}, nil
}

// DescribeDBClusters returns information about the specified DB cluster or lists
// all clusters when no identifier is provided.
func (s *NeptuneService) DescribeDBClusters(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	clusterID := request.GetStringParam(params, "DBClusterIdentifier")
	if clusterID != "" {
		cluster, err := store.GetCluster(clusterID)
		if err != nil {
			return nil, translateStoreError(err)
		}
		return map[string]interface{}{
			"DBClusters": protocol.XMLElements{ElementName: "DBCluster", Items: []interface{}{enrichClusterWithTags(store, cluster)}},
		}, nil
	}

	clusters, err := store.ListClusters()
	if err != nil {
		return nil, translateStoreError(err)
	}

	items := make([]interface{}, 0, len(clusters))
	for _, c := range clusters {
		items = append(items, enrichClusterWithTags(store, c))
	}

	marker := request.GetStringParam(params, "Marker")
	maxRecords := request.GetIntParam(params, "MaxRecords")
	resultItems, nextMarker, isTruncated := paginateItems(items, marker, maxRecords, func(item interface{}) string {
		if m, ok := item.(map[string]interface{}); ok {
			if v, ok := m["DBClusterIdentifier"].(string); ok {
				return v
			}
		}
		return ""
	})

	result := map[string]interface{}{
		"DBClusters": protocol.XMLElements{ElementName: "DBCluster", Items: resultItems},
	}
	if isTruncated {
		result["Marker"] = nextMarker
	}
	return result, nil
}

// StartDBCluster restarts a stopped Neptune DB cluster.
func (s *NeptuneService) StartDBCluster(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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

	if cluster.Status != "stopped" {
		return nil, fmt.Errorf("neptune: DBCluster %s is not in stopped state (current: %s)", id, cluster.Status)
	}

	cluster.Status = "available"
	if err := store.UpdateCluster(cluster); err != nil {
		return nil, translateStoreError(err)
	}

	return map[string]interface{}{
		"DBCluster": cluster,
	}, nil
}

// StopDBCluster stops a running Neptune DB cluster.
func (s *NeptuneService) StopDBCluster(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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

	if cluster.Status != "available" {
		return nil, fmt.Errorf("neptune: DBCluster %s is not in available state (current: %s)", id, cluster.Status)
	}

	cluster.Status = "stopped"
	if err := store.UpdateCluster(cluster); err != nil {
		return nil, translateStoreError(err)
	}

	return map[string]interface{}{
		"DBCluster": cluster,
	}, nil
}

// FailoverDBCluster forces a failover for the specified DB cluster.
func (s *NeptuneService) FailoverDBCluster(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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

	cluster.Status = "available"
	if err := store.UpdateCluster(cluster); err != nil {
		return nil, translateStoreError(err)
	}

	return map[string]interface{}{
		"DBCluster": cluster,
	}, nil
}

// AddRoleToDBCluster associates an IAM role with the specified DB cluster.
func (s *NeptuneService) AddRoleToDBCluster(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	id := request.GetStringParam(params, "DBClusterIdentifier")
	if id == "" {
		return nil, fmt.Errorf("neptune: DBClusterIdentifier is required")
	}
	roleArn := request.GetStringParam(params, "RoleArn")
	if roleArn == "" {
		return nil, fmt.Errorf("neptune: RoleArn is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	cluster, err := store.GetCluster(id)
	if err != nil {
		return nil, translateStoreError(err)
	}

	featureName := request.GetStringParam(params, "FeatureName")
	for _, r := range cluster.AssociatedRoles {
		if r.RoleArn == roleArn {
			return nil, fmt.Errorf("neptune: IAM role %s is already associated with cluster %s", roleArn, id)
		}
	}
	cluster.AssociatedRoles = append(cluster.AssociatedRoles, neptunestore.DBClusterRole{
		RoleArn:     roleArn,
		FeatureName: featureName,
		Status:      "ACTIVE",
	})
	if err := store.UpdateCluster(cluster); err != nil {
		return nil, translateStoreError(err)
	}

	return map[string]interface{}{}, nil
}

// RemoveRoleFromDBCluster disassociates an IAM role from the specified DB cluster.
func (s *NeptuneService) RemoveRoleFromDBCluster(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	id := request.GetStringParam(params, "DBClusterIdentifier")
	if id == "" {
		return nil, fmt.Errorf("neptune: DBClusterIdentifier is required")
	}
	roleArn := request.GetStringParam(params, "RoleArn")
	if roleArn == "" {
		return nil, fmt.Errorf("neptune: RoleArn is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	cluster, err := store.GetCluster(id)
	if err != nil {
		return nil, translateStoreError(err)
	}

	found := false
	filtered := make([]neptunestore.DBClusterRole, 0, len(cluster.AssociatedRoles))
	for _, r := range cluster.AssociatedRoles {
		if r.RoleArn == roleArn {
			found = true
			continue
		}
		filtered = append(filtered, r)
	}
	if !found {
		return nil, awserrors.NewAWSError("InvalidParameterValue", fmt.Sprintf("role %s is not associated with cluster %s", roleArn, id), http.StatusBadRequest)
	}
	cluster.AssociatedRoles = filtered
	if err := store.UpdateCluster(cluster); err != nil {
		return nil, translateStoreError(err)
	}

	return map[string]interface{}{}, nil
}

// cascadeDeleteClusterResources removes all instances, cluster endpoints, and tags
// associated with the given cluster. Errors are logged but not returned so that
// the cluster deletion itself always succeeds.
func cascadeDeleteClusterResources(store neptunestore.NeptuneStoreInterface, cluster *neptunestore.DBCluster) {
	clusterID := cluster.DBClusterIdentifier

	instances, err := store.ListInstances()
	if err != nil {
		logs.Warn("cascade: failed to list instances", logs.Err(err))
		return
	}
	for _, inst := range instances {
		if inst.DBClusterIdentifier == clusterID {
			if delErr := store.DeleteInstance(inst.DBInstanceIdentifier); delErr != nil {
				logs.Warn("cascade: failed to delete instance", logs.String("instance", inst.DBInstanceIdentifier), logs.Err(delErr))
			} else {
				tags, _ := store.GetTags(inst.DBInstanceArn)
				if len(tags) > 0 {
					keys := make([]string, len(tags))
					for i, t := range tags {
						keys[i] = t.Key
					}
					if tagErr := store.RemoveTags(inst.DBInstanceArn, keys); tagErr != nil {
						logs.Warn("cascade: failed to remove instance tags", logs.String("instance", inst.DBInstanceIdentifier), logs.Err(tagErr))
					}
				}
			}
		}
	}

	endpoints, err := store.ListClusterEndpoints(clusterID)
	if err != nil {
		logs.Warn("cascade: failed to list cluster endpoints", logs.Err(err))
		return
	}
	for _, ep := range endpoints {
		if delErr := store.DeleteClusterEndpoint(ep.DBClusterEndpointIdentifier); delErr != nil {
			logs.Warn("cascade: failed to delete cluster endpoint", logs.String("endpoint", ep.DBClusterEndpointIdentifier), logs.Err(delErr))
		}
	}

	tags, _ := store.GetTags(cluster.DBClusterArn)
	if len(tags) > 0 {
		keys := make([]string, len(tags))
		for i, t := range tags {
			keys[i] = t.Key
		}
		if tagErr := store.RemoveTags(cluster.DBClusterArn, keys); tagErr != nil {
			logs.Warn("cascade: failed to remove cluster tags", logs.String("cluster", clusterID), logs.Err(tagErr))
		}
	}
}
