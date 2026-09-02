package wafv2

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/protobuf/proto"
	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/services/aws/wafv2/inspection"
)

// ManagedRuleGroupSummary represents a high-level summary of an AWS-managed rule group.
type ManagedRuleGroupSummary struct {
	Name                *string `json:"Name,omitempty"`
	VendorName          *string `json:"VendorName,omitempty"`
	Description         *string `json:"Description,omitempty"`
	VersioningSupported *bool   `json:"VersioningSupported,omitempty"`
}

// managedRuleGroupDetail holds the WCU capacity, rule list, and label
// information for a managed rule group. The authoritative source is the
// inspection package's managed rules catalog, which also drives
// evaluation; the rule listings, label namespaces and capacities are
// reproduced from the AWS WAF Managed Rules documentation.
type managedRuleGroupDetail struct {
	WCU             int64
	Rules           []managedRule
	AvailableLabels []string
	ConsumedLabels  []string
}

type managedRule struct {
	Name   string
	Action string
}

// managedRuleGroupDetailFromCatalog projects the catalog entry into the
// shape the DescribeManagedRuleGroup response builds on. Every rule of
// the group contributes its label; the group's extra labels — the ones
// the documentation attributes to the group rather than to a rule —
// follow them.
func managedRuleGroupDetailFromCatalog(group *inspection.ManagedRuleGroup) managedRuleGroupDetail {
	detail := managedRuleGroupDetail{WCU: group.WCU}
	for _, rule := range group.Rules {
		detail.Rules = append(detail.Rules, managedRule{
			Name:   rule.Name,
			Action: rule.Action,
		})
		if rule.Label != "" {
			detail.AvailableLabels = append(detail.AvailableLabels, rule.Label)
		}
	}
	detail.AvailableLabels = append(detail.AvailableLabels, group.ExtraLabels...)
	return detail
}

// awsManagedRuleGroups is the authoritative list of AWS-managed rule
// groups, sourced from the official AWS WAF Managed Rules documentation
// (https://docs.aws.amazon.com/waf/latest/developerguide/aws-managed-rule-groups-list.html).
var awsManagedRuleGroups = []ManagedRuleGroupSummary{
	{Name: proto.String("AWSManagedRulesCommonRuleSet"), VendorName: proto.String("AWS"), Description: proto.String("Core rule set providing baseline protection against common web application vulnerabilities (OWASP Top 10)."), VersioningSupported: proto.Bool(true)},
	{Name: proto.String("AWSManagedRulesAdminProtectionRuleSet"), VendorName: proto.String("AWS"), Description: proto.String("Blocks external access to exposed administrative pages."), VersioningSupported: proto.Bool(true)},
	{Name: proto.String("AWSManagedRulesKnownBadInputsRuleSet"), VendorName: proto.String("AWS"), Description: proto.String("Blocks request patterns known to be invalid and associated with exploitation or discovery of vulnerabilities."), VersioningSupported: proto.Bool(true)},
	{Name: proto.String("AWSManagedRulesSQLiRuleSet"), VendorName: proto.String("AWS"), Description: proto.String("Blocks request patterns associated with SQL injection attacks."), VersioningSupported: proto.Bool(true)},
	{Name: proto.String("AWSManagedRulesLinuxRuleSet"), VendorName: proto.String("AWS"), Description: proto.String("Blocks exploitation of vulnerabilities specific to Linux, including Linux-specific LFI attacks."), VersioningSupported: proto.Bool(true)},
	{Name: proto.String("AWSManagedRulesUnixRuleSet"), VendorName: proto.String("AWS"), Description: proto.String("Blocks exploitation of vulnerabilities specific to POSIX and POSIX-like operating systems."), VersioningSupported: proto.Bool(true)},
	{Name: proto.String("AWSManagedRulesWindowsRuleSet"), VendorName: proto.String("AWS"), Description: proto.String("Blocks exploitation of vulnerabilities specific to Windows, including PowerShell and shell command injection."), VersioningSupported: proto.Bool(true)},
	{Name: proto.String("AWSManagedRulesPHPRuleSet"), VendorName: proto.String("AWS"), Description: proto.String("Blocks exploitation of vulnerabilities specific to the PHP programming language."), VersioningSupported: proto.Bool(true)},
	{Name: proto.String("AWSManagedRulesWordPressRuleSet"), VendorName: proto.String("AWS"), Description: proto.String("Blocks exploitation of vulnerabilities specific to WordPress sites."), VersioningSupported: proto.Bool(true)},
	{Name: proto.String("AWSManagedRulesAmazonIpReputationList"), VendorName: proto.String("AWS"), Description: proto.String("Blocks IP addresses identified by Amazon internal threat intelligence as associated with bots or malicious activity."), VersioningSupported: proto.Bool(false)},
	{Name: proto.String("AWSManagedRulesAnonymousIpList"), VendorName: proto.String("AWS"), Description: proto.String("Blocks requests from services that permit obfuscation of viewer identity, including VPNs, proxies, and Tor nodes."), VersioningSupported: proto.Bool(false)},
	{Name: proto.String("AWSManagedRulesBotControlRuleSet"), VendorName: proto.String("AWS"), Description: proto.String("Detects and manages bot traffic with common and targeted protection levels."), VersioningSupported: proto.Bool(true)},
	{Name: proto.String("AWSManagedRulesATPRuleSet"), VendorName: proto.String("AWS"), Description: proto.String("AWS WAF Fraud Control account takeover prevention. Inspects login attempts for stolen credentials and anomalous patterns."), VersioningSupported: proto.Bool(true)},
	{Name: proto.String("AWSManagedRulesACFPRuleSet"), VendorName: proto.String("AWS"), Description: proto.String("AWS WAF Fraud Control account creation fraud prevention. Inspects registration requests for fraudulent patterns."), VersioningSupported: proto.Bool(true)},
	{Name: proto.String("AWSManagedRulesAntiDDoSRuleSet"), VendorName: proto.String("AWS"), Description: proto.String("Detects and mitigates Layer 7 DDoS attacks with soft and hard mitigations."), VersioningSupported: proto.Bool(true)},
}

// managedRuleGroupVersions provides known version identifiers for
// versioned managed rule groups. Sourced from the AWS Managed Rules
// changelog.
var managedRuleGroupVersions = map[string][]string{
	"AWSManagedRulesCommonRuleSet":          {"Version_2025_09_01", "Version_2024_06_01", "Version_2023_10_01"},
	"AWSManagedRulesAdminProtectionRuleSet": {"Version_2024_06_01", "Version_2023_10_01"},
	"AWSManagedRulesKnownBadInputsRuleSet":  {"Version_2025_09_01", "Version_2024_06_01", "Version_2023_10_01"},
	"AWSManagedRulesSQLiRuleSet":            {"Version_2025_09_01", "Version_2024_06_01", "Version_2023_10_01"},
	"AWSManagedRulesLinuxRuleSet":           {"Version_2024_06_01", "Version_2023_10_01"},
	"AWSManagedRulesUnixRuleSet":            {"Version_2024_06_01", "Version_2023_10_01"},
	"AWSManagedRulesWindowsRuleSet":         {"Version_2024_06_01", "Version_2023_10_01"},
	"AWSManagedRulesPHPRuleSet":             {"Version_2024_06_01", "Version_2023_10_01"},
	"AWSManagedRulesWordPressRuleSet":       {"Version_2024_06_01", "Version_2023_10_01"},
	"AWSManagedRulesBotControlRuleSet":      {"Version_2025_11_01", "Version_2025_06_01", "Version_2024_06_01"},
	"AWSManagedRulesATPRuleSet":             {"Version_2025_11_01", "Version_2025_06_01", "Version_2024_06_01"},
	"AWSManagedRulesACFPRuleSet":            {"Version_2025_11_01", "Version_2025_06_01", "Version_2024_06_01"},
	"AWSManagedRulesAntiDDoSRuleSet":        {"Version_2025_09_01", "Version_2025_06_01"},
}

// ListAvailableManagedRuleGroups returns a paginated list of all
// available AWS-managed rule groups.
func (s *WAFv2Service) ListAvailableManagedRuleGroups(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	limit := pagination.GetMaxItems(req.Parameters, 100, "Limit")
	nextMarker := pagination.GetMarker(req.Parameters, "NextMarker")

	startIdx := 0
	if nextMarker != "" {
		for i, rg := range awsManagedRuleGroups {
			if rg.Name != nil && *rg.Name == nextMarker {
				startIdx = i
				break
			}
		}
	}

	endIdx := startIdx + limit
	var respNextMarker *string
	if endIdx < len(awsManagedRuleGroups) {
		respNextMarker = awsManagedRuleGroups[endIdx].Name
	} else {
		endIdx = len(awsManagedRuleGroups)
	}

	resp := map[string]interface{}{
		"ManagedRuleGroups": awsManagedRuleGroups[startIdx:endIdx],
	}
	var nextMarkerVal string
	if respNextMarker != nil {
		nextMarkerVal = *respNextMarker
	}
	pagination.SetNextToken(resp, "NextMarker", nextMarkerVal)
	return resp, nil
}

// DescribeManagedRuleGroup provides details about the specified
// AWS-managed rule group: the capacity, label namespace, rule list and
// available labels, all sourced from the managed rules catalog. The
// response carries no statement bodies — DescribeManagedRuleGroup
// documents the rules by name, priority and action only.
func (s *WAFv2Service) DescribeManagedRuleGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetStringParam(req.Parameters, "Name")
	vendorName := request.GetStringParam(req.Parameters, "VendorName")
	versionName := request.GetStringParam(req.Parameters, "VersionName")

	if name == "" || vendorName == "" {
		return nil, invalidParamError("Name and VendorName are required")
	}

	group, ok := inspection.LookupManagedRuleGroup(vendorName, name)
	if !ok {
		return nil, notFoundError("ManagedRuleGroup")
	}
	detail := managedRuleGroupDetailFromCatalog(group)

	rules := make([]map[string]interface{}, 0, len(detail.Rules))
	for _, r := range detail.Rules {
		// RuleSummary carries the rule's name and action only; the
		// action union member name is the PascalCase action.
		rules = append(rules, map[string]interface{}{
			"Name":   r.Name,
			"Action": map[string]interface{}{r.Action: map[string]interface{}{}},
		})
	}

	availableLabels := make([]map[string]interface{}, 0, len(detail.AvailableLabels))
	for _, l := range detail.AvailableLabels {
		availableLabels = append(availableLabels, map[string]interface{}{"Name": l})
	}

	consumedLabels := make([]map[string]interface{}, 0, len(detail.ConsumedLabels))
	for _, l := range detail.ConsumedLabels {
		consumedLabels = append(consumedLabels, map[string]interface{}{"Name": l})
	}

	var versionNameResp *string
	if versionName != "" {
		versionNameResp = &versionName
	}

	return map[string]interface{}{
		"Capacity":        detail.WCU,
		"LabelNamespace":  group.Namespace + ":",
		"VersionName":     versionNameResp,
		"SnsTopicArn":     fmt.Sprintf("arn:aws:sns:us-east-1:123456789012:aws-managed-waf-%s", strings.ToLower(name)),
		"Rules":           rules,
		"AvailableLabels": availableLabels,
		"ConsumedLabels":  consumedLabels,
	}, nil
}

// ListAvailableManagedRuleGroupVersions returns the available versions
// for the specified AWS-managed rule group. IP reputation lists do not
// support versioning.
func (s *WAFv2Service) ListAvailableManagedRuleGroupVersions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetStringParam(req.Parameters, "Name")
	vendorName := request.GetStringParam(req.Parameters, "VendorName")

	if name == "" || vendorName == "" {
		return nil, invalidParamError("Name and VendorName are required")
	}

	found := false
	for _, rg := range awsManagedRuleGroups {
		if rg.Name != nil && rg.VendorName != nil && *rg.Name == name && *rg.VendorName == vendorName {
			found = true
			break
		}
	}

	if !found {
		return nil, notFoundError("ManagedRuleGroup")
	}

	versions, ok := managedRuleGroupVersions[name]
	if !ok {
		return map[string]interface{}{
			"Versions": []interface{}{},
		}, nil
	}

	versionList := make([]map[string]interface{}, 0, len(versions))
	for _, v := range versions {
		versionList = append(versionList, map[string]interface{}{
			"Name": v,
		})
	}

	return map[string]interface{}{
		"Versions": versionList,
	}, nil
}
