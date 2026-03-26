package wafv2

import (
	"context"
	"time"

	"vorpalstacks/internal/services/aws/common/request"
)

// ManagedRuleGroupSummary represents a high-level summary of an AWS-managed rule group.
type ManagedRuleGroupSummary struct {
	Name                *string `json:"Name,omitempty"`
	VendorName          *string `json:"VendorName,omitempty"`
	Description         *string `json:"Description,omitempty"`
	VersioningSupported *bool   `json:"VersioningSupported,omitempty"`
}

// awsManagedRuleGroups contains the list of available AWS-managed rule groups.
var awsManagedRuleGroups = []ManagedRuleGroupSummary{
	{Name: strPtr("AWSManagedRulesCommonRuleSet"), VendorName: strPtr("AWS"), Description: strPtr("Core rule set containing rules generally applicable to web applications."), VersioningSupported: boolPtr(true)},
	{Name: strPtr("AWSManagedRulesAmazonIpReputationList"), VendorName: strPtr("AWS"), Description: strPtr("Rules based on Amazon threat intelligence."), VersioningSupported: boolPtr(true)},
	{Name: strPtr("AWSManagedRulesKnownBadInputsRuleSet"), VendorName: strPtr("AWS"), Description: strPtr("Rules to block request patterns that are known to be invalid."), VersioningSupported: boolPtr(true)},
	{Name: strPtr("AWSManagedRulesSQLiRuleSet"), VendorName: strPtr("AWS"), Description: strPtr("Rules to block SQL injection attacks."), VersioningSupported: boolPtr(true)},
	{Name: strPtr("AWSManagedRulesLinuxRuleSet"), VendorName: strPtr("AWS"), Description: strPtr("Rules to block exploitation of Linux specific vulnerabilities."), VersioningSupported: boolPtr(true)},
	{Name: strPtr("AWSManagedRulesUnixRuleSet"), VendorName: strPtr("AWS"), Description: strPtr("Rules to block exploitation of Unix specific vulnerabilities."), VersioningSupported: boolPtr(true)},
	{Name: strPtr("AWSManagedRulesWindowsRuleSet"), VendorName: strPtr("AWS"), Description: strPtr("Rules to block exploitation of Windows specific vulnerabilities."), VersioningSupported: boolPtr(true)},
	{Name: strPtr("AWSManagedRulesPHPRuleSet"), VendorName: strPtr("AWS"), Description: strPtr("Rules to block exploitation of PHP specific vulnerabilities."), VersioningSupported: boolPtr(true)},
	{Name: strPtr("AWSManagedRulesWordPressRuleSet"), VendorName: strPtr("AWS"), Description: strPtr("Rules to block exploitation of WordPress specific vulnerabilities."), VersioningSupported: boolPtr(true)},
	{Name: strPtr("AWSManagedRulesBotControlRuleSet"), VendorName: strPtr("AWS"), Description: strPtr("Rules to detect and block bot traffic."), VersioningSupported: boolPtr(true)},
	{Name: strPtr("AWSManagedRulesATPRuleSet"), VendorName: strPtr("AWS"), Description: strPtr("Account Takeover Protection rules."), VersioningSupported: boolPtr(true)},
}

// ListAvailableManagedRuleGroups returns a paginated list of all available AWS-managed rule groups.
func (s *WAFv2Service) ListAvailableManagedRuleGroups(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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

	resp := map[string]interface{}{
		"ManagedRuleGroups": awsManagedRuleGroups[startIdx:endIdx],
	}
	if respNextMarker != nil {
		resp["NextMarker"] = respNextMarker
	}
	return resp, nil
}

// DescribeManagedRuleGroup provides details about the specified AWS-managed rule group, including capacity and available labels.
func (s *WAFv2Service) DescribeManagedRuleGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetStringParam(req.Parameters, "Name")
	vendorName := request.GetStringParam(req.Parameters, "VendorName")
	versionName := request.GetStringParam(req.Parameters, "VersionName")

	if name == "" || vendorName == "" {
		return nil, validationError("Name and VendorName are required")
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

	return nil, notFoundError("ManagedRuleGroup")
}

// ListAvailableManagedRuleGroupVersions returns the available versions for the specified AWS-managed rule group.
func (s *WAFv2Service) ListAvailableManagedRuleGroupVersions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetStringParam(req.Parameters, "Name")
	vendorName := request.GetStringParam(req.Parameters, "VendorName")

	if name == "" || vendorName == "" {
		return nil, validationError("Name and VendorName are required")
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
