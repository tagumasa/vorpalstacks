package wafv2

import (
	"encoding/json"
	"fmt"
	"net"
	"regexp"
	"strings"
	"unicode/utf8"

	"vorpalstacks/internal/services/aws/wafv2/inspection"
	"vorpalstacks/internal/store/aws/waf"

	"golang.org/x/crypto/sha3"
)

// maxRegexPatternStringLength is the Smithy RegexPatternString @length
// maximum, counted in Unicode characters like every @length trait; the
// shape's ".*" pattern admits multibyte regex patterns.
const maxRegexPatternStringLength = 512

// entityNamePattern mirrors the Smithy @pattern trait on
// com.amazonaws.wafv2#EntityName: ^[\w\-]+$ (also the pattern of the
// CustomResponseBodies map keys, whose key shape is EntityName).
var entityNamePattern = regexp.MustCompile(`^[\w\-]+$`)

// entityDescriptionPattern mirrors the Smithy @pattern trait on
// com.amazonaws.wafv2#EntityDescription. The middle class gains \v because
// Go's \s does not include the vertical tab that the Java-flavoured model
// pattern accepts.
var entityDescriptionPattern = regexp.MustCompile(`^[\w+=:#@/\-,.][\w+=:#@/\-,\.\s\v]+[\w+=:#@/\-,.]$`)

// tokenDomainPattern mirrors the Smithy @pattern trait on
// com.amazonaws.wafv2#TokenDomain: ^[\w./-]+$.
var tokenDomainPattern = regexp.MustCompile(`^[\w./-]+$`)

// validateScope checks that the Scope parameter is a valid Smithy enum
// value (CLOUDFRONT or REGIONAL).
func validateScope(scope string) error {
	if scope != "CLOUDFRONT" && scope != "REGIONAL" {
		return invalidParamError("Scope must be CLOUDFRONT or REGIONAL")
	}
	return nil
}

// metricNamePattern mirrors the Smithy @pattern trait on
// com.amazonaws.wafv2#MetricName: ^[\w#:\.\-/]+$
var metricNamePattern = regexp.MustCompile(`^[\w#:\.\-/]+$`)

// validateVisibilityConfig validates a VisibilityConfig against the
// Smithy model: VisibilityConfig itself is required on
// Create/UpdateWebACL and Create/UpdateRuleGroup, and its MetricName is
// required, 1-255 characters, and must match the MetricName pattern.
// The two Boolean members default to false when omitted, which the
// typed struct already represents. A nil VisibilityConfig is a
// contract violation, not an omission to tolerate.
func validateVisibilityConfig(vc *waf.VisibilityConfig) error {
	if vc == nil {
		return invalidParamError("VisibilityConfig is required")
	}
	if len(vc.MetricName) < 1 || len(vc.MetricName) > 255 {
		return invalidParamError("VisibilityConfig MetricName must be between 1 and 255 characters")
	}
	if !metricNamePattern.MatchString(vc.MetricName) {
		return invalidParamError(fmt.Sprintf("VisibilityConfig MetricName contains invalid characters: %s", vc.MetricName))
	}
	return nil
}

// validateDefaultAction checks that the DefaultAction only contains Allow
// or Block. Per the Smithy model, DefaultAction is a separate shape from
// RuleAction, only supports terminating actions (Allow, Block), and is
// required wherever it is validated (Create/UpdateWebACL).
func validateDefaultAction(action *waf.Action) error {
	if action == nil {
		return invalidParamError("DefaultAction is required")
	}
	if action.Allow != nil {
		return nil
	}
	if action.Block != nil {
		return nil
	}
	return invalidParamError("DefaultAction must be Allow or Block")
}

// validateIPAddressVersion checks that IPAddressVersion is explicitly
// provided and is a valid enum value (IPV4 or IPV6). Per the Smithy
// model, IPAddressVersion is @required on CreateIPSetRequest.
func validateIPAddressVersion(v string) error {
	if v == "" {
		return invalidParamError("IPAddressVersion is required")
	}
	if v != "IPV4" && v != "IPV6" {
		return invalidParamError(fmt.Sprintf("Invalid IPAddressVersion: %s (must be IPV4 or IPV6)", v))
	}
	return nil
}

// validateIPAddress validates a single CIDR address against the specified
// IP address version. Uses net.ParseCIDR for syntax verification and
// checks that the IP family (IPv4 vs IPv6) matches IPAddressVersion.
func validateIPAddress(addr, version string) error {
	_, ipNet, err := net.ParseCIDR(addr)
	if err != nil {
		return invalidParamError(fmt.Sprintf("Invalid CIDR notation: %s", addr))
	}
	if version == "IPV4" && ipNet.IP.To4() == nil {
		return invalidParamError(fmt.Sprintf("Address %s is not a valid IPv4 CIDR", addr))
	}
	if version == "IPV6" && ipNet.IP.To4() != nil {
		return invalidParamError(fmt.Sprintf("Address %s is not a valid IPv6 CIDR", addr))
	}
	return nil
}

// validateEntityName validates a resource name against the Smithy
// EntityName shape constraints: length 1-128 and pattern ^[\w\-]+$. The
// @length trait counts Unicode characters; the pattern only admits pure
// ASCII, so both counts agree for every accepted name, but counting
// characters keeps the reported length accurate for rejected input.
func validateEntityName(name string) error {
	if n := utf8.RuneCountInString(name); n < 1 || n > 128 {
		return invalidParamError(fmt.Sprintf("Name length must be between 1 and 128 characters (got %d)", n))
	}
	if !entityNamePattern.MatchString(name) {
		return invalidParamError(fmt.Sprintf("Name must match %s", entityNamePattern.String()))
	}
	return nil
}

// validateEntityDescription validates a resource description against the
// Smithy EntityDescription shape constraints: length 0-256 and the
// EntityDescription pattern. The pattern is only applied to non-empty
// values: the protocol layer cannot distinguish an omitted optional
// member from an explicitly empty one. The @length trait counts Unicode
// characters; the pattern only admits pure ASCII, so both counts agree
// for every accepted description, but counting characters keeps the
// reported length accurate for rejected input.
func validateEntityDescription(desc string) error {
	if n := utf8.RuneCountInString(desc); n > 256 {
		return invalidParamError(fmt.Sprintf("Description length must not exceed 256 characters (got %d)", n))
	}
	if desc != "" && !entityDescriptionPattern.MatchString(desc) {
		return invalidParamError("Description must start and end with an allowed character ([\\w+=:#@/-,.]) and contain at least 3 characters")
	}
	return nil
}

// validateRuleAction validates that a Rule.Action contains exactly one
// valid action type. Per the Smithy model, RuleAction supports six
// member types: Allow, Block, Count, Captcha, Challenge, Monetize.
// Rule.Action is stored as interface{} (it may be *waf.Action after
// conversion or a raw map after JSON deserialisation from storage).
func validateRuleAction(action interface{}) error {
	if action == nil {
		return nil
	}
	switch a := action.(type) {
	case *waf.Action:
		if a == nil {
			return nil
		}
		count := 0
		if a.Allow != nil {
			count++
		}
		if a.Block != nil {
			count++
		}
		if a.Count != nil {
			count++
		}
		if a.Captcha != nil {
			count++
		}
		if a.Challenge != nil {
			count++
		}
		if a.Monetize != nil {
			count++
			if a.Monetize.PriceMultiplier != "" && !priceMultiplierPattern.MatchString(a.Monetize.PriceMultiplier) {
				return invalidParamError(fmt.Sprintf("Invalid Monetize PriceMultiplier: %s (must be a string integer 1 to 100)", a.Monetize.PriceMultiplier))
			}
		}
		if count == 0 {
			return invalidParamError("RuleAction must specify at least one of Allow, Block, Count, Captcha, Challenge, or Monetize")
		}
		if count > 1 {
			return invalidParamError("RuleAction must specify exactly one action type")
		}
	}
	return nil
}

// validateOverrideAction validates that an OverrideAction contains only
// valid member types. Per the Smithy model, OverrideAction supports only
// Count and None. None is represented as an empty Action (all nil).
func validateOverrideAction(action *waf.Action) error {
	if action == nil {
		return nil
	}
	if action.Allow != nil {
		return invalidParamError("OverrideAction does not support Allow")
	}
	if action.Block != nil {
		return invalidParamError("OverrideAction does not support Block")
	}
	if action.Captcha != nil {
		return invalidParamError("OverrideAction does not support Captcha")
	}
	if action.Challenge != nil {
		return invalidParamError("OverrideAction does not support Challenge")
	}
	if action.Monetize != nil {
		return invalidParamError("OverrideAction does not support Monetize")
	}
	return nil
}

// validateLabelMatchScope validates the Scope field of a
// LabelMatchStatement. Per the Smithy model, Scope must be LABEL or
// NAMESPACE.
func validateLabelMatchScope(scope string) error {
	if scope != "LABEL" && scope != "NAMESPACE" {
		return invalidParamError(fmt.Sprintf("Invalid LabelMatchScope: %s (must be LABEL or NAMESPACE)", scope))
	}
	return nil
}

// validateRateLimit validates the Limit field of a RateBasedStatement
// against the Smithy RateLimit shape constraint: range 10-2000000000.
func validateRateLimit(limit int64) error {
	if limit < 10 || limit > 2000000000 {
		return invalidParamError(fmt.Sprintf("RateBasedStatement Limit must be between 10 and 2000000000 (got %d)", limit))
	}
	return nil
}

// validateRateCustomKeys validates the custom aggregation keys of a
// rate-based statement against the model requirements: CUSTOM_KEYS
// must carry at least one aggregation key in CustomKeys, and an IP or
// ForwardedIP custom key cannot be the only key — the model requires
// at least one other key alongside it (to aggregate on only the
// (forwarded) IP address, the AggregateKeyType IP or FORWARDED_IP is
// used instead). Every entry must set exactly one union member.
func validateRateCustomKeys(rate *waf.RateBasedStatement) error {
	if rate.AggregateKeyType != "CUSTOM_KEYS" {
		return nil
	}
	if len(rate.CustomKeys) == 0 {
		return invalidParamError("RateBasedStatement with AggregateKeyType CUSTOM_KEYS must specify at least one custom aggregation key in CustomKeys")
	}
	addressKeys, otherKeys := 0, 0
	for i, custom := range rate.CustomKeys {
		if custom == nil {
			return invalidParamError(fmt.Sprintf("CustomKeys[%d] must set exactly one aggregation key", i))
		}
		set := 0
		addressKey := false
		for _, member := range []struct {
			present bool
			isAddr  bool
		}{
			{custom.Header != nil, false},
			{custom.Cookie != nil, false},
			{custom.QueryArgument != nil, false},
			{custom.QueryString != nil, false},
			{custom.HTTPMethod != nil, false},
			{custom.ForwardedIP != nil, true},
			{custom.IP != nil, true},
			{custom.LabelNamespace != nil, false},
			{custom.UriPath != nil, false},
			{custom.JA3Fingerprint != nil, false},
			{custom.JA4Fingerprint != nil, false},
			{custom.ASN != nil, false},
		} {
			if member.present {
				set++
				addressKey = member.isAddr
			}
		}
		if set != 1 {
			return invalidParamError(fmt.Sprintf("CustomKeys[%d] must set exactly one aggregation key, got %d", i, set))
		}
		if addressKey {
			addressKeys++
		} else {
			otherKeys++
		}
	}
	if addressKeys > 0 && otherKeys == 0 {
		return invalidParamError("CustomKeys cannot aggregate on only the IP or forwarded IP address; use AggregateKeyType IP or FORWARDED_IP instead")
	}
	return nil
}

// validateSizeValue validates the Size field of a SizeConstraintStatement
// against the Smithy Size shape constraint: range 0-21474836480.
func validateSizeValue(size int64) error {
	if size < 0 || size > 21474836480 {
		return invalidParamError(fmt.Sprintf("SizeConstraintStatement Size must be between 0 and 21474836480 (got %d)", size))
	}
	return nil
}

// validateLogDestinationARN validates that a LogDestinationConfigs entry
// is a well-formed ARN for one of the three AWS-supported logging
// destination types: Kinesis Firehose, CloudWatch Logs, or an S3 bucket.
func validateLogDestinationARN(arn string) error {
	for _, prefix := range []string{"arn:aws:firehose:", "arn:aws:logs:", "arn:aws:s3:", "arn:aws-cn:", "arn:aws-us-gov:"} {
		if strings.HasPrefix(arn, prefix) {
			return nil
		}
	}
	return invalidParamError(fmt.Sprintf("LogDestinationConfigs must be a Kinesis Firehose, CloudWatch Logs, or S3 bucket ARN (got %s)", arn))
}

// validateFilterBehavior validates the Behavior field of a logging Filter.
// Per the AWS WAF API, Behavior must be KEEP or DROP.
func validateFilterBehavior(b string) error {
	if b != "KEEP" && b != "DROP" {
		return invalidParamError(fmt.Sprintf("Invalid Filter Behavior: %s (must be KEEP or DROP)", b))
	}
	return nil
}

// validateFilterRequirement validates the Requirement field of a logging
// Filter. Per the AWS WAF API, Requirement must be MEETS_ALL or
// MEETS_ANY.
func validateFilterRequirement(r string) error {
	if r != "MEETS_ALL" && r != "MEETS_ANY" {
		return invalidParamError(fmt.Sprintf("Invalid Filter Requirement: %s (must be MEETS_ALL or MEETS_ANY)", r))
	}
	return nil
}

// validateTokenDomains validates the TokenDomains field of a WebACL.
// Per the Smithy model, TokenDomain has length 1-253, pattern
// ^[\w./-]+$, and a WebACL may specify at most 5 token domains.
func validateTokenDomains(domains interface{}) error {
	if domains == nil {
		return nil
	}
	arr, ok := domains.([]interface{})
	if !ok {
		return invalidParamError("TokenDomains must be an array of strings")
	}
	if len(arr) > 5 {
		return invalidParamError(fmt.Sprintf("TokenDomains must not exceed 5 entries (got %d)", len(arr)))
	}
	for _, d := range arr {
		s, ok := d.(string)
		if !ok {
			return invalidParamError("TokenDomains entries must be strings")
		}
		if len(s) < 1 || len(s) > 253 {
			return invalidParamError(fmt.Sprintf("TokenDomain length must be between 1 and 253 characters (got %d)", len(s)))
		}
		if !tokenDomainPattern.MatchString(s) {
			return invalidParamError(fmt.Sprintf("TokenDomain must match %s (got %q)", tokenDomainPattern.String(), s))
		}
	}
	return nil
}

// validateCustomResponseBodies validates the CustomResponseBodies field
// of a WebACL. Keys follow the EntityName shape: length 1-128 and pattern
// ^[\w\-]+$.
func validateCustomResponseBodies(bodies interface{}) error {
	if bodies == nil {
		return nil
	}
	m, ok := bodies.(map[string]interface{})
	if !ok {
		return invalidParamError("CustomResponseBodies must be a map")
	}
	for k := range m {
		if len(k) < 1 || len(k) > 128 {
			return invalidParamError(fmt.Sprintf("CustomResponseBodies key length must be between 1 and 128 characters (got %d for key %q)", len(k), k))
		}
		if !entityNamePattern.MatchString(k) {
			return invalidParamError(fmt.Sprintf("CustomResponseBodies key must match %s (got %q)", entityNamePattern.String(), k))
		}
	}
	return nil
}

// --- Enum validators ---

func isValidTextTransformationType(t string) bool {
	switch t {
	case "NONE", "COMPRESS_WHITE_SPACE", "HTML_ENTITY_DECODE", "LOWERCASE",
		"CMD_LINE", "URL_DECODE", "BASE64_DECODE", "HEX_DECODE", "MD5",
		"REPLACE_COMMENTS", "ESCAPE_SEQ_DECODE", "SQL_HEX_DECODE",
		"CSS_DECODE", "JS_DECODE", "NORMALIZE_PATH", "NORMALIZE_PATH_WIN",
		"REMOVE_NULLS", "REPLACE_NULLS", "BASE64_DECODE_EXT",
		"URL_DECODE_UNI", "UTF8_TO_UNICODE", "REMOVE_WHITESPACE",
		"TRIM", "TRIM_LEFT", "TRIM_RIGHT", "REMOVE_COMMENTS_CHAR",
		"UPPERCASE", "CMD_LINE_WIN", "CMD_LINE_UNIX", "JS_DECODE_EXT",
		"SHA256":
		return true
	}
	return false
}

func isValidPositionalConstraint(p string) bool {
	switch p {
	case "EXACTLY", "STARTS_WITH", "ENDS_WITH", "CONTAINS", "CONTAINS_WORD":
		return true
	}
	return false
}

func isValidComparisonOperator(op string) bool {
	switch op {
	case "EQ", "NE", "LE", "LT", "GE", "GT":
		return true
	}
	return false
}

func isValidOversizeHandling(h string) bool {
	switch h {
	case "CONTINUE", "MATCH", "NO_MATCH":
		return true
	}
	return false
}

func isValidMapMatchScope(s string) bool {
	switch s {
	case "ALL", "KEY", "VALUE":
		return true
	}
	return false
}

func isValidBodyParsingFallbackBehavior(b string) bool {
	switch b {
	case "MATCH", "NO_MATCH", "EVALUATE_AS_STRING":
		return true
	}
	return false
}

func isValidAggregateKeyType(k string) bool {
	switch k {
	case "IP", "FORWARDED_IP", "CUSTOM_KEYS", "CONSTANT":
		return true
	}
	return false
}

// validateTextTransformations validates an array of TextTransformation.
func validateTextTransformations(tts []*waf.TextTransformation) error {
	for _, tt := range tts {
		if tt == nil {
			continue
		}
		if !isValidTextTransformationType(tt.Type) {
			return invalidParamError(fmt.Sprintf("Invalid TextTransformationType: %s", tt.Type))
		}
	}
	return nil
}

// validateFieldToMatch recursively validates all enum fields inside a
// FieldToMatch.
func validateFieldToMatch(ftm *waf.FieldToMatch) error {
	if ftm == nil {
		return nil
	}
	if ftm.Body != nil && ftm.Body.OversizeHandling != "" && !isValidOversizeHandling(ftm.Body.OversizeHandling) {
		return invalidParamError(fmt.Sprintf("Invalid OversizeHandling: %s", ftm.Body.OversizeHandling))
	}
	if ftm.Cookies != nil {
		if ftm.Cookies.MatchScope != "" && !isValidMapMatchScope(ftm.Cookies.MatchScope) {
			return invalidParamError(fmt.Sprintf("Invalid MatchScope: %s", ftm.Cookies.MatchScope))
		}
		if ftm.Cookies.OversizeHandling != "" && !isValidOversizeHandling(ftm.Cookies.OversizeHandling) {
			return invalidParamError(fmt.Sprintf("Invalid OversizeHandling: %s", ftm.Cookies.OversizeHandling))
		}
	}
	if ftm.Headers != nil {
		if ftm.Headers.MatchScope != "" && !isValidMapMatchScope(ftm.Headers.MatchScope) {
			return invalidParamError(fmt.Sprintf("Invalid MatchScope: %s", ftm.Headers.MatchScope))
		}
		if ftm.Headers.OversizeHandling != "" && !isValidOversizeHandling(ftm.Headers.OversizeHandling) {
			return invalidParamError(fmt.Sprintf("Invalid OversizeHandling: %s", ftm.Headers.OversizeHandling))
		}
	}
	if ftm.HeaderOrder != nil && ftm.HeaderOrder.OversizeHandling != "" && !isValidOversizeHandling(ftm.HeaderOrder.OversizeHandling) {
		return invalidParamError(fmt.Sprintf("Invalid OversizeHandling: %s", ftm.HeaderOrder.OversizeHandling))
	}
	if ftm.JsonBody != nil {
		if ftm.JsonBody.MatchScope != "" && !isValidMapMatchScope(ftm.JsonBody.MatchScope) {
			return invalidParamError(fmt.Sprintf("Invalid JsonMatchScope: %s", ftm.JsonBody.MatchScope))
		}
		if ftm.JsonBody.InvalidFallbackBehavior != "" && !isValidBodyParsingFallbackBehavior(ftm.JsonBody.InvalidFallbackBehavior) {
			return invalidParamError(fmt.Sprintf("Invalid InvalidFallbackBehavior: %s", ftm.JsonBody.InvalidFallbackBehavior))
		}
		if ftm.JsonBody.OversizeHandling != "" && !isValidOversizeHandling(ftm.JsonBody.OversizeHandling) {
			return invalidParamError(fmt.Sprintf("Invalid OversizeHandling: %s", ftm.JsonBody.OversizeHandling))
		}
	}
	return nil
}

// validateStatement recursively validates all enum and required fields in
// a Statement tree. All sixteen Smithy statement types are covered.
// Returns WAFInvalidParameterException for any invalid enum or missing
// required field. nested marks statements reached through a logical or
// scope-down statement: a ManagedRuleGroupStatement is legal only as a
// web ACL rule's own top-level statement, never nested, per the API's
// statement placement rules.
func validateStatement(stmt *waf.Statement, nested bool) error {
	if stmt == nil {
		return nil
	}
	if nested && stmt.ManagedRuleGroupStatement != nil {
		return invalidParamError("A ManagedRuleGroupStatement can only be used as a rule's top-level statement in a web ACL")
	}

	if stmt.ByteMatchStatement != nil {
		if len(stmt.ByteMatchStatement.SearchString) == 0 {
			return invalidParamError("ByteMatchStatement SearchString is required")
		}
		if stmt.ByteMatchStatement.PositionalConstraint != "" && !isValidPositionalConstraint(stmt.ByteMatchStatement.PositionalConstraint) {
			return invalidParamError(fmt.Sprintf("Invalid PositionalConstraint: %s", stmt.ByteMatchStatement.PositionalConstraint))
		}
		if err := validateFieldToMatch(stmt.ByteMatchStatement.FieldToMatch); err != nil {
			return err
		}
		if err := validateTextTransformations(stmt.ByteMatchStatement.TextTransformations); err != nil {
			return err
		}
	}

	if stmt.SqliMatchStatement != nil {
		if err := validateFieldToMatch(stmt.SqliMatchStatement.FieldToMatch); err != nil {
			return err
		}
		if err := validateTextTransformations(stmt.SqliMatchStatement.TextTransformations); err != nil {
			return err
		}
	}

	if stmt.XssMatchStatement != nil {
		if err := validateFieldToMatch(stmt.XssMatchStatement.FieldToMatch); err != nil {
			return err
		}
		if err := validateTextTransformations(stmt.XssMatchStatement.TextTransformations); err != nil {
			return err
		}
	}

	if stmt.SizeConstraintStatement != nil {
		if err := validateSizeValue(stmt.SizeConstraintStatement.Size); err != nil {
			return err
		}
		if stmt.SizeConstraintStatement.ComparisonOperator != "" && !isValidComparisonOperator(stmt.SizeConstraintStatement.ComparisonOperator) {
			return invalidParamError(fmt.Sprintf("Invalid ComparisonOperator: %s", stmt.SizeConstraintStatement.ComparisonOperator))
		}
		if err := validateFieldToMatch(stmt.SizeConstraintStatement.FieldToMatch); err != nil {
			return err
		}
		if err := validateTextTransformations(stmt.SizeConstraintStatement.TextTransformations); err != nil {
			return err
		}
	}

	if stmt.GeoMatchStatement != nil {
		if len(stmt.GeoMatchStatement.CountryCodes) == 0 {
			return invalidParamError("GeoMatchStatement CountryCodes is required")
		}
	}

	if stmt.RegexMatchStatement != nil {
		if stmt.RegexMatchStatement.RegexString == "" {
			return invalidParamError("RegexMatchStatement RegexString is required")
		}
		if n := utf8.RuneCountInString(stmt.RegexMatchStatement.RegexString); n > maxRegexPatternStringLength {
			return invalidParamError(fmt.Sprintf("RegexMatchStatement RegexString must not exceed %d characters (got %d)", maxRegexPatternStringLength, n))
		}
		if err := validateFieldToMatch(stmt.RegexMatchStatement.FieldToMatch); err != nil {
			return err
		}
		if err := validateTextTransformations(stmt.RegexMatchStatement.TextTransformations); err != nil {
			return err
		}
	}

	if stmt.RegexPatternSetRefStatement != nil {
		if stmt.RegexPatternSetRefStatement.ARN == "" {
			return invalidParamError("RegexPatternSetReferenceStatement ARN is required")
		}
		if err := validateFieldToMatch(stmt.RegexPatternSetRefStatement.FieldToMatch); err != nil {
			return err
		}
		if err := validateTextTransformations(stmt.RegexPatternSetRefStatement.TextTransformations); err != nil {
			return err
		}
	}

	if stmt.IPSetReferenceStatement != nil {
		if stmt.IPSetReferenceStatement.ARN == "" {
			return invalidParamError("IPSetReferenceStatement ARN is required")
		}
	}

	if stmt.LabelMatchStatement != nil {
		if err := validateLabelMatchScope(stmt.LabelMatchStatement.Scope); err != nil {
			return err
		}
		if stmt.LabelMatchStatement.Key == "" {
			return invalidParamError("LabelMatchStatement Key is required")
		}
		if len(stmt.LabelMatchStatement.Key) > 1024 {
			return invalidParamError(fmt.Sprintf("LabelMatchStatement Key must not exceed 1024 characters (got %d)", len(stmt.LabelMatchStatement.Key)))
		}
	}

	if stmt.AsnMatchStatement != nil {
		if len(stmt.AsnMatchStatement.AsnList) == 0 {
			return invalidParamError("AsnMatchStatement AsnList is required")
		}
	}

	if stmt.AndStatement != nil {
		if len(stmt.AndStatement.Statements) < 1 {
			return invalidParamError("AndStatement must contain at least one Statement")
		}
		for _, s := range stmt.AndStatement.Statements {
			if err := validateStatement(s, true); err != nil {
				return err
			}
		}
	}

	if stmt.OrStatement != nil {
		if len(stmt.OrStatement.Statements) < 1 {
			return invalidParamError("OrStatement must contain at least one Statement")
		}
		for _, s := range stmt.OrStatement.Statements {
			if err := validateStatement(s, true); err != nil {
				return err
			}
		}
	}

	if stmt.NotStatement != nil {
		if err := validateStatement(stmt.NotStatement.Statement, true); err != nil {
			return err
		}
	}

	if stmt.RateBasedStatement != nil {
		if err := validateRateLimit(stmt.RateBasedStatement.Limit); err != nil {
			return err
		}
		if stmt.RateBasedStatement.AggregateKeyType != "" && !isValidAggregateKeyType(stmt.RateBasedStatement.AggregateKeyType) {
			return invalidParamError(fmt.Sprintf("Invalid AggregateKeyType: %s", stmt.RateBasedStatement.AggregateKeyType))
		}
		if err := validateRateCustomKeys(stmt.RateBasedStatement); err != nil {
			return err
		}
		if err := validateStatement(stmt.RateBasedStatement.ScopeDownStatement, true); err != nil {
			return err
		}
	}

	if stmt.ManagedRuleGroupStatement != nil {
		if stmt.ManagedRuleGroupStatement.Name == "" {
			return invalidParamError("ManagedRuleGroupStatement Name is required")
		}
		if stmt.ManagedRuleGroupStatement.VendorName == "" {
			return invalidParamError("ManagedRuleGroupStatement VendorName is required")
		}
		if err := validateStatement(stmt.ManagedRuleGroupStatement.ScopeDownStatement, true); err != nil {
			return err
		}
	}

	if stmt.RuleGroupReferenceStatement != nil {
		if stmt.RuleGroupReferenceStatement.ARN == "" {
			return invalidParamError("RuleGroupReferenceStatement ARN is required")
		}
	}

	return nil
}

// validateRules validates all statements and action/override fields in a
// rule list. This is the single entry point for rule validation called by
// parseRules before any rule is persisted to the store. Web ACL rules may
// reference a managed rule group as their top-level statement; rule
// group rules may not — the API forbids using a managed rule group
// inside another rule group.
func validateRules(rules []*waf.Rule, allowManagedRuleGroups bool) error {
	for _, r := range rules {
		if r == nil {
			continue
		}
		if err := validateRuleAction(r.Action); err != nil {
			return err
		}
		if err := validateOverrideAction(r.OverrideAction); err != nil {
			return err
		}
		if err := validateImmunityConfig(r.CaptchaConfig, "Captcha"); err != nil {
			return err
		}
		if err := validateImmunityConfig(r.ChallengeConfig, "Challenge"); err != nil {
			return err
		}
		if !allowManagedRuleGroups && statementReferencesManagedRuleGroup(r.Statement) {
			return invalidParamError(fmt.Sprintf("Rule %s: a ManagedRuleGroupStatement cannot be used inside a rule group", r.Name))
		}
		if err := validateStatement(r.Statement, false); err != nil {
			return err
		}
	}
	return nil
}

// statementReferencesManagedRuleGroup reports whether a statement tree
// contains a managed rule group reference at any depth.
func statementReferencesManagedRuleGroup(stmt *waf.Statement) bool {
	if stmt == nil {
		return false
	}
	if stmt.ManagedRuleGroupStatement != nil {
		return true
	}
	if stmt.AndStatement != nil {
		for _, sub := range stmt.AndStatement.Statements {
			if statementReferencesManagedRuleGroup(sub) {
				return true
			}
		}
	}
	if stmt.OrStatement != nil {
		for _, sub := range stmt.OrStatement.Statements {
			if statementReferencesManagedRuleGroup(sub) {
				return true
			}
		}
	}
	if stmt.NotStatement != nil && statementReferencesManagedRuleGroup(stmt.NotStatement.Statement) {
		return true
	}
	if stmt.RateBasedStatement != nil && statementReferencesManagedRuleGroup(stmt.RateBasedStatement.ScopeDownStatement) {
		return true
	}
	return false
}

// The ImmunityTime setting's bounds come from the Smithy
// TimeWindowSecond range, 60 to 259200 seconds.
const (
	immunityTimeMinSeconds = 60
	immunityTimeMaxSeconds = 259200
)

// validateImmunityConfig validates one raw CaptchaConfig or
// ChallengeConfig: ImmunityTime is required inside a present
// ImmunityTimeProperty, must stay within the TimeWindowSecond range,
// and the Challenge action documents a minimum of 300 seconds.
func validateImmunityConfig(raw interface{}, action string) error {
	if raw == nil {
		return nil
	}
	var config struct {
		ImmunityTimeProperty *struct {
			ImmunityTime int64 `json:"ImmunityTime"`
		} `json:"ImmunityTimeProperty"`
	}
	data, err := json.Marshal(raw)
	if err != nil {
		return invalidParamError("immunity configuration is not serialisable")
	}
	if err := json.Unmarshal(data, &config); err != nil {
		return invalidParamError(fmt.Sprintf("%sConfig must be an object", action))
	}
	if config.ImmunityTimeProperty == nil {
		return nil
	}
	immunity := config.ImmunityTimeProperty.ImmunityTime
	if immunity == 0 {
		return invalidParamError("ImmunityTimeProperty ImmunityTime is required")
	}
	if immunity < immunityTimeMinSeconds || immunity > immunityTimeMaxSeconds {
		return invalidParamError(fmt.Sprintf("ImmunityTime must be between %d and %d seconds (got %d)", immunityTimeMinSeconds, immunityTimeMaxSeconds, immunity))
	}
	if action == "Challenge" && immunity < waf.ChallengeImmunityTimeMin {
		return invalidParamError(fmt.Sprintf("The minimum Challenge immunity time is %d seconds (got %d)", waf.ChallengeImmunityTimeMin, immunity))
	}
	return nil
}

// priceMultiplierPattern mirrors the Smithy pattern trait on
// PriceMultiplier: the string form of an integer 1 to 100.
var priceMultiplierPattern = regexp.MustCompile(`^([1-9][0-9]?|100)$`)

// managedProductChains is the Smithy BlockchainChain enum.
var managedProductChains = map[string]bool{
	"BASE":          true,
	"SOLANA":        true,
	"BASE_SEPOLIA":  true,
	"SOLANA_DEVNET": true,
}

// testProductChains are the BlockchainChain enum's test networks; a
// MonetizationConfig's payment networks must be all production or all
// test.
var testProductChains = map[string]bool{
	"BASE_SEPOLIA":  true,
	"SOLANA_DEVNET": true,
}

// base58Alphabet is the Base58 alphabet Solana wallet addresses use.
const base58Alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

// validateMonetizationConfig validates the raw MonetizationConfig of a
// web ACL or rule group: one or two payment networks, each with a valid
// chain, a checksum-valid wallet address of the chain's format and
// exactly one USDC price inside the documented amount bounds, with all
// networks in the same production or test environment.
func validateMonetizationConfig(raw interface{}) error {
	if raw == nil {
		return nil
	}
	var config struct {
		CryptoConfig *struct {
			PaymentNetworks []struct {
				Chain         string `json:"Chain"`
				WalletAddress string `json:"WalletAddress"`
				Prices        []struct {
					Amount   string `json:"Amount"`
					Currency string `json:"Currency"`
				} `json:"Prices"`
			} `json:"PaymentNetworks"`
		} `json:"CryptoConfig"`
		CurrencyMode string `json:"CurrencyMode"`
	}
	data, err := json.Marshal(raw)
	if err != nil {
		return invalidParamError("MonetizationConfig is not serialisable")
	}
	if err := json.Unmarshal(data, &config); err != nil {
		return invalidParamError("MonetizationConfig must be an object")
	}
	if config.CurrencyMode != "" && config.CurrencyMode != "REAL" && config.CurrencyMode != "TEST" {
		return invalidParamError(fmt.Sprintf("Invalid CurrencyMode: %s (must be REAL or TEST)", config.CurrencyMode))
	}
	if config.CryptoConfig == nil {
		return nil
	}
	networks := config.CryptoConfig.PaymentNetworks
	if len(networks) < 1 || len(networks) > 2 {
		return invalidParamError("CryptoConfig PaymentNetworks must contain 1 to 2 networks")
	}
	sawTest := false
	sawProduction := false
	for i, network := range networks {
		if !managedProductChains[network.Chain] {
			return invalidParamError(fmt.Sprintf("PaymentNetwork %d: Invalid Chain: %s", i, network.Chain))
		}
		if testProductChains[network.Chain] {
			sawTest = true
		} else {
			sawProduction = true
		}
		if !validWalletAddress(network.Chain, network.WalletAddress) {
			return invalidParamError(fmt.Sprintf("PaymentNetwork %d: Invalid WalletAddress for chain %s", i, network.Chain))
		}
		if len(network.Prices) != 1 {
			return invalidParamError(fmt.Sprintf("PaymentNetwork %d: Prices must contain exactly one entry", i))
		}
		price := network.Prices[0]
		if price.Currency != "USDC" {
			return invalidParamError(fmt.Sprintf("PaymentNetwork %d: Currency must be USDC (got %s)", i, price.Currency))
		}
		millis, err := inspection.ParsePriceMillis(price.Amount)
		if err != nil || millis < waf.PriceAmountMinMillis || millis > waf.PriceAmountMaxMillis {
			return invalidParamError(fmt.Sprintf("PaymentNetwork %d: Amount must be a decimal between 0.001 and 999999999.999 with at most three decimal places (got %s)", i, price.Amount))
		}
	}
	if sawTest && sawProduction {
		return invalidParamError("CryptoConfig PaymentNetworks must all be production networks or all test networks")
	}
	return nil
}

// validWalletAddress reports whether the address has the format of its
// chain: EVM chains take a 0x-prefixed 20-byte hex address whose EIP-55
// checksum must hold whenever the letters mix case, and Solana chains
// take a 32 to 44 character Base58 public key.
func validWalletAddress(chain, address string) bool {
	switch chain {
	case "BASE", "BASE_SEPOLIA":
		return validEVMAddress(address)
	case "SOLANA", "SOLANA_DEVNET":
		if len(address) < 32 || len(address) > 44 {
			return false
		}
		for i := 0; i < len(address); i++ {
			if !strings.ContainsRune(base58Alphabet, rune(address[i])) {
				return false
			}
		}
		return true
	}
	return false
}

// validEVMAddress checks the 0x-prefixed 40-hex-digit form and, for
// mixed-case addresses, the EIP-55 checksum. An address set entirely in
// one case skips the checksum, the bypass the WalletAddress
// documentation describes.
func validEVMAddress(address string) bool {
	if len(address) != 42 || address[0] != '0' || address[1] != 'x' {
		return false
	}
	lower := strings.ToLower(address)
	for i := 2; i < len(lower); i++ {
		c := lower[i]
		if (c < '0' || c > '9') && (c < 'a' || c > 'f') {
			return false
		}
	}
	if address == lower || address == "0x"+strings.ToUpper(address[2:]) {
		return true
	}
	digest := sha3.NewLegacyKeccak256()
	digest.Write([]byte(lower[2:]))
	hash := digest.Sum(nil)
	for i := 2; i < len(address); i++ {
		c := address[i]
		if c < 'a' || c > 'f' {
			continue
		}
		nibble := hash[(i-2)/2]
		if (i-2)%2 == 0 {
			nibble >>= 4
		} else {
			nibble &= 0x0f
		}
		uppercase := c >= 'A' && c <= 'F'
		if (nibble >= 8) != uppercase {
			return false
		}
	}
	return true
}

// validateMonetizeRules enforces the Monetize action's configuration
// constraints across a rule list: the action requires a
// MonetizationConfig on the owning web ACL or rule group, is available
// only for the CloudFront scope, and cannot be used by rate-based
// rules.
func validateMonetizeRules(rules []*waf.Rule, scope string, monetizationConfig interface{}) error {
	if !rulesContainMonetize(rules) {
		return nil
	}
	if monetizationConfig == nil {
		return invalidParamError("The Monetize action requires a MonetizationConfig on the web ACL or rule group")
	}
	if !isGlobalScope(scope) {
		return invalidParamError("The Monetize action is available only for web ACLs associated with Amazon CloudFront distributions (Scope CLOUDFRONT)")
	}
	for _, rule := range rules {
		if rule == nil {
			continue
		}
		action := typedRuleAction(rule.Action)
		if action == nil || action.Monetize == nil {
			continue
		}
		if rule.Statement != nil && rule.Statement.RateBasedStatement != nil {
			return invalidParamError(fmt.Sprintf("Rule %s: the Monetize action cannot be used for rate-based rules", rule.Name))
		}
	}
	return nil
}

// rulesContainMonetize reports whether any rule carries the Monetize
// action.
func rulesContainMonetize(rules []*waf.Rule) bool {
	for _, rule := range rules {
		if rule == nil {
			continue
		}
		if action := typedRuleAction(rule.Action); action != nil && action.Monetize != nil {
			return true
		}
	}
	return false
}

// typedRuleAction converts a stored rule action — a typed pointer
// before storage, a raw map after a JSON round-trip — to its typed
// form.
func typedRuleAction(action interface{}) *waf.Action {
	if action == nil {
		return nil
	}
	if typed, ok := action.(*waf.Action); ok {
		return typed
	}
	data, err := json.Marshal(action)
	if err != nil {
		return nil
	}
	var out waf.Action
	if err := json.Unmarshal(data, &out); err != nil {
		return nil
	}
	return &out
}
