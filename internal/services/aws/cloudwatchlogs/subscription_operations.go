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

	if logGroupName == "" || filterName == "" || destinationArn == "" {
		return nil, ErrMissingParameter
	}

	if !arn.IsLambdaARN(destinationArn) && !arn.IsKinesisARN(destinationArn) &&
		!isFirehoseARN(destinationArn) {
		return nil, NewLogsError("InvalidParameterException",
			fmt.Sprintf("Invalid destinationArn: %s. Must be a Lambda, Kinesis, or Firehose ARN", destinationArn), 400)
	}

	if roleArn != "" {
		if s.bus == nil {
			return nil, NewLogsError("InvalidParameterException",
				"RoleArn validation is not available (event bus not configured)", 400)
		}
		rr := s.bus.RoleResolver()
		if rr == nil {
			return nil, NewLogsError("InvalidParameterException",
				"RoleArn validation is not available (role resolver not configured)", 400)
		}
		if err := rr.ValidateRole(ctx, roleArn); err != nil {
			return nil, awserrors.NewAWSError("InvalidParameterException", fmt.Sprintf("Invalid role ARN: %s", roleArn), 400)
		}
	}

	if distribution == "" {
		distribution = "ByLogStream"
	}

	if distribution != "Random" && distribution != "ByLogStream" {
		return nil, NewLogsError("InvalidParameterException",
			fmt.Sprintf("Invalid distribution: %s. Allowed values: Random, ByLogStream", distribution), 400)
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	_, err = store.GetLogGroup(logGroupName)
	if err != nil {
		return nil, mapStoreError(err)
	}

	filter := &logsstore.SubscriptionFilter{
		LogGroupName:   logGroupName,
		FilterName:     filterName,
		FilterPattern:  filterPattern,
		DestinationArn: destinationArn,
		RoleArn:        roleArn,
		Distribution:   distribution,
	}

	if err := store.PutSubscriptionFilterWithLimitCheck(filter, 2); err != nil {
		return nil, mapStoreError(err)
	}

	return response.EmptyResponse(), nil
}

// DeleteSubscriptionFilter deletes the specified subscription filter from the CloudWatch Logs log group.
func (s *LogsService) DeleteSubscriptionFilter(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	logGroupName := request.GetParamLowerFirst(req.Parameters, "LogGroupName")
	filterName := request.GetParamLowerFirst(req.Parameters, "FilterName")

	if logGroupName == "" || filterName == "" {
		return nil, ErrMissingParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.DeleteSubscriptionFilter(logGroupName, filterName); err != nil {
		return nil, mapStoreError(err)
	}

	return response.EmptyResponse(), nil
}

// DescribeSubscriptionFilters returns a list of subscription filters for the specified CloudWatch Logs log group.
func (s *LogsService) DescribeSubscriptionFilters(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	logGroupName := request.GetParamLowerFirst(req.Parameters, "LogGroupName")
	filterNamePrefix := request.GetParamLowerFirst(req.Parameters, "FilterNamePrefix")
	nextToken := request.GetParamLowerFirst(req.Parameters, "NextToken")
	limit := int32(request.GetIntParam(req.Parameters, "Limit"))
	if limit <= 0 {
		limit = 50
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	filters, err := store.ListSubscriptionFilters(logGroupName, filterNamePrefix)
	if err != nil {
		return nil, mapStoreError(err)
	}

	result := pagination.PaginateSlice(filters, nextToken, int(limit), func(f *logsstore.SubscriptionFilter) string {
		return f.FilterName
	})

	subscriptionFilters := make([]map[string]interface{}, len(result.Items))
	for i, f := range result.Items {
		subscriptionFilters[i] = formatSubscriptionFilter(f)
	}

	resp := map[string]interface{}{
		"subscriptionFilters": subscriptionFilters,
	}
	if result.NextMarker != "" {
		resp["nextToken"] = result.NextMarker
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
	return result
}

func isFirehoseARN(ar string) bool {
	return arn.GetServiceFromARN(ar) == "firehose"
}
