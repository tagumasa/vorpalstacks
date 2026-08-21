package ec2

import (
	"strings"
	"testing"

	awserrors "vorpalstacks/internal/common/errors"
	ec2store "vorpalstacks/internal/store/aws/ec2"
)

func TestParseInt64Param(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		param     string
		want      int64
		wantError bool
	}{
		{"empty returns -1", "", "FromPort", -1, false},
		{"valid port", "80", "FromPort", 80, false},
		{"valid max results", "50", "MaxResults", 50, false},
		{"invalid number", "abc", "MaxResults", 0, true},
		{"error names parameter", "xx", "ToPort", 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseInt64Param(tt.input, tt.param)
			if tt.wantError {
				if err == nil {
					t.Fatalf("expected error for input %q", tt.input)
				}
				if !strings.Contains(err.Error(), tt.param) {
					t.Errorf("error %q does not name parameter %q", err.Error(), tt.param)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("got %d, want %d", got, tt.want)
			}
		})
	}
}

func TestValidateInstanceTenancy(t *testing.T) {
	for _, valid := range []string{"", "default", "dedicated", "host"} {
		if err := validateInstanceTenancy(valid); err != nil {
			t.Errorf("tenancy %q should be valid, got %v", valid, err)
		}
	}
	if err := validateInstanceTenancy("bogus"); err == nil {
		t.Error("tenancy bogus should be rejected")
	}
}

func TestParsePaginationParams(t *testing.T) {
	tests := []struct {
		name      string
		params    map[string]interface{}
		maxWant   int
		tokenWant string
		wantError bool
	}{
		{"omitted returns no limit", map[string]interface{}{}, 0, "", false},
		{"valid value", map[string]interface{}{"MaxResults": "100"}, 100, "", false},
		{"lower bound", map[string]interface{}{"MaxResults": "5"}, 5, "", false},
		{"upper bound", map[string]interface{}{"MaxResults": "1000"}, 1000, "", false},
		{"below range", map[string]interface{}{"MaxResults": "4"}, 0, "", true},
		{"above range", map[string]interface{}{"MaxResults": "1001"}, 0, "", true},
		{"non numeric", map[string]interface{}{"MaxResults": "abc"}, 0, "", true},
		{"next token", map[string]interface{}{"NextToken": "vpc-1"}, 0, "vpc-1", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, max, err := parsePaginationParams(tt.params)
			if tt.wantError {
				if err == nil {
					t.Fatalf("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if max != tt.maxWant {
				t.Errorf("max = %d, want %d", max, tt.maxWant)
			}
			if token != tt.tokenWant {
				t.Errorf("token = %q, want %q", token, tt.tokenWant)
			}
		})
	}
}

func TestValidateFilterNames(t *testing.T) {
	if err := validateFilterNames([]ec2Filter{{Name: "vpc-id", Values: []string{"x"}}}, allowedVPCFilters); err != nil {
		t.Errorf("vpc-id should be allowed: %v", err)
	}
	if err := validateFilterNames([]ec2Filter{{Name: "no-such-filter", Values: []string{"x"}}}, allowedVPCFilters); err == nil {
		t.Error("unknown filter should be rejected")
	}
	if err := validateFilterNames([]ec2Filter{{Name: "ingress.ip-permission.cidr", Values: []string{"0.0.0.0/0"}}}, allowedSGFilters); err != nil {
		t.Errorf("ingress.ip-permission.cidr should be allowed: %v", err)
	}
	if err := validateFilterNames([]ec2Filter{{Name: "ingress.ip-permission.nonsense", Values: []string{"x"}}}, allowedSGFilters); err == nil {
		t.Error("unknown ip-permission sub-filter should be rejected")
	}
}

func TestParseLegacyFlatRule(t *testing.T) {
	rules, err := parseLegacyFlatRule(map[string]interface{}{
		"IpProtocol": "tcp",
		"FromPort":   "80",
		"ToPort":     "80",
		"CidrIp":     "0.0.0.0/0",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(rules))
	}
	r := rules[0]
	if r.IpProtocol != "tcp" || r.FromPort != 80 || r.ToPort != 80 {
		t.Errorf("unexpected rule: %+v", r)
	}
	if len(r.IpRanges) != 1 || r.IpRanges[0].CidrIp != "0.0.0.0/0" {
		t.Errorf("unexpected ipRanges: %+v", r.IpRanges)
	}

	// Flat parameters absent entirely: no rule, no error.
	rules, err = parseLegacyFlatRule(map[string]interface{}{})
	if err != nil || len(rules) != 0 {
		t.Errorf("empty params should produce no rule, got %v, %v", rules, err)
	}

	// Protocol without any source is malformed.
	_, err = parseLegacyFlatRule(map[string]interface{}{"IpProtocol": "tcp"})
	if err == nil {
		t.Error("permission without a source should be rejected")
	}

	// Invalid CIDR is rejected.
	_, err = parseLegacyFlatRule(map[string]interface{}{
		"IpProtocol": "tcp", "CidrIp": "999.0.0.0/24",
	})
	if err == nil {
		t.Error("invalid CIDR should be rejected")
	}
}

func TestParseIPRulesRequiresSource(t *testing.T) {
	// A permission with a protocol and ports but no source at all is
	// rejected: AWS requires exactly one of an IP range, a prefix list or
	// a security group.
	_, err := parseIPRules(map[string]interface{}{
		"IpPermissions.1.IpProtocol": "tcp",
		"IpPermissions.1.FromPort":   "80",
		"IpPermissions.1.ToPort":     "80",
	}, "IpPermissions")
	if err == nil {
		t.Fatal("permission without any source should be rejected")
	}
	if apiErr, ok := err.(*awserrors.AWSError); !ok || apiErr.Code != "InvalidParameterValue" {
		t.Errorf("expected InvalidParameterValue, got %v", err)
	}

	// The same permission with an IPv4 range parses.
	rules, err := parseIPRules(map[string]interface{}{
		"IpPermissions.1.IpProtocol":        "tcp",
		"IpPermissions.1.FromPort":          "80",
		"IpPermissions.1.ToPort":            "80",
		"IpPermissions.1.IpRanges.1.CidrIp": "0.0.0.0/0",
	}, "IpPermissions")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 1 || len(rules[0].IpRanges) != 1 {
		t.Errorf("unexpected rules: %+v", rules)
	}
}

func TestParseSecurityGroupRuleIDs(t *testing.T) {
	params := map[string]interface{}{
		"SecurityGroupRuleId.1": "sgr-aaa",
		"SecurityGroupRuleId.2": "sgr-bbb",
	}
	ids := parseSecurityGroupRuleIDs(params)
	if len(ids) != 2 || ids[0] != "sgr-aaa" || ids[1] != "sgr-bbb" {
		t.Errorf("unexpected ids: %v", ids)
	}
}

func TestFindDuplicateEntry(t *testing.T) {
	existing := []ec2store.IPRule{{
		IpProtocol: "tcp",
		FromPort:   80,
		ToPort:     80,
		IpRanges:   []ec2store.IPRange{{CidrIp: "10.0.0.0/8"}},
	}}
	requestedSame := []ec2store.IPRule{{
		IpProtocol: "tcp",
		FromPort:   80,
		ToPort:     80,
		IpRanges:   []ec2store.IPRange{{CidrIp: "10.0.0.0/8"}},
	}}
	requestedNew := []ec2store.IPRule{{
		IpProtocol: "tcp",
		FromPort:   80,
		ToPort:     80,
		IpRanges:   []ec2store.IPRange{{CidrIp: "192.168.0.0/16"}},
	}}
	if findDuplicateEntry(existing, requestedSame) == "" {
		t.Error("identical entry should be detected as duplicate")
	}
	if findDuplicateEntry(existing, requestedNew) != "" {
		t.Error("new CIDR should not be a duplicate")
	}
}

func TestParsePaginationParamsIDExclusion(t *testing.T) {
	// MaxResults together with a list of resource IDs is rejected with
	// InvalidParameterCombination per the EC2 API Reference pagination notes.
	bad := map[string]interface{}{"MaxResults": "100", "VpcId.1": "vpc-1"}
	if _, _, err := parsePaginationParams(bad, "VpcId"); err == nil {
		t.Error("MaxResults + VpcId.1 should be rejected")
	}
	// The documented rule covers ID lists only; a name list such as
	// GroupName must remain combinable with MaxResults.
	names := map[string]interface{}{"MaxResults": "100", "GroupName.1": "default"}
	if _, _, err := parsePaginationParams(names, "GroupId"); err != nil {
		t.Errorf("MaxResults + GroupName.1 should be accepted: %v", err)
	}
	okIDs := map[string]interface{}{"VpcId.1": "vpc-1"}
	if _, max, err := parsePaginationParams(okIDs, "VpcId"); err != nil || max != 0 {
		t.Errorf("IDs without MaxResults should pass: err=%v max=%d", err, max)
	}
	okLimit := map[string]interface{}{"MaxResults": "100"}
	if _, max, err := parsePaginationParams(okLimit, "VpcId"); err != nil || max != 100 {
		t.Errorf("MaxResults without IDs should pass: err=%v max=%d", err, max)
	}
}

func TestAllowedFiltersDocumentedNames(t *testing.T) {
	// Documented filter names must be accepted rather than rejected as
	// unknown, including alternative spellings and families backed by
	// attributes this implementation does not carry (they match nothing).
	cases := []struct {
		name    string
		allowed func(string) bool
	}{
		{"enable-dns64", allowedSubnetFilters},
		{"enable-lni-at-device-index", allowedSubnetFilters},
		{"outpost-arn", allowedSubnetFilters},
		{"customer-owned-ipv4-pool", allowedSubnetFilters},
		{"availabilityZone", allowedSubnetFilters},
		{"cidrBlock", allowedSubnetFilters},
		{"defaultForAz", allowedSubnetFilters},
		{"ipv6-cidr-block-association.state", allowedSubnetFilters},
		{"private-dns-name-options-on-launch.hostname-type", allowedSubnetFilters},
		{"ipv6-cidr-block-association.ipv6-pool", allowedVPCFilters},
		{"ipv6-cidr-block-association.state", allowedVPCFilters},
		{"ip-permission.cidr", allowedSGFilters},
		{"ip-permission.protocol", allowedSGFilters},
		{"egress.ip-permission.user-id", allowedSGFilters},
		{"ingress.ip-permission.ipv6-cidr", allowedSGFilters},
	}
	for _, c := range cases {
		if !c.allowed(c.name) {
			t.Errorf("filter %q should be allowed", c.name)
		}
	}
	for _, name := range []string{"no-such-filter", "ip-permission.nonsense"} {
		if allowedSGFilters(name) {
			t.Errorf("filter %q should be rejected", name)
		}
	}
	if allowedVPCFilters("ip-permission.cidr") {
		t.Error("ip-permission.cidr must not be allowed for DescribeVpcs")
	}
}

func TestHasFlatRuleParams(t *testing.T) {
	if hasFlatRuleParams(map[string]interface{}{}) {
		t.Error("empty params should report no flat rule")
	}
	if !hasFlatRuleParams(map[string]interface{}{"CidrIp": "0.0.0.0/0"}) {
		t.Error("CidrIp should be detected as a flat rule parameter")
	}
	if !hasFlatRuleParams(map[string]interface{}{"FromPort": "22"}) {
		t.Error("FromPort should be detected as a flat rule parameter")
	}
}

func TestMatchesSubnetFiltersDocumentedNames(t *testing.T) {
	sn := &ec2store.Subnet{
		SubnetId:         "subnet-1",
		AvailabilityZone: "us-east-1a",
		CidrBlock:        "10.0.1.0/24",
		EnableDns64:      true,
	}
	if !matchesSubnetFilters(sn, []ec2Filter{{Name: "enable-dns64", Values: []string{"true"}}}) {
		t.Error("enable-dns64=true should match a subnet with EnableDns64")
	}
	if matchesSubnetFilters(sn, []ec2Filter{{Name: "enable-dns64", Values: []string{"false"}}}) {
		t.Error("enable-dns64=false should not match a subnet with EnableDns64")
	}
	// Attributes not carried by the implementation never match.
	if matchesSubnetFilters(sn, []ec2Filter{{Name: "outpost-arn", Values: []string{"arn:aws:outposts:..."}}}) {
		t.Error("outpost-arn should never match")
	}
	if matchesSubnetFilters(sn, []ec2Filter{{Name: "ipv6-cidr-block-association.state", Values: []string{"associated"}}}) {
		t.Error("ipv6-cidr-block-association.* should never match")
	}
	// Alternative spellings behave like the canonical names.
	if !matchesSubnetFilters(sn, []ec2Filter{{Name: "availabilityZone", Values: []string{"us-east-1a"}}}) {
		t.Error("availabilityZone alias should match")
	}
	if !matchesSubnetFilters(sn, []ec2Filter{{Name: "cidrBlock", Values: []string{"10.0.1.0/24"}}}) {
		t.Error("cidrBlock alias should match")
	}
	if !matchesSubnetFilters(sn, []ec2Filter{{Name: "private-dns-name-options-on-launch.hostname-type", Values: []string{"ip-name"}}}) {
		t.Error("hostname-type=ip-name should match IPv4-only subnets")
	}
}

func TestMatchesSGFiltersPlainPermissionForm(t *testing.T) {
	sg := &ec2store.SecurityGroup{
		GroupId: "sg-1",
		IpPermissions: []ec2store.IPRule{{
			IpProtocol: "tcp", FromPort: 22, ToPort: 22,
			IpRanges: []ec2store.IPRange{{CidrIp: "0.0.0.0/0"}},
		}},
		IpPermissionsEgress: []ec2store.IPRule{{
			IpProtocol:       "-1",
			UserIdGroupPairs: []ec2store.GroupPair{{GroupId: "sg-peer", UserId: "111122223333"}},
		}},
	}
	// The plain ip-permission.* form (used by the AWS documentation
	// examples) filters inbound rules.
	if !matchesSGFilters(sg, []ec2Filter{{Name: "ip-permission.cidr", Values: []string{"0.0.0.0/0"}}}) {
		t.Error("ip-permission.cidr should match the inbound IPv4 range")
	}
	if matchesSGFilters(sg, []ec2Filter{{Name: "ip-permission.from-port", Values: []string{"443"}}}) {
		t.Error("ip-permission.from-port=443 should not match")
	}
	if !matchesSGFilters(sg, []ec2Filter{{Name: "egress.ip-permission.user-id", Values: []string{"111122223333"}}}) {
		t.Error("egress.ip-permission.user-id should match the referenced account")
	}
}

func TestPaginateEC2(t *testing.T) {
	items := []string{"a", "b", "c", "d", "e", "f"}
	key := func(s string) string { return s }

	// No MaxResults: all items returned.
	res, err := paginateEC2(items, "", 0, key)
	if err != nil || len(res.Items) != 6 {
		t.Errorf("no-limit pagination failed: %v, %d items", err, len(res.Items))
	}

	// Page of 2 with valid marker.
	res, err = paginateEC2(items, "b", 2, key)
	if err != nil || len(res.Items) != 2 || res.Items[0] != "c" || !res.IsTruncated {
		t.Errorf("marker pagination failed: %v %+v", err, res)
	}

	// Invalid marker with limit: the underlying helper returns an empty page.
	res, err = paginateEC2(items, "zzz", 2, key)
	if err == nil || len(res.Items) != 0 {
		t.Errorf("invalid marker should produce InvalidNextToken, got %v %+v", err, res)
	}

	// Invalid marker without limit.
	_, err = paginateEC2(items, "zzz", 0, key)
	if err == nil {
		t.Error("invalid marker without limit should produce InvalidNextToken")
	}
}

func TestPermissionExists(t *testing.T) {
	existing := []ec2store.IPRule{{
		IpProtocol: "tcp",
		FromPort:   443,
		ToPort:     443,
		IpRanges:   []ec2store.IPRange{{CidrIp: "10.0.0.0/8"}},
	}}
	same := ec2store.IPRule{IpProtocol: "tcp", FromPort: 443, ToPort: 443,
		IpRanges: []ec2store.IPRange{{CidrIp: "10.0.0.0/8"}}}
	otherCidr := ec2store.IPRule{IpProtocol: "tcp", FromPort: 443, ToPort: 443,
		IpRanges: []ec2store.IPRange{{CidrIp: "192.168.0.0/16"}}}
	wrongPorts := ec2store.IPRule{IpProtocol: "tcp", FromPort: 80, ToPort: 80,
		IpRanges: []ec2store.IPRange{{CidrIp: "10.0.0.0/8"}}}
	if !permissionExists(existing, same) {
		t.Error("identical permission should exist")
	}
	if permissionExists(existing, otherCidr) {
		t.Error("different CIDR should not exist")
	}
	if permissionExists(existing, wrongPorts) {
		t.Error("different ports should not exist")
	}
}
