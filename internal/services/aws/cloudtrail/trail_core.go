package cloudtrail

import (
	"context"
	"fmt"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/iam"
	tags "vorpalstacks/internal/common/tags"
	cloudtrailstore "vorpalstacks/internal/store/aws/cloudtrail"
	storecommon "vorpalstacks/internal/store/aws/common"
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
	Tags                       []tags.Tag
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

// TrailNameInput carries the trail name or ARN for single-trail operations
// (StartLogging, StopLogging).
type TrailNameInput struct {
	Name string
}

// UpdateTrailInput carries the raw update members for UpdateTrail. The
// members are presence-checked by the Core so that explicitly-provided empty
// strings clear the corresponding field (AWS spec behaviour), which requires
// the raw wire values rather than pre-formatted strings.
type UpdateTrailInput struct {
	Name                  string
	CloudWatchLogsRoleArn string
	IAMValidator          *iam.IAMValidator
	Params                map[string]interface{}
}

// DescribeTrailsInput carries the optional TrailNameList filter. When
// NamesProvided is true (even with an empty list) only the named trails are
// resolved; otherwise every trail is listed.
type DescribeTrailsInput struct {
	Names         []string
	NamesProvided bool
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
	// The model marks Name as required: an omitted name is a client error
	// on both planes and must be rejected before the store lookup so the
	// admin console does not surface a not-found for it.
	if in.NameOrARN == "" {
		return ErrInvalidParameter
	}
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

// resolveTrailCore resolves a trail by name or ARN, rejecting an empty
// selector with InvalidParameterException before the store lookup.
func (s *CloudTrailService) resolveTrailCore(store cloudtrailstore.CloudTrailStoreInterface, name string) (*cloudtrailstore.Trail, error) {
	if name == "" {
		return nil, ErrInvalidParameter
	}
	trail, err := store.ResolveTrail(name)
	if err != nil {
		return nil, s.mapStoreError(err)
	}
	return trail, nil
}

// updateTrailCore is the single entry point for UpdateTrail: it resolves the
// trail, validates the CloudWatchLogs role against IAM when one is supplied,
// applies the presence-checked update members, and persists the result.
func (s *CloudTrailService) updateTrailCore(ctx context.Context, store cloudtrailstore.CloudTrailStoreInterface, in UpdateTrailInput) (*cloudtrailstore.Trail, error) {
	if in.Name == "" {
		return nil, ErrInvalidParameter
	}

	trail, err := store.ResolveTrail(in.Name)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	if in.IAMValidator != nil && in.CloudWatchLogsRoleArn != "" {
		if err := in.IAMValidator.ValidateRoleForService(ctx, in.CloudWatchLogsRoleArn, iam.ServicePrincipalCloudTrail); err != nil {
			return nil, err
		}
	}

	applyTrailUpdates(trail, in.Params)

	if err := store.UpdateTrail(trail); err != nil {
		return nil, s.mapStoreError(err)
	}

	return trail, nil
}

// describeTrailsCore is the single entry point for DescribeTrails: when a
// TrailNameList was supplied only those trails are resolved (unresolvable
// names are silently skipped, per AWS behaviour), otherwise every trail in
// the store is returned.
func (s *CloudTrailService) describeTrailsCore(store cloudtrailstore.CloudTrailStoreInterface, in DescribeTrailsInput) ([]*cloudtrailstore.Trail, error) {
	if in.NamesProvided {
		var trails []*cloudtrailstore.Trail
		for _, name := range in.Names {
			trail, err := store.ResolveTrail(name)
			if err != nil {
				continue
			}
			trails = append(trails, trail)
		}
		return trails, nil
	}
	trails, err := listAllTrails(store)
	if err != nil {
		return nil, s.mapStoreError(err)
	}
	return trails, nil
}

// startLoggingCore is the single entry point for StartLogging.
func (s *CloudTrailService) startLoggingCore(store cloudtrailstore.CloudTrailStoreInterface, in TrailNameInput) error {
	if in.Name == "" {
		return ErrInvalidParameter
	}
	if err := store.StartLogging(in.Name); err != nil {
		return s.mapStoreError(err)
	}
	return nil
}

// stopLoggingCore is the single entry point for StopLogging.
func (s *CloudTrailService) stopLoggingCore(store cloudtrailstore.CloudTrailStoreInterface, in TrailNameInput) error {
	if in.Name == "" {
		return ErrInvalidParameter
	}
	if err := store.StopLogging(in.Name); err != nil {
		return s.mapStoreError(err)
	}
	return nil
}

// listAllTrails paginates through all trails across multiple pages.
func listAllTrails(store cloudtrailstore.CloudTrailStoreInterface) ([]*cloudtrailstore.Trail, error) {
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

// applyTrailUpdates applies UpdateTrail parameters using existence checks so
// that explicitly-provided empty strings clear the field (AWS spec behaviour).
func applyTrailUpdates(trail *cloudtrailstore.Trail, params map[string]interface{}) {
	if v, ok := params["S3BucketName"]; ok {
		trail.S3BucketName = fmt.Sprintf("%v", v)
	}
	if v, ok := params["S3KeyPrefix"]; ok {
		trail.S3KeyPrefix = fmt.Sprintf("%v", v)
	}
	if v, ok := params["SnsTopicName"]; ok {
		trail.SnsTopicName = fmt.Sprintf("%v", v)
	}
	if v, ok := params["SnsTopicArn"]; ok {
		trail.SnsTopicARN = fmt.Sprintf("%v", v)
	}
	if b := resolveBool(params, "IncludeGlobalServiceEvents"); b != nil {
		trail.IncludeGlobalServiceEvents = *b
	}
	if b := resolveBool(params, "IsMultiRegionTrail"); b != nil {
		trail.IsMultiRegionTrail = *b
	}
	if b := resolveBool(params, "IsOrganizationTrail"); b != nil {
		trail.IsOrganizationTrail = *b
	}
	if b := resolveBool(params, "EnableLogFileValidation"); b != nil {
		trail.LogFileValidationEnabled = *b
	}
	if v, ok := params["CloudWatchLogsLogGroupArn"]; ok {
		trail.CloudWatchLogsLogGroupARN = fmt.Sprintf("%v", v)
	}
	if v, ok := params["CloudWatchLogsRoleArn"]; ok {
		trail.CloudWatchLogsRoleARN = fmt.Sprintf("%v", v)
	}
	if v, ok := params["KmsKeyId"]; ok {
		trail.KMSKeyID = fmt.Sprintf("%v", v)
	}
}

// parseTagsFromParams converts a raw TagsList parameter (as produced by the
// request parser) into a []tags.Tag.
func parseTagsFromParams(params map[string]interface{}) []tags.Tag {
	return tags.ParseTags(params, "TagsList")
}
