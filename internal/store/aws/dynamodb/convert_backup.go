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
	}
}
