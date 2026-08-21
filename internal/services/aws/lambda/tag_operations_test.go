package lambda

import (
	"testing"
)

// TestSplitTaggableResource pins the TaggableResource pattern parsing:
// function resources with and without a qualifier, event source
// mappings, and non-taggable shapes.
func TestSplitTaggableResource(t *testing.T) {
	cases := []struct {
		arn                 string
		resourceType, id, q string
	}{
		{"arn:aws:lambda:us-east-1:123456789012:function:myfn", "function", "myfn", ""},
		{"arn:aws:lambda:us-east-1:123456789012:function:myfn:PROD", "function", "myfn", "PROD"},
		{"arn:aws:lambda:us-east-1:123456789012:function:myfn:$LATEST", "function", "myfn", "$LATEST"},
		{"arn:aws:lambda:us-east-1:123456789012:event-source-mapping:12345678-1234-1234-1234-123456789012", "event-source-mapping", "12345678-1234-1234-1234-123456789012", ""},
		{"arn:aws:lambda:us-east-1:123456789012:layer:mylayer", "layer", "mylayer", ""},
		{"not-an-arn", "", "", ""},
	}
	for _, tc := range cases {
		rt, id, q := splitTaggableResource(tc.arn)
		if rt != tc.resourceType || id != tc.id || q != tc.q {
			t.Fatalf("splitTaggableResource(%s) = (%q, %q, %q), want (%q, %q, %q)",
				tc.arn, rt, id, q, tc.resourceType, tc.id, tc.q)
		}
	}
}

// TestTagResourceKey_Transform pins the tag store key derivation: ESM
// ARNs map to the namespaced key, qualified function ARNs keep the
// qualifier for validation to reject, plain functions stay bare.
func TestTagResourceKey_Transform(t *testing.T) {
	cfg := lambdaTagConfig(nil, nil)
	cases := []struct {
		arn  string
		want string
	}{
		{"arn:aws:lambda:us-east-1:123456789012:function:myfn", "myfn"},
		{"arn:aws:lambda:us-east-1:123456789012:function:myfn:PROD", "myfn:PROD"},
		{"arn:aws:lambda:us-east-1:123456789012:event-source-mapping:12345678-1234-1234-1234-123456789012",
			"event-source-mapping/12345678-1234-1234-1234-123456789012"},
	}
	for _, tc := range cases {
		if got := cfg.ResourceKey(tc.arn); got != tc.want {
			t.Fatalf("ResourceKey(%s) = %q, want %q", tc.arn, got, tc.want)
		}
	}
	if got := esmTagResourceKey("u"); got != "event-source-mapping/u" {
		t.Fatalf("esmTagResourceKey = %q", got)
	}
}
