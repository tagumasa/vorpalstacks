package wafv2

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/protobuf/proto"
	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
)

// ManagedRuleGroupSummary represents a high-level summary of an AWS-managed rule group.
type ManagedRuleGroupSummary struct {
	Name                *string `json:"Name,omitempty"`
	VendorName          *string `json:"VendorName,omitempty"`
	Description         *string `json:"Description,omitempty"`
	VersioningSupported *bool   `json:"VersioningSupported,omitempty"`
}

// managedRuleGroupDetail holds the WCU capacity, rule list, and label
// information for a managed rule group. The data is sourced from the
// official AWS WAF Managed Rules documentation.
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

// managedRuleGroupDetails provides per-group WCU, rule lists, and label
// information sourced from the AWS WAF Managed Rules documentation.
var managedRuleGroupDetails = map[string]managedRuleGroupDetail{
	"AWSManagedRulesCommonRuleSet": {
		WCU: 700,
		Rules: []managedRule{
			{Name: "NoUserAgent_HEADER", Action: "Block"},
			{Name: "UserAgent_BadBots_HEADER", Action: "Block"},
			{Name: "SizeRestrictions_QUERYSTRING", Action: "Block"},
			{Name: "SizeRestrictions_COOKIE_HEADER", Action: "Block"},
			{Name: "SizeRestrictions_BODY", Action: "Block"},
			{Name: "SizeRestrictions_URIPATH", Action: "Block"},
			{Name: "EC2MetaDataSSRF_BODY", Action: "Block"},
			{Name: "EC2MetaDataSSRF_COOKIE", Action: "Block"},
			{Name: "EC2MetaDataSSRF_URIPATH", Action: "Block"},
			{Name: "EC2MetaDataSSRF_QUERYARGUMENTS", Action: "Block"},
			{Name: "GenericLFI_QUERYARGUMENTS", Action: "Block"},
			{Name: "GenericLFI_URIPATH", Action: "Block"},
			{Name: "GenericLFI_BODY", Action: "Block"},
			{Name: "RestrictedExtensions_URIPATH", Action: "Block"},
			{Name: "RestrictedExtensions_QUERYARGUMENTS", Action: "Block"},
			{Name: "GenericRFI_QUERYARGUMENTS", Action: "Block"},
			{Name: "GenericRFI_BODY", Action: "Block"},
			{Name: "GenericRFI_URIPATH", Action: "Block"},
			{Name: "CrossSiteScripting_COOKIE", Action: "Block"},
			{Name: "CrossSiteScripting_QUERYARGUMENTS", Action: "Block"},
			{Name: "CrossSiteScripting_BODY", Action: "Block"},
			{Name: "CrossSiteScripting_URIPATH", Action: "Block"},
		},
		AvailableLabels: []string{
			"awswaf:managed:aws:core-rule-set:NoUserAgent_Header",
			"awswaf:managed:aws:core-rule-set:BadBots_Header",
			"awswaf:managed:aws:core-rule-set:SizeRestrictions_QueryString",
			"awswaf:managed:aws:core-rule-set:SizeRestrictions_Cookie_Header",
			"awswaf:managed:aws:core-rule-set:SizeRestrictions_Body",
			"awswaf:managed:aws:core-rule-set:SizeRestrictions_URIPath",
			"awswaf:managed:aws:core-rule-set:EC2MetaDataSSRF_Body",
			"awswaf:managed:aws:core-rule-set:EC2MetaDataSSRF_Cookie",
			"awswaf:managed:aws:core-rule-set:EC2MetaDataSSRF_URIPath",
			"awswaf:managed:aws:core-rule-set:EC2MetaDataSSRF_QueryArguments",
			"awswaf:managed:aws:core-rule-set:GenericLFI_QueryArguments",
			"awswaf:managed:aws:core-rule-set:GenericLFI_URIPath",
			"awswaf:managed:aws:core-rule-set:GenericLFI_Body",
			"awswaf:managed:aws:core-rule-set:RestrictedExtensions_URIPath",
			"awswaf:managed:aws:core-rule-set:RestrictedExtensions_QueryArguments",
			"awswaf:managed:aws:core-rule-set:GenericRFI_QueryArguments",
			"awswaf:managed:aws:core-rule-set:GenericRFI_Body",
			"awswaf:managed:aws:core-rule-set:GenericRFI_URIPath",
			"awswaf:managed:aws:core-rule-set:CrossSiteScripting_Cookie",
			"awswaf:managed:aws:core-rule-set:CrossSiteScripting_QueryArguments",
			"awswaf:managed:aws:core-rule-set:CrossSiteScripting_Body",
			"awswaf:managed:aws:core-rule-set:CrossSiteScripting_URIPath",
		},
	},
	"AWSManagedRulesAdminProtectionRuleSet": {
		WCU: 100,
		Rules: []managedRule{
			{Name: "AdminProtection_URIPATH", Action: "Block"},
		},
		AvailableLabels: []string{
			"awswaf:managed:aws:admin-protection:AdminProtection_URIPath",
		},
	},
	"AWSManagedRulesKnownBadInputsRuleSet": {
		WCU: 200,
		Rules: []managedRule{
			{Name: "JavaDeserializationRCE_HEADER", Action: "Block"},
			{Name: "JavaDeserializationRCE_BODY", Action: "Block"},
			{Name: "JavaDeserializationRCE_URIPATH", Action: "Block"},
			{Name: "JavaDeserializationRCE_QUERYSTRING", Action: "Block"},
			{Name: "Host_localhost_HEADER", Action: "Block"},
			{Name: "PROPFIND_METHOD", Action: "Block"},
			{Name: "ExploitablePaths_URIPATH", Action: "Block"},
			{Name: "Log4JRCE_HEADER", Action: "Block"},
			{Name: "Log4JRCE_QUERYSTRING", Action: "Block"},
			{Name: "Log4JRCE_BODY", Action: "Block"},
			{Name: "Log4JRCE_URIPATH", Action: "Block"},
			{Name: "ReactJSRCE_BODY", Action: "Block"},
		},
		AvailableLabels: []string{
			"awswaf:managed:aws:known-bad-inputs:JavaDeserializationRCE_Header",
			"awswaf:managed:aws:known-bad-inputs:JavaDeserializationRCE_Body",
			"awswaf:managed:aws:known-bad-inputs:JavaDeserializationRCE_URIPath",
			"awswaf:managed:aws:known-bad-inputs:JavaDeserializationRCE_QueryString",
			"awswaf:managed:aws:known-bad-inputs:Host_Localhost_Header",
			"awswaf:managed:aws:known-bad-inputs:Propfind_Method",
			"awswaf:managed:aws:known-bad-inputs:ExploitablePaths_URIPath",
			"awswaf:managed:aws:known-bad-inputs:Log4JRCE_Header",
			"awswaf:managed:aws:known-bad-inputs:Log4JRCE_QueryString",
			"awswaf:managed:aws:known-bad-inputs:Log4JRCE_Body",
			"awswaf:managed:aws:known-bad-inputs:Log4JRCE_URIPath",
			"awswaf:managed:aws:known-bad-inputs:ReactJSRCE_Body",
		},
	},
	"AWSManagedRulesSQLiRuleSet": {
		WCU: 200,
		Rules: []managedRule{
			{Name: "SQLi_QUERYARGUMENTS", Action: "Block"},
			{Name: "SQLiExtendedPatterns_QUERYARGUMENTS", Action: "Block"},
			{Name: "SQLi_BODY", Action: "Block"},
			{Name: "SQLiExtendedPatterns_BODY", Action: "Block"},
			{Name: "SQLiExtendedPatterns_HEADER", Action: "Block"},
			{Name: "SQLiExtendedPatterns_URIPATH", Action: "Block"},
			{Name: "SQLi_COOKIE", Action: "Block"},
			{Name: "SQLi_URIPATH", Action: "Block"},
		},
		AvailableLabels: []string{
			"awswaf:managed:aws:sql-database:SQLi_QueryArguments",
			"awswaf:managed:aws:sql-database:SQLiExtendedPatterns_QueryArguments",
			"awswaf:managed:aws:sql-database:SQLi_Body",
			"awswaf:managed:aws:sql-database:SQLiExtendedPatterns_Body",
			"awswaf:managed:aws:sql-database:SQLiExtendedPatterns_Header",
			"awswaf:managed:aws:sql-database:SQLiExtendedPatterns_UriPath",
			"awswaf:managed:aws:sql-database:SQLi_Cookie",
			"awswaf:managed:aws:sql-database:SQLi_URIPath",
		},
	},
	"AWSManagedRulesLinuxRuleSet": {
		WCU: 200,
		Rules: []managedRule{
			{Name: "LFI_URIPATH", Action: "Block"},
			{Name: "LFI_QUERYSTRING", Action: "Block"},
			{Name: "LFI_HEADER", Action: "Block"},
		},
		AvailableLabels: []string{
			"awswaf:managed:aws:linux-os:LFI_URIPath",
			"awswaf:managed:aws:linux-os:LFI_QueryString",
			"awswaf:managed:aws:linux-os:LFI_Header",
		},
	},
	"AWSManagedRulesUnixRuleSet": {
		WCU: 100,
		Rules: []managedRule{
			{Name: "UNIXShellCommandsVariables_QUERYSTRING", Action: "Block"},
			{Name: "UNIXShellCommandsVariables_BODY", Action: "Block"},
			{Name: "UNIXShellCommandsVariables_HEADER", Action: "Block"},
		},
		AvailableLabels: []string{
			"awswaf:managed:aws:posix-os:UNIXShellCommandsVariables_QueryString",
			"awswaf:managed:aws:posix-os:UNIXShellCommandsVariables_Body",
			"awswaf:managed:aws:posix-os:UNIXShellCommandsVariables_Header",
		},
	},
	"AWSManagedRulesWindowsRuleSet": {
		WCU: 200,
		Rules: []managedRule{
			{Name: "WindowsShellCommands_HEADER", Action: "Block"},
			{Name: "WindowsShellCommands_QUERYARGUMENTS", Action: "Block"},
			{Name: "WindowsShellCommands_QUERYSTRING", Action: "Block"},
			{Name: "WindowsShellCommands_URIPATH", Action: "Block"},
			{Name: "WindowsShellCommands_BODY", Action: "Block"},
			{Name: "PowerShellCommands_COOKIE", Action: "Block"},
			{Name: "PowerShellCommands_QUERYARGUMENTS", Action: "Block"},
			{Name: "PowerShellCommands_BODY", Action: "Block"},
		},
		AvailableLabels: []string{
			"awswaf:managed:aws:windows-os:WindowsShellCommands_Header",
			"awswaf:managed:aws:windows-os:WindowsShellCommands_QueryArguments",
			"awswaf:managed:aws:windows-os:WindowsShellCommands_QueryString",
			"awswaf:managed:aws:windows-os:WindowsShellCommands_UriPath",
			"awswaf:managed:aws:windows-os:WindowsShellCommands_Body",
			"awswaf:managed:aws:windows-os:PowerShellCommands_Cookie",
			"awswaf:managed:aws:windows-os:PowerShellCommands_QueryArguments",
			"awswaf:managed:aws:windows-os:PowerShellCommands_Body",
		},
	},
	"AWSManagedRulesPHPRuleSet": {
		WCU: 100,
		Rules: []managedRule{
			{Name: "PHPHighRiskMethodsVariables_HEADER", Action: "Block"},
			{Name: "PHPHighRiskMethodsVariables_QUERYSTRING", Action: "Block"},
			{Name: "PHPHighRiskMethodsVariables_BODY", Action: "Block"},
			{Name: "PHPHighRiskMethodsVariables_URIPATH", Action: "Block"},
		},
		AvailableLabels: []string{
			"awswaf:managed:aws:php-app:PHPHighRiskMethodsVariables_Header",
			"awswaf:managed:aws:php-app:PHPHighRiskMethodsVariables_QueryString",
			"awswaf:managed:aws:php-app:PHPHighRiskMethodsVariables_Body",
			"awswaf:managed:aws:php-app:PHPHighRiskMethodsVariables_URIPath",
		},
	},
	"AWSManagedRulesWordPressRuleSet": {
		WCU: 100,
		Rules: []managedRule{
			{Name: "WordPressExploitableCommands_QUERYSTRING", Action: "Block"},
			{Name: "WordPressExploitablePaths_URIPATH", Action: "Block"},
		},
		AvailableLabels: []string{
			"awswaf:managed:aws:wordpress-app:WordPressExploitableCommands_QUERYSTRING",
			"awswaf:managed:aws:wordpress-app:WordPressExploitablePaths_URIPATH",
		},
	},
	"AWSManagedRulesAmazonIpReputationList": {
		WCU: 25,
		Rules: []managedRule{
			{Name: "AWSManagedIPReputationList", Action: "Block"},
			{Name: "AWSManagedReconnaissanceList", Action: "Block"},
			{Name: "AWSManagedIPDDoSList", Action: "Count"},
		},
		AvailableLabels: []string{
			"awswaf:managed:aws:amazon-ip-list:AWSManagedIPReputationList",
			"awswaf:managed:aws:amazon-ip-list:AWSManagedReconnaissanceList",
			"awswaf:managed:aws:amazon-ip-list:AWSManagedIPDDoSList",
		},
	},
	"AWSManagedRulesAnonymousIpList": {
		WCU: 50,
		Rules: []managedRule{
			{Name: "AnonymousIPList", Action: "Block"},
			{Name: "HostingProviderIPList", Action: "Block"},
		},
		AvailableLabels: []string{
			"awswaf:managed:aws:anonymous-ip-list:AnonymousIPList",
			"awswaf:managed:aws:anonymous-ip-list:HostingProviderIPList",
		},
	},
	"AWSManagedRulesBotControlRuleSet": {
		WCU: 50,
		Rules: []managedRule{
			{Name: "CategoryAdvertising", Action: "Block"},
			{Name: "CategoryArchiver", Action: "Block"},
			{Name: "CategoryContentFetcher", Action: "Block"},
			{Name: "CategoryEmailClient", Action: "Block"},
			{Name: "CategoryHttpLibrary", Action: "Block"},
			{Name: "CategoryLinkChecker", Action: "Block"},
			{Name: "CategoryMiscellaneous", Action: "Block"},
			{Name: "CategoryMonitoring", Action: "Block"},
			{Name: "CategoryPagePreview", Action: "Block"},
			{Name: "CategoryScrapingFramework", Action: "Block"},
			{Name: "CategorySearchEngine", Action: "Block"},
			{Name: "CategorySecurity", Action: "Block"},
			{Name: "CategorySeo", Action: "Block"},
			{Name: "CategorySocialMedia", Action: "Block"},
			{Name: "CategoryWebhooks", Action: "Block"},
			{Name: "CategoryAI", Action: "Block"},
			{Name: "SignalAutomatedBrowser", Action: "Block"},
			{Name: "SignalKnownBotDataCenter", Action: "Block"},
			{Name: "SignalNonBrowserUserAgent", Action: "Block"},
			{Name: "TGT_VolumetricIpTokenAbsent", Action: "Challenge"},
			{Name: "TGT_TokenAbsent", Action: "Count"},
			{Name: "TGT_VolumetricSession", Action: "Captcha"},
			{Name: "TGT_VolumetricSessionMaximum", Action: "Block"},
			{Name: "TGT_SignalAutomatedBrowser", Action: "Captcha"},
			{Name: "TGT_SignalBrowserAutomationExtension", Action: "Captcha"},
			{Name: "TGT_SignalBrowserInconsistency", Action: "Captcha"},
			{Name: "TGT_ML_CoordinatedActivityLow", Action: "Challenge"},
			{Name: "TGT_ML_CoordinatedActivityMedium", Action: "Captcha"},
			{Name: "TGT_ML_CoordinatedActivityHigh", Action: "Captcha"},
			{Name: "TGT_TokenReuseIpLow", Action: "Count"},
			{Name: "TGT_TokenReuseIpMedium", Action: "Captcha"},
			{Name: "TGT_TokenReuseIpHigh", Action: "Block"},
			{Name: "TGT_TokenReuseCountryLow", Action: "Count"},
			{Name: "TGT_TokenReuseCountryMedium", Action: "Captcha"},
			{Name: "TGT_TokenReuseCountryHigh", Action: "Block"},
			{Name: "TGT_TokenReuseAsnLow", Action: "Count"},
			{Name: "TGT_TokenReuseAsnMedium", Action: "Captcha"},
			{Name: "TGT_TokenReuseAsnHigh", Action: "Block"},
		},
		AvailableLabels: []string{
			"awswaf:managed:aws:bot-control:CategoryAdvertising",
			"awswaf:managed:aws:bot-control:CategorySearchEngine",
			"awswaf:managed:aws:bot-control:bot:verified",
			"awswaf:managed:aws:bot-control:SignalAutomatedBrowser",
			"awswaf:managed:aws:bot-control:TGT_VolumetricIpTokenAbsent",
			"awswaf:managed:aws:bot-control:TGT_VolumetricSession",
			"awswaf:managed:aws:bot-control:TGT_VolumetricSessionMaximum",
		},
	},
	"AWSManagedRulesATPRuleSet": {
		WCU: 50,
		Rules: []managedRule{
			{Name: "UnsupportedCognitoIDP", Action: "Block"},
			{Name: "VolumetricIpHigh", Action: "Block"},
			{Name: "VolumetricSession", Action: "Block"},
			{Name: "AttributeCompromisedCredentials", Action: "Block"},
			{Name: "AttributeUsernameTraversal", Action: "Block"},
			{Name: "AttributePasswordTraversal", Action: "Block"},
			{Name: "AttributeLongSession", Action: "Block"},
			{Name: "TokenRejected", Action: "Block"},
			{Name: "SignalMissingCredential", Action: "Block"},
			{Name: "VolumetricIpFailedLoginResponseHigh", Action: "Block"},
			{Name: "VolumetricSessionFailedLoginResponseHigh", Action: "Block"},
		},
		AvailableLabels: []string{
			"awswaf:managed:aws:atp:VolumetricIpHigh",
			"awswaf:managed:aws:atp:VolumetricSession",
			"awswaf:managed:aws:atp:AttributeCompromisedCredentials",
			"awswaf:managed:aws:atp:AttributeUsernameTraversal",
			"awswaf:managed:aws:atp:AttributePasswordTraversal",
			"awswaf:managed:aws:atp:AttributeLongSession",
			"awswaf:managed:aws:atp:SignalMissingCredential",
			"awswaf:managed:aws:atp:VolumetricIpFailedLoginResponseHigh",
			"awswaf:managed:aws:atp:VolumetricSessionFailedLoginResponseHigh",
			"awswaf:managed:aws:atp:signal:credential_compromised",
		},
	},
	"AWSManagedRulesACFPRuleSet": {
		WCU: 50,
		Rules: []managedRule{
			{Name: "UnsupportedCognitoIDP", Action: "Block"},
			{Name: "AllRequests", Action: "Challenge"},
			{Name: "RiskScoreHigh", Action: "Block"},
			{Name: "SignalCredentialCompromised", Action: "Block"},
			{Name: "SignalClientHumanInteractivityAbsentLow", Action: "Captcha"},
			{Name: "AutomatedBrowser", Action: "Block"},
			{Name: "BrowserInconsistency", Action: "Captcha"},
			{Name: "VolumetricIpHigh", Action: "Captcha"},
			{Name: "VolumetricSessionHigh", Action: "Block"},
			{Name: "AttributeUsernameTraversalHigh", Action: "Block"},
			{Name: "VolumetricPhoneNumberHigh", Action: "Block"},
			{Name: "VolumetricAddressHigh", Action: "Block"},
			{Name: "VolumetricAddressLow", Action: "Captcha"},
			{Name: "VolumetricIPSuccessfulResponse", Action: "Block"},
			{Name: "VolumetricSessionSuccessfulResponse", Action: "Block"},
			{Name: "VolumetricSessionTokenReuseIp", Action: "Block"},
		},
		AvailableLabels: []string{
			"awswaf:managed:aws:acfp:RiskScoreHigh",
			"awswaf:managed:aws:acfp:SignalCredentialCompromised",
			"awswaf:managed:aws:acfp:AutomatedBrowser",
			"awswaf:managed:aws:acfp:VolumetricIpHigh",
			"awswaf:managed:aws:acfp:VolumetricSessionHigh",
			"awswaf:managed:aws:acfp:AttributeUsernameTraversalHigh",
			"awswaf:managed:aws:acfp:VolumetricPhoneNumberHigh",
			"awswaf:managed:aws:acfp:VolumetricAddressHigh",
			"awswaf:managed:aws:acfp:VolumetricIPSuccessfulResponse",
			"awswaf:managed:aws:acfp:VolumetricSessionSuccessfulResponse",
			"awswaf:managed:aws:acfp:VolumetricSessionTokenReuseIp",
			"awswaf:managed:aws:acfp:signal:credential_compromised",
		},
	},
	"AWSManagedRulesAntiDDoSRuleSet": {
		WCU: 50,
		Rules: []managedRule{
			{Name: "ChallengeAllDuringEvent", Action: "Challenge"},
			{Name: "ChallengeDDoSRequests", Action: "Challenge"},
			{Name: "DDoSRequests", Action: "Block"},
		},
		AvailableLabels: []string{
			"awswaf:managed:aws:anti-ddos:ChallengeAllDuringEvent",
			"awswaf:managed:aws:anti-ddos:ChallengeDDoSRequests",
			"awswaf:managed:aws:anti-ddos:DDoSRequests",
			"awswaf:managed:aws:anti-ddos:event-detected",
			"awswaf:managed:aws:anti-ddos:ddos-request",
		},
	},
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
// AWS-managed rule group, including accurate WCU capacity, rule names,
// actions, and available labels sourced from the AWS WAF Managed Rules
// documentation.
func (s *WAFv2Service) DescribeManagedRuleGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetStringParam(req.Parameters, "Name")
	vendorName := request.GetStringParam(req.Parameters, "VendorName")
	versionName := request.GetStringParam(req.Parameters, "VersionName")

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

	detail, ok := managedRuleGroupDetails[name]
	if !ok {
		return nil, notFoundError("ManagedRuleGroup")
	}

	labelNamespace := "awswaf:" + vendorName + ":" + name + ":"

	rules := make([]map[string]interface{}, 0, len(detail.Rules))
	for _, r := range detail.Rules {
		action := strings.ToLower(r.Action)
		rules = append(rules, map[string]interface{}{
			"Name":   r.Name,
			"Action": map[string]interface{}{action: map[string]interface{}{}},
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
		"LabelNamespace":  labelNamespace,
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
