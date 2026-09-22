package cloudwatchlogs

import (
	"context"
	"fmt"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
	"vorpalstacks/internal/utils/aws/arn"
	"vorpalstacks/pkg/filterpattern"
)

// PutSubscriptionFilter creates or updates a subscription filter for the specified CloudWatch Logs log group.
func (s *LogsService) PutSubscriptionFilter(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	logGroupName := request.GetParamLowerFirst(req.Parameters, "LogGroupName")
	filterName := request.GetParamLowerFirst(req.Parameters, "FilterName")
	filterPattern := request.GetParamLowerFirst(req.Parameters, "FilterPattern")
	destinationArn := request.GetParamLowerFirst(req.Parameters, "DestinationArn")
	roleArn := request.GetParamLowerFirst(req.Parameters, "RoleArn")
	distribution := request.GetParamLowerFirst(req.Parameters, "Distribution")

	fieldSelectionCriteria := request.GetParamLowerFirst(req.Parameters, "FieldSelectionCriteria")

	store, err := s.getLogsStoreByRegion(reqCtx.GetRegion())
	if err != nil {
		return nil, err
	}

	if err := s.putSubscriptionFilterCore(ctx, store, &PutSubscriptionFilterInput{
		LogGroupName:           logGroupName,
		FilterName:             filterName,
		FilterPattern:          filterPattern,
		FilterPatternSet:       request.HasParamLowerFirst(req.Parameters, "FilterPattern"),
		DestinationArn:         destinationArn,
		RoleArn:                roleArn,
		Distribution:           distribution,
		ApplyOnTransformedLogs: request.GetBoolParam(req.Parameters, "ApplyOnTransformedLogs"),
		FieldSelectionCriteria: fieldSelectionCriteria,
		EmitSystemFields:       request.GetStringList(req.Parameters, "EmitSystemFields"),
		Region:                 reqCtx.GetRegion(),
	}); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// PutSubscriptionFilterInput holds parameters for PutSubscriptionFilter.
type PutSubscriptionFilterInput struct {
	LogGroupName  string
	FilterName    string
	FilterPattern string
	// FilterPatternSet is the wire-presence flag for the required
	// filterPattern member, whose shape allows the empty string: a present
	// empty pattern matches every event, an absent one rejects.
	FilterPatternSet       bool
	DestinationArn         string
	RoleArn                string
	Distribution           string
	ApplyOnTransformedLogs bool
	FieldSelectionCriteria string
	EmitSystemFields       []string
	Region                 string
}

// putSubscriptionFilterCore validates input and creates or updates a
// subscription filter.
func (s *LogsService) putSubscriptionFilterCore(ctx context.Context, store *logsstore.Store, input *PutSubscriptionFilterInput) error {
	if err := validateLogGroupName(input.LogGroupName); err != nil {
		return err
	}
	if err := validateFilterName(input.FilterName); err != nil {
		return err
	}
	if input.DestinationArn == "" {
		return errRequiredMember("destinationArn")
	}
	// filterPattern is a required member whose shape allows the empty
	// string, so presence — not value — is the requiredness test.
	if !input.FilterPatternSet {
		return errRequiredMember("filterPattern")
	}
	if err := validateFilterPattern(input.FilterPattern); err != nil {
		return err
	}

	// The accepted forms are the destinationArn member documentation's
	// list: a Kinesis stream, a Lambda function, or a CloudWatch Logs
	// destination (the logical-destination ARN the platform's own
	// PutDestination mints). Firehose is a documented AWS form but the
	// platform Firehose service does not exist yet, so its ARNs are
	// rejected here rather than accepted into a filter that would discard
	// every matched batch at delivery.
	switch {
	case arn.IsLambdaARN(input.DestinationArn), arn.IsKinesisARN(input.DestinationArn), isCloudWatchLogsDestinationARN(input.DestinationArn):
	case isFirehoseARN(input.DestinationArn):
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("Invalid destinationArn: %s. Firehose delivery streams are not supported until the platform Firehose service exists", input.DestinationArn), 400)
	default:
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("Invalid destinationArn: %s. Must be a Lambda, Kinesis, or CloudWatch Logs destination ARN", input.DestinationArn), 400)
	}

	if input.RoleArn != "" {
		if s.eventBus() == nil {
			return NewLogsError("InvalidParameterException",
				"RoleArn validation is not available (event bus not configured)", 400)
		}
		rr := s.eventBus().RoleResolver()
		if rr == nil {
			return NewLogsError("InvalidParameterException",
				"RoleArn validation is not available (role resolver not configured)", 400)
		}
		if err := rr.ValidateRole(ctx, input.RoleArn); err != nil {
			return NewLogsError("InvalidParameterException", fmt.Sprintf("Invalid role ARN: %s", input.RoleArn), 400)
		}
	}

	distribution := input.Distribution
	if distribution == "" {
		distribution = "ByLogStream"
	}
	if !validateDistribution(distribution) {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("Invalid distribution: %s. Allowed values: Random, ByLogStream", distribution), 400)
	}

	if _, err := store.GetLogGroup(input.LogGroupName); err != nil {
		return mapStoreError(err)
	}

	// The regex-pattern census spans both filter families (see
	// refuseRegexPatternQuotaExceeded); the filter being replaced keeps
	// its slot.
	if filterpattern.PatternContainsRegex(input.FilterPattern) {
		if err := refuseRegexPatternQuotaExceeded(store, input.LogGroupName, "", input.FilterName); err != nil {
			return err
		}
	}

	if err := validateFieldSelectionCriteria(input.FieldSelectionCriteria); err != nil {
		return err
	}
	if err := validateFieldSelectionCriteriaSyntax(input.FieldSelectionCriteria); err != nil {
		return err
	}
	if err := validateEmitSystemFields(input.EmitSystemFields); err != nil {
		return err
	}

	filter := &logsstore.SubscriptionFilter{
		LogGroupName:           input.LogGroupName,
		FilterName:             input.FilterName,
		FilterPattern:          input.FilterPattern,
		DestinationArn:         input.DestinationArn,
		RoleArn:                input.RoleArn,
		Distribution:           distribution,
		ApplyOnTransformedLogs: input.ApplyOnTransformedLogs,
		FieldSelectionCriteria: input.FieldSelectionCriteria,
		EmitSystemFields:       input.EmitSystemFields,
	}

	if err := store.PutSubscriptionFilterWithLimitCheck(filter, logsstore.MaxSubscriptionFiltersPerLogGroup); err != nil {
		return mapStoreError(err)
	}

	// The reachability probe a stored filter triggers: the documented
	// CONTROL_MESSAGE record the destination receives "mainly for
	// checking if the destination is reachable".
	s.emitControlMessage(input.Region, input.DestinationArn, input.LogGroupName, input.FilterName, distribution)
	return nil
}

// DeleteSubscriptionFilterInput holds parameters for DeleteSubscriptionFilter.
type DeleteSubscriptionFilterInput struct {
	LogGroupName string
	FilterName   string
}

// deleteSubscriptionFilterCore validates input and deletes a subscription filter.
func (s *LogsService) deleteSubscriptionFilterCore(store *logsstore.Store, input *DeleteSubscriptionFilterInput) error {
	if err := validateLogGroupName(input.LogGroupName); err != nil {
		return err
	}
	if err := validateFilterName(input.FilterName); err != nil {
		return err
	}
	if err := store.DeleteSubscriptionFilter(input.LogGroupName, input.FilterName); err != nil {
		return mapStoreError(err)
	}
	return nil
}

// DescribeSubscriptionFiltersInput holds parameters for DescribeSubscriptionFilters.
type DescribeSubscriptionFiltersInput struct {
	LogGroupName     string
	FilterNamePrefix string
	NextToken        string
	Limit            int32
}

// describeSubscriptionFiltersCore validates input and lists subscription filters.
func (s *LogsService) describeSubscriptionFiltersCore(store *logsstore.Store, input *DescribeSubscriptionFiltersInput) ([]*logsstore.SubscriptionFilter, string, error) {
	if err := validateLogGroupName(input.LogGroupName); err != nil {
		return nil, "", err
	}

	limit, err := validateListLimit(input.Limit, logsstore.DefaultDescribeLimit, logsstore.MaxDescribeLimit)
	if err != nil {
		return nil, "", err
	}

	// The operation declares ResourceNotFoundException: a listing against
	// a group that does not exist fails rather than serving an empty
	// page (the sibling stream listing carries the same check).
	if _, err := store.GetLogGroup(input.LogGroupName); err != nil {
		return nil, "", mapStoreError(err)
	}

	filters, err := store.ListSubscriptionFilters(input.LogGroupName, input.FilterNamePrefix)
	if err != nil {
		return nil, "", mapStoreError(err)
	}

	// The listing pages through the scoped token vocabulary: the group
	// and prefix digest into the scope, so a group-A token replayed
	// against group B rejects instead of repositioning its walk.
	scope := listingScope("subscriptionfilters", input.LogGroupName, input.FilterNamePrefix)
	result, err := paginateScopedListing(scope, input.NextToken, filters, int(limit), func(f *logsstore.SubscriptionFilter) string {
		return f.FilterName
	})
	if err != nil {
		return nil, "", err
	}
	return result.Items, result.NextMarker, nil
}

// DeleteSubscriptionFilter deletes the specified subscription filter from the CloudWatch Logs log group.
func (s *LogsService) DeleteSubscriptionFilter(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.getLogsStoreByRegion(reqCtx.GetRegion())
	if err != nil {
		return nil, err
	}

	if err := s.deleteSubscriptionFilterCore(store, &DeleteSubscriptionFilterInput{
		LogGroupName: request.GetParamLowerFirst(req.Parameters, "LogGroupName"),
		FilterName:   request.GetParamLowerFirst(req.Parameters, "FilterName"),
	}); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// DescribeSubscriptionFilters returns a list of subscription filters for the specified CloudWatch Logs log group.
func (s *LogsService) DescribeSubscriptionFilters(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return describeStoreFamilyHandler(s, reqCtx, req, "subscriptionFilters",
		func(store *logsstore.Store, nextToken string, limit int32) ([]*logsstore.SubscriptionFilter, string, error) {
			return s.describeSubscriptionFiltersCore(store, &DescribeSubscriptionFiltersInput{
				LogGroupName:     request.GetParamLowerFirst(req.Parameters, "LogGroupName"),
				FilterNamePrefix: request.GetParamLowerFirst(req.Parameters, "FilterNamePrefix"),
				NextToken:        nextToken,
				Limit:            limit,
			})
		}, formatSubscriptionFilter)
}

func formatSubscriptionFilter(f *logsstore.SubscriptionFilter) map[string]interface{} {
	result := map[string]interface{}{
		"filterName":     f.FilterName,
		"logGroupName":   f.LogGroupName,
		"filterPattern":  f.FilterPattern,
		"destinationArn": f.DestinationArn,
		"distribution":   f.Distribution,
		"creationTime":   f.CreationTime.UnixMilli(),
	}
	if f.RoleArn != "" {
		result["roleArn"] = f.RoleArn
	}
	if f.ApplyOnTransformedLogs {
		result["applyOnTransformedLogs"] = f.ApplyOnTransformedLogs
	}
	if f.FieldSelectionCriteria != "" {
		result["fieldSelectionCriteria"] = f.FieldSelectionCriteria
	}
	if len(f.EmitSystemFields) > 0 {
		result["emitSystemFields"] = f.EmitSystemFields
	}
	return result
}

func isFirehoseARN(ar string) bool {
	return arn.GetServiceFromARN(ar) == "firehose"
}
