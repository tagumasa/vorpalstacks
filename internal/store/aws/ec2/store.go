package ec2

import (
	"errors"
	"fmt"

	"vorpalstacks/internal/core/storage"
	storecommon "vorpalstacks/internal/store/aws/common"
)

const (
	vpcBucket           = "ec2-vpcs"
	subnetBucket        = "ec2-subnets"
	securityGroupBucket = "ec2-security-groups"
)

var (
	// ErrVPCNotFound is returned when a requested VPC does not exist.
	ErrVPCNotFound = errors.New("ec2: VPC not found")
	// ErrVPCAlreadyExists is returned when creating a VPC with a duplicate VpcId.
	ErrVPCAlreadyExists = errors.New("ec2: VPC already exists")
	// ErrSubnetNotFound is returned when a requested subnet does not exist.
	ErrSubnetNotFound = errors.New("ec2: subnet not found")
	// ErrSubnetAlreadyExists is returned when creating a subnet with a duplicate SubnetId.
	ErrSubnetAlreadyExists = errors.New("ec2: subnet already exists")
	// ErrSGNotFound is returned when a requested security group does not exist.
	ErrSGNotFound = errors.New("ec2: security group not found")
	// ErrSGAlreadyExists is returned when creating a security group with a duplicate GroupId.
	ErrSGAlreadyExists = errors.New("ec2: security group already exists")
	// ErrInvalidVPCState is returned when an operation is attempted on a VPC in an invalid state.
	ErrInvalidVPCState = errors.New("ec2: invalid VPC state")
	// ErrInvalidSubnetState is returned when an operation is attempted on a subnet in an invalid state.
	ErrInvalidSubnetState = errors.New("ec2: invalid subnet state")
)

// EC2Store provides persistent storage for EC2 VPC, Subnet, and SecurityGroup
// resources, backed by region-scoped Pebble buckets.
type EC2Store struct {
	vpcs           *storecommon.BaseStore
	subnets        *storecommon.BaseStore
	securityGroups *storecommon.BaseStore
	accountID      string
	region         string
}

// NewEC2Store creates a new EC2Store backed by the given regional storage.
func NewEC2Store(store storage.TransactionalStorageWith2PC, accountID, region string) *EC2Store {
	return &EC2Store{
		vpcs:           storecommon.NewBaseStore(store.Bucket(vpcBucket+"-"+region), "ec2-vpc"),
		subnets:        storecommon.NewBaseStore(store.Bucket(subnetBucket+"-"+region), "ec2-subnet"),
		securityGroups: storecommon.NewBaseStore(store.Bucket(securityGroupBucket+"-"+region), "ec2-sg"),
		accountID:      accountID,
		region:         region,
	}
}

// CreateVPC persists a new VPC. Returns ErrVPCAlreadyExists if the VpcId is already in use.
func (s *EC2Store) CreateVPC(vpc *VPC) error {
	if s.vpcs.Exists(vpc.VpcId) {
		return ErrVPCAlreadyExists
	}
	return s.vpcs.Put(vpc.VpcId, vpc)
}

// GetVPC retrieves a VPC by its ID. Returns ErrVPCNotFound if it does not exist.
func (s *EC2Store) GetVPC(vpcId string) (*VPC, error) {
	var vpc VPC
	if err := s.vpcs.Get(vpcId, &vpc); err != nil {
		if errors.Is(err, storecommon.ErrNotFound) {
			return nil, ErrVPCNotFound
		}
		return nil, fmt.Errorf("ec2: failed to get VPC %s: %w", vpcId, err)
	}
	return &vpc, nil
}

// DeleteVPC removes a VPC by its ID.
func (s *EC2Store) DeleteVPC(vpcId string) error {
	return s.vpcs.Delete(vpcId)
}

// ListVPCs returns all VPCs in this region.
func (s *EC2Store) ListVPCs() ([]*VPC, error) {
	result, err := storecommon.List[VPC](s.vpcs, storecommon.ListOptions{}, nil)
	if err != nil {
		return nil, err
	}
	return result.Items, nil
}

// CreateSubnet persists a new subnet. Returns ErrSubnetAlreadyExists if the SubnetId is already in use.
func (s *EC2Store) CreateSubnet(subnet *Subnet) error {
	if s.subnets.Exists(subnet.SubnetId) {
		return ErrSubnetAlreadyExists
	}
	return s.subnets.Put(subnet.SubnetId, subnet)
}

// GetSubnet retrieves a subnet by its ID. Returns ErrSubnetNotFound if it does not exist.
func (s *EC2Store) GetSubnet(subnetId string) (*Subnet, error) {
	var subnet Subnet
	if err := s.subnets.Get(subnetId, &subnet); err != nil {
		if errors.Is(err, storecommon.ErrNotFound) {
			return nil, ErrSubnetNotFound
		}
		return nil, fmt.Errorf("ec2: failed to get subnet %s: %w", subnetId, err)
	}
	return &subnet, nil
}

// DeleteSubnet removes a subnet by its ID.
func (s *EC2Store) DeleteSubnet(subnetId string) error {
	return s.subnets.Delete(subnetId)
}

// ListSubnets returns all subnets in this region.
func (s *EC2Store) ListSubnets() ([]*Subnet, error) {
	result, err := storecommon.List[Subnet](s.subnets, storecommon.ListOptions{}, nil)
	if err != nil {
		return nil, err
	}
	return result.Items, nil
}

// ListSubnetsByVPC returns all subnets belonging to the specified VPC.
func (s *EC2Store) ListSubnetsByVPC(vpcId string) ([]*Subnet, error) {
	result, err := storecommon.List(s.subnets, storecommon.ListOptions{},
		func(subnet *Subnet) bool { return subnet.VpcId == vpcId })
	if err != nil {
		return nil, err
	}
	return result.Items, nil
}

// CreateSecurityGroup persists a new security group. Returns ErrSGAlreadyExists if the GroupId is already in use.
func (s *EC2Store) CreateSecurityGroup(sg *SecurityGroup) error {
	if s.securityGroups.Exists(sg.GroupId) {
		return ErrSGAlreadyExists
	}
	return s.securityGroups.Put(sg.GroupId, sg)
}

// GetSecurityGroup retrieves a security group by its ID. Returns ErrSGNotFound if it does not exist.
func (s *EC2Store) GetSecurityGroup(groupId string) (*SecurityGroup, error) {
	var sg SecurityGroup
	if err := s.securityGroups.Get(groupId, &sg); err != nil {
		if errors.Is(err, storecommon.ErrNotFound) {
			return nil, ErrSGNotFound
		}
		return nil, fmt.Errorf("ec2: failed to get security group %s: %w", groupId, err)
	}
	return &sg, nil
}

// UpdateSecurityGroup replaces the persisted security group.
func (s *EC2Store) UpdateSecurityGroup(sg *SecurityGroup) error {
	return s.securityGroups.Put(sg.GroupId, sg)
}

// DeleteSecurityGroup removes a security group by its ID.
func (s *EC2Store) DeleteSecurityGroup(groupId string) error {
	return s.securityGroups.Delete(groupId)
}

// ListSecurityGroups returns all security groups in this region.
func (s *EC2Store) ListSecurityGroups() ([]*SecurityGroup, error) {
	result, err := storecommon.List[SecurityGroup](s.securityGroups, storecommon.ListOptions{}, nil)
	if err != nil {
		return nil, err
	}
	return result.Items, nil
}

// ListSecurityGroupsByVPC returns all security groups belonging to the specified VPC.
func (s *EC2Store) ListSecurityGroupsByVPC(vpcId string) ([]*SecurityGroup, error) {
	result, err := storecommon.List(s.securityGroups, storecommon.ListOptions{},
		func(sg *SecurityGroup) bool { return sg.VpcId == vpcId })
	if err != nil {
		return nil, err
	}
	return result.Items, nil
}
