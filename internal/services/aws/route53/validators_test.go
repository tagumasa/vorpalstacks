package route53

import (
	"strings"
	"testing"

	route53store "vorpalstacks/internal/store/aws/route53"
)

// TestValidateCommentUnicodeLengths pins that hosted-zone comments follow
// the Smithy ResourceDescription @length(0, 256) trait counted in Unicode
// characters; the shape carries no pattern, so multibyte comments are valid
// input and must not be rejected on byte length.
func TestValidateCommentUnicodeLengths(t *testing.T) {
	cjk := "\u65e5" // one CJK character, 3 bytes

	if err := validateComment(strings.Repeat(cjk, 128)); err != nil {
		t.Errorf("128-character CJK comment rejected: %v", err)
	}
	if err := validateComment(strings.Repeat(cjk, 256)); err != nil {
		t.Errorf("256-character CJK comment rejected: %v", err)
	}
	if err := validateComment(strings.Repeat(cjk, 257)); err == nil {
		t.Error("257-character CJK comment accepted")
	}
	if err := validateComment(""); err != nil {
		t.Errorf("empty comment rejected: %v", err)
	}
}

// TestValidateHealthCheckConfigUnicodeLengths pins that the health-check
// SearchString / ResourcePath @length(0,255) traits are counted in Unicode
// characters; neither shape carries a pattern, and SearchString in
// particular matches response bodies where multibyte text is realistic.
func TestValidateHealthCheckConfigUnicodeLengths(t *testing.T) {
	cjk := "\u65e5" // one CJK character, 3 bytes

	base := func() *route53store.HealthCheckConfig {
		return &route53store.HealthCheckConfig{
			Type:             "HTTP_STR_MATCH",
			IPAddress:        "192.0.2.10",
			RequestInterval:  30,
			FailureThreshold: 3,
		}
	}

	cfg := base()
	cfg.SearchString = strings.Repeat(cjk, 255)
	if err := validateHealthCheckConfig(cfg); err != nil {
		t.Errorf("255-character CJK search string rejected: %v", err)
	}
	cfg = base()
	cfg.SearchString = strings.Repeat(cjk, 256)
	if err := validateHealthCheckConfig(cfg); err == nil {
		t.Error("256-character CJK search string accepted")
	}

	cfg = base()
	cfg.Type = "HTTP"
	cfg.ResourcePath = strings.Repeat(cjk, 255)
	if err := validateHealthCheckConfig(cfg); err != nil {
		t.Errorf("255-character CJK resource path rejected: %v", err)
	}
	cfg = base()
	cfg.Type = "HTTP"
	cfg.ResourcePath = strings.Repeat(cjk, 256)
	if err := validateHealthCheckConfig(cfg); err == nil {
		t.Error("256-character CJK resource path accepted")
	}
}

// TestCallerReferenceValidators pins the three nonce contracts behind the
// required CallerReference members: Nonce (hosted zone, 1 to 128
// characters), CidrNonce (CIDR collection, 1 to 64 characters plus the
// ASCII-only pattern) and HealthCheckNonce (health check, 1 to 64
// characters). Lengths count Unicode characters like every @length trait,
// so a multibyte hosted-zone token is valid up to 128 characters.
func TestCallerReferenceValidators(t *testing.T) {
	cjk := "\u65e5" // one CJK character, 3 bytes

	// Hosted zone — Nonce @length(1, 128), no pattern.
	if err := validateHostedZoneCallerReference(""); err == nil {
		t.Error("empty hosted-zone CallerReference accepted")
	}
	if err := validateHostedZoneCallerReference(strings.Repeat("a", 128)); err != nil {
		t.Errorf("128-character CallerReference rejected: %v", err)
	}
	if err := validateHostedZoneCallerReference(strings.Repeat(cjk, 128)); err != nil {
		t.Errorf("128-character CJK CallerReference rejected: %v", err)
	}
	if err := validateHostedZoneCallerReference(strings.Repeat("a", 129)); err == nil {
		t.Error("129-character hosted-zone CallerReference accepted")
	}

	// CIDR collection — CidrNonce @length(1, 64) and @pattern ASCII only.
	if err := validateCidrCallerReference(""); err == nil {
		t.Error("empty CIDR-collection CallerReference accepted")
	}
	if err := validateCidrCallerReference(strings.Repeat("a", 64)); err != nil {
		t.Errorf("64-character CIDR-collection CallerReference rejected: %v", err)
	}
	if err := validateCidrCallerReference(strings.Repeat("a", 65)); err == nil {
		t.Error("65-character CIDR-collection CallerReference accepted")
	}
	if err := validateCidrCallerReference(cjk); err == nil {
		t.Error("non-ASCII CIDR-collection CallerReference accepted")
	}

	// Health check — HealthCheckNonce @length(1, 64), no pattern.
	if err := validateHealthCheckCallerReference(""); err == nil {
		t.Error("empty health-check CallerReference accepted")
	}
	if err := validateHealthCheckCallerReference(strings.Repeat("a", 64)); err != nil {
		t.Errorf("64-character health-check CallerReference rejected: %v", err)
	}
	if err := validateHealthCheckCallerReference(strings.Repeat("a", 65)); err == nil {
		t.Error("65-character health-check CallerReference accepted")
	}
}

// TestResourceRecordValueContract pins the ResourceRecord element contract:
// the Value member is required on the shape, the RData @length(0, 4000)
// bound counts characters, and an explicitly empty value stays valid
// because the RData minimum is 0.
func TestResourceRecordValueContract(t *testing.T) {
	// A raw element without the Value member must be reported absent so the
	// change batch is rejected rather than persisting an empty RData record.
	if _, present := resourceRecordValueFromMap(map[string]interface{}{"NotValue": "x"}); present {
		t.Error("element without the Value member reported as present")
	}
	v, present := resourceRecordValueFromMap(map[string]interface{}{"Value": "10.0.0.1"})
	if !present || v != "10.0.0.1" {
		t.Errorf("Value element not extracted: %q %v", v, present)
	}
	// The case fallback mirrors GetStringParam.
	v, present = resourceRecordValueFromMap(map[string]interface{}{"value": "10.0.0.2"})
	if !present || v != "10.0.0.2" {
		t.Errorf("lower-case value element not extracted: %q %v", v, present)
	}

	if err := validateResourceRecordValue(false, ""); err == nil {
		t.Error("member-absent element accepted")
	} else if !strings.Contains(err.Error(), "InvalidChangeBatch") {
		t.Errorf("member-absent element rejected with %v, want InvalidChangeBatch", err)
	}
	if err := validateResourceRecordValue(true, strings.Repeat("a", 4000)); err != nil {
		t.Errorf("4000-character value rejected: %v", err)
	}
	if err := validateResourceRecordValue(true, strings.Repeat("a", 4001)); err == nil {
		t.Error("4001-character value accepted")
	} else if !strings.Contains(err.Error(), "InvalidChangeBatch") {
		t.Errorf("4001-character value rejected with %v, want InvalidChangeBatch", err)
	}
	// The @length minimum is 0, so an explicitly empty value stays valid.
	if err := validateResourceRecordValue(true, ""); err != nil {
		t.Errorf("explicitly empty value rejected: %v", err)
	}
}
