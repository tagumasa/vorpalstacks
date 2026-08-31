// Package ssm provides input validation for Systems Manager Parameter Store
// operations. Each validator maps a Smithy constraint or AWS API contract to
// a corresponding AWS-shaped error type.
package ssm

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/core/logs"
	ssmstore "vorpalstacks/internal/store/aws/ssm"
)

// maxParameterNameListLen is the Smithy ParameterNameList length cap (1-10).
// SSM API ops that take a Name list (GetParameters, DeleteParameters) share it.
const maxParameterNameListLen = 10

// validDataTypes is the AWS-documented DataType value set. Smithy declares
// ParameterDataType as a free string with a length cap, but the public API
// contract restricts it to these values (the enum membership check also
// subsumes the 128-character Smithy length cap).
var validDataTypes = map[string]struct{}{
	"text":                {},
	"aws:ec2:image":       {},
	"aws:ssm:integration": {},
}

// validTiers mirrors the Smithy ParameterTier enum members.
var validTiers = map[ssmstore.ParameterTier]struct{}{
	ssmstore.ParameterTierStandard:           {},
	ssmstore.ParameterTierAdvanced:           {},
	ssmstore.ParameterTierIntelligentTiering: {},
}

// keyIDRegex is the Smithy ParameterKeyId pattern.
var keyIDRegex = regexp.MustCompile(`^[a-zA-Z0-9:/_-]+$`)

// validateParameterNameList enforces the Smithy ParameterNameList contract
// (required, @length 1-10) on operations that accept a Name list. AWS
// returns ValidationException for a missing, empty, or oversized list.
func validateParameterNameList(names []string) error {
	if len(names) == 0 || len(names) > maxParameterNameListLen {
		return ErrValidationException
	}
	return nil
}

// validateDataType enforces the AWS DataType value set. AWS returns
// ValidationException with ParameterValidationException for unknown values.
func validateDataType(dataType string) error {
	if dataType == "" {
		return nil
	}
	if _, ok := validDataTypes[dataType]; !ok {
		return ErrInvalidParameterValue
	}
	return nil
}

// validateTier enforces the Smithy ParameterTier enum. AWS returns
// ValidationException for unknown tiers.
func validateTier(tier ssmstore.ParameterTier) error {
	if tier == "" {
		return nil
	}
	if _, ok := validTiers[tier]; !ok {
		return ErrInvalidParameterValue
	}
	return nil
}

// validateKeyID enforces the Smithy ParameterKeyId pattern and length.
// Empty key is allowed (uses the AWS-managed default key).
func validateKeyID(keyID string) error {
	if keyID == "" {
		return nil
	}
	if len(keyID) > ssmstore.MaxParameterKeyIdLength {
		return ErrInvalidParameterValue
	}
	if !keyIDRegex.MatchString(keyID) {
		return ErrInvalidParameterValue
	}
	return nil
}

// validateAllowedPattern compiles the AllowedPattern regex at PutParameter
// time and enforces the Smithy length cap, counted in Unicode characters
// (the shape carries no pattern trait of its own, and regex sources may
// contain multibyte character classes). AWS returns
// InvalidAllowedPatternException for an invalid regex.
func validateAllowedPattern(pattern string) error {
	if pattern == "" {
		return nil
	}
	if utf8.RuneCountInString(pattern) > ssmstore.MaxAllowedPatternLength {
		return ErrInvalidParameterValue
	}
	if _, err := regexp.Compile(pattern); err != nil {
		return ErrInvalidAllowedPattern
	}
	return nil
}

// enforceAllowedPattern was moved to store/aws/ssm/store.go — the store is the
// single point of Value-vs-Pattern enforcement for both new and overwrite
// paths. The service-layer copy was dead code: callers cannot invoke it
// because PutParameter already routes through the store.

// validateHierarchyPath enforces the AWS contract that hierarchy paths start
// with '/'. Used by GetParametersByPath.
func validateHierarchyPath(path string) error {
	if path == "" {
		return ErrInvalidParameterName
	}
	if !strings.HasPrefix(path, "/") {
		return fmt.Errorf("%w: hierarchy paths must start with '/'", ErrInvalidParameterName)
	}
	return nil
}

// validateMaxResultsForPath enforces the GetParametersByPath Smithy range
// (1-10, default 10 when zero).
func validateMaxResultsForPath(maxResults int32) (int32, error) {
	v, err := pagination.ResolveMaxItems(int(maxResults), 10, 1, 10,
		func(int) error { return ErrInvalidParameterValue })
	return int32(v), err
}

// validateMaxResultsForPage enforces the DescribeParameters/GetParameterHistory
// Smithy range (1-50, default 50 when zero).
func validateMaxResultsForPage(maxResults int32) (int32, error) {
	v, err := pagination.ResolveMaxItems(int(maxResults), ssmstore.MaxPageResults, 1, ssmstore.MaxPageResults,
		func(int) error { return ErrInvalidParameterValue })
	return int32(v), err
}

// validateLabels enforces the Smithy ParameterLabel / ParameterLabelList
// constraints on LabelParameterVersion / UnlabelParameterVersion input.
// Length must be 1-100, the request list 1-10, and labels must not start
// with a digit (numeric prefixes are version selectors, not labels).
func validateLabels(labels []string) error {
	if len(labels) == 0 {
		return ErrInvalidParameterLabel
	}
	if len(labels) > 10 {
		return ErrInvalidParameterLabel
	}
	for _, label := range labels {
		// ParameterLabel @length(1,100) counts Unicode characters (the
		// shape carries no pattern).
		if n := utf8.RuneCountInString(label); n < 1 || n > 100 {
			return ErrInvalidParameterLabel
		}
		if label[0] >= '0' && label[0] <= '9' {
			return ErrInvalidParameterLabel
		}
	}
	return nil
}

// validPolicyTypes enumerates the Smithy-defined ParameterPolicyType values.
// Any entry outside this set is rejected at PutParameter time.
var validPolicyTypes = map[string]struct{}{
	"Expiration":                    {},
	"ExpirationNotification":        {},
	"NoNotificationAfterExpiration": {},
}

// policyRequiresVersion lists the policy types that must carry a Version
// field. NoNotificationAfterExpiration has no Version.
var policyRequiresVersion = map[string]struct{}{
	"Expiration":             {},
	"ExpirationNotification": {},
}

// validatePolicies parses the PutParameter Policies JSON document at write
// time. The Smithy ParameterPolicies member is a string of length 1-4096
// counted in Unicode characters (no pattern), holding a JSON array of
// policy entries. We accept the canonical shape and reject malformed input;
// the stored text is replayed verbatim on output.
func validatePolicies(policies string) error {
	if policies == "" {
		return nil
	}
	if utf8.RuneCountInString(policies) > 4096 {
		return ErrInvalidParameterValue
	}
	var probe []map[string]string
	if err := json.Unmarshal([]byte(policies), &probe); err != nil {
		return ErrInvalidParameterValue
	}
	for _, p := range probe {
		policyType := p["Type"]
		if _, ok := validPolicyTypes[policyType]; !ok {
			return ErrInvalidParameterValue
		}
		if _, ok := policyRequiresVersion[policyType]; ok && p["Version"] == "" {
			return ErrInvalidParameterValue
		}
	}
	return nil
}

// policiesToResponse converts a stored Policies JSON document into the
// AWS ParameterInlinePolicy output shape: a list of {PolicyText, PolicyType,
// PolicyStatus}. PolicyText is the JSON serialization of each entry
// (preserving the full policy document including Version and Attributes).
// PolicyStatus is always "Finished" — vorpalstacks does not run policy
// timers, so a stored policy is treated as already-applied for
// output-shape compatibility.
func policiesToResponse(policies string) []interface{} {
	if policies == "" {
		return []interface{}{}
	}
	var entries []json.RawMessage
	if err := json.Unmarshal([]byte(policies), &entries); err != nil {
		// Stored policies are validated at write time; a parse failure here
		// means the stored document no longer matches its validated shape.
		// Surface the anomaly in the log rather than silently dropping it.
		logs.Warn("Stored parameter policies failed to parse", logs.String("error", err.Error()))
		return []interface{}{}
	}
	out := make([]interface{}, 0, len(entries))
	for _, raw := range entries {
		var meta struct {
			Type string `json:"Type"`
		}
		if err := json.Unmarshal(raw, &meta); err != nil || meta.Type == "" {
			logs.Warn("Stored parameter policy entry missing Type; skipping entry")
			continue
		}
		out = append(out, map[string]interface{}{
			"PolicyText":   string(raw),
			"PolicyType":   meta.Type,
			"PolicyStatus": "Finished",
		})
	}
	return out
}
