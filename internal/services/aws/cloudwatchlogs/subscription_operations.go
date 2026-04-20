package cloudwatchlogs

import (
	"context"
	"fmt"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
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

	if roleArn != "" {
		if s.bus != nil {
			if rr := s.bus.RoleResolver(); rr != nil {
				if err := rr.ValidateRole(ctx, roleArn); err != nil {
					return nil, awserrors.NewAWSError("InvalidParameterException", fmt.Sprintf("Invalid role ARN: %s", roleArn), 400)
				}
			}
		}
	}

	if distribution == "" {
		distribution = "ByLogStream"
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	_, err = store.GetLogGroup(logGroupName)
	if err != nil {
		return nil, mapStoreError(err)
	}

	filters, err := store.ListSubscriptionFilters(logGroupName, "")
	if err != nil {
		return nil, mapStoreError(err)
	}

	existingCount := 0
	for _, f := range filters {
		if f.FilterName != filterName {
			existingCount++
		}
	}
	if existingCount >= 2 {
		return nil, ErrLimitExceeded
	}

	filter := &logsstore.SubscriptionFilter{
		LogGroupName:   logGroupName,
		FilterName:     filterName,
		FilterPattern:  filterPattern,
		DestinationArn: destinationArn,
		RoleArn:        roleArn,
		Distribution:   distribution,
	}

	if err := store.PutSubscriptionFilter(filter); err != nil {
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
		"roleArn":        f.RoleArn,
		"distribution":   f.Distribution,
		"creationTime":   f.CreationTime.UnixMilli(),
	}
	return result
}
