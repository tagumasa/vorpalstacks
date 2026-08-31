package route53

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	awserrors "vorpalstacks/internal/common/errors"
	route53store "vorpalstacks/internal/store/aws/route53"
)

// ---------------------------------------------------------------------------
// Health Check Type enum (Smithy: com.amazonaws.route53#HealthCheckType)
// ---------------------------------------------------------------------------

var validHealthCheckTypes = map[string]bool{
	"HTTP":              true,
	"HTTPS":             true,
	"HTTP_STR_MATCH":    true,
	"HTTPS_STR_MATCH":   true,
	"TCP":               true,
	"CALCULATED":        true,
	"CLOUDWATCH_METRIC": true,
	"RECOVERY_CONTROL":  true,
}

func validateHealthCheckType(t string) bool {
	return validHealthCheckTypes[t]
}

// ---------------------------------------------------------------------------
// Record Type enum (Smithy: com.amazonaws.route53#RRType — 17 values)
// ---------------------------------------------------------------------------

var validRecordTypes = map[string]bool{
	"SOA":   true,
	"A":     true,
	"TXT":   true,
	"NS":    true,
	"CNAME": true,
	"MX":    true,
	"NAPTR": true,
	"PTR":   true,
	"SRV":   true,
	"SPF":   true,
	"AAAA":  true,
	"CAA":   true,
	"DS":    true,
	"TLSA":  true,
	"SSHFP": true,
	"SVCB":  true,
	"HTTPS": true,
}

func validateRecordType(t string) bool {
	return validRecordTypes[t]
}

// ---------------------------------------------------------------------------
// Hosted Zone Type enum (Smithy: com.amazonaws.route53#HostedZoneType)
// Only one enum member: PRIVATE_HOSTED_ZONE (enumValue: PrivateHostedZone).
// Absence of the filter = public zones.
// ---------------------------------------------------------------------------

func isPrivateHostedZoneFilter(v string) bool {
	if v == "" {
		return false
	}
	switch strings.ToLower(v) {
	case "privatehostedzone", "private_hosted_zone", "private":
		return true
	}
	return false
}

// ---------------------------------------------------------------------------
// Change Action enum (Smithy: com.amazonaws.route53#ChangeAction)
// ---------------------------------------------------------------------------

var validChangeActions = map[string]bool{
	"CREATE": true,
	"DELETE": true,
	"UPSERT": true,
}

func validateChangeAction(action string) bool {
	return validChangeActions[action]
}

// ---------------------------------------------------------------------------
// InsufficientDataHealthStatus enum
// (Smithy: com.amazonaws.route53#InsufficientDataHealthStatus)
// ---------------------------------------------------------------------------

var validInsufficientDataHealthStatus = map[string]bool{
	"Healthy":         true,
	"Unhealthy":       true,
	"LastKnownStatus": true,
}

func validateInsufficientDataHealthStatus(v string) bool {
	return validInsufficientDataHealthStatus[v]
}

// ---------------------------------------------------------------------------
// Comment length (Smithy: com.amazonaws.route53#ResourceDescription, 0-256
// counted in Unicode characters; the shape carries no pattern)
// ---------------------------------------------------------------------------

// maxCommentLength is the ResourceDescription @length maximum, counted in
// Unicode characters like every @length trait.
const maxCommentLength = 256

func validateComment(comment string) error {
	if utf8.RuneCountInString(comment) > maxCommentLength {
		return awserrors.NewAWSError("InvalidInput", "Comment must not exceed 256 characters", 400)
	}
	return nil
}

// ---------------------------------------------------------------------------
// CallerReference (nonce) validation — every Route 53 Create* request carries
// a required CallerReference idempotency token. The token must come from the
// caller: a server-synthesised value changes on every retry, defeating the
// execute-once semantics the member exists to provide, so an omitted member
// is rejected instead of being filled in.
// ---------------------------------------------------------------------------

// maxHostedZoneCallerReferenceLen is the Nonce @length maximum for
// CreateHostedZoneRequest.CallerReference (1 to 128 characters).
const maxHostedZoneCallerReferenceLen = 128

// maxCidrCallerReferenceLen is the CidrNonce @length maximum for
// CreateCidrCollectionRequest.CallerReference (1 to 64 characters).
const maxCidrCallerReferenceLen = 64

// maxHealthCheckCallerReferenceLen is the HealthCheckNonce @length maximum
// for CreateHealthCheckRequest.CallerReference (1 to 64 characters).
const maxHealthCheckCallerReferenceLen = 64

func validateHostedZoneCallerReference(ref string) error {
	if ref == "" {
		return awserrors.NewAWSError("InvalidInput", "CallerReference is required", 400)
	}
	if utf8.RuneCountInString(ref) > maxHostedZoneCallerReferenceLen {
		return awserrors.NewAWSError("InvalidInput", "CallerReference must be 1 to 128 characters long", 400)
	}
	return nil
}

// validateCidrCallerReference also enforces the CidrNonce @pattern (ASCII
// characters only) on top of the shared required-and-length rule.
func validateCidrCallerReference(ref string) error {
	if ref == "" {
		return awserrors.NewAWSError("InvalidInput", "CallerReference is required", 400)
	}
	if utf8.RuneCountInString(ref) > maxCidrCallerReferenceLen {
		return awserrors.NewAWSError("InvalidInput", "CallerReference must be 1 to 64 characters long", 400)
	}
	for _, c := range ref {
		if c > unicode.MaxASCII {
			return awserrors.NewAWSError("InvalidInput", "CallerReference may contain ASCII characters only", 400)
		}
	}
	return nil
}

func validateHealthCheckCallerReference(ref string) error {
	if ref == "" {
		return awserrors.NewAWSError("InvalidInput", "CallerReference is required", 400)
	}
	if utf8.RuneCountInString(ref) > maxHealthCheckCallerReferenceLen {
		return awserrors.NewAWSError("InvalidInput", "CallerReference must be 1 to 64 characters long", 400)
	}
	return nil
}

// ---------------------------------------------------------------------------
// Resource record values (Smithy: com.amazonaws.route53#RData)
// ---------------------------------------------------------------------------

// maxRDataLen is the RData @length maximum for ResourceRecord.Value
// (0 to 4000 characters; the member itself is required on the shape).
const maxRDataLen = 4000

// validateResourceRecordValue enforces the ResourceRecord element contract:
// the Value member must be present — an element without it is not a resource
// record and invalidates the whole change batch — and must stay within the
// RData length bound. The RData minimum is 0, so an explicitly empty value
// remains shape-valid and accepted.
func validateResourceRecordValue(present bool, value string) error {
	if !present {
		return awserrors.NewAWSError("InvalidChangeBatch",
			"Invalid Resource Record: Value is required for every resource record", 400)
	}
	if utf8.RuneCountInString(value) > maxRDataLen {
		return awserrors.NewAWSError("InvalidChangeBatch",
			"Invalid Resource Record: the record value must not exceed 4000 characters", 400)
	}
	return nil
}

// ---------------------------------------------------------------------------
// Domain name validation (RFC 1035)
// ---------------------------------------------------------------------------

func validateDomainName(name string) error {
	n := strings.TrimSuffix(name, ".")
	if n == "" {
		return awserrors.NewAWSError("InvalidDomainName", "Hosted zone name is required", 400)
	}
	if len(n) > 253 {
		return awserrors.NewAWSError("InvalidDomainName", "Hosted zone name must not exceed 253 characters", 400)
	}
	labels := strings.Split(n, ".")
	for _, label := range labels {
		if label == "" {
			return awserrors.NewAWSError("InvalidDomainName", "Hosted zone name contains an empty label", 400)
		}
		if len(label) > 63 {
			return awserrors.NewAWSError("InvalidDomainName", "DNS label must not exceed 63 characters", 400)
		}
		for i, c := range label {
			isAlnum := (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
			isHyphen := c == '-'
			if !isAlnum && !isHyphen {
				return awserrors.NewAWSError("InvalidDomainName",
					fmt.Sprintf("Invalid character in DNS label %q", label), 400)
			}
			if i == 0 && isHyphen {
				return awserrors.NewAWSError("InvalidDomainName",
					fmt.Sprintf("DNS label %q must not start with a hyphen", label), 400)
			}
			if i == len(label)-1 && isHyphen {
				return awserrors.NewAWSError("InvalidDomainName",
					fmt.Sprintf("DNS label %q must not end with a hyphen", label), 400)
			}
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// Health Check Config validation
// Validates Type enum, numeric ranges, and string length limits per
// Smithy traits and AWS documentation.
// ---------------------------------------------------------------------------

func validateHealthCheckConfig(config *route53store.HealthCheckConfig) error {
	if config == nil {
		return awserrors.NewAWSError("InvalidInput", "HealthCheckConfig is required", 400)
	}

	if !validateHealthCheckType(config.Type) {
		return awserrors.NewAWSError("InvalidInput",
			fmt.Sprintf("Invalid or missing health check type: %q. Must be one of: HTTP, HTTPS, HTTP_STR_MATCH, HTTPS_STR_MATCH, TCP, CALCULATED, CLOUDWATCH_METRIC, RECOVERY_CONTROL", config.Type), 400)
	}

	if config.Port > 65535 {
		return awserrors.NewAWSError("InvalidInput", "Port must be between 1 and 65535", 400)
	}
	if config.FailureThreshold > 10 {
		return awserrors.NewAWSError("InvalidInput", "FailureThreshold must be between 1 and 10", 400)
	}
	if config.RequestInterval > 0 && (config.RequestInterval < 10 || config.RequestInterval > 30) {
		return awserrors.NewAWSError("InvalidInput", "RequestInterval must be between 10 and 30", 400)
	}
	if config.HealthThreshold > 256 {
		return awserrors.NewAWSError("InvalidInput", "HealthThreshold must be between 0 and 256", 400)
	}

	// ResourcePath / SearchString / FullyQualifiedDomainName /
	// RoutingControlArn carry @length(0,255) or (1,255) with no pattern, so
	// lengths count Unicode characters (SearchString in particular searches
	// response bodies, where multibyte text is realistic).
	if utf8.RuneCountInString(config.ResourcePath) > 255 {
		return awserrors.NewAWSError("InvalidInput", "ResourcePath must not exceed 255 characters", 400)
	}
	if utf8.RuneCountInString(config.SearchString) > 255 {
		return awserrors.NewAWSError("InvalidInput", "SearchString must not exceed 255 characters", 400)
	}
	if utf8.RuneCountInString(config.FullyQualifiedDomainName) > 255 {
		return awserrors.NewAWSError("InvalidInput", "FullyQualifiedDomainName must not exceed 255 characters", 400)
	}
	if len(config.IPAddress) > 45 {
		return awserrors.NewAWSError("InvalidInput", "IPAddress must not exceed 45 characters", 400)
	}
	if utf8.RuneCountInString(config.RoutingControlArn) > 255 {
		return awserrors.NewAWSError("InvalidInput", "RoutingControlArn must not exceed 255 characters", 400)
	}

	if config.InsufficientDataHealthStatus != "" && !validateInsufficientDataHealthStatus(config.InsufficientDataHealthStatus) {
		return awserrors.NewAWSError("InvalidInput",
			fmt.Sprintf("Invalid InsufficientDataHealthStatus: %q. Must be one of: Healthy, Unhealthy, LastKnownStatus", config.InsufficientDataHealthStatus), 400)
	}

	return nil
}
