package ec2

import (
	"testing"

	ec2store "vorpalstacks/internal/store/aws/ec2"
)

// TestRemoveByRuleIDsIPv6 guard-tests the IPv6 branch: entries must be read
// from Ipv6Ranges and reported with cidrIpv6, and the kept IPv4 ranges must
// be unaffected. A previous copy of this logic iterated IpRanges instead,
// corrupting the persisted rule set.
func TestRemoveByRuleIDsIPv6(t *testing.T) {
	rule := ec2store.IPRule{
		IpProtocol: "tcp",
		FromPort:   443,
		ToPort:     443,
		IpRanges:   []ec2store.IPRange{{CidrIp: "203.0.113.0/24", RuleId: "sgr-v4"}},
		Ipv6Ranges: []ec2store.IPRange{{CidrIp: "2001:db8::/64", RuleId: "sgr-v6"}},
	}

	keep, revoked := removeByRuleIDs(rule, []string{"sgr-v6"}, nil, false, "sg-1")

	if len(revoked) != 1 {
		t.Fatalf("revoked count = %d, want 1", len(revoked))
	}
	if revoked[0].RuleId != "sgr-v6" {
		t.Errorf("revoked rule ID = %q, want sgr-v6", revoked[0].RuleId)
	}
	if revoked[0].CidrIpv6 != "2001:db8::/64" {
		t.Errorf("revoked cidrIpv6 = %q, want 2001:db8::/64", revoked[0].CidrIpv6)
	}
	if revoked[0].CidrIpv4 != "" {
		t.Errorf("revoked cidrIpv4 = %q, want empty", revoked[0].CidrIpv4)
	}
	if len(keep.Ipv6Ranges) != 0 {
		t.Errorf("kept Ipv6Ranges = %v, want empty", keep.Ipv6Ranges)
	}
	if len(keep.IpRanges) != 1 || keep.IpRanges[0].CidrIp != "203.0.113.0/24" {
		t.Errorf("kept IpRanges = %v, want the untouched IPv4 range", keep.IpRanges)
	}
}

// TestRemoveByRuleIDsIPv4 verifies the IPv4 branch reports cidrIpv4 and does
// not duplicate the entry into the IPv6 family.
func TestRemoveByRuleIDsIPv4(t *testing.T) {
	rule := ec2store.IPRule{
		IpProtocol: "tcp",
		FromPort:   22,
		ToPort:     22,
		IpRanges:   []ec2store.IPRange{{CidrIp: "203.0.113.4/32", RuleId: "sgr-v4"}},
	}

	keep, revoked := removeByRuleIDs(rule, []string{"sgr-v4"}, nil, true, "sg-1")

	if len(revoked) != 1 {
		t.Fatalf("revoked count = %d, want 1", len(revoked))
	}
	if revoked[0].CidrIpv4 != "203.0.113.4/32" {
		t.Errorf("revoked cidrIpv4 = %q, want 203.0.113.4/32", revoked[0].CidrIpv4)
	}
	if revoked[0].CidrIpv6 != "" {
		t.Errorf("revoked cidrIpv6 = %q, want empty (no phantom IPv6 entry)", revoked[0].CidrIpv6)
	}
	if revoked[0].IsEgress != true {
		t.Errorf("revoked isEgress = %v, want true", revoked[0].IsEgress)
	}
	if len(keep.IpRanges) != 0 {
		t.Errorf("kept IpRanges = %v, want empty", keep.IpRanges)
	}
	if len(keep.Ipv6Ranges) != 0 {
		t.Errorf("kept Ipv6Ranges = %v, want empty", keep.Ipv6Ranges)
	}
}

// TestRemoveByRuleIDsMultipleEntriesWithinOneRule covers the entry-count
// bookkeeping that previously cancelled itself: a rule with one IPv4 and one
// IPv6 entry, removing both IDs must remove 2 entries, not 0.
func TestRemoveByRuleIDsMultipleEntriesWithinOneRule(t *testing.T) {
	rule := ec2store.IPRule{
		IpProtocol: "tcp",
		FromPort:   80,
		ToPort:     80,
		IpRanges:   []ec2store.IPRange{{CidrIp: "192.0.2.0/24", RuleId: "sgr-a"}},
		Ipv6Ranges: []ec2store.IPRange{{CidrIp: "2001:db8:1::/64", RuleId: "sgr-b"}},
	}

	keep, revoked := removeByRuleIDs(rule, []string{"sgr-a", "sgr-b"}, nil, false, "sg-1")

	if len(revoked) != 2 {
		t.Fatalf("revoked count = %d, want 2", len(revoked))
	}
	removed := ruleEntryCount(rule) - ruleEntryCount(keep)
	if removed != 2 {
		t.Errorf("removed entry count = %d, want 2 (must not cancel out)", removed)
	}
	for _, r := range revoked {
		switch r.RuleId {
		case "sgr-a":
			if r.CidrIpv4 != "192.0.2.0/24" || r.CidrIpv6 != "" {
				t.Errorf("sgr-a family mismatch: v4=%q v6=%q", r.CidrIpv4, r.CidrIpv6)
			}
		case "sgr-b":
			if r.CidrIpv6 != "2001:db8:1::/64" || r.CidrIpv4 != "" {
				t.Errorf("sgr-b family mismatch: v4=%q v6=%q", r.CidrIpv4, r.CidrIpv6)
			}
		default:
			t.Errorf("unexpected revoked rule ID %q", r.RuleId)
		}
	}
}

// TestRemoveByRuleIDsGroupPairAndPrefixList covers the remaining source
// families to ensure the generic traversal reports the right reference.
func TestRemoveByRuleIDsGroupPairAndPrefixList(t *testing.T) {
	rule := ec2store.IPRule{
		IpProtocol: "-1",
		UserIdGroupPairs: []ec2store.GroupPair{{
			GroupId: "sg-peer", UserId: "111122223333", RuleId: "sgr-g",
		}},
		PrefixListIds: []ec2store.PrefixListId{{
			PrefixListId: "pl-1", RuleId: "sgr-p",
		}},
	}

	keep, revoked := removeByRuleIDs(rule, []string{"sgr-g", "sgr-p"}, nil, false, "sg-1")

	if len(revoked) != 2 {
		t.Fatalf("revoked count = %d, want 2", len(revoked))
	}
	for _, r := range revoked {
		switch r.RuleId {
		case "sgr-g":
			if r.ReferencedGroupId != "sg-peer" {
				t.Errorf("referenced group ID = %q, want sg-peer", r.ReferencedGroupId)
			}
		case "sgr-p":
			if r.PrefixListId != "pl-1" {
				t.Errorf("prefix list ID = %q, want pl-1", r.PrefixListId)
			}
		}
	}
	if len(keep.UserIdGroupPairs) != 0 || len(keep.PrefixListIds) != 0 {
		t.Error("kept entry families must be empty after full removal")
	}
}

// TestRemoveByRuleIDsNoMatch keeps everything when no ID matches.
func TestRemoveByRuleIDsNoMatch(t *testing.T) {
	rule := ec2store.IPRule{
		IpProtocol: "tcp",
		IpRanges:   []ec2store.IPRange{{CidrIp: "10.0.0.0/8", RuleId: "sgr-x"}},
	}
	keep, revoked := removeByRuleIDs(rule, []string{"sgr-other"}, nil, false, "sg-1")
	if len(revoked) != 0 {
		t.Errorf("revoked count = %d, want 0", len(revoked))
	}
	if len(keep.IpRanges) != 1 {
		t.Errorf("kept IpRanges = %d entries, want 1", len(keep.IpRanges))
	}
}
