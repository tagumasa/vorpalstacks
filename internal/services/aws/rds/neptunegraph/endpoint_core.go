package neptunegraph

// Private graph endpoint Core functions: the single validation and
// persistence path for the private endpoint operations.

import (
	"context"
	"fmt"

	"vorpalstacks/internal/core/logs"
	ngstore "vorpalstacks/internal/store/aws/rds/neptunegraph"
)

// CreatePrivateGraphEndpointInput carries the wire-parsed
// CreatePrivateGraphEndpoint request.
type CreatePrivateGraphEndpointInput struct {
	GraphIdentifier string
	VpcId           string
	SubnetIds       []string
	Region          string
}

// GetPrivateGraphEndpointInput carries the wire-parsed
// GetPrivateGraphEndpoint request.
type GetPrivateGraphEndpointInput struct {
	GraphIdentifier string
	VpcId           string
}

// ListPrivateGraphEndpointsInput carries the wire-parsed
// ListPrivateGraphEndpoints request.
type ListPrivateGraphEndpointsInput struct {
	GraphIdentifier string
}

// DeletePrivateGraphEndpointInput carries the wire-parsed
// DeletePrivateGraphEndpoint request.
type DeletePrivateGraphEndpointInput struct {
	GraphIdentifier string
	VpcId           string
}

// createPrivateGraphEndpointCore validates the VPC/subnet membership through
// EC2 and persists a private endpoint for the graph.
func (s *NeptuneGraphService) createPrivateGraphEndpointCore(ctx context.Context, store *ngstore.NeptuneGraphStore, in *CreatePrivateGraphEndpointInput) (*ngstore.PrivateGraphEndpoint, error) {
	graphID := in.GraphIdentifier
	if graphID == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "graphIdentifier")
	}

	graph, err := s.resolveGraphIdentifier(store, graphID)
	if err != nil {
		if ngstore.IsNotFound(err) {
			return nil, newResourceNotFoundException("graph", graphID)
		}
		return nil, err
	}

	vpcID := in.VpcId

	if vpcID == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "vpcId")
	}

	if len(in.SubnetIds) == 0 {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "subnetIds must not be empty")
	}

	if s.eventBus == nil {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "EC2 service not available for VPC/Subnet validation")
	}
	ec2 := s.eventBus.EC2Invoker()
	if ec2 == nil {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "EC2 service not available for VPC/Subnet validation")
	}
	for _, subnetId := range in.SubnetIds {
		subnetVpcId, _, err := ec2.LookupSubnet(ctx, in.Region, subnetId)
		if err != nil {
			return nil, newValidationException("ILLEGAL_ARGUMENT", fmt.Sprintf("subnet %s not found: %v", subnetId, err))
		}
		if subnetVpcId != vpcID {
			return nil, newValidationException("ILLEGAL_ARGUMENT", fmt.Sprintf("subnet %s belongs to VPC %s, not %s", subnetId, subnetVpcId, vpcID))
		}
	}

	ep := &ngstore.PrivateGraphEndpoint{
		GraphId:       graph.Id,
		VpcId:         vpcID,
		VpcEndpointId: generateID("vpce-"),
		SubnetIds:     in.SubnetIds,
		Status:        "AVAILABLE",
		AccountID:     s.accountID,
		Region:        in.Region,
	}

	if err := store.CreateEndpoint(ep); err != nil {
		if ngstore.IsAlreadyExists(err) {
			return nil, newConflictException("CONCURRENT_MODIFICATION")
		}
		return nil, err
	}

	return ep, nil
}

// getPrivateGraphEndpointCore retrieves a private endpoint by graph and VPC
// identifiers.
func (s *NeptuneGraphService) getPrivateGraphEndpointCore(store *ngstore.NeptuneGraphStore, in *GetPrivateGraphEndpointInput) (*ngstore.PrivateGraphEndpoint, error) {
	graphID := in.GraphIdentifier
	if graphID == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "graphIdentifier")
	}
	vpcID := in.VpcId
	if vpcID == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "vpcId")
	}

	ep, err := store.GetEndpoint(graphID, vpcID)
	if err != nil {
		if ngstore.IsNotFound(err) {
			return nil, newResourceNotFoundException("endpoint", vpcID)
		}
		return nil, err
	}

	return ep, nil
}

// listPrivateGraphEndpointsCore returns all private endpoints for the graph.
func (s *NeptuneGraphService) listPrivateGraphEndpointsCore(store *ngstore.NeptuneGraphStore, in *ListPrivateGraphEndpointsInput) ([]*ngstore.PrivateGraphEndpoint, error) {
	graphID := in.GraphIdentifier
	if graphID == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "graphIdentifier")
	}

	return store.ListEndpoints(graphID)
}

// deletePrivateGraphEndpointCore removes a private endpoint identified by
// graph and VPC identifiers.
func (s *NeptuneGraphService) deletePrivateGraphEndpointCore(store *ngstore.NeptuneGraphStore, in *DeletePrivateGraphEndpointInput) (*ngstore.PrivateGraphEndpoint, error) {
	graphID := in.GraphIdentifier
	if graphID == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "graphIdentifier")
	}
	vpcID := in.VpcId
	if vpcID == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "vpcId")
	}

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

	return ep, nil
}
