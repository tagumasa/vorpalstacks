package ec2

import (
	"testing"

	"vorpalstacks/internal/common/protocol"
	ec2store "vorpalstacks/internal/store/aws/ec2"
	"vorpalstacks/internal/utils/aws/types"
)

// TestSecurityGroupWireNames pins the EC2 Query protocol wire names of the
// DescribeSecurityGroups response: the SDK deserialiser matches
// groupDescription (not description), tagSet (not tags) and groups (not
// userIdGroupPairs), and all nested sets must use <item> elements.
func TestSecurityGroupWireNames(t *testing.T) {
	sg := &ec2store.SecurityGroup{
		GroupId:     "sg-0123456789abcdef0",
		GroupName:   "web",
		Description: "web tier",
		VpcId:       "vpc-0123456789abcdef0",
		OwnerId:     "123456789012",
		IpPermissions: []ec2store.IPRule{{
			IpProtocol: "tcp",
			FromPort:   80,
			ToPort:     80,
			IpRanges: []ec2store.IPRange{{
				RuleId:      "sgr-0123456789abcdef0",
				CidrIp:      "0.0.0.0/0",
				Description: "open web",
			}},
			UserIdGroupPairs: []ec2store.GroupPair{{GroupId: "sg-fedcba9876543210", UserId: "123456789012"}},
		}},
		IpPermissionsEgress: []ec2store.IPRule{{
			IpProtocol: "-1",
			FromPort:   -1,
			ToPort:     -1,
			IpRanges:   []ec2store.IPRange{{CidrIp: "0.0.0.0/0"}},
		}},
		Tags: []types.Tag{{Key: "env", Value: "test"}},
	}

	m := securityGroupToXMLMap(sg)
	for _, key := range []string{"groupId", "groupName", "groupDescription", "vpcId", "ownerId"} {
		if _, ok := m[key]; !ok {
			t.Errorf("missing wire key %q in %v", key, m)
		}
	}
	if _, ok := m["description"]; ok {
		t.Error("description key must not be present (wire name is groupDescription)")
	}

	ipPerms := m["ipPermissions"].(protocol.XMLElements)
	if ipPerms.ElementName != "item" {
		t.Errorf("ipPermissions element name = %q, want item", ipPerms.ElementName)
	}
	first := ipPerms.Items[0].(map[string]interface{})
	if first["fromPort"] != int64(80) {
		t.Errorf("fromPort = %v, want 80", first["fromPort"])
	}
	ranges := first["ipRanges"].(protocol.XMLElements)
	if ranges.ElementName != "item" {
		t.Errorf("ipRanges element name = %q, want item", ranges.ElementName)
	}
	if ranges.Items[0].(map[string]interface{})["cidrIp"] != "0.0.0.0/0" {
		t.Error("cidrIp not set correctly")
	}
	groups := first["groups"].(protocol.XMLElements)
	if groups.Items[0].(map[string]interface{})["groupId"] != "sg-fedcba9876543210" {
		t.Error("groups entry missing groupId")
	}

	egress := m["ipPermissionsEgress"].(protocol.XMLElements)
	egressRule := egress.Items[0].(map[string]interface{})
	if _, ok := egressRule["fromPort"]; ok {
		t.Error("protocol -1 rule must omit fromPort")
	}

	tags := m["tagSet"].(protocol.XMLElements)
	if tags.ElementName != "item" {
		t.Errorf("tagSet element name = %q, want item", tags.ElementName)
	}
}

// TestVPCWireNames pins the VPC response wire names, including the
// cidrBlockAssociationSet structure.
func TestVPCWireNames(t *testing.T) {
	vpc := &ec2store.VPC{
		VpcId:     "vpc-0123456789abcdef0",
		CidrBlock: "10.0.0.0/16",
		State:     "available",
		CidrBlockAssociationSet: []ec2store.VpcCidrBlockAssociation{{
			AssociationId:  "vpc-cidr-assoc-0123456789abcdef0",
			CidrBlock:      "10.0.0.0/16",
			CidrBlockState: ec2store.VpcCidrBlockAssociationState{State: "associated"},
		}},
		Tags: []types.Tag{{Key: "Name", Value: "main"}},
	}
	m := vpcToXMLMap(vpc)
	assocSet := m["cidrBlockAssociationSet"].(protocol.XMLElements)
	if assocSet.ElementName != "item" {
		t.Errorf("cidrBlockAssociationSet element = %q, want item", assocSet.ElementName)
	}
	assoc := assocSet.Items[0].(map[string]interface{})
	if assoc["associationId"] != "vpc-cidr-assoc-0123456789abcdef0" {
		t.Error("associationId missing")
	}
	inner := assoc["cidrBlockState"].(map[string]interface{})
	if inner["state"] != "associated" {
		t.Errorf("cidrBlockState.state = %v", inner["state"])
	}
	if _, ok := m["tagSet"]; !ok {
		t.Error("tagSet missing from VPC wire map")
	}
}

// TestIPv6RangeWireName pins that IPv6 ranges serialise under cidrIpv6 (the
// store reuses the CidrIp field for IPv6 addresses).
func TestIPv6RangeWireName(t *testing.T) {
	m := ipRuleToXMLMap(ec2store.IPRule{
		IpProtocol: "tcp",
		FromPort:   443,
		ToPort:     443,
		Ipv6Ranges: []ec2store.IPRange{{CidrIp: "::/0"}},
	})
	ranges := m["ipv6Ranges"].(protocol.XMLElements)
	entry := ranges.Items[0].(map[string]interface{})
	if entry["cidrIpv6"] != "::/0" {
		t.Errorf("cidrIpv6 = %v, want ::/0", entry["cidrIpv6"])
	}
	if _, ok := entry["cidrIp"]; ok {
		t.Error("IPv6 range must not emit cidrIp")
	}
}

// TestSGRuleXMLMapOmitsPortsForAllProtocols pins the Authorize response rule
// shape for protocol -1.
func TestSGRuleXMLMapOmitsPortsForAllProtocols(t *testing.T) {
	m := sgRuleToXMLMap(SecurityGroupRule{
		RuleId:     "sgr-x",
		IpProtocol: "-1",
		FromPort:   -1,
		ToPort:     -1,
		CidrIpv4:   "0.0.0.0/0",
	})
	if _, ok := m["FromPort"]; ok {
		t.Error("protocol -1 must omit FromPort")
	}
	if _, ok := m["ToPort"]; ok {
		t.Error("protocol -1 must omit ToPort")
	}
	if m["CidrIpv4"] != "0.0.0.0/0" {
		t.Error("CidrIpv4 missing")
	}
}

// TestRevokedRuleXMLMap pins the RevokedSecurityGroupRule shape: rule ID and
// referencedGroupId are present, no ARN.
func TestRevokedRuleXMLMap(t *testing.T) {
	m := revokedRuleToXMLMap(SecurityGroupRule{
		RuleId:            "sgr-x",
		IpProtocol:        "tcp",
		FromPort:          80,
		ToPort:            80,
		ReferencedGroupId: "sg-y",
	})
	if m["SecurityGroupRuleId"] != "sgr-x" {
		t.Error("rule id missing")
	}
	if m["ReferencedGroupId"] != "sg-y" {
		t.Error("referencedGroupId missing")
	}
	if _, ok := m["SecurityGroupRuleArn"]; ok {
		t.Error("RevokedSecurityGroupRule must not carry an ARN")
	}
}
