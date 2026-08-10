package ec2

import (
	"context"
	"net/http"

	"connectrpc.com/connect"

	svccommon "vorpalstacks/internal/common"
	svcerrors "vorpalstacks/internal/common/errors"

	pbcommon "vorpalstacks/internal/pb/aws/common"
	pb "vorpalstacks/internal/pb/aws/ec2"
	ec2connect "vorpalstacks/internal/pb/aws/ec2/ec2connect"
)

// AdminHandler implements the EC2 admin console gRPC-Web handler.
// It delegates to shared Core functions so that the same validation logic
// and store access are used by both the HTTP API handlers and the admin
// console gRPC-Web handlers.
type AdminHandler struct {
	ec2connect.UnimplementedEC2ServiceHandler
	service *EC2Service
}

var _ ec2connect.EC2ServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new EC2 admin handler.
func NewAdminHandler(svc *EC2Service) *AdminHandler {
	return &AdminHandler{service: svc}
}

// NewConnectHandler creates a gRPC-Web connect handler for the EC2 admin console.
func NewConnectHandler(svc *EC2Service) (string, http.Handler) {
	return ec2connect.NewEC2ServiceHandler(NewAdminHandler(svc))
}

func (h *AdminHandler) getStoreFromHeaders(headers http.Header) (*ec2Stores, error) {
	region := svccommon.GetRegionFromHeader(headers)
	store, err := h.service.GetStoreForRegion(region)
	if err != nil {
		return nil, err
	}
	return &ec2Stores{store: store}, nil
}

// CreateVpc creates a VPC via the admin console.
func (h *AdminHandler) CreateVpc(ctx context.Context, req *connect.Request[pb.CreateVpcRequest]) (*connect.Response[pb.CreateVpcResult], error) {
	stores, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	input := CreateVpcInput{
		CidrBlock: req.Msg.GetCidrblock(),
	}
	if t := req.Msg.GetInstancetenancy(); t > 0 {
		input.InstanceTenancy = t.String()
	}

	result, err := h.service.createVpcCore(stores.store, input)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.CreateVpcResult{
		Vpc: toPbVpc(result.Vpc),
	}), nil
}

// DescribeVpcs lists VPCs via the admin console.
func (h *AdminHandler) DescribeVpcs(ctx context.Context, req *connect.Request[pb.DescribeVpcsRequest]) (*connect.Response[pb.DescribeVpcsResult], error) {
	stores, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	var filters []ec2Filter
	for _, f := range req.Msg.GetFilters() {
		filters = append(filters, ec2Filter{
			Name:   f.GetName(),
			Values: f.GetValues(),
		})
	}

	result, err := h.service.describeVpcsCore(stores.store, req.Msg.GetVpcids(), filters, req.Msg.GetNexttoken(), int(req.Msg.GetMaxresults()))
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	pbVpcs := make([]*pb.Vpc, 0, len(result.Vpcs))
	for _, v := range result.Vpcs {
		pbVpcs = append(pbVpcs, toPbVpc(v))
	}

	resp := &pb.DescribeVpcsResult{
		Vpcs: pbVpcs,
	}
	if result.NextToken != "" {
		resp.Nexttoken = result.NextToken
	}
	return connect.NewResponse(resp), nil
}

// DeleteVpc deletes a VPC via the admin console.
func (h *AdminHandler) DeleteVpc(ctx context.Context, req *connect.Request[pb.DeleteVpcRequest]) (*connect.Response[pbcommon.Empty], error) {
	stores, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	if err := h.service.deleteVpcCore(stores.store, req.Msg.GetVpcid()); err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(&pbcommon.Empty{}), nil
}

// CreateSubnet creates a subnet via the admin console.
func (h *AdminHandler) CreateSubnet(ctx context.Context, req *connect.Request[pb.CreateSubnetRequest]) (*connect.Response[pb.CreateSubnetResult], error) {
	stores, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	region := svccommon.GetRegionFromHeader(req.Header())
	result, err := h.service.createSubnetCore(stores.store, CreateSubnetInput{
		VpcId:            req.Msg.GetVpcid(),
		CidrBlock:        req.Msg.GetCidrblock(),
		AvailabilityZone: req.Msg.GetAvailabilityzone(),
		Region:           region,
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.CreateSubnetResult{
		Subnet: toPbSubnet(result.Subnet),
	}), nil
}

// DescribeSubnets lists subnets via the admin console.
func (h *AdminHandler) DescribeSubnets(ctx context.Context, req *connect.Request[pb.DescribeSubnetsRequest]) (*connect.Response[pb.DescribeSubnetsResult], error) {
	stores, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	var filters []ec2Filter
	for _, f := range req.Msg.GetFilters() {
		filters = append(filters, ec2Filter{
			Name:   f.GetName(),
			Values: f.GetValues(),
		})
	}

	result, err := h.service.describeSubnetsCore(stores.store, req.Msg.GetSubnetids(), filters, req.Msg.GetNexttoken(), int(req.Msg.GetMaxresults()))
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	pbSubnets := make([]*pb.Subnet, 0, len(result.Subnets))
	for _, sn := range result.Subnets {
		pbSubnets = append(pbSubnets, toPbSubnet(sn))
	}

	resp := &pb.DescribeSubnetsResult{
		Subnets: pbSubnets,
	}
	if result.NextToken != "" {
		resp.Nexttoken = result.NextToken
	}
	return connect.NewResponse(resp), nil
}

// DeleteSubnet deletes a subnet via the admin console.
func (h *AdminHandler) DeleteSubnet(ctx context.Context, req *connect.Request[pb.DeleteSubnetRequest]) (*connect.Response[pbcommon.Empty], error) {
	stores, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	region := svccommon.GetRegionFromHeader(req.Header())
	if err := h.service.deleteSubnetCore(ctx, stores.store, region, req.Msg.GetSubnetid()); err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(&pbcommon.Empty{}), nil
}

// CreateSecurityGroup creates a security group via the admin console.
func (h *AdminHandler) CreateSecurityGroup(ctx context.Context, req *connect.Request[pb.CreateSecurityGroupRequest]) (*connect.Response[pb.CreateSecurityGroupResult], error) {
	stores, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	region := svccommon.GetRegionFromHeader(req.Header())
	result, err := h.service.createSecurityGroupCore(stores.store, CreateSecurityGroupInput{
		GroupName:   req.Msg.GetGroupname(),
		Description: req.Msg.GetDescription(),
		VpcId:       req.Msg.GetVpcid(),
	}, region)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.CreateSecurityGroupResult{
		Groupid:          result.GroupId,
		Securitygrouparn: result.Arn,
	}), nil
}

// DescribeSecurityGroups lists security groups via the admin console.
func (h *AdminHandler) DescribeSecurityGroups(ctx context.Context, req *connect.Request[pb.DescribeSecurityGroupsRequest]) (*connect.Response[pb.DescribeSecurityGroupsResult], error) {
	stores, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	var filters []ec2Filter
	for _, f := range req.Msg.GetFilters() {
		filters = append(filters, ec2Filter{
			Name:   f.GetName(),
			Values: f.GetValues(),
		})
	}

	result, err := h.service.describeSecurityGroupsCore(stores.store, req.Msg.GetGroupids(), req.Msg.GetGroupnames(), filters, req.Msg.GetNexttoken(), int(req.Msg.GetMaxresults()))
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	pbSGs := make([]*pb.SecurityGroup, 0, len(result.SecurityGroups))
	for _, sg := range result.SecurityGroups {
		pbSGs = append(pbSGs, toPbSecurityGroup(sg))
	}

	resp := &pb.DescribeSecurityGroupsResult{
		Securitygroups: pbSGs,
	}
	if result.NextToken != "" {
		resp.Nexttoken = result.NextToken
	}
	return connect.NewResponse(resp), nil
}

// DeleteSecurityGroup deletes a security group via the admin console.
func (h *AdminHandler) DeleteSecurityGroup(ctx context.Context, req *connect.Request[pb.DeleteSecurityGroupRequest]) (*connect.Response[pb.DeleteSecurityGroupResult], error) {
	stores, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	region := svccommon.GetRegionFromHeader(req.Header())
	if err := h.service.deleteSecurityGroupCore(ctx, stores.store, region, req.Msg.GetGroupid()); err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(&pb.DeleteSecurityGroupResult{}), nil
}
