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

	dest, err := s.putDestinationCore(name, roleArn, targetArn, accessPolicy, tags, reqCtx.GetRegion())
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"destination": formatDestination(dest),
	}, nil
}

func (s *LogsService) putDestinationCore(name, roleArn, targetArn, accessPolicy string, tags map[string]string, region string) (*logsstore.Destination, error) {
	if err := validateDestinationName(name); err != nil {
		return nil, err
	}
	if err := validateKinesisOrFirehoseArn(targetArn); err != nil {
		return nil, err
	}
	if err := validateIAMRoleArn(roleArn); err != nil {
		return nil, err
	}

	store, err := s.getLogsStoreByRegion(region)
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
	return dest, nil
}

// DeleteDestination deletes a CloudWatch Logs destination.
func (s *LogsService) DeleteDestination(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamLowerFirst(req.Parameters, "DestinationName")

	if err := s.deleteDestinationCore(name, reqCtx.GetRegion()); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

func (s *LogsService) deleteDestinationCore(name, region string) error {
	if err := validateDestinationName(name); err != nil {
		return err
	}

	store, err := s.getLogsStoreByRegion(region)
	if err != nil {
		return err
	}

	if err := store.DeleteDestination(name); err != nil {
		return mapStoreError(err)
	}
	return nil
}

// DescribeDestinations lists CloudWatch Logs destinations, optionally filtered by prefix.
func (s *LogsService) DescribeDestinations(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	prefix := request.GetParamLowerFirst(req.Parameters, "DestinationNamePrefix")
	nextToken := request.GetParamLowerFirst(req.Parameters, "NextToken")
	limit := int32(request.GetIntParam(req.Parameters, "Limit"))

	destinations, nextMarker, err := s.describeDestinationsCore(prefix, nextToken, reqCtx.GetRegion(), limit)
	if err != nil {
		return nil, err
	}

	items := make([]map[string]interface{}, len(destinations))
	for i, d := range destinations {
		items[i] = formatDestination(d)
	}

	resp := map[string]interface{}{
		"destinations": items,
	}
	if nextMarker != "" {
		resp["nextToken"] = nextMarker
	}

	return resp, nil
}

func (s *LogsService) describeDestinationsCore(prefix, nextToken, region string, limit int32) ([]*logsstore.Destination, string, error) {
	l, err := validateListLimit(limit, 50, 50)
	if err != nil {
		return nil, "", err
	}

	store, err := s.getLogsStoreByRegion(region)
	if err != nil {
		return nil, "", err
	}

	destinations, err := store.ListDestinations(prefix)
	if err != nil {
		return nil, "", mapStoreError(err)
	}

	result := pagination.PaginateSlice(destinations, nextToken, int(l), func(d *logsstore.Destination) string {
		return d.Name
	})

	return result.Items, result.NextMarker, nil
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

	if err := s.putDestinationPolicyCore(name, accessPolicy, reqCtx.GetRegion()); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

func (s *LogsService) putDestinationPolicyCore(name, accessPolicy, region string) error {
	if err := validateDestinationName(name); err != nil {
		return err
	}
	if err := validateAccessPolicyJSON(accessPolicy); err != nil {
		return err
	}

	store, err := s.getLogsStoreByRegion(region)
	if err != nil {
		return err
	}

	if err := store.PutDestinationPolicy(name, accessPolicy); err != nil {
		return mapStoreError(err)
	}
	return nil
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
