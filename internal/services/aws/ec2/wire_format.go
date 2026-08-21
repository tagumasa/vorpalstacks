package ec2

import (
	"vorpalstacks/internal/common/protocol"
	types "vorpalstacks/internal/common/tags"
	ec2store "vorpalstacks/internal/store/aws/ec2"
)

// This file converts store structs into response maps keyed with the exact
// EC2 Query protocol wire names (Smithy xmlName traits). The generic
// struct encoder cannot be used directly because several wire names differ
// from the Go field names (Tags -> tagSet, Description -> groupDescription,
// UserIdGroupPairs -> groups) and all list sets must be wrapped in <item>
// elements, which protocol.XMLElements provides explicitly.

// vpcToXMLMap converts a store VPC to a wire-named response map.
func vpcToXMLMap(v *ec2store.VPC) map[string]interface{} {
	m := map[string]interface{}{
		"vpcId":           v.VpcId,
		"cidrBlock":       v.CidrBlock,
		"state":           v.State,
		"ownerId":         v.OwnerId,
		"instanceTenancy": v.InstanceTenancy,
		"dhcpOptionsId":   v.DhcpOptionsId,
		"isDefault":       v.IsDefault,
		"cidrBlockAssociationSet": protocol.XMLElements{
			ElementName: "item",
			Items:       vpcCidrAssocsToXMLItems(v.CidrBlockAssociationSet),
		},
		"tagSet": protocol.XMLElements{
			ElementName: "item",
			Items:       tagsToXMLItems(v.Tags),
		},
	}
	return m
}

// vpcCidrAssocsToXMLItems converts CIDR associations to wire-named maps.
func vpcCidrAssocsToXMLItems(assocs []ec2store.VpcCidrBlockAssociation) []interface{} {
	items := make([]interface{}, 0, len(assocs))
	for _, a := range assocs {
		items = append(items, map[string]interface{}{
			"associationId": a.AssociationId,
			"cidrBlock":     a.CidrBlock,
			"cidrBlockState": map[string]interface{}{
				"state": a.CidrBlockState.State,
			},
		})
	}
	return items
}

// subnetToXMLMap converts a store Subnet to a wire-named response map.
func subnetToXMLMap(sn *ec2store.Subnet) map[string]interface{} {
	return map[string]interface{}{
		"subnetId":                    sn.SubnetId,
		"vpcId":                       sn.VpcId,
		"cidrBlock":                   sn.CidrBlock,
		"availabilityZone":            sn.AvailabilityZone,
		"availableIpAddressCount":     sn.AvailableIpAddressCount,
		"state":                       sn.State,
		"ownerId":                     sn.OwnerId,
		"subnetArn":                   sn.SubnetArn,
		"defaultForAz":                sn.DefaultForAz,
		"mapPublicIpOnLaunch":         sn.MapPublicIpOnLaunch,
		"assignIpv6AddressOnCreation": sn.AssignIpv6AddressOnCreation,
		"ipv6Native":                  sn.Ipv6Native,
		"enableDns64":                 sn.EnableDns64,
		"tagSet": protocol.XMLElements{
			ElementName: "item",
			Items:       tagsToXMLItems(sn.Tags),
		},
	}
}

// securityGroupToXMLMap converts a store SecurityGroup to a wire-named
// response map. The description wire name is groupDescription.
func securityGroupToXMLMap(sg *ec2store.SecurityGroup) map[string]interface{} {
	return map[string]interface{}{
		"groupId":          sg.GroupId,
		"groupName":        sg.GroupName,
		"groupDescription": sg.Description,
		"vpcId":            sg.VpcId,
		"ownerId":          sg.OwnerId,
		"securityGroupArn": sg.SecurityGroupArn,
		"ipPermissions": protocol.XMLElements{
			ElementName: "item",
			Items:       ipRulesToXMLItems(sg.IpPermissions),
		},
		"ipPermissionsEgress": protocol.XMLElements{
			ElementName: "item",
			Items:       ipRulesToXMLItems(sg.IpPermissionsEgress),
		},
		"tagSet": protocol.XMLElements{
			ElementName: "item",
			Items:       tagsToXMLItems(sg.Tags),
		},
	}
}

// ipRulesToXMLItems converts IP permission rules to wire-named maps.
func ipRulesToXMLItems(rules []ec2store.IPRule) []interface{} {
	items := make([]interface{}, 0, len(rules))
	for _, r := range rules {
		items = append(items, ipRuleToXMLMap(r))
	}
	return items
}

// ipRuleToXMLMap converts a single IP permission rule to a wire-named map.
// Port fields are omitted for the all-protocols rule (-1), matching the AWS
// wire format.
func ipRuleToXMLMap(r ec2store.IPRule) map[string]interface{} {
	m := map[string]interface{}{
		"ipProtocol": r.IpProtocol,
	}
	if r.IpProtocol != "-1" {
		m["fromPort"] = r.FromPort
		m["toPort"] = r.ToPort
	}
	if len(r.UserIdGroupPairs) > 0 {
		groups := make([]interface{}, 0, len(r.UserIdGroupPairs))
		for _, g := range r.UserIdGroupPairs {
			gm := map[string]interface{}{}
			if g.GroupId != "" {
				gm["groupId"] = g.GroupId
			}
			if g.GroupName != "" {
				gm["groupName"] = g.GroupName
			}
			if g.UserId != "" {
				gm["userId"] = g.UserId
			}
			if g.VpcId != "" {
				gm["vpcId"] = g.VpcId
			}
			if g.PeeringStatus != "" {
				gm["peeringStatus"] = g.PeeringStatus
			}
			if g.Description != "" {
				gm["description"] = g.Description
			}
			if g.VpcPeeringConnectionId != "" {
				gm["vpcPeeringConnectionId"] = g.VpcPeeringConnectionId
			}
			groups = append(groups, gm)
		}
		m["groups"] = protocol.XMLElements{ElementName: "item", Items: groups}
	}
	if len(r.IpRanges) > 0 {
		ranges := make([]interface{}, 0, len(r.IpRanges))
		for _, ir := range r.IpRanges {
			rm := map[string]interface{}{"cidrIp": ir.CidrIp}
			if ir.Description != "" {
				rm["description"] = ir.Description
			}
			ranges = append(ranges, rm)
		}
		m["ipRanges"] = protocol.XMLElements{ElementName: "item", Items: ranges}
	}
	if len(r.Ipv6Ranges) > 0 {
		ranges := make([]interface{}, 0, len(r.Ipv6Ranges))
		for _, ir := range r.Ipv6Ranges {
			rm := map[string]interface{}{"cidrIpv6": ir.CidrIp}
			if ir.Description != "" {
				rm["description"] = ir.Description
			}
			ranges = append(ranges, rm)
		}
		m["ipv6Ranges"] = protocol.XMLElements{ElementName: "item", Items: ranges}
	}
	if len(r.PrefixListIds) > 0 {
		pls := make([]interface{}, 0, len(r.PrefixListIds))
		for _, pl := range r.PrefixListIds {
			pm := map[string]interface{}{"prefixListId": pl.PrefixListId}
			if pl.Description != "" {
				pm["description"] = pl.Description
			}
			pls = append(pls, pm)
		}
		m["prefixListIds"] = protocol.XMLElements{ElementName: "item", Items: pls}
	}
	return m
}

// tagsToXMLItems converts store tags to wire-named key/value maps.
func tagsToXMLItems(tags []types.Tag) []interface{} {
	items := make([]interface{}, 0, len(tags))
	for _, t := range tags {
		items = append(items, map[string]interface{}{"key": t.Key, "value": t.Value})
	}
	return items
}
