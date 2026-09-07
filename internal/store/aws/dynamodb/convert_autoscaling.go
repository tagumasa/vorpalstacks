package dynamodb

import (
	pb "vorpalstacks/internal/pb/storage/storage_dynamodb"
)

func boolPtrToProto(v *bool) *bool {
	if v == nil {
		return nil
	}
	value := *v
	return &value
}

func targetTrackingToProto(t *TargetTrackingScalingPolicyConfiguration) *pb.TargetTrackingScalingPolicyConfiguration {
	if t == nil {
		return nil
	}
	return &pb.TargetTrackingScalingPolicyConfiguration{
		DisableScaleIn:   boolPtrToProto(t.DisableScaleIn),
		ScaleInCooldown:  t.ScaleInCooldown,
		ScaleOutCooldown: t.ScaleOutCooldown,
		TargetValue:      t.TargetValue,
	}
}

func protoToTargetTracking(p *pb.TargetTrackingScalingPolicyConfiguration) *TargetTrackingScalingPolicyConfiguration {
	if p == nil {
		return nil
	}
	return &TargetTrackingScalingPolicyConfiguration{
		DisableScaleIn:   boolPtrToProto(p.DisableScaleIn),
		ScaleInCooldown:  p.ScaleInCooldown,
		ScaleOutCooldown: p.ScaleOutCooldown,
		TargetValue:      p.TargetValue,
	}
}

func autoScalingPolicyToProto(p AutoScalingPolicyDescription) *pb.AutoScalingPolicyDescription {
	return &pb.AutoScalingPolicyDescription{
		PolicyName:                               p.PolicyName,
		TargetTrackingScalingPolicyConfiguration: targetTrackingToProto(p.TargetTrackingScalingPolicyConfiguration),
	}
}

func protoToAutoScalingPolicy(p *pb.AutoScalingPolicyDescription) AutoScalingPolicyDescription {
	if p == nil {
		return AutoScalingPolicyDescription{}
	}
	return AutoScalingPolicyDescription{
		PolicyName:                               p.PolicyName,
		TargetTrackingScalingPolicyConfiguration: protoToTargetTracking(p.TargetTrackingScalingPolicyConfiguration),
	}
}

// autoScalingSettingsToProto converts one capacity dimension's auto-scaling
// description to its protobuf form.
func autoScalingSettingsToProto(s *AutoScalingSettingsDescription) *pb.AutoScalingSettingsDescription {
	if s == nil {
		return nil
	}
	policies := make([]*pb.AutoScalingPolicyDescription, len(s.ScalingPolicies))
	for i, policy := range s.ScalingPolicies {
		policies[i] = autoScalingPolicyToProto(policy)
	}
	return &pb.AutoScalingSettingsDescription{
		MinimumUnits:        s.MinimumUnits,
		MaximumUnits:        s.MaximumUnits,
		AutoScalingDisabled: boolPtrToProto(s.AutoScalingDisabled),
		AutoScalingRoleArn:  s.AutoScalingRoleArn,
		ScalingPolicies:     policies,
	}
}

func protoToAutoScalingSettings(p *pb.AutoScalingSettingsDescription) *AutoScalingSettingsDescription {
	if p == nil {
		return nil
	}
	policies := make([]AutoScalingPolicyDescription, len(p.ScalingPolicies))
	for i, policy := range p.ScalingPolicies {
		policies[i] = protoToAutoScalingPolicy(policy)
	}
	return &AutoScalingSettingsDescription{
		MinimumUnits:        p.MinimumUnits,
		MaximumUnits:        p.MaximumUnits,
		AutoScalingDisabled: boolPtrToProto(p.AutoScalingDisabled),
		AutoScalingRoleArn:  p.AutoScalingRoleArn,
		ScalingPolicies:     policies,
	}
}

// indexAutoScalingSettingsToProto converts one index's settings to protobuf.
func indexAutoScalingSettingsToProto(s IndexAutoScalingSettings) *pb.IndexAutoScalingSettings {
	return &pb.IndexAutoScalingSettings{
		IndexName:                     s.IndexName,
		ProvisionedReadCapacityUnits:  s.ProvisionedReadCapacityUnits,
		ProvisionedWriteCapacityUnits: s.ProvisionedWriteCapacityUnits,
		Read:                          autoScalingSettingsToProto(s.Read),
		Write:                         autoScalingSettingsToProto(s.Write),
	}
}

func protoToIndexAutoScalingSettings(p *pb.IndexAutoScalingSettings) IndexAutoScalingSettings {
	if p == nil {
		return IndexAutoScalingSettings{}
	}
	return IndexAutoScalingSettings{
		IndexName:                     p.IndexName,
		ProvisionedReadCapacityUnits:  p.ProvisionedReadCapacityUnits,
		ProvisionedWriteCapacityUnits: p.ProvisionedWriteCapacityUnits,
		Read:                          protoToAutoScalingSettings(p.Read),
		Write:                         protoToAutoScalingSettings(p.Write),
	}
}

// replicaAutoScalingToProto converts one replica description to protobuf.
func replicaAutoScalingToProto(r ReplicaAutoScalingDescription) *pb.ReplicaAutoScalingDescription {
	indexes := make([]*pb.IndexAutoScalingSettings, len(r.GlobalSecondaryIndexes))
	for i, index := range r.GlobalSecondaryIndexes {
		indexes[i] = indexAutoScalingSettingsToProto(index)
	}
	return &pb.ReplicaAutoScalingDescription{
		RegionName:             r.RegionName,
		Read:                   autoScalingSettingsToProto(r.Read),
		Write:                  autoScalingSettingsToProto(r.Write),
		GlobalSecondaryIndexes: indexes,
	}
}

func protoToReplicaAutoScaling(p *pb.ReplicaAutoScalingDescription) ReplicaAutoScalingDescription {
	if p == nil {
		return ReplicaAutoScalingDescription{}
	}
	indexes := make([]IndexAutoScalingSettings, len(p.GlobalSecondaryIndexes))
	for i, index := range p.GlobalSecondaryIndexes {
		indexes[i] = protoToIndexAutoScalingSettings(index)
	}
	return ReplicaAutoScalingDescription{
		RegionName:             p.RegionName,
		Read:                   protoToAutoScalingSettings(p.Read),
		Write:                  protoToAutoScalingSettings(p.Write),
		GlobalSecondaryIndexes: indexes,
	}
}

// tableReplicaAutoScalingToProto converts the table-level auto-scaling
// record to its protobuf form.
func tableReplicaAutoScalingToProto(s *TableReplicaAutoScalingSettings) *pb.TableReplicaAutoScalingSettings {
	if s == nil {
		return nil
	}
	replicas := make([]*pb.ReplicaAutoScalingDescription, len(s.Replicas))
	for i, replica := range s.Replicas {
		replicas[i] = replicaAutoScalingToProto(replica)
	}
	return &pb.TableReplicaAutoScalingSettings{Replicas: replicas}
}

func protoToTableReplicaAutoScaling(p *pb.TableReplicaAutoScalingSettings) *TableReplicaAutoScalingSettings {
	if p == nil {
		return nil
	}
	replicas := make([]ReplicaAutoScalingDescription, len(p.Replicas))
	for i, replica := range p.Replicas {
		replicas[i] = protoToReplicaAutoScaling(replica)
	}
	return &TableReplicaAutoScalingSettings{Replicas: replicas}
}
