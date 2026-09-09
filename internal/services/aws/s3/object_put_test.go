package s3

import (
	"testing"
)

// The x-amz-tagging header is URL query-parameter encoded per the API
// contract: keys and values are percent-decoded, and a malformed escape is
// an invalid argument rather than a verbatim store.
func TestParseTaggingHeader(t *testing.T) {
	tags, err := parseTaggingHeader("a%20b=c%21d&plain=pair")
	if err != nil {
		t.Fatalf("parseTaggingHeader: %v", err)
	}
	if len(tags) != 2 {
		t.Fatalf("tags = %d, want 2", len(tags))
	}
	if tags[0].Key != "a b" || tags[0].Value != "c!d" {
		t.Fatalf("first tag = %q=%q, want a b=c!d", tags[0].Key, tags[0].Value)
	}
	if tags[1].Key != "plain" || tags[1].Value != "pair" {
		t.Fatalf("second tag = %q=%q, want plain=pair", tags[1].Key, tags[1].Value)
	}

	if _, err := parseTaggingHeader("bad%E=c"); err == nil {
		t.Fatal("malformed escape must be rejected")
	}
	if _, err := parseTaggingHeader("ok=bad%ZZ"); err == nil {
		t.Fatal("malformed escape in value must be rejected")
	}
	if got, err := parseTaggingHeader(""); got != nil || err != nil {
		t.Fatalf("empty header = %v, %v; want nil, nil", got, err)
	}
	if got, err := parseTaggingHeader("noequals"); got != nil || err != nil {
		t.Fatalf("pair without '=' = %v, %v; want nil, nil", got, err)
	}
}
