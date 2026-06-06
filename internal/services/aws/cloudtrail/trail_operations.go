package cloudtrail

import (
	"context"
	"time"

	"vorpalstacks/internal/common/iam"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tags "vorpalstacks/internal/common/tags"
	cloudtrailstore "vorpalstacks/internal/store/aws/cloudtrail"
	storecommon "vorpalstacks/internal/store/aws/common"
)

// listAllTrails paginates through all trails across multiple pages.
func (s *CloudTrailService) listAllTrails(store cloudtrailstore.CloudTrailStoreInterface) ([]*cloudtrailstore.Trail, error) {
	var allTrails []*cloudtrailstore.Trail
	var marker string
	for {
		opts := storecommon.ListOptions{MaxItems: 1000}
		if marker != "" {
			opts.Marker = marker
		}
		result, err := store.ListTrails(opts)
		if err != nil {
			return nil, err
		}
		allTrails = append(allTrails, result.Items...)
		if result.NextMarker == "" {
			break
		}
		marker = result.NextMarker
	}
	return allTrails, nil
}

// resolveTrailFromRequest resolves a trail by the Name or TrailArn parameter.
func (s *CloudTrailService) resolveTrailFromRequest(reqCtx *request.RequestContext, req *request.ParsedRequest, paramName string) (cloudtrailstore.CloudTrailStoreInterface, *cloudtrailstore.Trail, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, nil, s.mapStoreError(err)
	}

	name := req.GetParam(paramName)
	if name == "" {
		return nil, nil, ErrInvalidParameter
	}

	trail, err := store.ResolveTrail(name)
	if err != nil {
		return nil, nil, s.mapStoreError(err)
	}
	return store, trail, nil
}

// resolveBool extracts a boolean from request parameters, accepting both
// bool and string ("true"/"false") representations.
func resolveBool(params map[string]interface{}, key string) *bool {
	v := params[key]
	if v == nil {
		return nil
	}
	if b, ok := v.(bool); ok {
		return &b
	}
	if s, ok := v.(string); ok {
		val := s == "true"
		return &val
	}
	return nil
}

// applyTrailUpdates copies non-empty fields from request parameters into the trail.
func applyTrailUpdates(trail *cloudtrailstore.Trail, req *request.ParsedRequest) {
	if v := req.GetParam("S3BucketName"); v != "" {
		trail.S3BucketName = v
	}
	if v := req.GetParam("S3KeyPrefix"); v != "" {
		trail.S3KeyPrefix = v
	}
	if v := req.GetParam("SnsTopicName"); v != "" {
		trail.SnsTopicName = v
	}
	if v := req.GetParam("SnsTopicArn"); v != "" {
		trail.SnsTopicARN = v
	}
	if b := resolveBool(req.Parameters, "IncludeGlobalServiceEvents"); b != nil {
		trail.IncludeGlobalServiceEvents = *b
	}
	if b := resolveBool(req.Parameters, "IsMultiRegionTrail"); b != nil {
		trail.IsMultiRegionTrail = *b
	}
	if b := resolveBool(req.Parameters, "IsOrganizationTrail"); b != nil {
		trail.IsOrganizationTrail = *b
	}
	if b := resolveBool(req.Parameters, "EnableLogFileValidation"); b != nil {
		trail.LogFileValidationEnabled = *b
	}
	if v := req.GetParam("CloudWatchLogsLogGroupArn"); v != "" {
		trail.CloudWatchLogsLogGroupARN = v
	}
	if v := req.GetParam("CloudWatchLogsRoleArn"); v != "" {
		trail.CloudWatchLogsRoleARN = v
	}
	if v := req.GetParam("KmsKeyId"); v != "" {
		trail.KMSKeyID = v
	}
}

// CreateTrail creates a new CloudTrail trail.
func (s *CloudTrailService) CreateTrail(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := req.GetParam("Name")
	if name == "" {
		return nil, ErrInvalidParameter
	}

	trail := cloudtrailstore.NewTrail(name, req.GetParam("S3BucketName"), reqCtx.GetRegion())

	if cwRoleARN := req.GetParam("CloudWatchLogsRoleArn"); cwRoleARN != "" {
		validator := reqCtx.GetIAMValidator()
		if err := validator.ValidateRoleForService(ctx, cwRoleARN, iam.ServicePrincipalCloudTrail); err != nil {
			return nil, err
		}
	}

	applyTrailUpdates(trail, req)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	created, err := store.CreateTrail(trail)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	if tagList := tags.ParseTags(req.Parameters, "TagsList"); len(tagList) > 0 {
		tagMap := make(map[string]string)
		for _, t := range tagList {
			tagMap[t.Key] = t.Value
		}
		if err := store.Tag(created.Name, tagMap); err != nil {
			return nil, s.mapStoreError(err)
		}
	}

	return s.formatTrail(created), nil
}

// DeleteTrail deletes the specified CloudTrail trail.
func (s *CloudTrailService) DeleteTrail(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := req.GetParam("Name")
	if name == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	if err := store.DeleteTrail(name); err != nil {
		return nil, s.mapStoreError(err)
	}

	return response.EmptyResponse(), nil
}

// UpdateTrail updates the settings for a CloudTrail trail.
func (s *CloudTrailService) UpdateTrail(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := req.GetParam("Name")
	if name == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	trail, err := store.ResolveTrail(name)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	if cwRoleARN := req.GetParam("CloudWatchLogsRoleArn"); cwRoleARN != "" {
		validator := reqCtx.GetIAMValidator()
		if err := validator.ValidateRoleForService(ctx, cwRoleARN, iam.ServicePrincipalCloudTrail); err != nil {
			return nil, err
		}
	}

	applyTrailUpdates(trail, req)

	if err := store.UpdateTrail(trail); err != nil {
		return nil, s.mapStoreError(err)
	}

	return s.formatTrail(trail), nil
}

// DescribeTrails retrieves information about the specified CloudTrail trails.
func (s *CloudTrailService) DescribeTrails(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	var trails []*cloudtrailstore.Trail

	trailNameListRaw := req.Parameters["TrailNameList"]
	if trailNameListRaw == nil {
		trailNameListRaw = req.Parameters["trailNameList"]
	}

	if trailNameListRaw != nil {
		if arr, ok := trailNameListRaw.([]interface{}); ok && len(arr) > 0 {
			for _, name := range arr {
				if nameStr, ok := name.(string); ok {
					trail, err := store.ResolveTrail(nameStr)
					if err != nil {
						continue
					}
					trails = append(trails, trail)
				}
			}
		}
	} else {
		trails, err = s.listAllTrails(store)
		if err != nil {
			return nil, s.mapStoreError(err)
		}
	}

	formattedTrails := make([]map[string]interface{}, 0)
	for _, t := range trails {
		formattedTrails = append(formattedTrails, s.formatTrail(t))
	}

	return map[string]interface{}{
		"trailList": formattedTrails,
	}, nil
}

// GetTrail retrieves the settings for the specified CloudTrail trail.
func (s *CloudTrailService) GetTrail(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_, trail, err := s.resolveTrailFromRequest(reqCtx, req, "Name")
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Trail": s.formatTrail(trail),
	}, nil
}

// GetTrailStatus retrieves the status of the specified CloudTrail trail.
func (s *CloudTrailService) GetTrailStatus(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_, trail, err := s.resolveTrailFromRequest(reqCtx, req, "Name")
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"IsLogging":           trail.IsLogging,
		"LatestDeliveryError": "",
	}

	if trail.StartedLoggingAt != nil {
		result["StartLoggingTime"] = float64(trail.StartedLoggingAt.Unix())
	}
	if trail.StoppedLoggingAt != nil {
		result["StopLoggingTime"] = float64(trail.StoppedLoggingAt.Unix())
	}
	if trail.IsLogging {
		result["LatestDeliveryTime"] = float64(time.Now().UTC().Unix())
	}

	return result, nil
}

// ListTrails retrieves all CloudTrail trails for the account.
func (s *CloudTrailService) ListTrails(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	opts := storecommon.ListOptions{MaxItems: 1000}
	if nextToken := req.GetParam("NextToken"); nextToken != "" {
		opts.Marker = nextToken
	}

	ctResult, err := store.ListTrails(opts)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	formattedTrails := make([]map[string]interface{}, 0, len(ctResult.Items))
	for _, t := range ctResult.Items {
		formattedTrails = append(formattedTrails, map[string]interface{}{
			"TrailARN":   t.TrailARN,
			"Name":       t.Name,
			"HomeRegion": t.HomeRegion,
		})
	}

	result := map[string]interface{}{
		"Trails": formattedTrails,
	}
	if ctResult.NextMarker != "" {
		result["NextToken"] = ctResult.NextMarker
	}

	return result, nil
}

// StartLogging starts recording AWS API calls for a trail.
func (s *CloudTrailService) StartLogging(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := req.GetParam("Name")
	if name == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	if err := store.StartLogging(name); err != nil {
		return nil, s.mapStoreError(err)
	}

	return response.EmptyResponse(), nil
}

// StopLogging stops recording AWS API calls for a trail.
func (s *CloudTrailService) StopLogging(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := req.GetParam("Name")
	if name == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	if err := store.StopLogging(name); err != nil {
		return nil, s.mapStoreError(err)
	}

	return response.EmptyResponse(), nil
}

func (s *CloudTrailService) formatTrail(t *cloudtrailstore.Trail) map[string]interface{} {
	result := map[string]interface{}{
		"Name":                       t.Name,
		"TrailARN":                   t.TrailARN,
		"IncludeGlobalServiceEvents": t.IncludeGlobalServiceEvents,
		"IsMultiRegionTrail":         t.IsMultiRegionTrail,
		"HomeRegion":                 t.HomeRegion,
		"HasCustomEventSelectors":    t.HasCustomEventSelectors,
		"HasInsightSelectors":        t.HasInsightSelectors,
		"IsOrganizationTrail":        t.IsOrganizationTrail,
		"LogFileValidationEnabled":   t.LogFileValidationEnabled,
		"IsLogging":                  t.IsLogging,
	}

	if t.S3BucketName != "" {
		result["S3BucketName"] = t.S3BucketName
	}
	if t.S3KeyPrefix != "" {
		result["S3KeyPrefix"] = t.S3KeyPrefix
	}
	if t.SnsTopicName != "" {
		result["SnsTopicName"] = t.SnsTopicName
	}
	if t.SnsTopicARN != "" {
		result["SnsTopicArn"] = t.SnsTopicARN
	}
	if t.CloudWatchLogsLogGroupARN != "" {
		result["CloudWatchLogsLogGroupArn"] = t.CloudWatchLogsLogGroupARN
	}
	if t.CloudWatchLogsRoleARN != "" {
		result["CloudWatchLogsRoleArn"] = t.CloudWatchLogsRoleARN
	}
	if t.KMSKeyID != "" {
		result["KmsKeyId"] = t.KMSKeyID
	}

	return result
}

func (s *CloudTrailService) mapStoreError(err error) error {
	if err == nil {
		return nil
	}
	switch err {
	case cloudtrailstore.ErrTrailNotFound:
		return ErrTrailNotFound
	case cloudtrailstore.ErrTrailAlreadyExists:
		return ErrTrailAlreadyExists
	default:
		return ErrInternalError
	}
}
