package sts

import (
	"encoding/base64"
	"encoding/xml"
)

// samlAssertion is the minimal SAML 2.0 structure required to extract the
// fields that AssumeRoleWithSAML surfaces back to the caller. We only care
// about the assertion node and a handful of sub-elements; anything else is
// dropped by the XML decoder. Real AWS SAML responses use the namespace
// urn:oasis:names:tc:SAML:2.0:assertion; the encoding/xml package treats
// namespaces transparently when no namespace prefix is required.
type samlAssertion struct {
	XMLName   xml.Name    `xml:"Assertion"`
	Issuer    string      `xml:"Issuer"`
	Subject   samlSubject `xml:"Subject"`
	Audience  string      `xml:"Conditions>AudienceRestriction>Audience"`
	Assertion string      `xml:"-"`
}

// samlSubject captures the NameID used as the SAML subject identifier. In
// real AWS responses this would be the federated user's persistent or
// transient identifier.
type samlSubject struct {
	NameID string `xml:"NameID"`
}

// decodeSAMLFields decodes the base64-encoded SAML assertion and extracts
// the Issuer, Subject (NameID) and Audience values. Returns empty strings
// when the assertion cannot be parsed — VorpalStacks does not validate
// SAML signatures in TEST_MODE and the SDK tests use a dummy base64 token
// ("VGhpcyBpcyBhIGR1bW15IFNBTUwgYXNzZXJ0aW9u") that decodes to plain
// text rather than XML, so the fallback to legacy defaults preserves
// existing test behaviour.
func decodeSAMLFields(assertion string) (issuer, subject, audience string) {
	raw, err := base64.StdEncoding.DecodeString(assertion)
	if err != nil {
		raw, err = base64.URLEncoding.DecodeString(assertion)
		if err != nil {
			return "", "", ""
		}
	}
	var parsed samlAssertion
	if err := xml.Unmarshal(raw, &parsed); err != nil {
		return "", "", ""
	}
	return parsed.Issuer, parsed.Subject.NameID, parsed.Audience
}
