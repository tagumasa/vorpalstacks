package dynamodb

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	pb "vorpalstacks/internal/pb/storage/storage_dynamodb"
)

// BackupToProto converts a Backup to its protobuf representation.
func BackupToProto(b *Backup) *pb.Backup {
	if b == nil {
		return nil
	}
	return &pb.Backup{
		BackupName:              b.BackupName,
		BackupArn:               b.BackupArn,
		SourceTableName:         b.SourceTableName,
		SourceTableArn:          b.SourceTableArn,
		SourceTableId:           b.SourceTableId,
		SourceTableCreationTime: timestamppb.New(b.SourceTableCreationTime),
		SourceTableSizeBytes:    b.SourceTableSizeBytes,
		SourceTableItemCount:    b.SourceTableItemCount,
		BackupStatus:            backupStatusToProto(b.BackupStatus),
		BackupType:              backupTypeToProto(b.BackupType),
		BackupCreationDateTime:  timestamppb.New(b.BackupCreationDateTime),
		BackupSizeBytes:         b.BackupSizeBytes,
		BackupExpiryDateTime:    timestamppb.New(b.BackupExpiryDateTime),
		KeySchema:               keySchemaToProto(b.KeySchema),
		AttributeDefinitions:    attributeDefinitionsToProto(b.AttributeDefinitions),
		BillingMode:             billingModeToProto(b.BillingMode),
		ProvisionedThroughput:   provisionedThroughputToProto(b.ProvisionedThroughput),
		GlobalSecondaryIndexes:  globalSecondaryIndexesToProto(b.GlobalSecondaryIndexes),
		LocalSecondaryIndexes:   localSecondaryIndexesToProto(b.LocalSecondaryIndexes),
		VectorIndexes:           vectorIndexesToProto(b.VectorIndexes),
	}
}

// ProtoToBackup converts a protobuf Backup to its internal representation.
func ProtoToBackup(p *pb.Backup) *Backup {
	if p == nil {
		return nil
	}
	return &Backup{
		BackupName:              p.BackupName,
		BackupArn:               p.BackupArn,
		SourceTableName:         p.SourceTableName,
		SourceTableArn:          p.SourceTableArn,
		SourceTableId:           p.SourceTableId,
		SourceTableCreationTime: p.SourceTableCreationTime.AsTime(),
		SourceTableSizeBytes:    p.SourceTableSizeBytes,
		SourceTableItemCount:    p.SourceTableItemCount,
		BackupStatus:            protoToBackupStatus(p.BackupStatus),
		BackupType:              protoToBackupType(p.BackupType),
		BackupCreationDateTime:  p.BackupCreationDateTime.AsTime(),
		BackupSizeBytes:         p.BackupSizeBytes,
		BackupExpiryDateTime:    p.BackupExpiryDateTime.AsTime(),
		KeySchema:               protoToKeySchema(p.KeySchema),
		AttributeDefinitions:    protoToAttributeDefinitions(p.AttributeDefinitions),
		BillingMode:             protoToBillingMode(p.BillingMode),
		ProvisionedThroughput:   protoToProvisionedThroughput(p.ProvisionedThroughput),
		GlobalSecondaryIndexes:  protoToGlobalSecondaryIndexes(p.GlobalSecondaryIndexes),
		LocalSecondaryIndexes:   protoToLocalSecondaryIndexes(p.LocalSecondaryIndexes),
		VectorIndexes:           protoToVectorIndexes(p.VectorIndexes),
	}
}

// backupSnapshotToProto converts a backup's item snapshot to its protobuf
// form; the key map and the attribute map are both preserved per item.
func backupSnapshotToProto(items []*Item) *pb.BackupSnapshot {
	if items == nil {
		return nil
	}
	pbItems := make([]*pb.Item, len(items))
	for i, item := range items {
		pbItems[i] = &pb.Item{
			TableName:  item.TableName,
			Key:        attributeValueMapToProtoDirect(item.Key),
			Attributes: attributeValueMapToProtoDirect(item.Attributes),
		}
	}
	return &pb.BackupSnapshot{Items: pbItems}
}

// protoToBackupSnapshot converts a protobuf item snapshot back to items.
func protoToBackupSnapshot(p *pb.BackupSnapshot) []*Item {
	if p == nil {
		return nil
	}
	items := make([]*Item, len(p.Items))
	for i := range p.Items {
		items[i] = itemFromProto(p.Items[i])
	}
	return items
}
