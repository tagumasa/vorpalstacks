package neptunegraph

import (
	"context"
	"fmt"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	ngstore "vorpalstacks/internal/store/aws/neptunegraph"
)

// CreatePrivateGraphEndpoint creates a VPC-private endpoint for accessing the specified graph.
func (s *NeptuneGraphService) CreatePrivateGraphEndpoint(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	graphID := request.GetStringParam(req.Parameters, "graphIdentifier")
	if graphID == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "graphIdentifier")
	}

	_, err = store.GetGraph(graphID)
	if err != nil {
		if ngstore.IsNotFound(err) {
			return nil, newResourceNotFoundException("graph", graphID)
		}
		return nil, err
	}

	vpcID := request.GetStringParam(req.Parameters, "vpcId")

	if vpcID == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "vpcId")
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

	if len(subnetIds) == 0 {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "subnetIds must not be empty")
	}

	if s.eventBus == nil {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "EC2 service not available for VPC/Subnet validation")
	}
	ec2 := s.eventBus.EC2Invoker()
	if ec2 == nil {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "EC2 service not available for VPC/Subnet validation")
	}
	for _, subnetId := range subnetIds {
		subnetVpcId, _, err := ec2.LookupSubnet(ctx, reqCtx.GetRegion(), subnetId)
		if err != nil {
			return nil, newValidationException("ILLEGAL_ARGUMENT", fmt.Sprintf("subnet %s not found: %v", subnetId, err))
		}
		if subnetVpcId != vpcID {
			return nil, newValidationException("ILLEGAL_ARGUMENT", fmt.Sprintf("subnet %s belongs to VPC %s, not %s", subnetId, subnetVpcId, vpcID))
		}
	}

	ep := &ngstore.PrivateGraphEndpoint{
		GraphId:   graphID,
		VpcId:     vpcID,
		SubnetIds: subnetIds,
		Status:    "AVAILABLE",
		AccountID: s.accountID,
		Region:    reqCtx.GetRegion(),
	}

	if err := store.CreateEndpoint(ep); err != nil {
		if ngstore.IsAlreadyExists(err) {
			return nil, newConflictException("CONCURRENT_MODIFICATION")
		}
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

	graphID := request.GetStringParam(req.Parameters, "graphIdentifier")
	vpcID := request.GetStringParam(req.Parameters, "vpcId")

	ep, err := store.GetEndpoint(graphID, vpcID)
	if err != nil {
		if ngstore.IsNotFound(err) {
			return nil, newResourceNotFoundException("endpoint", vpcID)
		}
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

	graphID := request.GetStringParam(req.Parameters, "graphIdentifier")
	if graphID == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "graphIdentifier")
	}

	endpoints, err := store.ListEndpoints(graphID)
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0, len(endpoints))
	for _, ep := range endpoints {
		items = append(items, endpointToResponse(ep))
	}

	return map[string]interface{}{
		"privateGraphEndpoints": items,
	}, nil
}

// DeletePrivateGraphEndpoint removes a private graph endpoint identified by graph and VPC identifiers.
func (s *NeptuneGraphService) DeletePrivateGraphEndpoint(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	graphID := request.GetStringParam(req.Parameters, "graphIdentifier")
	vpcID := request.GetStringParam(req.Parameters, "vpcId")

	ep, err := store.GetEndpoint(graphID, vpcID)
	if err != nil {
		if ngstore.IsNotFound(err) {
			return nil, newResourceNotFoundException("endpoint", vpcID)
		}
		return nil, err
	}

	if err := store.DeleteEndpoint(graphID, vpcID); err != nil {
		logs.Warn("failed to delete endpoint", logs.String("graphId", graphID), logs.String("vpcId", vpcID), logs.Err(err))
	}

	return endpointToResponse(ep), nil
}
