package acm

import (
	"crypto/x509"
	"encoding/pem"
	"math/big"
	"regexp"
	"testing"
)

// The serial wire form is the Smithy SerialNumber shape: colon-separated
// lowercase hex byte pairs.
var serialWirePattern = regexp.MustCompile(`^[0-9a-f]{2}(:[0-9a-f]{2}){1,19}$`)

func TestFormatSerialNumberHex(t *testing.T) {
	cases := []struct {
		name string
		in   *big.Int
		want string
	}{
		{"zero serial padded to two pairs", big.NewInt(0), "00:00"},
		{"single byte padded to two pairs", big.NewInt(0x05), "00:05"},
		{"multi byte", new(big.Int).SetBytes([]byte{0xe5, 0x87, 0xef, 0x34, 0x7a, 0x4a, 0x0f, 0xde}), "e5:87:ef:34:7a:4a:0f:de"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := formatSerialNumberHex(tc.in); got != tc.want {
				t.Fatalf("formatSerialNumberHex(%v) = %q, want %q", tc.in, got, tc.want)
			}
			if !serialWirePattern.MatchString(tc.want) {
				t.Fatalf("expected form %q violates the wire pattern", tc.want)
			}
		})
	}
}

func TestIssuedMaterialSerialMatchesWirePattern(t *testing.T) {
	material, err := issueCertificateMaterial("serial-wire.example.com", nil, "RSA_2048")
	if err != nil {
		t.Fatalf("issueCertificateMaterial: %v", err)
	}
	if !serialWirePattern.MatchString(material.serial) {
		t.Fatalf("issued serial %q does not match the wire pattern", material.serial)
	}
	block, _ := pem.Decode([]byte(material.certPEM))
	if block == nil {
		t.Fatal("issued PEM did not decode")
	}
	parsed, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatalf("parse issued certificate: %v", err)
	}
	if want := formatSerialNumberHex(parsed.SerialNumber); material.serial != want {
		t.Fatalf("stored serial %q != certificate serial %q", material.serial, want)
	}
}

func TestIssuedMaterialKeyUsagesAndImportParity(t *testing.T) {
	material, err := issueCertificateMaterial("keyusage.example.com", nil, "RSA_2048")
	if err != nil {
		t.Fatalf("issueCertificateMaterial: %v", err)
	}
	want := []string{"DIGITAL_SIGNATURE", "KEY_ENCIPHERMENT"}
	if len(material.keyUsages) != len(want) {
		t.Fatalf("material.keyUsages = %v, want %v", material.keyUsages, want)
	}
	for i, name := range want {
		if material.keyUsages[i] != name {
			t.Fatalf("material.keyUsages = %v, want %v", material.keyUsages, want)
		}
	}

	// The import path derives usages from the same PEM, so an imported
	// copy of the issued material persists identical usage fields.
	block, _ := pem.Decode([]byte(material.certPEM))
	if block == nil {
		t.Fatal("issued PEM did not decode")
	}
	parsed, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatalf("parse issued certificate: %v", err)
	}
	kus, ekus := x509UsageFields(parsed)
	if len(kus) != len(want) {
		t.Fatalf("x509UsageFields key usages = %v, want %v", kus, want)
	}
	for i, name := range want {
		if kus[i].Name != name {
			t.Fatalf("x509UsageFields key usages = %v, want %v", kus, want)
		}
	}
	if len(ekus) != 0 {
		t.Fatalf("x509UsageFields extended key usages = %v, want none", ekus)
	}
}

func TestX509UsageFieldsAnyExtendedKeyUsage(t *testing.T) {
	parsed := &x509.Certificate{
		KeyUsage:    x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageAny},
	}
	_, ekus := x509UsageFields(parsed)
	if len(ekus) != 1 || ekus[0].Name != "ANY" || ekus[0].OID != "" {
		t.Fatalf("x509UsageFields extended key usages = %v, want exactly [Name:ANY]", ekus)
	}
	names := extKeyUsageEnumsToStrings(parsed.ExtKeyUsage, nil)
	if len(names) != 1 || names[0] != "ANY" {
		t.Fatalf("extKeyUsageEnumsToStrings = %v, want [ANY]", names)
	}
}
