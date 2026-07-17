package athena

import (
	"context"
	"regexp"
	"strconv"
	"strings"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/core/logs"
	athenastore "vorpalstacks/internal/store/aws/athena"
	storecommon "vorpalstacks/internal/store/aws/common"
)

var workGroupNameRegex = regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)

var arnRegex = regexp.MustCompile(`^arn:aws:athena:[^:]+:[^:]*:(workgroup|datacatalog|namedquery|preparedstatement)/(.+)$`)

func normalizeAthenaARN(arn string, accountID string) string {
	parts := strings.SplitN(arn, ":", 6)
	if len(parts) >= 5 && parts[4] == "" {
		parts[4] = accountID
		return strings.Join(parts, ":")
	}
	return arn
}

func validateWorkGroupName(name string) error {
	if len(name) < 1 || len(name) > 32 {
		return ErrInvalidParameterException
	}
	if !workGroupNameRegex.MatchString(name) {
		return ErrInvalidParameterException
	}
	return nil
}

// CreateWorkGroup creates a new workgroup in Athena.
// The workgroup is a container for queries and query results.
func (s *AthenaService) CreateWorkGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "Name")
	if name == "" {
		return nil, ErrInvalidRequestException
	}

	if err := validateWorkGroupName(name); err != nil {
		return nil, err
	}

	description := request.GetParamCaseInsensitive(req.Parameters, "Description")

	configMap := request.GetMapParamCaseInsensitive(req.Parameters, "Configuration")
	var configuration *athenastore.WorkGroupConfiguration
	if configMap != nil {
		var err error
		configuration, err = s.parseWorkGroupConfiguration(configMap)
		if err != nil {
			return nil, err
		}
	}

	if configuration == nil {
		configuration = &athenastore.WorkGroupConfiguration{
			ResultConfiguration: &athenastore.ResultConfiguration{
				OutputLocation: "s3://aws-athena-query-results-" + s.accountID + "-" + reqCtx.GetRegion() + "/",
			},
		}
	}

	tags := tagutil.ToMap(tagutil.ParseTagsWithQueryFallback(req.Parameters, "Tags"))

	workGroup := &athenastore.WorkGroup{
		Name:          name,
		Description:   description,
		Configuration: configuration,
		State:         athenastore.WorkGroupStateEnabled,
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := stores.workGroupStore.CreateWorkGroup(workGroup); err != nil {
		if err == athenastore.ErrWorkGroupAlreadyExists {
			return nil, ErrResourceAlreadyExistsException
		}
		return nil, err
	}

	if len(tags) > 0 {
		arn := stores.workGroupStore.GetARN(name)
		if err := stores.workGroupStore.Tag(arn, tags); err != nil {
			logs.Warn("Failed to tag workgroup", logs.String("workgroup", name), logs.Err(err))
		}
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
		workGroup.Description = description
	}

	state := request.GetParamCaseInsensitive(req.Parameters, "State")
	if state != "" {
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
	if name == "" {
		return nil, ErrInvalidRequestException
	}

	if name == "primary" {
		return nil, ErrInvalidRequestException
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	stores.namedQueryStore.DeleteNamedQueriesByWorkGroup(name)
	stores.preparedStatementStore.DeletePreparedStatementsByWorkGroup(name)
	deletedQEIds, _ := stores.queryExecutionStore.DeleteQueryExecutionsByWorkGroup(name)
	stores.resultStore.DeleteResultsByIDs(deletedQEIds)

	if err := stores.workGroupStore.DeleteWorkGroup(name); err != nil {
		if err == athenastore.ErrWorkGroupNotFound {
			return nil, workGroupNotFound(name)
		}
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

	maxResults := 50
	if maxStr := request.GetParamCaseInsensitive(req.Parameters, "MaxResults"); maxStr != "" {
		if val, err := strconv.Atoi(maxStr); err == nil && val > 0 {
			maxResults = val
		}
	}

	var marker string
	if nextToken := request.GetParamCaseInsensitive(req.Parameters, "NextToken"); nextToken != "" {
		marker = nextToken
	}

	result, err := stores.workGroupStore.ListWorkGroups(storecommon.ListOptions{
		Marker:   marker,
		MaxItems: maxResults,
	})
	if err != nil {
		return nil, err
	}

	workGroupSummaries := make([]map[string]interface{}, len(result.Items))
	for i, wg := range result.Items {
		workGroupSummaries[i] = map[string]interface{}{
			"Name":         wg.Name,
			"State":        wg.State,
			"Description":  wg.Description,
			"CreationTime": float64(wg.CreatedTime.Unix()) + float64(wg.CreatedTime.Nanosecond())/1e9,
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
	default:
		return nil, ErrInvalidRequestException
	}
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Tags": tagutil.MapToResponse(tags),
	}, nil
}

func (s *AthenaService) parseWorkGroupConfiguration(config map[string]interface{}) (*athenastore.WorkGroupConfiguration, error) {
	cfg := &athenastore.WorkGroupConfiguration{}

	if resultConfigRaw, ok := config["ResultConfiguration"]; ok {
		if resultConfig, ok := resultConfigRaw.(map[string]interface{}); ok {
			rc, err := s.parseResultConfiguration(resultConfig)
			if err != nil {
				return nil, err
			}
			cfg.ResultConfiguration = rc
		}
	}

	if enforce, ok := config["EnforceWorkGroupConfiguration"].(bool); ok {
		cfg.EnforceWorkGroupConfiguration = enforce
	}

	if publish, ok := config["PublishCloudWatchMetricsEnabled"].(bool); ok {
		cfg.PublishCloudWatchMetricsEnabled = publish
	}

	if bytesScanned, ok := config["BytesScannedCutoffPerQuery"].(float64); ok {
		cfg.BytesScannedCutoffPerQuery = int64(bytesScanned)
	}

	if requesterPays, ok := config["RequesterPaysEnabled"].(bool); ok {
		cfg.RequesterPaysEnabled = requesterPays
	}

	if engineVersionRaw, ok := config["EngineVersion"]; ok {
		if engineVersion, ok := engineVersionRaw.(map[string]interface{}); ok {
			cfg.EngineVersion = s.parseEngineVersion(engineVersion)
		}
	}

	if additional, ok := config["AdditionalConfiguration"].(string); ok {
		cfg.AdditionalConfiguration = additional
	}

	if executionRole, ok := config["ExecutionRole"].(string); ok {
		cfg.ExecutionRole = executionRole
	}

	if custEncMap, ok := config["CustomerContentEncryptionConfiguration"].(map[string]interface{}); ok {
		cfg.CustomerContentEncryptionConfiguration = &athenastore.CustomerContentEncryptionConfiguration{}
		if kmsKey, ok := custEncMap["KmsKey"].(string); ok {
			cfg.CustomerContentEncryptionConfiguration.KmsKey = kmsKey
		}
	}

	if enableMin, ok := config["EnableMinimumEncryptionConfiguration"].(bool); ok {
		cfg.EnableMinimumEncryptionConfiguration = enableMin
	}

	return cfg, nil
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
		workGroup.Configuration.AdditionalConfiguration = additional
	}

	if executionRole, ok := updates["ExecutionRole"].(string); ok {
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
