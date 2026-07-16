package neptunegraph

import (
	pb "vorpalstacks/internal/pb/storage/storage_neptunegraph"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func graphToProto(g *Graph) *pb.Graph {
	if g == nil {
		return nil
	}
	p := &pb.Graph{
		Id:                 g.Id,
		Name:               g.Name,
		Arn:                g.Arn,
		Status:             g.Status,
		Endpoint:           g.Endpoint,
		DeletionProtection: g.DeletionProtection,
		PublicConnectivity: g.PublicConnectivity,
		KmsKeyIdentifier:   g.KmsKeyIdentifier,
		BuildNumber:        g.BuildNumber,
		SourceSnapshotId:   g.SourceSnapshotId,
		StatusReason:       g.StatusReason,
		AccountId:          g.AccountID,
		Region:             g.Region,
	}
	if g.ProvisionedMemory != nil {
		p.ProvisionedMemory = *g.ProvisionedMemory
	}
	if g.ReplicaCount != nil {
		p.ReplicaCount = *g.ReplicaCount
	}
	if g.VectorSearchConfiguration != nil {
		p.VectorSearchConfiguration = &pb.VectorSearchConfiguration{
			Dimension: g.VectorSearchConfiguration.Dimension,
		}
	}
	if g.CreateTime != nil {
		p.CreateTime = timestamppb.New(*g.CreateTime)
	}
	return p
}

func protoToGraph(p *pb.Graph) *Graph {
	if p == nil {
		return nil
	}
	g := &Graph{
		Id:                 p.GetId(),
		Name:               p.GetName(),
		Arn:                p.GetArn(),
		Status:             p.GetStatus(),
		Endpoint:           p.GetEndpoint(),
		DeletionProtection: p.GetDeletionProtection(),
		PublicConnectivity: p.GetPublicConnectivity(),
		KmsKeyIdentifier:   p.GetKmsKeyIdentifier(),
		BuildNumber:        p.GetBuildNumber(),
		SourceSnapshotId:   p.GetSourceSnapshotId(),
		StatusReason:       p.GetStatusReason(),
		AccountID:          p.GetAccountId(),
		Region:             p.GetRegion(),
	}
	if v := p.GetProvisionedMemory(); v != 0 {
		g.ProvisionedMemory = int32Ptr(v)
	}
	if v := p.GetReplicaCount(); v != 0 {
		g.ReplicaCount = int32Ptr(v)
	}
	if p.VectorSearchConfiguration != nil {
		g.VectorSearchConfiguration = &VectorSearchConfig{
			Dimension: p.VectorSearchConfiguration.GetDimension(),
		}
	}
	if p.CreateTime != nil {
		t := p.CreateTime.AsTime()
		g.CreateTime = &t
	}
	return g
}

func snapshotToProto(s *GraphSnapshot) *pb.GraphSnapshot {
	if s == nil {
		return nil
	}
	p := &pb.GraphSnapshot{
		Id:               s.Id,
		Name:             s.Name,
		Arn:              s.Arn,
		Status:           s.Status,
		KmsKeyIdentifier: s.KmsKeyIdentifier,
		SourceGraphId:    s.SourceGraphId,
		AccountId:        s.AccountID,
		Region:           s.Region,
	}
	if s.SnapshotCreateTime != nil {
		p.SnapshotCreateTime = timestamppb.New(*s.SnapshotCreateTime)
	}
	return p
}

func protoToSnapshot(p *pb.GraphSnapshot) *GraphSnapshot {
	if p == nil {
		return nil
	}
	s := &GraphSnapshot{
		Id:               p.GetId(),
		Name:             p.GetName(),
		Arn:              p.GetArn(),
		Status:           p.GetStatus(),
		KmsKeyIdentifier: p.GetKmsKeyIdentifier(),
		SourceGraphId:    p.GetSourceGraphId(),
		AccountID:        p.GetAccountId(),
		Region:           p.GetRegion(),
	}
	if p.SnapshotCreateTime != nil {
		t := p.SnapshotCreateTime.AsTime()
		s.SnapshotCreateTime = &t
	}
	return s
}

func endpointToProto(e *PrivateGraphEndpoint) *pb.PrivateGraphEndpoint {
	if e == nil {
		return nil
	}
	return &pb.PrivateGraphEndpoint{
		GraphId:       e.GraphId,
		VpcId:         e.VpcId,
		VpcEndpointId: e.VpcEndpointId,
		SubnetIds:     e.SubnetIds,
		Status:        e.Status,
		AccountId:     e.AccountID,
		Region:        e.Region,
	}
}

func protoToEndpoint(p *pb.PrivateGraphEndpoint) *PrivateGraphEndpoint {
	if p == nil {
		return nil
	}
	return &PrivateGraphEndpoint{
		GraphId:       p.GetGraphId(),
		VpcId:         p.GetVpcId(),
		VpcEndpointId: p.GetVpcEndpointId(),
		SubnetIds:     p.SubnetIds,
		Status:        p.GetStatus(),
		AccountID:     p.GetAccountId(),
		Region:        p.GetRegion(),
	}
}

func queryToProto(q *QueryRecord) *pb.QueryRecord {
	if q == nil {
		return nil
	}
	return &pb.QueryRecord{
		Id:          q.Id,
		QueryString: q.QueryString,
		Language:    q.Language,
		State:       q.State,
		GraphId:     q.GraphId,
		Elapsed:     q.Elapsed,
		Waited:      q.Waited,
		StartedAt:   timestamppb.New(q.StartedAt),
	}
}

func protoToQuery(p *pb.QueryRecord) *QueryRecord {
	if p == nil {
		return nil
	}
	q := &QueryRecord{
		Id:          p.GetId(),
		QueryString: p.GetQueryString(),
		Language:    p.GetLanguage(),
		State:       p.GetState(),
		GraphId:     p.GetGraphId(),
		Elapsed:     p.GetElapsed(),
		Waited:      p.GetWaited(),
	}
	if p.GetStartedAt() != nil {
		q.StartedAt = p.GetStartedAt().AsTime()
	}
	return q
}

func importTaskToProto(t *ImportTask) *pb.ImportTaskRecord {
	if t == nil {
		return nil
	}
	p := &pb.ImportTaskRecord{
		TaskId:             t.TaskId,
		GraphId:            t.GraphId,
		Status:             t.Status,
		Format:             t.Format,
		ParquetType:        t.ParquetType,
		Source:             t.Source,
		RoleArn:            t.RoleArn,
		AttemptNumber:      t.AttemptNumber,
		StatusReason:       t.StatusReason,
		GraphName:          t.GraphName,
		DeletionProtection: t.DeletionProtection,
		KmsKeyIdentifier:   t.KmsKeyIdentifier,
		PublicConnectivity: t.PublicConnectivity,
		BlankNodeHandling:  t.BlankNodeHandling,
		FailOnError:        t.FailOnError,
	}
	if t.ReplicaCount != nil {
		p.ReplicaCount = *t.ReplicaCount
	}
	if t.MinProvisionedMemory != nil {
		p.MinProvisionedMemory = *t.MinProvisionedMemory
	}
	if t.MaxProvisionedMemory != nil {
		p.MaxProvisionedMemory = *t.MaxProvisionedMemory
	}
	if t.VectorSearchConfiguration != nil {
		p.VectorSearchConfiguration = &pb.VectorSearchConfiguration{
			Dimension: t.VectorSearchConfiguration.Dimension,
		}
	}
	if t.ImportOptions != nil {
		p.ImportOptions = importOptionsToProto(t.ImportOptions)
	}
	if t.ImportTaskDetails != nil {
		p.ImportTaskDetails = importTaskDetailsToProto(t.ImportTaskDetails)
	}
	if t.StartTime != nil {
		p.StartTime = timestamppb.New(*t.StartTime)
	}
	return p
}

func protoToImportTask(p *pb.ImportTaskRecord) *ImportTask {
	if p == nil {
		return nil
	}
	t := &ImportTask{
		TaskId:             p.GetTaskId(),
		GraphId:            p.GetGraphId(),
		Status:             p.GetStatus(),
		Format:             p.GetFormat(),
		ParquetType:        p.GetParquetType(),
		Source:             p.GetSource(),
		RoleArn:            p.GetRoleArn(),
		AttemptNumber:      p.GetAttemptNumber(),
		StatusReason:       p.GetStatusReason(),
		GraphName:          p.GetGraphName(),
		DeletionProtection: p.GetDeletionProtection(),
		KmsKeyIdentifier:   p.GetKmsKeyIdentifier(),
		PublicConnectivity: p.GetPublicConnectivity(),
		BlankNodeHandling:  p.GetBlankNodeHandling(),
		FailOnError:        p.GetFailOnError(),
	}
	if v := p.GetReplicaCount(); v != 0 {
		t.ReplicaCount = int32Ptr(v)
	}
	if v := p.GetMinProvisionedMemory(); v != 0 {
		t.MinProvisionedMemory = int32Ptr(v)
	}
	if v := p.GetMaxProvisionedMemory(); v != 0 {
		t.MaxProvisionedMemory = int32Ptr(v)
	}
	if p.VectorSearchConfiguration != nil {
		t.VectorSearchConfiguration = &VectorSearchConfig{
			Dimension: p.VectorSearchConfiguration.GetDimension(),
		}
	}
	if p.ImportOptions != nil {
		t.ImportOptions = protoToImportOptions(p.ImportOptions)
	}
	if p.ImportTaskDetails != nil {
		t.ImportTaskDetails = protoToImportTaskDetails(p.ImportTaskDetails)
	}
	if p.StartTime != nil {
		startTime := p.StartTime.AsTime()
		t.StartTime = &startTime
	}
	return t
}

func importOptionsToProto(o *ImportOptions) *pb.ImportOptions {
	if o == nil || o.Neptune == nil {
		return nil
	}
	p := &pb.ImportOptions{
		Neptune: &pb.NeptuneImportOptions{
			S3ExportPath:     o.Neptune.S3ExportPath,
			S3ExportKmsKeyId: o.Neptune.S3ExportKmsKeyId,
		},
	}
	if o.Neptune.PreserveDefaultVertexLabels != nil {
		p.Neptune.PreserveDefaultVertexLabels = *o.Neptune.PreserveDefaultVertexLabels
	}
	if o.Neptune.PreserveEdgeIds != nil {
		p.Neptune.PreserveEdgeIds = *o.Neptune.PreserveEdgeIds
	}
	return p
}

func protoToImportOptions(p *pb.ImportOptions) *ImportOptions {
	if p == nil || p.Neptune == nil {
		return nil
	}
	return &ImportOptions{
		Neptune: &NeptuneImportOptions{
			S3ExportPath:                p.Neptune.GetS3ExportPath(),
			S3ExportKmsKeyId:            p.Neptune.GetS3ExportKmsKeyId(),
			PreserveDefaultVertexLabels: proto.Bool(p.Neptune.GetPreserveDefaultVertexLabels()),
			PreserveEdgeIds:             proto.Bool(p.Neptune.GetPreserveEdgeIds()),
		},
	}
}

func importTaskDetailsToProto(d *ImportTaskDetails) *pb.ImportTaskDetails {
	if d == nil {
		return nil
	}
	p := &pb.ImportTaskDetails{}
	if d.ProgressPercentage != nil {
		p.ProgressPercentage = *d.ProgressPercentage
	}
	if d.TimeElapsedSeconds != nil {
		p.TimeElapsedSeconds = *d.TimeElapsedSeconds
	}
	if d.StatementCount != nil {
		p.StatementCount = *d.StatementCount
	}
	if d.DictionaryEntryCount != nil {
		p.DictionaryEntryCount = *d.DictionaryEntryCount
	}
	if d.ErrorCount != nil {
		p.ErrorCount = *d.ErrorCount
	}
	if d.ErrorDetails != nil {
		p.ErrorDetails = *d.ErrorDetails
	}
	if d.Status != nil {
		p.Status = *d.Status
	}
	if d.StartTime != nil {
		p.StartTime = timestamppb.New(*d.StartTime)
	}
	return p
}

func protoToImportTaskDetails(p *pb.ImportTaskDetails) *ImportTaskDetails {
	if p == nil {
		return nil
	}
	d := &ImportTaskDetails{
		ProgressPercentage:   int32Ptr(p.GetProgressPercentage()),
		TimeElapsedSeconds:   int64Ptr(p.GetTimeElapsedSeconds()),
		StatementCount:       int64Ptr(p.GetStatementCount()),
		DictionaryEntryCount: int64Ptr(p.GetDictionaryEntryCount()),
		ErrorCount:           int32Ptr(p.GetErrorCount()),
		ErrorDetails:         stringPtr(p.GetErrorDetails()),
		Status:               stringPtr(p.GetStatus()),
	}
	if p.StartTime != nil {
		t := p.StartTime.AsTime()
		d.StartTime = &t
	}
	return d
}

func exportTaskToProto(t *ExportTask) *pb.ExportTaskRecord {
	if t == nil {
		return nil
	}
	p := &pb.ExportTaskRecord{
		TaskId:           t.TaskId,
		GraphId:          t.GraphId,
		Status:           t.Status,
		Format:           t.Format,
		ParquetType:      t.ParquetType,
		Destination:      t.Destination,
		RoleArn:          t.RoleArn,
		KmsKeyIdentifier: t.KmsKeyIdentifier,
		StatusReason:     t.StatusReason,
	}
	if t.ExportFilter != nil {
		p.ExportFilter = exportFilterToProto(t.ExportFilter)
	}
	if t.ExportTaskDetails != nil {
		p.ExportTaskDetails = exportTaskDetailsToProto(t.ExportTaskDetails)
	}
	if t.StartTime != nil {
		p.StartTime = timestamppb.New(*t.StartTime)
	}
	return p
}

func protoToExportTask(p *pb.ExportTaskRecord) *ExportTask {
	if p == nil {
		return nil
	}
	t := &ExportTask{
		TaskId:           p.GetTaskId(),
		GraphId:          p.GetGraphId(),
		Status:           p.GetStatus(),
		Format:           p.GetFormat(),
		ParquetType:      p.GetParquetType(),
		Destination:      p.GetDestination(),
		RoleArn:          p.GetRoleArn(),
		KmsKeyIdentifier: p.GetKmsKeyIdentifier(),
		StatusReason:     p.GetStatusReason(),
	}
	if p.ExportFilter != nil {
		t.ExportFilter = protoToExportFilter(p.ExportFilter)
	}
	if p.ExportTaskDetails != nil {
		t.ExportTaskDetails = protoToExportTaskDetails(p.ExportTaskDetails)
	}
	if p.StartTime != nil {
		startTime := p.StartTime.AsTime()
		t.StartTime = &startTime
	}
	return t
}

func exportFilterToProto(f *ExportFilter) *pb.ExportFilter {
	if f == nil {
		return nil
	}
	p := &pb.ExportFilter{
		EdgeFilter:   make(map[string]*pb.ExportFilterElement),
		VertexFilter: make(map[string]*pb.ExportFilterElement),
	}
	for k, v := range f.EdgeFilter {
		p.EdgeFilter[k] = exportFilterElementToProto(&v)
	}
	for k, v := range f.VertexFilter {
		p.VertexFilter[k] = exportFilterElementToProto(&v)
	}
	return p
}

func protoToExportFilter(p *pb.ExportFilter) *ExportFilter {
	if p == nil {
		return nil
	}
	f := &ExportFilter{
		EdgeFilter:   make(map[string]ExportFilterElement),
		VertexFilter: make(map[string]ExportFilterElement),
	}
	for k, v := range p.EdgeFilter {
		f.EdgeFilter[k] = *protoToExportFilterElement(v)
	}
	for k, v := range p.VertexFilter {
		f.VertexFilter[k] = *protoToExportFilterElement(v)
	}
	return f
}

func exportFilterElementToProto(e *ExportFilterElement) *pb.ExportFilterElement {
	if e == nil {
		return nil
	}
	p := &pb.ExportFilterElement{
		Properties: make(map[string]*pb.ExportFilterPropertyAttributes),
	}
	for k, v := range e.Properties {
		pa := &pb.ExportFilterPropertyAttributes{
			MultiValueHandling: v.MultiValueHandling,
		}
		if v.OutputType != nil {
			pa.OutputType = *v.OutputType
		}
		if v.SourcePropertyName != nil {
			pa.SourcePropertyName = *v.SourcePropertyName
		}
		p.Properties[k] = pa
	}
	return p
}

func protoToExportFilterElement(p *pb.ExportFilterElement) *ExportFilterElement {
	if p == nil {
		return nil
	}
	e := &ExportFilterElement{
		Properties: make(map[string]ExportFilterPropertyAttributes),
	}
	for k, v := range p.Properties {
		pa := ExportFilterPropertyAttributes{
			MultiValueHandling: v.GetMultiValueHandling(),
			OutputType:         stringPtr(v.GetOutputType()),
			SourcePropertyName: stringPtr(v.GetSourcePropertyName()),
		}
		e.Properties[k] = pa
	}
	return e
}

func exportTaskDetailsToProto(d *ExportTaskDetails) *pb.ExportTaskDetails {
	if d == nil {
		return nil
	}
	p := &pb.ExportTaskDetails{}
	if d.ProgressPercentage != nil {
		p.ProgressPercentage = *d.ProgressPercentage
	}
	if d.TimeElapsedSeconds != nil {
		p.TimeElapsedSeconds = *d.TimeElapsedSeconds
	}
	if d.NumEdgesWritten != nil {
		p.NumEdgesWritten = *d.NumEdgesWritten
	}
	if d.NumVerticesWritten != nil {
		p.NumVerticesWritten = *d.NumVerticesWritten
	}
	if d.StartTime != nil {
		p.StartTime = timestamppb.New(*d.StartTime)
	}
	return p
}

func protoToExportTaskDetails(p *pb.ExportTaskDetails) *ExportTaskDetails {
	if p == nil {
		return nil
	}
	d := &ExportTaskDetails{
		ProgressPercentage: int32Ptr(p.GetProgressPercentage()),
		TimeElapsedSeconds: int64Ptr(p.GetTimeElapsedSeconds()),
		NumEdgesWritten:    int64Ptr(p.GetNumEdgesWritten()),
		NumVerticesWritten: int64Ptr(p.GetNumVerticesWritten()),
	}
	if p.StartTime != nil {
		t := p.StartTime.AsTime()
		d.StartTime = &t
	}
	return d
}

func tagsToProto(tags map[string]string) *pb.TagMap {
	if tags == nil {
		return nil
	}
	return &pb.TagMap{Tags: tags}
}

func protoToTags(p *pb.TagMap) map[string]string {
	if p == nil {
		return nil
	}
	return p.Tags
}

func int32Ptr(v int32) *int32    { return &v }
func int64Ptr(v int64) *int64    { return &v }
func stringPtr(v string) *string { return &v }
