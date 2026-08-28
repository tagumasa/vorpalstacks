package neptunegraph

import (
	"bufio"
	"context"
	"encoding/csv"
	"os"
	"strings"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage/graphengine"
	ngstore "vorpalstacks/internal/store/aws/rds/neptunegraph"
	"vorpalstacks/internal/utils/ntriples"
)

const importBatchSize = 500

// importBatcher manages batched node and edge writes during bulk import,
// flushing automatically when the batch reaches importBatchSize.
type importBatcher struct {
	db            *graphengine.DB
	extToInternal map[string]graphengine.NodeID
	nodeBatch     []struct {
		Labels []string
		Props  graphengine.Props
	}
	edgeBatch []graphengine.Edge
	dictCount int64
	errCount  int64
}

func newImportBatcher(db *graphengine.DB) *importBatcher {
	return &importBatcher{
		db:            db,
		extToInternal: make(map[string]graphengine.NodeID),
	}
}

func (b *importBatcher) queueNode(labels []string, props graphengine.Props) {
	b.nodeBatch = append(b.nodeBatch, struct {
		Labels []string
		Props  graphengine.Props
	}{Labels: labels, Props: props})
	if len(b.nodeBatch) >= importBatchSize {
		b.flushNodes()
	}
}

func (b *importBatcher) queueEdge(from, to graphengine.NodeID, label string, props graphengine.Props) {
	b.edgeBatch = append(b.edgeBatch, graphengine.Edge{
		From: from, To: to, Label: label, Props: props,
	})
	if len(b.edgeBatch) >= importBatchSize {
		b.flushEdges()
	}
}

func (b *importBatcher) flushNodes() {
	if len(b.nodeBatch) == 0 {
		return
	}
	ids, err := b.db.AddNodeBatch(b.nodeBatch)
	if err != nil {
		b.errCount += int64(len(b.nodeBatch))
		logs.Warn("AddNodeBatch failed", logs.Err(err))
	} else {
		for i, id := range ids {
			if i < len(b.nodeBatch) {
				if extID, ok := b.nodeBatch[i].Props["~id"].(string); ok {
					b.extToInternal[extID] = id
					delete(b.nodeBatch[i].Props, "~id")
				}
			}
		}
		b.dictCount += int64(len(ids))
	}
	b.nodeBatch = nil
}

func (b *importBatcher) flushEdges() {
	if len(b.edgeBatch) == 0 {
		return
	}
	_, err := b.db.AddEdgeBatch(b.edgeBatch)
	if err != nil {
		b.errCount += int64(len(b.edgeBatch))
		logs.Warn("AddEdgeBatch failed", logs.Err(err))
	} else {
		b.dictCount += int64(len(b.edgeBatch))
	}
	b.edgeBatch = nil
}

func (b *importBatcher) flush() {
	b.flushNodes()
	b.flushEdges()
}

// CreateGraphUsingImportTask creates a new graph and initiates a bulk import task from the specified source.
func (s *NeptuneGraphService) CreateGraphUsingImportTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	in := &CreateGraphUsingImportTaskInput{
		GraphName:               request.GetStringParam(req.Parameters, "graphName"),
		RoleArn:                 request.GetStringParam(req.Parameters, "roleArn"),
		Source:                  request.GetStringParam(req.Parameters, "source"),
		Format:                  strings.ToUpper(request.GetStringParam(req.Parameters, "format")),
		ParquetType:             strings.ToUpper(request.GetStringParam(req.Parameters, "parquetType")),
		BlankNodeHandling:       request.GetStringParam(req.Parameters, "blankNodeHandling"),
		KmsKeyIdentifier:        request.GetStringParam(req.Parameters, "kmsKeyIdentifier"),
		DeletionProtection:      request.GetBoolParam(req.Parameters, "deletionProtection"),
		PublicConnectivity:      request.GetBoolParam(req.Parameters, "publicConnectivity"),
		FailOnError:             request.GetBoolParam(req.Parameters, "failOnError"),
		VectorSearchConfig:      req.Parameters["vectorSearchConfiguration"],
		HasReplicaCount:         request.HasParam(req.Parameters, "replicaCount"),
		ReplicaCount:            request.GetIntParam(req.Parameters, "replicaCount"),
		HasMinProvisionedMemory: request.HasParam(req.Parameters, "minProvisionedMemory"),
		MinProvisionedMemory:    request.GetIntParam(req.Parameters, "minProvisionedMemory"),
		HasMaxProvisionedMemory: request.HasParam(req.Parameters, "maxProvisionedMemory"),
		MaxProvisionedMemory:    request.GetIntParam(req.Parameters, "maxProvisionedMemory"),
		HasImportOptions:        request.HasParam(req.Parameters, "importOptions"),
		ImportOptions:           parseImportOptions(req.Parameters),
		Tags:                    parseTagsFromParams(req.Parameters),
		Region:                  reqCtx.GetRegion(),
	}

	task, err := s.createGraphUsingImportTaskCore(ctx, store, in)
	if err != nil {
		return nil, err
	}
	return importTaskToResponse(task), nil
}

// GetImportTask retrieves an import task by its identifier.
func (s *NeptuneGraphService) GetImportTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	task, err := s.getImportTaskCore(store, request.GetStringParam(req.Parameters, "taskIdentifier"))
	if err != nil {
		return nil, err
	}

	return importTaskToResponse(task), nil
}

// ListImportTasks returns a paginated list of import task summaries.
func (s *NeptuneGraphService) ListImportTasks(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	in := &ListImportTasksInput{
		MaxItems: clampMaxResults(request.GetIntParam(req.Parameters, "maxResults")),
		Marker:   request.GetStringParam(req.Parameters, "nextToken"),
	}

	res, err := s.listImportTasksCore(store, in)
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0, len(res.Tasks))
	for _, t := range res.Tasks {
		items = append(items, importTaskSummaryToResponse(t))
	}

	result := map[string]interface{}{
		"tasks": items,
	}
	if res.Truncated {
		result["nextToken"] = res.NextToken
	}
	return result, nil
}

// CancelImportTask cancels an in-progress import task, transitioning it to CANCELLED state.
func (s *NeptuneGraphService) CancelImportTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	in := &CancelImportTaskInput{
		TaskIdentifier: request.GetStringParam(req.Parameters, "taskIdentifier"),
	}

	task, err := s.cancelImportTaskCore(store, in)
	if err != nil {
		return nil, err
	}
	return importTaskSummaryToResponse(task), nil
}

// StartImportTask initiates a bulk import task on an existing graph from the specified source.
func (s *NeptuneGraphService) StartImportTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	in := &StartImportTaskInput{
		GraphIdentifier:   request.GetStringParam(req.Parameters, "graphIdentifier"),
		RoleArn:           request.GetStringParam(req.Parameters, "roleArn"),
		Source:            request.GetStringParam(req.Parameters, "source"),
		Format:            strings.ToUpper(request.GetStringParam(req.Parameters, "format")),
		ParquetType:       strings.ToUpper(request.GetStringParam(req.Parameters, "parquetType")),
		BlankNodeHandling: request.GetStringParam(req.Parameters, "blankNodeHandling"),
		FailOnError:       request.GetBoolParam(req.Parameters, "failOnError"),
		HasImportOptions:  request.HasParam(req.Parameters, "importOptions"),
		ImportOptions:     parseImportOptions(req.Parameters),
	}

	task, err := s.startImportTaskCore(store, in)
	if err != nil {
		return nil, err
	}
	return importTaskToResponse(task), nil
}

func (s *NeptuneGraphService) importCSV(db *graphengine.DB, filePath string) (statementCount, dictionaryCount, errorCount int64) {
	f, err := os.Open(filePath)
	if err != nil {
		errorCount++
		logs.Warn("failed to open import file", logs.String("path", filePath), logs.Err(err))
		return
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.LazyQuotes = true

	header, err := r.Read()
	if err != nil {
		errorCount++
		logs.Warn("failed to read CSV header", logs.Err(err))
		return
	}

	headerLower := make([]string, len(header))
	for i, h := range header {
		headerLower[i] = strings.ToLower(strings.TrimSpace(h))
	}

	hasFromTo := false
	for _, h := range headerLower {
		if h == "~from" || h == "~to" {
			hasFromTo = true
			break
		}
	}
	isEdge := hasFromTo

	idIndex := -1
	labelIndex := -1
	fromIndex := -1
	toIndex := -1
	for i, h := range headerLower {
		switch h {
		case "~id":
			idIndex = i
		case "~label":
			labelIndex = i
		case "~from":
			fromIndex = i
		case "~to":
			toIndex = i
		}
	}

	propIndices := make(map[string]int)
	for i, h := range headerLower {
		if !strings.HasPrefix(h, "~") {
			propIndices[h] = i
		}
	}

	b := newImportBatcher(db)

	if isEdge {
		// Edge-only CSV files have no node records to populate the
		// ext-to-internal ID mapping. Build it from existing DB nodes
		// so that ~from/~to external IDs can be resolved.
		_ = db.ForEachNode(func(node *graphengine.Node) error {
			if extID, ok := node.Props["~id"].(string); ok {
				b.extToInternal[extID] = node.ID
			}
			return nil
		})
	}

	for {
		record, err := r.Read()
		if err != nil {
			if err.Error() != "EOF" {
				errorCount++
			}
			break
		}

		statementCount++

		if isEdge {
			b.flushNodes()

			if fromIndex < 0 || toIndex < 0 {
				errorCount++
				continue
			}

			fromExt := strings.TrimSpace(record[fromIndex])
			toExt := strings.TrimSpace(record[toIndex])

			fromID, ok := b.extToInternal[fromExt]
			if !ok {
				errorCount++
				continue
			}
			toID, ok := b.extToInternal[toExt]
			if !ok {
				errorCount++
				continue
			}

			label := ""
			if labelIndex >= 0 && labelIndex < len(record) {
				label = strings.TrimSpace(record[labelIndex])
			}

			props := make(graphengine.Props)
			for k, idx := range propIndices {
				if idx < len(record) && record[idx] != "" {
					props[k] = strings.TrimSpace(record[idx])
				}
			}

			b.queueEdge(fromID, toID, label, props)
		} else {
			var labels []string
			if labelIndex >= 0 && labelIndex < len(record) {
				for _, l := range strings.Split(strings.TrimSpace(record[labelIndex]), ":") {
					if l != "" {
						labels = append(labels, strings.TrimSpace(l))
					}
				}
			}

			props := make(graphengine.Props)
			if idIndex >= 0 && idIndex < len(record) {
				props["~id"] = strings.TrimSpace(record[idIndex])
			}
			for k, idx := range propIndices {
				if idx < len(record) && record[idx] != "" {
					props[k] = strings.TrimSpace(record[idx])
				}
			}

			b.queueNode(labels, props)
		}
	}

	b.flush()
	return statementCount, b.dictCount, errorCount + b.errCount
}

func (s *NeptuneGraphService) importRDF(db *graphengine.DB, filePath string) (statementCount, dictionaryCount, errorCount int64) {
	f, err := os.Open(filePath)
	if err != nil {
		errorCount++
		logs.Warn("failed to open import file", logs.String("path", filePath), logs.Err(err))
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 1024*1024), 10*1024*1024)

	extToInternal := make(map[string]graphengine.NodeID)
	b := newImportBatcher(db)
	b.extToInternal = extToInternal

	ensureNode := func(extID string) graphengine.NodeID {
		if id, ok := extToInternal[extID]; ok {
			return id
		}
		nodeID, err := db.AddNode([]string{"Resource"}, graphengine.Props{"uri": extID})
		if err != nil {
			return 0
		}
		extToInternal[extID] = nodeID
		dictionaryCount++
		return nodeID
	}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		subject, predicate, obj, ok := ntriples.ParseLine(line)
		if !ok {
			errorCount++
			continue
		}

		statementCount++

		subjID := ensureNode(subject)

		isURI := strings.HasPrefix(obj, "<")
		if isURI {
			objURI := strings.Trim(obj, "<>")
			objID := ensureNode(objURI)

			predLabel := ntriples.ExtractLocalName(predicate)

			b.queueEdge(subjID, objID, predLabel, nil)
		} else {
			predKey := ntriples.ExtractLocalName(predicate)

			node, err := db.GetNode(subjID)
			if err == nil && node != nil {
				updatedProps := make(graphengine.Props)
				for k, v := range node.Props {
					updatedProps[k] = v
				}
				if strings.HasPrefix(obj, "\"") {
					closing := strings.Index(obj[1:], "\"")
					if closing >= 0 {
						obj = obj[1 : closing+1]
					}
				}
				updatedProps[predKey] = obj
				_ = db.UpdateNode(subjID, updatedProps)
			}
		}
	}

	b.flush()
	return statementCount, dictionaryCount + b.dictCount, errorCount + b.errCount
}

func parseImportOptions(params map[string]interface{}) *ngstore.ImportOptions {
	v, ok := params["importOptions"]
	if !ok || v == nil {
		return nil
	}
	m, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	neptune, ok := m["neptune"]
	if !ok || neptune == nil {
		return nil
	}
	nm, ok := neptune.(map[string]interface{})
	if !ok {
		return nil
	}
	opts := &ngstore.ImportOptions{
		Neptune: &ngstore.NeptuneImportOptions{},
	}
	if s, ok := nm["s3ExportPath"].(string); ok {
		opts.Neptune.S3ExportPath = s
	}
	if s, ok := nm["s3ExportKmsKeyId"].(string); ok {
		opts.Neptune.S3ExportKmsKeyId = s
	}
	if b, ok := nm["preserveDefaultVertexLabels"].(bool); ok {
		opts.Neptune.PreserveDefaultVertexLabels = &b
	}
	if b, ok := nm["preserveEdgeIds"].(bool); ok {
		opts.Neptune.PreserveEdgeIds = &b
	}
	return opts
}
