package rds

import (
	pb "vorpalstacks/internal/pb/storage/storage_rds"
	"vorpalstacks/internal/utils/aws/types"

	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

func clusterRoleToProto(r *DBClusterRole) *pb.DBClusterRole {
	if r == nil {
		return nil
	}
	return &pb.DBClusterRole{
		RoleArn:     r.RoleArn,
		FeatureName: r.FeatureName,
		Status:      r.Status,
	}
}

func protoToClusterRole(r *pb.DBClusterRole) *DBClusterRole {
	if r == nil {
		return nil
	}
	return &DBClusterRole{
		RoleArn:     r.GetRoleArn(),
		FeatureName: r.GetFeatureName(),
		Status:      r.GetStatus(),
	}
}

func scalingConfigToProto(c *ServerlessV2ScalingConfiguration) *pb.ServerlessV2ScalingConfiguration {
	if c == nil {
		return nil
	}
	return &pb.ServerlessV2ScalingConfiguration{
		MinCapacity: c.MinCapacity,
		MaxCapacity: c.MaxCapacity,
	}
}

func protoToScalingConfig(c *pb.ServerlessV2ScalingConfiguration) *ServerlessV2ScalingConfiguration {
	if c == nil {
		return nil
	}
	return &ServerlessV2ScalingConfiguration{
		MinCapacity: c.GetMinCapacity(),
		MaxCapacity: c.GetMaxCapacity(),
	}
}

func dbClusterMemberToProto(m *DBClusterMember) *pb.DBClusterMember {
	if m == nil {
		return nil
	}
	return &pb.DBClusterMember{
		DbInstanceIdentifier:          m.DBInstanceIdentifier,
		IsClusterWriter:               m.IsClusterWriter,
		DbClusterParameterGroupStatus: m.DBClusterParameterGroupStatus,
		PromotionTier:                 m.PromotionTier,
	}
}

func protoToDBClusterMember(m *pb.DBClusterMember) *DBClusterMember {
	if m == nil {
		return nil
	}
	return &DBClusterMember{
		DBInstanceIdentifier:          m.GetDbInstanceIdentifier(),
		IsClusterWriter:               m.GetIsClusterWriter(),
		DBClusterParameterGroupStatus: m.GetDbClusterParameterGroupStatus(),
		PromotionTier:                 m.GetPromotionTier(),
	}
}

func clusterPendingModifiedValuesToProto(v *ClusterPendingModifiedValues) *pb.ClusterPendingModifiedValues {
	if v == nil {
		return nil
	}
	p := &pb.ClusterPendingModifiedValues{
		DbClusterIdentifier: v.DBClusterIdentifier,
		EngineVersion:       v.EngineVersion,
		StorageType:         v.StorageType,
		NetworkType:         v.NetworkType,
	}
	if v.IAMDatabaseAuthenticationEnabled != nil {
		p.IamDatabaseAuthenticationEnabled = *v.IAMDatabaseAuthenticationEnabled
	}
	if v.BackupRetentionPeriod != nil {
		p.BackupRetentionPeriod = int32(*v.BackupRetentionPeriod)
	}
	if v.AllocatedStorage != nil {
		p.AllocatedStorage = *v.AllocatedStorage
	}
	if v.Iops != nil {
		p.Iops = *v.Iops
	}
	return p
}

func protoToClusterPendingModifiedValues(p *pb.ClusterPendingModifiedValues) *ClusterPendingModifiedValues {
	if p == nil {
		return nil
	}
	v := &ClusterPendingModifiedValues{
		DBClusterIdentifier: p.GetDbClusterIdentifier(),
		EngineVersion:       p.GetEngineVersion(),
		StorageType:         p.GetStorageType(),
		NetworkType:         p.GetNetworkType(),
	}
	iamAuth := p.GetIamDatabaseAuthenticationEnabled()
	v.IAMDatabaseAuthenticationEnabled = &iamAuth
	brp := int(p.GetBackupRetentionPeriod())
	v.BackupRetentionPeriod = &brp
	alloc := p.GetAllocatedStorage()
	v.AllocatedStorage = &alloc
	iops := p.GetIops()
	v.Iops = &iops
	return v
}

func ClusterToProto(c *DBCluster) *pb.DBCluster {
	if c == nil {
		return nil
	}
	p := &pb.DBCluster{
		DbClusterIdentifier:              c.DBClusterIdentifier,
		Engine:                           c.Engine,
		EngineVersion:                    c.EngineVersion,
		Status:                           c.Status,
		MasterUsername:                   c.MasterUsername,
		DatabaseName:                     c.DatabaseName,
		Port:                             int32(c.Port),
		BackupRetentionPeriod:            int32(c.BackupRetentionPeriod),
		PreferredBackupWindow:            c.PreferredBackupWindow,
		PreferredMaintenanceWindow:       c.PreferredMaintenanceWindow,
		AvailabilityZones:                c.AvailabilityZones,
		MultiAz:                          c.MultiAZ,
		VpcSecurityGroupIds:              c.VpcSecurityGroupIds,
		DbSubnetGroupName:                c.DBSubnetGroupName,
		DbClusterParameterGroupName:      c.DBClusterParameterGroupName,
		StorageEncrypted:                 c.StorageEncrypted,
		KmsKeyId:                         c.KmsKeyId,
		CopyTagsToSnapshot:               c.CopyTagsToSnapshot,
		DeletionProtection:               c.DeletionProtection,
		EnableCloudwatchLogsExports:      c.EnabledCloudwatchLogsExports,
		EnableIamDatabaseAuthentication:  c.IAMDatabaseAuthenticationEnabled,
		ReplicationSourceIdentifier:      c.ReplicationSourceIdentifier,
		GlobalClusterIdentifier:          c.GlobalClusterIdentifier,
		StorageType:                      c.StorageType,
		ServerlessV2ScalingConfiguration: scalingConfigToProto(c.ServerlessV2ScalingConfiguration),
		AccountId:                        c.AccountID,
		Region:                           c.Region,
		DbClusterArn:                     c.DBClusterArn,
		MasterUserPasswordHash:           c.MasterUserPasswordHash,
		AllocatedStorage:                 c.AllocatedStorage,
		DbClusterResourceId:              c.DbClusterResourceId,
		NetworkType:                      c.NetworkType,
		PercentProgress:                  c.PercentProgress,
		PendingModifiedValues:            clusterPendingModifiedValuesToProto(c.PendingModifiedValues),
		HostedZoneId:                     c.HostedZoneId,
		ReadReplicaIdentifiers:           c.ReadReplicaIdentifiers,
	}
	if c.ClusterCreateTime != nil {
		p.ClusterCreateTime = timestamppb.New(*c.ClusterCreateTime)
	}
	if c.EarliestRestorableTime != nil {
		p.EarliestRestorableTime = timestamppb.New(*c.EarliestRestorableTime)
	}
	if c.LatestRestorableTime != nil {
		p.LatestRestorableTime = timestamppb.New(*c.LatestRestorableTime)
	}
	for _, role := range c.AssociatedRoles {
		p.AssociatedRoles = append(p.AssociatedRoles, clusterRoleToProto(&role))
	}
	if c.Endpoint != nil {
		p.Endpoint = &pb.ClusterEndpoint{Address: c.Endpoint.Address, Port: int32(c.Endpoint.Port)}
	}
	for _, m := range c.DBClusterMembers {
		p.DbClusterMembers = append(p.DbClusterMembers, dbClusterMemberToProto(&m))
	}
	if c.ReaderEndpoint != nil {
		p.ReaderEndpoint = &pb.ClusterEndpoint{Address: c.ReaderEndpoint.Address, Port: int32(c.ReaderEndpoint.Port)}
	}
	if c.AutomaticRestartTime != nil {
		p.AutomaticRestartTime = timestamppb.New(*c.AutomaticRestartTime)
	}
	return p
}

func ProtoToCluster(p *pb.DBCluster) *DBCluster {
	if p == nil {
		return nil
	}
	c := &DBCluster{
		DBClusterIdentifier:              p.GetDbClusterIdentifier(),
		Engine:                           p.GetEngine(),
		EngineVersion:                    p.GetEngineVersion(),
		Status:                           p.GetStatus(),
		MasterUsername:                   p.GetMasterUsername(),
		DatabaseName:                     p.GetDatabaseName(),
		Port:                             int(p.GetPort()),
		BackupRetentionPeriod:            int(p.GetBackupRetentionPeriod()),
		PreferredBackupWindow:            p.GetPreferredBackupWindow(),
		PreferredMaintenanceWindow:       p.GetPreferredMaintenanceWindow(),
		AvailabilityZones:                p.GetAvailabilityZones(),
		MultiAZ:                          p.GetMultiAz(),
		VpcSecurityGroupIds:              p.GetVpcSecurityGroupIds(),
		DBSubnetGroupName:                p.GetDbSubnetGroupName(),
		DBClusterParameterGroupName:      p.GetDbClusterParameterGroupName(),
		StorageEncrypted:                 p.GetStorageEncrypted(),
		KmsKeyId:                         p.GetKmsKeyId(),
		CopyTagsToSnapshot:               p.GetCopyTagsToSnapshot(),
		DeletionProtection:               p.GetDeletionProtection(),
		EnabledCloudwatchLogsExports:     p.GetEnableCloudwatchLogsExports(),
		IAMDatabaseAuthenticationEnabled: p.GetEnableIamDatabaseAuthentication(),
		ReplicationSourceIdentifier:      p.GetReplicationSourceIdentifier(),
		GlobalClusterIdentifier:          p.GetGlobalClusterIdentifier(),
		StorageType:                      p.GetStorageType(),
		ServerlessV2ScalingConfiguration: protoToScalingConfig(p.GetServerlessV2ScalingConfiguration()),
		AccountID:                        p.GetAccountId(),
		Region:                           p.GetRegion(),
		DBClusterArn:                     p.GetDbClusterArn(),
		MasterUserPasswordHash:           p.GetMasterUserPasswordHash(),
		AllocatedStorage:                 p.GetAllocatedStorage(),
		DbClusterResourceId:              p.GetDbClusterResourceId(),
		NetworkType:                      p.GetNetworkType(),
		PercentProgress:                  p.GetPercentProgress(),
		PendingModifiedValues:            protoToClusterPendingModifiedValues(p.GetPendingModifiedValues()),
		HostedZoneId:                     p.GetHostedZoneId(),
		ReadReplicaIdentifiers:           p.GetReadReplicaIdentifiers(),
	}
	if p.ClusterCreateTime != nil {
		t := p.ClusterCreateTime.AsTime()
		c.ClusterCreateTime = &t
	}
	if p.EarliestRestorableTime != nil {
		t := p.EarliestRestorableTime.AsTime()
		c.EarliestRestorableTime = &t
	}
	if p.LatestRestorableTime != nil {
		t := p.LatestRestorableTime.AsTime()
		c.LatestRestorableTime = &t
	}
	for _, role := range p.GetAssociatedRoles() {
		c.AssociatedRoles = append(c.AssociatedRoles, *protoToClusterRole(role))
	}
	if ep := p.GetEndpoint(); ep != nil {
		c.Endpoint = &Endpoint{Address: ep.GetAddress(), Port: int(ep.GetPort())}
	}
	for _, m := range p.GetDbClusterMembers() {
		c.DBClusterMembers = append(c.DBClusterMembers, *protoToDBClusterMember(m))
	}
	if ep := p.GetReaderEndpoint(); ep != nil {
		c.ReaderEndpoint = &Endpoint{Address: ep.GetAddress(), Port: int(ep.GetPort())}
	}
	if p.AutomaticRestartTime != nil {
		t := p.AutomaticRestartTime.AsTime()
		c.AutomaticRestartTime = &t
	}
	return c
}

func InstanceToProto(i *DBInstance) *pb.DBInstance {
	if i == nil {
		return nil
	}
	p := &pb.DBInstance{
		DbInstanceIdentifier:            i.DBInstanceIdentifier,
		DbClusterIdentifier:             i.DBClusterIdentifier,
		Engine:                          i.Engine,
		EngineVersion:                   i.EngineVersion,
		DbInstanceClass:                 i.DBInstanceClass,
		Status:                          i.DBInstanceStatus,
		AvailabilityZone:                i.AvailabilityZone,
		PreferredMaintenanceWindow:      i.PreferredMaintenanceWindow,
		DbParameterGroupName:            i.DBParameterGroupName,
		DbSecurityGroups:                i.DBSecurityGroups,
		VpcSecurityGroupIds:             i.VpcSecurityGroupIds,
		DbSubnetGroupName:               i.DBSubnetGroupName,
		EnableCloudwatchLogsExports:     i.EnabledCloudwatchLogsExports,
		EnableIamDatabaseAuthentication: i.IAMDatabaseAuthenticationEnabled,
		PubliclyAccessible:              i.PubliclyAccessible,
		AutoMinorVersionUpgrade:         i.AutoMinorVersionUpgrade,
		CopyTagsToSnapshot:              i.CopyTagsToSnapshot,
		AccountId:                       i.AccountID,
		Region:                          i.Region,
		DbInstanceArn:                   i.DBInstanceArn,

		// AWS-standard DBInstance fields (RDS-5/RDS-20).
		AllocatedStorage:                   i.AllocatedStorage,
		MasterUsername:                     i.MasterUsername,
		StorageType:                        i.StorageType,
		BackupRetentionPeriod:              i.BackupRetentionPeriod,
		LicenseModel:                       i.LicenseModel,
		StorageEncrypted:                   i.StorageEncrypted,
		KmsKeyId:                           i.KmsKeyId,
		DeletionProtection:                 i.DeletionProtection,
		MultiAz:                            i.MultiAZ,
		SecondaryAvailabilityZone:          i.SecondaryAvailabilityZone,
		Port:                               i.Port,
		OptionGroupName:                    i.OptionGroupName,
		VpcId:                              i.VpcId,
		Iops:                               i.Iops,
		MaxAllocatedStorage:                i.MaxAllocatedStorage,
		StorageThroughput:                  i.StorageThroughput,
		MonitoringInterval:                 i.MonitoringInterval,
		EnhancedMonitoringResourceArn:      i.EnhancedMonitoringResourceArn,
		PerformanceInsightsEnabled:         i.PerformanceInsightsEnabled,
		PerformanceInsightsKmsKeyId:        i.PerformanceInsightsKMSKeyId,
		PerformanceInsightsRetentionPeriod: i.PerformanceInsightsRetentionPeriod,
		CaCertificateIdentifier:            i.CACertificateIdentifier,
		DbiResourceId:                      i.DbiResourceId,
		PreferredBackupWindow:              i.PreferredBackupWindow,

		// MasterUserPasswordHash: write-only, persisted for verification
		// (H1 fix). json:"-" on the Go struct prevents it from appearing
		// in API responses.
		MasterUserPasswordHash: i.MasterUserPasswordHash,
	}
	if i.InstanceCreateTime != nil {
		p.InstanceCreateTime = timestamppb.New(*i.InstanceCreateTime)
	}
	if i.LatestRestorableTime != nil {
		p.LatestRestorableTime = timestamppb.New(*i.LatestRestorableTime)
	}
	if i.Endpoint != nil {
		p.Endpoint = &pb.ClusterEndpoint{Address: i.Endpoint.Address, Port: int32(i.Endpoint.Port)}
	}
	return p
}

func ProtoToInstance(p *pb.DBInstance) *DBInstance {
	if p == nil {
		return nil
	}
	i := &DBInstance{
		DBInstanceIdentifier:             p.GetDbInstanceIdentifier(),
		DBClusterIdentifier:              p.GetDbClusterIdentifier(),
		Engine:                           p.GetEngine(),
		EngineVersion:                    p.GetEngineVersion(),
		DBInstanceClass:                  p.GetDbInstanceClass(),
		DBInstanceStatus:                 p.GetStatus(),
		AvailabilityZone:                 p.GetAvailabilityZone(),
		PreferredMaintenanceWindow:       p.GetPreferredMaintenanceWindow(),
		DBParameterGroupName:             p.GetDbParameterGroupName(),
		DBSecurityGroups:                 p.GetDbSecurityGroups(),
		VpcSecurityGroupIds:              p.GetVpcSecurityGroupIds(),
		DBSubnetGroupName:                p.GetDbSubnetGroupName(),
		EnabledCloudwatchLogsExports:     p.GetEnableCloudwatchLogsExports(),
		IAMDatabaseAuthenticationEnabled: p.GetEnableIamDatabaseAuthentication(),
		PubliclyAccessible:               p.GetPubliclyAccessible(),
		AutoMinorVersionUpgrade:          p.GetAutoMinorVersionUpgrade(),
		CopyTagsToSnapshot:               p.GetCopyTagsToSnapshot(),
		AccountID:                        p.GetAccountId(),
		Region:                           p.GetRegion(),
		DBInstanceArn:                    p.GetDbInstanceArn(),

		// AWS-standard DBInstance fields (RDS-5/RDS-20).
		AllocatedStorage:                   p.GetAllocatedStorage(),
		MasterUsername:                     p.GetMasterUsername(),
		StorageType:                        p.GetStorageType(),
		BackupRetentionPeriod:              p.GetBackupRetentionPeriod(),
		LicenseModel:                       p.GetLicenseModel(),
		StorageEncrypted:                   p.GetStorageEncrypted(),
		KmsKeyId:                           p.GetKmsKeyId(),
		DeletionProtection:                 p.GetDeletionProtection(),
		MultiAZ:                            p.GetMultiAz(),
		SecondaryAvailabilityZone:          p.GetSecondaryAvailabilityZone(),
		Port:                               p.GetPort(),
		OptionGroupName:                    p.GetOptionGroupName(),
		VpcId:                              p.GetVpcId(),
		Iops:                               p.GetIops(),
		MaxAllocatedStorage:                p.GetMaxAllocatedStorage(),
		StorageThroughput:                  p.GetStorageThroughput(),
		MonitoringInterval:                 p.GetMonitoringInterval(),
		EnhancedMonitoringResourceArn:      p.GetEnhancedMonitoringResourceArn(),
		PerformanceInsightsEnabled:         p.GetPerformanceInsightsEnabled(),
		PerformanceInsightsKMSKeyId:        p.GetPerformanceInsightsKmsKeyId(),
		PerformanceInsightsRetentionPeriod: p.GetPerformanceInsightsRetentionPeriod(),
		CACertificateIdentifier:            p.GetCaCertificateIdentifier(),
		DbiResourceId:                      p.GetDbiResourceId(),
		PreferredBackupWindow:              p.GetPreferredBackupWindow(),

		// MasterUserPasswordHash: write-only (H1 fix).
		MasterUserPasswordHash: p.GetMasterUserPasswordHash(),
	}
	if p.InstanceCreateTime != nil {
		t := p.InstanceCreateTime.AsTime()
		i.InstanceCreateTime = &t
	}
	if p.LatestRestorableTime != nil {
		t := p.LatestRestorableTime.AsTime()
		i.LatestRestorableTime = &t
	}
	if ep := p.GetEndpoint(); ep != nil {
		i.Endpoint = &Endpoint{Address: ep.GetAddress(), Port: int(ep.GetPort())}
	}
	return i
}

func SnapshotToProto(s *DBClusterSnapshot) *pb.DBClusterSnapshot {
	if s == nil {
		return nil
	}
	p := &pb.DBClusterSnapshot{
		DbClusterSnapshotIdentifier: s.DBClusterSnapshotIdentifier,
		DbClusterIdentifier:         s.DBClusterIdentifier,
		Engine:                      s.Engine,
		EngineVersion:               s.EngineVersion,
		SnapshotType:                s.SnapshotType,
		Status:                      s.Status,
		Port:                        int32(s.Port),
		VpcId:                       s.VpcId,
		StorageEncrypted:            s.StorageEncrypted,
		KmsKeyId:                    s.KmsKeyId,
		DbSnapshotArn:               s.DBSnapshotArn,
		AccountId:                   s.AccountID,
		Region:                      s.Region,
		RestoreAttributeValues:      s.RestoreAttributeValues,

		// AWS-standard fields captured at snapshot time (RDS-2).
		MasterUsername:                   s.MasterUsername,
		AllocatedStorage:                 s.AllocatedStorage,
		StorageType:                      s.StorageType,
		LicenseModel:                     s.LicenseModel,
		IamDatabaseAuthenticationEnabled: s.IAMDatabaseAuthenticationEnabled,
	}
	if s.SnapshotCreateTime != nil {
		p.SnapshotCreateTime = timestamppb.New(*s.SnapshotCreateTime)
	}
	if s.ClusterCreateTime != nil {
		p.ClusterCreateTime = timestamppb.New(*s.ClusterCreateTime)
	}
	return p
}

func ProtoToSnapshot(p *pb.DBClusterSnapshot) *DBClusterSnapshot {
	if p == nil {
		return nil
	}
	s := &DBClusterSnapshot{
		DBClusterSnapshotIdentifier: p.GetDbClusterSnapshotIdentifier(),
		DBClusterIdentifier:         p.GetDbClusterIdentifier(),
		Engine:                      p.GetEngine(),
		EngineVersion:               p.GetEngineVersion(),
		SnapshotType:                p.GetSnapshotType(),
		Status:                      p.GetStatus(),
		Port:                        int(p.GetPort()),
		VpcId:                       p.GetVpcId(),
		StorageEncrypted:            p.GetStorageEncrypted(),
		KmsKeyId:                    p.GetKmsKeyId(),
		DBSnapshotArn:               p.GetDbSnapshotArn(),
		AccountID:                   p.GetAccountId(),
		Region:                      p.GetRegion(),
		RestoreAttributeValues:      p.GetRestoreAttributeValues(),

		MasterUsername:                   p.GetMasterUsername(),
		AllocatedStorage:                 p.GetAllocatedStorage(),
		StorageType:                      p.GetStorageType(),
		LicenseModel:                     p.GetLicenseModel(),
		IAMDatabaseAuthenticationEnabled: p.GetIamDatabaseAuthenticationEnabled(),
	}
	if p.SnapshotCreateTime != nil {
		t := p.SnapshotCreateTime.AsTime()
		s.SnapshotCreateTime = &t
	}
	if p.ClusterCreateTime != nil {
		t := p.ClusterCreateTime.AsTime()
		s.ClusterCreateTime = &t
	}
	return s
}

func ClusterParameterGroupToProto(pg *DBClusterParameterGroup) *pb.DBClusterParameterGroup {
	if pg == nil {
		return nil
	}
	p := &pb.DBClusterParameterGroup{
		DbClusterParameterGroupName: pg.DBClusterParameterGroupName,
		DbParameterGroupFamily:      pg.DBParameterGroupFamily,
		Description:                 pg.Description,
		Arn:                         pg.ARN,
	}
	for _, param := range pg.Parameters {
		p.Parameters = append(p.Parameters, ParameterToProto(&param))
	}
	return p
}

func ProtoToClusterParameterGroup(p *pb.DBClusterParameterGroup) *DBClusterParameterGroup {
	if p == nil {
		return nil
	}
	pg := &DBClusterParameterGroup{
		DBClusterParameterGroupName: p.GetDbClusterParameterGroupName(),
		DBParameterGroupFamily:      p.GetDbParameterGroupFamily(),
		Description:                 p.GetDescription(),
		ARN:                         p.GetArn(),
	}
	for _, param := range p.Parameters {
		pg.Parameters = append(pg.Parameters, *ProtoToParameter(param))
	}
	return pg
}

func ParameterGroupToProto(pg *DBParameterGroup) *pb.DBParameterGroup {
	if pg == nil {
		return nil
	}
	p := &pb.DBParameterGroup{
		DbParameterGroupName:   pg.DBParameterGroupName,
		DbParameterGroupFamily: pg.DBParameterGroupFamily,
		Description:            pg.Description,
		Arn:                    pg.ARN,
	}
	for _, param := range pg.Parameters {
		p.Parameters = append(p.Parameters, ParameterToProto(&param))
	}
	return p
}

func ProtoToParameterGroup(p *pb.DBParameterGroup) *DBParameterGroup {
	if p == nil {
		return nil
	}
	pg := &DBParameterGroup{
		DBParameterGroupName:   p.GetDbParameterGroupName(),
		DBParameterGroupFamily: p.GetDbParameterGroupFamily(),
		Description:            p.GetDescription(),
		ARN:                    p.GetArn(),
	}
	for _, param := range p.Parameters {
		pg.Parameters = append(pg.Parameters, *ProtoToParameter(param))
	}
	return pg
}

func subnetToProto(s *Subnet) *pb.Subnet {
	if s == nil {
		return nil
	}
	return &pb.Subnet{
		SubnetIdentifier:       s.SubnetIdentifier,
		SubnetAvailabilityZone: s.SubnetAvailabilityZone,
		SubnetOutpost:          s.SubnetOutpost,
		SubnetStatus:           s.SubnetStatus,
	}
}

func protoToSubnet(s *pb.Subnet) *Subnet {
	if s == nil {
		return nil
	}
	return &Subnet{
		SubnetIdentifier:       s.GetSubnetIdentifier(),
		SubnetAvailabilityZone: s.GetSubnetAvailabilityZone(),
		SubnetOutpost:          s.GetSubnetOutpost(),
		SubnetStatus:           s.GetSubnetStatus(),
	}
}

func SubnetGroupToProto(sg *DBSubnetGroup) *pb.DBSubnetGroup {
	if sg == nil {
		return nil
	}
	p := &pb.DBSubnetGroup{
		DbSubnetGroupName:        sg.DBSubnetGroupName,
		DbSubnetGroupDescription: sg.DBSubnetGroupDescription,
		VpcId:                    sg.VpcId,
		SubnetGroupStatus:        sg.SubnetGroupStatus,
		Arn:                      sg.ARN,
	}
	for i := range sg.Subnets {
		p.Subnets = append(p.Subnets, subnetToProto(&sg.Subnets[i]))
	}
	return p
}

func ProtoToSubnetGroup(p *pb.DBSubnetGroup) *DBSubnetGroup {
	if p == nil {
		return nil
	}
	sg := &DBSubnetGroup{
		DBSubnetGroupName:        p.GetDbSubnetGroupName(),
		DBSubnetGroupDescription: p.GetDbSubnetGroupDescription(),
		VpcId:                    p.GetVpcId(),
		SubnetGroupStatus:        p.GetSubnetGroupStatus(),
		ARN:                      p.GetArn(),
	}
	for _, s := range p.GetSubnets() {
		sg.Subnets = append(sg.Subnets, *protoToSubnet(s))
	}
	return sg
}

func globalClusterMemberToProto(m *GlobalClusterMember) *pb.GlobalClusterMember {
	if m == nil {
		return nil
	}
	return &pb.GlobalClusterMember{
		DbClusterArn:            m.DBClusterArn,
		Readers:                 m.Readers,
		IsWriter:                m.IsWriter,
		GlobalClusterIdentifier: m.GlobalClusterIdentifier,
	}
}

func protoToGlobalClusterMember(m *pb.GlobalClusterMember) *GlobalClusterMember {
	if m == nil {
		return nil
	}
	return &GlobalClusterMember{
		DBClusterArn:            m.GetDbClusterArn(),
		Readers:                 m.GetReaders(),
		IsWriter:                m.GetIsWriter(),
		GlobalClusterIdentifier: m.GetGlobalClusterIdentifier(),
	}
}

func GlobalClusterToProto(gc *GlobalCluster) *pb.GlobalCluster {
	if gc == nil {
		return nil
	}
	p := &pb.GlobalCluster{
		GlobalClusterIdentifier: gc.GlobalClusterIdentifier,
		GlobalClusterResourceId: gc.GlobalClusterResourceId,
		GlobalClusterArn:        gc.GlobalClusterArn,
		Engine:                  gc.Engine,
		EngineVersion:           gc.EngineVersion,
		Status:                  gc.Status,
		StorageEncrypted:        gc.StorageEncrypted,
		DeletionProtection:      gc.DeletionProtection,
		AccountId:               gc.AccountID,
		Region:                  gc.Region,
	}
	for i := range gc.GlobalClusterMembers {
		p.GlobalClusterMembers = append(p.GlobalClusterMembers, globalClusterMemberToProto(&gc.GlobalClusterMembers[i]))
	}
	return p
}

func ProtoToGlobalCluster(p *pb.GlobalCluster) *GlobalCluster {
	if p == nil {
		return nil
	}
	gc := &GlobalCluster{
		GlobalClusterIdentifier: p.GetGlobalClusterIdentifier(),
		GlobalClusterResourceId: p.GetGlobalClusterResourceId(),
		GlobalClusterArn:        p.GetGlobalClusterArn(),
		Engine:                  p.GetEngine(),
		EngineVersion:           p.GetEngineVersion(),
		Status:                  p.GetStatus(),
		StorageEncrypted:        p.GetStorageEncrypted(),
		DeletionProtection:      p.GetDeletionProtection(),
		AccountID:               p.GetAccountId(),
		Region:                  p.GetRegion(),
	}
	for _, m := range p.GetGlobalClusterMembers() {
		gc.GlobalClusterMembers = append(gc.GlobalClusterMembers, *protoToGlobalClusterMember(m))
	}
	return gc
}

func EventSubscriptionToProto(sub *EventSubscription) *pb.EventSubscription {
	if sub == nil {
		return nil
	}
	p := &pb.EventSubscription{
		CustSubscriptionId:  sub.CustSubscriptionId,
		SnsTopicArn:         sub.SnsTopicArn,
		Status:              sub.Status,
		SourceType:          sub.SourceType,
		SourceIdsList:       sub.SourceIdsList,
		EventCategoriesList: sub.EventCategoriesList,
		Enabled:             sub.Enabled,
		CustSubscriptionArn: sub.CustSubscriptionArn,
	}
	if sub.SubscriptionCreationTime != nil {
		p.SubscriptionCreationTime = timestamppb.New(*sub.SubscriptionCreationTime)
	}
	return p
}

func ProtoToEventSubscription(p *pb.EventSubscription) *EventSubscription {
	if p == nil {
		return nil
	}
	sub := &EventSubscription{
		CustSubscriptionId:  p.GetCustSubscriptionId(),
		SnsTopicArn:         p.GetSnsTopicArn(),
		Status:              p.GetStatus(),
		SourceType:          p.GetSourceType(),
		SourceIdsList:       p.GetSourceIdsList(),
		EventCategoriesList: p.GetEventCategoriesList(),
		Enabled:             p.GetEnabled(),
		CustSubscriptionArn: p.GetCustSubscriptionArn(),
	}
	if p.SubscriptionCreationTime != nil {
		t := p.SubscriptionCreationTime.AsTime()
		sub.SubscriptionCreationTime = &t
	}
	return sub
}

func ParameterToProto(p *Parameter) *pb.Parameter {
	if p == nil {
		return nil
	}
	return &pb.Parameter{
		ParameterName:        p.ParameterName,
		ParameterValue:       p.ParameterValue,
		Description:          p.Description,
		Source:               p.Source,
		ApplyType:            p.ApplyType,
		DataType:             p.DataType,
		AllowedValues:        p.AllowedValues,
		IsModifiable:         p.IsModifiable,
		MinimumEngineVersion: p.MinimumEngineVersion,
		ApplyMethod:          p.ApplyMethod,
	}
}

func ProtoToParameter(p *pb.Parameter) *Parameter {
	if p == nil {
		return nil
	}
	return &Parameter{
		ParameterName:        p.GetParameterName(),
		ParameterValue:       p.GetParameterValue(),
		Description:          p.GetDescription(),
		Source:               p.GetSource(),
		ApplyType:            p.GetApplyType(),
		DataType:             p.GetDataType(),
		AllowedValues:        p.GetAllowedValues(),
		IsModifiable:         p.GetIsModifiable(),
		MinimumEngineVersion: p.GetMinimumEngineVersion(),
		ApplyMethod:          p.GetApplyMethod(),
	}
}

func TagToProto(t *types.Tag) *pb.Tag {
	if t == nil {
		return nil
	}
	return &pb.Tag{
		Key:   t.Key,
		Value: t.Value,
	}
}

func ProtoToTag(t *pb.Tag) *types.Tag {
	if t == nil {
		return nil
	}
	return &types.Tag{
		Key:   t.GetKey(),
		Value: t.GetValue(),
	}
}

func TagsToProto(tags []types.Tag) []*pb.Tag {
	if tags == nil {
		return nil
	}
	result := make([]*pb.Tag, 0, len(tags))
	for i := range tags {
		result = append(result, TagToProto(&tags[i]))
	}
	return result
}

func ProtoToTags(tags []*pb.Tag) []types.Tag {
	if tags == nil {
		return nil
	}
	result := make([]types.Tag, 0, len(tags))
	for _, t := range tags {
		result = append(result, *ProtoToTag(t))
	}
	return result
}

func EventToProto(e *Event) *pb.Event {
	if e == nil {
		return nil
	}
	p := &pb.Event{
		EventId:          e.EventID,
		EventCategories:  e.EventCategories,
		Message:          e.Message,
		SourceArn:        e.SourceArn,
		SourceIdentifier: e.SourceIdentifier,
		SourceType:       e.SourceType,
	}
	if !e.Date.IsZero() {
		p.Date = timestamppb.New(e.Date)
	}
	return p
}

func ProtoToEvent(p *pb.Event) *Event {
	if p == nil {
		return nil
	}
	e := &Event{
		EventID:          p.GetEventId(),
		EventCategories:  p.GetEventCategories(),
		Message:          p.GetMessage(),
		SourceArn:        p.GetSourceArn(),
		SourceIdentifier: p.GetSourceIdentifier(),
		SourceType:       p.GetSourceType(),
	}
	if p.Date != nil {
		e.Date = p.Date.AsTime()
	}
	return e
}
