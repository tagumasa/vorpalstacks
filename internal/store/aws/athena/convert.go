package athena

import (
	"encoding/json"
	"log"

	"google.golang.org/protobuf/types/known/timestamppb"
	pb "vorpalstacks/internal/pb/storage/storage_athena"
)

// WorkGroupToProto converts a WorkGroup to its protobuf representation.
func WorkGroupToProto(wg *WorkGroup) *pb.WorkGroup {
	if wg == nil {
		return nil
	}
	p := &pb.WorkGroup{
		Name:          wg.Name,
		Description:   wg.Description,
		Configuration: WorkGroupConfigurationToProto(wg.Configuration),
		State:         string(wg.State),
		CreatedTime:   timestamppb.New(wg.CreatedTime),
	}
	return p
}

// ProtoToWorkGroup converts a protobuf WorkGroup to its internal representation.
func ProtoToWorkGroup(p *pb.WorkGroup) *WorkGroup {
	if p == nil {
		return nil
	}
	wg := &WorkGroup{
		Name:          p.Name,
		Description:   p.Description,
		Configuration: ProtoToWorkGroupConfiguration(p.Configuration),
		State:         WorkGroupState(p.State),
	}
	if p.CreatedTime != nil {
		wg.CreatedTime = p.CreatedTime.AsTime()
	}
	return wg
}

// WorkGroupConfigurationToProto converts a WorkGroupConfiguration to its protobuf representation.
func WorkGroupConfigurationToProto(c *WorkGroupConfiguration) *pb.WorkGroupConfiguration {
	if c == nil {
		return nil
	}
	return &pb.WorkGroupConfiguration{
		ResultConfiguration:                    ResultConfigurationToProto(c.ResultConfiguration),
		EnforceWorkGroupConfiguration:          c.EnforceWorkGroupConfiguration,
		PublishCloudWatchMetricsEnabled:        c.PublishCloudWatchMetricsEnabled,
		BytesScannedCutoffPerQuery:             c.BytesScannedCutoffPerQuery,
		RequesterPaysEnabled:                   c.RequesterPaysEnabled,
		EngineVersion:                          EngineVersionToProto(c.EngineVersion),
		AdditionalConfiguration:                c.AdditionalConfiguration,
		ExecutionRole:                          c.ExecutionRole,
		CustomerContentEncryptionConfiguration: CustomerContentEncryptionConfigurationToProto(c.CustomerContentEncryptionConfiguration),
		EnableMinimumEncryptionConfiguration:   c.EnableMinimumEncryptionConfiguration,
	}
}

// ProtoToWorkGroupConfiguration converts a protobuf WorkGroupConfiguration to its internal representation.
func ProtoToWorkGroupConfiguration(p *pb.WorkGroupConfiguration) *WorkGroupConfiguration {
	if p == nil {
		return nil
	}
	return &WorkGroupConfiguration{
		ResultConfiguration:                    ProtoToResultConfiguration(p.ResultConfiguration),
		EnforceWorkGroupConfiguration:          p.EnforceWorkGroupConfiguration,
		PublishCloudWatchMetricsEnabled:        p.PublishCloudWatchMetricsEnabled,
		BytesScannedCutoffPerQuery:             p.BytesScannedCutoffPerQuery,
		RequesterPaysEnabled:                   p.RequesterPaysEnabled,
		EngineVersion:                          ProtoToEngineVersion(p.EngineVersion),
		AdditionalConfiguration:                p.AdditionalConfiguration,
		ExecutionRole:                          p.ExecutionRole,
		CustomerContentEncryptionConfiguration: ProtoToCustomerContentEncryptionConfiguration(p.CustomerContentEncryptionConfiguration),
		EnableMinimumEncryptionConfiguration:   p.EnableMinimumEncryptionConfiguration,
	}
}

// ResultConfigurationToProto converts a ResultConfiguration to its protobuf representation.
func ResultConfigurationToProto(c *ResultConfiguration) *pb.ResultConfiguration {
	if c == nil {
		return nil
	}
	return &pb.ResultConfiguration{
		OutputLocation:          c.OutputLocation,
		EncryptionConfiguration: EncryptionConfigurationToProto(c.EncryptionConfiguration),
		ExpectedBucketOwner:     c.ExpectedBucketOwner,
		AclConfiguration:        ACLConfigurationToProto(c.ACLConfiguration),
	}
}

// ProtoToResultConfiguration converts a protobuf ResultConfiguration to its internal representation.
func ProtoToResultConfiguration(p *pb.ResultConfiguration) *ResultConfiguration {
	if p == nil {
		return nil
	}
	return &ResultConfiguration{
		OutputLocation:          p.OutputLocation,
		EncryptionConfiguration: ProtoToEncryptionConfiguration(p.EncryptionConfiguration),
		ExpectedBucketOwner:     p.ExpectedBucketOwner,
		ACLConfiguration:        ProtoToACLConfiguration(p.AclConfiguration),
	}
}

// EncryptionConfigurationToProto converts an EncryptionConfiguration to its protobuf representation.
func EncryptionConfigurationToProto(c *EncryptionConfiguration) *pb.EncryptionConfiguration {
	if c == nil {
		return nil
	}
	return &pb.EncryptionConfiguration{
		EncryptionOption: c.EncryptionOption,
		KmsKey:           c.KmsKey,
	}
}

// ProtoToEncryptionConfiguration converts a protobuf EncryptionConfiguration to its internal representation.
func ProtoToEncryptionConfiguration(p *pb.EncryptionConfiguration) *EncryptionConfiguration {
	if p == nil {
		return nil
	}
	return &EncryptionConfiguration{
		EncryptionOption: p.EncryptionOption,
		KmsKey:           p.KmsKey,
	}
}

// ACLConfigurationToProto converts an ACLConfiguration to its protobuf representation.
func ACLConfigurationToProto(c *ACLConfiguration) *pb.ACLConfiguration {
	if c == nil {
		return nil
	}
	return &pb.ACLConfiguration{
		S3AclOption: c.S3ACLOption,
	}
}

// ProtoToACLConfiguration converts a protobuf ACLConfiguration to its internal representation.
func ProtoToACLConfiguration(p *pb.ACLConfiguration) *ACLConfiguration {
	if p == nil {
		return nil
	}
	return &ACLConfiguration{
		S3ACLOption: p.S3AclOption,
	}
}

// EngineVersionToProto converts an EngineVersion to its protobuf representation.
func EngineVersionToProto(e *EngineVersion) *pb.EngineVersion {
	if e == nil {
		return nil
	}
	return &pb.EngineVersion{
		SelectedEngineVersion:  e.SelectedEngineVersion,
		EffectiveEngineVersion: e.EffectiveEngineVersion,
	}
}

// ProtoToEngineVersion converts a protobuf EngineVersion to its internal representation.
func ProtoToEngineVersion(p *pb.EngineVersion) *EngineVersion {
	if p == nil {
		return nil
	}
	return &EngineVersion{
		SelectedEngineVersion:  p.SelectedEngineVersion,
		EffectiveEngineVersion: p.EffectiveEngineVersion,
	}
}

// CustomerContentEncryptionConfigurationToProto converts a CustomerContentEncryptionConfiguration to its protobuf representation.
func CustomerContentEncryptionConfigurationToProto(c *CustomerContentEncryptionConfiguration) *pb.CustomerContentEncryptionConfiguration {
	if c == nil {
		return nil
	}
	return &pb.CustomerContentEncryptionConfiguration{
		KmsKey: c.KmsKey,
	}
}

// ProtoToCustomerContentEncryptionConfiguration converts a protobuf CustomerContentEncryptionConfiguration to its internal representation.
func ProtoToCustomerContentEncryptionConfiguration(p *pb.CustomerContentEncryptionConfiguration) *CustomerContentEncryptionConfiguration {
	if p == nil {
		return nil
	}
	return &CustomerContentEncryptionConfiguration{
		KmsKey: p.KmsKey,
	}
}

// QueryExecutionToProto converts a QueryExecution to its protobuf representation.
func QueryExecutionToProto(qe *QueryExecution) *pb.QueryExecution {
	if qe == nil {
		return nil
	}
	return &pb.QueryExecution{
		QueryExecutionId:      qe.QueryExecutionId,
		Query:                 qe.Query,
		StatementType:         string(qe.StatementType),
		WorkGroup:             qe.WorkGroup,
		QueryExecutionContext: QueryExecutionContextToProto(qe.QueryExecutionContext),
		ResultConfiguration:   ResultConfigurationToProto(qe.ResultConfiguration),
		Status:                QueryExecutionStatusToProto(qe.Status),
		Statistics:            QueryExecutionStatisticsToProto(qe.Statistics),
		SubstatementType:      qe.SubstatementType,
	}
}

// ProtoToQueryExecution converts a protobuf QueryExecution to its internal representation.
func ProtoToQueryExecution(p *pb.QueryExecution) *QueryExecution {
	if p == nil {
		return nil
	}
	return &QueryExecution{
		QueryExecutionId:      p.QueryExecutionId,
		Query:                 p.Query,
		StatementType:         StatementType(p.StatementType),
		WorkGroup:             p.WorkGroup,
		QueryExecutionContext: ProtoToQueryExecutionContext(p.QueryExecutionContext),
		ResultConfiguration:   ProtoToResultConfiguration(p.ResultConfiguration),
		Status:                ProtoToQueryExecutionStatus(p.Status),
		Statistics:            ProtoToQueryExecutionStatistics(p.Statistics),
		SubstatementType:      p.SubstatementType,
	}
}

// QueryExecutionContextToProto converts a QueryExecutionContext to its protobuf representation.
func QueryExecutionContextToProto(c *QueryExecutionContext) *pb.QueryExecutionContext {
	if c == nil {
		return nil
	}
	return &pb.QueryExecutionContext{
		Database: c.Database,
		Catalog:  c.Catalog,
	}
}

// ProtoToQueryExecutionContext converts a protobuf QueryExecutionContext to its internal representation.
func ProtoToQueryExecutionContext(p *pb.QueryExecutionContext) *QueryExecutionContext {
	if p == nil {
		return nil
	}
	return &QueryExecutionContext{
		Database: p.Database,
		Catalog:  p.Catalog,
	}
}

// QueryExecutionStatusToProto converts a QueryExecutionStatus to its protobuf representation.
func QueryExecutionStatusToProto(s *QueryExecutionStatus) *pb.QueryExecutionStatus {
	if s == nil {
		return nil
	}
	p := &pb.QueryExecutionStatus{
		State:              string(s.State),
		StateChangeReason:  s.StateChangeReason,
		SubmissionDateTime: timestamppb.New(s.SubmissionDateTime),
		AthenaError:        AthenaErrorToProto(s.AthenaError),
	}
	if !s.CompletionDateTime.IsZero() {
		p.CompletionDateTime = timestamppb.New(s.CompletionDateTime)
	}
	return p
}

// ProtoToQueryExecutionStatus converts a protobuf QueryExecutionStatus to its internal representation.
func ProtoToQueryExecutionStatus(p *pb.QueryExecutionStatus) *QueryExecutionStatus {
	if p == nil {
		return nil
	}
	s := &QueryExecutionStatus{
		State:             QueryExecutionState(p.State),
		StateChangeReason: p.StateChangeReason,
		AthenaError:       ProtoToAthenaError(p.AthenaError),
	}
	if p.SubmissionDateTime != nil {
		s.SubmissionDateTime = p.SubmissionDateTime.AsTime()
	}
	if p.CompletionDateTime != nil {
		s.CompletionDateTime = p.CompletionDateTime.AsTime()
	}
	return s
}

// AthenaErrorToProto converts an AthenaError to its protobuf representation.
func AthenaErrorToProto(e *AthenaError) *pb.AthenaError {
	if e == nil {
		return nil
	}
	return &pb.AthenaError{
		ErrorCategory:     int32(e.ErrorCategory),
		ErrorType:         e.ErrorType,
		Retryable:         e.Retryable,
		ErrorMessage:      e.ErrorMessage,
		SyntaxErrorRow:    int32(e.SyntaxErrorRow),
		SyntaxErrorColumn: int32(e.SyntaxErrorColumn),
	}
}

// ProtoToAthenaError converts a protobuf AthenaError to its internal representation.
func ProtoToAthenaError(p *pb.AthenaError) *AthenaError {
	if p == nil {
		return nil
	}
	return &AthenaError{
		ErrorCategory:     int(p.ErrorCategory),
		ErrorType:         p.ErrorType,
		Retryable:         p.Retryable,
		ErrorMessage:      p.ErrorMessage,
		SyntaxErrorRow:    int(p.SyntaxErrorRow),
		SyntaxErrorColumn: int(p.SyntaxErrorColumn),
	}
}

// QueryExecutionStatisticsToProto converts a QueryExecutionStatistics to its protobuf representation.
func QueryExecutionStatisticsToProto(s *QueryExecutionStatistics) *pb.QueryExecutionStatistics {
	if s == nil {
		return nil
	}
	return &pb.QueryExecutionStatistics{
		EngineExecutionTimeInMillis:   s.EngineExecutionTimeInMillis,
		DataScannedInBytes:            s.DataScannedInBytes,
		DataManifestLocation:          s.DataManifestLocation,
		TotalExecutionTimeInMillis:    s.TotalExecutionTimeInMillis,
		QueryQueueTimeInMillis:        s.QueryQueueTimeInMillis,
		QueryPlanningTimeInMillis:     s.QueryPlanningTimeInMillis,
		ServiceProcessingTimeInMillis: s.ServiceProcessingTimeInMillis,
		ResultReuseInformation:        ResultReuseInformationToProto(s.ResultReuseInformation),
	}
}

// ProtoToQueryExecutionStatistics converts a protobuf QueryExecutionStatistics to its internal representation.
func ProtoToQueryExecutionStatistics(p *pb.QueryExecutionStatistics) *QueryExecutionStatistics {
	if p == nil {
		return nil
	}
	return &QueryExecutionStatistics{
		EngineExecutionTimeInMillis:   p.EngineExecutionTimeInMillis,
		DataScannedInBytes:            p.DataScannedInBytes,
		DataManifestLocation:          p.DataManifestLocation,
		TotalExecutionTimeInMillis:    p.TotalExecutionTimeInMillis,
		QueryQueueTimeInMillis:        p.QueryQueueTimeInMillis,
		QueryPlanningTimeInMillis:     p.QueryPlanningTimeInMillis,
		ServiceProcessingTimeInMillis: p.ServiceProcessingTimeInMillis,
		ResultReuseInformation:        ProtoToResultReuseInformation(p.ResultReuseInformation),
	}
}

// ResultReuseInformationToProto converts a ResultReuseInformation to its protobuf representation.
func ResultReuseInformationToProto(r *ResultReuseInformation) *pb.ResultReuseInformation {
	if r == nil {
		return nil
	}
	return &pb.ResultReuseInformation{
		ReusedPreviousResult: r.ReusedPreviousResult,
	}
}

// ProtoToResultReuseInformation converts a protobuf ResultReuseInformation to its internal representation.
func ProtoToResultReuseInformation(p *pb.ResultReuseInformation) *ResultReuseInformation {
	if p == nil {
		return nil
	}
	return &ResultReuseInformation{
		ReusedPreviousResult: p.ReusedPreviousResult,
	}
}

// QueryResultToProto converts a QueryResult to its protobuf representation.
func QueryResultToProto(qr *QueryResult) *pb.QueryResult {
	if qr == nil {
		return nil
	}
	return &pb.QueryResult{
		QueryExecutionId: qr.QueryExecutionId,
		ResultSet:        ResultSetToProto(qr.ResultSet),
		NextToken:        qr.NextToken,
		UpdateCount:      qr.UpdateCount,
	}
}

// ProtoToQueryResult converts a protobuf QueryResult to its internal representation.
func ProtoToQueryResult(p *pb.QueryResult) *QueryResult {
	if p == nil {
		return nil
	}
	return &QueryResult{
		QueryExecutionId: p.QueryExecutionId,
		ResultSet:        ProtoToResultSet(p.ResultSet),
		NextToken:        p.NextToken,
		UpdateCount:      p.UpdateCount,
	}
}

// ResultSetToProto converts a ResultSet to its protobuf representation.
func ResultSetToProto(rs *ResultSet) *pb.ResultSet {
	if rs == nil {
		return nil
	}
	rows := make([]*pb.Row, len(rs.Rows))
	for i, r := range rs.Rows {
		rows[i] = RowToProto(&r)
	}
	return &pb.ResultSet{
		Rows:              rows,
		ResultSetMetadata: ResultSetMetadataToProto(rs.ResultSetMetadata),
	}
}

// ProtoToResultSet converts a protobuf ResultSet to its internal representation.
func ProtoToResultSet(p *pb.ResultSet) *ResultSet {
	if p == nil {
		return nil
	}
	rows := make([]Row, len(p.Rows))
	for i, r := range p.Rows {
		rows[i] = ProtoToRow(r)
	}
	return &ResultSet{
		Rows:              rows,
		ResultSetMetadata: ProtoToResultSetMetadata(p.ResultSetMetadata),
	}
}

// RowToProto converts a Row to its protobuf representation.
func RowToProto(r *Row) *pb.Row {
	if r == nil {
		return nil
	}
	data := make([]*pb.Datum, len(r.Data))
	for i, d := range r.Data {
		data[i] = DatumToProto(d)
	}
	return &pb.Row{Data: data}
}

// ProtoToRow converts a protobuf Row to its internal representation.
func ProtoToRow(p *pb.Row) Row {
	if p == nil {
		return Row{}
	}
	data := make([]Datum, len(p.Data))
	for i, d := range p.Data {
		data[i] = ProtoToDatum(d)
	}
	return Row{Data: data}
}

// DatumToProto converts a Datum to its protobuf representation.
func DatumToProto(d Datum) *pb.Datum {
	return &pb.Datum{
		VarCharValue: d.VarCharValue,
	}
}

// ProtoToDatum converts a protobuf Datum to its internal representation.
func ProtoToDatum(p *pb.Datum) Datum {
	return Datum{
		VarCharValue: p.VarCharValue,
	}
}

// ResultSetMetadataToProto converts a ResultSetMetadata to its protobuf representation.
func ResultSetMetadataToProto(m *ResultSetMetadata) *pb.ResultSetMetadata {
	if m == nil {
		return nil
	}
	cols := make([]*pb.ColumnInfo, len(m.ColumnInfo))
	for i, c := range m.ColumnInfo {
		cols[i] = ColumnInfoToProto(c)
	}
	return &pb.ResultSetMetadata{ColumnInfo: cols}
}

// ProtoToResultSetMetadata converts a protobuf ResultSetMetadata to its internal representation.
func ProtoToResultSetMetadata(p *pb.ResultSetMetadata) *ResultSetMetadata {
	if p == nil {
		return nil
	}
	cols := make([]ColumnInfo, len(p.ColumnInfo))
	for i, c := range p.ColumnInfo {
		cols[i] = ProtoToColumnInfo(c)
	}
	return &ResultSetMetadata{ColumnInfo: cols}
}

// ColumnInfoToProto converts a ColumnInfo to its protobuf representation.
func ColumnInfoToProto(c ColumnInfo) *pb.ColumnInfo {
	return &pb.ColumnInfo{
		Label:         c.Label,
		Name:          c.Name,
		Type:          c.Type,
		Precision:     c.Precision,
		Scale:         c.Scale,
		Nullable:      c.Nullable,
		CaseSensitive: c.CaseSensitive,
		SchemaName:    c.SchemaName,
		TableName:     c.TableName,
		CatalogName:   c.CatalogName,
	}
}

// ProtoToColumnInfo converts a protobuf ColumnInfo to its internal representation.
func ProtoToColumnInfo(p *pb.ColumnInfo) ColumnInfo {
	return ColumnInfo{
		Label:         p.Label,
		Name:          p.Name,
		Type:          p.Type,
		Precision:     p.Precision,
		Scale:         p.Scale,
		Nullable:      p.Nullable,
		CaseSensitive: p.CaseSensitive,
		SchemaName:    p.SchemaName,
		TableName:     p.TableName,
		CatalogName:   p.CatalogName,
	}
}

// NamedQueryToProto converts a NamedQuery to its protobuf representation.
func NamedQueryToProto(nq *NamedQuery) *pb.NamedQuery {
	if nq == nil {
		return nil
	}
	return &pb.NamedQuery{
		NamedQueryId: nq.NamedQueryId,
		Name:         nq.Name,
		Description:  nq.Description,
		Database:     nq.Database,
		QueryString:  nq.QueryString,
		WorkGroup:    nq.WorkGroup,
		CreatedTime:  timestamppb.New(nq.CreatedTime),
	}
}

// ProtoToNamedQuery converts a protobuf NamedQuery to its internal representation.
func ProtoToNamedQuery(p *pb.NamedQuery) *NamedQuery {
	if p == nil {
		return nil
	}
	nq := &NamedQuery{
		NamedQueryId: p.NamedQueryId,
		Name:         p.Name,
		Description:  p.Description,
		Database:     p.Database,
		QueryString:  p.QueryString,
		WorkGroup:    p.WorkGroup,
	}
	if p.CreatedTime != nil {
		nq.CreatedTime = p.CreatedTime.AsTime()
	}
	return nq
}

// PreparedStatementToProto converts a PreparedStatement to its protobuf representation.
func PreparedStatementToProto(ps *PreparedStatement) *pb.PreparedStatement {
	if ps == nil {
		return nil
	}
	return &pb.PreparedStatement{
		StatementName:    ps.StatementName,
		WorkGroupName:    ps.WorkGroupName,
		QueryStatement:   ps.QueryStatement,
		Description:      ps.Description,
		LastModifiedTime: timestamppb.New(ps.LastModifiedTime),
		CreatedTime:      timestamppb.New(ps.CreatedTime),
	}
}

// ProtoToPreparedStatement converts a protobuf PreparedStatement to its internal representation.
func ProtoToPreparedStatement(p *pb.PreparedStatement) *PreparedStatement {
	if p == nil {
		return nil
	}
	ps := &PreparedStatement{
		StatementName:  p.StatementName,
		WorkGroupName:  p.WorkGroupName,
		QueryStatement: p.QueryStatement,
		Description:    p.Description,
	}
	if p.LastModifiedTime != nil {
		ps.LastModifiedTime = p.LastModifiedTime.AsTime()
	}
	if p.CreatedTime != nil {
		ps.CreatedTime = p.CreatedTime.AsTime()
	}
	return ps
}

// DataCatalogToProto converts a DataCatalog to its protobuf representation.
func DataCatalogToProto(dc *DataCatalog) *pb.DataCatalog {
	if dc == nil {
		return nil
	}
	return &pb.DataCatalog{
		Name:        dc.Name,
		Description: dc.Description,
		Type:        dc.Type,
		Parameters:  dc.Parameters,
	}
}

// ProtoToDataCatalog converts a protobuf DataCatalog to its internal representation.
func ProtoToDataCatalog(p *pb.DataCatalog) *DataCatalog {
	if p == nil {
		return nil
	}
	return &DataCatalog{
		Name:        p.Name,
		Description: p.Description,
		Type:        p.Type,
		Parameters:  p.Parameters,
	}
}

// DatabaseToProto converts a Database to its protobuf representation.
func DatabaseToProto(db *Database) *pb.Database {
	if db == nil {
		return nil
	}
	return &pb.Database{
		Name:        db.Name,
		Description: db.Description,
		Parameters:  db.Parameters,
	}
}

// ProtoToDatabase converts a protobuf Database to its internal representation.
func ProtoToDatabase(p *pb.Database) *Database {
	if p == nil {
		return nil
	}
	return &Database{
		Name:        p.Name,
		Description: p.Description,
		Parameters:  p.Parameters,
	}
}

// TableMetadataToProto converts a TableMetadata to its protobuf representation.
func TableMetadataToProto(tm *TableMetadata) *pb.TableMetadata {
	if tm == nil {
		return nil
	}
	columns := make([]*pb.Column, len(tm.Columns))
	for i, c := range tm.Columns {
		columns[i] = ColumnToProto(c)
	}
	partitionKeys := make([]*pb.Column, len(tm.PartitionKeys))
	for i, c := range tm.PartitionKeys {
		partitionKeys[i] = ColumnToProto(c)
	}
	return &pb.TableMetadata{
		Name:          tm.Name,
		DatabaseName:  tm.DatabaseName,
		Description:   tm.Description,
		TableType:     tm.TableType,
		Columns:       columns,
		PartitionKeys: partitionKeys,
		Parameters:    tm.Parameters,
	}
}

// ProtoToTableMetadata converts a protobuf TableMetadata to its internal representation.
func ProtoToTableMetadata(p *pb.TableMetadata) *TableMetadata {
	if p == nil {
		return nil
	}
	columns := make([]Column, len(p.Columns))
	for i, c := range p.Columns {
		columns[i] = ProtoToColumn(c)
	}
	partitionKeys := make([]Column, len(p.PartitionKeys))
	for i, c := range p.PartitionKeys {
		partitionKeys[i] = ProtoToColumn(c)
	}
	return &TableMetadata{
		Name:          p.Name,
		DatabaseName:  p.DatabaseName,
		Description:   p.Description,
		TableType:     p.TableType,
		Columns:       columns,
		PartitionKeys: partitionKeys,
		Parameters:    p.Parameters,
	}
}

// ColumnToProto converts a Column to its protobuf representation.
func ColumnToProto(c Column) *pb.Column {
	return &pb.Column{
		Name:    c.Name,
		Type:    c.Type,
		Comment: c.Comment,
	}
}

// ProtoToColumn converts a protobuf Column to its internal representation.
func ProtoToColumn(p *pb.Column) Column {
	return Column{
		Name:    p.Name,
		Type:    p.Type,
		Comment: p.Comment,
	}
}

// StoredRowToProto converts a StoredRow to its protobuf representation.
func StoredRowToProto(sr *StoredRow) *pb.StoredRow {
	if sr == nil {
		return nil
	}
	jsonBytes, _ := json.Marshal(sr.Values)
	return &pb.StoredRow{
		ValuesJson: jsonBytes,
	}
}

// ProtoToStoredRow converts a protobuf StoredRow to its internal representation.
func ProtoToStoredRow(p *pb.StoredRow) *StoredRow {
	if p == nil {
		return nil
	}
	sr := &StoredRow{
		Values: make(map[string]interface{}),
	}
	if len(p.ValuesJson) > 0 {
		if err := json.Unmarshal(p.ValuesJson, &sr.Values); err != nil {
			log.Printf("Failed to unmarshal ValuesJson: %v", err)
		}
	}
	return sr
}

// StoredTableToProto converts a StoredTable to its protobuf representation.
func StoredTableToProto(st *StoredTable) *pb.StoredTable {
	if st == nil {
		return nil
	}
	columns := make([]*pb.Column, len(st.Columns))
	for i, c := range st.Columns {
		columns[i] = ColumnToProto(c)
	}
	rows := make([]*pb.StoredRow, len(st.Rows))
	for i, r := range st.Rows {
		rows[i] = StoredRowToProto(r)
	}
	return &pb.StoredTable{
		DatabaseName: st.DatabaseName,
		TableName:    st.TableName,
		Columns:      columns,
		Rows:         rows,
	}
}

// ProtoToStoredTable converts a protobuf StoredTable to its internal representation.
func ProtoToStoredTable(p *pb.StoredTable) *StoredTable {
	if p == nil {
		return nil
	}
	columns := make([]Column, len(p.Columns))
	for i, c := range p.Columns {
		columns[i] = ProtoToColumn(c)
	}
	rows := make([]*StoredRow, len(p.Rows))
	for i, r := range p.Rows {
		rows[i] = ProtoToStoredRow(r)
	}
	return &StoredTable{
		DatabaseName: p.DatabaseName,
		TableName:    p.TableName,
		Columns:      columns,
		Rows:         rows,
	}
}
