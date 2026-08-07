package sesv2

import (
	"encoding/json"
	"fmt"
	"net/mail"
	"regexp"
	"strings"

	sesv2store "vorpalstacks/internal/store/aws/sesv2"
)

// ---------------------------------------------------------------------------
// Enum validators — values sourced from Smithy enum shapes
// (third_party/api-models-aws/models/sesv2/service/2019-09-27)
// ---------------------------------------------------------------------------

var validMailTypes = map[string]bool{
	"MARKETING":     true,
	"TRANSACTIONAL": true,
}

func validateMailType(v string) bool {
	return validMailTypes[v]
}

var validContactLanguages = map[string]bool{
	"EN": true,
	"JA": true,
}

func validateContactLanguage(v string) bool {
	return validContactLanguages[v]
}

var validTlsPolicies = map[string]bool{
	"REQUIRE":  true,
	"OPTIONAL": true,
}

func validateTlsPolicy(v string) bool {
	return validTlsPolicies[v]
}

var validHttpsPolicies = map[string]bool{
	"REQUIRE":           true,
	"REQUIRE_OPEN_ONLY": true,
	"OPTIONAL":          true,
}

func validateHttpsPolicy(v string) bool {
	return validHttpsPolicies[v]
}

var validBehaviorOnMxFailures = map[string]bool{
	"USE_DEFAULT_VALUE": true,
	"REJECT_MESSAGE":    true,
}

func validateBehaviorOnMxFailure(v string) bool {
	return validBehaviorOnMxFailures[v]
}

var validScalingModes = map[string]bool{
	"STANDARD": true,
	"MANAGED":  true,
}

func validateScalingMode(v string) bool {
	return validScalingModes[v]
}

var validFeatureStatuses = map[string]bool{
	"ENABLED":  true,
	"DISABLED": true,
}

func validateFeatureStatus(v string) bool {
	return validFeatureStatuses[v]
}

var validSuppressionListReasons = map[string]bool{
	"BOUNCE":    true,
	"COMPLAINT": true,
}

func validateSuppressionListReason(v string) bool {
	return validSuppressionListReasons[v]
}

var validSuppressionListScopes = map[string]bool{
	"ACCOUNT": true,
	"TENANT":  true,
}

func validateSuppressionListScope(v string) bool {
	return validSuppressionListScopes[v]
}

var validSuppressionConfidenceVerdictThresholds = map[string]bool{
	"MEDIUM":  true,
	"HIGH":    true,
	"MANAGED": true,
}

func validateSuppressionConfidenceVerdictThreshold(v string) bool {
	return validSuppressionConfidenceVerdictThresholds[v]
}

var validSubscriptionStatuses = map[string]bool{
	"OPT_IN":  true,
	"OPT_OUT": true,
}

func validateSubscriptionStatus(v string) bool {
	return validSubscriptionStatuses[v]
}

// DkimSigningAttributesOrigin enum has 30+ regional AWS_SES_* values
// plus AWS_SES and EXTERNAL.  Accept any string that starts with
// "AWS_SES" or equals "EXTERNAL".
func validateDkimSigningAttributesOrigin(v string) bool {
	if v == "EXTERNAL" {
		return true
	}
	if strings.HasPrefix(v, "AWS_SES") {
		return true
	}
	return false
}

var validEventTypes = map[string]bool{
	"SEND":              true,
	"REJECT":            true,
	"BOUNCE":            true,
	"COMPLAINT":         true,
	"DELIVERY":          true,
	"OPEN":              true,
	"CLICK":             true,
	"RENDERING_FAILURE": true,
	"DELIVERY_DELAY":    true,
	"SUBSCRIPTION":      true,
}

func validateEventType(v string) bool {
	return validEventTypes[v]
}

func validateEventTypes(types []string) bool {
	for _, t := range types {
		if !validateEventType(t) {
			return false
		}
	}
	return true
}

// ---------------------------------------------------------------------------
// Range validators
// ---------------------------------------------------------------------------

// validateMaxDeliverySeconds checks the Smithy @range(min=300, max=50400).
func validateMaxDeliverySeconds(v int32) bool {
	return v >= 300 && v <= 50400
}

// ---------------------------------------------------------------------------
// Length / pattern validators
// ---------------------------------------------------------------------------

var configSetNamePattern = regexp.MustCompile(`^[a-zA-Z0-9_-]{1,64}$`)

func validateConfigurationSetName(name string) bool {
	return configSetNamePattern.MatchString(name)
}

var emailTemplateNamePattern = regexp.MustCompile(`^[a-zA-Z0-9_-]{1,64}$`)

func validateEmailTemplateName(name string) bool {
	return emailTemplateNamePattern.MatchString(name)
}

var poolNamePattern = regexp.MustCompile(`^[a-zA-Z0-9_-]{1,63}$`)

func validatePoolName(name string) bool {
	return poolNamePattern.MatchString(name)
}

var contactListNamePattern = regexp.MustCompile(`^[a-zA-Z0-9_-]{1,128}$`)

func validateContactListName(name string) bool {
	return contactListNamePattern.MatchString(name)
}

// ---------------------------------------------------------------------------
// Email / address validators
// ---------------------------------------------------------------------------

// validateFromAddress validates that fromEmailAddress is a syntactically
// valid email address or domain. A bare domain (no "@") is accepted because
// SES allows the FromEmailAddress to reference a verified domain identity.
func validateFromAddress(addr string) bool {
	if addr == "" {
		return false
	}
	if !strings.Contains(addr, "@") {
		return true
	}
	parts := strings.Split(addr, "@")
	if len(parts) != 2 {
		return false
	}
	return parts[0] != "" && parts[1] != ""
}

// validateEmailAddress validates a recipient email address. Unlike
// validateFromAddress, bare domains are rejected — recipients must be
// individual mailboxes.
func validateEmailAddress(addr string) bool {
	if addr == "" {
		return false
	}
	parts := strings.Split(addr, "@")
	if len(parts) != 2 {
		return false
	}
	return parts[0] != "" && parts[1] != ""
}

// validateEmailAddressRFC5321 validates an email address using Go's
// net/mail parser. This rejects obviously malformed addresses that
// would otherwise be silently accepted.
func validateEmailAddressRFC5321(addr string) bool {
	if addr == "" {
		return false
	}
	_, err := mail.ParseAddress(addr)
	return err == nil
}

// ---------------------------------------------------------------------------
// Composite validators
// ---------------------------------------------------------------------------

// validateSingleEventDestination enforces the AWS constraint that only one
// of the five destination types may be specified per event destination.
// Returns the count of non-nil destinations.
func countEventDestinations(sns, kinesis, cloudwatch, pinpoint, eventbridge bool) int {
	count := 0
	if sns {
		count++
	}
	if kinesis {
		count++
	}
	if cloudwatch {
		count++
	}
	if pinpoint {
		count++
	}
	if eventbridge {
		count++
	}
	return count
}

// validatePolicyJSON validates that the policy string is valid JSON.
// AWS expects an IAM-style policy document for SES sending authorisation.
func validatePolicyJSON(policy string) error {
	if policy == "" {
		return fmt.Errorf("Policy is required")
	}
	var v interface{}
	if err := json.Unmarshal([]byte(policy), &v); err != nil {
		return fmt.Errorf("Policy must be valid JSON: %w", err)
	}
	return nil
}

// ---------------------------------------------------------------------------
// Identity format validator — checks that an identity is either a valid email
// address or a domain name with at least one dot. Whitespace and control
// characters are rejected.
func validateIdentityFormat(identity string) bool {
	if identity == "" {
		return false
	}
	// Reject whitespace and control characters.
	for _, r := range identity {
		if r <= 0x20 || r == 0x7f {
			return false
		}
	}
	// Domain identity: no "@" — must contain at least one dot.
	if !strings.Contains(identity, "@") {
		return strings.Contains(identity, ".")
	}
	// Email address: exactly one "@" with non-empty local and domain parts.
	parts := strings.Split(identity, "@")
	if len(parts) != 2 {
		return false
	}
	return parts[0] != "" && parts[1] != ""
}

// ---------------------------------------------------------------------------
// Raw message size cap
// ---------------------------------------------------------------------------

// maxRawMessageSize is the maximum size of a raw email message in bytes.
// AWS SES accepts messages up to 10 MB.
const maxRawMessageSize = 10 * 1024 * 1024

// maxTemplateDataSize is the Smithy @length max for EmailTemplateData
// (com.amazonaws.sesv2#EmailTemplateData: length max 262144).
const maxTemplateDataSize = 262144

// maxAttributesDataSize caps the AttributesData JSON blob on contacts.
// Smithy does not declare a formal length trait; the cap mirrors
// TemplateData to prevent unbounded JSON payload DoS.
const maxAttributesDataSize = 262144

// maxSimpleBodySize caps the total size of Simple email body content
// (subject + html + text). AWS SES imposes a 10 MB message limit; this
// cap is the same order as maxRawMessageSize for consistency.
const maxSimpleBodySize = 10 * 1024 * 1024

// maxFilteredSuppressionScan is the safety upper bound for the in-memory
// filter walk in ListSuppressedDestinations, mirroring the
// maxFilteredContactScan guard used by ListContacts.
const maxFilteredSuppressionScan = 50000

// ---------------------------------------------------------------------------
// Message-tag validators
// ---------------------------------------------------------------------------

var messageTagFieldPattern = regexp.MustCompile(`^[a-zA-Z0-9_-]{1,256}$`)

// validateMessageTagField checks a single Name or Value against the
// AWS-documented charset (ASCII letters, digits, underscores, dashes)
// and length (max 256) constraint for MessageTag fields.
func validateMessageTagField(v string) bool {
	return messageTagFieldPattern.MatchString(v)
}

const maxMessageTags = 50

// validateMessageTags validates every tag in a slice of MessageTag. Each
// Name and Value must satisfy validateMessageTagField, and the total
// number of tags must not exceed maxMessageTags.
func validateMessageTags(tags []sesv2store.MessageTag) error {
	if len(tags) > maxMessageTags {
		return fmt.Errorf("too many EmailTags (max %d)", maxMessageTags)
	}
	for _, t := range tags {
		if !validateMessageTagField(t.Name) {
			return fmt.Errorf("invalid EmailTag Name: %q", t.Name)
		}
		if !validateMessageTagField(t.Value) {
			return fmt.Errorf("invalid EmailTag Value: %q", t.Value)
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// Reply-to address validator
// ---------------------------------------------------------------------------

// validateReplyToAddresses validates every address in a ReplyToAddresses
// list using the RFC 5321 parser. Display-name forms are rejected because
// reply-to addresses are individual mailboxes, not identities.
func validateReplyToAddresses(addrs []string) error {
	for _, a := range addrs {
		if !validateSuppressionEmailAddress(a) {
			return fmt.Errorf("invalid ReplyToAddresses entry: %q", a)
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// Strict email validator (addr-spec only, no display-name rewriting)
// ---------------------------------------------------------------------------

// validateSuppressionEmailAddress validates that the supplied string is a
// bare RFC 5321 addr-spec (e.g. user@example.com). net/mail.ParseAddress
// also accepts display-name forms ("Name" <addr@domain>) and silently
// rewrites them, which is wrong for suppression-list keys: the original
// string becomes the store lookup key, so a subsequent
// GetSuppressedDestination("user@example.com") would miss an entry stored
// as '"Name" <user@example.com>'. Reject when the parser rewrites the
// input.
func validateSuppressionEmailAddress(addr string) bool {
	if addr == "" {
		return false
	}
	parsed, err := mail.ParseAddress(addr)
	if err != nil {
		return false
	}
	return parsed.Address == addr
}

// ---------------------------------------------------------------------------
// Email content validator
// ---------------------------------------------------------------------------

// validateContent checks that the supplied EmailContent has meaningful
// payload. An empty Simple (nil Body), empty Raw (zero-length Data), or
// empty Template (empty TemplateName) is rejected because AWS rejects
// content-less emails with MessageRejected / InvalidContent.
func validateContent(content *sesv2store.EmailContent) error {
	if content == nil {
		return fmt.Errorf("Content is required")
	}
	if content.Simple != nil {
		totalSize := 0
		if content.Simple.Body != nil {
			if content.Simple.Body.Html != nil {
				totalSize += len(content.Simple.Body.Html.Data)
			}
			if content.Simple.Body.Text != nil {
				totalSize += len(content.Simple.Body.Text.Data)
			}
		}
		if totalSize == 0 {
			return fmt.Errorf("Simple content must have a non-empty body")
		}
		if totalSize > maxSimpleBodySize {
			return fmt.Errorf("Simple content exceeds size limit")
		}
		return nil
	}
	if content.Raw != nil {
		if len(content.Raw.Data) == 0 {
			return fmt.Errorf("Raw content must have non-empty Data")
		}
		return nil
	}
	if content.Template != nil {
		if content.Template.TemplateName == "" {
			return fmt.Errorf("Template content must have a TemplateName")
		}
		return nil
	}
	return fmt.Errorf("Content must specify Simple, Raw, or Template")
}
