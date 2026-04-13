package cloudwatchlogs

import (
	"context"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// PutDestination creates a CloudWatch Logs destination.
func (s *LogsService) PutDestination(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamLowerFirst(req.Parameters, "DestinationName")
	roleArn := request.GetParamLowerFirst(req.Parameters, "RoleArn")
	targetArn := request.GetParamLowerFirst(req.Parameters, "TargetArn")
	accessPolicy := request.GetParamLowerFirst(req.Parameters, "AccessPolicy")

	if name == "" {
		return nil, ErrMissingParameter
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

	if name == "" {
		return nil, ErrMissingParameter
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
	limit := int32(request.GetIntParam(req.Parameters, "Limit"))
	if limit <= 0 {
		limit = 50
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

// PutDestinationPolicy sets the resource-based policy for a CloudWatch Logs destination.
func (s *LogsService) PutDestinationPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamLowerFirst(req.Parameters, "DestinationName")
	accessPolicy := request.GetParamLowerFirst(req.Parameters, "AccessPolicy")

	if name == "" {
		return nil, ErrMissingParameter
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
		"accessPolicy":    d.AccessPolicy,
		"creationTime":    d.CreationTime,
	}
	return result
}
