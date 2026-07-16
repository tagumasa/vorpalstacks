package iotevents

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	tagutil "vorpalstacks/internal/common/tags"
	storecommon "vorpalstacks/internal/store/aws/common"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

// CreateDetectorModel creates a new detector model.
func (s *IoTEventsService) CreateDetectorModel(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "detectorModelName")
	if name == "" {
		return nil, iotstore.ErrInvalidRequest
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if _, err := store.GetDetectorModel(name); err == nil {
		return nil, iotstore.ErrDetectorModelAlreadyExists
	}

	if request.GetParamCaseInsensitive(req.Parameters, "roleArn") == "" {
		return nil, iotstore.ErrInvalidRequest
	}

	evaluationMethod := request.GetParamCaseInsensitive(req.Parameters, "evaluationMethod")
	if evaluationMethod == "" {
		evaluationMethod = "BATCH"
	}
	now := time.Now().UTC()
	def, err := parseStructParam(req.Parameters, "detectorModelDefinition")
	if err != nil {
		return nil, iotstore.ErrInvalidRequest
	}
	dm := &iotstore.DetectorModel{
		DetectorModelName:        name,
		DetectorModelDescription: request.GetParamCaseInsensitive(req.Parameters, "detectorModelDescription"),
		RoleARN:                  request.GetParamCaseInsensitive(req.Parameters, "roleArn"),
		DetectorModelDefinition:  def,
		Key:                      request.GetParamCaseInsensitive(req.Parameters, "key"),
		EvaluationMethod:         evaluationMethod,
		DetectorModelVersion:     "1",
		Status:                   "ACTIVE",
		CreationDate:             now,
		LastModifiedDate:         now,
	}

	created, err := store.CreateDetectorModel(dm)
	if err != nil {
		return nil, err
	}

	// Load the model into the state machine for event evaluation.
	if concrete, ok := store.(*iotstore.IotStore); ok {
		concrete.LoadDetectorModel(created)
	}

	if tags := tagutil.ParseTagsAsMap(req.Parameters, "tags"); len(tags) > 0 {
		if err := store.TagResource(created.DetectorModelARN, tags); err != nil {
			slog.Warn("failed to tag DetectorModel on create", "arn", created.DetectorModelARN, "error", err)
		}
	}

	return map[string]interface{}{
		"detectorModelConfiguration": detectorModelConfig(created),
	}, nil
}

// DescribeDetectorModel describes a detector model.
func (s *IoTEventsService) DescribeDetectorModel(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "detectorModelName")
	if name == "" {
		return nil, iotstore.ErrInvalidRequest
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	dm, err := store.GetDetectorModel(name)
	if err != nil {
		return nil, iotstore.ErrDetectorModelNotFound
	}

	return map[string]interface{}{
		"detectorModel": map[string]interface{}{
			"detectorModelConfiguration": detectorModelConfig(dm),
			"detectorModelDefinition":    dm.DetectorModelDefinition,
		},
	}, nil
}

// UpdateDetectorModel updates an existing detector model.
func (s *IoTEventsService) UpdateDetectorModel(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "detectorModelName")
	if name == "" {
		return nil, iotstore.ErrInvalidRequest
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	existing, err := store.GetDetectorModel(name)
	if err != nil {
		return nil, iotstore.ErrDetectorModelNotFound
	}

	if request.HasParam(req.Parameters, "detectorModelDescription") {
		existing.DetectorModelDescription = request.GetParamCaseInsensitive(req.Parameters, "detectorModelDescription")
	}
	if request.HasParam(req.Parameters, "roleArn") {
		existing.RoleARN = request.GetParamCaseInsensitive(req.Parameters, "roleArn")
	}
	if def, err := parseStructParam(req.Parameters, "detectorModelDefinition"); err != nil {
		return nil, iotstore.ErrInvalidRequest
	} else if len(def) > 0 {
		existing.DetectorModelDefinition = def
	}
	if request.HasParam(req.Parameters, "key") {
		existing.Key = request.GetParamCaseInsensitive(req.Parameters, "key")
	}
	if request.HasParam(req.Parameters, "evaluationMethod") {
		existing.EvaluationMethod = request.GetParamCaseInsensitive(req.Parameters, "evaluationMethod")
	}
	existing.LastModifiedDate = time.Now().UTC()

	if err := store.UpdateDetectorModel(existing); err != nil {
		return nil, err
	}

	if concrete, ok := store.(*iotstore.IotStore); ok {
		concrete.LoadDetectorModel(existing)
	}

	return map[string]interface{}{
		"detectorModelConfiguration": detectorModelConfig(existing),
	}, nil
}

// DeleteDetectorModel deletes a detector model.
func (s *IoTEventsService) DeleteDetectorModel(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "detectorModelName")
	if name == "" {
		return nil, iotstore.ErrInvalidRequest
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if _, err := store.GetDetectorModel(name); err != nil {
		return nil, iotstore.ErrDetectorModelNotFound
	}

	dm, _ := store.GetDetectorModel(name)
	if dm != nil && dm.DetectorModelARN != "" {
		_ = store.DeleteAllTags(dm.DetectorModelARN)
	}

	if err := store.DeleteDetectorModel(name); err != nil {
		return nil, err
	}

	if concrete, ok := store.(*iotstore.IotStore); ok {
		concrete.UnloadModel(name)
	}

	return map[string]interface{}{}, nil
}

// ListDetectorModels lists detector models.
func (s *IoTEventsService) ListDetectorModels(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	opts := storecommon.ListOptions{}
	if token := request.GetParamCaseInsensitive(req.Parameters, "nextToken"); token != "" {
		opts.Marker = token
	}
	if max := request.GetIntParam(req.Parameters, "maxResults"); max > 0 {
		opts.MaxItems = max
	}

	result, err := store.ListDetectorModels(opts)
	if err != nil {
		return nil, err
	}

	summaries := make([]map[string]interface{}, 0, len(result.Items))
	for _, dm := range result.Items {
		summaries = append(summaries, map[string]interface{}{
			"detectorModelName":        dm.DetectorModelName,
			"detectorModelDescription": dm.DetectorModelDescription,
			"creationTime":             dm.CreationDate.Unix(),
		})
	}

	resp := map[string]interface{}{
		"detectorModelSummaries": summaries,
	}
	if result.NextMarker != "" {
		resp["nextToken"] = result.NextMarker
	}
	return resp, nil
}

// CreateInput creates a new input.
func (s *IoTEventsService) CreateInput(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "inputName")
	if name == "" {
		return nil, iotstore.ErrInvalidRequest
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if _, err := store.GetInput(name); err == nil {
		return nil, iotstore.ErrInputAlreadyExists
	}

	now := time.Now().UTC()
	def, err := parseStructParam(req.Parameters, "inputDefinition")
	if err != nil {
		return nil, iotstore.ErrInvalidRequest
	}
	input := &iotstore.Input{
		InputName:        name,
		InputDescription: request.GetParamCaseInsensitive(req.Parameters, "inputDescription"),
		InputDefinition:  def,
		Status:           "ACTIVE",
		CreationDate:     now,
		LastModifiedDate: now,
	}

	created, err := store.CreateInput(input)
	if err != nil {
		return nil, err
	}

	if tags := tagutil.ParseTagsAsMap(req.Parameters, "tags"); len(tags) > 0 {
		if err := store.TagResource(created.InputARN, tags); err != nil {
			slog.Warn("failed to tag Input on create", "arn", created.InputARN, "error", err)
		}
	}

	return map[string]interface{}{
		"inputConfiguration": inputConfig(created),
	}, nil
}

// DescribeInput describes an input.
func (s *IoTEventsService) DescribeInput(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "inputName")
	if name == "" {
		return nil, iotstore.ErrInvalidRequest
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	input, err := store.GetInput(name)
	if err != nil {
		return nil, iotstore.ErrInputNotFound
	}

	return map[string]interface{}{
		"input": map[string]interface{}{
			"inputConfiguration": inputConfig(input),
			"inputDefinition":    input.InputDefinition,
		},
	}, nil
}

// UpdateInput updates an existing input.
func (s *IoTEventsService) UpdateInput(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "inputName")
	if name == "" {
		return nil, iotstore.ErrInvalidRequest
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	existing, err := store.GetInput(name)
	if err != nil {
		return nil, iotstore.ErrInputNotFound
	}

	if request.HasParam(req.Parameters, "inputDescription") {
		existing.InputDescription = request.GetParamCaseInsensitive(req.Parameters, "inputDescription")
	}
	if def, err := parseStructParam(req.Parameters, "inputDefinition"); err != nil {
		return nil, iotstore.ErrInvalidRequest
	} else if len(def) > 0 {
		existing.InputDefinition = def
	}
	existing.LastModifiedDate = time.Now().UTC()

	if err := store.UpdateInput(existing); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"inputConfiguration": inputConfig(existing),
	}, nil
}

// DeleteInput deletes an input.
func (s *IoTEventsService) DeleteInput(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "inputName")
	if name == "" {
		return nil, iotstore.ErrInvalidRequest
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if _, err := store.GetInput(name); err != nil {
		return nil, iotstore.ErrInputNotFound
	}

	input, _ := store.GetInput(name)
	if input != nil && input.InputARN != "" {
		_ = store.DeleteAllTags(input.InputARN)
	}

	if err := store.DeleteInput(name); err != nil {
		return nil, err
	}

	return map[string]interface{}{}, nil
}

// ListInputs lists inputs.
func (s *IoTEventsService) ListInputs(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	opts := storecommon.ListOptions{}
	if token := request.GetParamCaseInsensitive(req.Parameters, "nextToken"); token != "" {
		opts.Marker = token
	}
	if max := request.GetIntParam(req.Parameters, "maxResults"); max > 0 {
		opts.MaxItems = max
	}

	result, err := store.ListInputs(opts)
	if err != nil {
		return nil, err
	}

	summaries := make([]map[string]interface{}, 0, len(result.Items))
	for _, input := range result.Items {
		summaries = append(summaries, map[string]interface{}{
			"inputName":    input.InputName,
			"inputArn":     input.InputARN,
			"creationTime": input.CreationDate.Unix(),
		})
	}

	resp := map[string]interface{}{
		"inputSummaries": summaries,
	}
	if result.NextMarker != "" {
		resp["nextToken"] = result.NextMarker
	}
	return resp, nil
}

// validateEventsResourceARN checks that the referenced IoT Events resource
// actually exists. ARN format: arn:aws:iotevents:region:account:type/name.
func validateEventsResourceARN(store iotstore.IotStoreInterface, arn string) error {
	parts := strings.SplitN(arn, "/", 2)
	if len(parts) != 2 {
		return iotstore.ErrInvalidRequest
	}
	resourceType := parts[len(parts)-2]
	resourceType = resourceType[strings.LastIndex(resourceType, ":")+1:]
	name := parts[len(parts)-1]

	switch resourceType {
	case "detectorModel":
		if _, err := store.GetDetectorModel(name); err != nil {
			return iotstore.ErrDetectorModelNotFound
		}
	case "input":
		if _, err := store.GetInput(name); err != nil {
			return iotstore.ErrInputNotFound
		}
	default:
		return iotstore.ErrInvalidRequest
	}
	return nil
}

// TagResource adds or overwrites tags on an IoT Events resource.
// resourceArn is read from query parameters per Smithy model.
func (s *IoTEventsService) TagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	resourceARN := request.GetParamCaseInsensitive(req.Parameters, "resourceArn")
	if resourceARN == "" {
		return nil, iotstore.ErrInvalidRequest
	}

	if err := validateEventsResourceARN(store, resourceARN); err != nil {
		return nil, err
	}

	tags := tagutil.ParseTagsAsMap(req.Parameters, "tags")
	if len(tags) > 0 {
		if err := store.TagResource(resourceARN, tags); err != nil {
			return nil, err
		}
	}

	return map[string]interface{}{}, nil
}

// UntagResource removes tags from an IoT Events resource.
func (s *IoTEventsService) UntagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	resourceARN := request.GetParamCaseInsensitive(req.Parameters, "resourceArn")
	if resourceARN == "" {
		return nil, iotstore.ErrInvalidRequest
	}

	if err := validateEventsResourceARN(store, resourceARN); err != nil {
		return nil, err
	}

	tagKeys := tagutil.ParseTagKeysAsSlice(req.Parameters, "tagKeys")
	if len(tagKeys) > 0 {
		if err := store.UntagResource(resourceARN, tagKeys); err != nil {
			return nil, err
		}
	}

	return map[string]interface{}{}, nil
}

// ListTagsForResource lists tags attached to an IoT Events resource.
func (s *IoTEventsService) ListTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	resourceARN := request.GetParamCaseInsensitive(req.Parameters, "resourceArn")
	if resourceARN == "" {
		return nil, iotstore.ErrInvalidRequest
	}

	if err := validateEventsResourceARN(store, resourceARN); err != nil {
		return nil, err
	}

	tags, err := store.ListTags(resourceARN)
	if err != nil {
		return nil, err
	}

	tagList := make([]map[string]interface{}, 0, len(tags))
	for k, v := range tags {
		tagList = append(tagList, map[string]interface{}{
			"key":   k,
			"value": v,
		})
	}

	return paginatedMaps("tags", tagList, req.Parameters), nil
}

func detectorModelConfig(dm *iotstore.DetectorModel) map[string]interface{} {
	return map[string]interface{}{
		"detectorModelName":        dm.DetectorModelName,
		"detectorModelArn":         dm.DetectorModelARN,
		"detectorModelDescription": dm.DetectorModelDescription,
		"roleArn":                  dm.RoleARN,
		"key":                      dm.Key,
		"evaluationMethod":         dm.EvaluationMethod,
		"detectorModelVersion":     dm.DetectorModelVersion,
		"status":                   dm.Status,
		"creationTime":             dm.CreationDate.Unix(),
		"lastUpdateTime":           dm.LastModifiedDate.Unix(),
	}
}

func inputConfig(input *iotstore.Input) map[string]interface{} {
	return map[string]interface{}{
		"inputName":        input.InputName,
		"inputArn":         input.InputARN,
		"inputDescription": input.InputDescription,
		"status":           input.Status,
		"creationTime":     input.CreationDate.Unix(),
		"lastUpdateTime":   input.LastModifiedDate.Unix(),
	}
}

func parseStructParam(params map[string]interface{}, key string) (map[string]interface{}, error) {
	if m := request.GetMapParamCaseInsensitive(params, key); m != nil {
		return m, nil
	}
	raw := request.GetParamCaseInsensitive(params, key)
	if raw == "" {
		return nil, nil
	}
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		return nil, awserrors.NewValidationException(fmt.Sprintf("invalid JSON in parameter %q: %v", key, err))
	}
	return result, nil
}
