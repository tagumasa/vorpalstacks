// Package acm provides ACM (AWS Certificate Manager) storage functionality for vorpalstacks.
package acm

import (
	"time"

	"vorpalstacks/internal/utils/aws/types"
)

// Certificate represents an ACM certificate.
type Certificate struct {
	CertificateArn          string              `json:"certificateArn"`
	DomainName              string              `json:"domainName"`
	SubjectAlternativeNames []string            `json:"subjectAlternativeNames,omitempty"`
	DomainValidationOptions []*DomainValidation `json:"domainValidationOptions,omitempty"`
	Serial                  string              `json:"serial"`
	Subject                 string              `json:"subject"`
	Issuer                  string              `json:"issuer"`
	CreatedAt               time.Time           `json:"createdAt"`
	IssuedAt                time.Time           `json:"issuedAt"`
	NotBefore               time.Time           `json:"notBefore"`
	NotAfter                time.Time           `json:"notAfter"`
	KeyAlgorithm            string              `json:"keyAlgorithm"`
	SignatureAlgorithm      string              `json:"signatureAlgorithm"`
	Status                  string              `json:"status"`
	Type                    string              `json:"type"`
	RenewalEligibility      string              `json:"renewalEligibility"`
	RenewalSummary          *RenewalSummary     `json:"renewalSummary,omitempty"`
	Options                 *CertificateOptions `json:"options,omitempty"`
	KeyUsages               []*KeyUsage         `json:"keyUsages,omitempty"`
	ExtendedKeyUsages       []*ExtendedKeyUsage `json:"extendedKeyUsages,omitempty"`
	InUseBy                 []string            `json:"inUseBy,omitempty"`
	ImportedAt              time.Time           `json:"importedAt,omitempty"`
	RevokedAt               time.Time           `json:"revokedAt,omitempty"`
	RevocationReason        string              `json:"revocationReason,omitempty"`
	FailureReason           string              `json:"failureReason,omitempty"`
	CertificateAuthorityArn string              `json:"certificateAuthorityArn,omitempty"`
	Certificate             string              `json:"certificate,omitempty"`
	CertificateChain        string              `json:"certificateChain,omitempty"`
	PrivateKey              string              `json:"privateKey,omitempty"`
	Tags                    []types.Tag         `json:"tags,omitempty"`
	AccountID               string              `json:"accountId"`
	Region                  string              `json:"region"`
}

// DomainValidation represents domain validation details for a certificate.
type DomainValidation struct {
	DomainName       string          `json:"domainName"`
	ValidationEmails []string        `json:"validationEmails,omitempty"`
	ValidationDomain string          `json:"validationDomain"`
	ValidationMethod string          `json:"validationMethod"`
	ValidationStatus string          `json:"validationStatus"`
	ResourceRecord   *ResourceRecord `json:"resourceRecord,omitempty"`
}

// ResourceRecord represents a DNS resource record for domain validation.
type ResourceRecord struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

// CertificateOptions represents options for an ACM certificate.
type CertificateOptions struct {
	CertificateTransparencyLoggingPreference string `json:"certificateTransparencyLoggingPreference"`
	Export                                   string `json:"export"`
}

// RenewalSummary represents the renewal summary for a certificate.
type RenewalSummary struct {
	RenewalStatus           string              `json:"renewalStatus"`
	RenewalStatusReason     string              `json:"renewalStatusReason"`
	DomainValidationOptions []*DomainValidation `json:"domainValidationOptions,omitempty"`
	UpdatedAt               time.Time           `json:"updatedAt"`
}

// KeyUsage represents the key usage extension for a certificate.
type KeyUsage struct {
	Name string `json:"name"`
}

// ExtendedKeyUsage represents an extended key usage extension for a certificate.
type ExtendedKeyUsage struct {
	Name string `json:"name"`
	OID  string `json:"oid"`
}

// CertificateSummary represents a summary of an ACM certificate.
type CertificateSummary struct {
	CertificateArn                       string   `json:"certificateArn"`
	DomainName                           string   `json:"domainName"`
	SubjectAlternativeNameSummaries      []string `json:"subjectAlternativeNameSummaries,omitempty"`
	HasAdditionalSubjectAlternativeNames bool     `json:"hasAdditionalSubjectAlternativeNames"`
	Status                               string   `json:"status"`
	Type                                 string   `json:"type"`
	RenewalEligibility                   string   `json:"renewalEligibility"`
	KeyAlgorithm                         string   `json:"keyAlgorithm"`
	KeyUsages                            []string `json:"keyUsages,omitempty"`
	ExtendedKeyUsages                    []string `json:"extendedKeyUsages,omitempty"`
	InUse                                bool     `json:"inUse"`
	NotBefore                            float64  `json:"notBefore"`
	NotAfter                             float64  `json:"notAfter"`
	CreatedAt                            float64  `json:"createdAt"`
	IssuedAt                             float64  `json:"issuedAt"`
	ImportedAt                           float64  `json:"importedAt"`
	Exported                             bool     `json:"exported"`
	ExportOption                         string   `json:"exportOption"`
}

// CertificateListResult represents the result of listing ACM certificates.
type CertificateListResult struct {
	Certificates []*CertificateSummary
	IsTruncated  bool
	NextToken    string
}

// AccountConfiguration represents the ACM account-level configuration.
type AccountConfiguration struct {
	ExpiryEvents             ExpiryEventsConfiguration `json:"ExpiryEvents"`
	CertVisibilityPreference string                    `json:"CertVisibilityPreference,omitempty"`
}

// ExpiryEventsConfiguration represents the configuration for expiration events.
type ExpiryEventsConfiguration struct {
	DaysBeforeExpiry int `json:"DaysBeforeExpiry,omitempty"`
}
