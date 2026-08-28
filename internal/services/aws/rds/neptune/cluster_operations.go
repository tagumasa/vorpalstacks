package neptune

import (
	"context"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"vorpalstacks/internal/common/protocol"
	"vorpalstacks/internal/common/request"
	rdssvc "vorpalstacks/internal/services/aws/rds"
	neptunestore "vorpalstacks/internal/store/aws/rds/neptune"
	arnutil "vorpalstacks/internal/utils/aws/arn"
)

// hashMasterPassword bcrypts the given plaintext password and returns the
// hash string. An empty input yields an empty hash (no password set).
func hashMasterPassword(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintext), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash master password: %w", err)
	}
	return string(hash), nil
}

// isValidIAMRoleArn validates that the given string is a well-formed IAM
// role ARN (arn:aws:iam::<account>:role/<name>).
func isValidIAMRoleArn(arn string) bool {
	partition, service, _, _, resource := arnutil.SplitARN(arn)
	if (partition != "aws" && partition != "aws-cn" && partition != "aws-us-gov") || service != "iam" {
		return false
	}
	return strings.HasPrefix(resource, "role/")
}

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
		// Previously dropped output fields.
		"AllocatedStorage": c.AllocatedStorage,
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
		// Use Smithy-correct output key "VpcSecurityGroups" (was "VpcSecurityGroupIds").
		m["VpcSecurityGroups"] = protocol.XMLElements{ElementName: "VpcSecurityGroupMembership", Items: vpcSecurityToInterface(c.VpcSecurityGroupIds, c.Status)}
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
	// Previously dropped output fields (conditional).
	if c.DbClusterResourceId != "" {
		m["DbClusterResourceId"] = c.DbClusterResourceId
	}
	if c.NetworkType != "" {
		m["NetworkType"] = c.NetworkType
	}
	if c.PercentProgress != "" {
		m["PercentProgress"] = c.PercentProgress
	}
	if c.HostedZoneId != "" {
		m["HostedZoneId"] = c.HostedZoneId
	}
	if len(c.ReadReplicaIdentifiers) > 0 {
		m["ReadReplicaIdentifiers"] = protocol.XMLElements{ElementName: "ReadReplicaIdentifier", Items: stringSliceToInterface(c.ReadReplicaIdentifiers)}
	}
	if c.AutomaticRestartTime != nil {
		m["AutomaticRestartTime"] = c.AutomaticRestartTime.Format(time.RFC3339)
	}
	if c.PendingModifiedValues != nil {
		m["PendingModifiedValues"] = c.PendingModifiedValues
	}
	if len(c.DBClusterMembers) > 0 {
		members := make([]interface{}, 0, len(c.DBClusterMembers))
		for _, mem := range c.DBClusterMembers {
			members = append(members, map[string]interface{}{
				"DBInstanceIdentifier":          mem.DBInstanceIdentifier,
				"IsClusterWriter":               mem.IsClusterWriter,
				"DBClusterParameterGroupStatus": mem.DBClusterParameterGroupStatus,
				"PromotionTier":                 mem.PromotionTier,
			})
		}
		m["DBClusterMembers"] = protocol.XMLElements{ElementName: "DBClusterMember", Items: members}
	}
	if c.ReaderEndpoint != nil {
		m["ReaderEndpoint"] = fmt.Sprintf("%s:%d", c.ReaderEndpoint.Address, c.ReaderEndpoint.Port)
	}
	return m
}

// vpcSecurityToInterface converts VPC security group IDs to the AWS response
// format, deriving the Status from the cluster's lifecycle state instead of
// hardcoding "active".
func vpcSecurityToInterface(ids []string, clusterStatus string) []interface{} {
	sgStatus := "active"
	switch clusterStatus {
	case "creating":
		sgStatus = "adding"
	case "deleting":
		sgStatus = "removing"
	}
	result := make([]interface{}, len(ids))
	for i, id := range ids {
		result[i] = map[string]interface{}{"VpcSecurityGroupId": id, "Status": sgStatus}
	}
	return result
}

// CreateDBCluster creates a new Neptune DB cluster with the specified configuration.
func (s *NeptuneService) CreateDBCluster(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	in := &CreateDBClusterInput{
		CreateClusterParams: rdssvc.CreateClusterParams{
			DBClusterIdentifier:          request.GetStringParam(params, "DBClusterIdentifier"),
			Engine:                       request.GetStringParam(params, "Engine"),
			EngineVersion:                rdssvc.ResolveEngineVersion(request.GetStringParam(params, "Engine"), request.GetStringParam(params, "EngineVersion")),
			DatabaseName:                 request.GetStringParam(params, "DatabaseName"),
			MasterUsername:               request.GetStringParam(params, "MasterUsername"),
			Port:                         request.GetIntParam(params, "Port"),
			BackupRetentionPeriod:        request.GetIntParam(params, "BackupRetentionPeriod"),
			AvailabilityZones:            getNeptuneStringList(params, "AvailabilityZones", "AvailabilityZone", "member"),
			DBSubnetGroupName:            request.GetStringParam(params, "DBSubnetGroupName"),
			DBClusterParameterGroupName:  request.GetStringParam(params, "DBClusterParameterGroupName"),
			StorageEncrypted:             request.GetBoolParam(params, "StorageEncrypted"),
			CopyTagsToSnapshot:           request.GetBoolParam(params, "CopyTagsToSnapshot"),
			DeletionProtection:           request.GetBoolParam(params, "DeletionProtection"),
			IAMDatabaseAuthentication:    request.GetBoolParam(params, "EnableIAMDatabaseAuthentication"),
			EnabledCloudwatchLogsExports: request.GetStringList(params, "EnableCloudwatchLogsExports"),
			AccountID:                    reqCtx.GetAccountID(),
			Region:                       reqCtx.GetRegion(),
		},
		MasterUserPassword:          request.GetStringParam(params, "MasterUserPassword"),
		PreferredBackupWindow:       request.GetStringParam(params, "PreferredBackupWindow"),
		PreferredMaintenanceWindow:  request.GetStringParam(params, "PreferredMaintenanceWindow"),
		KmsKeyId:                    request.GetStringParam(params, "KmsKeyId"),
		ReplicationSourceIdentifier: request.GetStringParam(params, "ReplicationSourceIdentifier"),
		GlobalClusterIdentifier:     request.GetStringParam(params, "GlobalClusterIdentifier"),
		StorageType:                 request.GetStringParam(params, "StorageType"),
		VpcSecurityGroupIds:         getNeptuneStringList(params, "VpcSecurityGroupIds", "VpcSecurityGroupId", "member"),
		Tags:                        parseNeptuneTags(params),
	}
	if svsc := request.GetMapParam(params, "ServerlessV2ScalingConfiguration"); svsc != nil {
		in.HasServerlessV2Scaling = true
		in.ServerlessV2MinCapacity = request.GetFloatParam(svsc, "MinCapacity")
		in.ServerlessV2MaxCapacity = request.GetFloatParam(svsc, "MaxCapacity")
	}
	return s.createDBClusterCore(ctx, store, in)
}

// DeleteDBCluster deletes the specified Neptune DB cluster.
func (s *NeptuneService) DeleteDBCluster(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	in := &DeleteDBClusterInput{
		DeleteClusterParams: rdssvc.DeleteClusterParams{
			DBClusterIdentifier:       request.GetStringParam(params, "DBClusterIdentifier"),
			SkipFinalSnapshot:         request.GetBoolParam(params, "SkipFinalSnapshot"),
			FinalDBSnapshotIdentifier: request.GetStringParam(params, "FinalDBSnapshotIdentifier"),
			AccountID:                 reqCtx.GetAccountID(),
			Region:                    reqCtx.GetRegion(),
		},
	}
	return s.deleteDBClusterCore(ctx, store, in)
}

// ModifyDBCluster updates the configuration of the specified DB cluster.
func (s *NeptuneService) ModifyDBCluster(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	in := &ModifyDBClusterInput{
		DBClusterIdentifier:             request.GetStringParam(params, "DBClusterIdentifier"),
		EngineVersion:                   request.GetStringParam(params, "EngineVersion"),
		DBClusterParameterGroupName:     request.GetStringParam(params, "DBClusterParameterGroupName"),
		Port:                            request.GetIntParam(params, "Port"),
		BackupRetentionPeriod:           request.GetIntParam(params, "BackupRetentionPeriod"),
		PreferredBackupWindow:           request.GetStringParam(params, "PreferredBackupWindow"),
		PreferredMaintenanceWindow:      request.GetStringParam(params, "PreferredMaintenanceWindow"),
		StorageType:                     request.GetStringParam(params, "StorageType"),
		MasterUserPassword:              request.GetStringParam(params, "MasterUserPassword"),
		NetworkType:                     request.GetStringParam(params, "NetworkType"),
		HasDeletionProtection:           request.HasParam(params, "DeletionProtection"),
		DeletionProtection:              request.GetBoolParam(params, "DeletionProtection"),
		HasEnableIAMDatabaseAuth:        request.HasParam(params, "EnableIAMDatabaseAuthentication"),
		EnableIAMDatabaseAuthentication: request.GetBoolParam(params, "EnableIAMDatabaseAuthentication"),
		VpcSecurityGroupIds:             getNeptuneStringList(params, "VpcSecurityGroupIds", "VpcSecurityGroupId", "member"),
		EnabledCloudwatchLogsExports:    request.GetStringList(params, "EnableCloudwatchLogsExports"),
		NewDBClusterIdentifier:          request.GetStringParam(params, "NewDBClusterIdentifier"),
		Region:                          reqCtx.GetRegion(),
	}
	if svsc := request.GetMapParam(params, "ServerlessV2ScalingConfiguration"); svsc != nil {
		in.HasServerlessV2Scaling = true
		in.ServerlessV2MinCapacity = request.GetFloatParam(svsc, "MinCapacity")
		in.ServerlessV2MaxCapacity = request.GetFloatParam(svsc, "MaxCapacity")
	}
	return s.modifyDBClusterCore(ctx, store, in)
}

// DescribeDBClusters returns information about the specified DB cluster or lists
// all clusters when no identifier is provided.
func (s *NeptuneService) DescribeDBClusters(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	clusters, nextMarker, err := rdssvc.QueryClusters(store, rdssvc.DescribeDBClustersInput{
		DBClusterIdentifier: request.GetStringParam(params, "DBClusterIdentifier"),
		Filters:             nil,
		Marker:              request.GetStringParam(params, "Marker"),
		MaxRecords:          int32(request.GetIntParam(params, "MaxRecords")),
	})
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0, len(clusters))
	for _, c := range clusters {
		items = append(items, enrichClusterWithTags(store, c))
	}

	result := map[string]interface{}{
		"DBClusters": protocol.XMLElements{ElementName: "DBCluster", Items: items},
	}
	if nextMarker != "" {
		result["Marker"] = nextMarker
	}
	return result, nil
}

// StartDBCluster restarts a stopped Neptune DB cluster.
func (s *NeptuneService) StartDBCluster(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	in := &StartDBClusterInput{
		DBClusterIdentifier: request.GetStringParam(req.Parameters, "DBClusterIdentifier"),
		Region:              reqCtx.GetRegion(),
	}
	return s.startDBClusterCore(ctx, store, in)
}

// StopDBCluster stops a running Neptune DB cluster.
func (s *NeptuneService) StopDBCluster(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	in := &StopDBClusterInput{
		DBClusterIdentifier: request.GetStringParam(req.Parameters, "DBClusterIdentifier"),
	}
	return s.stopDBClusterCore(ctx, store, in)
}

// FailoverDBCluster forces a failover for the specified DB cluster.
func (s *NeptuneService) FailoverDBCluster(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	in := &FailoverDBClusterInput{
		DBClusterIdentifier: request.GetStringParam(req.Parameters, "DBClusterIdentifier"),
	}
	return s.failoverDBClusterCore(ctx, store, in)
}

// AddRoleToDBCluster associates an IAM role with the specified DB cluster.
func (s *NeptuneService) AddRoleToDBCluster(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	in := &AddRoleToDBClusterInput{
		DBClusterIdentifier: request.GetStringParam(req.Parameters, "DBClusterIdentifier"),
		RoleArn:             request.GetStringParam(req.Parameters, "RoleArn"),
		FeatureName:         request.GetStringParam(req.Parameters, "FeatureName"),
	}
	return s.addRoleToDBClusterCore(ctx, store, in)
}

// RemoveRoleFromDBCluster disassociates an IAM role from the specified DB cluster.
func (s *NeptuneService) RemoveRoleFromDBCluster(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	in := &RemoveRoleFromDBClusterInput{
		DBClusterIdentifier: request.GetStringParam(req.Parameters, "DBClusterIdentifier"),
		RoleArn:             request.GetStringParam(req.Parameters, "RoleArn"),
		FeatureName:         request.GetStringParam(req.Parameters, "FeatureName"),
	}
	return s.removeRoleFromDBClusterCore(ctx, store, in)
}
