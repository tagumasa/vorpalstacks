package cloudtrail

import (
	"context"
	"strings"
	"time"

	"vorpalstacks/internal/common/iam"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tags "vorpalstacks/internal/common/tags"
	cloudtrailstore "vorpalstacks/internal/store/aws/cloudtrail"
)

func getParam(req *request.ParsedRequest, key string) string {
	return request.GetParamLowerFirst(req.Parameters, key)
}

// CreateTrail creates a new CloudTrail trail.
func (s *CloudTrailService) CreateTrail(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := getParam(req, "Name")
	s3BucketName := getParam(req, "S3BucketName")

	if name == "" {
		return nil, ErrInvalidParameter
	}

	trail := cloudtrailstore.NewTrail(name, s3BucketName, reqCtx.GetRegion())

	if s3KeyPrefix := getParam(req, "S3KeyPrefix"); s3KeyPrefix != "" {
		trail.S3KeyPrefix = s3KeyPrefix
	}
	if snsTopicName := getParam(req, "SnsTopicName"); snsTopicName != "" {
		trail.SnsTopicName = snsTopicName
	}
	if snsTopicARN := getParam(req, "SnsTopicArn"); snsTopicARN != "" {
		trail.SnsTopicARN = snsTopicARN
	}
	if v := req.Parameters["IncludeGlobalServiceEvents"]; v != nil {
		if b, ok := v.(bool); ok {
			trail.IncludeGlobalServiceEvents = b
		} else if s := getParam(req, "IncludeGlobalServiceEvents"); s == "false" {
			trail.IncludeGlobalServiceEvents = false
		}
	}
	if v := req.Parameters["IsMultiRegionTrail"]; v != nil {
		if b, ok := v.(bool); ok {
			trail.IsMultiRegionTrail = b
		} else if s := getParam(req, "IsMultiRegionTrail"); s == "true" {
			trail.IsMultiRegionTrail = true
		}
	}
	if v := req.Parameters["IsOrganizationTrail"]; v != nil {
		if b, ok := v.(bool); ok {
			trail.IsOrganizationTrail = b
		} else if s := getParam(req, "IsOrganizationTrail"); s == "true" {
			trail.IsOrganizationTrail = true
		}
	}
	if v := req.Parameters["EnableLogFileValidation"]; v != nil {
		if b, ok := v.(bool); ok {
			trail.LogFileValidationEnabled = b
		} else if s := getParam(req, "EnableLogFileValidation"); s == "true" {
			trail.LogFileValidationEnabled = true
		}
	}
	if cwLogGroupARN := getParam(req, "CloudWatchLogsLogGroupArn"); cwLogGroupARN != "" {
		trail.CloudWatchLogsLogGroupARN = cwLogGroupARN
	}
	if cwRoleARN := getParam(req, "CloudWatchLogsRoleArn"); cwRoleARN != "" {
		validator := reqCtx.GetIAMValidator()
		if err := validator.ValidateRoleForService(ctx, cwRoleARN, iam.ServicePrincipalCloudTrail); err != nil {
			return nil, err
		}
		trail.CloudWatchLogsRoleARN = cwRoleARN
	}
	if kmsKeyID := getParam(req, "KmsKeyId"); kmsKeyID != "" {
		trail.KMSKeyID = kmsKeyID
	}

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
		if err := store.TagResource(created.Name, tagMap); err != nil {
			return nil, s.mapStoreError(err)
		}
	}

	return s.formatTrail(created), nil
}

// DeleteTrail deletes the specified CloudTrail trail.
func (s *CloudTrailService) DeleteTrail(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := getParam(req, "Name")

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
	name := getParam(req, "Name")

	if name == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	trail, err := store.GetTrail(name)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	if s3BucketName := getParam(req, "S3BucketName"); s3BucketName != "" {
		trail.S3BucketName = s3BucketName
	}
	if s3KeyPrefix := getParam(req, "S3KeyPrefix"); s3KeyPrefix != "" {
		trail.S3KeyPrefix = s3KeyPrefix
	}
	if snsTopicName := getParam(req, "SnsTopicName"); snsTopicName != "" {
		trail.SnsTopicName = snsTopicName
	}
	if snsTopicARN := getParam(req, "SnsTopicArn"); snsTopicARN != "" {
		trail.SnsTopicARN = snsTopicARN
	}
	if v := req.Parameters["IncludeGlobalServiceEvents"]; v != nil {
		if b, ok := v.(bool); ok {
			trail.IncludeGlobalServiceEvents = b
		} else if s := getParam(req, "IncludeGlobalServiceEvents"); s != "" {
			trail.IncludeGlobalServiceEvents = s == "true"
		}
	}
	if v := req.Parameters["IsMultiRegionTrail"]; v != nil {
		if b, ok := v.(bool); ok {
			trail.IsMultiRegionTrail = b
		} else if s := getParam(req, "IsMultiRegionTrail"); s != "" {
			trail.IsMultiRegionTrail = s == "true"
		}
	}
	if v := req.Parameters["IsOrganizationTrail"]; v != nil {
		if b, ok := v.(bool); ok {
			trail.IsOrganizationTrail = b
		} else if s := getParam(req, "IsOrganizationTrail"); s != "" {
			trail.IsOrganizationTrail = s == "true"
		}
	}
	if v := req.Parameters["EnableLogFileValidation"]; v != nil {
		if b, ok := v.(bool); ok {
			trail.LogFileValidationEnabled = b
		} else if s := getParam(req, "EnableLogFileValidation"); s != "" {
			trail.LogFileValidationEnabled = s == "true"
		}
	}
	if cwLogGroupARN := getParam(req, "CloudWatchLogsLogGroupArn"); cwLogGroupARN != "" {
		trail.CloudWatchLogsLogGroupARN = cwLogGroupARN
	}
	if cwRoleARN := getParam(req, "CloudWatchLogsRoleArn"); cwRoleARN != "" {
		validator := reqCtx.GetIAMValidator()
		if err := validator.ValidateRoleForService(ctx, cwRoleARN, iam.ServicePrincipalCloudTrail); err != nil {
			return nil, err
		}
		trail.CloudWatchLogsRoleARN = cwRoleARN
	}
	if kmsKeyID := getParam(req, "KmsKeyId"); kmsKeyID != "" {
		trail.KMSKeyID = kmsKeyID
	}

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
					var trail *cloudtrailstore.Trail
					var err error
					if strings.Contains(nameStr, ":trail/") {
						trail, err = store.GetTrailByARN(nameStr)
					} else {
						trail, err = store.GetTrail(nameStr)
					}
					if err != nil {
						continue
					}
					trails = append(trails, trail)
				}
			}
		}
	} else {
		trails, err = store.ListTrails()
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
	name := getParam(req, "Name")

	if name == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	var trail *cloudtrailstore.Trail
	if strings.Contains(name, ":trail/") {
		trail, err = store.GetTrailByARN(name)
	} else {
		trail, err = store.GetTrail(name)
	}
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	return map[string]interface{}{
		"Trail": s.formatTrail(trail),
	}, nil
}

// GetTrailStatus retrieves the status of the specified CloudTrail trail.
func (s *CloudTrailService) GetTrailStatus(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := getParam(req, "Name")

	if name == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	var trail *cloudtrailstore.Trail
	if strings.Contains(name, ":trail/") {
		trail, err = store.GetTrailByARN(name)
	} else {
		trail, err = store.GetTrail(name)
	}
	if err != nil {
		return nil, s.mapStoreError(err)
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

	trails, err := store.ListTrails()
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	resumeAfter := ""
	if nextToken := getParam(req, "NextToken"); nextToken != "" {
		resumeAfter = nextToken
	}

	formattedTrails := make([]map[string]interface{}, 0)
	skipping := resumeAfter != ""
	for _, t := range trails {
		if skipping {
			if t.Name == resumeAfter {
				skipping = false
			}
			continue
		}
		formattedTrails = append(formattedTrails, map[string]interface{}{
			"TrailARN":   t.TrailARN,
			"Name":       t.Name,
			"HomeRegion": t.HomeRegion,
		})
	}

	nextToken := ""
	if len(formattedTrails) >= 1000 {
		nextToken = formattedTrails[len(formattedTrails)-1]["Name"].(string)
		formattedTrails = formattedTrails[:1000]
	}

	result := map[string]interface{}{
		"Trails": formattedTrails,
	}
	if nextToken != "" {
		result["NextToken"] = nextToken
	}

	return result, nil
}

// StartLogging starts recording AWS API calls for a trail.
func (s *CloudTrailService) StartLogging(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := getParam(req, "Name")

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
	name := getParam(req, "Name")

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
