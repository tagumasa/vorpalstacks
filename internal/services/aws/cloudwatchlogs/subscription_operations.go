package cloudwatchlogs

import (
	"context"
	"fmt"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
	"vorpalstacks/internal/utils/aws/arn"
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
		DestinationArn:         destinationArn,
		RoleArn:                roleArn,
		Distribution:           distribution,
		ApplyOnTransformedLogs: request.GetBoolParam(req.Parameters, "ApplyOnTransformedLogs"),
		FieldSelectionCriteria: fieldSelectionCriteria,
		EmitSystemFields:       request.GetStringList(req.Parameters, "EmitSystemFields"),
	}); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// PutSubscriptionFilterInput holds parameters for PutSubscriptionFilter.
type PutSubscriptionFilterInput struct {
	LogGroupName           string
	FilterName             string
	FilterPattern          string
	DestinationArn         string
	RoleArn                string
	Distribution           string
	ApplyOnTransformedLogs bool
	FieldSelectionCriteria string
	EmitSystemFields       []string
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
		return ErrMissingParameter
	}
	if err := validateFilterPattern(input.FilterPattern); err != nil {
		return err
	}

	if !arn.IsLambdaARN(input.DestinationArn) && !arn.IsKinesisARN(input.DestinationArn) &&
		!isFirehoseARN(input.DestinationArn) {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("Invalid destinationArn: %s. Must be a Lambda, Kinesis, or Firehose ARN", input.DestinationArn), 400)
	}

	if input.RoleArn != "" {
		if s.bus == nil {
			return NewLogsError("InvalidParameterException",
				"RoleArn validation is not available (event bus not configured)", 400)
		}
		rr := s.bus.RoleResolver()
		if rr == nil {
			return NewLogsError("InvalidParameterException",
				"RoleArn validation is not available (role resolver not configured)", 400)
		}
		if err := rr.ValidateRole(ctx, input.RoleArn); err != nil {
			return awserrors.NewAWSError("InvalidParameterException", fmt.Sprintf("Invalid role ARN: %s", input.RoleArn), 400)
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

	if err := validateFieldSelectionCriteria(input.FieldSelectionCriteria); err != nil {
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

	if err := store.PutSubscriptionFilterWithLimitCheck(filter, 2); err != nil {
		return mapStoreError(err)
	}
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

	limit := input.Limit
	if limit <= 0 {
		limit = 50
	}

	filters, err := store.ListSubscriptionFilters(input.LogGroupName, input.FilterNamePrefix)
	if err != nil {
		return nil, "", mapStoreError(err)
	}

	result := pagination.PaginateSlice(filters, input.NextToken, int(limit), func(f *logsstore.SubscriptionFilter) string {
		return f.FilterName
	})
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
	limit, err := validateListLimit(int32(request.GetIntParam(req.Parameters, "Limit")), 50, 50)
	if err != nil {
		return nil, err
	}

	store, err := s.getLogsStoreByRegion(reqCtx.GetRegion())
	if err != nil {
		return nil, err
	}

	items, nextMarker, err := s.describeSubscriptionFiltersCore(store, &DescribeSubscriptionFiltersInput{
		LogGroupName:     request.GetParamLowerFirst(req.Parameters, "LogGroupName"),
		FilterNamePrefix: request.GetParamLowerFirst(req.Parameters, "FilterNamePrefix"),
		NextToken:        request.GetParamLowerFirst(req.Parameters, "NextToken"),
		Limit:            limit,
	})
	if err != nil {
		return nil, err
	}

	subscriptionFilters := make([]map[string]interface{}, len(items))
	for i, f := range items {
		subscriptionFilters[i] = formatSubscriptionFilter(f)
	}

	resp := map[string]interface{}{
		"subscriptionFilters": subscriptionFilters,
	}
	if nextMarker != "" {
		resp["nextToken"] = nextMarker
	}

	return resp, nil
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
