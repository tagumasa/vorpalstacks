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
