package neptune

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/protocol"
	"vorpalstacks/internal/common/request"
	types "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/core/logs"
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
	// ReaderEndpoint mirrors the cluster endpoint for Neptune's
	// single-writer topology. AWS Neptune surfaces both endpoints.
	if cluster.ReaderEndpoint == nil {
		cluster.ReaderEndpoint = &neptunestore.Endpoint{Address: addr, Port: enginePort}
	}
	if err := store.UpdateCluster(cluster); err != nil {
		logs.Warn("failed to persist cluster endpoint", logs.String("cluster", cluster.DBClusterIdentifier), logs.Err(err))
	}
}

// CreateDBCluster creates a new Neptune DB cluster with the specified configuration.
func (s *NeptuneService) CreateDBCluster(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	engineVersion := rdssvc.ResolveEngineVersion(request.GetStringParam(params, "Engine"), request.GetStringParam(params, "EngineVersion"))

	createParams := rdssvc.CreateClusterParams{
		DBClusterIdentifier:          request.GetStringParam(params, "DBClusterIdentifier"),
		Engine:                       request.GetStringParam(params, "Engine"),
		EngineVersion:                engineVersion,
		DatabaseName:                 request.GetStringParam(params, "DatabaseName"),
		MasterUsername:               request.GetStringParam(params, "MasterUsername"),
		Port:                         request.GetIntParam(params, "Port"),
		BackupRetentionPeriod:        request.GetIntParam(params, "BackupRetentionPeriod"),
		AvailabilityZones:            request.GetStringList(params, "AvailabilityZones"),
		DBSubnetGroupName:            request.GetStringParam(params, "DBSubnetGroupName"),
		DBClusterParameterGroupName:  request.GetStringParam(params, "DBClusterParameterGroupName"),
		StorageEncrypted:             request.GetBoolParam(params, "StorageEncrypted"),
		CopyTagsToSnapshot:           request.GetBoolParam(params, "CopyTagsToSnapshot"),
		DeletionProtection:           request.GetBoolParam(params, "DeletionProtection"),
		IAMDatabaseAuthentication:    request.GetBoolParam(params, "EnableIAMDatabaseAuthentication"),
		EnabledCloudwatchLogsExports: request.GetStringList(params, "EnableCloudwatchLogsExports"),
		AccountID:                    reqCtx.GetAccountID(),
		Region:                       reqCtx.GetRegion(),
	}

	if err := rdssvc.ValidateCreateClusterParams(store, createParams); err != nil {
		return nil, neptuneTranslateError(err)
	}

	replicationSource := request.GetStringParam(params, "ReplicationSourceIdentifier")
	if replicationSource != "" {
		if _, err := store.GetCluster(replicationSource); err != nil {
			return nil, awserrors.NewAWSError("InvalidParameterValue", fmt.Sprintf("replication source cluster %s not found", replicationSource), http.StatusBadRequest)
		}
	}

	backupRetention := createParams.BackupRetentionPeriod
	if backupRetention == 0 {
		backupRetention = 1
	}
	if err := validateBackupRetentionPeriod(backupRetention); err != nil {
		return nil, awserrors.NewAWSError("InvalidParameterValue", err.Error(), http.StatusBadRequest)
	}

	masterPassword := request.GetStringParam(params, "MasterUserPassword")
	masterPasswordHash, err := hashMasterPassword(masterPassword)
	if err != nil {
		return nil, awserrors.NewAWSError("InvalidParameterValue", err.Error(), http.StatusBadRequest)
	}

	cluster := rdssvc.BuildCluster(createParams)
	cluster.BackupRetentionPeriod = backupRetention
	cluster.PreferredBackupWindow = request.GetStringParam(params, "PreferredBackupWindow")
	cluster.PreferredMaintenanceWindow = request.GetStringParam(params, "PreferredMaintenanceWindow")
	cluster.KmsKeyId = request.GetStringParam(params, "KmsKeyId")
	cluster.ReplicationSourceIdentifier = replicationSource
	cluster.GlobalClusterIdentifier = request.GetStringParam(params, "GlobalClusterIdentifier")
	cluster.StorageType = request.GetStringParam(params, "StorageType")
	cluster.MasterUserPasswordHash = masterPasswordHash
	cluster.DbClusterResourceId = fmt.Sprintf("cluster-%s", cluster.DBClusterIdentifier)
	cluster.NetworkType = "IPV4"
	cluster.EarliestRestorableTime = cluster.ClusterCreateTime
	cluster.LatestRestorableTime = cluster.ClusterCreateTime

	if sgList := request.GetStringList(params, "VpcSecurityGroupIds"); len(sgList) > 0 {
		if _, err := s.resolveSecurityGroups(ctx, reqCtx.GetRegion(), sgList); err != nil {
			return nil, translateStoreError(err)
		}
		cluster.VpcSecurityGroupIds = sgList
	}
	if svsc := request.GetMapParam(params, "ServerlessV2ScalingConfiguration"); svsc != nil {
		minCap := request.GetFloatParam(svsc, "MinCapacity")
		maxCap := request.GetFloatParam(svsc, "MaxCapacity")
		if minCap < 0.5 || minCap > 128 || maxCap < 1 || maxCap > 256 || minCap >= maxCap {
			return nil, awserrors.NewAWSError("InvalidParameterValue", "ServerlessV2ScalingConfiguration: MinCapacity must be 0.5-128, MaxCapacity 1-256, and MinCapacity < MaxCapacity", http.StatusBadRequest)
		}
		cluster.ServerlessV2ScalingConfiguration = &neptunestore.ServerlessV2ScalingConfiguration{
			MinCapacity: minCap,
			MaxCapacity: maxCap,
		}
	}

	if err := store.CreateCluster(cluster); err != nil {
		return nil, translateStoreError(err)
	}

	if cluster.GlobalClusterIdentifier != "" {
		if gc, err := store.GetGlobalCluster(cluster.GlobalClusterIdentifier); err == nil {
			isWriter := len(gc.GlobalClusterMembers) == 0
			if !isWriter {
				hasWriter := false
				for _, m := range gc.GlobalClusterMembers {
					if m.IsWriter {
						hasWriter = true
						break
					}
				}
				if !hasWriter {
					isWriter = true
				}
			}
			gc.GlobalClusterMembers = append(gc.GlobalClusterMembers, neptunestore.GlobalClusterMember{
				DBClusterArn:            cluster.DBClusterArn,
				IsWriter:                isWriter,
				GlobalClusterIdentifier: gc.GlobalClusterIdentifier,
			})
			if err := store.UpdateGlobalCluster(gc); err != nil {
				logs.Warn("failed to register cluster as global cluster member", logs.String("cluster", cluster.DBClusterIdentifier), logs.Err(err))
			}
		}
	}

	var enginePort int
	if eng := s.engineFor(cluster.Engine); eng != nil {
		if port, err := eng.Open(reqCtx.GetRegion(), cluster.DBClusterIdentifier); err != nil {
			logs.Warn("failed to open cluster engine", logs.String("cluster", cluster.DBClusterIdentifier), logs.Err(err))
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
			if eng := s.engineFor(cluster.Engine); eng != nil {
				eng.Close(cluster.DBClusterIdentifier)
			}
			removeClusterFromGlobal(store, cluster)
			store.DeleteCluster(cluster.DBClusterIdentifier)
			return nil, awserrors.NewAWSError("InvalidParameterValue", fmt.Sprintf("Failed to tag cluster: %v", err), http.StatusBadRequest)
		}
	}

	recordEvent(store, "db-cluster", cluster.DBClusterIdentifier, cluster.DBClusterArn,
		fmt.Sprintf("DB cluster %s created", cluster.DBClusterIdentifier), []string{"creation"})

	cluster.Status = "available"
	if err := store.UpdateCluster(cluster); err != nil {
		logs.Warn("failed to transition cluster to available", logs.String("cluster", cluster.DBClusterIdentifier), logs.Err(err))
	}
	s.scheduleTransition(reqCtx.GetRegion(), 500*time.Millisecond, func(st neptunestore.NeptuneStoreInterface) error {
		c, err := st.GetCluster(cluster.DBClusterIdentifier)
		if err != nil || c.Status != "creating" {
			return nil
		}
		c.Status = "available"
		return st.UpdateCluster(c)
	})

	return map[string]interface{}{
		"DBCluster": enrichClusterWithTags(store, cluster),
	}, nil
}

// DeleteDBCluster deletes the specified Neptune DB cluster.
func (s *NeptuneService) DeleteDBCluster(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	cluster, err := rdssvc.ValidateDeleteClusterParams(store, rdssvc.DeleteClusterParams{
		DBClusterIdentifier:       request.GetStringParam(params, "DBClusterIdentifier"),
		SkipFinalSnapshot:         request.GetBoolParam(params, "SkipFinalSnapshot"),
		FinalDBSnapshotIdentifier: request.GetStringParam(params, "FinalDBSnapshotIdentifier"),
		AccountID:                 reqCtx.GetAccountID(),
		Region:                    reqCtx.GetRegion(),
	})
	if err != nil {
		return nil, neptuneTranslateError(err)
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
		snapshot := rdssvc.BuildFinalSnapshot(cluster, finalSnapshotID, reqCtx.GetAccountID(), reqCtx.GetRegion())
		if err := store.CreateSnapshot(snapshot); err != nil {
			cluster.Status = "available"
			store.UpdateCluster(cluster)
			return nil, translateStoreError(err)
		}
	}

	if err := store.DeleteCluster(cluster.DBClusterIdentifier); err != nil {
		cluster.Status = "available"
		if rbErr := store.UpdateCluster(cluster); rbErr != nil {
			logs.Warn("failed to roll back cluster status after delete failure", logs.String("cluster", cluster.DBClusterIdentifier), logs.Err(rbErr))
		}
		if !skipFinal {
			if delErr := store.DeleteSnapshot(finalSnapshotID); delErr != nil {
				logs.Warn("failed to clean up orphaned snapshot after delete failure", logs.String("snapshot", finalSnapshotID), logs.Err(delErr))
			}
		}
		return nil, translateStoreError(err)
	}

	cascadeDeleteClusterResources(store, cluster)
	removeClusterFromGlobal(store, cluster)

	if eng := s.engineFor(cluster.Engine); eng != nil {
		if err := eng.Close(cluster.DBClusterIdentifier); err != nil {
			logs.Warn("failed to close cluster engine", logs.String("cluster", cluster.DBClusterIdentifier), logs.Err(err))
		}
	}

	recordEvent(store, "db-cluster", cluster.DBClusterIdentifier, cluster.DBClusterArn,
		fmt.Sprintf("DB cluster %s deleted", cluster.DBClusterIdentifier), []string{"deletion"})

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
		// Validate referenced parameter group exists before assigning.
		if _, err := store.GetClusterParameterGroup(v); err != nil {
			return nil, awserrors.NewAWSError("DBClusterParameterGroupNotFoundFault", fmt.Sprintf("DB Cluster Parameter Group not found: %s", v), http.StatusNotFound)
		}
		cluster.DBClusterParameterGroupName = v
	}
	if v := request.GetIntParam(params, "Port"); v > 0 {
		// Validate Port range on modify.
		if err := validatePort(v); err != nil {
			return nil, awserrors.NewAWSError("InvalidParameterValue", err.Error(), http.StatusBadRequest)
		}
		cluster.Port = v
		if cluster.Endpoint != nil {
			cluster.Endpoint.Port = v
		}
	}
	if v := request.GetIntParam(params, "BackupRetentionPeriod"); v > 0 {
		// Validate BackupRetentionPeriod range on modify.
		if err := validateBackupRetentionPeriod(v); err != nil {
			return nil, awserrors.NewAWSError("InvalidParameterValue", err.Error(), http.StatusBadRequest)
		}
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
	// Accept MasterUserPassword on modify and store as bcrypt hash.
	if pwd := request.GetStringParam(params, "MasterUserPassword"); pwd != "" {
		hash, hashErr := hashMasterPassword(pwd)
		if hashErr != nil {
			return nil, awserrors.NewAWSError("InvalidParameterValue", hashErr.Error(), http.StatusBadRequest)
		}
		cluster.MasterUserPasswordHash = hash
	}
	// Handle NetworkType on modify (Smithy ModifyDBClusterMessage member).
	if nt := request.GetStringParam(params, "NetworkType"); nt != "" {
		cluster.NetworkType = nt
	}
	// Handle ServerlessV2ScalingConfiguration on modify.
	if svsc := request.GetMapParam(params, "ServerlessV2ScalingConfiguration"); svsc != nil {
		minCap := request.GetFloatParam(svsc, "MinCapacity")
		maxCap := request.GetFloatParam(svsc, "MaxCapacity")
		if minCap < 0.5 || minCap > 128 || maxCap < 1 || maxCap > 256 || minCap >= maxCap {
			return nil, awserrors.NewAWSError("InvalidParameterValue", "ServerlessV2ScalingConfiguration: MinCapacity must be 0.5-128, MaxCapacity 1-256, and MinCapacity < MaxCapacity", http.StatusBadRequest)
		}
		cluster.ServerlessV2ScalingConfiguration = &neptunestore.ServerlessV2ScalingConfiguration{
			MinCapacity: minCap,
			MaxCapacity: maxCap,
		}
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
		if err := reparentClusterResources(store, oldArn, cluster.DBClusterArn, oldID, newID); err != nil {
			store.DeleteCluster(newID)
			cluster.DBClusterIdentifier = oldID
			cluster.DBClusterArn = oldArn
			return nil, awserrors.NewAWSError("InvalidDBClusterStateFault",
				fmt.Sprintf("cluster rename failed during resource reparenting: %v", err), http.StatusBadRequest)
		}
		if err := store.DeleteCluster(oldID); err != nil {
			logs.Error("failed to delete old cluster record after successful rename",
				logs.String("oldID", oldID), logs.String("newID", newID), logs.Err(err))
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

	// FailoverDBCluster requires the cluster to be in 'available' state.
	if cluster.Status != "available" {
		return nil, awserrors.NewAWSError("InvalidDBClusterStateFault", fmt.Sprintf("DBCluster %s is not in available state (current: %s)", id, cluster.Status), http.StatusBadRequest)
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
	// Validate IAM role ARN format.
	if !isValidIAMRoleArn(roleArn) {
		return nil, awserrors.NewAWSError("InvalidParameterValue", fmt.Sprintf("Invalid IAM role ARN: %s", roleArn), http.StatusBadRequest)
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
// a cluster is renamed. Tag-copy, instance-list, and instance-update failures
// are propagated so the caller can roll back the rename. Old-tag removal is
// best-effort because the tags have already been copied to the new ARN.
func reparentClusterResources(store neptunestore.NeptuneStoreInterface, oldArn, newArn, oldID, newID string) error {
	tags, err := store.GetTags(oldArn)
	if err != nil {
		return fmt.Errorf("reparent: failed to get tags from %s: %w", oldArn, err)
	}
	if len(tags) > 0 {
		if err := store.AddTags(newArn, tags); err != nil {
			return fmt.Errorf("reparent: failed to copy tags to %s: %w", newArn, err)
		}
		keys := make([]string, len(tags))
		for i, t := range tags {
			keys[i] = t.Key
		}
		if err := store.RemoveTags(oldArn, keys); err != nil {
			logs.Warn("reparent: failed to remove old cluster tags after copy", logs.Err(err))
		}
	}

	instances, err := store.ListInstances()
	if err != nil {
		return fmt.Errorf("reparent: failed to list instances: %w", err)
	}
	for _, inst := range instances {
		if inst.DBClusterIdentifier == oldID {
			inst.DBClusterIdentifier = newID
			if err := store.UpdateInstance(inst); err != nil {
				return fmt.Errorf("reparent: failed to update instance %s cluster ref: %w", inst.DBInstanceIdentifier, err)
			}
		}
	}
	return nil
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

// removeClusterFromGlobal removes the cluster's membership entry from its
// parent global cluster, if any. Used by DeleteDBCluster and by the
// CreateDBCluster tag-failure rollback path.
func removeClusterFromGlobal(store neptunestore.NeptuneStoreInterface, cluster *neptunestore.DBCluster) {
	if cluster.GlobalClusterIdentifier == "" {
		return
	}
	gc, err := store.GetGlobalCluster(cluster.GlobalClusterIdentifier)
	if err != nil {
		return
	}
	filtered := make([]neptunestore.GlobalClusterMember, 0, len(gc.GlobalClusterMembers))
	for _, m := range gc.GlobalClusterMembers {
		if m.DBClusterArn != cluster.DBClusterArn {
			filtered = append(filtered, m)
		}
	}
	gc.GlobalClusterMembers = filtered
	if err := store.UpdateGlobalCluster(gc); err != nil {
		logs.Warn("failed to remove cluster from global cluster members",
			logs.String("cluster", cluster.DBClusterIdentifier), logs.Err(err))
	}
}
