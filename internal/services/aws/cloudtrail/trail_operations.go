package cloudtrail

import (
	"context"
	"fmt"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/iam"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	cloudtrailstore "vorpalstacks/internal/store/aws/cloudtrail"
)

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

// CreateTrail creates a new CloudTrail trail.
func (s *CloudTrailService) CreateTrail(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	in := CreateTrailInput{
		Name:         req.GetParam("Name"),
		S3BucketName: req.GetParam("S3BucketName"),
		Region:       reqCtx.GetRegion(),
		Tags:         parseTagsFromParams(req.Parameters),
	}
	if v, ok := req.Parameters["S3KeyPrefix"]; ok {
		in.S3KeyPrefix = fmt.Sprintf("%v", v)
	}
	if v, ok := req.Parameters["SnsTopicName"]; ok {
		in.SnsTopicName = fmt.Sprintf("%v", v)
	}
	in.IncludeGlobalServiceEvents = resolveBool(req.Parameters, "IncludeGlobalServiceEvents")
	in.IsMultiRegionTrail = resolveBool(req.Parameters, "IsMultiRegionTrail")
	in.IsOrganizationTrail = resolveBool(req.Parameters, "IsOrganizationTrail")
	in.EnableLogFileValidation = resolveBool(req.Parameters, "EnableLogFileValidation")
	if v, ok := req.Parameters["CloudWatchLogsLogGroupArn"]; ok {
		in.CloudWatchLogsLogGroupARN = fmt.Sprintf("%v", v)
	}
	if v, ok := req.Parameters["CloudWatchLogsRoleArn"]; ok {
		in.CloudWatchLogsRoleARN = fmt.Sprintf("%v", v)
	}
	if v, ok := req.Parameters["KmsKeyId"]; ok {
		in.KMSKeyID = fmt.Sprintf("%v", v)
	}

	// IAM role validation (HTTP API only — requires request context).
	if in.CloudWatchLogsRoleARN != "" {
		validator := reqCtx.GetIAMValidator()
		if err := validator.ValidateRoleForService(ctx, in.CloudWatchLogsRoleARN, iam.ServicePrincipalCloudTrail); err != nil {
			return nil, err
		}
	}

	created, err := s.createTrailCore(store, in)
	if err != nil {
		return nil, err
	}

	return s.formatTrail(created), nil
}

// DeleteTrail deletes the specified CloudTrail trail by name or ARN.
func (s *CloudTrailService) DeleteTrail(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := req.GetParam("Name")
	if name == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	if err := s.deleteTrailCore(store, DeleteTrailInput{NameOrARN: name}); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// UpdateTrail updates the settings for a CloudTrail trail.
func (s *CloudTrailService) UpdateTrail(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	trail, err := s.updateTrailCore(ctx, store, UpdateTrailInput{
		Name:                  req.GetParam("Name"),
		CloudWatchLogsRoleArn: req.GetParam("CloudWatchLogsRoleArn"),
		IAMValidator:          reqCtx.GetIAMValidator(),
		Params:                req.Parameters,
	})
	if err != nil {
		return nil, err
	}

	return s.formatTrail(trail), nil
}

// DescribeTrails retrieves information about the specified CloudTrail trails.
func (s *CloudTrailService) DescribeTrails(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	var names []string
	namesProvided := false

	trailNameListRaw := req.Parameters["TrailNameList"]
	if trailNameListRaw == nil {
		trailNameListRaw = req.Parameters["trailNameList"]
	}

	if trailNameListRaw != nil {
		namesProvided = true
		if arr, ok := trailNameListRaw.([]interface{}); ok && len(arr) > 0 {
			for _, name := range arr {
				if nameStr, ok := name.(string); ok {
					names = append(names, nameStr)
				}
			}
		}
	}

	trails, err := s.describeTrailsCore(store, DescribeTrailsInput{
		Names:         names,
		NamesProvided: namesProvided,
	})
	if err != nil {
		return nil, err
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
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	trail, err := s.resolveTrailCore(store, req.GetParam("Name"))
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Trail": s.formatTrail(trail),
	}, nil
}

// GetTrailStatus retrieves the status of the specified CloudTrail trail.
func (s *CloudTrailService) GetTrailStatus(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	trail, err := s.resolveTrailCore(store, req.GetParam("Name"))
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

	coreResult, err := s.listTrailsCore(store, ListTrailsInput{
		NextToken: req.GetParam("NextToken"),
	})
	if err != nil {
		return nil, err
	}

	formattedTrails := make([]map[string]interface{}, 0, len(coreResult.Items))
	for _, t := range coreResult.Items {
		formattedTrails = append(formattedTrails, map[string]interface{}{
			"TrailARN":   t.TrailARN,
			"Name":       t.Name,
			"HomeRegion": t.HomeRegion,
		})
	}

	result := map[string]interface{}{
		"Trails": formattedTrails,
	}
	if coreResult.NextToken != "" {
		result["NextToken"] = coreResult.NextToken
	}

	return result, nil
}

// StartLogging starts recording AWS API calls for a trail.
func (s *CloudTrailService) StartLogging(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	if err := s.startLoggingCore(store, TrailNameInput{Name: req.GetParam("Name")}); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// StopLogging stops recording AWS API calls for a trail.
func (s *CloudTrailService) StopLogging(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	if err := s.stopLoggingCore(store, TrailNameInput{Name: req.GetParam("Name")}); err != nil {
		return nil, err
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
	mapped := awserrors.MapStoreError(err, storeErrorMappings)
	if mapped != err {
		return mapped
	}
	return ErrInternalError
}
