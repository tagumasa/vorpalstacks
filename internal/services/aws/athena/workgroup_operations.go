package athena

import (
	"context"
	"regexp"
	"strings"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	athenastore "vorpalstacks/internal/store/aws/athena"
)

// arnRegex matches the three taggable Athena resource types per the Smithy
// TagResource documentation: "workgroups, data catalogs, or capacity
// reservations". Named queries and prepared statements are NOT taggable.
var arnRegex = regexp.MustCompile(`^arn:aws:athena:[^:]+:[^:]*:(workgroup|datacatalog|capacityreservation)/(.+)$`)

func normalizeAthenaARN(arn string, accountID string) string {
	parts := strings.SplitN(arn, ":", 6)
	if len(parts) >= 5 && parts[4] == "" {
		parts[4] = accountID
		return strings.Join(parts, ":")
	}
	return arn
}

// CreateWorkGroup creates a new workgroup in Athena.
// The workgroup is a container for queries and query results.
func (s *AthenaService) CreateWorkGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "Name")

	description := request.GetParamCaseInsensitive(req.Parameters, "Description")

	tags := tagutil.ToMap(tagutil.ParseTagsWithQueryFallback(req.Parameters, "Tags"))

	var cfg *WorkGroupConfigInput
	configMap := request.GetMapParamCaseInsensitive(req.Parameters, "Configuration")
	if configMap != nil {
		cfg = parseConfigMapToInput(configMap)
	} else {
		cfg = &WorkGroupConfigInput{
			OutputLocation: "s3://aws-athena-query-results-" + s.accountID + "-" + reqCtx.GetRegion() + "/",
		}
	}

	input := WorkGroupCreateInput{
		Name:        name,
		Description: description,
		Config:      cfg,
		Tags:        tags,
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := createWorkGroupCore(stores, input); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// GetWorkGroup retrieves the details of a specific workgroup.
func (s *AthenaService) GetWorkGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "WorkGroup")

	workGroup, err := s.getWorkGroupCore(reqCtx, name)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"WorkGroup": s.workGroupToResponse(workGroup),
	}, nil
}

// UpdateWorkGroup updates an existing workgroup's configuration.
func (s *AthenaService) UpdateWorkGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	input := UpdateWorkGroupInput{
		WorkGroup:            request.GetParamCaseInsensitive(req.Parameters, "WorkGroup"),
		Description:          request.GetParamCaseInsensitive(req.Parameters, "Description"),
		State:                request.GetParamCaseInsensitive(req.Parameters, "State"),
		ConfigurationUpdates: request.GetMapParamCaseInsensitive(req.Parameters, "ConfigurationUpdates"),
	}

	if err := s.updateWorkGroupCore(reqCtx, input); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// DeleteWorkGroup deletes a workgroup from Athena.
func (s *AthenaService) DeleteWorkGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "WorkGroup")

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := deleteWorkGroupCore(stores, name); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ListWorkGroups returns a list of workgroups in the account.
func (s *AthenaService) ListWorkGroups(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	maxResults, err := validateMaxResults(req.Parameters, 50, 1, 50)
	if err != nil {
		return nil, err
	}

	marker := request.GetParamCaseInsensitive(req.Parameters, "NextToken")

	result, err := listWorkGroupsCore(stores, maxResults, marker)
	if err != nil {
		return nil, err
	}

	workGroupSummaries := make([]map[string]interface{}, len(result.Items))
	for i, wg := range result.Items {
		workGroupSummaries[i] = map[string]interface{}{
			"Name":         wg.Name,
			"State":        wg.State,
			"Description":  wg.Description,
			"CreationTime": float64(wg.CreationTime.Unix()) + float64(wg.CreationTime.Nanosecond())/1e9,
		}
	}

	resp := map[string]interface{}{
		"WorkGroups": workGroupSummaries,
	}

	if result.NextMarker != "" {
		resp["NextToken"] = result.NextMarker
	}

	return resp, nil
}

// TagResource adds tags to an Athena resource.
func (s *AthenaService) TagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	input := TagResourceInput{
		ResourceARN: request.GetParamCaseInsensitive(req.Parameters, "ResourceARN"),
		Tags:        tagutil.ToMap(tagutil.ParseTagsWithQueryFallback(req.Parameters, "Tags")),
	}

	if err := s.tagResourceCore(reqCtx, input); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// UntagResource removes tags from an Athena resource.
func (s *AthenaService) UntagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	input := UntagResourceInput{
		ResourceARN: request.GetParamCaseInsensitive(req.Parameters, "ResourceARN"),
		TagKeys:     request.GetStringList(req.Parameters, "TagKeys"),
	}

	if err := s.untagResourceCore(reqCtx, input); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ListTagsForResource lists the tags associated with an Athena resource.
func (s *AthenaService) ListTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	maxResults, hasMaxResults := request.GetIntParamCaseInsensitive(req.Parameters, "MaxResults")
	input := ListTagsForResourceInput{
		ResourceARN:   request.GetParamCaseInsensitive(req.Parameters, "ResourceARN"),
		MaxResults:    maxResults,
		HasMaxResults: hasMaxResults,
		NextToken:     pagination.GetMarker(req.Parameters, "NextToken"),
	}

	tagList, nextMarker, err := s.listTagsForResourceCore(reqCtx, input)
	if err != nil {
		return nil, err
	}

	return pagination.BuildListResponse("Tags", tagList, nextMarker), nil
}

// parseConfigMapToInput converts the raw HTTP Configuration map into a
// WorkGroupConfigInput DTO. Validation is deferred to createWorkGroupCore
// so that both HTTP and admin handler share a single validation path.
func parseConfigMapToInput(config map[string]interface{}) *WorkGroupConfigInput {
	cfg := &WorkGroupConfigInput{}

	if resultConfigRaw, ok := config["ResultConfiguration"]; ok {
		if resultConfig, ok := resultConfigRaw.(map[string]interface{}); ok {
			if outputLocation, ok := resultConfig["OutputLocation"].(string); ok {
				cfg.OutputLocation = outputLocation
			}
		}
	}

	if enforce, ok := config["EnforceWorkGroupConfiguration"].(bool); ok {
		cfg.EnforceConfig = enforce
	}

	if publish, ok := config["PublishCloudWatchMetricsEnabled"].(bool); ok {
		cfg.PublishMetrics = publish
	}

	if bytesScanned, ok := config["BytesScannedCutoffPerQuery"].(float64); ok {
		v := int64(bytesScanned)
		cfg.BytesScannedCutoff = &v
	}

	if requesterPays, ok := config["RequesterPaysEnabled"].(bool); ok {
		cfg.RequesterPaysEnabled = requesterPays
	}

	if engineVersionRaw, ok := config["EngineVersion"]; ok {
		if engineVersion, ok := engineVersionRaw.(map[string]interface{}); ok {
			if selected, ok := engineVersion["SelectedEngineVersion"].(string); ok {
				cfg.EngineVersionSelected = selected
			}
			if effective, ok := engineVersion["EffectiveEngineVersion"].(string); ok {
				cfg.EngineVersionEffective = effective
			}
		}
	}

	if additional, ok := config["AdditionalConfiguration"].(string); ok {
		cfg.AdditionalConfiguration = additional
	}

	if executionRole, ok := config["ExecutionRole"].(string); ok {
		cfg.ExecutionRole = executionRole
	}

	if custEncMap, ok := config["CustomerContentEncryptionConfiguration"].(map[string]interface{}); ok {
		if kmsKey, ok := custEncMap["KmsKey"].(string); ok {
			cfg.CustomerContentEncryptionKmsKey = kmsKey
		}
	}

	if enableMin, ok := config["EnableMinimumEncryptionConfiguration"].(bool); ok {
		cfg.EnableMinimumEncryptionConfiguration = enableMin
	}

	return cfg
}

func (s *AthenaService) workGroupToResponse(wg *athenastore.WorkGroup) map[string]interface{} {
	response := map[string]interface{}{
		"Name":         wg.Name,
		"State":        wg.State,
		"Description":  wg.Description,
		"CreationTime": float64(wg.CreatedTime.Unix()) + float64(wg.CreatedTime.Nanosecond())/1e9,
	}

	if wg.Configuration != nil {
		response["Configuration"] = s.configurationToResponse(wg.Configuration)
	}

	return response
}

func (s *AthenaService) configurationToResponse(cfg *athenastore.WorkGroupConfiguration) map[string]interface{} {
	response := map[string]interface{}{
		"EnforceWorkGroupConfiguration":   cfg.EnforceWorkGroupConfiguration,
		"PublishCloudWatchMetricsEnabled": cfg.PublishCloudWatchMetricsEnabled,
		"RequesterPaysEnabled":            cfg.RequesterPaysEnabled,
	}

	if cfg.ResultConfiguration != nil {
		resultConfig := map[string]interface{}{
			"OutputLocation": cfg.ResultConfiguration.OutputLocation,
		}
		if cfg.ResultConfiguration.EncryptionConfiguration != nil {
			resultConfig["EncryptionConfiguration"] = map[string]interface{}{
				"EncryptionOption": cfg.ResultConfiguration.EncryptionConfiguration.EncryptionOption,
				"KmsKey":           cfg.ResultConfiguration.EncryptionConfiguration.KmsKey,
			}
		}
		if cfg.ResultConfiguration.ExpectedBucketOwner != "" {
			resultConfig["ExpectedBucketOwner"] = cfg.ResultConfiguration.ExpectedBucketOwner
		}
		if cfg.ResultConfiguration.ACLConfiguration != nil {
			resultConfig["AclConfiguration"] = map[string]interface{}{
				"S3AclOption": cfg.ResultConfiguration.ACLConfiguration.S3ACLOption,
			}
		}
		response["ResultConfiguration"] = resultConfig
	}

	if cfg.BytesScannedCutoffPerQuery > 0 {
		response["BytesScannedCutoffPerQuery"] = cfg.BytesScannedCutoffPerQuery
	}

	if cfg.EngineVersion != nil {
		response["EngineVersion"] = map[string]interface{}{
			"SelectedEngineVersion":  cfg.EngineVersion.SelectedEngineVersion,
			"EffectiveEngineVersion": cfg.EngineVersion.EffectiveEngineVersion,
		}
	}

	if cfg.AdditionalConfiguration != "" {
		response["AdditionalConfiguration"] = cfg.AdditionalConfiguration
	}

	if cfg.ExecutionRole != "" {
		response["ExecutionRole"] = cfg.ExecutionRole
	}

	if cfg.CustomerContentEncryptionConfiguration != nil {
		response["CustomerContentEncryptionConfiguration"] = map[string]interface{}{
			"KmsKey": cfg.CustomerContentEncryptionConfiguration.KmsKey,
		}
	}

	if cfg.EnableMinimumEncryptionConfiguration {
		response["EnableMinimumEncryptionConfiguration"] = cfg.EnableMinimumEncryptionConfiguration
	}

	return response
}
