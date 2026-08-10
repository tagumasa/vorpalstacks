package ec2

import (
	"context"
	"fmt"
	"net/http"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/pagination"
	ec2store "vorpalstacks/internal/store/aws/ec2"
)

// CreateSubnetInput is the transport-agnostic input for CreateSubnet.
type CreateSubnetInput struct {
	VpcId               string
	CidrBlock           string
	AvailabilityZone    string
	MapPublicIpOnLaunch bool
	Tags                []ec2Tag
	Region              string
}

// SubnetResult is the transport-agnostic result for subnet operations.
type SubnetResult struct {
	Subnet *ec2store.Subnet
}

// SubnetListResult is the transport-agnostic result for DescribeSubnets.
type SubnetListResult struct {
	Subnets     []*ec2store.Subnet
	NextToken   string
	IsTruncated bool
}

// createSubnetCore contains the business logic for CreateSubnet.
func (s *EC2Service) createSubnetCore(store *ec2store.EC2Store, input CreateSubnetInput) (*SubnetResult, error) {
	if input.VpcId == "" {
		return nil, awserrors.NewMissingParameter("VpcId is required")
	}
	if input.CidrBlock == "" {
		return nil, awserrors.NewMissingParameter("CidrBlock is required")
	}
	canonCIDR, err := canonicalizeCIDR(input.CidrBlock)
	if err != nil {
		return nil, err
	}

	vpc, err := store.GetVPC(input.VpcId)
	if err != nil {
		return nil, translateStoreError(err)
	}
	if err := validateSubnetInVPC(canonCIDR, vpc.CidrBlock); err != nil {
		return nil, err
	}

	existingSubnets, err := store.ListSubnetsByVPC(input.VpcId)
	if err != nil {
		return nil, translateStoreError(err)
	}
	if err := validateSubnetCIDROverlap(canonCIDR, existingSubnets); err != nil {
		return nil, err
	}

	subnetID, err := GenerateSubnetID()
	if err != nil {
		return nil, err
	}

	az := input.AvailabilityZone
	if az == "" {
		az = input.Region + "a"
	}

	subnet := &ec2store.Subnet{
		SubnetId:                subnetID,
		VpcId:                   input.VpcId,
		CidrBlock:               canonCIDR,
		AvailabilityZone:        az,
		AvailableIpAddressCount: calculateAvailableIPs(canonCIDR),
		State:                   "available",
		OwnerId:                 s.accountID,
		SubnetArn:               fmt.Sprintf("arn:aws:ec2:%s:%s:subnet/%s", input.Region, s.accountID, subnetID),
		MapPublicIpOnLaunch:     input.MapPublicIpOnLaunch,
		SubnetType:              "ipv4",
		Tags:                    toStoreTags(input.Tags),
	}

	if err := store.CreateSubnet(subnet); err != nil {
		return nil, translateStoreError(err)
	}

	return &SubnetResult{Subnet: subnet}, nil
}

// describeSubnetsCore contains the business logic for DescribeSubnets.
func (s *EC2Service) describeSubnetsCore(store *ec2store.EC2Store, subnetIDs []string, filters []ec2Filter, nextToken string, maxResults int) (*SubnetListResult, error) {
	if len(subnetIDs) > 0 {
		subnets := make([]*ec2store.Subnet, 0, len(subnetIDs))
		for _, id := range subnetIDs {
			subnet, err := store.GetSubnet(id)
			if err != nil {
				return nil, translateStoreError(err)
			}
			subnets = append(subnets, subnet)
		}
		return &SubnetListResult{Subnets: subnets}, nil
	}

	allSubnets, err := store.ListSubnets()
	if err != nil {
		return nil, translateStoreError(err)
	}

	if len(filters) > 0 {
		filtered := make([]*ec2store.Subnet, 0, len(allSubnets))
		for _, sn := range allSubnets {
			if matchesSubnetFilters(sn, filters) {
				filtered = append(filtered, sn)
			}
		}
		allSubnets = filtered
	}

	if maxResults <= 0 {
		maxResults = 100
	}
	page := pagination.PaginateSlice(allSubnets, nextToken, maxResults, func(sn *ec2store.Subnet) string {
		return sn.SubnetId
	})
	return &SubnetListResult{
		Subnets:     page.Items,
		NextToken:   page.NextMarker,
		IsTruncated: page.IsTruncated,
	}, nil
}

// deleteSubnetCore contains the business logic for DeleteSubnet.
func (s *EC2Service) deleteSubnetCore(ctx context.Context, store *ec2store.EC2Store, region, subnetID string) error {
	if subnetID == "" {
		return awserrors.NewMissingParameter("SubnetId is required")
	}
	if _, err := store.GetSubnet(subnetID); err != nil {
		return translateStoreError(err)
	}

	if s.bus != nil {
		for _, checker := range s.bus.SubnetUsageCheckers() {
			if checker.IsSubnetInUse(ctx, region, subnetID) {
				return awserrors.NewAWSError(
					"DependencyViolation",
					fmt.Sprintf("The subnet '%s' has dependencies and cannot be deleted", subnetID),
					http.StatusBadRequest,
				)
			}
		}
	}

	return translateStoreError(store.DeleteSubnet(subnetID))
}
