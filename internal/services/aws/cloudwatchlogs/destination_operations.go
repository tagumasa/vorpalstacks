package cloudwatchlogs

import (
	"context"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// PutDestination creates a CloudWatch Logs destination.
func (s *LogsService) PutDestination(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamLowerFirst(req.Parameters, "DestinationName")
	roleArn := request.GetParamLowerFirst(req.Parameters, "RoleArn")
	targetArn := request.GetParamLowerFirst(req.Parameters, "TargetArn")
	accessPolicy := request.GetParamLowerFirst(req.Parameters, "AccessPolicy")
	tags := tagutil.ToMap(tagutil.ParseTagsWithQueryFallback(req.Parameters, "Tags"))

	if err := validateDestinationName(name); err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	arn := store.ARNBuilder().CloudWatch().Destination(name)

	dest := &logsstore.Destination{
		Name:         name,
		ARN:          arn,
		RoleArn:      roleArn,
		TargetArn:    targetArn,
		AccessPolicy: accessPolicy,
		Tags:         tags,
	}

	if err := store.PutDestination(dest); err != nil {
		return nil, mapStoreError(err)
	}

	return map[string]interface{}{
		"destination": formatDestination(dest),
	}, nil
}

// DeleteDestination deletes a CloudWatch Logs destination.
func (s *LogsService) DeleteDestination(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamLowerFirst(req.Parameters, "DestinationName")

	if err := validateDestinationName(name); err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.DeleteDestination(name); err != nil {
		return nil, mapStoreError(err)
	}

	return response.EmptyResponse(), nil
}

// DescribeDestinations lists CloudWatch Logs destinations, optionally filtered by prefix.
func (s *LogsService) DescribeDestinations(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	prefix := request.GetParamLowerFirst(req.Parameters, "DestinationNamePrefix")
	nextToken := request.GetParamLowerFirst(req.Parameters, "NextToken")
	limit, err := validateListLimit(int32(request.GetIntParam(req.Parameters, "Limit")), 50, 50)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	destinations, err := store.ListDestinations(prefix)
	if err != nil {
		return nil, mapStoreError(err)
	}

	result := pagination.PaginateSlice(destinations, nextToken, int(limit), func(d *logsstore.Destination) string {
		return d.Name
	})

	items := make([]map[string]interface{}, len(result.Items))
	for i, d := range result.Items {
		items[i] = formatDestination(d)
	}

	resp := map[string]interface{}{
		"destinations": items,
	}
	if result.NextMarker != "" {
		resp["nextToken"] = result.NextMarker
	}

	return resp, nil
}

// PutDestinationPolicy sets the resource-based policy for a CloudWatch Logs
// destination. AWS always overwrites the existing policy regardless of the
// forceUpdate parameter — forceUpdate is only an affirmation that subscription
// filters have been updated when migrating from individual AWS account
// permissions to an Organisation ID, which is not applicable to this
// edge/on-premises platform.
func (s *LogsService) PutDestinationPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamLowerFirst(req.Parameters, "DestinationName")
	accessPolicy := request.GetParamLowerFirst(req.Parameters, "AccessPolicy")

	if err := validateDestinationName(name); err != nil {
		return nil, err
	}
	if err := validateAccessPolicy(accessPolicy); err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.PutDestinationPolicy(name, accessPolicy); err != nil {
		return nil, mapStoreError(err)
	}

	return response.EmptyResponse(), nil
}

func formatDestination(d *logsstore.Destination) map[string]interface{} {
	result := map[string]interface{}{
		"destinationName": d.Name,
		"arn":             d.ARN,
		"roleArn":         d.RoleArn,
		"targetArn":       d.TargetArn,
		"creationTime":    d.CreationTime,
	}
	if d.AccessPolicy != "" {
		result["accessPolicy"] = d.AccessPolicy
	}
	if len(d.Tags) > 0 {
		result["tags"] = d.Tags
	}
	return result
}
