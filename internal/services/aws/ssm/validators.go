// Package ssm provides input validation for Systems Manager Parameter Store
// operations. Each validator maps a Smithy constraint or AWS API contract to
// a corresponding AWS-shaped error type.
package ssm

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	ssmstore "vorpalstacks/internal/store/aws/ssm"
)

// maxParameterNameListLen is the Smithy ParameterNameList length cap (1-10).
// SSM API ops that take a Name list (GetParameters, DeleteParameters) share it.
const maxParameterNameListLen = 10

// validDataTypes is the AWS-documented DataType value set. Smithy declares
// ParameterDataType as a free string with a length cap, but the public API
// contract restricts it to these values.
var validDataTypes = map[string]struct{}{
	"text":          {},
	"aws:ec2:image": {},
}

// validTiers mirrors the Smithy ParameterTier enum members.
var validTiers = map[ssmstore.ParameterTier]struct{}{
	ssmstore.ParameterTierStandard:           {},
	ssmstore.ParameterTierAdvanced:           {},
	ssmstore.ParameterTierIntelligentTiering: {},
}

// keyIDRegex is the Smithy ParameterKeyId pattern.
var keyIDRegex = regexp.MustCompile(`^[a-zA-Z0-9:/_-]+$`)

// validateParameterNameList enforces the Smithy length cap on operations that
// accept a Name list. AWS returns InvalidParameter for overflow.
func validateParameterNameList(names []string) error {
	if len(names) > maxParameterNameListLen {
		return ErrInvalidParameterName
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

// validateKeyID enforces the Smithy ParameterKeyId pattern. Empty key is
// allowed (uses the AWS-managed default key).
func validateKeyID(keyID string) error {
	if keyID == "" {
		return nil
	}
	if !keyIDRegex.MatchString(keyID) {
		return ErrInvalidParameterValue
	}
	return nil
}

// validateAllowedPattern compiles the AllowedPattern regex at PutParameter
// time. AWS returns InvalidAllowedPatternException for invalid regex.
func validateAllowedPattern(pattern string) error {
	if pattern == "" {
		return nil
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
	switch {
	case maxResults == 0:
		return 10, nil
	case maxResults < 1 || maxResults > 10:
		return 0, ErrInvalidParameterValue
	default:
		return maxResults, nil
	}
}

// validateMaxResultsForPage enforces the DescribeParameters/GetParameterHistory
// Smithy range (1-50, default 50 when zero).
func validateMaxResultsForPage(maxResults int32) (int32, error) {
	switch {
	case maxResults == 0:
		return 50, nil
	case maxResults < 1 || maxResults > 50:
		return 0, ErrInvalidParameterValue
	default:
		return maxResults, nil
	}
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
		if len(label) < 1 || len(label) > 100 {
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
// holding a JSON array of policy entries. We accept the canonical shape and
// reject malformed input; the stored text is replayed verbatim on output.
func validatePolicies(policies string) error {
	if policies == "" {
		return nil
	}
	if len(policies) > 4096 {
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

// ParameterPutFields holds the raw string inputs of a PutParameter call.
// Both the HTTP query-protocol path and the admin gRPC path populate this
// from their respective wire formats and then call normalisePutParameter so
// validation, defaulting and Tier auto-promotion happen once in a single place.
type ParameterPutFields struct {
	Name           string
	Value          string
	Type           string
	Description    string
	KeyID          string
	AllowedPattern string
	DataType       string
	Tier           string
	Policies       string
}

// normalisePutParameter validates every PutParameter input field and returns
// a fully-populated Parameter. Defaults (DataType="text", Tier="Standard") and
// Tier auto-promotion (Standard -> Advanced on >4KB or Policies) are applied
// here so callers cannot diverge. The returned Parameter is ready for the
// store; callers only need to supply the LastModifiedBy before storing.
func normalisePutParameter(in ParameterPutFields) (*ssmstore.Parameter, error) {
	if in.Name == "" {
		return nil, ErrInvalidParameterName
	}

	paramType := ssmstore.ParameterType(in.Type)
	if paramType == "" {
		paramType = ssmstore.ParameterTypeString
	}
	switch paramType {
	case ssmstore.ParameterTypeString, ssmstore.ParameterTypeStringList, ssmstore.ParameterTypeSecureString:
	default:
		return nil, ErrInvalidParameterType
	}

	if err := validateKeyID(in.KeyID); err != nil {
		return nil, err
	}
	if err := validateAllowedPattern(in.AllowedPattern); err != nil {
		return nil, err
	}

	dataType := in.DataType
	if dataType == "" {
		dataType = "text"
	}
	if err := validateDataType(dataType); err != nil {
		return nil, err
	}

	tier := ssmstore.ParameterTier(in.Tier)
	if in.Tier != "" {
		if err := validateTier(tier); err != nil {
			return nil, err
		}
	}

	if err := validatePolicies(in.Policies); err != nil {
		return nil, err
	}

	param := ssmstore.NewParameter(in.Name, in.Value, paramType)
	param.Description = in.Description
	param.KeyID = in.KeyID
	param.AllowedPattern = in.AllowedPattern
	param.DataType = dataType
	param.Tier = tier
	param.Policies = in.Policies

	// AWS auto-promotes Standard-tier parameters to Advanced when the value
	// exceeds the 4KB Standard-tier limit or when any Policies are attached.
	// An empty Tier means the caller omitted it — AWS treats that as
	// Standard-equivalent, so the same promotion rule must apply.
	if param.Tier == "" || param.Tier == ssmstore.ParameterTierStandard {
		if len(param.Value) > 4096 || param.Policies != "" {
			param.Tier = ssmstore.ParameterTierAdvanced
		}
	}

	return param, nil
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
		return []interface{}{}
	}
	out := make([]interface{}, 0, len(entries))
	for _, raw := range entries {
		var meta struct {
			Type string `json:"Type"`
		}
		if err := json.Unmarshal(raw, &meta); err != nil || meta.Type == "" {
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
