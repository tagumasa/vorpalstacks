package dynamodb

import (
	"time"

	pb "vorpalstacks/internal/pb/storage/storage_dynamodb"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// GlobalTable conversion

// GlobalTableToProto converts a GlobalTable to its protobuf representation.
func GlobalTableToProto(g *GlobalTable) *pb.GlobalTable {
	if g == nil {
		return nil
	}
	gsiWrite := make([]*pb.IndexAutoScalingSettings, len(g.GlobalSecondaryIndexWriteSettings))
	for i, index := range g.GlobalSecondaryIndexWriteSettings {
		gsiWrite[i] = indexAutoScalingSettingsToProto(index)
	}
	return &pb.GlobalTable{
		GlobalTableName:                   g.GlobalTableName,
		GlobalTableArn:                    g.GlobalTableArn,
		GlobalTableStatus:                 globalTableStatusToProto(g.GlobalTableStatus),
		CreationDateTime:                  timestamppb.New(g.CreationDateTime),
		ReplicationGroup:                  replicasToProto(g.ReplicationGroup),
		WriteAutoScalingSettings:          autoScalingSettingsToProto(g.WriteAutoScalingSettings),
		GlobalSecondaryIndexWriteSettings: gsiWrite,
	}
}

// ProtoToGlobalTable converts a protobuf GlobalTable to its internal representation.
func ProtoToGlobalTable(p *pb.GlobalTable) *GlobalTable {
	if p == nil {
		return nil
	}
	gsiWrite := make([]IndexAutoScalingSettings, len(p.GlobalSecondaryIndexWriteSettings))
	for i, index := range p.GlobalSecondaryIndexWriteSettings {
		gsiWrite[i] = protoToIndexAutoScalingSettings(index)
	}
	return &GlobalTable{
		GlobalTableName:                   p.GlobalTableName,
		GlobalTableArn:                    p.GlobalTableArn,
		GlobalTableStatus:                 protoToGlobalTableStatus(p.GlobalTableStatus),
		CreationDateTime:                  p.CreationDateTime.AsTime(),
		ReplicationGroup:                  protoToReplicas(p.ReplicationGroup),
		WriteAutoScalingSettings:          protoToAutoScalingSettings(p.WriteAutoScalingSettings),
		GlobalSecondaryIndexWriteSettings: gsiWrite,
	}
}

func replicasToProto(r []*Replica) []*pb.Replica {
	if r == nil {
		return nil
	}
	result := make([]*pb.Replica, len(r))
	for i, replica := range r {
		gsiRead := make([]*pb.IndexAutoScalingSettings, len(replica.GlobalSecondaryIndexReadSettings))
		for j, index := range replica.GlobalSecondaryIndexReadSettings {
			gsiRead[j] = indexAutoScalingSettingsToProto(index)
		}
		result[i] = &pb.Replica{
			RegionName:                       replica.RegionName,
			ReplicaStatus:                    replicaStatusToProto(replica.ReplicaStatus),
			BillingMode:                      billingModeToProto(BillingMode(replica.BillingMode)),
			ProvisionedReadCapacityUnits:     replica.ProvisionedReadCapacityUnits,
			ProvisionedWriteCapacityUnits:    replica.ProvisionedWriteCapacityUnits,
			ReadAutoScalingSettings:          autoScalingSettingsToProto(replica.ReadAutoScalingSettings),
			GlobalSecondaryIndexReadSettings: gsiRead,
			TableClass:                       tableClassToProto(replica.TableClass),
			KmsMasterKeyId:                   replica.KMSMasterKeyId,
			OnDemandThroughputOverride:       onDemandThroughputToProto(replica.OnDemandThroughputOverride),
			GlobalSecondaryIndexOverrides:    replicaGSIOverridesToProto(replica.GlobalSecondaryIndexOverrides),
		}
		if replica.TableClassLastUpdated != nil {
			result[i].TableClassLastUpdated = timestamppb.New(*replica.TableClassLastUpdated)
		}
	}
	return result
}

func replicaGSIOverridesToProto(list []*ReplicaGlobalSecondaryIndex) []*pb.ReplicaGlobalSecondaryIndex {
	if list == nil {
		return nil
	}
	result := make([]*pb.ReplicaGlobalSecondaryIndex, len(list))
	for i, entry := range list {
		result[i] = &pb.ReplicaGlobalSecondaryIndex{
			IndexName:                    entry.IndexName,
			ProvisionedReadCapacityUnits: entry.ProvisionedReadCapacityUnits,
			OnDemandThroughputOverride:   onDemandThroughputToProto(entry.OnDemandThroughputOverride),
		}
	}
	return result
}

func protoToReplicaGSIOverrides(list []*pb.ReplicaGlobalSecondaryIndex) []*ReplicaGlobalSecondaryIndex {
	if list == nil {
		return nil
	}
	result := make([]*ReplicaGlobalSecondaryIndex, len(list))
	for i, entry := range list {
		result[i] = &ReplicaGlobalSecondaryIndex{
			IndexName:                    entry.IndexName,
			ProvisionedReadCapacityUnits: entry.ProvisionedReadCapacityUnits,
			OnDemandThroughputOverride:   protoToOnDemandThroughput(entry.OnDemandThroughputOverride),
		}
	}
	return result
}

func protoToReplicas(p []*pb.Replica) []*Replica {
	if p == nil {
		return nil
	}
	result := make([]*Replica, len(p))
	for i, replica := range p {
		gsiRead := make([]IndexAutoScalingSettings, len(replica.GlobalSecondaryIndexReadSettings))
		for j, index := range replica.GlobalSecondaryIndexReadSettings {
			gsiRead[j] = protoToIndexAutoScalingSettings(index)
		}
		var tableClassUpdated *time.Time
		if replica.TableClassLastUpdated != nil {
			updated := replica.TableClassLastUpdated.AsTime()
			tableClassUpdated = &updated
		}
		result[i] = &Replica{
			RegionName:                       replica.RegionName,
			ReplicaStatus:                    protoToReplicaStatus(replica.ReplicaStatus),
			BillingMode:                      protoToBillingMode(replica.BillingMode),
			ProvisionedReadCapacityUnits:     replica.ProvisionedReadCapacityUnits,
			ProvisionedWriteCapacityUnits:    replica.ProvisionedWriteCapacityUnits,
			ReadAutoScalingSettings:          protoToAutoScalingSettings(replica.ReadAutoScalingSettings),
			GlobalSecondaryIndexReadSettings: gsiRead,
			TableClass:                       protoToTableClass(replica.TableClass),
			TableClassLastUpdated:            tableClassUpdated,
			KMSMasterKeyId:                   replica.KmsMasterKeyId,
			OnDemandThroughputOverride:       protoToOnDemandThroughput(replica.OnDemandThroughputOverride),
			GlobalSecondaryIndexOverrides:    protoToReplicaGSIOverrides(replica.GlobalSecondaryIndexOverrides),
		}
	}
	return result
}
