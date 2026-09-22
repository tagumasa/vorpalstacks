package cloudwatchlogs

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"unicode/utf8"

	tagutil "vorpalstacks/internal/common/tags"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// The policy-document and tag validators: the document forms' JSON
// and structure rules, the selection-criteria shapes, and the TagKey and
// TagValue traits.

// --- Policy document JSON validators ---

// validatePolicyDocumentJSON validates that a policy document is non-empty,
// within the Smithy length limit, and parseable as valid JSON.
func validatePolicyDocumentJSON(doc string) error {
	if err := validatePolicyDocument(doc); err != nil {
		return err
	}
	if !json.Valid([]byte(doc)) {
		return NewLogsError("InvalidParameterException",
			"Policy document must be valid JSON", 400)
	}
	return nil
}

// validateAccessPolicyJSON validates that an access policy is non-empty,
// parseable as valid JSON, and within the operation's byte cap ("This can
// be up to 5120 bytes").
func validateAccessPolicyJSON(policy string) error {
	if err := validateAccessPolicy(policy); err != nil {
		return err
	}
	if len(policy) > logsstore.MaxAccessPolicyBytes {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("Access policy must be at most %d bytes; the request carries %d",
				logsstore.MaxAccessPolicyBytes, len(policy)), 400)
	}
	if !json.Valid([]byte(policy)) {
		return NewLogsError("InvalidParameterException",
			"Access policy must be valid JSON", 400)
	}
	return nil
}

// --- Reserved name prefix validators ---

// validatePolicyNamePrefix rejects policy names that start with the reserved
// AWS prefix. AWS documentation states that policy names must not begin with
// "aws/" or "AWS:".
func validatePolicyNamePrefix(name string) error {
	lower := strings.ToLower(name)
	if strings.HasPrefix(lower, "aws/") || strings.HasPrefix(lower, "aws:") {
		return NewLogsError("InvalidParameterException",
			"Policy names starting with 'aws/' are reserved and not allowed", 400)
	}
	return nil
}

// --- Selection criteria length validator (25 KB max) ---

const maxSelectionCriteriaBytes = 25 * 1024

// validateSelectionCriteria checks that the selection criteria string does
// not exceed the AWS-documented maximum of 25 KB.
func validateSelectionCriteria(sc string) error {
	if len(sc) > maxSelectionCriteriaBytes {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("SelectionCriteria must not exceed %d bytes", maxSelectionCriteriaBytes), 400)
	}
	return nil
}

// validateSubscriptionSelectionCriteria enforces the only supported form
// for a SUBSCRIPTION_FILTER_POLICY: "If policyType is
// SUBSCRIPTION_FILTER_POLICY, the only supported selectionCriteria filter
// is LogGroupName NOT IN []" — a (possibly empty) quoted log-group-name
// exclusion list.
func validateSubscriptionSelectionCriteria(sc string) error {
	if sc == "" {
		return nil
	}
	s := strings.TrimSpace(sc)
	const head = "LogGroupName NOT IN "
	if !strings.HasPrefix(s, head) {
		return NewLogsError("InvalidParameterException",
			"The only supported selectionCriteria for a SUBSCRIPTION_FILTER_POLICY is LogGroupName NOT IN []", 400)
	}
	list := strings.TrimSpace(strings.TrimPrefix(s, head))
	if !strings.HasPrefix(list, "[") || !strings.HasSuffix(list, "]") {
		return NewLogsError("InvalidParameterException",
			"The LogGroupName NOT IN selection must carry a bracketed list of quoted log group names", 400)
	}
	inner := strings.TrimSpace(list[1 : len(list)-1])
	if inner == "" {
		return nil
	}
	for _, part := range strings.Split(inner, ",") {
		name := strings.TrimSpace(part)
		if len(name) < 2 || name[0] != '"' || name[len(name)-1] != '"' {
			return NewLogsError("InvalidParameterException",
				"The LogGroupName NOT IN selection must carry a bracketed list of quoted log group names", 400)
		}
		if err := validateLogGroupName(name[1 : len(name)-1]); err != nil {
			return err
		}
	}
	return nil
}

// --- Tag member validators (the TagKey / TagValue shapes) ---

var (
	// tagKeyCWLPattern is the Smithy pattern trait on TagKey.
	tagKeyCWLPattern = regexp.MustCompile(`^([\p{L}\p{Z}\p{N}_.:/=+\-@]+)$`)
	// tagValueCWLPattern is the Smithy pattern trait on TagValue (the
	// empty value is a valid entry).
	tagValueCWLPattern = regexp.MustCompile(`^([\p{L}\p{Z}\p{N}_.:/=+\-@]*)$`)
)

// tagEntryViolation walks the tag map in sorted key order and returns the
// formatted trait violation of the first offending entry — key length or
// character set, the reserved aws: prefix, then value length or character
// set — or the empty string when every entry satisfies the TagKey and
// TagValue traits. The sorted walk keeps the reported offender
// deterministic on maps holding several violators. The log plane and the
// vended-delivery family share this trait check and raise its findings
// under their own error identities.
func tagEntryViolation(tags map[string]string) string {
	keys := make([]string, 0, len(tags))
	for k := range tags {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := tags[k]
		if utf8.RuneCountInString(k) > tagutil.MaxTagKeyLength || !tagKeyCWLPattern.MatchString(k) {
			return fmt.Sprintf("Invalid tag key: %s. Must be 1-%d characters matching the allowed character set",
				k, tagutil.MaxTagKeyLength)
		}
		// "The aws: prefix is reserved for AWS use. If a tag has a tag key
		// with this prefix, then you can't edit or delete the tag's key or
		// value" (AWS tagging restrictions) — the pattern trait admits the
		// prefix, so the reservation rejects here: a customer write can
		// never introduce a key the reservation would make uneditable.
		if strings.HasPrefix(k, "aws:") {
			return fmt.Sprintf("Tag keys can't start with the reserved prefix aws: (%s)", k)
		}
		if utf8.RuneCountInString(v) > tagutil.MaxTagValueLength || !tagValueCWLPattern.MatchString(v) {
			return fmt.Sprintf("Invalid tag value for key %s: must be 0-%d characters matching the allowed character set",
				k, tagutil.MaxTagValueLength)
		}
	}
	return ""
}

// validateTagEntries enforces the model's TagKey and TagValue traits on
// every entry of a tag map under the log plane's error identity: key
// 1-128 characters matching the key pattern, value 0-256 characters
// matching the value pattern. The Tags map's own 1-50 length is enforced
// by the callers, which choose the error identity their operation
// declares (TooManyTagsException on TagResource alone).
func validateTagEntries(tags map[string]string) error {
	if violation := tagEntryViolation(tags); violation != "" {
		return NewLogsError("InvalidParameterException", violation, 400)
	}
	return nil
}

// validateTagKeyElements enforces the TagKey traits on every element of an
// untag key list.
func validateTagKeyElements(keys []string) error {
	for _, k := range keys {
		if utf8.RuneCountInString(k) > tagutil.MaxTagKeyLength || utf8.RuneCountInString(k) < 1 || !tagKeyCWLPattern.MatchString(k) {
			return NewLogsError("InvalidParameterException",
				fmt.Sprintf("Invalid tag key: %s. Must be 1-%d characters matching the allowed character set", k, tagutil.MaxTagKeyLength), 400)
		}
	}
	return nil
}

// dppDocument is the typed form of a data protection policy document: the
// two Statement blocks the member documentation mandates, over the
// optional Name/Description/Version metadata fields.
type dppDocument struct {
	Name        string         `json:"Name"`
	Description string         `json:"Description"`
	Version     string         `json:"Version"`
	Statement   []dppStatement `json:"Statement"`
}

type dppStatement struct {
	// DataIdentifer carries the documentation's own spelling of the key.
	DataIdentifer []string                   `json:"DataIdentifer"`
	Operation     map[string]json.RawMessage `json:"Operation"`
}

// validateDataProtectionPolicyDocument enforces the structure the
// policyDocument member documentation mandates: "This policy must include
// two JSON blocks: The first block must include both a DataIdentifer
// array and an Operation property with an Audit action... This Audit
// action must contain a FindingsDestination object... The second block
// must include both a DataIdentifer array and an Operation property with
// an Deidentify action... it must contain the "MaskConfig": {} object.
// The "MaskConfig": {} object must be empty... The contents of the two
// DataIdentifer arrays must match exactly", within "up to 30,720
// characters". The masking the second block performs is the platform's
// recorded unimplemented surface; the document's shape is still the
// operation's contract.
func validateDataProtectionPolicyDocument(policyDocument string) error {
	if utf8.RuneCountInString(policyDocument) > logsstore.MaxDataProtectionPolicyDocumentBytes {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("The JSON specified in policyDocument can be up to %d characters", logsstore.MaxDataProtectionPolicyDocumentBytes), 400)
	}
	var doc dppDocument
	if err := json.Unmarshal([]byte(policyDocument), &doc); err != nil {
		return NewLogsError("InvalidParameterException",
			"The policyDocument must be a JSON object", 400)
	}
	if len(doc.Statement) != 2 {
		return NewLogsError("InvalidParameterException",
			"The data protection policy must include two JSON blocks: one with an Audit action and one with a Deidentify action", 400)
	}
	audit, deidentify := doc.Statement[0], doc.Statement[1]

	if len(audit.DataIdentifer) == 0 {
		return NewLogsError("InvalidParameterException",
			"The first block must include a DataIdentifer array", 400)
	}
	if len(audit.Operation) != 1 {
		return NewLogsError("InvalidParameterException",
			"The first block's Operation must carry the Audit action alone", 400)
	}
	auditRaw, ok := audit.Operation["Audit"]
	if !ok {
		return NewLogsError("InvalidParameterException",
			"The first block's Operation must carry the Audit action", 400)
	}
	var auditAction struct {
		FindingsDestination map[string]json.RawMessage `json:"FindingsDestination"`
	}
	if err := json.Unmarshal(auditRaw, &auditAction); err != nil || auditAction.FindingsDestination == nil {
		return NewLogsError("InvalidParameterException",
			"The Audit action must contain a FindingsDestination object", 400)
	}

	if len(deidentify.DataIdentifer) == 0 {
		return NewLogsError("InvalidParameterException",
			"The second block must include a DataIdentifer array", 400)
	}
	if len(audit.DataIdentifer) != len(deidentify.DataIdentifer) {
		return NewLogsError("InvalidParameterException",
			"The contents of the two DataIdentifer arrays must match exactly", 400)
	}
	for i, id := range audit.DataIdentifer {
		if deidentify.DataIdentifer[i] != id {
			return NewLogsError("InvalidParameterException",
				"The contents of the two DataIdentifer arrays must match exactly", 400)
		}
	}
	if len(deidentify.Operation) != 1 {
		return NewLogsError("InvalidParameterException",
			"The second block's Operation must carry the Deidentify action alone", 400)
	}
	deidRaw, ok := deidentify.Operation["Deidentify"]
	if !ok {
		return NewLogsError("InvalidParameterException",
			"The second block's Operation must carry the Deidentify action", 400)
	}
	var deidentifyAction struct {
		MaskConfig map[string]json.RawMessage `json:"MaskConfig"`
	}
	if err := json.Unmarshal(deidRaw, &deidentifyAction); err != nil || deidentifyAction.MaskConfig == nil {
		return NewLogsError("InvalidParameterException",
			`The Deidentify action must contain the "MaskConfig": {} object`, 400)
	}
	if len(deidentifyAction.MaskConfig) != 0 {
		return NewLogsError("InvalidParameterException",
			`The "MaskConfig": {} object must be empty`, 400)
	}
	return nil
}
