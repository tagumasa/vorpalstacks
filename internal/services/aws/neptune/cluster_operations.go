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
	neptunestore "vorpalstacks/internal/store/aws/neptune"
	arnutil "vorpalstacks/internal/utils/aws/arn"
	"vorpalstacks/internal/utils/aws/types"
)

func clusterToResponseMap(c *neptunestore.DBCluster) map[string]interface{} {
	m := map[string]interface{}{
		"DBClusterIdentifier":              c.DBClusterIdentifier,
		"Engine":                           c.Engine,
		"Status":                           c.Status,
		"Port":                             c.Port,
		"BackupRetentionPeriod":            c.BackupRetentionPeriod,
		"MultiAZ":                          c.MultiAZ,
		"StorageEncrypted":                 c.StorageEncrypted,
		"CopyTagsToSnapshot":               c.CopyTagsToSnapshot,
		"DeletionProtection":               c.DeletionProtection,
		"IAMDatabaseAuthenticationEnabled": c.IAMDatabaseAuthenticationEnabled,
		"DBClusterArn":                     c.DBClusterArn,
	}
	if c.EngineVersion != "" {
		m["EngineVersion"] = c.EngineVersion
	}
	if c.MasterUsername != "" {
		m["MasterUsername"] = c.MasterUsername
	}
	if c.DatabaseName != "" {
		m["DatabaseName"] = c.DatabaseName
	}
	if c.PreferredBackupWindow != "" {
		m["PreferredBackupWindow"] = c.PreferredBackupWindow
	}
	if c.PreferredMaintenanceWindow != "" {
		m["PreferredMaintenanceWindow"] = c.PreferredMaintenanceWindow
	}
	if len(c.AvailabilityZones) > 0 {
		m["AvailabilityZones"] = protocol.XMLElements{ElementName: "AvailabilityZone", Items: stringSliceToInterface(c.AvailabilityZones)}
	}
	if len(c.VpcSecurityGroupIds) > 0 {
		m["VpcSecurityGroupIds"] = protocol.XMLElements{ElementName: "VpcSecurityGroupMembership", Items: vpcSecurityToInterface(c.VpcSecurityGroupIds)}
	}
	if c.DBSubnetGroupName != "" {
		m["DBSubnetGroup"] = c.DBSubnetGroupName
	}
	if c.DBClusterParameterGroupName != "" {
		m["DBClusterParameterGroup"] = c.DBClusterParameterGroupName
	}
	if c.KmsKeyId != "" {
		m["KmsKeyId"] = c.KmsKeyId
	}
	if len(c.EnabledCloudwatchLogsExports) > 0 {
		m["EnabledCloudwatchLogsExports"] = protocol.XMLElements{ElementName: "member", Items: stringSliceToInterface(c.EnabledCloudwatchLogsExports)}
	}
	if c.ClusterCreateTime != nil {
		m["ClusterCreateTime"] = c.ClusterCreateTime.Format(time.RFC3339)
	}
	if c.EarliestRestorableTime != nil {
		m["EarliestRestorableTime"] = c.EarliestRestorableTime.Format(time.RFC3339)
	}
	if c.LatestRestorableTime != nil {
		m["LatestRestorableTime"] = c.LatestRestorableTime.Format(time.RFC3339)
	}
	if len(c.AssociatedRoles) > 0 {
		roles := make([]interface{}, 0, len(c.AssociatedRoles))
		for _, r := range c.AssociatedRoles {
			roles = append(roles, map[string]interface{}{"RoleArn": r.RoleArn, "FeatureName": r.FeatureName, "Status": r.Status})
		}
		m["AssociatedRoles"] = roles
	}
	if c.ReplicationSourceIdentifier != "" {
		m["ReplicationSourceIdentifier"] = c.ReplicationSourceIdentifier
	}
	if c.GlobalClusterIdentifier != "" {
		m["GlobalClusterIdentifier"] = c.GlobalClusterIdentifier
	}
	if c.StorageType != "" {
		m["StorageType"] = c.StorageType
	}
	if c.ServerlessV2ScalingConfiguration != nil {
		m["ServerlessV2ScalingConfiguration"] = map[string]interface{}{
			"MinCapacity": c.ServerlessV2ScalingConfiguration.MinCapacity,
			"MaxCapacity": c.ServerlessV2ScalingConfiguration.MaxCapacity,
		}
	}
	if c.Endpoint != nil {
		m["Endpoint"] = map[string]interface{}{
			"Address": c.Endpoint.Address,
			"Port":    c.Endpoint.Port,
		}
	}
	return m
}

func vpcSecurityToInterface(ids []string) []interface{} {
	result := make([]interface{}, len(ids))
	for i, id := range ids {
		result[i] = map[string]interface{}{"VpcSecurityGroupId": id, "Status": "active"}
	}
	return result
}

func enrichClusterWithTags(store neptunestore.NeptuneStoreInterface, cluster *neptunestore.DBCluster) map[string]interface{} {
	m := clusterToResponseMap(cluster)
	tags, err := store.GetTags(cluster.DBClusterArn)
	if err != nil {
		logs.Warn("failed to get tags for cluster", logs.String("cluster", cluster.DBClusterIdentifier), logs.Err(err))
	}
	if len(tags) > 0 {
		tagItems := make([]interface{}, 0, len(tags))
		for _, t := range tags {
			tagItems = append(tagItems, map[string]interface{}{"Key": t.Key, "Value": t.Value})
		}
		m["TagList"] = protocol.XMLElements{ElementName: "Tag", Items: tagItems}
	}
	if roles, ok := m["AssociatedRoles"].([]interface{}); ok && len(roles) > 0 {
		m["AssociatedRoles"] = protocol.XMLElements{ElementName: "DBClusterRole", Items: roles}
	}
	return m
}

func enrichClusterWithTagsAndEndpoint(store neptunestore.NeptuneStoreInterface, cluster *neptunestore.DBCluster, endpointAddr string, endpointPort int) map[string]interface{} {
	m := enrichClusterWithTags(store, cluster)
	if endpointAddr != "" && endpointPort > 0 {
		m["Endpoint"] = fmt.Sprintf("%s:%d", endpointAddr, endpointPort)
	} else if endpointAddr != "" {
		m["Endpoint"] = endpointAddr
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
	backupRetention := request.GetIntParam(params, "BackupRetentionPeriod")
	if backupRetention == 0 {
		backupRetention = 1
	}
	cluster := &neptunestore.DBCluster{
		DBClusterIdentifier:              id,
		Engine:                           engine,
		EngineVersion:                    request.GetStringParam(params, "EngineVersion"),
		Status:                           "available",
		Port:                             port,
		BackupRetentionPeriod:            backupRetention,
		PreferredBackupWindow:            request.GetStringParam(params, "PreferredBackupWindow"),
		PreferredMaintenanceWindow:       request.GetStringParam(params, "PreferredMaintenanceWindow"),
		MasterUsername:                   request.GetStringParam(params, "MasterUsername"),
		DatabaseName:                     request.GetStringParam(params, "DatabaseName"),
		DBClusterParameterGroupName:      request.GetStringParam(params, "DBClusterParameterGroupName"),
		DBSubnetGroupName:                request.GetStringParam(params, "DBSubnetGroupName"),
		StorageEncrypted:                 request.GetBoolParam(params, "StorageEncrypted"),
		KmsKeyId:                         request.GetStringParam(params, "KmsKeyId"),
		CopyTagsToSnapshot:               request.GetBoolParam(params, "CopyTagsToSnapshot"),
		DeletionProtection:               request.GetBoolParam(params, "DeletionProtection"),
		IAMDatabaseAuthenticationEnabled: request.GetBoolParam(params, "EnableIAMDatabaseAuthentication"),
		ClusterCreateTime:                &now,
		EarliestRestorableTime:           &now,
		LatestRestorableTime:             &now,
		ReplicationSourceIdentifier:      replicationSource,
		GlobalClusterIdentifier:          request.GetStringParam(params, "GlobalClusterIdentifier"),
		StorageType:                      request.GetStringParam(params, "StorageType"),
		AccountID:                        reqCtx.GetAccountID(),
		Region:                           reqCtx.GetRegion(),
		DBClusterArn:                     arnutil.NewARNBuilder(reqCtx.GetAccountID(), reqCtx.GetRegion()).RDS().Cluster(id),
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
		cluster.EnabledCloudwatchLogsExports = logExports
	}

	if err := store.CreateCluster(cluster); err != nil {
		return nil, translateStoreError(err)
	}

	var enginePort int
	if s.dataPlaneService != nil {
		if port, err := s.dataPlaneService.OpenClusterEngine(id); err != nil {
			logs.Warn("failed to open cluster engine", logs.String("cluster", id), logs.Err(err))
		} else {
			enginePort = port
		}
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
		"DBCluster": enrichClusterWithTagsAndEndpoint(store, cluster, s.endpointAddress(id), enginePort),
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

	if s.dataPlaneService != nil {
		if err := s.dataPlaneService.CloseClusterEngine(id); err != nil {
			logs.Warn("failed to close cluster engine", logs.String("cluster", id), logs.Err(err))
		}
	}

	if err := store.DeleteCluster(id); err != nil {
		return nil, translateStoreError(err)
	}

	recordEvent(store, "db-cluster", id, cluster.DBClusterArn,
		fmt.Sprintf("DB cluster %s deleted", id), []string{"deletion"})

	return map[string]interface{}{
		"DBCluster": enrichClusterWithTags(store, cluster),
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
		cluster.IAMDatabaseAuthenticationEnabled = request.GetBoolParam(params, "EnableIAMDatabaseAuthentication")
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
		cluster.EnabledCloudwatchLogsExports = request.GetStringList(params, "EnableCloudwatchLogsExports")
	}

	if err := store.UpdateCluster(cluster); err != nil {
		return nil, translateStoreError(err)
	}

	newID := request.GetStringParam(params, "NewDBClusterIdentifier")
	if newID != "" && newID != id {
		oldArn := cluster.DBClusterArn
		oldID := cluster.DBClusterIdentifier
		cluster.DBClusterIdentifier = newID
		cluster.DBClusterArn = arnutil.NewARNBuilder(cluster.AccountID, cluster.Region).RDS().Cluster(newID)
		reparentClusterResources(store, oldArn, cluster.DBClusterArn, oldID, newID)
		if err := store.CreateCluster(cluster); err != nil {
			return nil, translateStoreError(err)
		}
		if err := store.DeleteCluster(oldID); err != nil {
			_ = store.DeleteCluster(newID)
			return nil, translateStoreError(err)
		}
	}

	addr, port := s.clusterEndpointInfo(cluster.DBClusterIdentifier)
	return map[string]interface{}{
		"DBCluster": enrichClusterWithTagsAndEndpoint(store, cluster, addr, port),
	}, nil
}

func (s *NeptuneService) clusterEndpointInfo(clusterID string) (addr string, port int) {
	if s.dataPlaneService == nil {
		return "", 0
	}
	p, err := s.dataPlaneService.GetClusterPort(clusterID)
	if err != nil || p == 0 {
		return "", 0
	}
	return s.endpointAddress(clusterID), p
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
		addr, port := s.clusterEndpointInfo(clusterID)
		return map[string]interface{}{
			"DBClusters": protocol.XMLElements{ElementName: "DBCluster", Items: []interface{}{enrichClusterWithTagsAndEndpoint(store, cluster, addr, port)}},
		}, nil
	}

	clusters, err := store.ListClusters()
	if err != nil {
		return nil, translateStoreError(err)
	}

	items := make([]interface{}, 0, len(clusters))
	for _, c := range clusters {
		addr, port := s.clusterEndpointInfo(c.DBClusterIdentifier)
		items = append(items, enrichClusterWithTagsAndEndpoint(store, c, addr, port))
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

	var enginePort int
	if s.dataPlaneService != nil {
		if port, err := s.dataPlaneService.OpenClusterEngine(id); err != nil {
			logs.Warn("failed to open cluster engine on start", logs.String("cluster", id), logs.Err(err))
		} else {
			enginePort = port
		}
	}

	return map[string]interface{}{
		"DBCluster": enrichClusterWithTagsAndEndpoint(store, cluster, s.endpointAddress(id), enginePort),
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

	if s.dataPlaneService != nil {
		if err := s.dataPlaneService.CloseClusterEngine(id); err != nil {
			logs.Warn("failed to close cluster engine on stop", logs.String("cluster", id), logs.Err(err))
		}
	}

	cluster.Status = "stopped"
	if err := store.UpdateCluster(cluster); err != nil {
		return nil, translateStoreError(err)
	}

	return map[string]interface{}{
		"DBCluster": enrichClusterWithTagsAndEndpoint(store, cluster, "", 0),
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

	addr, port := s.clusterEndpointInfo(id)
	return map[string]interface{}{
		"DBCluster": enrichClusterWithTagsAndEndpoint(store, cluster, addr, port),
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
// buildRestoredCluster constructs a DBCluster for restore-from-snapshot or
// point-in-time restore, deriving defaults from the source when not overridden.
func buildRestoredCluster(clusterID, engine, engineVersion string, params map[string]interface{}, now *time.Time, reqCtx *request.RequestContext) *neptunestore.DBCluster {
	backupRetention := request.GetIntParam(params, "BackupRetentionPeriod")
	if backupRetention == 0 {
		backupRetention = 1
	}
	return &neptunestore.DBCluster{
		DBClusterIdentifier:         clusterID,
		Engine:                      engine,
		EngineVersion:               engineVersion,
		Status:                      "available",
		Port:                        request.GetIntParam(params, "Port"),
		BackupRetentionPeriod:       backupRetention,
		DBClusterParameterGroupName: request.GetStringParam(params, "DBClusterParameterGroupName"),
		DBSubnetGroupName:           request.GetStringParam(params, "DBSubnetGroupName"),
		StorageEncrypted:            request.GetBoolParam(params, "StorageEncrypted"),
		DeletionProtection:          request.GetBoolParam(params, "DeletionProtection"),
		ClusterCreateTime:           now,
		EarliestRestorableTime:      now,
		LatestRestorableTime:        now,
		AccountID:                   reqCtx.GetAccountID(),
		Region:                      reqCtx.GetRegion(),
		DBClusterArn:                arnutil.NewARNBuilder(reqCtx.GetAccountID(), reqCtx.GetRegion()).RDS().Cluster(clusterID),
	}
}

// reparentClusterResources migrates tags and updates instance references when
// a cluster is renamed. Errors are logged but not propagated so the rename
// itself always succeeds.
func reparentClusterResources(store neptunestore.NeptuneStoreInterface, oldArn, newArn, oldID, newID string) {
	tags, err := store.GetTags(oldArn)
	if err == nil && len(tags) > 0 {
		if err := store.AddTags(newArn, tags); err != nil {
			logs.Warn("reparent: failed to copy tags to new cluster ARN", logs.Err(err))
		} else {
			keys := make([]string, len(tags))
			for i, t := range tags {
				keys[i] = t.Key
			}
			if err := store.RemoveTags(oldArn, keys); err != nil {
				logs.Warn("reparent: failed to remove old cluster tags", logs.Err(err))
			}
		}
	}

	instances, err := store.ListInstances()
	if err != nil {
		logs.Warn("reparent: failed to list instances", logs.Err(err))
	} else {
		for _, inst := range instances {
			if inst.DBClusterIdentifier == oldID {
				inst.DBClusterIdentifier = newID
				if err := store.UpdateInstance(inst); err != nil {
					logs.Warn("reparent: failed to update instance cluster ref", logs.String("instance", inst.DBInstanceIdentifier), logs.Err(err))
				}
			}
		}
	}
}

func cascadeDeleteClusterResources(store neptunestore.NeptuneStoreInterface, cluster *neptunestore.DBCluster) {
	clusterID := cluster.DBClusterIdentifier

	instances, err := store.ListInstances()
	if err != nil {
		logs.Warn("cascade: failed to list instances", logs.Err(err))
	} else {
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
	}

	endpoints, err := store.ListClusterEndpoints(clusterID)
	if err != nil {
		logs.Warn("cascade: failed to list cluster endpoints", logs.Err(err))
	} else {
		for _, ep := range endpoints {
			if delErr := store.DeleteClusterEndpoint(ep.DBClusterEndpointIdentifier); delErr != nil {
				logs.Warn("cascade: failed to delete cluster endpoint", logs.String("endpoint", ep.DBClusterEndpointIdentifier), logs.Err(delErr))
			}
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
