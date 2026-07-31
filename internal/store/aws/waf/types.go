package waf

// Package waf provides WAF (Web Application Firewall) data store implementations
// for vorpalstacks.

import (
	"time"

	"vorpalstacks/internal/utils/aws/types"
)

// WebACL represents a WAF Web Access Control List.
type WebACL struct {
	ID                     string            `json:"id"`
	Name                   string            `json:"name"`
	ARN                    string            `json:"arn"`
	MetricName             string            `json:"metricName"`
	Capacity               int64             `json:"capacity"`
	Rules                  []*Rule           `json:"rules"`
	DefaultAction          interface{}       `json:"defaultAction"`
	VisibilityConfig       *VisibilityConfig `json:"visibilityConfig"`
	Scope                  string            `json:"scope"`
	Description            string            `json:"description"`
	LockToken              string            `json:"lockToken"`
	Tags                   []types.Tag       `json:"tags"`
	CustomResponseBodies   interface{}       `json:"customResponseBodies,omitempty"`
	CaptchaConfig          interface{}       `json:"captchaConfig,omitempty"`
	ChallengeConfig        interface{}       `json:"challengeConfig,omitempty"`
	TokenDomains           interface{}       `json:"tokenDomains,omitempty"`
	LabelNamespace         string            `json:"labelNamespace,omitempty"`
	AssociationConfig      interface{}       `json:"associationConfig,omitempty"`
	ApplicationConfig      interface{}       `json:"applicationConfig,omitempty"`
	MonetizationConfig     interface{}       `json:"monetizationConfig,omitempty"`
	DataProtectionConfig   interface{}       `json:"dataProtectionConfig,omitempty"`
	OnSourceDDoSProtection interface{}       `json:"onSourceDDoSProtectionConfig,omitempty"`
	CreatedAt              time.Time         `json:"createdAt"`
	ModifiedAt             time.Time         `json:"modifiedAt"`
}

// Rule represents a WAF rule.
type Rule struct {
	ID               string            `json:"ruleId"`
	Name             string            `json:"name"`
	ARN              string            `json:"arn"`
	MetricName       string            `json:"metricName"`
	Priority         int32             `json:"priority"`
	Action           interface{}       `json:"action"`
	Statement        *Statement        `json:"statement"`
	OverrideAction   *Action           `json:"overrideAction"`
	VisibilityConfig *VisibilityConfig `json:"visibilityConfig"`
	RuleLabels       interface{}       `json:"ruleLabels,omitempty"`
	CaptchaConfig    interface{}       `json:"captchaConfig,omitempty"`
	ChallengeConfig  interface{}       `json:"challengeConfig,omitempty"`
	Tags             []types.Tag       `json:"tags"`
}

// Statement represents a WAF rule statement. Statement is a union —
// exactly one of the following fields should be non-nil at any time,
// matching the AWS WAF API contract. All sixteen Smithy statement types
// are modelled here to support future rule evaluation engine work.
type Statement struct {
	AndStatement                *AndStatement                `json:"AndStatement,omitempty"`
	AsnMatchStatement           *AsnMatchStatement           `json:"AsnMatchStatement,omitempty"`
	ByteMatchStatement          *ByteMatchStatement          `json:"ByteMatchStatement,omitempty"`
	GeoMatchStatement           *GeoMatchStatement           `json:"GeoMatchStatement,omitempty"`
	IPSetReferenceStatement     *IPSetReferenceStatement     `json:"IPSetReferenceStatement,omitempty"`
	LabelMatchStatement         *LabelMatchStatement         `json:"LabelMatchStatement,omitempty"`
	ManagedRuleGroupStatement   *ManagedRuleGroupStatement   `json:"ManagedRuleGroupStatement,omitempty"`
	NotStatement                *NotStatement                `json:"NotStatement,omitempty"`
	OrStatement                 *OrStatement                 `json:"OrStatement,omitempty"`
	RateBasedStatement          *RateBasedStatement          `json:"RateBasedStatement,omitempty"`
	RegexMatchStatement         *RegexMatchStatement         `json:"RegexMatchStatement,omitempty"`
	RegexPatternSetRefStatement *RegexPatternSetRefStatement `json:"RegexPatternSetReferenceStatement,omitempty"`
	RuleGroupReferenceStatement *RuleGroupReferenceStatement `json:"RuleGroupReferenceStatement,omitempty"`
	SizeConstraintStatement     *SizeConstraintStatement     `json:"SizeConstraintStatement,omitempty"`
	SqliMatchStatement          *SqliMatchStatement          `json:"SqliMatchStatement,omitempty"`
	XssMatchStatement           *XssMatchStatement           `json:"XssMatchStatement,omitempty"`
}

// ByteMatchStatement represents a byte match rule statement.
type ByteMatchStatement struct {
	SearchString         []byte                `json:"SearchString"`
	FieldToMatch         *FieldToMatch         `json:"FieldToMatch"`
	TextTransformations  []*TextTransformation `json:"TextTransformations"`
	PositionalConstraint string                `json:"PositionalConstraint"`
}

// SqliMatchStatement represents an SQL injection match rule statement.
type SqliMatchStatement struct {
	FieldToMatch        *FieldToMatch         `json:"FieldToMatch"`
	TextTransformations []*TextTransformation `json:"TextTransformations"`
	SensitivityLevel    string                `json:"SensitivityLevel,omitempty"`
}

// XssMatchStatement represents an XSS match rule statement.
type XssMatchStatement struct {
	FieldToMatch        *FieldToMatch         `json:"FieldToMatch"`
	TextTransformations []*TextTransformation `json:"TextTransformations"`
}

// SizeConstraintStatement represents a size constraint rule statement.
type SizeConstraintStatement struct {
	FieldToMatch        *FieldToMatch         `json:"FieldToMatch"`
	TextTransformations []*TextTransformation `json:"TextTransformations"`
	ComparisonOperator  string                `json:"ComparisonOperator"`
	Size                int64                 `json:"Size"`
}

// GeoMatchStatement represents a geographic match rule statement.
type GeoMatchStatement struct {
	CountryCodes      []string           `json:"CountryCodes,omitempty"`
	ForwardedIPConfig *ForwardedIPConfig `json:"ForwardedIPConfig,omitempty"`
}

// RateBasedStatement represents a rate-based rule statement.
type RateBasedStatement struct {
	Limit               int64                    `json:"Limit"`
	EvaluationWindowSec int64                    `json:"EvaluationWindowSec,omitempty"`
	AggregateKeyType    string                   `json:"AggregateKeyType"`
	ScopeDownStatement  *Statement               `json:"ScopeDownStatement,omitempty"`
	ForwardedIPConfig   *ForwardedIPConfig       `json:"ForwardedIPConfig,omitempty"`
	CustomKeys          []map[string]interface{} `json:"CustomKeys,omitempty"`
}

// ManagedRuleGroupStatement represents a managed rule group statement.
type ManagedRuleGroupStatement struct {
	Name                    string                   `json:"Name"`
	VendorName              string                   `json:"VendorName"`
	Version                 string                   `json:"Version,omitempty"`
	ExcludedRules           []ExcludedRule           `json:"ExcludedRules,omitempty"`
	ScopeDownStatement      *Statement               `json:"ScopeDownStatement,omitempty"`
	ManagedRuleGroupConfigs []map[string]interface{} `json:"ManagedRuleGroupConfigs,omitempty"`
	RuleActionOverrides     map[string]interface{}   `json:"RuleActionOverrides,omitempty"`
}

// AndStatement represents a logical AND statement.
type AndStatement struct {
	Statements []*Statement `json:"Statements"`
}

// OrStatement represents a logical OR statement.
type OrStatement struct {
	Statements []*Statement `json:"Statements"`
}

// NotStatement represents a logical NOT statement.
type NotStatement struct {
	Statement *Statement `json:"Statement"`
}

// RegexMatchStatement matches a web request component against a single
// regular expression pattern that you provide inline in the statement.
type RegexMatchStatement struct {
	RegexString         string                `json:"RegexString"`
	FieldToMatch        *FieldToMatch         `json:"FieldToMatch"`
	TextTransformations []*TextTransformation `json:"TextTransformations"`
}

// RegexPatternSetRefStatement references a RegexPatternSet by ARN. This
// is distinct from RegexMatchStatement, which stores the regex pattern
// inline.
type RegexPatternSetRefStatement struct {
	ARN                 string                `json:"ARN"`
	FieldToMatch        *FieldToMatch         `json:"FieldToMatch"`
	TextTransformations []*TextTransformation `json:"TextTransformations"`
}

// IPSetReferenceStatement represents an IP set reference statement.
type IPSetReferenceStatement struct {
	ARN                    string                  `json:"ARN"`
	IPSetForwardedIPConfig *IPSetForwardedIPConfig `json:"IPSetForwardedIPConfig,omitempty"`
}

// LabelMatchStatement represents a label match statement.
type LabelMatchStatement struct {
	Scope string `json:"Scope"`
	Key   string `json:"Key"`
}

// AsnMatchStatement represents an ASN match statement.
type AsnMatchStatement struct {
	AsnList           []string           `json:"AsnList,omitempty"`
	ForwardedIPConfig *ForwardedIPConfig `json:"ForwardedIPConfig,omitempty"`
}

// RuleGroupReferenceStatement references a rule group by ARN.
type RuleGroupReferenceStatement struct {
	ARN           string         `json:"ARN"`
	ExcludedRules []ExcludedRule `json:"ExcludedRules,omitempty"`
}

// ExcludedRule identifies a rule to exclude from a managed rule group
// or rule group reference.
type ExcludedRule struct {
	Name string `json:"Name"`
}

// ForwardedIPConfig holds the configuration for forwarded IP processing.
type ForwardedIPConfig struct {
	HeaderName       string `json:"HeaderName"`
	FallbackBehavior string `json:"FallbackBehavior"`
}

// IPSetForwardedIPConfig holds the configuration for forwarded IP in an IP set.
type IPSetForwardedIPConfig struct {
	HeaderName       string `json:"HeaderName"`
	Position         string `json:"Position"`
	FallbackBehavior string `json:"FallbackBehavior"`
}

// FieldToMatch is a union that specifies the part of the web request
// that you want WAF to inspect. Exactly one field should be non-nil.
type FieldToMatch struct {
	AllQueryArguments   *All                 `json:"AllQueryArguments,omitempty"`
	Body                *Body                `json:"Body,omitempty"`
	Cookies             *Cookies             `json:"Cookies,omitempty"`
	HeaderOrder         *HeaderOrderMatch    `json:"HeaderOrder,omitempty"`
	Headers             *Headers             `json:"Headers,omitempty"`
	JA3Fingerprint      *JA3Fingerprint      `json:"JA3Fingerprint,omitempty"`
	JA4Fingerprint      *JA4Fingerprint      `json:"JA4Fingerprint,omitempty"`
	JsonBody            *JsonBody            `json:"JsonBody,omitempty"`
	Method              *All                 `json:"Method,omitempty"`
	QueryString         *All                 `json:"QueryString,omitempty"`
	SingleHeader        *SingleHeader        `json:"SingleHeader,omitempty"`
	SingleQueryArgument *SingleQueryArgument `json:"SingleQueryArgument,omitempty"`
	UriFragment         *UriFragment         `json:"UriFragment,omitempty"`
	UriPath             *All                 `json:"UriPath,omitempty"`
}

// All is an empty marker struct used by several FieldToMatch variants
// that do not carry configuration data.
type All struct{}

// Body represents the request body field to match.
type Body struct {
	OversizeHandling string `json:"OversizeHandling,omitempty"`
}

// Cookies represents the cookies field to match.
type Cookies struct {
	MatchPattern     CookieMatchPattern `json:"MatchPattern"`
	MatchScope       string             `json:"MatchScope"`
	OversizeHandling string             `json:"OversizeHandling"`
}

// HeaderOrderMatch represents the header order field to match.
type HeaderOrderMatch struct {
	OversizeHandling string `json:"OversizeHandling,omitempty"`
}

// Headers represents the headers field to match.
type Headers struct {
	MatchPattern     HeaderMatchPattern `json:"MatchPattern"`
	MatchScope       string             `json:"MatchScope"`
	OversizeHandling string             `json:"OversizeHandling"`
}

// JA3Fingerprint represents the JA3 fingerprint field to match.
type JA3Fingerprint struct {
	FallbackBehavior string `json:"FallbackBehavior"`
}

// JA4Fingerprint represents the JA4 fingerprint field to match.
type JA4Fingerprint struct {
	FallbackBehavior string `json:"FallbackBehavior"`
}

// JsonBody represents the JSON body field to match.
type JsonBody struct {
	MatchPattern            JsonMatchPattern `json:"MatchPattern"`
	MatchScope              string           `json:"MatchScope"`
	InvalidFallbackBehavior string           `json:"InvalidFallbackBehavior,omitempty"`
	OversizeHandling        string           `json:"OversizeHandling,omitempty"`
}

// SingleHeader identifies a single header by name.
type SingleHeader struct {
	Name string `json:"Name"`
}

// SingleQueryArgument identifies a single query argument by name.
type SingleQueryArgument struct {
	Name string `json:"Name"`
}

// UriFragment represents the URI fragment field to match.
type UriFragment struct {
	FallbackBehavior string `json:"FallbackBehavior"`
}

// CookieMatchPattern specifies the cookies to inspect.
type CookieMatchPattern struct {
	All             *All     `json:"All,omitempty"`
	IncludedCookies []string `json:"IncludedCookies,omitempty"`
	ExcludedCookies []string `json:"ExcludedCookies,omitempty"`
}

// HeaderMatchPattern specifies the headers to inspect.
type HeaderMatchPattern struct {
	All             *All     `json:"All,omitempty"`
	IncludedHeaders []string `json:"IncludedHeaders,omitempty"`
	ExcludedHeaders []string `json:"ExcludedHeaders,omitempty"`
}

// JsonMatchPattern specifies the JSON paths to inspect.
type JsonMatchPattern struct {
	All           *All     `json:"All,omitempty"`
	IncludedPaths []string `json:"IncludedPaths,omitempty"`
}

// TextTransformation represents a text transformation for matching.
type TextTransformation struct {
	Priority int    `json:"Priority"`
	Type     string `json:"Type"`
}

// Action represents a WAF rule action. For DefaultAction on a WebACL,
// only Allow and Block are valid (terminating actions). For RuleAction
// on individual rules, all six types are valid.
type Action struct {
	Allow     *AllowAction     `json:"Allow,omitempty"`
	Block     *BlockAction     `json:"Block,omitempty"`
	Count     *CountAction     `json:"Count,omitempty"`
	Captcha   *CaptchaAction   `json:"Captcha,omitempty"`
	Challenge *ChallengeAction `json:"Challenge,omitempty"`
	Monetize  *MonetizeAction  `json:"Monetize,omitempty"`
}

// AllowAction instructs WAF to allow the web request. The optional
// CustomRequestHandling lets you insert custom headers into the request.
type AllowAction struct {
	CustomRequestHandling *CustomRequestHandling `json:"CustomRequestHandling,omitempty"`
}

// BlockAction instructs WAF to block the web request. The optional
// CustomResponse lets you customise the HTTP response returned to the client.
type BlockAction struct {
	CustomResponse *CustomResponse `json:"CustomResponse,omitempty"`
}

// CountAction instructs WAF to count the request and continue evaluating
// remaining rules. The optional CustomRequestHandling lets you insert
// custom headers into the request.
type CountAction struct {
	CustomRequestHandling *CustomRequestHandling `json:"CustomRequestHandling,omitempty"`
}

// CaptchaAction instructs WAF to run a CAPTCHA check. The optional
// CustomRequestHandling lets you insert custom headers into the request.
type CaptchaAction struct {
	CustomRequestHandling *CustomRequestHandling `json:"CustomRequestHandling,omitempty"`
}

// ChallengeAction instructs WAF to run a silent challenge check. The
// optional CustomRequestHandling lets you insert custom headers.
type ChallengeAction struct {
	CustomRequestHandling *CustomRequestHandling `json:"CustomRequestHandling,omitempty"`
}

// MonetizeAction instructs WAF to return an HTTP 402 Payment Required
// response. PriceMultiplier is a string ("1"–"100") per the Smithy
// PriceMultiplier shape (pattern ^([1-9][0-9]?|100)$). This action is
// available only for web ACLs associated with Amazon CloudFront.
type MonetizeAction struct {
	PriceMultiplier string `json:"PriceMultiplier,omitempty"`
}

// CustomRequestHandling defines custom header injection for allow, count,
// captcha, and challenge actions.
type CustomRequestHandling struct {
	InsertHeaders []CustomHTTPHeader `json:"InsertHeaders,omitempty"`
}

// CustomResponse defines a custom HTTP response for block actions.
type CustomResponse struct {
	ResponseCode          int                `json:"ResponseCode,omitempty"`
	CustomResponseBodyKey string             `json:"CustomResponseBodyKey,omitempty"`
	ResponseHeaders       []CustomHTTPHeader `json:"ResponseHeaders,omitempty"`
}

// CustomHTTPHeader is a single name/value HTTP header used by
// CustomRequestHandling and CustomResponse.
type CustomHTTPHeader struct {
	Name  string `json:"Name,omitempty"`
	Value string `json:"Value,omitempty"`
}

// VisibilityConfig holds the visibility configuration for a rule or web ACL.
type VisibilityConfig struct {
	SampledRequestsEnabled   bool   `json:"SampledRequestsEnabled"`
	CloudWatchMetricsEnabled bool   `json:"CloudWatchMetricsEnabled"`
	MetricName               string `json:"MetricName"`
}

// RuleGroup represents a WAF rule group.
type RuleGroup struct {
	ID                   string            `json:"id"`
	Name                 string            `json:"name"`
	ARN                  string            `json:"arn"`
	Scope                string            `json:"scope"`
	Capacity             int64             `json:"capacity"`
	Rules                []*Rule           `json:"rules"`
	VisibilityConfig     *VisibilityConfig `json:"visibilityConfig"`
	Description          string            `json:"description"`
	LockToken            string            `json:"lockToken"`
	Tags                 []types.Tag       `json:"tags"`
	CustomResponseBodies interface{}       `json:"customResponseBodies,omitempty"`
	LabelNamespace       string            `json:"labelNamespace,omitempty"`
	AvailableLabels      interface{}       `json:"availableLabels,omitempty"`
	ConsumedLabels       interface{}       `json:"consumedLabels,omitempty"`
	MonetizationConfig   interface{}       `json:"monetizationConfig,omitempty"`
	CreatedAt            time.Time         `json:"createdAt"`
	ModifiedAt           time.Time         `json:"modifiedAt"`
}

// IPSet represents a WAF IP set.
type IPSet struct {
	ID               string      `json:"id"`
	Name             string      `json:"name"`
	ARN              string      `json:"arn"`
	Scope            string      `json:"scope"`
	IPAddressVersion string      `json:"ipAddressVersion"`
	Description      string      `json:"description"`
	Addresses        []string    `json:"addresses"`
	LockToken        string      `json:"lockToken"`
	Tags             []types.Tag `json:"tags"`
	CreatedAt        time.Time   `json:"createdAt"`
	ModifiedAt       time.Time   `json:"modifiedAt"`
}

// RegexPatternSet represents a WAF regular expression pattern set.
type RegexPatternSet struct {
	ID              string      `json:"id"`
	Name            string      `json:"name"`
	ARN             string      `json:"arn"`
	Scope           string      `json:"scope"`
	Description     string      `json:"description"`
	RegularPatterns []string    `json:"regularPatterns"`
	LockToken       string      `json:"lockToken"`
	Tags            []types.Tag `json:"tags"`
	CreatedAt       time.Time   `json:"createdAt"`
	ModifiedAt      time.Time   `json:"modifiedAt"`
}

// WebACLAssociation represents an association between a web ACL and a resource.
type WebACLAssociation struct {
	WebACLArn   string `json:"webAclArn"`
	ResourceArn string `json:"resourceArn"`
}

// WebACLListResult represents the result of listing web ACLs.
type WebACLListResult struct {
	WebACLs     []*WebACL
	IsTruncated bool
	NextMarker  string
}

// RuleGroupListResult represents the result of listing rule groups.
type RuleGroupListResult struct {
	RuleGroups  []*RuleGroup
	IsTruncated bool
	NextMarker  string
}

// IPSetListResult represents the result of listing IP sets.
type IPSetListResult struct {
	IPSets      []*IPSet
	IsTruncated bool
	NextMarker  string
}

// RegexPatternSetListResult represents the result of listing regex pattern sets.
type RegexPatternSetListResult struct {
	RegexPatternSets []*RegexPatternSet
	IsTruncated      bool
	NextMarker       string
}

// LoggingConfiguration represents a WAF logging configuration.
type LoggingConfiguration struct {
	ResourceArn              string         `json:"resourceArn"`
	LogDestinationConfigs    []string       `json:"logDestinationConfigs"`
	LogScope                 string         `json:"logScope"`
	LogType                  string         `json:"logType"`
	LoggingFilter            *LoggingFilter `json:"loggingFilter,omitempty"`
	ManagedByFirewallManager bool           `json:"managedByFirewallManager"`
	RedactedFields           []interface{}  `json:"redactedFields,omitempty"`
	CreatedAt                time.Time      `json:"createdAt"`
}

// LoggingFilter represents filtering for logging configuration.
type LoggingFilter struct {
	DefaultBehavior string   `json:"defaultBehavior"`
	Filters         []Filter `json:"filters"`
}

// Filter represents a single filter in logging configuration.
type Filter struct {
	Behavior    string            `json:"behavior"`
	Conditions  []FilterCondition `json:"conditions"`
	Requirement string            `json:"requirement"`
}

// FilterCondition represents a condition in a logging filter.
type FilterCondition struct {
	ActionCondition    *ActionCondition    `json:"actionCondition,omitempty"`
	LabelNameCondition *LabelNameCondition `json:"labelNameCondition,omitempty"`
}

// ActionCondition represents an action condition for logging filter.
type ActionCondition struct {
	Action string `json:"action"`
}

// LabelNameCondition represents a label name condition for logging filter.
type LabelNameCondition struct {
	LabelName string `json:"labelName"`
}

// LoggingConfigurationListResult represents the result of listing logging configurations.
type LoggingConfigurationListResult struct {
	LoggingConfigurations []*LoggingConfiguration
	IsTruncated           bool
	NextMarker            string
}
