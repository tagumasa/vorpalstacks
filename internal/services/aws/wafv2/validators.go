package wafv2

import (
	"fmt"
	"net"
	"strings"

	"vorpalstacks/internal/store/aws/waf"
)

// validateScope checks that the Scope parameter is a valid Smithy enum
// value (CLOUDFRONT or REGIONAL).
func validateScope(scope string) error {
	if scope != "CLOUDFRONT" && scope != "REGIONAL" {
		return invalidParamError("Scope must be CLOUDFRONT or REGIONAL")
	}
	return nil
}

// validateDefaultAction checks that the DefaultAction only contains Allow
// or Block. Per the Smithy model, DefaultAction is a separate shape from
// RuleAction and only supports terminating actions (Allow, Block).
func validateDefaultAction(action *waf.Action) error {
	if action == nil {
		return nil
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
// EntityName shape constraint: length 1-128.
func validateEntityName(name string) error {
	if len(name) < 1 || len(name) > 128 {
		return invalidParamError(fmt.Sprintf("Name length must be between 1 and 128 characters (got %d)", len(name)))
	}
	return nil
}

// validateEntityDescription validates a resource description against the
// Smithy EntityDescription shape constraint: length 0-256.
func validateEntityDescription(desc string) error {
	if len(desc) > 256 {
		return invalidParamError(fmt.Sprintf("Description length must not exceed 256 characters (got %d)", len(desc)))
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
// Per the Smithy model, TokenDomain has length 1-253 and a WebACL may
// specify at most 5 token domains.
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
	}
	return nil
}

// validateCustomResponseBodies validates the CustomResponseBodies field
// of a WebACL. Keys must be 1-128 characters.
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
		"URL_DECODE_UNI", "UTF8_TO_UNICODE":
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
// required field.
func validateStatement(stmt *waf.Statement) error {
	if stmt == nil {
		return nil
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
		if len(stmt.RegexMatchStatement.RegexString) > 512 {
			return invalidParamError(fmt.Sprintf("RegexMatchStatement RegexString must not exceed 512 characters (got %d)", len(stmt.RegexMatchStatement.RegexString)))
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
			if err := validateStatement(s); err != nil {
				return err
			}
		}
	}

	if stmt.OrStatement != nil {
		if len(stmt.OrStatement.Statements) < 1 {
			return invalidParamError("OrStatement must contain at least one Statement")
		}
		for _, s := range stmt.OrStatement.Statements {
			if err := validateStatement(s); err != nil {
				return err
			}
		}
	}

	if stmt.NotStatement != nil {
		if err := validateStatement(stmt.NotStatement.Statement); err != nil {
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
		if err := validateStatement(stmt.RateBasedStatement.ScopeDownStatement); err != nil {
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
		if err := validateStatement(stmt.ManagedRuleGroupStatement.ScopeDownStatement); err != nil {
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
// parseRules before any rule is persisted to the store.
func validateRules(rules []*waf.Rule) error {
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
		if err := validateStatement(r.Statement); err != nil {
			return err
		}
	}
	return nil
}
