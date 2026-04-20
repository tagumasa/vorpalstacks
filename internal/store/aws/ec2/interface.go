package ec2

// EC2StoreInterface defines the operations provided by the EC2 store.
type EC2StoreInterface interface {
	// VPC operations
	CreateVPC(vpc *VPC) error
	GetVPC(vpcId string) (*VPC, error)
	DeleteVPC(vpcId string) error
	ListVPCs() ([]*VPC, error)

	// Subnet operations
	CreateSubnet(subnet *Subnet) error
	GetSubnet(subnetId string) (*Subnet, error)
	DeleteSubnet(subnetId string) error
	ListSubnets() ([]*Subnet, error)
	ListSubnetsByVPC(vpcId string) ([]*Subnet, error)

	// SecurityGroup operations
	CreateSecurityGroup(sg *SecurityGroup) error
	GetSecurityGroup(groupId string) (*SecurityGroup, error)
	UpdateSecurityGroup(sg *SecurityGroup) error
	DeleteSecurityGroup(groupId string) error
	ListSecurityGroups() ([]*SecurityGroup, error)
	ListSecurityGroupsByVPC(vpcId string) ([]*SecurityGroup, error)
}
