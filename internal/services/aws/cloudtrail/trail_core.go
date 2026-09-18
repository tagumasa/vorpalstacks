package cloudtrail

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"vorpalstacks/internal/common/iam"
	tags "vorpalstacks/internal/common/tags"
	cloudtrailstore "vorpalstacks/internal/store/aws/cloudtrail"
	storecommon "vorpalstacks/internal/store/aws/common"
	"vorpalstacks/internal/utils/aws/arn"
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
	RecursiveLogging           *bool
	CloudWatchLogsLogGroupARN  string
	CloudWatchLogsRoleARN      string
	KMSKeyID                   string
	Tags                       []tags.Tag
	Region                     string
	IAMValidator               *iam.IAMValidator
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
// EventSelectors), validates tags, and persists to the store. The
// CloudWatchLogsRoleArn trust validation runs here when a validator is
// injected: both planes supply one (the HTTP API from its request context,
// the admin console from the service's role provider), and a nil validator
// leaves the member unvalidated, mirroring updateTrailCore.
func (s *CloudTrailService) createTrailCore(ctx context.Context, store cloudtrailstore.CloudTrailStoreInterface, in CreateTrailInput) (*cloudtrailstore.Trail, error) {
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
		// The member carries a topic name or ARN; the reference form is
		// validated here with the other syntax checks, the topic's
		// existence and policy in the destination validation below.
		if err := validateSnsTopicReference(in.SnsTopicName); err != nil {
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
	if in.IAMValidator != nil && in.CloudWatchLogsRoleARN != "" {
		if err := in.IAMValidator.ValidateRoleForService(ctx, in.CloudWatchLogsRoleARN, iam.ServicePrincipalCloudTrail); err != nil {
			return nil, err
		}
	}

	// The trail quota is enforced inside CreateTrail (count and write under
	// the store mutex): "Trails per Region — 5" (Quotas in AWS CloudTrail);
	// the model declares MaximumNumberOfTrailsExceededException on
	// CreateTrail for the sixth trail.

	// Destination validation: the S3 bucket must exist with a sufficient
	// policy, and a configured SNS topic must resolve (name or ARN) to an
	// existing topic whose policy grants CloudTrail publish access. The
	// resolved topic ARN is what the trail carries and the responses echo.
	accountID := store.GetAccountID()
	region := store.GetRegion()
	trailARN := arn.NewARNBuilder(accountID, region).CloudTrail().Trail(in.Name)
	if err := s.validateTrailS3Destination(ctx, region, accountID, in.S3BucketName, in.S3KeyPrefix, trailARN); err != nil {
		return nil, err
	}

	trail := cloudtrailstore.NewTrail(in.Name, in.S3BucketName, in.Region)

	if in.S3KeyPrefix != "" {
		trail.S3KeyPrefix = in.S3KeyPrefix
	}
	if in.SnsTopicName != "" {
		topicARN, err := s.resolveTrailSnsTopic(ctx, region, accountID, in.SnsTopicName, trailARN)
		if err != nil {
			return nil, err
		}
		trail.SnsTopicName = in.SnsTopicName
		trail.SnsTopicARN = topicARN
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
	if in.RecursiveLogging != nil {
		trail.RecursiveLogging = *in.RecursiveLogging
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
		if errors.Is(err, cloudtrailstore.ErrTrailQuotaExceeded) {
			return nil, newMaximumNumberOfTrailsExceededException(
				fmt.Sprintf("The maximum number of trails per Region (%d) has been reached", cloudtrailstore.MaxTrailsPerRegion))
		}
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
// by the HTTP API and the admin gRPC handler. The DeleteTrail reference
// documents no logging precondition — deleting a multi-Region trail stops
// logging of events in all Regions — so an actively logging trail deletes
// like any other. It cleans up the associated resource policy and deletes
// the trail.
func (s *CloudTrailService) deleteTrailCore(store cloudtrailstore.CloudTrailStoreInterface, in DeleteTrailInput) error {
	// The model marks Name as required: an omitted name is a client error
	// on both planes and must be rejected before the store lookup so the
	// admin console does not surface a not-found for it. DeleteTrail
	// declares InvalidTrailNameException, not the generic parameter error.
	if in.NameOrARN == "" {
		return ErrInvalidTrailName
	}
	trail, err := s.resolveTrailCore(store, in.NameOrARN)
	if err != nil {
		return err
	}

	// Best-effort cleanup of the associated resource policy, public keys
	// and event configuration so that no orphaned material lingers after
	// the trail is gone; a cleanup failure is logged but never blocks the
	// deletion.
	if err := store.DeleteResourcePolicy(trail.TrailARN); err != nil {
		slog.Warn("cloudtrail: trail resource-policy cleanup failed",
			"trail", trail.Name, "error", err)
	}
	if err := store.DeletePublicKeysByTrail(trail.Name); err != nil {
		slog.Warn("cloudtrail: trail public-key cleanup failed",
			"trail", trail.Name, "error", err)
	}
	if err := store.DeleteEventConfiguration(trail.Name, ""); err != nil {
		slog.Warn("cloudtrail: trail event-configuration cleanup failed",
			"trail", trail.Name, "error", err)
	}

	if err := store.DeleteTrail(trail.Name); err != nil {
		return s.mapStoreError(err)
	}

	return nil
}

// listTrailsCore is the single entry point for listing trails shared by the
// HTTP API and the admin gRPC handler.
func (s *CloudTrailService) listTrailsCore(store cloudtrailstore.CloudTrailStoreInterface, in ListTrailsInput) (*ListTrailsResult, error) {
	maxItems := in.MaxItems
	if maxItems <= 0 || maxItems > cloudtrailstore.MaxListTrailsMaxItems {
		maxItems = cloudtrailstore.DefaultListTrailsMaxItems
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
// selector with InvalidTrailNameException before the store lookup — the
// error the single-trail operations declare for a bad name. A selector
// that carries the ARN prefix but is not a well-formed trail ARN answers
// the declared CloudTrailARNInvalidException instead of surfacing as a
// not-found trail.
func (s *CloudTrailService) resolveTrailCore(store cloudtrailstore.CloudTrailStoreInterface, name string) (*cloudtrailstore.Trail, error) {
	if name == "" {
		return nil, ErrInvalidTrailName
	}
	if strings.HasPrefix(name, "arn:") && !isCloudTrailResourceARN(name, "trail/") {
		return nil, newCloudTrailARNInvalidException(
			fmt.Sprintf("The specified trail ARN is not valid: %s", name))
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
		return nil, ErrInvalidTrailName
	}

	trail, err := s.resolveTrailCore(store, in.Name)
	if err != nil {
		return nil, err
	}

	if in.IAMValidator != nil && in.CloudWatchLogsRoleArn != "" {
		if err := in.IAMValidator.ValidateRoleForService(ctx, in.CloudWatchLogsRoleArn, iam.ServicePrincipalCloudTrail); err != nil {
			return nil, err
		}
	}

	// Destination validation for the touched members: a bucket or prefix
	// change re-validates the S3 policy against the effective post-update
	// delivery path, and a topic member (or a topic carried without a
	// resolved ARN) resolves the topic and checks its policy. The external
	// checks run before the mutation so a refusal persists nothing. The
	// members' own syntax runs first — a malformed bucket name is rejected
	// as a name, never as a missing destination.
	accountID := store.GetAccountID()
	region := store.GetRegion()
	if v, ok := in.Params["S3BucketName"]; ok {
		if err := validateS3BucketName(fmt.Sprintf("%v", v)); err != nil {
			return nil, err
		}
	}
	if v, ok := in.Params["S3KeyPrefix"]; ok {
		if prefix := fmt.Sprintf("%v", v); prefix != "" {
			if err := validateS3KeyPrefix(prefix); err != nil {
				return nil, err
			}
		}
	}
	_, bucketTouched := in.Params["S3BucketName"]
	_, prefixTouched := in.Params["S3KeyPrefix"]
	if bucketTouched || prefixTouched {
		bucket := trail.S3BucketName
		if bucketTouched {
			bucket = fmt.Sprintf("%v", in.Params["S3BucketName"])
		}
		prefix := trail.S3KeyPrefix
		if prefixTouched {
			prefix = fmt.Sprintf("%v", in.Params["S3KeyPrefix"])
		}
		if err := s.validateTrailS3Destination(ctx, region, accountID, bucket, prefix, trail.TrailARN); err != nil {
			return nil, err
		}
	}
	topicName := ""
	resolveTopic := false
	if raw, touched := in.Params["SnsTopicName"]; touched {
		topicName = fmt.Sprintf("%v", raw)
		resolveTopic = topicName != ""
	} else if trail.SnsTopicName != "" && trail.SnsTopicARN == "" {
		// A topic recorded before SNS resolution existed carries no ARN;
		// any update heals it.
		topicName = trail.SnsTopicName
		resolveTopic = true
	}
	resolvedSnsTopicARN := ""
	if resolveTopic {
		if err := validateSnsTopicReference(topicName); err != nil {
			return nil, err
		}
		arnStr, err := s.resolveTrailSnsTopic(ctx, region, accountID, topicName, trail.TrailARN)
		if err != nil {
			return nil, err
		}
		resolvedSnsTopicARN = arnStr
	}

	// The mutation runs as one load-apply-persist step under the store
	// lock, so it cannot lose a concurrent writer's change (the delivery
	// worker's bookkeeping). The pre-mutation validation flag is captured
	// inside the closure: the mutex serialises the load, so a second
	// enable-validation update sees the first's write and generates no
	// second key.
	validationEnabledBefore := false
	updated, err := store.MutateTrail(trail.Name, func(t *cloudtrailstore.Trail) error {
		validationEnabledBefore = t.LogFileValidationEnabled
		return applyTrailUpdates(t, in.Params, resolvedSnsTopicARN)
	})
	if err != nil {
		return nil, err
	}

	// Enabling log file validation provisions the key material the flag
	// promises: CreateTrail generates the key at creation, so the update
	// path generates one on the same transition (false→true). Re-asserting
	// the flag on an already-validated trail keeps the existing key.
	if !validationEnabledBefore && updated.LogFileValidationEnabled {
		if _, err := store.GenerateAndStorePublicKey(updated.Name); err != nil {
			return nil, fmt.Errorf("failed to generate public key for trail: %w", err)
		}
	}

	return updated, nil
}

// describeTrailsCore is the single entry point for DescribeTrails: when a
// TrailNameList was supplied only those trails are resolved (unresolvable
// names are silently skipped, per AWS behaviour), otherwise every trail in
// the store is returned. An entry that carries the ARN prefix but is not
// a well-formed trail ARN answers the declared CloudTrailARNInvalidException.
func (s *CloudTrailService) describeTrailsCore(store cloudtrailstore.CloudTrailStoreInterface, in DescribeTrailsInput) ([]*cloudtrailstore.Trail, error) {
	if in.NamesProvided {
		var trails []*cloudtrailstore.Trail
		for _, name := range in.Names {
			if strings.HasPrefix(name, "arn:") && !isCloudTrailResourceARN(name, "trail/") {
				return nil, newCloudTrailARNInvalidException(
					fmt.Sprintf("The specified trail ARN is not valid: %s", name))
			}
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
		return ErrInvalidTrailName
	}
	if err := store.StartLogging(in.Name); err != nil {
		return s.mapStoreError(err)
	}
	return nil
}

// stopLoggingCore is the single entry point for StopLogging.
func (s *CloudTrailService) stopLoggingCore(store cloudtrailstore.CloudTrailStoreInterface, in TrailNameInput) error {
	if in.Name == "" {
		return ErrInvalidTrailName
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
		opts := storecommon.ListOptions{MaxItems: cloudtrailstore.MaxListTrailsMaxItems}
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
// that explicitly-provided empty strings clear the optional fields (AWS spec
// behaviour). Every present member runs through the same validator the
// create path uses, so an update can never persist a value create would
// reject; the required-value members (S3BucketName) are validated even when
// provided empty, because the trail must always carry a usable bucket.
// resolvedSnsTopicARN carries the topic ARN the caller resolved from the
// SnsTopicName member (name or ARN form); it is applied whenever a topic
// stays configured, covering both the member-update and the heal of records
// predating SNS resolution.
func applyTrailUpdates(trail *cloudtrailstore.Trail, params map[string]interface{}, resolvedSnsTopicARN string) error {
	if v, ok := params["S3BucketName"]; ok {
		name := fmt.Sprintf("%v", v)
		if err := validateS3BucketName(name); err != nil {
			return err
		}
		trail.S3BucketName = name
	}
	if v, ok := params["S3KeyPrefix"]; ok {
		prefix := fmt.Sprintf("%v", v)
		if prefix != "" {
			if err := validateS3KeyPrefix(prefix); err != nil {
				return err
			}
		}
		trail.S3KeyPrefix = prefix
	}
	if v, ok := params["SnsTopicName"]; ok {
		name := fmt.Sprintf("%v", v)
		if name != "" {
			if err := validateSnsTopicReference(name); err != nil {
				return err
			}
		}
		trail.SnsTopicName = name
		// The ARN is resolved from the name by the destination validation
		// before the mutation and applied below.
		trail.SnsTopicARN = ""
	}
	if b := boolParam(params, "IncludeGlobalServiceEvents"); b != nil {
		trail.IncludeGlobalServiceEvents = *b
	}
	if b := boolParam(params, "IsMultiRegionTrail"); b != nil {
		trail.IsMultiRegionTrail = *b
	}
	if b := boolParam(params, "IsOrganizationTrail"); b != nil {
		trail.IsOrganizationTrail = *b
	}
	if b := boolParam(params, "EnableLogFileValidation"); b != nil {
		trail.LogFileValidationEnabled = *b
	}
	if b := boolParam(params, "RecursiveLogging"); b != nil {
		trail.RecursiveLogging = *b
	}
	if v, ok := params["CloudWatchLogsLogGroupArn"]; ok {
		arn := fmt.Sprintf("%v", v)
		if arn != "" {
			if err := validateCloudWatchLogsLogGroupARN(arn); err != nil {
				return err
			}
		}
		trail.CloudWatchLogsLogGroupARN = arn
	}
	if v, ok := params["CloudWatchLogsRoleArn"]; ok {
		arn := fmt.Sprintf("%v", v)
		if arn != "" {
			if err := validateCloudWatchLogsRoleARN(arn); err != nil {
				return err
			}
		}
		trail.CloudWatchLogsRoleARN = arn
	}
	if v, ok := params["KmsKeyId"]; ok {
		keyID := fmt.Sprintf("%v", v)
		if keyID != "" {
			if err := validateKMSKeyID(keyID); err != nil {
				return err
			}
		}
		trail.KMSKeyID = keyID
	}
	if resolvedSnsTopicARN != "" {
		trail.SnsTopicARN = resolvedSnsTopicARN
	}
	return nil
}

// parseTagsFromParams converts a raw TagsList parameter (as produced by the
// request parser) into a []tags.Tag.
func parseTagsFromParams(params map[string]interface{}) []tags.Tag {
	return tags.ParseTags(params, "TagsList")
}
