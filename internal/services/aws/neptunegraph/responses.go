package neptunegraph

import (
	ngstore "vorpalstacks/internal/store/aws/neptunegraph"
)

func graphToResponse(g *ngstore.Graph, isCreate bool) map[string]interface{} {
	r := map[string]interface{}{
		"id":                 g.Id,
		"name":               g.Name,
		"arn":                g.Arn,
		"status":             g.Status,
		"endpoint":           g.Endpoint,
		"deletionProtection": g.DeletionProtection,
		"publicConnectivity": g.PublicConnectivity,
		"buildNumber":        g.BuildNumber,
		"sourceSnapshotId":   g.SourceSnapshotId,
		"accountId":          g.AccountID,
		"region":             g.Region,
	}
	if g.ProvisionedMemory != nil {
		r["provisionedMemory"] = *g.ProvisionedMemory
	}
	if g.ReplicaCount != nil {
		r["replicaCount"] = *g.ReplicaCount
	}
	if g.VectorSearchConfiguration != nil {
		r["vectorSearchConfiguration"] = map[string]interface{}{
			"dimension": g.VectorSearchConfiguration.Dimension,
		}
	}
	if g.KmsKeyIdentifier != "" {
		r["kmsKeyIdentifier"] = g.KmsKeyIdentifier
	}
	if g.StatusReason != "" {
		r["statusReason"] = g.StatusReason
	}
	if g.CreateTime != nil {
		r["createTime"] = g.CreateTime.Unix()
	}
	return r
}

func graphSummaryToResponse(g *ngstore.Graph) map[string]interface{} {
	r := map[string]interface{}{
		"id":                 g.Id,
		"name":               g.Name,
		"arn":                g.Arn,
		"status":             g.Status,
		"endpoint":           g.Endpoint,
		"deletionProtection": g.DeletionProtection,
		"publicConnectivity": g.PublicConnectivity,
	}
	if g.ProvisionedMemory != nil {
		r["provisionedMemory"] = *g.ProvisionedMemory
	}
	if g.ReplicaCount != nil {
		r["replicaCount"] = *g.ReplicaCount
	}
	if g.KmsKeyIdentifier != "" {
		r["kmsKeyIdentifier"] = g.KmsKeyIdentifier
	}
	return r
}

func snapshotToResponse(s *ngstore.GraphSnapshot) map[string]interface{} {
	r := map[string]interface{}{
		"id":            s.Id,
		"name":          s.Name,
		"arn":           s.Arn,
		"status":        s.Status,
		"sourceGraphId": s.SourceGraphId,
		"accountId":     s.AccountID,
		"region":        s.Region,
	}
	if s.KmsKeyIdentifier != "" {
		r["kmsKeyIdentifier"] = s.KmsKeyIdentifier
	}
	if s.SnapshotCreateTime != nil {
		r["snapshotCreateTime"] = s.SnapshotCreateTime.Unix()
	}
	return r
}

func endpointToResponse(ep *ngstore.PrivateGraphEndpoint) map[string]interface{} {
	r := map[string]interface{}{
		"status":    ep.Status,
		"subnetIds": ep.SubnetIds,
		"accountId": ep.AccountID,
		"region":    ep.Region,
	}
	if ep.VpcId != "" {
		r["vpcId"] = ep.VpcId
	}
	if ep.VpcEndpointId != "" {
		r["vpcEndpointId"] = ep.VpcEndpointId
	}
	return r
}

func queryToResponse(q *ngstore.QueryRecord) map[string]interface{} {
	r := map[string]interface{}{
		"id":          q.Id,
		"queryString": q.QueryString,
		"state":       q.State,
	}
	if q.Elapsed > 0 {
		r["elapsed"] = q.Elapsed
	}
	if q.Waited > 0 {
		r["waited"] = q.Waited
	}
	return r
}

func graphDataSummaryToResponse(s *ngstore.GraphDataSummary) map[string]interface{} {
	r := map[string]interface{}{
		"numNodes":      s.NumNodes,
		"numEdges":      s.NumEdges,
		"numNodeLabels": s.NumNodeLabels,
		"numEdgeLabels": s.NumEdgeLabels,
		"nodeLabels":    s.NodeLabels,
		"edgeLabels":    s.EdgeLabels,
	}
	if s.NumNodeProperties != nil {
		r["numNodeProperties"] = *s.NumNodeProperties
	}
	if s.NumEdgeProperties != nil {
		r["numEdgeProperties"] = *s.NumEdgeProperties
	}
	if s.TotalNodePropertyValues != nil {
		r["totalNodePropertyValues"] = *s.TotalNodePropertyValues
	}
	if s.TotalEdgePropertyValues != nil {
		r["totalEdgePropertyValues"] = *s.TotalEdgePropertyValues
	}
	if len(s.NodeStructures) > 0 {
		ns := make([]interface{}, 0, len(s.NodeStructures))
		for _, n := range s.NodeStructures {
			ns = append(ns, nodeStructureToResponse(&n))
		}
		r["nodeStructures"] = ns
	}
	if len(s.EdgeStructures) > 0 {
		es := make([]interface{}, 0, len(s.EdgeStructures))
		for _, e := range s.EdgeStructures {
			es = append(es, edgeStructureToResponse(&e))
		}
		r["edgeStructures"] = es
	}
	return r
}

func nodeStructureToResponse(n *ngstore.NodeStructure) map[string]interface{} {
	r := map[string]interface{}{
		"nodeProperties": n.NodeProperties,
	}
	if n.Count != nil {
		r["count"] = *n.Count
	}
	if len(n.DistinctOutgoingEdgeLabels) > 0 {
		r["distinctOutgoingEdgeLabels"] = n.DistinctOutgoingEdgeLabels
	}
	return r
}

func edgeStructureToResponse(e *ngstore.EdgeStructure) map[string]interface{} {
	r := map[string]interface{}{
		"edgeProperties": e.EdgeProperties,
	}
	if e.Count != nil {
		r["count"] = *e.Count
	}
	return r
}

func importTaskToResponse(t *ngstore.ImportTask) map[string]interface{} {
	r := map[string]interface{}{
		"taskId":  t.TaskId,
		"graphId": t.GraphId,
		"status":  t.Status,
		"source":  t.Source,
		"roleArn": t.RoleArn,
	}
	if t.Format != "" {
		r["format"] = t.Format
	}
	if t.ParquetType != "" {
		r["parquetType"] = t.ParquetType
	}
	if t.StatusReason != "" {
		r["statusReason"] = t.StatusReason
	}
	if t.ImportTaskDetails != nil {
		r["importTaskDetails"] = importTaskDetailsToResponse(t.ImportTaskDetails)
	}
	if t.ImportOptions != nil {
		r["importOptions"] = importOptionsToResponse(t.ImportOptions)
	}
	return r
}

func importTaskSummaryToResponse(t *ngstore.ImportTask) map[string]interface{} {
	r := map[string]interface{}{
		"taskId":  t.TaskId,
		"graphId": t.GraphId,
		"status":  t.Status,
		"source":  t.Source,
		"roleArn": t.RoleArn,
	}
	if t.Format != "" {
		r["format"] = t.Format
	}
	if t.ParquetType != "" {
		r["parquetType"] = t.ParquetType
	}
	return r
}

func importTaskDetailsToResponse(d *ngstore.ImportTaskDetails) map[string]interface{} {
	r := map[string]interface{}{}
	if d.ProgressPercentage != nil {
		r["progressPercentage"] = *d.ProgressPercentage
	}
	if d.StartTime != nil {
		r["startTime"] = d.StartTime.Unix()
	}
	if d.TimeElapsedSeconds != nil {
		r["timeElapsedSeconds"] = *d.TimeElapsedSeconds
	}
	if d.StatementCount != nil {
		r["statementCount"] = *d.StatementCount
	}
	if d.DictionaryEntryCount != nil {
		r["dictionaryEntryCount"] = *d.DictionaryEntryCount
	}
	if d.ErrorCount != nil {
		r["errorCount"] = *d.ErrorCount
	}
	if d.ErrorDetails != nil {
		r["errorDetails"] = *d.ErrorDetails
	}
	if d.Status != nil {
		r["status"] = *d.Status
	}
	return r
}

func importOptionsToResponse(o *ngstore.ImportOptions) map[string]interface{} {
	if o == nil || o.Neptune == nil {
		return nil
	}
	r := map[string]interface{}{
		"s3ExportPath":     o.Neptune.S3ExportPath,
		"s3ExportKmsKeyId": o.Neptune.S3ExportKmsKeyId,
	}
	if o.Neptune.PreserveDefaultVertexLabels != nil {
		r["preserveDefaultVertexLabels"] = *o.Neptune.PreserveDefaultVertexLabels
	}
	if o.Neptune.PreserveEdgeIds != nil {
		r["preserveEdgeIds"] = *o.Neptune.PreserveEdgeIds
	}
	return map[string]interface{}{"neptune": r}
}

func exportTaskToResponse(t *ngstore.ExportTask) map[string]interface{} {
	r := map[string]interface{}{
		"taskId":           t.TaskId,
		"graphId":          t.GraphId,
		"status":           t.Status,
		"format":           t.Format,
		"destination":      t.Destination,
		"roleArn":          t.RoleArn,
		"kmsKeyIdentifier": t.KmsKeyIdentifier,
	}
	if t.ParquetType != "" {
		r["parquetType"] = t.ParquetType
	}
	if t.StatusReason != "" {
		r["statusReason"] = t.StatusReason
	}
	if t.ExportFilter != nil {
		r["exportFilter"] = exportFilterToResponse(t.ExportFilter)
	}
	if t.ExportTaskDetails != nil {
		r["exportTaskDetails"] = exportTaskDetailsToResponse(t.ExportTaskDetails)
	}
	return r
}

func exportTaskSummaryToResponse(t *ngstore.ExportTask) map[string]interface{} {
	r := map[string]interface{}{
		"taskId":           t.TaskId,
		"graphId":          t.GraphId,
		"status":           t.Status,
		"format":           t.Format,
		"destination":      t.Destination,
		"roleArn":          t.RoleArn,
		"kmsKeyIdentifier": t.KmsKeyIdentifier,
	}
	if t.ParquetType != "" {
		r["parquetType"] = t.ParquetType
	}
	if t.StatusReason != "" {
		r["statusReason"] = t.StatusReason
	}
	return r
}

func exportTaskDetailsToResponse(d *ngstore.ExportTaskDetails) map[string]interface{} {
	r := map[string]interface{}{}
	if d.ProgressPercentage != nil {
		r["progressPercentage"] = *d.ProgressPercentage
	}
	if d.StartTime != nil {
		r["startTime"] = d.StartTime.Unix()
	}
	if d.TimeElapsedSeconds != nil {
		r["timeElapsedSeconds"] = *d.TimeElapsedSeconds
	}
	if d.NumEdgesWritten != nil {
		r["numEdgesWritten"] = *d.NumEdgesWritten
	}
	if d.NumVerticesWritten != nil {
		r["numVerticesWritten"] = *d.NumVerticesWritten
	}
	return r
}

func exportFilterToResponse(f *ngstore.ExportFilter) map[string]interface{} {
	if f == nil {
		return nil
	}
	r := map[string]interface{}{}
	if len(f.EdgeFilter) > 0 {
		ef := make(map[string]interface{})
		for k, v := range f.EdgeFilter {
			ef[k] = exportFilterElementToResponse(&v)
		}
		r["edgeFilter"] = ef
	}
	if len(f.VertexFilter) > 0 {
		vf := make(map[string]interface{})
		for k, v := range f.VertexFilter {
			vf[k] = exportFilterElementToResponse(&v)
		}
		r["vertexFilter"] = vf
	}
	return r
}

func exportFilterElementToResponse(e *ngstore.ExportFilterElement) map[string]interface{} {
	if e == nil {
		return nil
	}
	r := map[string]interface{}{}
	for k, v := range e.Properties {
		props := map[string]interface{}{
			"multiValueHandling": v.MultiValueHandling,
		}
		if v.OutputType != nil {
			props["outputType"] = *v.OutputType
		}
		if v.SourcePropertyName != nil {
			props["sourcePropertyName"] = *v.SourcePropertyName
		}
		r[k] = props
	}
	return r
}
