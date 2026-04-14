package timestream

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	pb "vorpalstacks/internal/pb/storage/storage_timestream"
)

// Database conversion

// DatabaseToProto converts a Database to its protobuf representation.
func DatabaseToProto(d *Database) *pb.Database {
	if d == nil {
		return nil
	}
	return &pb.Database{
		DatabaseName:    d.DatabaseName,
		Arn:             d.ARN,
		TableCount:      d.TableCount,
		KmsKeyId:        d.KmsKeyId,
		CreationTime:    timestamppb.New(d.CreationTime),
		LastUpdatedTime: timestamppb.New(d.LastUpdatedTime),
	}
}

// ProtoToDatabase converts a protobuf Database to its internal representation.
func ProtoToDatabase(p *pb.Database) *Database {
	if p == nil {
		return nil
	}
	return &Database{
		DatabaseName:    p.DatabaseName,
		ARN:             p.Arn,
		TableCount:      p.TableCount,
		KmsKeyId:        p.KmsKeyId,
		CreationTime:    p.CreationTime.AsTime(),
		LastUpdatedTime: p.LastUpdatedTime.AsTime(),
	}
}

// RetentionProperties conversion

// RetentionPropertiesToProto converts a RetentionProperties to its protobuf representation.
func RetentionPropertiesToProto(r *RetentionProperties) *pb.RetentionProperties {
	if r == nil {
		return nil
	}
	return &pb.RetentionProperties{
		MemoryStoreRetentionPeriodInHours:  r.MemoryStoreRetentionPeriodInHours,
		MagneticStoreRetentionPeriodInDays: r.MagneticStoreRetentionPeriodInDays,
	}
}

// ProtoToRetentionProperties converts a protobuf RetentionProperties to its internal representation.
func ProtoToRetentionProperties(p *pb.RetentionProperties) *RetentionProperties {
	if p == nil {
		return nil
	}
	return &RetentionProperties{
		MemoryStoreRetentionPeriodInHours:  p.MemoryStoreRetentionPeriodInHours,
		MagneticStoreRetentionPeriodInDays: p.MagneticStoreRetentionPeriodInDays,
	}
}

// PartitionKey conversion

// PartitionKeyToProto converts a PartitionKey to its protobuf representation.
func PartitionKeyToProto(pk *PartitionKey) *pb.PartitionKey {
	if pk == nil {
		return nil
	}
	return &pb.PartitionKey{
		Type:                string(pk.Type),
		Name:                pk.Name,
		EnforcementInRecord: string(pk.EnforcementInRecord),
	}
}

// ProtoToPartitionKey converts a protobuf PartitionKey to its internal representation.
func ProtoToPartitionKey(p *pb.PartitionKey) *PartitionKey {
	if p == nil {
		return nil
	}
	return &PartitionKey{
		Type:                PartitionKeyType(p.Type),
		Name:                p.Name,
		EnforcementInRecord: EnforcementInRecord(p.EnforcementInRecord),
	}
}

func partitionKeysToProto(pks []PartitionKey) []*pb.PartitionKey {
	if pks == nil {
		return nil
	}
	result := make([]*pb.PartitionKey, len(pks))
	for i, pk := range pks {
		result[i] = PartitionKeyToProto(&pk)
	}
	return result
}

func protoToPartitionKeys(pbs []*pb.PartitionKey) []PartitionKey {
	if pbs == nil {
		return nil
	}
	result := make([]PartitionKey, len(pbs))
	for i, pb := range pbs {
		result[i] = *ProtoToPartitionKey(pb)
	}
	return result
}

// Schema conversion

// SchemaToProto converts a Schema to its protobuf representation.
func SchemaToProto(s *Schema) *pb.Schema {
	if s == nil {
		return nil
	}
	return &pb.Schema{
		CompositePartitionKey: partitionKeysToProto(s.CompositePartitionKey),
	}
}

// ProtoToSchema converts a protobuf Schema to its internal representation.
func ProtoToSchema(p *pb.Schema) *Schema {
	if p == nil {
		return nil
	}
	return &Schema{
		CompositePartitionKey: protoToPartitionKeys(p.CompositePartitionKey),
	}
}

// Table conversion

// TableToProto converts a Table to its protobuf representation.
func TableToProto(t *Table) *pb.Table {
	if t == nil {
		return nil
	}
	return &pb.Table{
		TableName:           t.TableName,
		DatabaseName:        t.DatabaseName,
		Arn:                 t.ARN,
		TableStatus:         string(t.TableStatus),
		RetentionProperties: RetentionPropertiesToProto(t.RetentionProperties),
		Schema:              SchemaToProto(t.Schema),
		CreationTime:        timestamppb.New(t.CreationTime),
		LastUpdatedTime:     timestamppb.New(t.LastUpdatedTime),
	}
}

// ProtoToTable converts a protobuf Table to its internal representation.
func ProtoToTable(p *pb.Table) *Table {
	if p == nil {
		return nil
	}
	return &Table{
		TableName:           p.TableName,
		DatabaseName:        p.DatabaseName,
		ARN:                 p.Arn,
		TableStatus:         TableStatus(p.TableStatus),
		RetentionProperties: ProtoToRetentionProperties(p.RetentionProperties),
		Schema:              ProtoToSchema(p.Schema),
		CreationTime:        p.CreationTime.AsTime(),
		LastUpdatedTime:     p.LastUpdatedTime.AsTime(),
	}
}

// ScheduleConfiguration conversion

// ScheduleConfigurationToProto converts a ScheduleConfiguration to its protobuf representation.
func ScheduleConfigurationToProto(s *ScheduleConfiguration) *pb.ScheduleConfiguration {
	if s == nil {
		return nil
	}
	return &pb.ScheduleConfiguration{
		ScheduleExpression: s.ScheduleExpression,
	}
}

// ProtoToScheduleConfiguration converts a protobuf ScheduleConfiguration to its internal representation.
func ProtoToScheduleConfiguration(p *pb.ScheduleConfiguration) *ScheduleConfiguration {
	if p == nil {
		return nil
	}
	return &ScheduleConfiguration{
		ScheduleExpression: p.ScheduleExpression,
	}
}

// SNSConfiguration conversion

// SNSConfigurationToProto converts an SNSConfiguration to its protobuf representation.
func SNSConfigurationToProto(s *SNSConfiguration) *pb.SNSConfiguration {
	if s == nil {
		return nil
	}
	return &pb.SNSConfiguration{
		TopicArn: s.TopicARN,
	}
}

// ProtoToSNSConfiguration converts a protobuf SNSConfiguration to its internal representation.
func ProtoToSNSConfiguration(p *pb.SNSConfiguration) *SNSConfiguration {
	if p == nil {
		return nil
	}
	return &SNSConfiguration{
		TopicARN: p.TopicArn,
	}
}

// NotificationConfiguration conversion

// NotificationConfigurationToProto converts a NotificationConfiguration to its protobuf representation.
func NotificationConfigurationToProto(n *NotificationConfiguration) *pb.NotificationConfiguration {
	if n == nil {
		return nil
	}
	return &pb.NotificationConfiguration{
		SnsConfiguration: SNSConfigurationToProto(n.SNSConfiguration),
	}
}

// ProtoToNotificationConfiguration converts a protobuf NotificationConfiguration to its internal representation.
func ProtoToNotificationConfiguration(p *pb.NotificationConfiguration) *NotificationConfiguration {
	if p == nil {
		return nil
	}
	return &NotificationConfiguration{
		SNSConfiguration: ProtoToSNSConfiguration(p.SnsConfiguration),
	}
}

// S3ErrorReportConfiguration conversion

// S3ErrorReportConfigurationToProto converts an S3ErrorReportConfiguration to its protobuf representation.
func S3ErrorReportConfigurationToProto(s *S3ErrorReportConfiguration) *pb.S3ErrorReportConfiguration {
	if s == nil {
		return nil
	}
	return &pb.S3ErrorReportConfiguration{
		BucketName:      s.BucketName,
		ObjectKeyPrefix: s.ObjectKeyPrefix,
	}
}

// ProtoToS3ErrorReportConfiguration converts a protobuf S3ErrorReportConfiguration to its internal representation.
func ProtoToS3ErrorReportConfiguration(p *pb.S3ErrorReportConfiguration) *S3ErrorReportConfiguration {
	if p == nil {
		return nil
	}
	return &S3ErrorReportConfiguration{
		BucketName:      p.BucketName,
		ObjectKeyPrefix: p.ObjectKeyPrefix,
	}
}

// ErrorReportConfiguration conversion

// ErrorReportConfigurationToProto converts an ErrorReportConfiguration to its protobuf representation.
func ErrorReportConfigurationToProto(e *ErrorReportConfiguration) *pb.ErrorReportConfiguration {
	if e == nil {
		return nil
	}
	return &pb.ErrorReportConfiguration{
		S3Configuration: S3ErrorReportConfigurationToProto(e.S3Configuration),
	}
}

// ProtoToErrorReportConfiguration converts a protobuf ErrorReportConfiguration to its internal representation.
func ProtoToErrorReportConfiguration(p *pb.ErrorReportConfiguration) *ErrorReportConfiguration {
	if p == nil {
		return nil
	}
	return &ErrorReportConfiguration{
		S3Configuration: ProtoToS3ErrorReportConfiguration(p.S3Configuration),
	}
}

// SourceColumn conversion

// SourceColumnToProto converts a SourceColumn to its protobuf representation.
func SourceColumnToProto(s *SourceColumn) *pb.SourceColumn {
	if s == nil {
		return nil
	}
	return &pb.SourceColumn{
		Name: s.Name,
	}
}

// ProtoToSourceColumn converts a protobuf SourceColumn to its internal representation.
func ProtoToSourceColumn(p *pb.SourceColumn) *SourceColumn {
	if p == nil {
		return nil
	}
	return &SourceColumn{
		Name: p.Name,
	}
}

// DestinationColumn conversion

// DestinationColumnToProto converts a DestinationColumn to its protobuf representation.
func DestinationColumnToProto(d *DestinationColumn) *pb.DestinationColumn {
	if d == nil {
		return nil
	}
	return &pb.DestinationColumn{
		Name: d.Name,
	}
}

// ProtoToDestinationColumn converts a protobuf DestinationColumn to its internal representation.
func ProtoToDestinationColumn(p *pb.DestinationColumn) *DestinationColumn {
	if p == nil {
		return nil
	}
	return &DestinationColumn{
		Name: p.Name,
	}
}

// DimensionMapping conversion

// DimensionMappingToProto converts a DimensionMapping to its protobuf representation.
func DimensionMappingToProto(d *DimensionMapping) *pb.DimensionMapping {
	if d == nil {
		return nil
	}
	return &pb.DimensionMapping{
		SourceColumn:      SourceColumnToProto(d.SourceColumn),
		DestinationColumn: DestinationColumnToProto(d.DestinationColumn),
	}
}

// ProtoToDimensionMapping converts a protobuf DimensionMapping to its internal representation.
func ProtoToDimensionMapping(p *pb.DimensionMapping) *DimensionMapping {
	if p == nil {
		return nil
	}
	return &DimensionMapping{
		SourceColumn:      ProtoToSourceColumn(p.SourceColumn),
		DestinationColumn: ProtoToDestinationColumn(p.DestinationColumn),
	}
}

func dimensionMappingsToProto(ds []DimensionMapping) []*pb.DimensionMapping {
	if ds == nil {
		return nil
	}
	result := make([]*pb.DimensionMapping, len(ds))
	for i, d := range ds {
		result[i] = DimensionMappingToProto(&d)
	}
	return result
}

func protoToDimensionMappings(pbs []*pb.DimensionMapping) []DimensionMapping {
	if pbs == nil {
		return nil
	}
	result := make([]DimensionMapping, len(pbs))
	for i, pb := range pbs {
		result[i] = *ProtoToDimensionMapping(pb)
	}
	return result
}

// MultiMeasureAttributeMapping conversion

// MultiMeasureAttributeMappingToProto converts a MultiMeasureAttributeMapping to its protobuf representation.
func MultiMeasureAttributeMappingToProto(m *MultiMeasureAttributeMapping) *pb.MultiMeasureAttributeMapping {
	if m == nil {
		return nil
	}
	return &pb.MultiMeasureAttributeMapping{
		SourceColumn:                    SourceColumnToProto(m.SourceColumn),
		TargetMultiMeasureAttributeName: m.TargetMultiMeasureAttributeName,
		MeasureValueType:                string(m.MeasureValueMeasureValueType),
	}
}

// ProtoToMultiMeasureAttributeMapping converts a protobuf MultiMeasureAttributeMapping to its internal representation.
func ProtoToMultiMeasureAttributeMapping(p *pb.MultiMeasureAttributeMapping) *MultiMeasureAttributeMapping {
	if p == nil {
		return nil
	}
	return &MultiMeasureAttributeMapping{
		SourceColumn:                    ProtoToSourceColumn(p.SourceColumn),
		TargetMultiMeasureAttributeName: p.TargetMultiMeasureAttributeName,
		MeasureValueMeasureValueType:    MeasureValueType(p.MeasureValueType),
	}
}

func multiMeasureAttributeMappingsToProto(ms []MultiMeasureAttributeMapping) []*pb.MultiMeasureAttributeMapping {
	if ms == nil {
		return nil
	}
	result := make([]*pb.MultiMeasureAttributeMapping, len(ms))
	for i, m := range ms {
		result[i] = MultiMeasureAttributeMappingToProto(&m)
	}
	return result
}

func protoToMultiMeasureAttributeMappings(pbs []*pb.MultiMeasureAttributeMapping) []MultiMeasureAttributeMapping {
	if pbs == nil {
		return nil
	}
	result := make([]MultiMeasureAttributeMapping, len(pbs))
	for i, pb := range pbs {
		result[i] = *ProtoToMultiMeasureAttributeMapping(pb)
	}
	return result
}

// MultiMeasureMappings conversion

// MultiMeasureMappingsToProto converts a MultiMeasureMappings to its protobuf representation.
func MultiMeasureMappingsToProto(m *MultiMeasureMappings) *pb.MultiMeasureMappings {
	if m == nil {
		return nil
	}
	return &pb.MultiMeasureMappings{
		MultiMeasureAttributeMappings: multiMeasureAttributeMappingsToProto(m.MultiMeasureAttributeMappings),
	}
}

// ProtoToMultiMeasureMappings converts a protobuf MultiMeasureMappings to its internal representation.
func ProtoToMultiMeasureMappings(p *pb.MultiMeasureMappings) *MultiMeasureMappings {
	if p == nil {
		return nil
	}
	return &MultiMeasureMappings{
		MultiMeasureAttributeMappings: protoToMultiMeasureAttributeMappings(p.MultiMeasureAttributeMappings),
	}
}

// MixedMeasureMapping conversion

// MixedMeasureMappingToProto converts a MixedMeasureMapping to its protobuf representation.
func MixedMeasureMappingToProto(m *MixedMeasureMapping) *pb.MixedMeasureMapping {
	if m == nil {
		return nil
	}
	return &pb.MixedMeasureMapping{
		MeasureName:                   m.MeasureName,
		SourceColumn:                  m.SourceColumn,
		TargetMeasureName:             m.TargetMeasureName,
		MeasureValueType:              string(m.MeasureValueMeasureValueType),
		MultiMeasureAttributeMappings: multiMeasureAttributeMappingsToProto(m.MultiMeasureAttributeMappings),
	}
}

// ProtoToMixedMeasureMapping converts a protobuf MixedMeasureMapping to its internal representation.
func ProtoToMixedMeasureMapping(p *pb.MixedMeasureMapping) *MixedMeasureMapping {
	if p == nil {
		return nil
	}
	return &MixedMeasureMapping{
		MeasureName:                   p.MeasureName,
		SourceColumn:                  p.SourceColumn,
		TargetMeasureName:             p.TargetMeasureName,
		MeasureValueMeasureValueType:  MeasureValueType(p.MeasureValueType),
		MultiMeasureAttributeMappings: protoToMultiMeasureAttributeMappings(p.MultiMeasureAttributeMappings),
	}
}

func mixedMeasureMappingsToProto(ms []MixedMeasureMapping) []*pb.MixedMeasureMapping {
	if ms == nil {
		return nil
	}
	result := make([]*pb.MixedMeasureMapping, len(ms))
	for i, m := range ms {
		result[i] = MixedMeasureMappingToProto(&m)
	}
	return result
}

func protoToMixedMeasureMappings(pbs []*pb.MixedMeasureMapping) []MixedMeasureMapping {
	if pbs == nil {
		return nil
	}
	result := make([]MixedMeasureMapping, len(pbs))
	for i, pb := range pbs {
		result[i] = *ProtoToMixedMeasureMapping(pb)
	}
	return result
}

// TimestreamConfiguration conversion

// TimestreamConfigurationToProto converts a TimestreamConfiguration to its protobuf representation.
func TimestreamConfigurationToProto(t *TimestreamConfiguration) *pb.TimestreamConfiguration {
	if t == nil {
		return nil
	}
	return &pb.TimestreamConfiguration{
		DatabaseName:         t.DatabaseName,
		TableName:            t.TableName,
		TimeColumn:           t.TimeColumn,
		DimensionMappings:    dimensionMappingsToProto(t.DimensionMappings),
		MultiMeasureMappings: MultiMeasureMappingsToProto(t.MultiMeasureMappings),
		MixedMeasureMappings: mixedMeasureMappingsToProto(t.MixedMeasureMappings),
	}
}

// ProtoToTimestreamConfiguration converts a protobuf TimestreamConfiguration to its internal representation.
func ProtoToTimestreamConfiguration(p *pb.TimestreamConfiguration) *TimestreamConfiguration {
	if p == nil {
		return nil
	}
	return &TimestreamConfiguration{
		DatabaseName:         p.DatabaseName,
		TableName:            p.TableName,
		TimeColumn:           p.TimeColumn,
		DimensionMappings:    protoToDimensionMappings(p.DimensionMappings),
		MultiMeasureMappings: ProtoToMultiMeasureMappings(p.MultiMeasureMappings),
		MixedMeasureMappings: protoToMixedMeasureMappings(p.MixedMeasureMappings),
	}
}

// TargetConfiguration conversion

// TargetConfigurationToProto converts a TargetConfiguration to its protobuf representation.
func TargetConfigurationToProto(t *TargetConfiguration) *pb.TargetConfiguration {
	if t == nil {
		return nil
	}
	return &pb.TargetConfiguration{
		TimestreamConfiguration: TimestreamConfigurationToProto(t.TimestreamConfiguration),
	}
}

// ProtoToTargetConfiguration converts a protobuf TargetConfiguration to its internal representation.
func ProtoToTargetConfiguration(p *pb.TargetConfiguration) *TargetConfiguration {
	if p == nil {
		return nil
	}
	return &TargetConfiguration{
		TimestreamConfiguration: ProtoToTimestreamConfiguration(p.TimestreamConfiguration),
	}
}

// ScheduledQuery conversion

// ScheduledQueryToProto converts a ScheduledQuery to its protobuf representation.
func ScheduledQueryToProto(s *ScheduledQuery) *pb.ScheduledQuery {
	if s == nil {
		return nil
	}
	return &pb.ScheduledQuery{
		Arn:                            s.ARN,
		Name:                           s.Name,
		QueryString:                    s.QueryString,
		ScheduleConfiguration:          ScheduleConfigurationToProto(s.ScheduleConfiguration),
		NotificationConfiguration:      NotificationConfigurationToProto(s.NotificationConfiguration),
		ScheduledQueryExecutionRoleArn: s.ScheduledQueryExecutionRoleARN,
		KmsKeyId:                       s.KmsKeyID,
		ErrorReportConfiguration:       ErrorReportConfigurationToProto(s.ErrorReportConfiguration),
		TargetConfiguration:            TargetConfigurationToProto(s.TargetConfiguration),
		ClientToken:                    s.ClientToken,
		ScheduledQueryStatus:           string(s.ScheduledQueryStatus),
		CreationTime:                   timestamppb.New(s.CreationTime),
		LastRunStatus:                  s.LastRunStatus,
		PreviousRunTime:                timestamppb.New(s.PreviousRunTime),
		NextRunTime:                    timestamppb.New(s.NextRunTime),
	}
}

// ProtoToScheduledQuery converts a protobuf ScheduledQuery to its internal representation.
func ProtoToScheduledQuery(p *pb.ScheduledQuery) *ScheduledQuery {
	if p == nil {
		return nil
	}
	return &ScheduledQuery{
		ARN:                            p.Arn,
		Name:                           p.Name,
		QueryString:                    p.QueryString,
		ScheduleConfiguration:          ProtoToScheduleConfiguration(p.ScheduleConfiguration),
		NotificationConfiguration:      ProtoToNotificationConfiguration(p.NotificationConfiguration),
		ScheduledQueryExecutionRoleARN: p.ScheduledQueryExecutionRoleArn,
		KmsKeyID:                       p.KmsKeyId,
		ErrorReportConfiguration:       ProtoToErrorReportConfiguration(p.ErrorReportConfiguration),
		TargetConfiguration:            ProtoToTargetConfiguration(p.TargetConfiguration),
		ClientToken:                    p.ClientToken,
		ScheduledQueryStatus:           ScheduledQueryStatus(p.ScheduledQueryStatus),
		CreationTime:                   p.CreationTime.AsTime(),
		LastRunStatus:                  p.LastRunStatus,
		PreviousRunTime:                p.PreviousRunTime.AsTime(),
		NextRunTime:                    p.NextRunTime.AsTime(),
	}
}

// ExecutionStats conversion

// ExecutionStatsToProto converts an ExecutionStats to its protobuf representation.
func ExecutionStatsToProto(e *ExecutionStats) *pb.ExecutionStats {
	if e == nil {
		return nil
	}
	return &pb.ExecutionStats{
		TotalBytesProcessed:    e.TotalBytesProcessed,
		DataWrites:             e.DataWrites,
		TotalRecordsProcessed:  e.TotalRecordsProcessed,
		BytesMetered:           e.BytesMetered,
		QueryResultRows:        e.QueryResultRows,
		CumulativeBytesScanned: e.CumulativeBytesScanned,
		ExecutionTimeInMillis:  e.ExecutionTimeInMillis,
	}
}

// ProtoToExecutionStats converts a protobuf ExecutionStats to its internal representation.
func ProtoToExecutionStats(p *pb.ExecutionStats) *ExecutionStats {
	if p == nil {
		return nil
	}
	return &ExecutionStats{
		TotalBytesProcessed:    p.TotalBytesProcessed,
		DataWrites:             p.DataWrites,
		TotalRecordsProcessed:  p.TotalRecordsProcessed,
		BytesMetered:           p.BytesMetered,
		QueryResultRows:        p.QueryResultRows,
		CumulativeBytesScanned: p.CumulativeBytesScanned,
		ExecutionTimeInMillis:  p.ExecutionTimeInMillis,
	}
}

// ScheduledQueryRun conversion

// ScheduledQueryRunToProto converts a ScheduledQueryRun to its protobuf representation.
func ScheduledQueryRunToProto(s *ScheduledQueryRun) *pb.ScheduledQueryRun {
	if s == nil {
		return nil
	}
	return &pb.ScheduledQueryRun{
		Arn:               s.ARN,
		ScheduledQueryArn: s.ScheduledQueryARN,
		InvocationTime:    timestamppb.New(s.InvocationTime),
		TriggerTime:       timestamppb.New(s.TriggerTime),
		RunStatus:         string(s.RunStatus),
		ExecutionStats:    ExecutionStatsToProto(s.ExecutionStats),
		Error:             s.Error,
		CompletionTime:    timestamppb.New(s.CompletionTime),
	}
}

// ProtoToScheduledQueryRun converts a protobuf ScheduledQueryRun to its internal representation.
func ProtoToScheduledQueryRun(p *pb.ScheduledQueryRun) *ScheduledQueryRun {
	if p == nil {
		return nil
	}
	return &ScheduledQueryRun{
		ARN:               p.Arn,
		ScheduledQueryARN: p.ScheduledQueryArn,
		InvocationTime:    p.InvocationTime.AsTime(),
		TriggerTime:       p.TriggerTime.AsTime(),
		RunStatus:         ScheduleRunStatus(p.RunStatus),
		ExecutionStats:    ProtoToExecutionStats(p.ExecutionStats),
		Error:             p.Error,
		CompletionTime:    p.CompletionTime.AsTime(),
	}
}

// BatchLoadTask conversion

// BatchLoadTaskToProto converts a BatchLoadTask to its protobuf representation.
func BatchLoadTaskToProto(b *BatchLoadTask) *pb.BatchLoadTask {
	if b == nil {
		return nil
	}
	return &pb.BatchLoadTask{
		TaskId:          b.TaskId,
		DatabaseName:    b.DatabaseName,
		TableName:       b.TableName,
		TaskStatus:      string(b.TaskStatus),
		CreationTime:    timestamppb.New(b.CreationTime),
		LastUpdatedTime: timestamppb.New(b.LastUpdatedTime),
		ResumableUntil:  timestamppb.New(b.ResumableUntil),
	}
}

// ProtoToBatchLoadTask converts a protobuf BatchLoadTask to its internal representation.
func ProtoToBatchLoadTask(p *pb.BatchLoadTask) *BatchLoadTask {
	if p == nil {
		return nil
	}
	return &BatchLoadTask{
		TaskId:          p.TaskId,
		DatabaseName:    p.DatabaseName,
		TableName:       p.TableName,
		TaskStatus:      BatchLoadStatus(p.TaskStatus),
		CreationTime:    p.CreationTime.AsTime(),
		LastUpdatedTime: p.LastUpdatedTime.AsTime(),
		ResumableUntil:  p.ResumableUntil.AsTime(),
	}
}

// BatchLoadProgressReport conversion

// BatchLoadProgressReportToProto converts a BatchLoadProgressReport to its protobuf representation.
func BatchLoadProgressReportToProto(b *BatchLoadProgressReport) *pb.BatchLoadProgressReport {
	if b == nil {
		return nil
	}
	return &pb.BatchLoadProgressReport{
		BytesMetered:            b.BytesMetered,
		FileFailures:            b.FileFailures,
		ParseFailures:           b.ParseFailures,
		RecordIngestionFailures: b.RecordIngestionFailures,
		RecordsIngested:         b.RecordsIngested,
		RecordsProcessed:        b.RecordsProcessed,
	}
}

// ProtoToBatchLoadProgressReport converts a protobuf BatchLoadProgressReport to its internal representation.
func ProtoToBatchLoadProgressReport(p *pb.BatchLoadProgressReport) *BatchLoadProgressReport {
	if p == nil {
		return nil
	}
	return &BatchLoadProgressReport{
		BytesMetered:            p.BytesMetered,
		FileFailures:            p.FileFailures,
		ParseFailures:           p.ParseFailures,
		RecordIngestionFailures: p.RecordIngestionFailures,
		RecordsIngested:         p.RecordsIngested,
		RecordsProcessed:        p.RecordsProcessed,
	}
}

// DataSourceS3Configuration conversion

// DataSourceS3ConfigurationToProto converts a DataSourceS3Configuration to its protobuf representation.
func DataSourceS3ConfigurationToProto(d *DataSourceS3Configuration) *pb.DataSourceS3Configuration {
	if d == nil {
		return nil
	}
	return &pb.DataSourceS3Configuration{
		BucketName:      d.BucketName,
		ObjectKeyPrefix: d.ObjectKeyPrefix,
	}
}

// ProtoToDataSourceS3Configuration converts a protobuf DataSourceS3Configuration to its internal representation.
func ProtoToDataSourceS3Configuration(p *pb.DataSourceS3Configuration) *DataSourceS3Configuration {
	if p == nil {
		return nil
	}
	return &DataSourceS3Configuration{
		BucketName:      p.BucketName,
		ObjectKeyPrefix: p.ObjectKeyPrefix,
	}
}

// CsvConfiguration conversion

// CsvConfigurationToProto converts a CsvConfiguration to its protobuf representation.
func CsvConfigurationToProto(c *CsvConfiguration) *pb.CsvConfiguration {
	if c == nil {
		return nil
	}
	return &pb.CsvConfiguration{
		ColumnSeparator: c.ColumnSeparator,
		EscapeChar:      c.EscapeChar,
		NullValue:       c.NullValue,
		QuoteChar:       c.QuoteChar,
		TrimWhiteSpace:  c.TrimWhiteSpace,
	}
}

// ProtoToCsvConfiguration converts a protobuf CsvConfiguration to its internal representation.
func ProtoToCsvConfiguration(p *pb.CsvConfiguration) *CsvConfiguration {
	if p == nil {
		return nil
	}
	return &CsvConfiguration{
		ColumnSeparator: p.ColumnSeparator,
		EscapeChar:      p.EscapeChar,
		NullValue:       p.NullValue,
		QuoteChar:       p.QuoteChar,
		TrimWhiteSpace:  p.TrimWhiteSpace,
	}
}

// DataSourceConfiguration conversion

// DataSourceConfigurationToProto converts a DataSourceConfiguration to its protobuf representation.
func DataSourceConfigurationToProto(d *DataSourceConfiguration) *pb.DataSourceConfiguration {
	if d == nil {
		return nil
	}
	return &pb.DataSourceConfiguration{
		DataFormat:                string(d.DataFormat),
		DataSourceS3Configuration: DataSourceS3ConfigurationToProto(d.DataSourceS3Configuration),
		CsvConfiguration:          CsvConfigurationToProto(d.CsvConfiguration),
	}
}

// ProtoToDataSourceConfiguration converts a protobuf DataSourceConfiguration to its internal representation.
func ProtoToDataSourceConfiguration(p *pb.DataSourceConfiguration) *DataSourceConfiguration {
	if p == nil {
		return nil
	}
	return &DataSourceConfiguration{
		DataFormat:                BatchLoadDataFormat(p.DataFormat),
		DataSourceS3Configuration: ProtoToDataSourceS3Configuration(p.DataSourceS3Configuration),
		CsvConfiguration:          ProtoToCsvConfiguration(p.CsvConfiguration),
	}
}

// DataModelS3Configuration conversion

// DataModelS3ConfigurationToProto converts a DataModelS3Configuration to its protobuf representation.
func DataModelS3ConfigurationToProto(d *DataModelS3Configuration) *pb.DataModelS3Configuration {
	if d == nil {
		return nil
	}
	return &pb.DataModelS3Configuration{
		BucketName: d.BucketName,
		ObjectKey:  d.ObjectKey,
	}
}

// ProtoToDataModelS3Configuration converts a protobuf DataModelS3Configuration to its internal representation.
func ProtoToDataModelS3Configuration(p *pb.DataModelS3Configuration) *DataModelS3Configuration {
	if p == nil {
		return nil
	}
	return &DataModelS3Configuration{
		BucketName: p.BucketName,
		ObjectKey:  p.ObjectKey,
	}
}

// DataModel conversion

// DataModelToProto converts a DataModel to its protobuf representation.
func DataModelToProto(d *DataModel) *pb.DataModel {
	if d == nil {
		return nil
	}
	return &pb.DataModel{
		DimensionMappings:    dimensionMappingsToProto(d.DimensionMappings),
		MeasureNameColumn:    d.MeasureNameColumn,
		MixedMeasureMappings: mixedMeasureMappingsToProto(d.MixedMeasureMappings),
		MultiMeasureMappings: MultiMeasureMappingsToProto(d.MultiMeasureMappings),
		TimeColumn:           d.TimeColumn,
		TimeUnit:             string(d.TimeUnit),
	}
}

// ProtoToDataModel converts a protobuf DataModel to its internal representation.
func ProtoToDataModel(p *pb.DataModel) *DataModel {
	if p == nil {
		return nil
	}
	return &DataModel{
		DimensionMappings:    protoToDimensionMappings(p.DimensionMappings),
		MeasureNameColumn:    p.MeasureNameColumn,
		MixedMeasureMappings: protoToMixedMeasureMappings(p.MixedMeasureMappings),
		MultiMeasureMappings: ProtoToMultiMeasureMappings(p.MultiMeasureMappings),
		TimeColumn:           p.TimeColumn,
		TimeUnit:             TimeUnit(p.TimeUnit),
	}
}

// DataModelConfiguration conversion

// DataModelConfigurationToProto converts a DataModelConfiguration to its protobuf representation.
func DataModelConfigurationToProto(d *DataModelConfiguration) *pb.DataModelConfiguration {
	if d == nil {
		return nil
	}
	return &pb.DataModelConfiguration{
		DataModel:                DataModelToProto(d.DataModel),
		DataModelS3Configuration: DataModelS3ConfigurationToProto(d.DataModelS3Configuration),
	}
}

// ProtoToDataModelConfiguration converts a protobuf DataModelConfiguration to its internal representation.
func ProtoToDataModelConfiguration(p *pb.DataModelConfiguration) *DataModelConfiguration {
	if p == nil {
		return nil
	}
	return &DataModelConfiguration{
		DataModel:                ProtoToDataModel(p.DataModel),
		DataModelS3Configuration: ProtoToDataModelS3Configuration(p.DataModelS3Configuration),
	}
}

// ReportS3Configuration conversion

// ReportS3ConfigurationToProto converts a ReportS3Configuration to its protobuf representation.
func ReportS3ConfigurationToProto(r *ReportS3Configuration) *pb.ReportS3Configuration {
	if r == nil {
		return nil
	}
	return &pb.ReportS3Configuration{
		BucketName:       r.BucketName,
		EncryptionOption: string(r.EncryptionOption),
		KmsKeyId:         r.KmsKeyId,
		ObjectKeyPrefix:  r.ObjectKeyPrefix,
	}
}

// ProtoToReportS3Configuration converts a protobuf ReportS3Configuration to its internal representation.
func ProtoToReportS3Configuration(p *pb.ReportS3Configuration) *ReportS3Configuration {
	if p == nil {
		return nil
	}
	return &ReportS3Configuration{
		BucketName:       p.BucketName,
		EncryptionOption: S3EncryptionOption(p.EncryptionOption),
		KmsKeyId:         p.KmsKeyId,
		ObjectKeyPrefix:  p.ObjectKeyPrefix,
	}
}

// ReportConfiguration conversion

// ReportConfigurationToProto converts a ReportConfiguration to its protobuf representation.
func ReportConfigurationToProto(r *ReportConfiguration) *pb.ReportConfiguration {
	if r == nil {
		return nil
	}
	return &pb.ReportConfiguration{
		ReportS3Configuration: ReportS3ConfigurationToProto(r.ReportS3Configuration),
	}
}

// ProtoToReportConfiguration converts a protobuf ReportConfiguration to its internal representation.
func ProtoToReportConfiguration(p *pb.ReportConfiguration) *ReportConfiguration {
	if p == nil {
		return nil
	}
	return &ReportConfiguration{
		ReportS3Configuration: ProtoToReportS3Configuration(p.ReportS3Configuration),
	}
}

// BatchLoadTaskDescription conversion

// BatchLoadTaskDescriptionToProto converts a BatchLoadTaskDescription to its protobuf representation.
func BatchLoadTaskDescriptionToProto(b *BatchLoadTaskDescription) *pb.BatchLoadTaskDescription {
	if b == nil {
		return nil
	}
	return &pb.BatchLoadTaskDescription{
		TaskId:                  b.TaskId,
		TargetDatabaseName:      b.TargetDatabaseName,
		TargetTableName:         b.TargetTableName,
		TaskStatus:              string(b.TaskStatus),
		CreationTime:            timestamppb.New(b.CreationTime),
		LastUpdatedTime:         timestamppb.New(b.LastUpdatedTime),
		ResumableUntil:          timestamppb.New(b.ResumableUntil),
		ErrorMessage:            b.ErrorMessage,
		RecordVersion:           b.RecordVersion,
		DataSourceConfiguration: DataSourceConfigurationToProto(b.DataSourceConfiguration),
		DataModelConfiguration:  DataModelConfigurationToProto(b.DataModelConfiguration),
		ReportConfiguration:     ReportConfigurationToProto(b.ReportConfiguration),
		ProgressReport:          BatchLoadProgressReportToProto(b.ProgressReport),
	}
}

// ProtoToBatchLoadTaskDescription converts a protobuf BatchLoadTaskDescription to its internal representation.
func ProtoToBatchLoadTaskDescription(p *pb.BatchLoadTaskDescription) *BatchLoadTaskDescription {
	if p == nil {
		return nil
	}
	return &BatchLoadTaskDescription{
		TaskId:                  p.TaskId,
		TargetDatabaseName:      p.TargetDatabaseName,
		TargetTableName:         p.TargetTableName,
		TaskStatus:              BatchLoadStatus(p.TaskStatus),
		CreationTime:            p.CreationTime.AsTime(),
		LastUpdatedTime:         p.LastUpdatedTime.AsTime(),
		ResumableUntil:          p.ResumableUntil.AsTime(),
		ErrorMessage:            p.ErrorMessage,
		RecordVersion:           p.RecordVersion,
		DataSourceConfiguration: ProtoToDataSourceConfiguration(p.DataSourceConfiguration),
		DataModelConfiguration:  ProtoToDataModelConfiguration(p.DataModelConfiguration),
		ReportConfiguration:     ProtoToReportConfiguration(p.ReportConfiguration),
		ProgressReport:          ProtoToBatchLoadProgressReport(p.ProgressReport),
	}
}

// AccountSettings conversion

// AccountSettingsToProto converts AccountSettings to its protobuf representation.
func AccountSettingsToProto(a *AccountSettings) *pb.AccountSettings {
	if a == nil {
		return nil
	}
	return &pb.AccountSettings{
		MaxQueryTcu:             a.MaxQueryTCU,
		QueryPricingMode:        a.QueryPricingMode,
		EncryptionConfiguration: a.EncryptionConfiguration,
		QueryComputeType:        a.QueryComputeType,
	}
}

// ProtoToAccountSettings converts protobuf AccountSettings to its internal representation.
func ProtoToAccountSettings(p *pb.AccountSettings) *AccountSettings {
	if p == nil {
		return nil
	}
	return &AccountSettings{
		MaxQueryTCU:             p.MaxQueryTcu,
		QueryPricingMode:        p.QueryPricingMode,
		EncryptionConfiguration: p.EncryptionConfiguration,
		QueryComputeType:        p.QueryComputeType,
	}
}

// Helper functions for zero time checks
