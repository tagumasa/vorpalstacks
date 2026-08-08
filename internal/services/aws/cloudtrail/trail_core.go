package cloudtrail

import (
	awserrors "vorpalstacks/internal/common/errors"
	tags "vorpalstacks/internal/common/tags"
	cloudtrailstore "vorpalstacks/internal/store/aws/cloudtrail"
	storecommon "vorpalstacks/internal/store/aws/common"
	"vorpalstacks/internal/utils/aws/types"
)

// ---------------------------------------------------------------------------
// Transport-agnostic Input structs
// ---------------------------------------------------------------------------

// CreateTrailInput carries every field that CreateTrail needs, in a format
// independent of the wire protocol (HTTP Query vs gRPC-Web). Both the HTTP
// API handler (trail_operations.go) and the admin gRPC handler
// (admin_handler.go) build this struct and delegate to createTrailCore.
type CreateTrailInput struct {
	Name                       string
	S3BucketName               string
	S3KeyPrefix                string
	SnsTopicName               string
	IncludeGlobalServiceEvents *bool
	IsMultiRegionTrail         *bool
	IsOrganizationTrail        *bool
	EnableLogFileValidation    *bool
	CloudWatchLogsLogGroupARN  string
	CloudWatchLogsRoleARN      string
	KMSKeyID                   string
	Tags                       []types.Tag
	Region                     string
}

// DeleteTrailInput carries the name or ARN of the trail to delete.
type DeleteTrailInput struct {
	NameOrARN string
}

// ListTrailsInput carries pagination parameters for listing trails.
type ListTrailsInput struct {
	NextToken string
	MaxItems  int
}

// TrailResult is the transport-agnostic result of createTrailCore.
type TrailResult struct {
	Name                       string
	TrailARN                   string
	S3BucketName               string
	S3KeyPrefix                string
	SnsTopicName               string
	SnsTopicARN                string
	IncludeGlobalServiceEvents bool
	IsMultiRegionTrail         bool
	IsOrganizationTrail        bool
	LogFileValidationEnabled   bool
	CloudWatchLogsLogGroupARN  string
	CloudWatchLogsRoleARN      string
	KMSKeyID                   string
	HomeRegion                 string
}

// ListTrailsResult is the transport-agnostic result of listTrailsCore.
type ListTrailsResult struct {
	Items     []TrailInfo
	NextToken string
}

// TrailInfo is the minimal trail information returned by listTrailsCore.
type TrailInfo struct {
	Name       string
	TrailARN   string
	HomeRegion string
}

// ---------------------------------------------------------------------------
// Core functions — single validation + persistence path
// ---------------------------------------------------------------------------

// createTrailCore is the single entry point for trail creation logic shared
// by the HTTP API and the admin gRPC handler. It performs all AWS-spec
// validation, constructs the trail via NewTrail (ensuring default
// EventSelectors), validates tags, and persists to the store.
//
// CloudWatchLogsRoleArn IAM validation is performed by the HTTP API handler
// before calling this function (it requires the request-context IAM
// validator which is not available in the admin handler).
func (s *CloudTrailService) createTrailCore(store cloudtrailstore.CloudTrailStoreInterface, in CreateTrailInput) (*cloudtrailstore.Trail, error) {
	if err := validateTrailName(in.Name); err != nil {
		return nil, err
	}
	if err := validateS3BucketName(in.S3BucketName); err != nil {
		return nil, err
	}
	if in.S3KeyPrefix != "" {
		if err := validateS3KeyPrefix(in.S3KeyPrefix); err != nil {
			return nil, err
		}
	}
	if in.SnsTopicName != "" {
		if err := validateSnsTopicName(in.SnsTopicName); err != nil {
			return nil, err
		}
	}
	if in.KMSKeyID != "" {
		if err := validateKMSKeyID(in.KMSKeyID); err != nil {
			return nil, err
		}
	}
	if in.CloudWatchLogsLogGroupARN != "" {
		if err := validateCloudWatchLogsLogGroupARN(in.CloudWatchLogsLogGroupARN); err != nil {
			return nil, err
		}
	}
	if in.CloudWatchLogsRoleARN != "" {
		if err := validateCloudWatchLogsRoleARN(in.CloudWatchLogsRoleARN); err != nil {
			return nil, err
		}
	}

	trail := cloudtrailstore.NewTrail(in.Name, in.S3BucketName, in.Region)

	if in.S3KeyPrefix != "" {
		trail.S3KeyPrefix = in.S3KeyPrefix
	}
	if in.SnsTopicName != "" {
		trail.SnsTopicName = in.SnsTopicName
	}
	if in.IncludeGlobalServiceEvents != nil {
		trail.IncludeGlobalServiceEvents = *in.IncludeGlobalServiceEvents
	}
	if in.IsMultiRegionTrail != nil {
		trail.IsMultiRegionTrail = *in.IsMultiRegionTrail
	}
	if in.IsOrganizationTrail != nil {
		trail.IsOrganizationTrail = *in.IsOrganizationTrail
	}
	if in.EnableLogFileValidation != nil {
		trail.LogFileValidationEnabled = *in.EnableLogFileValidation
	}
	if in.CloudWatchLogsLogGroupARN != "" {
		trail.CloudWatchLogsLogGroupARN = in.CloudWatchLogsLogGroupARN
	}
	if in.CloudWatchLogsRoleARN != "" {
		trail.CloudWatchLogsRoleARN = in.CloudWatchLogsRoleARN
	}
	if in.KMSKeyID != "" {
		trail.KMSKeyID = in.KMSKeyID
	}

	if len(in.Tags) > 0 {
		if err := validateCloudTrailTags(in.Tags); err != nil {
			return nil, err
		}
	}

	created, err := store.CreateTrail(trail)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	if len(in.Tags) > 0 {
		tagMap := make(map[string]string)
		for _, t := range in.Tags {
			tagMap[t.Key] = t.Value
		}
		if err := store.Tag(created.Name, tagMap); err != nil {
			return nil, s.mapStoreError(err)
		}
	}

	return created, nil
}

// deleteTrailCore is the single entry point for trail deletion logic shared
// by the HTTP API and the admin gRPC handler. It enforces the IsLogging
// precondition (AWS spec: "If the trail is currently logging, you must first
// call StopLogging"), cleans up the associated resource policy, and deletes
// the trail.
func (s *CloudTrailService) deleteTrailCore(store cloudtrailstore.CloudTrailStoreInterface, in DeleteTrailInput) error {
	trail, err := store.ResolveTrail(in.NameOrARN)
	if err != nil {
		return s.mapStoreError(err)
	}

	if trail.IsLogging {
		return awserrors.NewAWSError("OperationNotPermittedException",
			"Cannot delete a trail that is currently logging. Call StopLogging first.", 400)
	}

	// Best-effort cleanup of the associated resource policy and public
	// keys so that no orphaned material lingers after the trail is gone.
	_ = store.DeleteResourcePolicy(trail.TrailARN)
	_ = store.DeletePublicKeysByTrail(trail.Name)

	if err := store.DeleteTrail(trail.Name); err != nil {
		return s.mapStoreError(err)
	}

	return nil
}

// listTrailsCore is the single entry point for listing trails shared by the
// HTTP API and the admin gRPC handler.
func (s *CloudTrailService) listTrailsCore(store cloudtrailstore.CloudTrailStoreInterface, in ListTrailsInput) (*ListTrailsResult, error) {
	maxItems := in.MaxItems
	if maxItems <= 0 || maxItems > 1000 {
		maxItems = 1000
	}

	opts := storecommon.ListOptions{MaxItems: maxItems}
	if in.NextToken != "" {
		opts.Marker = in.NextToken
	}

	result, err := store.ListTrails(opts)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	items := make([]TrailInfo, 0, len(result.Items))
	for _, t := range result.Items {
		items = append(items, TrailInfo{
			Name:       t.Name,
			TrailARN:   t.TrailARN,
			HomeRegion: t.HomeRegion,
		})
	}

	return &ListTrailsResult{
		Items:     items,
		NextToken: result.NextMarker,
	}, nil
}

// trailToTrailResult converts a store Trail to a TrailResult for transport-
// agnostic consumption.
func trailToTrailResult(t *cloudtrailstore.Trail) *TrailResult {
	return &TrailResult{
		Name:                       t.Name,
		TrailARN:                   t.TrailARN,
		S3BucketName:               t.S3BucketName,
		S3KeyPrefix:                t.S3KeyPrefix,
		SnsTopicName:               t.SnsTopicName,
		SnsTopicARN:                t.SnsTopicARN,
		IncludeGlobalServiceEvents: t.IncludeGlobalServiceEvents,
		IsMultiRegionTrail:         t.IsMultiRegionTrail,
		IsOrganizationTrail:        t.IsOrganizationTrail,
		LogFileValidationEnabled:   t.LogFileValidationEnabled,
		CloudWatchLogsLogGroupARN:  t.CloudWatchLogsLogGroupARN,
		CloudWatchLogsRoleARN:      t.CloudWatchLogsRoleARN,
		KMSKeyID:                   t.KMSKeyID,
		HomeRegion:                 t.HomeRegion,
	}
}

// parseTagsFromParams converts a raw TagsList parameter (as produced by the
// request parser) into a []types.Tag.
func parseTagsFromParams(params map[string]interface{}) []types.Tag {
	return tags.ParseTags(params, "TagsList")
}
