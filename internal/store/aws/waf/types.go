package waf

// Package waf provides WAF (Web Application Firewall) data store implementations
// for vorpalstacks.

import (
	"time"

	"vorpalstacks/internal/utils/aws/types"
)

// WebACL represents a WAF Web Access Control List.
type WebACL struct {
	ID               string            `json:"id"`
	Name             string            `json:"name"`
	ARN              string            `json:"arn"`
	MetricName       string            `json:"metricName"`
	Capacity         int64             `json:"capacity"`
	Rules            []*Rule           `json:"rules"`
	DefaultAction    interface{}       `json:"defaultAction"`
	VisibilityConfig *VisibilityConfig `json:"visibilityConfig"`
	Scope            string            `json:"scope"`
	Description      string            `json:"description"`
	LockToken        string            `json:"lockToken"`
	Tags             []types.Tag       `json:"tags"`
	CreatedAt        time.Time         `json:"createdAt"`
	ModifiedAt       time.Time         `json:"modifiedAt"`
}

// Rule represents a WAF rule.
type Rule struct {
	ID               string            `json:"ruleId"`
	Name             string            `json:"name"`
	ARN              string            `json:"arn"`
	MetricName       string            `json:"metricName"`
	Priority         int32             `json:"priority"`
	Predicates       []interface{}     `json:"predicates"`
	Action           interface{}       `json:"action"`
	Statement        *Statement        `json:"statement"`
	OverrideAction   *Action           `json:"overrideAction"`
	VisibilityConfig *VisibilityConfig `json:"visibilityConfig"`
	Tags             []types.Tag       `json:"tags"`
}

// Statement represents a WAF rule statement.
type Statement struct {
	ByteMatchStatement         *ByteMatchStatement         `json:"byteMatchStatement,omitempty"`
	SQLInjectionMatchStatement *SQLInjectionMatchStatement `json:"sqlInjectionMatchStatement,omitempty"`
	XSSMatchStatement          *XSSMatchStatement          `json:"xssMatchStatement,omitempty"`
	SizeConstraintStatement    *SizeConstraintStatement    `json:"sizeConstraintStatement,omitempty"`
	GeoMatchStatement          *GeoMatchStatement          `json:"geoMatchStatement,omitempty"`
	RateBasedStatement         *RateBasedStatement         `json:"rateBasedStatement,omitempty"`
	ManagedRuleGroupStatement  *ManagedRuleGroupStatement  `json:"managedRuleGroupStatement,omitempty"`
	AndStatement               *AndStatement               `json:"andStatement,omitempty"`
	OrStatement                *OrStatement                `json:"orStatement,omitempty"`
	NotStatement               *NotStatement               `json:"notStatement,omitempty"`
	RegexMatchStatement        *RegexMatchStatement        `json:"regexMatchStatement,omitempty"`
	IPSetReferenceStatement    *IPSetReferenceStatement    `json:"ipSetReferenceStatement,omitempty"`
}

// ByteMatchStatement represents a byte match rule statement.
type ByteMatchStatement struct {
	SearchString         []byte                `json:"searchString"`
	FieldToMatch         *FieldToMatch         `json:"fieldToMatch"`
	TextTransformations  []*TextTransformation `json:"textTransformations"`
	PositionalConstraint string                `json:"positionalConstraint"`
}

// SQLInjectionMatchStatement represents an SQL injection match rule statement.
type SQLInjectionMatchStatement struct {
	FieldToMatch        *FieldToMatch         `json:"fieldToMatch"`
	TextTransformations []*TextTransformation `json:"textTransformations"`
	SensitivityLevel    string                `json:"sensitivityLevel,omitempty"`
}

// XSSMatchStatement represents an XSS match rule statement.
type XSSMatchStatement struct {
	FieldToMatch        *FieldToMatch         `json:"fieldToMatch"`
	TextTransformations []*TextTransformation `json:"textTransformations"`
}

// SizeConstraintStatement represents a size constraint rule statement.
type SizeConstraintStatement struct {
	FieldToMatch        *FieldToMatch         `json:"fieldToMatch"`
	TextTransformations []*TextTransformation `json:"textTransformations"`
	ComparisonOperator  string                `json:"comparisonOperator"`
	Size                int64                 `json:"size"`
}

// GeoMatchStatement represents a geographic match rule statement.
type GeoMatchStatement struct {
	FieldToMatch *FieldToMatch `json:"fieldToMatch"`
	CountryCodes []string      `json:"countryCodes"`
}

// RateBasedStatement represents a rate-based rule statement.
type RateBasedStatement struct {
	Limit                   int64                    `json:"limit"`
	EvaluationWindowSec     int64                    `json:"evaluationWindowSec"`
	AggregateKeyType        string                   `json:"aggregateKeyType"`
	ScopeDownStatement      *Statement               `json:"scopeDownStatement,omitempty"`
	IPSetReferenceStatement *IPSetReferenceStatement `json:"ipSetReferenceStatement,omitempty"`
}

// ManagedRuleGroupStatement represents a managed rule group statement.
type ManagedRuleGroupStatement struct {
	Name          string   `json:"name"`
	VendorName    string   `json:"vendorName"`
	Version       string   `json:"version,omitempty"`
	ExcludedRules []string `json:"excludedRules,omitempty"`
}

// AndStatement represents a logical AND statement.
type AndStatement struct {
	Statements []*Statement `json:"statements"`
}

// OrStatement represents a logical OR statement.
type OrStatement struct {
	Statements []*Statement `json:"statements"`
}

// NotStatement represents a logical NOT statement.
type NotStatement struct {
	Statement *Statement `json:"statement"`
}

// RegexMatchStatement represents a regular expression match statement.
type RegexMatchStatement struct {
	RegexPatternSetID   string                `json:"regexPatternSetId"`
	FieldToMatch        *FieldToMatch         `json:"fieldToMatch"`
	TextTransformations []*TextTransformation `json:"textTransformations"`
}

// IPSetReferenceStatement represents an IP set reference statement.
type IPSetReferenceStatement struct {
	IPSetArn               string                  `json:"ipSetArn"`
	IPSetForwardedIPConfig *IPSetForwardedIPConfig `json:"ipSetForwardedIPConfig,omitempty"`
}

// IPSetForwardedIPConfig holds the configuration for forwarded IP in an IP set.
type IPSetForwardedIPConfig struct {
	HeaderName       string `json:"headerName"`
	Position         string `json:"position"`
	FallbackBehavior string `json:"fallbackBehavior"`
}

// FieldToMatch represents the field to match in a rule statement.
type FieldToMatch struct {
	Type           string `json:"type"`
	SingleQueryArg string `json:"singleQueryArg,omitempty"`
	UriPath        string `json:"uriPath,omitempty"`
	QueryString    string `json:"queryString,omitempty"`
	Body           string `json:"body,omitempty"`
	Method         string `json:"method,omitempty"`
	HeaderOrder    string `json:"headerOrder,omitempty"`
	Host           string `json:"host,omitempty"`
	Cookie         string `json:"cookie,omitempty"`
	AllQueryArgs   string `json:"allQueryArgs,omitempty"`
	Header         string `json:"header,omitempty"`
	IP             string `json:"ip,omitempty"`
}

// TextTransformation represents a text transformation for matching.
type TextTransformation struct {
	Priority int    `json:"priority"`
	Type     string `json:"type"`
}

// Action represents a WAF rule action.
type Action struct {
	Allow     *AllowAction     `json:"allow,omitempty"`
	Block     *BlockAction     `json:"block,omitempty"`
	Count     *CountAction     `json:"count,omitempty"`
	Captcha   *CaptchaAction   `json:"captcha,omitempty"`
	Challenge *ChallengeAction `json:"challenge,omitempty"`
}

// AllowAction represents an allow action.
type AllowAction struct{}

// BlockAction represents a block action.
type BlockAction struct{}

// CountAction represents a count action.
type CountAction struct{}

// CaptchaAction represents a CAPTCHA action.
type CaptchaAction struct{}

// ChallengeAction represents a challenge action.
type ChallengeAction struct{}

// VisibilityConfig holds the visibility configuration for a rule or web ACL.
type VisibilityConfig struct {
	SampledRequestsEnabled   bool   `json:"sampledRequestsEnabled"`
	CloudWatchMetricsEnabled bool   `json:"cloudWatchMetricsEnabled"`
	MetricName               string `json:"metricName"`
}

// RuleGroup represents a WAF rule group.
type RuleGroup struct {
	ID               string            `json:"id"`
	Name             string            `json:"name"`
	ARN              string            `json:"arn"`
	Capacity         int64             `json:"capacity"`
	Rules            []*Rule           `json:"rules"`
	VisibilityConfig *VisibilityConfig `json:"visibilityConfig"`
	Description      string            `json:"description"`
	LockToken        string            `json:"lockToken"`
	Tags             []types.Tag       `json:"tags"`
	CreatedAt        time.Time         `json:"createdAt"`
	ModifiedAt       time.Time         `json:"modifiedAt"`
}

// IPSet represents a WAF IP set.
type IPSet struct {
	ID               string        `json:"id"`
	Name             string        `json:"name"`
	ARN              string        `json:"arn"`
	IPAddressVersion string        `json:"ipAddressVersion"`
	Description      string        `json:"description"`
	Addresses        []string      `json:"addresses"`
	IPSetDescriptors []interface{} `json:"ipSetDescriptors"`
	LockToken        string        `json:"lockToken"`
	Tags             []types.Tag   `json:"tags"`
	CreatedAt        time.Time     `json:"createdAt"`
	ModifiedAt       time.Time     `json:"modifiedAt"`
}

// RegexPatternSet represents a WAF regular expression pattern set.
type RegexPatternSet struct {
	ID                  string      `json:"id"`
	Name                string      `json:"name"`
	ARN                 string      `json:"arn"`
	Description         string      `json:"description"`
	RegularPatterns     []string    `json:"regularPatterns"`
	RegexPatternStrings []string    `json:"regexPatternStrings"`
	LockToken           string      `json:"lockToken"`
	Tags                []types.Tag `json:"tags"`
	CreatedAt           time.Time   `json:"createdAt"`
	ModifiedAt          time.Time   `json:"modifiedAt"`
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
	ResourceArn              string          `json:"resourceArn"`
	LogDestinationConfigs    []string        `json:"logDestinationConfigs"`
	LogScope                 string          `json:"logScope"`
	LogType                  string          `json:"logType"`
	LoggingFilter            *LoggingFilter  `json:"loggingFilter,omitempty"`
	ManagedByFirewallManager bool            `json:"managedByFirewallManager"`
	RedactedFields           []*FieldToMatch `json:"redactedFields,omitempty"`
	CreatedAt                time.Time       `json:"createdAt"`
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
