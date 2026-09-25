package dynamodb

import (
	"time"

	pb "vorpalstacks/internal/pb/storage/storage_dynamodb"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// restoreSummaryToProto converts a RestoreSummary to its protobuf form.
func restoreSummaryToProto(r *RestoreSummary) *pb.RestoreSummary {
	if r == nil {
		return nil
	}
	return &pb.RestoreSummary{
		SourceBackupArn:   r.SourceBackupArn,
		SourceTableArn:    r.SourceTableArn,
		RestoreDateTime:   timestamppb.New(r.RestoreDateTime),
		RestoreInProgress: r.RestoreInProgress,
	}
}

func protoToRestoreSummary(p *pb.RestoreSummary) *RestoreSummary {
	if p == nil {
		return nil
	}
	return &RestoreSummary{
		SourceBackupArn:   p.SourceBackupArn,
		SourceTableArn:    p.SourceTableArn,
		RestoreDateTime:   p.RestoreDateTime.AsTime(),
		RestoreInProgress: p.RestoreInProgress,
	}
}

func TableToProto(t *Table) *pb.Table {
	if t == nil {
		return nil
	}

	return &pb.Table{
		Name:                          t.Name,
		Arn:                           t.ARN,
		TableId:                       t.TableId,
		Status:                        tableStatusToProto(t.Status),
		CreationDateTime:              timestamppb.New(t.CreationDateTime),
		LastUpdatedDateTime:           timestamppb.New(t.LastUpdatedDateTime),
		KeySchema:                     keySchemaToProto(t.KeySchema),
		AttributeDefinitions:          attributeDefinitionsToProto(t.AttributeDefinitions),
		ProvisionedThroughput:         provisionedThroughputToProto(t.ProvisionedThroughput),
		BillingMode:                   billingModeToProto(t.BillingMode),
		GlobalSecondaryIndexes:        globalSecondaryIndexesToProto(t.GlobalSecondaryIndexes),
		LocalSecondaryIndexes:         localSecondaryIndexesToProto(t.LocalSecondaryIndexes),
		VectorIndexes:                 vectorIndexesToProto(t.VectorIndexes),
		StreamSpecification:           streamSpecificationToProto(t.StreamSpecification),
		SseDescription:                sseDescriptionToProto(t.SSEDescription),
		TableSizeBytes:                t.TableSizeBytes,
		ItemCount:                     t.ItemCount,
		DeletionProtectionEnabled:     t.DeletionProtectionEnabled,
		StreamArn:                     t.StreamArn,
		LatestStreamLabel:             t.LatestStreamLabel,
		TimeToLive:                    timeToLiveToProto(t.TimeToLive),
		PointInTimeRecovery:           pointInTimeRecoveryToProto(t.PointInTimeRecovery),
		ResourcePolicy:                t.ResourcePolicy,
		KinesisDataStreamDestinations: kinesisDataStreamsToProto(t.KinesisDataStreamDestinations),
		ContributorInsightsEnabled:    t.ContributorInsightsEnabled,
		ContributorInsightsMode:       contributorInsightsModeToProto(t.ContributorInsightsMode),
		TableClass:                    tableClassToProto(t.TableClass),
		WarmThroughput:                warmThroughputToProto(t.WarmThroughput),
		OnDemandThroughput:            onDemandThroughputToProto(t.OnDemandThroughput),
		RestoreSummary:                restoreSummaryToProto(t.RestoreSummary),
		ResourcePolicyRevisionId:      int32(t.ResourcePolicyRevisionId),
		ContributorInsightsUpdatedAt:  timestamppb.New(t.ContributorInsightsUpdatedAt),
		GlobalTableSourceArn:          t.GlobalTableSourceArn,
		BillingModeSwitches:           timesToProto(t.BillingModeSwitches),
	}
}

// ProtoToTable converts storage proto to API Table.
func ProtoToTable(p *pb.Table) *Table {
	if p == nil {
		return nil
	}

	table := &Table{
		Name:                          p.Name,
		TableId:                       p.TableId,
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
		VectorIndexes:                 protoToVectorIndexes(p.VectorIndexes),
		StreamSpecification:           protoToStreamSpecification(p.StreamSpecification),
		SSEDescription:                protoToSSEDescription(p.SseDescription),
		TableSizeBytes:                p.TableSizeBytes,
		ItemCount:                     p.ItemCount,
		DeletionProtectionEnabled:     p.DeletionProtectionEnabled,
		StreamArn:                     p.StreamArn,
		LatestStreamLabel:             p.LatestStreamLabel,
		TimeToLive:                    protoToTimeToLive(p.TimeToLive),
		PointInTimeRecovery:           protoToPointInTimeRecovery(p.PointInTimeRecovery),
		ResourcePolicy:                p.ResourcePolicy,
		KinesisDataStreamDestinations: protoToKinesisDataStreams(p.KinesisDataStreamDestinations),
		ContributorInsightsEnabled:    p.ContributorInsightsEnabled,
		ContributorInsightsMode:       protoToContributorInsightsMode(p.ContributorInsightsMode),
		TableClass:                    protoToTableClass(p.GetTableClass()),
		WarmThroughput:                protoToWarmThroughput(p.WarmThroughput),
		OnDemandThroughput:            protoToOnDemandThroughput(p.OnDemandThroughput),
		RestoreSummary:                protoToRestoreSummary(p.RestoreSummary),
		ResourcePolicyRevisionId:      int(p.ResourcePolicyRevisionId),
		GlobalTableSourceArn:          p.GlobalTableSourceArn,
		BillingModeSwitches:           protoToTimes(p.BillingModeSwitches),
	}
	if p.ContributorInsightsUpdatedAt != nil {
		table.ContributorInsightsUpdatedAt = p.ContributorInsightsUpdatedAt.AsTime()
	}
	return table
}

// timesToProto converts a timestamp slice to its protobuf form, returning
// nil for an empty slice so an absent window serialises without the field.
func timesToProto(ts []time.Time) []*timestamppb.Timestamp {
	if len(ts) == 0 {
		return nil
	}
	result := make([]*timestamppb.Timestamp, len(ts))
	for i, t := range ts {
		result[i] = timestamppb.New(t)
	}
	return result
}

// protoToTimes converts a protobuf timestamp slice back, returning nil for
// an empty slice.
func protoToTimes(ts []*timestamppb.Timestamp) []time.Time {
	if len(ts) == 0 {
		return nil
	}
	result := make([]time.Time, len(ts))
	for i, t := range ts {
		result[i] = t.AsTime()
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

func sseDescriptionToProto(s *SSEDescription) *pb.SSEDescription {
	if s == nil {
		return nil
	}
	return &pb.SSEDescription{
		Status:                         sseStatusToProto(s.Status),
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
		Status:                         protoToSSEStatus(s.Status),
		SSEType:                        protoToSSEType(s.SseType),
		KMSMasterKeyArn:                s.KmsMasterKeyArn,
		InaccessibleEncryptionDateTime: s.InaccessibleEncryptionDateTime.AsTime(),
	}
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

func pointInTimeRecoveryToProto(p *PointInTimeRecoveryDescription) *pb.PointInTimeRecoveryDescription {
	if p == nil {
		return nil
	}
	result := &pb.PointInTimeRecoveryDescription{
		Status:                     pointInTimeRecoveryStatusToProto(p.Status),
		EarliestRestorableDateTime: timestamppb.New(p.EarliestRestorableDateTime),
		LatestRestorableDateTime:   timestamppb.New(p.LatestRestorableDateTime),
	}
	if p.RecoveryPeriodInDays > 0 {
		result.RecoveryPeriodInDays = int32(p.RecoveryPeriodInDays)
	}
	return result
}

func protoToPointInTimeRecovery(p *pb.PointInTimeRecoveryDescription) *PointInTimeRecoveryDescription {
	if p == nil {
		return nil
	}
	result := &PointInTimeRecoveryDescription{
		Status:                     protoToPointInTimeRecoveryStatus(p.Status),
		EarliestRestorableDateTime: p.EarliestRestorableDateTime.AsTime(),
		LatestRestorableDateTime:   p.LatestRestorableDateTime.AsTime(),
	}
	if p.RecoveryPeriodInDays > 0 {
		result.RecoveryPeriodInDays = int(p.RecoveryPeriodInDays)
	}
	return result
}

func kinesisDataStreamsToProto(k []*KinesisDataStreamDestination) []*pb.KinesisDataStreamDestination {
	if k == nil {
		return nil
	}
	result := make([]*pb.KinesisDataStreamDestination, len(k))
	for i, d := range k {
		result[i] = &pb.KinesisDataStreamDestination{
			StreamArn:                            d.StreamArn,
			DestinationStatus:                    destinationStatusToProto(d.DestinationStatus),
			DestinationStatusDescription:         d.DestinationStatusDescription,
			ApproximateCreationDateTimePrecision: approximateCreationDateTimePrecisionToProto(d.ApproximateCreationDateTimePrecision),
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
			StreamArn:                            d.StreamArn,
			DestinationStatus:                    protoToDestinationStatus(d.DestinationStatus),
			DestinationStatusDescription:         d.DestinationStatusDescription,
			ApproximateCreationDateTimePrecision: protoToApproximateCreationDateTimePrecision(d.ApproximateCreationDateTimePrecision),
		}
	}
	return result
}
