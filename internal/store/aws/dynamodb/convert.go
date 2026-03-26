package dynamodb

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	pb "vorpalstacks/internal/pb/storage/storage_dynamodb"
	"vorpalstacks/internal/store/aws/common"
)

// GlobalTable conversion

// GlobalTableToProto converts a GlobalTable to its protobuf representation.
func GlobalTableToProto(g *GlobalTable) *pb.GlobalTable {
	if g == nil {
		return nil
	}
	return &pb.GlobalTable{
		GlobalTableName:   g.GlobalTableName,
		GlobalTableArn:    g.GlobalTableArn,
		GlobalTableStatus: g.GlobalTableStatus,
		CreationDateTime:  timestamppb.New(g.CreationDateTime),
		ReplicationGroup:  replicasToProto(g.ReplicationGroup),
	}
}

// ProtoToGlobalTable converts a protobuf GlobalTable to its internal representation.
func ProtoToGlobalTable(p *pb.GlobalTable) *GlobalTable {
	if p == nil {
		return nil
	}
	return &GlobalTable{
		GlobalTableName:   p.GlobalTableName,
		GlobalTableArn:    p.GlobalTableArn,
		GlobalTableStatus: p.GlobalTableStatus,
		CreationDateTime:  p.CreationDateTime.AsTime(),
		ReplicationGroup:  protoToReplicas(p.ReplicationGroup),
	}
}

func replicasToProto(r []*Replica) []*pb.Replica {
	if r == nil {
		return nil
	}
	result := make([]*pb.Replica, len(r))
	for i, replica := range r {
		result[i] = &pb.Replica{
			RegionName:                    replica.RegionName,
			ReplicaStatus:                 replica.ReplicaStatus,
			BillingMode:                   replica.BillingMode,
			ProvisionedReadCapacityUnits:  replica.ProvisionedReadCapacityUnits,
			ProvisionedWriteCapacityUnits: replica.ProvisionedWriteCapacityUnits,
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
		result[i] = &Replica{
			RegionName:                    replica.RegionName,
			ReplicaStatus:                 replica.ReplicaStatus,
			BillingMode:                   replica.BillingMode,
			ProvisionedReadCapacityUnits:  replica.ProvisionedReadCapacityUnits,
			ProvisionedWriteCapacityUnits: replica.ProvisionedWriteCapacityUnits,
		}
	}
	return result
}

// TableToProto converts API Table to storage proto.
func TableToProto(t *Table) *pb.Table {
	if t == nil {
		return nil
	}

	return &pb.Table{
		Name:                          t.Name,
		Arn:                           t.ARN,
		Status:                        tableStatusToProto(t.Status),
		CreationDateTime:              timestamppb.New(t.CreationDateTime),
		LastUpdatedDateTime:           timestamppb.New(t.LastUpdatedDateTime),
		KeySchema:                     keySchemaToProto(t.KeySchema),
		AttributeDefinitions:          attributeDefinitionsToProto(t.AttributeDefinitions),
		ProvisionedThroughput:         provisionedThroughputToProto(t.ProvisionedThroughput),
		BillingMode:                   billingModeToProto(t.BillingMode),
		GlobalSecondaryIndexes:        globalSecondaryIndexesToProto(t.GlobalSecondaryIndexes),
		LocalSecondaryIndexes:         localSecondaryIndexesToProto(t.LocalSecondaryIndexes),
		StreamSpecification:           streamSpecificationToProto(t.StreamSpecification),
		SseDescription:                sseDescriptionToProto(t.SSEDescription),
		TableSizeBytes:                t.TableSizeBytes,
		ItemCount:                     t.ItemCount,
		Tags:                          tagsToProto(t.Tags),
		DeletionProtectionEnabled:     t.DeletionProtectionEnabled,
		StreamArn:                     t.StreamArn,
		LatestStreamLabel:             t.LatestStreamLabel,
		TimeToLive:                    timeToLiveToProto(t.TimeToLive),
		PointInTimeRecovery:           pointInTimeRecoveryToProto(t.PointInTimeRecovery),
		ResourcePolicy:                t.ResourcePolicy,
		KinesisDataStreamDestinations: kinesisDataStreamsToProto(t.KinesisDataStreamDestinations),
		ContributorInsightsEnabled:    t.ContributorInsightsEnabled,
	}
}

// ProtoToTable converts storage proto to API Table.
func ProtoToTable(p *pb.Table) *Table {
	if p == nil {
		return nil
	}

	return &Table{
		Name:                          p.Name,
		ARN:                           p.Arn,
		Status:                        protoToTableStatus(p.Status),
		CreationDateTime:              p.CreationDateTime.AsTime(),
		LastUpdatedDateTime:           p.LastUpdatedDateTime.AsTime(),
		KeySchema:                     protoToKeySchema(p.KeySchema),
		AttributeDefinitions:          protoToAttributeDefinitions(p.AttributeDefinitions),
		ProvisionedThroughput:         protoToProvisionedThroughput(p.ProvisionedThroughput),
		BillingMode:                   protoToBillingMode(p.BillingMode),
		GlobalSecondaryIndexes:        protoToGlobalSecondaryIndexes(p.GlobalSecondaryIndexes),
		LocalSecondaryIndexes:         protoToLocalSecondaryIndexes(p.LocalSecondaryIndexes),
		StreamSpecification:           protoToStreamSpecification(p.StreamSpecification),
		SSEDescription:                protoToSSEDescription(p.SseDescription),
		TableSizeBytes:                p.TableSizeBytes,
		ItemCount:                     p.ItemCount,
		Tags:                          protoToTags(p.Tags),
		DeletionProtectionEnabled:     p.DeletionProtectionEnabled,
		StreamArn:                     p.StreamArn,
		LatestStreamLabel:             p.LatestStreamLabel,
		TimeToLive:                    protoToTimeToLive(p.TimeToLive),
		PointInTimeRecovery:           protoToPointInTimeRecovery(p.PointInTimeRecovery),
		ResourcePolicy:                p.ResourcePolicy,
		KinesisDataStreamDestinations: protoToKinesisDataStreams(p.KinesisDataStreamDestinations),
		ContributorInsightsEnabled:    p.ContributorInsightsEnabled,
	}
}

// ItemToProto converts API Item map to proto Item.
func ItemToProto(tableName string, attrs map[string]*AttributeValue) *pb.Item {
	if attrs == nil {
		return nil
	}
	return &pb.Item{
		TableName:  tableName,
		Attributes: attributeValueMapToProtoDirect(attrs),
	}
}

// ProtoToItem converts proto Item to API Item map.
func ProtoToItem(p *pb.Item) map[string]*AttributeValue {
	if p == nil {
		return nil
	}
	return protoToAttributeValueMapDirect(p.Attributes)
}

// Enum conversion functions

func tableStatusToProto(s TableStatus) pb.TableStatus {
	switch s {
	case TableStatusCreating:
		return pb.TableStatus_TABLE_STATUS_CREATING
	case TableStatusActive:
		return pb.TableStatus_TABLE_STATUS_ACTIVE
	case TableStatusUpdating:
		return pb.TableStatus_TABLE_STATUS_UPDATING
	case TableStatusDeleting:
		return pb.TableStatus_TABLE_STATUS_DELETING
	default:
		return pb.TableStatus_TABLE_STATUS_UNSPECIFIED
	}
}

func protoToTableStatus(s pb.TableStatus) TableStatus {
	switch s {
	case pb.TableStatus_TABLE_STATUS_CREATING:
		return TableStatusCreating
	case pb.TableStatus_TABLE_STATUS_ACTIVE:
		return TableStatusActive
	case pb.TableStatus_TABLE_STATUS_UPDATING:
		return TableStatusUpdating
	case pb.TableStatus_TABLE_STATUS_DELETING:
		return TableStatusDeleting
	default:
		return ""
	}
}

func billingModeToProto(b BillingMode) pb.BillingMode {
	switch b {
	case BillingModeProvisioned:
		return pb.BillingMode_BILLING_MODE_PROVISIONED
	case BillingModePayPerRequest:
		return pb.BillingMode_BILLING_MODE_PAY_PER_REQUEST
	default:
		return pb.BillingMode_BILLING_MODE_UNSPECIFIED
	}
}

func protoToBillingMode(b pb.BillingMode) BillingMode {
	switch b {
	case pb.BillingMode_BILLING_MODE_PROVISIONED:
		return BillingModeProvisioned
	case pb.BillingMode_BILLING_MODE_PAY_PER_REQUEST:
		return BillingModePayPerRequest
	default:
		return ""
	}
}

func keyTypeToProto(k KeyType) pb.KeyType {
	switch k {
	case KeyTypeHash:
		return pb.KeyType_KEY_TYPE_HASH
	case KeyTypeRange:
		return pb.KeyType_KEY_TYPE_RANGE
	default:
		return pb.KeyType_KEY_TYPE_UNSPECIFIED
	}
}

func protoToKeyType(k pb.KeyType) KeyType {
	switch k {
	case pb.KeyType_KEY_TYPE_HASH:
		return KeyTypeHash
	case pb.KeyType_KEY_TYPE_RANGE:
		return KeyTypeRange
	default:
		return ""
	}
}

func scalarAttributeTypeToProto(s ScalarAttributeType) pb.ScalarAttributeType {
	switch s {
	case ScalarAttributeTypeS:
		return pb.ScalarAttributeType_SCALAR_ATTRIBUTE_TYPE_S
	case ScalarAttributeTypeN:
		return pb.ScalarAttributeType_SCALAR_ATTRIBUTE_TYPE_N
	case ScalarAttributeTypeB:
		return pb.ScalarAttributeType_SCALAR_ATTRIBUTE_TYPE_B
	default:
		return pb.ScalarAttributeType_SCALAR_ATTRIBUTE_TYPE_UNSPECIFIED
	}
}

func protoToScalarAttributeType(s pb.ScalarAttributeType) ScalarAttributeType {
	switch s {
	case pb.ScalarAttributeType_SCALAR_ATTRIBUTE_TYPE_S:
		return ScalarAttributeTypeS
	case pb.ScalarAttributeType_SCALAR_ATTRIBUTE_TYPE_N:
		return ScalarAttributeTypeN
	case pb.ScalarAttributeType_SCALAR_ATTRIBUTE_TYPE_B:
		return ScalarAttributeTypeB
	default:
		return ""
	}
}

func indexStatusToProto(s IndexStatus) pb.IndexStatus {
	switch s {
	case IndexStatusCreating:
		return pb.IndexStatus_INDEX_STATUS_CREATING
	case IndexStatusActive:
		return pb.IndexStatus_INDEX_STATUS_ACTIVE
	case IndexStatusUpdating:
		return pb.IndexStatus_INDEX_STATUS_UPDATING
	case IndexStatusDeleting:
		return pb.IndexStatus_INDEX_STATUS_DELETING
	default:
		return pb.IndexStatus_INDEX_STATUS_UNSPECIFIED
	}
}

func protoToIndexStatus(s pb.IndexStatus) IndexStatus {
	switch s {
	case pb.IndexStatus_INDEX_STATUS_CREATING:
		return IndexStatusCreating
	case pb.IndexStatus_INDEX_STATUS_ACTIVE:
		return IndexStatusActive
	case pb.IndexStatus_INDEX_STATUS_UPDATING:
		return IndexStatusUpdating
	case pb.IndexStatus_INDEX_STATUS_DELETING:
		return IndexStatusDeleting
	default:
		return ""
	}
}

func backupStatusToProto(s BackupStatus) pb.BackupStatus {
	switch s {
	case BackupStatusCreating:
		return pb.BackupStatus_BACKUP_STATUS_CREATING
	case BackupStatusAvailable:
		return pb.BackupStatus_BACKUP_STATUS_AVAILABLE
	case BackupStatusDeleted:
		return pb.BackupStatus_BACKUP_STATUS_DELETED
	default:
		return pb.BackupStatus_BACKUP_STATUS_UNSPECIFIED
	}
}

func protoToBackupStatus(s pb.BackupStatus) BackupStatus {
	switch s {
	case pb.BackupStatus_BACKUP_STATUS_CREATING:
		return BackupStatusCreating
	case pb.BackupStatus_BACKUP_STATUS_AVAILABLE:
		return BackupStatusAvailable
	case pb.BackupStatus_BACKUP_STATUS_DELETED:
		return BackupStatusDeleted
	default:
		return ""
	}
}

func backupTypeToProto(t BackupType) pb.BackupType {
	switch t {
	case BackupTypeUser:
		return pb.BackupType_BACKUP_TYPE_USER
	case BackupTypeSystem:
		return pb.BackupType_BACKUP_TYPE_SYSTEM
	case BackupTypeAWSBackup:
		return pb.BackupType_BACKUP_TYPE_AWS_BACKUP
	default:
		return pb.BackupType_BACKUP_TYPE_UNSPECIFIED
	}
}

func protoToBackupType(t pb.BackupType) BackupType {
	switch t {
	case pb.BackupType_BACKUP_TYPE_USER:
		return BackupTypeUser
	case pb.BackupType_BACKUP_TYPE_SYSTEM:
		return BackupTypeSystem
	case pb.BackupType_BACKUP_TYPE_AWS_BACKUP:
		return BackupTypeAWSBackup
	default:
		return ""
	}
}

// Struct conversion functions

func keySchemaToProto(ks []*KeySchemaElement) []*pb.KeySchemaElement {
	if ks == nil {
		return nil
	}
	result := make([]*pb.KeySchemaElement, len(ks))
	for i, k := range ks {
		result[i] = &pb.KeySchemaElement{
			AttributeName: k.AttributeName,
			KeyType:       keyTypeToProto(k.KeyType),
		}
	}
	return result
}

func protoToKeySchema(ks []*pb.KeySchemaElement) []*KeySchemaElement {
	if ks == nil {
		return nil
	}
	result := make([]*KeySchemaElement, len(ks))
	for i, k := range ks {
		result[i] = &KeySchemaElement{
			AttributeName: k.AttributeName,
			KeyType:       protoToKeyType(k.KeyType),
		}
	}
	return result
}

func attributeDefinitionsToProto(ad []*AttributeDefinition) []*pb.AttributeDefinition {
	if ad == nil {
		return nil
	}
	result := make([]*pb.AttributeDefinition, len(ad))
	for i, a := range ad {
		result[i] = &pb.AttributeDefinition{
			AttributeName: a.AttributeName,
			AttributeType: scalarAttributeTypeToProto(a.AttributeType),
		}
	}
	return result
}

func protoToAttributeDefinitions(ad []*pb.AttributeDefinition) []*AttributeDefinition {
	if ad == nil {
		return nil
	}
	result := make([]*AttributeDefinition, len(ad))
	for i, a := range ad {
		result[i] = &AttributeDefinition{
			AttributeName: a.AttributeName,
			AttributeType: protoToScalarAttributeType(a.AttributeType),
		}
	}
	return result
}

func provisionedThroughputToProto(pt *ProvisionedThroughput) *pb.ProvisionedThroughput {
	if pt == nil {
		return nil
	}
	return &pb.ProvisionedThroughput{
		ReadCapacityUnits:      pt.ReadCapacityUnits,
		WriteCapacityUnits:     pt.WriteCapacityUnits,
		LastDecreaseDateTime:   timestamppb.New(pt.LastDecreaseDateTime),
		LastIncreaseDateTime:   timestamppb.New(pt.LastIncreaseDateTime),
		NumberOfDecreasesToday: pt.NumberOfDecreasesToday,
	}
}

func protoToProvisionedThroughput(pt *pb.ProvisionedThroughput) *ProvisionedThroughput {
	if pt == nil {
		return nil
	}
	return &ProvisionedThroughput{
		ReadCapacityUnits:      pt.ReadCapacityUnits,
		WriteCapacityUnits:     pt.WriteCapacityUnits,
		LastDecreaseDateTime:   pt.LastDecreaseDateTime.AsTime(),
		LastIncreaseDateTime:   pt.LastIncreaseDateTime.AsTime(),
		NumberOfDecreasesToday: pt.NumberOfDecreasesToday,
	}
}

func projectionToProto(p *Projection) *pb.Projection {
	if p == nil {
		return nil
	}
	return &pb.Projection{
		ProjectionType:   p.ProjectionType,
		NonKeyAttributes: p.NonKeyAttributes,
	}
}

func protoToProjection(p *pb.Projection) *Projection {
	if p == nil {
		return nil
	}
	return &Projection{
		ProjectionType:   p.ProjectionType,
		NonKeyAttributes: p.NonKeyAttributes,
	}
}

func globalSecondaryIndexesToProto(idx []*GlobalSecondaryIndex) []*pb.GlobalSecondaryIndex {
	if idx == nil {
		return nil
	}
	result := make([]*pb.GlobalSecondaryIndex, len(idx))
	for i, g := range idx {
		result[i] = &pb.GlobalSecondaryIndex{
			IndexName:             g.IndexName,
			IndexArn:              g.IndexArn,
			KeySchema:             keySchemaToProto(g.KeySchema),
			Projection:            projectionToProto(g.Projection),
			ProvisionedThroughput: provisionedThroughputToProto(g.ProvisionedThroughput),
			IndexStatus:           indexStatusToProto(g.IndexStatus),
			IndexSizeBytes:        g.IndexSizeBytes,
			ItemCount:             g.ItemCount,
		}
	}
	return result
}

func protoToGlobalSecondaryIndexes(idx []*pb.GlobalSecondaryIndex) []*GlobalSecondaryIndex {
	if idx == nil {
		return nil
	}
	result := make([]*GlobalSecondaryIndex, len(idx))
	for i, g := range idx {
		result[i] = &GlobalSecondaryIndex{
			IndexName:             g.IndexName,
			IndexArn:              g.IndexArn,
			KeySchema:             protoToKeySchema(g.KeySchema),
			Projection:            protoToProjection(g.Projection),
			ProvisionedThroughput: protoToProvisionedThroughput(g.ProvisionedThroughput),
			IndexStatus:           protoToIndexStatus(g.IndexStatus),
			IndexSizeBytes:        g.IndexSizeBytes,
			ItemCount:             g.ItemCount,
		}
	}
	return result
}

func localSecondaryIndexesToProto(idx []*LocalSecondaryIndex) []*pb.LocalSecondaryIndex {
	if idx == nil {
		return nil
	}
	result := make([]*pb.LocalSecondaryIndex, len(idx))
	for i, l := range idx {
		result[i] = &pb.LocalSecondaryIndex{
			IndexName:      l.IndexName,
			KeySchema:      keySchemaToProto(l.KeySchema),
			Projection:     projectionToProto(l.Projection),
			IndexSizeBytes: l.IndexSizeBytes,
			ItemCount:      l.ItemCount,
		}
	}
	return result
}

func protoToLocalSecondaryIndexes(idx []*pb.LocalSecondaryIndex) []*LocalSecondaryIndex {
	if idx == nil {
		return nil
	}
	result := make([]*LocalSecondaryIndex, len(idx))
	for i, l := range idx {
		result[i] = &LocalSecondaryIndex{
			IndexName:      l.IndexName,
			KeySchema:      protoToKeySchema(l.KeySchema),
			Projection:     protoToProjection(l.Projection),
			IndexSizeBytes: l.IndexSizeBytes,
			ItemCount:      l.ItemCount,
		}
	}
	return result
}

func streamSpecificationToProto(s *StreamSpecification) *pb.StreamSpecification {
	if s == nil {
		return nil
	}
	return &pb.StreamSpecification{
		StreamEnabled:  s.StreamEnabled,
		StreamViewType: streamViewTypeToProto(s.StreamViewType),
	}
}

func protoToStreamSpecification(s *pb.StreamSpecification) *StreamSpecification {
	if s == nil {
		return nil
	}
	return &StreamSpecification{
		StreamEnabled:  s.StreamEnabled,
		StreamViewType: protoToStreamViewType(s.StreamViewType),
	}
}

func streamViewTypeToProto(s StreamViewType) pb.StreamViewType {
	switch s {
	case StreamViewTypeNewImage:
		return pb.StreamViewType_STREAM_VIEW_TYPE_NEW_IMAGE
	case StreamViewTypeOldImage:
		return pb.StreamViewType_STREAM_VIEW_TYPE_OLD_IMAGE
	case StreamViewTypeNewAndOldImages:
		return pb.StreamViewType_STREAM_VIEW_TYPE_NEW_AND_OLD_IMAGES
	case StreamViewTypeKeysOnly:
		return pb.StreamViewType_STREAM_VIEW_TYPE_KEYS_ONLY
	default:
		return pb.StreamViewType_STREAM_VIEW_TYPE_UNSPECIFIED
	}
}

func protoToStreamViewType(s pb.StreamViewType) StreamViewType {
	switch s {
	case pb.StreamViewType_STREAM_VIEW_TYPE_NEW_IMAGE:
		return StreamViewTypeNewImage
	case pb.StreamViewType_STREAM_VIEW_TYPE_OLD_IMAGE:
		return StreamViewTypeOldImage
	case pb.StreamViewType_STREAM_VIEW_TYPE_NEW_AND_OLD_IMAGES:
		return StreamViewTypeNewAndOldImages
	case pb.StreamViewType_STREAM_VIEW_TYPE_KEYS_ONLY:
		return StreamViewTypeKeysOnly
	default:
		return ""
	}
}

func sseDescriptionToProto(s *SSEDescription) *pb.SSEDescription {
	if s == nil {
		return nil
	}
	return &pb.SSEDescription{
		Status:                         s.Status,
		SseType:                        sseTypeToProto(s.SSEType),
		KmsMasterKeyArn:                s.KMSMasterKeyArn,
		InaccessibleEncryptionDateTime: timestamppb.New(s.InaccessibleEncryptionDateTime),
	}
}

func protoToSSEDescription(s *pb.SSEDescription) *SSEDescription {
	if s == nil {
		return nil
	}
	return &SSEDescription{
		Status:                         s.Status,
		SSEType:                        protoToSSEType(s.SseType),
		KMSMasterKeyArn:                s.KmsMasterKeyArn,
		InaccessibleEncryptionDateTime: s.InaccessibleEncryptionDateTime.AsTime(),
	}
}

func sseTypeToProto(s SSEType) pb.SSEType {
	switch s {
	case SSETypeAES256:
		return pb.SSEType_SSE_TYPE_AES256
	case SSETypeKMS:
		return pb.SSEType_SSE_TYPE_KMS
	default:
		return pb.SSEType_SSE_TYPE_UNSPECIFIED
	}
}

func protoToSSEType(s pb.SSEType) SSEType {
	switch s {
	case pb.SSEType_SSE_TYPE_AES256:
		return SSETypeAES256
	case pb.SSEType_SSE_TYPE_KMS:
		return SSETypeKMS
	default:
		return ""
	}
}

func tagsToProto(t []common.Tag) []*pb.Tag {
	if t == nil {
		return nil
	}
	result := make([]*pb.Tag, len(t))
	for i, tag := range t {
		result[i] = &pb.Tag{
			Key:   tag.Key,
			Value: tag.Value,
		}
	}
	return result
}

func protoToTags(t []*pb.Tag) []common.Tag {
	if t == nil {
		return nil
	}
	result := make([]common.Tag, len(t))
	for i, tag := range t {
		result[i] = common.Tag{
			Key:   tag.Key,
			Value: tag.Value,
		}
	}
	return result
}

func timeToLiveToProto(t *TimeToLiveSpecification) *pb.TimeToLiveSpecification {
	if t == nil {
		return nil
	}
	return &pb.TimeToLiveSpecification{
		Enabled:       t.Enabled,
		AttributeName: t.AttributeName,
		Status:        ttlStatusToProto(t.Status),
	}
}

func protoToTimeToLive(t *pb.TimeToLiveSpecification) *TimeToLiveSpecification {
	if t == nil {
		return nil
	}
	return &TimeToLiveSpecification{
		Enabled:       t.Enabled,
		AttributeName: t.AttributeName,
		Status:        protoToTTLStatus(t.Status),
	}
}

func ttlStatusToProto(s TTLStatus) pb.TTLStatus {
	switch s {
	case TTLStatusEnabling:
		return pb.TTLStatus_TTL_STATUS_ENABLING
	case TTLStatusEnabled:
		return pb.TTLStatus_TTL_STATUS_ENABLED
	case TTLStatusDisabling:
		return pb.TTLStatus_TTL_STATUS_DISABLING
	case TTLStatusDisabled:
		return pb.TTLStatus_TTL_STATUS_DISABLED
	default:
		return pb.TTLStatus_TTL_STATUS_UNSPECIFIED
	}
}

func protoToTTLStatus(s pb.TTLStatus) TTLStatus {
	switch s {
	case pb.TTLStatus_TTL_STATUS_ENABLING:
		return TTLStatusEnabling
	case pb.TTLStatus_TTL_STATUS_ENABLED:
		return TTLStatusEnabled
	case pb.TTLStatus_TTL_STATUS_DISABLING:
		return TTLStatusDisabling
	case pb.TTLStatus_TTL_STATUS_DISABLED:
		return TTLStatusDisabled
	default:
		return ""
	}
}

func pointInTimeRecoveryToProto(p *PointInTimeRecoveryDescription) *pb.PointInTimeRecoveryDescription {
	if p == nil {
		return nil
	}
	return &pb.PointInTimeRecoveryDescription{
		Status:                     pointInTimeRecoveryStatusToProto(p.Status),
		EarliestRestorableDateTime: timestamppb.New(p.EarliestRestorableDateTime),
		LatestRestorableDateTime:   timestamppb.New(p.LatestRestorableDateTime),
	}
}

func protoToPointInTimeRecovery(p *pb.PointInTimeRecoveryDescription) *PointInTimeRecoveryDescription {
	if p == nil {
		return nil
	}
	return &PointInTimeRecoveryDescription{
		Status:                     protoToPointInTimeRecoveryStatus(p.Status),
		EarliestRestorableDateTime: p.EarliestRestorableDateTime.AsTime(),
		LatestRestorableDateTime:   p.LatestRestorableDateTime.AsTime(),
	}
}

func pointInTimeRecoveryStatusToProto(s PointInTimeRecoveryStatus) pb.PointInTimeRecoveryStatus {
	switch s {
	case PITRStatusEnabled:
		return pb.PointInTimeRecoveryStatus_POINT_IN_TIME_RECOVERY_STATUS_ENABLED
	case PITRStatusDisabled:
		return pb.PointInTimeRecoveryStatus_POINT_IN_TIME_RECOVERY_STATUS_DISABLED
	default:
		return pb.PointInTimeRecoveryStatus_POINT_IN_TIME_RECOVERY_STATUS_UNSPECIFIED
	}
}

func protoToPointInTimeRecoveryStatus(s pb.PointInTimeRecoveryStatus) PointInTimeRecoveryStatus {
	switch s {
	case pb.PointInTimeRecoveryStatus_POINT_IN_TIME_RECOVERY_STATUS_ENABLED:
		return PITRStatusEnabled
	case pb.PointInTimeRecoveryStatus_POINT_IN_TIME_RECOVERY_STATUS_DISABLED:
		return PITRStatusDisabled
	default:
		return ""
	}
}

func kinesisDataStreamsToProto(k []*KinesisDataStreamDestination) []*pb.KinesisDataStreamDestination {
	if k == nil {
		return nil
	}
	result := make([]*pb.KinesisDataStreamDestination, len(k))
	for i, d := range k {
		result[i] = &pb.KinesisDataStreamDestination{
			StreamArn:                    d.StreamArn,
			DestinationStatus:            d.DestinationStatus,
			DestinationStatusDescription: d.DestinationStatusDescription,
		}
	}
	return result
}

func protoToKinesisDataStreams(k []*pb.KinesisDataStreamDestination) []*KinesisDataStreamDestination {
	if k == nil {
		return nil
	}
	result := make([]*KinesisDataStreamDestination, len(k))
	for i, d := range k {
		result[i] = &KinesisDataStreamDestination{
			StreamArn:                    d.StreamArn,
			DestinationStatus:            d.DestinationStatus,
			DestinationStatusDescription: d.DestinationStatusDescription,
		}
	}
	return result
}

// AttributeValue conversion functions (CRITICAL: Recursive structure handling)

func attributeValueToProto(av *AttributeValue) *pb.AttributeValue {
	if av == nil {
		return nil
	}

	switch {
	case av.S != nil:
		return &pb.AttributeValue{Value: &pb.AttributeValue_S{S: *av.S}}
	case av.N != nil:
		return &pb.AttributeValue{Value: &pb.AttributeValue_N{N: *av.N}}
	case av.B != nil:
		return &pb.AttributeValue{Value: &pb.AttributeValue_B{B: av.B}}
	case av.SS != nil:
		return &pb.AttributeValue{Value: &pb.AttributeValue_Ss{Ss: &pb.StringSet{Values: av.SS}}}
	case av.NS != nil:
		return &pb.AttributeValue{Value: &pb.AttributeValue_Ns{Ns: &pb.NumberSet{Values: av.NS}}}
	case av.BS != nil:
		return &pb.AttributeValue{Value: &pb.AttributeValue_Bs{Bs: &pb.BytesSet{Values: av.BS}}}
	case av.M != nil:
		return &pb.AttributeValue{Value: &pb.AttributeValue_M{M: attributeValueMapToProto(av.M)}}
	case av.L != nil:
		return &pb.AttributeValue{Value: &pb.AttributeValue_L{L: attributeValueListToProto(av.L)}}
	case av.NULL != nil && *av.NULL:
		return &pb.AttributeValue{Value: &pb.AttributeValue_Null{Null: &pb.NullValue{}}}
	case av.BOOL != nil:
		return &pb.AttributeValue{Value: &pb.AttributeValue_BoolVal{BoolVal: *av.BOOL}}
	default:
		return nil
	}
}

func protoToAttributeValue(p *pb.AttributeValue) *AttributeValue {
	if p == nil {
		return nil
	}

	switch v := p.Value.(type) {
	case *pb.AttributeValue_S:
		return &AttributeValue{S: &v.S}
	case *pb.AttributeValue_N:
		return &AttributeValue{N: &v.N}
	case *pb.AttributeValue_B:
		return &AttributeValue{B: v.B}
	case *pb.AttributeValue_Ss:
		return &AttributeValue{SS: v.Ss.Values}
	case *pb.AttributeValue_Ns:
		return &AttributeValue{NS: v.Ns.Values}
	case *pb.AttributeValue_Bs:
		return &AttributeValue{BS: v.Bs.Values}
	case *pb.AttributeValue_M:
		return &AttributeValue{M: protoToAttributeValueMap(v.M)}
	case *pb.AttributeValue_L:
		return &AttributeValue{L: protoToAttributeValueList(v.L)}
	case *pb.AttributeValue_Null:
		t := true
		return &AttributeValue{NULL: &t}
	case *pb.AttributeValue_BoolVal:
		return &AttributeValue{BOOL: &v.BoolVal}
	default:
		return nil
	}
}

func attributeValueMapToProto(m map[string]*AttributeValue) *pb.MapValue {
	if m == nil {
		return nil
	}
	entries := make(map[string]*pb.AttributeValue, len(m))
	for k, v := range m {
		entries[k] = attributeValueToProto(v)
	}
	return &pb.MapValue{Entries: entries}
}

func protoToAttributeValueMap(p *pb.MapValue) map[string]*AttributeValue {
	if p == nil {
		return nil
	}
	m := make(map[string]*AttributeValue, len(p.Entries))
	for k, v := range p.Entries {
		m[k] = protoToAttributeValue(v)
	}
	return m
}

func attributeValueListToProto(l []*AttributeValue) *pb.ListValue {
	if l == nil {
		return nil
	}
	values := make([]*pb.AttributeValue, len(l))
	for i, v := range l {
		values[i] = attributeValueToProto(v)
	}
	return &pb.ListValue{Values: values}
}

func protoToAttributeValueList(p *pb.ListValue) []*AttributeValue {
	if p == nil {
		return nil
	}
	l := make([]*AttributeValue, len(p.Values))
	for i, v := range p.Values {
		l[i] = protoToAttributeValue(v)
	}
	return l
}

// Direct map conversion for Item's Key/Attributes fields (proto uses map<string, AttributeValue>)

func attributeValueMapToProtoDirect(m map[string]*AttributeValue) map[string]*pb.AttributeValue {
	if m == nil {
		return nil
	}
	result := make(map[string]*pb.AttributeValue, len(m))
	for k, v := range m {
		result[k] = attributeValueToProto(v)
	}
	return result
}

func protoToAttributeValueMapDirect(m map[string]*pb.AttributeValue) map[string]*AttributeValue {
	if m == nil {
		return nil
	}
	result := make(map[string]*AttributeValue, len(m))
	for k, v := range m {
		result[k] = protoToAttributeValue(v)
	}
	return result
}

// Backup conversion

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

// ExportDescription conversion

// ExportDescriptionToProto converts an ExportDescription to its protobuf representation.
func ExportDescriptionToProto(e *ExportDescription) *pb.ExportDescription {
	if e == nil {
		return nil
	}
	return &pb.ExportDescription{
		ExportArn:         e.ExportArn,
		ExportStatus:      e.ExportStatus,
		StartTime:         timestamppb.New(e.StartTime),
		EndTime:           timestamppb.New(e.EndTime),
		ManifestFilesSize: e.ManifestFilesSize,
		ItemCount:         e.ItemCount,
		TableArn:          e.TableArn,
		TableId:           e.TableId,
		ExportFormat:      e.ExportFormat,
		S3Bucket:          e.S3Bucket,
		S3Prefix:          e.S3Prefix,
		FailureCode:       e.FailureCode,
		FailureMessage:    e.FailureMessage,
	}
}

// ProtoToExportDescription converts a protobuf ExportDescription to its internal representation.
func ProtoToExportDescription(p *pb.ExportDescription) *ExportDescription {
	if p == nil {
		return nil
	}
	return &ExportDescription{
		ExportArn:         p.ExportArn,
		ExportStatus:      p.ExportStatus,
		StartTime:         p.StartTime.AsTime(),
		EndTime:           p.EndTime.AsTime(),
		ManifestFilesSize: p.ManifestFilesSize,
		ItemCount:         p.ItemCount,
		TableArn:          p.TableArn,
		TableId:           p.TableId,
		ExportFormat:      p.ExportFormat,
		S3Bucket:          p.S3Bucket,
		S3Prefix:          p.S3Prefix,
		FailureCode:       p.FailureCode,
		FailureMessage:    p.FailureMessage,
	}
}

// ImportTableDescription conversion

// ImportTableDescriptionToProto converts an ImportTableDescription to its protobuf representation.
func ImportTableDescriptionToProto(i *ImportTableDescription) *pb.ImportTableDescription {
	if i == nil {
		return nil
	}
	return &pb.ImportTableDescription{
		ImportArn:          i.ImportArn,
		ImportStatus:       i.ImportStatus,
		TableArn:           i.TableArn,
		TableId:            i.TableId,
		StartTime:          timestamppb.New(i.StartTime),
		EndTime:            timestamppb.New(i.EndTime),
		ProcessedItemCount: i.ProcessedItemCount,
		ProcessedSizeBytes: i.ProcessedSizeBytes,
		InputFormat:        i.InputFormat,
		S3BucketSource:     s3BucketSourceToProto(i.S3BucketSource),
		FailureCode:        i.FailureCode,
		FailureMessage:     i.FailureMessage,
	}
}

// ProtoToImportTableDescription converts a protobuf ImportTableDescription to its internal representation.
func ProtoToImportTableDescription(p *pb.ImportTableDescription) *ImportTableDescription {
	if p == nil {
		return nil
	}
	return &ImportTableDescription{
		ImportArn:          p.ImportArn,
		ImportStatus:       p.ImportStatus,
		TableArn:           p.TableArn,
		TableId:            p.TableId,
		StartTime:          p.StartTime.AsTime(),
		EndTime:            p.EndTime.AsTime(),
		ProcessedItemCount: p.ProcessedItemCount,
		ProcessedSizeBytes: p.ProcessedSizeBytes,
		InputFormat:        p.InputFormat,
		S3BucketSource:     protoToS3BucketSource(p.S3BucketSource),
		FailureCode:        p.FailureCode,
		FailureMessage:     p.FailureMessage,
	}
}

// S3BucketSource conversion

func s3BucketSourceToProto(s *S3BucketSource) *pb.S3BucketSource {
	if s == nil {
		return nil
	}
	return &pb.S3BucketSource{
		S3Bucket:      s.S3Bucket,
		S3Prefix:      s.S3Prefix,
		S3BucketOwner: s.S3BucketOwner,
	}
}

func protoToS3BucketSource(p *pb.S3BucketSource) *S3BucketSource {
	if p == nil {
		return nil
	}
	return &S3BucketSource{
		S3Bucket:      p.S3Bucket,
		S3Prefix:      p.S3Prefix,
		S3BucketOwner: p.S3BucketOwner,
	}
}
