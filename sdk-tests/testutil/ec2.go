package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"vorpalstacks-sdk-tests/config"
)

// RunEC2Tests exercises the EC2 VPC/Subnet/SecurityGroup data plane with the
// AWS SDK, verifying response wire shapes (securityGroupRuleSet, tagSet,
// groupDescription, cidrBlockAssociationSet, ipPermissions) that were
// previously invisible to SDK clients.
func (r *TestRunner) RunEC2Tests() []TestResult {
	var results []TestResult

	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return append(results, TestResult{
			Service:  "ec2",
			TestName: "Setup",
			Status:   "FAIL",
			Error:    fmt.Sprintf("Failed to load config: %v", err),
		})
	}

	client := ec2.NewFromConfig(cfg)
	ctx := context.Background()
	ts := fmt.Sprintf("%d", time.Now().UnixNano())

	results = append(results, r.ec2VPCTests(ctx, client, ts)...)
	results = append(results, r.ec2SubnetTests(ctx, client, ts)...)
	results = append(results, r.ec2SecurityGroupTests(ctx, client, ts)...)
	results = append(results, r.ec2RuleIDTests(ctx, client, ts)...)
	results = append(results, r.ec2ValidationTests(ctx, client, ts)...)

	return results
}

func ec2Name(ts, name string) string {
	return fmt.Sprintf("ec2-%s-%s", name, ts)
}

func (r *TestRunner) ec2VPCTests(ctx context.Context, client *ec2.Client, ts string) []TestResult {
	var results []TestResult
	pass := func(name string) TestResult {
		return TestResult{Service: "ec2", TestName: name, Status: "PASS"}
	}
	fail := func(name string, err error) TestResult {
		return TestResult{Service: "ec2", TestName: name, Status: "FAIL", Error: err.Error()}
	}
	record := func(name string, err error) TestResult {
		if err != nil {
			return fail(name, err)
		}
		return pass(name)
	}

	// CreateVpc with tags: response must expose vpcId, cidrBlockAssociationSet and tagSet.
	vpcName := ec2Name(ts, "vpc")
	createOut, err := client.CreateVpc(ctx, &ec2.CreateVpcInput{
		CidrBlock: aws.String("10.88.0.0/16"),
		TagSpecifications: []ec2types.TagSpecification{{
			ResourceType: ec2types.ResourceTypeVpc,
			Tags: []ec2types.Tag{
				{Key: aws.String("Name"), Value: aws.String(vpcName)},
			},
		}},
	})
	if err != nil {
		return append(results, fail("CreateVpc", err))
	}
	results = append(results, pass("CreateVpc"))
	vpcID := *createOut.Vpc.VpcId

	defer func() {
		_, _ = client.DeleteVpc(ctx, &ec2.DeleteVpcInput{VpcId: aws.String(vpcID)})
	}()

	results = append(results, record("CreateVpc_CidrAssociationSetVisible", func() error {
		if len(createOut.Vpc.CidrBlockAssociationSet) == 0 {
			return fmt.Errorf("cidrBlockAssociationSet empty in CreateVpc response")
		}
		assoc := createOut.Vpc.CidrBlockAssociationSet[0]
		if assoc.CidrBlock == nil || *assoc.CidrBlock != "10.88.0.0/16" {
			return fmt.Errorf("association cidrBlock = %v", assoc.CidrBlock)
		}
		if assoc.AssociationId == nil || *assoc.AssociationId == "" {
			return fmt.Errorf("associationId missing")
		}
		return nil
	}()))

	results = append(results, record("CreateVpc_TagsVisible", func() error {
		if len(createOut.Vpc.Tags) == 0 {
			return fmt.Errorf("tagSet empty in CreateVpc response")
		}
		if *createOut.Vpc.Tags[0].Key != "Name" || *createOut.Vpc.Tags[0].Value != vpcName {
			return fmt.Errorf("unexpected tags: %+v", createOut.Vpc.Tags)
		}
		return nil
	}()))

	// DescribeVpcs must return the tag set (previously invisible: tags vs tagSet).
	results = append(results, record("DescribeVpcs_TagSetVisible", func() error {
		out, err := client.DescribeVpcs(ctx, &ec2.DescribeVpcsInput{VpcIds: []string{vpcID}})
		if err != nil {
			return err
		}
		if len(out.Vpcs) != 1 {
			return fmt.Errorf("expected 1 VPC, got %d", len(out.Vpcs))
		}
		if len(out.Vpcs[0].Tags) != 1 {
			return fmt.Errorf("expected 1 tag on described VPC, got %d", len(out.Vpcs[0].Tags))
		}
		if len(out.Vpcs[0].CidrBlockAssociationSet) == 0 {
			return fmt.Errorf("cidrBlockAssociationSet empty in DescribeVpcs response")
		}
		return nil
	}()))

	// Filter by tag-key must narrow results.
	results = append(results, record("DescribeVpcs_FilterTagKey", func() error {
		out, err := client.DescribeVpcs(ctx, &ec2.DescribeVpcsInput{
			Filters: []ec2types.Filter{{Name: aws.String("tag:Name"), Values: []string{vpcName}}},
		})
		if err != nil {
			return err
		}
		found := false
		for _, v := range out.Vpcs {
			if *v.VpcId == vpcID {
				found = true
			}
		}
		if !found {
			return fmt.Errorf("tag filter did not match the created VPC (%d results)", len(out.Vpcs))
		}
		return nil
	}()))

	// Pagination: MaxResults below 5 must be rejected.
	results = append(results, record("DescribeVpcs_MaxResultsBelowRangeRejected", func() error {
		_, err := client.DescribeVpcs(ctx, &ec2.DescribeVpcsInput{MaxResults: aws.Int32(3)})
		if err == nil {
			return fmt.Errorf("MaxResults=3 should be rejected")
		}
		return nil
	}()))

	// Unknown filter must be rejected.
	results = append(results, record("DescribeVpcs_UnknownFilterRejected", func() error {
		_, err := client.DescribeVpcs(ctx, &ec2.DescribeVpcsInput{
			Filters: []ec2types.Filter{{Name: aws.String("not-a-real-filter"), Values: []string{"x"}}},
		})
		if err == nil {
			return fmt.Errorf("unknown filter should be rejected")
		}
		return nil
	}()))

	// Documented filters for attributes without stored data are accepted and
	// simply match nothing.
	results = append(results, record("DescribeVpcs_IPv6AssociationFilterAccepted", func() error {
		out, err := client.DescribeVpcs(ctx, &ec2.DescribeVpcsInput{
			Filters: []ec2types.Filter{{Name: aws.String("ipv6-cidr-block-association.state"), Values: []string{"associated"}}},
		})
		if err != nil {
			return err
		}
		if len(out.Vpcs) != 0 {
			return fmt.Errorf("ipv6 association filter matched %d VPCs, want 0", len(out.Vpcs))
		}
		return nil
	}()))

	// MaxResults combined with explicit IDs is rejected.
	results = append(results, record("DescribeVpcs_MaxResultsWithIDsRejected", func() error {
		_, err := client.DescribeVpcs(ctx, &ec2.DescribeVpcsInput{
			VpcIds:     []string{vpcID},
			MaxResults: aws.Int32(10),
		})
		if err == nil {
			return fmt.Errorf("MaxResults + VpcIds accepted")
		}
		return nil
	}()))

	return results
}

func (r *TestRunner) ec2SubnetTests(ctx context.Context, client *ec2.Client, ts string) []TestResult {
	var results []TestResult
	pass := func(name string) TestResult {
		return TestResult{Service: "ec2", TestName: name, Status: "PASS"}
	}
	fail := func(name string, err error) TestResult {
		return TestResult{Service: "ec2", TestName: name, Status: "FAIL", Error: err.Error()}
	}

	vpcOut, err := client.CreateVpc(ctx, &ec2.CreateVpcInput{CidrBlock: aws.String("10.89.0.0/16")})
	if err != nil {
		return append(results, fail("SubnetSetup_CreateVpc", err))
	}
	vpcID := *vpcOut.Vpc.VpcId
	defer func() {
		_, _ = client.DeleteVpc(ctx, &ec2.DeleteVpcInput{VpcId: aws.String(vpcID)})
	}()

	subnetOut, err := client.CreateSubnet(ctx, &ec2.CreateSubnetInput{
		VpcId:            aws.String(vpcID),
		CidrBlock:        aws.String("10.89.1.0/24"),
		AvailabilityZone: aws.String(r.region + "a"),
	})
	if err != nil {
		return append(results, fail("CreateSubnet", err))
	}
	results = append(results, pass("CreateSubnet"))
	subnetID := *subnetOut.Subnet.SubnetId
	defer func() {
		_, _ = client.DeleteSubnet(ctx, &ec2.DeleteSubnetInput{SubnetId: aws.String(subnetID)})
	}()

	results = append(results, func() TestResult {
		out, err := client.DescribeSubnets(ctx, &ec2.DescribeSubnetsInput{SubnetIds: []string{subnetID}})
		if err != nil {
			return fail("DescribeSubnets_RoundTrip", err)
		}
		if len(out.Subnets) != 1 {
			return fail("DescribeSubnets_RoundTrip", fmt.Errorf("expected 1 subnet, got %d", len(out.Subnets)))
		}
		sn := out.Subnets[0]
		if sn.AvailableIpAddressCount == nil || *sn.AvailableIpAddressCount != 251 {
			return fail("DescribeSubnets_RoundTrip", fmt.Errorf("availableIpAddressCount = %v, want 251", sn.AvailableIpAddressCount))
		}
		if sn.SubnetArn == nil || *sn.SubnetArn == "" {
			return fail("DescribeSubnets_RoundTrip", fmt.Errorf("subnetArn missing"))
		}
		return pass("DescribeSubnets_RoundTrip")
	}())

	// Deleting the VPC with a live subnet must fail with DependencyViolation.
	results = append(results, func() TestResult {
		_, err := client.DeleteVpc(ctx, &ec2.DeleteVpcInput{VpcId: aws.String(vpcID)})
		if err == nil {
			return fail("DeleteVpc_DependencyViolation", fmt.Errorf("VPC with subnet deleted"))
		}
		return pass("DeleteVpc_DependencyViolation")
	}())

	// Documented enable-dns64 filter is accepted and matches the default.
	results = append(results, func() TestResult {
		out, err := client.DescribeSubnets(ctx, &ec2.DescribeSubnetsInput{
			Filters: []ec2types.Filter{{Name: aws.String("enable-dns64"), Values: []string{"false"}}},
		})
		if err != nil {
			return fail("DescribeSubnets_Dns64FilterAccepted", err)
		}
		found := false
		for _, sn := range out.Subnets {
			if *sn.SubnetId == subnetID {
				found = true
			}
		}
		if !found {
			return fail("DescribeSubnets_Dns64FilterAccepted", fmt.Errorf("created subnet not matched by enable-dns64=false"))
		}
		return pass("DescribeSubnets_Dns64FilterAccepted")
	}())

	// MaxResults combined with explicit IDs is rejected.
	results = append(results, func() TestResult {
		_, err := client.DescribeSubnets(ctx, &ec2.DescribeSubnetsInput{
			SubnetIds:  []string{subnetID},
			MaxResults: aws.Int32(10),
		})
		if err == nil {
			return fail("DescribeSubnets_MaxResultsWithIDsRejected", fmt.Errorf("MaxResults + SubnetIds accepted"))
		}
		return pass("DescribeSubnets_MaxResultsWithIDsRejected")
	}())

	return results
}

func (r *TestRunner) ec2SecurityGroupTests(ctx context.Context, client *ec2.Client, ts string) []TestResult {
	var results []TestResult
	pass := func(name string) TestResult {
		return TestResult{Service: "ec2", TestName: name, Status: "PASS"}
	}
	fail := func(name string, err error) TestResult {
		return TestResult{Service: "ec2", TestName: name, Status: "FAIL", Error: err.Error()}
	}
	sgName := ec2Name(ts, "sg")

	vpcOut, err := client.CreateVpc(ctx, &ec2.CreateVpcInput{CidrBlock: aws.String("10.90.0.0/16")})
	if err != nil {
		return append(results, fail("SGSetup_CreateVpc", err))
	}
	vpcID := *vpcOut.Vpc.VpcId
	defer func() {
		_, _ = client.DeleteVpc(ctx, &ec2.DeleteVpcInput{VpcId: aws.String(vpcID)})
	}()

	sgOut, err := client.CreateSecurityGroup(ctx, &ec2.CreateSecurityGroupInput{
		GroupName:   aws.String(sgName),
		Description: aws.String("ec2 sdk test group"),
		VpcId:       aws.String(vpcID),
	})
	if err != nil {
		return append(results, fail("CreateSecurityGroup", err))
	}
	results = append(results, pass("CreateSecurityGroup"))
	sgID := *sgOut.GroupId
	defer func() {
		_, _ = client.DeleteSecurityGroup(ctx, &ec2.DeleteSecurityGroupInput{GroupId: aws.String(sgID)})
	}()

	// AuthorizeSecurityGroupIngress must return the securityGroupRuleSet with
	// sgr- rule IDs (previously invisible: securityGroupRules wire mismatch).
	var ruleID string
	results = append(results, func() TestResult {
		out, err := client.AuthorizeSecurityGroupIngress(ctx, &ec2.AuthorizeSecurityGroupIngressInput{
			GroupId: aws.String(sgID),
			IpPermissions: []ec2types.IpPermission{{
				IpProtocol: aws.String("tcp"),
				FromPort:   aws.Int32(80),
				ToPort:     aws.Int32(80),
				IpRanges:   []ec2types.IpRange{{CidrIp: aws.String("203.0.113.0/24"), Description: aws.String("web")}},
			}},
		})
		if err != nil {
			return fail("AuthorizeIngress_RuleSetVisible", err)
		}
		if len(out.SecurityGroupRules) == 0 {
			return fail("AuthorizeIngress_RuleSetVisible", fmt.Errorf("securityGroupRules empty in response"))
		}
		rule := out.SecurityGroupRules[0]
		if rule.SecurityGroupRuleId == nil || *rule.SecurityGroupRuleId == "" {
			return fail("AuthorizeIngress_RuleSetVisible", fmt.Errorf("securityGroupRuleId missing"))
		}
		ruleID = *rule.SecurityGroupRuleId
		if rule.CidrIpv4 == nil || *rule.CidrIpv4 != "203.0.113.0/24" {
			return fail("AuthorizeIngress_RuleSetVisible", fmt.Errorf("cidrIpv4 = %v", rule.CidrIpv4))
		}
		if rule.SecurityGroupRuleArn == nil || *rule.SecurityGroupRuleArn == "" {
			return fail("AuthorizeIngress_RuleSetVisible", fmt.Errorf("rule ARN missing"))
		}
		return pass("AuthorizeIngress_RuleSetVisible")
	}())

	// Duplicate authorisation must be rejected.
	results = append(results, func() TestResult {
		_, err := client.AuthorizeSecurityGroupIngress(ctx, &ec2.AuthorizeSecurityGroupIngressInput{
			GroupId: aws.String(sgID),
			IpPermissions: []ec2types.IpPermission{{
				IpProtocol: aws.String("tcp"),
				FromPort:   aws.Int32(80),
				ToPort:     aws.Int32(80),
				IpRanges:   []ec2types.IpRange{{CidrIp: aws.String("203.0.113.0/24")}},
			}},
		})
		if err == nil {
			return fail("AuthorizeIngress_DuplicateRejected", fmt.Errorf("duplicate rule accepted"))
		}
		return pass("AuthorizeIngress_DuplicateRejected")
	}())

	// A permission with no source at all (no IP range, prefix list or
	// security group) must be rejected: AWS requires exactly one source.
	results = append(results, func() TestResult {
		_, err := client.AuthorizeSecurityGroupIngress(ctx, &ec2.AuthorizeSecurityGroupIngressInput{
			GroupId: aws.String(sgID),
			IpPermissions: []ec2types.IpPermission{{
				IpProtocol: aws.String("tcp"),
				FromPort:   aws.Int32(443),
				ToPort:     aws.Int32(443),
			}},
		})
		if err := AssertErrorContains(err, "InvalidParameterValue"); err != nil {
			return fail("AuthorizeIngress_SourcelessRejected", err)
		}
		return pass("AuthorizeIngress_SourcelessRejected")
	}())

	// The egress direction is symmetric: destinations are required.
	results = append(results, func() TestResult {
		_, err := client.AuthorizeSecurityGroupEgress(ctx, &ec2.AuthorizeSecurityGroupEgressInput{
			GroupId: aws.String(sgID),
			IpPermissions: []ec2types.IpPermission{{
				IpProtocol: aws.String("tcp"),
				FromPort:   aws.Int32(8443),
				ToPort:     aws.Int32(8443),
			}},
		})
		if err := AssertErrorContains(err, "InvalidParameterValue"); err != nil {
			return fail("AuthorizeEgress_SourcelessRejected", err)
		}
		return pass("AuthorizeEgress_SourcelessRejected")
	}())

	// DescribeSecurityGroups must expose groupDescription, ipPermissions with
	// item-wrapped ranges, and the default egress rule.
	results = append(results, func() TestResult {
		out, err := client.DescribeSecurityGroups(ctx, &ec2.DescribeSecurityGroupsInput{GroupIds: []string{sgID}})
		if err != nil {
			return fail("DescribeSecurityGroups_PermissionsVisible", err)
		}
		if len(out.SecurityGroups) != 1 {
			return fail("DescribeSecurityGroups_PermissionsVisible", fmt.Errorf("expected 1 SG, got %d", len(out.SecurityGroups)))
		}
		sg := out.SecurityGroups[0]
		if sg.Description == nil || *sg.Description != "ec2 sdk test group" {
			return fail("DescribeSecurityGroups_PermissionsVisible", fmt.Errorf("groupDescription = %v", sg.Description))
		}
		if len(sg.IpPermissions) != 1 {
			return fail("DescribeSecurityGroups_PermissionsVisible", fmt.Errorf("expected 1 ingress permission, got %d", len(sg.IpPermissions)))
		}
		perm := sg.IpPermissions[0]
		if len(perm.IpRanges) != 1 || *perm.IpRanges[0].CidrIp != "203.0.113.0/24" {
			return fail("DescribeSecurityGroups_PermissionsVisible", fmt.Errorf("ipRanges = %+v", perm.IpRanges))
		}
		if perm.IpRanges[0].Description == nil || *perm.IpRanges[0].Description != "web" {
			return fail("DescribeSecurityGroups_PermissionsVisible", fmt.Errorf("range description = %v", perm.IpRanges[0].Description))
		}
		if len(sg.IpPermissionsEgress) != 1 {
			return fail("DescribeSecurityGroups_PermissionsVisible", fmt.Errorf("default egress rule invisible, got %d", len(sg.IpPermissionsEgress)))
		}
		return pass("DescribeSecurityGroups_PermissionsVisible")
	}())

	// RevokeSecurityGroupEgress by rule ID (the default all-traffic rule):
	// response must carry revokedSecurityGroupRuleSet entries.
	results = append(results, func() TestResult {
		out, err := client.RevokeSecurityGroupEgress(ctx, &ec2.RevokeSecurityGroupEgressInput{
			GroupId: aws.String(sgID),
			IpPermissions: []ec2types.IpPermission{{
				IpProtocol: aws.String("-1"),
				IpRanges:   []ec2types.IpRange{{CidrIp: aws.String("0.0.0.0/0")}},
			}},
		})
		if err != nil {
			return fail("RevokeEgress_RuleSetVisible", err)
		}
		if len(out.RevokedSecurityGroupRules) == 0 {
			return fail("RevokeEgress_RuleSetVisible", fmt.Errorf("revokedSecurityGroupRules empty"))
		}
		rev := out.RevokedSecurityGroupRules[0]
		if rev.CidrIpv4 == nil || *rev.CidrIpv4 != "0.0.0.0/0" {
			return fail("RevokeEgress_RuleSetVisible", fmt.Errorf("revoked cidrIpv4 = %v", rev.CidrIpv4))
		}
		return pass("RevokeEgress_RuleSetVisible")
	}())

	// The plain ip-permission.* filter form (inbound) used by the AWS
	// documentation examples must be accepted and match.
	results = append(results, func() TestResult {
		out, err := client.DescribeSecurityGroups(ctx, &ec2.DescribeSecurityGroupsInput{
			Filters: []ec2types.Filter{{Name: aws.String("ip-permission.cidr"), Values: []string{"203.0.113.0/24"}}},
		})
		if err != nil {
			return fail("DescribeSecurityGroups_PlainPermissionFilter", err)
		}
		found := false
		for _, g := range out.SecurityGroups {
			if *g.GroupId == sgID {
				found = true
			}
		}
		if !found {
			return fail("DescribeSecurityGroups_PlainPermissionFilter", fmt.Errorf("plain ip-permission.cidr filter did not match (%d results)", len(out.SecurityGroups)))
		}
		return pass("DescribeSecurityGroups_PlainPermissionFilter")
	}())

	// RevokeSecurityGroupIngress by SecurityGroupRuleIds.
	if ruleID != "" {
		results = append(results, func() TestResult {
			out, err := client.RevokeSecurityGroupIngress(ctx, &ec2.RevokeSecurityGroupIngressInput{
				GroupId:              aws.String(sgID),
				SecurityGroupRuleIds: []string{ruleID},
			})
			if err != nil {
				return fail("RevokeIngress_ByRuleID", err)
			}
			if len(out.RevokedSecurityGroupRules) == 0 {
				return fail("RevokeIngress_ByRuleID", fmt.Errorf("revoked list empty for rule-ID revoke"))
			}
			if out.RevokedSecurityGroupRules[0].SecurityGroupRuleId == nil ||
				*out.RevokedSecurityGroupRules[0].SecurityGroupRuleId != ruleID {
				return fail("RevokeIngress_ByRuleID", fmt.Errorf("revoked rule ID mismatch"))
			}
			return pass("RevokeIngress_ByRuleID")
		}())
	}

	// Revoke of a non-existent permission must fail InvalidPermission.NotFound.
	results = append(results, func() TestResult {
		_, err := client.RevokeSecurityGroupIngress(ctx, &ec2.RevokeSecurityGroupIngressInput{
			GroupId: aws.String(sgID),
			IpPermissions: []ec2types.IpPermission{{
				IpProtocol: aws.String("tcp"),
				FromPort:   aws.Int32(22),
				ToPort:     aws.Int32(22),
				IpRanges:   []ec2types.IpRange{{CidrIp: aws.String("198.51.100.0/24")}},
			}},
		})
		if err == nil {
			return fail("RevokeIngress_NotFound", fmt.Errorf("revoking absent rule succeeded"))
		}
		return pass("RevokeIngress_NotFound")
	}())

	return results
}

func (r *TestRunner) ec2RuleIDTests(ctx context.Context, client *ec2.Client, ts string) []TestResult {
	var results []TestResult
	pass := func(name string) TestResult {
		return TestResult{Service: "ec2", TestName: name, Status: "PASS"}
	}
	fail := func(name string, err error) TestResult {
		return TestResult{Service: "ec2", TestName: name, Status: "FAIL", Error: err.Error()}
	}
	record := func(name string, err error) TestResult {
		if err != nil {
			return fail(name, err)
		}
		return pass(name)
	}

	vpcOut, err := client.CreateVpc(ctx, &ec2.CreateVpcInput{CidrBlock: aws.String("10.94.0.0/16")})
	if err != nil {
		return append(results, fail("RuleIDSetup_CreateVpc", err))
	}
	vpcID := *vpcOut.Vpc.VpcId
	defer func() {
		_, _ = client.DeleteVpc(ctx, &ec2.DeleteVpcInput{VpcId: aws.String(vpcID)})
	}()

	sgOut, err := client.CreateSecurityGroup(ctx, &ec2.CreateSecurityGroupInput{
		GroupName:   aws.String(ec2Name(ts, "sg-ruleid")),
		Description: aws.String("rule id tests"),
		VpcId:       aws.String(vpcID),
	})
	if err != nil {
		return append(results, fail("RuleIDSetup_CreateSecurityGroup", err))
	}
	sgID := *sgOut.GroupId
	defer func() {
		_, _ = client.DeleteSecurityGroup(ctx, &ec2.DeleteSecurityGroupInput{GroupId: aws.String(sgID)})
	}()

	// IPv6-only rule: revoke by rule ID must report the IPv6 CIDR (the
	// revoked list was previously empty because the IPv6 removal helper
	// iterated the IPv4 ranges).
	var v6RuleID string
	results = append(results, func() TestResult {
		out, err := client.AuthorizeSecurityGroupIngress(ctx, &ec2.AuthorizeSecurityGroupIngressInput{
			GroupId: aws.String(sgID),
			IpPermissions: []ec2types.IpPermission{{
				IpProtocol: aws.String("tcp"),
				FromPort:   aws.Int32(443),
				ToPort:     aws.Int32(443),
				Ipv6Ranges: []ec2types.Ipv6Range{{CidrIpv6: aws.String("2001:db8:1234:1a00::/64")}},
			}},
		})
		if err != nil {
			return fail("AuthorizeIngress_IPv6RuleID", err)
		}
		if len(out.SecurityGroupRules) != 1 {
			return fail("AuthorizeIngress_IPv6RuleID", fmt.Errorf("expected 1 rule, got %d", len(out.SecurityGroupRules)))
		}
		v6RuleID = *out.SecurityGroupRules[0].SecurityGroupRuleId
		return pass("AuthorizeIngress_IPv6RuleID")
	}())

	results = append(results, record("RevokeIngress_IPv6ByRuleID", func() error {
		out, err := client.RevokeSecurityGroupIngress(ctx, &ec2.RevokeSecurityGroupIngressInput{
			GroupId:              aws.String(sgID),
			SecurityGroupRuleIds: []string{v6RuleID},
		})
		if err != nil {
			return err
		}
		if len(out.RevokedSecurityGroupRules) != 1 {
			return fmt.Errorf("revoked list = %d entries, want 1", len(out.RevokedSecurityGroupRules))
		}
		rev := out.RevokedSecurityGroupRules[0]
		if rev.CidrIpv6 == nil || *rev.CidrIpv6 != "2001:db8:1234:1a00::/64" {
			return fmt.Errorf("revoked cidrIpv6 = %v", rev.CidrIpv6)
		}
		if rev.CidrIpv4 != nil {
			return fmt.Errorf("revoked cidrIpv4 = %v, want nil", rev.CidrIpv4)
		}
		return nil
	}()))

	// Mixed-family permission: revoking the IPv4 entry by rule ID must leave
	// the IPv6 entry intact (the removal previously cancelled itself out and
	// corrupted the stored families).
	var mixedV4ID string
	results = append(results, func() TestResult {
		out, err := client.AuthorizeSecurityGroupIngress(ctx, &ec2.AuthorizeSecurityGroupIngressInput{
			GroupId: aws.String(sgID),
			IpPermissions: []ec2types.IpPermission{{
				IpProtocol: aws.String("tcp"),
				FromPort:   aws.Int32(8080),
				ToPort:     aws.Int32(8080),
				IpRanges:   []ec2types.IpRange{{CidrIp: aws.String("198.51.100.0/24")}},
				Ipv6Ranges: []ec2types.Ipv6Range{{CidrIpv6: aws.String("2001:db8:abcd::/64")}},
			}},
		})
		if err != nil {
			return fail("AuthorizeIngress_MixedFamily", err)
		}
		for _, rule := range out.SecurityGroupRules {
			if rule.CidrIpv4 != nil && *rule.CidrIpv4 == "198.51.100.0/24" {
				mixedV4ID = *rule.SecurityGroupRuleId
			}
		}
		if mixedV4ID == "" {
			return fail("AuthorizeIngress_MixedFamily", fmt.Errorf("IPv4 rule ID missing from %d rules", len(out.SecurityGroupRules)))
		}
		return pass("AuthorizeIngress_MixedFamily")
	}())

	results = append(results, record("RevokeIngress_MixedFamilyIPv4ByID", func() error {
		out, err := client.RevokeSecurityGroupIngress(ctx, &ec2.RevokeSecurityGroupIngressInput{
			GroupId:              aws.String(sgID),
			SecurityGroupRuleIds: []string{mixedV4ID},
		})
		if err != nil {
			return err
		}
		if len(out.RevokedSecurityGroupRules) != 1 {
			return fmt.Errorf("revoked list = %d entries, want 1", len(out.RevokedSecurityGroupRules))
		}
		if out.RevokedSecurityGroupRules[0].CidrIpv4 == nil || *out.RevokedSecurityGroupRules[0].CidrIpv4 != "198.51.100.0/24" {
			return fmt.Errorf("revoked cidrIpv4 = %v", out.RevokedSecurityGroupRules[0].CidrIpv4)
		}
		return nil
	}()))

	results = append(results, record("RevokeIngress_MixedFamilyIPv6Survives", func() error {
		out, err := client.DescribeSecurityGroups(ctx, &ec2.DescribeSecurityGroupsInput{GroupIds: []string{sgID}})
		if err != nil {
			return err
		}
		for _, perm := range out.SecurityGroups[0].IpPermissions {
			if perm.FromPort == nil || *perm.FromPort != 8080 {
				continue
			}
			if len(perm.IpRanges) != 0 {
				return fmt.Errorf("IPv4 range still present after rule-ID revoke")
			}
			if len(perm.Ipv6Ranges) != 1 || *perm.Ipv6Ranges[0].CidrIpv6 != "2001:db8:abcd::/64" {
				return fmt.Errorf("IPv6 range lost or corrupted: %+v", perm.Ipv6Ranges)
			}
			return nil
		}
		return fmt.Errorf("permission for port 8080 disappeared entirely")
	}()))

	// Specifying both SecurityGroupRuleIds and IpPermissions is rejected.
	results = append(results, record("RevokeIngress_RuleIDAndPermissionsRejected", func() error {
		_, err := client.RevokeSecurityGroupIngress(ctx, &ec2.RevokeSecurityGroupIngressInput{
			GroupId: aws.String(sgID),
			IpPermissions: []ec2types.IpPermission{{
				IpProtocol: aws.String("tcp"),
				FromPort:   aws.Int32(8080),
				ToPort:     aws.Int32(8080),
				Ipv6Ranges: []ec2types.Ipv6Range{{CidrIpv6: aws.String("2001:db8:abcd::/64")}},
			}},
			SecurityGroupRuleIds: []string{mixedV4ID},
		})
		if err == nil {
			return fmt.Errorf("rule IDs + permissions accepted")
		}
		return nil
	}()))

	// Authorize with neither IpPermissions nor legacy parameters is rejected.
	results = append(results, record("AuthorizeIngress_NoPermissionsRejected", func() error {
		_, err := client.AuthorizeSecurityGroupIngress(ctx, &ec2.AuthorizeSecurityGroupIngressInput{
			GroupId: aws.String(sgID),
		})
		if err == nil {
			return fmt.Errorf("authorize without permissions accepted")
		}
		return nil
	}()))

	// Revoke with neither rule IDs nor permissions is rejected.
	results = append(results, record("RevokeIngress_NoParametersRejected", func() error {
		_, err := client.RevokeSecurityGroupIngress(ctx, &ec2.RevokeSecurityGroupIngressInput{
			GroupId: aws.String(sgID),
		})
		if err == nil {
			return fmt.Errorf("revoke without parameters accepted")
		}
		return nil
	}()))

	return results
}

func (r *TestRunner) ec2ValidationTests(ctx context.Context, client *ec2.Client, ts string) []TestResult {
	var results []TestResult
	expectFail := func(name string, err error, wantCode string) TestResult {
		if err == nil {
			return TestResult{Service: "ec2", TestName: name, Status: "FAIL", Error: "expected error, got success"}
		}
		return TestResult{Service: "ec2", TestName: name, Status: "PASS"}
	}

	// CreateVpc without CidrBlock is rejected.
	_, err := client.CreateVpc(ctx, &ec2.CreateVpcInput{})
	results = append(results, expectFail("CreateVpc_EmptyCidrRejected", err, ""))

	// CreateVpc rejects an invalid tenancy.
	_, err = client.CreateVpc(ctx, &ec2.CreateVpcInput{
		CidrBlock:       aws.String("10.91.0.0/16"),
		InstanceTenancy: ec2types.Tenancy("bogus"),
	})
	results = append(results, expectFail("CreateVpc_InvalidTenancyRejected", err, ""))

	// CreateSubnet with a CIDR outside the VPC range is rejected.
	vpcOut, verr := client.CreateVpc(ctx, &ec2.CreateVpcInput{CidrBlock: aws.String("10.92.0.0/16")})
	if verr == nil {
		vpcID := *vpcOut.Vpc.VpcId
		defer func() {
			_, _ = client.DeleteVpc(ctx, &ec2.DeleteVpcInput{VpcId: aws.String(vpcID)})
		}()
		_, err = client.CreateSubnet(ctx, &ec2.CreateSubnetInput{
			VpcId:     aws.String(vpcID),
			CidrBlock: aws.String("192.0.2.0/24"),
		})
		results = append(results, expectFail("CreateSubnet_CidrOutsideVpcRejected", err, ""))
	}

	// CreateSecurityGroup without a description is rejected client-side;
	// exercise the empty-VpcId rejection instead by omitting VpcId.
	_, err = client.CreateSecurityGroup(ctx, &ec2.CreateSecurityGroupInput{
		GroupName:   aws.String(ec2Name(ts, "sg-novpc")),
		Description: aws.String("no vpc"),
	})
	results = append(results, expectFail("CreateSecurityGroup_EmptyVpcIdRejected", err, ""))

	// Authorize with an invalid port is rejected.
	vpc2, verr2 := client.CreateVpc(ctx, &ec2.CreateVpcInput{CidrBlock: aws.String("10.93.0.0/16")})
	if verr2 == nil {
		vpcID2 := *vpc2.Vpc.VpcId
		defer func() {
			_, _ = client.DeleteVpc(ctx, &ec2.DeleteVpcInput{VpcId: aws.String(vpcID2)})
		}()
		sgOut, sgErr := client.CreateSecurityGroup(ctx, &ec2.CreateSecurityGroupInput{
			GroupName:   aws.String(ec2Name(ts, "sg-ports")),
			Description: aws.String("ports"),
			VpcId:       aws.String(vpcID2),
		})
		if sgErr == nil {
			sgID := *sgOut.GroupId
			defer func() {
				_, _ = client.DeleteSecurityGroup(ctx, &ec2.DeleteSecurityGroupInput{GroupId: aws.String(sgID)})
			}()
			_, err = client.AuthorizeSecurityGroupIngress(ctx, &ec2.AuthorizeSecurityGroupIngressInput{
				GroupId: aws.String(sgID),
				IpPermissions: []ec2types.IpPermission{{
					IpProtocol: aws.String("tcp"),
					FromPort:   aws.Int32(70000),
					ToPort:     aws.Int32(70000),
					IpRanges:   []ec2types.IpRange{{CidrIp: aws.String("0.0.0.0/0")}},
				}},
			})
			results = append(results, expectFail("AuthorizeIngress_InvalidPortRejected", err, ""))
		}
	}

	return results
}
