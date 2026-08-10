package ec2

import (
	"google.golang.org/protobuf/proto"

	ec2store "vorpalstacks/internal/store/aws/ec2"

	pb "vorpalstacks/internal/pb/aws/ec2"
	"vorpalstacks/internal/utils/aws/types"
)

// ec2Stores is a thin wrapper for the admin handler's store access.
type ec2Stores struct {
	store *ec2store.EC2Store
}

// toPbVpc converts a store VPC to a proto Vpc message.
func toPbVpc(v *ec2store.VPC) *pb.Vpc {
	pbVpc := &pb.Vpc{
		Vpcid:         v.VpcId,
		Cidrblock:     v.CidrBlock,
		Dhcpoptionsid: v.DhcpOptionsId,
		Ownerid:       v.OwnerId,
		Isdefault:     proto.Bool(v.IsDefault),
	}
	if v.State != "" {
		if val, ok := pb.VpcState_value[v.State]; ok {
			pbVpc.State = pb.VpcState(val)
		}
	}
	if v.InstanceTenancy != "" {
		if val, ok := pb.Tenancy_value[v.InstanceTenancy]; ok {
			pbVpc.Instancetenancy = pb.Tenancy(val)
		}
	}
	if len(v.CidrBlockAssociationSet) > 0 {
		pbVpc.Cidrblockassociationset = make([]*pb.VpcCidrBlockAssociation, len(v.CidrBlockAssociationSet))
		for i, assoc := range v.CidrBlockAssociationSet {
			pbAssoc := &pb.VpcCidrBlockAssociation{
				Associationid: assoc.AssociationId,
				Cidrblock:     assoc.CidrBlock,
			}
			if assoc.CidrBlockState.State != "" {
				if val, ok := pb.VpcCidrBlockStateCode_value[assoc.CidrBlockState.State]; ok {
					pbAssoc.Cidrblockstate = &pb.VpcCidrBlockState{
						State: pb.VpcCidrBlockStateCode(val),
					}
				}
			}
			pbVpc.Cidrblockassociationset[i] = pbAssoc
		}
	}
	pbVpc.Tags = toPbTags(v.Tags)
	return pbVpc
}

// toPbSubnet converts a store Subnet to a proto Subnet message.
func toPbSubnet(sn *ec2store.Subnet) *pb.Subnet {
	pbSn := &pb.Subnet{
		Subnetid:                    sn.SubnetId,
		Vpcid:                       sn.VpcId,
		Cidrblock:                   sn.CidrBlock,
		Availabilityzone:            sn.AvailabilityZone,
		Availabilityzoneid:          sn.AvailabilityZoneId,
		Availableipaddresscount:     proto.Int32(int32(sn.AvailableIpAddressCount)),
		Ownerid:                     sn.OwnerId,
		Subnetarn:                   sn.SubnetArn,
		Defaultforaz:                proto.Bool(sn.DefaultForAz),
		Mappubliciponlaunch:         proto.Bool(sn.MapPublicIpOnLaunch),
		Assignipv6Addressoncreation: proto.Bool(sn.AssignIpv6AddressOnCreation),
		Ipv6Native:                  proto.Bool(sn.Ipv6Native),
		Enabledns64:                 proto.Bool(sn.EnableDns64),
	}
	if sn.State != "" {
		if val, ok := pb.SubnetState_value[sn.State]; ok {
			pbSn.State = pb.SubnetState(val)
		}
	}
	pbSn.Tags = toPbTags(sn.Tags)
	return pbSn
}

// toPbSecurityGroup converts a store SecurityGroup to a proto SecurityGroup.
func toPbSecurityGroup(sg *ec2store.SecurityGroup) *pb.SecurityGroup {
	pbSg := &pb.SecurityGroup{
		Groupid:          sg.GroupId,
		Groupname:        sg.GroupName,
		Description:      sg.Description,
		Vpcid:            sg.VpcId,
		Ownerid:          sg.OwnerId,
		Securitygrouparn: sg.SecurityGroupArn,
	}
	if len(sg.IpPermissions) > 0 {
		pbSg.Ippermissions = toPbIpPermissions(sg.IpPermissions)
	}
	if len(sg.IpPermissionsEgress) > 0 {
		pbSg.Ippermissionsegress = toPbIpPermissions(sg.IpPermissionsEgress)
	}
	pbSg.Tags = toPbTags(sg.Tags)
	return pbSg
}

// toPbIpPermissions converts store IPRule slice to proto IpPermission slice.
func toPbIpPermissions(rules []ec2store.IPRule) []*pb.IpPermission {
	result := make([]*pb.IpPermission, len(rules))
	for i, rule := range rules {
		pp := &pb.IpPermission{
			Ipprotocol: rule.IpProtocol,
			Fromport:   proto.Int32(int32(rule.FromPort)),
			Toport:     proto.Int32(int32(rule.ToPort)),
		}
		for _, ip := range rule.IpRanges {
			pp.Ipranges = append(pp.Ipranges, &pb.IpRange{
				Cidrip:      ip.CidrIp,
				Description: ip.Description,
			})
		}
		for _, ip := range rule.Ipv6Ranges {
			pp.Ipv6Ranges = append(pp.Ipv6Ranges, &pb.Ipv6Range{
				Cidripv6:    ip.CidrIp,
				Description: ip.Description,
			})
		}
		for _, pair := range rule.UserIdGroupPairs {
			pp.Useridgrouppairs = append(pp.Useridgrouppairs, &pb.UserIdGroupPair{
				Groupid:                pair.GroupId,
				Groupname:              pair.GroupName,
				Userid:                 pair.UserId,
				Vpcid:                  pair.VpcId,
				Peeringstatus:          pair.PeeringStatus,
				Description:            pair.Description,
				Vpcpeeringconnectionid: pair.VpcPeeringConnectionId,
			})
		}
		for _, pl := range rule.PrefixListIds {
			pp.Prefixlistids = append(pp.Prefixlistids, &pb.PrefixListId{
				Prefixlistid: pl.PrefixListId,
				Description:  pl.Description,
			})
		}
		result[i] = pp
	}
	return result
}

// toPbTags converts store types.Tag slice to proto Tag slice.
func toPbTags(tags []types.Tag) []*pb.Tag {
	if len(tags) == 0 {
		return nil
	}
	result := make([]*pb.Tag, len(tags))
	for i, t := range tags {
		result[i] = &pb.Tag{Key: t.Key, Value: t.Value}
	}
	return result
}
