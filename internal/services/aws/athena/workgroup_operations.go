package athena

import (
	"context"
	"regexp"
	"sort"
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
	if name == "" {
		return nil, ErrInvalidRequestException
	}

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
	if name == "" {
		return nil, ErrInvalidRequestException
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	workGroup, err := stores.workGroupStore.GetWorkGroup(name)
	if err != nil {
		if err == athenastore.ErrWorkGroupNotFound {
			return nil, workGroupNotFound(name)
		}
		return nil, err
	}

	return map[string]interface{}{
		"WorkGroup": s.workGroupToResponse(workGroup),
	}, nil
}

// UpdateWorkGroup updates an existing workgroup's configuration.
func (s *AthenaService) UpdateWorkGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "WorkGroup")
	if name == "" {
		return nil, ErrInvalidRequestException
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	workGroup, err := stores.workGroupStore.GetWorkGroup(name)
	if err != nil {
		if err == athenastore.ErrWorkGroupNotFound {
			return nil, workGroupNotFound(name)
		}
		return nil, err
	}

	description := request.GetParamCaseInsensitive(req.Parameters, "Description")
	if description != "" {
		if err := validateWorkGroupDescriptionString(description); err != nil {
			return nil, err
		}
		workGroup.Description = description
	}

	state := request.GetParamCaseInsensitive(req.Parameters, "State")
	if state != "" {
		if err := validateWorkGroupState(state); err != nil {
			return nil, err
		}
		workGroup.State = athenastore.WorkGroupState(state)
	}

	configUpdatesMap := request.GetMapParamCaseInsensitive(req.Parameters, "ConfigurationUpdates")
	if configUpdatesMap != nil {
		if err := s.applyConfigurationUpdates(workGroup, configUpdatesMap); err != nil {
			return nil, err
		}
	}

	if err := stores.workGroupStore.UpdateWorkGroup(workGroup); err != nil {
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

func (s *AthenaService) validateResourceExists(stores *athenaStores, resourceType, resourceName string) error {
	switch resourceType {
	case "workgroup":
		_, err := stores.workGroupStore.GetWorkGroup(resourceName)
		if err != nil {
			return workGroupNotFound(resourceName)
		}
	case "datacatalog":
		_, err := stores.dataCatalogStore.GetDataCatalog(resourceName)
		if err != nil {
			return dataCatalogNotFound(resourceName)
		}
	case "capacityreservation":
		_, err := stores.capacityReservationStore.GetCapacityReservation(resourceName)
		if err != nil {
			return capacityReservationNotFound(resourceName)
		}
	default:
		return ErrInvalidRequestException
	}
	return nil
}

// TagResource adds tags to an Athena resource.
func (s *AthenaService) TagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceArn := request.GetParamCaseInsensitive(req.Parameters, "ResourceARN")
	if resourceArn == "" {
		return nil, ErrInvalidRequestException
	}

	matches := arnRegex.FindStringSubmatch(resourceArn)
	if matches == nil {
		return nil, ErrInvalidRequestException
	}

	resourceType := matches[1]
	resourceName := matches[2]

	tags := tagutil.ToMap(tagutil.ParseTagsWithQueryFallback(req.Parameters, "Tags"))
	if err := validateTags(tags); err != nil {
		return nil, err
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.validateResourceExists(stores, resourceType, resourceName); err != nil {
		return nil, err
	}

	resourceArn = normalizeAthenaARN(resourceArn, s.accountID)

	if len(tags) > 0 {
		switch resourceType {
		case "workgroup":
			if err := stores.workGroupStore.Tag(resourceArn, tags); err != nil {
				return nil, err
			}
		case "datacatalog":
			if err := stores.dataCatalogStore.Tag(resourceArn, tags); err != nil {
				return nil, err
			}
		case "capacityreservation":
			if err := stores.capacityReservationStore.Tag(resourceArn, tags); err != nil {
				return nil, err
			}
		default:
			return nil, ErrInvalidRequestException
		}
	}

	return response.EmptyResponse(), nil
}

// UntagResource removes tags from an Athena resource.
func (s *AthenaService) UntagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceArn := request.GetParamCaseInsensitive(req.Parameters, "ResourceARN")
	if resourceArn == "" {
		return nil, ErrInvalidRequestException
	}

	matches := arnRegex.FindStringSubmatch(resourceArn)
	if matches == nil {
		return nil, ErrInvalidRequestException
	}

	tagKeys := request.GetStringList(req.Parameters, "TagKeys")

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.validateResourceExists(stores, matches[1], matches[2]); err != nil {
		return nil, err
	}

	resourceArn = normalizeAthenaARN(resourceArn, s.accountID)

	if len(tagKeys) > 0 {
		switch matches[1] {
		case "workgroup":
			if err := stores.workGroupStore.Untag(resourceArn, tagKeys); err != nil {
				return nil, err
			}
		case "datacatalog":
			if err := stores.dataCatalogStore.Untag(resourceArn, tagKeys); err != nil {
				return nil, err
			}
		case "capacityreservation":
			if err := stores.capacityReservationStore.Untag(resourceArn, tagKeys); err != nil {
				return nil, err
			}
		default:
			return nil, ErrInvalidRequestException
		}
	}

	return response.EmptyResponse(), nil
}

// ListTagsForResource lists the tags associated with an Athena resource.
func (s *AthenaService) ListTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceArn := request.GetParamCaseInsensitive(req.Parameters, "ResourceARN")
	if resourceArn == "" {
		return nil, ErrInvalidRequestException
	}

	matches := arnRegex.FindStringSubmatch(resourceArn)
	if matches == nil {
		return nil, ErrInvalidRequestException
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.validateResourceExists(stores, matches[1], matches[2]); err != nil {
		return nil, err
	}

	resourceArn = normalizeAthenaARN(resourceArn, s.accountID)

	var tags map[string]string
	switch matches[1] {
	case "workgroup":
		tags, err = stores.workGroupStore.List(resourceArn)
	case "datacatalog":
		tags, err = stores.dataCatalogStore.List(resourceArn)
	case "capacityreservation":
		tags, err = stores.capacityReservationStore.List(resourceArn)
	default:
		return nil, ErrInvalidRequestException
	}
	if err != nil {
		return nil, err
	}

	tagList := tagutil.MapToResponse(tags)
	sort.Slice(tagList, func(i, j int) bool {
		return tagList[i]["Key"].(string) < tagList[j]["Key"].(string)
	})

	maxResults, err := validateMaxResults(req.Parameters, len(tagList), 75, pagination.AbsoluteMaxItems)
	if err != nil {
		return nil, err
	}

	marker := pagination.GetMarker(req.Parameters, "NextToken")
	pageResult := pagination.PaginateSlice(tagList, marker, maxResults, func(item map[string]interface{}) string {
		return item["Key"].(string)
	})

	return pagination.BuildListResponse("Tags", pageResult.Items, pageResult.NextMarker), nil
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

func (s *AthenaService) parseEngineVersion(engineVersion map[string]interface{}) *athenastore.EngineVersion {
	ev := &athenastore.EngineVersion{}
	if selected, ok := engineVersion["SelectedEngineVersion"].(string); ok {
		ev.SelectedEngineVersion = selected
	}
	if effective, ok := engineVersion["EffectiveEngineVersion"].(string); ok {
		ev.EffectiveEngineVersion = effective
	}
	if ev.SelectedEngineVersion == "" {
		ev.SelectedEngineVersion = "AUTO"
	}
	if ev.EffectiveEngineVersion == "" {
		ev.EffectiveEngineVersion = "Athena engine version 3"
	}
	return ev
}

func (s *AthenaService) applyConfigurationUpdates(workGroup *athenastore.WorkGroup, updates map[string]interface{}) error {
	if workGroup.Configuration == nil {
		workGroup.Configuration = &athenastore.WorkGroupConfiguration{}
	}

	if resultConfigUpdatesRaw, ok := updates["ResultConfigurationUpdates"]; ok {
		if resultConfigUpdates, ok := resultConfigUpdatesRaw.(map[string]interface{}); ok {
			if workGroup.Configuration.ResultConfiguration == nil {
				workGroup.Configuration.ResultConfiguration = &athenastore.ResultConfiguration{}
			}
			rc := workGroup.Configuration.ResultConfiguration

			if outputLocation, ok := resultConfigUpdates["OutputLocation"].(string); ok {
				rc.OutputLocation = outputLocation
			}
			if encConfigMap, ok := resultConfigUpdates["EncryptionConfiguration"].(map[string]interface{}); ok {
				rc.EncryptionConfiguration = &athenastore.EncryptionConfiguration{}
				if encOption, ok := encConfigMap["EncryptionOption"].(string); ok {
					rc.EncryptionConfiguration.EncryptionOption = encOption
				}
				if kmsKey, ok := encConfigMap["KmsKey"].(string); ok {
					rc.EncryptionConfiguration.KmsKey = kmsKey
				}
			}
			if expectedBucketOwner, ok := resultConfigUpdates["ExpectedBucketOwner"].(string); ok {
				rc.ExpectedBucketOwner = expectedBucketOwner
			}
			if aclConfigMap, ok := resultConfigUpdates["AclConfiguration"].(map[string]interface{}); ok {
				aclOption, _ := aclConfigMap["S3AclOption"].(string)
				if aclOption != "BUCKET_OWNER_FULL_CONTROL" {
					return ErrInvalidParameterException
				}
				rc.ACLConfiguration = &athenastore.ACLConfiguration{S3ACLOption: aclOption}
			}

			if remove, ok := resultConfigUpdates["RemoveOutputLocation"].(bool); ok && remove {
				rc.OutputLocation = ""
			}
			if remove, ok := resultConfigUpdates["RemoveEncryptionConfiguration"].(bool); ok && remove {
				rc.EncryptionConfiguration = nil
			}
			if remove, ok := resultConfigUpdates["RemoveExpectedBucketOwner"].(bool); ok && remove {
				rc.ExpectedBucketOwner = ""
			}
			if remove, ok := resultConfigUpdates["RemoveAclConfiguration"].(bool); ok && remove {
				rc.ACLConfiguration = nil
			}
		}
	}

	if enforce, ok := updates["EnforceWorkGroupConfiguration"].(bool); ok {
		workGroup.Configuration.EnforceWorkGroupConfiguration = enforce
	}

	if bytesScanned, ok := updates["BytesScannedCutoffPerQuery"].(float64); ok {
		workGroup.Configuration.BytesScannedCutoffPerQuery = int64(bytesScanned)
		if err := validateBytesScannedCutoff(workGroup.Configuration.BytesScannedCutoffPerQuery); err != nil {
			return err
		}
	}

	if remove, ok := updates["RemoveBytesScannedCutoffPerQuery"].(bool); ok && remove {
		workGroup.Configuration.BytesScannedCutoffPerQuery = 0
	}

	if requesterPays, ok := updates["RequesterPaysEnabled"].(bool); ok {
		workGroup.Configuration.RequesterPaysEnabled = requesterPays
	}

	if publish, ok := updates["PublishCloudWatchMetricsEnabled"].(bool); ok {
		workGroup.Configuration.PublishCloudWatchMetricsEnabled = publish
	}

	if engineVersionRaw, ok := updates["EngineVersion"]; ok {
		if engineVersion, ok := engineVersionRaw.(map[string]interface{}); ok {
			workGroup.Configuration.EngineVersion = s.parseEngineVersion(engineVersion)
		}
	}

	if additional, ok := updates["AdditionalConfiguration"].(string); ok {
		if err := validateAdditionalConfiguration(additional); err != nil {
			return err
		}
		workGroup.Configuration.AdditionalConfiguration = additional
	}

	if executionRole, ok := updates["ExecutionRole"].(string); ok {
		if err := validateExecutionRole(executionRole); err != nil {
			return err
		}
		workGroup.Configuration.ExecutionRole = executionRole
	}

	if custEncMap, ok := updates["CustomerContentEncryptionConfiguration"].(map[string]interface{}); ok {
		workGroup.Configuration.CustomerContentEncryptionConfiguration = &athenastore.CustomerContentEncryptionConfiguration{}
		if kmsKey, ok := custEncMap["KmsKey"].(string); ok {
			workGroup.Configuration.CustomerContentEncryptionConfiguration.KmsKey = kmsKey
		}
	}

	if remove, ok := updates["RemoveCustomerContentEncryptionConfiguration"].(bool); ok && remove {
		workGroup.Configuration.CustomerContentEncryptionConfiguration = nil
	}

	if enableMin, ok := updates["EnableMinimumEncryptionConfiguration"].(bool); ok {
		workGroup.Configuration.EnableMinimumEncryptionConfiguration = enableMin
	}

	return nil
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
