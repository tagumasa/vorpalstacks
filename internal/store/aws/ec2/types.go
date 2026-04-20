package ec2

// VPC represents an Amazon VPC.
type VPC struct {
	VpcId              string `json:"VpcId"`
	CidrBlock          string `json:"CidrBlock"`
	State              string `json:"State"`
	OwnerId            string `json:"OwnerId"`
	InstanceTenancy    string `json:"InstanceTenancy"`
	EnableDnsSupport   bool   `json:"EnableDnsSupport"`
	EnableDnsHostnames bool   `json:"EnableDnsHostnames"`
	Tags               []Tag  `json:"Tags"`
}

// Subnet represents a VPC subnet.
type Subnet struct {
	SubnetId            string `json:"SubnetId"`
	VpcId               string `json:"VpcId"`
	CidrBlock           string `json:"CidrBlock"`
	AvailabilityZone    string `json:"AvailabilityZone"`
	State               string `json:"State"`
	OwnerId             string `json:"OwnerId"`
	MapPublicIpOnLaunch bool   `json:"MapPublicIpOnLaunch"`
	Tags                []Tag  `json:"Tags"`
}

// SecurityGroup represents a VPC security group.
type SecurityGroup struct {
	GroupId             string   `json:"GroupId"`
	GroupName           string   `json:"GroupName"`
	Description         string   `json:"Description"`
	VpcId               string   `json:"VpcId"`
	OwnerId             string   `json:"OwnerId"`
	Tags                []Tag    `json:"Tags"`
	IpPermissions       []IPRule `json:"IpPermissions"`
	IpPermissionsEgress []IPRule `json:"IpPermissionsEgress"`
}

// IPRule represents an IP permission rule for a security group.
type IPRule struct {
	IpProtocol       string      `json:"IpProtocol"`
	FromPort         int64       `json:"FromPort"`
	ToPort           int64       `json:"ToPort"`
	UserIdGroupPairs []GroupPair `json:"UserIdGroupPairs"`
	IpRanges         []IPRange   `json:"IpRanges"`
	Ipv6Ranges       []IPRange   `json:"Ipv6Ranges"`
}

// GroupPair represents a user ID/group pair reference in a security group rule.
type GroupPair struct {
	GroupId       string `json:"GroupId,omitempty"`
	GroupName     string `json:"GroupName,omitempty"`
	UserId        string `json:"UserId,omitempty"`
	VpcId         string `json:"VpcId,omitempty"`
	PeeringStatus string `json:"PeeringStatus,omitempty"`
}

// IPRange represents a CIDR IP range in a security group rule.
type IPRange struct {
	CidrIp      string `json:"CidrIp,omitempty"`
	Description string `json:"Description,omitempty"`
}

// Tag represents an EC2 resource tag.
type Tag struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}
