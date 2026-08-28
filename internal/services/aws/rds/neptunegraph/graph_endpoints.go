package neptunegraph

import (
	"context"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	ngstore "vorpalstacks/internal/store/aws/rds/neptunegraph"
)

// CreatePrivateGraphEndpoint creates a VPC-private endpoint for accessing the specified graph.
func (s *NeptuneGraphService) CreatePrivateGraphEndpoint(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	var subnetIds []string
	if v, ok := req.Parameters["subnetIds"]; ok {
		if arr, ok := v.([]interface{}); ok {
			for _, item := range arr {
				if str, ok := item.(string); ok {
					subnetIds = append(subnetIds, str)
				}
			}
		}
	}

	in := &CreatePrivateGraphEndpointInput{
		GraphIdentifier: request.GetStringParam(req.Parameters, "graphIdentifier"),
		VpcId:           request.GetStringParam(req.Parameters, "vpcId"),
		SubnetIds:       subnetIds,
		Region:          reqCtx.GetRegion(),
	}

	ep, err := s.createPrivateGraphEndpointCore(ctx, store, in)
	if err != nil {
		return nil, err
	}
	return endpointToResponse(ep), nil
}

// GetPrivateGraphEndpoint retrieves a private graph endpoint by graph and VPC identifiers.
func (s *NeptuneGraphService) GetPrivateGraphEndpoint(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	in := &GetPrivateGraphEndpointInput{
		GraphIdentifier: request.GetStringParam(req.Parameters, "graphIdentifier"),
		VpcId:           request.GetStringParam(req.Parameters, "vpcId"),
	}

	ep, err := s.getPrivateGraphEndpointCore(store, in)
	if err != nil {
		return nil, err
	}
	return endpointToResponse(ep), nil
}

// ListPrivateGraphEndpoints returns all private endpoints for the specified graph.
func (s *NeptuneGraphService) ListPrivateGraphEndpoints(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	in := &ListPrivateGraphEndpointsInput{
		GraphIdentifier: request.GetStringParam(req.Parameters, "graphIdentifier"),
	}

	endpoints, err := s.listPrivateGraphEndpointsCore(store, in)
	if err != nil {
		return nil, err
	}

	maxResults := clampMaxResults(request.GetIntParam(req.Parameters, "maxResults"))
	nextToken := request.GetStringParam(req.Parameters, "nextToken")

	result := pagination.PaginateSlice(endpoints, nextToken, maxResults, func(ep *ngstore.PrivateGraphEndpoint) string {
		return ep.GraphId + ":" + ep.VpcId
	})

	items := make([]interface{}, 0, len(result.Items))
	for _, ep := range result.Items {
		items = append(items, endpointToResponse(ep))
	}

	resp := map[string]interface{}{
		"privateGraphEndpoints": items,
	}
	if result.IsTruncated {
		resp["nextToken"] = result.NextMarker
	}
	return resp, nil
}

// DeletePrivateGraphEndpoint removes a private graph endpoint identified by graph and VPC identifiers.
func (s *NeptuneGraphService) DeletePrivateGraphEndpoint(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	in := &DeletePrivateGraphEndpointInput{
		GraphIdentifier: request.GetStringParam(req.Parameters, "graphIdentifier"),
		VpcId:           request.GetStringParam(req.Parameters, "vpcId"),
	}

	ep, err := s.deletePrivateGraphEndpointCore(store, in)
	if err != nil {
		return nil, err
	}
	return endpointToResponse(ep), nil
}
