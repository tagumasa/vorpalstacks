package neptune

import (
	"context"
	"fmt"
	"time"

	"vorpalstacks/internal/services/aws/common/protocol"
	"vorpalstacks/internal/services/aws/common/request"
	neptunestore "vorpalstacks/internal/store/aws/neptune"
)

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

	if _, err := store.GetCluster(id); err == nil {
		return nil, fmt.Errorf("neptune: DBCluster %s already exists", id)
	}

	now := time.Now()
	cluster := &neptunestore.DBCluster{
		DBClusterIdentifier:             id,
		Engine:                          engine,
		EngineVersion:                   request.GetStringParam(params, "EngineVersion"),
		Status:                          "creating",
		Port:                            request.GetIntParam(params, "Port"),
		BackupRetentionPeriod:           request.GetIntParam(params, "BackupRetentionPeriod"),
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
		ReplicationSourceIdentifier:     request.GetStringParam(params, "ReplicationSourceIdentifier"),
		GlobalClusterIdentifier:         request.GetStringParam(params, "GlobalClusterIdentifier"),
		StorageType:                     request.GetStringParam(params, "StorageType"),
		AccountID:                       reqCtx.GetAccountID(),
		Region:                          reqCtx.GetRegion(),
	}

	if azList := request.GetStringList(params, "AvailabilityZones"); len(azList) > 0 {
		cluster.AvailabilityZones = azList
	}
	if sgList := request.GetStringList(params, "VpcSecurityGroupIds"); len(sgList) > 0 {
		cluster.VpcSecurityGroupIds = sgList
	}
	if logExports := request.GetStringList(params, "EnableCloudwatchLogsExports"); len(logExports) > 0 {
		cluster.EnableCloudwatchLogsExports = logExports
	}

	if err := store.CreateCluster(cluster); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"DBCluster": cluster,
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
		return nil, err
	}

	if request.GetBoolParam(params, "SkipFinalSnapshot") {
		cluster.Status = "deleting"
	} else {
		cluster.Status = "deleting"
	}

	if err := store.UpdateCluster(cluster); err != nil {
		return nil, err
	}

	_ = store.DeleteCluster(id)

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
		return nil, err
	}

	if v := request.GetStringParam(params, "NewDBClusterIdentifier"); v != "" {
		cluster.DBClusterIdentifier = v
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
	if v := request.GetStringParam(params, "DeletionProtection"); v != "" {
		cluster.DeletionProtection = request.GetBoolParam(params, "DeletionProtection")
	}
	if request.GetBoolParam(params, "EnableIAMDatabaseAuthentication") {
		cluster.EnableIAMDatabaseAuthentication = true
	}

	if err := store.UpdateCluster(cluster); err != nil {
		return nil, err
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
			return nil, err
		}
		return map[string]interface{}{
			"DBClusters": protocol.XMLElements{ElementName: "DBCluster", Items: []interface{}{cluster}},
		}, nil
	}

	clusters, err := store.ListClusters()
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0, len(clusters))
	for _, c := range clusters {
		items = append(items, c)
	}

	return map[string]interface{}{
		"DBClusters": protocol.XMLElements{ElementName: "DBCluster", Items: items},
	}, nil
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
		return nil, err
	}

	cluster.Status = "starting"
	if err := store.UpdateCluster(cluster); err != nil {
		return nil, err
	}

	cluster.Status = "available"
	_ = store.UpdateCluster(cluster)

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
		return nil, err
	}

	cluster.Status = "stopped"
	if err := store.UpdateCluster(cluster); err != nil {
		return nil, err
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
		return nil, err
	}

	cluster.Status = "failing-over"
	_ = store.UpdateCluster(cluster)
	cluster.Status = "available"
	_ = store.UpdateCluster(cluster)

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
		return nil, err
	}

	featureName := request.GetStringParam(params, "FeatureName")
	cluster.AssociatedRoles = append(cluster.AssociatedRoles, neptunestore.DBClusterRole{
		RoleArn:     roleArn,
		FeatureName: featureName,
		Status:      "ACTIVE",
	})
	_ = store.UpdateCluster(cluster)

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
		return nil, err
	}

	filtered := cluster.AssociatedRoles[:0]
	for _, r := range cluster.AssociatedRoles {
		if r.RoleArn != roleArn {
			filtered = append(filtered, r)
		}
	}
	cluster.AssociatedRoles = filtered
	_ = store.UpdateCluster(cluster)

	return map[string]interface{}{}, nil
}
