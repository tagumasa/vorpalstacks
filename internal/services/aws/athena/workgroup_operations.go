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
func (s *Service) CreateWorkGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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
		configuration = s.parseWorkGroupConfiguration(configMap)
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
		if err := stores.workGroupStore.TagResource(arn, tags); err != nil {
			logs.Warn("Failed to tag workgroup", logs.String("workgroup", name), logs.Err(err))
		}
	}

	return response.EmptyResponse(), nil
}

// GetWorkGroup retrieves the details of a specific workgroup.
func (s *Service) GetWorkGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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
			return nil, ErrResourceNotFoundException
		}
		return nil, err
	}

	return map[string]interface{}{
		"WorkGroup": s.workGroupToResponse(workGroup),
	}, nil
}

// UpdateWorkGroup updates an existing workgroup's configuration.
func (s *Service) UpdateWorkGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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
			return nil, ErrResourceNotFoundException
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
		s.applyConfigurationUpdates(workGroup, configUpdatesMap)
	}

	if err := stores.workGroupStore.UpdateWorkGroup(workGroup); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// DeleteWorkGroup deletes a workgroup from Athena.
func (s *Service) DeleteWorkGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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
	if err := stores.workGroupStore.DeleteWorkGroup(name); err != nil {
		if err == athenastore.ErrWorkGroupNotFound {
			return nil, ErrResourceNotFoundException
		}
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ListWorkGroups returns a list of workgroups in the account.
func (s *Service) ListWorkGroups(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	workGroups, err := stores.workGroupStore.ListWorkGroups()
	if err != nil {
		return nil, err
	}

	maxResults := 50
	if maxStr := request.GetParamCaseInsensitive(req.Parameters, "MaxResults"); maxStr != "" {
		if val, err := strconv.Atoi(maxStr); err == nil && val > 0 {
			maxResults = val
		}
	}

	offset := 0
	if nextToken := request.GetParamCaseInsensitive(req.Parameters, "NextToken"); nextToken != "" {
		if val, err := strconv.Atoi(nextToken); err == nil && val >= 0 {
			offset = val
		}
	}

	var workGroupSummaries []map[string]interface{}
	for i, wg := range workGroups {
		if i < offset {
			continue
		}
		if len(workGroupSummaries) >= maxResults {
			break
		}
		workGroupSummaries = append(workGroupSummaries, map[string]interface{}{
			"Name":         wg.Name,
			"State":        wg.State,
			"Description":  wg.Description,
			"CreationTime": float64(wg.CreatedTime.Unix()) + float64(wg.CreatedTime.Nanosecond())/1e9,
		})
	}

	result := map[string]interface{}{
		"WorkGroups": workGroupSummaries,
	}

	if offset+maxResults < len(workGroups) {
		result["NextToken"] = strconv.Itoa(offset + maxResults)
	}

	return result, nil
}

func (s *Service) validateResourceExists(stores *athenaStores, resourceType, resourceName string) error {
	switch resourceType {
	case "workgroup":
		_, err := stores.workGroupStore.GetWorkGroup(resourceName)
		if err != nil {
			return ErrResourceNotFoundException
		}
	case "datacatalog":
		_, err := stores.dataCatalogStore.GetDataCatalog(resourceName)
		if err != nil {
			return ErrResourceNotFoundException
		}
	default:
		return ErrInvalidRequestException
	}
	return nil
}

// TagResource adds tags to an Athena resource.
func (s *Service) TagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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
			if err := stores.workGroupStore.TagResource(resourceArn, tags); err != nil {
				return nil, err
			}
		case "datacatalog":
			if err := stores.dataCatalogStore.TagResource(resourceArn, tags); err != nil {
				return nil, err
			}
		default:
			return nil, ErrInvalidRequestException
		}
	}

	return response.EmptyResponse(), nil
}

// UntagResource removes tags from an Athena resource.
func (s *Service) UntagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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
		if err := stores.workGroupStore.UntagResource(resourceArn, tagKeys); err != nil {
			return nil, err
		}
	}

	return response.EmptyResponse(), nil
}

// ListTagsForResource lists the tags associated with an Athena resource.
func (s *Service) ListTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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
		tags, err = stores.workGroupStore.ListTags(resourceArn)
	case "datacatalog":
		tags, err = stores.dataCatalogStore.ListTags(resourceArn)
	default:
		return nil, ErrInvalidRequestException
	}
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Tags":      tagutil.MapToResponse(tags),
		"NextToken": "",
	}, nil
}

func (s *Service) parseWorkGroupConfiguration(config map[string]interface{}) *athenastore.WorkGroupConfiguration {
	cfg := &athenastore.WorkGroupConfiguration{}

	if resultConfigRaw, ok := config["ResultConfiguration"]; ok {
		if resultConfig, ok := resultConfigRaw.(map[string]interface{}); ok {
			cfg.ResultConfiguration = s.parseResultConfiguration(resultConfig)
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

	return cfg
}

func (s *Service) parseEngineVersion(engineVersion map[string]interface{}) *athenastore.EngineVersion {
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

func (s *Service) applyConfigurationUpdates(workGroup *athenastore.WorkGroup, updates map[string]interface{}) {
	if workGroup.Configuration == nil {
		workGroup.Configuration = &athenastore.WorkGroupConfiguration{}
	}

	if resultConfigUpdatesRaw, ok := updates["ResultConfigurationUpdates"]; ok {
		if resultConfigUpdates, ok := resultConfigUpdatesRaw.(map[string]interface{}); ok {
			if workGroup.Configuration.ResultConfiguration == nil {
				workGroup.Configuration.ResultConfiguration = &athenastore.ResultConfiguration{}
			}
			if outputLocation, ok := resultConfigUpdates["OutputLocation"].(string); ok {
				workGroup.Configuration.ResultConfiguration.OutputLocation = outputLocation
			}
		}
	}

	if enforce, ok := updates["EnforceWorkGroupConfiguration"].(bool); ok {
		workGroup.Configuration.EnforceWorkGroupConfiguration = enforce
	}

	if bytesScanned, ok := updates["BytesScannedCutoffPerQuery"].(float64); ok {
		workGroup.Configuration.BytesScannedCutoffPerQuery = int64(bytesScanned)
	}

	if requesterPays, ok := updates["RequesterPaysEnabled"].(bool); ok {
		workGroup.Configuration.RequesterPaysEnabled = requesterPays
	}
}

func (s *Service) workGroupToResponse(wg *athenastore.WorkGroup) map[string]interface{} {
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

func (s *Service) configurationToResponse(cfg *athenastore.WorkGroupConfiguration) map[string]interface{} {
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

	return response
}

// GetDefaultWorkGroup returns the default primary workgroup, creating it if it does not exist.
func (s *Service) GetDefaultWorkGroup(reqCtx *request.RequestContext) (*athenastore.WorkGroup, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	wg, err := stores.workGroupStore.GetWorkGroup("primary")
	if err == nil {
		return wg, nil
	}

	if err == athenastore.ErrWorkGroupNotFound {
		defaultWg := &athenastore.WorkGroup{
			Name:        "primary",
			Description: "The default workgroup",
			Configuration: &athenastore.WorkGroupConfiguration{
				ResultConfiguration: &athenastore.ResultConfiguration{
					OutputLocation: "s3://aws-athena-query-results-" + s.accountID + "-" + reqCtx.GetRegion() + "/",
				},
				EngineVersion: &athenastore.EngineVersion{
					SelectedEngineVersion:  "AUTO",
					EffectiveEngineVersion: "Athena engine version 3",
				},
			},
			State: athenastore.WorkGroupStateEnabled,
		}
		if err := stores.workGroupStore.CreateWorkGroup(defaultWg); err != nil {
			logs.Warn("failed to create default primary workgroup", logs.Err(err))
		}
		return defaultWg, nil
	}

	return nil, err
}
