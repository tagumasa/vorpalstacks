package ec2

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/handler"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/eventbus"
	storecommon "vorpalstacks/internal/store/aws/common"
	ec2store "vorpalstacks/internal/store/aws/ec2"
)

// EC2Service handles incoming EC2 API requests for VPC, Subnet, and SecurityGroup resources.
type EC2Service struct {
	accountID      string
	region         string
	storageManager *storage.RegionStorageManager
	stores         sync.Map
	bus            eventbus.Bus
}

// NewEC2Service creates a new EC2Service for the specified account and region.
func NewEC2Service(accountID, region string) *EC2Service {
	return &EC2Service{accountID: accountID, region: region}
}

// SetStorageManager injects the region storage manager.
func (s *EC2Service) SetStorageManager(sm *storage.RegionStorageManager) {
	s.storageManager = sm
}

// SetEventBus injects the event bus for cross-service dependency checks
// (e.g. checking whether Lambda or Neptune reference a subnet/SG before
// deletion).
func (s *EC2Service) SetEventBus(bus eventbus.Bus) {
	s.bus = bus
}

// SetEC2Store pre-populates the store cache for the given region.
func (s *EC2Service) SetEC2Store(region string, store *ec2store.EC2Store) {
	s.stores.Store(region, store)
}

// storeForRegion returns the EC2 store for the given region, creating it lazily.
func (s *EC2Service) storeForRegion(region string) (*ec2store.EC2Store, error) {
	return storecommon.GetOrCreateStoreE(&s.stores, region, func() (*ec2store.EC2Store, error) {
		if s.storageManager == nil {
			return nil, fmt.Errorf("storage manager not set")
		}
		rs, err := s.storageManager.GetStorage(region)
		if err != nil {
			return nil, err
		}
		tstore, ok := rs.(storage.TransactionalStorageWith2PC)
		if !ok {
			return nil, fmt.Errorf("storage does not support 2PC")
		}
		return ec2store.NewEC2Store(tstore, s.accountID, region), nil
	})
}

// store returns the EC2 store for the request's region.
func (s *EC2Service) store(reqCtx *request.RequestContext) (*ec2store.EC2Store, error) {
	return s.storeForRegion(reqCtx.GetRegion())
}

// GetStoreForRegion returns the cached EC2 store for the given region.
func (s *EC2Service) GetStoreForRegion(region string) (*ec2store.EC2Store, error) {
	return s.storeForRegion(region)
}

// LookupSubnet resolves a subnet ID to its VPC ID and availability zone for
// the given region. Implements eventbus.EC2SubnetLookup for cross-service
// subnet resolution (e.g. Neptune DB subnet groups).
func (s *EC2Service) LookupSubnet(ctx context.Context, region string, subnetId string) (string, string, error) {
	store, err := s.storeForRegion(region)
	if err != nil {
		return "", "", err
	}
	subnet, err := store.GetSubnet(subnetId)
	if err != nil {
		return "", "", err
	}
	return subnet.VpcId, subnet.AvailabilityZone, nil
}

// LookupSecurityGroup resolves a security group ID to its VPC ID for the
// given region. Implements eventbus.EC2SubnetLookup for cross-service
// security group validation (e.g. Neptune DB clusters, Lambda VPC config).
func (s *EC2Service) LookupSecurityGroup(ctx context.Context, region string, groupId string) (string, error) {
	store, err := s.storeForRegion(region)
	if err != nil {
		return "", err
	}
	sg, err := store.GetSecurityGroup(groupId)
	if err != nil {
		return "", err
	}
	return sg.VpcId, nil
}

// LookupVPC validates that a VPC exists in the given region. Returns nil if
// the VPC exists, or an error otherwise. Implements eventbus.EC2SubnetLookup
// for cross-service VPC validation (e.g. Route53 private hosted zone
// associations).
func (s *EC2Service) LookupVPC(ctx context.Context, region string, vpcId string) error {
	store, err := s.storeForRegion(region)
	if err != nil {
		return err
	}
	_, err = store.GetVPC(vpcId)
	return err
}

// RegisterHandlers registers all EC2 API handlers with the dispatcher.
func (s *EC2Service) RegisterHandlers(d handler.Registrar) {
	d.RegisterHandlerForService("ec2", "CreateVpc", s.CreateVpc)
	d.RegisterHandlerForService("ec2", "DescribeVpcs", s.DescribeVpcs)
	d.RegisterHandlerForService("ec2", "DeleteVpc", s.DeleteVpc)

	d.RegisterHandlerForService("ec2", "CreateSubnet", s.CreateSubnet)
	d.RegisterHandlerForService("ec2", "DescribeSubnets", s.DescribeSubnets)
	d.RegisterHandlerForService("ec2", "DeleteSubnet", s.DeleteSubnet)

	d.RegisterHandlerForService("ec2", "CreateSecurityGroup", s.CreateSecurityGroup)
	d.RegisterHandlerForService("ec2", "DescribeSecurityGroups", s.DescribeSecurityGroups)
	d.RegisterHandlerForService("ec2", "DeleteSecurityGroup", s.DeleteSecurityGroup)
	d.RegisterHandlerForService("ec2", "AuthorizeSecurityGroupIngress", s.AuthorizeSecurityGroupIngress)
	d.RegisterHandlerForService("ec2", "AuthorizeSecurityGroupEgress", s.AuthorizeSecurityGroupEgress)
	d.RegisterHandlerForService("ec2", "RevokeSecurityGroupIngress", s.RevokeSecurityGroupIngress)
	d.RegisterHandlerForService("ec2", "RevokeSecurityGroupEgress", s.RevokeSecurityGroupEgress)
}

// storeErrorMappings maps EC2 store sentinel errors to AWS-compatible API errors.
var storeErrorMappings = []awserrors.StoreErrorMapping{
	{Store: ec2store.ErrVPCNotFound, AWS: awserrors.NewAWSError("InvalidVpcID.NotFound", "The VPC ID does not exist", http.StatusNotFound)},
	{Store: ec2store.ErrVPCAlreadyExists, AWS: awserrors.NewAWSError("InvalidVpcID.Duplicate", "The VPC already exists", http.StatusConflict)},
	{Store: ec2store.ErrSubnetNotFound, AWS: awserrors.NewAWSError("InvalidSubnetID.NotFound", "The subnet ID does not exist", http.StatusNotFound)},
	{Store: ec2store.ErrSubnetAlreadyExists, AWS: awserrors.NewAWSError("InvalidSubnetID.Duplicate", "The subnet already exists", http.StatusConflict)},
	{Store: ec2store.ErrSGNotFound, AWS: awserrors.NewAWSError("InvalidGroup.NotFound", "The security group does not exist", http.StatusNotFound)},
	{Store: ec2store.ErrSGAlreadyExists, AWS: awserrors.NewAWSError("InvalidGroup.Duplicate", "The security group already exists", http.StatusConflict)},
}

// translateStoreError maps EC2 store sentinel errors to AWS-compatible API errors.
func translateStoreError(err error) error {
	return awserrors.MapStoreError(err, storeErrorMappings)
}
