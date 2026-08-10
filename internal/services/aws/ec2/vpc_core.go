package ec2

import (
	"net/http"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/pagination"
	ec2store "vorpalstacks/internal/store/aws/ec2"
	"vorpalstacks/internal/utils/aws/generators"
	"vorpalstacks/internal/utils/aws/types"
)

const (
	cidrAssociationPrefix = "vpc-cidr-assoc-"
)

// CreateVpcInput is the transport-agnostic input for CreateVpc.
type CreateVpcInput struct {
	CidrBlock          string
	InstanceTenancy    string
	EnableDnsSupport   *bool
	EnableDnsHostnames *bool
	Tags               []ec2Tag
}

// VpcResult is the transport-agnostic result for VPC operations.
type VpcResult struct {
	Vpc *ec2store.VPC
}

// VpcListResult is the transport-agnostic result for DescribeVpcs.
type VpcListResult struct {
	Vpcs        []*ec2store.VPC
	NextToken   string
	IsTruncated bool
}

// ec2Tag is the transport-agnostic tag type shared by all core inputs.
type ec2Tag struct {
	Key   string
	Value string
}

// generateCidrAssociationID generates a unique CIDR association ID.
func generateCidrAssociationID() (string, error) {
	return generators.GenerateIDWithPrefix(cidrAssociationPrefix, 17)
}

// createVpcCore contains the business logic for creating a VPC. Both the HTTP
// handler and the admin handler delegate to this function so that validation,
// CIDR canonicalisation, and CIDR association generation are shared.
func (s *EC2Service) createVpcCore(store *ec2store.EC2Store, input CreateVpcInput) (*VpcResult, error) {
	if input.CidrBlock == "" {
		return nil, awserrors.NewAWSError("MissingParameter",
			"CidrBlock is required",
			http.StatusBadRequest)
	}
	canonCIDR, err := canonicalizeCIDR(input.CidrBlock)
	if err != nil {
		return nil, err
	}

	vpcID, err := GenerateVpcID()
	if err != nil {
		return nil, err
	}

	instanceTenancy := input.InstanceTenancy
	if instanceTenancy == "" {
		instanceTenancy = "default"
	}

	enableDnsSupport := true
	if input.EnableDnsSupport != nil {
		enableDnsSupport = *input.EnableDnsSupport
	}
	enableDnsHostnames := false
	if input.EnableDnsHostnames != nil {
		enableDnsHostnames = *input.EnableDnsHostnames
	}

	assocID, err := generateCidrAssociationID()
	if err != nil {
		return nil, err
	}

	vpc := &ec2store.VPC{
		VpcId:              vpcID,
		CidrBlock:          canonCIDR,
		State:              "available",
		OwnerId:            s.accountID,
		InstanceTenancy:    instanceTenancy,
		DhcpOptionsId:      "dopt-default",
		IsDefault:          false,
		EnableDnsSupport:   enableDnsSupport,
		EnableDnsHostnames: enableDnsHostnames,
		CidrBlockAssociationSet: []ec2store.VpcCidrBlockAssociation{
			{
				AssociationId:  assocID,
				CidrBlock:      canonCIDR,
				CidrBlockState: ec2store.VpcCidrBlockAssociationState{State: "associated"},
			},
		},
		Tags: toStoreTags(input.Tags),
	}

	if err := store.CreateVPC(vpc); err != nil {
		return nil, translateStoreError(err)
	}

	return &VpcResult{Vpc: vpc}, nil
}

// describeVpcsCore contains the business logic for DescribeVpcs. It supports
// filtering, VpcId lookup, and pagination.
func (s *EC2Service) describeVpcsCore(store *ec2store.EC2Store, vpcIDs []string, filters []ec2Filter, nextToken string, maxResults int) (*VpcListResult, error) {
	if len(vpcIDs) > 0 {
		vpcs := make([]*ec2store.VPC, 0, len(vpcIDs))
		for _, id := range vpcIDs {
			vpc, err := store.GetVPC(id)
			if err != nil {
				return nil, translateStoreError(err)
			}
			vpcs = append(vpcs, vpc)
		}
		return &VpcListResult{Vpcs: vpcs}, nil
	}

	allVpcs, err := store.ListVPCs()
	if err != nil {
		return nil, translateStoreError(err)
	}

	if len(filters) > 0 {
		filtered := make([]*ec2store.VPC, 0, len(allVpcs))
		for _, v := range allVpcs {
			if matchesVPCFilters(v, filters) {
				filtered = append(filtered, v)
			}
		}
		allVpcs = filtered
	}

	if maxResults <= 0 {
		maxResults = 100
	}
	page := pagination.PaginateSlice(allVpcs, nextToken, maxResults, func(v *ec2store.VPC) string {
		return v.VpcId
	})
	return &VpcListResult{
		Vpcs:        page.Items,
		NextToken:   page.NextMarker,
		IsTruncated: page.IsTruncated,
	}, nil
}

// deleteVpcCore contains the business logic for DeleteVpc.
func (s *EC2Service) deleteVpcCore(store *ec2store.EC2Store, vpcID string) error {
	if vpcID == "" {
		return awserrors.NewMissingParameter("VpcId is required")
	}

	subnets, err := store.ListSubnetsByVPC(vpcID)
	if err != nil {
		return err
	}
	if len(subnets) > 0 {
		return awserrors.NewAWSError("DependencyViolation",
			"The vpc '"+vpcID+"' has dependencies and cannot be deleted",
			http.StatusBadRequest)
	}

	sgs, err := store.ListSecurityGroupsByVPC(vpcID)
	if err != nil {
		return err
	}
	if len(sgs) > 0 {
		return awserrors.NewAWSError("DependencyViolation",
			"The vpc '"+vpcID+"' has dependencies and cannot be deleted",
			http.StatusBadRequest)
	}

	return translateStoreError(store.DeleteVPC(vpcID))
}

// toStoreTags converts core ec2Tag slice to store types.Tag slice.
func toStoreTags(tags []ec2Tag) []types.Tag {
	if len(tags) == 0 {
		return nil
	}
	result := make([]types.Tag, len(tags))
	for i, t := range tags {
		result[i] = types.Tag{Key: t.Key, Value: t.Value}
	}
	return result
}
