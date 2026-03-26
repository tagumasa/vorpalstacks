package waf

import (
	"context"
	"time"

	"vorpalstacks/internal/services/aws/common/request"
)

// ManagedRuleGroupSummary represents a summary of a managed rule group.
type ManagedRuleGroupSummary struct {
	Name                *string `json:"Name,omitempty"`
	VendorName          *string `json:"VendorName,omitempty"`
	Description         *string `json:"Description,omitempty"`
	VersioningSupported *bool   `json:"VersioningSupported,omitempty"`
}

// RuleSummary represents a summary of a rule.
type RuleSummary struct {
	Name   *string `json:"Name,omitempty"`
	Action *string `json:"Action,omitempty"`
}

// LabelSummary represents a summary of a label.
type LabelSummary struct {
	Name *string `json:"Name,omitempty"`
}

// ManagedRuleGroupVersion represents a version of a managed rule group.
type ManagedRuleGroupVersion struct {
	Name        *string    `json:"Name,omitempty"`
	LastUpdated *time.Time `json:"LastUpdated,omitempty"`
}

// awsManagedRuleGroups contains the list of AWS-managed rule groups available for use.
var awsManagedRuleGroups = []ManagedRuleGroupSummary{
	{
		Name:                strPtr("AWSManagedRulesCommonRuleSet"),
		VendorName:          strPtr("AWS"),
		Description:         strPtr("Core rule set containing rules generally applicable to web applications. This rule set provides protection against exploitation of a wide range of vulnerabilities, including those described in OWASP publications."),
		VersioningSupported: boolPtr(true),
	},
	{
		Name:                strPtr("AWSManagedRulesAmazonIpReputationList"),
		VendorName:          strPtr("AWS"),
		Description:         strPtr("Rules based on Amazon threat intelligence. This rule set blocks requests from IP addresses that are associated with bots or other malicious threats."),
		VersioningSupported: boolPtr(true),
	},
	{
		Name:                strPtr("AWSManagedRulesKnownBadInputsRuleSet"),
		VendorName:          strPtr("AWS"),
		Description:         strPtr("Rules to block request patterns that are known to be invalid and are associated with exploitation or discovery of vulnerabilities. This can help reduce false positives when you need to block a specific attack but don't want to block all requests that contain a specific string."),
		VersioningSupported: boolPtr(true),
	},
	{
		Name:                strPtr("AWSManagedRulesSQLiRuleSet"),
		VendorName:          strPtr("AWS"),
		Description:         strPtr("Rules to block SQL injection attacks. This rule set provides protection against exploitation of SQL injection vulnerabilities."),
		VersioningSupported: boolPtr(true),
	},
	{
		Name:                strPtr("AWSManagedRulesLinuxRuleSet"),
		VendorName:          strPtr("AWS"),
		Description:         strPtr("Rules to block exploitation of Linux specific vulnerabilities. This rule set provides protection against exploitation of Linux operating system vulnerabilities."),
		VersioningSupported: boolPtr(true),
	},
	{
		Name:                strPtr("AWSManagedRulesUnixRuleSet"),
		VendorName:          strPtr("AWS"),
		Description:         strPtr("Rules to block exploitation of Unix specific vulnerabilities. This rule set provides protection against exploitation of Unix operating system vulnerabilities."),
		VersioningSupported: boolPtr(true),
	},
	{
		Name:                strPtr("AWSManagedRulesWindowsRuleSet"),
		VendorName:          strPtr("AWS"),
		Description:         strPtr("Rules to block exploitation of Windows specific vulnerabilities. This rule set provides protection against exploitation of Windows operating system vulnerabilities."),
		VersioningSupported: boolPtr(true),
	},
	{
		Name:                strPtr("AWSManagedRulesPHPRuleSet"),
		VendorName:          strPtr("AWS"),
		Description:         strPtr("Rules to block exploitation of PHP specific vulnerabilities. This rule set provides protection against exploitation of PHP vulnerabilities."),
		VersioningSupported: boolPtr(true),
	},
	{
		Name:                strPtr("AWSManagedRulesWordPressRuleSet"),
		VendorName:          strPtr("AWS"),
		Description:         strPtr("Rules to block exploitation of WordPress specific vulnerabilities. This rule set provides protection against exploitation of WordPress vulnerabilities."),
		VersioningSupported: boolPtr(true),
	},
	{
		Name:                strPtr("AWSManagedRulesBotControlRuleSet"),
		VendorName:          strPtr("AWS"),
		Description:         strPtr("Rules to detect and block bot traffic. This rule set provides protection against bots, scrapers, and scanners."),
		VersioningSupported: boolPtr(true),
	},
	{
		Name:                strPtr("AWSManagedRulesATPRuleSet"),
		VendorName:          strPtr("AWS"),
		Description:         strPtr("Account Takeover Protection rules. This rule set provides protection against account takeover attacks."),
		VersioningSupported: boolPtr(true),
	},
}

// ListAvailableManagedRuleGroups returns a list of available managed rule groups.
func (s *WAFService) ListAvailableManagedRuleGroups(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	scope := request.GetStringParam(req.Parameters, "Scope")
	_ = scope

	limit := 100
	if limitParam := request.GetIntParam(req.Parameters, "Limit"); limitParam > 0 {
		limit = limitParam
	}

	nextMarker := request.GetStringParam(req.Parameters, "NextMarker")

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

	return map[string]interface{}{
		"ManagedRuleGroups": awsManagedRuleGroups[startIdx:endIdx],
		"NextMarker":        respNextMarker,
	}, nil
}

// DescribeManagedRuleGroup returns the details of a managed rule group.
func (s *WAFService) DescribeManagedRuleGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetStringParam(req.Parameters, "Name")
	vendorName := request.GetStringParam(req.Parameters, "VendorName")
	versionName := request.GetStringParam(req.Parameters, "VersionName")

	if name == "" || vendorName == "" {
		return nil, NewValidationException("Name and VendorName are required")
	}

	for _, rg := range awsManagedRuleGroups {
		if rg.Name != nil && rg.VendorName != nil && *rg.Name == name && *rg.VendorName == vendorName {
			capacity := int64(100)
			if name == "AWSManagedRulesBotControlRuleSet" {
				capacity = 50
			} else if name == "AWSManagedRulesATPRuleSet" {
				capacity = 50
			}

			labelNamespace := "awswaf:" + vendorName + ":" + name + ":"

			var versionNameResp *string
			if versionName != "" {
				versionNameResp = &versionName
			}

			return map[string]interface{}{
				"Capacity":       capacity,
				"LabelNamespace": labelNamespace,
				"VersionName":    versionNameResp,
				"Rules": []map[string]interface{}{
					{"Name": "Rule1", "Action": map[string]interface{}{"Block": map[string]interface{}{}}},
					{"Name": "Rule2", "Action": map[string]interface{}{"Block": map[string]interface{}{}}},
				},
				"AvailableLabels": []map[string]interface{}{
					{"Name": labelNamespace + "label1"},
					{"Name": labelNamespace + "label2"},
				},
			}, nil
		}
	}

	return nil, NewResourceNotFoundException("ManagedRuleGroup not found: " + name)
}

// ListAvailableManagedRuleGroupVersions returns the versions of a managed rule group.
func (s *WAFService) ListAvailableManagedRuleGroupVersions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetStringParam(req.Parameters, "Name")
	vendorName := request.GetStringParam(req.Parameters, "VendorName")

	if name == "" || vendorName == "" {
		return nil, NewValidationException("Name and VendorName are required")
	}

	found := false
	for _, rg := range awsManagedRuleGroups {
		if rg.Name != nil && rg.VendorName != nil && *rg.Name == name && *rg.VendorName == vendorName {
			found = true
			break
		}
	}

	if !found {
		return nil, NewResourceNotFoundException("ManagedRuleGroup not found: " + name)
	}

	now := time.Now()
	versions := []map[string]interface{}{
		{"Name": "Version_2026_03_01", "LastUpdated": now.Unix()},
		{"Name": "Version_2026_02_01", "LastUpdated": now.Unix()},
		{"Name": "Version_2026_01_01", "LastUpdated": now.Unix()},
	}

	return map[string]interface{}{
		"Versions": versions,
	}, nil
}

func strPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}
