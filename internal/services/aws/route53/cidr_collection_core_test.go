package route53

import (
	"encoding/base64"
	"testing"
)

func TestDecodeCidrBlocksToken(t *testing.T) {
	valid := base64.StdEncoding.EncodeToString([]byte("loc-a\x0010.0.0.0/24"))
	loc, cidr, err := decodeCidrBlocksToken(valid)
	if err != nil || loc != "loc-a" || cidr != "10.0.0.0/24" {
		t.Fatalf("round trip: got (%q, %q, %v)", loc, cidr, err)
	}

	// An absent token means "start from the beginning" and is not an error.
	loc, cidr, err = decodeCidrBlocksToken("")
	if err != nil || loc != "" || cidr != "" {
		t.Fatalf("absent token: got (%q, %q, %v)", loc, cidr, err)
	}

	invalid := map[string]string{
		"not base64":     "!!!not-base64!!!",
		"no separator":   base64.StdEncoding.EncodeToString([]byte("loc-only")),
		"empty location": base64.StdEncoding.EncodeToString([]byte("\x0010.0.0.0/24")),
		"empty cidr":     base64.StdEncoding.EncodeToString([]byte("loc-a\x00")),
	}
	for name, token := range invalid {
		if _, _, err := decodeCidrBlocksToken(token); err == nil {
			t.Fatalf("%s: expected rejection", name)
		}
	}
}

func TestDecodeCidrLocationsToken(t *testing.T) {
	name, err := decodeCidrLocationsToken(encodeCidrLocationsToken("loc-a"))
	if err != nil || name != "loc-a" {
		t.Fatalf("round trip: got (%q, %v)", name, err)
	}

	invalid := map[string]string{
		"not base64": "!!!not-base64!!!",
		"empty name": base64.StdEncoding.EncodeToString([]byte("")),
	}
	for name, token := range invalid {
		if _, err := decodeCidrLocationsToken(token); err == nil {
			t.Fatalf("%s: expected rejection", name)
		}
	}
}

// TestValidateCidrLocationName pins the location name member constraints:
// 1-16 characters from [0-9A-Za-z_-]. Longer names and forbidden
// characters are rejected with InvalidInput.
func TestValidateCidrLocationName(t *testing.T) {
	for _, name := range []string{"a", "loc-a", "Site_01", "0123456789abcdef"} {
		if err := validateCidrLocationName(name); err != nil {
			t.Fatalf("%q must be valid: %v", name, err)
		}
	}
	invalid := map[string]string{
		"17 characters":  "0123456789abcdefg",
		"space":          "loc a",
		"dot":            "loc.a",
		"slash":          "loc/a",
		"unicode":        "場所",
		"empty":          "",
		"empty is blank": " ",
	}
	for name, value := range invalid {
		if err := validateCidrLocationName(value); err == nil {
			t.Fatalf("%s: %q must be rejected", name, value)
		}
	}
}

// TestValidateCidrEntry pins the Cidr member constraints: 1-50 characters
// with at least one non-whitespace character.
func TestValidateCidrEntry(t *testing.T) {
	for _, cidr := range []string{"10.0.0.0/24", "2001:db8::/32", "a"} {
		if err := validateCidrEntry(cidr); err != nil {
			t.Fatalf("%q must be valid: %v", cidr, err)
		}
	}
	invalid := map[string]string{
		"empty":    "",
		"blank":    "   ",
		"too long": "1234567890123456789012345678901234567890123456789012345",
	}
	for name, value := range invalid {
		if err := validateCidrEntry(value); err == nil {
			t.Fatalf("%s: %q must be rejected", name, value)
		}
	}
}
