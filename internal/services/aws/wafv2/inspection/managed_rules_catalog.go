package inspection

import (
	"regexp"
	"strings"

	wafstore "vorpalstacks/internal/store/aws/waf"
)

// ManagedRule is one rule of an AWS-managed rule group as this platform
// models it. Statement is the local inspection equivalent, derived from
// the rule's published description in the AWS Managed Rules documentation
// (rule names, actions, labels and capacities are reproduced from the
// listings; AWS does not publish the exact match patterns, so each local
// statement covers the documented example patterns together with the
// canonical public signatures of the same threat class). A nil Statement
// marks a rule whose inspection input exists only inside AWS — threat
// intelligence feeds, device fingerprints, ML models or unpublished
// advisory patterns — which therefore never matches locally and is
// reported as unenforced.
type ManagedRule struct {
	Name      string
	Priority  int32
	Action    string
	Label     string
	Statement *wafstore.Statement
}

// ManagedRuleGroup is one AWS-managed rule group in the catalog. Namespace
// is the group's label namespace, the prefix of every label its rules add;
// ExtraLabels are the labels the documentation attributes to the group as
// a whole rather than to an individual rule.
type ManagedRuleGroup struct {
	Name        string
	Namespace   string
	WCU         int64
	Rules       []ManagedRule
	ExtraLabels []string
}

// managedRuleGroups is the platform's catalog of the AWS managed rule
// groups. Priorities follow the documented listing order, which is the
// group's rule evaluation order; AWS does not publish the numeric
// priorities the DescribeManagedRuleGroup API returns.
var managedRuleGroups = []*ManagedRuleGroup{
	{
		Name:      "AWSManagedRulesCommonRuleSet",
		Namespace: "awswaf:managed:aws:core-rule-set",
		WCU:       700,
		Rules: []ManagedRule{
			managedCatalogRule(0, "NoUserAgent_HEADER", "Block", "NoUserAgent_Header",
				notStatement(regexStatement(singleHeader("user-agent"), ttNone(), "^."))),
			managedCatalogRule(1, "UserAgent_BadBots_HEADER", "Block", "BadBots_Header",
				orContains(singleHeader("user-agent"), ttLower(),
					"nessus", "nmap", "sqlmap", "nikto", "masscan", "zgrab",
					"dirbuster", "wpscan", "hydra", "gobuster", "acunetix", "havij")),
			managedCatalogRule(2, "SizeRestrictions_QUERYSTRING", "Block", "SizeRestrictions_QueryString",
				sizeAbove(queryStringField(), 2048)),
			managedCatalogRule(3, "SizeRestrictions_COOKIE_HEADER", "Block", "SizeRestrictions_Cookie_Header",
				sizeAbove(singleHeader("cookie"), 10240)),
			managedCatalogRule(4, "SizeRestrictions_BODY", "Block", "SizeRestrictions_Body",
				sizeAbove(bodyField(), 8192)),
			managedCatalogRule(5, "SizeRestrictions_URIPATH", "Block", "SizeRestrictions_URIPath",
				sizeAbove(uriPathField(), 1024)),
			managedCatalogRule(6, "EC2MetaDataSSRF_BODY", "Block", "EC2MetaDataSSRF_Body",
				ec2MetaDataStatement(bodyField())),
			managedCatalogRule(7, "EC2MetaDataSSRF_COOKIE", "Block", "EC2MetaDataSSRF_Cookie",
				ec2MetaDataStatement(cookiesField())),
			managedCatalogRule(8, "EC2MetaDataSSRF_URIPATH", "Block", "EC2MetaDataSSRF_URIPath",
				ec2MetaDataStatement(uriPathField())),
			managedCatalogRule(9, "EC2MetaDataSSRF_QUERYARGUMENTS", "Block", "EC2MetaDataSSRF_QueryArguments",
				ec2MetaDataStatement(allQueryArgumentsField())),
			managedCatalogRule(10, "GenericLFI_QUERYARGUMENTS", "Block", "GenericLFI_QueryArguments",
				lfiStatement(allQueryArgumentsField())),
			managedCatalogRule(11, "GenericLFI_URIPATH", "Block", "GenericLFI_URIPath",
				lfiStatement(uriPathField())),
			managedCatalogRule(12, "GenericLFI_BODY", "Block", "GenericLFI_Body",
				lfiStatement(bodyField())),
			managedCatalogRule(13, "RestrictedExtensions_URIPATH", "Block", "RestrictedExtensions_URIPath",
				restrictedExtensionsStatement(uriPathField())),
			managedCatalogRule(14, "RestrictedExtensions_QUERYARGUMENTS", "Block", "RestrictedExtensions_QueryArguments",
				restrictedExtensionsStatement(allQueryArgumentsField())),
			managedCatalogRule(15, "GenericRFI_QUERYARGUMENTS", "Block", "GenericRFI_QueryArguments",
				rfiStatement(allQueryArgumentsField())),
			managedCatalogRule(16, "GenericRFI_BODY", "Block", "GenericRFI_Body",
				rfiStatement(bodyField())),
			managedCatalogRule(17, "GenericRFI_URIPATH", "Block", "GenericRFI_URIPath",
				rfiStatement(uriPathField())),
			managedCatalogRule(18, "CrossSiteScripting_COOKIE", "Block", "CrossSiteScripting_Cookie",
				xssStatement(cookiesField())),
			managedCatalogRule(19, "CrossSiteScripting_QUERYARGUMENTS", "Block", "CrossSiteScripting_QueryArguments",
				xssStatement(allQueryArgumentsField())),
			managedCatalogRule(20, "CrossSiteScripting_BODY", "Block", "CrossSiteScripting_Body",
				xssStatement(bodyField())),
			managedCatalogRule(21, "CrossSiteScripting_URIPATH", "Block", "CrossSiteScripting_URIPath",
				xssStatement(uriPathField())),
		},
	},
	{
		Name:      "AWSManagedRulesAdminProtectionRuleSet",
		Namespace: "awswaf:managed:aws:admin-protection",
		WCU:       100,
		Rules: []ManagedRule{
			managedCatalogRule(0, "AdminProtection_URIPATH", "Block", "AdminProtection_URIPath",
				orContains(uriPathField(), ttURLLower(),
					"sqlmanager", "phpmyadmin", "adminer", "manager/html")),
		},
	},
	{
		Name:      "AWSManagedRulesKnownBadInputsRuleSet",
		Namespace: "awswaf:managed:aws:known-bad-inputs",
		WCU:       200,
		Rules: []ManagedRule{
			managedCatalogRule(0, "JavaDeserializationRCE_HEADER", "Block", "JavaDeserializationRCE_Header",
				javaDeserializationStatement(headersField())),
			managedCatalogRule(1, "JavaDeserializationRCE_BODY", "Block", "JavaDeserializationRCE_Body",
				javaDeserializationStatement(bodyField())),
			managedCatalogRule(2, "JavaDeserializationRCE_URIPATH", "Block", "JavaDeserializationRCE_URIPath",
				javaDeserializationStatement(uriPathField())),
			managedCatalogRule(3, "JavaDeserializationRCE_QUERYSTRING", "Block", "JavaDeserializationRCE_QueryString",
				javaDeserializationStatement(queryStringField())),
			managedCatalogRule(4, "Host_localhost_HEADER", "Block", "Host_Localhost_Header",
				containsStatement(singleHeader("host"), ttLower(), "localhost")),
			managedCatalogRule(5, "PROPFIND_METHOD", "Block", "Propfind_Method",
				&wafstore.Statement{ByteMatchStatement: &wafstore.ByteMatchStatement{
					FieldToMatch:         methodField(),
					SearchString:         []byte("PROPFIND"),
					PositionalConstraint: "EXACTLY",
					TextTransformations:  ttNone(),
				}}),
			managedCatalogRule(6, "ExploitablePaths_URIPATH", "Block", "ExploitablePaths_URIPath",
				orContains(uriPathField(), ttURLLower(), "web-inf", "/.env", "/.git")),
			managedCatalogRule(7, "Log4JRCE_HEADER", "Block", "Log4JRCE_Header",
				log4jStatement(headersField())),
			managedCatalogRule(8, "Log4JRCE_QUERYSTRING", "Block", "Log4JRCE_QueryString",
				log4jStatement(queryStringField())),
			managedCatalogRule(9, "Log4JRCE_BODY", "Block", "Log4JRCE_Body",
				log4jStatement(bodyField())),
			managedCatalogRule(10, "Log4JRCE_URIPATH", "Block", "Log4JRCE_URIPath",
				log4jStatement(uriPathField())),
			// The advisory patterns for this rule are not published in the
			// AWS documentation, so no local statement can be derived.
			managedCatalogRule(11, "ReactJSRCE_BODY", "Block", "ReactJSRCE_Body", nil),
		},
	},
	{
		Name:      "AWSManagedRulesSQLiRuleSet",
		Namespace: "awswaf:managed:aws:sql-database",
		WCU:       200,
		Rules: []ManagedRule{
			managedCatalogRule(0, "SQLi_QUERYARGUMENTS", "Block", "SQLi_QueryArguments",
				sqliStatement(allQueryArgumentsField())),
			managedCatalogRule(1, "SQLiExtendedPatterns_QUERYARGUMENTS", "Block", "SQLiExtendedPatterns_QueryArguments",
				sqliExtendedStatement(allQueryArgumentsField())),
			managedCatalogRule(2, "SQLi_BODY", "Block", "SQLi_Body",
				sqliStatement(bodyField())),
			managedCatalogRule(3, "SQLiExtendedPatterns_BODY", "Block", "SQLiExtendedPatterns_Body",
				sqliExtendedStatement(bodyField())),
			managedCatalogRule(4, "SQLiExtendedPatterns_HEADER", "Block", "SQLiExtendedPatterns_Header",
				sqliExtendedStatement(headersField())),
			managedCatalogRule(5, "SQLiExtendedPatterns_URIPATH", "Block", "SQLiExtendedPatterns_UriPath",
				sqliExtendedStatement(uriPathField())),
			managedCatalogRule(6, "SQLi_COOKIE", "Block", "SQLi_Cookie",
				sqliStatement(cookiesField())),
			managedCatalogRule(7, "SQLi_URIPATH", "Block", "SQLi_URIPath",
				sqliStatement(uriPathField())),
		},
	},
	{
		Name:      "AWSManagedRulesLinuxRuleSet",
		Namespace: "awswaf:managed:aws:linux-os",
		WCU:       200,
		Rules: []ManagedRule{
			managedCatalogRule(0, "LFI_URIPATH", "Block", "LFI_URIPath",
				linuxLFIStatement(uriPathField())),
			managedCatalogRule(1, "LFI_QUERYSTRING", "Block", "LFI_QueryString",
				linuxLFIStatement(queryStringField())),
			managedCatalogRule(2, "LFI_HEADER", "Block", "LFI_Header",
				linuxLFIStatement(headersField())),
		},
	},
	{
		Name:      "AWSManagedRulesUnixRuleSet",
		Namespace: "awswaf:managed:aws:posix-os",
		WCU:       100,
		Rules: []ManagedRule{
			managedCatalogRule(0, "UNIXShellCommandsVariables_QUERYSTRING", "Block", "UNIXShellCommandsVariables_QueryString",
				unixShellStatement(queryStringField())),
			managedCatalogRule(1, "UNIXShellCommandsVariables_BODY", "Block", "UNIXShellCommandsVariables_Body",
				unixShellStatement(bodyField())),
			managedCatalogRule(2, "UNIXShellCommandsVariables_HEADER", "Block", "UNIXShellCommandsVariables_Header",
				unixShellStatement(headersField())),
		},
	},
	{
		Name:      "AWSManagedRulesWindowsRuleSet",
		Namespace: "awswaf:managed:aws:windows-os",
		WCU:       200,
		Rules: []ManagedRule{
			managedCatalogRule(0, "WindowsShellCommands_HEADER", "Block", "WindowsShellCommands_Header",
				windowsShellStatement(headersField())),
			managedCatalogRule(1, "WindowsShellCommands_QUERYARGUMENTS", "Block", "WindowsShellCommands_QueryArguments",
				windowsShellStatement(allQueryArgumentsField())),
			managedCatalogRule(2, "WindowsShellCommands_QUERYSTRING", "Block", "WindowsShellCommands_QueryString",
				windowsShellStatement(queryStringField())),
			managedCatalogRule(3, "WindowsShellCommands_URIPATH", "Block", "WindowsShellCommands_UriPath",
				windowsShellStatement(uriPathField())),
			managedCatalogRule(4, "WindowsShellCommands_BODY", "Block", "WindowsShellCommands_Body",
				windowsShellStatement(bodyField())),
			managedCatalogRule(5, "PowerShellCommands_COOKIE", "Block", "PowerShellCommands_Cookie",
				powerShellStatement(cookiesField())),
			managedCatalogRule(6, "PowerShellCommands_QUERYARGUMENTS", "Block", "PowerShellCommands_QueryArguments",
				powerShellStatement(allQueryArgumentsField())),
			managedCatalogRule(7, "PowerShellCommands_BODY", "Block", "PowerShellCommands_Body",
				powerShellStatement(bodyField())),
		},
	},
	{
		Name:      "AWSManagedRulesPHPRuleSet",
		Namespace: "awswaf:managed:aws:php-app",
		WCU:       100,
		Rules: []ManagedRule{
			managedCatalogRule(0, "PHPHighRiskMethodsVariables_HEADER", "Block", "PHPHighRiskMethodsVariables_Header",
				phpStatement(headersField())),
			managedCatalogRule(1, "PHPHighRiskMethodsVariables_QUERYSTRING", "Block", "PHPHighRiskMethodsVariables_QueryString",
				phpStatement(queryStringField())),
			managedCatalogRule(2, "PHPHighRiskMethodsVariables_BODY", "Block", "PHPHighRiskMethodsVariables_Body",
				phpStatement(bodyField())),
			managedCatalogRule(3, "PHPHighRiskMethodsVariables_URIPATH", "Block", "PHPHighRiskMethodsVariables_URIPath",
				phpStatement(uriPathField())),
		},
	},
	{
		Name:      "AWSManagedRulesWordPressRuleSet",
		Namespace: "awswaf:managed:aws:wordpress-app",
		WCU:       100,
		Rules: []ManagedRule{
			managedCatalogRule(0, "WordPressExploitableCommands_QUERYSTRING", "Block", "WordPressExploitableCommands_QUERYSTRING",
				containsStatement(queryStringField(), ttURLLower(), "do-reset-wordpress")),
			managedCatalogRule(1, "WordPressExploitablePaths_URIPATH", "Block", "WordPressExploitablePaths_URIPATH",
				orContains(uriPathField(), ttURLLower(), "xmlrpc.php", "wp-config.php")),
		},
	},
	// The rule groups below match on data that exists only inside AWS:
	// threat-intelligence IP feeds, anonymity-service and hosting-provider
	// ranges, device fingerprints, ML models and fraud telemetry. Their
	// rules are catalogued so the managed rule group APIs describe them
	// faithfully, but no local statement can be derived, so they never
	// match and each is reported as unenforced when the group runs.
	{
		Name:      "AWSManagedRulesAmazonIpReputationList",
		Namespace: "awswaf:managed:aws:amazon-ip-list",
		WCU:       25,
		Rules: []ManagedRule{
			managedCatalogRule(0, "AWSManagedIPReputationList", "Block", "AWSManagedIPReputationList", nil),
			managedCatalogRule(1, "AWSManagedReconnaissanceList", "Block", "AWSManagedReconnaissanceList", nil),
			managedCatalogRule(2, "AWSManagedIPDDoSList", "Count", "AWSManagedIPDDoSList", nil),
		},
	},
	{
		Name:      "AWSManagedRulesAnonymousIpList",
		Namespace: "awswaf:managed:aws:anonymous-ip-list",
		WCU:       50,
		Rules: []ManagedRule{
			managedCatalogRule(0, "AnonymousIPList", "Block", "AnonymousIPList", nil),
			managedCatalogRule(1, "HostingProviderIPList", "Block", "HostingProviderIPList", nil),
		},
	},
	{
		Name:      "AWSManagedRulesBotControlRuleSet",
		Namespace: "awswaf:managed:aws:bot-control",
		WCU:       50,
		ExtraLabels: []string{
			"awswaf:managed:aws:bot-control:bot:verified",
		},
		Rules: []ManagedRule{
			managedCatalogRule(0, "CategoryAdvertising", "Block", "CategoryAdvertising", nil),
			managedCatalogRule(1, "CategoryArchiver", "Block", "CategoryArchiver", nil),
			managedCatalogRule(2, "CategoryContentFetcher", "Block", "CategoryContentFetcher", nil),
			managedCatalogRule(3, "CategoryEmailClient", "Block", "CategoryEmailClient", nil),
			managedCatalogRule(4, "CategoryHttpLibrary", "Block", "CategoryHttpLibrary", nil),
			managedCatalogRule(5, "CategoryLinkChecker", "Block", "CategoryLinkChecker", nil),
			managedCatalogRule(6, "CategoryMiscellaneous", "Block", "CategoryMiscellaneous", nil),
			managedCatalogRule(7, "CategoryMonitoring", "Block", "CategoryMonitoring", nil),
			managedCatalogRule(8, "CategoryPagePreview", "Block", "CategoryPagePreview", nil),
			managedCatalogRule(9, "CategoryScrapingFramework", "Block", "CategoryScrapingFramework", nil),
			managedCatalogRule(10, "CategorySearchEngine", "Block", "CategorySearchEngine", nil),
			managedCatalogRule(11, "CategorySecurity", "Block", "CategorySecurity", nil),
			managedCatalogRule(12, "CategorySeo", "Block", "CategorySeo", nil),
			managedCatalogRule(13, "CategorySocialMedia", "Block", "CategorySocialMedia", nil),
			managedCatalogRule(14, "CategoryWebhooks", "Block", "CategoryWebhooks", nil),
			managedCatalogRule(15, "CategoryAI", "Block", "CategoryAI", nil),
			managedCatalogRule(16, "SignalAutomatedBrowser", "Block", "SignalAutomatedBrowser", nil),
			managedCatalogRule(17, "SignalKnownBotDataCenter", "Block", "SignalKnownBotDataCenter", nil),
			managedCatalogRule(18, "SignalNonBrowserUserAgent", "Block", "SignalNonBrowserUserAgent", nil),
			managedCatalogRule(19, "TGT_VolumetricIpTokenAbsent", "Challenge", "TGT_VolumetricIpTokenAbsent", nil),
			managedCatalogRule(20, "TGT_TokenAbsent", "Count", "TGT_TokenAbsent", nil),
			managedCatalogRule(21, "TGT_VolumetricSession", "Captcha", "TGT_VolumetricSession", nil),
			managedCatalogRule(22, "TGT_VolumetricSessionMaximum", "Block", "TGT_VolumetricSessionMaximum", nil),
			managedCatalogRule(23, "TGT_SignalAutomatedBrowser", "Captcha", "TGT_SignalAutomatedBrowser", nil),
			managedCatalogRule(24, "TGT_SignalBrowserAutomationExtension", "Captcha", "TGT_SignalBrowserAutomationExtension", nil),
			managedCatalogRule(25, "TGT_SignalBrowserInconsistency", "Captcha", "TGT_SignalBrowserInconsistency", nil),
			managedCatalogRule(26, "TGT_ML_CoordinatedActivityLow", "Challenge", "TGT_ML_CoordinatedActivityLow", nil),
			managedCatalogRule(27, "TGT_ML_CoordinatedActivityMedium", "Captcha", "TGT_ML_CoordinatedActivityMedium", nil),
			managedCatalogRule(28, "TGT_ML_CoordinatedActivityHigh", "Captcha", "TGT_ML_CoordinatedActivityHigh", nil),
			managedCatalogRule(29, "TGT_TokenReuseIpLow", "Count", "TGT_TokenReuseIpLow", nil),
			managedCatalogRule(30, "TGT_TokenReuseIpMedium", "Captcha", "TGT_TokenReuseIpMedium", nil),
			managedCatalogRule(31, "TGT_TokenReuseIpHigh", "Block", "TGT_TokenReuseIpHigh", nil),
			managedCatalogRule(32, "TGT_TokenReuseCountryLow", "Count", "TGT_TokenReuseCountryLow", nil),
			managedCatalogRule(33, "TGT_TokenReuseCountryMedium", "Captcha", "TGT_TokenReuseCountryMedium", nil),
			managedCatalogRule(34, "TGT_TokenReuseCountryHigh", "Block", "TGT_TokenReuseCountryHigh", nil),
			managedCatalogRule(35, "TGT_TokenReuseAsnLow", "Count", "TGT_TokenReuseAsnLow", nil),
			managedCatalogRule(36, "TGT_TokenReuseAsnMedium", "Captcha", "TGT_TokenReuseAsnMedium", nil),
			managedCatalogRule(37, "TGT_TokenReuseAsnHigh", "Block", "TGT_TokenReuseAsnHigh", nil),
		},
	},
	{
		Name:      "AWSManagedRulesATPRuleSet",
		Namespace: "awswaf:managed:aws:atp",
		WCU:       50,
		ExtraLabels: []string{
			"awswaf:managed:aws:atp:signal:credential_compromised",
		},
		Rules: []ManagedRule{
			managedCatalogRule(0, "UnsupportedCognitoIDP", "Block", "UnsupportedCognitoIDP", nil),
			managedCatalogRule(1, "VolumetricIpHigh", "Block", "VolumetricIpHigh", nil),
			managedCatalogRule(2, "VolumetricSession", "Block", "VolumetricSession", nil),
			managedCatalogRule(3, "AttributeCompromisedCredentials", "Block", "AttributeCompromisedCredentials", nil),
			managedCatalogRule(4, "AttributeUsernameTraversal", "Block", "AttributeUsernameTraversal", nil),
			managedCatalogRule(5, "AttributePasswordTraversal", "Block", "AttributePasswordTraversal", nil),
			managedCatalogRule(6, "AttributeLongSession", "Block", "AttributeLongSession", nil),
			managedCatalogRule(7, "TokenRejected", "Block", "TokenRejected", nil),
			managedCatalogRule(8, "SignalMissingCredential", "Block", "SignalMissingCredential", nil),
			managedCatalogRule(9, "VolumetricIpFailedLoginResponseHigh", "Block", "VolumetricIpFailedLoginResponseHigh", nil),
			managedCatalogRule(10, "VolumetricSessionFailedLoginResponseHigh", "Block", "VolumetricSessionFailedLoginResponseHigh", nil),
		},
	},
	{
		Name:      "AWSManagedRulesACFPRuleSet",
		Namespace: "awswaf:managed:aws:acfp",
		WCU:       50,
		ExtraLabels: []string{
			"awswaf:managed:aws:acfp:signal:credential_compromised",
		},
		Rules: []ManagedRule{
			managedCatalogRule(0, "UnsupportedCognitoIDP", "Block", "UnsupportedCognitoIDP", nil),
			managedCatalogRule(1, "AllRequests", "Challenge", "AllRequests", nil),
			managedCatalogRule(2, "RiskScoreHigh", "Block", "RiskScoreHigh", nil),
			managedCatalogRule(3, "SignalCredentialCompromised", "Block", "SignalCredentialCompromised", nil),
			managedCatalogRule(4, "SignalClientHumanInteractivityAbsentLow", "Captcha", "SignalClientHumanInteractivityAbsentLow", nil),
			managedCatalogRule(5, "AutomatedBrowser", "Block", "AutomatedBrowser", nil),
			managedCatalogRule(6, "BrowserInconsistency", "Captcha", "BrowserInconsistency", nil),
			managedCatalogRule(7, "VolumetricIpHigh", "Captcha", "VolumetricIpHigh", nil),
			managedCatalogRule(8, "VolumetricSessionHigh", "Block", "VolumetricSessionHigh", nil),
			managedCatalogRule(9, "AttributeUsernameTraversalHigh", "Block", "AttributeUsernameTraversalHigh", nil),
			managedCatalogRule(10, "VolumetricPhoneNumberHigh", "Block", "VolumetricPhoneNumberHigh", nil),
			managedCatalogRule(11, "VolumetricAddressHigh", "Block", "VolumetricAddressHigh", nil),
			managedCatalogRule(12, "VolumetricAddressLow", "Captcha", "VolumetricAddressLow", nil),
			managedCatalogRule(13, "VolumetricIPSuccessfulResponse", "Block", "VolumetricIPSuccessfulResponse", nil),
			managedCatalogRule(14, "VolumetricSessionSuccessfulResponse", "Block", "VolumetricSessionSuccessfulResponse", nil),
			managedCatalogRule(15, "VolumetricSessionTokenReuseIp", "Block", "VolumetricSessionTokenReuseIp", nil),
		},
	},
	{
		Name:      "AWSManagedRulesAntiDDoSRuleSet",
		Namespace: "awswaf:managed:aws:anti-ddos",
		WCU:       50,
		ExtraLabels: []string{
			"awswaf:managed:aws:anti-ddos:event-detected",
			"awswaf:managed:aws:anti-ddos:ddos-request",
		},
		Rules: []ManagedRule{
			managedCatalogRule(0, "ChallengeAllDuringEvent", "Challenge", "ChallengeAllDuringEvent", nil),
			managedCatalogRule(1, "ChallengeDDoSRequests", "Challenge", "ChallengeDDoSRequests", nil),
			managedCatalogRule(2, "DDoSRequests", "Block", "DDoSRequests", nil),
		},
	},
}

// managedRuleGroupIndex indexes the catalog by group name. Its
// initialisation also qualifies every rule's label with the group's
// namespace: the catalog entries above declare the label suffixes the
// documentation lists, and the qualified form is what a matching rule adds
// to the request.
var managedRuleGroupIndex = func() map[string]*ManagedRuleGroup {
	index := make(map[string]*ManagedRuleGroup, len(managedRuleGroups))
	for _, group := range managedRuleGroups {
		for i := range group.Rules {
			group.Rules[i].Label = group.Namespace + ":" + group.Rules[i].Label
		}
		index[group.Name] = group
	}
	return index
}()

// LookupManagedRuleGroup returns the catalog entry for the vendor and
// group name. AWS is the only vendor.
func LookupManagedRuleGroup(vendorName, name string) (*ManagedRuleGroup, bool) {
	if !strings.EqualFold(vendorName, "AWS") {
		return nil, false
	}
	group, ok := managedRuleGroupIndex[name]
	return group, ok
}

// managedCatalogRule builds one catalog entry.
func managedCatalogRule(priority int32, name, action, label string, stmt *wafstore.Statement) ManagedRule {
	return ManagedRule{Name: name, Priority: priority, Action: action, Label: label, Statement: stmt}
}

// --- request components ---

func uriPathField() *wafstore.FieldToMatch {
	return &wafstore.FieldToMatch{UriPath: &wafstore.All{}}
}

func queryStringField() *wafstore.FieldToMatch {
	return &wafstore.FieldToMatch{QueryString: &wafstore.All{}}
}

func allQueryArgumentsField() *wafstore.FieldToMatch {
	return &wafstore.FieldToMatch{AllQueryArguments: &wafstore.All{}}
}

// bodyField returns the body component with the oversize handling the
// managed rule listings document: inspection continues past the
// component's size limit.
func bodyField() *wafstore.FieldToMatch {
	return &wafstore.FieldToMatch{Body: &wafstore.Body{OversizeHandling: "CONTINUE"}}
}

// headersField returns every request header, keys and values, continuing
// past the header inspection limits.
func headersField() *wafstore.FieldToMatch {
	return &wafstore.FieldToMatch{Headers: &wafstore.Headers{
		MatchPattern:     wafstore.HeaderMatchPattern{All: &wafstore.All{}},
		MatchScope:       "ALL",
		OversizeHandling: "CONTINUE",
	}}
}

// cookiesField returns every cookie, names and values, continuing past the
// cookie inspection limits.
func cookiesField() *wafstore.FieldToMatch {
	return &wafstore.FieldToMatch{Cookies: &wafstore.Cookies{
		MatchPattern:     wafstore.CookieMatchPattern{All: &wafstore.All{}},
		MatchScope:       "ALL",
		OversizeHandling: "CONTINUE",
	}}
}

func singleHeader(name string) *wafstore.FieldToMatch {
	return &wafstore.FieldToMatch{SingleHeader: &wafstore.SingleHeader{Name: name}}
}

func methodField() *wafstore.FieldToMatch {
	return &wafstore.FieldToMatch{Method: &wafstore.All{}}
}

// --- text transformations ---

func ttNone() []*wafstore.TextTransformation {
	return []*wafstore.TextTransformation{{Priority: 0, Type: "NONE"}}
}

func ttLower() []*wafstore.TextTransformation {
	return []*wafstore.TextTransformation{{Priority: 0, Type: "LOWERCASE"}}
}

// ttURLLower decodes URL encoding before lowercasing, so encoded
// attack payloads match their plain forms.
func ttURLLower() []*wafstore.TextTransformation {
	return []*wafstore.TextTransformation{
		{Priority: 0, Type: "URL_DECODE"},
		{Priority: 1, Type: "LOWERCASE"},
	}
}

// --- statement builders ---

func containsStatement(field *wafstore.FieldToMatch, tts []*wafstore.TextTransformation, needle string) *wafstore.Statement {
	return &wafstore.Statement{ByteMatchStatement: &wafstore.ByteMatchStatement{
		FieldToMatch:         field,
		SearchString:         []byte(needle),
		PositionalConstraint: "CONTAINS",
		TextTransformations:  tts,
	}}
}

// orContains builds the disjunction of substring matches over the given
// needles.
func orContains(field *wafstore.FieldToMatch, tts []*wafstore.TextTransformation, needles ...string) *wafstore.Statement {
	statements := make([]*wafstore.Statement, 0, len(needles))
	for _, needle := range needles {
		statements = append(statements, containsStatement(field, tts, needle))
	}
	return &wafstore.Statement{OrStatement: &wafstore.OrStatement{Statements: statements}}
}

func regexStatement(field *wafstore.FieldToMatch, tts []*wafstore.TextTransformation, pattern string) *wafstore.Statement {
	return &wafstore.Statement{RegexMatchStatement: &wafstore.RegexMatchStatement{
		RegexString:         pattern,
		FieldToMatch:        field,
		TextTransformations: tts,
	}}
}

func notStatement(inner *wafstore.Statement) *wafstore.Statement {
	return &wafstore.Statement{NotStatement: &wafstore.NotStatement{Statement: inner}}
}

func sizeAbove(field *wafstore.FieldToMatch, size int64) *wafstore.Statement {
	return &wafstore.Statement{SizeConstraintStatement: &wafstore.SizeConstraintStatement{
		FieldToMatch:        field,
		ComparisonOperator:  "GT",
		Size:                size,
		TextTransformations: ttNone(),
	}}
}

func xssStatement(field *wafstore.FieldToMatch) *wafstore.Statement {
	return &wafstore.Statement{XssMatchStatement: &wafstore.XssMatchStatement{
		FieldToMatch:        field,
		TextTransformations: ttNone(),
	}}
}

func sqliStatement(field *wafstore.FieldToMatch) *wafstore.Statement {
	return &wafstore.Statement{SqliMatchStatement: &wafstore.SqliMatchStatement{
		FieldToMatch:        field,
		SensitivityLevel:    "LOW",
		TextTransformations: ttURLLower(),
	}}
}

// --- signature sets ---

// ec2MetaDataStatement matches attempts to reach the EC2 instance
// metadata service: its link-local address and the IMDS and user-data
// paths.
func ec2MetaDataStatement(field *wafstore.FieldToMatch) *wafstore.Statement {
	return orContains(field, ttURLLower(),
		"169.254.169.254", "latest/meta-data", "latest/user-data")
}

// lfiStatement matches generic path traversal: the parent-directory
// sequences in plain and URL-encoded form.
func lfiStatement(field *wafstore.FieldToMatch) *wafstore.Statement {
	return orContains(field, ttURLLower(), "../", "..\\")
}

// restrictedExtensionsStatement matches system file extensions that are
// unsafe to serve or run.
func restrictedExtensionsStatement(field *wafstore.FieldToMatch) *wafstore.Statement {
	return orContains(field, ttURLLower(), ".log", ".ini", ".conf", ".cfg", ".bak", ".old")
}

// rfiStatement matches remote file inclusion payloads: a URL scheme
// followed by a literal IPv4 host.
func rfiStatement(field *wafstore.FieldToMatch) *wafstore.Statement {
	return regexStatement(field, ttURLLower(), `(https?|ftps?|file)://[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}`)
}

// javaDeserializationStatement matches Java deserialization and Spring
// expression-language remote code execution payloads.
func javaDeserializationStatement(field *wafstore.FieldToMatch) *wafstore.Statement {
	return orContains(field, ttURLLower(),
		"java.lang.runtime", "getruntime(", "class.module.classloader", "processbuilder")
}

// log4jStatement matches Log4j lookup injection in plain and URL-encoded
// form, case-insensitively through the lowercase transformation.
func log4jStatement(field *wafstore.FieldToMatch) *wafstore.Statement {
	return regexStatement(field, ttURLLower(), `\$\{\s*jndi\s*:`)
}

// sqliExtendedStatement matches SQL injection tokens beyond the built-in
// SQLi detector: the union-select exfiltration form, schema enumeration,
// and the time-based and file-access functions.
func sqliExtendedStatement(field *wafstore.FieldToMatch) *wafstore.Statement {
	return orContains(field, ttURLLower(),
		"union select", "information_schema", "sleep(", "benchmark(", "waitfor delay",
		"load_file(", "into outfile", "or 1=1", "' or '")
}

// linuxLFIStatement matches Linux-specific file disclosure targets: the
// procfs tree and the account and network configuration files.
func linuxLFIStatement(field *wafstore.FieldToMatch) *wafstore.Statement {
	return orContains(field, ttURLLower(),
		"/proc/", "/etc/passwd", "/etc/shadow", "/etc/hosts")
}

// unixShellStatement matches Unix shell command injection and variable
// expansion attempts.
func unixShellStatement(field *wafstore.FieldToMatch) *wafstore.Statement {
	return orContains(field, ttURLLower(),
		"$home", "$path", "$(", "${", "/bin/sh", "/bin/bash")
}

// windowsShellStatement matches Windows command-shell injection
// sequences: command chaining and continuation operators around shell
// interpreters and diagnostic commands.
func windowsShellStatement(field *wafstore.FieldToMatch) *wafstore.Statement {
	return orContains(field, ttURLLower(),
		"||nslookup", "&&nslookup", ";cmd", "|cmd", "&cmd", "%comspec%")
}

// powerShellStatement matches PowerShell command injection attempts.
func powerShellStatement(field *wafstore.FieldToMatch) *wafstore.Statement {
	return orContains(field, ttURLLower(),
		"invoke-expression", "powershell", "new-object", "start-process",
		"downloadstring", "iex ")
}

// phpStatement matches PHP code injection: high-risk functions and the
// superglobal variables.
func phpStatement(field *wafstore.FieldToMatch) *wafstore.Statement {
	return orContains(field, ttURLLower(),
		"fsockopen", "$_get", "$_post", "$_request", "$_server", "$_env",
		"eval(", "base64_decode", "file_get_contents", "php://")
}

// catalogRegexCheck compiles every inline regex of the catalog once at
// initialisation, so an invalid pattern fails at startup rather than at
// match time.
var catalogRegexCheck = func() struct{} {
	validateCatalogRegexes()
	return struct{}{}
}()

func validateCatalogRegexes() {
	var compiled []*regexp.Regexp
	for _, group := range managedRuleGroups {
		for _, rule := range group.Rules {
			regexOfStatement(rule.Statement, &compiled)
		}
	}
}

func regexOfStatement(stmt *wafstore.Statement, compiled *[]*regexp.Regexp) {
	if stmt == nil {
		return
	}
	switch {
	case stmt.RegexMatchStatement != nil:
		re, err := regexp.Compile(stmt.RegexMatchStatement.RegexString)
		if err != nil {
			panic("managed rule catalog regex: " + err.Error())
		}
		*compiled = append(*compiled, re)
	case stmt.AndStatement != nil:
		for _, sub := range stmt.AndStatement.Statements {
			regexOfStatement(sub, compiled)
		}
	case stmt.OrStatement != nil:
		for _, sub := range stmt.OrStatement.Statements {
			regexOfStatement(sub, compiled)
		}
	case stmt.NotStatement != nil:
		regexOfStatement(stmt.NotStatement.Statement, compiled)
	}
}
