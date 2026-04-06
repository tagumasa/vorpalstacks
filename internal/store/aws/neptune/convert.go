package neptune

import (
	pb "vorpalstacks/internal/pb/storage/storage_neptune"

	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

// clusterRoleToProto converts a DBClusterRole to its protobuf representation.
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

// protoToClusterRole converts a protobuf DBClusterRole to its domain representation.
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

// scalingConfigToProto converts a ServerlessV2ScalingConfiguration to protobuf.
func scalingConfigToProto(c *ServerlessV2ScalingConfiguration) *pb.ServerlessV2ScalingConfiguration {
	if c == nil {
		return nil
	}
	return &pb.ServerlessV2ScalingConfiguration{
		MinCapacity: c.MinCapacity,
		MaxCapacity: c.MaxCapacity,
	}
}

// protoToScalingConfig converts a protobuf ServerlessV2ScalingConfiguration to its domain representation.
func protoToScalingConfig(c *pb.ServerlessV2ScalingConfiguration) *ServerlessV2ScalingConfiguration {
	if c == nil {
		return nil
	}
	return &ServerlessV2ScalingConfiguration{
		MinCapacity: c.GetMinCapacity(),
		MaxCapacity: c.GetMaxCapacity(),
	}
}

// ClusterToProto converts a DBCluster to its protobuf representation.
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
		EnableCloudwatchLogsExports:      c.EnableCloudwatchLogsExports,
		EnableIamDatabaseAuthentication:  c.EnableIAMDatabaseAuthentication,
		ReplicationSourceIdentifier:      c.ReplicationSourceIdentifier,
		GlobalClusterIdentifier:          c.GlobalClusterIdentifier,
		StorageType:                      c.StorageType,
		ServerlessV2ScalingConfiguration: scalingConfigToProto(c.ServerlessV2ScalingConfiguration),
		AccountId:                        c.AccountID,
		Region:                           c.Region,
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
	return p
}

// ProtoToCluster converts a protobuf DBCluster to its domain representation.
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
		EnableCloudwatchLogsExports:      p.GetEnableCloudwatchLogsExports(),
		EnableIAMDatabaseAuthentication:  p.GetEnableIamDatabaseAuthentication(),
		ReplicationSourceIdentifier:      p.GetReplicationSourceIdentifier(),
		GlobalClusterIdentifier:          p.GetGlobalClusterIdentifier(),
		StorageType:                      p.GetStorageType(),
		ServerlessV2ScalingConfiguration: protoToScalingConfig(p.GetServerlessV2ScalingConfiguration()),
		AccountID:                        p.GetAccountId(),
		Region:                           p.GetRegion(),
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
	return c
}

// InstanceToProto converts a DBInstance to its protobuf representation.
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
		Status:                          i.Status,
		AvailabilityZone:                i.AvailabilityZone,
		PreferredMaintenanceWindow:      i.PreferredMaintenanceWindow,
		DbParameterGroupName:            i.DBParameterGroupName,
		DbSecurityGroups:                i.DBSecurityGroups,
		VpcSecurityGroupIds:             i.VpcSecurityGroupIds,
		DbSubnetGroupName:               i.DBSubnetGroupName,
		EnableCloudwatchLogsExports:     i.EnableCloudwatchLogsExports,
		EnableIamDatabaseAuthentication: i.EnableIAMDatabaseAuthentication,
		PubliclyAccessible:              i.PubliclyAccessible,
		AutoMinorVersionUpgrade:         i.AutoMinorVersionUpgrade,
		CopyTagsToSnapshot:              i.CopyTagsToSnapshot,
		AccountId:                       i.AccountID,
		Region:                          i.Region,
	}
	if i.InstanceCreateTime != nil {
		p.InstanceCreateTime = timestamppb.New(*i.InstanceCreateTime)
	}
	return p
}

// ProtoToInstance converts a protobuf DBInstance to its domain representation.
func ProtoToInstance(p *pb.DBInstance) *DBInstance {
	if p == nil {
		return nil
	}
	i := &DBInstance{
		DBInstanceIdentifier:            p.GetDbInstanceIdentifier(),
		DBClusterIdentifier:             p.GetDbClusterIdentifier(),
		Engine:                          p.GetEngine(),
		EngineVersion:                   p.GetEngineVersion(),
		DBInstanceClass:                 p.GetDbInstanceClass(),
		Status:                          p.GetStatus(),
		AvailabilityZone:                p.GetAvailabilityZone(),
		PreferredMaintenanceWindow:      p.GetPreferredMaintenanceWindow(),
		DBParameterGroupName:            p.GetDbParameterGroupName(),
		DBSecurityGroups:                p.GetDbSecurityGroups(),
		VpcSecurityGroupIds:             p.GetVpcSecurityGroupIds(),
		DBSubnetGroupName:               p.GetDbSubnetGroupName(),
		EnableCloudwatchLogsExports:     p.GetEnableCloudwatchLogsExports(),
		EnableIAMDatabaseAuthentication: p.GetEnableIamDatabaseAuthentication(),
		PubliclyAccessible:              p.GetPubliclyAccessible(),
		AutoMinorVersionUpgrade:         p.GetAutoMinorVersionUpgrade(),
		CopyTagsToSnapshot:              p.GetCopyTagsToSnapshot(),
		AccountID:                       p.GetAccountId(),
		Region:                          p.GetRegion(),
	}
	if p.InstanceCreateTime != nil {
		t := p.InstanceCreateTime.AsTime()
		i.InstanceCreateTime = &t
	}
	return i
}

// SnapshotToProto converts a DBClusterSnapshot to its protobuf representation.
func SnapshotToProto(s *DBClusterSnapshot) *pb.DBClusterSnapshot {
	if s == nil {
		return nil
	}
	p := &pb.DBClusterSnapshot{
		DbClusterSnapshotIdentifier: s.DBClusterSnapshotIdentifier,
		DbClusterIdentifier:         s.DBClusterIdentifier,
		Engine:                      s.Engine,
		EngineVersion:               s.EngineVersion,
		Status:                      s.Status,
		Port:                        int32(s.Port),
		VpcId:                       s.VpcId,
		StorageEncrypted:            s.StorageEncrypted,
		KmsKeyId:                    s.KmsKeyId,
		DbSnapshotArn:               s.DBSnapshotArn,
		AccountId:                   s.AccountID,
		Region:                      s.Region,
		RestoreAttributeValues:      s.RestoreAttributeValues,
	}
	if s.SnapshotCreateTime != nil {
		p.SnapshotCreateTime = timestamppb.New(*s.SnapshotCreateTime)
	}
	if s.ClusterCreateTime != nil {
		p.ClusterCreateTime = timestamppb.New(*s.ClusterCreateTime)
	}
	return p
}

// ProtoToSnapshot converts a protobuf DBClusterSnapshot to its domain representation.
func ProtoToSnapshot(p *pb.DBClusterSnapshot) *DBClusterSnapshot {
	if p == nil {
		return nil
	}
	s := &DBClusterSnapshot{
		DBClusterSnapshotIdentifier: p.GetDbClusterSnapshotIdentifier(),
		DBClusterIdentifier:         p.GetDbClusterIdentifier(),
		Engine:                      p.GetEngine(),
		EngineVersion:               p.GetEngineVersion(),
		Status:                      p.GetStatus(),
		Port:                        int(p.GetPort()),
		VpcId:                       p.GetVpcId(),
		StorageEncrypted:            p.GetStorageEncrypted(),
		KmsKeyId:                    p.GetKmsKeyId(),
		DBSnapshotArn:               p.GetDbSnapshotArn(),
		AccountID:                   p.GetAccountId(),
		Region:                      p.GetRegion(),
		RestoreAttributeValues:      p.GetRestoreAttributeValues(),
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

// ClusterParameterGroupToProto converts a DBClusterParameterGroup to protobuf.
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

// ProtoToClusterParameterGroup converts a protobuf DBClusterParameterGroup to its domain representation.
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

// ParameterGroupToProto converts a DBParameterGroup to protobuf.
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

// ProtoToParameterGroup converts a protobuf DBParameterGroup to its domain representation.
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

// subnetToProto converts a Subnet to its protobuf representation.
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

// protoToSubnet converts a protobuf Subnet to its domain representation.
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

// SubnetGroupToProto converts a DBSubnetGroup to its protobuf representation.
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

// ProtoToSubnetGroup converts a protobuf DBSubnetGroup to its domain representation.
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

// globalClusterMemberToProto converts a GlobalClusterMember to protobuf.
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

// protoToGlobalClusterMember converts a protobuf GlobalClusterMember to its domain representation.
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

// GlobalClusterToProto converts a GlobalCluster to its protobuf representation.
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

// ProtoToGlobalCluster converts a protobuf GlobalCluster to its domain representation.
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

// EventSubscriptionToProto converts an EventSubscription to its protobuf representation.
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

// ProtoToEventSubscription converts a protobuf EventSubscription to its domain representation.
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

// ParameterToProto converts a Parameter to its protobuf representation.
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

// ProtoToParameter converts a protobuf Parameter to its domain representation.
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

// TagToProto converts a Tag to its protobuf representation.
func TagToProto(t *Tag) *pb.Tag {
	if t == nil {
		return nil
	}
	return &pb.Tag{
		Key:   t.Key,
		Value: t.Value,
	}
}

// ProtoToTag converts a protobuf Tag to its domain representation.
func ProtoToTag(t *pb.Tag) *Tag {
	if t == nil {
		return nil
	}
	return &Tag{
		Key:   t.GetKey(),
		Value: t.GetValue(),
	}
}

// TagsToProto converts a slice of Tags to a slice of protobuf Tags.
func TagsToProto(tags []Tag) []*pb.Tag {
	if tags == nil {
		return nil
	}
	result := make([]*pb.Tag, 0, len(tags))
	for i := range tags {
		result = append(result, TagToProto(&tags[i]))
	}
	return result
}

// ProtoToTags converts a slice of protobuf Tags to a slice of domain Tags.
func ProtoToTags(tags []*pb.Tag) []Tag {
	if tags == nil {
		return nil
	}
	result := make([]Tag, 0, len(tags))
	for _, t := range tags {
		result = append(result, *ProtoToTag(t))
	}
	return result
}

// ClusterEndpointToProto converts a DBClusterEndpoint to its protobuf representation.
func ClusterEndpointToProto(ep *DBClusterEndpoint) *pb.DBClusterEndpoint {
	if ep == nil {
		return nil
	}
	return &pb.DBClusterEndpoint{
		DbClusterEndpointIdentifier: ep.DBClusterEndpointIdentifier,
		DbClusterIdentifier:         ep.DBClusterIdentifier,
		Endpoint:                    ep.Endpoint,
		Status:                      ep.Status,
		EndpointType:                ep.EndpointType,
		ExcludedMembers:             ep.ExcludedMembers,
		StaticMembers:               ep.StaticMembers,
		DbClusterEndpointArn:        ep.DBClusterEndpointArn,
	}
}

// ProtoToClusterEndpoint converts a protobuf DBClusterEndpoint to its domain representation.
func ProtoToClusterEndpoint(p *pb.DBClusterEndpoint) *DBClusterEndpoint {
	if p == nil {
		return nil
	}
	return &DBClusterEndpoint{
		DBClusterEndpointIdentifier: p.GetDbClusterEndpointIdentifier(),
		DBClusterIdentifier:         p.GetDbClusterIdentifier(),
		Endpoint:                    p.GetEndpoint(),
		Status:                      p.GetStatus(),
		EndpointType:                p.GetEndpointType(),
		ExcludedMembers:             p.GetExcludedMembers(),
		StaticMembers:               p.GetStaticMembers(),
		DBClusterEndpointArn:        p.GetDbClusterEndpointArn(),
	}
}
