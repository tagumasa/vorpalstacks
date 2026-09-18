package cloudtrail

import (
	"context"
	"errors"
	"fmt"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	cloudtrailstore "vorpalstacks/internal/store/aws/cloudtrail"
)

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
		IAMValidator: reqCtx.GetIAMValidator(),
	}
	if v, ok := req.Parameters["S3KeyPrefix"]; ok {
		in.S3KeyPrefix = fmt.Sprintf("%v", v)
	}
	if v, ok := req.Parameters["SnsTopicName"]; ok {
		in.SnsTopicName = fmt.Sprintf("%v", v)
	}
	in.IncludeGlobalServiceEvents = boolParam(req.Parameters, "IncludeGlobalServiceEvents")
	in.IsMultiRegionTrail = boolParam(req.Parameters, "IsMultiRegionTrail")
	in.IsOrganizationTrail = boolParam(req.Parameters, "IsOrganizationTrail")
	in.EnableLogFileValidation = boolParam(req.Parameters, "EnableLogFileValidation")
	in.RecursiveLogging = boolParam(req.Parameters, "RecursiveLogging")
	if v, ok := req.Parameters["CloudWatchLogsLogGroupArn"]; ok {
		in.CloudWatchLogsLogGroupARN = fmt.Sprintf("%v", v)
	}
	if v, ok := req.Parameters["CloudWatchLogsRoleArn"]; ok {
		in.CloudWatchLogsRoleARN = fmt.Sprintf("%v", v)
	}
	if v, ok := req.Parameters["KmsKeyId"]; ok {
		in.KMSKeyID = fmt.Sprintf("%v", v)
	}

	created, err := s.createTrailCore(ctx, store, in)
	if err != nil {
		return nil, err
	}

	return s.formatTrailMutationResponse(created), nil
}

// DeleteTrail deletes the specified CloudTrail trail by name or ARN.
func (s *CloudTrailService) DeleteTrail(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := req.GetParam("Name")

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

	return s.formatTrailMutationResponse(trail), nil
}

// DescribeTrails retrieves information about the specified CloudTrail trails.
func (s *CloudTrailService) DescribeTrails(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	var names []string

	// The model spells the wire member "trailNameList" (lowercase first
	// letter); the shared list extractor owns the protocol spellings, so
	// presence is decided by the wire member itself, with the extracted
	// entries covering the flattened query-protocol forms as well.
	names = request.GetStringList(req.Parameters, "trailNameList")
	_, namesProvided := req.Parameters["trailNameList"]
	namesProvided = namesProvided || len(names) > 0

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

	// The Latest* delivery members are model-optional and report the
	// delivery machinery's recorded outcomes — absent until a real
	// delivery happened, never synthesised from the logging state. The
	// Date members serialise as epoch times; the Attempt*/TimeLogging*
	// members are model-typed String and carry RFC 3339 timestamps.
	result := map[string]interface{}{
		"IsLogging": trail.IsLogging,
	}

	if trail.StartedLoggingAt != nil {
		result["StartLoggingTime"] = float64(trail.StartedLoggingAt.Unix())
		result["TimeLoggingStarted"] = trail.StartedLoggingAt.UTC().Format(time.RFC3339)
	}
	if trail.StoppedLoggingAt != nil {
		result["StopLoggingTime"] = float64(trail.StoppedLoggingAt.Unix())
		result["TimeLoggingStopped"] = trail.StoppedLoggingAt.UTC().Format(time.RFC3339)
	}
	if trail.LatestDeliveryTime != nil {
		result["LatestDeliveryTime"] = float64(trail.LatestDeliveryTime.Unix())
	}
	if trail.LatestDeliveryError != "" {
		result["LatestDeliveryError"] = trail.LatestDeliveryError
	}
	if trail.LatestDeliveryAttemptTime != nil {
		result["LatestDeliveryAttemptTime"] = trail.LatestDeliveryAttemptTime.UTC().Format(time.RFC3339)
		if trail.LatestDeliveryAttemptSuccess {
			result["LatestDeliveryAttemptSucceeded"] = "true"
		} else {
			result["LatestDeliveryAttemptSucceeded"] = "false"
		}
	}
	if trail.LatestDigestTime != nil {
		result["LatestDigestDeliveryTime"] = float64(trail.LatestDigestTime.Unix())
	}
	if trail.LatestDigestError != "" {
		result["LatestDigestDeliveryError"] = trail.LatestDigestError
	}
	if trail.LatestCWLogsDeliveryTime != nil {
		result["LatestCloudWatchLogsDeliveryTime"] = float64(trail.LatestCWLogsDeliveryTime.Unix())
	}
	if trail.LatestCWLogsDeliveryError != "" {
		result["LatestCloudWatchLogsDeliveryError"] = trail.LatestCWLogsDeliveryError
	}
	if trail.LatestNotificationTime != nil {
		result["LatestNotificationTime"] = float64(trail.LatestNotificationTime.Unix())
	}
	if trail.LatestNotificationError != "" {
		result["LatestNotificationError"] = trail.LatestNotificationError
	}
	if trail.LatestNotificationAttemptTime != nil {
		result["LatestNotificationAttemptTime"] = trail.LatestNotificationAttemptTime.UTC().Format(time.RFC3339)
		if trail.LatestNotificationAttemptSuccess {
			result["LatestNotificationAttemptSucceeded"] = "true"
		} else {
			result["LatestNotificationAttemptSucceeded"] = "false"
		}
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

// formatTrail renders the Trail description shape (GetTrail,
// DescribeTrails items). It is NOT the CreateTrail/UpdateTrail response
// shape: the mutation responses carry the trail's mutable members only,
// and IsLogging belongs to GetTrailStatus, never to a Trail.
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
		"RecursiveLogging":           t.RecursiveLogging,
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
		result["SnsTopicARN"] = t.SnsTopicARN
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

// formatTrailMutationResponse renders the CreateTrail/UpdateTrail response
// shape: the trail's mutable members, without the HomeRegion and
// HasCustomEventSelectors/HasInsightSelectors description members.
func (s *CloudTrailService) formatTrailMutationResponse(t *cloudtrailstore.Trail) map[string]interface{} {
	result := s.formatTrail(t)
	delete(result, "HomeRegion")
	delete(result, "HasCustomEventSelectors")
	delete(result, "HasInsightSelectors")
	return result
}

func (s *CloudTrailService) mapStoreError(err error) error {
	if err == nil {
		return nil
	}
	// An already wire-shaped error — a Core guard refusal surfaced through
	// a store call — passes through unchanged.
	var awsErr *awserrors.AWSError
	if errors.As(err, &awsErr) {
		return err
	}
	mapped := awserrors.MapStoreError(err, storeErrorMappings)
	if mapped != err {
		return mapped
	}
	return ErrInternalError
}
