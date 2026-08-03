package sts

import (
	"encoding/base64"
	"encoding/xml"
	"time"
)

// samlAssertion is the minimal SAML 2.0 structure required to extract the
// fields that AssumeRoleWithSAML surfaces back to the caller. We only care
// about the assertion node and a handful of sub-elements; anything else is
// dropped by the XML decoder. Real AWS SAML responses use the namespace
// urn:oasis:names:tc:SAML:2.0:assertion; the encoding/xml package treats
// namespaces transparently when no namespace prefix is required.
type samlAssertion struct {
	XMLName    xml.Name       `xml:"Assertion"`
	Issuer     string         `xml:"Issuer"`
	Subject    samlSubject    `xml:"Subject"`
	Conditions samlConditions `xml:"Conditions"`
	Assertion  string         `xml:"-"`
}

// samlConditions captures the NotOnOrAfter attribute and the audience
// restriction from the Conditions element.
type samlConditions struct {
	NotOnOrAfter string `xml:"NotOnOrAfter,attr"`
	Audience     string `xml:"AudienceRestriction>Audience"`
}

// samlSubject captures the NameID used as the SAML subject identifier. In
// real AWS responses this would be the federated user's persistent or
// transient identifier.
type samlSubject struct {
	NameID string `xml:"NameID"`
}

// decodeSAMLFields decodes the base64-encoded SAML assertion and extracts
// the Issuer, Subject (NameID), Audience and Conditions/NotOnOrAfter
// values. Returns empty strings when the assertion cannot be parsed —
// VorpalStacks does not validate SAML signatures in TEST_MODE and the SDK
// tests use a dummy base64 token that decodes to plain text rather than
// XML, so the fallback to legacy defaults preserves existing test
// behaviour.
func decodeSAMLFields(assertion string) (issuer, subject, audience, notOnOrAfter string) {
	raw, err := base64.StdEncoding.DecodeString(assertion)
	if err != nil {
		raw, err = base64.URLEncoding.DecodeString(assertion)
		if err != nil {
			return "", "", "", ""
		}
	}
	var parsed samlAssertion
	if err := xml.Unmarshal(raw, &parsed); err != nil {
		return "", "", "", ""
	}
	return parsed.Issuer, parsed.Subject.NameID, parsed.Conditions.Audience, parsed.Conditions.NotOnOrAfter
}

// isSAMLAssertionExpired checks whether a base64-encoded SAML assertion
// has passed its Conditions/NotOnOrAfter boundary. Returns false when the
// assertion cannot be parsed as XML — this preserves TEST_MODE
// compatibility for dummy tokens. A parseable SAML assertion without a
// NotOnOrAfter attribute is treated as malformed and fail-closed (true),
// because SAML assertions must declare conditions.
//
// This function is only reached in TEST_MODE; in production,
// AssumeRoleWithSAML returns ErrIDPCommunicationError before reaching
// this check.
func isSAMLAssertionExpired(assertion string) bool {
	issuer, subject, audience, notOnOrAfter := decodeSAMLFields(assertion)
	if notOnOrAfter == "" {
		// If the assertion was parseable XML (any of issuer/subject/
		// audience was extracted), the missing NotOnOrAfter is a
		// malformed assertion — fail-closed. If nothing was extracted
		// at all, it is a dummy token — return false to preserve
		// TEST_MODE compatibility.
		if issuer == "" && subject == "" && audience == "" {
			return false
		}
		return true
	}
	t, err := time.Parse(time.RFC3339, notOnOrAfter)
	if err != nil {
		return false
	}
	// SAML "NotOnOrAfter" semantics: the assertion is invalid at the exact
	// boundary time, so we use !Before (equivalent to >=) rather than After.
	return !time.Now().UTC().Before(t)
}
