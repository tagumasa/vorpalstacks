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
		// Edge platform convention: Endpoint includes port as "address:port"
		// so that clients can discover the dynamically assigned engine port.
		// This is consistent with MySQL instances and all other services that
		// use dynamic port allocation.
		m["Endpoint"] = fmt.Sprintf("%s:%d", c.Endpoint.Address, c.Endpoint.Port)
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

// setClusterEndpoint persists the connection endpoint on the cluster.  It
// constructs the address from endpointAddress, stores it in cluster.Endpoint,
// and writes the update through to Pebble.
func (s *NeptuneService) setClusterEndpoint(store neptunestore.NeptuneStoreInterface, cluster *neptunestore.DBCluster, enginePort int) {
	addr := s.endpointAddressFor(cluster.DBClusterIdentifier, cluster.Engine)
	if addr == "" || enginePort <= 0 {
		return
	}
	cluster.Endpoint = &neptunestore.Endpoint{Address: addr, Port: enginePort}
	if err := store.UpdateCluster(cluster); err != nil {
		logs.Warn("failed to persist cluster endpoint", logs.String("cluster", cluster.DBClusterIdentifier), logs.Err(err))
	}
}

// CreateDBCluster creates a new Neptune DB cluster with the specified configuration.
func (s *NeptuneService) CreateDBCluster(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	id := request.GetStringParam(params, "DBClusterIdentifier")
	if id == "" {
		return nil, awserrors.NewMissingParameter("DBClusterIdentifier is required")
	}
	if err := rdssvc.ValidateDBClusterIdentifier(id); err != nil {
		return nil, awserrors.NewAWSError("InvalidParameterValue", err.Error(), http.StatusBadRequest)
	}

	engine := request.GetStringParam(params, "Engine")
	if engine == "" {
		engine = "neptune"
	}
	if err := rdssvc.ValidateEngine(engine); err != nil {
		return nil, awserrors.NewAWSError("InvalidParameterValue", err.Error(), http.StatusBadRequest)
	}
	engineVersion := request.GetStringParam(params, "EngineVersion")
	if engineVersion == "" {
		engineVersion = rdssvc.DefaultEngineVersion(engine)
	}
	if err := rdssvc.ValidateEngineVersion(engine, engineVersion); err != nil {
		return nil, awserrors.NewAWSError("InvalidParameterValue", err.Error(), http.StatusBadRequest)
	}
	if dbName := request.GetStringParam(params, "DatabaseName"); dbName != "" {
		if err := rdssvc.ValidateDatabaseName(dbName); err != nil {
			return nil, awserrors.NewAWSError("InvalidParameterValue", err.Error(), http.StatusBadRequest)
		}
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
		EngineVersion:                    engineVersion,
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

	if cluster.GlobalClusterIdentifier != "" {
		if gc, err := store.GetGlobalCluster(cluster.GlobalClusterIdentifier); err == nil {
			isWriter := len(gc.GlobalClusterMembers) == 0
			gc.GlobalClusterMembers = append(gc.GlobalClusterMembers, neptunestore.GlobalClusterMember{
				DBClusterArn:            cluster.DBClusterArn,
				IsWriter:                isWriter,
				GlobalClusterIdentifier: gc.GlobalClusterIdentifier,
			})
			if err := store.UpdateGlobalCluster(gc); err != nil {
				logs.Warn("failed to register cluster as global cluster member", logs.String("cluster", id), logs.Err(err))
			}
		}
	}

	var enginePort int
	if eng := s.engineFor(engine); eng != nil {
		if port, err := eng.Open(reqCtx.GetRegion(), id); err != nil {
			logs.Warn("failed to open cluster engine", logs.String("cluster", id), logs.Err(err))
		} else {
			enginePort = port
		}
	}

	s.setClusterEndpoint(store, cluster, enginePort)

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
	if err := store.UpdateCluster(cluster); err != nil {
		return nil, translateStoreError(err)
	}

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
				logs.Warn("failed to remove cluster from global cluster members", logs.String("cluster", id), logs.Err(err))
			}
		}
	}

	if eng := s.engineFor(cluster.Engine); eng != nil {
		if err := eng.Close(id); err != nil {
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

	if v := request.GetStringParam(params, "EngineVersion"); v != "" {
		cluster.EngineVersion = v
	}
	if v := request.GetStringParam(params, "DBClusterParameterGroupName"); v != "" {
		cluster.DBClusterParameterGroupName = v
	}
	if v := request.GetIntParam(params, "Port"); v > 0 {
		cluster.Port = v
		if cluster.Endpoint != nil {
			cluster.Endpoint.Port = v
		}
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

	if newPort := request.GetIntParam(params, "Port"); newPort > 0 {
		s.setClusterEndpoint(store, cluster, newPort)
	}

	newID := request.GetStringParam(params, "NewDBClusterIdentifier")
	if newID != "" && newID != id {
		oldArn := cluster.DBClusterArn
		oldID := cluster.DBClusterIdentifier
		cluster.DBClusterIdentifier = newID
		cluster.DBClusterArn = arnutil.NewARNBuilder(cluster.AccountID, cluster.Region).RDS().Cluster(newID)
		if err := store.CreateCluster(cluster); err != nil {
			cluster.DBClusterIdentifier = oldID
			cluster.DBClusterArn = oldArn
			return nil, translateStoreError(err)
		}
		reparentClusterResources(store, oldArn, cluster.DBClusterArn, oldID, newID)
		if err := store.DeleteCluster(oldID); err != nil {
			logs.Warn("failed to delete old cluster after rename", logs.String("oldID", oldID), logs.Err(err))
		}
	}

	return map[string]interface{}{
		"DBCluster": enrichClusterWithTags(store, cluster),
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

	if cluster.Status != "stopped" {
		return nil, awserrors.NewAWSError("InvalidDBClusterStateFault", fmt.Sprintf("DBCluster %s is not in stopped state (current: %s)", id, cluster.Status), http.StatusBadRequest)
	}

	cluster.Status = "available"
	if err := store.UpdateCluster(cluster); err != nil {
		return nil, translateStoreError(err)
	}

	var enginePort int
	if eng := s.engineFor(cluster.Engine); eng != nil {
		if port, err := eng.Open(reqCtx.GetRegion(), id); err != nil {
			logs.Warn("failed to open cluster engine on start", logs.String("cluster", id), logs.Err(err))
		} else {
			enginePort = port
		}
	}

	s.setClusterEndpoint(store, cluster, enginePort)

	return map[string]interface{}{
		"DBCluster": enrichClusterWithTags(store, cluster),
	}, nil
}

// StopDBCluster stops a running Neptune DB cluster.
func (s *NeptuneService) StopDBCluster(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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

	if cluster.Status != "available" {
		return nil, awserrors.NewAWSError("InvalidDBClusterStateFault", fmt.Sprintf("DBCluster %s is not in available state (current: %s)", id, cluster.Status), http.StatusBadRequest)
	}

	if eng := s.engineFor(cluster.Engine); eng != nil {
		if err := eng.Close(id); err != nil {
			logs.Warn("failed to close cluster engine on stop", logs.String("cluster", id), logs.Err(err))
		}
	}

	cluster.Status = "stopped"
	if err := store.UpdateCluster(cluster); err != nil {
		return nil, translateStoreError(err)
	}

	return map[string]interface{}{
		"DBCluster": enrichClusterWithTags(store, cluster),
	}, nil
}

// FailoverDBCluster forces a failover for the specified DB cluster.
func (s *NeptuneService) FailoverDBCluster(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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

	cluster.Status = "available"
	if err := store.UpdateCluster(cluster); err != nil {
		return nil, translateStoreError(err)
	}

	return map[string]interface{}{
		"DBCluster": enrichClusterWithTags(store, cluster),
	}, nil
}

// AddRoleToDBCluster associates an IAM role with the specified DB cluster.
func (s *NeptuneService) AddRoleToDBCluster(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	id := request.GetStringParam(params, "DBClusterIdentifier")
	if id == "" {
		return nil, awserrors.NewMissingParameter("DBClusterIdentifier is required")
	}
	roleArn := request.GetStringParam(params, "RoleArn")
	if roleArn == "" {
		return nil, awserrors.NewMissingParameter("RoleArn is required")
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
			return nil, awserrors.NewAWSError("DBClusterRoleAlreadyExistsFault", fmt.Sprintf("IAM role %s is already associated with cluster %s", roleArn, id), http.StatusBadRequest)
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
		return nil, awserrors.NewMissingParameter("DBClusterIdentifier is required")
	}
	roleArn := request.GetStringParam(params, "RoleArn")
	if roleArn == "" {
		return nil, awserrors.NewMissingParameter("RoleArn is required")
	}
	featureName := request.GetStringParam(params, "FeatureName")

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
		if r.RoleArn == roleArn && (featureName == "" || r.FeatureName == featureName) {
			found = true
			continue
		}
		filtered = append(filtered, r)
	}
	if !found {
		return nil, awserrors.NewAWSError("DBClusterRoleNotFoundFault", fmt.Sprintf("role %s is not associated with cluster %s", roleArn, id), http.StatusBadRequest)
	}
	cluster.AssociatedRoles = filtered
	if err := store.UpdateCluster(cluster); err != nil {
		return nil, translateStoreError(err)
	}

	return map[string]interface{}{}, nil
}

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

// removeTagsForResource removes all tags associated with the given resource ARN.
func removeTagsForResource(store neptunestore.NeptuneStoreInterface, resourceArn string) {
	tags, err := store.GetTags(resourceArn)
	if err != nil || len(tags) == 0 {
		return
	}
	keys := make([]string, len(tags))
	for i, t := range tags {
		keys[i] = t.Key
	}
	if err := store.RemoveTags(resourceArn, keys); err != nil {
		logs.Warn("failed to remove tags on delete", logs.String("arn", resourceArn), logs.Err(err))
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
					removeTagsForResource(store, inst.DBInstanceArn)
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

	removeTagsForResource(store, cluster.DBClusterArn)
}
